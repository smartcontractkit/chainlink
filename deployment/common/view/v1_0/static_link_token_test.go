package v1_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestStaticLinkTokenView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]]
	_, tx, lt, err := link_token_interface.DeployLinkToken(chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	v, err := GenerateStaticLinkTokenView(lt)
	require.NoError(t, err)

	assert.Equal(t, v.Owner, common.HexToAddress("0x0")) // Ownerless
	assert.Equal(t, "StaticLinkToken 1.0.0", v.TypeAndVersion)
	assert.Equal(t, uint8(18), v.Decimals)
	assert.Equal(t, "1000000000000000000000000000", v.Supply.String())
}
