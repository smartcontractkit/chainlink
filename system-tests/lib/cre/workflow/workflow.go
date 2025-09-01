package workflow

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libnet "github.com/smartcontractkit/chainlink/system-tests/lib/net"
)

func RegisterWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, donID uint64, workflowName, binaryURL string, configURL, secretsURL *string, artifactsDirInContainer *string) (string, error) {
	workFlowData, workFlowErr := libnet.DownloadAndDecodeBase64(ctx, binaryURL)
	if workFlowErr != nil {
		return "", errors.Wrap(workFlowErr, "failed to download and decode workflow binary")
	}

	var binaryURLToUse string
	if artifactsDirInContainer != nil {
		binaryURLToUse = fmt.Sprintf("file://%s/%s", *artifactsDirInContainer, filepath.Base(binaryURL))
	} else {
		binaryURLToUse = binaryURL
	}

	var configData []byte
	var configErr error
	configURLToUse := ""
	if configURL != nil && *configURL != "" {
		configData, configErr = libnet.Download(ctx, *configURL)
		if configErr != nil {
			return "", errors.Wrap(configErr, "failed to download workflow config")
		}

		if artifactsDirInContainer != nil {
			configURLToUse = fmt.Sprintf("file://%s/%s", *artifactsDirInContainer, filepath.Base(*configURL))
		} else {
			configURLToUse = *configURL
		}
	}

	secretsURLToUse := ""
	if secretsURL != nil && *secretsURL != "" {
		if artifactsDirInContainer != nil {
			secretsURLToUse = fmt.Sprintf("file://%s/%s", *artifactsDirInContainer, filepath.Base(*secretsURL))
		} else {
			secretsURLToUse = *secretsURL
		}
	}

	// use non-encoded workflow name
	workflowID, idErr := generateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, secretsURLToUse)
	if idErr != nil {
		return "", errors.Wrap(idErr, "failed to generate workflow ID")
	}

	workflowRegistryInstance, instanceErr := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if instanceErr != nil {
		return "", errors.Wrap(instanceErr, "failed to create workflow registry instance")
	}

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), libc.MustSafeUint32FromUint64(donID), uint8(0), binaryURLToUse, configURLToUse, secretsURLToUse))
	if decodeErr != nil {
		return "", errors.Wrap(decodeErr, "failed to register workflow")
	}

	return workflowID, nil
}

func GetWorkflowNames(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address) ([]string, error) {
	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create workflow registry instance")
	}

	metadataList, metadataListErr := workflowRegistryInstance.GetWorkflowMetadataListByOwner(sc.NewCallOpts(), sc.MustGetRootKeyAddress(), big.NewInt(0), big.NewInt(10))
	if metadataListErr != nil {
		return nil, errors.Wrap(metadataListErr, "failed to get workflow metadata list")
	}

	workflows := make([]string, 0)

	for _, metadata := range metadataList {
		workflows = append(workflows, metadata.WorkflowName)
	}

	return workflows, nil
}

func DeleteAllWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address) error {
	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return errors.Wrap(err, "failed to create workflow registry instance")
	}

	metadataList, metadataListErr := workflowRegistryInstance.GetWorkflowMetadataListByOwner(sc.NewCallOpts(), sc.MustGetRootKeyAddress(), big.NewInt(0), big.NewInt(10))
	if metadataListErr != nil {
		return errors.Wrap(metadataListErr, "failed to get workflow metadata list")
	}

	var computeHashKey = func(owner common.Address, workflowName string) [32]byte {
		ownerBytes := owner.Bytes()
		nameBytes := []byte(workflowName)
		data := make([]byte, len(ownerBytes)+len(nameBytes))
		copy(data, ownerBytes)
		copy(data[len(ownerBytes):], nameBytes)

		return crypto.Keccak256Hash(data)
	}

	for _, metadata := range metadataList {
		workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), metadata.WorkflowName)
		_, deleteErr := sc.Decode(workflowRegistryInstance.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey))
		if deleteErr != nil {
			return errors.Wrap(deleteErr, "failed to delete workflow named "+metadata.WorkflowName)
		}
	}

	return nil
}

func DeleteWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, workflowName string) error {
	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return errors.Wrap(err, "failed to create workflow registry instance")
	}

	var computeHashKey = func(owner common.Address, workflowName string) [32]byte {
		ownerBytes := owner.Bytes()
		nameBytes := []byte(workflowName)
		data := make([]byte, len(ownerBytes)+len(nameBytes))
		copy(data, ownerBytes)
		copy(data[len(ownerBytes):], nameBytes)

		return crypto.Keccak256Hash(data)
	}

	workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflowName)
	_, deleteErr := sc.Decode(workflowRegistryInstance.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey))
	if deleteErr != nil {
		return errors.Wrap(deleteErr, "failed to delete workflow named "+workflowName)
	}

	return nil
}

func RemoveWorkflowArtifactsFromLocalEnv(workflowArtifactsLocations ...string) error {
	for _, artifactLocation := range workflowArtifactsLocations {
		if artifactLocation == "" {
			continue
		}

		err := os.Remove(artifactLocation)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to remove workflow artifact located at %s: %s", artifactLocation, err.Error()))
		}
	}
	return nil
}

func generateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := pkgworkflows.GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}
