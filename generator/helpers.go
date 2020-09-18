package generator

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

func getExportedType(pkg *packages.Package, fieldSet *ast.FieldList) map[string]struct{} {
	importSet := make(map[string]struct{})
	for _, input := range fieldSet.List {
		paramObj := qualifiedIdentObject(pkg.TypesInfo, input.Type)
		if paramObj == nil {
			continue
		}
		if paramObj.Exported() {
			_, ok := importSet[paramObj.Pkg().Path()]
			if !ok {
				importSet[paramObj.Pkg().Path()] = struct{}{}
			}
		}
	}
	return importSet
}

func mergeImportSets(dst, src map[string]struct{}) {
	for k, v := range src {
		_, ok := dst[k]
		if !ok {
			dst[k] = v
		}
	}
}

func genInitializeImplementFunc(newImpleFuncName, impleName *ast.Ident, allDeps []*ast.Field, assignExp []ast.Expr) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: newImpleFuncName,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: allDeps,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: impleName,
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							X: &ast.CompositeLit{
								Type: impleName,
								Elts: assignExp,
							},
							Op: token.AND,
						},
					},
				},
			},
		},
	}
}

func newFlowFunc(settingLocalInterfaceVarName string, interfaceName, funcName, typeFuncName *ast.Ident, funcLit *ast.FuncLit) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: funcName,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent(settingLocalInterfaceVarName)},
						Type:  interfaceName,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: typeFuncName,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						funcLit,
					},
				},
			},
		},
	}
}

func qualifiedIdentObject(info *types.Info, expr ast.Expr) types.Object {
	switch expr := expr.(type) {
	case *ast.Ident:
		return info.ObjectOf(expr)
	case *ast.SelectorExpr:
		pkgName, ok := expr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		if _, ok := info.ObjectOf(pkgName).(*types.PkgName); !ok {
			return nil
		}
		return info.ObjectOf(expr.Sel)
	default:
		return nil
	}
}

func isEffeImport(path string) bool {
	// TODO(light): This is depending on details of the current loader.
	const vendorPart = "vendor/"
	if i := strings.LastIndex(path, vendorPart); i != -1 && (i == 0 || path[i-1] == '/') {
		path = path[i+len(vendorPart):]
	}
	return path == "github.com/GettEngineering/effe"
}

func findExprInBody(fn *ast.FuncDecl, pkg *packages.Package, extrType string) *ast.CallExpr {
	for _, stmt := range fn.Body.List {
		stmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := stmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		if qualifiedIdentObject(pkg.TypesInfo, call.Fun) == types.Universe.Lookup("panic") {
			if len(call.Args) != 1 {
				continue
			}
			call, ok = call.Args[0].(*ast.CallExpr)
			if !ok {
				continue
			}
		}
		buildObj := qualifiedIdentObject(pkg.TypesInfo, call.Fun)
		if buildObj == nil || buildObj.Pkg() == nil || !isEffeImport(buildObj.Pkg().Path()) || buildObj.Name() != extrType {
			continue
		}
		return call
	}
	return nil
}

func addUniqueString(a []string, val string) []string {
	for _, v := range a {
		if v == val {
			return a
		}
	}
	a = append(a, val)
	return a
}
