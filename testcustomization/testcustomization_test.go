package testcustomization

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/GettEngineering/effe/generator"
	"github.com/GettEngineering/effe/loaders"
	"github.com/GettEngineering/effe/strategies"
	eTesting "github.com/GettEngineering/effe/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomization(t *testing.T) {
	settings := generator.DefaultSettigs()
	strategy := strategies.NewChain(strategies.WithServiceObjectName(settings.LocalInterfaceVarname()))
	err := strategy.Register("PostRequestComponent", GenPostRequestComponent)
	require.NoError(t, err)

	loader := loaders.NewLoader(loaders.WithPackages([]string{"effe", "testcustomization"}))
	err = loader.Register("POST", LoadPostRequestComponent)
	require.NoError(t, err)

	gen := generator.NewGenerator(
		generator.WithSetttings(settings),
		generator.WithLoader(loader),
		generator.WithStrategy(strategy),
	)

	dslGo, err := ioutil.ReadFile(filepath.Join("testcustomization.go"))
	assert.NoError(t, err)

	eTesting.RunTests(t, gen, "testdata", map[string][]byte{"github.com/GettEngineering/effe/testcustomization/testcustomization.go": dslGo}, []string{"github.com/GettEngineering/effe/testcustomization"})
}
