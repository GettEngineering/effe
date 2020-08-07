package effe_test

import (
	"testing"

	"github.com/GettEngineering/effe/generator"
	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/strategies"
	tEffe "github.com/GettEngineering/effe/testing"
)

func TestEffe(t *testing.T) {
	settings := generator.DefaultSettigs()
	gen := generator.NewGenerator(
		generator.WithSetttings(settings),
		generator.WithLoader(loaders.NewLoader(loaders.WithPackages([]string{"effe"}))),
		generator.WithStrategy(
			strategies.NewChain(strategies.WithServiceObjectName(settings.LocalInterfaceVarname())),
		),
	)
	tEffe.RunTests(t, gen, "testdata", nil, []string{})
}
