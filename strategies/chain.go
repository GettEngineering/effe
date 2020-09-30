package strategies

import (
	"go/ast"
	goTypes "go/types"

	"github.com/GettEngineering/effe/plugin"
	"github.com/GettEngineering/effe/types"
	"github.com/pkg/errors"
)

type Chain interface {
	BuildFlow([]types.Component, types.Component, *goTypes.Info) (ast.Expr, []string, error)
	Register(string, Generator) error
}

type chain struct {
	plugins           []plugin.Plugin
	serviceObjectName string
	generators        map[string]Generator
}

type Option func(c *chain)

func WithServiceObjectName(newServiceObjectName string) Option {
	return func(c *chain) {
		c.serviceObjectName = newServiceObjectName
	}
}

func Use(p plugin.Plugin) Option {
	return func(c *chain) {
		c.plugins = append(c.plugins, p)
	}
}

type Generator func(FlowGen, types.Component) (ComponentCall, error)

func Default() map[string]Generator {
	return map[string]Generator{
		"DecisionComponent": GenDecisionComponentCall,
		"SimpleComponent":   GenSimpleComponentCall,
		"CaseComponent":     GenCaseComponentCall,
		"WrapComponent":     GenWrapComponentCall,
	}
}

func (c *chain) Register(t string, gen Generator) error {
	_, ok := c.generators[t]
	if ok {
		return errors.Errorf("can't register generator for type %s: already registered", t)
	}

	c.generators[t] = gen
	return nil
}

func (c *chain) BuildFlow(components []types.Component, failure types.Component, typesInfo *goTypes.Info) (ast.Expr, []string, error) {
	f := &flowGen{
		globalVarNamesCounter: make(map[string]int),
		importSet:             make(map[string]struct{}),
		chain:                 c,
		serviceObjectName:     c.serviceObjectName,
		typesInfo:             typesInfo,
	}

	f.plugins = c.plugins
	imports := make([]string, 0)
	calls, err := GenComponentCalls(f, components...)
	if err != nil {
		return nil, imports, err
	}
	var call ComponentCall
	if failure != nil {
		failureCall, err := f.GenComponentCall(failure)
		if err != nil {
			return nil, imports, err
		}
		call = BuildMultiComponentCall(f, calls, failureCall)
	} else {
		call = BuildMultiComponentCall(f, calls, nil)
	}

	for impr := range f.importSet {
		imports = append(imports, impr)
	}

	return call.Fn(), imports, nil
}

func NewChain(opts ...Option) Chain {
	c := &chain{
		plugins:           make([]plugin.Plugin, 0),
		serviceObjectName: "service",
		generators:        Default(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
