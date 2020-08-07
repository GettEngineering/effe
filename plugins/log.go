package plugins

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/GettEngineering/effe/plugin"
)

type logPlugin struct{}

func (l logPlugin) Success(ctx plugin.Context, name string, fields []*ast.Field) []ast.Stmt {
	return []ast.Stmt{}
}

func (l logPlugin) Change(ctx plugin.Context, c plugin.ComponentStmt) bool {
	return false
}

func (l logPlugin) Before(ctx plugin.Context, name string, fields []*ast.Field) []ast.Stmt {
	firstArg := fmt.Sprintf("call step %s", name)

	vars := []ast.Expr{}
	for _, field := range fields {
		vars = append(vars, field.Names[0])
	}

	for _, field := range fields {
		firstArg += strings.Join([]string{field.Names[0].Name, "%v"}, ": ")
	}

	exprs := []ast.Expr{
		&ast.BasicLit{Value: fmt.Sprintf("\"%s\\n\"", firstArg)},
	}
	exprs = append(exprs, vars...)

	return []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("log"),
					Sel: ast.NewIdent("Printf"),
				},
				Args: exprs,
			},
		},
	}
}

func (l logPlugin) Imports() []string {
	return []string{
		"log",
	}
}

// Initializes LogPlugin
//
// Example:
//
//      func BuildComponent2(service BuildComponent2Service) BuildComponent2Func {
//          return func() error {
//              log.Printf("call step Step1\n")
//                  err := service.Step1()
//                  if err != nil {
//                      return err
//                  }
//                  log.Printf("call step Step2\n")
//                  err = service.Step2()
//                  if err != nil {
//                      return err
//                  }
//          }
//      }
func NewLogPlugin() plugin.Plugin {
	return &logPlugin{}
}
