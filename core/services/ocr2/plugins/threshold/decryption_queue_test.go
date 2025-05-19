package threshold

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func Test_decryptionQueue_NewThresholdDecryptor(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1001, 64, 1002, lggr)

	assert.Equal(t, 5, dq.maxQueueLength)
	assert.Equal(t, 1001, dq.maxCiphertextBytes)
	assert.Equal(t, time.Duration(1002), dq.completedRequestsCacheTimeout)
}

func Test_decryptionQueue_Decrypt_ReturnResultAfterCallingDecrypt(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1000, 64, testutils.WaitTimeout(t), lggr)

	go func() {
		waitForPendingRequestToBeAdded(t, dq, []byte("1"))
		dq.SetResult([]byte("1"), []byte("decrypted"), nil)
	}()

	ctx := testutils.Context(t)

	pt, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	require.NoError(t, err)
	if !reflect.DeepEqual(pt, []byte("decrypted")) {
		t.Error("did not get expected result")
	}
}

func Test_decryptionQueue_Decrypt_CiphertextIdTooLarge(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 16, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	_, err := dq.Decrypt(ctx, []byte("largeCiphertextId"), []byte("ciphertext"))
	assert.Equal(t, "ciphertextId too large", err.Error())
}

func Test_decryptionQueue_Decrypt_EmptyCiphertextId(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	_, err := dq.Decrypt(ctx, []byte(""), []byte("ciphertext"))
	assert.Equal(t, "ciphertextId is empty", err.Error())
}

func Test_decryptionQueue_Decrypt_CiphertextTooLarge(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 10, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("largeciphertext"))
	assert.Equal(t, "ciphertext too large", err.Error())
}

func Test_decryptionQueue_Decrypt_EmptyCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	_, err := dq.Decrypt(ctx, []byte("1"), []byte(""))
	assert.Equal(t, "ciphertext is empty", err.Error())
}

func Test_decryptionQueue_Decrypt_DuplicateCiphertextId(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	go func() {
		_, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("1"))

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	assert.Equal(t, "ciphertextId must be unique", err.Error())
}

func Test_decryptionQueue_Decrypt_ContextCancelled(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 64, 100, lggr)

	ctx, cancel := context.WithTimeout(testutils.Context(t), time.Duration(100)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("2"), []byte("encrypted"))
	assert.Equal(t, "context provided by caller was cancelled", err.Error())
}

func Test_decryptionQueue_Decrypt_QueueFull(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(1, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("4"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("4"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	_, err := dq.Decrypt(ctx2, []byte("3"), []byte("encrypted"))
	assert.Equal(t, "queue is full", err.Error())
}

func Test_decryptionQueue_GetRequests(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(3, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("5"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("5"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	go func() {
		_, err := dq.Decrypt(ctx2, []byte("6"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("6"))

	requests := dq.GetRequests(2, 1000)
	expected := []decryptionPlugin.DecryptionRequest{
		{CiphertextId: []byte("5"), Ciphertext: []byte("encrypted")},
		{CiphertextId: []byte("6"), Ciphertext: []byte("encrypted")},
	}

	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get the expected requests")
	}
}

func Test_decryptionQueue_GetCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(3, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	go func() {
		_, err := dq.Decrypt(ctx, []byte("7"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
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
	dq := NewDecryptionQueue(3, 1000, 64, testutils.WaitTimeout(t), lggr)

	_, err := dq.GetCiphertext([]byte{0xa5})
	assert.ErrorIs(t, err, decryptionPlugin.ErrNotFound)
}

func Test_decryptionQueue_Decrypt_DecryptCalledAfterReadyResult(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(2, 1000, 64, testutils.WaitTimeout(t), lggr)

	dq.SetResult([]byte("9"), []byte("decrypted"), nil)

	ctx := testutils.Context(t)

	pt, err := dq.Decrypt(ctx, []byte("9"), []byte("encrypted"))
	require.NoError(t, err)
	if !reflect.DeepEqual(pt, []byte("decrypted")) {
		t.Error("did not get expected plaintext")
	}
}

func Test_decryptionQueue_ReadyResult_ExpireRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(2, 1000, 64, 100, lggr)

	dq.SetResult([]byte("9"), []byte("decrypted"), nil)

	waitForCompletedRequestToBeAdded(t, dq, []byte("9"))

	ctx, cancel := context.WithTimeout(testutils.Context(t), time.Duration(100)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("9"), []byte("encrypted"))
	assert.Equal(t, "context provided by caller was cancelled", err.Error())
}

func Test_decryptionQueue_Decrypt_CleanupSuccessfulRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(2, 1000, 64, testutils.WaitTimeout(t), lggr)

	dq.SetResult([]byte("10"), []byte("decrypted"), nil)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	_, err1 := dq.Decrypt(ctx1, []byte("10"), []byte("encrypted")) // This will remove the decrypted result to completedRequests
	require.NoError(t, err1)

	ctx2, cancel2 := context.WithTimeout(testutils.Context(t), time.Duration(100)*time.Millisecond)
	defer cancel2()

	_, err2 := dq.Decrypt(ctx2, []byte("10"), []byte("encrypted"))
	assert.Equal(t, "context provided by caller was cancelled", err2.Error())
}

func Test_decryptionQueue_Decrypt_UserErrorDuringDecryption(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1000, 64, testutils.WaitTimeout(t), lggr)
	ciphertextId := []byte{0x12, 0x0f}

	go func() {
		waitForPendingRequestToBeAdded(t, dq, ciphertextId)
		dq.SetResult(ciphertextId, nil, decryptionPlugin.ErrAggregation)
	}()

	ctx := testutils.Context(t)

	_, err := dq.Decrypt(ctx, ciphertextId, []byte("encrypted"))
	assert.Equal(t, "pending decryption request for ciphertextId 0x120f was closed without a response", err.Error())
}

func Test_decryptionQueue_Decrypt_HandleClosedChannelWithoutPlaintextResponse(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(5, 1000, 64, testutils.WaitTimeout(t), lggr)
	ciphertextId := []byte{0x00, 0xff}

	go func() {
		waitForPendingRequestToBeAdded(t, dq, ciphertextId)
		close(dq.pendingRequests[string(ciphertextId)].chPlaintext)
	}()

	ctx := testutils.Context(t)

	_, err := dq.Decrypt(ctx, ciphertextId, []byte("encrypted"))
	assert.Equal(t, "pending decryption request for ciphertextId 0x00ff was closed without a response", err.Error())
}

func Test_decryptionQueue_GetRequests_RequestsCountLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("11"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	go func() {
		_, err := dq.Decrypt(ctx2, []byte("12"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("12"))

	ctx3, cancel3 := context.WithCancel(testutils.Context(t))
	defer cancel3()

	go func() {
		_, err := dq.Decrypt(ctx3, []byte("13"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("13"))

	requests := dq.GetRequests(2, 1000)
	expected := []decryptionPlugin.DecryptionRequest{
		{CiphertextId: []byte("11"), Ciphertext: []byte("encrypted")},
		{CiphertextId: []byte("12"), Ciphertext: []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_GetRequests_TotalBytesLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 10, 64, testutils.WaitTimeout(t), lggr)

	ctx1, cancel1 := context.WithCancel(testutils.Context(t))
	defer cancel1()

	go func() {
		_, err := dq.Decrypt(ctx1, []byte("11"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	ctx2, cancel2 := context.WithCancel(testutils.Context(t))
	defer cancel2()

	go func() {
		_, err := dq.Decrypt(ctx2, []byte("12"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("12"))

	ctx3, cancel3 := context.WithCancel(testutils.Context(t))
	defer cancel3()

	go func() {
		_, err := dq.Decrypt(ctx3, []byte("13"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("13"))

	requests := dq.GetRequests(4, 30)
	expected := []decryptionPlugin.DecryptionRequest{
		{CiphertextId: []byte("11"), Ciphertext: []byte("encrypted")},
		{CiphertextId: []byte("12"), Ciphertext: []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_GetRequests_PendingRequestQueueShorterThanRequestCountLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	go func() {
		_, err := dq.Decrypt(ctx, []byte("11"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	requests := dq.GetRequests(2, 1000)
	expected := []decryptionPlugin.DecryptionRequest{
		{CiphertextId: []byte("11"), Ciphertext: []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_GetRequests_ExpiredRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx, cancel := context.WithCancel(testutils.Context(t))

	go func() {
		_, err := dq.Decrypt(ctx, []byte("11"), []byte("encrypted"))
		require.Equal(t, "context provided by caller was cancelled", err.Error())
	}()

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))
	cancel() // Context cancellation should expire the pending request
	waitForPendingRequestToBeRemoved(t, dq, []byte("11"))

	requests := dq.GetRequests(2, 1000)
	expected := []decryptionPlugin.DecryptionRequest{}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("did not get expected requests")
	}
}

func Test_decryptionQueue_Start(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, 64, testutils.WaitTimeout(t), lggr)

	ctx := testutils.Context(t)

	err := dq.Start(ctx)

	require.NoError(t, err)
}

func Test_decryptionQueue_Close(t *testing.T) {
	lggr := logger.TestLogger(t)
	dq := NewDecryptionQueue(4, 1000, 64, testutils.WaitTimeout(t), lggr)

	dq.SetResult([]byte("14"), []byte("decrypted"), nil)

	err := dq.Close()

	require.NoError(t, err)
}

func waitForPendingRequestToBeAdded(t *testing.T, dq *decryptionQueue, ciphertextId decryptionPlugin.CiphertextId) {
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.RLock()
		_, exists := dq.pendingRequests[string(ciphertextId)]
		dq.mu.RUnlock()
		return exists
	}, testutils.WaitTimeout(t), "10ms").Should(gomega.BeTrue(), "pending request should be added")
}

func waitForPendingRequestToBeRemoved(t *testing.T, dq *decryptionQueue, ciphertextId decryptionPlugin.CiphertextId) {
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.RLock()
		_, exists := dq.pendingRequests[string(ciphertextId)]
		dq.mu.RUnlock()
		return exists
	}, testutils.WaitTimeout(t), "10ms").Should(gomega.BeFalse(), "pending request should be removed")
}

func waitForCompletedRequestToBeAdded(t *testing.T, dq *decryptionQueue, ciphertextId decryptionPlugin.CiphertextId) {
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.RLock()
		_, exists := dq.completedRequests[string(ciphertextId)]
		dq.mu.RUnlock()
		return exists
	}, testutils.WaitTimeout(t), "10ms").Should(gomega.BeFalse(), "completed request should be removed")
}
