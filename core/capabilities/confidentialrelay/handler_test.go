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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/teeattestation/passthrough"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// capExecExecutionID is the canonical 32-byte hex execution ID used by the
// capability-execute fixtures. The handler looks up the execution helper by
// (WorkflowID, ExecutionID), so tests must register their helper under it.
const capExecExecutionID = "32c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1"

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

type mockGatewayConnector struct {
	core.UnimplementedGatewayConnector
	// mu guards resps: the handler serves requests on its own goroutines, so
	// sends race with test assertions without it.
	mu           sync.Mutex
	resps        []*jsonrpc.Response[json.RawMessage]
	addedMethods []string
	removed      bool
}

func (m *mockGatewayConnector) SendToGateway(_ context.Context, _ string, resp *jsonrpc.Response[json.RawMessage]) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resps = append(m.resps, resp)
	return nil
}

func (m *mockGatewayConnector) respCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.resps)
}

// waitResp blocks until the handler has sent at least one response and returns
// the most recent. HandleGatewayMessage dispatches and returns without serving,
// so assertions have to wait for the send rather than read straight after it.
func (m *mockGatewayConnector) waitResp(t *testing.T) *jsonrpc.Response[json.RawMessage] {
	t.Helper()
	require.Eventually(t, func() bool { return m.respCount() > 0 },
		testWaitTimeout, time.Millisecond, "handler did not send a response to the gateway")
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.resps[len(m.resps)-1]
}

// testWaitTimeout bounds waits for the handler's asynchronous sends.
const testWaitTimeout = 30 * time.Second

func (m *mockGatewayConnector) AddHandler(_ context.Context, methods []string, _ core.GatewayConnectorHandler) error {
	m.addedMethods = methods
	return nil
}
func (m *mockGatewayConnector) RemoveHandler(_ context.Context, _ []string) error {
	m.removed = true
	return nil
}

// mockExecutionHelper is a host.ExecutionHelperWithRawSecrets stub. The relay handler only
// calls CallCapability (capability execute) and GetRawSecrets (secrets get); the embedded
// interface satisfies the rest of the surface and panics if anything else is invoked, which
// keeps the mock minimal while still flagging unexpected calls.
type mockExecutionHelper struct {
	host.ExecutionHelperWithRawSecrets

	capResp    *sdkpb.CapabilityResponse
	capErr     error
	rawSecrets []*vault.SecretResponse
	secretsErr error

	lastCapabilityRequest *sdkpb.CapabilityRequest
	lastSecretsRequest    *sdkpb.GetSecretsRequest
	// lastCapabilityCRE records the CRE tenants on the context CallCapability was
	// invoked with. The tenant-scoped limiters downstream fail closed without them.
	lastCapabilityCRE contexts.CRE
}

func (m *mockExecutionHelper) CallCapability(ctx context.Context, req *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	m.lastCapabilityRequest = req
	m.lastCapabilityCRE = contexts.CREValue(ctx)
	return m.capResp, m.capErr
}

func (m *mockExecutionHelper) GetRawSecrets(_ context.Context, req *sdkpb.GetSecretsRequest, _ host.EncryptionKeyFetcher) ([]*vault.SecretResponse, error) {
	m.lastSecretsRequest = req
	return m.rawSecrets, m.secretsErr
}

type mockCapRegistry struct {
	core.UnimplementedCapabilitiesRegistry
	configs   map[string]capabilities.CapabilityConfiguration
	dons      map[string][]capabilities.DONWithNodes
	localNode capabilities.Node
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
	validator, err := passthrough.New()
	require.NoError(t, err)
	h, err := NewHandler(registry, &ExecutionHandlers{}, gwConn, newRelayResponseSigner(key), lggr, limits.Factory{Logger: lggr}, validator, true)
	require.NoError(t, err)
	// Keep the not-yet-registered wait short so the "handler not found" cases do
	// not stall the suite for the production default.
	h.getExecutionWait = 50 * time.Millisecond
	return h
}

// capExecHelper returns a helper that succeeds with a fixed capability response payload.
func capExecHelper() *mockExecutionHelper {
	return &mockExecutionHelper{
		capResp: &sdkpb.CapabilityResponse{
			Response: &sdkpb.CapabilityResponse_Payload{Payload: &anypb.Any{Value: []byte("result-proto-bytes")}},
		},
	}
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

// secretsGetTestRegistry builds a mock registry whose enclave config, DON membership and
// local node satisfy the relay's attestation, enclave-config and workflow-authorization
// checks. The vault data is returned by the execution helper (secretsGetTestHelper), not
// the registry.
func secretsGetTestRegistry(t *testing.T) *mockCapRegistry {
	t.Helper()
	// withEnclaveConfig wires WorkflowDON.Members and F to match testEnclaveConfig (so both
	// verifyEnclaveConfigMatchesDON and verifyWorkflowAuthorization pass); we only set the
	// DON identity used in the vault request metadata here.
	return withEnclaveConfig(&mockCapRegistry{
		localNode: capabilities.Node{
			WorkflowDON: capabilities.DON{
				ID:            42,
				ConfigVersion: 7,
			},
		},
	})
}

// secretsGetTestHelper returns the execution helper that serves the "API_KEY" secret for
// the secrets-get fixtures. The handler calls GetRawSecrets and translates the result.
func secretsGetTestHelper() *mockExecutionHelper {
	// Must match secretsGetTestParams(t).EnclavePublicKey and pass the
	// confidentialrelay.validateEnclavePublicKey hex check (#2032).
	enclaveKey := "aabbcc"
	return &mockExecutionHelper{
		rawSecrets: []*vault.SecretResponse{
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
		workflowID      string
		executionID     string
		helper          func(t *testing.T) *mockExecutionHelper
		checkResp       func(t *testing.T, resp *jsonrpc.Response[json.RawMessage])
		checkExecutable func(t *testing.T, helper *mockExecutionHelper)
	}{
		{
			name: "capability execute success",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			workflowID:  "wf-1",
			executionID: capExecExecutionID,
			helper:      func(_ *testing.T) *mockExecutionHelper { return capExecHelper() },
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					Owner:         testOwner, // chainlink-common#2032 requires 0x-prefixed 20-byte hex
					ExecutionID:   capExecExecutionID,
					OrgID:         "org-1",
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
					ExecutionID:   capExecExecutionID,
					OrgID:         "org-1",
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
			checkExecutable: func(t *testing.T, helper *mockExecutionHelper) {
				require.NotNil(t, helper.lastCapabilityRequest, "CallCapability should have been called")
				assert.Equal(t, "my-cap@1.0.0", helper.lastCapabilityRequest.Id)
				assert.Equal(t, "Execute", helper.lastCapabilityRequest.Method)
				// Without the CRE tenants on the context, every tenant-scoped
				// limiter downstream fails closed rather than reading a limit.
				// Normalized(): WithCRE strips the owner's 0x prefix and lowercases
				// it, so the tenant key matches the one the engine path produces.
				assert.Equal(t, contexts.CRE{
					Org:      "org-1",
					Owner:    testOwner,
					Workflow: "wf-1",
				}.Normalized(), helper.lastCapabilityCRE)
			},
		},
		{
			name: "capability execute attestation failure",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			workflowID:  "wf-1",
			executionID: "",
			helper:      func(_ *testing.T) *mockExecutionHelper { return capExecHelper() },
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
			name: "capability execute execution handler not found",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			// No helper registered. The lookup now runs after attestation and
			// enclave-config checks, so the request must pass those (valid
			// attestation, config, and a decodable payload) to reach the not-found
			// path rather than failing earlier.
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					ExecutionID:   capExecExecutionID,
					CapabilityID:  "missing-cap@1.0.0",
					Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
					EnclaveConfig: testEnclaveConfigPtr(),
					Attestation:   testAttestationB64,
				})
			},
			checkResp: func(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) {
				require.NotNil(t, resp.Error)
				assert.Equal(t, jsonrpc.ErrInvalidParams, resp.Error.Code)
			},
		},
		{
			name: "capability execute error returned in result",
			registry: func(_ *testing.T) *mockCapRegistry {
				return withEnclaveConfig(&mockCapRegistry{})
			},
			workflowID:  "wf-1",
			executionID: capExecExecutionID,
			helper: func(_ *testing.T) *mockExecutionHelper {
				return &mockExecutionHelper{capErr: errors.New("execution failed")}
			},
			req: func(t *testing.T) *jsonrpc.Request[json.RawMessage] {
				sdkReq := &sdkpb.CapabilityRequest{Id: "fail-cap@1.0.0", Method: "Execute"}
				b, err := proto.Marshal(sdkReq)
				require.NoError(t, err)
				return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
					WorkflowID:    "wf-1",
					Owner:         testOwner,
					ExecutionID:   capExecExecutionID,
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
					ExecutionID:   capExecExecutionID,
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
			name:        "secrets get invokes vault execute with stable capability metadata",
			registry:    secretsGetTestRegistry,
			req:         secretsGetTestRequest,
			workflowID:  "wf-secrets-1",
			executionID: "0000000000000000000000000000000000000000000000000000000000000001",
			helper:      func(_ *testing.T) *mockExecutionHelper { return secretsGetTestHelper() },
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
			checkExecutable: func(t *testing.T, helper *mockExecutionHelper) {
				require.NotNil(t, helper.lastSecretsRequest, "GetRawSecrets should have been called")
				require.Len(t, helper.lastSecretsRequest.Requests, 1)
				assert.Equal(t, "API_KEY", helper.lastSecretsRequest.Requests[0].Id)
				assert.Equal(t, "main", helper.lastSecretsRequest.Requests[0].Namespace)
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
			var helper *mockExecutionHelper
			if tt.helper != nil {
				helper = tt.helper(t)
				h.executionHandlers.AddExecution(tt.workflowID, tt.executionID, helper)
			}
			err := h.HandleGatewayMessage(t.Context(), "gw-1", tt.req(t))
			require.NoError(t, err)
			resp := gwConn.waitResp(t)
			tt.checkResp(t, resp)
			if tt.checkExecutable != nil {
				tt.checkExecutable(t, helper)
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
// EnclaveConfig and have it accepted. The check runs on both the capability
// execute and secrets get paths.
func TestHandler_VerifyEnclaveConfig(t *testing.T) {
	t.Parallel()
	// capExecHandler builds a handler whose execution helper succeeds, so any rejection
	// observed comes from the enclave-config check rather than a missing execution.
	capExecHandler := func(t *testing.T) (*Handler, *mockGatewayConnector) {
		t.Helper()
		reg := withEnclaveConfig(&mockCapRegistry{})
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		h.executionHandlers.AddExecution("wf-1", capExecExecutionID, capExecHelper())
		return h, gwConn
	}

	capExecReq := func(t *testing.T, cfg *confidentialrelaytypes.EnclaveConfig) *jsonrpc.Request[json.RawMessage] {
		t.Helper()
		return makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
			WorkflowID:    "wf-1",
			Owner:         testOwner,
			ExecutionID:   capExecExecutionID,
			ReferenceID:   "1",
			CapabilityID:  "my-cap@1.0.0",
			Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
			EnclaveConfig: cfg,
			Attestation:   testAttestationB64,
		})
	}

	t.Run("matching config accepted on capability execute", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, testEnclaveConfigPtr()))
		require.NoError(t, err)
		require.Nil(t, gwConn.waitResp(t).Error)
	})

	t.Run("nil config rejected on capability execute (required)", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		// missing config cannot be checked against DON state
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, nil))
		require.NoError(t, err)
		require.NotNil(t, gwConn.waitResp(t).Error)
	})

	t.Run("F below DON minimum rejected on capability execute", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		badCfg := testEnclaveConfig()
		badCfg.F = testEnclaveF - 1 // below the DON's minimum F
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, &badCfg))
		require.NoError(t, err)
		require.NotNil(t, gwConn.waitResp(t).Error)
	})

	t.Run("F above DON minimum accepted on capability execute", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		cfg := testEnclaveConfig()
		cfg.F = testEnclaveF + 1 // a higher F is a stricter quorum; floor check accepts it
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, &cfg))
		require.NoError(t, err)
		require.Nil(t, gwConn.waitResp(t).Error)
	})

	t.Run("signers count mismatch rejected on capability execute", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		badCfg := testEnclaveConfig()
		badCfg.Signers = badCfg.Signers[:2]
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, &badCfg))
		require.NoError(t, err)
		require.NotNil(t, gwConn.waitResp(t).Error)
	})

	t.Run("signer value mismatch rejected on capability execute", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		badCfg := testEnclaveConfig()
		badCfg.Signers = [][]byte{
			make32Byte(0xa1),
			make32Byte(0xb1),
			make32Byte(0xc1),
			make32Byte(0xff), // last signer differs
		}
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, &badCfg))
		require.NoError(t, err)
		require.NotNil(t, gwConn.waitResp(t).Error)
	})

	t.Run("matching is order-independent on capability execute", func(t *testing.T) {
		t.Parallel()
		h, gwConn := capExecHandler(t)
		shuffled := testEnclaveConfig()
		// Reverse Signers; the comparison must still pass.
		n := len(shuffled.Signers)
		rev := make([][]byte, n)
		for i, s := range shuffled.Signers {
			rev[n-1-i] = s
		}
		shuffled.Signers = rev
		err := h.HandleGatewayMessage(context.Background(), "gw-1", capExecReq(t, &shuffled))
		require.NoError(t, err)
		require.Nil(t, gwConn.waitResp(t).Error)
	})

	t.Run("F below DON minimum rejected on secrets get", func(t *testing.T) {
		t.Parallel()
		reg := secretsGetTestRegistry(t)
		gwConn := &mockGatewayConnector{}
		h := newTestHandler(t, reg, gwConn)
		params := secretsGetTestParams(t)
		h.executionHandlers.AddExecution(params.WorkflowID, params.ExecutionID, secretsGetTestHelper())
		params.EnclaveConfig.F = testEnclaveF - 1 // below the DON's minimum F
		req := makeRequest(t, confidentialrelaytypes.MethodSecretsGet, params)
		err := h.HandleGatewayMessage(context.Background(), "gw-1", req)
		require.NoError(t, err)
		require.NotNil(t, gwConn.waitResp(t).Error)
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

	result, err := translateVaultResponse(vaultResp.Responses, enclaveKey)
	require.NoError(t, err)
	require.Len(t, result.Secrets, 1)
	require.Equal(t, base64.StdEncoding.EncodeToString(shareBytes), result.Secrets[0].EncryptedShares[0])
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

	result, err := translateVaultResponse(vaultResp.Responses, enclaveKey)
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

// blockingExecutionHelper reports when CallCapability is entered and holds there
// until released, so a test can observe how many requests the handler serves at
// the same time.
type blockingExecutionHelper struct {
	host.ExecutionHelperWithRawSecrets
	entered chan struct{}
	release chan struct{}
}

func newBlockingExecutionHelper(capacity int) *blockingExecutionHelper {
	return &blockingExecutionHelper{
		entered: make(chan struct{}, capacity),
		release: make(chan struct{}),
	}
}

func (b *blockingExecutionHelper) CallCapability(ctx context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	b.entered <- struct{}{}
	select {
	case <-b.release:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return &sdkpb.CapabilityResponse{
		Response: &sdkpb.CapabilityResponse_Payload{Payload: &anypb.Any{Value: []byte("result-proto-bytes")}},
	}, nil
}

// blockingCapExecHandler wires a handler whose capability calls block until the
// returned helper is released.
func blockingCapExecHandler(t *testing.T, capacity int) (*Handler, *mockGatewayConnector, *blockingExecutionHelper) {
	t.Helper()
	helper := newBlockingExecutionHelper(capacity)
	gwConn := &mockGatewayConnector{}
	h := newTestHandler(t, withEnclaveConfig(&mockCapRegistry{}), gwConn)
	h.executionHandlers.AddExecution("wf-1", capExecExecutionID, helper)
	return h, gwConn, helper
}

func blockingCapExecRequest(t *testing.T, id string) *jsonrpc.Request[json.RawMessage] {
	t.Helper()
	req := makeRequest(t, confidentialrelaytypes.MethodCapabilityExec, confidentialrelaytypes.CapabilityRequestParams{
		WorkflowID:    "wf-1",
		Owner:         testOwner,
		ExecutionID:   capExecExecutionID,
		ReferenceID:   "1",
		CapabilityID:  "my-cap@1.0.0",
		Payload:       makeCapabilityPayload(t, map[string]any{"key": "val"}),
		EnclaveConfig: testEnclaveConfigPtr(),
		Attestation:   testAttestationB64,
	})
	req.ID = id
	return req
}

// TestHandler_ServesRequestsConcurrently is the regression guard for the
// head-of-line block that made a burst of relay requests fail. The gateway
// connector calls HandleGatewayMessage from a single goroutine per connection
// and waits for it to return, so serving inline forced requests through one at a
// time and the tail of a burst outlived the caller's timeout. Each call must
// return promptly and the requests must overlap.
func TestHandler_ServesRequestsConcurrently(t *testing.T) {
	t.Parallel()
	const concurrent = 8
	h, gwConn, helper := blockingCapExecHandler(t, concurrent)

	for i := range concurrent {
		require.NoError(t, h.HandleGatewayMessage(t.Context(), "gw-1", blockingCapExecRequest(t, fmt.Sprintf("req-%d", i))))
	}

	// Every request reached the capability call while all the others were still
	// held, which is only possible if none of them blocked the dispatch path.
	for i := range concurrent {
		select {
		case <-helper.entered:
		case <-time.After(testWaitTimeout):
			t.Fatalf("only %d of %d requests were served concurrently; dispatch is serialized", i, concurrent)
		}
	}

	// Nothing can have answered yet: all of them are parked in the capability call.
	require.Zero(t, gwConn.respCount())

	close(helper.release)
	require.Eventually(t, func() bool { return gwConn.respCount() == concurrent },
		testWaitTimeout, time.Millisecond, "expected one response per request")
}
