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

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	config "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
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

	var defaults chainlink.Config
	require.NoError(t, cfgtest.DocDefaultsOnly(strings.NewReader(docsTOML), &defaults, config.DecodeTOML))

	t.Run("EVM", func(t *testing.T) {
		fallbackDefaults := evmcfg.Defaults(nil)
		docDefaults := defaults.EVM[0].Chain

		require.Equal(t, "", *docDefaults.ChainType)
		docDefaults.ChainType = nil

		// clean up KeySpecific as a special case
		require.Equal(t, 1, len(docDefaults.KeySpecific))
		ks := evmcfg.KeySpecific{Key: new(ethkey.EIP55Address),
			GasEstimator: evmcfg.KeySpecificGasEstimator{PriceMax: new(assets.Wei)}}
		require.Equal(t, ks, docDefaults.KeySpecific[0])
		docDefaults.KeySpecific = nil

		// per-job limits are nilable
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.OCR)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.DR)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.Keeper)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.VRF)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.FM)
		docDefaults.GasEstimator.LimitJobType = evmcfg.GasLimitJobType{}

		// EIP1559FeeCapBufferBlocks doesn't have a constant default - it is derived from another field
		require.Zero(t, *docDefaults.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks)
		docDefaults.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks = nil

		// addresses w/o global values
		require.Zero(t, *docDefaults.FlagsContractAddress)
		require.Zero(t, *docDefaults.LinkContractAddress)
		require.Zero(t, *docDefaults.OperatorFactoryAddress)
		docDefaults.FlagsContractAddress = nil
		docDefaults.LinkContractAddress = nil
		docDefaults.OperatorFactoryAddress = nil

		assertTOML(t, fallbackDefaults, docDefaults)
	})

	t.Run("Cosmos", func(t *testing.T) {
		var fallbackDefaults cosmos.CosmosConfig
		fallbackDefaults.SetDefaults()

		assertTOML(t, fallbackDefaults.Chain, defaults.Cosmos[0].Chain)
	})

	t.Run("Solana", func(t *testing.T) {
		var fallbackDefaults solana.SolanaConfig
		fallbackDefaults.SetDefaults()

		assertTOML(t, fallbackDefaults.Chain, defaults.Solana[0].Chain)
	})

	t.Run("Starknet", func(t *testing.T) {
		var fallbackDefaults starknet.StarknetConfig
		fallbackDefaults.SetDefaults()

		assertTOML(t, fallbackDefaults.Chain, defaults.Starknet[0].Chain)
	})
}

func assertTOML[T any](t *testing.T, fallback, docs T) {
	t.Helper()
	t.Logf("fallback: %#v", fallback)
	t.Logf("docs: %#v", docs)
	fb, err := toml.Marshal(fallback)
	require.NoError(t, err)
	db, err := toml.Marshal(docs)
	require.NoError(t, err)
	fs, ds := string(fb), string(db)
	assert.Equal(t, fs, ds, diff.Diff(fs, ds))
}

var (
	//go:embed testdata/example.toml
	exampleTOML string
	//go:embed testdata/example.md
	exampleMarkdown string
)

func Test_generateDocs(t *testing.T) {
	got, err := generateDocs(exampleTOML, `[//]: # (Generated - DO NOT EDIT.)
`, `Bar = 7 # Required
`)
	require.NoError(t, err)
	assert.Equal(t, exampleMarkdown, got)
}
