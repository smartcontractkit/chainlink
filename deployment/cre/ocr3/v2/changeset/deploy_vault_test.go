package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestDeployVault(t *testing.T) {
	env := test.SetupEnvV2(t, false)

	changesetOutput, err := DeployVault{}.Apply(*env.Env, DeployVaultInput{
		ChainSelector: env.RegistrySelector,
		Qualifier:     "vault",
	})
	if err != nil {
		t.Fatalf("changeset apply failed: %v", err)
	}

	// Verify the datastore contains the deployed contract
	addresses, err := changesetOutput.DataStore.Addresses().Fetch()
	require.NoError(t, err, "should fetch addresses without error")
	require.Len(t, addresses, 2, "expected two deployed contracts (Plugin and DKG")

	// Verify the address is for the correct chain
	for _, addr := range addresses {
		require.Equal(t, env.RegistrySelector, addr.ChainSelector, "deployed contract should be on the correct chain")
		require.NotEmpty(t, addr.Address, "deployed contract address should not be empty")

		// Verify the contract type is correct
		require.Equal(t, datastore.ContractType("OCR3Capability"), addr.Type, "contract type should be OCR3Capability")
		require.NotNil(t, addr.Version, "contract version should be set")
	}
}
