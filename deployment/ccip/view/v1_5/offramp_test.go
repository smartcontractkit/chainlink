package v1_5

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestOffRampView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))[0]]
	_, tx, c, err := commit_store.DeployCommitStore(
		chain.DeployerKey, chain.Client, commit_store.CommitStoreStaticConfig{
			ChainSelector:       chainsel.TEST_90000002.Selector,
			SourceChainSelector: chainsel.TEST_90000001.Selector,
			OnRamp:              common.HexToAddress("0x4"),
			RmnProxy:            common.HexToAddress("0x1"),
		})
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	sc := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
		ChainSelector:       chainsel.TEST_90000002.Selector,
		SourceChainSelector: chainsel.TEST_90000001.Selector,
		RmnProxy:            common.HexToAddress("0x1"),
		CommitStore:         c.Address(),
		TokenAdminRegistry:  common.HexToAddress("0x3"),
		OnRamp:              common.HexToAddress("0x4"),
	}
	rl := evm_2_evm_offramp.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  big.NewInt(100),
		Rate:      big.NewInt(10),
	}
	_, tx, c2, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
		chain.DeployerKey, chain.Client, sc, rl)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	v, err := GenerateOffRampView(c2)
	require.NoError(t, err)
	assert.Equal(t, v.StaticConfig, sc)
	assert.Equal(t, "EVM2EVMOffRamp 1.5.0", v.TypeAndVersion)
	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
