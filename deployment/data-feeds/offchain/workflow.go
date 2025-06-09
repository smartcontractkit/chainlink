package offchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"text/template"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/view/v1_0"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

const (
	workflowJobPath  = "workflow.tmpl"
	workflowSpecPath = "workflow_spec.yaml.tmpl"
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

type WorkflowSpecCfg struct {
	Feeds                            []v1_0.Feed
	WorkflowName                     string
	WorkflowOwner                    string
	TriggersMaxFrequencyMs           int
	ConsensusRef                     string
	ConsensusReportID                string
	ConsensusConfigKeyID             string
	ConsensusAggregationMethod       string
	ConsensusAllowedPartialStaleness string
	ConsensusEncoderABI              string
	DeltaStageSec                    int
	WriteTargetTrigger               string
	TargetsSchedule                  string
	TargetProcessor                  string
	CREStepTimeout                   int64
	DataFeedsCacheAddress            string
}

func JobSpecFromWorkflowSpec(workflowSpec string, workflowJobName string, workflowOwner string) (wfSpec string, err error) {
	externalID := uuid.New().String()

	wfCfg := WorkflowJobCfg{
		JobName:       workflowJobName,
		ExternalJobID: externalID,
		Workflow:      workflowSpec,
		WorkflowID:    getWorkflowID(workflowSpec),
		WorkflowOwner: workflowOwner,
	}

	workflowJobSpec, err := wfCfg.createSpec()
	if err != nil {
		return "", fmt.Errorf("failed to create workflow job spec: %w", err)
	}
	return workflowJobSpec, nil
}

func CreateWorkflowSpec(
	feeds []v1_0.Feed,
	workflowName string,
	workflowOwner string,
	triggerFrequency int,
	consensusRef string,
	consensusReportID string,
	consensusAggregationMethod string,
	consensusConfigKeyID string,
	consensusAllowedPartialStaleness string,
	consensusEncoderABI string,
	deltaStageSec int,
	writeTargetTrigger string,
	targetsSchedule string,
	creStepTimeout int64,
	targetProcessor string,
	cacheAddress string) (wfSpec string, err error) {
	wfCfg := WorkflowSpecCfg{
		Feeds:                            feeds,
		WorkflowName:                     workflowName,
		WorkflowOwner:                    workflowOwner,
		TriggersMaxFrequencyMs:           triggerFrequency,
		ConsensusRef:                     consensusRef,
		ConsensusReportID:                consensusReportID,
		ConsensusAggregationMethod:       consensusAggregationMethod,
		ConsensusConfigKeyID:             consensusConfigKeyID,
		ConsensusAllowedPartialStaleness: consensusAllowedPartialStaleness,
		ConsensusEncoderABI:              consensusEncoderABI,
		DeltaStageSec:                    deltaStageSec,
		WriteTargetTrigger:               writeTargetTrigger,
		TargetsSchedule:                  targetsSchedule,
		TargetProcessor:                  targetProcessor,
		CREStepTimeout:                   creStepTimeout,
		DataFeedsCacheAddress:            cacheAddress,
	}

	workflowSpec, err := wfCfg.createSpec()
	if err != nil {
		return "", fmt.Errorf("failed to create workflow job spec: %w", err)
	}
	return workflowSpec, nil
}

func (wfCfg *WorkflowJobCfg) createSpec() (string, error) {
	t, err := template.New("s").ParseFS(offchainFs, workflowJobPath)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, workflowJobPath, wfCfg)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (wfCfg *WorkflowSpecCfg) createSpec() (string, error) {
	t, err := template.New("s").ParseFS(offchainFs, workflowSpecPath)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, workflowSpecPath, wfCfg)
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
