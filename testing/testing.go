package testing

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name                 string
	pkg                  string
	goFiles              map[string][]byte
	wantEffeOutput       []byte
	wantEffeError        bool
	wantEffeErrorStrings []string
}

// nolint:gocognit
func loadTestCase(root string, extGoFiles map[string][]byte) (*testCase, error) {
	name := filepath.Base(root)
	pkg, err := ioutil.ReadFile(filepath.Join(root, "pkg"))
	if err != nil {
		return nil, fmt.Errorf("load test case %s: %v", name, err)
	}
	var wantEffeOutput []byte
	var wantEffeErrorStrings []string
	effeErrb, err := ioutil.ReadFile(filepath.Join(root, "want", "effe_errs.txt"))
	wantEffeError := err == nil
	if wantEffeError {
		wantEffeErrorStrings = strings.Split(string(effeErrb), "\n")
	} else {
		wantEffeOutput, err = ioutil.ReadFile(filepath.Join(root, "want", "effe_gen.go"))
		if err != nil {
			return nil, fmt.Errorf("load test case %s: %v, if this is a new testcase, run updater", name, err)
		}
	}

	goFiles := map[string][]byte{}
	for k, v := range extGoFiles {
		goFiles[k] = v
	}
	goFiles["github.com/GettEngineering/effe/effe.go"] = []byte(SourceDSL)
	err = filepath.Walk(root, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, src)
		if err != nil {
			return err // unlikely
		}
		if info.Mode().IsDir() && rel == "want" {
			// The "want" directory should not be included in goFiles.
			return filepath.SkipDir
		}
		if !info.Mode().IsRegular() || filepath.Ext(src) != ".go" {
			return nil
		}
		data, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		goFiles["example.com/"+filepath.ToSlash(rel)] = data
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load test case %s: %v", name, err)
	}
	return &testCase{
		name:                 name,
		pkg:                  string(bytes.TrimSpace(pkg)),
		wantEffeOutput:       wantEffeOutput,
		goFiles:              goFiles,
		wantEffeError:        wantEffeError,
		wantEffeErrorStrings: wantEffeErrorStrings,
	}, nil
}

// materialize creates a new GOPATH at the given directory, which may or
// may not exist.
func (test *testCase) materialize(gopath string, deps []string) error {
	for name, content := range test.goFiles {
		dst := filepath.Join(gopath, "src", filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
			return fmt.Errorf("materialize GOPATH: %v", err)
		}
		if err := ioutil.WriteFile(dst, content, 0600); err != nil {
			return fmt.Errorf("materialize GOPATH: %v", err)
		}
	}

	// Add go.mod files to example.com and github.com/GettEngineering/effe.
	const importPath = "example.com"
	deps = append(deps, "github.com/GettEngineering/effe")

	requireContent := fmt.Sprintf("module %s\n\nrequire (\n", importPath)
	replaceContent := "replace (\n"
	for _, dep := range deps {
		depLoc := filepath.Join(gopath, "src", filepath.FromSlash(dep))
		requireContent += fmt.Sprintf("  %s v0.1.0\n", dep)
		replaceContent += fmt.Sprintf("  %s => %s\n", dep, depLoc)
		if err := ioutil.WriteFile(filepath.Join(depLoc, "go.mod"), []byte("module "+dep+"\n"), 0600); err != nil {
			return fmt.Errorf("generate go.mod for %s: %v", dep, err)
		}
	}
	replaceContent += ")\n"
	requireContent += ")\n"
	requireContent += replaceContent

	gomod := filepath.Join(gopath, "src", filepath.FromSlash(importPath), "go.mod")
	if err := ioutil.WriteFile(gomod, []byte(requireContent), 0600); err != nil {
		return fmt.Errorf("generate go.mod for %s: %v", gomod, err)
	}
	return nil
}

// scrubError rewrites the given string to remove occurrences of GOPATH/src,
// rewrites OS-specific path separators to slashes, and any line/column
// information to a fixed ":x:y". For example, if the gopath parameter is
// "C:\GOPATH" and running on Windows, the string
// "C:\GOPATH\src\foo\bar.go:15:4" would be rewritten to "foo/bar.go:x:y".
func scrubError(gopath string, s string) string {
	sb := new(strings.Builder)
	query := gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator)
	for {
		// Find next occurrence of source root. This indicates the next path to
		// scrub.
		start := strings.Index(s, query)
		if start == -1 {
			sb.WriteString(s)
			break
		}

		// Find end of file name (extension ".go").
		fileStart := start + len(query)
		fileEnd := strings.Index(s[fileStart:], ".go")
		if fileEnd == -1 {
			// If no ".go" occurs to end of string, further searches will fail too.
			// Break the loop.
			sb.WriteString(s)
			break
		}
		fileEnd += fileStart + 3 // Advance to end of extension.

		// Write out file name and advance scrub position.
		file := s[fileStart:fileEnd]
		if os.PathSeparator != '/' {
			file = strings.Replace(file, string(os.PathSeparator), "/", -1)
		}
		sb.WriteString(s[:start])
		sb.WriteString(file)
		s = s[fileEnd:]

		// Peek past to see if there is line/column info.
		linecol, linecolLen := scrubLineColumn(s)
		sb.WriteString(linecol)
		s = s[linecolLen:]
	}
	return sb.String()
}

func scrubLineColumn(s string) (replacement string, n int) {
	if !strings.HasPrefix(s, ":") {
		return "", 0
	}
	// Skip first colon and run of digits.
	for n++; len(s) > n && '0' <= s[n] && s[n] <= '9'; {
		n++
	}
	if n == 1 {
		// No digits followed colon.
		return "", 0
	}

	// Start on column part.
	if !strings.HasPrefix(s[n:], ":") {
		return ":x", n
	}
	lineEnd := n
	// Skip second colon and run of digits.
	for n++; len(s) > n && '0' <= s[n] && s[n] <= '9'; {
		n++
	}
	if n == lineEnd+1 {
		// No digits followed second colon.
		return ":x", lineEnd
	}
	return ":x:y", n
}

//nolint:gocognit
func RunTests(t *testing.T, gen Generator, testRoot string, goFiles map[string][]byte, deps []string) {
	testdataEnts, err := ioutil.ReadDir(testRoot) // ReadDir sorts by name.
	if err != nil {
		t.Fatal(err)
	}
	tests := make([]*testCase, 0, len(testdataEnts))
	for _, ent := range testdataEnts {
		name := ent.Name()
		if !ent.IsDir() || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}

		test, err := loadTestCase(filepath.Join(testRoot, name), goFiles)
		if err != nil {
			t.Error(err)
			continue
		}
		tests = append(tests, test)
	}
	ctx := context.Background()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			gopath, err := ioutil.TempDir("", "effe_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(gopath)
			gopath, err = filepath.EvalSymlinks(gopath)
			if err != nil {
				t.Fatal(err)
			}
			if err := test.materialize(gopath, deps); err != nil {
				t.Fatal(err)
			}
			wd := filepath.Join(gopath, "src", "example.com")

			gens, errs := gen.Generate(ctx, wd, append(os.Environ(), "GOPATH="+gopath), []string{test.pkg})
			assert.Empty(t, errs)
			for _, gen := range gens {
				if test.wantEffeError {
					assert.Greater(t, len(gen.Errs), 0)
					gotErrStrings := make([]string, len(gen.Errs))
					for i, e := range gen.Errs {
						gotErrStrings[i] = scrubError(gopath, e.Error())
					}
					if diff := cmp.Diff(gotErrStrings, test.wantEffeErrorStrings); diff != "" {
						t.Errorf("Errors didn't match expected errors from effe_errors.txt:\n%s", diff)
					}
				} else {
					assert.Empty(t, gen.Errs)
					assert.NotEmpty(t, gen.OutputPath)
					genContent, err := ioutil.ReadFile(gen.OutputPath)
					assert.NoError(t, err)
					if !bytes.Equal(genContent, test.wantEffeOutput) {
						gotS, wantS := string(genContent), string(test.wantEffeOutput)
						diff := cmp.Diff(strings.Split(gotS, "\n"), strings.Split(wantS, "\n"))
						t.Fatalf("effe output differs from golden file. \n*** got:\n%s\n\n*** want:\n%s\n\n*** diff:\n%s", gotS, wantS, diff)
					}
				}
			}
		})
	}
}
