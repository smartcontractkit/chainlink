//go:build wasip1

// Command cronconsensus is a workflow that exercises cron and consensus together: on every cron tick each node
// observes its own random number, consensus agrees on their median, and consensus then signs a report of it.
package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type Config struct{}

type Result struct {
	Median        int64
	ReportCreated bool
}

func InitWorkflow(_ *Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}), onCronTrigger),
	}, nil
}

func onCronTrigger(config *Config, runtime cre.Runtime, _ *cron.Payload) (*Result, error) {
	median, err := cre.RunInNodeMode(config, runtime, observe, cre.ConsensusMedianAggregation[int64]()).Await()
	if err != nil {
		return nil, fmt.Errorf("median consensus failed: %w", err)
	}

	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: fmt.Appendf(nil, "median=%d", median),
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("report consensus failed: %w", err)
	}

	runtime.Logger().Info("Consensus reached", "median", median)
	return &Result{Median: median, ReportCreated: report != nil}, nil
}

// observe is this node's own random number in [1,100], so that the nodes disagree and the median is real consensus.
func observe(_ *Config, nodeRuntime cre.NodeRuntime) (int64, error) {
	r, err := nodeRuntime.Rand()
	if err != nil {
		return 0, err
	}

	return r.Int63n(100) + 1, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
