package main

import (
	"os"

	"github.com/GettEngineering/effe/generator"
	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/plugins"
	"github.com/GettEngineering/effe/strategies"
	"github.com/GettEngineering/effe/testing"
)

func main() {
	settings := generator.DefaultSettigs()
	strategy := strategies.NewChain(
		strategies.WithServiceObjectName(settings.LocalInterfaceVarname()),
		strategies.Use(plugins.NewLogPlugin()),
	)
	gen := generator.NewGenerator(
		generator.WithSetttings(settings),
		generator.WithLoader(loaders.NewLoader(loaders.WithPackages([]string{"effe"}))),
		generator.WithStrategy(strategy),
	)

	testing.UpdateExpectedResult(os.Args[2], gen, map[string][]byte{}, []string{})
}
