package synchronization

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// Verify interface implementation at compile time
var (
	_ ChipIngressService = (*chipIngressBatchClient)(nil)
	_ services.Service   = (*chipIngressBatchClient)(nil)
)

// chipIngressBatchClient is a stub implementation that wraps chipingress.Client
type chipIngressBatchClient struct {
	services.Service
	eng *services.Engine
}

// NewChipIngressBatchClient creates a new ChipIngressService that uses the chipingress.Client
// This is a stub implementation for now
func NewChipIngressBatchClient(endpoint interface{}, cfg interface{}, ks keystore.CSA, lggr logger.Logger, chipClient interface{}) ChipIngressService {
	c := &chipIngressBatchClient{}
	c.Service, c.eng = services.Config{
		Name:  "ChipIngressBatchClient",
		Start: c.start,
		Close: c.close,
	}.NewServiceEngine(lggr)
	return c
}

// start implements the start logic for the chip ingress batch client
func (c *chipIngressBatchClient) start(ctx context.Context) error {
	c.eng.Info("ChIP ingress batch client started")
	return nil
}

// close implements the close logic for the chip ingress batch client
func (c *chipIngressBatchClient) close() error {
	c.eng.Info("ChIP ingress batch client closed")
	return nil
}

// Send implements ChipIngressService
func (c *chipIngressBatchClient) Send(ctx context.Context, telemetry []byte, contractID string, telemType TelemetryType, chainSelector uint64, domain string, entity string) {
	// Stub implementation - logs that chip ingress is being used
	c.eng.Debugw("ChIP ingress send (stub)", "contractID", contractID, "telemType", telemType, "chainSelector", chainSelector, "domain", domain, "entity", entity)
}
