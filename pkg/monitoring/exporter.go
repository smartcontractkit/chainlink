package monitoring

import (
	"context"
)

// Exporter methods can be executed out of order and should be thread safe.
type Exporter interface {
	// Export is executed on each update on a monitored feed
	Export(ctx context.Context, data interface{})
	// Cleanup is executed once a monitor for a specific feed is terminated.
	Cleanup()
}
