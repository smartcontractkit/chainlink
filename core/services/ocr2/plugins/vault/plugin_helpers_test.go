package vault

import (
	"context"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type testPluginOption func(*testPluginBuildOpts)

type testPluginBuildOpts struct {
	lggr                                 logger.Logger
	store                                *requests.Store[*vaulttypes.Request]
	publicKey                            *tdh2easy.PublicKey
	privateKeyShare                      *tdh2easy.PrivateShare
	onchainCfg                           ocr3types.ReportingPluginConfig
	maxSecretsPerOwner                   int
	maxCiphertextLengthBytes             int
	maxIdentifierOwnerLengthBytes        int
	maxIdentifierNamespaceLengthBytes    int
	maxIdentifierKeyLengthBytes          int
	maxRequestBatchSize                  int
	batchSize                            int
	maxBlobPayloadBytes                  int
	vaultOptimizationsEnabled            bool
	marshalBlob                          func(ocr3_1types.BlobHandle) ([]byte, error)
	unmarshalBlob                        func([]byte) (ocr3_1types.BlobHandle, error)
	maxObservationBytesOverride          int
	maxReportsPlusPrecursorBytesOverride int
}

func withLggr(lggr logger.Logger) testPluginOption {
	return func(o *testPluginBuildOpts) { o.lggr = lggr }
}

func withStore(store *requests.Store[*vaulttypes.Request]) testPluginOption {
	return func(o *testPluginBuildOpts) { o.store = store }
}

func withKeys(pk *tdh2easy.PublicKey, share *tdh2easy.PrivateShare) testPluginOption {
	return func(o *testPluginBuildOpts) {
		o.publicKey = pk
		o.privateKeyShare = share
	}
}

func withMaxCiphertextLengthBytes(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.maxCiphertextLengthBytes = n }
}

func withMaxIdentifierLengths(owner, namespace, key int) testPluginOption {
	return func(o *testPluginBuildOpts) {
		o.maxIdentifierOwnerLengthBytes = owner
		o.maxIdentifierNamespaceLengthBytes = namespace
		o.maxIdentifierKeyLengthBytes = key
	}
}

func withMaxSecretsPerOwner(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.maxSecretsPerOwner = n }
}

func withVaultOptimizationsEnabled() testPluginOption {
	return func(o *testPluginBuildOpts) { o.vaultOptimizationsEnabled = true }
}

func withOnchainCfg(n int, f int) testPluginOption {
	return func(o *testPluginBuildOpts) {
		o.onchainCfg = ocr3types.ReportingPluginConfig{N: n, F: f}
	}
}

func withBatchSize(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.batchSize = n }
}

func withMaxRequestBatchSize(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.maxRequestBatchSize = n }
}

func withMaxBlobPayloadBytes(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.maxBlobPayloadBytes = n }
}

func withMarshalBlob(fn func(ocr3_1types.BlobHandle) ([]byte, error)) testPluginOption {
	return func(o *testPluginBuildOpts) { o.marshalBlob = fn }
}

func withMaxObservationBytes(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.maxObservationBytesOverride = n }
}

func withMaxReportsPlusPrecursorBytes(n int) testPluginOption {
	return func(o *testPluginBuildOpts) { o.maxReportsPlusPrecursorBytesOverride = n }
}

func newTestReportingPlugin(t *testing.T, opts ...testPluginOption) *ReportingPlugin {
	t.Helper()
	o := testPluginBuildOpts{
		lggr:                              logger.TestLogger(t),
		store:                             requests.NewStore[*vaulttypes.Request](),
		onchainCfg:                        ocr3types.ReportingPluginConfig{N: 0, F: 0},
		maxSecretsPerOwner:                1,
		maxCiphertextLengthBytes:          1024,
		maxIdentifierOwnerLengthBytes:     100,
		maxIdentifierNamespaceLengthBytes: 100,
		maxIdentifierKeyLengthBytes:       100,
		maxRequestBatchSize:               10,
		batchSize:                         10,
		marshalBlob:                       mockMarshalBlob,
		unmarshalBlob:                     mockUnmarshalBlob,
	}
	for _, opt := range opts {
		opt(&o)
	}
	cfg := makeReportingPluginConfig(t, o.batchSize, o.publicKey, o.privateKeyShare,
		o.maxSecretsPerOwner, o.maxCiphertextLengthBytes,
		o.maxIdentifierOwnerLengthBytes, o.maxIdentifierNamespaceLengthBytes,
		o.maxIdentifierKeyLengthBytes, o.maxRequestBatchSize, o.maxBlobPayloadBytes)
	if o.vaultOptimizationsEnabled {
		cfg.VaultOptimizationsEnabled = limits.NewGateLimiter(true)
	}
	ctx := context.Background()
	pl, err := initializePluginLimits(ctx, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)
	maxObs := pl.MaxObservationBytes
	maxPrec := pl.MaxReportsPlusPrecursorBytes
	if o.maxObservationBytesOverride > 0 {
		maxObs = o.maxObservationBytesOverride
	}
	if o.maxReportsPlusPrecursorBytesOverride > 0 {
		maxPrec = o.maxReportsPlusPrecursorBytesOverride
	}
	lc, err := vaultcap.NewRequestLifecycleTracker(o.lggr)
	require.NoError(t, err)
	return &ReportingPlugin{
		lggr:                         o.lggr,
		store:                        o.store,
		metrics:                      newTestMetrics(t),
		cfg:                          cfg,
		onchainCfg:                   o.onchainCfg,
		validator:                    makeTestValidator(cfg),
		lifecycle:                    lc,
		marshalBlob:                  o.marshalBlob,
		unmarshalBlob:                o.unmarshalBlob,
		maxObservationBytes:          maxObs,
		maxReportsPlusPrecursorBytes: maxPrec,
	}
}

func makeTestValidator(cfg *ReportingPluginConfig) *vaultcap.RequestValidator {
	return vaultcap.NewRequestValidator(
		cfg.MaxRequestBatchSize,
		cfg.MaxCiphertextLengthBytes,
		cfg.MaxIdentifierKeyLengthBytes,
		cfg.MaxIdentifierOwnerLengthBytes,
		cfg.MaxIdentifierNamespaceLengthBytes,
	)
}

func makeReportingPluginConfig(
	t *testing.T,
	batchSize int,
	publicKey *tdh2easy.PublicKey,
	privateKeyShare *tdh2easy.PrivateShare,
	maxSecretsPerOwner int,
	maxCipherTextLengthBytes int,
	maxIdentifierOwnerLengthBytes int,
	maxIdentifierNamespaceOwnerLengthBytes int,
	maxIdentifierKeyLengthBytes int,
	maxRequestBatchSize int,
	maxBlobPayloadBytes int,
) *ReportingPluginConfig {
	msl, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Int(maxSecretsPerOwner))
	require.NoError(t, err)

	cipherTextLimiter, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Size(pkgconfig.Size(maxCipherTextLengthBytes)*pkgconfig.Byte))
	require.NoError(t, err)

	shareLimiter, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, cresettings.Default.VaultShareSizeLimit)
	require.NoError(t, err)

	ownerLimiter, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Size(pkgconfig.Size(maxIdentifierOwnerLengthBytes)*pkgconfig.Byte))
	require.NoError(t, err)

	namespaceOwnerLimiter, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Size(pkgconfig.Size(maxIdentifierNamespaceOwnerLengthBytes)*pkgconfig.Byte))
	require.NoError(t, err)

	keyLimiter, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Size(pkgconfig.Size(maxIdentifierKeyLengthBytes)*pkgconfig.Byte))
	require.NoError(t, err)

	bsl, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Int(batchSize))
	require.NoError(t, err)

	requestBatchSizeLimiter, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Int(maxRequestBatchSize))
	require.NoError(t, err)

	var maxBlobPayloadLimiter limits.BoundLimiter[pkgconfig.Size]
	if maxBlobPayloadBytes > 0 {
		maxBlobPayloadLimiter, err = limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, settings.Size(pkgconfig.Size(maxBlobPayloadBytes)*pkgconfig.Byte))
	} else {
		maxBlobPayloadLimiter, err = limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, cresettings.Default.VaultMaxBlobPayloadSizeLimit)
	}
	require.NoError(t, err)

	return &ReportingPluginConfig{
		MaxBatchSize: bsl,

		PublicKey:                         publicKey,
		PrivateKeyShare:                   privateKeyShare,
		MaxSecretsPerOwner:                msl,
		MaxShareLengthBytes:               shareLimiter,
		MaxCiphertextLengthBytes:          cipherTextLimiter,
		MaxIdentifierOwnerLengthBytes:     ownerLimiter,
		MaxIdentifierNamespaceLengthBytes: namespaceOwnerLimiter,
		MaxIdentifierKeyLengthBytes:       keyLimiter,
		MaxRequestBatchSize:               requestBatchSizeLimiter,
		MaxBlobPayloadBytes:               maxBlobPayloadLimiter,
		VaultForceEmptyOCRRounds:          limits.NewGateLimiter(false),
		VaultOptimizationsEnabled:         limits.NewGateLimiter(false),
	}
}

func mockUnmarshalBlob(data []byte) (ocr3_1types.BlobHandle, error) {
	return ocr3_1types.BlobHandle{}, nil
}

func mockMarshalBlob(ocr3_1types.BlobHandle) ([]byte, error) {
	return []byte{}, nil
}
