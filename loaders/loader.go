package loaders

import (
	"go/ast"

	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
)

// Loader converts flow expressions to components
type Loader interface {
	// LoadFlow takes array of expressions and set of function declartions and returns an
	// array of components and a failure component.
	LoadFlow([]ast.Expr, map[string]*ast.FuncDecl) ([]types.Component, types.Component, error)

	// Adds new handler by type here.
	Register(apiExtType string, c ComponentLoadFunc) error

	// Returns a declaration of function by name
	GetFuncDecl(name string) (*ast.FuncDecl, bool)
}

type loader struct {
	loaders  map[string]ComponentLoadFunc
	packages []string
	decls    map[string]*ast.FuncDecl
}

// Initializes a new Loader
func NewLoader(opts ...Option) Loader {
	l := &loader{
		loaders:  Default(),
		decls:    make(map[string]*ast.FuncDecl),
		packages: []string{},
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type Option func(l *loader)

// WithLoaders is used for settings new loaders in a Loader.
// Default loaders are used from the method DefaultLoader
func WithLoaders(loaders map[string]ComponentLoadFunc) Option {
	return func(l *loader) {
		l.loaders = loaders
	}
}

// WithPackages is used for checking permitted packages which are used in DSL.
// It is necessary because you can add a custom expression type.
//
// Example
//
//      func BuildMyBusinessFlow(){
//            effe.BuildFLow(
//              effe.Step(step1),
//              custompkg.MyStep(step2),
//            )
//      }
//
//      // You need custompkg to loader packages
//      gen := generator.NewGenerator(
//          generator.WithSetttings(settings),
//          generator.WithLoader(loaders.NewLoader(loaders.WithPackages([]string{"effe", "custompkg"}))),
//          generator.WithDrawer(drawer.NewDrawer()),
//          generator.WithStrategy(
//              strategies.NewChain(strategies.WithServiceObjectName(settings.LocalInterfaceVarname())),
//          ),
//      )
func WithPackages(packages []string) Option {
	return func(l *loader) {
		l.packages = packages
	}
}

// Adds new handler by type here.
func (l *loader) Register(apiExtType string, c ComponentLoadFunc) error {
	_, ok := l.loaders[apiExtType]
	if ok {
		return errors.Errorf("api method %s already registered", apiExtType)
	}

	l.loaders[apiExtType] = c
	return nil
}

// LoadComponent converts expressions to a component with a specifig type
func (l *loader) LoadComponent(call *ast.CallExpr) (types.Component, error) {
	t, err := getComponentType(call, l)
	if err != nil {
		return nil, err
	}

	handler, _ := l.getLoader(t)
	return handler(call, l)
}

func (l loader) getLoader(name string) (ComponentLoadFunc, bool) {
	handler, ok := l.loaders[name]
	return handler, ok
}

// Returns a declaration of function by name
func (l loader) GetFuncDecl(name string) (*ast.FuncDecl, bool) {
	v, ok := l.decls[name]
	return v, ok
}

// LoadFlow takes array of expressions and set of function declartions and returns an
// array of components and a failure component. If component with type failure is not found
// LoadFlow returs nil in the second component.
func (l *loader) LoadFlow(args []ast.Expr, decls map[string]*ast.FuncDecl) ([]types.Component, types.Component, error) {
	l.decls = decls
	failureComponent, failureIndex, err := genComponentFromArgsWithType(args, FailureExprType, l)
	if err != nil && err != ErrNoExpr {
		return nil, nil, err
	} else if err == nil {
		args = RemoveExprByIndex(args, failureIndex)
	}

	components, err := NewComponentsFromArgs(args, l)
	if err != nil {
		return nil, nil, err
	}
	return components, failureComponent, nil
}

// FlowLoader provides functions which use in Loader.
// It's a helper for easy reusing existed code
type FlowLoader interface {
	LoadComponent(call *ast.CallExpr) (types.Component, error)
	GetFuncDecl(string) (*ast.FuncDecl, bool)
}

// Type for loaders uses in Loader.
type ComponentLoadFunc func(*ast.CallExpr, FlowLoader) (types.Component, error)

func getComponentType(stepFunc *ast.CallExpr, l *loader) (string, error) {
	effeStepFuncCall, ok := stepFunc.Fun.(*ast.SelectorExpr)

	if !ok {
		return "", &types.LoadError{
			Err: errors.New("function is not a seleted function"),
			Pos: stepFunc.Pos(),
		}
	}
	packageNameIdent, ok := effeStepFuncCall.X.(*ast.Ident)

	if !ok || !l.isRegisteredDSLPackage(packageNameIdent.Name) {
		return "", errors.Errorf("package %s is not registered", packageNameIdent.Name)
	}

	_, ok = l.getLoader(effeStepFuncCall.Sel.Name)
	if !ok {
		return "", &types.LoadError{
			Err: errors.Errorf("unsupported dsl method %s", effeStepFuncCall.Sel.Name),
			Pos: effeStepFuncCall.Pos(),
		}
	}
	return effeStepFuncCall.Sel.Name, nil
}

func (l loader) isRegisteredDSLPackage(name string) bool {
	for _, pkgName := range l.packages {
		if pkgName == name {
			return true
		}
	}
	return false
}
