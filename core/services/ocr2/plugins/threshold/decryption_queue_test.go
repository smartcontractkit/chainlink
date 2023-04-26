package decryption_queue

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestNewThresholdDecryptor(t *testing.T) {
	dq := NewThresholdDecryptor(5, 1000, 1000)
	if dq.maxQueueLength != 5 || dq.maxCiphertextBytes != 1000 || dq.completedRequestsCacheTimeoutMs != 1000 {
		t.Error("Init failed to set maxQueueSize and localCacheTimeoutMs")
	}
}

func TestDecryptAndResultReady(t *testing.T) {
	dq := NewThresholdDecryptor(5, 1000, 1000)

	go func() {
		waitForPendingRequestToBeAdded(t, dq, []byte("1"))
		dq.ResultReady([]byte("1"), []byte("decrypted"))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel()

	pt, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(pt, []byte("decrypted")) {
		t.Error("Decrypt did not return the expected plaintext")
	}
}

func TestCiphertextTooLarge(t *testing.T) {
	dq := NewThresholdDecryptor(1, 10, 1000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("largeciphertext"))
	if err == nil || err.Error() != "ciphertext too large" {
		t.Error("Decrypt did not return expected error when ciphertext is too large")
	}
}

func TestDuplicateCiphertextId(t *testing.T) {
	dq := NewThresholdDecryptor(1, 1000, 1000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel()

	//nolint:errcheck
	go dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("1"))

	_, err := dq.Decrypt(ctx, []byte("1"), []byte("encrypted"))
	if err == nil || err.Error() != "Decrypt ciphertextId must be unique" {
		t.Error("Decrypt did not return expected error for duplicate ciphertextId")
	}
}

func TestDecryptTimeout(t *testing.T) {
	dq := NewThresholdDecryptor(1, 1000, 100)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("2"), []byte("encrypted"))
	if err == nil || err.Error() != "cancelled" {
		t.Error("Decrypt did not timeout as expected")
	}
}

func TestDecryptQueueFull(t *testing.T) {
	dq := NewThresholdDecryptor(1, 1000, 1000)

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel1()

	//nolint:errcheck
	go dq.Decrypt(ctx1, []byte("4"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("4"))

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel2()

	_, err := dq.Decrypt(ctx2, []byte("3"), []byte("encrypted"))
	if err == nil || err.Error() != "queue is full" {
		t.Error("Decrypt did not return queue full error as expected")
	}
}

func TestGetRequests(t *testing.T) {
	dq := NewThresholdDecryptor(3, 1000, 1000)

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel1()

	//nolint:errcheck
	go dq.Decrypt(ctx1, []byte("5"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("5"))

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel2()

	//nolint:errcheck
	go dq.Decrypt(ctx2, []byte("6"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("6"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{
		{[]byte("5"), []byte("encrypted")},
		{[]byte("6"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("GetRequests did not return the expected requests")
	}
}

func TestGetCiphertext(t *testing.T) {
	dq := NewThresholdDecryptor(3, 1000, 1000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel()

	//nolint:errcheck
	go dq.Decrypt(ctx, []byte("7"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("7"))

	ct, err := dq.GetCiphertext([]byte("7"))
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(ct, []byte("encrypted")) {
		t.Error("GetCiphertext did not return the expected ciphertext")
	}
}

func TestGetCiphertextNotFound(t *testing.T) {
	dq := NewThresholdDecryptor(3, 1000, 1000)

	_, err := dq.GetCiphertext([]byte("8"))
	if err == nil || err.Error() != "ciphertext not found" {
		t.Error("GetCiphertext did not return the expected error")
	}
}

func TestDecryptFromCompletedRequests(t *testing.T) {
	dq := NewThresholdDecryptor(2, 1000, 1000)

	dq.ResultReady([]byte("9"), []byte("decrypted"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel()

	pt, err := dq.Decrypt(ctx, []byte("9"), []byte("encrypted"))
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(pt, []byte("decrypted")) {
		t.Error("Decrypt did not return the expected plaintext from completedRequests")
	}
}

func TestExpiredCompletedRequest(t *testing.T) {
	dq := NewThresholdDecryptor(2, 1000, 100)

	dq.ResultReady([]byte("9"), []byte("decrypted"))

	NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.Lock()
		_, exists := dq.completedRequests[string([]byte("9"))]
		dq.mu.Unlock()
		return exists
	}, "15s", "10ms").Should(BeFalse(), "Completed request should be removed")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)
	defer cancel()

	_, err := dq.Decrypt(ctx, []byte("9"), []byte("encrypted"))
	if err == nil || err.Error() != "cancelled" {
		t.Error("Decrypt did not return the expected timeout error for an expired completedRequest")
	}
}

func TestCompletedRequestCleanup(t *testing.T) {
	dq := NewThresholdDecryptor(2, 1000, 100)

	dq.ResultReady([]byte("10"), []byte("decrypted"))

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)
	defer cancel1()

	_, _ = dq.Decrypt(ctx1, []byte("10"), []byte("encrypted")) // This will remove the decrypted result to completedRequests

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)
	defer cancel2()

	_, err := dq.Decrypt(ctx2, []byte("10"), []byte("encrypted"))
	if err == nil || err.Error() != "cancelled" {
		t.Error("Decrypt did not return the expected timeout error")
	}
}

func TestGetRequestsCountLimit(t *testing.T) {
	dq := NewThresholdDecryptor(4, 1000, 1000)

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel1()

	//nolint:errcheck
	go dq.Decrypt(ctx1, []byte("11"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel2()

	//nolint:errcheck
	go dq.Decrypt(ctx2, []byte("12"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("12"))

	ctx3, cancel3 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel3()

	//nolint:errcheck
	go dq.Decrypt(ctx3, []byte("13"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("13"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{
		{[]byte("11"), []byte("encrypted")},
		{[]byte("12"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("GetRequests did not return the expected requests when requestCountlimit is reached")
	}
}

func TestGetRequestsBytesLimit(t *testing.T) {
	dq := NewThresholdDecryptor(4, 10, 1000)

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel1()

	//nolint:errcheck
	go dq.Decrypt(ctx1, []byte("11"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel2()

	//nolint:errcheck
	go dq.Decrypt(ctx2, []byte("12"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("12"))

	ctx3, cancel3 := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel3()

	//nolint:errcheck
	go dq.Decrypt(ctx3, []byte("13"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("13"))

	requests := dq.GetRequests(4, 30)
	expected := []DecryptionRequest{
		{[]byte("11"), []byte("encrypted")},
		{[]byte("12"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("GetRequests did not return the expected requests when totalBytesLimit is exceeded")
	}
}

func TestGetRequestsWithShortPendingRequestQueue(t *testing.T) {
	dq := NewThresholdDecryptor(4, 1000, 1000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	defer cancel()

	//nolint:errcheck
	go dq.Decrypt(ctx, []byte("11"), []byte("encrypted"))

	waitForPendingRequestToBeAdded(t, dq, []byte("11"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{
		{[]byte("11"), []byte("encrypted")},
	}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("GetRequests did not return the expected requests when limit is set")
	}
}

func TestGetRequestsWithExpiredRequest(t *testing.T) {
	dq := NewThresholdDecryptor(4, 1000, 1000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)
	defer cancel()

	//nolint:errcheck
	dq.Decrypt(ctx, []byte("11"), []byte("encrypted"))

	requests := dq.GetRequests(2, 1000)
	expected := []DecryptionRequest{}
	if !reflect.DeepEqual(requests, expected) {
		t.Error("GetRequests did not return the expected requests when limit is set")
	}
}

func TestStart(t *testing.T) {
	dq := NewThresholdDecryptor(4, 1000, 1000)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := dq.Start(ctx)

	if err != nil {
		t.Error("Start returned unexpected error")
	}
}

func TestClose(t *testing.T) {
	dq := NewThresholdDecryptor(4, 1000, 1000)

	dq.ResultReady([]byte("14"), []byte("decrypted"))

	err := dq.Close()

	if err != nil {
		t.Error("Close returned unexpected error")
	}
}

func waitForPendingRequestToBeAdded(t *testing.T, dq *decryptionQueue, ciphertextId CiphertextId) {
	NewGomegaWithT(t).Eventually(func() bool {
		dq.mu.Lock()
		_, exists := dq.pendingRequests[string(ciphertextId)]
		dq.mu.Unlock()
		return exists
	}, "15s", "10ms").Should(BeTrue(), "Pending request should be added")
}
