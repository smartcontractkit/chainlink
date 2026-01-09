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
	_ MultitypeMonitoringEndpoint    = (*ChipIngressAgent)(nil)
)

// ChipIngressAgent implements both commontypes.MonitoringEndpoint and MultitypeMonitoringEndpoint.
// When created with a TelemType, it uses that type for SendLog calls.
// When created without a TelemType (via NewChipIngressAgentMultitype), SendLog is a no-op
// and SendTypedLog must be used instead.
type ChipIngressAgent struct {
	Network       string
	ChainID       string
	ContractID    string
	ChainSelector uint64
	telemService  synchronization.ChipIngressService
	lggr          logger.Logger

	TelemType synchronization.TelemetryType // Empty for multitype endpoints
	Domain    string                        // Derived from TelemetryType (empty for multitype)
	Entity    string                        // Derived from TelemetryType (empty for multitype)
}

// NewChipIngressAgent creates a new agent for a single telemetry type endpoint.
// It derives the chain selector from the Network and ChainID.
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

// NewChipIngressAgentMultitype creates a new agent for multitype telemetry endpoints.
// Unlike NewChipIngressAgent, the telemetry type is not set at construction time
// and must be provided with each SendTypedLog call.
func NewChipIngressAgentMultitype(
	telemService synchronization.ChipIngressService,
	network string,
	chainID string,
	contractID string,
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

	return &ChipIngressAgent{
		Network:       network,
		ChainID:       chainID,
		ContractID:    contractID,
		ChainSelector: details.ChainSelector,
		telemService:  telemService,
		lggr:          lggr,
		// TelemType, Domain, Entity left empty for multitype
	}, nil
}

// SendLog implements commontypes.MonitoringEndpoint.
// It forwards the telemetry log to the TelemetryService using the TelemType set at construction.
// For multitype agents (created via NewChipIngressAgentMultitype), this is a no-op.
//
// Note: This method does not accept a context parameter because it implements
// the commontypes.MonitoringEndpoint interface from libocr, which defines
// SendLog(log []byte) without a context. A background context is used internally.
func (a *ChipIngressAgent) SendLog(log []byte) {
	if a.TelemType == "" {
		// Multitype agent - SendLog should not be called, use SendTypedLog instead
		a.lggr.Warnw("SendLog called on multitype agent, use SendTypedLog instead")
		return
	}
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

// SendTypedLog implements MultitypeMonitoringEndpoint.
// It forwards the telemetry log to the TelemetryService with the specified telemetry type.
//
// Note: This method does not accept a context parameter because it implements
// the MultitypeMonitoringEndpoint interface, which mirrors the libocr
// MonitoringEndpoint.SendLog signature for consistency. A background context is used internally.
func (a *ChipIngressAgent) SendTypedLog(telemType synchronization.TelemetryType, log []byte) {
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
	a.telemService.Send(ctx, payload)
}
