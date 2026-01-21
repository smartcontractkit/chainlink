//go:build wasip1

package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v3"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"

	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// Magic numbers embedded in different report formats to prove E2E connectivity:
//
// Format 5 (CapabilityTrigger) - protobuf encoded:
//
//	Stream 1: TEST/USD → MAGIC_NUMBER_FORMAT5 (424242)
//
// Format 7 (EVMABIEncodeUnpackedExpr) - ABI encoded:
//
//	Stream 4: DATA/USD → MAGIC_NUMBER_FORMAT7 (555555)
const (
	MAGIC_NUMBER_FORMAT5 = 424242 // For ReportFormat 5 (CapabilityTrigger)
	MAGIC_NUMBER_FORMAT7 = 555555 // For ReportFormat 7 (EVMABIEncodeUnpackedExpr)
)

type WorkflowConfig struct {
	StreamIDs      []uint32 `yaml:"stream_ids"`
	MaxFrequencyMs uint64   `yaml:"max_frequency_ms"`
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (WorkflowConfig, error) {
		cfg := WorkflowConfig{
			StreamIDs:      []uint32{1, 4}, // Subscribe to both Format 5 and Format 7 streams
			MaxFrequencyMs: 1000,
		}
		if len(configBytes) > 0 {
			if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
				return WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
			}
		}
		return cfg, nil
	}).Run(RunLLOConsumerWorkflow)
}

type streamsTrigger struct {
	config *anypb.Any
}

func NewStreamsTrigger(cfg *streams.Config) cre.Trigger[*streams.Report, *streams.Report] {
	configAny := &anypb.Any{}
	if cfg != nil {
		_ = anypb.MarshalFrom(configAny, cfg, proto.MarshalOptions{Deterministic: true})
	}
	return &streamsTrigger{config: configAny}
}

func (*streamsTrigger) IsTrigger()                {}
func (*streamsTrigger) NewT() *streams.Report     { return &streams.Report{} }
func (*streamsTrigger) CapabilityID() string      { return "streams-trigger@2.0.0" }
func (*streamsTrigger) Method() string            { return "" }
func (t *streamsTrigger) ConfigAsAny() *anypb.Any { return t.config }
func (t *streamsTrigger) Adapt(report *streams.Report) (*streams.Report, error) {
	if report == nil {
		return &streams.Report{}, nil
	}
	return report, nil
}

func RunLLOConsumerWorkflow(config WorkflowConfig, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[WorkflowConfig], error) {
	logger.Info(fmt.Sprintf("LLO_CONSUMER_STARTING: streams=%v, expecting MAGIC_FORMAT5=%d, MAGIC_FORMAT7=%d",
		config.StreamIDs, MAGIC_NUMBER_FORMAT5, MAGIC_NUMBER_FORMAT7))

	trigger := NewStreamsTrigger(&streams.Config{
		StreamIds:      config.StreamIDs,
		MaxFrequencyMs: config.MaxFrequencyMs,
	})

	return cre.Workflow[WorkflowConfig]{
		cre.Handler(trigger, onStreamsTrigger),
	}, nil
}

// onStreamsTrigger processes LLO reports from the Streams DON
//
// SECURITY NOTE:
// By the time data reaches this workflow, it has already been verified by the
// TriggerSubscriber's SignedReportRemoteAggregator which:
//   - Performs cryptographic signature verification (ecrecover)
//   - Checks signers against the on-chain capabilities registry
//   - Requires F+1 valid signatures for BFT consensus
//   - Deduplicates reports (only one event per SeqNr)
//   - Rejects stale reports (maxAgeSec)
//
// The workflow can trust the data it receives.
func onStreamsTrigger(config WorkflowConfig, runtime cre.Runtime, report *streams.Report) (string, error) {
	logger := runtime.Logger()

	if report == nil {
		return "", fmt.Errorf("LLO_E2E_ERROR: nil report")
	}

	// Log basic report info for observability
	sigCount := len(report.Sigs)
	configDigest := hex.EncodeToString(report.ConfigDigest)
	if len(configDigest) > 16 {
		configDigest = configDigest[:16] + "..."
	}
	reportLen := len(report.Report)

	logger.Info(fmt.Sprintf("LLO_REPORT_RECEIVED[SeqNr=%d]: %d signatures, configDigest=%s, reportLen=%d",
		report.SeqNr, sigCount, configDigest, reportLen))

	// Try both Format 5 (protobuf) and Format 7 (ABI) decoding
	// Log results for both so we can verify both formats are working
	var results []string

	// Try Format 5 (CapabilityTrigger - protobuf)
	decF5, errF5 := decodeFormat5Report(report.Report)
	if errF5 == nil {
		valueInt := decF5.IntPart()
		isMatch := valueInt == MAGIC_NUMBER_FORMAT5
		result := fmt.Sprintf("LLO_E2E_FORMAT5[SeqNr=%d]: Value=%d Expected=%d Match=%v Sigs=%d",
			report.SeqNr, valueInt, MAGIC_NUMBER_FORMAT5, isMatch, sigCount)
		logger.Info(result)
		results = append(results, result)
	} else {
		logger.Info(fmt.Sprintf("LLO_E2E_FORMAT5_SKIP[SeqNr=%d]: %v", report.SeqNr, errF5))
	}

	// Try Format 7 (EVMABIEncodeUnpackedExpr - ABI)
	valueF7, errF7 := decodeFormat7Report(report.Report)
	if errF7 == nil {
		isMatch := valueF7 == MAGIC_NUMBER_FORMAT7
		result := fmt.Sprintf("LLO_E2E_FORMAT7[SeqNr=%d]: Value=%d Expected=%d Match=%v Sigs=%d",
			report.SeqNr, valueF7, MAGIC_NUMBER_FORMAT7, isMatch, sigCount)
		logger.Info(result)
		results = append(results, result)
	} else {
		logger.Info(fmt.Sprintf("LLO_E2E_FORMAT7_SKIP[SeqNr=%d]: %v", report.SeqNr, errF7))
	}

	// Return results based on which format was decoded successfully
	// Format 5 (protobuf) takes priority if both decode (shouldn't happen in practice)
	if errF5 == nil {
		valueInt := decF5.IntPart()
		isMatch := valueInt == MAGIC_NUMBER_FORMAT5
		return "", fmt.Errorf("LLO_E2E_VALUE[SeqNr=%d]: Format=5 Value=%d Expected=%d Match=%v Sigs=%d",
			report.SeqNr, valueInt, MAGIC_NUMBER_FORMAT5, isMatch, sigCount)
	}

	if errF7 == nil {
		isMatch := valueF7 == MAGIC_NUMBER_FORMAT7
		return "", fmt.Errorf("LLO_E2E_VALUE[SeqNr=%d]: Format=7 Value=%d Expected=%d Match=%v Sigs=%d",
			report.SeqNr, valueF7, MAGIC_NUMBER_FORMAT7, isMatch, sigCount)
	}

	// Neither format decoded successfully
	return "", fmt.Errorf("LLO_E2E_ERROR[SeqNr=%d]: decode failed - Format5: %v, Format7: %v, reportLen=%d",
		report.SeqNr, errF5, errF7, reportLen)
}

// decodeFormat5Report extracts the decimal value from ReportFormat 5 (CapabilityTrigger)
// This is protobuf encoded with the structure:
//
//	OCRTriggerReport {
//	  Outputs: Map {
//	    "Payload": List[
//	      Map { "StreamID": int64, "Decimal": bytes }
//	    ]
//	  }
//	}
func decodeFormat5Report(data []byte) (decimal.Decimal, error) {
	// Unmarshal the OCRTriggerReport protobuf
	ocrReport := &capabilitiespb.OCRTriggerReport{}
	if err := proto.Unmarshal(data, ocrReport); err != nil {
		return decimal.Zero, fmt.Errorf("unmarshal OCRTriggerReport: %w", err)
	}

	// Get Outputs map
	outputs := ocrReport.GetOutputs()
	if outputs == nil || outputs.Fields == nil {
		return decimal.Zero, fmt.Errorf("no Outputs in report")
	}

	// Get Payload field (list of stream values)
	payloadValue, ok := outputs.Fields["Payload"]
	if !ok || payloadValue == nil {
		return decimal.Zero, fmt.Errorf("no Payload field in Outputs")
	}

	// Extract the list
	payloadList := payloadValue.GetListValue()
	if payloadList == nil {
		return decimal.Zero, fmt.Errorf("Payload is not a list")
	}

	fields := payloadList.GetFields()
	if len(fields) == 0 {
		return decimal.Zero, fmt.Errorf("Payload list is empty")
	}

	// Get first stream value
	firstStreamValue := fields[0]
	streamMap := firstStreamValue.GetMapValue()
	if streamMap == nil || streamMap.Fields == nil {
		return decimal.Zero, fmt.Errorf("first payload item is not a map")
	}

	// Get Decimal bytes
	decimalValue, ok := streamMap.Fields["Decimal"]
	if !ok || decimalValue == nil {
		return decimal.Zero, fmt.Errorf("no Decimal field in stream value")
	}

	decimalBytes := decimalValue.GetBytesValue()
	if len(decimalBytes) == 0 {
		return decimal.Zero, fmt.Errorf("Decimal bytes are empty")
	}

	// Unmarshal the decimal using shopspring/decimal binary format
	var dec decimal.Decimal
	if err := dec.UnmarshalBinary(decimalBytes); err != nil {
		return decimal.Zero, fmt.Errorf("unmarshal decimal: %w", err)
	}

	return dec, nil
}

// decodeFormat7Report extracts the value from ReportFormat 7 (EVMABIEncodeUnpackedExpr)
// This is ABI encoded with the structure:
//
//	Header (ABI-padded to 32-byte boundaries):
//	  - FeedID (32 bytes)
//	  - ValidFromTimestamp (32 bytes, padded)
//	  - Timestamp (32 bytes, padded)
//	  - NativeFee (32 bytes)
//	  - LinkFee (32 bytes)
//	  - ExpiresAt (32 bytes, padded)
//
//	Payload (ABI encoded):
//	  - Data values as int192
func decodeFormat7Report(data []byte) (int64, error) {
	// Format 7 header size: 6 x 32 bytes = 192 bytes (ABI-padded)
	const headerSize = 192

	if len(data) < headerSize {
		return 0, fmt.Errorf("report too short for Format 7: %d bytes (need at least %d)", len(data), headerSize)
	}

	// The payload follows the header
	payload := data[headerSize:]

	if len(payload) < 32 {
		return 0, fmt.Errorf("payload too short: %d bytes", len(payload))
	}

	// The payload contains ABI-encoded int192 values
	// int192 is encoded as a 32-byte value (left-padded for positive, sign-extended for negative)
	// We extract the first value which is DATA/USD with MAGIC_NUMBER_FORMAT7
	valueBytes := payload[:32]

	// Check if it's a negative number (sign bit set)
	isNegative := valueBytes[0]&0x80 != 0
	if isNegative {
		return 0, fmt.Errorf("negative values not supported")
	}

	// For positive numbers, use big.Int for proper handling
	value := new(big.Int).SetBytes(valueBytes)
	if value.IsInt64() {
		return value.Int64(), nil
	}

	// Fallback: try reading as big-endian uint64 from last 8 bytes
	if len(valueBytes) >= 8 {
		return int64(binary.BigEndian.Uint64(valueBytes[len(valueBytes)-8:])), nil
	}

	return 0, fmt.Errorf("unable to decode value from payload")
}
