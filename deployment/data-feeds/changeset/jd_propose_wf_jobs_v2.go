package changeset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/view/v1_0"
)

const (
	timeoutV2 = 240 * time.Second
)

// ProposeWFJobsToJDV2Changeset is a Durable Pipeline compatible changeset that reads a feed state file,
// creates a workflow job spec from it and proposes it to JD.
var ProposeWFJobsToJDV2Changeset = cldf.CreateChangeSet(proposeWFJobsToJDV2Logic, proposeWFJobsToJDV2Precondition)

func proposeWFJobsToJDV2Logic(env cldf.Environment, c types.ProposeWFJobsV2Config) (cldf.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), timeoutV2)
	defer cancel()

	chainInfo, _ := cldf_chain_utils.ChainInfo(c.ChainSelector)

	domain := getDomain(c.Domain)

	feedStatePath := filepath.Join("domains", domain, env.Name, "inputs", "feeds", c.NodeFilter.Zone, chainInfo.ChainName+".json")
	feedState, _ := readFeedStateFile(feedStatePath)

	// Only get feeds that are part of the workflow
	feeds := *getFeedsByWorkflow(&feedState.Feeds, c.WorkflowSpecConfig.WorkflowName)

	// Add extra padded zeros to the feed IDs
	for i := range feeds {
		extraPaddedZeros := strings.Repeat("0", 32)
		feeds[i].FeedID += extraPaddedZeros
	}

	workflowSpecConfig := c.WorkflowSpecConfig
	workflowState := feedState.Workflows[workflowSpecConfig.WorkflowName]

	// Addressbook is deprecated, but we still use it for the time being
	cacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, env.DataStore.Addresses(), c.ChainSelector, &c.CacheLabel)

	// default values
	consensusEncoderAbi, _ := _getWorkflowConsensusEncoderAbi(workflowSpecConfig.TargetContractEncoderType)
	consensusConfigKeyID := getConsensusConfigKey(workflowSpecConfig.ConsensusConfigKeyID)
	consensusRef := getConsensusRef(workflowSpecConfig.ConsensusRef)
	triggersMaxFrequencyMs := getMaxFrequencyMs(workflowSpecConfig.TriggersMaxFrequencyMs)
	deltaStageSec := getDeltaStage(workflowSpecConfig.DeltaStageSec)
	targetSchedule := getTargetSchedule(workflowSpecConfig.TargetsSchedule)

	// create the workflow YAML spec
	workflowSpec, err := offchain.CreateWorkflowSpec(
		feeds,
		workflowSpecConfig.WorkflowName,
		workflowState.Owner,
		triggersMaxFrequencyMs,
		consensusRef,
		workflowSpecConfig.ConsensusReportID,
		workflowSpecConfig.ConsensusAggregationMethod,
		consensusConfigKeyID,
		workflowSpecConfig.ConsensusAllowedPartialStaleness,
		consensusEncoderAbi,
		deltaStageSec,
		workflowSpecConfig.WriteTargetTrigger,
		targetSchedule,
		workflowSpecConfig.CREStepTimeout,
		workflowSpecConfig.TargetProcessor,
		cacheAddress,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create workflow spec: %w", err)
	}

	// log the workflow spec for debugging purposes.
	fmt.Println(workflowSpec)

	// create workflow job spec TOML
	workflowJobSpec, err := offchain.JobSpecFromWorkflowSpec(workflowSpec, c.WorkflowJobName, workflowState.Owner)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create workflow job spec: %w", err)
	}

	// propose workflow jobs to JD
	out, err := offchain.ProposeJobs(ctx, env, workflowJobSpec, &workflowSpecConfig.WorkflowName, c.NodeFilter)
	if err != nil {
		env.Logger.Debugf("%s", workflowJobSpec)
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose workflow job spec: %w", err)
	}

	// Save workflow spec in the datastore
	ds := datastore.NewMemoryDataStore()

	err = UpdateWorkflowMetadataDS(env, ds, workflowSpecConfig.WorkflowName, workflowSpec)
	if err != nil {
		env.Logger.Errorf("failed to set workflow spec in datastore: %s", err)
	}

	out.DataStore = ds

	return out, nil
}

func proposeWFJobsToJDV2Precondition(env cldf.Environment, c types.ProposeWFJobsV2Config) error {
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

	if c.NodeFilter == nil {
		return errors.New("missing node filter")
	}

	if c.NodeFilter.DONID == 0 {
		return errors.New("missing DON ID in node filter")
	}

	if c.NodeFilter.EnvLabel == "" {
		return errors.New("missing environment label in node filter")
	}

	if c.NodeFilter.Zone != "zone-a" && c.NodeFilter.Zone != "zone-b" {
		return errors.New("missing or invalid zone in node filter")
	}

	domain := getDomain(c.Domain)
	chainInfo, err := cldf_chain_utils.ChainInfo(c.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chain info for chain %d: %w", c.ChainSelector, err)
	}

	feedStatePath := filepath.Join("domains", domain, env.Name, "inputs", "feeds", c.NodeFilter.Zone, chainInfo.ChainName+".json")

	feedState, err := readFeedStateFile(feedStatePath)
	if err != nil {
		return fmt.Errorf("failed to read feed state file %s: %w", feedStatePath, err)
	}

	feeds := *getFeedsByWorkflow(&feedState.Feeds, c.WorkflowSpecConfig.WorkflowName)
	feedState.Feeds = feeds

	err = feedState.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate feeds: %w", err)
	}

	workflowState := feedState.Workflows[c.WorkflowSpecConfig.WorkflowName]
	if workflowState.Name == "" || workflowState.Owner == "" || workflowState.Forwarder == "" {
		return fmt.Errorf("no workflow found for hash %s in %s", c.WorkflowSpecConfig.WorkflowName, feedStatePath)
	}

	// Addressbook is deprecated, but we still use it for the time being
	cacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, env.DataStore.Addresses(), c.ChainSelector, &c.CacheLabel)
	if cacheAddress == "" {
		return errors.New("failed to get data feeds cache address")
	}

	return nil
}

func _getWorkflowConsensusEncoderAbi(targetContractEncoderType string) (string, error) {
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

func getDomain(domain string) string {
	if domain == "" {
		return "data-feeds"
	}
	return domain
}

func readFeedStateFile(inputFileName string) (*v1_0.FeedState, error) {
	content, err := os.ReadFile(inputFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load feed mapping from %s: %w", inputFileName, err)
	}

	var feedState *v1_0.FeedState

	err = json.Unmarshal(content, &feedState)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed state from %s: %w", inputFileName, err)
	}
	return feedState, nil
}

func getConsensusConfigKey(consensusConfigKeyID string) string {
	if consensusConfigKeyID == "" {
		return "evm"
	}
	return consensusConfigKeyID
}

func getConsensusRef(consensusRef string) string {
	if consensusRef == "" {
		return "data-feeds"
	}
	return consensusRef
}

func getMaxFrequencyMs(maxFrequencyMs *int) int {
	if maxFrequencyMs == nil {
		return 5000
	}
	return *maxFrequencyMs
}

func getDeltaStage(deltaStageSec *int) int {
	if deltaStageSec == nil {
		return 45
	}
	return *deltaStageSec
}

func getTargetSchedule(targetSchedule string) string {
	if targetSchedule == "" {
		return "oneAtATime"
	}
	return targetSchedule
}

func getFeedsByWorkflow(allFeeds *[]v1_0.Feed, workflowHash string) *[]v1_0.Feed {
	var result []v1_0.Feed
	for _, feed := range *allFeeds {
		for _, workflow := range feed.Workflows {
			if workflow == workflowHash {
				result = append(result, feed)
				break
			}
		}
	}

	return &result
}
