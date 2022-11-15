package evm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestUpdateKeySpecificMaxGasPrice_NewEntry(t *testing.T) {
	t.Parallel()

	address := common.HexToAddress("0x1234567890")
	price := assets.NewWeiI(12345)
	updater := evm.UpdateKeySpecificMaxGasPrice(address, price)
	config := types.ChainCfg{}

	err := updater(&config)

	require.NoError(t, err)
	require.NotNil(t, config.KeySpecific)
	require.Equal(t, price, config.KeySpecific[address.Hex()].EvmMaxGasPriceWei)
}

func TestUpdateKeySpecificMaxGasPrice_ExistingEntry(t *testing.T) {
	t.Parallel()

	address := common.HexToAddress("0x1234567890")
	price1 := assets.NewWeiI(12345)
	price2 := assets.NewWeiI(54321)
	updater := evm.UpdateKeySpecificMaxGasPrice(address, price2)
	config := types.ChainCfg{
		KeySpecific: map[string]types.ChainCfg{
			"0x1234567890": {
				EvmMaxGasPriceWei: (*assets.Wei)(price1),
			},
		},
	}

	err := updater(&config)

	require.NoError(t, err)
	require.NotNil(t, config.KeySpecific)
	require.Equal(t, price2, config.KeySpecific[address.Hex()].EvmMaxGasPriceWei)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestUpdateConfig(t *testing.T) {
	t.Parallel()

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg)
	require.NoError(t, kst.Unlock(cltest.Password))

	chainSet := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, KeyStore: kst.Eth(), GeneralConfig: cfg, Client: ethClient})
	address := common.HexToAddress("0x1234567890")
	price := assets.NewWeiI(12345)
	updater := evm.UpdateKeySpecificMaxGasPrice(address, price)

	chain, err := chainSet.Get(&cltest.FixtureChainID)
	require.NoError(t, err)

	err = chainSet.UpdateConfig(&cltest.FixtureChainID, updater)
	require.NoError(t, err)

	require.Equal(t, price, chain.Config().KeySpecificMaxGasPriceWei(address))
}

func TestAddClose(t *testing.T) {
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

	opts, _, _ := evmtest.NewChainSetOpts(t, evmtest.TestChainOpts{DB: db, KeyStore: kst.Eth(), GeneralConfig: cfg})
	opts.GenEthClient = func(*big.Int) evmclient.Client {
		return cltest.NewEthMocksWithStartupAssertions(t)
	}
	cfgs := cfg.(v2.HasEVMConfigs).EVMConfigs()
	chainSet, err := evm.NewTOMLChainSet(testutils.Context(t), opts, cfgs)
	require.NoError(t, err)

	require.NoError(t, chainSet.Start(testutils.Context(t)))
	require.NoError(t, chainSet.Chains()[0].Ready())

	chains := chainSet.Chains()
	require.Equal(t, 2, len(chains))
	require.NotEqual(t, chains[0].ID().String(), chains[1].ID().String())

	assert.NoError(t, chains[0].Ready())
	assert.NoError(t, chains[1].Ready())

	chainSet.Close()

	chains[0].Client().(*evmmocks.Client).AssertCalled(t, "Close")
	chains[1].Client().(*evmmocks.Client).AssertCalled(t, "Close")

	assert.Error(t, chains[0].Ready())
	assert.Error(t, chains[1].Ready())
}
