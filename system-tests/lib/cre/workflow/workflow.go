package workflow

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/Masterminds/semver/v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	workflowreg "github.com/smartcontractkit/cre-tools/pkg/blockchain/contracts/workflow_registry"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	libnet "github.com/smartcontractkit/chainlink/system-tests/lib/net"
)

const (
	// defaultWorkflowQueryLimit is the default limit for querying workflow lists
	defaultWorkflowQueryLimit = 100

	// File URL template for container artifacts
	fileURLTemplate = "file://%s/%s"

	// defaultWorkflowTag is the default tag for workflow versioning
	defaultWorkflowTag    = "latest"
	defaultWorkflowStatus = uint8(0)
)

// ============================================================================
// V2 Registry - Exported Functions
// ============================================================================

// RegisterWithContract registers a workflow in the registry.
// It supports both v1 and v2 workflow registry versions.
func RegisterWithContract(
	ctx context.Context,
	sc *seth.Client,
	workflowRegistryAddr common.Address,
	version *semver.Version,
	donID uint64, workflowName, binaryURL string,
	configURL, secretsURL *string,
	artifactsDirInContainer *string,
) (string, error) {
	// Download and decode workflow binary
	workflowData, err := libnet.DownloadAndDecodeBase64(ctx, binaryURL)
	if err != nil {
		return "", fmt.Errorf("failed to download and decode workflow binary: %w", err)
	}

	// Construct binary URL for container if needed
	binaryURLToUse := constructArtifactURL(binaryURL, artifactsDirInContainer)

	// Handle config URL if provided
	var configData []byte
	configURLToUse := ""
	if configURL != nil && *configURL != "" {
		configData, err = libnet.Download(ctx, *configURL)
		if err != nil {
			return "", fmt.Errorf("failed to download workflow config: %w", err)
		}
		configURLToUse = constructArtifactURL(*configURL, artifactsDirInContainer)
	}

	// Handle secrets URL if provided
	secretsURLToUse := ""
	if secretsURL != nil && *secretsURL != "" {
		secretsURLToUse = constructArtifactURL(*secretsURL, artifactsDirInContainer)
	}

	// Generate workflow ID
	workflowID, err := pkgworkflows.GenerateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workflowData, configData, secretsURLToUse)
	if err != nil {
		return "", fmt.Errorf("failed to generate workflow ID: %w", err)
	}

	// Register workflow based on version
	switch version.Major() {
	case 2:
		// Ensure owner is linked before registering
		if err := LinkOwner(sc, workflowRegistryAddr); err != nil {
			return "", fmt.Errorf("failed to link owner: %w", err)
		}

		// Register workflow
		if err := upsertWorkflow(sc, workflowRegistryAddr, workflowName, workflowID, binaryURLToUse, configURLToUse); err != nil {
			return "", fmt.Errorf("failed to register workflow: %w", err)
		}
	default:
		if err := registerWorkflowV1(sc, workflowRegistryAddr, donID, workflowName, workflowID, binaryURLToUse, configURLToUse, secretsURLToUse); err != nil {
			return "", err
		}
	}

	return workflowID, nil
}

// LinkOwner links a workflow owner address to the registry.
// Returns nil if owner is already linked.
func LinkOwner(sc *seth.Client, workflowRegistryAddr common.Address) error {
	if workflowRegistryAddr == (common.Address{}) {
		return fmt.Errorf("registry address cannot be zero")
	}

	ownerAddr := sc.MustGetRootKeyAddress()
	framework.L.Info().
		Str("owner", ownerAddr.Hex()).
		Str("registry", workflowRegistryAddr.Hex()).
		Msg("Linking owner to registry")

	wrc, err := NewWorkflowRegistryClient(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}
	defer wrc.Close()

	isLinked, err := wrc.IsOwnerLinked(ownerAddr)
	if err != nil {
		return fmt.Errorf("failed to check if owner linked: %w", err)
	}
	if isLinked {
		framework.L.Info().
			Str("owner", ownerAddr.Hex()).
			Msg("Owner already linked")
		return nil
	}

	// Generate and submit linking proof (24 hour validity)
	privateKey := sc.MustGetRootPrivateKey()
	validityDuration := int64(24 * 60 * 60)

	proof, err := wrc.GenerateLinkOwnerProof(context.Background(), privateKey, ownerAddr, validityDuration)
	if err != nil {
		return fmt.Errorf("failed to generate link owner proof: %w", err)
	}

	txOut, err := wrc.LinkOwner(proof.ValidityTimestamp, proof.Proof, proof.Signature)
	if err != nil {
		return fmt.Errorf("failed to link owner: %w", err)
	}

	framework.L.Info().
		Str("owner", ownerAddr.Hex()).
		Str("tx", txOut.Hash.Hex()).
		Msg("Linked owner")
	return nil
}

// upsertWorkflow registers or updates a workflow in the registry.
func upsertWorkflow(
	sc *seth.Client,
	workflowRegistryAddr common.Address,
	workflowName, workflowID, binaryURL, configURL string,
) error {
	if workflowName == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}
	if workflowID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}
	if workflowRegistryAddr == (common.Address{}) {
		return fmt.Errorf("registry address cannot be zero")
	}

	framework.L.Info().
		Str("workflow", workflowName).
		Str("workflowId", workflowID).
		Msg("Registering workflow")

	wrc, err := NewWorkflowRegistryClient(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}
	defer wrc.Close()

	var workflowIDBytes [32]byte
	copy(workflowIDBytes[:], common.Hex2Bytes(workflowID))

	params := workflowreg.RegisterWorkflowParameters{
		WorkflowName: workflowName,
		Tag:          defaultWorkflowTag,
		WorkflowID:   workflowIDBytes,
		Status:       defaultWorkflowStatus,
		DonFamily:    contracts.DonFamily,
		BinaryURL:    binaryURL,
		ConfigURL:    configURL,
		Attributes:   nil,
		KeepAlive:    false,
	}

	txOut, err := wrc.UpsertWorkflow(params)
	if err != nil {
		return fmt.Errorf("failed to upsert workflow: %w", err)
	}

	framework.L.Info().
		Str("workflow", workflowName).
		Str("tx", txOut.Hash.Hex()).
		Msg("Registered workflow")
	return nil
}

// deleteWorkflow deletes a workflow from the registry by name.
func deleteWorkflow(sc *seth.Client, workflowRegistryAddr common.Address, workflowName string) error {
	if workflowName == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}
	if workflowRegistryAddr == (common.Address{}) {
		return fmt.Errorf("registry address cannot be zero")
	}

	framework.L.Info().
		Str("workflow", workflowName).
		Str("registry", workflowRegistryAddr.Hex()).
		Msg("Deleting workflow")

	wrc, err := NewWorkflowRegistryClient(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}
	defer wrc.Close()

	txOut, err := wrc.DeleteWorkflowByNameAndTag(workflowName, defaultWorkflowTag)
	if err != nil {
		return fmt.Errorf("failed to delete workflow %q: %w", workflowName, err)
	}

	framework.L.Info().
		Str("workflow", workflowName).
		Str("tx", txOut.Hash.Hex()).
		Msg("Deleted workflow")
	return nil
}

// deleteAllWorkflows deletes all workflows from the registry.
func deleteAllWorkflows(sc *seth.Client, workflowRegistryAddr common.Address) error {
	if workflowRegistryAddr == (common.Address{}) {
		return fmt.Errorf("registry address cannot be zero")
	}

	framework.L.Info().
		Str("registry", workflowRegistryAddr.Hex()).
		Msg("Deleting all workflows")

	wrc, err := NewWorkflowRegistryClient(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}
	defer wrc.Close()

	deleted, err := wrc.DeleteAllWorkflows(contracts.DonFamily)
	if err != nil {
		return fmt.Errorf("failed to delete all workflows: %w", err)
	}

	framework.L.Info().
		Int("count", deleted).
		Msg("Deleted workflows")
	return nil
}

// DeleteWithContract removes a workflow from the workflow registry contract.
// It supports both v1 and v2 workflow registry versions.
func DeleteWithContract(
	ctx context.Context,
	sc *seth.Client,
	workflowRegistryAddr common.Address,
	version *semver.Version,
	workflowName string,
) error {
	switch version.Major() {
	case 2:
		return deleteWorkflow(sc, workflowRegistryAddr, workflowName)
	default:
		return deleteWorkflowV1(sc, workflowRegistryAddr, workflowName)
	}
}

// DeleteAllWithContract removes all workflows owned by the caller from the workflow registry contract.
// It supports both v1 and v2 workflow registry versions.
func DeleteAllWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) error {
	switch version.Major() {
	case 2:
		return deleteAllWorkflows(sc, workflowRegistryAddr)
	default:
		return deleteAllWorkflowsV1(sc, workflowRegistryAddr)
	}
}

// ============================================================================
// V2 Registry - Unexported Functions
// ============================================================================

// constructArtifactURL constructs the appropriate URL based on whether artifacts are in a container.
func constructArtifactURL(originalURL string, artifactsDirInContainer *string) string {
	if artifactsDirInContainer != nil {
		return fmt.Sprintf(fileURLTemplate, *artifactsDirInContainer, filepath.Base(originalURL))
	}
	return originalURL
}

// ============================================================================
// V1 Registry - Unexported Functions
// ============================================================================

// registerWorkflowV1 handles workflow registration for v1 registry contracts
func registerWorkflowV1(sc *seth.Client, workflowRegistryAddr common.Address, donID uint64,
	workflowName, workflowID, binaryURL, configURL, secretsURL string) error {
	registry, err := createRegistryV1Instance(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}

	// Register workflow
	_, err = sc.Decode(registry.RegisterWorkflow(
		sc.NewTXOpts(),
		workflowName,
		[32]byte(common.Hex2Bytes(workflowID)),
		libc.MustSafeUint32FromUint64(donID),
		defaultWorkflowStatus,
		binaryURL,
		configURL,
		secretsURL,
	))
	if err != nil {
		return fmt.Errorf("failed to register workflow: %w", err)
	}

	return nil
}

// deleteAllWorkflowsV1 removes all workflows for v1 registry contracts.
func deleteAllWorkflowsV1(sc *seth.Client, workflowRegistryAddr common.Address) error {
	// Create registry instance once for all operations
	registry, err := createRegistryV1Instance(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}

	// Get list of workflows to delete
	workflows, err := registry.GetWorkflowMetadataListByOwner(
		sc.NewCallOpts(),
		sc.MustGetRootKeyAddress(),
		big.NewInt(0),
		big.NewInt(defaultWorkflowQueryLimit),
	)
	if err != nil {
		return fmt.Errorf("failed to get workflow metadata list: %w", err)
	}

	// Delete each workflow using the same registry instance
	for _, workflow := range workflows {
		workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflow.WorkflowName)
		if _, err := sc.Decode(registry.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey)); err != nil {
			return fmt.Errorf("failed to delete workflow %q: %w", workflow.WorkflowName, err)
		}
	}

	return nil
}

// computeHashKey generates a Keccak256 hash from owner address and workflow name.
// This is used for v1 workflow registry contract operations.
func computeHashKey(owner common.Address, workflowName string) [32]byte {
	ownerBytes := owner.Bytes()
	nameBytes := []byte(workflowName)
	data := make([]byte, len(ownerBytes)+len(nameBytes))
	copy(data, ownerBytes)
	copy(data[len(ownerBytes):], nameBytes)

	return crypto.Keccak256Hash(data)
}

// deleteWorkflowV1 handles workflow deletion for v1 registry contracts.
func deleteWorkflowV1(sc *seth.Client, workflowRegistryAddr common.Address, workflowName string) error {
	registry, err := createRegistryV1Instance(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}

	workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflowName)
	if _, err := sc.Decode(registry.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey)); err != nil {
		return fmt.Errorf("failed to delete workflow %q: %w", workflowName, err)
	}

	return nil
}

// createRegistryV1Instance creates a new v1 workflow registry instance.
func createRegistryV1Instance(sc *seth.Client, workflowRegistryAddr common.Address) (*workflow_registry_wrapper.WorkflowRegistry, error) {
	registry, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow registry instance: %w", err)
	}

	// add contract ABI to Seth, so that it can decode transaction errors
	abi, aErr := workflow_registry_wrapper.WorkflowRegistryMetaData.GetAbi()
	if aErr != nil {
		return nil, fmt.Errorf("failed to get WorkflowRegistryV1 ABI: %w", aErr)
	}

	sc.ABIFinder.ContractStore.AddABI("WorkflowRegistryV1", *abi)

	return registry, nil
}
