package agent

import (
	"slices"
	"strings"
)

func (s *Server) appendComponentLogs(componentKey string, lines []string) {
	if strings.TrimSpace(componentKey) == "" || len(lines) == 0 {
		return
	}
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	if len(filtered) == 0 {
		return
	}

	s.logsMu.Lock()
	defer s.logsMu.Unlock()
	existing := append(s.componentLogs[componentKey], filtered...)
	if len(existing) > componentLogsRingSize {
		existing = existing[len(existing)-componentLogsRingSize:]
	}
	s.componentLogs[componentKey] = existing
}

func (s *Server) getComponentLogs(componentKey string, limit int) ([]string, int) {
	s.logsMu.Lock()
	defer s.logsMu.Unlock()
	lines := s.componentLogs[componentKey]
	total := len(lines)
	if total == 0 {
		return []string{}, 0
	}
	if limit <= 0 || limit > total {
		limit = total
	}
	out := append([]string{}, lines[total-limit:]...)
	return out, total
}

func (s *Server) componentLogKeys() []string {
	s.logsMu.Lock()
	defer s.logsMu.Unlock()
	keys := make([]string, 0, len(s.componentLogs))
	for k := range s.componentLogs {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
