package monitoring

import "go.opentelemetry.io/otel/attribute"

func KvMapToOtelAttributes(kvmap map[string]string) []attribute.KeyValue {
	otelKVs := make([]attribute.KeyValue, 0, len(kvmap))
	for k, v := range kvmap {
		otelKVs = append(otelKVs, attribute.String(k, v))
	}
	return otelKVs
}
