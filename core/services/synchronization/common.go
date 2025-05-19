package synchronization

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// TelemetryType defines supported telemetry types
type TelemetryType string

const (
	EnhancedEA        TelemetryType = "enhanced-ea"
	FunctionsRequests TelemetryType = "functions-requests"
	EnhancedEAMercury TelemetryType = "enhanced-ea-mercury"
	OCR               TelemetryType = "ocr"
	OCR2Automation    TelemetryType = "ocr2-automation"
	OCR2Functions     TelemetryType = "ocr2-functions"
	OCR2CCIPCommit    TelemetryType = "ocr2-ccip-commit"
	OCR2CCIPExec      TelemetryType = "ocr2-ccip-exec"
	OCR2Threshold     TelemetryType = "ocr2-threshold"
	OCR2S4            TelemetryType = "ocr2-s4"
	OCR2Median        TelemetryType = "ocr2-median"
	OCR3Mercury       TelemetryType = "ocr3-mercury"
	OCR3DataFeeds     TelemetryType = "ocr3-data-feeds"
	AutomationCustom  TelemetryType = "automation-custom"
	OCR3Automation    TelemetryType = "ocr3-automation"
	OCR3Rebalancer    TelemetryType = "ocr3-rebalancer"
	OCR3CCIPCommit    TelemetryType = "ocr3-ccip-commit"
	OCR3CCIPExec      TelemetryType = "ocr3-ccip-exec"
	OCR3CCIPBootstrap TelemetryType = "ocr3-bootstrap"
	HeadReport        TelemetryType = "head-report"

	PipelineBridge TelemetryType = "pipeline-bridge"
	LLOObservation TelemetryType = "llo-observation"
	LLOOutcome     TelemetryType = "llo-outcome"
	LLOReport      TelemetryType = "llo-report"
)

type TelemPayload struct {
	Telemetry  []byte
	TelemType  TelemetryType
	ContractID string
}

// TelemetryService encapsulates all the functionality needed to
// send telemetry to the ingress server using wsrpc
type TelemetryService interface {
	services.ServiceCtx
	Send(ctx context.Context, telemetry []byte, contractID string, telemType TelemetryType)
}
