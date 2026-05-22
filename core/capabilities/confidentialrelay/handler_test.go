package confidentialrelay

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vault "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	vaulttypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
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
	key, err := p2pkey.NewV2()
	require.NoError(t, err)
	h, err := NewHandler(registry, gwConn, newRelayResponseSigner(key), lggr, limits.Factory{Logger: lggr})
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

// secretsGetTestRegistry builds a mock registry with a vault executable that
// returns a valid GetSecretsResponse for the "API_KEY" secret.
func secretsGetTestRegistry(t *testing.T) *mockCapRegistry {
	t.Helper()
	// Must match secretsGetTestParams().EnclavePublicKey and pass the
	// confidentialrelay.validateEnclavePublicKey hex check (#2032).
	enclaveKey := "aabbcc"
	vaultResp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: &vault.SecretIdentifier{
					Key:       "API_KEY",
					Namespace: vaulttypes.DefaultNamespace,
					Owner:     "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
				},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: hex.EncodeToString([]byte("encrypted-value")),
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: enclaveKey,
								Shares:        []string{hex.EncodeToString([]byte("share-1"))},
							},
						},
					},
				},
			},
		},
	}
	payload, err := anypb.New(vaultResp)
	require.NoError(t, err)

	return withEnclaveConfig(&mockCapRegistry{
		executables: map[string]*mockExecutable{
			vault.CapabilityID: {
				execResult: capabilities.CapabilityResponse{Payload: payload},
			},
		},
		localNode: capabilities.Node{
			WorkflowDON: capabilities.DON{ID: 42, ConfigVersion: 7},
		},
	})
}

// testOwner is a 0x-prefixed 20-byte hex address that satisfies
// chainlink-common's confidentialrelay.validateOwnerAddress.
const testOwner = "0x0000000000000000000000000000000000000001"

// secretsGetTestRequest builds a secrets-get request with a known owner and org ID.
//
// Field formats are pinned by chainlink-common's confidentialrelay.Validate (introduced
// in chainlink-common#2032): Owner must be a 0x-prefixed 20-byte hex address,
// ExecutionID must be 32-byte hex (64 chars, no prefix), EnclavePublicKey must be
// hex-encoded, and Secrets entries need both Key and Namespace. Hash() returns an
// error if any of these are missing or malformed, so signing breaks if the fixture
// drifts.
func secretsGetTestRequest(t *testing.T) *jsonrpc.Request[json.RawMessage] {
	t.Helper()
	return makeRequest(t, confidentialrelaytypes.MethodSecretsGet, secretsGetTestParams())
}

// secretsGetTestParams returns the canonical valid params used by both the request
// builder and the response-signature-verification step.
func secretsGetTestParams() confidentialrelaytypes.SecretsRequestParams {
	return confidentialrelaytypes.SecretsRequestParams{
		WorkflowID:       "wf-secrets-1",
		Owner:            "0xab5801a7d398351b8be11c439e05c5b3259aec9b", // lowercase, should be normalized
		ExecutionID:      "0000000000000000000000000000000000000000000000000000000000000001",
		OrgID:            "org-123",
		EnclavePublicKey: "aabbcc",
		Secrets: []confidentialrelaytypes.SecretIdentifier{
			{Key: "API_KEY", Namespace: "main"},
		},
		Attestation: testAttestationB64,
	}
}

func TestHandler_HandleGatewayMessage(t *testing.T) {
	tests := []struct {
		name            string
		registry        func(t *testing.T) *mockCapRegistry
		req             func(t *testing.T) *jsonrpc.Request[json.RawMessage]
		modifyHandler   func(t *testing.T, h *Handler)
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
					Owner:        testOwner, // chainlink-common#2032 requires 0x-prefixed 20-byte hex
					ExecutionID:  "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:  "17",
					CapabilityID: "my-cap@1.0.0",
					Payload:      makeCapabilityPayload(t, map[string]any{"key": "val"}),
					Attestation:  testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				params := confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					Owner:        testOwner,
					ExecutionID:  "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:  "17",
					CapabilityID: "my-cap@1.0.0",
					Payload:      makeCapabilityPayload(t, map[string]any{"key": "val"}),
				}
				var result confidentialrelaytypes.SignedCapabilityResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				require.Len(t, result.Signatures, 1)
				assertValidCapabilitySignature(t, params, result)

				decoded, err := base64.StdEncoding.DecodeString(result.Result.Payload)
				require.NoError(t, err)
				var capResp sdkpb.CapabilityResponse
				require.NoError(t, proto.Unmarshal(decoded, &capResp))
				require.NotNil(t, capResp.GetPayload())
				assert.Equal(t, "result-proto-bytes", string(capResp.GetPayload().GetValue()))
				assert.Empty(t, result.Result.Error)
			},
			checkExecutable: func(t *testing.T, reg *mockCapRegistry) {
				exec := reg.executables["my-cap@1.0.0"]
				require.NotNil(t, exec.lastRequest, "Execute should have been called")
				assert.Equal(t, "wf-1", exec.lastRequest.Metadata.WorkflowID)
				assert.Equal(t, testOwner, exec.lastRequest.Metadata.WorkflowOwner)
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
					Owner:        testOwner,
					ExecutionID:  "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:  "17",
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
					Owner:        testOwner,
					ExecutionID:  "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:  "17",
					CapabilityID: "fail-cap@1.0.0",
					Payload:      base64.StdEncoding.EncodeToString(b),
					Attestation:  testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				params := confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:   "wf-1",
					Owner:        testOwner,
					ExecutionID:  "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:  "17",
					CapabilityID: "fail-cap@1.0.0",
					Payload:      base64.StdEncoding.EncodeToString(mustMarshalProto(t, &sdkpb.CapabilityRequest{Id: "fail-cap@1.0.0", Method: "Execute"})),
				}
				var result confidentialrelaytypes.SignedCapabilityResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				require.Len(t, result.Signatures, 1)
				assertValidCapabilitySignature(t, params, result)
				assert.Equal(t, "execution failed", result.Result.Error)
				assert.Empty(t, result.Result.Payload)
			},
		},
		{
			name:     "secrets get invokes vault execute with stable capability metadata",
			registry: secretsGetTestRegistry,
			req:      secretsGetTestRequest,
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				// signSecretsResponse hashes against the request params (no Attestation),
				// so we drop it here for the signature check.
				params := secretsGetTestParams()
				params.Attestation = ""
				var result confidentialrelaytypes.SignedSecretsResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				require.Len(t, result.Signatures, 1)
				assertValidSecretsSignature(t, params, result)
				require.Len(t, result.Result.Secrets, 1)
				assert.Equal(t, "API_KEY", result.Result.Secrets[0].ID.Key)
			},
			checkExecutable: func(t *testing.T, reg *mockCapRegistry) {
				exec := reg.executables[vault.CapabilityID]
				require.NotNil(t, exec.lastRequest, "vault Execute should have been called")

				var vaultReq vault.GetSecretsRequest
				require.NoError(t, exec.lastRequest.Payload.UnmarshalTo(&vaultReq))
				require.Len(t, vaultReq.Requests, 1)
				assert.Equal(t, "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", vaultReq.Requests[0].Id.Owner)
				assert.Equal(t, "0xab5801a7d398351b8be11c439e05c5b3259aec9b", exec.lastRequest.Metadata.WorkflowOwner)
				assert.Equal(t, "wf-secrets-1", exec.lastRequest.Metadata.WorkflowID)
				assert.Equal(t, uint32(42), exec.lastRequest.Metadata.WorkflowDonID)
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
			if tt.modifyHandler != nil {
				tt.modifyHandler(t, h)
			}
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

func mustMarshalProto(t *testing.T, msg proto.Message) []byte {
	t.Helper()
	b, err := proto.Marshal(msg)
	require.NoError(t, err)
	return b
}

func assertValidCapabilitySignature(
	t *testing.T,
	params confidentialrelaytypes.CapabilityRequestParams,
	result confidentialrelaytypes.SignedCapabilityResponseResult,
) {
	t.Helper()
	hash, err := result.Result.Hash(params)
	require.NoError(t, err)
	payload := confidentialrelaytypes.RelayResponseSignaturePayload(hash)
	pubKey := ed25519.PublicKey(result.Signatures[0].Signer)
	require.True(t, ed25519.Verify(pubKey, payload, result.Signatures[0].Signature))
}

func assertValidSecretsSignature(
	t *testing.T,
	params confidentialrelaytypes.SecretsRequestParams,
	result confidentialrelaytypes.SignedSecretsResponseResult,
) {
	t.Helper()
	hash, err := result.Result.Hash(params)
	require.NoError(t, err)
	payload := confidentialrelaytypes.RelayResponseSignaturePayload(hash)
	pubKey := ed25519.PublicKey(result.Signatures[0].Signer)
	require.True(t, ed25519.Verify(pubKey, payload, result.Signatures[0].Signature))
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
