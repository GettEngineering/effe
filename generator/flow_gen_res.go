package generator

import (
	"go/ast"
	"strings"

	"github.com/GettEngineering/effe/fields"
	"github.com/pkg/errors"
)

type implFieldInfo struct {
	serviceFuncName  *ast.Ident
	originalFuncName *ast.Ident
	input            *ast.FieldList
	output           *ast.FieldList
	deps             *ast.FieldList
}

type flowGenRes struct {
	implFuncDecls          []*ast.FuncDecl
	typeSpecs              []*ast.TypeSpec
	flowFuncDecl           *ast.FuncDecl
	depInitializerFuncDecl *ast.FuncDecl
	imports                []string
}

func (g Generator) genFlowFunc(funcName, interfaceName *ast.Ident, flowFunc *ast.FuncLit) (*ast.TypeSpec, *ast.FuncDecl) {
	typeFunc := &ast.TypeSpec{
		Name: ast.NewIdent(funcName.Name + g.settings.FlowFuncPostfix()),
		Type: flowFunc.Type,
	}
	return typeFunc, newFlowFunc(g.settings.LocalInterfaceVarname(), interfaceName, funcName, typeFunc.Name, flowFunc)
}

func (g Generator) genImplField(impleName *ast.Ident, field implFieldInfo, deps []*ast.Field) (*ast.Field, *ast.FuncDecl, *ast.KeyValueExpr) {
	structFieldIdent := ast.NewIdent(field.originalFuncName.Name + g.settings.ImplFieldPostfix())
	structField := &ast.Field{
		Names: []*ast.Ident{structFieldIdent},
		Type: &ast.FuncType{
			Params:  field.input,
			Results: field.output,
		},
	}

	depArgs := []ast.Expr{}
	for _, dep := range field.deps.List {
		flowDep := fields.FindFieldWithType(deps, dep.Type)
		depArgs = append(depArgs, flowDep.Names[0])
	}

	assignExpr := &ast.KeyValueExpr{
		Key: structFieldIdent,
		Value: &ast.CallExpr{
			Fun:  field.originalFuncName,
			Args: depArgs,
		},
	}

	callArgs := []ast.Expr{}
	for _, input := range field.input.List {
		callArgs = append(callArgs, input.Names[0])
	}

	impleNameIdent := ast.NewIdent(strings.ToLower(string([]rune(impleName.Name)[0])))
	implFunc := &ast.FuncDecl{
		Name: field.serviceFuncName,
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{impleNameIdent},
					Type: &ast.StarExpr{
						X: impleName,
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  field.input,
			Results: field.output,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}

	if field.output != nil && len(field.output.List) > 0 {
		implFunc.Body.List = append(implFunc.Body.List, &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   impleNameIdent,
						Sel: structFieldIdent,
					},
					Args: callArgs,
				},
			},
		})
	} else {
		implFunc.Body.List = append(implFunc.Body.List, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   impleNameIdent,
					Sel: structFieldIdent,
				},
				Args: callArgs,
			},
		})
	}

	return structField, implFunc, assignExpr
}

func (g Generator) genImplementation(impleName, newImpleFuncName *ast.Ident, f *flowGen) (*ast.TypeSpec, *ast.FuncDecl, []*ast.FuncDecl) {
	funcDecls := make([]*ast.FuncDecl, 0)
	structType := &ast.StructType{
		Fields: &ast.FieldList{},
	}
	allDeps := f.getSortedFlowDependecies()

	assignExprs := []ast.Expr{}

	for _, field := range f.sortedImplFields() {
		strField, implFunc, assignExp := g.genImplField(impleName, field, allDeps)
		assignExprs = append(assignExprs, assignExp)
		structType.Fields.List = append(structType.Fields.List, strField)
		funcDecls = append(funcDecls, implFunc)
	}
	interfaceServiceFunc := genInitializeImplementFunc(newImpleFuncName, impleName, allDeps, assignExprs)

	typeSpec := &ast.TypeSpec{
		Name: impleName,
		Type: structType,
	}
	return typeSpec, interfaceServiceFunc, funcDecls
}

func genInterface(interfaceName *ast.Ident, f *flowGen) *ast.TypeSpec {
	inter := &ast.InterfaceType{
		Methods: &ast.FieldList{},
	}

	for _, fieldInfo := range f.sortedImplFields() {
		fn := &ast.FuncType{
			Params:  fieldInfo.input,
			Results: fieldInfo.output,
		}

		inter.Methods.List = append(inter.Methods.List, &ast.Field{
			Names: []*ast.Ident{fieldInfo.serviceFuncName},
			Type:  fn,
		})
	}

	return &ast.TypeSpec{
		Name: interfaceName,
		Type: inter,
	}
}

func (g Generator) genFlow(flowFunc *ast.FuncDecl, buildFlowFuncCall *ast.CallExpr, f *flowGen) (*flowGenRes, error) {
	flowComponents, failureComponent, err := g.loader.LoadFlow(buildFlowFuncCall.Args, f.pkgFuncDecls)
	if err != nil {
		return nil, err
	}

	fn, imports, err := g.strategy.BuildFlow(flowComponents, failureComponent)
	if err != nil {
		return nil, err
	}
	resFunc, ok := fn.(*ast.FuncLit)
	if !ok {
		return nil, errors.New("something goes wrong")
	}

	for _, flowComponent := range flowComponents {
		f.genImplFields(flowComponent)
	}

	if failureComponent != nil {
		f.genImplFields(failureComponent)
	}

	interfaceName := ast.NewIdent(flowFunc.Name.Name + g.settings.InterfaceNamePostfix())
	implName := ast.NewIdent(flowFunc.Name.Name + g.settings.ImplPostfix())
	newImplFuncName := ast.NewIdent(g.settings.NewImplFuncPrefix() + implName.Name)

	flowDeclTypeSpec, flowDecl := g.genFlowFunc(flowFunc.Name, interfaceName, resFunc)
	implTypeSpec, implInitializationFunc, implMethods := g.genImplementation(implName, newImplFuncName, f)
	serviceInterfaceSpec := genInterface(interfaceName, f)
	res := &flowGenRes{
		imports:                imports,
		flowFuncDecl:           flowDecl,
		implFuncDecls:          implMethods,
		depInitializerFuncDecl: implInitializationFunc,
	}
	res.typeSpecs = append(res.typeSpecs, serviceInterfaceSpec)
	res.typeSpecs = append(res.typeSpecs, implTypeSpec)
	res.typeSpecs = append(res.typeSpecs, flowDeclTypeSpec)
	return res, nil
}
