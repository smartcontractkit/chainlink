package confidentialrelay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

type barrierDON struct {
	total       int
	mu          sync.Mutex
	started     int
	allStarted  chan struct{}
	releaseOnce sync.Once
}

func newBarrierDON(total int) *barrierDON {
	return &barrierDON{
		total:      total,
		allStarted: make(chan struct{}),
	}
}

func (d *barrierDON) SendToNode(ctx context.Context, _ string, _ *jsonrpc.Request[json.RawMessage]) error {
	d.mu.Lock()
	d.started++
	if d.started == d.total {
		d.releaseOnce.Do(func() { close(d.allStarted) })
	}
	ch := d.allStarted
	d.mu.Unlock()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

var nodeOne = config.NodeConfig{
	Name:    "node1",
	Address: "0x1234",
}

func setupHandler(t *testing.T, numNodes int) (*handler, *common.Callback, *mocks.DON, *clockwork.FakeClock) {
	t.Helper()
	lggr := logger.Test(t)
	don := mocks.NewDON(t)

	members := make([]config.NodeConfig, numNodes)
	for i := range numNodes {
		members[i] = config.NodeConfig{
			Name:    fmt.Sprintf("node%d", i),
			Address: fmt.Sprintf("0x%04d", i),
		}
	}

	donConfig := &config.DONConfig{
		DonId:   "test_relay_don",
		F:       1,
		Members: members,
	}
	handlerConfig := Config{
		RequestTimeoutSec: 30,
	}
	methodConfig, err := json.Marshal(handlerConfig)
	require.NoError(t, err)

	clock := clockwork.NewFakeClock()
	limitsFactory := limits.Factory{Settings: cresettings.DefaultGetter, Logger: lggr}
	h, err := NewHandler(methodConfig, donConfig, don, lggr, clock, limitsFactory)
	require.NoError(t, err)
	h.aggregator = &mockAggregator{}
	cb := common.NewCallback()
	return h, cb, don, clock
}

type mockAggregator struct {
	err error
}

func (m *mockAggregator) Aggregate(_ map[string]jsonrpc.Response[json.RawMessage], _ int, _ int, _ logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	return nil, m.err
}

type respondingMockAggregator struct{}

func (m *respondingMockAggregator) Aggregate(resps map[string]jsonrpc.Response[json.RawMessage], _ int, _ int, _ logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	if len(resps) == 0 {
		return nil, errInsufficientResponsesForQuorum
	}
	// Return the first response we find.
	for _, r := range resps {
		return &r, nil
	}
	return nil, errInsufficientResponsesForQuorum
}

func TestConfidentialRelayHandler_Methods(t *testing.T) {
	h, _, _, _ := setupHandler(t, 4)
	methods := h.Methods()
	assert.Equal(t, []string{MethodSecretsGet, MethodCapabilityExec}, methods)
}

func TestConfidentialRelayHandler_HandleLegacyUserMessage(t *testing.T) {
	h, cb, _, _ := setupHandler(t, 4)
	err := h.HandleLegacyUserMessage(t.Context(), nil, cb)
	require.ErrorContains(t, err, "confidential relay handler does not support legacy messages")
}

func TestConfidentialRelayHandler_RequestIDTooLong(t *testing.T) {
	h, cb, _, _ := setupHandler(t, 4)

	longID := strings.Repeat("x", 201)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     longID,
		Method: MethodCapabilityExec,
	}

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	expected := fmt.Sprintf("request ID is too long: %d. max is 200 characters", len(longID))
	require.EqualError(t, err, expected)
}

func TestConfidentialRelayHandler_EmptyRequestID(t *testing.T) {
	h, cb, _, _ := setupHandler(t, 4)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "",
		Method: MethodCapabilityExec,
	}

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.EqualError(t, err, "request ID cannot be empty")
}

func TestConfidentialRelayHandler_FanOutAndQuorumSuccess(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	h.aggregator = &respondingMockAggregator{}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := json.RawMessage(`{"workflow_id":"wf1","secrets":[{"key":"k","namespace":"ns"}],"enclave_public_key":"pk"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	resultData := json.RawMessage(`{"secrets":[],"master_public_key":"mpk","threshold":1}`)
	response := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  MethodCapabilityExec,
		Result:  &resultData,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.NoError, resp.ErrorCode)
		var jsonResp jsonrpc.Response[json.RawMessage]
		err = json.Unmarshal(resp.RawResponse, &jsonResp)
		assert.NoError(t, err)
		assert.Equal(t, "req-1", jsonResp.ID)
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	err = h.HandleNodeMessage(t.Context(), &response, "0x0000")
	require.NoError(t, err)
	wg.Wait()
}

func TestConfidentialRelayHandler_QuorumWithRealAggregator(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	// Use the real aggregator; DON F=1 so quorum = F+1 = 2
	h.aggregator = &aggregator{}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-quorum",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	resultData := json.RawMessage(`{"payload":"result"}`)
	makeResp := func() *jsonrpc.Response[json.RawMessage] {
		rd := make(json.RawMessage, len(resultData))
		copy(rd, resultData)
		return &jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      "req-quorum",
			Method:  MethodCapabilityExec,
			Result:  &rd,
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.NoError, resp.ErrorCode)
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	// Send 2 matching responses (F+1 = 2)
	for i := range 2 {
		err = h.HandleNodeMessage(t.Context(), makeResp(), fmt.Sprintf("0x%04d", i))
		require.NoError(t, err)
	}
	wg.Wait()
}

func TestConfidentialRelayHandler_QuorumWithDivergentResponses(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	h.aggregator = &aggregator{}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-diverge",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.NoError, resp.ErrorCode)
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	// One divergent response
	divergentResult := json.RawMessage(`{"secrets":[],"master_public_key":"DIFFERENT","threshold":1}`)
	divergentResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-diverge",
		Method:  MethodCapabilityExec,
		Result:  &divergentResult,
	}
	err = h.HandleNodeMessage(t.Context(), divergentResp, "0x0000")
	require.NoError(t, err)

	// Two matching responses (quorum = F+1 = 2)
	matchingResult := json.RawMessage(`{"secrets":[],"master_public_key":"mpk","threshold":1}`)
	for i := 1; i <= 2; i++ {
		rd := make(json.RawMessage, len(matchingResult))
		copy(rd, matchingResult)
		resp := &jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      "req-diverge",
			Method:  MethodCapabilityExec,
			Result:  &rd,
		}
		err = h.HandleNodeMessage(t.Context(), resp, fmt.Sprintf("0x%04d", i))
		require.NoError(t, err)
	}
	wg.Wait()
}

func TestConfidentialRelayHandler_QuorumUnobtainable(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	h.aggregator = &mockAggregator{err: errQuorumUnobtainable}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-unobtainable",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	response := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-unobtainable",
		Method:  MethodCapabilityExec,
		Error: &jsonrpc.WireError{
			Code:    -32603,
			Message: errQuorumUnobtainable.Error(),
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		var jsonResp jsonrpc.Response[json.RawMessage]
		err = json.Unmarshal(resp.RawResponse, &jsonResp)
		assert.NoError(t, err)
		assert.Equal(t, "req-unobtainable", jsonResp.ID)
		assert.NotNil(t, jsonResp.Error)
		assert.Contains(t, jsonResp.Error.Message, "quorum unobtainable")
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	err = h.HandleNodeMessage(t.Context(), &response, "0x0000")
	require.NoError(t, err)
	wg.Wait()
}

func TestConfidentialRelayHandler_RequestTimeout(t *testing.T) {
	h, cb, don, clock := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Use the real aggregator so responses are not immediately satisfied
	h.aggregator = &aggregator{}

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-timeout",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.RequestTimeoutError, resp.ErrorCode)
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	// Advance clock past the request timeout and trigger cleanup
	clock.Advance(31 * time.Second)
	h.removeExpiredRequests(t.Context())
	wg.Wait()
}

func TestConfidentialRelayHandler_DuplicateRequestID(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-dup",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	cb2 := common.NewCallback()
	err = h.HandleJSONRPCUserMessage(t.Context(), req, cb2)
	require.ErrorContains(t, err, "request ID already exists")
}

func TestConfidentialRelayHandler_RateLimitedNode(t *testing.T) {
	handlerConfig := Config{
		RequestTimeoutSec: 30,
	}
	methodConfig, err := json.Marshal(handlerConfig)
	require.NoError(t, err)

	lggr := logger.Test(t)
	don := mocks.NewDON(t)
	donConfig := &config.DONConfig{
		DonId:   "test_relay_don",
		F:       1,
		Members: []config.NodeConfig{nodeOne},
	}
	clock := clockwork.NewFakeClock()
	limitsFactory := limits.Factory{Settings: cresettings.DefaultGetter, Logger: lggr}
	h, err := NewHandler(methodConfig, donConfig, don, lggr, clock, limitsFactory)
	require.NoError(t, err)
	h.aggregator = &respondingMockAggregator{}
	h.globalNodeRateLimiter = limits.GlobalRateLimiter(rate.Limit(100), 100)
	h.perNodeRateLimiters[nodeOne.Address] = limits.GlobalRateLimiter(rate.Limit(0.001), 1)

	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	cb := common.NewCallback()
	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-ratelimit",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	err = h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	resultData := json.RawMessage(`{"secrets":[]}`)
	response := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-ratelimit",
		Method:  MethodCapabilityExec,
		Result:  &resultData,
	}

	// First response from node uses the burst allowance
	err = h.HandleNodeMessage(t.Context(), &response, nodeOne.Address)
	require.NoError(t, err)

	// Verify callback was called
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	resp, err := cb.Wait(ctx)
	require.NoError(t, err)
	assert.Equal(t, api.NoError, resp.ErrorCode)

	// Start a new request
	cb2 := common.NewCallback()
	req2 := jsonrpc.Request[json.RawMessage]{
		ID:     "req-ratelimit-2",
		Method: MethodCapabilityExec,
		Params: &params,
	}
	err = h.HandleJSONRPCUserMessage(t.Context(), req2, cb2)
	require.NoError(t, err)

	response2 := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-ratelimit-2",
		Method:  MethodCapabilityExec,
		Result:  &resultData,
	}

	// Second response should be rate limited (silently dropped)
	err = h.HandleNodeMessage(t.Context(), &response2, nodeOne.Address)
	require.NoError(t, err)

	// Callback should NOT be called - verify with timeout
	ctx2, cancel2 := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel2()
	_, err = cb2.Wait(ctx2)
	require.Error(t, err) // Should timeout
}

func TestConfidentialRelayHandler_LateNodeResponse(t *testing.T) {
	h, cb, _, _ := setupHandler(t, 4)

	resultData := json.RawMessage(`{"secrets":[]}`)
	staleResponse := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "nonexistent-request",
		Method:  MethodCapabilityExec,
		Result:  &resultData,
	}

	// This should not error, just silently ignore
	err := h.HandleNodeMessage(t.Context(), &staleResponse, "0x0000")
	require.NoError(t, err)

	// Verify callback was not triggered
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	_, err = cb.Wait(ctx)
	require.Error(t, err)
}

func TestConfidentialRelayHandler_AllNodesFanOutFail(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("connection refused"))

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-allfail",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.FatalError, resp.ErrorCode)
		var jsonResp jsonrpc.Response[json.RawMessage]
		err = json.Unmarshal(resp.RawResponse, &jsonResp)
		assert.NoError(t, err)
		assert.Contains(t, jsonResp.Error.Message, "failed to forward user request to nodes")
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)
	wg.Wait()
}

func TestConfidentialRelayHandler_FanOutWaitsWhileQuorumStillPossible(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(
		func(_ context.Context, nodeAddress string, _ *jsonrpc.Request[json.RawMessage]) error {
			switch nodeAddress {
			case "0x0000", "0x0001":
				return errors.New("connection refused")
			default:
				return nil
			}
		},
	)

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-still-possible",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	require.NotNil(t, h.getActiveRequest(req.ID), "request should remain active while quorum is still possible")

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	_, err = cb.Wait(ctx)
	require.Error(t, err)
}

func TestConfidentialRelayHandler_FanOutFailsWhenQuorumBecomesImpossible(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(
		func(_ context.Context, nodeAddress string, _ *jsonrpc.Request[json.RawMessage]) error {
			switch nodeAddress {
			case "0x0000", "0x0001", "0x0002":
				return errors.New("connection refused")
			default:
				return nil
			}
		},
	)

	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-quorum-impossible",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.FatalError, resp.ErrorCode)
		var jsonResp jsonrpc.Response[json.RawMessage]
		err = json.Unmarshal(resp.RawResponse, &jsonResp)
		assert.NoError(t, err)
		assert.Contains(t, jsonResp.Error.Message, "failed to forward user request to nodes")
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)
	wg.Wait()

	require.Nil(t, h.getActiveRequest(req.ID), "request should be cleaned up once quorum is impossible")
}

func TestConfidentialRelayHandler_FanOutToNodes_IsConcurrent(t *testing.T) {
	lggr := logger.Test(t)
	don := newBarrierDON(2)
	donConfig := &config.DONConfig{
		DonId: "test_relay_don",
		F:     1,
		Members: []config.NodeConfig{
			{Name: "node0", Address: "0x0000"},
			{Name: "node1", Address: "0x0001"},
		},
	}

	methodConfig, err := json.Marshal(Config{
		RequestTimeoutSec: 30,
	})
	require.NoError(t, err)

	limitsFactory := limits.Factory{Settings: cresettings.DefaultGetter, Logger: lggr}
	h, err := NewHandler(methodConfig, donConfig, don, lggr, clockwork.NewFakeClock(), limitsFactory)
	require.NoError(t, err)

	cb := common.NewCallback()
	params := json.RawMessage(`{"workflow_id":"wf1"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-concurrent-fanout",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- h.HandleJSONRPCUserMessage(ctx, req, cb)
	}()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		cancel()
		<-done
		t.Fatal("HandleJSONRPCUserMessage did not fan out to nodes concurrently")
	}

	don.mu.Lock()
	started := don.started
	don.mu.Unlock()
	assert.Equal(t, 2, started)
}

func TestConfidentialRelayHandler_CapabilityExecMethod(t *testing.T) {
	h, cb, don, _ := setupHandler(t, 4)
	h.aggregator = &respondingMockAggregator{}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := json.RawMessage(`{"workflow_id":"wf1","capability_id":"cap1","payload":"data"}`)
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-cap",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.NoError, resp.ErrorCode)
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	resultData := json.RawMessage(`{"payload":"result"}`)
	response := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-cap",
		Method:  MethodCapabilityExec,
		Result:  &resultData,
	}
	err = h.HandleNodeMessage(t.Context(), &response, "0x0000")
	require.NoError(t, err)
	wg.Wait()
	don.AssertCalled(t, "SendToNode", mock.Anything, mock.Anything, mock.Anything)
}
