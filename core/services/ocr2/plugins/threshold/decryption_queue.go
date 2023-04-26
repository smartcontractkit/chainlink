package decryption_queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type CiphertextId = []byte

type ThresholdDecryptor interface {
	Decrypt(ctx context.Context, ciphertextId CiphertextId, ciphertext []byte) ([]byte, error)
}

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
	pendingChan         chan []byte
	ciphertext          []byte
	expirationTimestamp uint64
}

type completedRequest struct {
	plaintext           []byte
	expirationTimestamp uint64
	timer               *time.Timer
}

type decryptionQueue struct {
	maxQueueLength                  uint32
	maxCiphertextBytes              uint32
	completedRequestsCacheTimeoutMs uint64
	pendingRequestQueue             []CiphertextId
	pendingRequests                 map[string]pendingRequest
	completedRequests               map[string]completedRequest
	mu                              sync.Mutex
}

var (
	_ ThresholdDecryptor       = &decryptionQueue{}
	_ DecryptionQueuingService = &decryptionQueue{}
	_ job.ServiceCtx           = &decryptionQueue{}
)

func NewThresholdDecryptor(maxQueueLength uint32, maxCiphertextBytes uint32, completedRequestsCacheTimeoutMs uint64) *decryptionQueue {
	dq := decryptionQueue{
		maxQueueLength,
		maxCiphertextBytes,
		completedRequestsCacheTimeoutMs,
		[]CiphertextId{},
		make(map[string]pendingRequest),
		make(map[string]completedRequest),
		sync.Mutex{},
	}
	return &dq
}

func (dq *decryptionQueue) Decrypt(ctx context.Context, ciphertextId CiphertextId, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) > int(dq.maxCiphertextBytes) {
		return nil, errors.New("ciphertext too large")
	}

	pendingChan, err := dq.getResult(ciphertextId, ciphertext)
	if err != nil {
		return nil, err
	}

	select {
	case pt := <-pendingChan:
		return pt, nil
	case <-ctx.Done():
		dq.mu.Lock()
		defer dq.mu.Unlock()
		delete(dq.pendingRequests, string(ciphertextId))
		return nil, errors.New("cancelled")
	}
}

func (dq *decryptionQueue) getResult(ciphertextId CiphertextId, ciphertext []byte) (chan []byte, error) {

	dq.mu.Lock()
	defer dq.mu.Unlock()

	pendingChan := make(chan []byte, 1)

	req, ok := dq.completedRequests[string(ciphertextId)]
	if ok {
		pendingChan <- req.plaintext
		req.timer.Stop()
		delete(dq.completedRequests, string(ciphertextId))
		return pendingChan, nil
	}

	_, isDuplicateId := dq.pendingRequests[string(ciphertextId)]
	if isDuplicateId {
		return nil, errors.New("Decrypt ciphertextId must be unique")
	}

	if uint32(len(dq.pendingRequestQueue)) >= dq.maxQueueLength {
		return nil, errors.New("queue is full")
	}
	dq.pendingRequestQueue = append(dq.pendingRequestQueue, ciphertextId)

	dq.pendingRequests[string(ciphertextId)] = pendingRequest{
		pendingChan,
		ciphertext,
		uint64(time.Now().Unix()) + dq.completedRequestsCacheTimeoutMs/1000,
	}

	return pendingChan, nil
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

		additionalBytes := len(requestId) + len(pendingRequest.ciphertext)

		if (totalBytes + additionalBytes) > int(totalBytesLimit) {
			continue
		}

		requests = append(requests, requestToAdd)
		totalBytes += additionalBytes
	}

	return requests
}

func (dq *decryptionQueue) GetCiphertext(ciphertextId CiphertextId) ([]byte, error) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

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
		req.pendingChan <- plaintext
		delete(dq.pendingRequests, string(ciphertextId))
	} else {
		timer := time.AfterFunc(time.Duration(dq.completedRequestsCacheTimeoutMs)*time.Millisecond, func() {
			dq.mu.Lock()
			delete(dq.completedRequests, string(ciphertextId))
			dq.mu.Unlock()
		})

		dq.completedRequests[string(ciphertextId)] = completedRequest{
			plaintext,
			uint64(time.Now().Unix()) + dq.completedRequestsCacheTimeoutMs/1000,
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
