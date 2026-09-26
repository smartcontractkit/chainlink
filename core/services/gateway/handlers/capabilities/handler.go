package capabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	MethodWorkflowSyncer = "workflow_syncer"

	// Error messages
	ErrTransformingMessageToRequest = "error transforming message to request"

	handlerName = "WebAPIHandler"
)

type handler struct {
	services.StateMachine
	config          HandlerConfig
	router          *handlers.NodeRouter
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *ratelimit.RateLimiter
	wg              sync.WaitGroup
	metrics         *metrics
}

type HandlerConfig struct {
	NodeRateLimiter ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
}

var _ handlers.Handler = (*handler)(nil)

func NewHandler(handlerConfig json.RawMessage, shardedDONs []config.ShardedDONConfig, shardsConnMgrs [][]handlers.DON, httpClient network.HTTPClient, lggr logger.Logger) (*handler, error) {
	var cfg HandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}

	router, err := handlers.BuildNodeRouter(shardedDONs, shardsConnMgrs)
	if err != nil {
		return nil, err
	}
	defaultDonID := ""
	if len(shardedDONs) > 0 {
		defaultDonID = shardedDONs[0].DonName
	}

	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, err
	}

	metrics, err := newMetrics()
	if err != nil {
		return nil, err
	}

	return &handler{
		config:          cfg,
		router:          router,
		lggr:            logger.Named(lggr, "WebAPIHandler."+defaultDonID),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		metrics:         metrics,
	}, nil
}

// sendHTTPMessageToClient is an outgoing message from the gateway to external endpoints
// returns message to be sent back to the capability node
func (h *handler) sendHTTPMessageToClient(ctx context.Context, req network.HTTPRequest, msg *api.Message) (*api.Message, error) {
	var payload Response
	resp, err := h.httpClient.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	payload = Response{
		ExecutionError: false,
		StatusCode:     resp.StatusCode,
		Headers:        resp.Headers,
		Body:           resp.Body,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &api.Message{
		Body: api.MessageBody{
			MessageID: msg.Body.MessageID,
			Method:    msg.Body.Method,
			DonID:     msg.Body.DonID,
			Payload:   payloadBytes,
		},
	}, nil
}

func (h *handler) handleWebAPIOutgoingMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.lggr.Debugw("handling webAPI outgoing message", "messageId", msg.Body.MessageID, "nodeAddr", nodeAddr)
	h.metrics.recordArtifactFetchRequestReceived(ctx, msg.Body.Method)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	var payload Request
	err := json.Unmarshal(msg.Body.Payload, &payload)
	if err != nil {
		return err
	}

	timeout := time.Duration(payload.TimeoutMs) * time.Millisecond
	req := network.HTTPRequest{
		Method:           payload.Method,
		URL:              payload.URL,
		Headers:          payload.Headers,
		Body:             payload.Body,
		MaxResponseBytes: payload.MaxResponseBytes,
		Timeout:          timeout,
	}

	// send response to node async
	h.wg.Go(func() {
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := logger.With(h.lggr, "url", payload.URL, "messageId", msg.Body.MessageID, "method", payload.Method, "timeout", payload.TimeoutMs, "workflowID", payload.WorkflowID)
		l.Debug("Sending request to client")
		respMsg, err := h.sendHTTPMessageToClient(newCtx, req, msg)
		switch {
		case err == nil:
			h.metrics.recordArtifactFetchOutcome(newCtx, fetchOutcomeSuccess)
		case errors.Is(err, context.DeadlineExceeded):
			h.metrics.recordArtifactFetchOutcome(newCtx, fetchOutcomeTimeout)
		default:
			h.metrics.recordArtifactFetchOutcome(newCtx, fetchOutcomeExternalError)
		}
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			payload := Response{
				ExecutionError: true,
				ErrorMessage:   err.Error(),
			}
			payloadBytes, err2 := json.Marshal(payload)
			if err2 != nil {
				// should not happen
				l.Errorw("error while marshalling payload", "err", err2)
				return
			}
			respMsg = &api.Message{
				Body: api.MessageBody{
					MessageID: msg.Body.MessageID,
					Method:    msg.Body.Method,
					DonID:     msg.Body.DonID,
					Payload:   payloadBytes,
				},
			}
		}

		// Work around the fact that the connection manager expects all messages
		// to have a valid signature by reusing the signature that came with the message.
		// This is OK to do because:
		// - our trust model for Gateways assumes that we can trust the Gateway node. This is a central assumption since
		// the Gateway node has access to plaintext secrets sent by DON nodes.
		// - the connection between the Gateway and DON Node is already authorized via a DON-side and Gateway-side
		// allowlist, and secured via TLS.
		respMsg.Signature = msg.Signature
		req, err := common.ValidatedRequestFromMessage(respMsg)
		if err != nil {
			l.Errorw(ErrTransformingMessageToRequest, "err", err)
			return
		}
		nodeDON, ok := h.router.DONFor(nodeAddr)
		if !ok {
			l.Errorw("no connection manager found for node", "to", nodeAddr)
			return
		}
		err = nodeDON.SendToNode(newCtx, nodeAddr, req)
		h.metrics.recordArtifactFetchResponseDelivery(newCtx, err == nil)
		if err != nil {
			l.Errorw("failed to send to node", "err", err, "to", nodeAddr)
			return
		}
		l.Debugw("sent response to node", "to", nodeAddr)
	})
	return nil
}

func (h *handler) Methods() []string {
	return []string{
		MethodWorkflowSyncer,
	}
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	msg, err := common.ValidatedMessageFromResp(resp)
	if err != nil {
		return err
	}
	if msg.Body.Sender != nodeAddr {
		return errors.New("message sender mismatch when reading from node ")
	}
	start := time.Now()
	switch msg.Body.Method {
	case MethodWorkflowSyncer:
		err = h.handleWebAPIOutgoingMessage(ctx, msg, nodeAddr)
	default:
		err = fmt.Errorf("unsupported method: %s", msg.Body.Method)
	}
	h.metrics.recordHandleDuration(ctx, time.Since(start), msg.Body.Method, err == nil)
	return err
}

func (h *handler) Start(context.Context) error {
	return h.StartOnce(handlerName, func() error {
		return nil
	})
}

func (h *handler) Close() error {
	return h.StopOnce(handlerName, func() error {
		h.wg.Wait()
		return nil
	})
}

func (h *handler) HandleJSONRPCUserMessage(_ context.Context, _ jsonrpc.Request[json.RawMessage], _ handlers.Callback) error {
	return errors.New("capabilities handler does not support JSON-RPC user messages")
}

func (h *handler) HandleLegacyUserMessage(_ context.Context, _ *api.Message, _ handlers.Callback) error {
	return errors.New("capabilities handler does not support legacy user messages")
}

const (
	fetchOutcomeSuccess       = "success"
	fetchOutcomeTimeout       = "timeout"
	fetchOutcomeExternalError = "external_error"
)

type metrics struct {
	handleDuration    metric.Int64Histogram
	requestsTotal     metric.Int64Counter
	fetchOutcomeTotal metric.Int64Counter
	deliveryTotal     metric.Int64Counter
}

func (m *metrics) recordHandleDuration(ctx context.Context, d time.Duration, method string, success bool) {
	successStr := "false"
	if success {
		successStr = "true"
	}
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", successStr),
		attribute.String("method", method),
	))
}

// recordArtifactFetchRequestReceived counts every artifact-fetch request the gateway
// receives from a node, independent of whether it is later rate-limited or fails.
func (m *metrics) recordArtifactFetchRequestReceived(ctx context.Context, method string) {
	m.requestsTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("method", method),
	))
}

// recordArtifactFetchOutcome counts the result of the outbound HTTP call the gateway
// makes to the external endpoint on behalf of the node.
func (m *metrics) recordArtifactFetchOutcome(ctx context.Context, outcome string) {
	m.fetchOutcomeTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("outcome", outcome),
	))
}

// recordArtifactFetchResponseDelivery counts whether the gateway succeeded in delivering
// the fetch result back to the requesting node over the DON connection.
func (m *metrics) recordArtifactFetchResponseDelivery(ctx context.Context, success bool) {
	successStr := "false"
	if success {
		successStr = "true"
	}
	m.deliveryTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("success", successStr),
	))
}

func newMetrics() (*metrics, error) {
	h, err := beholder.GetMeter().Int64Histogram("platform_gateway_capabilities_handle_node_message_duration_ms")
	if err != nil {
		return nil, err
	}

	requestsTotal, err := beholder.GetMeter().Int64Counter("platform_gateway_capabilities_artifact_fetch_requests_total")
	if err != nil {
		return nil, err
	}

	fetchOutcomeTotal, err := beholder.GetMeter().Int64Counter("platform_gateway_capabilities_artifact_fetch_outcome_total")
	if err != nil {
		return nil, err
	}

	deliveryTotal, err := beholder.GetMeter().Int64Counter("platform_gateway_capabilities_response_delivery_total")
	if err != nil {
		return nil, err
	}

	return &metrics{
		handleDuration:    h,
		requestsTotal:     requestsTotal,
		fetchOutcomeTotal: fetchOutcomeTotal,
		deliveryTotal:     deliveryTotal,
	}, nil
}
