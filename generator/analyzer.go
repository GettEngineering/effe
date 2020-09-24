package generator

import (
	"go/ast"

	"github.com/pkg/errors"
)

type analyzer struct {
	flowDecls []flowDecl
	visited   map[string]struct{}
	inProcess map[string]struct{}
}

func newAnayzer(flowDecls []flowDecl) *analyzer {
	return &analyzer{
		flowDecls: flowDecls,
		visited:   make(map[string]struct{}),
		inProcess: make(map[string]struct{}),
	}
}

func (a *analyzer) sortFlowDeclsByDependecies() ([]flowDecl, []error) {
	errs := []error{}
	sequence := make([]string, len(a.flowDecls))
	for _, flowDecl := range a.flowDecls {
		a.inProcess[flowDecl.FlowName()] = struct{}{}
		callPath, err := a.findFlowCallPath(flowDecl.buildFlowFuncCall)
		if err != nil {
			errs = append(errs, err)
		}
		callPath = append(callPath, flowDecl.FlowName())
		delete(a.inProcess, flowDecl.FlowName())
		for _, call := range callPath {
			sequence = addUniqueString(sequence, call)
		}
	}
	if len(errs) > 0 {
		return []flowDecl{}, errs
	}

	sortedFlowDecls := []flowDecl{}
	for _, p := range sequence {
		for _, flowDecl := range a.flowDecls {
			if p == flowDecl.FlowName() {
				sortedFlowDecls = append(sortedFlowDecls, flowDecl)
				break
			}
		}
	}
	return sortedFlowDecls, nil
}

func (a *analyzer) findFlowCallPath(funcCall *ast.CallExpr) ([]string, error) {
	callPath := []string{}
	for _, arg := range funcCall.Args {
		switch arg := arg.(type) {
		case *ast.CallExpr:
			childDependecies, err := a.findFlowCallPath(arg)
			if err != nil {
				return []string{}, err
			}
			callPath = append(childDependecies, callPath...)
		case *ast.Ident:
			flowDeclIndex := -1
			for i, flowDecl := range a.flowDecls {
				if flowDecl.FlowName() == arg.Name {
					flowDeclIndex = i
					break
				}
			}
			if flowDeclIndex == -1 {
				continue
			}
			flowName := a.flowDecls[flowDeclIndex].FlowName()
			callPath = append(callPath, flowName)
			if _, ok := a.inProcess[flowName]; ok {
				return []string{}, errors.Errorf("circular dependency found for %s", flowName)
			}

			a.inProcess[flowName] = struct{}{}
			defer delete(a.inProcess, flowName)
			dependeciesFordependecies, err := a.findFlowCallPath(a.flowDecls[flowDeclIndex].buildFlowFuncCall)
			if err != nil {
				return []string{}, err
			}
			a.visited[flowName] = struct{}{}

			callPath = append(dependeciesFordependecies, callPath...)
		}
	}
	return callPath, nil
}
