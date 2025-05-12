//go:build wasip1

package main

import (
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/targets/chainwriter"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/cli/cmd/testdata/fixtures/capabilities/basictrigger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

func BuildWorkflow(config []byte) *sdk.WorkflowSpecFactory {
	workflow := sdk.NewWorkflowSpecFactory()

	triggerCfg := basictrigger.TriggerConfig{Name: "trigger", Number: 100}
	trigger := triggerCfg.New(workflow)

	// Config elided because it's not relevant for the test.
	consensusInput := ocr3cap.ReduceConsensusInput[basictrigger.TriggerOutputs]{
		Observation: trigger,
	}
	consensus := ocr3cap.ReduceConsensusConfig[basictrigger.TriggerOutputs]{}.New(workflow, "consensus", consensusInput)

	input := chainwriter.TargetInput{
		SignedReport: consensus,
	}
	chainwriter.TargetConfig{
		Address: "0x0",
		CreStepTimeout: 0,
		DeltaStage: "30s",
		Schedule: "oneAtATime",
	}.New(workflow, "write_ethereum-testnet-sepolia@1.0.0", input)

	return workflow
}

func main() {
	runner := wasm.NewRunner()
	workflow := BuildWorkflow(runner.Config())
	runner.Run(workflow)
}
