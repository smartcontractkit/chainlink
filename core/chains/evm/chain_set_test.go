package evm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestChainSet(t *testing.T) {
	t.Parallel()

	newId := testutils.NewRandomEVMChainID()
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		one := uint32(1)
		c.EVM[0].MinIncomingConfirmations = &one
		t := true
		c.EVM = append(c.EVM, &v2.EVMConfig{ChainID: utils.NewBig(newId), Enabled: &t, Chain: v2.Defaults(nil)})
	})
	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg)
	require.NoError(t, kst.Unlock(cltest.Password))

	opts := evmtest.NewChainSetOpts(t, evmtest.TestChainOpts{DB: db, KeyStore: kst.Eth(), GeneralConfig: cfg})
	opts.GenEthClient = func(*big.Int) evmclient.Client {
		return cltest.NewEthMocksWithStartupAssertions(t)
	}
	chainSet, err := evm.NewTOMLChainSet(testutils.Context(t), opts)
	require.NoError(t, err)

	require.NoError(t, chainSet.Start(testutils.Context(t)))
	require.NoError(t, chainSet.Chains()[0].Ready())

	chains := chainSet.Chains()
	require.Equal(t, 2, len(chains))
	require.NotEqual(t, chains[0].ID().String(), chains[1].ID().String())

	assert.NoError(t, chains[0].Ready())
	assert.NoError(t, chains[1].Ready())

	chainSet.Close()

	chains[0].Client().(*evmclimocks.Client).AssertCalled(t, "Close")
	chains[1].Client().(*evmclimocks.Client).AssertCalled(t, "Close")

	assert.Error(t, chains[0].Ready())
	assert.Error(t, chains[1].Ready())
}
