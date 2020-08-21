package strategies

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"github.com/GettEngineering/effe/fields"
)

func getNamesFromFieldList(list []*ast.Field) []ast.Expr {
	if list == nil {
		return []ast.Expr{}
	}
	vars := make([]ast.Expr, len(list))
	for index, field := range list {
		vars[index] = field.Names[0]
	}
	return vars
}

func BuildReturnStmt(output *ast.FieldList, vars map[string]*ast.Ident) *ast.ReturnStmt {
	returnStmt := &ast.ReturnStmt{}
	for _, output := range output.List {
		typeStrName := fields.GetTypeStrName(output.Type)
		if typeStrName == "error" {
			returnStmt.Results = append(returnStmt.Results, &ast.BasicLit{Value: "nil"})
			continue
		}
		if vars != nil {
			v, ok := vars[typeStrName]
			if ok {
				returnStmt.Results = append(returnStmt.Results, &ast.Ident{
					Name: v.Name,
					Obj: &ast.Object{
						Type: output.Type,
					},
				})
				continue
			}
		}
		returnStmt.Results = append(returnStmt.Results, buildNilVarByType(output.Type))
	}
	return returnStmt
}

func buildTypeWithVar(t ast.Expr, name string) ast.Expr {
	switch t.(type) {
	case *ast.StructType:
		return &ast.CallExpr{
			Fun: ast.NewIdent("make"),
			Args: []ast.Expr{
				ast.NewIdent(name),
			},
		}
	case *ast.ArrayType:
	case *ast.BasicLit:
	case *ast.InterfaceType, *ast.FuncType:
		return &ast.BasicLit{Value: "nil"}
	}
	return t
}

func buildNilVarByType(t ast.Expr) ast.Expr {
	switch t := t.(type) {
	case *ast.StarExpr:
		return &ast.BasicLit{Value: "nil"}
	case *ast.Ident:
		if t.Obj != nil && t.Obj.Decl != nil {
			typeSpec, ok := t.Obj.Decl.(*ast.TypeSpec)
			if ok {
				return buildTypeWithVar(typeSpec.Type, t.Name)
			}
		}
		return &ast.BasicLit{Value: "nil"}
	case *ast.SelectorExpr:
		return &ast.CompositeLit{
			Type: t,
		}
	}

	return &ast.BasicLit{Value: "nil"}
}

func BuildFailureReturnStmt(output *ast.FieldList, vars map[string]*ast.Ident, defaultErrMsg string) (*ast.ReturnStmt, bool) {
	fmtUsed := false
	returnStmt := BuildReturnStmt(output, vars)
	for index, output := range output.List {
		typeStrName := fields.GetTypeStrName(output.Type)
		if typeStrName != "error" {
			continue
		}

		msg := &ast.BasicLit{
			Value: "\"c\"",
		}
		if defaultErrMsg != "" {
			msg.Value = fmt.Sprintf("\"%s\"", defaultErrMsg)
		}

		if !fmtUsed {
			fmtUsed = true
		}
		var lit ast.Expr

		v, ok := vars["error"]
		if ok {
			lit = v
		} else {
			lit = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(fmtLibrary),
					Sel: ast.NewIdent("Errorf"),
				},
				Args: []ast.Expr{msg},
			}
		}

		returnStmt.Results[index] = lit
	}
	return returnStmt, fmtUsed
}

func sortComponentOutput(result []*ast.Field) []*ast.Field {
	sort.SliceStable(result, func(i, j int) bool {
		if types.ExprString(result[i].Type) == errorExpr && types.ExprString(result[j].Type) != errorExpr {
			return false
		}
		return true
	})
	return result
}

func sortComponentInput(result []*ast.Field) []*ast.Field {
	sort.SliceStable(result, func(i, j int) bool {
		if types.ExprString(result[i].Type) == "context.Context" && types.ExprString(result[j].Type) != "context.Context" {
			return true
		}
		return false
	})
	return result
}

func BuildMultiComponentName(calls []ComponentCall) *ast.Ident {
	l := len(calls)
	if l == 1 {
		return calls[0].Name()
	}

	if calls[0].Name() != nil && calls[l-1].Name() != nil {
		return ast.NewIdent(strings.Join([]string{
			calls[0].Name().Name,
			calls[l-1].Name().Name,
		}, "-"))
	}

	return ast.NewIdent("Foo")
}

func BuildMultiComponentCall(f FlowGen, calls []ComponentCall, failureCall ComponentCall) ComponentCall {
	ctx := &BlockContext{
		Input:   new(ast.FieldList),
		Output:  new(ast.FieldList),
		Vars:    make(map[string]*ast.Ident),
		Builder: f.VarBuilder(),
	}

	ctx.CalculateInput(calls)
	ctx.CalculateOutput(calls)

	ctx.Output.List = sortComponentOutput(ctx.Output.List)
	block := &ast.BlockStmt{}

	for _, call := range calls {
		componentStmt := f.BuildComponentStmt(ctx, call, failureCall)
		componentBlock := f.ApplyPlugins(ctx, componentStmt)
		block.List = append(block.List, componentBlock.List...)
	}

	fnType := &ast.FuncType{
		Params:  ctx.Input,
		Results: ctx.Output,
	}
	ctx.Input.List = sortComponentInput(ctx.Input.List)
	name := BuildMultiComponentName(calls)
	returnStmt := BuildReturnStmt(ctx.Output, ctx.Vars)
	block.List = append(block.List, returnStmt)

	return &componentCall{
		input:  ctx.Input,
		output: ctx.Output,
		fn: &ast.FuncLit{
			Type: fnType,
			Body: block,
		},
		name: name,
	}
}
