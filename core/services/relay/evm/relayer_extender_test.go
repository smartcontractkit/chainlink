package evm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestChainRelayExtenders(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	newId := testutils.NewRandomEVMChainID()
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		one := uint32(1)
		c.EVM[0].MinIncomingConfirmations = &one
		t := true
		c.EVM = append(c.EVM, &toml.EVMConfig{ChainID: ubig.New(newId), Enabled: &t, Chain: toml.Defaults(nil)})
	})
	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db)
	require.NoError(t, kst.Unlock(ctx, cltest.Password))

	opts := evmtest.NewChainRelayExtOpts(t, evmtest.TestChainOpts{DB: db, KeyStore: kst.Eth(), GeneralConfig: cfg})
	opts.GenEthClient = func(*big.Int) evmclient.Client {
		return cltest.NewEthMocksWithStartupAssertions(t)
	}
	relayExtenders, err := evmrelay.NewChainRelayerExtenders(testutils.Context(t), opts)
	require.NoError(t, err)

	require.Equal(t, relayExtenders.Len(), 2)
	relayExtendersInstances := relayExtenders.Slice()
	for _, c := range relayExtendersInstances {
		require.NoError(t, c.Start(testutils.Context(t)))
		require.NoError(t, c.Ready())
	}

	require.NotEqual(t, relayExtendersInstances[0].Chain().ID().String(), relayExtendersInstances[1].Chain().ID().String())

	for _, c := range relayExtendersInstances {
		require.NoError(t, c.Close())
	}

	relayExtendersInstances[0].Chain().Client().(*evmclimocks.Client).AssertCalled(t, "Close")
	relayExtendersInstances[1].Chain().Client().(*evmclimocks.Client).AssertCalled(t, "Close")

	assert.Error(t, relayExtendersInstances[0].Chain().Ready())
	assert.Error(t, relayExtendersInstances[1].Chain().Ready())

	// test extender methods on single instance
	relayExt := relayExtendersInstances[0]
	s, err := relayExt.GetChainStatus(testutils.Context(t))
	assert.NotEmpty(t, s)
	assert.NoError(t, err)
}
