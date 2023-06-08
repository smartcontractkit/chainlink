package threshold

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type Decryptor interface {
	Decrypt(ctx context.Context, ciphertextId decryptionPlugin.CiphertextId, ciphertext []byte) ([]byte, error)
}

type pendingRequest struct {
	chPlaintext chan<- []byte
	ciphertext  []byte
}

type completedRequest struct {
	plaintext []byte
	timer     *time.Timer
}

type decryptionQueue struct {
	maxQueueLength                int
	maxCiphertextBytes            int
	maxCiphertextIdLen            int
	completedRequestsCacheTimeout time.Duration
	pendingRequestQueue           []decryptionPlugin.CiphertextId
	pendingRequests               map[string]pendingRequest
	completedRequests             map[string]completedRequest
	mu                            sync.RWMutex
	lggr                          logger.Logger
}

var (
	_ Decryptor                                 = &decryptionQueue{}
	_ decryptionPlugin.DecryptionQueuingService = &decryptionQueue{}
	_ job.ServiceCtx                            = &decryptionQueue{}
)

func NewDecryptionQueue(maxQueueLength int, maxCiphertextBytes int, maxCiphertextIdLen int, completedRequestsCacheTimeout time.Duration, lggr logger.Logger) *decryptionQueue {
	dq := decryptionQueue{
		maxQueueLength,
		maxCiphertextBytes,
		maxCiphertextIdLen,
		completedRequestsCacheTimeout,
		[]decryptionPlugin.CiphertextId{},
		make(map[string]pendingRequest),
		make(map[string]completedRequest),
		sync.RWMutex{},
		lggr.Named("decryptionQueue"),
	}
	return &dq
}

func (dq *decryptionQueue) Decrypt(ctx context.Context, ciphertextId decryptionPlugin.CiphertextId, ciphertext []byte) ([]byte, error) {
	if len(ciphertextId) > dq.maxCiphertextIdLen {
		return nil, errors.New("ciphertextId too large")
	}

	if len(ciphertextId) == 0 {
		return nil, errors.New("ciphertextId is empty")
	}

	if len(ciphertext) > dq.maxCiphertextBytes {
		return nil, errors.New("ciphertext too large")
	}

	if len(ciphertext) == 0 {
		return nil, errors.New("ciphertext is empty")
	}

	chPlaintext, err := dq.getResult(ciphertextId, ciphertext)
	if err != nil {
		return nil, err
	}

	select {
	case pt, ok := <-chPlaintext:
		if ok {
			return pt, nil
		}
		return nil, fmt.Errorf("pending decryption request for ciphertextId %s was closed without a response", ciphertextId)
	case <-ctx.Done():
		dq.mu.Lock()
		defer dq.mu.Unlock()
		delete(dq.pendingRequests, string(ciphertextId))
		return nil, errors.New("context provided by caller was cancelled")
	}
}

func (dq *decryptionQueue) getResult(ciphertextId decryptionPlugin.CiphertextId, ciphertext []byte) (<-chan []byte, error) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	chPlaintext := make(chan []byte, 1)

	req, ok := dq.completedRequests[string(ciphertextId)]
	if ok {
		dq.lggr.Debugf("ciphertextId %s was already decrypted by the DON", ciphertextId)
		chPlaintext <- req.plaintext
		req.timer.Stop()
		delete(dq.completedRequests, string(ciphertextId))
		return chPlaintext, nil
	}

	_, isDuplicateId := dq.pendingRequests[string(ciphertextId)]
	if isDuplicateId {
		return nil, errors.New("ciphertextId must be unique")
	}

	if len(dq.pendingRequestQueue) >= dq.maxQueueLength {
		return nil, errors.New("queue is full")
	}
	dq.pendingRequestQueue = append(dq.pendingRequestQueue, ciphertextId)

	dq.pendingRequests[string(ciphertextId)] = pendingRequest{
		chPlaintext,
		ciphertext,
	}
	dq.lggr.Debugf("ciphertextId %s added to pendingRequestQueue")

	return chPlaintext, nil
}

func (dq *decryptionQueue) GetRequests(requestCountLimit int, totalBytesLimit int) []decryptionPlugin.DecryptionRequest {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	requests := make([]decryptionPlugin.DecryptionRequest, 0, requestCountLimit)
	totalBytes := 0
	indicesToRemove := make(map[int]struct{})

	for i := 0; len(requests) < requestCountLimit; i++ {
		if i >= len(dq.pendingRequestQueue) {
			break
		}

		requestId := dq.pendingRequestQueue[i]
		pendingRequest, exists := dq.pendingRequests[string(requestId)]

		if !exists {
			dq.lggr.Debugf("pending decryption request for ciphertextId %s expired", requestId)
			indicesToRemove[i] = struct{}{}
			continue
		}

		requestToAdd := decryptionPlugin.DecryptionRequest{
			CiphertextId: requestId,
			Ciphertext:   pendingRequest.ciphertext,
		}

		requestTotalLen := len(requestId) + len(pendingRequest.ciphertext)

		if (totalBytes + requestTotalLen) > totalBytesLimit {
			dq.lggr.Debug("totalBytesLimit reached in GetRequests")
			break
		}

		requests = append(requests, requestToAdd)
		totalBytes += requestTotalLen
	}

	dq.pendingRequestQueue = removeMultipleIndices(dq.pendingRequestQueue, indicesToRemove)

	dq.lggr.Debug("returing %d of %d total requests awaiting decryption", requestCountLimit, len(dq.pendingRequestQueue))

	return requests
}

func removeMultipleIndices[T any](data []T, indicesToRemove map[int]struct{}) []T {
	filtered := make([]T, 0, len(data)-len(indicesToRemove))

	for i, v := range data {
		if _, exists := indicesToRemove[i]; !exists {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func (dq *decryptionQueue) GetCiphertext(ciphertextId decryptionPlugin.CiphertextId) ([]byte, error) {
	dq.mu.RLock()
	defer dq.mu.RUnlock()

	req, ok := dq.pendingRequests[string(ciphertextId)]
	if !ok {
		return nil, errors.New("ciphertext not found")
	}

	return req.ciphertext, nil
}

func (dq *decryptionQueue) SetResult(ciphertextId decryptionPlugin.CiphertextId, plaintext []byte) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	req, ok := dq.pendingRequests[string(ciphertextId)]
	if ok {
		dq.lggr.Debugf("responding with result for pending decryption request ciphertextId %s", ciphertextId)
		req.chPlaintext <- plaintext
		close(req.chPlaintext)
		delete(dq.pendingRequests, string(ciphertextId))
	} else {
		// Cache plaintext result in completedRequests map for cacheTimeoutMs to account for delayed Decrypt() calls
		timer := time.AfterFunc(dq.completedRequestsCacheTimeout, func() {
			dq.lggr.Debugf("expired decryption result for ciphertextId %s from completedRequests cache", ciphertextId)
			dq.mu.Lock()
			delete(dq.completedRequests, string(ciphertextId))
			dq.mu.Unlock()
		})

		dq.lggr.Debugf("adding decryption result for ciphertextId %s to completedRequests cache", ciphertextId)
		dq.completedRequests[string(ciphertextId)] = completedRequest{
			plaintext,
			timer,
		}
	}
}

func (dq *decryptionQueue) Start(ctx context.Context) error {
	return nil
}

func (dq *decryptionQueue) Close() error {
	for _, completedRequest := range dq.completedRequests {
		completedRequest.timer.Stop()
	}
	return nil
}
