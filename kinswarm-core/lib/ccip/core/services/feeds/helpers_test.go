package feeds

// SetConnectionsManager allows us to manually set the connections manager.
// Only used for testing.
func (s *service) SetConnectionsManager(cm ConnectionsManager) {
	s.connMgr = cm
}
