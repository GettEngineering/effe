package strategies

import "go/ast"

type ComponentCall interface {
	Fn() ast.Expr
	Input() *ast.FieldList
	Output() *ast.FieldList
	Name() *ast.Ident
}

type componentCall struct {
	fn     ast.Expr
	input  *ast.FieldList
	output *ast.FieldList
	name   *ast.Ident
}

func (c componentCall) Fn() ast.Expr {
	return c.fn
}

func (c componentCall) Input() *ast.FieldList {
	return c.input
}

func (c componentCall) Output() *ast.FieldList {
	return c.output
}

func (c componentCall) Name() *ast.Ident {
	return c.name
}
