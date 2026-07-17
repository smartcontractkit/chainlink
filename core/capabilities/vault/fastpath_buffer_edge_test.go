package vault

import (
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestFastPathBuffer_CompleteUnknownID(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Minute)

	// Completing an unknown request should not panic or affect unrelated channels.
	buf.Complete("nonexistent", &vaulttypes.Response{ID: "nonexistent", Payload: []byte("ok")})

	req := &vaultcommon.GetSecretsRequest{}
	respCh := buf.Submit("real", req)
	buf.Drain()
	buf.Complete("real", &vaulttypes.Response{ID: "real", Payload: []byte("ok")})

	select {
	case resp := <-respCh:
		require.Equal(t, []byte("ok"), resp.Payload)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for real response")
	}
}

func TestFastPathBuffer_CompleteAfterChannelConsumed(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	req := &vaultcommon.GetSecretsRequest{}
	respCh := buf.Submit("double-complete", req)

	buf.Drain()
	buf.Complete("double-complete", &vaulttypes.Response{ID: "double-complete", Payload: []byte("first")})

	select {
	case resp := <-respCh:
		require.Equal(t, []byte("first"), resp.Payload)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first response")
	}

	// Second completion should not panic or block because the channel is full/closed.
	buf.Complete("double-complete", &vaulttypes.Response{ID: "double-complete", Payload: []byte("second")})
}

func TestFastPathBuffer_ExpireInFlight(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Second)
	req := &vaultcommon.GetSecretsRequest{}
	respCh := buf.Submit("inflight-expire", req)

	// Move from pending to in-flight.
	drained := buf.Drain()
	require.Len(t, drained, 1)

	clock.Advance(2 * time.Second)
	buf.ExpireOlderThan(clock.Now())

	select {
	case resp := <-respCh:
		require.NotEmpty(t, resp.Error)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for in-flight expiry response")
	}
}

func TestFastPathBuffer_ConcurrentStress(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Minute)

	const numGoroutines = 20
	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // submitters, drainers, completers

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				buf.Submit(makeID(id, j), &vaultcommon.GetSecretsRequest{})
			}
		}(i)
	}

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for range iterations {
				_ = buf.Drain()
			}
		}(i)
	}

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				buf.Complete(makeID(id, j), &vaulttypes.Response{ID: makeID(id, j), Payload: []byte("ok")})
			}
		}(i)
	}

	wg.Wait()

	// Drain and complete any remaining pending requests to ensure no goroutine is leaked.
	remaining := buf.Drain()
	for _, fr := range remaining {
		buf.Complete(fr.ID, &vaulttypes.Response{ID: fr.ID, Payload: []byte("ok")})
	}
	buf.ExpireOlderThan(clock.Now().Add(2 * time.Minute))
}

func makeID(goroutine, iteration int) string {
	return "stress-" + string(rune('a'+goroutine%26)) + "-" + string(rune('0'+iteration%10))
}
