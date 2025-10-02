package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	workflow_registry_wrapper_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	libnet "github.com/smartcontractkit/chainlink/system-tests/lib/net"
)

const (
	// defaultWorkflowQueryLimit is the default limit for querying workflow lists
	defaultWorkflowQueryLimit = 100

	// File URL template for container artifacts
	fileURLTemplate = "file://%s/%s"

	// Default values for workflow registration
	defaultWorkflowTag    = "some-tag"
	defaultWorkflowStatus = uint8(0)

	// Common error message templates
	errCreateContractInstance  = "failed to create %s %s instance"
	errCreateRegistryInstance  = "failed to create workflow registry instance"
	errGetWorkflowMetadataList = "failed to get workflow metadata list"
	errDeleteWorkflow          = "failed to delete workflow %q"
)

func RegisterWithContract(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, typeVersion deployment.TypeAndVersion,
	donID uint64, workflowName, binaryURL string,
	configURL, secretsURL *string,
	artifactsDirInContainer *string,
) (string, error) {
	// Download and decode workflow binary
	workflowData, err := libnet.DownloadAndDecodeBase64(ctx, binaryURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to download and decode workflow binary")
	}

	// Construct binary URL for container if needed
	binaryURLToUse := constructArtifactURL(binaryURL, artifactsDirInContainer)

	// Handle config URL if provided
	var configData []byte
	configURLToUse := ""
	if configURL != nil && *configURL != "" {
		configData, err = libnet.Download(ctx, *configURL)
		if err != nil {
			return "", errors.Wrap(err, "failed to download workflow config")
		}
		configURLToUse = constructArtifactURL(*configURL, artifactsDirInContainer)
	}

	// Handle secrets URL if provided
	secretsURLToUse := ""
	if secretsURL != nil && *secretsURL != "" {
		secretsURLToUse = constructArtifactURL(*secretsURL, artifactsDirInContainer)
	}

	// Generate workflow ID
	workflowID, err := generateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workflowData, configData, secretsURLToUse)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate workflow ID")
	}

	// Register workflow based on version
	switch typeVersion.Version.Major() {
	case 2:
		if err := registerWorkflowV2(sc, workflowRegistryAddr, typeVersion, workflowName, workflowID, binaryURLToUse, configURLToUse); err != nil {
			return "", err
		}
	default:
		if err := registerWorkflowV1(sc, workflowRegistryAddr, donID, workflowName, workflowID, binaryURLToUse, configURLToUse, secretsURLToUse); err != nil {
			return "", err
		}
	}

	return workflowID, nil
}

func LinkOwner(sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) error {
	switch tv.Version.Major() {
	case 2:
		validity := time.Now().UTC().Add(time.Hour * 24)
		validityTimestamp := big.NewInt(validity.Unix())
		defaultOrgID := 22
		nonce := uuid.New().String()
		workflowOwner := sc.MustGetRootKeyAddress().Hex()
		data := fmt.Sprintf("%s%d%s", workflowOwner, defaultOrgID, nonce)
		hash := sha256.Sum256([]byte(data))
		ownershipProof := hex.EncodeToString(hash[:])
		linkRequestType := uint8(0)

		registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, tv)
		if err != nil {
			return err
		}

		version, versionErr := registry.TypeAndVersion(sc.NewCallOpts())
		if versionErr != nil {
			return versionErr
		}

		messageDigest, err := PreparePayloadForSigning(
			OwnershipProofSignaturePayload{
				RequestType:              linkRequestType,
				WorkflowOwnerAddress:     common.HexToAddress(workflowOwner),
				ChainID:                  strconv.FormatInt(sc.ChainID, 10),
				WorkflowRegistryContract: workflowRegistryAddr,
				Version:                  version,
				ValidityTimestamp:        validity,
				OwnershipProofHash:       common.HexToHash(ownershipProof),
			})
		if err != nil {
			return fmt.Errorf("failed to prepare payload for signing: %w", err)
		}

		signature, err := crypto.Sign(messageDigest, sc.MustGetRootPrivateKey())
		if err != nil {
			return fmt.Errorf("failed to sign ownership proof: %w", err)
		}

		signature[64] += 27

		_, err = sc.Decode(registry.LinkOwner(sc.NewTXOpts(), validityTimestamp, common.HexToHash(ownershipProof), signature))
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("invalid version for linking owner")
	}
}

// GetWorkflowNames retrieves all workflow names for the given registry contract.
// It supports both v1 and v2 workflow registry versions.
func GetWorkflowNames(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) ([]string, error) {
	switch tv.Version.Major() {
	case 2:
		return getWorkflowNamesV2(sc, workflowRegistryAddr, tv)
	default:
		return getWorkflowNamesV1(sc, workflowRegistryAddr)
	}
}

// DeleteWithContract removes a workflow from the workflow registry contract.
// It supports both v1 and v2 workflow registry versions.
func DeleteWithContract(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, tv deployment.TypeAndVersion,
	workflowName string,
) error {
	switch tv.Version.Major() {
	case 2:
		return deleteWorkflowV2(ctx, sc, workflowRegistryAddr, tv, workflowName)
	default:
		return deleteWorkflowV1(ctx, sc, workflowRegistryAddr, workflowName)
	}
}

// DeleteAllWithContract removes all workflows owned by the caller from the workflow registry contract.
// It supports both v1 and v2 workflow registry versions.
func DeleteAllWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) error {
	switch tv.Version.Major() {
	case 2:
		return deleteAllWorkflowsV2(ctx, sc, workflowRegistryAddr, tv)
	default:
		return deleteAllWorkflowsV1(ctx, sc, workflowRegistryAddr)
	}
}

// RemoveWorkflowArtifactsFromLocalEnv removes workflow artifact files from the local filesystem.
// Empty file paths are silently skipped.
func RemoveWorkflowArtifactsFromLocalEnv(artifactPaths ...string) error {
	for _, path := range artifactPaths {
		if path == "" {
			continue
		}

		if err := os.Remove(path); err != nil {
			return errors.Wrapf(err, "failed to remove workflow artifact at %q", path)
		}
	}
	return nil
}

// constructArtifactURL constructs the appropriate URL based on whether artifacts are in a container
func constructArtifactURL(originalURL string, artifactsDirInContainer *string) string {
	if artifactsDirInContainer != nil {
		return fmt.Sprintf(fileURLTemplate, *artifactsDirInContainer, filepath.Base(originalURL))
	}
	return originalURL
}

// registerWorkflowV2 handles workflow registration for v2 registry contracts
func registerWorkflowV2(sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion,
	workflowName, workflowID, binaryURL, configURL string) error {
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, tv)
	if err != nil {
		return err
	}

	// Check and link owner if needed using existing helper function
	if verifyErr := verifyOwnerLinkedWithRegistry(registry, sc, workflowName); verifyErr != nil {
		// If owner is not linked, try to link them
		if linkErr := LinkOwner(sc, workflowRegistryAddr, tv); linkErr != nil {
			return errors.Wrap(linkErr, "failed to link owner to org")
		}
	}

	// Register workflow
	_, err = sc.Decode(registry.UpsertWorkflow(
		sc.NewTXOpts(),
		workflowName,
		defaultWorkflowTag,
		[32]byte(common.Hex2Bytes(workflowID)),
		defaultWorkflowStatus,
		contracts.DonFamily,
		binaryURL,
		configURL,
		nil,
		false,
	))
	if err != nil {
		return errors.Wrap(err, "failed to register workflow")
	}

	return nil
}

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
		return errors.Wrap(err, "failed to register workflow")
	}

	return nil
}

// deleteAllWorkflowsV2 removes all workflows for v2 registry contracts.
func deleteAllWorkflowsV2(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) error {
	// Create registry instance once for all operations
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, tv)
	if err != nil {
		return err
	}

	// Verify owner linking before attempting any deletions
	if verifyErr := verifyOwnerLinkedWithRegistry(registry, sc, ""); verifyErr != nil {
		return verifyErr
	}

	// Get list of workflows to delete using the same registry instance
	workflows, err := getWorkflowListWithRegistryV2(registry, sc)
	if err != nil {
		return err
	}

	// Delete each workflow using the same registry instance
	for _, workflow := range workflows {
		if _, err := sc.Decode(registry.DeleteWorkflow(sc.NewTXOpts(), workflow.WorkflowId)); err != nil {
			return errors.Wrapf(err, errDeleteWorkflow, workflow.WorkflowName)
		}
	}

	return nil
}

// deleteAllWorkflowsV1 removes all workflows for v1 registry contracts.
func deleteAllWorkflowsV1(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address) error {
	// Create registry instance once for all operations
	registry, err := createRegistryV1Instance(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}

	// Get list of workflows to delete using the same registry instance
	workflows, err := getWorkflowListWithRegistryV1(registry, sc)
	if err != nil {
		return err
	}

	// Delete each workflow using the same registry instance
	for _, workflow := range workflows {
		workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflow.WorkflowName)
		if _, err := sc.Decode(registry.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey)); err != nil {
			return errors.Wrapf(err, errDeleteWorkflow, workflow.WorkflowName)
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

// deleteWorkflowV2 handles workflow deletion for v2 registry contracts.
func deleteWorkflowV2(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion, workflowName string,
) error {
	// Create registry instance once for all operations
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, tv)
	if err != nil {
		return err
	}

	// Find workflow using the same registry instance
	workflowID, err := findWorkflowByNameWithRegistry(registry, sc, workflowName)
	if err != nil {
		return errors.Wrapf(err, "failed to find workflow %q", workflowName)
	}

	// Verify owner linking using the same registry instance
	if err := verifyOwnerLinkedWithRegistry(registry, sc, workflowName); err != nil {
		return err
	}

	// Delete workflow using the same registry instance
	if _, err := sc.Decode(registry.DeleteWorkflow(sc.NewTXOpts(), workflowID)); err != nil {
		return errors.Wrapf(err, "failed to delete workflow %q (ID: %x)", workflowName, workflowID)
	}

	return nil
}

// deleteWorkflowV1 handles workflow deletion for v1 registry contracts.
func deleteWorkflowV1(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, workflowName string,
) error {
	registry, err := createRegistryV1Instance(sc, workflowRegistryAddr)
	if err != nil {
		return err
	}

	workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflowName)
	if _, err := sc.Decode(registry.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey)); err != nil {
		return errors.Wrapf(err, "failed to delete workflow %q", workflowName)
	}

	return nil
}

// findWorkflowByNameWithRegistry finds a workflow by name using an existing registry instance and returns its ID.
func findWorkflowByNameWithRegistry(registry *workflow_registry_wrapper_v2.WorkflowRegistry, sc *seth.Client, workflowName string) ([32]byte, error) {
	workflows, err := getWorkflowListWithRegistryV2(registry, sc)
	if err != nil {
		return [32]byte{}, err
	}

	for _, workflow := range workflows {
		if workflow.WorkflowName == workflowName {
			return workflow.WorkflowId, nil
		}
	}

	return [32]byte{}, errors.Errorf("workflow %q not found in registry", workflowName)
}

// verifyOwnerLinkedWithRegistry checks if the owner is properly linked using an existing registry instance.
func verifyOwnerLinkedWithRegistry(registry *workflow_registry_wrapper_v2.WorkflowRegistry, sc *seth.Client, workflowName string) error {
	ownerAddr := sc.MustGetRootKeyAddress()
	isLinked, err := registry.IsOwnerLinked(sc.NewCallOpts(), ownerAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to check owner link status for workflow %q", workflowName)
	}

	if !isLinked {
		return errors.Errorf("owner %s is not linked to an organization, cannot delete workflow %q",
			ownerAddr.Hex(), workflowName)
	}

	return nil
}

// getRegistryV2Instance creates a new v2 workflow registry instance.
func getRegistryV2Instance(sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) (*workflow_registry_wrapper_v2.WorkflowRegistry, error) {
	registry, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return nil, errors.Wrapf(err, errCreateContractInstance, tv.Type, tv.Version)
	}
	return registry, nil
}

// createRegistryV1Instance creates a new v1 workflow registry instance.
func createRegistryV1Instance(sc *seth.Client, workflowRegistryAddr common.Address) (*workflow_registry_wrapper.WorkflowRegistry, error) {
	registry, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return nil, errors.Wrap(err, errCreateRegistryInstance)
	}
	return registry, nil
}

// getWorkflowListWithRegistryV2 retrieves the full workflow list using an existing v2 registry instance.
func getWorkflowListWithRegistryV2(registry *workflow_registry_wrapper_v2.WorkflowRegistry, sc *seth.Client) ([]workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView, error) {
	workflows, err := registry.GetWorkflowListByOwner(
		sc.NewCallOpts(),
		sc.MustGetRootKeyAddress(),
		big.NewInt(0),
		big.NewInt(defaultWorkflowQueryLimit),
	)
	if err != nil {
		return nil, errors.Wrap(err, errGetWorkflowMetadataList)
	}

	return workflows, nil
}

// getWorkflowListV2 retrieves the full workflow list for v2 registry contracts.
func getWorkflowListV2(sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) ([]workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView, error) {
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, tv)
	if err != nil {
		return nil, err
	}

	return getWorkflowListWithRegistryV2(registry, sc)
}

// getWorkflowNamesV2 retrieves all workflow names for v2 registry contracts.
func getWorkflowNamesV2(sc *seth.Client, workflowRegistryAddr common.Address, tv deployment.TypeAndVersion) ([]string, error) {
	workflows, err := getWorkflowListV2(sc, workflowRegistryAddr, tv)
	if err != nil {
		return nil, err
	}

	workflowNames := make([]string, 0, len(workflows))
	for _, workflow := range workflows {
		workflowNames = append(workflowNames, workflow.WorkflowName)
	}

	return workflowNames, nil
}

// getWorkflowListWithRegistryV1 retrieves the full workflow list using an existing v1 registry instance.
func getWorkflowListWithRegistryV1(registry *workflow_registry_wrapper.WorkflowRegistry, sc *seth.Client) ([]workflow_registry_wrapper.WorkflowRegistryWorkflowMetadata, error) {
	workflows, err := registry.GetWorkflowMetadataListByOwner(
		sc.NewCallOpts(),
		sc.MustGetRootKeyAddress(),
		big.NewInt(0),
		big.NewInt(defaultWorkflowQueryLimit),
	)
	if err != nil {
		return nil, errors.Wrap(err, errGetWorkflowMetadataList)
	}

	return workflows, nil
}

// getWorkflowListV1 retrieves the full workflow list for v1 registry contracts.
func getWorkflowListV1(sc *seth.Client, workflowRegistryAddr common.Address) ([]workflow_registry_wrapper.WorkflowRegistryWorkflowMetadata, error) {
	registry, err := createRegistryV1Instance(sc, workflowRegistryAddr)
	if err != nil {
		return nil, err
	}

	return getWorkflowListWithRegistryV1(registry, sc)
}

// getWorkflowNamesV1 retrieves all workflow names for v1 registry contracts.
func getWorkflowNamesV1(sc *seth.Client, workflowRegistryAddr common.Address) ([]string, error) {
	workflows, err := getWorkflowListV1(sc, workflowRegistryAddr)
	if err != nil {
		return nil, err
	}

	workflowNames := make([]string, 0, len(workflows))
	for _, workflow := range workflows {
		workflowNames = append(workflowNames, workflow.WorkflowName)
	}

	return workflowNames, nil
}

// generateWorkflowIDFromStrings creates a workflow ID from string inputs.
// The owner address can be provided with or without the "0x" prefix.
func generateWorkflowIDFromStrings(owner, name string, workflow, config []byte, secretsURL string) (string, error) {
	// Remove "0x" prefix if present
	ownerHex := strings.TrimPrefix(owner, "0x")

	ownerBytes, err := hex.DecodeString(ownerHex)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode owner address")
	}

	workflowID, err := pkgworkflows.GenerateWorkflowID(ownerBytes, name, workflow, config, secretsURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate workflow ID")
	}

	return hex.EncodeToString(workflowID[:]), nil
}
