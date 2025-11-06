package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
)

func TestDeployWorkflowRegistry(t *testing.T) {
	// Create a minimal environment with one EVM chain
	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	task := runtime.ChangesetTask(DeployWorkflowRegistry{}, DeployWorkflowRegistryInput{
		ChainSelector: selector,
		Qualifier:     "test-workflow-registry-v2",
	})

	err = rt.Exec(task)
	require.NoError(t, err, "failed to deploy WorkflowRegistry")

	output := rt.State().Outputs[task.ID()]

	// Verify the datastore contains the deployed contract
	require.NotNil(t, rt.State().DataStore, "datastore should not be nil")
	addresses := rt.State().DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-workflow-registry-v2"))
	t.Logf("Found %d addresses with qualifier", len(addresses))
	require.Len(t, addresses, 1, "expected exactly one deployed contract with the test qualifier")

	// Verify the address is for the correct chain
	deployedAddress := addresses[0]
	require.Equal(t, selector, deployedAddress.ChainSelector, "deployed contract should be on the correct chain")
	require.NotEmpty(t, deployedAddress.Address, "deployed contract address should not be empty")

	// Verify the contract type is correct
	require.Equal(t, datastore.ContractType("WorkflowRegistry"), deployedAddress.Type, "contract type should be WorkflowRegistry")
	require.NotNil(t, deployedAddress.Version, "contract version should be set")

	// Verify reports are generated
	require.NotNil(t, output.Reports, "reports should be present")
	require.Len(t, output.Reports, 1, "should have exactly one operation report")
}
