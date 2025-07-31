package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/aggregators"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/targets/chainwriter"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm"
	types "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/web-trigger-based/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
)

func main() {
	runner := wasm.NewRunner()
	workflow, err := valueTriggeredWorkflow(runner)
	if err != nil {
		panic(err)
	}
	runner.Run(workflow)
}

func valueTriggeredWorkflow(r sdk.Runner) (*sdk.WorkflowSpecFactory, error) {
	var workflowConfig types.WorkflowConfig
	if err := yaml.Unmarshal(r.Config(), &workflowConfig); err != nil {
		return nil, err
	}

	w := sdk.NewWorkflowSpecFactory()

	// Web API Trigger: define allowed sender, rate limits, required parameters
	trigger := webapicap.TriggerConfig{
		AllowedSenders: []string{workflowConfig.AllowedTriggerSender},
		AllowedTopics:  []string{workflowConfig.AllowedTriggerTopic},
		RateLimiter: webapicap.RateLimiterConfig{
			GlobalBurst:    1000,
			GlobalRPS:      1000,
			PerSenderBurst: 1000,
			PerSenderRPS:   1000,
		},
		RequiredParams: []string{"value"},
	}.New(w)

	// Compute: get value from event
	contractValue := sdk.Compute1(
		w,
		"getValue",
		sdk.Compute1Inputs[webapicap.TriggerRequestPayloadParams]{Arg0: trigger.Params().(sdk.CapDefinition[webapicap.TriggerRequestPayloadParams])},
		func(SDK sdk.Runtime, o webapicap.TriggerRequestPayloadParams) (ValueOutput, error) {
			if len(o) == 0 {
				return ValueOutput{}, fmt.Errorf("no data found in event")
			}

			maybeValue, ok := o["value"]
			if !ok {
				return ValueOutput{}, fmt.Errorf("value with name 'value' not found in payload")
			}

			valueStr, ok := maybeValue.(string)
			if !ok {
				return ValueOutput{}, fmt.Errorf("value is not a string, but %T", maybeValue)
			}

			valueBigInt := new(big.Int)
			valueBigInt, ok = valueBigInt.SetString(valueStr, 10)
			if !ok {
				return ValueOutput{}, fmt.Errorf("failed to convert value %s to big.Int", valueStr)
			}

			// Convert the FeedID string to a byte array
			feedIDBytes, err := convertFeedIDtoBytes(workflowConfig.FeedID)
			if err != nil {
				return ValueOutput{}, fmt.Errorf("failed to convert FeedID to bytes: %w", err)
			}

			return ValueOutput{
				Price:     valueBigInt,
				Timestamp: time.Now().Unix(),
				FeedID:    feedIDBytes,
			}, nil
		},
	)

	// Consensus: all observations are aggregated, timestamps can be maximum 30 seconds apart
	// median of all values is used as the value price
	consensusInput := ocr3cap.ReduceConsensusInput[ValueOutput]{
		Observation: contractValue.Value(),
	}

	consensus := ocr3cap.ReduceConsensusConfig[ValueOutput]{
		Encoder: ocr3cap.EncoderEVM,
		EncoderConfig: map[string]any{
			"abi": "(bytes32 FeedID, uint32 Timestamp, uint224 Price)[] Reports",
		},
		ReportID: "0001",
		KeyID:    "evm",
		AggregationConfig: aggregators.ReduceAggConfig{
			Fields: []aggregators.AggregationField{
				{
					InputKey:  "FeedID",
					OutputKey: "FeedID",
					Method:    "mode",
				},
				{
					InputKey:      "Price",
					OutputKey:     "Price",
					Method:        "median",
					DeviationType: "any",
				},
				{
					InputKey:        "Timestamp",
					OutputKey:       "Timestamp",
					Method:          "median",
					DeviationString: "30",
					DeviationType:   "absolute",
				},
			},
			ReportFormat: aggregators.REPORT_FORMAT_ARRAY,
		},
	}.New(w, "consensus", consensusInput)

	// Write: write the median price to the Data Feeds Cache contract
	targetInput := chainwriter.TargetInput{
		SignedReport: consensus,
	}

	writeTargetName := "write_geth-testnet@1.0.0"
	if workflowConfig.WriteTargetName != "" {
		writeTargetName = workflowConfig.WriteTargetName
	}

	chainwriter.TargetConfig{
		CreStepTimeout: 40,                                   // 10 seconds
		Address:        workflowConfig.DataFeedsCacheAddress, // Data Feeds Cache contract address
		DeltaStage:     "15s",
		Schedule:       "oneAtATime",
	}.New(w, writeTargetName, targetInput)

	return w, nil
}

func Ptr[T any](v T) *T {
	return &v
}

func convertFeedIDtoBytes(feedIDStr string) ([32]byte, error) {
	b, err := hex.DecodeString(feedIDStr[2:])
	if err != nil {
		return [32]byte{}, err
	}

	if len(b) < 32 {
		nb := [32]byte{}
		copy(nb[:], b[:])
		return nb, err
	}

	return [32]byte(b), nil
}

type ValueOutput struct {
	Price     *big.Int
	Timestamp int64
	FeedID    [32]byte
}
