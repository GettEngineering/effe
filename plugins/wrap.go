package plugins

import (
	"fmt"
	"go/ast"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/plugin"
)

type wrapError struct {
	plugin.Plugin
}

func (w wrapError) Before(ctx plugin.Context, name string, fields []*ast.Field) []ast.Stmt {
	return []ast.Stmt{}
}

func (w wrapError) Success(ctx plugin.Context, name string, fields []*ast.Field) []ast.Stmt {
	return []ast.Stmt{}
}

func (w wrapError) findAllReturnStmts(block *ast.BlockStmt) []*ast.ReturnStmt {
	var returnStmts []*ast.ReturnStmt

	for _, stmt := range block.List {
		r, ok := stmt.(*ast.ReturnStmt)

		if ok {
			returnStmts = append(returnStmts, r)
			continue
		}

		ifStmt, ok := stmt.(*ast.IfStmt)
		if !ok {
			continue
		}

		for _, stmt := range ifStmt.Body.List {
			r, ok := stmt.(*ast.ReturnStmt)

			if ok {
				returnStmts = append(returnStmts, r)
				continue
			}
		}
	}
	return returnStmts
}

func (w wrapError) Change(ctx plugin.Context, componentStmt plugin.ComponentStmt) bool {
	if componentStmt.ErrStmt() == nil {
		return false
	}

	returnStmts := w.findAllReturnStmts(componentStmt.ErrStmt())

	errIndexes := make(map[int]struct{})

	for index, outputField := range ctx.OutputList() {
		if fields.GetTypeStrName(outputField.Type) == "error" {
			errIndexes[index] = struct{}{}
		}
	}

	pluginUsed := false

	for _, returnStmt := range returnStmts {
		for index := range errIndexes {
			var (
				sel  *ast.Ident
				args []ast.Expr
			)

			msg := &ast.BasicLit{
				Value: fmt.Sprintf("\"failure call %s\"", componentStmt.Name()),
			}

			field := returnStmt.Results[index]

			if fields.GetTypeStrName(field) == "nil" {
				sel = ast.NewIdent("New")
				args = append(args, msg)
			} else {
				sel = ast.NewIdent("Wrap")

				args = append(args, field)
				args = append(args, msg)
			}

			if !pluginUsed {
				pluginUsed = true
			}

			returnStmt.Results[index] = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("errors"),
					Sel: sel,
				},
				Args: args,
			}
		}
	}

	return pluginUsed
}

func (w wrapError) Imports() []string {
	return []string{
		"github.com/pkg/errors",
	}
}

// Initializes WrapPlugin
//
// Example:
//
//  func BuildComponent2(service BuildComponent2Service) BuildComponent2Func {
//      return func() error {
//          err := service.Step1()
//          if err != nil {
//              return errors.Wrap(err, "failure call Step1")
//          }
//          err = service.Step2()
//          if err != nil {
//              return errors.Wrap(err, "failure call Step2")
//          }
//      }
//  }
func NewWrapError() plugin.Plugin {
	return &wrapError{}
}
