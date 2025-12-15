//go:build wasip1

package main

import (
	"crypto/sha256"
	"fmt"
	"log/slog"

	solanago "github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

func RunSolWriteWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}), // every 30 seconds
			onTrigger,
		),
	}, nil

}
func onTrigger(config config.Config, runtime cre.Runtime, payload *cron.Payload) (string, error) {
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

	remainings, err := deriveRemaining(config, report)
	if err != nil {
		return "", fmt.Errorf("failed to derive remaining accounts: %w", err)
	}
	// 2. execute WriteReport
	output, err := solClient.WriteReport(runtime, &solana.WriteCreReportRequest{
		Receiver:          config.Receiver.Bytes(),
		Report:            report,
		RemainingAccounts: remainings,
	}).Await()
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("[logger] failed to write report on-chain: %v", err))
		return "", fmt.Errorf("failed to write report on solana chain: %w", err)
	}

	runtime.Logger().With().Info("Submitted report on-chain")

	var message = "Solana Workflow successfully completed"
	if output.ErrorMessage != nil {
		message = *output.ErrorMessage
	}

	return message, nil
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSolWriteWorkflow)
}

func deriveRemaining(config config.Config, report *cre.Report) ([]*solana.AccountMeta, error) {
	executionState, err := deriveExecutionState(config, report.X_GeneratedCodeOnly_Unwrap().RawReport)
	if err != nil {
		return nil, fmt.Errorf("failed to derive execution state: %w", err)
	}
	return []*solana.AccountMeta{
		{PublicKey: config.ForwarderState[:]},
		{PublicKey: executionState[:], IsWritable: true},
	}, nil
}

var (
	reportIDOffset    = 107
	reporIDSize       = 2
	executionIDOffset = 1
	executionIDSize   = 32
)

func deriveExecutionState(config config.Config, rawReport []byte) (solanago.PublicKey, error) {
	transmissionID, err := extractTransmissionID(config.Receiver, rawReport)
	if err != nil {
		return solanago.PublicKey{}, err
	}

	seeds := [][]byte{
		[]byte("execution_state"),
		config.ForwarderState.Bytes(),
		transmissionID[:],
	}

	ret, _, err := solanago.FindProgramAddress(seeds, config.ForwarderProgramID)

	return ret, err
}

func extractTransmissionID(receiver solanago.PublicKey, rawReport []byte) ([32]byte, error) {
	var data []byte

	if len(rawReport) <= reportIDOffset+reporIDSize {
		return [32]byte{}, fmt.Errorf("invalid len of raw report: %d", len(rawReport))
	}

	// 1. add receiver
	data = append(data, receiver.Bytes()...)

	// 2. add executionID
	executionID := rawReport[executionIDOffset : executionIDOffset+executionIDSize]
	data = append(data, executionID...)

	// 3. add reportID
	reportID := rawReport[reportIDOffset : reportIDOffset+reporIDSize]
	data = append(data, reportID...)

	return sha256.Sum256(data), nil
}
