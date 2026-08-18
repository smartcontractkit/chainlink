package vault

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func observeAtSeqNr(t *testing.T, r *ReportingPlugin, rdr *kv, seqNr uint64) *vaultcommon.Observations {
	t.Helper()
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)
	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	return obs
}

func attributedPendingQueueObservation(t *testing.T, observer commontypes.OracleID, obs *vaultcommon.Observations) types.AttributedObservation {
	t.Helper()
	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return types.AttributedObservation{Observer: observer, Observation: types.Observation(b)}
}

func stalledObservation() *vaultcommon.Observations {
	return &vaultcommon.Observations{
		SortNonce:               make([]byte, sortNonceLength),
		PendingQueueStallSignal: vaultcommon.PendingQueueStallSignal_PENDING_QUEUE_STALL_SIGNAL_STALLED,
	}
}

func continueObservation() *vaultcommon.Observations {
	return &vaultcommon.Observations{SortNonce: make([]byte, sortNonceLength)}
}

func TestPendingQueueStallTracker_Concurrent(t *testing.T) {
	t.Parallel()

	var tracker pendingQueueStallTracker
	const calls = 100
	var wg sync.WaitGroup
	for range calls {
		wg.Go(func() {
			tracker.record(7)
		})
	}
	wg.Wait()

	require.Equal(t, calls, tracker.record(7))
	require.Equal(t, 0, tracker.record(8))
}

func TestObservation_PendingQueueStallThreshold(t *testing.T) {
	t.Parallel()

	rdr := &kv{m: make(map[string]response)}
	r := newTestReportingPlugin(t, withVaultPendingQueueStallThreshold(1))

	require.False(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))
	require.True(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))
	require.True(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))
	require.False(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 2)))
}

func TestObservation_PendingQueueStallCounterAccumulatesWhileDisabled(t *testing.T) {
	t.Parallel()

	rdr := &kv{m: make(map[string]response)}
	r := newTestReportingPlugin(t)
	require.False(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))
	require.False(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))

	limiter, err := limits.MakeUpperBoundLimiter(
		limits.Factory{Settings: cresettings.DefaultGetter},
		settings.Int(2),
	)
	require.NoError(t, err)
	r.cfg.VaultPendingQueueStallThreshold = limiter

	require.True(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))
}

func TestObservation_ForceEmptyTakesPrecedenceOverStallSignal(t *testing.T) {
	t.Parallel()

	rdr := &kv{m: make(map[string]response)}
	r := newTestReportingPlugin(t, withVaultPendingQueueStallThreshold(2))
	r.cfg.VaultForceEmptyOCRRounds = limits.NewGateLimiter(true)

	require.True(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))

	r.cfg.VaultForceEmptyOCRRounds = limits.NewGateLimiter(false)
	require.False(t, observationSignalsPendingQueueStall(observeAtSeqNr(t, r, rdr, 1)))
}

func TestValidateObservation_PendingQueueStallSignal(t *testing.T) {
	t.Parallel()

	r := newTestReportingPlugin(t)
	rdr := &kv{m: make(map[string]response)}

	require.NoError(t, validatePendingQueueObservation(t, r, rdr, stalledObservation()))

	withContribution := stalledObservation()
	withContribution.Observations = []*vaultcommon.Observation{{Id: "request-1"}}
	require.ErrorContains(t, validatePendingQueueObservation(t, r, rdr, withContribution), "must not include contributions")

	withPendingItem := stalledObservation()
	withPendingItem.PendingQueueItems = [][]byte{{1}}
	require.ErrorContains(t, validatePendingQueueObservation(t, r, rdr, withPendingItem), "must not include contributions")

	unknown := continueObservation()
	unknown.PendingQueueStallSignal = vaultcommon.PendingQueueStallSignal(99)
	require.ErrorContains(t, validatePendingQueueObservation(t, r, rdr, unknown), "unknown pending queue stall signal")

	missingNonce := stalledObservation()
	missingNonce.SortNonce = nil
	require.ErrorContains(t, validatePendingQueueObservation(t, r, rdr, missingNonce), "sort nonce must be")
}

func TestObservationQuorum_PendingQueueStallSignal(t *testing.T) {
	t.Parallel()

	r := newTestReportingPlugin(t,
		withOnchainCfg(4, 1),
		withVaultIncludeInvalidPendingItemsEnabled(),
	)
	rdr := &kv{m: make(map[string]response)}
	writeDeleteSecretsPendingQueueItems(t, rdr, "request-1")

	aos := make([]types.AttributedObservation, 0, 3)
	aos = append(aos,
		attributedPendingQueueObservation(t, 0, stalledObservation()),
		attributedPendingQueueObservation(t, 1, stalledObservation()),
	)
	reached, err := r.ObservationQuorum(t.Context(), 1, types.AttributedQuery{}, aos, rdr, &blobber{})
	require.NoError(t, err)
	require.False(t, reached, "F+1 stall signals do not replace the overall 2F+1 observation requirement")

	aos = append(aos,
		attributedPendingQueueObservation(t, 2, continueObservation()),
	)
	reached, err = r.ObservationQuorum(t.Context(), 1, types.AttributedQuery{}, aos, rdr, &blobber{})
	require.NoError(t, err)
	require.True(t, reached)

	aos[1] = attributedPendingQueueObservation(t, 1, continueObservation())
	reached, err = r.ObservationQuorum(t.Context(), 1, types.AttributedQuery{}, aos, rdr, &blobber{})
	require.NoError(t, err)
	require.False(t, reached)
}

func TestStateTransition_PendingQueueStallSignalPurgesQueue(t *testing.T) {
	t.Parallel()

	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	rdr := &kv{m: make(map[string]response)}
	writeDeleteSecretsPendingQueueItems(t, rdr, "request-1", "request-2")

	localPayload := &vaultcommon.DeleteSecretsRequest{RequestId: "local-request"}
	require.NoError(t, r.store.Add(&vaulttypes.Request{
		IDVal:   "local-request",
		Payload: localPayload,
	}))

	normal := continueObservation()
	pending, err := anypb.New(localPayload)
	require.NoError(t, err)
	normal.PendingQueueItems = [][]byte{pending.Value}

	aos := []types.AttributedObservation{
		attributedPendingQueueObservation(t, 0, stalledObservation()),
		attributedPendingQueueObservation(t, 1, stalledObservation()),
		attributedPendingQueueObservation(t, 2, normal),
		attributedPendingQueueObservation(t, 3, normal),
	}
	out, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, rdr, &blobber{})
	require.NoError(t, err)

	outcomes := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(out, outcomes))
	require.Empty(t, outcomes.Outcomes)

	pendingQueue, err := newTestReadStore(t, rdr).GetPendingQueue(t.Context())
	require.NoError(t, err)
	require.Empty(t, pendingQueue)
	require.Equal(t, 1, r.store.Len(), "purge must leave node-local requests for normal TTL eviction")
}
