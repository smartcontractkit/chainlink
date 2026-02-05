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

// Magic numbers embedded in different report formats to prove E2E connectivity.
// E2E expected values per stream ID (business logic keys off these for the test).
var e2eExpectedValuesByStreamID = map[uint32]int64{
	1: 424242, // Stream 1: Format 5 (CapabilityTrigger) → MAGIC_NUMBER_FORMAT5
	4: 555555, // Stream 4: Format 7 (EVMABIEncodeUnpackedExpr) → MAGIC_NUMBER_FORMAT7
}

// DefaultStreamsTriggerCapabilityID is the capability ID used for LLO streams trigger (real DON).
const DefaultStreamsTriggerCapabilityID = "streams-trigger@2.0.0"

// MockStreamsTriggerCapabilityID is used when the capabilities DON runs only the mock (mock test).
const MockStreamsTriggerCapabilityID = "mock@1.0.0"

// Report format identifiers. The payload (streams.Report) includes report_format; the transmitter sets it.
const (
	ReportFormat5 = 5 // CapabilityTrigger (protobuf)
	ReportFormat7 = 7 // EVMABIEncodeUnpackedExpr (ABI)
)

// StreamValues is stream_id → value (int64). Format 5 decimals are truncated to int via IntPart().
type StreamValues map[uint32]int64

type WorkflowConfig struct {
	StreamIDs             []uint32 `yaml:"stream_ids"`
	MaxFrequencyMs        uint64   `yaml:"max_frequency_ms"`
	TriggerCapabilityID   string   `yaml:"trigger_capability_id"`    // optional; default streams-trigger@2.0.0. Use mock@1.0.0 for mock test.
	AcceptedReportFormats []uint32 `yaml:"accepted_report_formats"`  // Report formats the workflow accepts; DON filters by this. e.g. [5, 7]. Empty = accept all.
	TransmissionWindowMs  uint64   `yaml:"transmission_window_ms"`   // When > 0, DON delays pushing until next wall-clock boundary. 0 = use runner default or immediate.
	ExpectedReportFormat int      `yaml:"expected_report_format"`   // Fallback for decode when report does not carry report_format (e.g. mock). 5=protobuf, 7=ABI.
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
	config       *anypb.Any
	capabilityID string
}

func NewStreamsTrigger(cfg *streams.Config, capabilityID string) cre.Trigger[*streams.Report, *streams.Report] {
	configAny := &anypb.Any{}
	if cfg != nil {
		_ = anypb.MarshalFrom(configAny, cfg, proto.MarshalOptions{Deterministic: true})
	}
	if capabilityID == "" {
		capabilityID = DefaultStreamsTriggerCapabilityID
	}
	return &streamsTrigger{config: configAny, capabilityID: capabilityID}
}

func (*streamsTrigger) IsTrigger()                {}
func (*streamsTrigger) NewT() *streams.Report     { return &streams.Report{} }
func (t *streamsTrigger) CapabilityID() string    { return t.capabilityID }
func (*streamsTrigger) Method() string            { return "Trigger" }
func (t *streamsTrigger) ConfigAsAny() *anypb.Any { return t.config }
func (t *streamsTrigger) Adapt(report *streams.Report) (*streams.Report, error) {
	if report == nil {
		return &streams.Report{}, nil
	}
	return report, nil
}

func RunLLOConsumerWorkflow(config WorkflowConfig, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[WorkflowConfig], error) {
	logger.Info("LLO_CONSUMER_STARTING", "streams", config.StreamIDs, "triggerCapabilityID", config.TriggerCapabilityID)

	trigger := NewStreamsTrigger(&streams.Config{
		StreamIds:             config.StreamIDs,
		MaxFrequencyMs:        config.MaxFrequencyMs,
		AcceptedReportFormats: config.AcceptedReportFormats,
		TransmissionWindowMs:  config.TransmissionWindowMs,
	}, config.TriggerCapabilityID)

	return cre.Workflow[WorkflowConfig]{
		cre.Handler(trigger, onStreamsTrigger),
	}, nil
}

// ExtractReportValues decodes the report. Report type is in the payload: report.GetReportFormat() (5 or 7).
// If the payload has report_format set, that is used; else config.ExpectedReportFormat (e.g. for mock).
// Exactly one decoder is used—never "try both". Returns values keyed by stream ID.
func ExtractReportValues(report *streams.Report, config WorkflowConfig) (values StreamValues, usedFormat int, err error) {
	if report == nil || len(report.Report) == 0 {
		return nil, 0, fmt.Errorf("nil or empty report")
	}
	data := report.Report
	streamIDs := config.StreamIDs

	// Report type is in the payload: report.GetReportFormat(). Use it when set; else config.ExpectedReportFormat (e.g. mock).
	// Never try both decoders—require an explicit format (5 or 7).
	format := config.ExpectedReportFormat
	if r := report.GetReportFormat(); r != 0 {
		format = int(r)
	}
	if format != ReportFormat5 && format != ReportFormat7 {
		return nil, 0, fmt.Errorf("report format must be 5 or 7 (got %d); set expected_report_format in workflow config or ensure transmitter sets report_format on the report", format)
	}

	if format == ReportFormat5 {
		values, err = extractFormat5Values(data)
		if err != nil {
			return nil, 0, fmt.Errorf("format 5: %w", err)
		}
		return values, ReportFormat5, nil
	}
	values, err = extractFormat7Values(data, streamIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("format 7: %w", err)
	}
	return values, ReportFormat7, nil
}

// onStreamsTrigger processes LLO reports from the Streams DON
//
// SECURITY NOTE:
// On the v2 path (streams-trigger@2.0.0 with MethodConfig), the TriggerSubscriber
// uses DefaultModeAggregator: it requires MinResponsesToAggregate (2f+1) identical
// responses from capability DON nodes before forwarding. There is no per-report
// cryptographic signature verification on this path. The workflow receives the
// consensus-aggregated report after 2f+1 nodes have sent the same payload.
func onStreamsTrigger(config WorkflowConfig, runtime cre.Runtime, report *streams.Report) (string, error) {
	logger := runtime.Logger()

	if report == nil {
		return "", fmt.Errorf("LLO_E2E_ERROR: nil report")
	}

	// Report type is in the payload: streams.Report.report_format (5 = protobuf, 7 = ABI). Transmitter sets it.
	reportFormatFromPayload := report.GetReportFormat()
	sigCount := len(report.Sigs)
	configDigest := hex.EncodeToString(report.ConfigDigest)
	if len(configDigest) > 16 {
		configDigest = configDigest[:16] + "..."
	}
	logger.Info(fmt.Sprintf("LLO_REPORT_RECEIVED[SeqNr=%d]: report_format=%d (from payload), sigs=%d, configDigest=%s, reportLen=%d",
		report.SeqNr, reportFormatFromPayload, sigCount, configDigest, len(report.Report)))

	// Top-level extraction: decode report and get values keyed by stream ID (regardless of format).
	values, usedFormat, err := ExtractReportValues(report, config)
	if err != nil {
		return "", fmt.Errorf("LLO_E2E_ERROR[SeqNr=%d]: %w", report.SeqNr, err)
	}

	// Business logic: key off stream IDs. For E2E we log and validate—only fail when a present value doesn't match expected (reports can be single-stream).
	expectedByStream := e2eExpectedValuesByStreamID
	allMatch := true
	var resultParts []string
	for _, streamID := range config.StreamIDs {
		val, ok := values[streamID]
		expected, hasExpected := expectedByStream[streamID]
		match := ok && hasExpected && val == expected
		logger.Info(fmt.Sprintf("LLO_E2E_STREAM[SeqNr=%d]: streamID=%d value=%d (present=%v) expected=%d match=%v",
			report.SeqNr, streamID, val, ok, expected, match))
		resultParts = append(resultParts, fmt.Sprintf("streamID=%d Value=%d", streamID, val))
		if ok && hasExpected && val != expected {
			allMatch = false
		}
	}
	msg := fmt.Sprintf("LLO_E2E_VALUE[SeqNr=%d]: Format=%d Sigs=%d %s",
		report.SeqNr, usedFormat, sigCount, fmt.Sprint(resultParts))
	logger.Info(msg)
	if !allMatch {
		return "", fmt.Errorf("LLO_E2E_MISMATCH[SeqNr=%d]: %s", report.SeqNr, msg)
	}
	return msg, nil
}

// extractFormat5Values extracts all stream values from ReportFormat 5 (CapabilityTrigger).
// Structure: OCRTriggerReport.Outputs["Payload"] = list of maps, each { "StreamID": int64, "Decimal": bytes }.
func extractFormat5Values(data []byte) (StreamValues, error) {
	ocrReport := &capabilitiespb.OCRTriggerReport{}
	if err := proto.Unmarshal(data, ocrReport); err != nil {
		return nil, fmt.Errorf("unmarshal OCRTriggerReport: %w", err)
	}
	outputs := ocrReport.GetOutputs()
	if outputs == nil || outputs.Fields == nil {
		return nil, fmt.Errorf("no Outputs in report")
	}
	payloadValue, ok := outputs.Fields["Payload"]
	if !ok || payloadValue == nil {
		return nil, fmt.Errorf("no Payload field in Outputs")
	}
	payloadList := payloadValue.GetListValue()
	if payloadList == nil {
		return nil, fmt.Errorf("Payload is not a list")
	}
	fields := payloadList.GetFields()
	if len(fields) == 0 {
		return nil, fmt.Errorf("Payload list is empty")
	}
	out := make(StreamValues, len(fields))
	for i, item := range fields {
		streamMap := item.GetMapValue()
		if streamMap == nil || streamMap.Fields == nil {
			return nil, fmt.Errorf("payload item %d is not a map", i)
		}
		streamIDVal, ok := streamMap.Fields["StreamID"]
		if !ok || streamIDVal == nil {
			return nil, fmt.Errorf("payload item %d: no StreamID", i)
		}
		streamID := uint32(streamIDVal.GetInt64Value())
		decimalVal, ok := streamMap.Fields["Decimal"]
		if !ok || decimalVal == nil {
			return nil, fmt.Errorf("payload item %d: no Decimal", i)
		}
		decimalBytes := decimalVal.GetBytesValue()
		if len(decimalBytes) == 0 {
			return nil, fmt.Errorf("payload item %d: empty Decimal", i)
		}
		var dec decimal.Decimal
		if err := dec.UnmarshalBinary(decimalBytes); err != nil {
			return nil, fmt.Errorf("payload item %d unmarshal decimal: %w", i, err)
		}
		out[streamID] = dec.IntPart()
	}
	return out, nil
}

// extractFormat7Values extracts stream values from ReportFormat 7 (EVMABIEncodeUnpackedExpr).
// Header 192 bytes; payload is a sequence of 32-byte int192 values. Values are paired by index
// with streamIDs (payload[i] → streamIDs[i]).
func extractFormat7Values(data []byte, streamIDs []uint32) (StreamValues, error) {
	const headerSize = 192
	const valueSize = 32
	if len(data) < headerSize {
		return nil, fmt.Errorf("report too short for Format 7: %d bytes (need at least %d)", len(data), headerSize)
	}
	payload := data[headerSize:]
	needLen := len(streamIDs) * valueSize
	if len(payload) < needLen {
		return nil, fmt.Errorf("payload too short: %d bytes (need %d for %d streams)", len(payload), needLen, len(streamIDs))
	}
	out := make(StreamValues, len(streamIDs))
	for i, streamID := range streamIDs {
		valueBytes := payload[i*valueSize : (i+1)*valueSize]
		if valueBytes[0]&0x80 != 0 {
			return nil, fmt.Errorf("stream %d: negative values not supported", streamID)
		}
		value := new(big.Int).SetBytes(valueBytes)
		if value.IsInt64() {
			out[streamID] = value.Int64()
			continue
		}
		if len(valueBytes) >= 8 {
			out[streamID] = int64(binary.BigEndian.Uint64(valueBytes[len(valueBytes)-8:]))
			continue
		}
		return nil, fmt.Errorf("stream %d: unable to decode value", streamID)
	}
	return out, nil
}
