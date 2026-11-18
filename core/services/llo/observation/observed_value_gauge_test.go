package observation

import (
	"math/big"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
)

func newTestObservedValueVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_observed_value"}, []string{"stream_id", "value_type", "config_digest"})
}

// observedSeries returns every series in vec keyed by {stream_id, value_type, config_digest}. It reports errors with
// assert, not require, so it is safe in goroutines and in EventuallyWithT conditions.
func observedSeries(t assert.TestingT, vec *prometheus.GaugeVec) map[[3]string]float64 {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}
	ch := make(chan prometheus.Metric)
	go func() {
		vec.Collect(ch)
		close(ch)
	}()
	out := make(map[[3]string]float64)
	for m := range ch {
		var pb io_prometheus_client.Metric
		if !assert.NoError(t, m.Write(&pb)) {
			continue
		}
		labels := make(map[string]string, len(pb.GetLabel()))
		for _, lp := range pb.GetLabel() {
			labels[lp.GetName()] = lp.GetValue()
		}
		out[[3]string{labels["stream_id"], labels["value_type"], labels["config_digest"]}] = pb.GetGauge().GetValue()
	}
	return out
}

func dec(f float64) *lloprotocol.Decimal {
	return lloprotocol.ToDecimal(decimal.NewFromFloat(f))
}

func quote(bid, benchmark, ask float64) *lloprotocol.Quote {
	return &lloprotocol.Quote{Bid: decimal.NewFromFloat(bid), Benchmark: decimal.NewFromFloat(benchmark), Ask: decimal.NewFromFloat(ask)}
}

func Test_observedValueGauges(t *testing.T) {
	t.Parallel()
	dA := ocr2types.ConfigDigest{1}
	dB := ocr2types.ConfigDigest{2}
	t0 := time.Unix(1_700_000_000, 0)

	t.Run("decimal records benchmark", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(63421.07)}, t0)
		assert.Equal(t, map[[3]string]float64{
			{"1", "benchmark", dA.String()}: 63421.07,
		}, observedSeries(t, vec))
	})

	t.Run("legacy quote records bid, benchmark and ask", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{2: quote(99, 100, 101)}, t0)
		series := observedSeries(t, vec)
		assert.Equal(t, map[[3]string]float64{
			{"2", "bid", dA.String()}:       99,
			{"2", "benchmark", dA.String()}: 100,
			{"2", "ask", dA.String()}:       101,
		}, series)
		for k := range series {
			assert.Contains(t, []string{"benchmark", "bid", "ask"}, k[1])
		}
	})

	t.Run("timestamped values are unwrapped", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{
			1: &lloprotocol.TimestampedStreamValue{ObservedAtNanoseconds: 1, StreamValue: dec(63421.07)},
			2: &lloprotocol.TimestampedStreamValue{ObservedAtNanoseconds: 1, StreamValue: quote(99, 100, 101)},
		}, t0)
		assert.Equal(t, map[[3]string]float64{
			{"1", "benchmark", dA.String()}: 63421.07,
			{"2", "bid", dA.String()}:       99,
			{"2", "benchmark", dA.String()}: 100,
			{"2", "ask", dA.String()}:       101,
		}, observedSeries(t, vec))
	})

	t.Run("nil value deletes the stream's series", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(1), 2: dec(2)}, t0)
		g.record(dA, lloprotocol.StreamValues{1: nil, 2: dec(2)}, t0.Add(time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"2", "benchmark", dA.String()}: 2,
		}, observedSeries(t, vec))
	})

	t.Run("stream removed from the set deletes its series", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(1), 2: dec(2)}, t0)
		g.record(dA, lloprotocol.StreamValues{2: dec(2)}, t0.Add(time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"2", "benchmark", dA.String()}: 2,
		}, observedSeries(t, vec))
	})

	t.Run("new value updates the gauge", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(10)}, t0)
		g.record(dA, lloprotocol.StreamValues{1: dec(11)}, t0.Add(time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"1", "benchmark", dA.String()}: 11,
		}, observedSeries(t, vec))
	})

	t.Run("kind change replaces the series", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{3: dec(5)}, t0)
		g.record(dA, lloprotocol.StreamValues{3: quote(4, 5, 6)}, t0.Add(time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"3", "bid", dA.String()}:       4,
			{"3", "benchmark", dA.String()}: 5,
			{"3", "ask", dA.String()}:       6,
		}, observedSeries(t, vec))
		g.record(dA, lloprotocol.StreamValues{3: dec(7)}, t0.Add(2*time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"3", "benchmark", dA.String()}: 7,
		}, observedSeries(t, vec))
	})

	t.Run("config digest series expire after the TTL", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(1)}, t0)
		g.record(dB, lloprotocol.StreamValues{1: dec(2)}, t0.Add(5*time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"1", "benchmark", dA.String()}: 1,
			{"1", "benchmark", dB.String()}: 2,
		}, observedSeries(t, vec))
		g.record(dB, lloprotocol.StreamValues{1: dec(2)}, t0.Add(11*time.Second))
		assert.Equal(t, map[[3]string]float64{
			{"1", "benchmark", dB.String()}: 2,
		}, observedSeries(t, vec))
	})

	t.Run("deleteAll removes every series", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(1), 2: quote(1, 2, 3)}, t0)
		g.record(dB, lloprotocol.StreamValues{1: dec(1)}, t0)
		g.deleteAll()
		assert.Empty(t, observedSeries(t, vec))
	})

	t.Run("expire deletes the only digest after Observe stops", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		g.record(dA, lloprotocol.StreamValues{1: dec(1), 2: quote(1, 2, 3)}, t0)
		g.expire(t0.Add(observedValueDigestTTL))
		assert.Len(t, observedSeries(t, vec), 4, "series must survive up to the TTL")
		g.expire(t0.Add(observedValueDigestTTL + time.Nanosecond))
		assert.Empty(t, observedSeries(t, vec))
	})

	t.Run("concurrent Blue/Green records keep one series per stream and digest", func(t *testing.T) {
		t.Parallel()
		vec := newTestObservedValueVec()
		g := newObservedValueGauges(vec)
		var wg sync.WaitGroup
		for _, d := range []ocr2types.ConfigDigest{dA, dB} {
			wg.Go(func() {
				for i := range 100 {
					g.record(d, lloprotocol.StreamValues{1: dec(float64(i)), 2: quote(1, 2, 3)}, t0)
				}
			})
		}
		wg.Go(func() {
			for range 100 {
				observedSeries(t, vec)
				g.expire(t0)
			}
		})
		wg.Wait()
		assert.Equal(t, map[[3]string]float64{
			{"1", "benchmark", dA.String()}: 99,
			{"2", "bid", dA.String()}:       1,
			{"2", "benchmark", dA.String()}: 2,
			{"2", "ask", dA.String()}:       3,
			{"1", "benchmark", dB.String()}: 99,
			{"2", "bid", dB.String()}:       1,
			{"2", "benchmark", dB.String()}: 2,
			{"2", "ask", dB.String()}:       3,
		}, observedSeries(t, vec))
	})
}

func Test_decimalToFloat64(t *testing.T) {
	t.Parallel()
	// maxRelErr is 2 ulps: one rounding to 53 bits plus one for the division or multiplication by a power of ten.
	const maxRelErr = 2 * 0x1p-52
	check := func(t *testing.T, d decimal.Decimal) {
		t.Helper()
		want, _ := d.Float64()
		got := decimalToFloat64(d)
		if want == 0 {
			assert.Zero(t, got, "%s", d)
			return
		}
		assert.InDelta(t, 0, (got-want)/want, maxRelErr, "%s: got %v want %v", d, got, want)
	}

	cases := map[string]decimal.Decimal{
		"zero":                     decimal.Zero,
		"zero value":               {},
		"one":                      decimal.NewFromInt(1),
		"negative":                 decimal.RequireFromString("-63421.07"),
		"price":                    decimal.RequireFromString("63421.07"),
		"small fraction":           decimal.RequireFromString("0.000000081907646999640001"),
		"exponent -22":             decimal.New(123456789, -22),
		"exponent +22":             decimal.New(123456789, 22),
		"exponent -23 (fallback)":  decimal.New(123456789, -23),
		"exponent +23 (fallback)":  decimal.New(123456789, 23),
		"2^53+1 coefficient":       decimal.New(1<<53+1, 0),
		"max int64 coefficient":    decimal.New(1<<63-1, -4),
		"nanosecond timestamp":     decimal.NewFromInt(1794948135870123456),
		"65-bit coefficient":       decimal.NewFromBigInt(new(big.Int).Lsh(big.NewInt(1), 64), -10),
		"128-bit word boundary":    decimal.NewFromBigInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)), -20),
		"negative 100-bit":         decimal.NewFromBigInt(new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(12345), 86)), -16),
		"median of two (21 digit)": decimal.RequireFromString("63421.07").Add(decimal.RequireFromString("63421.08")).Div(decimal.NewFromInt(2)),
	}
	for name, d := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			check(t, d)
		})
	}

	t.Run("random parsed floats and medians", func(t *testing.T) {
		t.Parallel()
		r := rand.New(rand.NewPCG(1, 2))
		for range 10_000 {
			a := decimal.NewFromFloat((r.Float64() - 0.5) * 1e6)
			check(t, a)
			check(t, a.Add(decimal.NewFromFloat(r.Float64()*1e6)).Div(decimal.NewFromInt(2)))
		}
	})
}

// BenchmarkObservedValueGaugesRecord measures record() on the DON 4 shape: 2,724 decimal and 154 quote streams
// (3,186 series). Values mix the two coefficient shapes pipelines produce: a parsed float (about 17 digits) and an
// even-count median, whose Div(2) rounds to 16 decimals (about 21 digits).
func BenchmarkObservedValueGaugesRecord(b *testing.B) {
	const decimalStreams, quoteStreams = 2724, 154
	newValues := func(seed uint64) lloprotocol.StreamValues {
		r := rand.New(rand.NewPCG(seed, seed))
		price := func(i int) decimal.Decimal {
			d := decimal.NewFromFloat(r.Float64() * 100_000)
			if i%2 == 1 {
				d = d.Add(decimal.NewFromFloat(r.Float64() * 100_000)).Div(decimal.NewFromInt(2))
			}
			return d
		}
		sv := make(lloprotocol.StreamValues, decimalStreams+quoteStreams)
		for i := range decimalStreams {
			sv[llotypes.StreamID(i)] = lloprotocol.ToDecimal(price(i))
		}
		for i := range quoteStreams {
			p := price(i)
			sv[llotypes.StreamID(decimalStreams+i)] = &lloprotocol.Quote{Bid: p.Sub(decimal.NewFromInt(1)), Benchmark: p, Ask: p.Add(decimal.NewFromInt(1))}
		}
		return sv
	}
	digest := ocr2types.ConfigDigest{1}
	now := time.Unix(1_700_000_000, 0)

	b.Run("unchanged", func(b *testing.B) {
		g := newObservedValueGauges(newTestObservedValueVec())
		sv := newValues(1)
		g.record(digest, sv, now)
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			g.record(digest, sv, now)
		}
	})

	b.Run("all_changed", func(b *testing.B) {
		g := newObservedValueGauges(newTestObservedValueVec())
		svs := [2]lloprotocol.StreamValues{newValues(1), newValues(2)}
		g.record(digest, svs[1], now)
		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			g.record(digest, svs[i%2], now)
		}
	})
}
