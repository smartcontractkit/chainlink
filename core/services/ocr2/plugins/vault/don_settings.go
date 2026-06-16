package vault

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type donSettingsResolver struct {
	r              *ReportingPlugin
	ctx            context.Context
	store          ReadKVStore
	consensus      bool
	initialized    bool
	initErr        error
	cachedSettings *vaultcommon.NodeSettings
}

func (r *ReportingPlugin) newDonSettingsResolver(ctx context.Context, store ReadKVStore) *donSettingsResolver {
	return &donSettingsResolver{
		r:         r,
		ctx:       ctx,
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

func (r *ReportingPlugin) ensureActiveSettingsForRound(ctx context.Context, seqNr uint64, store ReadKVStore) error {
	if r.activeSettings != nil && r.activeSettingsSeqNr == seqNr {
		return nil
	}
	resolver := r.newDonSettingsResolver(ctx, store)
	if err := resolver.init(); err != nil {
		return err
	}
	r.activeSettings = resolver
	r.activeSettingsSeqNr = seqNr
	return nil
}

func (d *donSettingsResolver) init() error {
	if d.initialized {
		return d.initErr
	}
	d.initialized = true

	if !d.consensus {
		return nil
	}

	settings, err := d.store.GetDONSettings(d.ctx)
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

func (d *donSettingsResolver) donSettings() *vaultcommon.NodeSettings {
	return d.cachedSettings
}

func (r *ReportingPlugin) activeDonSettings() *vaultcommon.NodeSettings {
	if r.activeSettings == nil {
		return nil
	}
	return r.activeSettings.donSettings()
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
	allFieldsQuorum = applyDonSettingsQuorum(r.lggr, donSettingsUint64BaseFields(), merged, marshalledObs, threshold) && allFieldsQuorum

	if existing == nil {
		// Initial seeding of DON settings, reject if not all fields reached quorum
		if !allFieldsQuorum {
			r.lggr.Debugw("DON settings initial seed deferred: not all fields reached quorum", "threshold", threshold)
			return nil, nil
		}
		if err := store.WriteDONSettings(ctx, merged); err != nil {
			return nil, fmt.Errorf("failed to write DON settings: %w", err)
		}
		return merged, nil
	}

	if !proto.Equal(existing, merged) {
		if err := store.WriteDONSettings(ctx, merged); err != nil {
			return nil, fmt.Errorf("failed to write DON settings: %w", err)
		}
	}
	return merged, nil
}

// Wrappers for readability at call sites
func (d *donSettingsResolver) optimizationsEnabled() bool {
	return donSettingsBoolFieldByName(donSettingVaultOptimizationsEnabledName).resolve(d)
}

func (d *donSettingsResolver) forceEmptyOCRRounds() bool {
	return donSettingsBoolFieldByName(donSettingVaultForceEmptyOcrRoundsName).resolve(d)
}

func (d *donSettingsResolver) maxBlobPayloadBytes() (pkgconfig.Size, error) {
	return donSettingsUint64FieldByName(donSettingMaxBlobPayloadBytesName).resolveSize(d)
}

func (d *donSettingsResolver) checkMaxPendingQueueWriteSize(count int) error {
	return donSettingsUint64FieldByName(donSettingMaxPendingQueueWriteSizeName).resolveCheckInt(d, count)
}
