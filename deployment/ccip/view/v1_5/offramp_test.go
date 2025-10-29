package v1_5

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestOffRampView(t *testing.T) {
	t.Parallel()

	srcSelector := chainsel.TEST_90000001.Selector
	dstSelector := chainsel.TEST_90000002.Selector
	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{srcSelector, dstSelector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[srcSelector]

	_, tx, c, err := commit_store.DeployCommitStore(
		chain.DeployerKey, chain.Client, commit_store.CommitStoreStaticConfig{
			ChainSelector:       dstSelector,
			SourceChainSelector: srcSelector,
			OnRamp:              common.HexToAddress("0x4"),
			RmnProxy:            common.HexToAddress("0x1"),
		})
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	sc := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
		ChainSelector:       dstSelector,
		SourceChainSelector: srcSelector,
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
