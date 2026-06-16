package vault

import (
	"context"
	"fmt"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

const (
	donSettingVaultOptimizationsEnabledName   = "vault_optimizations_enabled"
	donSettingVaultForceEmptyOcrRoundsName   = "vault_force_empty_ocr_rounds"
	donSettingMaxBlobPayloadBytesName        = "max_blob_payload_bytes"
	donSettingMaxPendingQueueWriteSizeName   = "max_pending_queue_write_size"
)

type donSettingsField[T comparable] struct {
	name  string
	local func(ctx context.Context, r *ReportingPlugin) T
	get   func(*vaultcommon.NodeSettings) T
	set   func(*vaultcommon.NodeSettings, T)
}

func (f donSettingsField[T]) resolve(d *donSettingsResolver) T {
	if s := d.donSettings(); s != nil {
		return f.get(s)
	}
	return f.local(d.ctx, d.r)
}

type donSettingsUint64Field struct {
	base      donSettingsField[uint64]
	limitSize func(ctx context.Context, r *ReportingPlugin) (pkgconfig.Size, error)
	checkInt  func(ctx context.Context, r *ReportingPlugin, count int) error
}

func (f donSettingsUint64Field) resolveSize(d *donSettingsResolver) (pkgconfig.Size, error) {
	if s := d.donSettings(); s != nil {
		return pkgconfig.Size(f.base.get(s)), nil
	}
	if f.limitSize != nil {
		return f.limitSize(d.ctx, d.r)
	}
	return pkgconfig.Size(f.base.local(d.ctx, d.r)), nil
}

func (f donSettingsUint64Field) resolveCheckInt(d *donSettingsResolver, count int) error {
	if s := d.donSettings(); s != nil {
		limit := int(f.base.get(s))
		if count > limit {
			return limits.ErrorBoundLimited[int]{Limit: limit, Amount: count}
		}
		return nil
	}
	if f.checkInt != nil {
		return f.checkInt(d.ctx, d.r, count)
	}
	limit := int(f.base.local(d.ctx, d.r))
	if count > limit {
		return limits.ErrorBoundLimited[int]{Limit: limit, Amount: count}
	}
	return nil
}

func localGate(ctx context.Context, r *ReportingPlugin, gate limits.GateLimiter, gateName string) bool {
	return gateAllows(ctx, r.lggr, gate, gateName)
}

func localSizeLimit(ctx context.Context, limiter limits.BoundLimiter[pkgconfig.Size]) uint64 {
	v, _ := limiter.Limit(ctx)
	return uint64(v)
}

func localIntLimit(ctx context.Context, limiter limits.BoundLimiter[int]) uint64 {
	v, _ := limiter.Limit(ctx)
	return uint64(v)
}

var donSettingsBoolFields = []donSettingsField[bool]{
	{
		name: donSettingVaultOptimizationsEnabledName,
		local: func(ctx context.Context, r *ReportingPlugin) bool {
			return localGate(ctx, r, r.cfg.VaultOptimizationsEnabled, "VaultOptimizationsEnabled")
		},
		get: func(s *vaultcommon.NodeSettings) bool { return s.VaultOptimizationsEnabled },
		set: func(s *vaultcommon.NodeSettings, v bool) { s.VaultOptimizationsEnabled = v },
	},
	{
		name: donSettingVaultForceEmptyOcrRoundsName,
		local: func(ctx context.Context, r *ReportingPlugin) bool {
			return localGate(ctx, r, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds")
		},
		get: func(s *vaultcommon.NodeSettings) bool { return s.VaultForceEmptyOcrRounds },
		set: func(s *vaultcommon.NodeSettings, v bool) { s.VaultForceEmptyOcrRounds = v },
	},
}

var donSettingsUint64Fields = []donSettingsUint64Field{
	{
		base: donSettingsField[uint64]{
			name: "max_ciphertext_length_bytes",
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localSizeLimit(ctx, r.cfg.MaxCiphertextLengthBytes)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxCiphertextLengthBytes },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxCiphertextLengthBytes = v },
		},
	},
	{
		base: donSettingsField[uint64]{
			name: "max_identifier_key_length_bytes",
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localSizeLimit(ctx, r.cfg.MaxIdentifierKeyLengthBytes)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxIdentifierKeyLengthBytes },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxIdentifierKeyLengthBytes = v },
		},
	},
	{
		base: donSettingsField[uint64]{
			name: "max_identifier_owner_length_bytes",
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localSizeLimit(ctx, r.cfg.MaxIdentifierOwnerLengthBytes)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxIdentifierOwnerLengthBytes },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxIdentifierOwnerLengthBytes = v },
		},
	},
	{
		base: donSettingsField[uint64]{
			name: "max_identifier_namespace_length_bytes",
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localSizeLimit(ctx, r.cfg.MaxIdentifierNamespaceLengthBytes)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxIdentifierNamespaceLengthBytes },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxIdentifierNamespaceLengthBytes = v },
		},
	},
	{
		base: donSettingsField[uint64]{
			name: "max_share_length_bytes",
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localSizeLimit(ctx, r.cfg.MaxShareLengthBytes)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxShareLengthBytes },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxShareLengthBytes = v },
		},
	},
	{
		base: donSettingsField[uint64]{
			name: donSettingMaxBlobPayloadBytesName,
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localSizeLimit(ctx, r.cfg.MaxBlobPayloadBytes)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxBlobPayloadBytes },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxBlobPayloadBytes = v },
		},
		limitSize: func(ctx context.Context, r *ReportingPlugin) (pkgconfig.Size, error) {
			return r.cfg.MaxBlobPayloadBytes.Limit(ctx)
		},
	},
	{
		base: donSettingsField[uint64]{
			name: donSettingMaxPendingQueueWriteSizeName,
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localIntLimit(ctx, r.cfg.MaxPendingQueueWriteSize)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxPendingQueueWriteSize },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxPendingQueueWriteSize = v },
		},
		checkInt: func(ctx context.Context, r *ReportingPlugin, count int) error {
			return r.cfg.MaxPendingQueueWriteSize.Check(ctx, count)
		},
	},
	{
		base: donSettingsField[uint64]{
			name: "max_request_batch_size",
			local: func(ctx context.Context, r *ReportingPlugin) uint64 {
				return localIntLimit(ctx, r.cfg.MaxRequestBatchSize)
			},
			get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxRequestBatchSize },
			set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxRequestBatchSize = v },
		},
	},
}

var (
	donSettingsBoolFieldIndex   map[string]int
	donSettingsUint64FieldIndex map[string]int
)

func init() {
	donSettingsBoolFieldIndex = make(map[string]int, len(donSettingsBoolFields))
	for i, field := range donSettingsBoolFields {
		donSettingsBoolFieldIndex[field.name] = i
	}
	donSettingsUint64FieldIndex = make(map[string]int, len(donSettingsUint64Fields))
	for i, field := range donSettingsUint64Fields {
		donSettingsUint64FieldIndex[field.base.name] = i
	}
}

func donSettingsBoolFieldByName(name string) donSettingsField[bool] {
	i, ok := donSettingsBoolFieldIndex[name]
	if !ok {
		panic("unknown DON settings bool field: " + name)
	}
	return donSettingsBoolFields[i]
}

func donSettingsUint64FieldByName(name string) donSettingsUint64Field {
	i, ok := donSettingsUint64FieldIndex[name]
	if !ok {
		panic("unknown DON settings uint64 field: " + name)
	}
	return donSettingsUint64Fields[i]
}

func donSettingsUint64BaseFields() []donSettingsField[uint64] {
	bases := make([]donSettingsField[uint64], len(donSettingsUint64Fields))
	for i, field := range donSettingsUint64Fields {
		bases[i] = field.base
	}
	return bases
}

func populateLocalNodeSettings(ctx context.Context, r *ReportingPlugin) *vaultcommon.NodeSettings {
	settings := &vaultcommon.NodeSettings{}
	for _, field := range donSettingsBoolFields {
		field.set(settings, field.local(ctx, r))
	}
	for _, field := range donSettingsUint64Fields {
		field.base.set(settings, field.base.local(ctx, r))
	}
	return settings
}

func validateObservedNodeSettings(settings *vaultcommon.NodeSettings) error {
	for _, field := range donSettingsUint64Fields {
		if field.base.get(settings) == 0 {
			return fmt.Errorf("invalid observation: %s must be positive", field.base.name)
		}
	}
	return nil
}

func quorumValue[T comparable](counts map[T]int, threshold int) (T, bool) {
	for v, c := range counts {
		if c >= threshold {
			return v, true
		}
	}
	var zero T
	return zero, false
}

func applyDonSettingsQuorum[T comparable](
	lggr logger.Logger,
	fields []donSettingsField[T],
	merged *vaultcommon.NodeSettings,
	marshalledObs map[uint8]*vaultcommon.Observations,
	threshold int,
) bool {
	allFieldsQuorum := true
	for _, field := range fields {
		counts := map[T]int{}
		for _, obs := range marshalledObs {
			if obs.NodeSettings == nil {
				continue
			}
			counts[field.get(obs.NodeSettings)]++
		}
		if v, ok := quorumValue(counts, threshold); ok {
			field.set(merged, v)
			lggr.Debugw("DON settings field quorum reached", "field", field.name, "value", v, "counts", counts, "threshold", threshold)
		} else {
			allFieldsQuorum = false
			if len(counts) > 0 {
				lggr.Debugw("DON settings field quorum not reached", "field", field.name, "counts", counts, "threshold", threshold)
			}
		}
	}
	return allFieldsQuorum
}
