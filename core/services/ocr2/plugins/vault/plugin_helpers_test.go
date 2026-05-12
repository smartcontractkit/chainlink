package vault

import (
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
	lggr                              logger.Logger
	store                             *requests.Store[*vaulttypes.Request]
	publicKey                         *tdh2easy.PublicKey
	privateKeyShare                   *tdh2easy.PrivateShare
	onchainCfg                        ocr3types.ReportingPluginConfig
	maxSecretsPerOwner                int
	maxCiphertextLengthBytes          int
	maxIdentifierOwnerLengthBytes     int
	maxIdentifierNamespaceLengthBytes int
	maxIdentifierKeyLengthBytes       int
	maxRequestBatchSize               int
	batchSize                         int
	orgIDAsSecretOwnerEnabled         bool
	marshalBlob                       func(ocr3_1types.BlobHandle) ([]byte, error)
	unmarshalBlob                     func([]byte) (ocr3_1types.BlobHandle, error)
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

func withOrgIDEnabled() testPluginOption {
	return func(o *testPluginBuildOpts) { o.orgIDAsSecretOwnerEnabled = true }
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

func withMarshalBlob(fn func(ocr3_1types.BlobHandle) ([]byte, error)) testPluginOption {
	return func(o *testPluginBuildOpts) { o.marshalBlob = fn }
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
		o.maxIdentifierKeyLengthBytes, o.maxRequestBatchSize)
	if o.orgIDAsSecretOwnerEnabled {
		cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)
	}
	return &ReportingPlugin{
		lggr:          o.lggr,
		store:         o.store,
		metrics:       newTestMetrics(t),
		cfg:           cfg,
		onchainCfg:    o.onchainCfg,
		validator:     makeTestValidator(cfg),
		marshalBlob:   o.marshalBlob,
		unmarshalBlob: o.unmarshalBlob,
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
		OrgIDAsSecretOwnerEnabled:         limits.NewGateLimiter(false),
		VaultForceEmptyOCRRounds:          limits.NewGateLimiter(false),
	}
}

func mockUnmarshalBlob(data []byte) (ocr3_1types.BlobHandle, error) {
	return ocr3_1types.BlobHandle{}, nil
}

func mockMarshalBlob(ocr3_1types.BlobHandle) ([]byte, error) {
	return []byte{}, nil
}
