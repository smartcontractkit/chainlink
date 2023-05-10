package threshold

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func Test_decryptionQueue_NewThresholdDecryptor(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1001, 1002, lggr)

	assert.Equal(t, 5, dq.maxQueueLength)
	assert.Equal(t, 1001, dq.maxCiphertextBytes)
	assert.Equal(t, time.Duration(1002), dq.completedRequestsCacheTimeout)
}

func Test_decryptionQueue_Decrypt_ReturnResultAfterCallingDecrypt(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1000, testutils.WaitTimeout(t), lggr)

	go func() {
		waitForPendingRequestToBeAdded(t, dq, []byte("1"))
		dq.ReturnResult([]byte("1"), []byte("decrypted"))
	}()

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	pt, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	require.NoError(t, err)
	if !reflect.DeepEqual(pt, []byte("decrypted")) {
		t.Error("did not get expected result")
	}
}

func Test_decryptionQueue_Decrypt_CiphertextTooLarge(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 10, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("largeciphertext"))
	assert.Equal(t, err.Error(), "ciphertext too large")
}

func Test_decryptionQueue_Decrypt_DuplicateCiphertextId(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	go func() {
		_, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("1"))

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	assert.Equal(t, err.Error(), "ciphertextId must be unique")
}

func Test_decryptionQueue_Decrypt_ContextCancelled(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 100, lggr)

	ctx, cancel := context.WithTimeout(testutils.Context(t), time.Duration(100)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("2"), []byte("encrypted"))
	assert.Equal(t, err.Error(), "context provided by caller was cancelled")
}

func Test_decryptionQueue_Decrypt_QueueFull(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("4"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("4"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	_, err := dq.Decrypt(ctx2, []byte("3"), []byte("encrypted"))
	assert.Equal(t, err.Error(), "queue is full")
}

func Test_decryptionQueue_GetRequests(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(3, 1000, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("5"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("5"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	go func() {
		_, err := dq.Decrypt(ctx2, []byte("6"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("6"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{
		{[]byte("5"), []byte("encrypted")},
		{[]byte("6"), []byte("encrypted")},
	}

	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get the expected requests")
	}
}

func Test_decryptionQueue_GetCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(3, 1000, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	go func() {
		_, err := dq.Decrypt(ctx, []byte("7"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("7"))

	ct, err := dq.GetCiphertext([]byte("7"))
	require.NoError(t, err)
	if !reflect.DeepEqual(ct, []byte("encrypted")) {
		t.Error("did not get the expected requests")
	}
}

func Test_decryptionQueue_GetCiphertext_CiphertextNotFound(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(3, 1000, testutils.WaitTimeout(t), lggr)

	_, err := dq.GetCiphertext([]byte("8"))
	assert.Equal(t, err.Error(), "ciphertext not found")
}

func Test_decryptionQueue_Decrypt_DecryptCalledAfterReadyResult(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(2, 1000, testutils.WaitTimeout(t), lggr)

	dq.ReturnResult([]byte("9"), []byte("decrypted"))

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	pt, err := dq.Decrypt(ctx, []byte("9"), []byte("encrypted"))
	require.NoError(t, err)
	if !reflect.DeepEqual(pt, []byte("decrypted")) {
		t.Error("did not get expected plaintext")
	}
}

func Test_decryptionQueue_ReadyResult_ExpireRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(2, 1000, 100, lggr)

	dq.ReturnResult([]byte("9"), []byte("decrypted"))

	waitForCompletedRequestToBeAdded(t, dq, []byte("9"))

	ctx, cancel := context.WithTimeout(testutils.Context(t), time.Duration(100)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("9"), []byte("encrypted"))
	assert.Equal(t, err.Error(), "context provided by caller was cancelled")
}

func Test_decryptionQueue_Decrypt_CleanupSuccessfulRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(2, 1000, testutils.WaitTimeout(t), lggr)

	dq.ReturnResult([]byte("10"), []byte("decrypted"))

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	_, err1 := dq.Decrypt(ctx1, []byte("10"), []byte("encrypted")) // This will remove the decrypted result to completedRequests
	require.NoError(t, err1)

	ctx2, cancel2 := context.WithTimeout(testutils.Context(t), time.Duration(100)*time.Millisecond)
	defer cancel2()

	_, err2 := dq.Decrypt(ctx2, []byte("10"), []byte("encrypted"))
	assert.Equal(t, err2.Error(), "context provided by caller was cancelled")
}

func Test_decryptionQueue_Decrypt_HandleClosedChannelWithoutPlaintextResponse(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1000, testutils.WaitTimeout(t), lggr)

	go func() {
		waitForPendingRequestToBeAdded(t, dq, []byte("1"))
		close(dq.pendingRequests[string([]byte("1"))].chPlaintext)
	}()

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	assert.Equal(t, err.Error(), "pending decryption request for ciphertextId 1 was closed without a response")
}

func Test_decryptionQueue_GetRequests_RequestsCountLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("11"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	go func() {
		_, err := dq.Decrypt(ctx2, []byte("12"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("12"))

	ctx3, cancel3 := context.WithCancel(testutils.Context(t))
	defer cancel3()

	go func() {
		_, err := dq.Decrypt(ctx3, []byte("13"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("13"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{
		{[]byte("11"), []byte("encrypted")},
		{[]byte("12"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_GetRequests_TotalBytesLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 10, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("11"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	go func() {
		_, err := dq.Decrypt(ctx2, []byte("12"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("12"))

	ctx3, cancel3 := context.WithCancel(testutils.Context(t))
	defer cancel3()

	go func() {
		_, err := dq.Decrypt(ctx3, []byte("13"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("13"))

	requests := dq.GetRequests(4, 30)
	expected := []DecryptionRequest{
		{[]byte("11"), []byte("encrypted")},
		{[]byte("12"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_GetRequests_PendingRequestQueueShorterThanRequestCountLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	go func() {
		_, err := dq.Decrypt(ctx, []byte("11"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{
		{[]byte("11"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_GetRequests_ExpiredRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))

	go func() {
		_, err := dq.Decrypt(ctx, []byte("11"), []byte("encrypted"))
		require.Equal(t, err.Error(), "context provided by caller was cancelled")
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))
	cancel() // Context cancellation should expire the pending request
	waitForPendingRequestToBeRemoved(t, dq, []byte("11"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_Start(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	err := dq.Start(ctx)

	require.NoError(t, err)
}

func Test_decryptionQueue_Close(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, testutils.WaitTimeout(t), lggr)

	dq.ReturnResult([]byte("14"), []byte("decrypted"))

	err := dq.Close()

	require.NoError(t, err)
}

func waitForPendingRequestToBeAdded(t *testing.T, dq *decryptionQueue, ciphertextId CiphertextId) {
	NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.RLock()
		_, exists := dq.pendingRequests[string(ciphertextId)]
		dq.mu.RUnlock()
		return exists
	}, testutils.WaitTimeout(t), "10ms").Should(BeTrue(), "pending request should be added")
}

func waitForPendingRequestToBeRemoved(t *testing.T, dq *decryptionQueue, ciphertextId CiphertextId) {
	NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.RLock()
		_, exists := dq.pendingRequests[string(ciphertextId)]
		dq.mu.RUnlock()
		return exists
	}, testutils.WaitTimeout(t), "10ms").Should(BeFalse(), "pending request should be removed")
}

func waitForCompletedRequestToBeAdded(t *testing.T, dq *decryptionQueue, ciphertextId CiphertextId) {
	NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.RLock()
		_, exists := dq.completedRequests[string(ciphertextId)]
		dq.mu.RUnlock()
		return exists
	}, testutils.WaitTimeout(t), "10ms").Should(BeFalse(), "completed request should be removed")
}
