package workflow

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	workflow_registry_wrapper_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	chainlinkvalues "github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

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

	// vaultConfigStaticPropagationWait is used when node DB info is not available
	// (e.g. Kubernetes provider or missing state file) and we cannot actively poll
	// for registry syncer propagation.
	vaultConfigStaticPropagationWait = 15 * time.Second

	// latestRegistrySyncStateQuery retrieves the most recent capabilities registry snapshot stored by the syncer.
	latestRegistrySyncStateQuery = `SELECT data FROM registry_syncer_states ORDER BY id DESC LIMIT 1`
)

// vaultSyncStateDONCapConfig mirrors the JSON shape of registrysyncer.CapabilityConfiguration.
type vaultSyncStateDONCapConfig struct {
	Config []byte `json:"Config"`
}

// vaultSyncStateDON mirrors the JSON shape of registrysyncer.DON (CapabilityConfigurations only).
type vaultSyncStateDON struct {
	CapabilityConfigurations map[string]vaultSyncStateDONCapConfig `json:"CapabilityConfigurations"`
}

// vaultSyncStatePayload is a partial deserialisation of the registry_syncer_states.data JSON blob.
type vaultSyncStatePayload struct {
	IDsToDONs map[string]vaultSyncStateDON `json:"IDsToDONs"`
}

func RegisterWithContract(
	ctx context.Context,
	sc *seth.Client,
	workflowRegistryAddr common.Address,
	version *semver.Version,
	donID uint64, workflowName, binaryURL string,
	configURL, secretsURL *string,
	attributes []byte,
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
	switch version.Major() {
	case 2:
		if err := registerWorkflowV2(sc, workflowRegistryAddr, version, workflowName, workflowID, binaryURLToUse, configURLToUse, attributes); err != nil {
			return "", err
		}
	default:
		if err := registerWorkflowV1(sc, workflowRegistryAddr, donID, workflowName, workflowID, binaryURLToUse, configURLToUse, secretsURLToUse); err != nil {
			return "", err
		}
	}

	return workflowID, nil
}

func LinkOwner(sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) error {
	switch version.Major() {
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

		registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, version)
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
		return err
	default:
		return errors.New("invalid version for linking owner")
	}
}

// GetWorkflowNames retrieves all workflow names for the given registry contract.
// It supports both v1 and v2 workflow registry versions.
func GetWorkflowNames(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) ([]string, error) {
	switch version.Major() {
	case 2:
		return getWorkflowNamesV2(sc, workflowRegistryAddr, version)
	default:
		return getWorkflowNamesV1(sc, workflowRegistryAddr)
	}
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
		return deleteWorkflowV2(ctx, sc, workflowRegistryAddr, version, workflowName)
	default:
		return deleteWorkflowV1(ctx, sc, workflowRegistryAddr, workflowName)
	}
}

// DeleteAllWithContract removes all workflows owned by the caller from the workflow registry contract.
// It supports both v1 and v2 workflow registry versions.
func DeleteAllWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) error {
	switch version.Major() {
	case 2:
		return deleteAllWorkflowsV2(ctx, sc, workflowRegistryAddr, version)
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
func registerWorkflowV2(
	sc *seth.Client,
	workflowRegistryAddr common.Address,
	version *semver.Version,
	workflowName, workflowID, binaryURL, configURL string,
	attributes []byte,
) error {
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, version)
	if err != nil {
		return err
	}

	// Check and link owner if needed using existing helper function
	if verifyErr := verifyOwnerLinkedWithRegistry(registry, sc, workflowName); verifyErr != nil {
		// If owner is not linked, try to link them
		if linkErr := LinkOwner(sc, workflowRegistryAddr, version); linkErr != nil {
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
		attributes,
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
func deleteAllWorkflowsV2(_ context.Context, sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) error {
	// Create registry instance once for all operations
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, version)
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
func deleteAllWorkflowsV1(_ context.Context, sc *seth.Client, workflowRegistryAddr common.Address) error {
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
func deleteWorkflowV2(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version, workflowName string,
) error {
	// Create registry instance once for all operations
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, version)
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
func getRegistryV2Instance(sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) (*workflow_registry_wrapper_v2.WorkflowRegistry, error) {
	registry, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return nil, errors.Wrapf(err, errCreateContractInstance, "WorkflowRegistry", version)
	}

	// add contract ABI to Seth, so that it can decode transaction errors
	abi, aErr := workflow_registry_wrapper_v2.WorkflowRegistryMetaData.GetAbi()
	if aErr != nil {
		return nil, fmt.Errorf("failed to get WorkflowRegistryV2 ABI: %w", aErr)
	}

	sc.ABIFinder.ContractStore.AddABI("WorkflowRegistryV2", *abi)

	return registry, nil
}

// createRegistryV1Instance creates a new v1 workflow registry instance.
func createRegistryV1Instance(sc *seth.Client, workflowRegistryAddr common.Address) (*workflow_registry_wrapper.WorkflowRegistry, error) {
	registry, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return nil, errors.Wrap(err, errCreateRegistryInstance)
	}

	// add contract ABI to Seth, so that it can decode transaction errors
	abi, aErr := workflow_registry_wrapper.WorkflowRegistryMetaData.GetAbi()
	if aErr != nil {
		return nil, fmt.Errorf("failed to get WorkflowRegistryV1 ABI: %w", aErr)
	}

	sc.ABIFinder.ContractStore.AddABI("WorkflowRegistryV1", *abi)

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
func getWorkflowListV2(sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) ([]workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView, error) {
	registry, err := getRegistryV2Instance(sc, workflowRegistryAddr, version)
	if err != nil {
		return nil, err
	}

	return getWorkflowListWithRegistryV2(registry, sc)
}

// getWorkflowNamesV2 retrieves all workflow names for v2 registry contracts.
func getWorkflowNamesV2(sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) ([]string, error) {
	workflows, err := getWorkflowListV2(sc, workflowRegistryAddr, version)
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

// UpdateVaultCapabilityConfig merges the provided vaultPublicKey and threshold into the
// vault capability's DefaultConfig in the capabilities registry for the given DON,
// preserving any pre-existing fields. This is required so that workflow nodes can
// unwrap the capability config when calling runtime.GetSecret().
func UpdateVaultCapabilityConfig(ctx context.Context, sethClient *seth.Client, capabilitiesRegistryAddr string, don *capabilities_registry_v2.CapabilitiesRegistryDONInfo, vaultPublicKey string, threshold int) error {
	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(capabilitiesRegistryAddr), sethClient.Client,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create capabilities registry wrapper")
	}

	newConfigs := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, 0, len(don.CapabilityConfigurations))
	for _, cc := range don.CapabilityConfigurations {
		if cc.CapabilityId == vault_helpers.CapabilityID {
			existingCfg := &capabilitiespb.CapabilityConfig{}
			if len(cc.Config) > 0 {
				if unmarshalErr := proto.Unmarshal(cc.Config, existingCfg); unmarshalErr != nil {
					return errors.Wrap(unmarshalErr, "failed to unmarshal existing vault capability config")
				}
			}

			base := chainlinkvalues.EmptyMap()
			if existingCfg.DefaultConfig != nil {
				base, err = chainlinkvalues.FromMapValueProto(existingCfg.DefaultConfig)
				if err != nil {
					return errors.Wrap(err, "failed to convert existing vault capability config")
				}
			}
			newValues, wrapErr := chainlinkvalues.WrapMap(map[string]any{
				"VaultPublicKey": vaultPublicKey,
				"Threshold":      threshold,
			})
			if wrapErr != nil {
				return errors.Wrap(wrapErr, "failed to wrap vault capability config values")
			}
			for k, v := range newValues.Underlying {
				base.Underlying[k] = v
			}
			existingCfg.DefaultConfig = chainlinkvalues.ProtoMap(base)

			configBytes, marshalErr := proto.Marshal(existingCfg)
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "failed to marshal updated vault capability config")
			}

			cc.Config = configBytes
		}
		newConfigs = append(newConfigs, cc)
	}

	updateParams := capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
		Name:                     don.Name,
		Config:                   don.Config,
		CapabilityConfigurations: newConfigs,
		Nodes:                    don.NodeP2PIds,
		F:                        don.F,
		IsPublic:                 don.IsPublic,
	}

	_, updateErr := sethClient.Decode(capReg.UpdateDONByName(sethClient.NewTXOpts(), don.Name, updateParams))
	return errors.Wrap(updateErr, "UpdateDONByName tx failed")
}

// WaitForVaultConfigPropagation waits until every workflow node's local registry snapshot
// (registry_syncer_states) shows a non-nil DefaultConfig for the vault capability.
//
// It connects directly to each node's PostgreSQL database using the shared postgres server
// exposed at dbPort, with one database per node named db_0 … db_{nodeCount-1}.
//
// If dbPort or nodeCount is 0 (e.g. Kubernetes provider or missing state file), the function
// falls back to a static wait of vaultConfigStaticPropagationWait.
func WaitForVaultConfigPropagation(ctx context.Context, dbPort, nodeCount int) error {
	if dbPort == 0 || nodeCount == 0 {
		fmt.Printf("\n⚙️ Node DB info unavailable; waiting %s for vault config propagation\n", vaultConfigStaticPropagationWait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(vaultConfigStaticPropagationWait):
			return nil
		}
	}

	fmt.Printf("\n⚙️ Polling %d workflow node(s) on db port %d for vault capability config propagation\n", nodeCount, dbPort)

	const pollInterval = 2 * time.Second
	const pollTimeout = 2 * time.Minute

	timeoutCtx, cancel := context.WithTimeout(ctx, pollTimeout)
	defer cancel()

	pending := make(map[int]struct{}, nodeCount)
	for i := 0; i < nodeCount; i++ {
		pending[i] = struct{}{}
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		type checkResult struct {
			index int
			ready bool
			msg   string
		}

		results := make(chan checkResult, len(pending))
		eg, tickCtx := errgroup.WithContext(timeoutCtx)

		for nodeIndex := range pending {
			// capture for goroutine
			eg.Go(func() error {
				ready, msg := hasVaultCapabilityConfigOnNode(tickCtx, dbPort, nodeIndex)
				results <- checkResult{index: nodeIndex, ready: ready, msg: msg}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			close(results)
			return err
		}
		close(results)

		for r := range results {
			if r.ready {
				delete(pending, r.index)
				fmt.Printf("  ✅ node db_%d: vault config propagated\n", r.index)
			} else {
				fmt.Printf("  ⏳ node db_%d: %s\n", r.index, r.msg)
			}
		}

		if len(pending) == 0 {
			return nil
		}

		select {
		case <-timeoutCtx.Done():
			remaining := make([]int, 0, len(pending))
			for i := range pending {
				remaining = append(remaining, i)
			}
			return fmt.Errorf("timed out after %.0fs waiting for vault config propagation on nodes: %v", pollTimeout.Seconds(), remaining)
		case <-ticker.C:
		}
	}
}

// hasVaultCapabilityConfigOnNode queries db_{nodeIndex} at dbPort and returns true when the latest
// registry_syncer_states row contains a non-nil DefaultConfig for the vault capability.
func hasVaultCapabilityConfigOnNode(ctx context.Context, dbPort, nodeIndex int) (bool, string) {
	dsn := fmt.Sprintf(
		"host=127.0.0.1 port=%d user=%s password=%s dbname=db_%d sslmode=disable connect_timeout=3",
		dbPort, postgres.User, postgres.Password, nodeIndex,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return false, fmt.Sprintf("failed to open DB: %v", err)
	}
	defer db.Close()

	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var rawData []byte
	if err = db.GetContext(queryCtx, &rawData, latestRegistrySyncStateQuery); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "registry_syncer_states is empty"
		}
		return false, fmt.Sprintf("query failed: %v", err)
	}

	var state vaultSyncStatePayload
	if err = json.Unmarshal(rawData, &state); err != nil {
		return false, fmt.Sprintf("failed to unmarshal registry syncer state: %v", err)
	}

	for _, don := range state.IDsToDONs {
		capCfgEntry, ok := don.CapabilityConfigurations[vault_helpers.CapabilityID]
		if !ok || len(capCfgEntry.Config) == 0 {
			continue
		}

		capCfg := &capabilitiespb.CapabilityConfig{}
		if unmarshalErr := proto.Unmarshal(capCfgEntry.Config, capCfg); unmarshalErr != nil {
			return false, fmt.Sprintf("failed to unmarshal vault capability config protobuf: %v", unmarshalErr)
		}

		if capCfg.DefaultConfig != nil {
			return true, ""
		}
	}

	return false, fmt.Sprintf("vault capability %s DefaultConfig not yet set in any DON snapshot", vault_helpers.CapabilityID)
}

// VaultConfigHasPublicKey reports whether cfg.DefaultConfig already contains a
// VaultPublicKey field whose value equals publicKey. It is used to decide whether
// an update to the capabilities registry is necessary.
func VaultConfigHasPublicKey(cfg *capabilitiespb.CapabilityConfig, publicKey string) bool {
	if cfg == nil || cfg.DefaultConfig == nil {
		return false
	}
	existing, err := chainlinkvalues.FromMapValueProto(cfg.DefaultConfig)
	if err != nil {
		return false
	}
	v, ok := existing.Underlying["VaultPublicKey"]
	if !ok || v == nil {
		return false
	}
	val, err := v.Unwrap()
	if err != nil {
		return false
	}
	str, ok := val.(string)
	return ok && str == publicKey
}

// IsBase64File checks if the file at the given path is a base64-encoded file by reading a portion of it and attempting to decode it.
func IsBase64File(filename string) error {
	fileInfo, fErr := os.Stat(filename)
	if fErr != nil {
		return errors.Wrap(fErr, "failed to get file info")
	}

	readSize := min(fileInfo.Size(), 4*1024*1024) // 4MB

	file, oErr := os.Open(filename)
	if oErr != nil {
		return errors.Wrap(oErr, "failed to open file")
	}
	defer file.Close()

	buffer := make([]byte, readSize)
	n, rErr := file.Read(buffer)
	if rErr != nil && rErr != io.EOF {
		return errors.Wrap(rErr, "failed to read file")
	}

	if !isBase64Content(string(buffer[:n])) {
		return fmt.Errorf("❌ file %s is not a base64-encoded file", filename)
	}

	return nil
}

func isBase64Content(content string) bool {
	// Remove whitespace and newlines, just to be safe
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\t", "")

	if len(content) == 0 {
		return false
	}

	_, err := base64.StdEncoding.DecodeString(content)
	return err == nil
}
