package changeset_test

import (
	"embed"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

//go:embed testdata/*
var testFS embed.FS

func TestImportToAddressbook(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	resp, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			changeset.ImportToAddressbookChangeset,
			types.ImportAddressesConfig{
				ChainSelector: chainSelector,
				InputFileName: "testdata/import_addresses.json",
				InputFS:       testFS,
			},
		),
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
	tv, _ := resp.ExistingAddresses.AddressesForChain(chainSelector)
	require.Len(t, tv, 2)
}
