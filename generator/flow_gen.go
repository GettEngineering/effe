package generator

import (
	"go/ast"
	"sort"
	"strconv"
	"strings"

	"github.com/GettEngineering/effe/fields"
	"github.com/GettEngineering/effe/types"
)

type flowGen struct {
	pkgFuncDecls map[string]*ast.FuncDecl
	implFields   map[string]implFieldInfo
}

func (f *flowGen) genImplFields(c types.Component) {
	switch c := c.(type) {
	case *types.SimpleComponent:
		f.genImplField(c)
	case *types.WrapComponent:
		if c.Before != nil {
			f.genImplField(c.Before)
		}
		if c.Success != nil {
			f.genImplField(c.Success)
		}
		if c.Failure != nil {
			f.genImplField(c.Failure)
		}
		for _, child := range c.Children {
			f.genImplFields(child)
		}
	case *types.DecisionComponent:
		if c.Failure != nil {
			f.genImplFields(c.Failure)
		}

		for _, decisionCase := range c.Cases {
			f.genImplFields(decisionCase)
		}
	case *types.CaseComponent:
		for _, child := range c.Children {
			f.genImplFields(child)
		}
	}
}

func (f *flowGen) genImplField(simple *types.SimpleComponent) {
	_, ok := f.implFields[simple.FuncName.Name]
	if ok {
		return
	}
	f.implFields[simple.FuncName.Name] = implFieldInfo{
		input:            simple.Input,
		output:           simple.Output,
		serviceFuncName:  simple.FuncName,
		originalFuncName: simple.OriginalFuncName,
		deps:             simple.Deps,
	}
}

func (f *flowGen) sortedImplFields() []implFieldInfo {
	implFields := []implFieldInfo{}
	for _, field := range f.implFields {
		implFields = append(implFields, field)
	}
	sort.SliceStable(implFields, func(i, j int) bool {
		return implFields[i].serviceFuncName.Name <= implFields[j].serviceFuncName.Name
	})
	return implFields
}

func (f flowGen) getSortedFlowDependecies() []*ast.Field {
	depsSet := map[string]*ast.Field{}
	for _, fieldInfo := range f.implFields {
		for _, dep := range fieldInfo.deps.List {
			strTypeName := fields.GetTypeStrName(dep.Type)
			_, ok := depsSet[strTypeName]
			if ok {
				continue
			}
			depsSet[strTypeName] = &ast.Field{
				Type: dep.Type,
				Names: []*ast.Ident{
					{
						Name: dep.Names[0].Name,
					},
				},
			}
		}
	}

	allDeps := []*ast.Field{}

	for _, dep := range depsSet {
		allDeps = append(allDeps, dep)
	}

	sort.SliceStable(allDeps, func(i, j int) bool {
		return fields.GetTypeStrName(allDeps[i].Type) < fields.GetTypeStrName(allDeps[j].Type)
	})

	for i, dep := range allDeps {
		depCounter := 0
		for j := i + 1; j < len(allDeps); j++ {
			if dep.Names[0].Name != allDeps[j].Names[0].Name {
				break
			}
			depCounter++
			allDeps[j].Names[0].Name = strings.Join([]string{
				allDeps[j].Names[0].Name,
				strconv.Itoa(depCounter),
			}, "")
		}
	}
	return allDeps
}
