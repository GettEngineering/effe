package testcustomization

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/strategies"
	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
)

//go:generate go run ./cmd/updater/main.go -- ./testdata

func POST(args interface{}) interface{} {
	panic("implementation is not generated, run myeffe")
}

type PostRequestComponent struct {
	URI ast.Expr
}

func (p PostRequestComponent) Name() *ast.Ident {
	return ast.NewIdent(fmt.Sprintf("POST request to %s", p.URI))
}

type postRequestComponentCall struct {
	fn     ast.Expr
	input  *ast.FieldList
	output *ast.FieldList
	name   *ast.Ident
}

func (c postRequestComponentCall) Name() *ast.Ident {
	return c.name
}

func (c postRequestComponentCall) Input() *ast.FieldList {
	return c.input
}

func (c postRequestComponentCall) Output() *ast.FieldList {
	return c.output
}

func (c postRequestComponentCall) Fn() ast.Expr {
	return c.fn
}

func LoadPostRequestComponent(effeConditionCall *ast.CallExpr, f loaders.FlowLoader) (types.Component, error) {
	return &PostRequestComponent{
		URI: effeConditionCall.Args[0],
	}, nil
}

func GenPostRequestComponent(f strategies.FlowGen, c types.Component) (strategies.ComponentCall, error) {
	component, ok := c.(*PostRequestComponent)
	if !ok {
		return nil, errors.Errorf("component %s is not a component with type PostRequestComponent", component.Name())
	}

	f.AddImport("gopkg.in/h2non/gentleman.v2")

	cli := ast.NewIdent("cli")
	req := ast.NewIdent("req")
	output := &ast.FieldList{
		List: []*ast.Field{
			{
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("gentleman"),
						Sel: ast.NewIdent("Response"),
					},
				},
			},
			{
				Type: ast.NewIdent("error"),
			},
		},
	}
	fn := &ast.FuncLit{
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: output,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{cli},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("gentleman"),
								Sel: ast.NewIdent("New"),
							},
						},
					},
					Tok: token.DEFINE,
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							Sel: cli,
							X:   ast.NewIdent("URI"),
						},
						Args: []ast.Expr{
							component.URI,
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						req,
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   cli,
								Sel: ast.NewIdent("Request"),
							},
						},
					},
					Tok: token.DEFINE,
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   req,
							Sel: ast.NewIdent("Method"),
						},
						Args: []ast.Expr{
							ast.NewIdent("POST"),
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   cli,
								Sel: ast.NewIdent("Send"),
							},
						},
					},
				},
			},
		},
	}

	return &postRequestComponentCall{
		input:  &ast.FieldList{},
		output: output,
		fn:     fn,
	}, nil
}
