package offchain

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"text/template"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

const (
	workflowPath = "workflow.tmpl"
)

type WorkflowSpecAlias sdk.WorkflowSpec

type WorkflowJobCfg struct {
	JobName       string
	SpecFileName  string
	ExternalJobID string
	Workflow      string // yaml of the workflow
	WorkflowID    string
	WorkflowOwner string
}

func JobSpecFromWorkflow(inputFs embed.FS, inputFileName string, workflowJobName string) (string, error) {
	wfYaml, err := inputFs.ReadFile(inputFileName)
	if err != nil {
		return "", fmt.Errorf("failed to read workflow file: %w", err)
	}
	wfStr := string(wfYaml)
	wf, err := workflows.ParseWorkflowSpecYaml(wfStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse workflow spec: %w", err)
	}

	wfAlias := WorkflowSpecAlias(wf)
	if err := wfAlias.validate(); err != nil {
		return "", fmt.Errorf("workflow validation failed: %w", err)
	}

	externalID, err := createExternalJobID(wf.Name, wf.Owner)
	if err != nil {
		return "", fmt.Errorf("failed to get external job id: %w", err)
	}

	wfCfg := WorkflowJobCfg{
		JobName:       workflowJobName,
		ExternalJobID: externalID,
		Workflow:      wfStr,
		WorkflowID:    getWorkflowID(wfStr),
		WorkflowOwner: wf.Owner,
	}

	workflowJobSpec, err := wfCfg.createSpec()
	if err != nil {
		return "", fmt.Errorf("failed to create workflow job spec: %w", err)
	}
	return workflowJobSpec, nil
}

func (wf WorkflowSpecAlias) validate() error {
	configMap := wf.Consensus[0].Config
	aggregationConfig, ok := configMap["aggregation_config"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid aggregation_config type for workflow %s", wf.Name)
	}

	feeds, ok := aggregationConfig["feeds"]
	if !ok {
		return fmt.Errorf("feeds not found in aggregation_config for workflow %s", wf.Name)
	}
	for streamsID, feed := range feeds.(map[string]interface{}) {
		feedMap, feedMapExists := feed.(map[string]string)
		if !feedMapExists {
			return fmt.Errorf("invalid feed type %s", streamsID)
		}
		_, deviation := feedMap["deviation"]
		if !deviation {
			return fmt.Errorf("deviation not found in feed %s", streamsID)
		}
		_, heartbeat := feedMap["heartbeat"]
		if !heartbeat {
			return fmt.Errorf("heartbeat not found in feed %s", streamsID)
		}
		remmapedID, hasRemmapedID := feedMap["remappedID"]
		if !hasRemmapedID {
			return fmt.Errorf("remappedID not found in feed %s", streamsID)
		}
		if len(remmapedID) != 66 {
			return fmt.Errorf("invalid remappedID for feed %s", streamsID)
		}
	}

	encoderConfig, ok := configMap["encoder_config"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid encoder_config type for workflow %s", wf.Name)
	}
	encoderABI, ok := encoderConfig["abi"]
	if !ok {
		return fmt.Errorf("abi not found in encoder_config for workflow %s", wf.Name)
	}
	if encoderABI != "(bytes32 RemappedID, uint32 Timestamp, uint224 Price)[] Reports" {
		return fmt.Errorf("invalid encoder ABI for workflow %s", wf.Name)
	}
	return nil
}

func (wfCfg *WorkflowJobCfg) createSpec() (string, error) {
	t, err := template.New("s").ParseFS(offchainFs, workflowPath)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, workflowPath, wfCfg)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func getWorkflowID(wf string) string {
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(wf))
	cid := sha256Hash.Sum(nil)
	return hex.EncodeToString(cid)
}

func createExternalJobID(name, ownerAddress string) (string, error) {
	// this must be constant for a given logical wf so that the job distributor can
	// track the job
	if len(name) != 10 {
		return "", fmt.Errorf("workflow name must be 10 characters long, got %s", name)
	}
	if !gethcommon.IsHexAddress(ownerAddress) {
		return "", fmt.Errorf("invalid owner address %s", ownerAddress)
	}
	checksummed := gethcommon.HexToAddress(ownerAddress).Hex()
	id := []byte(name + checksummed)
	sha256Hash := sha256.New()
	sha256Hash.Write(id)
	id = sha256Hash.Sum(nil)

	return externalJobID(id, "workflow")
}

func externalJobID(wfid []byte, nodeID string) (string, error) {
	if len(wfid) == 0 {
		return "", errors.New("empty workflow id")
	}
	if len(wfid) < 16 {
		return "", fmt.Errorf("workflow id too short. must be at least 16 bytes got %d", len(wfid))
	}

	externalJobID := wfid[:16]
	// ensure deterministic uniqueness of the externalJobID
	nb := []byte(nodeID)
	sha256Hash := sha256.New()
	sha256Hash.Write(nb)
	nb = sha256Hash.Sum(nil)

	for i, b := range nb[:16] {
		externalJobID[i] ^= b
	}
	// tag as valid UUID v4 https://github.com/google/uuid/blob/0f11ee6918f41a04c201eceeadf612a377bc7fbc/version4.go#L53-L54
	externalJobID[6] = (externalJobID[6] & 0x0f) | 0x40 // Version 4
	externalJobID[8] = (externalJobID[8] & 0x3f) | 0x80 // Variant is 10

	id, err := uuid.FromBytes(externalJobID)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
