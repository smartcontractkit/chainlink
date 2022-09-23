package docs

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
)

func TestDoc(t *testing.T) {
	d := toml.NewDecoder(strings.NewReader(docsTOML))
	d.DisallowUnknownFields() // Ensure no extra fields
	var c chainlink.Config
	err := d.Decode(&c)
	var strict *toml.StrictMissingError
	if err != nil && strings.Contains(err.Error(), "undecoded keys: ") {
		t.Errorf("Docs contain extra fields: %v", err)
	} else if errors.As(err, &strict) {
		t.Fatal("StrictMissingError:", strict.String())
	} else {
		require.NoError(t, err)
	}

	cfgtest.AssertFieldsNotNil(t, c)

	t.Run("EVM", func(t *testing.T) {
		fallbackDefaults, _ := evmcfg.Defaults(nil)

		var defaults chainlink.Config
		require.NoError(t, cfgtest.DocDefaultsOnly(strings.NewReader(chainsEVMTOML), &defaults))
		docDefaults := defaults.EVM[0].Chain

		// clean up KeySpecific as a special case
		require.Equal(t, 1, len(docDefaults.KeySpecific))
		require.Equal(t, evmcfg.KeySpecific{}, docDefaults.KeySpecific[0])
		docDefaults.KeySpecific = nil

		fallback, err := toml.Marshal(fallbackDefaults)
		require.NoError(t, err)
		doc, err := toml.Marshal(docDefaults)
		require.NoError(t, err)
		fs, ds := string(fallback), string(doc)

		assert.Equal(t, fs, ds, diff.Diff(fs, ds))
	})
}

var (
	//go:embed testdata/example.toml
	exampleTOML string
	//go:embed testdata/example.md
	exampleMarkdown string
)

func Test_generateDocs(t *testing.T) {
	got, err := generateDocs(exampleTOML)
	require.NoError(t, err)
	assert.Equal(t, exampleMarkdown, got)
}
