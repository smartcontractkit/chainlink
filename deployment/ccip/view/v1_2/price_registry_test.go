package v1_2

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
)

func TestGeneratePriceRegistryView(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[selector]

	f1, f2 := common.HexToAddress("0x1"), common.HexToAddress("0x2")
	_, tx, c, err := price_registry_1_2_0.DeployPriceRegistry(
		chain.DeployerKey, chain.Client, []common.Address{chain.DeployerKey.From}, []common.Address{f1, f2}, uint32(10))
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	v, err := GeneratePriceRegistryView(c)
	require.NoError(t, err)
	assert.Equal(t, v.Owner, chain.DeployerKey.From)
	assert.Equal(t, "PriceRegistry 1.2.0", v.TypeAndVersion)
	assert.Equal(t, []common.Address{f1, f2}, v.FeeTokens)
	assert.Equal(t, "10", v.StalenessThreshold)
	assert.Equal(t, []common.Address{chain.DeployerKey.From}, v.Updaters)
	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
