package loaders

import (
	"go/ast"

	"github.com/GettEngineering/effe/types"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

// LoadSimpleComponent converts an expression declared with effe.Step to a component with type types.SimpleComponent
func LoadSimpleComponent(effeStepFuncCall *ast.CallExpr, f FlowLoader) (types.Component, error) {
	if len(effeStepFuncCall.Args) == 0 {
		return nil, &types.LoadError{
			Err: errors.New("no args to declare a step"),
			Pos: effeStepFuncCall.Pos(),
		}
	}
	stepFuncCallIdent, ok := effeStepFuncCall.Args[0].(*ast.Ident)
	if !ok {
		return nil, &types.LoadError{
			Err: errors.New("arg is not an identifier of function"),
			Pos: effeStepFuncCall.Pos(),
		}
	}
	stepFuncCallDecl, ok := f.GetFuncDecl(stepFuncCallIdent.Name)
	if !ok {
		return nil, &types.LoadError{
			Err: errors.Errorf("can't find a function with name %s", stepFuncCallIdent.Name),
			Pos: stepFuncCallIdent.Pos(),
		}
	}

	var returnStmt *ast.ReturnStmt
	for _, stmt := range stepFuncCallDecl.Body.List {
		var tmpReturnStmt *ast.ReturnStmt
		tmpReturnStmt, ok = stmt.(*ast.ReturnStmt)
		if ok {
			returnStmt = tmpReturnStmt
			break
		}
	}
	if returnStmt == nil {
		return nil, &types.LoadError{
			Err: errors.Errorf("function %s has incorrenct format: function must contain only return value", stepFuncCallIdent.Name),
			Pos: stepFuncCallIdent.Pos(),
		}
	}

	if len(returnStmt.Results) == 0 {
		return nil, &types.LoadError{
			Err: errors.Errorf("function %s has incorrenct format: return value should be a function", stepFuncCallIdent.Name),
			Pos: stepFuncCallIdent.Pos(),
		}
	}

	returnFuncLit, ok := returnStmt.Results[0].(*ast.FuncLit)
	if !ok {
		return nil, &types.LoadError{
			Err: errors.Errorf("function %s has incorrenct format: return value should be a function", stepFuncCallIdent.Name),
			Pos: stepFuncCallIdent.Pos(),
		}
	}

	serviceFuncName := ast.NewIdent(strcase.ToCamel(stepFuncCallIdent.Name))

	return &types.SimpleComponent{
		Deps:             stepFuncCallDecl.Type.Params,
		FuncName:         serviceFuncName,
		OriginalFuncName: stepFuncCallIdent,
		Input:            returnFuncLit.Type.Params,
		Output:           returnFuncLit.Type.Results,
	}, nil
}
