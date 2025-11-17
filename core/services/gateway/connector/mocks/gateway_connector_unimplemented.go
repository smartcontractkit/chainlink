package mocks

import (
	"context"
	"encoding/json"

	"github.com/stretchr/testify/mock"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

// Work around the fact that mockery doesn't support embedding interfaces in generated mocks.
// Any addition of a method to the interface will require adding a passthrough here
// to the mock. This is needed to disambiguate method calls between the `BaseGatewayConnector`
// and the `core.UnimplementedGatewayConnector` implementations.
type GatewayConnector struct {
	*BaseGatewayConnector
	core.UnimplementedGatewayConnector
}

func (g *GatewayConnector) AddHandler(ctx context.Context, methods []string, handler core.GatewayConnectorHandler) error {
	return g.BaseGatewayConnector.AddHandler(ctx, methods, handler)
}

func (g *GatewayConnector) RemoveHandler(ctx context.Context, methods []string) error {
	return g.BaseGatewayConnector.RemoveHandler(ctx, methods)
}

func (g *GatewayConnector) SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error {
	return g.BaseGatewayConnector.SendToGateway(ctx, gatewayID, resp)
}

func (g *GatewayConnector) SignMessage(ctx context.Context, msg []byte) ([]byte, error) {
	return g.BaseGatewayConnector.SignMessage(ctx, msg)
}

func (g *GatewayConnector) GatewayIDs(ctx context.Context) ([]string, error) {
	return g.BaseGatewayConnector.GatewayIDs(ctx)
}

func (g *GatewayConnector) DonID(ctx context.Context) (string, error) {
	return g.BaseGatewayConnector.DonID(ctx)
}

func (g *GatewayConnector) AwaitConnection(ctx context.Context, gatewayID string) error {
	return g.BaseGatewayConnector.AwaitConnection(ctx, gatewayID)
}

func NewGatewayConnector(t interface {
	mock.TestingT
	Cleanup(func())
}) *GatewayConnector {
	mock := NewBaseGatewayConnector(t)
	return &GatewayConnector{BaseGatewayConnector: mock}
}
