package strategies

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/plugin"
	"github.com/GettEngineering/effe/types"
)

const (
	errorExpr  = "error"
	fmtLibrary = "fmt"
)

type FlowGen interface {
	AddImport(string)
	BuildComponentStmt(ctx *BlockContext, cCall, failureCall ComponentCall) ComponentStmt
	ApplyPlugins(ctx *BlockContext, componentStmt ComponentStmt) *ast.BlockStmt
	ServiceName() string
	VarBuilder() VarBuilder
	GenComponentCall(types.Component) (ComponentCall, error)
}

func (f flowGen) ServiceName() string {
	return f.serviceObjectName
}

type flowGen struct {
	globalVarNamesCounter map[string]int
	plugins               []plugin.Plugin
	importSet             map[string]struct{}
	serviceObjectName     string
	chain                 *chain
}

func (f *flowGen) GenComponentCall(component types.Component) (ComponentCall, error) {
	cType := reflect.TypeOf(component).String()
	dotIndex := strings.Index(cType, ".")

	if dotIndex != -1 {
		cType = string([]byte(cType)[dotIndex+1:])
	}

	return f.chain.generators[cType](f, component)
}

func (f *flowGen) AddImport(impr string) {
	f.importSet[impr] = struct{}{}
}

func (f *flowGen) VarBuilder() VarBuilder {
	return f.buildVar
}

func (f *flowGen) buildVar(t ast.Expr) *ast.Ident {
	v := fields.NewIdentWithType(t)
	index := f.incrementGlovalVarNameCounter(v.Name)
	if index > 1 {
		v = ast.NewIdent(v.Name + strconv.Itoa(index))
	}
	return v
}

func (f *flowGen) incrementGlovalVarNameCounter(name string) int {
	f.globalVarNamesCounter[name]++
	return f.globalVarNamesCounter[name]
}

func (f *flowGen) buildFailureBlock(ctx *BlockContext, failureCall ComponentCall, output *ast.FieldList, cName *ast.Ident) *ast.BlockStmt {
	errVars := make([]string, 0)
	block := &ast.BlockStmt{}
	if output == nil {
		return nil
	}
	for _, output := range output.List {
		if fields.GetTypeStrName(output.Type) == errorExpr {
			v, ok := ctx.Vars[errorExpr]
			if ok {
				errVars = append(errVars, v.Name)
			}
		}
	}
	if len(errVars) == 0 {
		return nil
	}
	for _, errVarName := range errVars {
		ifBody := &ast.BlockStmt{}
		if failureCall != nil {
			failureStmt := f.BuildComponentStmt(ctx, failureCall, nil)
			ifBody.List = append(ifBody.List, failureStmt.Stmt())
		}

		returnStmt, usedFmtLibrary := BuildFailureReturnStmt(ctx.Output, ctx.Vars, fmt.Sprintf("failure call %s", cName))

		if usedFmtLibrary {
			f.AddImport(fmtLibrary)
		}

		ifBody.List = append(ifBody.List, returnStmt)
		ifStmt := &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent(errVarName),
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: ifBody,
		}
		block.List = append(block.List, ifStmt)
	}

	return block
}

func (f *flowGen) BuildComponentStmt(ctx *BlockContext, cCall, failureCall ComponentCall) ComponentStmt {
	var name string
	if cCall.Name() != nil {
		name = cCall.Name().Name
	} else {
		name = "call func"
	}
	componentStmt := &componentStmt{
		componentName: name,
	}

	componentStmt.inputFields = ctx.BuildInputVars(cCall.Input())

	call := &ast.CallExpr{
		Fun:  cCall.Fn(),
		Args: getNamesFromFieldList(componentStmt.inputFields),
	}

	if cCall.Output() != nil && len(cCall.Output().List) > 0 {
		outputFields, allVarsExist := ctx.BuildOutputVars(cCall.Output())

		assignStmt := &ast.AssignStmt{
			Rhs: []ast.Expr{
				call,
			},
			Lhs: getNamesFromFieldList(outputFields),
		}

		if allVarsExist {
			assignStmt.Tok = token.ASSIGN
		} else {
			assignStmt.Tok = token.DEFINE
		}

		componentStmt.outputFields = outputFields
		componentStmt.stmt = assignStmt
	} else {
		componentStmt.stmt = &ast.ExprStmt{
			X: call,
		}
	}
	componentStmt.failureBlock = f.buildFailureBlock(ctx, failureCall, cCall.Output(), cCall.Name())

	return componentStmt
}

func (f *flowGen) ApplyPlugins(ctx *BlockContext, componentStmt ComponentStmt) *ast.BlockStmt {
	usedPlugins := make(map[plugin.Plugin]struct{})
	block := &ast.BlockStmt{}

	for _, p := range f.plugins {
		beforeStmts := p.Before(ctx, componentStmt.Name(), componentStmt.InputFields())
		block.List = append(block.List, beforeStmts...)
		pluginUsed := p.Change(ctx, componentStmt)
		if len(beforeStmts) > 0 || pluginUsed {
			usedPlugins[p] = struct{}{}
		}
	}

	block.List = append(block.List, componentStmt.Stmt())
	if componentStmt.ErrStmt() != nil {
		block.List = append(block.List, componentStmt.ErrStmt().List...)
	}
	for _, p := range f.plugins {
		afterStmts := p.Success(ctx, componentStmt.Name(), componentStmt.OutputFields())
		if len(afterStmts) > 0 {
			usedPlugins[p] = struct{}{}
		}
		block.List = append(block.List, afterStmts...)
	}

	for p := range usedPlugins {
		for _, impr := range p.Imports() {
			f.AddImport(impr)
		}
	}

	return block
}
