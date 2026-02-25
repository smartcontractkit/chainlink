package agent

import (
	"slices"
	"strings"
	"time"
)

func (s *Server) lookupCachedStart(componentKey, payloadHash string) (*cachedStart, bool) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	start, ok := s.cache[componentKey]
	if !ok || start.PayloadHash != payloadHash {
		return nil, false
	}
	return &start, true
}

func (s *Server) cacheSuccessfulStart(componentKey, payloadHash string, output map[string]any) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache[componentKey] = cachedStart{
		PayloadHash: payloadHash,
		Output:      output,
	}
}

func (s *Server) deleteCachedStart(componentKey string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	delete(s.cache, componentKey)
}

func (s *Server) storeRuntime(componentKey string, state runtimeState) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.runtime[componentKey] = state
}

func (s *Server) takeRuntime(componentKey string) (runtimeState, bool) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	state, ok := s.runtime[componentKey]
	if ok {
		delete(s.runtime, componentKey)
	}
	return state, ok
}

func (s *Server) beginInFlight(id, scope string) {
	s.opsMu.Lock()
	defer s.opsMu.Unlock()
	s.inFlight[id] = inFlightOperation{
		ID:        id,
		Scope:     scope,
		StartedAt: time.Now(),
	}
}

func (s *Server) endInFlight(id string) {
	s.opsMu.Lock()
	defer s.opsMu.Unlock()
	delete(s.inFlight, id)
}

func (s *Server) inFlightSnapshot() ([]InFlightOperation, bool) {
	s.opsMu.Lock()
	defer s.opsMu.Unlock()

	out := make([]InFlightOperation, 0, len(s.inFlight))
	lifecycleBusy := false
	for _, op := range s.inFlight {
		if op.Scope == inFlightOperationScopeLifecycle {
			lifecycleBusy = true
		}
		out = append(out, InFlightOperation{
			ID:         op.ID,
			Scope:      op.Scope,
			StartedAt:  op.StartedAt.Format(time.RFC3339Nano),
			DurationMs: int64(time.Since(op.StartedAt) / time.Millisecond),
		})
	}
	slices.SortFunc(out, func(a, b InFlightOperation) int {
		return strings.Compare(a.ID, b.ID)
	})
	return out, lifecycleBusy
}

func (s *Server) cacheKeys() []string {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	keys := make([]string, 0, len(s.cache))
	for k := range s.cache {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func (s *Server) runtimeKeys() []string {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	keys := make([]string, 0, len(s.runtime))
	for k := range s.runtime {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func (s *Server) cacheSize() int {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	return len(s.cache)
}

func (s *Server) runtimeSize() int {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	return len(s.runtime)
}

func (s *Server) relayInfos() []RelayInfo {
	s.relayMu.Lock()
	defer s.relayMu.Unlock()

	out := make([]RelayInfo, 0, len(s.relays))
	for _, relay := range s.relays {
		if relay == nil {
			continue
		}
		out = append(out, RelayInfo{
			ID:            relay.ID,
			Name:          relay.Name,
			RequestedPort: relay.RequestedPort,
			BoundPort:     listenerPort(relay.Listener),
		})
	}
	slices.SortFunc(out, func(a, b RelayInfo) int {
		return strings.Compare(a.ID, b.ID)
	})
	return out
}

func (s *Server) relayCount() int {
	s.relayMu.Lock()
	defer s.relayMu.Unlock()
	return len(s.relays)
}
