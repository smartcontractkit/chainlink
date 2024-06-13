package requests

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type responseCacheEntry struct {
	response  *Response
	entryTime time.Time
}

type Handler struct {
	services.StateMachine

	lggr logger.Logger

	store *Store

	pendingRequests map[string]*Request

	responseCache   map[string]*responseCacheEntry
	cacheExpiryTime time.Duration

	responseCh chan *Response
	requestCh  chan *Request

	clock clockwork.Clock

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewHandler(lggr logger.Logger, s *Store, clock clockwork.Clock, responseExpiryTime time.Duration) *Handler {
	return &Handler{
		lggr:            lggr,
		store:           s,
		pendingRequests: map[string]*Request{},
		responseCache:   map[string]*responseCacheEntry{},
		responseCh:      make(chan *Response),
		requestCh:       make(chan *Request),
		clock:           clock,
		cacheExpiryTime: responseExpiryTime,
		stopCh:          make(services.StopChan),
	}
}

func (h *Handler) SendResponse(ctx context.Context, resp *Response) {
	select {
	case <-ctx.Done():
		return
	case h.responseCh <- resp:
	}
}

func (h *Handler) SendRequest(ctx context.Context, r *Request) {
	select {
	case <-ctx.Done():
		return
	case h.requestCh <- r:
	}
}

func (h *Handler) Start(_ context.Context) error {
	return h.StartOnce("RequestHandler", func() error {
		h.wg.Add(1)
		go func() {
			defer h.wg.Done()
			ctx, cancel := h.stopCh.NewCtx()
			defer cancel()
			h.worker(ctx)
		}()
		return nil
	})
}

func (h *Handler) Close() error {
	return h.StopOnce("RequestHandler", func() error {
		close(h.stopCh)
		h.wg.Wait()
		return nil
	})
}

func (h *Handler) worker(ctx context.Context) {
	responseCacheExpiryTicker := h.clock.NewTicker(h.cacheExpiryTime)
	defer responseCacheExpiryTicker.Stop()

	// Set to tick at 1 second as this is a sufficient resolution for expiring requests without causing too much overhead
	pendingRequestsExpiryTicker := h.clock.NewTicker(1 * time.Second)
	defer pendingRequestsExpiryTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pendingRequestsExpiryTicker.Chan():
			h.expirePendingRequests(ctx)
		case <-responseCacheExpiryTicker.Chan():
			h.expireCachedResponses()
		case req := <-h.requestCh:
			h.pendingRequests[req.WorkflowExecutionID] = req

			existingResponse := h.responseCache[req.WorkflowExecutionID]
			if existingResponse != nil {
				delete(h.responseCache, req.WorkflowExecutionID)
				h.lggr.Debugw("Found cached response for request", "workflowExecutionID", req.WorkflowExecutionID)
				h.sendResponse(ctx, req, existingResponse.response)
				continue
			}

			if err := h.store.Add(req); err != nil {
				h.lggr.Errorw("failed to add request to store", "err", err)
			}

		case resp := <-h.responseCh:
			req, wasPresent := h.store.evict(resp.WorkflowExecutionID)
			if !wasPresent {
				h.responseCache[resp.WorkflowExecutionID] = &responseCacheEntry{
					response:  resp,
					entryTime: h.clock.Now(),
				}
				h.lggr.Debugw("Caching response without request", "workflowExecutionID", resp.WorkflowExecutionID)
				continue
			}

			h.sendResponse(ctx, req, resp)
		}
	}
}

func (h *Handler) sendResponse(ctx context.Context, req *Request, resp *Response) {
	select {
	case <-ctx.Done():
		return
	case req.CallbackCh <- resp.CapabilityResponse:
		close(req.CallbackCh)
		delete(h.pendingRequests, req.WorkflowExecutionID)
	}
}

func (h *Handler) expirePendingRequests(ctx context.Context) {
	now := h.clock.Now()

	for _, req := range h.pendingRequests {
		if now.After(req.ExpiresAt) {
			resp := &Response{
				WorkflowExecutionID: req.WorkflowExecutionID,
				CapabilityResponse: capabilities.CapabilityResponse{
					Err:   fmt.Errorf("timeout exceeded: could not process request before expiry %s", req.WorkflowExecutionID),
					Value: nil,
				},
			}
			h.store.evict(req.WorkflowExecutionID)
			h.sendResponse(ctx, req, resp)
		}
	}
}

func (h *Handler) expireCachedResponses() {
	for k, v := range h.responseCache {
		if h.clock.Since(v.entryTime) > h.cacheExpiryTime {
			delete(h.responseCache, k)
			h.lggr.Debugw("Expired response", "workflowExecutionID", k)
		}
	}
}
