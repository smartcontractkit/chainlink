package v1_6

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCCIPHomeView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]]
	_, tx, cr, err := capabilities_registry.DeployCapabilitiesRegistry(
		chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	_, tx, ch, err := ccip_home.DeployCCIPHome(
		chain.DeployerKey, chain.Client, cr.Address())
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	v, err := GenerateCCIPHomeView(cr, ch)
	require.NoError(t, err)
	assert.Equal(t, "CCIPHome 1.6.0", v.TypeAndVersion)

	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
