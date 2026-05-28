package vault

import (
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/dkgrecipientkey"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPlugin_ReportingPluginFactory_UsesDefaultsIfNotProvidedInOffchainConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()

	_, orm := setupORM(t)
	dkgrecipientKey, err := dkgrecipientkey.New()
	require.NoError(t, err)
	instanceID := "instanceID"
	_ = writeDKGPackage(t, orm, dkgrecipientKey, instanceID)

	lpk := vaultcap.NewLazyPublicKey()
	rpf, err := NewReportingPluginFactory(lggr, store, orm, &dkgrecipientKey, lpk, limits.Factory{Settings: cresettings.DefaultGetter}, testRequestLifecycleTracker(t, lggr))
	require.NoError(t, err)

	cfg := vaultcommon.ReportingPluginConfig{
		DKGInstanceID: &instanceID,
	}
	cfgb, err := proto.Marshal(&cfg)
	require.NoError(t, err)
	rp, info, err := rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgb, N: 10, F: 3}, nil)
	require.NoError(t, err)

	typedRP := rp.(*ReportingPlugin)
	assertLimit(t, cresettings.Default.VaultPluginBatchSizeLimit.DefaultValue, typedRP.cfg.MaxBatchSize)
	assert.NotNil(t, typedRP.cfg.PublicKey)
	assert.NotNil(t, typedRP.cfg.PrivateKeyShare)
	assertLimit(t, 100, typedRP.cfg.MaxSecretsPerOwner)
	assertLimit(t, 2000, typedRP.cfg.MaxCiphertextLengthBytes)
	assertLimit(t, 64, typedRP.cfg.MaxIdentifierOwnerLengthBytes)
	assertLimit(t, 64, typedRP.cfg.MaxIdentifierNamespaceLengthBytes)
	assertLimit(t, 64, typedRP.cfg.MaxIdentifierKeyLengthBytes)

	infoObject, ok := info.(ocr3_1types.ReportingPluginInfo1)
	assert.True(t, ok, "ReportingPluginInfo not of type ReportingPluginInfo1")
	assert.Equal(t, "VaultReportingPlugin", infoObject.Name)
	assert.Equal(t, int(cresettings.Default.VaultMaxQuerySizeLimit.DefaultValue), infoObject.Limits.MaxQueryBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxObservationSizeLimit.DefaultValue), infoObject.Limits.MaxObservationBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxReportsPlusPrecursorSizeLimit.DefaultValue), infoObject.Limits.MaxReportsPlusPrecursorBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxReportSizeLimit.DefaultValue), infoObject.Limits.MaxReportBytes)
	assert.Equal(t, cresettings.Default.VaultMaxReportCount.DefaultValue, infoObject.Limits.MaxReportCount)
	assert.Equal(t, int(cresettings.Default.VaultMaxKeyValueModifiedKeysPlusValuesSizeLimit.DefaultValue), infoObject.Limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxBlobPayloadSizeLimit.DefaultValue), infoObject.Limits.MaxBlobPayloadBytes)

	// Verify that configProto overrides apply to MaxSecretsPerOwner,
	// while MaxBatchSize and other fields remain at cresettings defaults.
	cfg = vaultcommon.ReportingPluginConfig{
		BatchSize:                                     2,
		MaxSecretsPerOwner:                            2,
		MaxCiphertextLengthBytes:                      2,
		MaxIdentifierOwnerLengthBytes:                 2,
		MaxIdentifierNamespaceLengthBytes:             2,
		MaxIdentifierKeyLengthBytes:                   2,
		LimitsMaxQueryLength:                          2,
		LimitsMaxObservationLength:                    2,
		LimitsMaxReportsPlusPrecursorLength:           2,
		LimitsMaxReportLength:                         2,
		LimitsMaxReportCount:                          2,
		LimitsMaxKeyValueModifiedKeysPlusValuesLength: 2,
		LimitsMaxBlobPayloadLength:                    2,
		DKGInstanceID:                                 &instanceID,
	}
	cfgb, err = proto.Marshal(&cfg)
	require.NoError(t, err)

	rp, info, err = rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgb, N: 10, F: 3}, nil)
	require.NoError(t, err)

	typedRP = rp.(*ReportingPlugin)
	assertLimit(t, cresettings.Default.VaultPluginBatchSizeLimit.DefaultValue, typedRP.cfg.MaxBatchSize)
	assertLimit(t, 2, typedRP.cfg.MaxSecretsPerOwner)
	assertLimit(t, 2000, typedRP.cfg.MaxCiphertextLengthBytes)
	assertLimit(t, 64, typedRP.cfg.MaxIdentifierOwnerLengthBytes)
	assertLimit(t, 64, typedRP.cfg.MaxIdentifierNamespaceLengthBytes)
	assertLimit(t, 64, typedRP.cfg.MaxIdentifierKeyLengthBytes)

	infoObject, ok = info.(ocr3_1types.ReportingPluginInfo1)
	assert.True(t, ok, "ReportingPluginInfo not of type ReportingPluginInfo1")
	assert.Equal(t, "VaultReportingPlugin", infoObject.Name)
	assert.Equal(t, int(cresettings.Default.VaultMaxQuerySizeLimit.DefaultValue), infoObject.Limits.MaxQueryBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxObservationSizeLimit.DefaultValue), infoObject.Limits.MaxObservationBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxReportsPlusPrecursorSizeLimit.DefaultValue), infoObject.Limits.MaxReportsPlusPrecursorBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxReportSizeLimit.DefaultValue), infoObject.Limits.MaxReportBytes)
	assert.Equal(t, cresettings.Default.VaultMaxReportCount.DefaultValue, infoObject.Limits.MaxReportCount)
	assert.Equal(t, int(cresettings.Default.VaultMaxKeyValueModifiedKeysPlusValuesSizeLimit.DefaultValue), infoObject.Limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxBlobPayloadSizeLimit.DefaultValue), infoObject.Limits.MaxBlobPayloadBytes)
}

func TestPlugin_ReportingPluginFactory_PassesValidate(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()

	_, orm := setupORM(t)
	dkgrecipientKey, err := dkgrecipientkey.New()
	require.NoError(t, err)
	instanceID := "instanceID"
	_ = writeDKGPackage(t, orm, dkgrecipientKey, instanceID)

	lpk := vaultcap.NewLazyPublicKey()
	rpf, err := NewReportingPluginFactory(lggr, store, orm, &dkgrecipientKey, lpk, limits.Factory{Settings: cresettings.DefaultGetter}, testRequestLifecycleTracker(t, lggr))
	require.NoError(t, err)

	cfg := vaultcommon.ReportingPluginConfig{
		DKGInstanceID: &instanceID,
	}
	cfgb, err := proto.Marshal(&cfg)
	require.NoError(t, err)
	_, info, err := rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgb, N: 10, F: 3}, nil)
	require.NoError(t, err)

	infoObject, ok := info.(ocr3_1types.ReportingPluginInfo1)
	require.True(t, ok, "ReportingPluginInfo not of type ReportingPluginInfo1")
	validateErr := infoObject.Validate()
	require.NoError(t, validateErr)
}

func TestPlugin_ReportingPluginFactory_UseDKGResult(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()

	// Simulate DKG for a single recipient.
	_, orm := setupORM(t)
	dkgrecipientKey, err := dkgrecipientkey.New()
	require.NoError(t, err)

	instanceID := "instanceID"
	pkg := writeDKGPackage(t, orm, dkgrecipientKey, "instanceID")

	expectedTDH2MasterPublicKey, err := tdh2shim.TDH2PublicKeyFromDKGResult(pkg)
	require.NoError(t, err)
	expectedKeyShare, err := tdh2shim.TDH2PrivateShareFromDKGResult(pkg, dkgrecipientKey)
	require.NoError(t, err)

	lpk := vaultcap.NewLazyPublicKey()
	rpf, err := NewReportingPluginFactory(lggr, store, orm, &dkgrecipientKey, lpk, limits.Factory{Settings: cresettings.DefaultGetter}, testRequestLifecycleTracker(t, lggr))
	require.NoError(t, err)

	instanceIDString := instanceID
	rpCfg := vaultcommon.ReportingPluginConfig{
		DKGInstanceID: &instanceIDString,
	}
	cfgBytes, err := proto.Marshal(&rpCfg)
	require.NoError(t, err)
	rp, info, err := rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgBytes, N: 10, F: 3}, nil)
	require.NoError(t, err)

	typedRP := rp.(*ReportingPlugin)
	assertLimit(t, cresettings.Default.VaultPluginBatchSizeLimit.DefaultValue, typedRP.cfg.MaxBatchSize)

	pkBytes, err := typedRP.cfg.PublicKey.Marshal()
	require.NoError(t, err)
	pk := &tdh2.PublicKey{}
	err = pk.Unmarshal(pkBytes)
	require.NoError(t, err)
	assert.True(t, pk.Equal(expectedTDH2MasterPublicKey))

	ksBytes, err := typedRP.cfg.PrivateKeyShare.Marshal()
	require.NoError(t, err)
	ks := &tdh2.PrivateShare{}
	err = ks.Unmarshal(ksBytes)
	require.NoError(t, err)
	assert.Equal(t, expectedKeyShare, ks)

	infoObject, ok := info.(ocr3_1types.ReportingPluginInfo1)
	assert.True(t, ok, "ReportingPluginInfo not of type ReportingPluginInfo1")
	assert.Equal(t, "VaultReportingPlugin", infoObject.Name)

	key, err := lpk.Get().Marshal()
	require.NoError(t, err)
	assert.Equal(t, pkBytes, key)
}

func TestPlugin_ReportingPluginFactory_InvalidParams(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()

	lpk := vaultcap.NewLazyPublicKey()

	_, orm := setupORM(t)
	_, err := NewReportingPluginFactory(lggr, store, orm, nil, lpk, limits.Factory{Settings: cresettings.DefaultGetter}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "DKG recipient key cannot be nil when using result package db")

	_, err = NewReportingPluginFactory(lggr, store, nil, nil, lpk, limits.Factory{Settings: cresettings.DefaultGetter}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "result package db cannot be nil")

	dkgrecipientKey, err := dkgrecipientkey.New()
	require.NoError(t, err)
	_, err = NewReportingPluginFactory(lggr, store, orm, &dkgrecipientKey, lpk, limits.Factory{Settings: cresettings.DefaultGetter}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "request lifecycle tracker cannot be nil")
}
