package decryption_queue

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
	GetRequests(requestCountLimit uint32, totalBytesLimit uint32) []DecryptionRequest
	GetCiphertext(ciphertextId CiphertextId) ([]byte, error)
	ResultReady(ciphertextId CiphertextId, plaintext []byte)
}

type DecryptionRequest struct {
	ciphertextId CiphertextId
	ciphertext   []byte
}

type pendingRequest struct {
	chPlaintext chan []byte
	ciphertext  []byte
}

type completedRequest struct {
	plaintext []byte
	timer     *time.Timer
}

type decryptionQueue struct {
	maxQueueLength                  uint32
	maxCiphertextBytes              uint32
	completedRequestsCacheTimeoutMs uint64
	pendingRequestQueue             []CiphertextId
	pendingRequests                 map[string]pendingRequest
	completedRequests               map[string]completedRequest
	mu                              sync.RWMutex
	lggr                            logger.Logger
}

var (
	_ ThresholdDecryptor       = &decryptionQueue{}
	_ DecryptionQueuingService = &decryptionQueue{}
	_ job.ServiceCtx           = &decryptionQueue{}
)

func NewThresholdDecryptor(maxQueueLength uint32, maxCiphertextBytes uint32, completedRequestsCacheTimeoutMs uint64, lggr logger.Logger) *decryptionQueue {
	dq := decryptionQueue{
		maxQueueLength,
		maxCiphertextBytes,
		completedRequestsCacheTimeoutMs,
		[]CiphertextId{},
		make(map[string]pendingRequest),
		make(map[string]completedRequest),
		sync.RWMutex{},
		lggr,
	}
	return &dq
}

func (dq *decryptionQueue) Decrypt(ctx context.Context, ciphertextId CiphertextId, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) > int(dq.maxCiphertextBytes) {
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

func (dq *decryptionQueue) getResult(ciphertextId CiphertextId, ciphertext []byte) (chan []byte, error) {
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

	if uint32(len(dq.pendingRequestQueue)) >= dq.maxQueueLength {
		return nil, errors.New("queue is full")
	}
	dq.pendingRequestQueue = append(dq.pendingRequestQueue, ciphertextId)

	dq.pendingRequests[string(ciphertextId)] = pendingRequest{
		chPlaintext,
		ciphertext,
	}

	return chPlaintext, nil
}

func (dq *decryptionQueue) GetRequests(requestCountLimit uint32, totalBytesLimit uint32) []DecryptionRequest {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	requests := make([]DecryptionRequest, 0, requestCountLimit)
	totalBytes := 0

	for i := 0; len(requests) < int(requestCountLimit); i++ {
		if i >= len(dq.pendingRequestQueue) {
			break
		}

		requestId := dq.pendingRequestQueue[i]
		pendingRequest, exists := dq.pendingRequests[string(requestId)]

		if !exists {
			dq.pendingRequestQueue = append(dq.pendingRequestQueue[:i], dq.pendingRequestQueue[i+1:]...)
			i--
			continue
		}

		requestToAdd := DecryptionRequest{
			requestId,
			pendingRequest.ciphertext,
		}

		requestTotalLen := len(requestId) + len(pendingRequest.ciphertext)

		if (totalBytes + requestTotalLen) > int(totalBytesLimit) {
			break
		}

		requests = append(requests, requestToAdd)
		totalBytes += requestTotalLen
	}

	return requests
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

func (dq *decryptionQueue) ResultReady(ciphertextId CiphertextId, plaintext []byte) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	req, ok := dq.pendingRequests[string(ciphertextId)]
	if ok {
		req.chPlaintext <- plaintext
		delete(dq.pendingRequests, string(ciphertextId))
	} else {
		// Cache plaintext result in completedRequests map for cacheTimeoutMs to account for delayed Decrypt() calls
		timer := time.AfterFunc(time.Duration(dq.completedRequestsCacheTimeoutMs)*time.Millisecond, func() {
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
