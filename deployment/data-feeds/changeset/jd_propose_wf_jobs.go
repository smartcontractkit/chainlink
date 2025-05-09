package changeset

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/view/v1_0"
)

const (
	timeout = 120 * time.Second
)

// ProposeWFJobsToJDChangeset is a changeset that reads a feed state file, creates a workflow job spec from it and proposes it to JD.
var ProposeWFJobsToJDChangeset = cldf.CreateChangeSet(proposeWFJobsToJDLogic, proposeWFJobsToJDPrecondition)

func proposeWFJobsToJDLogic(env deployment.Environment, c types.ProposeWFJobsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), timeout)
	defer cancel()

	chainInfo, _ := deployment.ChainInfo(c.ChainSelector)

	feedStatePath := filepath.Join("feeds", chainInfo.ChainName+".json")
	feedState, _ := LoadJSON[*v1_0.FeedState](feedStatePath, c.InputFS)

	// Add extra padded zeros to the feed IDs
	for i := range feedState.Feeds {
		extraPaddedZeros := strings.Repeat("0", 32)
		feedState.Feeds[i].FeedID += extraPaddedZeros
	}

	workflowSpecConfig := c.WorkflowSpecConfig
	workflowState := feedState.Workflows[workflowSpecConfig.WorkflowName]

	//nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	cacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, c.ChainSelector, &c.CacheLabel)

	// default values
	consensusEncoderAbi, _ := getWorkflowConsensusEncoderAbi(workflowSpecConfig.TargetContractEncoderType)

	consensusConfigKeyID := workflowSpecConfig.ConsensusConfigKeyID
	if consensusConfigKeyID == "" {
		consensusConfigKeyID = "evm"
	}

	consensusRef := workflowSpecConfig.ConsensusRef
	if consensusRef == "" {
		consensusRef = "data-feeds"
	}

	var triggersMaxFrequencyMs int
	if workflowSpecConfig.TriggersMaxFrequencyMs == nil {
		triggersMaxFrequencyMs = 5000
	} else {
		triggersMaxFrequencyMs = *workflowSpecConfig.TriggersMaxFrequencyMs
	}

	var deltaStageSec int
	if workflowSpecConfig.DeltaStageSec == nil {
		deltaStageSec = 45
	} else {
		deltaStageSec = *workflowSpecConfig.DeltaStageSec
	}

	targetSchedule := workflowSpecConfig.TargetsSchedule
	if targetSchedule == "" {
		targetSchedule = "oneAtATime"
	}

	consensusPartialStaleness := workflowSpecConfig.ConsensusAllowedPartialStaleness
	if consensusPartialStaleness == "" {
		consensusPartialStaleness = "0.5"
	}

	// create the workflow YAML spec
	workflowSpec, err := offchain.CreateWorkflowSpec(
		feedState.Feeds,
		workflowSpecConfig.WorkflowName,
		workflowState.Owner,
		triggersMaxFrequencyMs,
		consensusRef,
		workflowSpecConfig.ConsensusReportID,
		workflowSpecConfig.ConsensusAggregationMethod,
		consensusConfigKeyID,
		consensusPartialStaleness,
		consensusEncoderAbi,
		deltaStageSec,
		workflowSpecConfig.WriteTargetTrigger,
		targetSchedule,
		workflowSpecConfig.CREStepTimeout,
		cacheAddress,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create workflow spec: %w", err)
	}

	// create workflow job spec TOML
	workflowJobSpec, err := offchain.JobSpecFromWorkflowSpec(workflowSpec, c.WorkflowJobName, workflowState.Owner)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create workflow job spec: %w", err)
	}

	// propose workflow jobs to JD
	out, err := offchain.ProposeJobs(ctx, env, workflowJobSpec, &workflowSpecConfig.WorkflowName, c.NodeFilter)
	if err != nil {
		env.Logger.Debugf("%s", workflowJobSpec)
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose workflow job spec: %w", err)
	}

	// write the workflow spec to artifacts/migration_name folder. Do not exit on error as the jobs are already proposed
	baseDir := ".."
	envName := env.Name

	wfSpecPath := filepath.Join(baseDir, envName, "artifacts", c.MigrationName, workflowState.Name+".yaml")
	err = os.MkdirAll(filepath.Dir(wfSpecPath), 0755)
	if err != nil {
		env.Logger.Errorf("failed to create directory for workflow file: %s", err)
		env.Logger.Debugf("%s", workflowJobSpec)
		return out, nil
	}

	err = os.WriteFile(wfSpecPath, []byte(workflowSpec), 0600)
	if err != nil {
		env.Logger.Errorf("failed to write workflow to file: %s", err)
		env.Logger.Debugf("%s", workflowJobSpec)
	}

	return out, nil
}

func proposeWFJobsToJDPrecondition(env deployment.Environment, c types.ProposeWFJobsConfig) error {
	if c.MigrationName == "" {
		return errors.New("migration name is required")
	}

	if c.WorkflowJobName == "" {
		return errors.New("workflow job name is required")
	}

	if c.WorkflowSpecConfig.WorkflowName == "" {
		return errors.New("workflow name is required")
	}

	validTargetEncoder := c.WorkflowSpecConfig.TargetContractEncoderType
	if validTargetEncoder != "data-feeds_decimal" && validTargetEncoder != "aptos" && validTargetEncoder != "ccip" {
		return fmt.Errorf("invalid consensus target encoder: %s", c.WorkflowSpecConfig.ConsensusAggregationMethod)
	}

	validMethod := c.WorkflowSpecConfig.ConsensusAggregationMethod
	if validMethod != "data_feeds" && validMethod != "llo_streams" {
		return fmt.Errorf("invalid consensus aggregation method: %s", c.WorkflowSpecConfig.ConsensusAggregationMethod)
	}

	if c.WorkflowSpecConfig.ConsensusReportID == "" {
		return errors.New("consensus report id is required")
	}

	if c.WorkflowSpecConfig.WriteTargetTrigger == "" {
		return errors.New("write target trigger is required")
	}
	if _, err := getWorkflowConsensusEncoderAbi(c.WorkflowSpecConfig.TargetContractEncoderType); err != nil {
		return fmt.Errorf("failed to get consensus encoder abi: %w", err)
	}

	chainInfo, err := deployment.ChainInfo(c.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chain info for chain %d: %w", c.ChainSelector, err)
	}
	inputFileName := filepath.Join("feeds", chainInfo.ChainName+".json")

	feedState, err := LoadJSON[*v1_0.FeedState](inputFileName, c.InputFS)
	if err != nil {
		return fmt.Errorf("failed to load feed mapping from %s: %w", inputFileName, err)
	}

	err = feedState.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate feeds: %w", err)
	}

	workflowState := feedState.Workflows[c.WorkflowSpecConfig.WorkflowName]
	if workflowState.Name == "" || workflowState.Owner == "" || workflowState.Forwarder == "" {
		return fmt.Errorf("no workflow found for hash %s in %s", c.WorkflowSpecConfig.WorkflowName, inputFileName)
	}

	//nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	cacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, c.ChainSelector, &c.CacheLabel)
	if cacheAddress == "" {
		return errors.New("failed to get data feeds cache address")
	}

	return nil
}

func getWorkflowConsensusEncoderAbi(targetContractEncoderType string) (string, error) {
	switch targetContractEncoderType {
	case "data-feeds_decimal":
		return "(bytes32 RemappedID, uint32 Timestamp, uint224 Price)[] Reports", nil
	case "ccip":
		return "(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports", nil
	case "aptos":
		return "(bytes32 RemappedID, bytes RawReport)[] Reports", nil
	default:
		return "", fmt.Errorf("unknown consensus encoder type %s", targetContractEncoderType)
	}
}
