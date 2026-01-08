package telemetry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	chainselector "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// Verify interface implementation at compile time
var (
	_ commontypes.MonitoringEndpoint = (*ChipIngressAgent)(nil)
	_ MultitypeMonitoringEndpoint    = (*ChipIngressAgentMultitype)(nil)
)

type ChipIngressAgent struct {
	Network       string
	ChainID       string
	ContractID    string
	ChainSelector uint64
	telemService  synchronization.ChipIngressService
	lggr          logger.Logger

	TelemType synchronization.TelemetryType // Stored for single-type SendLog
	Domain    string                        // Derived from TelemetryType (for single-type endpoints)
	Entity    string                        // Derived from TelemetryType (for single-type endpoints)
}

// NewChipIngressAgent creates a new adapter for a telemetryEndpoint
// It derives the chain selector from the Network and ChainID
func NewChipIngressAgent(
	telemService synchronization.ChipIngressService,
	network string,
	chainID string,
	contractID string,
	telemType synchronization.TelemetryType,
	lggr logger.Logger,
) (*ChipIngressAgent, error) {
	if telemService == nil {
		return nil, errors.New("telemetry service cannot be nil")
	}
	// Use chain-selectors package to get the ChainDetails which includes the selector
	details, err := chainselector.GetChainDetailsByChainIDAndFamily(chainID, strings.ToLower(network))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for chainID %s and network %s: %w", chainID, network, err)
	}

	domain, entity, err := synchronization.TelemetryTypeToDomainAndEntity(telemType)
	if err != nil {
		return nil, fmt.Errorf("failed to map telemetry type to domain/entity: %w", err)
	}

	return &ChipIngressAgent{
		Network:       network,
		ChainID:       chainID,
		ContractID:    contractID,
		ChainSelector: details.ChainSelector,
		Domain:        domain,
		Entity:        entity,
		TelemType:     telemType,
		telemService:  telemService,
		lggr:          lggr,
	}, nil
}

// SendLog implements commontypes.MonitoringEndpoint
// It forwards the telemetry log to the TelemetryService
func (a *ChipIngressAgent) SendLog(log []byte) {
	ctx := context.Background()
	payload := synchronization.TelemPayload{
		Telemetry:     log,
		TelemType:     a.TelemType,
		ContractID:    a.ContractID,
		ChainSelector: a.ChainSelector,
		Domain:        a.Domain,
		Entity:        a.Entity,
		Network:       a.Network,
	}
	a.telemService.Send(ctx, payload)
}

type ChipIngressAgentMultitype struct {
	Network            string
	ChainID            string
	ContractID         string
	ChainSelector      uint64
	chipIngressService synchronization.ChipIngressService
	lggr               logger.Logger
	// no TelemType
}

func NewChipIngressAgentMultitype(
	telemService synchronization.ChipIngressService,
	network string,
	chainID string,
	contractID string,
	lggr logger.Logger,
) (*ChipIngressAgentMultitype, error) {
	if telemService == nil {
		return nil, errors.New("telemetry service cannot be nil")
	}
	// Use chain-selectors package to get the ChainDetails which includes the selector
	details, err := chainselector.GetChainDetailsByChainIDAndFamily(chainID, strings.ToLower(network))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for chainID %s and network %s: %w", chainID, network, err)
	}

	return &ChipIngressAgentMultitype{
		Network:            network,
		ChainID:            chainID,
		ContractID:         contractID,
		ChainSelector:      details.ChainSelector,
		chipIngressService: telemService,
		lggr:               lggr,
	}, nil
}

// SendTypedLog implements MultitypeMonitoringEndpoint
// It forwards the telemetry log to the TelemetryService with the specified telemetry type
func (a *ChipIngressAgentMultitype) SendTypedLog(telemType synchronization.TelemetryType, log []byte) {
	ctx := context.Background()

	domain, entity, err := synchronization.TelemetryTypeToDomainAndEntity(telemType)
	if err != nil {
		a.lggr.Errorw("failed to map telemetry type to domain/entity", "error", err, "telemType", telemType)
		return
	}

	payload := synchronization.TelemPayload{
		Telemetry:     log,
		TelemType:     telemType,
		ContractID:    a.ContractID,
		ChainSelector: a.ChainSelector,
		Domain:        domain,
		Entity:        entity,
		Network:       a.Network,
	}
	a.chipIngressService.Send(ctx, payload)
}
