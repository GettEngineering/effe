package strategies

import "go/ast"

type componentStmt struct {
	componentName string
	inputFields   []*ast.Field
	outputFields  []*ast.Field
	stmt          ast.Stmt
	failureBlock  *ast.BlockStmt
}

func (c componentStmt) Stmt() ast.Stmt {
	return c.stmt
}

func (c componentStmt) ErrStmt() *ast.BlockStmt {
	return c.failureBlock
}

func (c componentStmt) InputFields() []*ast.Field {
	return c.inputFields
}

func (c componentStmt) OutputFields() []*ast.Field {
	return c.outputFields
}

func (c componentStmt) Name() string {
	return c.componentName
}

type ComponentStmt interface {
	Stmt() ast.Stmt
	InputFields() []*ast.Field
	OutputFields() []*ast.Field
	ErrStmt() *ast.BlockStmt
	Name() string
}
