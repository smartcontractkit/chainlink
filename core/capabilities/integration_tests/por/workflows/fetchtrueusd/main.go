package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/aggregators"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/targets/chainwriter"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type fetchTrueUSDConfig struct {
	ConsumerAddress         string
	WriteTargetCapabilityID string
	CronSchedule            string
}

type trueUSDResponse struct {
	AccountName string    `json:"accountName"`
	TotalTrust  float64   `json:"totalTrust"`
	TotalToken  float64   `json:"totalToken"`
	Ripcord     bool      `json:"ripcord"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type computeOutput struct {
	TotalTrust uint64
	TotalToken uint64
	Ripcord    bool
	FeedID     [32]byte
	Timestamp  int64
}

type computeConfig struct {
	FeedID string
}

func convertFeedIDtoBytes(feedIDStr string) ([32]byte, error) {
	b, err := hex.DecodeString(feedIDStr[2:])
	if err != nil {
		return [32]byte{}, err
	}

	if len(b) < 32 {
		nb := [32]byte{}
		copy(nb[:], b)
		return nb, err
	}

	return [32]byte(b), nil
}

func BuildWorkflow(configBytes []byte) (*sdk.WorkflowSpecFactory, error) {
	config := &fetchTrueUSDConfig{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	workflow := sdk.NewWorkflowSpecFactory()

	cron := croncap.Config{
		Schedule: config.CronSchedule,
	}.New(workflow)

	compConf := computeConfig{
		FeedID: "0xA1B2C3D4E5F600010203040506070809", // any random bytes16 string to track the feed
	}

	compute := sdk.Compute1WithConfig(
		workflow,
		"compute",
		&sdk.ComputeConfig[computeConfig]{Config: compConf},
		sdk.Compute1Inputs[croncap.Payload]{Arg0: cron},
		func(runtime sdk.Runtime, config computeConfig, outputs croncap.Payload) (computeOutput, error) {
			feedID, err := convertFeedIDtoBytes(config.FeedID)
			if err != nil {
				return computeOutput{}, errors.New("cannot convert feedID to bytes")
			}

			fresp, err := runtime.Fetch(sdk.FetchRequest{
				URL:       "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
				Method:    "GET",
				TimeoutMs: 5000,
			})
			if err != nil {
				return computeOutput{}, fmt.Errorf("not able to fetch API response: %w", err)
			}

			var resp trueUSDResponse
			err = json.Unmarshal(fresp.Body, &resp)
			if err != nil {
				return computeOutput{}, fmt.Errorf("not able to unmarshal response payload %s, err: %w", fresp.Body, err)
			}

			if resp.Ripcord {
				runtime.Emitter().With(
					"feedID", config.FeedID,
				).Emit("ripcord flag set for feed ID " + config.FeedID)
			}

			return computeOutput{
				TotalTrust: uint64(resp.TotalTrust * 100), // 2 decimal places
				TotalToken: uint64(resp.TotalToken * 100), // 2 decimal places
				Ripcord:    resp.Ripcord,                  // 0 decimal places
				FeedID:     feedID,
				Timestamp:  resp.UpdatedAt.Unix(),
			}, nil
		},
	)

	consensusInput := ocr3cap.ReduceConsensusInput[computeOutput]{
		Observation: compute.Value(),
	}

	consensus := ocr3cap.ReduceConsensusConfig[computeOutput]{
		Encoder: ocr3cap.EncoderEVM,
		EncoderConfig: map[string]any{
			"abi": "(bytes32 FeedID, uint32 Timestamp, bytes Bundle)[] Reports",
			"subabi": map[string]string{
				"Reports.Bundle": "uint256 TotalTrust, uint256 TotalToken, bool Ripcord",
			},
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
					InputKey:        "Timestamp",
					OutputKey:       "Timestamp",
					Method:          "median",
					DeviationString: "300",
					DeviationType:   "absolute",
				},
				{
					InputKey:        "TotalTrust",
					OutputKey:       "TotalTrust",
					Method:          "median",
					DeviationString: "1",
					DeviationType:   "percent",
					SubMapField:     true,
				},
				{
					InputKey:        "TotalToken",
					OutputKey:       "TotalToken",
					Method:          "median",
					DeviationString: "1",
					DeviationType:   "percent",
					SubMapField:     true,
				},
				{
					InputKey:    "Ripcord",
					OutputKey:   "Ripcord",
					Method:      "mode",
					SubMapField: true,
				},
			},
			ReportFormat: aggregators.REPORT_FORMAT_ARRAY,
			SubMapKey:    "Bundle",
		},
	}.New(workflow, "consensus", consensusInput)

	targetInput := chainwriter.TargetInput{
		SignedReport: consensus,
	}

	chainwriter.TargetConfig{
		Address:    config.ConsumerAddress, // Sepolia PoR Cache using DF 1.5
		DeltaStage: "15s",
		Schedule:   "oneAtATime",
	}.New(workflow, config.WriteTargetCapabilityID, targetInput)

	return workflow, nil
}

func main() {
	runner := wasm.NewRunner()
	workflow, err := BuildWorkflow(runner.Config())
	if err != nil {
		panic(fmt.Errorf("could not build workflow: %w", err))
	}
	runner.Run(workflow)
}
