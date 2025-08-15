package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployWorkflowRegistry(t *testing.T) {
	lggr := logger.Test(t)

	// Create a minimal environment with one EVM chain
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  0,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	t.Logf("Environment created with operations bundle")

	// Get the chain selector for the deployment
	chainSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.NotEmpty(t, chainSelectors, "should have at least one EVM chain")
	chainSelector := chainSelectors[0]
	t.Logf("Using chain selector: %d", chainSelector)

	// Apply the changeset to deploy the V2 workflow registry
	t.Log("Starting changeset application...")
	changesetOutput, err := DeployWorkflowRegistry{}.Apply(env, DeployWorkflowRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "test-workflow-registry-v2",
	})
	t.Logf("Changeset result: err=%v, output=%v", err, changesetOutput)

	if err != nil {
		t.Fatalf("changeset apply failed: %v", err)
	}
	require.NotNil(t, changesetOutput, "changeset output should not be nil")
	t.Logf("Changeset applied successfully")

	// Verify the datastore contains the deployed contract
	require.NotNil(t, changesetOutput.DataStore, "datastore should not be nil")
	addresses := changesetOutput.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-workflow-registry-v2"))
	t.Logf("Found %d addresses with qualifier", len(addresses))
	require.Len(t, addresses, 1, "expected exactly one deployed contract with the test qualifier")

	// Verify the address is for the correct chain
	deployedAddress := addresses[0]
	require.Equal(t, chainSelector, deployedAddress.ChainSelector, "deployed contract should be on the correct chain")
	require.NotEmpty(t, deployedAddress.Address, "deployed contract address should not be empty")

	// Verify the contract type is correct
	require.Equal(t, datastore.ContractType("WorkflowRegistry"), deployedAddress.Type, "contract type should be WorkflowRegistry")
	require.NotNil(t, deployedAddress.Version, "contract version should be set")

	// Verify reports are generated
	require.NotNil(t, changesetOutput.Reports, "reports should be present")
	require.Len(t, changesetOutput.Reports, 1, "should have exactly one operation report")
}
