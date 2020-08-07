package drawer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
)

// Drawer draws diagram fora  business flow
type Drawer interface {
	// Draws a flow. Takes array of components and a failure component and returns dsl in plantuml
	DrawFlow([]types.Component, types.Component) (string, error)

	// DrawComponent generates a statement from a component
	DrawComponent(types.Component) (ComponentStmt, error)

	// DrawBlock generates multi-statement from an array of components
	DrawBlock([]types.Component, types.Component) (ComponentStmt, error)

	// Register is a method for adding a new generator for a custom component type
	Register(string, Generator) error
}

type drawer struct {
	drawers map[string]Generator
}

// Register is a method for adding a new generator for a custom component type
func (d drawer) Register(apiExtType string, c Generator) error {
	_, ok := d.drawers[apiExtType]
	if ok {
		return errors.Errorf("api method %s already registered", apiExtType)
	}

	d.drawers[apiExtType] = c
	return nil
}

// Initializes a new Drawer
func NewDrawer() Drawer {
	return &drawer{
		drawers: Default(),
	}
}

// Represents a statement for adding in a diagram.
type ComponentStmt interface {
	Stmt() string
	ReturnError() bool
}

type componentStmt struct {
	stmt      string
	returnErr bool
}

func (c componentStmt) Stmt() string {
	return c.stmt
}

func (c componentStmt) ReturnError() bool {
	return c.returnErr
}

// Generator converts a component to a statement
type Generator func(Drawer, types.Component) (ComponentStmt, error)

// Default generators.
// It's possible to build custom map of generators or/and reuse existing.
func Default() map[string]Generator {
	return map[string]Generator{
		"DecisionComponent": DrawDecision,
		"SimpleComponent":   DrawSimple,
		"CaseComponent":     DrawCase,
		"WrapComponent":     DrawWrap,
	}
}

// DrawCase converts a component with type types.CaseComponent to a statement
func DrawCase(d Drawer, c types.Component) (ComponentStmt, error) {
	caseComponent, ok := c.(*types.CaseComponent)
	if !ok {
		return nil, errors.Errorf("can't draw component %s", c.Name())
	}
	childStmt, err := d.DrawBlock(caseComponent.Children, nil)
	if err != nil {
		return nil, err
	}

	return &componentStmt{
		returnErr: childStmt.ReturnError(),
		stmt:      childStmt.Stmt(),
	}, nil
}

func buildStmts(stmts []string) string {
	return strings.Join(stmts, "\n")
}

// DrawDecision converts a component with type types.DecisionComponent to a statement
func DrawDecision(d Drawer, c types.Component) (ComponentStmt, error) {
	dComponent, ok := c.(*types.DecisionComponent)
	if !ok {
		return nil, errors.Errorf("can't draw component %s", c.Name())
	}
	dStmt := &componentStmt{}
	var failureStmt ComponentStmt
	if dComponent.Failure != nil {
		var err error
		failureStmt, err = DrawSimple(d, dComponent.Failure)
		if err != nil {
			return nil, errors.Wrapf(err, "can't render failure component %s", dComponent.Failure.Name())
		}
	}
	stmtBlock := make([]string, 0)

	tagType := fields.GetTypeStrName(dComponent.TagType)

	for index, dCase := range dComponent.Cases {
		caseStmtType := "if"
		if index > 0 {
			caseStmtType = "elseif"
		}
		stmtBlock = append(stmtBlock, fmt.Sprintf("%s (%s equals %s) then (yes)", caseStmtType, tagType, fields.GetTypeStrName(dCase.Tag)))
		caseStmt, err := DrawCase(d, dCase)
		if err != nil {
			return nil, err
		}
		stmtBlock = append(stmtBlock, caseStmt.Stmt())
		if caseStmt.ReturnError() && !dStmt.returnErr {
			dStmt.returnErr = true
		}
	}

	if dStmt.returnErr {
		stmtBlock = append(stmtBlock, drawFailureStmts(fmt.Sprintf("failure decision %s", tagType), failureStmt)...)
	}
	stmtBlock = append(stmtBlock, "endif")

	dStmt.stmt = buildStmts(stmtBlock)
	return dStmt, nil
}

// DrawSimple converts a component with type types.SimpleComponent to a statement
func DrawSimple(d Drawer, c types.Component) (ComponentStmt, error) {
	sComponent, ok := c.(*types.SimpleComponent)
	if !ok {
		return nil, errors.Errorf("can't draw component %s", c.Name())
	}

	stmt := &componentStmt{
		stmt: fmt.Sprintf(":%s;", sComponent.Name()),
	}

	if sComponent.Output != nil {
		for _, output := range sComponent.Output.List {
			if fields.GetTypeStrName(output.Type) == "error" {
				stmt.returnErr = true
				break
			}
		}
	}
	return stmt, nil
}

// DrawWrap converts component with type types.DecisionComponent to a statement
func DrawWrap(d Drawer, c types.Component) (ComponentStmt, error) {
	wComponent, ok := c.(*types.WrapComponent)
	if !ok {
		return nil, errors.Errorf("can't draw component %s", c.Name())
	}
	components := make([]types.Component, 0)

	if wComponent.Before != nil {
		components = append(components, wComponent.Before)
	}

	components = append(components, wComponent.Children...)

	if wComponent.Success != nil {
		components = append(components, wComponent.Success)
	}

	return d.DrawBlock(components, wComponent.Failure)
}

// DrawComponent converts component with dynamic type to a statement.
// DrawComponent uses a generator for a component by type. If the generator is not found
// then the function returns an error.
func (d *drawer) DrawComponent(component types.Component) (ComponentStmt, error) {
	cType := reflect.TypeOf(component).String()
	dotIndex := strings.Index(cType, ".")

	if dotIndex != -1 {
		cType = string([]byte(cType)[dotIndex+1:])
	}

	handler, ok := d.drawers[cType]

	if !ok {
		return nil, errors.Errorf("unsupported type %s", cType)
	}

	return handler(d, component)
}

func drawFailureStmts(msg string, stmt ComponentStmt) []string {
	msg = fmt.Sprintf("if ( %s ) then (yes)", msg)
	if stmt != nil {
		return []string{
			msg,
			stmt.Stmt(),
			"stop",
			"endif",
		}
	}

	return []string{
		msg,
		"stop",
		"endif",
	}
}

// Helper for building multi-statement from an array of components with a specific error handler.
func (d drawer) DrawBlock(components []types.Component, failure types.Component) (ComponentStmt, error) {
	block := &componentStmt{}
	stmts := make([]string, 0)
	var failureStmt ComponentStmt
	if failure != nil {
		var err error
		failureStmt, err = d.DrawComponent(failure)
		if err != nil {
			return nil, err
		}
	}
	for _, component := range components {
		cStmt, err := d.DrawComponent(component)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, cStmt.Stmt())
		if cStmt.ReturnError() {
			if !block.returnErr {
				block.returnErr = true
			}
			stmts = append(stmts, drawFailureStmts("failure", failureStmt)...)
		}
	}
	block.stmt = buildStmts(stmts)
	return block, nil
}

// Generates string in plantuml DSL by array of components and an error handler
func (d *drawer) DrawFlow(components []types.Component, failure types.Component) (string, error) {
	stmts := []string{
		"start",
	}
	fStmts, err := d.DrawBlock(components, failure)
	if err != nil {
		return "", err
	}

	stmts = append(stmts, fStmts.Stmt())
	stmts = append(stmts, "end")
	return buildStmts(stmts), nil
}
