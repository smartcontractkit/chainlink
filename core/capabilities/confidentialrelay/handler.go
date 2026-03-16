package confidentialrelay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay/attestation"
)

var _ core.GatewayConnectorHandler = (*Handler)(nil)

const HandlerName = "EnclaveRelayHandler"

type handlerMetrics struct {
	requestInternalError metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*handlerMetrics, error) {
	requestInternalError, err := beholder.GetMeter().Int64Counter("enclave_relay_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}
	requestSuccess, err := beholder.GetMeter().Int64Counter("enclave_relay_request_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register success counter: %w", err)
	}
	return &handlerMetrics{
		requestInternalError: requestInternalError,
		requestSuccess:       requestSuccess,
	}, nil
}

type gatewayConnector interface {
	SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error
	AddHandler(ctx context.Context, methods []string, handler core.GatewayConnectorHandler) error
	RemoveHandler(ctx context.Context, methods []string) error
}

// attestationValidatorFunc validates a Nitro attestation document.
type attestationValidatorFunc func(attestation []byte, expectedUserData []byte, trustedMeasurements []byte, caRootsPEM string) error

// Handler processes enclave relay requests from the gateway.
// It validates Nitro attestations and proxies requests to VaultDON or capability DONs.
type Handler struct {
	services.Service
	eng *services.Engine

	capRegistry      core.CapabilitiesRegistry
	gatewayConnector gatewayConnector
	trustedPCRs      []byte
	lggr             logger.Logger
	metrics          *handlerMetrics

	// validateAttestation validates Nitro attestation documents.
	// Defaults to the real nitro validator; overridden in tests.
	validateAttestation attestationValidatorFunc

	// caRootsPEM is an optional PEM-encoded CA root certificate for attestation
	// validation. Empty string uses DefaultCARoots (production default).
	caRootsPEM string
}

func NewHandler(capRegistry core.CapabilitiesRegistry, conn gatewayConnector, trustedPCRs []byte, lggr logger.Logger, caRootsPEM ...string) (*Handler, error) {
	m, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	var roots string
	if len(caRootsPEM) > 0 {
		roots = caRootsPEM[0]
	}
	h := &Handler{
		capRegistry:         capRegistry,
		gatewayConnector:    conn,
		trustedPCRs:         trustedPCRs,
		lggr:                logger.Named(lggr, HandlerName),
		metrics:             m,
		validateAttestation: attestation.ValidateNitroAttestation,
		caRootsPEM:          roots,
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
	if err := h.gatewayConnector.RemoveHandler(context.Background(), h.Methods()); err != nil {
		return fmt.Errorf("failed to remove enclave relay handler from connector: %w", err)
	}
	return nil
}

func (h *Handler) ID(_ context.Context) (string, error) {
	return HandlerName, nil
}

func (h *Handler) Methods() []string {
	return []string{confidentialrelaytypes.MethodCapabilityExec}
}

func (h *Handler) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) error {
	h.lggr.Debugw("received message from gateway", "gatewayID", gatewayID, "requestID", req.ID)

	var response *jsonrpc.Response[json.RawMessage]
	switch req.Method {
	case confidentialrelaytypes.MethodCapabilityExec:
		response = h.handleCapabilityExecute(ctx, gatewayID, req)
	default:
		response = h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrMethodNotFound, errors.New("unsupported method: "+req.Method))
	}

	if err := h.gatewayConnector.SendToGateway(ctx, gatewayID, response); err != nil {
		h.lggr.Errorw("failed to send message to gateway", "gatewayID", gatewayID, "err", err)
		return err
	}

	h.lggr.Infow("sent message to gateway", "gatewayID", gatewayID, "requestID", req.ID)
	h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
	))
	return nil
}

func (h *Handler) handleCapabilityExecute(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	var params confidentialrelaytypes.CapabilityRequestParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, err)
	}

	att := params.Attestation
	params.Attestation = ""
	if err := h.verifyAttestationHash(att, params, confidentialrelaytypes.DomainCapabilityExec); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, err)
	}

	cap, err := h.capRegistry.GetExecutable(ctx, params.CapabilityID)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("capability not found: %w", err))
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(params.Payload)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, fmt.Errorf("failed to decode payload: %w", err))
	}

	var sdkReq sdkpb.CapabilityRequest
	if err := proto.Unmarshal(payloadBytes, &sdkReq); err != nil {
		return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInvalidParams, fmt.Errorf("failed to unmarshal capability request: %w", err))
	}

	capReq := capabilities.CapabilityRequest{
		Payload:      sdkReq.Payload,
		Method:       sdkReq.Method,
		CapabilityId: params.CapabilityID,
		Metadata: capabilities.RequestMetadata{
			WorkflowID: params.WorkflowID,
		},
	}

	// Backward compatibility: extract values.Map from Payload into Inputs
	// for old-style capabilities that only look at Inputs.
	if sdkReq.Payload != nil {
		var valPB valuespb.Value
		if sdkReq.Payload.UnmarshalTo(&valPB) == nil {
			if v, vErr := values.FromProto(&valPB); vErr == nil {
				if m, ok := v.(*values.Map); ok {
					capReq.Inputs = m
				}
			}
		}
	}

	capResp, execErr := cap.Execute(ctx, capReq)

	var result confidentialrelaytypes.CapabilityResponseResult
	if execErr != nil {
		result.Error = execErr.Error()
	} else {
		sdkResp, err := toSDKCapabilityResponse(capResp)
		if err != nil {
			return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("converting capability response: %w", err))
		}
		respBytes, err := proto.Marshal(sdkResp)
		if err != nil {
			return h.errorResponse(ctx, gatewayID, req, jsonrpc.ErrInternal, fmt.Errorf("marshalling capability response: %w", err))
		}
		result.Payload = base64.StdEncoding.EncodeToString(respBytes)
	}

	return h.jsonResponse(req, result)
}

func (h *Handler) verifyAttestationHash(attestationB64 string, cleanParams any, domainTag string) error {
	if attestationB64 == "" {
		return errors.New("missing attestation")
	}

	paramsJSON, err := json.Marshal(cleanParams)
	if err != nil {
		return fmt.Errorf("failed to marshal params for attestation: %w", err)
	}

	hash := attestation.DomainHash(domainTag, paramsJSON)

	attestationBytes, err := base64.StdEncoding.DecodeString(attestationB64)
	if err != nil {
		return fmt.Errorf("failed to decode attestation: %w", err)
	}

	return h.validateAttestation(attestationBytes, hash, h.trustedPCRs, h.caRootsPEM)
}

func toSDKCapabilityResponse(capResp capabilities.CapabilityResponse) (*sdkpb.CapabilityResponse, error) {
	if capResp.Payload != nil {
		return &sdkpb.CapabilityResponse{
			Response: &sdkpb.CapabilityResponse_Payload{Payload: capResp.Payload},
		}, nil
	}

	if capResp.Value != nil {
		valProto := values.Proto(capResp.Value)
		wrapped, err := anypb.New(valProto)
		if err != nil {
			return nil, fmt.Errorf("wrapping value map in Any: %w", err)
		}
		return &sdkpb.CapabilityResponse{
			Response: &sdkpb.CapabilityResponse_Payload{Payload: wrapped},
		}, nil
	}

	return &sdkpb.CapabilityResponse{}, nil
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
				Message: err.Error(),
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

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Error: &jsonrpc.WireError{
			Code:    errorCode,
			Message: err.Error(),
		},
	}
}
