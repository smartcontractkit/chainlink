package common

import (
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

// RequestCache is used to store pending requests and collect incoming responses as they arrive.
// It is parameterized by responseData, which is a service-specific type storing all data needed to aggregate responses.
// Client needs to implement a ResponseProcessor, which is called for every response (see below).
// Additionally, each request has a timeout, after which the netry will be removed from the cache and an error sent to the callback channel.
// All methods are thread-safe.
type RequestCache[T any] interface {
	NewRequest(request *api.Message, callbackCh chan<- handlers.UserCallbackPayload, responseData *T) error
	ProcessResponse(response *api.Message, process ResponseProcessor[T]) error
}

// If aggregated != nil then the aggregated response is ready and the entry will be deleted from RequestCache.
// Otherwise, state will be updated to newState and the entry will remain in cache, awaiting more responses from nodes.
type ResponseProcessor[T any] func(response *api.Message, state *T) (aggregated *handlers.UserCallbackPayload, newState *T, err error)

type requestCache[T any] struct {
	cache        map[globalId]*pendingRequest[T]
	maxCacheSize uint32
	timeout      time.Duration
	mu           sync.Mutex
}

type globalId struct {
	sender string
	id     string
}

type pendingRequest[T any] struct {
	callbackCh   chan<- handlers.UserCallbackPayload
	responseData *T
	timeoutTimer *time.Timer
	mu           sync.Mutex
}

func NewRequestCache[T any](timeout time.Duration, maxCacheSize uint32) RequestCache[T] {
	return &requestCache[T]{cache: make(map[globalId]*pendingRequest[T]), timeout: timeout, maxCacheSize: maxCacheSize}
}

func (c *requestCache[T]) NewRequest(request *api.Message, callbackCh chan<- handlers.UserCallbackPayload, responseData *T) error {
	if request == nil {
		return errors.New("request is nil")
	}
	if responseData == nil {
		return errors.New("responseData is nil")
	}
	key := globalId{request.Body.Sender, request.Body.MessageId}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[key]
	if ok {
		return errors.New("request already exists")
	}
	if len(c.cache) >= int(c.maxCacheSize) {
		return errors.New("request cache is full")
	}
	timer := time.AfterFunc(c.timeout, func() {
		c.deleteAndSendOnce(key, handlers.UserCallbackPayload{Msg: request, ErrMsg: "timeout", ErrCode: api.RequestTimeoutError})
	})
	c.cache[key] = &pendingRequest[T]{callbackCh: callbackCh, responseData: responseData, timeoutTimer: timer}
	return nil
}

// Call ResponseProcessor on the response.
// There are two possible outcomes:
//
//	(a) remove request from cache and send aggregated response to the user
//	(b) update request's responseData and keep it in cache, awaiting more responses from nodes
func (c *requestCache[T]) ProcessResponse(response *api.Message, process ResponseProcessor[T]) error {
	if response == nil {
		return errors.New("response is nil")
	}
	key := globalId{response.Body.Receiver, response.Body.MessageId}
	// retrieve entry
	c.mu.Lock()
	entry, ok := c.cache[key]
	c.mu.Unlock()
	if !ok {
		return errors.New("request not found")
	}
	// process under per-entry lock
	entry.mu.Lock()
	aggregated, newResponseData, err := process(response, entry.responseData)
	if newResponseData != nil {
		entry.responseData = newResponseData
	}
	entry.mu.Unlock()
	if err != nil {
		return err
	}
	if aggregated != nil {
		c.deleteAndSendOnce(key, *aggregated)
	}
	return nil
}

func (c *requestCache[T]) deleteAndSendOnce(key globalId, callbackResponse handlers.UserCallbackPayload) {
	c.mu.Lock()
	entry, deleted := c.cache[key]
	delete(c.cache, key)
	c.mu.Unlock()
	if deleted {
		entry.timeoutTimer.Stop()
		entry.callbackCh <- callbackResponse
		close(entry.callbackCh)
	}
}
