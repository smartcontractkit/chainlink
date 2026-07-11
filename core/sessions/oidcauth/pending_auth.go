package oidcauth

import (
	"sync"
	"time"
)

// pendingAuth holds server-side secrets for one browser authorization-code
// login. The PKCE verifier and OIDC nonce must not round-trip through the
// signed-but-not-encrypted session cookie (see OCD-004 / OCD-005). Only the
// anti-CSRF state value is stored client-side; the secrets stay here, keyed by
// that state, until a single successful take() at exchange time.
//
// Like deviceFlowStore this is process-local: multi-replica deployments need
// session affinity for /oidc-login → /oidc-exchange, or a shared short-TTL store.
type pendingAuth struct {
	verifier  string
	nonce     string
	expiresAt time.Time
}

type pendingAuthStore struct {
	mu      sync.Mutex
	entries map[string]pendingAuth
}

func newPendingAuthStore() *pendingAuthStore {
	return &pendingAuthStore{entries: make(map[string]pendingAuth)}
}

const pendingAuthTTL = 15 * time.Minute

func (s *pendingAuthStore) sweepLocked(now time.Time) {
	for k, e := range s.entries {
		if now.After(e.expiresAt) {
			delete(s.entries, k)
		}
	}
}

// put records verifier+nonce for state. Overwrites any prior entry for the same
// state (should not happen with 256-bit state values).
func (s *pendingAuthStore) put(state, verifier, nonce string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked(time.Now())
	s.entries[state] = pendingAuth{
		verifier:  verifier,
		nonce:     nonce,
		expiresAt: time.Now().Add(pendingAuthTTL),
	}
}

// take removes and returns the entry for state. Single-use: a second take fails.
func (s *pendingAuthStore) take(state string) (verifier, nonce string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked(time.Now())
	e, ok := s.entries[state]
	if !ok {
		return "", "", false
	}
	delete(s.entries, state)
	if time.Now().After(e.expiresAt) {
		return "", "", false
	}
	return e.verifier, e.nonce, true
}
