package gateway

import (
	"sync"
	"time"
)

const pingTokenLen = 8

// token aware ping tracker for o11y
type pingProbeTracker struct {
	mu      sync.Mutex
	token   string
	start   time.Time
	pending bool
}

// register stores the token and start time of the probe about to be written.
// It must be called before the ping is written, because the matching pong can
// arrive before Write returns. Registering a new probe supersedes any previous
// unclaimed one.
func (t *pingProbeTracker) register(token string, start time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.token = token
	t.start = start
	t.pending = true
}

// consume returns the probe start time iff token matches the latest registered
// probe, clearing it so a duplicate pong cannot record a second observation.
// Superseded or unknown tokens yield no observation.
func (t *pingProbeTracker) consume(token string) (time.Time, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.pending || t.token != token {
		return time.Time{}, false
	}
	t.pending = false
	t.token = ""
	return t.start, true
}
