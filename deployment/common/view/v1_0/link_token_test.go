package v1_0

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLinkTokenView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.Chains[e.AllChainSelectors()[0]]
	_, tx, lt, err := link_token.DeployLinkToken(chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
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
	tx, err = lt.GrantMintAndBurnRoles(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	tx, err = lt.Mint(chain.DeployerKey, chain.DeployerKey.From, big.NewInt(100))
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
