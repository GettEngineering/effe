package types

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

// Based component type declared with effe.Step.
// It must be a function with a specific format.
type SimpleComponent struct {
	Input            *ast.FieldList
	Output           *ast.FieldList
	FuncName         *ast.Ident
	OriginalFuncName *ast.Ident
	Deps             *ast.FieldList
}

func (s SimpleComponent) Name() *ast.Ident {
	return s.FuncName
}

// CaseComponent is a result of parsing an expression with type effe.Case
type CaseComponent struct {
	Children []Component
	Tag      ast.Expr
}

func (c CaseComponent) Name() *ast.Ident {
	caseVal := types.ExprString(c.Tag)
	caseVal = strings.Replace(caseVal, "'", "\\'", -1)
	caseVal = strings.Replace(caseVal, "\"", "", -1)
	return ast.NewIdent(caseVal)
}

// WrapComponent is a result of parsing an expression with type effe.Wrap
type WrapComponent struct {
	Before   *SimpleComponent
	Success  *SimpleComponent
	Failure  *SimpleComponent
	Children []Component
}

func (w WrapComponent) Name() *ast.Ident {
	nameParts := []string{"wrap"}
	if w.Before != nil && w.Success != nil {
		nameParts = append(nameParts, "[")
	}
	if w.Before != nil {
		nameParts = append(nameParts, w.Before.Name().Name)
	}
	if w.Success != nil {
		nameParts = append(nameParts, w.Success.Name().Name)
	}
	if w.Before != nil && w.Success != nil {
		nameParts = append(nameParts, "]")
	}
	return ast.NewIdent(strings.Join(nameParts, " "))
}

// DecisionComponent is a result of parsing an expression with type effe.Wrap
type DecisionComponent struct {
	Cases   []*CaseComponent
	Tag     ast.Expr
	TagName *ast.Ident
	TagType ast.Expr
	Failure *SimpleComponent
}

func (d DecisionComponent) Name() *ast.Ident {
	return ast.NewIdent(fmt.Sprintf("decision %s", d.TagName.Name))
}

// Every component is used by the business process must implement this interface.
type Component interface {
	Name() *ast.Ident
}

// GenerateResult stores the result for a package from a call to Generate.
type GenerateResult struct {
	// PkgPath is the package's PkgPath.
	PkgPath string
	// OutputPath is the path where the generated output should be written.
	// May be empty if there were errors.
	OutputPath string
	// Errs is a slice of errors identified during generation.
	Errs []error
}

type LoadError struct {
	Err error
	Pos token.Pos
}

func (l *LoadError) Error() string {
	return l.Err.Error()
}
