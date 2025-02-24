package contracts

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	libnet "github.com/smartcontractkit/chainlink/system-tests/lib/net"
)

func RegisterWorkflow(sc *seth.Client, workflowRegistryAddr common.Address, donID uint32, workflowName, binaryURL, configURL string) error {
	workFlowData, err := libnet.DownloadAndDecodeBase64(binaryURL)
	if err != nil {
		return errors.Wrap(err, "failed to download and decode workflow binary")
	}

	var configData []byte
	if configURL != "" {
		configData, err = libnet.Download(configURL)
		if err != nil {
			return errors.Wrap(err, "failed to download workflow config")
		}
	}

	// use non-encoded workflow name
	workflowID, idErr := generateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
	if idErr != nil {
		return errors.Wrap(idErr, "failed to generate workflow ID")
	}

	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return errors.Wrap(err, "failed to create workflow registry instance")
	}

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), binaryURL, configURL, ""))
	if decodeErr != nil {
		return errors.Wrap(decodeErr, "failed to register workflow")
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
