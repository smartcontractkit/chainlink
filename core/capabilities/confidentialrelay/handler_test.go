package confidentialrelay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func makeCapabilityPayload(t *testing.T, inputs map[string]any) string {
	t.Helper()
	wrapped, err := values.Wrap(inputs)
	require.NoError(t, err)
	payload, err := anypb.New(values.Proto(wrapped))
	require.NoError(t, err)
	sdkReq := &sdkpb.CapabilityRequest{
		Id:      "my-cap@1.0.0",
		Payload: payload,
		Method:  "Execute",
	}
	b, err := proto.Marshal(sdkReq)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(b)
}

const testAttestationB64 = "ZHVtbXktYXR0ZXN0YXRpb24=" // base64("dummy-attestation")

func noopValidator(_ []byte, _, _ []byte) error { return nil }

type mockGatewayConnector struct {
	core.UnimplementedGatewayConnector
	lastResp     *jsonrpc.Response[json.RawMessage]
	addedMethods []string
	removed      bool
}

func (m *mockGatewayConnector) SendToGateway(_ context.Context, _ string, resp *jsonrpc.Response[json.RawMessage]) error {
	m.lastResp = resp
	return nil
}
func (m *mockGatewayConnector) AddHandler(_ context.Context, methods []string, _ core.GatewayConnectorHandler) error {
	m.addedMethods = methods
	return nil
}
func (m *mockGatewayConnector) RemoveHandler(_ context.Context, _ []string) error {
	m.removed = true
	return nil
}

type mockExecutable struct {
	infoResult  capabilities.CapabilityInfo
	infoErr     error
	execResult  capabilities.CapabilityResponse
	execErr     error
	lastRequest *capabilities.CapabilityRequest
}

func (m *mockExecutable) Info(_ context.Context) (capabilities.CapabilityInfo, error) {
	return m.infoResult, m.infoErr
}
func (m *mockExecutable) Execute(_ context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	m.lastRequest = &req
	return m.execResult, m.execErr
}
func (m *mockExecutable) RegisterToWorkflow(_ context.Context, _ capabilities.RegisterToWorkflowRequest) error {
	return nil
}
func (m *mockExecutable) UnregisterFromWorkflow(_ context.Context, _ capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

type mockCapRegistry struct {
	core.UnimplementedCapabilitiesRegistry
	executables map[string]*mockExecutable
	configs     map[string]capabilities.CapabilityConfiguration
	dons        map[string][]capabilities.DONWithNodes
	localNode   capabilities.Node
}

func (m *mockCapRegistry) GetExecutable(_ context.Context, id string) (capabilities.ExecutableCapability, error) {
	if exec, ok := m.executables[id]; ok {
		return exec, nil
	}
	return nil, fmt.Errorf("capability not found: %s", id)
}
func (m *mockCapRegistry) ConfigForCapability(_ context.Context, capID string, _ uint32) (capabilities.CapabilityConfiguration, error) {
	if cfg, ok := m.configs[capID]; ok {
		return cfg, nil
	}
	return capabilities.CapabilityConfiguration{}, fmt.Errorf("config not found: %s", capID)
}
func (m *mockCapRegistry) DONsForCapability(_ context.Context, capID string) ([]capabilities.DONWithNodes, error) {
	if dons, ok := m.dons[capID]; ok {
		return dons, nil
	}
	return nil, fmt.Errorf("no DONs found for: %s", capID)
}
func (m *mockCapRegistry) LocalNode(_ context.Context) (capabilities.Node, error) {
	return m.localNode, nil
}

func newTestHandler(t *testing.T, registry core.CapabilitiesRegistry, gwConn core.GatewayConnector) *Handler {
	t.Helper()
	lggr, err := logger.New()
	require.NoError(t, err)
	h, err := NewHandler(registry, gwConn, lggr)
	require.NoError(t, err)
	h.validateAttestation = noopValidator
	return h
}

// withEnclaveConfig adds the default confidential-workflows enclave config
// to a mock registry so getEnclaveAttestationConfig succeeds during tests.
func withEnclaveConfig(reg *mockCapRegistry) *mockCapRegistry {
	enclaveConfig := enclavesList{
		Enclaves: []enclaveEntry{{TrustedValues: []json.RawMessage{json.RawMessage(`{}`)}}},
	}
	wrapped, _ := values.WrapMap(enclaveConfig)
	if reg.configs == nil {
		reg.configs = map[string]capabilities.CapabilityConfiguration{}
	}
	reg.configs[confidentialWorkflowsCapID] = capabilities.CapabilityConfiguration{
		DefaultConfig: wrapped,
	}
	if reg.dons == nil {
		reg.dons = map[string][]capabilities.DONWithNodes{}
	}
	reg.dons[confidentialWorkflowsCapID] = []capabilities.DONWithNodes{
		{DON: capabilities.DON{ID: 1}},
	}
	return reg
}

func makeRequest(t *testing.T, method string, params any) *jsonrpc.Request[json.RawMessage] {
	t.Helper()
	b, err := json.Marshal(params)
	require.NoError(t, err)
	raw := json.RawMessage(b)
	return &jsonrpc.Request[json.RawMessage]{
		Method: method,
		ID:     "req-1",
		Params: &raw,
	}
}

func TestHandler_HandleGatewayMessage(t *testing.T) {
	tests := []struct {
		name            string
		registry        func(t *testing.T) *mockCapRegistry
		req             func(t *testing.T) *jsonrpc.Request[json.RawMessage]
		checkResp       func(t *testing.T, resp *jsonrpc.Response[json.RawMessage])
		checkExecutable func(t *testing.T, reg *mockCapRegistry)
	}{
		{
			name: "capability execute success",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{
					executables: map[string]*mockExecutable{
						"my-cap@1.0.0": {
							execResult: capabilities.CapabilityResponse{
								Payload: &anypb.Any{Value: []byte("result-proto-bytes")},
							},
						},
					},
				})
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					Owner:        "0xowner",
					ExecutionID:  "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:  "17",
					CapabilityID: "my-cap@1.0.0",
					Payload:      makeCapabilityPayload(t, map[string]any{"key": "val"}),
					Attestation:  testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				var result confidentialrelaytypes.CapabilityResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				decoded, err := base64.StdEncoding.DecodeString(result.Payload)
				require.NoError(t, err)
				var capResp sdkpb.CapabilityResponse
				require.NoError(t, proto.Unmarshal(decoded, &capResp))
				require.NotNil(t, capResp.GetPayload())
				assert.Equal(t, "result-proto-bytes", string(capResp.GetPayload().GetValue()))
				assert.Empty(t, result.Error)
			},
			checkExecutable: func(t *testing.T, reg *mockCapRegistry) {
				exec := reg.executables["my-cap@1.0.0"]
				require.NotNil(t, exec.lastRequest, "Execute should have been called")
				assert.Equal(t, "wf-1", exec.lastRequest.Metadata.WorkflowID)
				assert.Equal(t, "0xowner", exec.lastRequest.Metadata.WorkflowOwner)
				assert.Equal(t, "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1", exec.lastRequest.Metadata.WorkflowExecutionID)
				assert.Equal(t, "17", exec.lastRequest.Metadata.ReferenceID)
			},
		},
		{
			name: "capability execute sets Inputs from Payload for backward compat",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{
					executables: map[string]*mockExecutable{
						"my-cap@1.0.0": {
							execResult: capabilities.CapabilityResponse{},
						},
					},
				})
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					CapabilityID: "my-cap@1.0.0",
					Payload:      makeCapabilityPayload(t, map[string]any{"echo": "hello"}),
					Attestation:  testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
			},
			checkExecutable: func(t *testing.T, reg *mockCapRegistry) {
				exec := reg.executables["my-cap@1.0.0"]
				require.NotNil(t, exec.lastRequest, "Execute should have been called")
				require.NotNil(t, exec.lastRequest.Payload)
				var valPB valuespb.Value
				require.NoError(t, exec.lastRequest.Payload.UnmarshalTo(&valPB))
				require.NotNil(t, exec.lastRequest.Inputs)
				unwrapped, err := exec.lastRequest.Inputs.Unwrap()
				require.NoError(t, err)
				m, ok := unwrapped.(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "hello", m["echo"])
			},
		},
		{
			name: "capability execute attestation failure",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					CapabilityID: "my-cap@1.0.0",
					Payload:      base64.StdEncoding.EncodeToString([]byte("payload")),
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.NotNil(t, resp.Error)
				assert.Equal(t, jsonrpc.ErrInternal, resp.Error.Code)
			},
		},
		{
			name: "capability execute not found",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{executables: map[string]*mockExecutable{}})
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					CapabilityID: "missing-cap@1.0.0",
					Payload:      base64.StdEncoding.EncodeToString([]byte("payload")),
					Attestation:  testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.NotNil(t, resp.Error)
				assert.Equal(t, jsonrpc.ErrInternal, resp.Error.Code)
				assert.Equal(t, internalErrorMessage, resp.Error.Message)
			},
		},
		{
			name: "capability execute error returned in result",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{
					executables: map[string]*mockExecutable{
						"fail-cap@1.0.0": {execErr: errors.New("execution failed")},
					},
				})
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				sdkReq := &sdkpb.CapabilityRequest{Id: "fail-cap@1.0.0", Method: "Execute"}
				b, err := proto.Marshal(sdkReq)
				require.NoError(t, err)
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					CapabilityID: "fail-cap@1.0.0",
					Payload:      base64.StdEncoding.EncodeToString(b),
					Attestation:  testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				var result confidentialrelaytypes.CapabilityResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				assert.Equal(t, "execution failed", result.Error)
				assert.Empty(t, result.Payload)
			},
		},
		{
			name: "unsupported method",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, "unknown.method", nil)
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.NotNil(t, resp.Error)
				assert.Equal(t, jsonrpc.ErrMethodNotFound, resp.Error.Code)
			},
		},
		{
			name: "invalid params JSON",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			req: func(_ *testing.T) *jsonrpc.Request[json.RawMessage] {
				raw := json.RawMessage([]byte(`{invalid json`))
				return &jsonrpc.Request[json.RawMessage]{
					Method: confidentialrelaytypes.MethodCapabilityExec,
					ID:     "req-1",
					Params: &raw,
				}
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.NotNil(t, resp.Error)
				assert.Equal(t, jsonrpc.ErrInvalidParams, resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gwConn := &mockGatewayConnector{}
			reg := tt.registry(t)
			h := newTestHandler(t, reg, gwConn)
			err := h.HandleGatewayMessage(t.Context(), "gw-1", tt.req(t))
			require.NoError(t, err)
			require.NotNil(t, gwConn.lastResp)
			tt.checkResp(t, gwConn.lastResp)
			if tt.checkExecutable != nil {
				tt.checkExecutable(t, reg)
			}
		})
	}
}

func TestHandler_Lifecycle(t *testing.T) {
	gwConn := &mockGatewayConnector{}
	h := newTestHandler(t, &mockCapRegistry{}, gwConn)

	t.Run("start registers handler", func(t *testing.T) {
		require.NoError(t, h.Start(t.Context()))
		assert.Equal(t, h.Methods(), gwConn.addedMethods)
	})

	t.Run("close removes handler", func(t *testing.T) {
		require.NoError(t, h.Close())
		assert.True(t, gwConn.removed)
	})

	t.Run("ID returns handler name", func(t *testing.T) {
		id, err := h.ID(t.Context())
		require.NoError(t, err)
		assert.Equal(t, HandlerName, id)
	})
}
