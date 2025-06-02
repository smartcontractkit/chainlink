package v1_5

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGenerateRMNView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]]
	cfg := rmn_contract.RMNConfig{
		Voters: []rmn_contract.RMNVoter{
			{
				BlessVoteAddr: chain.DeployerKey.From,
				CurseVoteAddr: common.HexToAddress("0x3"),
				BlessWeight:   1,
				CurseWeight:   1,
			},
			{
				BlessVoteAddr: common.HexToAddress("0x1"),
				CurseVoteAddr: common.HexToAddress("0x2"),
				BlessWeight:   1,
				CurseWeight:   1,
			},
		},
		BlessWeightThreshold: uint16(2),
		CurseWeightThreshold: uint16(1),
	}
	_, tx, c, err := rmn_contract.DeployRMNContract(
		chain.DeployerKey, chain.Client, cfg)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	v, err := GenerateRMNView(c)
	require.NoError(t, err)
	assert.Equal(t, v.Owner, chain.DeployerKey.From)
	assert.Equal(t, "RMN 1.5.0", v.TypeAndVersion)
	assert.Equal(t, uint32(1), v.ConfigDetails.Version)
	assert.Equal(t, v.ConfigDetails.Config, cfg)
	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
