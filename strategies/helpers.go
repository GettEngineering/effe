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

func BuildReturnStmt(output *ast.FieldList, vars map[string]*ast.Ident, typesInfo *types.Info) *ast.ReturnStmt {
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
		returnStmt.Results = append(returnStmt.Results, buildNilVarByType(output.Type, typesInfo))
	}
	return returnStmt
}

func buildNilVarByType(t ast.Expr, typesInfo *types.Info) ast.Expr {
	if typesInfo.TypeOf(t) == nil {
		switch exprType := t.(type) {
		case *ast.StarExpr:
			t = exprType.X
		case *ast.SelectorExpr:
			t = exprType.X
		}
		if typesInfo.TypeOf(t) == nil {
			return &ast.BasicLit{Value: "nil"}
		}
	}

	originalType := typesInfo.TypeOf(t).Underlying()
	switch castedType := originalType.(type) {
	case *types.Interface, *types.Chan, *types.Pointer:
		return &ast.BasicLit{Value: "nil"}
	case *types.Struct, *types.Array, *types.Slice:
		return &ast.CompositeLit{
			Type: t,
		}
	case *types.Basic:
		switch castedType.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return &ast.BasicLit{Value: "0"}
		case types.String:
			return &ast.BasicLit{Value: "\"\""}
		case types.Float32, types.Float64:
			return &ast.BasicLit{Value: "0.0"}
		case types.Bool:
			return &ast.BasicLit{Value: "false"}
		}
	}

	return &ast.BasicLit{Value: "nil"}
}

func BuildFailureReturnStmt(output *ast.FieldList, vars map[string]*ast.Ident, defaultErrMsg string, typesInfo *types.Info) (*ast.ReturnStmt, bool) {
	fmtUsed := false
	returnStmt := BuildReturnStmt(output, vars, typesInfo)
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
		var lit ast.Expr

		v, ok := vars["error"]
		if ok {
			lit = v
		} else {
			if !fmtUsed {
				fmtUsed = true
			}
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
	returnStmt := BuildReturnStmt(ctx.Output, ctx.Vars, f.TypesInfo())
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
