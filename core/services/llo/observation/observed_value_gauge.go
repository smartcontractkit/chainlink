package observation

import (
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"

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
		"production OCR instance calls Observe. value_type is benchmark for decimal streams (the stream's value, " +
		"whatever it measures) and bid, benchmark or ask for quote streams. Values are float64, so digits beyond " +
		"about 15 significant figures are lost. A series is deleted when its stream has no fresh cached value, " +
		"leaves the plugin's stream set, or its config digest stops observing.",
}, []string{"stream_id", "value_type", "config_digest", "don_id", "csa_public_key"})

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

	// Constant per dataSource: the DON this job belongs to, and this node's CSA public key in hex (the same value as
	// the Beholder resource attribute csa_public_key). They let dashboards disambiguate stream_id across DONs and
	// join with RDD metadata.
	donIDLabel        string
	csaPublicKeyLabel string
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

func newObservedValueGauges(vec *prometheus.GaugeVec, donID uint32, csaPublicKey string) *observedValueGauges {
	return &observedValueGauges{
		vec:               vec,
		digests:           make(map[ocr2types.ConfigDigest]*digestObservedValues),
		donIDLabel:        strconv.FormatUint(uint64(donID), 10),
		csaPublicKeyLabel: csaPublicKey,
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
				entry.benchmark.Set(decimalToFloat64(v.Decimal()))
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
				entry.bid.Set(decimalToFloat64(v.Bid))
				entry.benchmark.Set(decimalToFloat64(v.Benchmark))
				entry.ask.Set(decimalToFloat64(v.Ask))
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
	entry.benchmark = g.vec.WithLabelValues(entry.streamIDLabel, observedValueTypeBenchmark, dg.digestLabel, g.donIDLabel, g.csaPublicKeyLabel)
	if isQuote {
		entry.bid = g.vec.WithLabelValues(entry.streamIDLabel, observedValueTypeBid, dg.digestLabel, g.donIDLabel, g.csaPublicKeyLabel)
		entry.ask = g.vec.WithLabelValues(entry.streamIDLabel, observedValueTypeAsk, dg.digestLabel, g.donIDLabel, g.csaPublicKeyLabel)
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
	g.vec.DeleteLabelValues(entry.streamIDLabel, observedValueTypeBenchmark, digestLabel, g.donIDLabel, g.csaPublicKeyLabel)
	if entry.bid != nil {
		g.vec.DeleteLabelValues(entry.streamIDLabel, observedValueTypeBid, digestLabel, g.donIDLabel, g.csaPublicKeyLabel)
		g.vec.DeleteLabelValues(entry.streamIDLabel, observedValueTypeAsk, digestLabel, g.donIDLabel, g.csaPublicKeyLabel)
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

// exactPow10 holds the powers of ten that float64 represents exactly.
var exactPow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11,
	1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19, 1e20, 1e21, 1e22,
}

// decimalToFloat64 converts d to float64 within 2 ulps of the nearest value. decimal.Float64 rounds exactly, but it
// builds a big.Rat and divides on every call (about 20x slower, with about 11 allocations); this runs on the OCR
// observation path for every changed stream, and the gauge does not need the last bit. It falls back to Float64 for
// exponents outside [-22, 22] and on platforms without 64-bit words.
func decimalToFloat64(d decimal.Decimal) float64 {
	exp := d.Exponent()
	if exp < -22 || exp > 22 || bits.UintSize != 64 {
		f, _ := d.Float64()
		return f
	}
	f := bigIntToFloat64(d.Coefficient()) // Coefficient copies: the only allocation
	if exp < 0 {
		return f / exactPow10[-exp]
	}
	return f * exactPow10[exp]
}

// bigIntToFloat64 converts c using its 64 most significant bits; the bits below them can only change the last bit.
// It requires 64-bit words.
func bigIntToFloat64(c *big.Int) float64 {
	n := c.BitLen()
	var f float64
	if n <= 64 {
		f = float64(c.Uint64()) // Uint64 is the low 64 bits of |c|, which is all of it here
	} else {
		w := c.Bits() // little-endian 64-bit words of |c|
		hi := (n - 1) / 64
		var top uint64
		if s := uint(n % 64); s == 0 {
			top = uint64(w[hi])
		} else {
			top = uint64(w[hi])<<(64-s) | uint64(w[hi-1])>>s
		}
		f = math.Ldexp(float64(top), n-64)
	}
	if c.Sign() < 0 {
		f = -f
	}
	return f
}
