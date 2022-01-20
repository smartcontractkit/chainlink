package monitoring

import (
	"context"
)

type Exporter interface {
	// Export is executed on each update on a monitored feed
	Export(ctx context.Context, data interface{})
	// Cleanup is executed one a monitor for a specific feed is terminated.
	Cleanup()
}
