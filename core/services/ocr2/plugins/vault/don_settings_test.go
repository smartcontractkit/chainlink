package vault

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// perOwnerKeyLimiter is a BoundLimiter that resolves the key length limit per
// owner from context, simulating a privileged per-owner limit override.
type perOwnerKeyLimiter struct {
	defaultBound pkgconfig.Size
	overrides    map[string]pkgconfig.Size
}

func (o *perOwnerKeyLimiter) Close() error { return nil }

func (o *perOwnerKeyLimiter) Limit(ctx context.Context) (pkgconfig.Size, error) {
	if v, ok := o.overrides[contexts.CREValue(ctx).Owner]; ok {
		return v, nil
	}
	return o.defaultBound, nil
}

func (o *perOwnerKeyLimiter) Check(ctx context.Context, n pkgconfig.Size) error {
	bound, _ := o.Limit(ctx)
	if n > bound {
		return limits.ErrorBoundLimited[pkgconfig.Size]{Limit: bound, Amount: n}
	}
	return nil
}

func nodeSettings(maxShare uint64) *vaultcommon.NodeSettings {
	return &vaultcommon.NodeSettings{
		MaxIdentifierKeyLengthBytes:       64,
		MaxIdentifierOwnerLengthBytes:     64,
		MaxIdentifierNamespaceLengthBytes: 64,
		MaxShareLengthBytes:               maxShare,
		MaxBlobPayloadBytes:               25600,
		MaxPendingQueueWriteSize:          1000,
		MaxRequestBatchSize:               10,
	}
}

func marshalObservationsWithSettings(t *testing.T, settings *vaultcommon.NodeSettings) []byte {
	t.Helper()
	obs := &vaultcommon.Observations{SortNonce: make([]byte, sortNonceLength), NodeSettings: settings}
	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return b
}

func TestPlugin_Observation_PopulatesNodeSettings(t *testing.T) {
	t.Parallel()
	t.Run("flag on", func(t *testing.T) {
		t.Parallel()
		r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
		rdr := &kv{m: make(map[string]response)}
		data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, nil)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		require.NoError(t, proto.Unmarshal(data, obs))
		require.NotNil(t, obs.NodeSettings)
		assert.Positive(t, obs.NodeSettings.MaxShareLengthBytes)
		assert.Positive(t, obs.NodeSettings.MaxRequestBatchSize)
	})

	t.Run("flag off", func(t *testing.T) {
		t.Parallel()
		r := newTestReportingPlugin(t)
		rdr := &kv{m: make(map[string]response)}
		data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, nil)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		require.NoError(t, proto.Unmarshal(data, obs))
		assert.Nil(t, obs.NodeSettings)
	})
}

func TestKVStore_DONSettings_RoundTrip(t *testing.T) {
	t.Parallel()
	store := newTestWriteStore(t, &kv{m: make(map[string]response)})
	settings := nodeSettings(600)
	require.NoError(t, store.WriteDONSettings(t.Context(), settings))

	got, err := store.GetDONSettings(t.Context())
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, proto.Equal(settings, got))
}

func TestPlugin_StateTransition_DONSettings_NoWriteWhenInitialQuorumIncomplete(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// Fewer than 2F+1 observations include NodeSettings (Byzantine omit attack).
	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: nodeSettings(600)},
		1: {NodeSettings: nodeSettings(600)},
		2: {},
		3: {},
		4: {},
	}

	merged, err := r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	assert.Nil(t, merged)

	stored, err := writeKV.GetDONSettings(t.Context())
	require.NoError(t, err)
	assert.Nil(t, stored)
}

func TestPlugin_ValidateObservation_RequiresNodeSettingsWhenConsensusEnabled(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	obs := &vaultcommon.Observations{SortNonce: make([]byte, sortNonceLength)}
	b, err := proto.Marshal(obs)
	require.NoError(t, err)

	err = r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "node settings are required")
}

func TestPlugin_ValidateObservation_AcceptsValidNodeSettings(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	settings := r.localNodeSettings(t.Context())
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.NoError(t, err)
}

func TestPlugin_ValidateObservation_AcceptsPeerNodeSettingsAboveLocalCfg(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	settings := r.localNodeSettings(t.Context())
	settings.MaxShareLengthBytes += 100
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.NoError(t, err)
}

func TestPlugin_ValidateObservation_RejectsMalformedNodeSettingsZeroLimit(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	settings := r.localNodeSettings(t.Context())
	settings.MaxShareLengthBytes = 0
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max_share_length_bytes must be positive")
}

func TestPlugin_ValidateObservation_IgnoresNodeSettingsWhenConsensusDisabled(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t)
	rdr := &kv{m: make(map[string]response)}

	// Arbitrary settings above local ceiling should be ignored when flag is off.
	settings := nodeSettings(99999)
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.NoError(t, err)
}

func TestPlugin_StateTransition_DONSettings_PerFieldQuorum_AllAgree(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	settings := nodeSettings(600)
	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: settings},
		1: {NodeSettings: settings},
		2: {NodeSettings: settings},
	}

	merged, err := r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	require.NotNil(t, merged)
	assert.Equal(t, uint64(600), merged.MaxShareLengthBytes)

	stored, err := writeKV.GetDONSettings(t.Context())
	require.NoError(t, err)
	assert.True(t, proto.Equal(merged, stored))
}

// TestPlugin_StateTransition_DONSettings_CommitLoggedAtInfoLevel guards the
// log level of DON settings commit events: system tests scan container logs
// for these lines, and containers run at the default info log level, so a
// revert to debug would silently break that coverage.
func TestPlugin_StateTransition_DONSettings_CommitLoggedAtInfoLevel(t *testing.T) {
	t.Parallel()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	r := newTestReportingPlugin(t, withLggr(lggr), withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	initial := nodeSettings(600)
	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: initial},
		1: {NodeSettings: initial},
		2: {NodeSettings: initial},
	}

	_, err := r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)

	updated := nodeSettings(800)
	for i := range uint8(3) {
		marshalledObs[i] = &vaultcommon.Observations{NodeSettings: updated}
	}
	_, err = r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)

	var sawInitialSeed, sawUpdate bool
	for _, entry := range observedLogs.All() {
		if entry.Level != zapcore.InfoLevel {
			continue
		}
		switch {
		case strings.Contains(entry.Message, "DON settings committed from per-field observation quorum"):
			sawInitialSeed = true
		case strings.Contains(entry.Message, "DON settings updated from per-field observation quorum"):
			sawUpdate = true
		}
	}
	require.True(t, sawInitialSeed, "expected info-level log for initial DON settings commit")
	require.True(t, sawUpdate, "expected info-level log for DON settings update")
}

func TestPlugin_StateTransition_DONSettings_PerFieldQuorum_OneFieldSplit(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	existing := nodeSettings(500)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), existing))

	base := nodeSettings(600)
	splitA := nodeSettings(700)
	splitB := nodeSettings(800)

	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: base},
		1: {NodeSettings: base},
		2: {NodeSettings: base},
		3: {NodeSettings: splitA},
	}

	merged, err := r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	assert.Equal(t, uint64(600), merged.MaxShareLengthBytes)

	// Add another observer with splitB - still no quorum on max_share (only 1 each for 700 and 800)
	marshalledObs[4] = &vaultcommon.Observations{NodeSettings: splitB}
	merged, err = r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	assert.Equal(t, uint64(600), merged.MaxShareLengthBytes)
}

func TestPlugin_StateTransition_DONSettings_NoOpWhenFlagOff(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	settings := nodeSettings(600)
	_, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
			{Observer: 1, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
			{Observer: 2, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
		},
		kvStore,
		nil,
	)
	require.NoError(t, err)

	stored, err := writeKV.GetDONSettings(t.Context())
	require.NoError(t, err)
	assert.Nil(t, stored)
}

func TestPlugin_StateTransition_DONSettings_EnforcesKVSettingsOverLocalCfg(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t,
		withVaultNodeSettingsConsensusEnabled(),
		withOnchainCfg(4, 1),
	)

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// Committed DON settings raise the share limit above the node-local default.
	kvSettings := nodeSettings(1000)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))

	resolver, err := r.ensureActiveSettingsForRound(t.Context(), 1, writeKV)
	require.NoError(t, err)
	shareLimit, err := resolver.maxShareLengthBytes(t.Context())
	require.NoError(t, err)
	assert.Equal(t, pkgconfig.Size(1000), shareLimit)
}

func TestPlugin_StateTransition_DONSettings_UsesLocalCfgWhenKVEmpty(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t,
		withVaultNodeSettingsConsensusEnabled(),
	)

	kvStore := &kv{m: make(map[string]response)}
	readKV := newTestReadStore(t, kvStore)

	resolver, err := r.ensureActiveSettingsForRound(t.Context(), 1, readKV)
	require.NoError(t, err)
	shareLimit, err := resolver.maxShareLengthBytes(t.Context())
	require.NoError(t, err)
	// No DON settings committed yet: falls back to the node-local default.
	assert.Equal(t, pkgconfig.Size(600), shareLimit)
}

func TestPlugin_StateTransition_DONSettings_LocalCfgIgnoredWhenKVSet(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	// Committed DON settings tighten the share limit below the node-local default.
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(400)))

	resolver, err := r.ensureActiveSettingsForRound(t.Context(), 1, writeKV)
	require.NoError(t, err)
	shareLimit, err := resolver.maxShareLengthBytes(t.Context())
	require.NoError(t, err)
	assert.Equal(t, pkgconfig.Size(400), shareLimit)
}

func TestPlugin_EnsureActiveSettingsForRound_FailsClosedOnKVReadErrorWhenConsensusEnabled(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	kvStore := &kv{m: map[string]response{
		donSettingsKey: {err: errors.New("kv unavailable")},
	}}
	readKV := newTestReadStore(t, kvStore)

	_, err := r.ensureActiveSettingsForRound(t.Context(), 1, readKV)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read DON settings from KV")
}

func TestPlugin_Observation_FailsClosedOnKVReadErrorWhenConsensusEnabled(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	kvStore := &kv{m: map[string]response{
		donSettingsKey: {err: errors.New("kv unavailable")},
	}}

	_, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, kvStore, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read DON settings from KV")
}

func TestPlugin_ValidateObservation_FailsClosedOnKVReadErrorWhenConsensusEnabled(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	kvStore := &kv{m: map[string]response{
		donSettingsKey: {err: errors.New("kv unavailable")},
	}}

	settings := r.localNodeSettings(t.Context())
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, kvStore, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read DON settings from KV")
}

func TestPlugin_Observation_MarshaledSizeWithinCapWhenNodeSettingsConsensusEnabled(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}
	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)
	require.LessOrEqual(t, len(data), r.maxObservationBytes)
}

func TestPlugin_ActiveSettings_EnforcesKVIdentifierKeyLimit(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withMaxIdentifierLengths(100, 100, 100))
	kvSettings := nodeSettings(600)
	kvSettings.MaxIdentifierKeyLengthBytes = 4

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))

	resolver, err := r.ensureActiveSettingsForRound(t.Context(), 1, writeKV)
	require.NoError(t, err)

	limits, err := resolver.secretIdentifierLimits(t.Context())
	require.NoError(t, err)

	err = r.validator.ValidateSecretIdentifier(t.Context(), "longkey", "owner", "ns", &limits)
	require.Error(t, err)

	err = r.validator.ValidateSecretIdentifier(t.Context(), "key", "owner", "ns", &limits)
	require.NoError(t, err)
}

func TestPlugin_ActiveSettings_PrivilegedOwnerKeepsOverrideOverDONLimits(t *testing.T) {
	t.Parallel()
	const (
		defaultKeyLimit = 5 * pkgconfig.Byte
		privilegedOwner = "privilegedowner"
		privilegedLimit = 20 * pkgconfig.Byte
	)

	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	// Swap in a per-owner-aware key limiter: the privileged owner may use
	// longer keys than the configured default.
	r.validator = vaultcap.NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		&perOwnerKeyLimiter{
			defaultBound: defaultKeyLimit,
			overrides:    map[string]pkgconfig.Size{privilegedOwner: privilegedLimit},
		},
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// DON-wide key limit is committed at 4 bytes, below both the default (5)
	// and the privileged override (20).
	kvSettings := nodeSettings(600)
	kvSettings.MaxIdentifierKeyLengthBytes = 4
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))
	_, err := r.ensureActiveSettingsForRound(t.Context(), 1, writeKV)
	require.NoError(t, err)

	longKey := "averylongkeyname" // 16 bytes: above the DON limit (4) and default (5), within the privileged override (20)

	// Privileged owner keeps their per-owner allowance on top of the DON baseline.
	require.NoError(t, r.checkSecretIdentifier(t.Context(), longKey, privilegedOwner, "ns"))

	// Regular owner is held to the DON-wide limit.
	err = r.checkSecretIdentifier(t.Context(), longKey, "regularowner", "ns")
	require.Error(t, err)
	require.ErrorContains(t, err, "key exceeds maximum length")
}

func TestPlugin_StateTransition_DONSettings_PersistedAfterTransitionNotAppliedSameRound(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t,
		withVaultNodeSettingsConsensusEnabled(),
		withOnchainCfg(4, 1),
	)

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// Round starts with a 500-byte share limit committed in KV.
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(500)))

	// All oracles advertise an 800-byte share limit on this round.
	settings := nodeSettings(800)
	aos := []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
		{Observer: 1, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
		{Observer: 2, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
	}

	_, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, kvStore, nil)
	require.NoError(t, err)

	// Same-round enforcement still reflects committed settings from round start.
	shareLimit, err := r.activeSettings.Load().maxShareLengthBytes(t.Context())
	require.NoError(t, err)
	assert.Equal(t, pkgconfig.Size(500), shareLimit)

	stored, err := writeKV.GetDONSettings(t.Context())
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, uint64(800), stored.MaxShareLengthBytes)
}

// TestPlugin_ActiveSettings_ConcurrentAccess reproduces the libocr calling
// pattern that makes ValidateObservation impure: it runs on background
// goroutines concurrently with the round loop, all touching the active
// settings resolver. Run with -race, this guards the atomic publication of
// activeSettings.
func TestPlugin_ActiveSettings_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(600)))

	b := marshalObservationsWithSettings(t, r.localNodeSettings(t.Context()))

	const goroutines = 8
	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Go(func() {
			seqNr := uint64(1 + i%2)
			for range 25 {
				err := r.ValidateObservation(t.Context(), seqNr, types.AttributedQuery{}, types.AttributedObservation{
					Observer:    0,
					Observation: types.Observation(b),
				}, kvStore, nil)
				assert.NoError(t, err) //nolint:testifylint // require.NoError inside a goroutine is unsafe

				_, err = r.activeSettings.Load().maxShareLengthBytes(t.Context())
				assert.NoError(t, err)
			}
		})
	}
	wg.Wait()
}

// TestPlugin_ActiveSettings_RoundPinnedAgainstMidComputationSwap guards that
// a round's settings stay pinned for the whole call, even when a stale
// goroutine from an earlier round swaps the shared resolver slot mid-flight.
func TestPlugin_ActiveSettings_RoundPinnedAgainstMidComputationSwap(t *testing.T) {
	t.Parallel()
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// Round 1 pins the committed 500-byte share limit.
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(500)))
	round1, err := r.ensureActiveSettingsForRound(t.Context(), 1, writeKV)
	require.NoError(t, err)

	// Round 1's state transition commits a raised 800-byte share limit.
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(800)))
	round2, err := r.ensureActiveSettingsForRound(t.Context(), 2, writeKV)
	require.NoError(t, err)
	round2Limit, err := round2.maxShareLengthBytes(t.Context())
	require.NoError(t, err)
	require.Equal(t, pkgconfig.Size(800), round2Limit)

	// A stale round-1 goroutine swaps the shared slot while round 2's call is
	// in flight.
	r.activeSettings.Store(round1)

	// The slot now resolves the old limit...
	slotLimit, err := r.activeSettings.Load().maxShareLengthBytes(t.Context())
	require.NoError(t, err)
	require.Equal(t, pkgconfig.Size(500), slotLimit)

	// ...but the round-2 call, pinned via the context, accepts a 700-byte
	// share: above round 1's limit, within round 2's.
	ctx := withRoundDonSettings(t.Context(), round2)
	require.NoError(t, r.checkMaxShareLength(ctx, 700))

	// A caller with no pinned round falls back to the swapped slot.
	require.Error(t, r.checkMaxShareLength(t.Context(), 700))
}
