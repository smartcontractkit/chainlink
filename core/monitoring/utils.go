package monitoring

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// Deprecated: use [beholder.OtelAttributes.AsStringAttributes]
func KvMapToOtelAttributes(kvmap map[string]string) []attribute.KeyValue {
	return beholder.OtelAttributes(kvmap).AsStringAttributes()
}
