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

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
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

func secretRequestForID(req *vaultcommon.GetSecretsRequest, id *vaultcommon.SecretIdentifier) (*vaultcommon.SecretRequest, error) {
	if req == nil {
		return nil, errors.New("GetSecrets request is nil")
	}
	if id == nil {
		return nil, errors.New("secret identifier is nil")
	}
	key := vaulttypes.KeyFor(id)
	for _, secretReq := range req.Requests {
		if vaulttypes.KeyFor(secretReq.Id) == key {
			return secretReq, nil
		}
	}
	return nil, fmt.Errorf("no secret request for id %s", key)
}

func validateGetSecretsShareLabels(secretReq *vaultcommon.SecretRequest, data *vaultcommon.SecretData) error {
	if len(data.EncryptedDecryptionKeyShares) != len(secretReq.EncryptionKeys) {
		return fmt.Errorf("expected %d encrypted share entries, got %d", len(secretReq.EncryptionKeys), len(data.EncryptedDecryptionKeyShares))
	}

	allowed := make(map[string]struct{}, len(secretReq.EncryptionKeys))
	for _, pk := range secretReq.EncryptionKeys {
		allowed[pk] = struct{}{}
	}

	seen := make(map[string]bool, len(secretReq.EncryptionKeys))
	for _, es := range data.EncryptedDecryptionKeyShares {
		if es == nil {
			return errors.New("nil encrypted share entry")
		}
		if err := validateEncryptedSharesEntry(es); err != nil {
			return err
		}
		if _, ok := allowed[es.EncryptionKey]; !ok {
			return fmt.Errorf("unexpected encryption key in response: %s", es.EncryptionKey)
		}
		if seen[es.EncryptionKey] {
			return fmt.Errorf("duplicate encryption key in response: %s", es.EncryptionKey)
		}
		seen[es.EncryptionKey] = true
	}

	for _, pk := range secretReq.EncryptionKeys {
		if !seen[pk] {
			return fmt.Errorf("missing encryption key in response: %s", pk)
		}
	}

	return nil
}

func buildPendingGetSecretsByID(items []*vaultcommon.StoredPendingQueueItem) (map[string]*vaultcommon.GetSecretsRequest, error) {
	out := make(map[string]*vaultcommon.GetSecretsRequest, len(items))
	for _, item := range items {
		payload, err := item.Item.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("pending queue item %s: %w", item.Id, err)
		}
		req, ok := payload.(*vaultcommon.GetSecretsRequest)
		if !ok {
			continue
		}
		out[item.Id] = req
	}
	return out, nil
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

func (r *ReportingPlugin) roundLggr(seqNr uint64) logger.Logger {
	return r.lggr.With("seqNr", seqNr)
}

func (r *ReportingPlugin) requestLggr(seqNr uint64, requestID string) logger.Logger {
	return r.roundLggr(seqNr).With("requestID", requestID)
}

func (r *ReportingPlugin) typedRequestLggr(seqNr uint64, requestID, requestType string) logger.Logger {
	return r.requestLggr(seqNr, requestID).With("requestType", requestType)
}
