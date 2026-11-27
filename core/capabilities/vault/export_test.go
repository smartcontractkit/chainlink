package vault

import (
	"context"
	"encoding/json"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

// ServeGatewayMessage runs the request synchronously, as the goroutine that
// HandleGatewayMessage starts does, so tests can assert on the result directly.
func (h *GatewayHandler) ServeGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) error {
	return h.serveGatewayMessage(ctx, gatewayID, req)
}
