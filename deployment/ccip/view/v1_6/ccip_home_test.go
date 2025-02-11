package v1_6

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCCIPHomeView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.Chains[e.AllChainSelectors()[0]]
	_, tx, cr, err := capabilities_registry.DeployCapabilitiesRegistry(
		chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	_, tx, ch, err := ccip_home.DeployCCIPHome(
		chain.DeployerKey, chain.Client, cr.Address())
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	v, err := GenerateCCIPHomeView(cr, ch)
	require.NoError(t, err)
	assert.Equal(t, "CCIPHome 1.6.0", v.TypeAndVersion)

	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
