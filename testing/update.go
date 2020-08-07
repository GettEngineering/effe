package testing

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GettEngineering/effe/types"
)

type Generator interface {
	Generate(context.Context, string, []string, []string) ([]types.GenerateResult, []error)
}

//nolint:funlen,gocognit
func UpdateExpectedResult(testRoot string, gen Generator, goFiles map[string][]byte, deps []string) {
	testdataEnts, err := ioutil.ReadDir(testRoot) // ReadDir sorts by name.
	if err != nil {
		log.Fatal(err)
	}
	tests := make([]*testCase, 0, len(testdataEnts))
	for _, ent := range testdataEnts {
		name := ent.Name()
		if !ent.IsDir() || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}

		test, err := loadTestCase(filepath.Join(testRoot, name), goFiles)
		if err != nil {
			log.Println(err)
			continue
		}
		tests = append(tests, test)
	}
	ctx := context.Background()
	for _, test := range tests {
		test := test
		gopath, err := ioutil.TempDir("", "effe_test")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(gopath)

		gopath, err = filepath.EvalSymlinks(gopath)
		if err != nil {
			log.Fatal(err)
		}
		if err = test.materialize(gopath, deps); err != nil {
			log.Fatal(err)
		}
		wd := filepath.Join(gopath, "src", "example.com")
		errMsg := make([]string, 0)
		genResults, errs := gen.Generate(ctx, wd, append(os.Environ(), "GOPATH="+gopath), []string{test.pkg})
		for _, err := range errs {
			errMsg = append(errMsg, err.Error())
		}
		for _, res := range genResults {
			if len(res.Errs) > 0 {
				for _, err := range res.Errs {
					errMsg = append(errMsg, err.Error())
				}
			}
		}

		if len(errMsg) > 0 {
			effeErrsFile := filepath.Join(testRoot, test.name, "want", "effe_errs.txt")
			formattedErrs := make([]string, len(errMsg))
			for errIndex, err := range errMsg {
				index := strings.Index(err, gopath)
				if index == -1 {
					formattedErrs[errIndex] = err
					continue
				}
				formattedErrs[errIndex] = string([]byte(err)[:index]) + test.pkg + "/effe.go:x:y"
			}

			err = ioutil.WriteFile(effeErrsFile, []byte(strings.Join(formattedErrs, "\n")), 0600)
			if err != nil {
				fmt.Printf("can't create file with errors %s: %s\n", effeErrsFile, err)
				os.Exit(1)
			}
			continue
		}

		effeGenFile := filepath.Join(gopath, "src", test.pkg, "effe_gen.go")
		generatedCode, err := ioutil.ReadFile(effeGenFile)
		if err != nil {
			fmt.Printf("can't read generated file %s: %s\n", effeGenFile, err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(filepath.Join(testRoot, test.name, "want", "effe_gen.go"), []byte(string(generatedCode)), 0600)
		if err != nil {
			fmt.Printf("can't copy data from generated file %s: %s\n", effeGenFile, err)
			os.Exit(1)
		}
	}
}
