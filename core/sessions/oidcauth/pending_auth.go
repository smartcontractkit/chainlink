package oidcauth

import (
	"sync"
	"time"
)

// pendingAuth is the in-memory copy used by unit tests that have no database.
// Production sign-in and exchange write oidc_pending_auth. The cookie holds
// only the anti-CSRF state. The PKCE verifier and OIDC nonce stay server-side,
// keyed by that state, until a single successful take at exchange time.
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
