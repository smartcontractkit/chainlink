package workflowregistry

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type SetupTestWorkflowRegistryResponse struct {
	Registry         *workflow_registry.WorkflowRegistry
	Chain            cldf.Chain
	RegistrySelector uint64
	AddressBook      cldf.AddressBook
}

func SetupTestWorkflowRegistry(t *testing.T, lggr logger.Logger, chainSel uint64) *SetupTestWorkflowRegistryResponse {
	chain := testChain(t)

	deployer, err := newWorkflowRegistryDeployer()
	require.NoError(t, err)
	resp, err := deployer.Deploy(changeset.DeployRequest{Chain: chain})
	require.NoError(t, err)

	addressBook := cldf.NewMemoryAddressBookFromMap(
		map[uint64]map[string]cldf.TypeAndVersion{
			chainSel: map[string]cldf.TypeAndVersion{
				resp.Address.Hex(): resp.Tv,
			},
		},
	)

	return &SetupTestWorkflowRegistryResponse{
		Registry:         deployer.Contract(),
		Chain:            chain,
		RegistrySelector: chain.Selector,
		AddressBook:      addressBook,
	}
}

func testChain(t *testing.T) cldf.Chain {
	chains, _ := memory.NewMemoryChains(t, 1, 5)
	var chain cldf.Chain
	for _, c := range chains {
		chain = c
		break
	}
	require.NotEmpty(t, chain)
	return chain
}
