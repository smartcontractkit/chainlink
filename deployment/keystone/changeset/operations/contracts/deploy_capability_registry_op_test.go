package contracts_test

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
)

func Test_DeployRegistryOp(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	b := optest.NewBundle(t)
	deps := contracts.DeployCapabilityRegistryOpDeps{
		Env: &env,
	}
	input := contracts.DeployCapabilityRegistryInput{
		ChainSelector: registrySel,
	}

	got, err := operations.ExecuteOperation(b, contracts.DeployCapabilityRegistryOp, deps, input)
	require.NoError(t, err)
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, 1)

	fmt.Println(env.DataStore.Addresses())
}
