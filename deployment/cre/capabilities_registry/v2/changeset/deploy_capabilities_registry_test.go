package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/cre"
)

func TestDeployCapabilitiesRegistry(t *testing.T) {
	lggr := logger.Test(t)
	env, chainSelector := cre.BuildMinimalEnvironment(t, lggr)

	// Apply the changeset to deploy the V2 capabilities registry
	t.Log("Starting changeset application...")
	changesetOutput, err := DeployCapabilitiesRegistry{}.Apply(env, DeployCapabilitiesRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "test-capabilities-registry-v2",
	})
	t.Logf("Changeset result: err=%v, output=%v", err, changesetOutput)

	if err != nil {
		t.Fatalf("changeset apply failed: %v", err)
	}
	require.NotNil(t, changesetOutput, "changeset output should not be nil")
	t.Logf("Changeset applied successfully")

	// Verify the datastore contains the deployed contract
	require.NotNil(t, changesetOutput.DataStore, "datastore should not be nil")
	addresses := changesetOutput.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-capabilities-registry-v2"))
	t.Logf("Found %d addresses with qualifier", len(addresses))
	require.Len(t, addresses, 1, "expected exactly one deployed contract with the test qualifier")

	// Verify the address is for the correct chain
	deployedAddress := addresses[0]
	require.Equal(t, chainSelector, deployedAddress.ChainSelector, "deployed contract should be on the correct chain")
	require.NotEmpty(t, deployedAddress.Address, "deployed contract address should not be empty")

	// Verify the contract type is correct
	require.Equal(t, datastore.ContractType("CapabilitiesRegistry"), deployedAddress.Type, "contract type should be CapabilitiesRegistry")
	require.NotNil(t, deployedAddress.Version, "contract version should be set")

	// Verify reports are generated
	require.NotNil(t, changesetOutput.Reports, "reports should be present")
	require.Len(t, changesetOutput.Reports, 1, "should have exactly one operation report")
}
