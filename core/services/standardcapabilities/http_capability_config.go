package standardcapabilities

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// injectHTTPCapabilityConfig merges node TOML configuration for the
// http-trigger and http-action capabilities into the job-spec config JSON.
//
// Node TOML is the authoritative source: any value set in TOML overrides the
// corresponding job-spec value; values unset in TOML fall back to the job-spec
// value (or, when absent there too, to the capability binary's built-in
// default, since zero values are omitted from the merged JSON).
//
// This keeps capability job specs minimal while preserving backward
// compatibility with existing job specs that still carry these fields.
func injectHTTPCapabilityConfig(lggr logger.Logger, command, configJSON string, inject func(map[string]any)) string {
	merged := make(map[string]any)
	if configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), &merged); err != nil {
			// Malformed job-spec config: log and let the capability binary
			// surface the parse error, rather than failing here.
			lggr.Warnw("failed to parse capability config JSON; TOML values will still be injected",
				"command", command, "error", err)
			merged = make(map[string]any)
		}
	}

	inject(merged)

	b, err := json.Marshal(merged)
	if err != nil {
		// merged only contains JSON-compatible values, so this is unreachable
		// in practice; fall back to the original config rather than failing.
		lggr.Errorw("failed to marshal merged capability config; using original", "command", command, "error", err)
		return configJSON
	}
	return string(b)
}

// setIfNotZero sets key to value in m when value is non-zero.
// Zero values are omitted so the capability binary applies its own defaults.
func setIfNotZero(m map[string]any, key string, value any) {
	switch v := value.(type) {
	case uint16:
		if v != 0 {
			m[key] = v
		}
	case uint32:
		if v != 0 {
			m[key] = v
		}
	case int:
		if v != 0 {
			m[key] = v
		}
	case float64:
		if v != 0 {
			m[key] = v
		}
	case string:
		if v != "" {
			m[key] = v
		}
	case []string:
		if len(v) > 0 {
			m[key] = v
		}
	case []int:
		if len(v) > 0 {
			m[key] = v
		}
	case map[string]any:
		if len(v) > 0 {
			m[key] = v
		}
	default:
		panic(fmt.Sprintf("setIfNotZero: unsupported type %T", value))
	}
}

// nested returns the nested map at key, creating it if needed.
func nested(m map[string]any, key string) map[string]any {
	existing, ok := m[key].(map[string]any)
	if !ok {
		existing = make(map[string]any)
		m[key] = existing
	}
	return existing
}

// attachNested sets m[key] to sub only when sub is non-empty, so unset
// config sections are omitted entirely rather than emitted as empty objects.
func attachNested(m map[string]any, key string, sub map[string]any) {
	if len(sub) > 0 {
		m[key] = sub
	}
}
