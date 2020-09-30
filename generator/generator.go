package generator

import (
	"context"
	"fmt"
	"go/ast"
	goTypes "go/types"

	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

const (
	// The function from Effe, which declares flow
	BuildFLowExprType = "BuildFlow"
)

// Generator loads dsl, generates code or diagrams.
type Generator struct {
	settings Settings
	strategy Strategy
	loader   Loader
	drawer   Drawer
}

// Loader executes parsers for components by type. Generator gets arguments from
// expression BuildFlow and passes to Loader.
type Loader interface {
	// Convert expressions to Components. Returns array of Components and a component with type Failure.
	LoadFlow([]ast.Expr, map[string]*ast.FuncDecl) ([]types.Component, types.Component, error)
}

// Strategy generates the flow function.
type Strategy interface {
	// BuildFlow takes a list of components and a failure component and returns
	// the flow function and an array of imports.sss
	BuildFlow([]types.Component, types.Component, *goTypes.Info) (ast.Expr, []string, error)
}

// Drawe draws graphs for business flows.
type Drawer interface {
	// DrawBuild takes a list of compoments and a failure components, returns
	// string in plantuml dsl
	DrawFlow([]types.Component, types.Component) (string, error)
}

type Option func(g *Generator)

// WithSetttings is used for overriding settings
// Default is settings from a method DefaultSettings()
func WithSetttings(s Settings) Option {
	return func(g *Generator) {
		g.settings = s
	}
}

// WithDrawer is used for overriding a drawer
func WithDrawer(d Drawer) Option {
	return func(g *Generator) {
		g.drawer = d
	}
}

// WithDrawer is used for overriding a loader
func WithLoader(l Loader) Option {
	return func(g *Generator) {
		g.loader = l
	}
}

// WithDrawer is used for overriding a strategy
func WithStrategy(s Strategy) Option {
	return func(g *Generator) {
		g.strategy = s
	}
}

// Initialize a new generator with options
func NewGenerator(opts ...Option) *Generator {
	g := &Generator{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// This method generates diagrams for a package
func (g *Generator) GenerateDiagram(ctx context.Context, wd string, env []string, patterns []string, outputDir string) ([]types.GenerateResult, []error) {
	pkgs, errs := load(ctx, wd, env, patterns)
	if len(errs) > 0 {
		return nil, errs
	}
	generated := make([]types.GenerateResult, 0)
	for _, pkg := range pkgs {
		genRes := types.GenerateResult{
			PkgPath: pkg.PkgPath,
		}
		res, errs := g.generateDiagramForPkg(pkg)
		if errs != nil {
			genRes.Errs = append(genRes.Errs, errs...)
			generated = append(generated, genRes)
			continue
		}

		outputFlows, err := writeDiagrams(pkg, outputDir, res)
		if err != nil {
			genRes.Errs = append(genRes.Errs, errs...)
			generated = append(generated, genRes)
			continue
		}

		for _, output := range outputFlows {
			generated = append(generated, types.GenerateResult{
				PkgPath:    pkg.PkgPath,
				OutputPath: output,
			})
		}
	}
	return generated, nil
}

type drawFlowRes struct {
	name  string
	graph string
}

func (g *Generator) generateDiagramForPkg(pkg *packages.Package) ([]drawFlowRes, []error) {
	pkgFuncDecls, flowDecls := g.loadFuncsAndFlows(pkg)
	analyzer := newAnayzer(flowDecls)
	sortedFlowDecls, errs := analyzer.sortFlowDeclsByDependecies()
	if len(errs) > 0 {
		return nil, errs
	}

	flows := make([]drawFlowRes, 0)

	for _, flowDecl := range sortedFlowDecls {
		flowComponents, failureComponent, err := g.loader.LoadFlow(flowDecl.buildFlowFuncCall.Args, pkgFuncDecls)
		if err != nil {
			errs = append(errs, errors.Errorf("can't load flow: %s, error: %s", flowDecl.flowFunc.Name, err))
			continue
		}
		var flowGraph string
		flowGraph, err = g.drawer.DrawFlow(flowComponents, failureComponent)
		if err != nil {
			errs = append(errs, errors.Errorf("can't load flow: %s, error: %s", flowDecl.flowFunc.Name, err))
			continue
		}
		flows = append(flows, drawFlowRes{
			name:  flowDecl.flowFunc.Name.Name,
			graph: flowGraph,
		})
		pkgFuncDecls[flowDecl.flowFunc.Name.Name] = generateEmptyFlowFuncAsStepDeclaration(flowDecl.flowFunc)
	}

	return flows, errs
}

func generateEmptyFlowFuncAsStepDeclaration(flowFuncDecl *ast.FuncDecl) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: flowFuncDecl.Name,
		Type: &ast.FuncType{
			Params: flowFuncDecl.Type.Params,
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.FuncType{
							Params:  flowFuncDecl.Type.Params,
							Results: flowFuncDecl.Type.Results,
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.FuncLit{
							Type: &ast.FuncType{
								Params:  flowFuncDecl.Type.Params,
								Results: flowFuncDecl.Type.Results,
							},
							Body: &ast.BlockStmt{},
						},
					},
				},
			},
		},
	}
}

// This method generates code for a directory and environments.
func (g *Generator) Generate(ctx context.Context, wd string, env []string, patterns []string) ([]types.GenerateResult, []error) {
	pkgs, errs := load(ctx, wd, env, patterns)
	if len(errs) > 0 {
		return nil, errs
	}
	generated := make([]types.GenerateResult, len(pkgs))
	for i, pkg := range pkgs {
		generated[i].PkgPath = pkg.PkgPath
		p, errs := g.generateForPackage(pkg)
		if errs != nil {
			generated[i].Errs = append(generated[i].Errs, errs...)
			continue
		}

		outputFileName, err := writeGeneratedCode(pkg, p)
		if err != nil {
			generated[i].Errs = append(generated[i].Errs, err)
			continue
		}

		generated[i].OutputPath = outputFileName
	}

	return generated, nil
}

type flowDecl struct {
	flowFunc          *ast.FuncDecl
	buildFlowFuncCall *ast.CallExpr
}

func (f flowDecl) FlowName() string {
	return f.flowFunc.Name.Name
}

type pkgGen struct {
	flowFuncDecls           []*ast.FuncDecl
	depInitializerFuncDecls []*ast.FuncDecl
	implFuncDecls           []*ast.FuncDecl
	typeSpecs               []*ast.TypeSpec
	imports                 []string
}

func (g *Generator) loadFuncsAndFlows(pkg *packages.Package) (map[string]*ast.FuncDecl, []flowDecl) {
	pkgFuncDecls := make(map[string]*ast.FuncDecl)
	flowDecls := []flowDecl{}

	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fn.Body == nil {
				continue
			}
			pkgFuncDecls[fn.Name.String()] = fn
			buildFlowFuncCall := findExprInBody(fn, pkg, BuildFLowExprType)
			if buildFlowFuncCall == nil {
				continue
			}

			flowDecls = append(flowDecls, flowDecl{
				flowFunc:          fn,
				buildFlowFuncCall: buildFlowFuncCall,
			})
		}
	}

	return pkgFuncDecls, flowDecls
}

func (g *Generator) generateForPackage(pkg *packages.Package) (*pkgGen, []error) {
	pkgFuncDecls, flowDecls := g.loadFuncsAndFlows(pkg)
	analyzer := newAnayzer(flowDecls)
	sortedFlowDecls, errs := analyzer.sortFlowDeclsByDependecies()
	if len(errs) > 0 {
		return nil, errs
	}

	importSet := make(map[string]struct{})
	p := &pkgGen{}
	for _, flowDecl := range sortedFlowDecls {
		f := &flowGen{
			pkgFuncDecls: pkgFuncDecls,
			implFields:   make(map[string]implFieldInfo),
		}

		res, err := g.genFlow(flowDecl.flowFunc, flowDecl.buildFlowFuncCall, f, pkg.TypesInfo)
		if err != nil {
			loadErr, ok := err.(*types.LoadError)
			if ok {
				position := pkg.Fset.Position(loadErr.Pos)
				err = fmt.Errorf("%s %s", loadErr.Err, position.String())
			}

			errs = append(errs, err)
			continue
		}

		//Import types, which are used in flow
		for _, fieldInfo := range f.implFields {
			if fieldInfo.input != nil {
				mergeImportSets(importSet, getExportedType(pkg, fieldInfo.input))
			}
			if fieldInfo.output != nil {
				mergeImportSets(importSet, getExportedType(pkg, fieldInfo.output))
			}
		}

		for _, impr := range res.imports {
			importSet[impr] = struct{}{}
		}

		p.depInitializerFuncDecls = append(p.depInitializerFuncDecls, res.depInitializerFuncDecl)
		p.flowFuncDecls = append(p.flowFuncDecls, res.flowFuncDecl)
		p.implFuncDecls = append(p.implFuncDecls, res.implFuncDecls...)
		p.typeSpecs = append(p.typeSpecs, res.typeSpecs...)
		pkgFuncDecls[flowDecl.FlowName()] = res.flowFuncDecl
		continue
	}

	if len(errs) > 0 {
		return nil, errs
	}
	for k := range importSet {
		p.imports = append(p.imports, k)
	}

	return p, nil
}

func load(ctx context.Context, wd string, env []string, patterns []string) ([]*packages.Package, []error) {
	cfg := &packages.Config{
		Context:    ctx,
		Mode:       packages.LoadAllSyntax, //nolint:staticcheck
		Dir:        wd,
		Env:        env,
		BuildFlags: []string{"-tags=effeinject"},
		// TODO(light): Use ParseFile to skip function bodies and comments in indirect packages.
	}
	escaped := make([]string, len(patterns))
	for i := range patterns {
		escaped[i] = "pattern=" + patterns[i]
	}
	pkgs, err := packages.Load(cfg, escaped...)
	if err != nil {
		return nil, []error{err}
	}
	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return pkgs, nil
}
