package keystone

import (
	"fmt"
	"slices"
	"strings"
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
	var triggerFeedIDsSb60 strings.Builder
	for _, feedID := range feedIDs {
		fmt.Fprintf(&triggerFeedIDsSb60, "        - \"%s\"\n", feedID)
	}
	triggerFeedIDs += triggerFeedIDsSb60.String()

	aggregationFeeds := ""
	var aggregationFeedsSb65 strings.Builder
	for _, feedID := range feedIDs {
		fmt.Fprintf(&aggregationFeedsSb65, "          \"%s\":\n            deviation: \"0.001\"\n            heartbeat: 3600\n", feedID)
	}
	aggregationFeeds += aggregationFeedsSb65.String()

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
	// keys of the map are stream IDs
	streamIDs := make([]uint32, 0, len(streamIDremapped))
	for streamID := range streamIDremapped {
		streamIDs = append(streamIDs, streamID)
	}
	slices.Sort(streamIDs)
	var triggerFeedIDs strings.Builder
	for _, streamID := range streamIDs {
		fmt.Fprintf(&triggerFeedIDs, "        - \"%d\"\n", streamID)
	}

	var aggregationFeeds strings.Builder
	for _, streamID := range streamIDs {
		fmt.Fprintf(&aggregationFeeds, "          \"%d\":\n            deviation: \"0.001\"\n            heartbeat: 3600\n            remappedID: \"%s\"\n", streamID, streamIDremapped[streamID])
	}

	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, fmt.Sprintf(lloStreamsWorkflow, workflowName, workflowOwner, &triggerFeedIDs, &aggregationFeeds,
		consumerAddr.String()))
	return workflowJobSpec.Job()
}
