package plugins

import (
	"fmt"
	"go/ast"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/plugin"
)

type history struct {
	plugin.Plugin
}

type historyCall struct {
	name string
	vars []interface{}
}

// This type is used by history plugin. Default implementation can be initialized with NewHistoryImpl
type History interface {
	Write(name string, vars ...interface{})
}

// Default implementation
func NewHistoryImpl() History {
	return &historyImpl{}
}

type historyImpl struct {
	data []historyCall
}

func (h *historyImpl) Write(name string, vars ...interface{}) {
	h.data = append(h.data, historyCall{
		name: name,
		vars: vars,
	})
}

func (h history) buildCall(prefix, name string, hVariable *ast.Field, vars []*ast.Field) *ast.CallExpr {
	args := []ast.Expr{
		&ast.BasicLit{Value: fmt.Sprintf("\"%s call step %s\"", prefix, name)},
	}
	for _, field := range vars {
		firstNmae := field.Names[0]
		if firstNmae != nil && fields.GetTypeStrName(field.Type) != "plugins.History" {
			args = append(args, &ast.BasicLit{Value: fmt.Sprintf("\"%s\"", firstNmae)})
			args = append(args, &ast.BasicLit{Value: firstNmae.Name})
		}
	}

	callExpr := &ast.CallExpr{
		Args: args,
		Fun: &ast.SelectorExpr{
			X:   hVariable.Names[0],
			Sel: ast.NewIdent("Write"),
		},
	}

	return callExpr
}

func (h history) buildHistoryType() ast.Expr {
	return &ast.StarExpr{
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("plugins"),
			Sel: ast.NewIdent("History"),
		},
	}
}

func (h history) Before(ctx plugin.Context, name string, fields []*ast.Field) []ast.Stmt {
	hVariable := ctx.FindInputByType(h.buildHistoryType())
	if hVariable == nil {
		hVariable = ctx.AddInput(h.buildHistoryType())
	}

	return []ast.Stmt{
		&ast.ExprStmt{X: h.buildCall("before", name, hVariable, fields)},
	}
}

func (h history) Change(ctx plugin.Context, c plugin.ComponentStmt) bool {
	return false
}

func (h history) Success(ctx plugin.Context, name string, fields []*ast.Field) []ast.Stmt {
	hVariable := ctx.FindInputByType(h.buildHistoryType())
	if hVariable == nil {
		hVariable = ctx.AddInput(h.buildHistoryType())
	}
	return []ast.Stmt{
		&ast.ExprStmt{X: h.buildCall("successful", name, hVariable, fields)},
	}
}

func (h history) Imports() []string {
	return []string{
		"github.com/GettEngineering/effe/plugins",
	}
}

// Initializes history plugin.
//
// Example of generated code:
//
//      func BuildComponent2(service BuildComponent2Service) BuildComponent2Func {
//          return func(historyPtrVal *plugins.History) error {
//              historyPtrVal.Write("before call step Step1")
//              err := service.Step1()
//              if err != nil {
//                  return err
//              }
//              historyPtrVal.Write("successful call step Step1", "err", err)
//              historyPtrVal.Write("before call step Step2")
//              err = service.Step2()
//              if err != nil {
//                  return err
//              }
//              historyPtrVal.Write("successful call step Step2", "err", err)
//          }
//      }
func NewHistory() plugin.Plugin {
	return &history{}
}
