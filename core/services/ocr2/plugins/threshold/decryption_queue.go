package threshold

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type CiphertextId = []byte

type ThresholdDecryptor interface {
	Decrypt(ctx context.Context, ciphertextId CiphertextId, ciphertext []byte) ([]byte, error)
}

// This interface will be replaced by the Threshold Decryption Plugin when it is ready
type DecryptionQueuingService interface {
	GetRequests(requestCountLimit int, totalBytesLimit int) []DecryptionRequest
	GetCiphertext(ciphertextId CiphertextId) ([]byte, error)
	ReturnResult(ciphertextId CiphertextId, plaintext []byte)
}

type DecryptionRequest struct {
	ciphertextId CiphertextId
	ciphertext   []byte
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
	completedRequestsCacheTimeout time.Duration
	pendingRequestQueue           []CiphertextId
	pendingRequests               map[string]pendingRequest
	completedRequests             map[string]completedRequest
	mu                            sync.RWMutex
	lggr                          logger.Logger
}

var (
	_ ThresholdDecryptor       = &decryptionQueue{}
	_ DecryptionQueuingService = &decryptionQueue{}
	_ job.ServiceCtx           = &decryptionQueue{}
)

func NewDecryptionQueue(maxQueueLength int, maxCiphertextBytes int, completedRequestsCacheTimeout time.Duration, lggr logger.Logger) *decryptionQueue {
	dq := decryptionQueue{
		maxQueueLength,
		maxCiphertextBytes,
		completedRequestsCacheTimeout,
		[]CiphertextId{},
		make(map[string]pendingRequest),
		make(map[string]completedRequest),
		sync.RWMutex{},
		lggr,
	}
	return &dq
}

func (dq *decryptionQueue) Decrypt(ctx context.Context, ciphertextId CiphertextId, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) > dq.maxCiphertextBytes {
		return nil, errors.New("ciphertext too large")
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
		return nil, fmt.Errorf("pending decryption request for ciphertextId %s was closed without a response", string(ciphertextId))
	case <-ctx.Done():
		dq.mu.Lock()
		defer dq.mu.Unlock()
		delete(dq.pendingRequests, string(ciphertextId))
		return nil, errors.New("context provided by caller was cancelled")
	}
}

func (dq *decryptionQueue) getResult(ciphertextId CiphertextId, ciphertext []byte) (<-chan []byte, error) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	chPlaintext := make(chan []byte, 1)

	req, ok := dq.completedRequests[string(ciphertextId)]
	if ok {
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

	return chPlaintext, nil
}

func (dq *decryptionQueue) GetRequests(requestCountLimit int, totalBytesLimit int) []DecryptionRequest {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	requests := make([]DecryptionRequest, 0, requestCountLimit)
	totalBytes := 0
	indicesToRemove := make(map[int]struct{})

	for i := 0; len(requests) < requestCountLimit; i++ {
		if i >= len(dq.pendingRequestQueue) {
			break
		}

		requestId := dq.pendingRequestQueue[i]
		pendingRequest, exists := dq.pendingRequests[string(requestId)]

		if !exists {
			indicesToRemove[i] = struct{}{}
			continue
		}

		requestToAdd := DecryptionRequest{
			requestId,
			pendingRequest.ciphertext,
		}

		requestTotalLen := len(requestId) + len(pendingRequest.ciphertext)

		if (totalBytes + requestTotalLen) > totalBytesLimit {
			break
		}

		requests = append(requests, requestToAdd)
		totalBytes += requestTotalLen
	}

	dq.pendingRequestQueue = removeMultipleIndices(dq.pendingRequestQueue, indicesToRemove)

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

func (dq *decryptionQueue) GetCiphertext(ciphertextId CiphertextId) ([]byte, error) {
	dq.mu.RLock()
	defer dq.mu.RUnlock()

	req, ok := dq.pendingRequests[string(ciphertextId)]
	if !ok {
		return nil, errors.New("ciphertext not found")
	}

	return req.ciphertext, nil
}

func (dq *decryptionQueue) ReturnResult(ciphertextId CiphertextId, plaintext []byte) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	req, ok := dq.pendingRequests[string(ciphertextId)]
	if ok {
		req.chPlaintext <- plaintext
		close(req.chPlaintext)
		delete(dq.pendingRequests, string(ciphertextId))
	} else {
		// Cache plaintext result in completedRequests map for cacheTimeoutMs to account for delayed Decrypt() calls
		timer := time.AfterFunc(dq.completedRequestsCacheTimeout, func() {
			dq.mu.Lock()
			delete(dq.completedRequests, string(ciphertextId))
			dq.mu.Unlock()
		})

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
