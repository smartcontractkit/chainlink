package v2

// Tests for the in-flight concurrency bounds on outbound HTTP actions, global and per node.

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	httpmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

var (
	globalOutboundLimit  = cresettings.Default.GatewayHTTPActionOutboundConcurrencyLimit.DefaultValue
	perNodeOutboundLimit = cresettings.Default.GatewayHTTPActionOutboundPerNodeConcurrencyLimit.DefaultValue
)

func actionMessage(t *testing.T, url string) *jsonrpc.Response[json.RawMessage] {
	t.Helper()
	reqBytes, err := json.Marshal(gateway_common.OutboundHTTPRequest{
		Method:           "GET",
		URL:              url,
		TimeoutMs:        5000,
		MaxResponseBytes: 1000,
	})
	require.NoError(t, err)
	raw := json.RawMessage(reqBytes)
	return &jsonrpc.Response[json.RawMessage]{
		ID:     fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String()),
		Result: &raw,
	}
}

func okHTTPResponse() *network.HTTPResponse {
	return &network.HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       []byte(`{"ok": true}`),
	}
}

// With every global slot taken, a further request must be refused before any outbound call.
func TestMakeOutgoingRequest_RejectsWhenGlobalConcurrencyExhausted(t *testing.T) {
	t.Parallel()
	handler := createTestHandler(t)
	for range globalOutboundLimit {
		require.NoError(t, handler.outboundConcurrencyLimiter.Use(t.Context(), 1))
	}

	err := handler.HandleNodeMessage(t.Context(), actionMessage(t, "https://example.com/api"), "node1")
	require.ErrorContains(t, err, "gateway outbound concurrency limit reached")
}

// One node must not be able to occupy more than its share
func TestMakeOutgoingRequest_RejectsWhenNodeLimitReached(t *testing.T) {
	t.Parallel()
	handler := createTestHandler(t)
	for range perNodeOutboundLimit {
		require.NoError(t, handler.perNodeOutboundLimiters["node1"].Use(t.Context(), 1))
	}

	err := handler.HandleNodeMessage(t.Context(), actionMessage(t, "https://example.com/api"), "node1")
	require.ErrorContains(t, err, "outbound concurrency limit reached for node node1")

	// The global slot taken before the per-node bound rejected must have been rolled back.
	available, err := handler.outboundConcurrencyLimiter.Available(t.Context())
	require.NoError(t, err)
	require.Equal(t, globalOutboundLimit, available, "global reservation must be rolled back")
}

// A node at its per-node bound must not consume capacity belonging to its peers.
func TestMakeOutgoingRequest_NodeLimitIsolatesPeers(t *testing.T) {
	t.Parallel()
	handler := createTestHandler(t)
	mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
	mockDon := handler.shards[0].connMgr.(*handlermocks.DON)

	mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(okHTTPResponse(), nil).Once()
	mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)

	for range perNodeOutboundLimit {
		require.NoError(t, handler.perNodeOutboundLimiters["node1"].Use(t.Context(), 1))
	}

	err := handler.HandleNodeMessage(t.Context(), actionMessage(t, "https://example.com/api"), "node1")
	require.ErrorContains(t, err, "outbound concurrency limit reached for node node1")

	require.NoError(t, handler.HandleNodeMessage(t.Context(), actionMessage(t, "https://example.com/ok"), "node2"))
	handler.wg.Wait()

	mockHTTPClient.AssertExpectations(t)
}

// Slots must be returned when a request finishes
func TestMakeOutgoingRequest_ReleasesSlotOnCompletion(t *testing.T) {
	t.Parallel()
	handler := createTestHandler(t)
	mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
	mockDon := handler.shards[0].connMgr.(*handlermocks.DON)

	mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(okHTTPResponse(), nil).Times(3)
	mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

	// Leave exactly one global slot free, so each request depends on its predecessor releasing.
	for range globalOutboundLimit - 1 {
		require.NoError(t, handler.outboundConcurrencyLimiter.Use(t.Context(), 1))
	}

	for i := range 3 {
		msg := actionMessage(t, fmt.Sprintf("https://example.com/%d", i))
		require.NoError(t, handler.HandleNodeMessage(t.Context(), msg, "node1"), "request %d", i)
		handler.wg.Wait()
	}

	available, err := handler.outboundConcurrencyLimiter.Available(t.Context())
	require.NoError(t, err)
	require.Equal(t, 1, available, "slot must be returned")
	mockHTTPClient.AssertExpectations(t)
}

// A failed outbound request must still release its slots on both bounds
func TestMakeOutgoingRequest_ReleasesSlotOnFailure(t *testing.T) {
	t.Parallel()
	handler := createTestHandler(t)
	mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
	mockDon := handler.shards[0].connMgr.(*handlermocks.DON)

	mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(nil, errors.New("upstream down")).Once()
	mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

	require.NoError(t, handler.HandleNodeMessage(t.Context(), actionMessage(t, "https://example.com/fail"), "node1"))
	handler.wg.Wait()

	globalAvailable, err := handler.outboundConcurrencyLimiter.Available(t.Context())
	require.NoError(t, err)
	require.Equal(t, globalOutboundLimit, globalAvailable, "global slot must be returned on failure")

	nodeAvailable, err := handler.perNodeOutboundLimiters["node1"].Available(t.Context())
	require.NoError(t, err)
	require.Equal(t, perNodeOutboundLimit, nodeAvailable, "per-node slot must be returned on failure")
}

// A node with no per-node limiter is not a DON member and must be refused
func TestAcquireOutboundSlot_RejectsUnknownNode(t *testing.T) {
	t.Parallel()
	handler := createTestHandler(t)

	release, err := handler.acquireOutboundSlot(t.Context(), "not-a-member")
	require.ErrorContains(t, err, "unexpected node")
	require.Nil(t, release)

	available, err := handler.outboundConcurrencyLimiter.Available(t.Context())
	require.NoError(t, err)
	require.Equal(t, globalOutboundLimit, available, "no global slot may be consumed")
}
