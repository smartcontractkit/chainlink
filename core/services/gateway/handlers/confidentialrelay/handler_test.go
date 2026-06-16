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

	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
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

// validCapParamsJSON returns JSON-RPC params for MethodCapabilityExec that satisfy
// chainlink-common's CapabilityRequestParams.Validate (testOwner / testExecutionID
// constants live in bundler_test.go, same package).
func validCapParamsJSON(workflowID string) json.RawMessage {
	raw, err := json.Marshal(validCapParams(workflowID))
	if err != nil {
		panic(err)
	}
	return json.RawMessage(raw)
}

func setupHandler(t *testing.T, numNodes int) (*handler, *common.Callback, *mocks.DON, *clockwork.FakeClock) {
	t.Helper()
	return setupHandlerWithF(t, numNodes, 1)
}

func setupHandlerWithF(t *testing.T, numNodes, f int) (*handler, *common.Callback, *mocks.DON, *clockwork.FakeClock) {
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
		F:       f,
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
	// The default bundler is the real one; tests that need a failure inject a mock.
	cb := common.NewCallback()
	return h, cb, don, clock
}

// mockBundler lets a test force the bundler error path.
type mockBundler struct {
	err error
}

func (m *mockBundler) Bundle(_ jsonrpc.Request[json.RawMessage], _ map[string]jsonrpc.Response[json.RawMessage], _ logger.Logger) (*jsonrpc.Response[json.RawMessage], int, error) {
	return nil, 0, m.err
}

func TestConfidentialRelayHandler_Methods(t *testing.T) {
	t.Parallel()
	h, _, _, _ := setupHandler(t, 4)
	methods := h.Methods()
	assert.Equal(t, []string{MethodSecretsGet, MethodCapabilityExec}, methods)
}

func TestConfidentialRelayHandler_HandleLegacyUserMessage(t *testing.T) {
	t.Parallel()
	h, cb, _, _ := setupHandler(t, 4)
	err := h.HandleLegacyUserMessage(t.Context(), nil, cb)
	require.ErrorContains(t, err, "confidential relay handler does not support legacy messages")
}

func TestConfidentialRelayHandler_RequestIDTooLong(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	h, cb, _, _ := setupHandler(t, 4)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "",
		Method: MethodCapabilityExec,
	}

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.EqualError(t, err, "request ID cannot be empty")
}

// At F=1 the gateway forwards once it has collected 2F+1 = 3 responses, so that
// (under <=F faulty) at least F+1 honest matching responses are guaranteed present
// for the enclave to verify. The gateway itself makes no signature/quorum decision.
func TestConfidentialRelayHandler_ForwardsBundleAtQuorum(t *testing.T) {
	t.Parallel()
	h, cb, don, _ := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := validCapParamsJSON("wf1")
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-quorum",
		Method: MethodCapabilityExec,
		Params: &params,
	}
	result := relaytypes.CapabilityResponseResult{Payload: "result"}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := cb.Wait(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, api.NoError, resp.ErrorCode)
		var jsonResp jsonrpc.Response[json.RawMessage]
		assert.NoError(t, json.Unmarshal(resp.RawResponse, &jsonResp))
		assert.NotNil(t, jsonResp.Result)
		var bundle relaytypes.SignedCapabilityResponseBundle
		assert.NoError(t, json.Unmarshal(*jsonResp.Result, &bundle))
		assert.Len(t, bundle.Responses, 3, "the gateway forwards every collected response")
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	// First two responses do not reach the 2F+1 threshold: no callback yet.
	for i := range 2 {
		err = h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-quorum", result, fmt.Appendf(nil, "signer-%d", i)), fmt.Sprintf("0x%04d", i))
		require.NoError(t, err)
	}
	require.NotNil(t, h.getActiveRequest(req.ID), "request stays active below the 2F+1 threshold")

	// Third response reaches 2F+1 and triggers the forward.
	err = h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-quorum", result, []byte("signer-2")), "0x0002")
	require.NoError(t, err)
	wg.Wait()
}

// The gateway does not resolve divergent results; it forwards them all and lets the
// enclave verify signatures and pick the result backed by F+1 valid signers.
func TestConfidentialRelayHandler_ForwardsAllDivergentResponses(t *testing.T) {
	t.Parallel()
	h, cb, don, _ := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := validCapParamsJSON("wf1")
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
		var jsonResp jsonrpc.Response[json.RawMessage]
		assert.NoError(t, json.Unmarshal(resp.RawResponse, &jsonResp))
		var bundle relaytypes.SignedCapabilityResponseBundle
		assert.NoError(t, json.Unmarshal(*jsonResp.Result, &bundle))
		assert.Len(t, bundle.Responses, 3, "divergent and matching responses are all forwarded untouched")
	}()

	err := h.HandleJSONRPCUserMessage(t.Context(), req, cb)
	require.NoError(t, err)

	divergent := relaytypes.CapabilityResponseResult{Payload: "DIVERGENT"}
	match := relaytypes.CapabilityResponseResult{Payload: "match"}
	err = h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-diverge", divergent, []byte("signer-0")), "0x0000")
	require.NoError(t, err)
	err = h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-diverge", match, []byte("signer-1")), "0x0001")
	require.NoError(t, err)
	err = h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-diverge", match, []byte("signer-2")), "0x0002")
	require.NoError(t, err)
	wg.Wait()
}

func TestConfidentialRelayHandler_BundlerErrorReturnsFatal(t *testing.T) {
	t.Parallel()
	h, cb, don, _ := setupHandler(t, 4)
	h.bundler = &mockBundler{err: errors.New("boom")}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := validCapParamsJSON("wf1")
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-bundler-err",
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
	}()

	require.NoError(t, h.HandleJSONRPCUserMessage(t.Context(), req, cb))
	result := relaytypes.CapabilityResponseResult{Payload: "x"}
	for i := range 3 {
		err := h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-bundler-err", result, fmt.Appendf(nil, "s-%d", i)), fmt.Sprintf("0x%04d", i))
		require.NoError(t, err)
	}
	wg.Wait()
}

// On timeout with fewer than 2F+1 but at least F+1 responses, the gateway forwards
// what it collected (the bundle may still carry F+1 valid signatures, e.g. if faulty
// nodes stayed silent) and lets the enclave decide. Here F=1, so 2 responses meet the
// F+1 floor.
func TestConfidentialRelayHandler_TimeoutForwardsPartialBundle(t *testing.T) {
	t.Parallel()
	h, cb, don, clock := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := validCapParamsJSON("wf1")
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-partial",
		Method: MethodCapabilityExec,
		Params: &params,
	}

	require.NoError(t, h.HandleJSONRPCUserMessage(t.Context(), req, cb))
	result := relaytypes.CapabilityResponseResult{Payload: "x"}
	for i := range 2 { // below 2F+1=3
		err := h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-partial", result, fmt.Appendf(nil, "s-%d", i)), fmt.Sprintf("0x%04d", i))
		require.NoError(t, err)
	}
	require.NotNil(t, h.getActiveRequest(req.ID))

	// The expiry sweep delivers the partial bundle to the callback synchronously,
	// so we can read it on the main goroutine afterwards.
	clock.Advance(31 * time.Second)
	h.removeExpiredRequests(t.Context())

	resp, err := cb.Wait(t.Context())
	require.NoError(t, err)
	require.Equal(t, api.NoError, resp.ErrorCode)
	var jsonResp jsonrpc.Response[json.RawMessage]
	require.NoError(t, json.Unmarshal(resp.RawResponse, &jsonResp))
	var bundle relaytypes.SignedCapabilityResponseBundle
	require.NoError(t, json.Unmarshal(*jsonResp.Result, &bundle))
	require.Len(t, bundle.Responses, 2)
}

// On timeout with fewer than F+1 signed responses, the enclave could never reach
// its F+1-valid-signature quorum, so the gateway skips the futile forward and
// returns a timeout error directly. Here F=1, so a single response is below the
// floor.
func TestConfidentialRelayHandler_TimeoutBelowQuorumFloorReturnsTimeout(t *testing.T) {
	t.Parallel()
	h, cb, don, clock := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := validCapParamsJSON("wf1")
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-below-floor",
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

	require.NoError(t, h.HandleJSONRPCUserMessage(t.Context(), req, cb))
	result := relaytypes.CapabilityResponseResult{Payload: "x"}
	// One response, below F+1=2.
	require.NoError(t, h.HandleNodeMessage(t.Context(), capExecSignedRespPtr(t, "req-below-floor", result, []byte("s-0")), "0x0000"))
	require.NotNil(t, h.getActiveRequest(req.ID))

	clock.Advance(31 * time.Second)
	h.removeExpiredRequests(t.Context())
	wg.Wait()
}

func TestConfidentialRelayHandler_TimeoutNoResponses(t *testing.T) {
	t.Parallel()
	h, cb, don, clock := setupHandler(t, 4)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	params := validCapParamsJSON("wf1")
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

	clock.Advance(31 * time.Second)
	h.removeExpiredRequests(t.Context())
	wg.Wait()
}

func TestConfidentialRelayHandler_DuplicateRequestID(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	handlerConfig := Config{
		RequestTimeoutSec: 30,
	}
	methodConfig, err := json.Marshal(handlerConfig)
	require.NoError(t, err)

	lggr := logger.Test(t)
	don := mocks.NewDON(t)
	// F=0 so the forward threshold (2F+1) is 1: a single response from the one-node
	// DON forwards immediately, isolating the rate-limit behavior under test.
	donConfig := &config.DONConfig{
		DonId:   "test_relay_don",
		F:       0,
		Members: []config.NodeConfig{nodeOne},
	}
	clock := clockwork.NewFakeClock()
	limitsFactory := limits.Factory{Settings: cresettings.DefaultGetter, Logger: lggr}
	h, err := NewHandler(methodConfig, donConfig, don, lggr, clock, limitsFactory)
	require.NoError(t, err)
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

	resultData := json.RawMessage(`{"result":{"payload":"r"},"signature":{}}`)
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
	t.Parallel()
	h, cb, _, _ := setupHandler(t, 4)

	resultData := json.RawMessage(`{"result":{},"signature":{}}`)
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func capExecSignedRespPtr(t *testing.T, id string, result relaytypes.CapabilityResponseResult, signer []byte) *jsonrpc.Response[json.RawMessage] {
	r := capExecSignedResponse(t, id, result, signer)
	return &r
}
