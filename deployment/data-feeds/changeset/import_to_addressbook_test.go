package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestImportToAddressbook(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	resp, err := changeset.ImportToAddressbookChangeset(env, types.ImportToAddressbookConfig{
		ChainSelector: chainSelector,
		InputFileName: "testdata/import_addresses.json",
		InputFS:       testFS,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	tv, _ := resp.AddressBook.AddressesForChain(chainSelector)
	require.Len(t, tv, 2)
}
