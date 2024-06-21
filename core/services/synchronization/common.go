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
	OCR2Threshold     TelemetryType = "ocr2-threshold"
	OCR2S4            TelemetryType = "ocr2-s4"
	OCR2Median        TelemetryType = "ocr2-median"
	OCR3Mercury       TelemetryType = "ocr3-mercury"
	AutomationCustom  TelemetryType = "automation-custom"
	OCR3Automation    TelemetryType = "ocr3-automation"
	HeadReport        TelemetryType = "head-report"
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
