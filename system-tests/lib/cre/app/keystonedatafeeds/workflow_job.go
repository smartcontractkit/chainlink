package keystonedatafeeds

import (
	"bytes"
	"text/template"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
)

type FeedConfig struct {
	FeedIDsIndex int32  `json:"feedIDsIndex"`
	Deviation    string `json:"deviation"`
	Heartbeat    int32  `json:"heartbeat"`
	RemappedID   string `json:"remappedID"`
}

// TODO shouldn't consumer address be configurable?
func WorkflowsJob(nodeID string, workflowName string, feeds []FeedConfig) *jobv1.ProposeJobRequest {
	const workflowTemplateLoad = `
 type = "workflow"
 schemaVersion = 1
 name = "{{ .WorkflowName }}"
 externalJobID = "{{ .JobID }}"
 workflow = """
 name: "{{ .WorkflowName }}"
 owner: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
 triggers:
  - id: streams-trigger@2.0.0
    config:
      feedIds:
 {{- range .Feeds }}
        - '{{ .FeedIDsIndex }}'
 {{- end }}
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
{{- range .Feeds }}
           "{{ .FeedIDsIndex }}":
             deviation: "{{ .Deviation }}"
             heartbeat: {{ .Heartbeat }}
             remappedID: {{ .RemappedID }}
{{- end }}
       encoder: "EVM"
       encoder_config:
         abi: "(bytes32 RemappedID, uint224 Price, uint32 Timestamp)[] Reports"
 targets:
   - id: write_ethereum_mock@1.0.0
     inputs:
       signed_report: "$(evm_median.outputs)"
     config:
       address: "0xEB739A9641938934D21A325A0C6b26126D48926A"
       params: ["$(report)"]
       abi: "receive(report bytes)"
       deltaStage: 2s
       schedule: allAtOnce
 """
 `

	tmpl, err := template.New("workflow").Parse(workflowTemplateLoad)

	if err != nil {
		panic(err)
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, map[string]interface{}{
		"WorkflowName": workflowName,
		"Feeds":        feeds,
		"JobID":        uuid.NewString(),
	})
	if err != nil {
		panic(err)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String()}
}
