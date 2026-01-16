package synchronization

import (
	"context"
	"fmt"

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
	Telemetry     []byte
	TelemType     TelemetryType
	ContractID    string
	Domain        string
	Entity        string
	ChainSelector uint64
	Network       string
	ChainID       string
}

// TelemetryService encapsulates all the functionality needed to
// send telemetry to the ingress server using wsrpc
type TelemetryService interface {
	services.ServiceCtx
	Send(ctx context.Context, telemetry []byte, contractID string, telemType TelemetryType)
}

type ChipIngressService interface {
	services.ServiceCtx
	Send(ctx context.Context, payload TelemPayload)
}

// TelemetryTypeToDomainAndEntity maps TelemetryType to (domain, entity) pairs for beholder ingestion.
// Based on atlas/ingress mappings.
func TelemetryTypeToDomainAndEntity(telemType TelemetryType) (domain, entity string, err error) {
	ocr2Entity := "offchainreporting2.TelemetryWrapper"
	ocr3Entity := "offchainreporting3.TelemetryWrapper"

	switch telemType {
	case OCR:
		return "data-feeds.telemetry.ocr", "offchainreporting.TelemetryWrapper", nil
	case OCR2Median:
		return "data-feeds.telemetry.ocr2-median", ocr2Entity, nil
	case OCR2Automation:
		return "automations.telemetry.ocr2-automation", ocr2Entity, nil
	case OCR2CCIPCommit:
		return "ccip.telemetry.ocr2", ocr2Entity, nil
	case OCR2CCIPExec:
		return "ccip.telemetry.ocr2", ocr2Entity, nil
	case OCR2Functions:
		return "functions.telemetry.ocr2-functions", ocr2Entity, nil
	case OCR3Automation:
		return "automations.telemetry.ocr3-automation", ocr3Entity, nil
	case OCR3Mercury:
		return "data-streams.telemetry.ocr3-mercury", ocr3Entity, nil
	case OCR3CCIPCommit:
		return "ccip.telemetry.ocr3", ocr3Entity, nil
	case OCR3CCIPExec:
		return "ccip.telemetry.ocr3", ocr3Entity, nil
	case OCR3DataFeeds:
		return "data-feeds.telemetry.ocr3-data-feeds", ocr3Entity, nil
	case EnhancedEA:
		return "data-feeds.telemetry.enhanced-ea", "telem.EnhancedEA", nil
	case EnhancedEAMercury:
		return "data-streams.telemetry.enhanced-ea-mercury", "telem.EnhancedEAMercury", nil
	case LLOObservation:
		return "data-streams.telemetry.llo-observation", "telem.LLOObservationTelemetry", nil
	case LLOOutcome:
		return "data-streams.telemetry.llo-outcome", "telem.LLOOutcomeTelemetry", nil
	case AutomationCustom:
		return "automations.telemetry.automation-custom", "telem.AutomationTelemWrapper", nil
	case FunctionsRequests:
		return "functions.telemetry.functions-requests", "telem.FunctionsRequest", nil
	case HeadReport:
		return "ccip.telemetry.head-report", "telem.HeadReportRequest", nil
	case PipelineBridge:
		return "data-streams.telemetry.pipeline-bridge", "telem.LLOBridgeTelemetry", nil
	default:
		return "", "", fmt.Errorf("could not resolve domain and entity from telem type, unsupported telem type: %s", telemType)
	}
}
