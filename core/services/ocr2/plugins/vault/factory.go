package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/smdkg/dkgocr"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/dkgrecipientkey"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ReportingPluginConfig struct {
	LazyPublicKey *vaultcap.LazyPublicKey
	// Sourced from the DKG DB results package
	PublicKey       *tdh2easy.PublicKey
	PrivateKeyShare *tdh2easy.PrivateShare

	// Sourced from the offchain config
	MaxSecretsPerOwner                limits.BoundLimiter[int]
	MaxCiphertextLengthBytes          limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierKeyLengthBytes       limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierOwnerLengthBytes     limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierNamespaceLengthBytes limits.BoundLimiter[pkgconfig.Size]
	MaxShareLengthBytes               limits.BoundLimiter[pkgconfig.Size]
	MaxRequestBatchSize               limits.BoundLimiter[int]
	MaxBatchSize                      limits.BoundLimiter[int]
	MaxBlobPayloadBytes               limits.BoundLimiter[pkgconfig.Size]
	VaultForceEmptyOCRRounds          limits.GateLimiter
	VaultOptimizationsEnabled         limits.GateLimiter
}

func NewReportingPluginFactory(
	lggr logger.Logger,
	store *requests.Store[*vaulttypes.Request],
	db dkgocrtypes.ResultPackageDatabase,
	recipientKey *dkgrecipientkey.Key,
	lazyPublicKey *vaultcap.LazyPublicKey,
	limitsFactory limits.Factory,
	lifecycle *vaultcap.RequestLifecycleTracker,
) (*ReportingPluginFactory, error) {
	if db == nil {
		return nil, errors.New("result package db cannot be nil")
	}

	if recipientKey == nil {
		return nil, errors.New("DKG recipient key cannot be nil when using result package db")
	}

	if lifecycle == nil {
		return nil, errors.New("request lifecycle tracker cannot be nil")
	}

	cfg := &ReportingPluginConfig{
		LazyPublicKey: lazyPublicKey,
	}

	return &ReportingPluginFactory{
		lggr:          lggr.Named("VaultReportingPluginFactory"),
		store:         store,
		cfg:           cfg,
		db:            db,
		recipientKey:  recipientKey,
		limitsFactory: limitsFactory,
		lifecycle:     lifecycle,
	}, nil
}

type ReportingPluginFactory struct {
	lggr          logger.Logger
	store         *requests.Store[*vaulttypes.Request]
	cfg           *ReportingPluginConfig
	db            dkgocrtypes.ResultPackageDatabase
	recipientKey  *dkgrecipientkey.Key
	limitsFactory limits.Factory
	lifecycle     *vaultcap.RequestLifecycleTracker
}

func (r *ReportingPluginFactory) getKeyMaterial(ctx context.Context, instanceID string) (publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, err error) {
	pack, err := r.db.ReadResultPackage(ctx, dkgocrtypes.InstanceID(instanceID))
	if err != nil {
		return nil, nil, fmt.Errorf("could not read result package from db: %w", err)
	}
	if pack == nil {
		return nil, nil, fmt.Errorf("no result package found in db for instance ID %s", instanceID)
	}
	rP := dkgocr.NewResultPackage()
	err = rP.UnmarshalBinary(pack.ReportWithResultPackage)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal result package: %w", err)
	}

	tdh2PubKey, err := tdh2shim.TDH2PublicKeyFromDKGResult(rP)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get tdh2 public key from DKG result: %w", err)
	}
	publicKey, err = tdh2ToTDH2EasyPK(tdh2PubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not convert to tdh2easy public key: %w", err)
	}

	tdh2PrivateKeyShare, err := tdh2shim.TDH2PrivateShareFromDKGResult(rP, r.recipientKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get tdh2 private key share from DKG result: %w", err)
	}
	privateKeyShare, err = tdh2ToTDH2EasyKS(tdh2PrivateKeyShare)
	if err != nil {
		return nil, nil, fmt.Errorf("could not convert to tdh2easy private key share: %w", err)
	}

	return publicKey, privateKeyShare, nil
}

const dkgPollInterval = 2 * time.Second

// pollForKeyMaterial polls the DKG result package database until the key
// material for the given instance ID is available or the context is cancelled.
// This avoids returning an immediate error when the DKG protocol hasn't
// completed yet, which would trigger libocr's exponential backoff (up to 2
// minutes between retries). By polling here within the MaxDurationInitialization
// window, the vault oracle can start as soon as the DKG result is written.
func (r *ReportingPluginFactory) pollForKeyMaterial(ctx context.Context, instanceID string) (publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, err error) {
	for {
		publicKey, privateKeyShare, err = r.getKeyMaterial(ctx, instanceID)
		if err == nil {
			return publicKey, privateKeyShare, nil
		}

		r.lggr.Debugw("DKG result package not yet available, will retry", "instanceID", instanceID, "error", err)

		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("context cancelled while waiting for DKG key material (instanceID=%s): %w", instanceID, err)
		case <-time.After(dkgPollInterval):
		}
	}
}

func newReportingPluginConfigLimiters(factory limits.Factory) (*ReportingPluginConfig, error) {
	maxCiphertextLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultCiphertextSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultCiphertextSizeLimit: %w", err)
	}

	maxIdentifierKeyLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultIdentifierKeySizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultIdentifierKeySizeLimit: %w", err)
	}

	maxIdentifierOwnerLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultIdentifierOwnerSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultIdentifierOwnerSizeLimit: %w", err)
	}

	maxIdentifierNamespaceLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultIdentifierNamespaceSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultIdentifierNamespaceSizeLimit: %w", err)
	}

	maxShareLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultShareSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultShareSizeLimit: %w", err)
	}

	maxRequestBatchSizeLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultRequestBatchSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultRequestBatchSizeLimit: %w", err)
	}

	vaultForceEmptyOCRRounds, err := limits.MakeGateLimiter(factory, cresettings.Default.VaultForceEmptyOCRRounds)
	if err != nil {
		return nil, fmt.Errorf("VaultForceEmptyOCRRounds: %w", err)
	}

	vaultOptimizationsEnabled, err := limits.MakeGateLimiter(factory, cresettings.Default.VaultOptimizationsEnabled)
	if err != nil {
		return nil, fmt.Errorf("VaultOptimizationsEnabled: %w", err)
	}

	maxBlobPayloadBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultMaxBlobPayloadSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultMaxBlobPayloadSizeLimit: %w", err)
	}

	return &ReportingPluginConfig{
		MaxShareLengthBytes:               maxShareLengthBytesLimiter,
		MaxRequestBatchSize:               maxRequestBatchSizeLimiter,
		MaxCiphertextLengthBytes:          maxCiphertextLengthBytesLimiter,
		MaxIdentifierKeyLengthBytes:       maxIdentifierKeyLengthBytesLimiter,
		MaxIdentifierOwnerLengthBytes:     maxIdentifierOwnerLengthBytesLimiter,
		MaxIdentifierNamespaceLengthBytes: maxIdentifierNamespaceLengthBytesLimiter,
		MaxBlobPayloadBytes:               maxBlobPayloadBytesLimiter,
		VaultForceEmptyOCRRounds:          vaultForceEmptyOCRRounds,
		VaultOptimizationsEnabled:         vaultOptimizationsEnabled,
	}, nil
}

func logLimit[N limits.Number](ctx context.Context, lggr logger.Logger, limiter limits.BoundLimiter[N]) N {
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: "DUMMY-OWNER-FOR-LOGGING"})
	limit, err := limiter.Limit(ctx)
	if err != nil {
		lggr.Errorw("could not fetch limit", "error", err)
	}
	return limit
}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	var configProto vaultcommon.ReportingPluginConfig
	if err := proto.Unmarshal(config.OffchainConfig, &configProto); err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not unmarshal reporting plugin config: %w", err)
	}

	cfg, err := newReportingPluginConfigLimiters(r.limitsFactory)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create reporting plugin config limiters: %w", err)
	}

	maxSecretsPerOwnerLimit := cresettings.Default.PerOwner.VaultSecretsLimit
	if configProto.MaxSecretsPerOwner != 0 {
		maxSecretsPerOwnerLimit.DefaultValue = int(configProto.MaxSecretsPerOwner)
	}

	cfg.MaxSecretsPerOwner, err = limits.MakeUpperBoundLimiter(r.limitsFactory, maxSecretsPerOwnerLimit)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create max secrets per owner limiter: %w", err)
	}

	cfg.MaxBatchSize, err = limits.MakeUpperBoundLimiter(r.limitsFactory, cresettings.Default.VaultPluginBatchSizeLimit)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create max batch size limiter: %w", err)
	}

	if configProto.DKGInstanceID == nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, errors.New("DKG instance ID cannot be nil")
	}

	r.lggr.Debugw("fetching key material for instance id", "instanceID", *configProto.DKGInstanceID)
	publicKey, privateKeyShare, err := r.pollForKeyMaterial(ctx, *configProto.DKGInstanceID)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not get key material from DB: %w", err)
	}

	r.cfg.LazyPublicKey.Set(publicKey)

	cfg.PublicKey = publicKey
	cfg.PrivateKeyShare = privateKeyShare

	metrics, err := newPluginMetrics(config.ConfigDigest.String())
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create plugin metrics: %w", err)
	}

	pluginLimits, err := initializePluginLimits(ctx, r.limitsFactory)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not resolve plugin limits: %w", err)
	}

	r.lggr.Debugw("instantiating VaultReportingPlugin with config",
		"maxSecretsPerOwner", logLimit(ctx, r.lggr, cfg.MaxSecretsPerOwner),
		"maxCiphertextLengthBytes", logLimit(ctx, r.lggr, cfg.MaxCiphertextLengthBytes),
		"maxIdentifierKeyLengthBytes", logLimit(ctx, r.lggr, cfg.MaxIdentifierKeyLengthBytes),
		"maxIdentifierOwnerLengthBytes", logLimit(ctx, r.lggr, cfg.MaxIdentifierOwnerLengthBytes),
		"maxIdentifierNamespaceLengthBytes", logLimit(ctx, r.lggr, cfg.MaxIdentifierNamespaceLengthBytes),
		"maxRequestBatchSize", logLimit(ctx, r.lggr, cfg.MaxRequestBatchSize),
		"maxShareLengthBytes", logLimit(ctx, r.lggr, cfg.MaxShareLengthBytes),
		"batchSize", logLimit(ctx, r.lggr, cfg.MaxBatchSize),
		"maxBlobPayloadBytes", logLimit(ctx, r.lggr, cfg.MaxBlobPayloadBytes),
		"maxQueryBytes", pluginLimits.MaxQueryBytes,
		"maxObservationBytes", pluginLimits.MaxObservationBytes,
		"maxReportsPlusPrecursorBytes", pluginLimits.MaxReportsPlusPrecursorBytes,
		"maxReportBytes", pluginLimits.MaxReportBytes,
		"maxReportCount", pluginLimits.MaxReportCount,
		"maxKeyValueModifiedKeysPlusValuesBytes", pluginLimits.MaxKeyValueModifiedKeysPlusValuesBytes,
		"maxKeyValueModifiedKeys", pluginLimits.MaxKeyValueModifiedKeys,
		"reportingPluginLimitsMaxBlobPayloadBytes", pluginLimits.MaxBlobPayloadBytes,
		"maxPerOracleUnexpiredBlobCumulativePayloadBytes", pluginLimits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes,
		"maxPerOracleUnexpiredBlobCount", pluginLimits.MaxPerOracleUnexpiredBlobCount,
	)

	validator := vaultcap.NewRequestValidator(
		cfg.MaxRequestBatchSize,
		cfg.MaxCiphertextLengthBytes,
		cfg.MaxIdentifierKeyLengthBytes,
		cfg.MaxIdentifierOwnerLengthBytes,
		cfg.MaxIdentifierNamespaceLengthBytes,
	)

	r.lifecycle.SetConfigDigest(config.ConfigDigest.String())

	return &ReportingPlugin{
			lggr:                         r.lggr.Named("VaultReportingPlugin"),
			store:                        r.store,
			cfg:                          cfg,
			metrics:                      metrics,
			onchainCfg:                   config,
			validator:                    validator,
			lifecycle:                    r.lifecycle,
			maxObservationBytes:          pluginLimits.MaxObservationBytes,
			maxReportsPlusPrecursorBytes: pluginLimits.MaxReportsPlusPrecursorBytes,
			unmarshalBlob: func(data []byte) (ocr3_1types.BlobHandle, error) {
				handle := ocr3_1types.BlobHandle{}
				err := handle.UnmarshalBinary(data)
				return handle, err
			},
			marshalBlob: func(handle ocr3_1types.BlobHandle) ([]byte, error) {
				return handle.MarshalBinary()
			},
		}, ocr3_1types.ReportingPluginInfo1{
			Name:   "VaultReportingPlugin",
			Limits: pluginLimits,
		}, nil
}
