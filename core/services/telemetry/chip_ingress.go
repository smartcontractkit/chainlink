package telemetry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	chainselector "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// Verify interface implementation at compile time
var _ commontypes.MonitoringEndpoint = (*ChipIngressAdapter)(nil)

// Emitter is the interface for the beholder.Emitter
// https://github.com/smartcontractkit/chainlink-common/blob/main/pkg/beholder/client.go#L27
type Emitter interface {
	// Emit is async and does not block the main thread
	Emit(ctx context.Context, body []byte, attrKVs ...any) error
}

// ChipIngressAdapter implements commontypes.MonitoringEndpoint
// It derives the chain selector from the Network and ChainID and forwards
// SendLog calls to the beholder.Emitter
// Used for migration from OTI to ChIP Ingress
type ChipIngressAdapter struct {
	Network       string
	ChainID       string
	ContractID    string
	ChainSelector uint64

	Domain string // Derived from TelemetryType
	Entity string // Derived from TelemetryType

	emitter Emitter
	lggr    logger.SugaredLogger
}

// NewChipIngressAdapter creates a new adapter for a telemetryEndpoint
// It derives the chain selector from the Network and ChainID
func NewChipIngressAdapter(
	network string,
	chainID string,
	contractID string,
	telemType synchronization.TelemetryType,
	emitter Emitter,
	lggr logger.SugaredLogger,
) (*ChipIngressAdapter, error) {
	if emitter == nil {
		return nil, errors.New("beholder emitter cannot be nil")
	}

	// Use chain-selectors package to get the ChainDetails which includes the selector
	details, err := chainselector.GetChainDetailsByChainIDAndFamily(chainID, strings.ToLower(network))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for chainID %s and network %s: %w", chainID, network, err)
	}

	domain, entity, err := telemTypeToDomainAndEntity(telemType)
	if err != nil {
		return nil, fmt.Errorf("failed to map telemetry type to domain/entity: %w", err)
	}

	return &ChipIngressAdapter{
		Network:       network,
		ChainID:       chainID,
		ContractID:    contractID,
		ChainSelector: details.ChainSelector,
		Domain:        domain,
		Entity:        entity,
		emitter:       emitter,
		lggr:          lggr,
	}, nil
}

// SendLog implements commontypes.MonitoringEndpoint
// It forwards the telemetry log to the beholder emitter with proper domain/entity attributes
func (a *ChipIngressAdapter) SendLog(log []byte) {
	// Not need to use context.WithTimeout because Emit is async and uses context.WithoutCancel(ctx)
	ctx := context.Background()
	// Emit is asyc and does not block the main thread
	err := a.emitter.Emit(ctx, log,
		"beholder_domain", a.Domain,
		"beholder_entity", a.Entity,
		"chain_id", a.ChainID,
		"network_name", a.Network,
		"chain_selector", a.ChainSelector,
		"contract_id", a.ContractID,
	)
	if err != nil {
		a.lggr.Errorw("failed to emit telemetry to beholder", "error", err)
	}
}

// telemTypeToDomainAndEntity maps TelemetryType to (domain, entity) pairs for beholder
// This function is based on the mapping from:
// https://github.com/smartcontractkit/atlas/blob/e0dfd7dbd28fc79890e8d0bcae6b9c8eddfba01b/ingress/ocr-telemetry/app/chip_ingress_batcher.go#L232-L278
func telemTypeToDomainAndEntity(telemType synchronization.TelemetryType) (domain, entity string, err error) {
	switch telemType {
	case synchronization.OCR:
		return "data-feeds", "ocr.v1.telemetry", nil
	case synchronization.OCR2Median:
		return "data-feeds", "ocr.v2.median.telemetry", nil
	case synchronization.OCR2Automation:
		return "automation", "ocr.v2.automation.telemetry", nil
	case synchronization.OCR3Automation:
		return "automation", "ocr.v3.automation.telemetry", nil
	case synchronization.AutomationCustom:
		return "automation", "automation.custom.telemetry", nil
	case synchronization.OCR2Functions:
		return "functions", "ocr.v2.functions.telemetry", nil
	case synchronization.OCR2Threshold:
		return "functions", "ocr.v2.threshold.telemetry", nil
	case synchronization.OCR2S4:
		return "functions", "ocr.v2.s4.telemetry", nil
	case synchronization.FunctionsRequests:
		return "functions", "functions.requests.telemetry", nil
	case synchronization.OCR2CCIPCommit:
		return "ccip", "ocr.v2.ccip.commit.telemetry", nil
	case synchronization.OCR2CCIPExec:
		return "ccip", "ocr.v2.ccip.exec.telemetry", nil
	case synchronization.OCR3CCIPCommit:
		return "ccip", "ocr.v3.ccip.commit.telemetry", nil
	case synchronization.OCR3CCIPExec:
		return "ccip", "ocr.v3.ccip.exec.telemetry", nil
	case synchronization.OCR3CCIPBootstrap:
		return "ccip", "ocr.v3.ccip.bootstrap.telemetry", nil
	case synchronization.OCR3Rebalancer:
		return "ccip", "ocr.v3.rebalancer.telemetry", nil
	case synchronization.OCR3Mercury:
		return "data-streams", "ocr.v3.mercury.telemetry", nil
	case synchronization.OCR3DataFeeds:
		return "data-streams", "ocr.v3.data-feeds.telemetry", nil
	case synchronization.EnhancedEA:
		return "data-feeds", "enhanced.ea.telemetry", nil
	case synchronization.EnhancedEAMercury:
		return "data-streams", "enhanced.ea.mercury.telemetry", nil
	case synchronization.HeadReport:
		return "platform", "head.report.telemetry", nil
	case synchronization.PipelineBridge:
		return "data-feeds", "pipeline.bridge.telemetry", nil
	case synchronization.LLOObservation:
		return "data-streams", "llo.observation.telemetry", nil
	case synchronization.LLOOutcome:
		return "data-streams", "llo.outcome.telemetry", nil
	case synchronization.LLOReport:
		return "data-streams", "llo.report.telemetry", nil
	default:
		return "", "", fmt.Errorf("unknown telemetry type: %s", telemType)
	}
}
