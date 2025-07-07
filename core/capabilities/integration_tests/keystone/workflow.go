package keystone

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

const hardcodedWorkflow = `
name: "%s"
owner: "0x%s"
triggers:
  - id: "streams-trigger@1.0.0"
    config:
      feedIds:
%s

consensus:
  - id: "offchain_reporting@1.0.0"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      report_id: "0001"
      key_id: "evm"	
      aggregation_method: "data_feeds"
      aggregation_config:
        feeds:
%s
      encoder: "EVM"
      encoder_config:
        abi: "(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports"

targets:
  - id: "write_geth-testnet@1.0.0"
    inputs:
      signed_report: "$(evm_median.outputs)"
    config:
      address: "%s"
      params: ["$(report)"]
      abi: "receive(report bytes)"
      deltaStage: %s
      schedule: %s
`

func createKeystoneWorkflowJob(t *testing.T,
	workflowName string,
	workflowOwner string,
	feedIDs []string,
	consumerAddr common.Address,
	deltaStage string,
	schedule string) job.Job {
	triggerFeedIDs := ""
	for _, feedID := range feedIDs {
		triggerFeedIDs += fmt.Sprintf("        - \"%s\"\n", feedID)
	}

	aggregationFeeds := ""
	for _, feedID := range feedIDs {
		aggregationFeeds += fmt.Sprintf("          \"%s\":\n            deviation: \"0.001\"\n            heartbeat: 3600\n", feedID)
	}

	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, fmt.Sprintf(hardcodedWorkflow, workflowName, workflowOwner, triggerFeedIDs, aggregationFeeds,
		consumerAddr.String(), deltaStage, schedule))
	return workflowJobSpec.Job()
}

const lloStreamsWorkflow = `
name: "%s"
owner: "0x%s"
triggers:
  - id: "streams-trigger:don_16nodes@2.0.0"
    config:
      feedIds:
%s

consensus:
  - id: "offchain_reporting@1.0.0"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      report_id: "0001"
      key_id: "evm"	
      aggregation_method: "llo_streams"
      aggregation_config:
        streams:
%s
      encoder: "EVM"
      encoder_config:
        abi: "(bytes32 RemappedID, uint224 Price, uint32 Timestamp)[] Reports"

targets:
  - id: "write_geth-testnet@1.0.0"
    inputs:
      signed_report: "$(evm_median.outputs)"
    config:
      address: "%s"
      params: ["$(report)"]
      abi: "receive(report bytes)"
      deltaStage: 1s
      schedule: oneAtATime
`

func createLLOStreamWorkflowJob(t *testing.T,
	workflowName string,
	workflowOwner string,
	streamIDremapped map[uint32]string,
	consumerAddr common.Address) job.Job {
	triggerFeedIDs := ""
	// keys of the map are stream IDs
	streamIDs := make([]uint32, 0, len(streamIDremapped))
	for streamID := range streamIDremapped {
		streamIDs = append(streamIDs, streamID)
	}
	slices.Sort(streamIDs)
	for _, streamID := range streamIDs {
		triggerFeedIDs += fmt.Sprintf("        - \"%d\"\n", streamID)
	}

	aggregationFeeds := ""
	for _, streamID := range streamIDs {
		aggregationFeeds += fmt.Sprintf("          \"%d\":\n            deviation: \"0.001\"\n            heartbeat: 3600\n            remappedID: \"%s\"\n", streamID, streamIDremapped[streamID])
	}

	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, fmt.Sprintf(lloStreamsWorkflow, workflowName, workflowOwner, triggerFeedIDs, aggregationFeeds,
		consumerAddr.String()))
	return workflowJobSpec.Job()
}

/* Keith's example from the doc:
name: "secure-mint-ex"
owner: "0xFFFsomeaddr"

triggers:
  - id: "secure_mint@1.0.0" # NEW TRIGGER
    config:
      maxFrequencyMs: 5000
      feedIds:
        - "1020001001"
        - "1020000101"
        - "1020000102"

consensus:
  - id: "offchain_reporting@1.0.0"
    ref: "secure-mint-reattest"
    inputs:
      observations:
        - $(trigger.outputs)
    config:
      report_id: "0003"
      key_id: "evm"
      aggregation_method: "secure_mint" #NEW AGGREGRATION METHOD
      aggregation_config:
        chain_id:
          <CHAIN_ID_FOR_WRITE_TARGET> # NEW Param, to match write target

      encoder: "EVM"
      encoder_config:
        abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports # Existing feed abi

targets:
  - id: "write_ethereum-testnet-sepolia-linea-1@1.0.0"
    inputs:
      signed_report: $(secure-mint-reattest.outputs)
    config:
      address: "0x3524AbD1923402484852E6De6d656965aB37767A"
      deltaStage: "45s"
      schedule: "oneAtATime"
*/
//TODO(gg): are config.feedIds needed? Is config.report_id needed?
const secureMintWorkflow = `
name: "%s"
owner: "0x%s"
triggers:
  - id: "securemint-trigger@1.0.0"
    config:
      maxFrequencyMs: 5000
      feedIds:  
        - "1020001001"  

consensus:
  - id: "offchain_reporting@1.0.0"
    ref: "secure-mint-consensus"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      report_id: "0003"  
      key_id: "evm"	
      aggregation_method: "secure_mint" #NEW AGGREGRATION METHOD
      aggregation_config:
        targetChainSelector:
          %d # CHAIN_ID_FOR_WRITE_TARGET: NEW Param, to match write target
      encoder: "EVM"
      encoder_config:
        abi: "(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports"

targets:
  - id: "write_geth-testnet@1.0.0"
    inputs:
      signed_report: $(secure-mint-consensus.outputs)
    config:
      address: "%s"
      params: ["$(report)"]
      abi: "receive(report bytes)"
      deltaStage: 1s
      schedule: oneAtATime
`

func createSecureMintWorkflowJob(t *testing.T,
	workflowName string,
	workflowOwner string,
	chainSelector int64,
	consumerAddr common.Address) job.Job {
	spec := fmt.Sprintf(secureMintWorkflow, workflowName, workflowOwner, chainSelector, consumerAddr.String())
	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, spec)
	return workflowJobSpec.Job()
}
