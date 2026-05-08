package vault

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// forceEmptyOCRRounds reports whether the VaultForceEmptyOCRRounds gate allows skipping pending-queue reads this round.
// When the gate allows (setting true), it logs a warning. When evaluation errors for reasons other than ErrorNotAllowed, it logs an error and returns false.
func forceEmptyOCRRounds(ctx context.Context, lggr logger.Logger, vaultForceEmptyOCRRounds limits.GateLimiter) bool {
	err := vaultForceEmptyOCRRounds.AllowErr(ctx)
	if err == nil {
		return true
	}
	if errors.Is(err, limits.ErrorNotAllowed{}) {
		return false
	}
	lggr.Errorw("unexpected error evaluating VaultForceEmptyOCRRounds gate; pending queue will be read normally", "error", err)
	return false
}

// resolveVaultOCRBoundLimitInt builds a short-lived BoundLimiter for an integer-sized CRE setting, reads Limit once, and closes the limiter.
func resolveVaultOCRBoundLimitInt[I constraints.Integer](
	ctx context.Context,
	factory limits.Factory,
	spec settings.IsSetting[I],
	settingKey string,
) (int, error) {
	limiter, err := limits.MakeUpperBoundLimiter(factory, spec)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", settingKey, err)
	}
	defer func() { _ = limiter.Close() }()

	v, err := limiter.Limit(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", settingKey, err)
	}
	return int(v), nil
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
