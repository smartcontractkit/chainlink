//go:build wasip1

package main

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/targets/chainwriter"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type FeedConfig struct {
	RemappedID string `json:"remapped_id"`
}

type ConfigFile struct {
	Feeds map[string]FeedConfig `json:"feeds"`
}

func BuildDataFeedsWorkflow(configData []byte) *sdk.WorkflowSpecFactory {
	workflow := sdk.NewWorkflowSpecFactory()

	var config ConfigFile
	if err := json.Unmarshal(configData, &config); err != nil {
		panic(err)
	}

	var feedIds []streams.FeedId
	for feedId := range config.Feeds {
		feedIds = append(feedIds, streams.FeedId(feedId))
	}

	triggerCfg := streams.TriggerConfig{
		FeedIds:        feedIds,
		MaxFrequencyMs: 5000,
	}
	trigger := triggerCfg.New(workflow)

	ocr3Config := ocr3cap.DataFeedsConsensusConfig{
		AggregationMethod: "data_feeds",
		Encoder:           "EVM",
		EncoderConfig: map[string]any{
			"abi": "(bytes32 RemappedID, bytes RawReport)[] Reports",
		},
		ReportId: "0001",
		KeyId:    "evm",
		AggregationConfig: ocr3cap.DataFeedsConsensusConfigAggregationConfig{
			AllowedPartialStaleness: "0.5",
		},
	}

	ocr3Config.AggregationConfig.Feeds = ocr3cap.DataFeedsConsensusConfigAggregationConfigFeeds{}

	// Configure feeds from JSON
	for feedId, feedConfig := range config.Feeds {
		ocr3Config.AggregationConfig.Feeds[feedId] = ocr3cap.FeedValue{
			Deviation:  "0.001",
			Heartbeat:  1800,
			RemappedID: &feedConfig.RemappedID,
		}
	}

	consensus := ocr3Config.New(workflow, "aptos_feeds", ocr3cap.DataFeedsConsensusInput{
		Observations: sdk.ListOf(trigger),
	})

	chainWriterConfig := chainwriter.TargetConfig{
		Address:    "0xf1099f135ddddad1c065203431be328a408b0ca452ada70374ce26bd2b32fdd3",
		DeltaStage: "10s",
		Schedule:   "oneAtATime",
	}
	chainWriterConfig.New(workflow, "write_aptos-testnet@1.0.0", chainwriter.TargetInput{SignedReport: consensus})

	return workflow
}

func main() {
	runner := wasm.NewRunner()
	workflow := BuildDataFeedsWorkflow(runner.Config())
	runner.Run(workflow)
}
