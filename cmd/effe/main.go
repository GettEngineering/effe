package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/GettEngineering/effe/drawer"
	"github.com/GettEngineering/effe/generator"
	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/strategies"
	"github.com/GettEngineering/effe/types"
)

const (
	version = "0.1.0"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("effe: ")
	log.SetOutput(os.Stderr)

	showVerstionPtr := flag.Bool("v", false, "show current version of effe")
	drawPtr := flag.Bool("d", false, "draw diagrams for business flows")
	drawOutPtr := flag.String("out", "graphs", "draw output directory")
	flag.Parse()
	if showVerstionPtr != nil && *showVerstionPtr {
		showVersion()
		return
	}

	d, err := os.Getwd()
	if err != nil {
		log.Printf("can't get path of current directory: %s", err)
		os.Exit(2)
	}

	settings := generator.DefaultSettigs()
	gen := generator.NewGenerator(
		generator.WithSetttings(settings),
		generator.WithLoader(loaders.NewLoader(loaders.WithPackages([]string{"effe"}))),
		generator.WithDrawer(drawer.NewDrawer()),
		generator.WithStrategy(
			strategies.NewChain(strategies.WithServiceObjectName(settings.LocalInterfaceVarname())),
		),
	)

	var (
		errs       []error
		genResults []types.GenerateResult
	)

	if drawPtr != nil && *drawPtr {
		if drawOutPtr == nil {
			log.Println("directory for output is not set")
			os.Exit(2)
		}
		genResults, errs = gen.GenerateDiagram(context.Background(), d, os.Environ(), []string{"."}, *drawOutPtr)
	} else {
		genResults, errs = gen.Generate(context.Background(), d, os.Environ(), []string{"."})
	}

	for index, err := range errs {
		log.Printf("failed generate: %s\n", err)
		if index == len(errs)-1 {
			os.Exit(2)
		}
	}

	for _, res := range genResults {
		if len(res.Errs) > 0 {
			for _, err := range res.Errs {
				log.Printf("failed generate: %s\n", err)
			}
		} else {
			log.Printf("wrote %s\n", res.OutputPath)
		}
	}
}

func showVersion() {
	log.Println(version)
}
