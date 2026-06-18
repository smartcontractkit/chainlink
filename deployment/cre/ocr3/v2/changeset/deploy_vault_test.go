package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestDeployVault(t *testing.T) {
	h := test.NewTestHarness(t)

	task := runtime.ChangesetTask(DeployVault{}, DeployVaultInput{
		ChainSelector: h.RegistrySelector,
		Qualifier:     "vault",
	})
	err := h.Runtime.Exec(task)
	require.NoError(t, err)
	output := h.Runtime.State().Outputs[task.ID()]

	// Verify the datastore contains the deployed contract
	addresses, err := output.DataStore.Addresses().Fetch()
	require.NoError(t, err, "should fetch addresses without error")
	require.Len(t, addresses, 2, "expected two deployed contracts (Plugin and DKG")

	// Verify the address is for the correct chain
	for _, addr := range addresses {
		require.Equal(t, h.RegistrySelector, addr.ChainSelector, "deployed contract should be on the correct chain")
		require.NotEmpty(t, addr.Address, "deployed contract address should not be empty")

		// Verify the contract type is correct
		require.Equal(t, datastore.ContractType("OCR3Capability"), addr.Type, "contract type should be OCR3Capability")
		require.NotNil(t, addr.Version, "contract version should be set")
	}
}
