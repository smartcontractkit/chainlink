package por

import "context"

type ReportMarshaler interface {
	Serialize(ctx context.Context, chain ChainSelector, report PorReport) ([]byte, error)

	MaxReportSize(ctx context.Context) int // The maximum size of the serialized report in bytes.
}
