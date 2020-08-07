package main

import (
	"bytes"
	"go/build"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	effePath := build.Default.GOPATH + "/src/github.com/GettEngineering/effe/"
	cPath, err := filepath.EvalSymlinks(filepath.Join(effePath, "effe.go"))
	if err != nil {
		log.Fatalf("EvalSymlinks error: %s", err)
	}
	if cPath == "" {
		log.Fatal("source path can't be empty")
	}
	sourceData, err := ioutil.ReadFile(cPath)
	if err != nil {
		log.Fatalf("error reading source file %s: %s", cPath, err)
	}

	if len(sourceData) == 0 {
		log.Fatalf("source file %s is empty", cPath)
	}

	var data bytes.Buffer

	data.WriteString("// Code generated automatically. DO NOT EDIT.\n")
	data.WriteString("\npackage testing\n\n")
	data.WriteString("var SourceDSL = `")
	data.WriteString(string(sourceData))
	data.WriteString("`")

	err = ioutil.WriteFile(filepath.Join(effePath, "testing/dsl_source_gen.go"), data.Bytes(), 0600)
	if err != nil {
		log.Fatalf("error writing dst file testing/dsl_source_gen.go: %s", err)
	}
}
