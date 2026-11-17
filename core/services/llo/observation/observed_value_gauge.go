package observation

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
)

// promStreamObservedValue uses snake_case labels (unlike the other metrics in this package) so config_digest joins
// directly with libocr's config_digest label, which carries the same ConfigDigest.String() value.
var promStreamObservedValue = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "llo",
	Subsystem: "datasource",
	Name:      "stream_observed_value",
	Help: "Latest value this node observed for an LLO stream, read from the observation cache when the " +
		"production OCR instance calls Observe. value_type is benchmark for decimal streams and bid, " +
		"benchmark or ask for quote streams. A series is deleted when its stream has no fresh cached value, " +
		"leaves the plugin's stream set, or its config digest stops observing.",
}, []string{"stream_id", "value_type", "config_digest"})

const (
	observedValueTypeBenchmark = "benchmark"
	observedValueTypeBid       = "bid"
	observedValueTypeAsk       = "ask"
	// observedValueDigestTTL is how long a config digest's series survive after its last production Observe.
	// A TTL (instead of deleting on every digest change) avoids churning every series while Blue and Green
	// instances are both production during a handover.
	observedValueDigestTTL = 10 * time.Second
)

// observedValueGauges owns the llo_datasource_stream_observed_value series that one dataSource sets,
// so it can delete them when a stream goes stale, a config digest stops observing, or the data source closes.
type observedValueGauges struct {
	vec     *prometheus.GaugeVec
	mu      sync.Mutex
	digests map[ocr2types.ConfigDigest]*digestObservedValues
}

type digestObservedValues struct {
	digestLabel  string    // ConfigDigest.String(), same format as libocr's config_digest label
	lastRecorded time.Time // time of the last record() for this digest
	round        uint64    // incremented on each record() for this digest
	streams      map[llotypes.StreamID]*streamObservedValue
}

type streamObservedValue struct {
	streamIDLabel string
	last          lloprotocol.StreamValue // value already written to the gauges; pointer compare skips Float64
	round         uint64                  // digestObservedValues.round when last seen with a value
	benchmark     prometheus.Gauge
	bid, ask      prometheus.Gauge // nil for decimal streams
}

func newObservedValueGauges(vec *prometheus.GaugeVec) *observedValueGauges {
	return &observedValueGauges{
		vec:     vec,
		digests: make(map[ocr2types.ConfigDigest]*digestObservedValues),
	}
}

// record sets one series per value component for every stream in streamValues that has a decimal or quote value,
// and deletes the series of streams that no longer have one. It also deletes the series of other config digests
// that have not recorded for observedValueDigestTTL.
func (g *observedValueGauges) record(digest ocr2types.ConfigDigest, streamValues lloprotocol.StreamValues, now time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()

	dg, ok := g.digests[digest]
	if !ok {
		dg = &digestObservedValues{
			digestLabel: digest.String(),
			streams:     make(map[llotypes.StreamID]*streamObservedValue, len(streamValues)),
		}
		g.digests[digest] = dg
	}
	dg.lastRecorded = now
	dg.round++

	seen := 0
	for id, orig := range streamValues {
		switch v := unwrapStreamValue(orig).(type) {
		case *lloprotocol.Decimal:
			if v == nil {
				continue
			}
			entry := g.getOrCreateStream(dg, id, false)
			if entry.last != orig {
				f, _ := v.Decimal().Float64()
				entry.benchmark.Set(f)
				entry.last = orig
			}
			entry.round = dg.round
			seen++
		case *lloprotocol.Quote:
			if v == nil {
				continue
			}
			entry := g.getOrCreateStream(dg, id, true)
			if entry.last != orig {
				bid, _ := v.Bid.Float64()
				benchmark, _ := v.Benchmark.Float64()
				ask, _ := v.Ask.Float64()
				entry.bid.Set(bid)
				entry.benchmark.Set(benchmark)
				entry.ask.Set(ask)
				entry.last = orig
			}
			entry.round = dg.round
			seen++
		}
	}

	if seen < len(dg.streams) {
		for id, entry := range dg.streams {
			if entry.round != dg.round {
				g.deleteStream(dg.digestLabel, entry)
				delete(dg.streams, id)
			}
		}
	}

	// dg.lastRecorded is now, so this only expires other digests (a retired Blue/Green instance).
	g.expireLocked(now)
}

// expire deletes the series of every config digest that has not recorded for observedValueDigestTTL. It also
// covers the only production digest when Observe stops (for example a DON stall), which record cannot detect.
func (g *observedValueGauges) expire(now time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.expireLocked(now)
}

func (g *observedValueGauges) expireLocked(now time.Time) {
	for d, dg := range g.digests {
		if now.Sub(dg.lastRecorded) > observedValueDigestTTL {
			g.deleteDigest(dg)
			delete(g.digests, d)
		}
	}
}

// deleteAll deletes every series this recorder set.
func (g *observedValueGauges) deleteAll() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, dg := range g.digests {
		g.deleteDigest(dg)
	}
	g.digests = make(map[ocr2types.ConfigDigest]*digestObservedValues)
}

// getOrCreateStream returns the entry for id, replacing it when its kind (decimal or quote) differs from isQuote.
func (g *observedValueGauges) getOrCreateStream(dg *digestObservedValues, id llotypes.StreamID, isQuote bool) *streamObservedValue {
	entry, ok := dg.streams[id]
	if ok && (entry.bid != nil) == isQuote {
		return entry
	}
	if ok {
		g.deleteStream(dg.digestLabel, entry)
	}
	entry = &streamObservedValue{streamIDLabel: strconv.FormatUint(uint64(id), 10)}
	entry.benchmark = g.vec.WithLabelValues(entry.streamIDLabel, observedValueTypeBenchmark, dg.digestLabel)
	if isQuote {
		entry.bid = g.vec.WithLabelValues(entry.streamIDLabel, observedValueTypeBid, dg.digestLabel)
		entry.ask = g.vec.WithLabelValues(entry.streamIDLabel, observedValueTypeAsk, dg.digestLabel)
	}
	dg.streams[id] = entry
	return entry
}

func (g *observedValueGauges) deleteDigest(dg *digestObservedValues) {
	for _, entry := range dg.streams {
		g.deleteStream(dg.digestLabel, entry)
	}
}

func (g *observedValueGauges) deleteStream(digestLabel string, entry *streamObservedValue) {
	g.vec.DeleteLabelValues(entry.streamIDLabel, observedValueTypeBenchmark, digestLabel)
	if entry.bid != nil {
		g.vec.DeleteLabelValues(entry.streamIDLabel, observedValueTypeBid, digestLabel)
		g.vec.DeleteLabelValues(entry.streamIDLabel, observedValueTypeAsk, digestLabel)
	}
}

// unwrapStreamValue returns the value inside any TimestampedStreamValue wrappers, or nil for a nil wrapper.
func unwrapStreamValue(v lloprotocol.StreamValue) lloprotocol.StreamValue {
	for {
		ts, ok := v.(*lloprotocol.TimestampedStreamValue)
		if !ok {
			return v
		}
		if ts == nil {
			return nil
		}
		v = ts.StreamValue
	}
}
