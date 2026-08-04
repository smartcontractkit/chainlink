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
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/teeattestation"
	"github.com/smartcontractkit/chainlink-common/pkg/teeattestation/nitro"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

var _ core.GatewayConnectorHandler = (*Handler)(nil)

const (
	HandlerName          = "EnclaveRelayHandler"
	internalErrorMessage = "internal error"

	// confidentialWorkflowsCapID is the capability ID for the confidential
	// workflows enclave pool. The relay handler uses it to look up trusted
	// enclave measurements from the capabilities registry.
	confidentialWorkflowsCapID = "confidential-workflows@1.0.0-alpha"

	// defaultGetExecutionWait bounds how long a relay callback waits for a
	// not-yet-registered execution handler to appear before giving up. The
	// enclave runs one DON-shared execution and only needs a relay quorum, so its
	// callback can reach a node before that node has started its own copy of the
	// execution and registered a handler (the start-edge race). Waiting briefly
	// lets a straggling node register and still sign, instead of dropping below
	// quorum. Kept well under the gateway relay request timeout (RequestTimeoutSec,
	// default 30s) so the node still answers in time to be counted.
	defaultGetExecutionWait = 5 * time.Second
)

// enclaveEntry mirrors the enclave config shape stored in the capabilities
// registry. Only the fields needed for attestation validation are included.
type enclaveEntry struct {
	TrustedValues []json.RawMessage `json:"trustedValues"`
	CARootsPEM    string            `json:"caRootsPEM,omitempty"`
}

type enclavesList struct {
	Enclaves []enclaveEntry
}

type handlerMetrics struct {
	requestCount         metric.Int64Counter
	requestLatency       metric.Int64Histogram
	requestInternalError metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*handlerMetrics, error) {
	requestCount, err := beholder.GetMeter().Int64Counter("enclave_relay_request_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register request count counter: %w", err)
	}
	requestLatency, err := beholder.GetMeter().Int64Histogram("enclave_relay_request_latency_ms", metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to register request latency histogram: %w", err)
	}
	requestInternalError, err := beholder.GetMeter().Int64Counter("enclave_relay_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}
	requestSuccess, err := beholder.GetMeter().Int64Counter("enclave_relay_request_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register success counter: %w", err)
	}
	return &handlerMetrics{
		requestCount:         requestCount,
		requestLatency:       requestLatency,
		requestInternalError: requestInternalError,
		requestSuccess:       requestSuccess,
	}, nil
}

// AttestationValidator validates TEE attestation documents, with or without a
// custom CA root. Both the production Nitro validator and the insecure
// passthrough validator implement it, so the handler validates the same way
// regardless of which is configured.
type AttestationValidator interface {
	ValidateAttestation(attestation, expectedUserData, trustedMeasurements []byte) error
	ValidateAttestationWithRoots(attestation, expectedUserData, trustedMeasurements []byte, caRootsPEM string) error
}

// nitroValidator is the production validator backed by AWS Nitro.
type nitroValidator struct{}

func (nitroValidator) ValidateAttestation(attestation, expectedUserData, trustedMeasurements []byte) error {
	return nitro.ValidateAttestation(attestation, expectedUserData, trustedMeasurements)
}

func (nitroValidator) ValidateAttestationWithRoots(attestation, expectedUserData, trustedMeasurements []byte, caRootsPEM string) error {
	return nitro.ValidateAttestationWithRoots(attestation, expectedUserData, trustedMeasurements, caRootsPEM)
}

// NewAttestationValidator returns the production validator backed by AWS Nitro.
func NewAttestationValidator() AttestationValidator {
	return nitroValidator{}
}

// Handler processes enclave relay requests from the gateway.
// It validates attestations and proxies requests to VaultDON or capability DONs.
type Handler struct {
	services.Service
	eng *services.Engine

	capRegistry       core.CapabilitiesRegistry
	executionHandlers *ExecutionHandlers
	gatewayConnector  core.GatewayConnector
	responseSigner    relayResponseSigner
	lggr              logger.Logger
	metrics           *handlerMetrics

	// validator validates TEE attestation documents.
	validator AttestationValidator
	// requireBFTQuorum selects the required signature quorum: when true the relay
	// demands a Byzantine quorum of 2*F+1 unique signers, otherwise a crash-fault
	// quorum of F+1.
	requireBFTQuorum bool
	limitsFactory    limits.Factory
	// getExecutionWait is how long a relay callback waits for a not-yet-registered
	// execution handler before failing (see defaultGetExecutionWait).
	getExecutionWait time.Duration

	// serveTime bounds how long one gateway request may be served for, from
	// ConfidentialCompute.ConfidentialRelayHandlerTimeout. The connector's context
	// carries no deadline of its own, so this is what keeps a stalled vault call
	// from outliving the request it belongs to.
	serveTime limits.TimeLimiter
}

func NewHandler(capRegistry core.CapabilitiesRegistry, executionHandlers *ExecutionHandlers, conn core.GatewayConnector, responseSigner relayResponseSigner, lggr logger.Logger, lf limits.Factory, validator AttestationValidator, requireBFTQuorum bool) (*Handler, error) {
	if responseSigner == nil {
		return nil, errors.New("response signer is required")
	}
	m, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}
	serveTime, err := lf.MakeTimeLimiter(cresettings.Default.ConfidentialCompute.ConfidentialRelayHandlerTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create serve time limiter: %w", err)
	}

	h := &Handler{
		capRegistry:       capRegistry,
		executionHandlers: executionHandlers,
		gatewayConnector:  conn,
		responseSigner:    responseSigner,
		lggr:              logger.Named(lggr, HandlerName),
		metrics:           m,
		validator:         validator,
		requireBFTQuorum:  requireBFTQuorum,
		limitsFactory:     lf,
		serveTime:         serveTime,
		getExecutionWait:  defaultGetExecutionWait,
	}
	h.Service, h.eng = services.Config{
		Name:  HandlerName,
		Start: h.start,
		Close: h.close,
	}.NewServiceEngine(lggr)
	return h, nil
}

func (h *Handler) start(ctx context.Context) error {
	if err := h.gatewayConnector.AddHandler(ctx, h.Methods(), h); err != nil {
		return fmt.Errorf("failed to add enclave relay handler to connector: %w", err)
	}
	return nil
}

func (h *Handler) close() error {
	err := h.serveTime.Close()
	if rmErr := h.gatewayConnector.RemoveHandler(context.Background(), h.Methods()); rmErr != nil {
		err = errors.Join(err, fmt.Errorf("failed to remove enclave relay handler from connector: %w", rmErr))
	}
	return err
}

func (h *Handler) ID(_ context.Context) (string, error) {
	return HandlerName, nil
}

func (h *Handler) Methods() []string {
	return []string{confidentialrelaytypes.MethodSecretsGet, confidentialrelaytypes.MethodCapabilityExec}
}

func (h *Handler) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) error {
	h.lggr.Debugw("received message from gateway", "gatewayID", gatewayID, "requestID", req.ID)

	// GoCtx ties the goroutine to the handler's service lifecycle, so Close waits
	// for in-flight requests instead of abandoning them mid-vault-call.
	h.eng.GoCtx(ctx, func(ctx context.Context) {
		ctx, done, err := h.serveTime.WithTimeout(ctx)
		if err != nil {
			h.lggr.Errorw("failed to apply serve timeout, dropping request", "gatewayID", gatewayID, "requestID", req.ID, "err", err)
			return
		}
		defer done()
		h.serveGatewayMessage(ctx, gatewayID, req)
	})

	return nil
}

func (h *Handler) serveGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) {
	startTime := time.Now()
	outcome := "success"
	var errorCode int64
	defer func() {
		attrs := []attribute.KeyValue{
			attribute.String("gateway_id", gatewayID),
			attribute.String("method", req.Method),
			attribute.String("outcome", outcome),
		}
		if errorCode != 0 {
			attrs = append(attrs, attribute.Int64("error_code", errorCode))
		}
		h.metrics.requestCount.Add(ctx, 1, metric.WithAttributes(attrs...))
		h.metrics.requestLatency.Record(ctx, time.Since(startTime).Milliseconds(), metric.WithAttributes(attrs...))
	}()

	var response *jsonrpc.Response[json.RawMessage]
	switch req.Method {
	case confidentialrelaytypes.MethodSecretsGet:
		response = h.handleSecretsGet(ctx, gatewayID, req)
	case confidentialrelaytypes.MethodCapabilityExec:
		response = h.handleCapabilityExecute(ctx, gatewayID, req)
	default:
		response = h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrMethodNotFound, errors.New("unsupported method: "+req.Method))
	}
	if response != nil && response.Error != nil {
		outcome = "error"
		errorCode = response.Error.Code
	}

	if err := h.gatewayConnector.SendToGateway(ctx, gatewayID, response); err != nil {
		outcome = "send_error"
		h.lggr.Errorw("failed to send message to gateway", "gatewayID", gatewayID, "err", err)
		return
	}

	h.lggr.Infow("sent message to gateway", "gatewayID", gatewayID, "requestID", req.ID)
	if response != nil && response.Error == nil {
		h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
			attribute.String("gateway_id", gatewayID),
		))
	}
}

func (h *Handler) handleSecretsGet(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	if req.Params == nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, errors.New("missing params"))
	}
	var params confidentialrelaytypes.SecretsRequestParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, err)
	}

	att := params.Attestation
	params.Attestation = ""
	if err := h.verifyAttestationHash(ctx, att, params, confidentialrelaytypes.DomainSecretsGet); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}
	// Fetch the local node once: it provides the WorkflowDON snapshot for both the
	// enclave-config check and the Workflow-DON authorization check below, plus the DON
	// metadata on the vault request. A registry read failure is node-side, so ErrInternal.
	localNode, err := h.capRegistry.LocalNode(ctx)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("failed to get local node: %w", err))
	}
	// Verify the enclave's reported config matches the onchain DON state
	// before treating the attested request as trusted: the Nitro attestation
	// binds the request hash, but a malicious host can produce a
	// genuinely-attested request over a forged enclave config unless we
	// compare the config value against the DON reference.
	if err = h.verifyEnclaveConfigMatchesDON(localNode, params.EnclaveConfig); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}

	// Beyond attestation, verify the Workflow DON authorized this request: the enclave
	// forwards the Workflow-DON-signed compute requests (a 2*F+1 quorum), whose PublicData
	// names the authorized owner. A TEE breach passes attestation but cannot forge a Workflow
	// DON quorum over a different owner (PRIV-433).
	if err = h.verifyWorkflowAuthorization(localNode.WorkflowDON, params); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, err)
	}

	secretsRequest := &sdkpb.GetSecretsRequest{Requests: make([]*sdkpb.SecretRequest, 0, len(params.Secrets))}

	for _, s := range params.Secrets {
		secretsRequest.Requests = append(secretsRequest.Requests, &sdkpb.SecretRequest{
			Id:        s.Key,
			Namespace: s.Namespace,
		})
	}

	// Resolve the execution handler only after attestation and Workflow-DON
	// authorization have passed, so an unverified callback cannot make the node
	// park a waiter. The enclave's callback can arrive before this node has started
	// its own copy of the execution (start-edge race); wait briefly for the handler
	// to register rather than failing and dropping below relay quorum. Bounded so a
	// callback for an execution this node never runs still fails in time to respond
	// within the gateway's relay request window.
	waitCtx, cancel := context.WithTimeout(ctx, h.getExecutionWait)
	defer cancel()
	handler, ok := h.executionHandlers.GetExecutionWithWait(waitCtx, params.WorkflowID, params.ExecutionID)
	if !ok {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, fmt.Errorf("execution handler for workflow %s execution %s not found", params.WorkflowID, params.ExecutionID))
	}

	vaultResp, err := handler.GetRawSecrets(ctx, secretsRequest, teeKeyFetcher(params.EnclavePublicKey))
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}

	result, err := translateVaultResponse(vaultResp, params.EnclavePublicKey)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}

	signedResult, err := h.signSecretsResponse(params, result)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("failed to sign secrets response: %w", err))
	}

	return h.jsonResponse(req, signedResult)
}

// resolveDONID determines the DON ID for a capability.
// Keeping for potential future use by handleCapabilityExecute.
func (h *Handler) resolveDONID(ctx context.Context, capability capabilities.ExecutableCapability) (uint32, error) { //nolint:unused // reserved for future multi-DON routing in handleCapabilityExecute
	info, err := capability.Info(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get capability info: %w", err)
	}
	if info.IsLocal {
		localNode, err := h.capRegistry.LocalNode(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get local node: %w", err)
		}
		return localNode.WorkflowDON.ID, nil
	}
	if info.DON == nil {
		return 0, errors.New("capability is not associated with any DON")
	}
	return info.DON.ID, nil
}

// translateVaultResponse converts a vault GetSecretsResponse to the enclave relay protocol format.
// Encoding conversion: ciphertext hex (vault) -> base64 (enclave relay); encrypted shares may be
// binary or hex (vault) -> base64 (enclave relay).
func translateVaultResponse(vaultResp []*vault.SecretResponse, enclaveKey string) (*confidentialrelaytypes.SecretsResponseResult, error) {
	result := &confidentialrelaytypes.SecretsResponseResult{}

	for _, sr := range vaultResp {
		if sr.GetError() != "" {
			return nil, fmt.Errorf("vault error for secret %s/%s: %s", sr.Id.GetNamespace(), sr.Id.GetKey(), sr.GetError())
		}

		data := sr.GetData()
		if data == nil {
			return nil, fmt.Errorf("vault returned no data for secret %s/%s", sr.Id.GetNamespace(), sr.Id.GetKey())
		}

		encryptedBytes, err := hex.DecodeString(data.EncryptedValue)
		if err != nil {
			return nil, fmt.Errorf("failed to decode encrypted value for %s: %w", sr.Id.GetKey(), err)
		}

		var shares []string
		for _, es := range data.EncryptedDecryptionKeyShares {
			if es.EncryptionKey == enclaveKey {
				if len(es.BinaryShares) > 0 {
					for _, shareBytes := range es.BinaryShares {
						shares = append(shares, base64.StdEncoding.EncodeToString(shareBytes))
					}
				} else {
					for _, share := range es.Shares {
						shareBytes, err := hex.DecodeString(share)
						if err != nil {
							return nil, fmt.Errorf("failed to decode share: %w", err)
						}
						shares = append(shares, base64.StdEncoding.EncodeToString(shareBytes))
					}
				}
				break
			}
		}
		if len(shares) == 0 {
			return nil, fmt.Errorf("no shares found for enclave key in secret %s/%s", sr.Id.GetNamespace(), sr.Id.GetKey())
		}

		result.Secrets = append(result.Secrets, confidentialrelaytypes.SecretEntry{
			ID: confidentialrelaytypes.SecretIdentifier{
				Key:       sr.Id.GetKey(),
				Namespace: sr.Id.GetNamespace(),
			},
			Ciphertext:      base64.StdEncoding.EncodeToString(encryptedBytes),
			EncryptedShares: shares,
		})
	}

	return result, nil
}

func (h *Handler) handleCapabilityExecute(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	if req.Params == nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, errors.New("missing params"))
	}
	var params confidentialrelaytypes.CapabilityRequestParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, err)
	}

	// The enclave's capability calls arrive as fresh gateway messages rather than
	// through the workflow engine, so ctx carries none of the CRE tenants the engine
	// seeds.
	//
	// Seeded as soon as the params are parsed so every ctx use below carries the
	// tenant.
	ctx = contexts.WithCRE(ctx, contexts.CRE{
		Org:      params.OrgID,
		Owner:    params.Owner,
		Workflow: params.WorkflowID,
	})

	att := params.Attestation
	params.Attestation = ""
	if err := h.verifyAttestationHash(ctx, att, params, confidentialrelaytypes.DomainCapabilityExec); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}

	// Verify the enclave's reported config matches the onchain DON state, same as
	// handleSecretsGet (PRIV-458): the Nitro attestation binds the request hash but
	// not the config value, so a malicious host could otherwise produce a
	// genuinely-attested request over a forged enclave config.
	localNode, err := h.capRegistry.LocalNode(ctx)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("failed to get local node: %w", err))
	}
	if err = h.verifyEnclaveConfigMatchesDON(localNode, params.EnclaveConfig); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(params.Payload)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, fmt.Errorf("failed to decode payload: %w", err))
	}

	sdkReq := &sdkpb.CapabilityRequest{}
	if err = proto.Unmarshal(payloadBytes, sdkReq); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, fmt.Errorf("failed to unmarshal capability request: %w", err))
	}

	// Resolve the execution handler only after attestation and enclave-config
	// verification, so an unverified callback cannot make the node park a waiter.
	// The enclave's callback can beat this node's own execution start (start-edge
	// race); a bounded wait lets a straggler register and sign instead of dropping
	// below relay quorum (see handleSecretsGet).
	waitCtx, cancel := context.WithTimeout(ctx, h.getExecutionWait)
	defer cancel()
	handler, ok := h.executionHandlers.GetExecutionWithWait(waitCtx, params.WorkflowID, params.ExecutionID)
	if !ok {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, fmt.Errorf("execution handler for workflow %s execution %s not found", params.WorkflowID, params.ExecutionID))
	}

	capResp, execErr := handler.CallCapability(ctx, sdkReq)

	var result confidentialrelaytypes.CapabilityResponseResult
	if execErr != nil {
		result.Error = execErr.Error()
	} else {
		// Deterministic marshal so every relay node emits byte-identical
		// response payloads; relay aggregation requires identical bytes to
		// reach quorum. [CL112-05]
		var respBytes []byte
		respBytes, err = proto.MarshalOptions{Deterministic: true}.Marshal(capResp)
		if err != nil {
			return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("marshalling capability response: %w", err))
		}
		result.Payload = base64.StdEncoding.EncodeToString(respBytes)
	}

	signedResult, err := h.signCapabilityResponse(params, result)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("failed to sign capability response: %w", err))
	}

	return h.jsonResponse(req, signedResult)
}

// verifyEnclaveConfigMatchesDON compares the enclave's reported EnclaveConfig
// against the local node's WorkflowDON membership and fault tolerance. The
// relay DON runs on the same nodes as the workflow DON, so
// localNode.WorkflowDON.Members is the right comparison target.
//
// PRIV-458: the Nitro attestation binds the request hash but does not on its
// own prove the config matches the DON, so a malicious host could produce a
// genuinely-attested request over a forged config. Comparing the attested
// config against onchain DON state closes that gap.
//
// localNode is passed in so each request fetches it once (it feeds request
// metadata too); the caller's lookup is an O(1) in-memory read populated by
// the registry syncer, so this stays off the RPC hot path.
//
// cfg is required: a nil EnclaveConfig is rejected. The wire field stays
// optional in chainlink-common so older senders compile, but the relay is the
// security boundary and will not service a request that omits the config, since
// an absent config cannot be checked against onchain DON state.
func (h *Handler) verifyEnclaveConfigMatchesDON(localNode capabilities.Node, cfg *confidentialrelaytypes.EnclaveConfig) error {
	if cfg == nil {
		return errors.New("enclave config is required")
	}
	expectedF := uint32(localNode.WorkflowDON.F)
	// F is a floor, matching the CC-side check (validateEnclaveSigners): the
	// enclave's reported F must meet or exceed the DON's fault tolerance. A
	// higher F is a stricter quorum and is acceptable; a lower one is not.
	if cfg.F < expectedF {
		return fmt.Errorf("enclave config F %d does not meet the minimum required DON F %d", cfg.F, expectedF)
	}
	if len(cfg.Signers) != len(localNode.WorkflowDON.Members) {
		return fmt.Errorf("enclave config signers count mismatch: enclave reports %d, expected %d",
			len(cfg.Signers), len(localNode.WorkflowDON.Members))
	}
	expected := make([][]byte, len(localNode.WorkflowDON.Members))
	for i := range localNode.WorkflowDON.Members {
		expected[i] = localNode.WorkflowDON.Members[i][:]
	}
	actual := append([][]byte(nil), cfg.Signers...)
	sort.Slice(actual, func(i, j int) bool { return bytes.Compare(actual[i], actual[j]) < 0 })
	sort.Slice(expected, func(i, j int) bool { return bytes.Compare(expected[i], expected[j]) < 0 })
	for i := range actual {
		if !bytes.Equal(actual[i], expected[i]) {
			return fmt.Errorf("enclave config signer mismatch at sorted index %d: enclave reports %x, expected %x",
				i, actual[i], expected[i])
		}
	}
	return nil
}

// getEnclaveAttestationConfig reads the enclave pool configuration from the
// capabilities registry and returns trusted measurement sets and CA roots
// for attestation validation. Called per-request so the config stays fresh
// (same pattern as CC's EnsureFreshEnclaves).
func (h *Handler) getEnclaveAttestationConfig(ctx context.Context) ([]json.RawMessage, string, error) {
	dons, err := h.capRegistry.DONsForCapability(ctx, confidentialWorkflowsCapID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find DON for %s: %w", confidentialWorkflowsCapID, err)
	}
	if len(dons) == 0 {
		return nil, "", fmt.Errorf("no DON found hosting %s", confidentialWorkflowsCapID)
	}

	capConfig, err := h.capRegistry.ConfigForCapability(ctx, confidentialWorkflowsCapID, dons[0].DON.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get config for %s: %w", confidentialWorkflowsCapID, err)
	}

	if capConfig.DefaultConfig == nil {
		return nil, "", fmt.Errorf("no default config for %s", confidentialWorkflowsCapID)
	}

	var enclaves enclavesList
	if err := capConfig.DefaultConfig.UnwrapTo(&enclaves); err != nil {
		return nil, "", fmt.Errorf("failed to unwrap enclave config: %w", err)
	}

	var measurements []json.RawMessage
	var caRootsPEM string
	for _, e := range enclaves.Enclaves {
		measurements = append(measurements, e.TrustedValues...)
		if caRootsPEM == "" && e.CARootsPEM != "" {
			caRootsPEM = e.CARootsPEM
		}
	}
	return measurements, caRootsPEM, nil
}

func (h *Handler) verifyAttestationHash(ctx context.Context, attestationB64 string, cleanParams any, domainTag string) error {
	if attestationB64 == "" {
		return errors.New("missing attestation")
	}

	paramsJSON, err := json.Marshal(cleanParams)
	if err != nil {
		return fmt.Errorf("failed to marshal params for attestation: %w", err)
	}

	hash := teeattestation.DomainHash(domainTag, paramsJSON)

	attestationBytes, err := base64.StdEncoding.DecodeString(attestationB64)
	if err != nil {
		return fmt.Errorf("failed to decode attestation: %w", err)
	}

	// Look up trusted measurements from the capabilities registry.
	// Each enclave instance may have different PCR values, so we try each set
	// and succeed if any match (same approach as pool.go's
	// validateAttestationAgainstMultipleMeasurements).
	measurements, caRootsPEM, err := h.getEnclaveAttestationConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get enclave attestation config: %w", err)
	}

	if len(measurements) == 0 {
		return errors.New("no trusted measurements found in capabilities registry")
	}

	var validationErr error
	for _, m := range measurements {
		var err error
		if caRootsPEM != "" {
			err = h.validator.ValidateAttestationWithRoots(attestationBytes, hash, m, caRootsPEM)
		} else {
			err = h.validator.ValidateAttestation(attestationBytes, hash, m)
		}
		if err == nil {
			return nil
		}
		validationErr = errors.Join(validationErr, err)
	}
	return fmt.Errorf("no trusted measurement set matched: %w", validationErr)
}

// verifyWorkflowAuthorization is the PRIV-433 check beyond attestation. Attestation only
// proves the request came from genuine enclave code; it does not prove the Workflow DON
// authorized fetching this owner's secrets. A compromised TEE would still pass attestation
// while self-asserting a victim's owner.
//
// The enclave forwards the Workflow-DON-signed compute requests it executed (a 2*F+1 quorum,
// where F is the Workflow DON fault tolerance). Each node signs the same ComputeRequest.Hash();
// we reconstruct that hash, verify each signature against the onchain Workflow DON signer set,
// and require the quorum of unique signers. The signed PublicData names the authorized owner
// and workflow, which must match this request. A breached enclave cannot forge a Workflow DON
// quorum over a different owner.
//
// All failures here are client errors: the request is unauthorized. The caller fetches the
// Workflow DON (a server-side concern) and passes it in, so registry failures stay internal.
func (h *Handler) verifyWorkflowAuthorization(don capabilities.DON, params confidentialrelaytypes.SecretsRequestParams) error {
	if len(params.SignedComputeRequests) == 0 {
		return errors.New("missing signed compute requests")
	}

	// Match the enclave's own quorum. With requireBFTQuorum the relay demands a
	// Byzantine quorum of 2*F+1 unique signers; otherwise a crash-fault quorum of
	// F+1 suffices.
	threshold := int(don.F) + 1
	if h.requireBFTQuorum {
		threshold = 2*int(don.F) + 1
	}

	// The forwarded requests differ only in their signature; they all sign one shared
	// ComputeRequest hash. Reconstruct that hash once and verify each signature over it.
	hash := params.SignedComputeRequests[0].Hash()
	payload := confidentialrelaytypes.SignedComputeRequestSignaturePayload(hash)

	signers := make(map[string]struct{})
	for _, scr := range params.SignedComputeRequests {
		if scr.Hash() != hash {
			return errors.New("forwarded signed compute requests do not share one compute request")
		}
		for _, member := range don.Members {
			if ed25519.Verify(ed25519.PublicKey(member[:]), payload, scr.Signature) {
				signers[member.String()] = struct{}{}
				break
			}
		}
	}
	if len(signers) < threshold {
		return fmt.Errorf("insufficient Workflow DON signatures: %d unique signers, need %d", len(signers), threshold)
	}

	// The signed request authorizes a specific owner and workflow; the secrets request must
	// match both, or a breached enclave could fetch another owner's secrets.
	var execution confidentialworkflow.WorkflowExecution
	if err := proto.Unmarshal(params.SignedComputeRequests[0].PublicData, &execution); err != nil {
		return fmt.Errorf("failed to unmarshal workflow execution from public data: %w", err)
	}
	if !common.IsHexAddress(params.Owner) || !common.IsHexAddress(execution.GetOwner()) {
		return errors.New("invalid owner address")
	}
	if common.HexToAddress(execution.GetOwner()) != common.HexToAddress(params.Owner) {
		return fmt.Errorf("owner not authorized: request %q vs signed %q", params.Owner, execution.GetOwner())
	}
	if execution.GetWorkflowId() != params.WorkflowID {
		return fmt.Errorf("workflow_id not authorized: request %q vs signed %q", params.WorkflowID, execution.GetWorkflowId())
	}
	return nil
}

func (h *Handler) jsonResponse(req *jsonrpc.Request[json.RawMessage], result any) *jsonrpc.Response[json.RawMessage] {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		h.lggr.Errorw("failed to marshal response", "err", err)
		return &jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      req.ID,
			Method:  req.Method,
			Error: &jsonrpc.WireError{
				Code:    jsonrpc.ErrInternal,
				Message: internalErrorMessage,
			},
		}
	}
	resultJSON := json.RawMessage(resultBytes)
	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Result:  &resultJSON,
	}
}

func (h *Handler) signSecretsResponse(
	params confidentialrelaytypes.SecretsRequestParams,
	result *confidentialrelaytypes.SecretsResponseResult,
) (*confidentialrelaytypes.SignedSecretsResponseResult, error) {
	if h.responseSigner == nil {
		return nil, errors.New("response signer not configured")
	}

	hash, err := result.Hash(params)
	if err != nil {
		return nil, fmt.Errorf("hash secrets response: %w", err)
	}
	signature, err := h.responseSigner.Sign(confidentialrelaytypes.RelayResponseSignaturePayload(hash))
	if err != nil {
		return nil, err
	}

	sig := confidentialrelaytypes.RelayResponseSignature{
		Signer:    h.responseSigner.PublicKey(),
		Signature: signature,
	}
	return &confidentialrelaytypes.SignedSecretsResponseResult{
		Result:    *result,
		Signature: sig,
		// Keep populating Signatures too: there are still readers of the array. The
		// cleanup order is readers first, then writers (this), then the field itself,
		// or builds go red.
		Signatures: []confidentialrelaytypes.RelayResponseSignature{sig},
	}, nil
}

func (h *Handler) signCapabilityResponse(
	params confidentialrelaytypes.CapabilityRequestParams,
	result confidentialrelaytypes.CapabilityResponseResult,
) (*confidentialrelaytypes.SignedCapabilityResponseResult, error) {
	if h.responseSigner == nil {
		return nil, errors.New("response signer not configured")
	}

	hash, err := result.Hash(params)
	if err != nil {
		return nil, fmt.Errorf("hash capability response: %w", err)
	}
	signature, err := h.responseSigner.Sign(confidentialrelaytypes.RelayResponseSignaturePayload(hash))
	if err != nil {
		return nil, err
	}

	sig := confidentialrelaytypes.RelayResponseSignature{
		Signer:    h.responseSigner.PublicKey(),
		Signature: signature,
	}
	return &confidentialrelaytypes.SignedCapabilityResponseResult{
		Result:    result,
		Signature: sig,
		// Keep populating Signatures too: there are still readers of the array. The
		// cleanup order is readers first, then writers (this), then the field itself,
		// or builds go red.
		Signatures: []confidentialrelaytypes.RelayResponseSignature{sig},
	}, nil
}

func (h *Handler) errorResponse(
	ctx context.Context,
	gatewayID string,
	req *jsonrpc.Request[json.RawMessage],
	errorCode int64,
	err error,
) *jsonrpc.Response[json.RawMessage] {
	h.lggr.Errorw("request error", "errorCode", errorCode, "err", err)
	h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
		attribute.Int64("error_code", errorCode),
	))

	message := err.Error()
	if errorCode == jsonrpc.ErrInternal {
		message = internalErrorMessage
	}

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Error: &jsonrpc.WireError{
			Code:    errorCode,
			Message: message,
		},
	}
}
