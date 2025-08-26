package keystone

import (
	"bytes"
	"fmt"
	"slices"
	"testing"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

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

const secureMintWorkflowTemplate = `
name: "{{.WorkflowName}}"
owner: "0x{{.WorkflowOwner}}"
triggers:
  - id: "securemint-trigger@1.0.0"
    config:
      maxFrequencyMs: 5000

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
        targetChainSelector: "{{.ChainSelector}}" # CHAIN_ID_FOR_WRITE_TARGET: NEW Param, to match write target
        dataID: "{{.DataID}}"
      encoder: "EVM"
      encoder_config:
        abi: "(bytes16 DataID, uint32 Timestamp, uint224 Answer)[] Reports"

targets:
  - id: "write_geth-testnet@1.0.0"
    inputs:
      signed_report: $(secure-mint-consensus.outputs)
    config:
      address: "{{.DFCacheAddress}}"
      params: ["$(report)"]
      abi: "receive(report bytes)"
      deltaStage: 1s
      schedule: oneAtATime
`

type secureMintWorkflowTemplateData struct {
	ChainSelector  uint64
	DataID         string
	DFCacheAddress string
	WorkflowName   string
	WorkflowOwner  string
}

func createSecureMintWorkflowJob(t *testing.T,
	workflowName string,
	workflowOwner string,
	chainSelector uint64,
	dataID string,
	dfCacheAddress common.Address) job.Job {
	tmpl, err := template.New("secureMintWorkflow").Parse(secureMintWorkflowTemplate)
	require.NoError(t, err)

	data := secureMintWorkflowTemplateData{
		ChainSelector:  chainSelector,
		DataID:         dataID,
		DFCacheAddress: dfCacheAddress.String(),
		WorkflowName:   workflowName,
		WorkflowOwner:  workflowOwner,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	spec := buf.String()
	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, spec)
	t.Logf("Generated workflow job spec: %s", spec)
	return workflowJobSpec.Job()
}
