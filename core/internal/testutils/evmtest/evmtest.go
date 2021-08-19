package evmtest

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func verifyMatchingChainIDs(t testing.TB, n *big.Int, m *big.Int) {
	require.Equal(t, n.Cmp(m), 0, "expected chain IDs to match")
}

type TestChainOpts struct {
	Client         eth.Client
	LogBroadcaster log.Broadcaster
	GeneralConfig  config.GeneralConfig
	ChainCfg       evmtypes.ChainCfg
	DB             *gorm.DB
}

func NewChainScopedConfig(t testing.TB, cfg config.GeneralConfig) evmconfig.ChainScopedConfig {
	return evmconfig.NewChainScopedConfig(nil, logger.Default, cfg, evmtypes.Chain{ID: *utils.NewBigI(0)})
}

// NewChainCollection returns a simple chain collection with one chain and
// allows to mock client/config on that chain
func NewChainCollection(t testing.TB, testopts TestChainOpts) evm.ChainCollection {
	opts := evm.ChainCollectionOpts{
		Config: testopts.GeneralConfig,
		Logger: logger.Default.With("testname", t.Name()),
		DB:     testopts.DB,
	}

	if testopts.Client != nil {
		opts.GenEthClient = func(c evmtypes.Chain) eth.Client {
			return testopts.Client
		}
	}
	if testopts.LogBroadcaster != nil {
		opts.GenLogBroadcaster = func(c evmtypes.Chain) log.Broadcaster {
			return testopts.LogBroadcaster
		}
	}

	chains := []evmtypes.Chain{
		{
			ID:  *utils.NewBigI(0),
			Cfg: testopts.ChainCfg,
			Nodes: []evmtypes.Node{{
				Name:       "evm-test-only-0",
				EVMChainID: *utils.NewBigI(0),
				WSURL:      "ws://example.invalid",
			}},
		},
	}

	cc, err := evm.NewChainCollection(opts, chains)
	require.NoError(t, err)
	return cc
}

func MustGetDefaultChain(t testing.TB, cc evm.ChainCollection) evm.Chain {
	chain, err := cc.Default()
	require.NoError(t, err)
	return chain
}

func MustInsertChainWithNode(t testing.TB, db *gorm.DB, chain evmtypes.Chain) evmtypes.Chain {
	err := db.Create(&chain).Error
	require.NoError(t, err)
	return chain
}
