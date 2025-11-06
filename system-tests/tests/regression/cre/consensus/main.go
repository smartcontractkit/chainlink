//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/consensus/config"
)

const (
	defaultPrice     = int64(100)
	defaultTimestamp = uint32(1759491269)
	maxRandomPrice   = int64(1000000)
	randomTimeOffset = int64(7200) // ±1 hour in seconds
)

type priceOutput struct {
	FeedID    [32]byte
	Timestamp uint32
	Price     *big.Int
}

func main() {
	wasm.NewRunner(parseConfig).Run(RunConsensusNegativeWorkflow)
}

// parseConfig unmarshals the YAML configuration
func parseConfig(b []byte) (config.Config, error) {
	var wfCfg config.Config
	if err := yaml.Unmarshal(b, &wfCfg); err != nil {
		return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return wfCfg, nil
}

func RunConsensusNegativeWorkflow(wfCfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onConsensusNegativeTrigger,
		),
	}, nil
}

func onConsensusNegativeTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (_ any, _ error) {
	runtime.Logger().Info("onConsensusNegativeTrigger called", "payload", payload)

	switch wfCfg.CaseToTrigger {
	case "Consensus - random timestamps":
		return runConsensusGenerateReportWithRandomTimestamps(runtime, wfCfg)
	case "Consensus - inconsistent feedIDs":
		return runConsensusGenerateReportWithInconsistentFeedIDs(runtime, wfCfg)
	case "Consensus - inconsistent prices":
		return runConsensusGenerateReportWithInconsistentPrices(runtime, wfCfg)
	default:
		runtime.Logger().Warn("The provided name for function to test in regression Consensus Workflow did not match any known functions", "functionToTest", wfCfg.CaseToTrigger)
		return nil, fmt.Errorf("the provided name for function to test in regression Consensus Workflow did not match any known functions: %s", wfCfg.CaseToTrigger)
	}
}

// runConsensusGenerateReportWithRandomTimestamps writes a report with different timestamps
// different timestamps should cause transaction to fail
func runConsensusGenerateReportWithRandomTimestamps(runtime cre.Runtime, wfCfg config.Config) (*cre.Report, error) {
	runtime.Logger().Info("Attempting to write report with different timestamps")

	priceOutputWithRandomTimestamp, err := createPriceOutputWithRandomTimestamp(runtime, wfCfg)
	if err != nil {
		runtime.Logger().Error("failed to create price output with random timestamps", "error", err)
		return nil, fmt.Errorf("failed to create price output with random timestamps: %w", err)
	}

	reportWithDifferentTimestamps, err := generateReports(runtime, priceOutputWithRandomTimestamp)
	if err != nil {
		runtime.Logger().Error("got expected error for WriteReport with random timestamps", "error", err)
		return nil, fmt.Errorf("expected error for WriteReport with random timestamps: %w", err)
	}

	runtime.Logger().Info("this is not expected: WriteReport with different timestamps should return an error", "generated_report", reportWithDifferentTimestamps)
	return reportWithDifferentTimestamps, nil
}

// runConsensusGenerateReportWithInconsistentFeedIDs writes a report with inconsistent feedIDs
// inconsistent feedIDs should cause consensus to fail
func runConsensusGenerateReportWithInconsistentFeedIDs(runtime cre.Runtime, wfCfg config.Config) (*cre.Report, error) {
	runtime.Logger().Info("Attempting to generate report with inconsistent feedIDs")

	priceOutputWithRandomFeedID, err := createPriceOutputWithRandomFeedID(runtime, wfCfg)
	if err != nil {
		runtime.Logger().Error("failed to create price output with random feedID", "error", err)
		return nil, fmt.Errorf("failed to create price output with random feedID: %w", err)
	}

	reportWithInconsistentFeedIDs, err := generateReports(runtime, priceOutputWithRandomFeedID)
	if err != nil {
		runtime.Logger().Error("got expected error for GenerateReport with inconsistent feedIDs", "error", err)
		return nil, fmt.Errorf("expected error for GenerateReport with inconsistent feedIDs: %w", err)
	}

	runtime.Logger().Info("this is not expected: GenerateReport with inconsistent feedIDs should return an error", "generated_report", reportWithInconsistentFeedIDs)
	return reportWithInconsistentFeedIDs, nil
}

// runConsensusGenerateReportWithInconsistentPrices writes a report with inconsistent prices
// inconsistent prices should cause consensus to fail
func runConsensusGenerateReportWithInconsistentPrices(runtime cre.Runtime, wfCfg config.Config) (*cre.Report, error) {
	runtime.Logger().Info("Attempting to generate report with inconsistent prices")

	priceOutputWithRandomPrice, err := createPriceOutputWithRandomPrice(runtime, wfCfg)
	if err != nil {
		runtime.Logger().Error("failed to create price output with random price", "error", err)
		return nil, fmt.Errorf("failed to create price output with random price: %w", err)
	}

	reportWithInconsistentPrices, err := generateReports(runtime, priceOutputWithRandomPrice)
	if err != nil {
		runtime.Logger().Error("got expected error for GenerateReport with inconsistent prices", "error", err)
		return nil, fmt.Errorf("expected error for GenerateReport with inconsistent prices: %w", err)
	}

	runtime.Logger().Info("this is not expected: GenerateReport with inconsistent prices should return an error", "generated_report", reportWithInconsistentPrices)
	return reportWithInconsistentPrices, nil
}

// createPriceOutputWithRandomTimestamp creates multiple price outputs with different random timestamps
func createPriceOutputWithRandomTimestamp(runtime cre.Runtime, wfCfg config.Config) (priceOutput, error) {
	runtime.Logger().Info("creating price outputs with different timestamps")

	defaultFeedID, err := convertFeedIDtoBytes(wfCfg.FeedID)
	if err != nil {
		runtime.Logger().Error("failed to decode feed ID", "error", err)
		return priceOutput{}, fmt.Errorf("failed to decode feed ID: %w", err)
	}

	// Generate random timestamp variations (±1 hour from base time)
	randomTimestamp := newRandomTimestamp(runtime)
	runtime.Logger().Info("creating priceOutput")
	outputWithRandomTimestamp := priceOutput{
		FeedID:    defaultFeedID,
		Timestamp: randomTimestamp,
		Price:     big.NewInt(defaultPrice),
	}
	runtime.Logger().Info("priceOutput with random timestamp created")

	return outputWithRandomTimestamp, nil
}

// createPriceOutputWithRandomFeedID creates price output with random feedID
func createPriceOutputWithRandomFeedID(runtime cre.Runtime, wfCfg config.Config) (priceOutput, error) {
	runtime.Logger().Info("creating price output with random feedID")

	randomFeedID := newRandomFeedID(runtime)
	runtime.Logger().Info("creating priceOutput with random feedID")
	outputWithRandomFeedID := priceOutput{
		FeedID:    randomFeedID,
		Timestamp: defaultTimestamp,
		Price:     big.NewInt(defaultPrice),
	}
	runtime.Logger().Info("priceOutput with random feedID created")

	return outputWithRandomFeedID, nil
}

// createPriceOutputWithRandomPrice creates price output with random price
func createPriceOutputWithRandomPrice(runtime cre.Runtime, wfCfg config.Config) (priceOutput, error) {
	runtime.Logger().Info("creating price output with random price")

	feedID, err := convertFeedIDtoBytes(wfCfg.FeedID)
	if err != nil {
		runtime.Logger().Error("failed to decode feed ID", "error", err)
		return priceOutput{}, fmt.Errorf("failed to decode feed ID: %w", err)
	}

	randomPrice := newRandomPrice(runtime)
	runtime.Logger().Info("creating priceOutput with random price")
	outputWithRandomPrice := priceOutput{
		FeedID:    feedID,
		Timestamp: defaultTimestamp,
		Price:     randomPrice,
	}
	runtime.Logger().Info("priceOutput with random price created")

	return outputWithRandomPrice, nil
}

func newRandomTimestamp(runtime cre.Runtime) uint32 {
	baseTime := time.Now().Unix()
	randomOffset := rand.Int63n(randomTimeOffset) - (randomTimeOffset / 2) // ±1 hour in seconds
	randomTimestamp := uint32(baseTime + randomOffset)
	runtime.Logger().Info("new random timestamp created", "random_timestamp", randomTimestamp)
	return randomTimestamp
}

// newRandomFeedID generates a random 32-byte feedID
func newRandomFeedID(runtime cre.Runtime) [32]byte {
	var randomFeedID [32]byte
	for i := range randomFeedID {
		randomFeedID[i] = byte(rand.Intn(256))
	}
	runtime.Logger().Info("new random feedID created", "random_feedID", hex.EncodeToString(randomFeedID[:]))
	return randomFeedID
}

// newRandomPrice generates a random price between 1 and maxRandomPrice
func newRandomPrice(runtime cre.Runtime) *big.Int {
	randomPrice := big.NewInt(rand.Int63n(maxRandomPrice) + 1) // Random price between 1 and maxRandomPrice
	runtime.Logger().Info("new random price created", "random_price", randomPrice.String())
	return randomPrice
}

// generateReports encodes price outputs and generates a report using consensus
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
