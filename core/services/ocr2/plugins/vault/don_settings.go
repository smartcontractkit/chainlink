package vault

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
)

type donSettingsResolver struct {
	r              *ReportingPlugin
	store          ReadKVStore
	consensus      bool
	initialized    bool
	initErr        error
	cachedSettings *vaultcommon.NodeSettings

	// seqNr is the OCR round this resolver was installed for. Set before
	// publication and never mutated afterwards.
	seqNr uint64
}

func (r *ReportingPlugin) newDonSettingsResolver(ctx context.Context, store ReadKVStore) *donSettingsResolver {
	return &donSettingsResolver{
		r:         r,
		store:     store,
		consensus: r.nodeSettingsConsensusEnabled(ctx),
	}
}

func (r *ReportingPlugin) nodeSettingsConsensusEnabled(ctx context.Context) bool {
	return gateAllows(ctx, r.lggr, r.cfg.VaultNodeSettingsConsensusEnabled, "VaultNodeSettingsConsensusEnabled")
}

func (r *ReportingPlugin) localNodeSettings(ctx context.Context) *vaultcommon.NodeSettings {
	return populateLocalNodeSettings(ctx, r)
}

// ensureActiveSettingsForRound installs and returns the resolver holding the
// round's DON-agreed settings. Callers must use the returned resolver for the
// whole call via withRoundDonSettings: the shared slot may be swapped by a
// concurrent round mid-computation.
func (r *ReportingPlugin) ensureActiveSettingsForRound(ctx context.Context, seqNr uint64, store ReadKVStore) (*donSettingsResolver, error) {
	if active := r.activeSettings.Load(); active.initialized && active.seqNr == seqNr {
		return active, nil
	}
	resolver := r.newDonSettingsResolver(ctx, store)
	if err := resolver.init(ctx); err != nil {
		return nil, err
	}
	resolver.seqNr = seqNr
	r.activeSettings.Store(resolver)
	return resolver, nil
}

// roundDonSettingsCtxKey keys the round-pinned settings resolver in a call's
// context.
type roundDonSettingsCtxKey struct{}

// withRoundDonSettings pins the round's resolver in the context so settings
// resolved within the call are unaffected by concurrent slot swaps.
func withRoundDonSettings(ctx context.Context, d *donSettingsResolver) context.Context {
	return context.WithValue(ctx, roundDonSettingsCtxKey{}, d)
}

// donSettingsForCall returns the round-pinned resolver from the context, or
// the most recently installed resolver when no round is pinned.
func (r *ReportingPlugin) donSettingsForCall(ctx context.Context) *donSettingsResolver {
	if d, ok := ctx.Value(roundDonSettingsCtxKey{}).(*donSettingsResolver); ok && d != nil {
		return d
	}
	return r.activeSettings.Load()
}

func (d *donSettingsResolver) init(ctx context.Context) error {
	if d.initialized {
		return d.initErr
	}
	d.initialized = true

	if !d.consensus {
		return nil
	}

	settings, err := d.store.GetDONSettings(ctx)
	if err != nil {
		d.initErr = fmt.Errorf("failed to read DON settings from KV: %w", err)
		d.r.lggr.Errorw("failed to read DON settings from KV", "error", err)
		return d.initErr
	}
	if settings != nil {
		d.cachedSettings = settings
	}
	return nil
}

// donSettings returns the DON-agreed settings, or nil when no resolver is
// active (or none have been committed yet), in which case all resolution
// falls back to node-local configuration.
func (d *donSettingsResolver) donSettings() *vaultcommon.NodeSettings {
	if d == nil {
		return nil
	}
	return d.cachedSettings
}

func (r *ReportingPlugin) mergeAndPersistDONSettingsFromObservationQuorum(
	ctx context.Context,
	store WriteKVStore,
	marshalledObs map[uint8]*vaultcommon.Observations,
) (*vaultcommon.NodeSettings, error) {
	threshold := 2*r.onchainCfg.F + 1

	existing, err := store.GetDONSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read existing DON settings: %w", err)
	}
	var merged *vaultcommon.NodeSettings
	if existing != nil {
		merged = proto.Clone(existing).(*vaultcommon.NodeSettings)
	} else {
		merged = &vaultcommon.NodeSettings{}
	}

	allFieldsQuorum := applyDonSettingsQuorum(r.lggr, donSettingsBoolFields, merged, marshalledObs, threshold)
	allFieldsQuorum = applyDonSettingsQuorum(r.lggr, donSettingsUint64BaseFields, merged, marshalledObs, threshold) && allFieldsQuorum

	if existing == nil {
		// Initial seeding of DON settings, reject if not all fields reached quorum
		if !allFieldsQuorum {
			r.lggr.Debugw("DON settings initial seed deferred: not all fields reached quorum", "threshold", threshold)
			return nil, nil
		}
		if err := store.WriteDONSettings(ctx, merged); err != nil {
			return nil, fmt.Errorf("failed to write DON settings: %w", err)
		}
		// Logged at info level: DON settings commits are rare, DON-wide events
		// and are asserted on by system tests scanning container logs.
		r.lggr.Infow("DON settings committed from per-field observation quorum (initial seed)", "threshold", threshold)
		return merged, nil
	}

	if !proto.Equal(existing, merged) {
		if err := store.WriteDONSettings(ctx, merged); err != nil {
			return nil, fmt.Errorf("failed to write DON settings: %w", err)
		}
		r.lggr.Infow("DON settings updated from per-field observation quorum", "threshold", threshold)
	}
	return merged, nil
}

// Wrappers for readability at call sites
func (d *donSettingsResolver) maxBlobPayloadBytes(ctx context.Context) (pkgconfig.Size, error) {
	return donUint64FieldMaxBlobPayloadBytes.resolveSize(ctx, d)
}

func (d *donSettingsResolver) maxShareLengthBytes(ctx context.Context) (pkgconfig.Size, error) {
	return donUint64FieldMaxShareLengthBytes.resolveSize(ctx, d)
}

func (d *donSettingsResolver) checkMaxPendingQueueWriteSize(ctx context.Context, count int) error {
	return donUint64FieldMaxPendingQueueWriteSize.resolveCheckInt(ctx, d, count)
}

func (d *donSettingsResolver) secretIdentifierLimits(ctx context.Context) (vaultcap.SecretIdentifierLimits, error) {
	owner, err := donUint64FieldMaxIdentifierOwnerLengthBytes.resolveSize(ctx, d)
	if err != nil {
		return vaultcap.SecretIdentifierLimits{}, err
	}
	namespace, err := donUint64FieldMaxIdentifierNamespaceLengthBytes.resolveSize(ctx, d)
	if err != nil {
		return vaultcap.SecretIdentifierLimits{}, err
	}
	key, err := donUint64FieldMaxIdentifierKeyLengthBytes.resolveSize(ctx, d)
	if err != nil {
		return vaultcap.SecretIdentifierLimits{}, err
	}
	return vaultcap.SecretIdentifierLimits{
		MaxOwnerLength:     owner,
		MaxNamespaceLength: namespace,
		MaxKeyLength:       key,
	}, nil
}

func (d *donSettingsResolver) maxRequestBatchSize(ctx context.Context) int {
	return donUint64FieldMaxRequestBatchSize.resolveInt(ctx, d)
}
