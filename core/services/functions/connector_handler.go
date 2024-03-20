package functions

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/multierr"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	fallow "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist"
	fsub "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

type functionsConnectorHandler struct {
	services.StateMachine

	connector                  connector.GatewayConnector
	signerKey                  *ecdsa.PrivateKey
	nodeAddress                string
	storage                    s4.Storage
	allowlist                  fallow.OnchainAllowlist
	rateLimiter                *hc.RateLimiter
	subscriptions              fsub.OnchainSubscriptions
	minimumBalance             assets.Link
	listener                   FunctionsListener
	offchainTransmitter        OffchainTransmitter
	allowedHeartbeatInitiators map[string]struct{}
	heartbeatRequests          map[RequestID]*HeartbeatResponse
	requestTimeoutSec          uint32
	orderedRequests            []RequestID
	mu                         sync.Mutex
	chStop                     services.StopChan
	shutdownWaitGroup          sync.WaitGroup
	lggr                       logger.Logger
}

const HeartbeatCacheSize = 1000

var (
	_ connector.Signer                  = &functionsConnectorHandler{}
	_ connector.GatewayConnectorHandler = &functionsConnectorHandler{}
)

var (
	promStorageUserUpdatesCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "storage_user_updates",
		Help: "Number of storage updates performed by users",
	}, []string{})
)

// internal request ID is a hash of (sender, requestID)
func InternalId(sender []byte, requestId []byte) RequestID {
	return RequestID(crypto.Keccak256Hash(append(sender, requestId...)).Bytes())
}

func NewFunctionsConnectorHandler(pluginConfig *config.PluginConfig, signerKey *ecdsa.PrivateKey, storage s4.Storage, allowlist fallow.OnchainAllowlist, rateLimiter *hc.RateLimiter, subscriptions fsub.OnchainSubscriptions, listener FunctionsListener, offchainTransmitter OffchainTransmitter, lggr logger.Logger) (*functionsConnectorHandler, error) {
	if signerKey == nil || storage == nil || allowlist == nil || rateLimiter == nil || subscriptions == nil || listener == nil || offchainTransmitter == nil {
		return nil, fmt.Errorf("all dependencies must be non-nil")
	}
	allowedHeartbeatInitiators := make(map[string]struct{})
	for _, initiator := range pluginConfig.AllowedHeartbeatInitiators {
		allowedHeartbeatInitiators[strings.ToLower(initiator)] = struct{}{}
	}
	return &functionsConnectorHandler{
		nodeAddress:                pluginConfig.GatewayConnectorConfig.NodeAddress,
		signerKey:                  signerKey,
		storage:                    storage,
		allowlist:                  allowlist,
		rateLimiter:                rateLimiter,
		subscriptions:              subscriptions,
		minimumBalance:             pluginConfig.MinimumSubscriptionBalance,
		listener:                   listener,
		offchainTransmitter:        offchainTransmitter,
		allowedHeartbeatInitiators: allowedHeartbeatInitiators,
		heartbeatRequests:          make(map[RequestID]*HeartbeatResponse),
		requestTimeoutSec:          pluginConfig.RequestTimeoutSec,
		chStop:                     make(services.StopChan),
		lggr:                       lggr.Named("FunctionsConnectorHandler"),
	}, nil
}

func (h *functionsConnectorHandler) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
}

func (h *functionsConnectorHandler) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(h.signerKey, data...)
}

func (h *functionsConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	body := &msg.Body
	fromAddr := ethCommon.HexToAddress(body.Sender)
	if !h.allowlist.Allow(fromAddr) {
		h.lggr.Errorw("allowlist prevented the request from this address", "id", gatewayId, "address", fromAddr)
		return
	}
	if !h.rateLimiter.Allow(body.Sender) {
		h.lggr.Errorw("request rate-limited", "id", gatewayId, "address", fromAddr)
		return
	}
	h.lggr.Debugw("handling gateway request", "id", gatewayId, "method", body.Method)

	switch body.Method {
	case functions.MethodSecretsList:
		h.handleSecretsList(ctx, gatewayId, body, fromAddr)
	case functions.MethodSecretsSet:
		if balance, err := h.subscriptions.GetMaxUserBalance(fromAddr); err != nil || balance.Cmp(h.minimumBalance.ToInt()) < 0 {
			h.lggr.Errorw("user subscription has insufficient balance", "id", gatewayId, "address", fromAddr, "balance", balance, "minBalance", h.minimumBalance)
			response := functions.ResponseBase{
				Success:      false,
				ErrorMessage: "user subscription has insufficient balance",
			}
			h.sendResponseAndLog(ctx, gatewayId, body, response)
			return
		}
		h.handleSecretsSet(ctx, gatewayId, body, fromAddr)
	case functions.MethodHeartbeat:
		h.handleHeartbeat(ctx, gatewayId, body, fromAddr)
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayId, "method", body.Method)
	}
}

func (h *functionsConnectorHandler) Start(ctx context.Context) error {
	return h.StartOnce("FunctionsConnectorHandler", func() error {
		if err := h.allowlist.Start(ctx); err != nil {
			return err
		}
		if err := h.subscriptions.Start(ctx); err != nil {
			return err
		}
		h.shutdownWaitGroup.Add(1)
		go h.reportLoop()
		return nil
	})
}

func (h *functionsConnectorHandler) Close() error {
	return h.StopOnce("FunctionsConnectorHandler", func() (err error) {
		close(h.chStop)
		err = multierr.Combine(err, h.allowlist.Close())
		err = multierr.Combine(err, h.subscriptions.Close())
		h.shutdownWaitGroup.Wait()
		return
	})
}

func (h *functionsConnectorHandler) handleSecretsList(ctx context.Context, gatewayId string, body *api.MessageBody, fromAddr ethCommon.Address) {
	var response functions.SecretsListResponse
	snapshot, err := h.storage.List(ctx, fromAddr)
	if err == nil {
		response.Success = true
		response.Rows = make([]functions.SecretsListRow, len(snapshot))
		for i, row := range snapshot {
			response.Rows[i] = functions.SecretsListRow{
				SlotID:     row.SlotId,
				Version:    row.Version,
				Expiration: row.Expiration,
			}
		}
	} else {
		response.ErrorMessage = fmt.Sprintf("Failed to list secrets: %v", err)
	}
	h.sendResponseAndLog(ctx, gatewayId, body, response)
}

func (h *functionsConnectorHandler) handleSecretsSet(ctx context.Context, gatewayId string, body *api.MessageBody, fromAddr ethCommon.Address) {
	var request functions.SecretsSetRequest
	var response functions.SecretsSetResponse
	err := json.Unmarshal(body.Payload, &request)
	if err == nil {
		key := s4.Key{
			Address: fromAddr,
			SlotId:  request.SlotID,
			Version: request.Version,
		}
		record := s4.Record{
			Expiration: request.Expiration,
			Payload:    request.Payload,
		}
		h.lggr.Debugw("handling a secrets_set request", "address", fromAddr, "slotId", request.SlotID, "payloadVersion", request.Version, "expiration", request.Expiration)
		err = h.storage.Put(ctx, &key, &record, request.Signature)
		if err == nil {
			response.Success = true
			promStorageUserUpdatesCount.WithLabelValues().Inc()
		} else {
			response.ErrorMessage = fmt.Sprintf("Failed to set secret: %v", err)
		}
	} else {
		response.ErrorMessage = fmt.Sprintf("Bad request to set secret: %v", err)
	}
	h.sendResponseAndLog(ctx, gatewayId, body, response)
}

func (h *functionsConnectorHandler) handleHeartbeat(ctx context.Context, gatewayId string, requestBody *api.MessageBody, fromAddr ethCommon.Address) {
	var request *OffchainRequest
	err := json.Unmarshal(requestBody.Payload, &request)
	if err != nil {
		h.sendResponseAndLog(ctx, gatewayId, requestBody, internalErrorResponse(fmt.Sprintf("failed to unmarshal request: %v", err)))
		return
	}
	if _, ok := h.allowedHeartbeatInitiators[requestBody.Sender]; !ok {
		h.sendResponseAndLog(ctx, gatewayId, requestBody, internalErrorResponse("sender not allowed to send heartbeat requests"))
		return
	}
	if !bytes.Equal(request.RequestInitiator, fromAddr.Bytes()) {
		h.sendResponseAndLog(ctx, gatewayId, requestBody, internalErrorResponse("RequestInitiator doesn't match sender"))
		return
	}
	if !bytes.Equal(request.SubscriptionOwner, fromAddr.Bytes()) {
		h.sendResponseAndLog(ctx, gatewayId, requestBody, internalErrorResponse("SubscriptionOwner doesn't match sender"))
		return
	}
	if request.Timestamp < uint64(time.Now().Unix())-uint64(h.requestTimeoutSec) {
		h.sendResponseAndLog(ctx, gatewayId, requestBody, internalErrorResponse("Request is too old"))
		return
	}

	internalId := InternalId(fromAddr.Bytes(), request.RequestId)
	request.RequestId = internalId[:]
	h.lggr.Infow("handling offchain heartbeat", "messageId", requestBody.MessageId, "internalId", internalId, "sender", requestBody.Sender)
	h.mu.Lock()
	response, ok := h.heartbeatRequests[internalId]
	if !ok { // new request
		response = &HeartbeatResponse{
			Status:     RequestStatePending,
			ReceivedTs: uint64(time.Now().Unix()),
		}
		h.cacheNewRequestLocked(internalId, response)
		h.shutdownWaitGroup.Add(1)
		go h.handleOffchainRequest(request)
	}
	responseToSend := *response
	h.mu.Unlock()
	requestBody.Receiver = requestBody.Sender
	h.sendResponseAndLog(ctx, gatewayId, requestBody, responseToSend)
}

func internalErrorResponse(internalError string) HeartbeatResponse {
	return HeartbeatResponse{
		Status:        RequestStateInternalError,
		InternalError: internalError,
	}
}

func (h *functionsConnectorHandler) handleOffchainRequest(request *OffchainRequest) {
	defer h.shutdownWaitGroup.Done()
	stopCtx, _ := h.chStop.NewCtx()
	ctx, cancel := context.WithTimeout(stopCtx, time.Duration(h.requestTimeoutSec)*time.Second)
	defer cancel()
	err := h.listener.HandleOffchainRequest(ctx, request)
	if err != nil {
		h.lggr.Errorw("internal error while processing", "id", request.RequestId, "err", err)
		h.mu.Lock()
		defer h.mu.Unlock()
		state, ok := h.heartbeatRequests[RequestID(request.RequestId)]
		if !ok {
			h.lggr.Errorw("request unexpectedly disappeared from local cache", "id", request.RequestId)
			return
		}
		state.CompletedTs = uint64(time.Now().Unix())
		state.Status = RequestStateInternalError
		state.InternalError = err.Error()
	} else {
		// no error - results will be sent to OCR aggregation and returned via reportLoop()
		h.lggr.Infow("request processed successfully, waiting for aggregation ...", "id", request.RequestId)
	}
}

// Listen to OCR reports passed from the plugin and process them against a local cache of requests.
func (h *functionsConnectorHandler) reportLoop() {
	defer h.shutdownWaitGroup.Done()
	for {
		select {
		case report := <-h.offchainTransmitter.ReportChannel():
			h.lggr.Infow("received report", "requestId", report.RequestId, "resultLen", len(report.Result), "errorLen", len(report.Error))
			if len(report.RequestId) != RequestIDLength {
				h.lggr.Errorw("report has invalid requestId", "requestId", report.RequestId)
				continue
			}
			h.mu.Lock()
			cachedResponse, ok := h.heartbeatRequests[RequestID(report.RequestId)]
			if !ok {
				h.lggr.Infow("received report for unknown request, caching it", "id", report.RequestId)
				cachedResponse = &HeartbeatResponse{}
				h.cacheNewRequestLocked(RequestID(report.RequestId), cachedResponse)
			}
			cachedResponse.CompletedTs = uint64(time.Now().Unix())
			cachedResponse.Status = RequestStateComplete
			cachedResponse.Response = report
			h.mu.Unlock()
		case <-h.chStop:
			h.lggr.Info("exiting reportLoop")
			return
		}
	}
}

func (h *functionsConnectorHandler) cacheNewRequestLocked(requestId RequestID, response *HeartbeatResponse) {
	// remove oldest requests
	for len(h.orderedRequests) >= HeartbeatCacheSize {
		delete(h.heartbeatRequests, h.orderedRequests[0])
		h.orderedRequests = h.orderedRequests[1:]
	}
	h.heartbeatRequests[requestId] = response
	h.orderedRequests = append(h.orderedRequests, requestId)
}

func (h *functionsConnectorHandler) sendResponseAndLog(ctx context.Context, gatewayId string, requestBody *api.MessageBody, payload any) {
	err := h.sendResponse(ctx, gatewayId, requestBody, payload)
	if err != nil {
		h.lggr.Errorw("failed to send response to gateway", "id", gatewayId, "err", err)
	} else {
		h.lggr.Debugw("sent to gateway", "id", gatewayId, "messageId", requestBody.MessageId, "donId", requestBody.DonId, "method", requestBody.Method)
	}
}

func (h *functionsConnectorHandler) sendResponse(ctx context.Context, gatewayId string, requestBody *api.MessageBody, payload any) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: requestBody.MessageId,
			DonId:     requestBody.DonId,
			Method:    requestBody.Method,
			Receiver:  requestBody.Sender,
			Payload:   payloadJson,
		},
	}
	if err = msg.Sign(h.signerKey); err != nil {
		return err
	}
	return h.connector.SendToGateway(ctx, gatewayId, msg)
}
