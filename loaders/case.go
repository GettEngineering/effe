package loaders

import (
	"go/ast"

	types "github.com/GettEngineering/effe/types"
)

// LoadCaseComponent converts an expression declared with effe.Case to a component with type types.CaseComponent
func LoadCaseComponent(
	caseCallExpr *ast.CallExpr,
	f FlowLoader,
) (types.Component, error) {
	caseComponents, err := NewComponentsFromArgs(caseCallExpr.Args[1:], f)
	if err != nil {
		return nil, err
	}

	return &types.CaseComponent{
		Tag:      caseCallExpr.Args[0],
		Children: caseComponents,
	}, nil
}
