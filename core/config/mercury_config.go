package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-data-streams/llo/transmitter/de"
)

// Mercury is the node-side view of Mercury config. It embeds the data-streams
// transmitter's Mercury interface and adds node-only sections (e.g. DataSource).
type Mercury interface {
	de.Mercury
	DataSource() MercuryDataSource
}

// MercuryDataSource exposes node-side tuning for the LLO observation data source.
type MercuryDataSource interface {
	// ObservationTimingBase returns the base duration T used to size the LLO observation
	// loop timing: cache entry TTL, stale-refresh threshold, loop pacing, and background
	// pipeline timeout. When zero, the plugin Observe deadline remainder is used instead
	// (legacy behavior tied to the on-chain MaxDurationObservation).
	ObservationTimingBase() time.Duration
}
