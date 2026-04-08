package pipeline

// TestingSetBridgeRequiredJSONPaths sets required JSON paths on a bridge task
// for tests in external test packages (e.g. pipeline_test).
func TestingSetBridgeRequiredJSONPaths(t *BridgeTask, paths [][]string) {
	t.requiredJSONPaths = paths
}
