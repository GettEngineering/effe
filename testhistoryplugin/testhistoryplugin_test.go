package testhistoryplugin

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/GettEngineering/effe/generator"
	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/plugins"
	"github.com/GettEngineering/effe/strategies"
	eTesting "github.com/GettEngineering/effe/testing"
	"github.com/stretchr/testify/assert"
)

func TestHistoryPlugin(t *testing.T) {
	settings := generator.DefaultSettigs()
	strategy := strategies.NewChain(
		strategies.WithServiceObjectName(settings.LocalInterfaceVarname()),
		strategies.Use(plugins.NewHistory()),
	)
	gen := generator.NewGenerator(
		generator.WithSetttings(settings),
		generator.WithLoader(loaders.NewLoader(loaders.WithPackages([]string{"effe"}))),
		generator.WithStrategy(strategy),
	)

	dslGo, err := ioutil.ReadFile(filepath.Join("gen.go"))
	assert.NoError(t, err)

	eTesting.RunTests(t, gen, "testdata", map[string][]byte{"github.com/GettEngineering/effe/testhistoryplugin/gen.go": dslGo}, []string{"github.com/GettEngineering/effe/testhistoryplugin"})
}
