package v2

import "context"

// TestPut is a test-only escape hatch onto put, exported here purely so engine_test.go's
// admission-denial coverage can reach it from the external v2_test package.
func (e *Engine) TestPut(ctx context.Context, event RoutedTriggerEvent) error {
	return e.put(ctx, event)
}
