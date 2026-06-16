package vault

import (
	"context"
	"fmt"
	"math"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

type donSettingsField[T comparable] struct {
	name  string
	local func(ctx context.Context, r *ReportingPlugin) T
	get   func(*vaultcommon.NodeSettings) T
	set   func(*vaultcommon.NodeSettings, T)
}

func (f donSettingsField[T]) resolve(ctx context.Context, d *donSettingsResolver) T {
	if s := d.donSettings(); s != nil {
		return f.get(s)
	}
	return f.local(ctx, d.r)
}

type donSettingsUint64Field struct {
	base      donSettingsField[uint64]
	limitSize func(ctx context.Context, r *ReportingPlugin) (pkgconfig.Size, error)
	checkInt  func(ctx context.Context, r *ReportingPlugin, count int) error
}

func (f donSettingsUint64Field) resolveSize(ctx context.Context, d *donSettingsResolver) (pkgconfig.Size, error) {
	if s := d.donSettings(); s != nil {
		return nodeSettingsUint64AsSize(f.base.get(s)), nil
	}
	if f.limitSize != nil {
		return f.limitSize(ctx, d.r)
	}
	return nodeSettingsUint64AsSize(f.base.local(ctx, d.r)), nil
}

func (f donSettingsUint64Field) resolveInt(ctx context.Context, d *donSettingsResolver) int {
	if s := d.donSettings(); s != nil {
		return nodeSettingsUint64AsInt(f.base.get(s))
	}
	return nodeSettingsUint64AsInt(f.base.local(ctx, d.r))
}

func (f donSettingsUint64Field) resolveCheckInt(ctx context.Context, d *donSettingsResolver, count int) error {
	if s := d.donSettings(); s != nil {
		limit := nodeSettingsUint64AsInt(f.base.get(s))
		if count > limit {
			return limits.ErrorBoundLimited[int]{Limit: limit, Amount: count}
		}
		return nil
	}
	if f.checkInt != nil {
		return f.checkInt(ctx, d.r, count)
	}
	limit := nodeSettingsUint64AsInt(f.base.local(ctx, d.r))
	if count > limit {
		return limits.ErrorBoundLimited[int]{Limit: limit, Amount: count}
	}
	return nil
}

func nodeSettingsUint64AsInt(v uint64) int {
	if v > uint64(math.MaxInt) {
		return math.MaxInt
	}
	return int(v)
}

func nodeSettingsUint64AsSize(v uint64) pkgconfig.Size {
	return pkgconfig.Size(nodeSettingsUint64AsInt(v))
}

func limitValueAsUint64(v int) uint64 {
	if v < 0 {
		return 0
	}
	return uint64(v)
}

func localGate(ctx context.Context, r *ReportingPlugin, gate limits.GateLimiter, gateName string) bool {
	return gateAllows(ctx, r.lggr, gate, gateName)
}

func localSizeLimit(ctx context.Context, limiter limits.BoundLimiter[pkgconfig.Size]) uint64 {
	v, _ := limiter.Limit(ctx)
	return limitValueAsUint64(int(v))
}

func localIntLimit(ctx context.Context, limiter limits.BoundLimiter[int]) uint64 {
	v, _ := limiter.Limit(ctx)
	return limitValueAsUint64(v)
}

var donFieldVaultOptimizationsEnabled = donSettingsField[bool]{
	name: "vault_optimizations_enabled",
	local: func(ctx context.Context, r *ReportingPlugin) bool {
		return localGate(ctx, r, r.cfg.VaultOptimizationsEnabled, "VaultOptimizationsEnabled")
	},
	get: func(s *vaultcommon.NodeSettings) bool { return s.VaultOptimizationsEnabled },
	set: func(s *vaultcommon.NodeSettings, v bool) { s.VaultOptimizationsEnabled = v },
}

var donFieldVaultForceEmptyOcrRounds = donSettingsField[bool]{
	name: "vault_force_empty_ocr_rounds",
	local: func(ctx context.Context, r *ReportingPlugin) bool {
		return localGate(ctx, r, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds")
	},
	get: func(s *vaultcommon.NodeSettings) bool { return s.VaultForceEmptyOcrRounds },
	set: func(s *vaultcommon.NodeSettings, v bool) { s.VaultForceEmptyOcrRounds = v },
}

var donSettingsBoolFields = []donSettingsField[bool]{
	donFieldVaultOptimizationsEnabled,
	donFieldVaultForceEmptyOcrRounds,
}

var donFieldMaxIdentifierKeyLengthBytes = donSettingsField[uint64]{
	name: "max_identifier_key_length_bytes",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localSizeLimit(ctx, r.cfg.MaxIdentifierKeyLengthBytes)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxIdentifierKeyLengthBytes },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxIdentifierKeyLengthBytes = v },
}

var donFieldMaxIdentifierOwnerLengthBytes = donSettingsField[uint64]{
	name: "max_identifier_owner_length_bytes",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localSizeLimit(ctx, r.cfg.MaxIdentifierOwnerLengthBytes)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxIdentifierOwnerLengthBytes },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxIdentifierOwnerLengthBytes = v },
}

var donFieldMaxIdentifierNamespaceLengthBytes = donSettingsField[uint64]{
	name: "max_identifier_namespace_length_bytes",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localSizeLimit(ctx, r.cfg.MaxIdentifierNamespaceLengthBytes)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxIdentifierNamespaceLengthBytes },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxIdentifierNamespaceLengthBytes = v },
}

var donFieldMaxShareLengthBytes = donSettingsField[uint64]{
	name: "max_share_length_bytes",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localSizeLimit(ctx, r.cfg.MaxShareLengthBytes)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxShareLengthBytes },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxShareLengthBytes = v },
}

var donFieldMaxBlobPayloadBytes = donSettingsField[uint64]{
	name: "max_blob_payload_bytes",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localSizeLimit(ctx, r.cfg.MaxBlobPayloadBytes)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxBlobPayloadBytes },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxBlobPayloadBytes = v },
}

var donFieldMaxPendingQueueWriteSize = donSettingsField[uint64]{
	name: "max_pending_queue_write_size",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localIntLimit(ctx, r.cfg.MaxPendingQueueWriteSize)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxPendingQueueWriteSize },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxPendingQueueWriteSize = v },
}

var donFieldMaxRequestBatchSize = donSettingsField[uint64]{
	name: "max_request_batch_size",
	local: func(ctx context.Context, r *ReportingPlugin) uint64 {
		return localIntLimit(ctx, r.cfg.MaxRequestBatchSize)
	},
	get: func(s *vaultcommon.NodeSettings) uint64 { return s.MaxRequestBatchSize },
	set: func(s *vaultcommon.NodeSettings, v uint64) { s.MaxRequestBatchSize = v },
}

var donUint64FieldMaxIdentifierKeyLengthBytes = donSettingsUint64Field{base: donFieldMaxIdentifierKeyLengthBytes}

var donUint64FieldMaxIdentifierOwnerLengthBytes = donSettingsUint64Field{base: donFieldMaxIdentifierOwnerLengthBytes}

var donUint64FieldMaxIdentifierNamespaceLengthBytes = donSettingsUint64Field{base: donFieldMaxIdentifierNamespaceLengthBytes}

var donUint64FieldMaxShareLengthBytes = donSettingsUint64Field{base: donFieldMaxShareLengthBytes}

var donUint64FieldMaxBlobPayloadBytes = donSettingsUint64Field{
	base: donFieldMaxBlobPayloadBytes,
	limitSize: func(ctx context.Context, r *ReportingPlugin) (pkgconfig.Size, error) {
		return r.cfg.MaxBlobPayloadBytes.Limit(ctx)
	},
}

var donUint64FieldMaxPendingQueueWriteSize = donSettingsUint64Field{
	base: donFieldMaxPendingQueueWriteSize,
	checkInt: func(ctx context.Context, r *ReportingPlugin, count int) error {
		return r.cfg.MaxPendingQueueWriteSize.Check(ctx, count)
	},
}

var donUint64FieldMaxRequestBatchSize = donSettingsUint64Field{base: donFieldMaxRequestBatchSize}

var donSettingsUint64Fields = []donSettingsUint64Field{
	donUint64FieldMaxIdentifierKeyLengthBytes,
	donUint64FieldMaxIdentifierOwnerLengthBytes,
	donUint64FieldMaxIdentifierNamespaceLengthBytes,
	donUint64FieldMaxShareLengthBytes,
	donUint64FieldMaxBlobPayloadBytes,
	donUint64FieldMaxPendingQueueWriteSize,
	donUint64FieldMaxRequestBatchSize,
}

var donSettingsUint64BaseFields = []donSettingsField[uint64]{
	donFieldMaxIdentifierKeyLengthBytes,
	donFieldMaxIdentifierOwnerLengthBytes,
	donFieldMaxIdentifierNamespaceLengthBytes,
	donFieldMaxShareLengthBytes,
	donFieldMaxBlobPayloadBytes,
	donFieldMaxPendingQueueWriteSize,
	donFieldMaxRequestBatchSize,
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
