package oidcauth

import (
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
)

func TestDeviceFlowStore_AddGetRemove(t *testing.T) {
	t.Parallel()
	s := newDeviceFlowStore()

	state := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}
	require.NoError(t, s.add("handle-1", state))

	got, ok := s.get("handle-1")
	require.True(t, ok)
	assert.Same(t, state, got)

	s.remove("handle-1")
	_, ok = s.get("handle-1")
	assert.False(t, ok, "handle must be single-use: gone after remove")
}

// TestDeviceFlowStore_ConsumeIfDoneSingleWinner asserts that when many polls
// race on one completed handle, exactly one observes the terminal result and
// consumes it. This is the double-spend guard: two polls must never both be
// handed the session cookie.
func TestDeviceFlowStore_ConsumeIfDoneSingleWinner(t *testing.T) {
	t.Parallel()
	s := newDeviceFlowStore()
	state := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}
	state.done = true
	state.sessionID = "sess-winner"
	require.NoError(t, s.add("race-handle", state))

	const pollers = 50
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		winners  int
		pendings int
	)
	wg.Add(pollers)
	for range pollers {
		go func() {
			defer wg.Done()
			sessionID, _, _, terminal, known, _ := s.consumeIfDone("race-handle")
			mu.Lock()
			defer mu.Unlock()
			switch {
			case terminal:
				winners++
				assert.Equal(t, "sess-winner", sessionID)
			case known:
				pendings++ // saw it before it was consumed, but not terminal: impossible here
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, winners, "exactly one poll may consume a completed handle")
	assert.Equal(t, 0, pendings, "a done flow must never report pending")
	_, ok := s.get("race-handle")
	assert.False(t, ok, "handle must be removed after it is consumed")
}

func TestDeviceFlowStore_ConcurrencyCap(t *testing.T) {
	t.Parallel()
	s := newDeviceFlowStore()
	for i := range maxConcurrentDeviceFlows {
		require.NoError(t, s.add(handleN(i), &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}))
	}
	// One past the cap must be refused, not silently dropped.
	err := s.add("overflow", &deviceFlowState{expiresAt: time.Now().Add(time.Minute)})
	require.ErrorIs(t, err, errTooManyDeviceFlows)
}

func TestDeviceFlowStore_ExpiredSweptOnAdd(t *testing.T) {
	t.Parallel()
	s := newDeviceFlowStore()
	// Fill with already-expired flows.
	for i := range maxConcurrentDeviceFlows {
		require.NoError(t, s.add(handleN(i), &deviceFlowState{expiresAt: time.Now().Add(-time.Second)}))
	}
	// A fresh add sweeps the expired entries, making room rather than erroring.
	require.NoError(t, s.add("fresh", &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}))
	_, ok := s.get("fresh")
	assert.True(t, ok)
}

func TestGenerateDeviceHandle_UniqueAndURLSafe(t *testing.T) {
	t.Parallel()
	oi := &oidcAuthenticator{}
	seen := make(map[string]struct{})
	for range 1000 {
		h, err := oi.generateDeviceHandle()
		require.NoError(t, err)
		require.NotEmpty(t, h)
		// RawURLEncoding of 32 bytes -> 43 chars, no padding, URL-safe alphabet.
		assert.Len(t, h, 43)
		assert.NotContains(t, h, "+")
		assert.NotContains(t, h, "/")
		assert.NotContains(t, h, "=")
		_, dup := seen[h]
		require.False(t, dup, "handle collision is a security failure")
		seen[h] = struct{}{}
	}
}

func TestFinishDeviceFlow_RecordsTerminalState(t *testing.T) {
	t.Parallel()
	oi := &oidcAuthenticator{lggr: logger.TestLogger(t)}

	okState := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}
	oi.finishDeviceFlow(okState, "sess-1", "user@example.com", clsessions.UserRoleEdit, nil)
	okState.mu.Lock()
	assert.True(t, okState.done)
	assert.Equal(t, "sess-1", okState.sessionID)
	assert.Equal(t, "user@example.com", okState.email)
	require.NoError(t, okState.err)
	okState.mu.Unlock()

	errState := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}
	wantErr := errors.New("denied")
	oi.finishDeviceFlow(errState, "", "", "", wantErr)
	errState.mu.Lock()
	assert.True(t, errState.done)
	assert.Empty(t, errState.sessionID)
	require.ErrorIs(t, errState.err, wantErr)
	errState.mu.Unlock()
}

// TestDeviceFlowState_ConcurrentAccess exercises the writer/reader lock that
// separates the background poll goroutine from the CLI poll handler.
func TestDeviceFlowState_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	oi := &oidcAuthenticator{lggr: logger.TestLogger(t)}
	state := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		oi.finishDeviceFlow(state, "s", "e", clsessions.UserRoleView, nil)
	}()
	go func() {
		defer wg.Done()
		for range 100 {
			state.mu.Lock()
			_ = state.done
			state.mu.Unlock()
		}
	}()
	wg.Wait()
}

func handleN(i int) string {
	return "h-" + strconv.Itoa(i)
}

// mutableConfig embeds the default *TestConfig and lets a test override the one
// field under test. TestConfig's methods use a pointer receiver, so embed the
// pointer to satisfy config.OIDC.
type mutableConfig struct {
	*TestConfig
	clientSecret *string
	clientID     *string
}

func (m mutableConfig) ClientSecret() string {
	if m.clientSecret != nil {
		return *m.clientSecret
	}
	return m.TestConfig.ClientSecret()
}

func (m mutableConfig) ClientID() string {
	if m.clientID != nil {
		return *m.clientID
	}
	return m.TestConfig.ClientID()
}

func TestValidateOIDCConfig_SecretOptional(t *testing.T) {
	t.Parallel()
	empty := ""
	base := &TestConfig{}

	// Public client: empty secret must be accepted (this is the change that
	// lets a Native app + PKCE work for the CLI device flow).
	require.NoError(t, validateOIDCConfig(mutableConfig{TestConfig: base, clientSecret: &empty}),
		"empty ClientSecret must be allowed for public/PKCE clients")

	// Confidential client: a secret is still fine (backwards compatibility).
	require.NoError(t, validateOIDCConfig(mutableConfig{TestConfig: base}),
		"a populated ClientSecret must still be accepted")

	// ClientID, by contrast, remains required.
	require.Error(t, validateOIDCConfig(mutableConfig{TestConfig: base, clientID: &empty}),
		"empty ClientID must still be rejected")
}
