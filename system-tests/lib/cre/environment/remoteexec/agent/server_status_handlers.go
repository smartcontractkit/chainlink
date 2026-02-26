package agent

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	runtimeKeys := s.runtimeKeys()
	cacheKeys := s.cacheKeys()
	relayInfos := s.relayInfos()
	componentLogKeys := s.componentLogKeys()
	inFlight, _ := s.inFlightSnapshot()
	chipSinkStatus := s.currentChipSinkStatus()

	s.respondJSONAny(w, http.StatusOK, AgentStatusResponse{
		AgentVersion:      agentVersion,
		ProtocolVersion:   protocolVersion,
		SupportedSchemas:  []string{SchemaVersionV1},
		Capabilities:      []string{capabilityStartComponent, capabilityDeployArtifacts, capabilityRelay, capabilityListCTFResources, capabilityLocks, capabilityComponentLogs, "chipSinkLifecycle"},
		UptimeSeconds:     int64(time.Since(s.startedAt).Seconds()),
		RuntimeComponents: runtimeKeys,
		CachedComponents:  cacheKeys,
		Relays:            relayInfos,
		ComponentLogKeys:  componentLogKeys,
		InFlight:          inFlight,
		ChipSink:          chipSinkStatus,
	})
}

func (s *Server) locks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	inFlight, lifecycleBusy := s.inFlightSnapshot()
	s.respondJSONAny(w, http.StatusOK, AgentLocksResponse{
		LifecycleBusy:    lifecycleBusy,
		CacheEntries:     s.cacheSize(),
		RuntimeEntries:   s.runtimeSize(),
		RelayCount:       s.relayCount(),
		ComponentLogKeys: len(s.componentLogKeys()),
		InFlight:         inFlight,
	})
}

func (s *Server) componentLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	componentKey := strings.TrimSpace(r.URL.Query().Get("componentKey"))
	if componentKey == "" {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "componentKey query parameter is required", nil)
		return
	}
	limit := defaultComponentLogsLimit
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, "limit query parameter must be a positive integer", nil)
			return
		}
		if parsed > maxComponentLogsLimit {
			parsed = maxComponentLogsLimit
		}
		limit = parsed
	}

	lines, total := s.getComponentLogs(componentKey, limit)
	s.respondJSONAny(w, http.StatusOK, ComponentLogsResponse{
		ComponentKey: componentKey,
		TotalLines:   total,
		Lines:        lines,
	})
}
