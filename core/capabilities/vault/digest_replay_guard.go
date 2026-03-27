package vault

import (
	"errors"
	"sync"
	"time"
)

var ErrDigestAlreadySeen = errors.New("request already authorized previously")

// DigestReplayGuard prevents replay of already-processed requests by tracking
// request digests with expiry timestamps. It is safe for concurrent use.
//
// Used by both the on-chain allowlist flow and the JWT auth flow to ensure
// that a given request digest is only accepted once.
type DigestReplayGuard struct {
	mu      sync.Mutex
	seen    map[string]int64 // digest → unix expiry timestamp
	nowFunc func() time.Time // injectable for testing
}

func NewDigestReplayGuard() *DigestReplayGuard {
	return &DigestReplayGuard{
		seen:    make(map[string]int64),
		nowFunc: time.Now,
	}
}

// CheckAndRecord returns ErrDigestAlreadySeen if the digest was previously
// recorded and has not yet expired. Otherwise it records the digest with
// the given expiry timestamp (unix seconds, UTC).
//
// Expired entries are cleaned up on every call.
func (g *DigestReplayGuard) CheckAndRecord(digest string, expiresAtUnix int64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.clearExpiredLocked()

	if _, exists := g.seen[digest]; exists {
		return ErrDigestAlreadySeen
	}

	g.seen[digest] = expiresAtUnix
	return nil
}

// ClearExpired removes all entries whose expiry timestamp is in the past.
// Call this to eagerly reclaim memory even when CheckAndRecord is not invoked.
func (g *DigestReplayGuard) ClearExpired() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.clearExpiredLocked()
}

func (g *DigestReplayGuard) clearExpiredLocked() {
	now := g.nowFunc().UTC().Unix()
	for digest, expiry := range g.seen {
		if now > expiry {
			delete(g.seen, digest)
		}
	}
}
