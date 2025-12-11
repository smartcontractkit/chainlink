//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

type WorkflowConfig struct {
}

func RunSolWriteWorkflow(cfg WorkflowConfig, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[WorkflowConfig], error) {
	return cre.Workflow[WorkflowConfig]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}), // every 30 seconds
			onTrigger,
		),
	}, nil

}
func onTrigger(config WorkflowConfig, runtime cre.Runtime, payload *cron.Payload) (string, error) {
	runtime.Logger().Info("Solana Write workflow started", "payload", payload)
	solClient := solana.Client{ChainSelector: chain_selectors.SOLANA_DEVNET.Selector}
	runtime.Logger().Info("Got Solana client", "chainSelector", solClient.ChainSelector)
	encodedPayload := []byte{1, 2, 3}
	// 1. encode report

	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encodedPayload,
		EncoderName:    "evm", // add borsh encoder later
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	// 2. execute Write
	output, err := solClient.WriteReport(runtime, &solana.WriteCreReportRequest{
		Report: report,
	}).Await()
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("[logger] failed to write report on-chain: %v", err))
		return "", fmt.Errorf("failed to write report on solana chain: %w", err)
	}

	runtime.Logger().With().Info("Submitted report on-chain")

	var message = "PoR Workflow successfully completed"
	if output.ErrorMessage != nil {
		message = *output.ErrorMessage
	}

	return message, nil
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (WorkflowConfig, error) {
		cfg := WorkflowConfig{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSolWriteWorkflow)
}
