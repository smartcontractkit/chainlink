package v1_2

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGeneratePriceRegistryView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.Chains[e.AllChainSelectors()[0]]
	f1, f2 := common.HexToAddress("0x1"), common.HexToAddress("0x2")
	_, tx, c, err := price_registry_1_2_0.DeployPriceRegistry(
		chain.DeployerKey, chain.Client, []common.Address{chain.DeployerKey.From}, []common.Address{f1, f2}, uint32(10))
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
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
