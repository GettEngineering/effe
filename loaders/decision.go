package loaders

import (
	"go/ast"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/types"

	"github.com/pkg/errors"
)

func genDecisionSwitchStatement(caseStatementExpr ast.Expr) (*ast.Ident, ast.Expr, ast.Expr, error) {
	switch expr := caseStatementExpr.(type) {
	case *ast.SelectorExpr:
		lit, ok := expr.X.(*ast.CompositeLit)
		if !ok {
			return nil, nil, nil, &types.LoadError{
				Pos: caseStatementExpr.Pos(),
				Err: errors.New("first arg must be an initialization of struct and getting an one field"),
			}
		}
		switchTagName := fields.NewIdentWithType(lit.Type)
		switchTag := &ast.SelectorExpr{
			X:   switchTagName,
			Sel: expr.Sel,
		}
		return switchTagName, lit.Type, switchTag, nil
	case *ast.CallExpr:
		funcIdent, ok := expr.Fun.(*ast.Ident)
		if !ok || funcIdent.Name != "new" {
			return nil, nil, nil, &types.LoadError{
				Pos: caseStatementExpr.Pos(),
				Err: errors.New("only function new is supported"),
			}
		}
		switchTagName := fields.NewIdentWithType(expr.Args[0])
		return switchTagName, expr.Args[0], switchTagName, nil
	default:
		return nil, nil, nil, &types.LoadError{
			Err: errors.New("unsupported format"),
			Pos: caseStatementExpr.Pos(),
		}
	}
}

// LoadDecisionComponent converts an expression declared with effe.Decision to a component with type types.DecisionComponent
func LoadDecisionComponent(effeDecisionFuncCall *ast.CallExpr, f FlowLoader) (types.Component, error) {
	if len(effeDecisionFuncCall.Args) < 2 {
		return nil, &types.LoadError{
			Pos: effeDecisionFuncCall.Pos(),
			Err: errors.New("args length must be more than 1"),
		}
	}
	decisionArgs := effeDecisionFuncCall.Args
	switchTagName, switchType, switchTag, err := genDecisionSwitchStatement(decisionArgs[0])
	if err != nil {
		return nil, err
	}
	decisionArgs = RemoveExprByIndex(decisionArgs, 0)
	decisionComponent := &types.DecisionComponent{
		Tag:     switchTag,
		TagType: switchType,
		TagName: switchTagName,
	}
	failureComponent, failureComponentIndex, err := genComponentFromArgsWithType(decisionArgs, FailureExprType, f)
	if err != nil && err != ErrNoExpr {
		return nil, err
	} else if err == nil {
		simpleFailureComponent, ok := failureComponent.(*types.SimpleComponent)
		if !ok {
			return nil, &types.LoadError{
				Err: errors.Errorf("component %s should be a simple component", failureComponent.Name()),
				Pos: effeDecisionFuncCall.Pos(),
			}
		}
		decisionComponent.Failure = simpleFailureComponent
		decisionArgs = RemoveExprByIndex(decisionArgs, failureComponentIndex)
	}

	for _, arg := range decisionArgs {
		caseCall, ok := arg.(*ast.CallExpr)
		if !ok {
			return nil, &types.LoadError{
				Pos: arg.Pos(),
				Err: errors.New("argument is not a call of function"),
			}
		}
		resComponent, err := f.LoadComponent(caseCall)
		if err != nil {
			return nil, err
		}

		caseComponent, ok := resComponent.(*types.CaseComponent)
		if !ok {
			return nil, &types.LoadError{
				Pos: arg.Pos(),
				Err: errors.New("child component is not an case expression"),
			}
		}
		decisionComponent.Cases = append(decisionComponent.Cases, caseComponent)
	}

	return decisionComponent, nil
}
