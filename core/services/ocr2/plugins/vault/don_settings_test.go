package vault

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func nodeSettings(optimizations bool, maxShare uint64) *vaultcommon.NodeSettings {
	return &vaultcommon.NodeSettings{
		VaultOptimizationsEnabled:         optimizations,
		VaultForceEmptyOcrRounds:          false,
		MaxCiphertextLengthBytes:          2048,
		MaxIdentifierKeyLengthBytes:       64,
		MaxIdentifierOwnerLengthBytes:     64,
		MaxIdentifierNamespaceLengthBytes: 64,
		MaxShareLengthBytes:               maxShare,
		MaxBlobPayloadBytes:               25600,
		MaxPendingQueueWriteSize:          1000,
		MaxRequestBatchSize:                 10,
	}
}

func marshalObservationsWithSettings(t *testing.T, settings *vaultcommon.NodeSettings) []byte {
	t.Helper()
	obs := &vaultcommon.Observations{NodeSettings: settings}
	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return b
}

func TestPlugin_Observation_PopulatesNodeSettings(t *testing.T) {
	t.Run("flag on", func(t *testing.T) {
		r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
		rdr := &kv{m: make(map[string]response)}
		data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, nil)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		require.NoError(t, proto.Unmarshal(data, obs))
		require.NotNil(t, obs.NodeSettings)
		assert.False(t, obs.NodeSettings.VaultOptimizationsEnabled)
		assert.Equal(t, uint64(1024), obs.NodeSettings.MaxCiphertextLengthBytes)
	})

	t.Run("flag off", func(t *testing.T) {
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
	store := newTestWriteStore(t, &kv{m: make(map[string]response)})
	settings := nodeSettings(true, 600)
	require.NoError(t, store.WriteDONSettings(t.Context(), settings))

	got, err := store.GetDONSettings(t.Context())
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, proto.Equal(settings, got))
}

func TestPlugin_StateTransition_DONSettings_NoWriteWhenInitialQuorumIncomplete(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// Fewer than 2F+1 observations include NodeSettings (Byzantine omit attack).
	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: nodeSettings(true, 600)},
		1: {NodeSettings: nodeSettings(true, 600)},
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
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	obs := &vaultcommon.Observations{}
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
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	settings := r.localNodeSettings(t.Context())
	settings.MaxShareLengthBytes = settings.MaxShareLengthBytes + 100
	settings.VaultOptimizationsEnabled = true
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.NoError(t, err)
}

func TestPlugin_ValidateObservation_RejectsMalformedNodeSettingsZeroLimit(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}

	settings := r.localNodeSettings(t.Context())
	settings.MaxCiphertextLengthBytes = 0
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max_ciphertext_length_bytes must be positive")
}

func TestPlugin_ValidateObservation_IgnoresNodeSettingsWhenConsensusDisabled(t *testing.T) {
	r := newTestReportingPlugin(t)
	rdr := &kv{m: make(map[string]response)}

	// Arbitrary settings above local ceiling should be ignored when flag is off.
	settings := nodeSettings(true, 99999)
	b := marshalObservationsWithSettings(t, settings)

	err := r.ValidateObservation(t.Context(), 1, types.AttributedQuery{}, types.AttributedObservation{
		Observer:    0,
		Observation: types.Observation(b),
	}, rdr, nil)
	require.NoError(t, err)
}

func TestPlugin_StateTransition_DONSettings_PerFieldQuorum_AllAgree(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	settings := nodeSettings(true, 600)
	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: settings},
		1: {NodeSettings: settings},
		2: {NodeSettings: settings},
	}

	merged, err := r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	require.NotNil(t, merged)
	assert.True(t, merged.VaultOptimizationsEnabled)

	stored, err := writeKV.GetDONSettings(t.Context())
	require.NoError(t, err)
	assert.True(t, proto.Equal(merged, stored))
}

func TestPlugin_StateTransition_DONSettings_PerFieldQuorum_OneFieldSplit(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	existing := nodeSettings(false, 500)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), existing))

	base := nodeSettings(true, 600)
	splitA := nodeSettings(true, 700)
	splitB := nodeSettings(true, 800)

	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: base},
		1: {NodeSettings: base},
		2: {NodeSettings: base},
		3: {NodeSettings: splitA},
	}

	merged, err := r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	assert.True(t, merged.VaultOptimizationsEnabled)
	assert.Equal(t, uint64(600), merged.MaxShareLengthBytes)

	// Add another observer with splitB - still no quorum on max_share (only 1 each for 700 and 800)
	marshalledObs[4] = &vaultcommon.Observations{NodeSettings: splitB}
	merged, err = r.mergeAndPersistDONSettingsFromObservationQuorum(t.Context(), writeKV, marshalledObs)
	require.NoError(t, err)
	assert.Equal(t, uint64(600), merged.MaxShareLengthBytes)
}

func TestPlugin_StateTransition_DONSettings_NoOpWhenFlagOff(t *testing.T) {
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	settings := nodeSettings(true, 600)
	marshalledObs := map[uint8]*vaultcommon.Observations{
		0: {NodeSettings: settings},
		1: {NodeSettings: settings},
		2: {NodeSettings: settings},
	}

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
	_ = marshalledObs
}

func TestPlugin_StateTransition_DONSettings_EnforcesKVSettingsOverLocalCfg(t *testing.T) {
	r := newTestReportingPlugin(t,
		withVaultNodeSettingsConsensusEnabled(),
		withOnchainCfg(4, 1),
	)
	r.cfg.VaultOptimizationsEnabled = limits.NewGateLimiter(false)

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	kvSettings := nodeSettings(true, 600)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))

	require.NoError(t, r.ensureActiveSettingsForRound(t.Context(), 1, writeKV))
	assert.True(t, r.activeSettings.optimizationsEnabled())
}

func TestPlugin_StateTransition_DONSettings_UsesLocalCfgWhenKVEmpty(t *testing.T) {
	r := newTestReportingPlugin(t,
		withVaultNodeSettingsConsensusEnabled(),
		withVaultOptimizationsEnabled(),
	)
	kvStore := &kv{m: make(map[string]response)}
	readKV := newTestReadStore(t, kvStore)

	require.NoError(t, r.ensureActiveSettingsForRound(t.Context(), 1, readKV))
	assert.True(t, r.activeSettings.optimizationsEnabled())
}

func TestPlugin_StateTransition_DONSettings_LocalCfgIgnoredWhenKVSet(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	r.cfg.VaultOptimizationsEnabled = limits.NewGateLimiter(false)

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(true, 600)))

	require.NoError(t, r.ensureActiveSettingsForRound(t.Context(), 1, writeKV))
	assert.True(t, r.activeSettings.optimizationsEnabled())
}

func TestPlugin_EnsureActiveSettingsForRound_FailsClosedOnKVReadErrorWhenConsensusEnabled(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	kvStore := &kv{m: map[string]response{
		donSettingsKey: {err: errors.New("kv unavailable")},
	}}
	readKV := newTestReadStore(t, kvStore)

	err := r.ensureActiveSettingsForRound(t.Context(), 1, readKV)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read DON settings from KV")
}

func TestPlugin_Observation_FailsClosedOnKVReadErrorWhenConsensusEnabled(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	kvStore := &kv{m: map[string]response{
		donSettingsKey: {err: errors.New("kv unavailable")},
	}}

	_, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, kvStore, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read DON settings from KV")
}

func TestPlugin_ValidateObservation_FailsClosedOnKVReadErrorWhenConsensusEnabled(t *testing.T) {
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
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled())
	rdr := &kv{m: make(map[string]response)}
	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)
	require.LessOrEqual(t, len(data), r.maxObservationBytes)
}

func TestPlugin_ActiveSettings_EnforcesKVCiphertextLimit(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withMaxCiphertextLengthBytes(1024))
	kvSettings := nodeSettings(false, 600)
	kvSettings.MaxCiphertextLengthBytes = 8

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))

	require.NoError(t, r.ensureActiveSettingsForRound(t.Context(), 1, writeKV))

	// 16 raw bytes -> 32 hex chars; limit is 8 bytes.
	tooLarge := "0102030405060708090a0b0c0d0e0f10"
	err := r.validator.ValidateCiphertextSize(t.Context(), "owner", tooLarge, r.activeSettings.donSettings())
	require.Error(t, err)

	small := "01020304"
	err = r.validator.ValidateCiphertextSize(t.Context(), "owner", small, r.activeSettings.donSettings())
	require.NoError(t, err)
}

func TestPlugin_ActiveSettings_EnforcesKVIdentifierKeyLimit(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withMaxIdentifierLengths(100, 100, 100))
	kvSettings := nodeSettings(false, 600)
	kvSettings.MaxIdentifierKeyLengthBytes = 4

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))

	require.NoError(t, r.ensureActiveSettingsForRound(t.Context(), 1, writeKV))

	err := r.validator.ValidateSecretIdentifier(t.Context(), "longkey", "owner", "ns", r.activeSettings.donSettings())
	require.Error(t, err)

	err = r.validator.ValidateSecretIdentifier(t.Context(), "key", "owner", "ns", r.activeSettings.donSettings())
	require.NoError(t, err)
}

func TestPlugin_ActiveSettings_UsesKVLimitAboveLocalCfg(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultNodeSettingsConsensusEnabled(), withMaxCiphertextLengthBytes(16))
	kvSettings := nodeSettings(false, 600)
	kvSettings.MaxCiphertextLengthBytes = 32

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), kvSettings))

	require.NoError(t, r.ensureActiveSettingsForRound(t.Context(), 1, writeKV))

	// 20 raw bytes; above local cfg (16) but within DON KV (32).
	withinDON := "0102030405060708090a0b0c0d0e0f1011121314"
	err := r.validator.ValidateCiphertextSize(t.Context(), "owner", withinDON, r.activeSettings.donSettings())
	require.NoError(t, err)

	// 33 raw bytes exceeds DON KV limit.
	tooLargeForDON := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122"
	err = r.validator.ValidateCiphertextSize(t.Context(), "owner", tooLargeForDON, r.activeSettings.donSettings())
	require.Error(t, err)
}

func TestPlugin_StateTransition_DONSettings_PersistedAfterTransitionNotAppliedSameRound(t *testing.T) {
	r := newTestReportingPlugin(t,
		withVaultNodeSettingsConsensusEnabled(),
		withOnchainCfg(4, 1),
	)
	r.cfg.VaultOptimizationsEnabled = limits.NewGateLimiter(false)

	kvStore := &kv{m: make(map[string]response)}
	writeKV := newTestWriteStore(t, kvStore)

	// Round starts with optimizations off in KV.
	require.NoError(t, writeKV.WriteDONSettings(t.Context(), nodeSettings(false, 600)))

	// All oracles advertise optimizations on this round.
	settings := nodeSettings(true, 600)
	aos := []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
		{Observer: 1, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
		{Observer: 2, Observation: types.Observation(marshalObservationsWithSettings(t, settings))},
	}

	_, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, kvStore, nil)
	require.NoError(t, err)

	// Same-round enforcement still reflects committed settings from round start.
	assert.False(t, r.activeSettings.optimizationsEnabled())

	stored, err := writeKV.GetDONSettings(t.Context())
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.True(t, stored.VaultOptimizationsEnabled)
}
