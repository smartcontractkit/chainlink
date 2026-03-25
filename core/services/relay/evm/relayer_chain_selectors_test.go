package evm_test

import (
	"math/big"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmmocks "github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	corekeystore "github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

// unsupportedEVMChainID is chosen so [chainselectors.SelectorFromChainId] fails (chain not in embedded selectors).
var unsupportedEVMChainID = new(big.Int).SetUint64(^uint64(0))

func TestNewRelayer_ChainIDNotInChainSelectors(t *testing.T) {
	t.Parallel()

	_, err := chainselectors.SelectorFromChainId(unsupportedEVMChainID.Uint64())
	require.Error(t, err, "precondition: chain ID must be absent from chain-selectors for this test")

	chain := evmmocks.NewChain(t)
	chain.On("ID").Return(unsupportedEVMChainID)
	evmCfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.ChainID = sqlutil.New(unsupportedEVMChainID)
	})
	chain.On("Config").Return(evmCfg)

	ds := sqltest.NewInMemoryDataSource(t)
	ethKS := keystoremocks.NewEth(t)
	ethKS.EXPECT().EnabledAddressesForChain(mock.Anything, mock.Anything).Return([]common.Address{}, nil).Maybe()

	lggr := logger.TestLogger(t)
	_, err = evm.NewRelayer(lggr, chain, evm.RelayerOpts{
		DS:                   ds,
		EVMKeystore:          keys.NewChainStore(corekeystore.NewEthSigner(ethKS, chain.ID()), chain.ID()),
		CSAKeystore:          &core.UnimplementedKeystore{},
		CapabilitiesRegistry: capabilities.NewRegistry(lggr),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain-selectors missing chain id")
}
