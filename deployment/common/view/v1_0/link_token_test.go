package v1_0

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLinkTokenView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]]
	_, tx, lt, err := link_token.DeployLinkToken(chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	testLinkTokenViewWithChain(t, chain, lt)
}

func TestLinkTokenViewZk(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		ZkChains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]]
	_, _, lt, err := link_token.DeployLinkTokenZk(nil, chain.ClientZkSyncVM, chain.DeployerKeyZkSyncVM, chain.Client)
	require.NoError(t, err)

	testLinkTokenViewWithChain(t, chain, lt)
}

func testLinkTokenViewWithChain(t *testing.T, chain cldf_evm.Chain, lt *link_token.LinkToken) {
	v, err := GenerateLinkTokenView(lt)
	require.NoError(t, err)

	assert.Equal(t, v.Owner, chain.DeployerKey.From)
	assert.Equal(t, "LinkToken 1.0.0", v.TypeAndVersion)
	assert.Equal(t, uint8(18), v.Decimals)
	// Initially nothing minted and no minters/burners.
	assert.Equal(t, "0", v.Supply.String())
	require.Empty(t, v.Minters)
	require.Empty(t, v.Burners)

	// Add some minters
	tx, err := lt.GrantMintAndBurnRoles(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	tx, err = lt.Mint(chain.DeployerKey, chain.DeployerKey.From, big.NewInt(100))
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	v, err = GenerateLinkTokenView(lt)
	require.NoError(t, err)

	assert.Equal(t, "100", v.Supply.String())
	require.Len(t, v.Minters, 1)
	require.Equal(t, v.Minters[0].String(), chain.DeployerKey.From.String())
	require.Len(t, v.Burners, 1)
	require.Equal(t, v.Burners[0].String(), chain.DeployerKey.From.String())
}
