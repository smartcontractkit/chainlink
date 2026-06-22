package confidentialrelay

import (
	"bytes"
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
	confidentialworkflow "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	vaulttypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
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
	// Wire WorkflowDON membership to match testEnclaveConfig so the relay-side
	// verifyEnclaveConfigMatchesDON check passes for fixtures that build
	// request params with testEnclaveConfig.
	reg.localNode.WorkflowDON.Members = testWorkflowDONMembers()
	reg.localNode.WorkflowDON.F = testEnclaveF
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

// make32Byte builds a 32-byte slice filled with the given byte, used by the
// enclave-config mismatch tests.
func make32Byte(b byte) []byte {
	s := make([]byte, 32)
	for i := range s {
		s[i] = b
	}
	return s
}

// testEnclaveF is the DON fault tolerance used across these tests. Untyped so it
// assigns cleanly to both EnclaveConfig.F (uint32) and WorkflowDON.F (uint8) without
// a narrowing conversion that would trip gosec G115.
const testEnclaveF = 1

// testWorkflowDONKeys returns the deterministic ed25519 keys backing the test Workflow
// DON. The public keys serve double duty: as EnclaveConfig.Signers (PRIV-458: the relay
// checks the enclave's reported signers match the onchain DON) and as the Workflow DON
// members whose signatures over the compute request the relay verifies (PRIV-433). They
// must therefore be real ed25519 public keys, not arbitrary bytes.
func testWorkflowDONKeys() ([]ed25519.PrivateKey, [][]byte) {
	const n = 4
	privs := make([]ed25519.PrivateKey, n)
	pubs := make([][]byte, n)
	for i := range privs {
		privs[i] = ed25519.NewKeyFromSeed(bytes.Repeat([]byte{byte(0x07 + i)}, ed25519.SeedSize))
		pubs[i] = append([]byte(nil), privs[i].Public().(ed25519.PublicKey)...)
	}
	return privs, pubs
}

// testEnclaveConfig is the canonical EnclaveConfig that handler tests put on outgoing
// request params. withEnclaveConfig wires the matching WorkflowDON membership into the
// mock CapabilitiesRegistry so verifyEnclaveConfigMatchesDON accepts requests built with
// it. Signers are real ed25519 public keys so they also back the PRIV-433 signature check.
func testEnclaveConfig() confidentialrelaytypes.EnclaveConfig {
	_, pubs := testWorkflowDONKeys()
	return confidentialrelaytypes.EnclaveConfig{
		Signers:         pubs,
		MasterPublicKey: []byte("test-master-public-key"),
		T:               3,
		F:               testEnclaveF,
	}
}

func testEnclaveConfigPtr() *confidentialrelaytypes.EnclaveConfig {
	c := testEnclaveConfig()
	return &c
}

// testWorkflowDONMembers returns []p2ptypes.PeerID whose [:] slices match
// testEnclaveConfig().Signers byte-for-byte.
func testWorkflowDONMembers() []p2ptypes.PeerID {
	cfg := testEnclaveConfig()
	members := make([]p2ptypes.PeerID, len(cfg.Signers))
	for i, s := range cfg.Signers {
		var pid p2ptypes.PeerID
		copy(pid[:], s)
		members[i] = pid
	}
	return members
}

// signedComputeRequestsForParams builds the Workflow-DON-signed compute requests the
// enclave forwards to the relay. It signs with 2*F+1 of the DON keys (the quorum
// verifyWorkflowAuthorization requires). PublicData carries the WorkflowExecution naming
// the owner and workflow the secrets request must match.
func signedComputeRequestsForParams(t *testing.T, params confidentialrelaytypes.SecretsRequestParams) []confidentialrelaytypes.SignedComputeRequest {
	t.Helper()
	privs, _ := testWorkflowDONKeys()
	publicData, err := proto.Marshal(&confidentialworkflow.WorkflowExecution{
		Owner:      params.Owner,
		WorkflowId: params.WorkflowID,
	})
	require.NoError(t, err)
	cr := confidentialrelaytypes.ComputeRequest{PublicData: publicData}
	payload := confidentialrelaytypes.SignedComputeRequestSignaturePayload(cr.Hash())
	quorum := 2*testEnclaveF + 1
	out := make([]confidentialrelaytypes.SignedComputeRequest, quorum)
	for i := range quorum {
		out[i] = confidentialrelaytypes.SignedComputeRequest{
			ComputeRequest: cr,
			Signature:      ed25519.Sign(privs[i], payload),
		}
	}
	return out
}

// secretsGetTestRegistry builds a mock registry with a vault executable that
// returns a valid GetSecretsResponse for the "API_KEY" secret.
func secretsGetTestRegistry(t *testing.T) *mockCapRegistry {
	t.Helper()
	// Must match secretsGetTestParams(t).EnclavePublicKey and pass the
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

	// withEnclaveConfig wires WorkflowDON.Members and F to match testEnclaveConfig (so both
	// verifyEnclaveConfigMatchesDON and verifyWorkflowAuthorization pass); we only set the
	// DON identity used in the vault request metadata here.
	return withEnclaveConfig(&mockCapRegistry{
		executables: map[string]*mockExecutable{
			vault.CapabilityID: {
				execResult: capabilities.CapabilityResponse{Payload: payload},
			},
		},
		localNode: capabilities.Node{
			WorkflowDON: capabilities.DON{
				ID:            42,
				ConfigVersion: 7,
			},
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
	return makeRequest(t, confidentialrelaytypes.MethodSecretsGet, secretsGetTestParams(t))
}

// secretsGetTestParams returns the canonical valid params used by both the request
// builder and the response-signature-verification step. EnclaveConfig and
// SignedComputeRequests are excluded from the response hash, but the relay validates
// both before serving (EnclaveConfig via Validate, SignedComputeRequests via the
// PRIV-433 Workflow DON authorization check), so the fixture must populate them.
func secretsGetTestParams(t *testing.T) confidentialrelaytypes.SecretsRequestParams {
	t.Helper()
	params := confidentialrelaytypes.SecretsRequestParams{
		WorkflowID:       "wf-secrets-1",
		Owner:            "0xab5801a7d398351b8be11c439e05c5b3259aec9b", // lowercase, should be normalized
		ExecutionID:      "0000000000000000000000000000000000000000000000000000000000000001",
		OrgID:            "org-123",
		EnclavePublicKey: "aabbcc",
		EnclaveConfig:    testEnclaveConfigPtr(),
		Secrets: []confidentialrelaytypes.SecretIdentifier{
			{Key: "API_KEY", Namespace: "main"},
		},
		Attestation: testAttestationB64,
	}
	params.SignedComputeRequests = signedComputeRequestsForParams(t, params)
	return params
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
					WorkflowID:    "wf-1",
					Owner:         testOwner, // chainlink-common#2032 requires 0x-prefixed 20-byte hex
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:   "17",
					CapabilityID:  "my-cap@1.0.0",
					Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
					EnclaveConfig: testEnclaveConfigPtr(),
					Attestation:   testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				params := confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:   "17",
					CapabilityID:  "my-cap@1.0.0",
					Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
					EnclaveConfig: testEnclaveConfigPtr(),
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
				assert.Empty(t, exec.lastRequest.Metadata.OrgID)
			},
		},
		{
			name: "capability execute sets metadata org id when PropagateOrgIDInRequestMetadata enabled",
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
			modifyHandler: func(t *testing.T, h *Handler) {
				getter, err := settings.NewJSONGetter([]byte(`{"global":{"PropagateOrgIDInRequestMetadata":"true"}}`))
				require.NoError(t, err)
				h.limitsFactory = limits.Factory{Logger: h.lggr, Settings: getter}
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					OrgID:         "org-relay-1",
					ReferenceID:   "17",
					CapabilityID:  "my-cap@1.0.0",
					Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
					EnclaveConfig: testEnclaveConfigPtr(),
					Attestation:   testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				params := confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					OrgID:         "org-relay-1",
					ReferenceID:   "17",
					CapabilityID:  "my-cap@1.0.0",
					Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
					EnclaveConfig: testEnclaveConfigPtr(),
				}
				var result confidentialrelaytypes.SignedCapabilityResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				require.Len(t, result.Signatures, 1)
				assertValidCapabilitySignature(t, params, result)
			},
			checkExecutable: func(t *testing.T, reg *mockCapRegistry) {
				exec := reg.executables["my-cap@1.0.0"]
				require.NotNil(t, exec.lastRequest)
				assert.Equal(t, "org-relay-1", exec.lastRequest.Metadata.OrgID)
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
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:   "17",
					CapabilityID:  "my-cap@1.0.0",
					Payload:       makeCapabilityPayload(t, map[string]any{"echo": "hello"}),
					EnclaveConfig: testEnclaveConfigPtr(),
					Attestation:   testAttestationB64,
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
					WorkflowID:    "wf-1",
					CapabilityID:  "missing-cap@1.0.0",
					Payload:       base64.StdEncoding.EncodeToString([]byte("payload")),
					EnclaveConfig: testEnclaveConfigPtr(),
					Attestation:   testAttestationB64,
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
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:   "17",
					CapabilityID:  "fail-cap@1.0.0",
					Payload:       base64.StdEncoding.EncodeToString(b),
					EnclaveConfig: testEnclaveConfigPtr(),
					Attestation:   testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				params := confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
					ReferenceID:   "17",
					CapabilityID:  "fail-cap@1.0.0",
					Payload:       base64.StdEncoding.EncodeToString(mustMarshalProto(t, &sdkpb.CapabilityRequest{Id: "fail-cap@1.0.0", Method: "Execute"})),
					EnclaveConfig: testEnclaveConfigPtr(),
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
				params := secretsGetTestParams(t)
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
				assert.Empty(t, exec.lastRequest.Metadata.OrgID, "org metadata requires PropagateOrgIDInRequestMetadata CRE setting")
			},
		},
		{
			name:     "secrets get sets metadata org id when PropagateOrgIDInRequestMetadata enabled",
			registry: secretsGetTestRegistry,
			req:      secretsGetTestRequest,
			modifyHandler: func(t *testing.T, h *Handler) {
				getter, err := settings.NewJSONGetter([]byte(`{"global":{"PropagateOrgIDInRequestMetadata":"true"}}`))
				require.NoError(t, err)
				h.limitsFactory = limits.Factory{Logger: h.lggr, Settings: getter}
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.Nil(t, resp.Error)
				params := secretsGetTestParams(t)
				params.Attestation = ""
				var result confidentialrelaytypes.SignedSecretsResponseResult
				require.NoError(t, json.Unmarshal(*resp.Result, &result))
				require.Len(t, result.Signatures, 1)
				assertValidSecretsSignature(t, params, result)
			},
			checkExecutable: func(t *testing.T, reg *mockCapRegistry) {
				exec := reg.executables[vault.CapabilityID]
				require.NotNil(t, exec.lastRequest)
				assert.Equal(t, "org-123", exec.lastRequest.Metadata.OrgID)
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

// TestHandler_VerifyEnclaveConfig covers the PRIV-458 / CL112-01 relay-side
// hardening: after the Nitro attestation cryptographically verifies the
// request hash, the handler must also compare the attested EnclaveConfig
// value against the local node's WorkflowDON state. Without this check, a
// malicious host can produce a genuinely-attested request over a forged
// EnclaveConfig and have it accepted.
func TestHandler_VerifyEnclaveConfig(t *testing.T) {
	t.Run("matching config accepted on capability execute", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: testEnclaveConfigPtr(),
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.Nil(t, resp.Error)
	})

	t.Run("nil config rejected on capability execute (required)", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: nil, // missing config cannot be checked against DON state
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.NotNil(t, resp.Error)
	})

	t.Run("F below DON minimum rejected on capability execute", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		badCfg := testEnclaveConfig()
		badCfg.F = testEnclaveF - 1 // below the DON's minimum F
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: &badCfg,
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.NotNil(t, resp.Error)
	})

	t.Run("F above DON minimum accepted on capability execute", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		cfg := testEnclaveConfig()
		cfg.F = testEnclaveF + 1 // a higher F is a stricter quorum; floor check accepts it
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: &cfg,
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.Nil(t, resp.Error)
	})

	t.Run("signers count mismatch rejected on capability execute", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		badCfg := testEnclaveConfig()
		badCfg.Signers = badCfg.Signers[:2]
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: &badCfg,
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.NotNil(t, resp.Error)
	})

	t.Run("signer value mismatch rejected on capability execute", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		badCfg := testEnclaveConfig()
		badCfg.Signers = [][]byte{
			make32Byte(0xa1),
			make32Byte(0xb1),
			make32Byte(0xc1),
			make32Byte(0xff), // last signer differs
		}
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: &badCfg,
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.NotNil(t, resp.Error)
	})

	t.Run("matching is order-independent on capability execute", func(t *testing.T) {
		reg := withEnclaveConfig(&mockCapRegistry{
			executables: map[string]*mockExecutable{
				"my-cap@1.0.0": {execResult: capabilities.CapabilityResponse{Payload: &anypb.Any{}}},
			},
		})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		shuffled := testEnclaveConfig()
		// Reverse Signers; the comparison must still pass.
		n := len(shuffled.Signers)
		rev := make([][]byte, n)
		for i, s := range shuffled.Signers {
			rev[n-1-i] = s
		}
		shuffled.Signers = rev
		req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1",
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: &shuffled,
			Attestation:   testAttestationB64,
		})
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.Nil(t, resp.Error)
	})

	t.Run("F below DON minimum rejected on secrets get", func(t *testing.T) {
		reg := secretsGetTestRegistry(t)
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		params := secretsGetTestParams(t)
		params.EnclaveConfig.F = testEnclaveF - 1 // below the DON's minimum F
		req := makeRequest(t, confidentialrelaytypes.MethodSecretsGet, params)
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		resp := gwConn.lastResp
		require.NotNil(t, resp.Error)
	})
}

func TestTranslateVaultResponse_BinaryShares(t *testing.T) {
	enclaveKey := "aabbcc"
	shareBytes := []byte("share-1")
	vaultResp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: &vault.SecretIdentifier{Key: "API_KEY", Namespace: vaulttypes.DefaultNamespace},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: hex.EncodeToString([]byte("encrypted-value")),
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: enclaveKey,
								BinaryShares:  [][]byte{shareBytes},
							},
						},
					},
				},
			},
		},
	}

	result, err := translateVaultResponse(vaultResp, enclaveKey)
	require.NoError(t, err)
	require.Len(t, result.Secrets, 1)
	require.Equal(t, base64.StdEncoding.EncodeToString(shareBytes), result.Secrets[0].EncryptedShares[0])
}

// Relay aggregation requires byte-identical responses from every node to reach
// quorum. A map-valued response is the worst case: protobuf map fields marshal
// in nondeterministic key order by default, so the same logical value would
// serialize to different bytes across nodes (and across calls) unless we force
// deterministic serialization. This guards both the Any payload construction in
// toSDKCapabilityResponse and the outer marshal. [CL112-05]
func TestToSDKCapabilityResponse_DeterministicSerialization(t *testing.T) {
	t.Parallel()
	val, err := values.Wrap(map[string]any{
		"alpha":   1,
		"bravo":   "two",
		"charlie": true,
		"delta":   []any{1, 2, 3},
		"echo":    map[string]any{"nested1": "x", "nested2": "y", "nested3": "z"},
		"foxtrot": "value-foxtrot",
		"golf":    "value-golf",
		"hotel":   "value-hotel",
	})
	require.NoError(t, err)
	valMap, ok := val.(*values.Map)
	require.True(t, ok)

	capResp := capabilities.CapabilityResponse{Value: valMap}

	var want []byte
	for i := range 50 {
		sdkResp, err := toSDKCapabilityResponse(capResp)
		require.NoError(t, err)
		got, err := proto.MarshalOptions{Deterministic: true}.Marshal(sdkResp)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		if i == 0 {
			want = got
			continue
		}
		require.Equal(t, want, got,
			"serialized capability response must be byte-identical across calls (iteration %d)", i)
	}
}

func TestTranslateVaultResponse_HexShares(t *testing.T) {
	enclaveKey := "aabbcc"
	shareBytes := []byte("share-1")
	vaultResp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: &vault.SecretIdentifier{Key: "API_KEY", Namespace: vaulttypes.DefaultNamespace},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: hex.EncodeToString([]byte("encrypted-value")),
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: enclaveKey,
								Shares:        []string{hex.EncodeToString(shareBytes)},
							},
						},
					},
				},
			},
		},
	}

	result, err := translateVaultResponse(vaultResp, enclaveKey)
	require.NoError(t, err)
	require.Len(t, result.Secrets, 1)
	require.Equal(t, base64.StdEncoding.EncodeToString(shareBytes), result.Secrets[0].EncryptedShares[0])
}

func TestVerifyWorkflowAuthorization(t *testing.T) {
	t.Parallel()
	const (
		owner      = "0xab5801a7d398351b8be11c439e05c5b3259aec9b"
		workflowID = "wf-secrets-1"
	)
	privs, _ := testWorkflowDONKeys()
	// The DON members are the public keys of privs; F=testEnclaveF => 2*F+1 = 3 quorum.
	don := capabilities.DON{Members: testWorkflowDONMembers(), F: testEnclaveF}

	// signedReqs builds compute requests naming o/wf, each signed by one of the given keys
	// over the shared request hash.
	signedReqs := func(t *testing.T, o, wf string, signers []ed25519.PrivateKey) []confidentialrelaytypes.SignedComputeRequest {
		t.Helper()
		publicData, err := proto.Marshal(&confidentialworkflow.WorkflowExecution{Owner: o, WorkflowId: wf})
		require.NoError(t, err)
		cr := confidentialrelaytypes.ComputeRequest{PublicData: publicData}
		payload := confidentialrelaytypes.SignedComputeRequestSignaturePayload(cr.Hash())
		out := make([]confidentialrelaytypes.SignedComputeRequest, len(signers))
		for i, s := range signers {
			out[i] = confidentialrelaytypes.SignedComputeRequest{ComputeRequest: cr, Signature: ed25519.Sign(s, payload)}
		}
		return out
	}

	// validParams: 2*F+1 = 3 distinct DON signers over a compute request naming owner/workflow.
	validParams := func(t *testing.T) confidentialrelaytypes.SecretsRequestParams {
		t.Helper()
		return confidentialrelaytypes.SecretsRequestParams{
			WorkflowID:            workflowID,
			Owner:                 owner,
			SignedComputeRequests: signedReqs(t, owner, workflowID, privs[:3]),
		}
	}

	h := newTestHandler(t, &mockCapRegistry{}, &mockGatewayConnector{})

	t.Run("valid 2F+1 quorum", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, h.verifyWorkflowAuthorization(don, validParams(t)))
	})

	t.Run("missing signed compute requests", func(t *testing.T) {
		t.Parallel()
		params := validParams(t)
		params.SignedComputeRequests = nil
		require.ErrorContains(t, h.verifyWorkflowAuthorization(don, params), "missing signed compute requests")
	})

	t.Run("insufficient signers for quorum", func(t *testing.T) {
		t.Parallel()
		params := validParams(t)
		// Only 2 signers; F=1 requires 2*1+1 = 3.
		params.SignedComputeRequests = signedReqs(t, owner, workflowID, privs[:2])
		require.ErrorContains(t, h.verifyWorkflowAuthorization(don, params), "insufficient Workflow DON signatures")
	})

	t.Run("signers not in Workflow DON", func(t *testing.T) {
		t.Parallel()
		stranger := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{0xfe}, ed25519.SeedSize))
		params := validParams(t)
		params.SignedComputeRequests = signedReqs(t, owner, workflowID, []ed25519.PrivateKey{stranger, stranger, stranger})
		require.ErrorContains(t, h.verifyWorkflowAuthorization(don, params), "insufficient Workflow DON signatures")
	})

	t.Run("owner mismatch", func(t *testing.T) {
		t.Parallel()
		params := validParams(t)
		params.Owner = "0x0000000000000000000000000000000000000002"
		require.ErrorContains(t, h.verifyWorkflowAuthorization(don, params), "owner not authorized")
	})

	t.Run("workflow id mismatch", func(t *testing.T) {
		t.Parallel()
		params := validParams(t)
		params.WorkflowID = "wf-other"
		require.ErrorContains(t, h.verifyWorkflowAuthorization(don, params), "workflow_id not authorized")
	})

	t.Run("forwarded requests disagree on compute request", func(t *testing.T) {
		t.Parallel()
		params := validParams(t)
		params.SignedComputeRequests = append(params.SignedComputeRequests, confidentialrelaytypes.SignedComputeRequest{
			ComputeRequest: confidentialrelaytypes.ComputeRequest{PublicData: []byte("different")},
			Signature:      []byte("irrelevant"),
		})
		require.ErrorContains(t, h.verifyWorkflowAuthorization(don, params), "do not share one compute request")
	})
}
