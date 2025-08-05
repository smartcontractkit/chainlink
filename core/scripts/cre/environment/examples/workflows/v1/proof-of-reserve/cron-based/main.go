package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/aggregators"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/targets/chainwriter"
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm"
	types "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
)

func main() {
	runner := wasm.NewRunner()
	workflow := BuildWorkflow(runner)
	runner.Run(workflow)
}

func BuildWorkflow(runner *wasm.Runner) *sdk.WorkflowSpecFactory {
	workflow := sdk.NewWorkflowSpecFactory()

	cron := croncap.Config{
		Schedule: "*/30 * * * * *", // Every 30 seconds
	}.New(workflow)

	var workflowConfig types.WorkflowConfig
	err := yaml.Unmarshal(runner.Config(), &workflowConfig)
	if err != nil {
		runner.ExitWithError(errors.New("cannot unmarshal config : %w"))
	}

	if workflowConfig.FeedID == "" {
		runner.ExitWithError(fmt.Errorf("feedID is empty in the config: %+v", workflowConfig))
	}

	computeConfig := types.ComputeConfig{
		FeedID:                workflowConfig.FeedID,
		URL:                   workflowConfig.URL,
		DataFeedsCacheAddress: workflowConfig.DataFeedsCacheAddress,
		WriteTargetName:       workflowConfig.WriteTargetName,
	}

	if workflowConfig.AuthKeySecretName != "" {
		// Secrets are only resolved by Custom Compute, if they are passed as config fields with type `sdk.SecretValue`
		// If we tried to call `sdk.Secret(workflowConfig.AuthKeySecretName)` directly inside the `Compute1WithConfig` function,
		// it would not be resolved and would be passed as a string to the compute function.
		computeConfig.AuthKey = sdk.Secret(workflowConfig.AuthKeySecretName)
	}

	compute := sdk.Compute1WithConfig(
		workflow,
		"compute",
		&sdk.ComputeConfig[types.ComputeConfig]{Config: computeConfig},
		sdk.Compute1Inputs[croncap.Payload]{Arg0: cron},
		func(runtime sdk.Runtime, config types.ComputeConfig, outputs croncap.Payload) (computeOutput, error) {
			feedID, err := convertFeedIDtoBytes(config.FeedID)
			if err != nil {
				return computeOutput{}, fmt.Errorf("cannot convert feedID to bytes : %w : %b", err, feedID)
			}

			fetchRequest := sdk.FetchRequest{
				URL:       config.URL + "?feedID=" + config.FeedID,
				Method:    "GET",
				TimeoutMs: 5000,
			}

			if string(config.AuthKey) != "" {
				fetchRequest.Headers = map[string]string{
					"Authorization": string(config.AuthKey),
				}
			}

			fresp, err := runtime.Fetch(fetchRequest)
			if err != nil {
				return computeOutput{}, err
			}

			var resp trueUSDResponse
			err = json.Unmarshal(fresp.Body, &resp)
			if err != nil {
				return computeOutput{}, err
			}

			runtime.Emitter().With(
				"feedID", config.FeedID,
			).Emit(fmt.Sprintf("TrueUSD price found: %.2f", resp.TotalTrust))

			if resp.Ripcord {
				runtime.Emitter().With(
					"feedID", config.FeedID,
				).Emit(fmt.Sprintf("ripcord flag set for feed ID %s", config.FeedID))
				return computeOutput{}, sdk.BreakErr
			}

			return computeOutput{
				Price:     int(resp.TotalTrust * 100),
				FeedID:    feedID, // TrueUSD
				Timestamp: resp.UpdatedAt.Unix(),
			}, nil
		},
	)

	consensusInput := ocr3cap.ReduceConsensusInput[computeOutput]{
		Observation: compute.Value(),
	}

	consensus := ocr3cap.ReduceConsensusConfig[computeOutput]{
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
	}.New(workflow, "consensus", consensusInput)

	targetInput := chainwriter.TargetInput{
		SignedReport: consensus,
	}

	writeTargetName := "write_geth-testnet@1.0.0"
	if workflowConfig.WriteTargetName != "" {
		writeTargetName = workflowConfig.WriteTargetName
	}

	chainwriter.TargetConfig{
		Address:    workflowConfig.DataFeedsCacheAddress, // KeystoneConsumer contract address
		DeltaStage: "15s",
		Schedule:   "oneAtATime",
	}.New(workflow, writeTargetName, targetInput)

	return workflow
}

type trueUSDResponse struct {
	AccountName string    `json:"accountName"`
	TotalTrust  float64   `json:"totalTrust"`
	Ripcord     bool      `json:"ripcord"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type computeOutput struct {
	Price     int
	FeedID    [32]byte
	Timestamp int64
}

func convertFeedIDtoBytes(feedIDStr string) ([32]byte, error) {
	if feedIDStr == "" {
		return [32]byte{}, fmt.Errorf("feedID string is empty")
	}

	if len(feedIDStr) < 2 {
		return [32]byte{}, fmt.Errorf("feedID string too short: %q", feedIDStr)
	}

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
