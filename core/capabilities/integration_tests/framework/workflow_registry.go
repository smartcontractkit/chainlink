package framework

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
)

type WorkflowRegistry struct {
	t        *testing.T
	backend  *EthBlockchain
	contract *workflow_registry_wrapper.WorkflowRegistry
	addr     common.Address
}

func NewWorkflowRegistry(ctx context.Context, t *testing.T, backend *EthBlockchain) *WorkflowRegistry {
	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backend.transactionOpts, backend.Backend.Client())
	backend.Backend.Commit()
	require.NoError(t, err)

	// setup contract state to allow the secrets to be updated
	updateAuthorizedAddress(t, backend, wfRegistryC, []common.Address{backend.transactionOpts.From}, true)

	return &WorkflowRegistry{t: t, addr: wfRegistryAddr, contract: wfRegistryC, backend: backend}
}

func (r *WorkflowRegistry) UpdateAllowedDons(donIDs []uint32) {
	updateAllowedDONs(r.t, r.backend, r.contract, donIDs, true)
}

func (r *WorkflowRegistry) RegisterWorkflow(input Workflow, donID uint32) {
	registerWorkflow(r.t, r.backend, r.contract, input, donID)
}

func (r *WorkflowRegistry) UpdateWorkflow(input UpdatedWorkflow, donID uint32) {
	updateWorkflow(r.t, r.backend, r.contract, input)
}

func (r *WorkflowRegistry) ComputeHashKey(owner string, field string) [32]byte {
	return computeHashKey(r.t, r.contract, owner, field)
}

func updateAuthorizedAddress(
	t *testing.T,
	th *EthBlockchain,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	addresses []common.Address,
	_ bool,
) {
	t.Helper()
	_, err := wfRegC.UpdateAuthorizedAddresses(th.transactionOpts, addresses, true)
	require.NoError(t, err, "failed to update authorised addresses")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
	gotAddresses, err := wfRegC.GetAllAuthorizedAddresses(&bind.CallOpts{
		From: th.transactionOpts.From,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, addresses, gotAddresses)
}

func updateAllowedDONs(
	t *testing.T,
	th *EthBlockchain,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	donIDs []uint32,
	allowed bool,
) {
	t.Helper()
	_, err := wfRegC.UpdateAllowedDONs(th.transactionOpts, donIDs, allowed)
	require.NoError(t, err, "failed to update DONs")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
	gotDons, err := wfRegC.GetAllAllowedDONs(&bind.CallOpts{
		From: th.transactionOpts.From,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, donIDs, gotDons)
}

type Workflow struct {
	Name       string
	ID         [32]byte
	Status     uint8
	BinaryURL  string
	ConfigURL  string
	SecretsURL string
}

func registerWorkflow(
	t *testing.T,
	th *EthBlockchain,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	input Workflow,
	donID uint32,
) {
	t.Helper()
	_, err := wfRegC.RegisterWorkflow(th.transactionOpts, input.Name, input.ID, donID,
		input.Status, input.BinaryURL, input.ConfigURL, input.SecretsURL)
	require.NoError(t, err, "failed to register workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

type UpdatedWorkflow struct {
	WorkflowKey [32]byte
	ID          [32]byte
	BinaryURL   string
	ConfigURL   string
	SecretsURL  string
}

func updateWorkflow(
	t *testing.T,
	th *EthBlockchain,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	input UpdatedWorkflow,
) {
	t.Helper()
	_, err := wfRegC.UpdateWorkflow(th.transactionOpts, input.WorkflowKey, input.ID, input.BinaryURL, input.ConfigURL, input.SecretsURL)
	require.NoError(t, err, "failed to update workflow")
	th.Backend.Commit()
	th.Backend.Commit()
	th.Backend.Commit()
}

func computeHashKey(
	t *testing.T,
	wfRegC *workflow_registry_wrapper.WorkflowRegistry,
	owner string,
	field string,
) [32]byte {
	t.Helper()
	hashKey, err := wfRegC.ComputeHashKey(&bind.CallOpts{}, common.HexToAddress(owner), field)
	require.NoError(t, err, "failed to compute hash key")
	return hashKey
}
