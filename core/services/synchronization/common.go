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
	Send(ctx context.Context, telemetry []byte, contractID string, telemType TelemetryType, chainSelector uint64, domain string, entity string)
}

// TelemetryTypeToDomainAndEntity maps TelemetryType to (domain, entity) pairs for beholder ingestion.
// Based on atlas/ingress mappings.
func TelemetryTypeToDomainAndEntity(telemType TelemetryType) (domain, entity string, err error) {
	switch telemType {
	case EnhancedEA:
		return "data-feeds", "enhanced.ea.telemetry", nil
	case FunctionsRequests:
		return "functions", "functions.requests.telemetry", nil
	case EnhancedEAMercury:
		return "data-streams", "enhanced.ea.mercury.telemetry", nil
	case OCR:
		return "data-feeds", "ocr.v1.telemetry", nil
	case OCR2Automation:
		return "automation", "ocr.v2.automation.telemetry", nil
	case OCR2Functions:
		return "functions", "ocr.v2.functions.telemetry", nil
	case OCR2CCIPCommit:
		return "ccip", "ocr.v2.ccip.commit.telemetry", nil
	case OCR2CCIPExec:
		return "ccip", "ocr.v2.ccip.exec.telemetry", nil
	case OCR2Threshold:
		return "functions", "ocr.v2.threshold.telemetry", nil
	case OCR2S4:
		return "functions", "ocr.v2.s4.telemetry", nil
	case OCR2Median:
		return "data-feeds", "ocr.v2.median.telemetry", nil
	case OCR3Mercury:
		return "data-streams", "ocr.v3.mercury.telemetry", nil
	case OCR3DataFeeds:
		return "data-streams", "ocr.v3.data-feeds.telemetry", nil
	case AutomationCustom:
		return "automation", "automation.custom.telemetry", nil
	case OCR3Automation:
		return "automation", "ocr.v3.automation.telemetry", nil
	case OCR3Rebalancer:
		return "ccip", "ocr.v3.rebalancer.telemetry", nil
	case OCR3CCIPCommit:
		return "ccip", "ocr.v3.ccip.commit.telemetry", nil
	case OCR3CCIPExec:
		return "ccip", "ocr.v3.ccip.exec.telemetry", nil
	case OCR3CCIPBootstrap:
		return "ccip", "ocr.v3.ccip.bootstrap.telemetry", nil
	case HeadReport:
		return "platform", "head.report.telemetry", nil
	case PipelineBridge:
		return "data-feeds", "pipeline.bridge.telemetry", nil
	case LLOObservation:
		return "data-streams", "llo.observation.telemetry", nil
	case LLOOutcome:
		return "data-streams", "llo.outcome.telemetry", nil
	case LLOReport:
		return "data-streams", "llo.report.telemetry", nil
	default:
		return "", "", fmt.Errorf("unknown telemetry type: %s", telemType)
	}
}
