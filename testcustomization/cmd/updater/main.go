package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"

	"github.com/GettEngineering/effe/generator"
	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/strategies"
	"github.com/GettEngineering/effe/testcustomization"
	"github.com/GettEngineering/effe/testing"
)

func main() {
	settings := generator.DefaultSettigs()
	strategy := strategies.NewChain(strategies.WithServiceObjectName(settings.LocalInterfaceVarname()))
	err := strategy.Register("PostRequestComponent", testcustomization.GenPostRequestComponent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loader := loaders.NewLoader(loaders.WithPackages([]string{"effe", "testcustomization"}))
	err = loader.Register("POST", testcustomization.LoadPostRequestComponent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gen := generator.NewGenerator(
		generator.WithSetttings(settings),
		generator.WithLoader(loader),
		generator.WithStrategy(strategy),
	)

	dslGo, err := ioutil.ReadFile(filepath.Join("testcustomization.go"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testing.UpdateExpectedResult(os.Args[2], gen, map[string][]byte{"github.com/GettEngineering/effe/testcustomization/testcustomization.go": dslGo}, []string{"github.com/GettEngineering/effe/testcustomization"})
}
