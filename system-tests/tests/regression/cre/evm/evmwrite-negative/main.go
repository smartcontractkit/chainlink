//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmwrite-negative/config"
)

const (
	defaultTimestamp = uint32(1759323120)
	defaultPrice     = int64(100)
	defaultGasLimit  = uint64(400000)
)

type priceOutput struct {
	FeedID    [32]byte
	Timestamp uint32
	Price     *big.Int
}

func main() {
	wasm.NewRunner(parseConfig).Run(RunWriteWorkflow)
}

// parseConfig unmarshals the YAML configuration
func parseConfig(b []byte) (config.Config, error) {
	var wfCfg config.Config
	if err := yaml.Unmarshal(b, &wfCfg); err != nil {
		return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return wfCfg, nil
}

func RunWriteWorkflow(wfCfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onEVMWriteTrigger,
		),
	}, nil
}

func onEVMWriteTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (_ any, _ error) {
	runtime.Logger().Info("onEVMWriteFailTrigger called", "payload", payload)

	evmClient := evm.Client{ChainSelector: wfCfg.ChainSelector}
	runtime.Logger().Info("EVM Client created", "chainSelector", wfCfg.ChainSelector)

	priceOutput, err := createPriceOutput(runtime, wfCfg)
	if err != nil {
		return "", fmt.Errorf("failed to create price output: %w", err)
	}

	report, err := generateReports(runtime, priceOutput)
	if err != nil {
		return "", fmt.Errorf("failed to generate reports: %w", err)
	}

	// cases with corrupt report are skipped as
	// the report generation and consensus on it
	// is a responsibility of the consensus capability
	// and should be tested in the consensus negative tests
	switch wfCfg.FunctionToTest {
	case "WriteReport - invalid receiver":
		return runWriteReportWithInvalidReceiver(evmClient, runtime, wfCfg, report)
	case "WriteReport - corrupt receiver address":
		return runWriteReportWithCorruptReceiverAddress(evmClient, runtime, wfCfg, report)
	case "WriteReport - invalid gas":
		return runWriteReportWithInvalidGas(evmClient, runtime, wfCfg, report)
	default:
		runtime.Logger().Warn("The provided name for function to test in regression EVM Write Workflow did not match any known functions", "functionToTest", wfCfg.FunctionToTest)
		return nil, fmt.Errorf("the provided name for function to test in regression EVM Write Workflow did not match any known functions: %s", wfCfg.FunctionToTest)
	}
}

// runWriteReportWithInvalidReceiver writes a report with an invalid receiver address
// invalid receiver should should return a receipt with reverted transaction
func runWriteReportWithInvalidReceiver(evmClient evm.Client, runtime cre.Runtime, wfCfg config.Config, report *cre.Report) (*evm.WriteReportReply, error) {
	runtime.Logger().Info("Attempting to write report with an invalid receiver", "invalid_receiver", wfCfg.InvalidInput)

	invalidReceiver := common.HexToAddress(wfCfg.InvalidInput).Bytes()
	gasConfig := &evm.GasConfig{GasLimit: defaultGasLimit}
	wrOutput, err := writeReport(runtime, evmClient, invalidReceiver, report, gasConfig)
	if err != nil || wrOutput.ErrorMessage != nil {
		runtime.Logger().Error("got expected error for WriteReport with invalid receiver", "invalid_receiver", invalidReceiver, "error", err)
		return nil, fmt.Errorf("expected error for WriteReport with invalid receiver '%s': %w", invalidReceiver, err)
	}

	runtime.Logger().Info("this is not expected: WriteReport with invalid receiver should return an error", "invalid_receiver", invalidReceiver, "wr_output", wrOutput)
	return wrOutput, nil
}

// runWriteReportWithCorruptReceiverAddress writes a report with a corrupt receiver address
// corrupt receiver should cause transaction to fail
func runWriteReportWithCorruptReceiverAddress(evmClient evm.Client, runtime cre.Runtime, wfCfg config.Config, report *cre.Report) (*evm.WriteReportReply, error) {
	runtime.Logger().Info("Attempting to write report with a corrupt receiver address", "corrupt_receiver", wfCfg.InvalidInput)

	// not using common.HexToAddress to simulate malformed address
	invalidReceiver := []byte(wfCfg.InvalidInput)
	gasConfig := &evm.GasConfig{GasLimit: defaultGasLimit}
	wrOutput, err := writeReport(runtime, evmClient, invalidReceiver, report, gasConfig)
	if err != nil {
		runtime.Logger().Error("got expected error for WriteReport with corrupt receiver address", "corrupt_receiver", invalidReceiver, "error", err)
		return nil, fmt.Errorf("expected error for WriteReport with corrupt receiver address '%s': %w", invalidReceiver, err)
	}

	runtime.Logger().Info("this is not expected: WriteReport with corrupt receiver address should return an error", "corrupt_receiver", invalidReceiver, "wr_output", wrOutput)
	return wrOutput, nil
}

// runWriteReportWithInvalidGas writes a report with an invalid gas
// invalid gas should cause transaction to fail
func runWriteReportWithInvalidGas(evmClient evm.Client, runtime cre.Runtime, wfCfg config.Config, report *cre.Report) (*evm.WriteReportReply, error) {
	runtime.Logger().Info("Attempting to write report with an invalid gas", "invalid_gas", wfCfg.InvalidInput)

	invalidGas, err := strconv.ParseUint(wfCfg.InvalidInput, 10, 64)
	if err != nil {
		runtime.Logger().Error("failed to parse gas limit", "error", err)
		return nil, fmt.Errorf("failed to parse gas limit: %w", err)
	}

	invalidGasConfig := &evm.GasConfig{GasLimit: invalidGas}
	receiver := wfCfg.DataFeedsCache.DataFeedsCacheAddress.Bytes() // valid receiver
	wrOutput, err := writeReport(runtime, evmClient, receiver, report, invalidGasConfig)
	if err != nil {
		runtime.Logger().Error("got expected error for WriteReport with invalid gas", "invalid_gas", invalidGas, "error", err)
		return nil, fmt.Errorf("expected error for WriteReport with invalid gas '%d': %w", invalidGas, err)
	}
	runtime.Logger().Info("this is not expected: WriteReport with invalid gas should return an error", "invalid_gas", invalidGas, "wr_output", wrOutput)
	return wrOutput, nil
}

// createPriceOutput creates a priceOutput struct from the workflow configuration
func createPriceOutput(runtime cre.Runtime, wfCfg config.Config) (priceOutput, error) {
	runtime.Logger().Info("converting feed ID to bytes", "feed_id", wfCfg.FeedID)
	feedID, err := convertFeedIDtoBytes(wfCfg.FeedID)
	if err != nil {
		runtime.Logger().Error("failed to decode feed ID", "error", err)
		return priceOutput{}, fmt.Errorf("failed to decode feed ID: %w", err)
	}

	runtime.Logger().Info("creating priceOutput")
	output := priceOutput{
		FeedID:    feedID,
		Timestamp: defaultTimestamp,
		Price:     big.NewInt(defaultPrice),
	}
	runtime.Logger().Info("priceOutput created")
	return output, nil
}

// generateReports encodes price outputs and generates a report
func generateReports(runtime cre.Runtime, output priceOutput) (*cre.Report, error) {
	outputs := make([]priceOutput, 1)
	outputs[0] = output
	runtime.Logger().Info("Encoding priceOutput...")
	encodedPrice, err := encodeReports(outputs)
	if err != nil {
		runtime.Logger().Error("failed to pack price report", "error", err)
		return nil, fmt.Errorf("failed to pack price report: %w", err)
	}
	runtime.Logger().Info("priceOutput encoded")

	runtime.Logger().Info("Generating report")
	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encodedPrice,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		runtime.Logger().Error("failed to generate report", "error", err)
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}
	runtime.Logger().Info("Report generated successfully")

	return report, nil
}

// writeReport writes the generated report to the blockchain
func writeReport(runtime cre.Runtime, evmClient evm.Client, receiver []byte, report *cre.Report, gasConfig *evm.GasConfig) (*evm.WriteReportReply, error) {
	// skip report and gasConfig validations for testing purposes
	runtime.Logger().Info("Writing report...", "receiver", receiver)
	wrOutput, err := evmClient.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  receiver,
		Report:    report,
		GasConfig: gasConfig,
	}).Await()
	if err != nil {
		runtime.Logger().Error("failed to write report on-chain", "error", err)
		return nil, fmt.Errorf("failed to write report on-chain: %w", err)
	}
	runtime.Logger().Info("Report successfully submitted on-chain", "write_output", wrOutput)

	return wrOutput, nil
}

func encodeReports(reports []priceOutput) ([]byte, error) {
	typ, err := abi.NewType("tuple[]", "",
		[]abi.ArgumentMarshaling{
			{Name: "FeedID", Type: "bytes32"},
			{Name: "Timestamp", Type: "uint32"},
			{Name: "Price", Type: "uint224"},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	args := abi.Arguments{
		{
			Name: "Reports",
			Type: typ,
		},
	}
	return args.Pack(reports)
}

// convertFeedIDtoBytes converts a hex string feed ID to a 32-byte array
func convertFeedIDtoBytes(feedID string) ([32]byte, error) {
	if feedID == "" {
		return [32]byte{}, fmt.Errorf("feedID string is empty")
	}

	// Remove hex prefix if present
	hexStr := feedID
	hexPrefix := "0x"
	if len(feedID) >= 2 && feedID[:2] == hexPrefix {
		hexStr = feedID[2:]
	}

	if len(hexStr) == 0 {
		return [32]byte{}, fmt.Errorf("feedID string contains no hex data: %q", feedID)
	}

	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode hex string %q: %w", feedID, err)
	}

	var result [32]byte
	if len(b) > 32 {
		// Truncate if too long
		copy(result[:], b[:32])
	} else {
		// Pad with zeros if too short
		copy(result[:], b)
	}

	return result, nil
}
