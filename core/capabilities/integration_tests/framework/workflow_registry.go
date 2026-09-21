package framework

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

const (
	// The v2 registry gates workflows on a per DON family limit rather than on the v1 allowlist of
	// DON IDs, so a family is only usable once these limits are set.
	workflowRegistryDONLimit         = uint32(10000)
	workflowRegistryUserDefaultLimit = uint32(1000)

	// Request type discriminator the registry expects in the payload signed for LinkOwner.
	workflowRegistryLinkRequestType = uint8(0)

	workflowRegistryOwnerLinkValidity = time.Hour
)

type WorkflowRegistry struct {
	t        *testing.T
	backend  *EthBlockchain
	contract *workflow_registry_wrapper.WorkflowRegistry
	addr     common.Address
}

func NewWorkflowRegistry(ctx context.Context, t *testing.T, backend *EthBlockchain) *WorkflowRegistry {
	// Deploy a test workflow_registry
	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper.DeployWorkflowRegistry(backend.transactionOpts, backend.Client())
	backend.Commit()
	require.NoError(t, err)

	r := &WorkflowRegistry{t: t, addr: wfRegistryAddr, contract: wfRegistryC, backend: backend}

	// The v2 registry only accepts workflows from a linked owner whose signer is allowlisted, so the
	// deployer is set up as both before any workflow can be registered.
	r.updateAllowedSigners([]common.Address{backend.transactionOpts.From})
	r.linkOwner(ctx, backend.transactionOpts.From)

	return r
}

// UpdateAllowedDons sets the workflow limits for the given DON families. Families are the v2
// registry's equivalent of the v1 allowlist of DON IDs, and match the families the DON is
// registered under in the capabilities registry.
func (r *WorkflowRegistry) UpdateAllowedDons(donFamilies []string) {
	for _, donFamily := range donFamilies {
		_, err := r.contract.SetDONLimit(r.backend.transactionOpts, donFamily, workflowRegistryDONLimit,
			workflowRegistryUserDefaultLimit)
		require.NoError(r.t, err, "failed to set limits for DON family %s", donFamily)
		r.commit()

		limits, err := r.contract.GetMaxWorkflowsPerDON(&bind.CallOpts{From: r.backend.transactionOpts.From}, donFamily)
		require.NoError(r.t, err)
		require.Equal(r.t, workflowRegistryDONLimit, limits.MaxWorkflows)
	}
}

func (r *WorkflowRegistry) RegisterWorkflow(input Workflow, donFamily string) {
	r.upsertWorkflow(input, donFamily)
}

// UpdateWorkflow updates a previously registered workflow. The v2 registry has a single upsert
// entrypoint, and rejects updates that change the status or the DON family of a workflow.
func (r *WorkflowRegistry) UpdateWorkflow(input Workflow, donFamily string) {
	r.upsertWorkflow(input, donFamily)
}

type Workflow struct {
	Name       string
	Tag        string
	ID         [32]byte
	Status     uint8
	BinaryURL  string
	ConfigURL  string
	Attributes []byte
	KeepAlive  bool
}

func (r *WorkflowRegistry) upsertWorkflow(input Workflow, donFamily string) {
	r.t.Helper()
	_, err := r.contract.UpsertWorkflow(r.backend.transactionOpts, input.Name, input.Tag, input.ID, input.Status,
		donFamily, input.BinaryURL, input.ConfigURL, input.Attributes, input.KeepAlive)
	require.NoError(r.t, err, "failed to upsert workflow")
	r.commit()
}

func (r *WorkflowRegistry) updateAllowedSigners(addresses []common.Address) {
	r.t.Helper()
	_, err := r.contract.UpdateAllowedSigners(r.backend.transactionOpts, addresses, true)
	require.NoError(r.t, err, "failed to update allowed signers")
	r.commit()

	for _, address := range addresses {
		allowed, err := r.contract.IsAllowedSigner(&bind.CallOpts{From: r.backend.transactionOpts.From}, address)
		require.NoError(r.t, err)
		require.True(r.t, allowed, "signer %s was not allowlisted", address)
	}
}

// linkOwner proves ownership of the given address to the registry, which the v2 contract requires
// before that address can register any workflow.
func (r *WorkflowRegistry) linkOwner(ctx context.Context, owner common.Address) {
	r.t.Helper()

	typeAndVersion, err := r.contract.TypeAndVersion(&bind.CallOpts{From: r.backend.transactionOpts.From})
	require.NoError(r.t, err)

	chainID, err := r.backend.Client().ChainID(ctx)
	require.NoError(r.t, err)

	validityTimestamp := big.NewInt(time.Now().UTC().Add(workflowRegistryOwnerLinkValidity).Unix())
	proof := ownershipProofHash(owner.String(), "integration-tests", "1")

	arguments, err := ownershipLinkABIArguments()
	require.NoError(r.t, err)

	packed, err := arguments.Pack(workflowRegistryLinkRequestType, owner, chainID, r.addr, typeAndVersion,
		validityTimestamp, proof)
	require.NoError(r.t, err)

	hash := crypto.Keccak256(packed)

	// The contract recovers the signer using EIP-191, so the digest has to carry the same prefix.
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	signature, err := r.backend.SignHash(crypto.Keccak256([]byte(prefixedMessage)))
	require.NoError(r.t, err)
	signature[64] += 27

	_, err = r.contract.LinkOwner(r.backend.transactionOpts, validityTimestamp, proof, signature)
	require.NoError(r.t, err, "failed to link owner")
	r.commit()

	linked, err := r.contract.IsOwnerLinked(&bind.CallOpts{From: r.backend.transactionOpts.From}, owner)
	require.NoError(r.t, err)
	require.True(r.t, linked, "owner %s was not linked", owner)
}

func ownershipProofHash(workflowOwnerAddress, organizationID, nonce string) [32]byte {
	return sha256.Sum256([]byte(workflowOwnerAddress + organizationID + nonce))
}

// ownershipLinkABIArguments describes, in order, the payload the registry hashes when verifying an
// owner link signature.
func ownershipLinkABIArguments() (abi.Arguments, error) {
	uint8Type, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create uint8 type: %w", err)
	}

	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create address type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create bytes32 type: %w", err)
	}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create uint256 type: %w", err)
	}

	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create string type: %w", err)
	}

	return abi.Arguments{
		{Type: uint8Type},   // request type
		{Type: addressType}, // owner address
		{Type: uint256Type}, // chain ID
		{Type: addressType}, // address of the contract
		{Type: stringType},  // type and version string
		{Type: uint256Type}, // validity timestamp
		{Type: bytes32Type}, // ownership proof hash
	}, nil
}

// commit mines enough blocks for the writes above to be visible to finalized reads.
func (r *WorkflowRegistry) commit() {
	r.backend.Commit()
	r.backend.Commit()
	r.backend.Commit()
}
