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

func LinkOwner(
	sc *seth.Client,
	workflowRegistryAddr common.Address,
	tv deployment.TypeAndVersion,
) error {
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

		wr, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(
			workflowRegistryAddr,
			sc.Client,
		)
		if err != nil {
			return errors.Wrapf(err, "could not get instance of %s %s", tv.Type, tv.Version)
		}

		version, err := wr.TypeAndVersion(sc.NewCallOpts())
		if err != nil {
			return err
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

		_, err = sc.Decode(wr.LinkOwner(sc.NewTXOpts(), validityTimestamp, common.HexToHash(ownershipProof), signature))
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("invalid version for linking owner")
	}
}

func RegisterWithContract(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, typeVersion deployment.TypeAndVersion,
	donID uint64, workflowName, binaryURL string,
	configURL, secretsURL *string,
	artifactsDirInContainer *string,
) (string, error) {
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

	switch typeVersion.Version.Major() {
	case 2:
		wr, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(
			workflowRegistryAddr,
			sc.Client,
		)
		if err != nil {
			return "", errors.Wrapf(err, "could not get instance of %s %s", typeVersion.Type, typeVersion.Version)
		}

		addr := sc.MustGetRootKeyAddress()
		linked, err := wr.IsOwnerLinked(sc.NewCallOpts(), addr)
		if err != nil {
			return "", errors.Wrap(err, "failed to check link status of owner")
		}

		if !linked {
			err := LinkOwner(sc, workflowRegistryAddr, typeVersion)
			if err != nil {
				return "", errors.Wrap(err, "failed to link owner to org")
			}
		}

		_, decodeErr := sc.Decode(wr.UpsertWorkflow(sc.NewTXOpts(), workflowName, "some-tag", [32]byte(common.Hex2Bytes(workflowID)), uint8(0), contracts.DonFamily, binaryURLToUse, configURLToUse, nil, false))
		if decodeErr != nil {
			return "", errors.Wrap(decodeErr, "failed to register workflow")
		}
		return workflowID, nil
	default:
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
}

func GetWorkflowNames(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, tv deployment.TypeAndVersion,
) ([]string, error) {
	workflows := make([]string, 0)
	switch tv.Version.Major() {
	case 2:
		wr, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(
			workflowRegistryAddr,
			sc.Client,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get instance of %s %s", tv.Type, tv.Version)
		}

		md, err := wr.GetWorkflowListByOwner(sc.NewCallOpts(), sc.MustGetRootKeyAddress(), big.NewInt(0), big.NewInt(10))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get workflow metadata list")
		}

		for _, m := range md {
			workflows = append(workflows, m.WorkflowName)
		}
		return workflows, nil
	default:
		workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create workflow registry instance")
		}
		metadataList, metadataListErr := workflowRegistryInstance.GetWorkflowMetadataListByOwner(sc.NewCallOpts(), sc.MustGetRootKeyAddress(), big.NewInt(0), big.NewInt(10))
		if metadataListErr != nil {
			return nil, errors.Wrap(metadataListErr, "failed to get workflow metadata list")
		}

		for _, metadata := range metadataList {
			workflows = append(workflows, metadata.WorkflowName)
		}

		return workflows, nil
	}
}

func DeleteAllWithContract(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, tv deployment.TypeAndVersion,
) error {
	switch tv.Version.Major() {
	case 2:
		wr, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(
			workflowRegistryAddr,
			sc.Client,
		)
		if err != nil {
			return errors.Wrapf(err, "could not get instance of %s %s", tv.Type, tv.Version)
		}

		md, err := wr.GetWorkflowListByOwner(sc.NewCallOpts(), sc.MustGetRootKeyAddress(), big.NewInt(0), big.NewInt(10))
		if err != nil {
			return errors.Wrap(err, "failed to get workflow metadata list")
		}

		for _, m := range md {
			workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), m.WorkflowName)
			if _, deleteErr := sc.Decode(wr.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey)); deleteErr != nil {
				return errors.Wrapf(deleteErr, "failed to delete workflow named %s", m.WorkflowName)
			}
		}
		return nil
	default:
		workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
		if err != nil {
			return errors.Wrap(err, "failed to create workflow registry instance")
		}

		metadataList, metadataListErr := workflowRegistryInstance.GetWorkflowMetadataListByOwner(sc.NewCallOpts(), sc.MustGetRootKeyAddress(), big.NewInt(0), big.NewInt(10))
		if metadataListErr != nil {
			return errors.Wrap(metadataListErr, "failed to get workflow metadata list")
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
}

func computeHashKey(owner common.Address, workflowName string) [32]byte {
	ownerBytes := owner.Bytes()
	nameBytes := []byte(workflowName)
	data := make([]byte, len(ownerBytes)+len(nameBytes))
	copy(data, ownerBytes)
	copy(data[len(ownerBytes):], nameBytes)

	return crypto.Keccak256Hash(data)
}

func DeleteWithContract(ctx context.Context, sc *seth.Client,
	workflowRegistryAddr common.Address, tv deployment.TypeAndVersion,
	workflowName string,
) error {
	switch tv.Version.Major() {
	case 2:
		wr, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(
			workflowRegistryAddr,
			sc.Client,
		)
		if err != nil {
			return errors.Wrapf(err, "could not get instance of %s %s", tv.Type, tv.Version)
		}

		workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflowName)
		if _, deleteErr := sc.Decode(wr.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey)); deleteErr != nil {
			return errors.Wrap(deleteErr, "failed to delete workflow named "+workflowName)
		}
		return nil
	default:
		workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
		if err != nil {
			return errors.Wrap(err, "failed to create workflow registry instance")
		}

		workflowHashKey := computeHashKey(sc.MustGetRootKeyAddress(), workflowName)
		_, deleteErr := sc.Decode(workflowRegistryInstance.DeleteWorkflow(sc.NewTXOpts(), workflowHashKey))
		if deleteErr != nil {
			return errors.Wrap(deleteErr, "failed to delete workflow named "+workflowName)
		}

		return nil
	}
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
