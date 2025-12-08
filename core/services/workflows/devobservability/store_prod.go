//go:build !dev

package devobservability

type noopStore struct{}

func init() {
	store := &noopStore{}
	globalStore = store
}

func (s *noopStore) GetWorkflows() []string { return nil }

func (s *noopStore) GetOrphanEvents(limit int) []EventEntry { return nil }

func (s *noopStore) GetWorkflowEvents(workflowID string, limit int, minSequence int64) []EventEntry {
	return nil
}

func (s *noopStore) Clear() {}

func (s *noopStore) Stats() map[string]interface{} { return nil }
