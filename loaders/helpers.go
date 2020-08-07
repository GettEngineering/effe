package loaders

import (
	"go/ast"

	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
)

// Loaders for components from github.com/GettEngineering/effe/types package
func Default() map[string]ComponentLoadFunc {
	return map[string]ComponentLoadFunc{
		StepExprType:     LoadSimpleComponent,
		WrapExprType:     LoadWrapComponent,
		DecisionExprType: LoadDecisionComponent,
		FailureExprType:  LoadSimpleComponent,
		BeforeExprType:   LoadSimpleComponent,
		SuccessExprType:  LoadSimpleComponent,
		CaseExprType:     LoadCaseComponent,
	}
}

// LoadComponentsWithTypes parses expressions in first argument and
// searches expressions with types in last argument. After that this method execute a method LoadComponent
// for each component and returns a set of components. Key is a string representation of the type.
func LoadComponentsWithTypes(exprs []ast.Expr, f FlowLoader, typeStrings ...string) (map[string]types.Component, []ast.Expr, error) {
	var (
		err                   error
		serviceComponentCalls map[string]*ast.CallExpr
	)

	serviceComponents := make(map[string]types.Component)

	serviceComponentCalls, exprs, err = ParseAndRemoveServiceComponents(exprs, typeStrings...)
	if err != nil {
		return nil, nil, err
	}

	var c types.Component
	if serviceComponentCalls == nil {
		return serviceComponents, exprs, nil
	}
	for k, v := range serviceComponentCalls {
		if v == nil {
			continue
		}
		c, err = f.LoadComponent(v)
		if err != nil {
			return nil, nil, err
		}
		serviceComponents[k] = c
	}
	return serviceComponents, exprs, nil
}

func ParseAndRemoveServiceComponents(exprs []ast.Expr, types ...string) (map[string]*ast.CallExpr, []ast.Expr, error) {
	exprMap := make(map[string]*ast.CallExpr)
	for _, t := range types {
		call, index, err := FindCallExprWithType(exprs, t)
		if err != nil && err != ErrNoExpr {
			return nil, exprs, err
		} else if err != nil && err == ErrNoExpr {
			continue
		}
		if index == -1 {
			continue
		}

		exprMap[t] = call

		exprs = RemoveExprByIndex(exprs, index)
	}

	return exprMap, exprs, nil
}

func RemoveExprByIndex(exprs []ast.Expr, index int) []ast.Expr {
	copy(exprs[index:], exprs[index+1:])
	exprs[len(exprs)-1] = nil
	return exprs[:len(exprs)-1]
}

func FindCallExprWithType(args []ast.Expr, exprType string) (*ast.CallExpr, int, error) {
	for index, arg := range args {
		arg, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}
		effeStepFuncCall, ok := arg.Fun.(*ast.SelectorExpr)

		if !ok {
			return nil, -1, &types.LoadError{
				Err: errors.New("function is not from external packages"),
				Pos: effeStepFuncCall.Pos(),
			}
		}

		if effeStepFuncCall.Sel.Name == exprType {
			return arg, index, nil
		}
	}
	return nil, -1, ErrNoExpr
}

func genComponentFromArgsWithType(args []ast.Expr, exprType string, f FlowLoader) (types.Component, int, error) {
	callExpr, index, err := FindCallExprWithType(args, exprType)
	var c types.Component
	if err != nil {
		return nil, index, err
	}
	c, err = f.LoadComponent(callExpr)
	return c, index, err
}

func NewComponentsFromArgs(args []ast.Expr, f FlowLoader) ([]types.Component, error) {
	components := []types.Component{}
	for i := 0; i < len(args); i++ {
		arg, ok := args[i].(*ast.CallExpr)
		if !ok {
			return nil, &types.LoadError{
				Err: errors.New("arg with index %d is not a call of function"),
				Pos: args[i].Pos(),
			}
		}

		simpleComponent, err := f.LoadComponent(arg)
		if err != nil {
			return nil, err
		}
		components = append(components, simpleComponent)
	}
	return components, nil
}
