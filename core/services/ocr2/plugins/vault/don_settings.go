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

func (r *ReportingPlugin) ensureActiveSettingsForRound(ctx context.Context, seqNr uint64, store ReadKVStore) error {
	if r.activeSettings != nil && r.activeSettingsSeqNr == seqNr {
		return nil
	}
	resolver := r.newDonSettingsResolver(ctx, store)
	if err := resolver.init(ctx); err != nil {
		return err
	}
	r.activeSettings = resolver
	r.activeSettingsSeqNr = seqNr
	return nil
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

func (d *donSettingsResolver) donSettings() *vaultcommon.NodeSettings {
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
func (d *donSettingsResolver) optimizationsEnabled(ctx context.Context) bool {
	return donFieldVaultOptimizationsEnabled.resolve(ctx, d)
}

func (d *donSettingsResolver) forceEmptyOCRRounds(ctx context.Context) bool {
	return donFieldVaultForceEmptyOcrRounds.resolve(ctx, d)
}

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
