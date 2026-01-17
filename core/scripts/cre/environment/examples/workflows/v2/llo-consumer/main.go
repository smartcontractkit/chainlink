//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v3"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"

	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// MAGIC_NUMBER is the value Mock EA returns for TEST/USD (Stream ID 1)
// If this exact value appears in the workflow output, it proves full E2E:
// Mock EA (424242) → Stream Jobs → LLO Plugin → CRE Transmitter → streams-trigger → Workflow
const MAGIC_NUMBER = 424242

type WorkflowConfig struct {
	StreamIDs      []uint32 `yaml:"stream_ids"`
	MaxFrequencyMs uint64   `yaml:"max_frequency_ms"`
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (WorkflowConfig, error) {
		cfg := WorkflowConfig{
			StreamIDs:      []uint32{1},
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
	logger.Info(fmt.Sprintf("LLO_CONSUMER_STARTING: streams=%v, expecting MAGIC_NUMBER=%d", config.StreamIDs, MAGIC_NUMBER))

	trigger := NewStreamsTrigger(&streams.Config{
		StreamIds:      config.StreamIDs,
		MaxFrequencyMs: config.MaxFrequencyMs,
	})

	return cre.Workflow[WorkflowConfig]{
		cre.Handler(trigger, onStreamsTrigger),
	}, nil
}

// onStreamsTrigger - the key function that proves we can read the value from Streams DON
func onStreamsTrigger(config WorkflowConfig, runtime cre.Runtime, report *streams.Report) (string, error) {
	logger := runtime.Logger()

	if report == nil {
		return "", fmt.Errorf("LLO_E2E_ERROR: nil report")
	}

	// Decode the value from the OCRTriggerReport
	value, err := decodeOCRTriggerReport(report.Report)
	if err != nil {
		return "", fmt.Errorf("LLO_E2E_ERROR[SeqNr=%d]: decode failed: %v", report.SeqNr, err)
	}

	// Check if value matches MAGIC_NUMBER
	valueInt := value.IntPart()
	isMatch := valueInt == MAGIC_NUMBER

	logger.Info(fmt.Sprintf("LLO_WORKFLOW_DECODED: Value=%d Expected=%d Match=%v", valueInt, MAGIC_NUMBER, isMatch))

	// OUTPUT - this is what the test asserts on
	return "", fmt.Errorf("LLO_E2E_VALUE[SeqNr=%d]: Value=%d Expected=%d Match=%v",
		report.SeqNr, valueInt, MAGIC_NUMBER, isMatch)
}

// decodeOCRTriggerReport extracts the decimal value from the OCRTriggerReport
// The report structure is:
//
//	OCRTriggerReport {
//	  Outputs: Map {
//	    "Payload": List[
//	      Map { "StreamID": int64, "Decimal": bytes }
//	    ]
//	  }
//	}
func decodeOCRTriggerReport(data []byte) (decimal.Decimal, error) {
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

	// Get first stream value (we're expecting stream ID 1)
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
