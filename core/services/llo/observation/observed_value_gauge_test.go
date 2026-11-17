package observation

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
)

func newTestObservedValueVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_observed_value"}, []string{"stream_id", "value_type", "config_digest"})
}

// observedSeries returns every series in vec keyed by {stream_id, value_type, config_digest}.
func observedSeries(t testing.TB, vec *prometheus.GaugeVec) map[[3]string]float64 {
	t.Helper()
	ch := make(chan prometheus.Metric)
	go func() {
		vec.Collect(ch)
		close(ch)
	}()
	out := make(map[[3]string]float64)
	for m := range ch {
		var pb io_prometheus_client.Metric
		require.NoError(t, m.Write(&pb))
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

// BenchmarkObservedValueGaugesRecord measures record() on the DON 4 shape: 2,724 decimal and 154 quote streams
// (3,186 series).
func BenchmarkObservedValueGaugesRecord(b *testing.B) {
	const decimalStreams, quoteStreams = 2724, 154
	newValues := func(seed float64) lloprotocol.StreamValues {
		sv := make(lloprotocol.StreamValues, decimalStreams+quoteStreams)
		for i := range decimalStreams {
			sv[llotypes.StreamID(i)] = dec(seed + float64(i))
		}
		for i := range quoteStreams {
			f := seed + float64(i)
			sv[llotypes.StreamID(decimalStreams+i)] = quote(f-1, f, f+1)
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
