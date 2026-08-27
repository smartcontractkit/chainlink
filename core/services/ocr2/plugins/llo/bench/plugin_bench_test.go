package bench

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo/dev/v31/llotest"
)

const (
	benchN = 4 // total oracles
	benchF = 1 // fault tolerance (2f+1 = 3 quorum)
)

// benchWorkloads is the scaling matrix. Total observed streams per round is
// numChannels*streamsPerChannel and is kept at or below the protocol limit
// (MaxObservationStreamValuesLength = 10_000).
var benchWorkloads = []workload{
	{numChannels: 10, streamsPerChannel: 1},
	{numChannels: 100, streamsPerChannel: 1},
	{numChannels: 100, streamsPerChannel: 10},
	{numChannels: 1000, streamsPerChannel: 1},
	{numChannels: 1000, streamsPerChannel: 10},
}

// ---------------------------------------------------------------------------
// Correctness gate
// ---------------------------------------------------------------------------

// TestParity asserts that, for the same workload, both plugins reach a
// report-producing steady state and emit the same number of reports in the
// same format. This guards the benchmark: if the two drivers diverge, the
// latency numbers are not comparing like for like.
func TestParity(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	for _, w := range benchWorkloads {
		t.Run(w.String(), func(t *testing.T) {
			t.Parallel()
			defs, _ := w.channelDefinitions()

			p30 := buildV30(t, defs, benchN, benchF)
			prev, seq := warmV30(t, p30, benchN, w.numChannels)
			_, reports30, _ := v30Round(t, p30, seq, prev, benchN)

			p31, db, bbf := buildV31(t, defs, benchN, benchF)
			seq31 := warmV31(t, p31, db, bbf, benchN, w.numChannels)
			reports31, _ := v31Round(t, p31, db, bbf, seq31, benchN)

			require.NotEmpty(t, reports30, "v30 produced no reports")
			require.Len(t, reports30, w.numChannels, "v30 should report every channel")
			require.Len(t, reports31, len(reports30), "v30 and v31 must produce the same number of reports")

			for i := range reports30 {
				require.Equal(t, llotypes.ReportFormatJSON, reports30[i].ReportWithInfo.Info.ReportFormat)
			}
			for i := range reports31 {
				require.Equal(t, llotypes.ReportFormatJSON, reports31[i].ReportWithInfo.Info.ReportFormat)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Full-round benchmark
// ---------------------------------------------------------------------------

// BenchmarkFullRound measures a complete steady-state round for each plugin:
// Observation → (Outcome | StateTransition) → Reports. Besides ns/op and the
// -benchmem allocation metrics, it reports per-round size characteristics:
//   - v30: outcome_B (the re-serialized state blob), obs_B, report_B, reports
//   - v31: precursor_B, obs_B, report_B, reports, and KeyValueState I/O volume
//     (kvread_B, kvwrite_B, kvkeys) — the incremental state cost v30 lacks.
func BenchmarkFullRound(b *testing.B) {
	for _, w := range benchWorkloads {
		defs, _ := w.channelDefinitions()

		b.Run(w.String()+"/v30", func(b *testing.B) {
			p := buildV30(b, defs, benchN, benchF)
			prev, seq := warmV30(b, p, benchN, w.numChannels)

			// Size probe (untimed). Reported after the loop because
			// b.ResetTimer() deletes user metrics.
			outcome, reports, obs := v30Round(b, p, seq, prev, benchN)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				v30Round(b, p, seq, prev, benchN)
			}
			b.StopTimer()
			b.ReportMetric(float64(len(outcome)), "outcome_B/op")
			b.ReportMetric(float64(len(obs)), "obs_B/op")
			b.ReportMetric(float64(reportBytes(reports)), "report_B/op")
			b.ReportMetric(float64(len(reports)), "reports/op")
		})

		b.Run(w.String()+"/v31", func(b *testing.B) {
			p, db, bbf := buildV31(b, defs, benchN, benchF)
			seq := warmV31(b, p, db, bbf, benchN, w.numChannels)

			// Size/IO probe (untimed). Reported after the loop because
			// b.ResetTimer() deletes user metrics.
			m := probeV31(b, p, db, bbf, seq, benchN)
			seq++

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				before := bbf.Broadcasts()
				v31Round(b, p, db, bbf, seq, benchN)
				seq++
				// Untimed: let the pump park the snapshot the next round
				// consumes, so every measured round carries stream values.
				// Without this the tight loop outruns the pump and measures
				// value-less rounds.
				b.StopTimer()
				waitForPump(b, bbf, before)
				b.StartTimer()
			}
			b.StopTimer()
			b.ReportMetric(float64(m.precBytes), "precursor_B/op")
			b.ReportMetric(float64(m.obsBytes), "obs_B/op")
			b.ReportMetric(float64(m.blobBytes), "blob_B/op")
			b.ReportMetric(float64(m.reportBytesTotal), "report_B/op")
			b.ReportMetric(float64(m.reports), "reports/op")
			b.ReportMetric(float64(m.kvReadBytes), "kvread_B/op")
			b.ReportMetric(float64(m.kvWriteBytes), "kvwrite_B/op")
			b.ReportMetric(float64(m.kvKeys), "kvkeys/op")
		})
	}
}

// ---------------------------------------------------------------------------
// Per-stage benchmarks (localize where the cost is spent)
// ---------------------------------------------------------------------------

// BenchmarkObservation measures only the Observation stage (data-source gather
// + observation encode; for v31 also the KeyValueState read).
func BenchmarkObservation(b *testing.B) {
	for _, w := range benchWorkloads {
		defs, _ := w.channelDefinitions()

		b.Run(w.String()+"/v30", func(b *testing.B) {
			p := buildV30(b, defs, benchN, benchF)
			prev, seq := warmV30(b, p, benchN, w.numChannels)
			outctx := ocr3types.OutcomeContext{SeqNr: seq, PreviousOutcome: prev}
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if _, err := p.Observation(ctx, outctx, nil); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(w.String()+"/v31", func(b *testing.B) {
			p, db, bbf := buildV31(b, defs, benchN, benchF)
			seq := warmV31(b, p, db, bbf, benchN, w.numChannels)
			ctx := context.Background()
			rtx, err := db.NewReadTransaction()
			require.NoError(b, err)
			defer rtx.Discard()
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if _, err := p.Observation(ctx, seq, ocrtypes.AttributedQuery{}, rtx, bbf); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkStateAdvance measures the state-generation stage: v30's Outcome
// (decode previous outcome, aggregate, re-encode full outcome) versus v31's
// StateTransition (incremental KeyValueState reads/writes) followed by the
// KeyValueDatabase commit that libocr performs after every StateTransition.
func BenchmarkStateAdvance(b *testing.B) {
	for _, w := range benchWorkloads {
		defs, _ := w.channelDefinitions()

		b.Run(w.String()+"/v30_Outcome", func(b *testing.B) {
			p := buildV30(b, defs, benchN, benchF)
			prev, seq := warmV30(b, p, benchN, w.numChannels)
			ctx := context.Background()
			outctx := ocr3types.OutcomeContext{SeqNr: seq, PreviousOutcome: prev}
			obs, err := p.Observation(ctx, outctx, nil)
			require.NoError(b, err)
			aos := replicate(obs, benchN)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if _, err := p.Outcome(ctx, outctx, nil, aos); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(w.String()+"/v31_StateTransition", func(b *testing.B) {
			p, db, bbf := buildV31(b, defs, benchN, benchF)
			seq := warmV31(b, p, db, bbf, benchN, w.numChannels)
			ctx := context.Background()
			rtx, err := db.NewReadTransaction()
			require.NoError(b, err)
			obs, err := p.Observation(ctx, seq, ocrtypes.AttributedQuery{}, rtx, bbf)
			rtx.Discard()
			require.NoError(b, err)
			aos := replicate(obs, benchN)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				wtx, err := db.NewReadWriteTransaction()
				if err != nil {
					b.Fatal(err)
				}
				if _, err := p.StateTransition(ctx, seq, ocrtypes.AttributedQuery{}, aos, wtx, bbf); err != nil {
					b.Fatal(err)
				}
				if err := wtx.Commit(); err != nil {
					b.Fatal(err)
				}
				seq++
			}
		})
	}
}

// BenchmarkReports measures only the Reports stage (turning a committed
// outcome/precursor into signed report payloads).
func BenchmarkReports(b *testing.B) {
	for _, w := range benchWorkloads {
		defs, _ := w.channelDefinitions()

		b.Run(w.String()+"/v30", func(b *testing.B) {
			p := buildV30(b, defs, benchN, benchF)
			prev, seq := warmV30(b, p, benchN, w.numChannels)
			outcome, reports, _ := v30Round(b, p, seq, prev, benchN)
			require.NotEmpty(b, reports)
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if _, err := p.Reports(ctx, seq, outcome); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(w.String()+"/v31", func(b *testing.B) {
			p, db, bbf := buildV31(b, defs, benchN, benchF)
			seq := warmV31(b, p, db, bbf, benchN, w.numChannels)
			// Produce a committed precursor to feed Reports repeatedly.
			ctx := context.Background()
			rtx, err := db.NewReadTransaction()
			require.NoError(b, err)
			obs, err := p.Observation(ctx, seq, ocrtypes.AttributedQuery{}, rtx, bbf)
			rtx.Discard()
			require.NoError(b, err)
			wtx, err := db.NewReadWriteTransaction()
			require.NoError(b, err)
			prec, err := p.StateTransition(ctx, seq, ocrtypes.AttributedQuery{}, replicate(obs, benchN), wtx, bbf)
			require.NoError(b, err)
			require.NoError(b, wtx.Commit())
			reports, err := p.Reports(ctx, seq, prec)
			require.NoError(b, err)
			require.NotEmpty(b, reports)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if _, err := p.Reports(ctx, seq, prec); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// v31 size probe
// ---------------------------------------------------------------------------

type v31Sizes struct {
	precBytes        int
	obsBytes         int
	reportBytesTotal int
	reports          int
	kvReadBytes      int
	kvWriteBytes     int
	kvKeys           int
	blobBytes        int
}

// probeV31 runs one instrumented round to capture per-round size and
// KeyValueState I/O characteristics. It commits, advancing the state by one
// seqNr.
func probeV31(tb testing.TB, p ocr3_1types.ReportingPlugin[llotypes.ReportInfo], db ocr3_1types.KeyValueDatabase, bbf *llotest.BlobBroadcastFetcher, seqNr uint64, n int) v31Sizes {
	tb.Helper()
	ctx := context.Background()

	rtx, err := db.NewReadTransaction()
	require.NoError(tb, err)
	cr := &countingReader{inner: rtx}
	// v31 observations reference stream values by blob handle, so obs_B/op alone
	// understates what a round puts on the wire; blobBytes captures the payload
	// the blob broadcast for this round added.
	blobBytesBefore := bbf.BroadcastBytes()
	obs, err := p.Observation(ctx, seqNr, ocrtypes.AttributedQuery{}, cr, bbf)
	rtx.Discard()
	require.NoError(tb, err)

	wtx, err := db.NewReadWriteTransaction()
	require.NoError(tb, err)
	crw := &countingRW{inner: wtx}
	prec, err := p.StateTransition(ctx, seqNr, ocrtypes.AttributedQuery{}, replicate(obs, n), crw, bbf)
	require.NoError(tb, err)
	require.NoError(tb, wtx.Commit())

	reports, err := p.Reports(ctx, seqNr, prec)
	require.NoError(tb, err)

	return v31Sizes{
		precBytes:        len(prec),
		obsBytes:         len(obs),
		reportBytesTotal: reportBytes(reports),
		reports:          len(reports),
		kvReadBytes:      cr.readBytes + crw.readBytes,
		kvWriteBytes:     crw.writeBytes,
		kvKeys:           crw.keys,
		blobBytes:        bbf.BroadcastBytes() - blobBytesBefore,
	}
}
