package keystone

import (
	"fmt"
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
