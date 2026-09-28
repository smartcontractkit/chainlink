package bench

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	llov31 "github.com/smartcontractkit/chainlink-data-streams/llo/dev/v31"
	"github.com/smartcontractkit/chainlink-data-streams/llo/dev/v31/llotest"
)

// The v31 blob path carries every stream value: the pump broadcasts a blob of
// observed values, and each round fetches the blobs the observations reference.
// These tests drive the plugin directly (as TestParity does) because neither
// end of that path can be made to fail from the node-level integration tests:
// the blob knobs are not plumbed through the LLO job spec, and the real
// transport does not drop blobs on demand.

const (
	blobFailureN = 4
	blobFailureF = 1
)

// blobFailureWorkload is deliberately tiny: these tests assert on whether
// reports appear at all, not on throughput.
var blobFailureWorkload = workload{numChannels: 2, streamsPerChannel: 1}

// fetchController wraps the in-memory broadcaster/fetcher so blob fetches can
// be failed on demand, modeling a blob that expired or that no peer served
// before the fetch deadline.
type fetchController struct {
	*llotest.BlobBroadcastFetcher
	failing  atomic.Bool
	attempts atomic.Uint64
}

func newFetchController() *fetchController {
	return &fetchController{BlobBroadcastFetcher: llotest.NewBlobBroadcastFetcher()}
}

func (f *fetchController) FetchBlob(ctx context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	f.attempts.Add(1)
	if f.failing.Load() {
		return nil, errors.New("blob unavailable")
	}
	return f.BlobBroadcastFetcher.FetchBlob(ctx, handle)
}

var _ ocr3_1types.BlobBroadcastFetcher = &fetchController{}

// newV31 builds a v31 plugin over the given fetcher, applying opts to the
// factory params.
func newV31(tb testing.TB, defs llotypes.ChannelDefinitions, bbf ocr3_1types.BlobBroadcastFetcher, opts ...func(*llov31.PluginFactoryParams)) (ocr3_1types.ReportingPlugin[llotypes.ReportInfo], ocr3_1types.KeyValueDatabase) {
	tb.Helper()
	p, _, err := newV31Factory(defs, opts...).NewReportingPlugin(context.Background(), pluginConfig(blobFailureN, blobFailureF), bbf)
	require.NoError(tb, err)
	return p, newKVDB(tb)
}

// driveRounds drives rounds from startSeq and returns the first seqNr that
// produced reports for every channel, or 0 if none did within the budget. It
// waits for the blob pump between rounds exactly as the benchmark drivers do.
func driveRounds(tb testing.TB, p ocr3_1types.ReportingPlugin[llotypes.ReportInfo], db ocr3_1types.KeyValueDatabase, bbf ocr3_1types.BlobBroadcastFetcher, pump *llotest.BlobBroadcastFetcher, startSeq uint64, rounds int, channels int) (reportingSeq uint64) {
	tb.Helper()
	for i := range rounds {
		seqNr := startSeq + uint64(i) //nolint:gosec // G115: small round index
		before := pump.Broadcasts()
		reports, _ := v31Round(tb, p, db, bbf, seqNr, blobFailureN)
		if seqNr > 1 {
			waitForPump(tb, pump, before)
		}
		if len(reports) >= channels {
			return seqNr
		}
	}
	return 0
}

// TestBlobBroadcastFailure covers the broadcast end of the blob path: while the
// pump cannot broadcast, no stream values reach the round, so no channel is
// reportable; the protocol must not fail, and reporting must resume once
// broadcasting works again.
func TestBlobBroadcastFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}
	t.Parallel()

	defs, _ := blobFailureWorkload.channelDefinitions()
	bbf := llotest.NewBlobBroadcastFetcher()
	bbf.SetBroadcastError(errors.New("broadcast unavailable"))
	p, db := newV31(t, defs, bbf)

	budget := int(warmupRounds(blobFailureWorkload.numChannels)) //nolint:gosec // G115: small round budget
	seq := driveRounds(t, p, db, bbf, bbf, 1, budget, blobFailureWorkload.numChannels)
	require.Zero(t, seq, "no channel should be reportable while blobs cannot be broadcast")
	require.Positive(t, bbf.Broadcasts(), "the pump should have attempted to broadcast")

	bbf.SetBroadcastError(nil)
	seq = driveRounds(t, p, db, bbf, bbf, uint64(budget+1), budget, blobFailureWorkload.numChannels) //nolint:gosec // G115: small round budget
	require.NotZero(t, seq, "reporting should resume once blobs can be broadcast again")
}

// v31RoundErr drives one v31 round like v31Round, but returns the
// StateTransition error instead of failing the test, so the abort path can be
// asserted on. On error nothing is committed, mirroring libocr discarding the
// transaction and retrying the same sequence number.
func v31RoundErr(tb testing.TB, p ocr3_1types.ReportingPlugin[llotypes.ReportInfo], db ocr3_1types.KeyValueDatabase, bbf ocr3_1types.BlobBroadcastFetcher, seqNr uint64) ([]ocr3types.ReportPlus[llotypes.ReportInfo], error) {
	tb.Helper()
	ctx := context.Background()

	var obs []byte
	if seqNr > 1 {
		rtx, err := db.NewReadTransaction()
		require.NoError(tb, err)
		obs, err = p.Observation(ctx, seqNr, ocrtypes.AttributedQuery{}, rtx, bbf)
		rtx.Discard()
		require.NoError(tb, err)
	}
	aos := bootOrReplicate(obs, blobFailureN, seqNr)

	wtx, err := db.NewReadWriteTransaction()
	require.NoError(tb, err)
	prec, err := p.StateTransition(ctx, seqNr, ocrtypes.AttributedQuery{}, aos, wtx, bbf)
	if err != nil {
		wtx.Discard()
		return nil, err
	}
	require.NoError(tb, wtx.Commit())

	reports, err := p.Reports(ctx, seqNr, prec)
	require.NoError(tb, err)
	return reports, nil
}

// TestBlobFetchFailure covers the fetch end: the blobs are broadcast, but the
// round cannot retrieve them (expired, or no peer served them in time). A
// fetch failure is node-local, so dropping the observation would make the state
// transition non-deterministic across oracles; the round must abort instead,
// and reporting must resume once fetches succeed.
func TestBlobFetchFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}
	t.Parallel()

	defs, _ := blobFailureWorkload.channelDefinitions()
	fc := newFetchController()
	fc.failing.Store(true)
	p, db := newV31(t, defs, fc)

	budget := int(warmupRounds(blobFailureWorkload.numChannels)) //nolint:gosec // G115: small round budget

	// Drive until a round references a blob and therefore aborts. Earlier
	// rounds only carry channel definitions, which travel inline.
	var seqNr uint64 = 1
	var abortErr error
	for range budget {
		reports, err := v31RoundErr(t, p, db, fc, seqNr)
		if err != nil {
			abortErr = err
			break
		}
		require.Empty(t, reports, "no channel should be reportable while blobs cannot be fetched")
		before := fc.Broadcasts()
		seqNr++
		waitForPump(t, fc.BlobBroadcastFetcher, before)
	}
	require.ErrorContains(t, abortErr, "failed to fetch blob for observation",
		"a round referencing an unfetchable blob should abort")
	require.Positive(t, fc.attempts.Load(), "the round should have attempted to fetch the blobs it referenced")

	// The aborted round committed nothing, so the same sequence number is
	// retried, which is what libocr does.
	fc.failing.Store(false)
	seq := driveRounds(t, p, db, fc, fc.BlobBroadcastFetcher, seqNr, budget, blobFailureWorkload.numChannels)
	require.NotZero(t, seq, "reporting should resume once blobs can be fetched again")
}

// TestBlobRoundBounds covers the blob round knobs the factory validates: a blob
// lifetime that does not outlive the snapshots referencing it is rejected, as
// is one past the cap, while the tightest valid pairing still reports.
func TestBlobRoundBounds(t *testing.T) {
	t.Parallel()

	// The tightest valid pairing: a snapshot gathered for one sequence number
	// is consumed by a later one, so the blob must stay fetchable for
	// BlobFetchMarginRounds past MaxSnapshotRounds.
	const tightSnapshotRounds = 2
	const tightLifetimeRounds = tightSnapshotRounds + llov31.BlobFetchMarginRounds - 1

	tcs := []struct {
		name              string
		maxSnapshotRounds uint64
		blobLifetime      uint64
		wantErr           string
	}{
		{
			name:              "lifetime leaving no fetch margin is rejected",
			maxSnapshotRounds: 4,
			blobLifetime:      4,
			wantErr:           "leaves less than",
		},
		{
			name:         "lifetime past the cap is rejected",
			blobLifetime: llov31.MaxBlobLifetimeRounds + 1,
			wantErr:      "exceeds MaxBlobLifetimeRounds",
		},
		{
			name:              "tightest valid pairing reports",
			maxSnapshotRounds: tightSnapshotRounds,
			blobLifetime:      tightLifetimeRounds,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defs, _ := blobFailureWorkload.channelDefinitions()
			opt := func(params *llov31.PluginFactoryParams) {
				params.MaxSnapshotRounds = tc.maxSnapshotRounds
				params.BlobLifetimeRounds = tc.blobLifetime
			}
			bbf := llotest.NewBlobBroadcastFetcher()

			if tc.wantErr != "" {
				_, _, err := newV31Factory(defs, opt).NewReportingPlugin(context.Background(), pluginConfig(blobFailureN, blobFailureF), bbf)
				require.ErrorContains(t, err, tc.wantErr)
				return
			}

			if testing.Short() {
				t.Skip("too slow for testing.Short")
			}
			p, db := newV31(t, defs, bbf, opt)
			budget := int(warmupRounds(blobFailureWorkload.numChannels)) //nolint:gosec // G115: small round budget
			seq := driveRounds(t, p, db, bbf, bbf, 1, budget, blobFailureWorkload.numChannels)
			require.NotZero(t, seq, "the tightest valid blob round bounds should still report")
		})
	}
}
