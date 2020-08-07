package loaders

import (
	"go/ast"

	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
)

// LoadWrapComponent converts an expression declared with effe.Wrap to a component with type types.WrapComponent
func LoadWrapComponent(effeWrapFuncCall *ast.CallExpr, f FlowLoader) (types.Component, error) {
	if len(effeWrapFuncCall.Args) < 1 {
		return nil, &types.LoadError{
			Err: errors.New("incorrect Wrap usage, args length should be more than 1"),
			Pos: effeWrapFuncCall.Pos(),
		}
	}

	serviceComponents, args, err := LoadComponentsWithTypes(effeWrapFuncCall.Args, f, SuccessExprType, FailureExprType, BeforeExprType)

	if err != nil {
		return nil, err
	}

	bodyComponents, err := NewComponentsFromArgs(args, f)
	if err != nil {
		return nil, err
	}

	wrap := &types.WrapComponent{
		Children: bodyComponents,
	}

	var (
		ok     bool
		simple *types.SimpleComponent
		c      types.Component
	)

	c, ok = serviceComponents[BeforeExprType]
	if ok {
		simple, ok = c.(*types.SimpleComponent)
		if !ok {
			return nil, &types.LoadError{
				Err: errors.New("before function must be a function with a format as for step"),
				Pos: effeWrapFuncCall.Pos(),
			}
		}
		wrap.Before = simple
	}

	c, ok = serviceComponents[SuccessExprType]
	if ok {
		simple, ok = c.(*types.SimpleComponent)
		if !ok {
			return nil, &types.LoadError{
				Err: errors.New("success function must be a function with a format as for step"),
				Pos: effeWrapFuncCall.Pos(),
			}
		}
		wrap.Success = simple
	}

	c, ok = serviceComponents[FailureExprType]
	if ok {
		simple, ok = c.(*types.SimpleComponent)
		if !ok {
			return nil, &types.LoadError{
				Err: errors.New("failure function must be a function with a format as for step"),
				Pos: effeWrapFuncCall.Pos(),
			}
		}
		wrap.Failure = simple
	}
	return wrap, nil
}
