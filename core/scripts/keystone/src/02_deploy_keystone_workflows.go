package src

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

type deployKeystoneWorkflows struct{}

func NewDeployKeystoneWorkflowsCommand() *deployKeystoneWorkflows {
	return &deployKeystoneWorkflows{}
}

func (g *deployKeystoneWorkflows) Name() string {
	return "deploy-keystone-workflows"
}

func (g *deployKeystoneWorkflows) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 1337, "chain id")
	nodeSetsPath := fs.String("nodesets", defaultNodeSetsPath, "Custom node sets location")
	nodeSetSize := fs.Int("nodesetsize", 5, "number of nodes in a nodeset")

	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	artefactsDir := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")
	err := fs.Parse(args)
	if err != nil ||
		*ethUrl == "" || ethUrl == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}
	workflowNodeSet := downloadNodeSets(*chainID, *nodeSetsPath, *nodeSetSize).Workflow
	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	env := helpers.SetupEnv(false)

	o := LoadOnchainMeta(*artefactsDir, env)
	deployKeystoneWorkflowsTo(workflowNodeSet, o.CapabilitiesRegistry, *chainID)
}

func deployKeystoneWorkflowsTo(nodeSet NodeSet, reg kcr.CapabilitiesRegistryInterface, chainID int64) {
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

	feedIds := []string{}
	for _, feed := range feeds {
		feedIds = append(feedIds, fmt.Sprintf("0x%x", feed.id))
	}
	workflowConfig := WorkflowJobSpecConfig{
		JobSpecName:          "keystone_workflow",
		WorkflowOwnerAddress: "0x1234567890abcdef1234567890abcdef12345678",
		FeedIDs:              feedIds,
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
 - id: streams-trigger@1.0.0
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
