package strategies

import (
	"fmt"
	"go/ast"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/types"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

func GenWrapComponentCall(f FlowGen, wComponent types.Component) (ComponentCall, error) {
	component, ok := wComponent.(*types.WrapComponent)
	if !ok {
		return nil, errors.Errorf("component %s is not a component with type WrapComponent", wComponent.Name())
	}
	calls := make([]ComponentCall, 0)

	if component.Before != nil {
		beforeCall, err := GenSimpleComponentCall(f, component.Before)
		if err != nil {
			return nil, err
		}
		calls = append(calls, beforeCall)
	}

	cCalls, err := GenComponentCalls(f, component.Children...)
	if err != nil {
		return nil, err
	}

	calls = append(calls, cCalls...)

	var successCall ComponentCall
	if component.Success != nil {
		successCall, err = GenSimpleComponentCall(f, component.Success)
		if err != nil {
			return nil, err
		}
		calls = append(calls, successCall)
	}

	var failureCall ComponentCall
	if component.Failure != nil {
		failureCall, err = GenSimpleComponentCall(f, component.Failure)
		if err != nil {
			return nil, err
		}
	}

	return BuildMultiComponentCall(f, calls, failureCall), nil
}

func GenSimpleComponentCall(f FlowGen, sComponent types.Component) (ComponentCall, error) {
	component, ok := sComponent.(*types.SimpleComponent)
	if !ok {
		return nil, errors.Errorf("component %s is not a component with type SimpleComponent", sComponent.Name())
	}
	return &componentCall{
		fn: &ast.SelectorExpr{
			X:   ast.NewIdent(f.ServiceName()),
			Sel: ast.NewIdent(strcase.ToCamel(component.FuncName.Name)),
		},
		input:  component.Input,
		output: component.Output,
		name:   component.Name(),
	}, nil
}

func GenCaseComponentCall(f FlowGen, cComponent types.Component) (ComponentCall, error) {
	component, ok := cComponent.(*types.CaseComponent)
	if !ok {
		return nil, errors.Errorf("component %s is not a component with type CaseComponent", cComponent.Name())
	}

	childCalls, err := GenComponentCalls(f, component.Children...)
	if err != nil {
		return nil, err
	}
	return BuildMultiComponentCall(f, childCalls, nil), nil
}

func buildSwitchStmt(tag ast.Expr, v *ast.Ident) *ast.SwitchStmt {
	switch typedTag := tag.(type) {
	case *ast.Ident:
		tag = v
	case *ast.SelectorExpr:
		typedTag.X = v
		tag = typedTag
	}

	return &ast.SwitchStmt{
		Tag:  tag,
		Body: &ast.BlockStmt{},
	}
}

func GenDecisionComponentCall(f FlowGen, dComponent types.Component) (ComponentCall, error) {
	component, ok := dComponent.(*types.DecisionComponent)
	if !ok {
		return nil, errors.Errorf("component %s is not a component with type DecisionComponent", dComponent.Name())
	}

	calls := make([]ComponentCall, 0)
	for _, caseComponent := range component.Cases {
		call, err := GenCaseComponentCall(f, caseComponent)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}

	ctx := &BlockContext{
		Input:   &ast.FieldList{},
		Output:  &ast.FieldList{},
		Vars:    make(map[string]*ast.Ident),
		Builder: f.VarBuilder(),
	}

	ctx.AddInput(component.TagType)
	for _, call := range calls {
		ctx.CalculateInput([]ComponentCall{call})
		ctx.CalculateOutput([]ComponentCall{call})
	}

	ctx.Output.List = sortComponentOutput(ctx.Output.List)

	sharedVars := make(map[string]*ast.Ident)
	for k, v := range ctx.Vars {
		sharedVars[k] = v
	}

	tagVar, ok := sharedVars[fields.GetTypeStrName(component.TagType)]
	if !ok {
		return nil, errors.Errorf("can't find a variable for switch in component %s", component.Name())
	}

	switchStmt := buildSwitchStmt(component.Tag, tagVar)

	for index, caseCall := range calls {
		ctx.Vars = make(map[string]*ast.Ident)
		for k, v := range sharedVars {
			ctx.Vars[k] = v
		}

		componentStmt := f.BuildComponentStmt(ctx, caseCall, nil)
		block := f.ApplyPlugins(ctx, componentStmt)

		returnStmt := BuildReturnStmt(ctx.Output, ctx.Vars, f.TypesInfo())
		block.List = append(block.List, returnStmt)

		switchStmt.Body.List = append(switchStmt.Body.List, &ast.CaseClause{
			Body: block.List,
			List: []ast.Expr{component.Cases[index].Tag},
		})
	}
	failMsg := fmt.Sprintf("unsupported logic by %s", fields.GetTypeStrName(component.Tag))
	returnStmt, fmtUsed := BuildFailureReturnStmt(ctx.Output, nil, failMsg, f.TypesInfo())
	if fmtUsed {
		f.AddImport(fmtLibrary)
	}

	switchStmt.Body.List = append(switchStmt.Body.List, &ast.CaseClause{
		Body: []ast.Stmt{
			returnStmt,
		},
	})
	ctx.Input.List = sortComponentInput(ctx.Input.List)

	return componentCall{
		fn: &ast.FuncLit{
			Type: &ast.FuncType{
				Params:  ctx.Input,
				Results: ctx.Output,
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					switchStmt,
				},
			},
		},
		name:   dComponent.Name(),
		input:  ctx.Input,
		output: ctx.Output,
	}, nil
}

func GenComponentCalls(f FlowGen, components ...types.Component) ([]ComponentCall, error) {
	calls := make([]ComponentCall, 0)
	for _, component := range components {
		call, err := f.GenComponentCall(component)
		if err != nil {
			return calls, err
		}
		calls = append(calls, call)
	}
	return calls, nil
}
