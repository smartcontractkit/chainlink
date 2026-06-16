package vault

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// gateAllows reports whether the given CRE gate allows the gated behavior.
// When evaluation errors for reasons other than ErrorNotAllowed, it logs an error and returns false.
func gateAllows(ctx context.Context, lggr logger.Logger, gate limits.GateLimiter, gateName string) bool {
	err := gate.AllowErr(ctx)
	if err == nil {
		return true
	}
	if errors.Is(err, limits.ErrorNotAllowed{}) {
		return false
	}
	lggr.Errorw("unexpected error evaluating CRE gate", "gate", gateName, "error", err)
	return false
}

// resolveVaultOCRBoundLimitInt builds a short-lived BoundLimiter for an integer-sized CRE setting, reads Limit once, and closes the limiter.
func resolveVaultOCRBoundLimitInt[I constraints.Integer](
	ctx context.Context,
	factory limits.Factory,
	spec settings.IsSetting[I],
	settingKey string,
) (int, error) {
	v, err := spec.GetSpec().GetOrDefault(ctx, factory.Settings)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", settingKey, err)
	}
	return int(v), nil
}

func validateEncryptedSharesEntry(es *vaultcommon.EncryptedShares) error {
	nBinary := len(es.BinaryShares)
	nString := len(es.Shares)
	if nBinary == 1 && nString == 0 {
		return nil
	}
	if nString == 1 && nBinary == 0 {
		return nil
	}
	return errors.New("observation must have exactly 1 share per encryption key")
}

func encryptedShareSizeForLimit(es *vaultcommon.EncryptedShares) (int, error) {
	if len(es.BinaryShares) == 1 {
		return len(es.BinaryShares[0]), nil
	}
	if len(es.Shares) == 1 {
		return len(es.Shares[0]), nil
	}
	return 0, errors.New("no share to measure")
}

func appendEncryptedShareEntry(dst, src *vaultcommon.EncryptedShares) {
	if len(src.BinaryShares) == 1 {
		dst.BinaryShares = append(dst.BinaryShares, src.BinaryShares[0])
		return
	}
	dst.Shares = append(dst.Shares, src.Shares[0])
}

func initializePluginLimits(ctx context.Context, limitsFactory limits.Factory) (ocr3_1types.ReportingPluginLimits, error) {
	maxQueryBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxQuerySizeLimit, "VaultMaxQuerySizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxObservationBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxObservationSizeLimit, "VaultMaxObservationSizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxReportsPlusPrecursorBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxReportsPlusPrecursorSizeLimit, "VaultMaxReportsPlusPrecursorSizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxReportBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxReportSizeLimit, "VaultMaxReportSizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxReportCount, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxReportCount, "VaultMaxReportCount")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxKVModifiedKeysPlusValuesBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxKeyValueModifiedKeysPlusValuesSizeLimit, "VaultMaxKeyValueModifiedKeysPlusValuesSizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxKVModifiedKeys, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxKeyValueModifiedKeys, "VaultMaxKeyValueModifiedKeys")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxBlobPayloadBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxBlobPayloadSizeLimit, "VaultMaxBlobPayloadSizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxPerOracleUnexpiredBlobCumulativePayloadBytes, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxPerOracleUnexpiredBlobCumulativePayloadSizeLimit, "VaultMaxPerOracleUnexpiredBlobCumulativePayloadSizeLimit")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}
	maxPerOracleUnexpiredBlobCount, err := resolveVaultOCRBoundLimitInt(ctx, limitsFactory, cresettings.Default.VaultMaxPerOracleUnexpiredBlobCount, "VaultMaxPerOracleUnexpiredBlobCount")
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, err
	}

	return ocr3_1types.ReportingPluginLimits{
		MaxQueryBytes:                                   maxQueryBytes,
		MaxObservationBytes:                             maxObservationBytes,
		MaxReportsPlusPrecursorBytes:                    maxReportsPlusPrecursorBytes,
		MaxReportBytes:                                  maxReportBytes,
		MaxReportCount:                                  maxReportCount,
		MaxKeyValueModifiedKeysPlusValuesBytes:          maxKVModifiedKeysPlusValuesBytes,
		MaxKeyValueModifiedKeys:                         maxKVModifiedKeys,
		MaxBlobPayloadBytes:                             maxBlobPayloadBytes,
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: maxPerOracleUnexpiredBlobCumulativePayloadBytes,
		MaxPerOracleUnexpiredBlobCount:                  maxPerOracleUnexpiredBlobCount,
	}, nil
}
