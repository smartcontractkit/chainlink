package ccipdeployment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDeployMockRMN(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      1,
	})
	ab := deployment.NewMemoryAddressBook()

	contract, err := DeployMockRMN(e, e.Chains[e.AllChainSelectors()[0]], ab)
	require.NoError(t, err)
	require.NotNil(t, contract)
	require.NotNil(t, contract.Address)
}

func TestDeployWETH9(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      1,
	})
	ab := deployment.NewMemoryAddressBook()
	contract, err := DeployWrappedNative(e, e.Chains[e.AllChainSelectors()[0]], ab)
	require.NoError(t, err)
	require.NotNil(t, contract)
	require.NotNil(t, contract.Address)
}

func TestDeployLegacyContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      1,
	})
	ab := deployment.NewMemoryAddressBook()

	err := DeployLegacyContracts(e, e.Chains[e.AllChainSelectors()[0]], ab)
	require.NoError(t, err)
	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.View(e.AllChainSelectors())
	require.NoError(t, err)
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}
