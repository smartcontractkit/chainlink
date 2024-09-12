package integration_tests

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

const workflowTemplateStreams = `
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

func addWorkflowJobStreams(t *testing.T, app *cltest.TestApplication,
	workflowName string,
	workflowOwner string,
	feedIDs []string,
	consumerAddr common.Address,
	deltaStage string,
	schedule string) {
	triggerFeedIDs := ""
	for _, feedID := range feedIDs {
		triggerFeedIDs += fmt.Sprintf("        - \"%s\"\n", feedID)
	}

	aggregationFeeds := ""
	for _, feedID := range feedIDs {
		aggregationFeeds += fmt.Sprintf("          \"%s\":\n            deviation: \"0.001\"\n            heartbeat: 3600\n", feedID)
	}

	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, fmt.Sprintf(workflowTemplateStreams, workflowName, workflowOwner, triggerFeedIDs, aggregationFeeds,
		consumerAddr.String(), deltaStage, schedule))
	job := workflowJobSpec.Job()

	err := app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
}

const workflowTemplatePoR = `
name: "%s"
owner: "0x%s"
triggers:
  - id: "cron-trigger@1.0.0"
    ref: "trigger"
    config:
      schedule:
        \"%s\"

# Compute

# Consensus

targets:
  - id: "kv-store-target@1.0.0"
    inputs:
      signed_report: "$(trigger.outputs)"
    config:
      deltaStage: %s
      schedule: %s
`

// targets:
//   - id: "log-target@1.0.0"
//     ref: "target"
//     inputs:
//       signed_report: $(trigger.outputs)
//     config:
//       deltaStage: %s
//       schedule: %s

func addWorkflowJobPoR(
	t *testing.T,
	app *cltest.TestApplication,
	workflowName string,
	workflowOwner string,
	cronSchedule string,
	consumerAddr common.Address,
	deltaStage string,
	schedule string,
) {
	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(
		t,
		fmt.Sprintf(
			workflowTemplatePoR,
			workflowName,
			workflowOwner,
			cronSchedule,
			deltaStage,
			schedule,
		),
	)

	job := workflowJobSpec.Job()

	err := app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
}

const standardCapabilityTemplateCron = `
type = "standardcapabilities"
schemaVersion = 1
name = "cron-capabilities"
command="./cron"
config=""
`

func addStandardCapabilityCron(
	t *testing.T,
	app *cltest.TestApplication,
) {
	// Add Cron
	job, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(standardCapabilityTemplateCron)
	require.NoError(t, err)

	err = app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
}

const standardCapabilityTemplateKeyValue = `
type = "standardcapabilities"
schemaVersion = 1
name = "kvstore-capabilities"
command="./kvstore"
config=""
`

func addStandardCapabilityKV(
	t *testing.T,
	app *cltest.TestApplication,
) {
	// Add KVStore
	job, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(standardCapabilityTemplateKeyValue)
	require.NoError(t, err)

	err = app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
}
