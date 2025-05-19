package src

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

func deployKeystoneWorkflowsTo(nodeSet NodeSet, reg kcr.CapabilitiesRegistryInterface) {
	fmt.Println("Deploying Keystone workflow jobs")
	caps, err := reg.GetCapabilities(&bind.CallOpts{})
	PanicErr(err)

	streams := NewStreamsTriggerV1Capability()
	ocr3 := NewOCR3V1ConsensusCapability()
	testnetWrite := NewEthereumGethTestnetV1WriteCapability()

	capSet := NewCapabilitySet(streams, ocr3, testnetWrite)
	expectedHashedCIDs := capSet.HashedIDs(reg)

	// Check that the capabilities are registered
	for _, c := range caps {
		found := false
		for _, expected := range expectedHashedCIDs {
			if c.HashedId == expected {
				found = true
				break
			}
		}

		if !found {
			panic(fmt.Sprintf("Capability %s not found in registry", c.HashedId))
		}
	}

	feedIDs := []string{}
	for _, feed := range feeds {
		feedIDs = append(feedIDs, fmt.Sprintf("0x%x", feed.id))
	}
	workflowConfig := WorkflowJobSpecConfig{
		JobSpecName:          "keystone_workflow",
		WorkflowOwnerAddress: "0x1234567890abcdef1234567890abcdef12345678",
		FeedIDs:              feedIDs,
		TargetID:             testnetWrite.GetID(),
		ConsensusID:          ocr3.GetID(),
		TriggerID:            streams.GetID(),
		TargetAddress:        "0x1234567890abcdef1234567890abcdef12345678",
	}
	jobSpecStr := createKeystoneWorkflowJob(workflowConfig)
	for _, n := range nodeSet.Nodes[1:] { // skip the bootstrap node
		api := newNodeAPI(n)
		upsertJob(api, workflowConfig.JobSpecName, jobSpecStr)
	}
}

type WorkflowJobSpecConfig struct {
	JobSpecName          string
	WorkflowOwnerAddress string
	FeedIDs              []string
	TargetID             string
	ConsensusID          string
	TriggerID            string
	TargetAddress        string
}

func createKeystoneWorkflowJob(workflowConfig WorkflowJobSpecConfig) string {
	const keystoneWorkflowTemplate = `
type = "workflow"
schemaVersion = 1
name = "{{ .JobSpecName }}"
workflow = """
name: "ccip_kiab1" 
owner: '{{ .WorkflowOwnerAddress }}'
triggers:
 - id: streams-trigger@1.1.0
   config:
     maxFrequencyMs: 10000
     feedIds:
{{- range .FeedIDs }}
       - '{{ . }}'
{{- end }}

consensus:
 - id: offchain_reporting@1.0.0
   ref: ccip_feeds
   inputs:
     observations:
       - $(trigger.outputs)
   config:
     report_id: '0001'
     key_id: 'evm'
     aggregation_method: data_feeds
     aggregation_config:
       feeds:
{{- range .FeedIDs }}
        '{{ . }}':
          deviation: '0.05'
          heartbeat: 1800
{{- end }}
     encoder: EVM
     encoder_config:
       abi: "(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports"
       abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports

targets:
 - id: {{ .TargetID }} 
   inputs:
     signed_report: $(ccip_feeds.outputs)
   config:
     address: '{{ .TargetAddress }}'
     deltaStage: 5s
     schedule: oneAtATime

"""
workflowOwner = "{{ .WorkflowOwnerAddress }}"
`

	tmpl, err := template.New("workflow").Parse(keystoneWorkflowTemplate)

	if err != nil {
		panic(err)
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, workflowConfig)
	if err != nil {
		panic(err)
	}

	return renderedTemplate.String()
}
