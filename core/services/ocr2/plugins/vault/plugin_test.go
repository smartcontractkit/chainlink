package vault

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/smdkg/dummydkg"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/dkgrecipientkey"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func writeDKGPackage(t *testing.T, orm dkgocrtypes.ResultPackageDatabase, key dkgocrtypes.P256Keyring, instanceID string) dkgocrtypes.ResultPackage {
	pkg, err := dummydkg.NewResultPackage(dkgocrtypes.InstanceID(instanceID), dkgocrtypes.ReportingPluginConfig{
		DealerPublicKeys:    []dkgocrtypes.P256ParticipantPublicKey{key.PublicKey()},
		RecipientPublicKeys: []dkgocrtypes.P256ParticipantPublicKey{key.PublicKey()},
		T:                   1,
	}, []dkgocrtypes.P256Keyring{key})
	require.NoError(t, err)

	pkgBin, err := pkg.MarshalBinary()
	require.NoError(t, err)
	require.NoError(t, orm.WriteResultPackage(t.Context(), dkgocrtypes.InstanceID(instanceID), dkgocrtypes.ResultPackageDatabaseValue{
		ConfigDigest:            [32]byte{0x1, 0x2, 0x3, 0x4},
		SeqNr:                   1,
		ReportWithResultPackage: pkgBin,
		Signatures: []types.AttributedOnchainSignature{
			{
				Signature: []byte{0x5, 0x6, 0x7, 0x8},
				Signer:    1,
			},
		},
	}))

	return pkg
}

func assertLimit[N limits.Number](t *testing.T, expected int, limiter limits.BoundLimiter[N]) {
	ctx := contexts.WithCRE(t.Context(), contexts.CRE{Owner: "foo"})
	l, err := limiter.Limit(ctx)
	require.NoError(t, err)

	assert.Equal(t, expected, int(l))
}

func TestPlugin_ReportingPluginFactory_UsesDefaultsIfNotProvidedInOffchainConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()

	_, orm := setupORM(t)
	dkgrecipientKey, err := dkgrecipientkey.New()
	require.NoError(t, err)
	instanceID := "instanceID"
	_ = writeDKGPackage(t, orm, dkgrecipientKey, instanceID)

	lpk := vaultcap.NewLazyPublicKey()
	rpf, err := NewReportingPluginFactory(lggr, store, orm, &dkgrecipientKey, lpk, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	cfg := vaultcommon.ReportingPluginConfig{
		DKGInstanceID: &instanceID,
	}
	cfgb, err := proto.Marshal(&cfg)
	require.NoError(t, err)
	rp, info, err := rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgb}, nil)
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
	assert.Equal(t, cresettings.Default.VaultMaxKeyValueModifiedKeys.DefaultValue, infoObject.Limits.MaxKeyValueModifiedKeys)
	assert.Equal(t, int(cresettings.Default.VaultMaxBlobPayloadSizeLimit.DefaultValue), infoObject.Limits.MaxBlobPayloadBytes)
	assert.Equal(t, int(cresettings.Default.VaultMaxPerOracleUnexpiredBlobCumulativePayloadSizeLimit.DefaultValue), infoObject.Limits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes)
	assert.Equal(t, cresettings.Default.VaultMaxPerOracleUnexpiredBlobCount.DefaultValue, infoObject.Limits.MaxPerOracleUnexpiredBlobCount)

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

	rp, info, err = rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgb}, nil)
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
	rpf, err := NewReportingPluginFactory(lggr, store, orm, &dkgrecipientKey, lpk, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	instanceIDString := instanceID
	rpCfg := vaultcommon.ReportingPluginConfig{
		DKGInstanceID: &instanceIDString,
	}
	cfgBytes, err := proto.Marshal(&rpCfg)
	require.NoError(t, err)
	rp, info, err := rpf.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{OffchainConfig: cfgBytes}, nil)
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
	_, err := NewReportingPluginFactory(lggr, store, orm, nil, lpk, limits.Factory{Settings: cresettings.DefaultGetter})
	require.Error(t, err)
	require.Contains(t, err.Error(), "DKG recipient key cannot be nil when using result package db")

	_, err = NewReportingPluginFactory(lggr, store, nil, nil, lpk, limits.Factory{Settings: cresettings.DefaultGetter})
	require.Error(t, err)
	require.Contains(t, err.Error(), "result package db cannot be nil")
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
	}
}

func TestPlugin_Observation_NothingInBatch(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Empty(t, obs.Observations)
}

func TestPlugin_Observation_PendingQueueEnabled_EmptyPendingQueue(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
		marshalBlob:   mockMarshalBlob,
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	expectedID := "request-1"
	err = store.Add(&vaulttypes.Request{Payload: p, IDVal: expectedID})
	require.NoError(t, err)

	expectedID2 := "request-2"
	err = store.Add(&vaulttypes.Request{Payload: p, IDVal: expectedID2})
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	bf := &blobber{}

	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	// There was no pending queue in the KV store, so this will be empty
	assert.Empty(t, obs.Observations)

	// We expect the pending queue observation to contain the request in the local queue.
	assert.Len(t, obs.PendingQueueItems, 2)

	assertPendingQueueItemsContain(t, bf.blobs, map[string]proto.Message{
		expectedID:  p,
		expectedID2: p,
	})

	assert.NotEmpty(t, obs.SortNonce)
}

func TestPlugin_Observation_PendingQueueEnabled_WithPendingQueueProvided(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	expectedID := "request-1"
	err = store.Add(&vaulttypes.Request{Payload: p, IDVal: expectedID})
	require.NoError(t, err)

	expectedID2 := "request-2"
	err = store.Add(&vaulttypes.Request{Payload: p, IDVal: expectedID2})
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	d := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-3",
		Ids:       []*vaultcommon.SecretIdentifier{id},
	}
	anyd, err := anypb.New(d)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-3", Item: anyd},
		},
	)
	require.NoError(t, err)

	bf := &blobber{}

	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	// We'll observe what we found in the KV store's pending queue, i.e. `d`.
	// This key doesn't exist in the store, so we should get a key does not exist error.
	assert.Len(t, obs.Observations, 1)
	gotO := obs.Observations[0]
	assert.True(t, proto.Equal(gotO.GetDeleteSecretsRequest(), d))

	assert.Len(t, gotO.GetDeleteSecretsResponse().Responses, 1)
	gotResp := gotO.GetDeleteSecretsResponse().Responses[0]
	assert.Equal(t, "key does not exist", gotResp.Error)

	// We expect the pending queue observation to contain the request in the local queue.
	assert.Len(t, obs.PendingQueueItems, 2)

	assertPendingQueueItemsContain(t, bf.blobs, map[string]proto.Message{
		expectedID:  p,
		expectedID2: p,
	})

	assert.NotEmpty(t, obs.SortNonce)
}

func TestPlugin_Observation_PendingQueueEnabled_ItemBothInPendingQueueAndLocalQueue(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	expectedID := "request-1"
	err = store.Add(&vaulttypes.Request{Payload: p, IDVal: expectedID})
	require.NoError(t, err)

	expectedID2 := "request-2"
	err = store.Add(&vaulttypes.Request{Payload: p, IDVal: expectedID2})
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	bf := &blobber{}

	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-2", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	// We'll observe what we found in the KV store's pending queue, i.e. `d`.
	// This key doesn't exist in the store, so we should get a key does not exist error.
	assert.Len(t, obs.Observations, 1)
	gotO := obs.Observations[0]
	assert.True(t, proto.Equal(gotO.GetGetSecretsRequest(), p))

	assert.Len(t, gotO.GetGetSecretsResponse().Responses, 1)
	gotResp := gotO.GetGetSecretsResponse().Responses[0]
	assert.Equal(t, "key does not exist", gotResp.GetError())

	// We expect the pending queue observation to contain the request in the local queue,
	// however it should exclude anything that's already in the pending queue. That'll get
	// processed this round.
	assert.Len(t, obs.PendingQueueItems, 1)

	assertPendingQueueItemsEqual(t, expectedID, bf.blobs[0], p)

	assert.NotEmpty(t, obs.SortNonce)
}

func assertPendingQueueItemsEqual(t *testing.T, expectedID string, got []byte, expectedPayload proto.Message) {
	t.Helper()
	gotMsg := &vaultcommon.StoredPendingQueueItem{}
	err := proto.Unmarshal(got, gotMsg)
	require.NoError(t, err)

	assert.Equal(t, expectedID, gotMsg.Id)
	gotm, err := gotMsg.Item.UnmarshalNew()
	require.NoError(t, err)

	assert.True(t, proto.Equal(expectedPayload, gotm))
}

func assertPendingQueueItemsContain(t *testing.T, gotItems [][]byte, expected map[string]proto.Message) {
	t.Helper()

	require.Len(t, gotItems, len(expected))

	remaining := make(map[string]proto.Message, len(expected))
	for id, payload := range expected {
		remaining[id] = payload
	}

	for _, got := range gotItems {
		gotMsg := &vaultcommon.StoredPendingQueueItem{}
		err := proto.Unmarshal(got, gotMsg)
		require.NoError(t, err)

		expectedPayload, ok := remaining[gotMsg.Id]
		require.True(t, ok, "unexpected pending queue item id %q", gotMsg.Id)

		gotPayload, err := gotMsg.Item.UnmarshalNew()
		require.NoError(t, err)
		assert.True(t, proto.Equal(expectedPayload, gotPayload))

		delete(remaining, gotMsg.Id)
	}

	assert.Empty(t, remaining)
}

type blockingBlobBroadcastFetcher struct {
	targetStarts int32
	started      atomic.Int32
	maxInFlight  atomic.Int32
	inFlight     atomic.Int32
	allStarted   chan struct{}
	release      chan struct{}
	once         sync.Once
}

func (b *blockingBlobBroadcastFetcher) BroadcastBlob(ctx context.Context, _ []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	currentInFlight := b.inFlight.Add(1)
	defer b.inFlight.Add(-1)

	for {
		maxInFlight := b.maxInFlight.Load()
		if currentInFlight <= maxInFlight || b.maxInFlight.CompareAndSwap(maxInFlight, currentInFlight) {
			break
		}
	}

	if b.started.Add(1) == b.targetStarts {
		b.once.Do(func() { close(b.allStarted) })
	}

	select {
	case <-b.release:
		return ocr3_1types.BlobHandle{}, nil
	case <-ctx.Done():
		return ocr3_1types.BlobHandle{}, ctx.Err()
	}
}

func (b *blockingBlobBroadcastFetcher) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	panic("FetchBlob should not be called in Observation tests")
}

type errorBlobBroadcastFetcher struct {
	err error
}

func (e *errorBlobBroadcastFetcher) BroadcastBlob(context.Context, []byte, ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	return ocr3_1types.BlobHandle{}, e.err
}

func (e *errorBlobBroadcastFetcher) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	panic("FetchBlob should not be called in Observation tests")
}

func TestPlugin_Observation_PendingQueueEnabled_BroadcastsPendingQueueBlobsInParallel(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}

	require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: "request-1"}))
	require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: "request-2"}))

	rdr := &kv{m: make(map[string]response)}
	bf := &blockingBlobBroadcastFetcher{
		targetStarts: 2,
		allStarted:   make(chan struct{}),
		release:      make(chan struct{}),
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, bf)
		errCh <- err
	}()

	select {
	case <-bf.allStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for concurrent blob broadcasts")
	}

	close(bf.release)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Observation to finish")
	}

	assert.Equal(t, int32(2), bf.maxInFlight.Load())
}

func TestPlugin_Observation_PendingQueueEnabled_BroadcastBlobError(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}

	require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: "request-1"}))
	rdr := &kv{m: make(map[string]response)}

	obs, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &errorBlobBroadcastFetcher{err: errors.New("boom")})
	require.NoError(t, err)
	require.NotNil(t, obs)

	warnLogs := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
	assert.Equal(t, 1, warnLogs.Len())
	fields := warnLogs.All()[0].ContextMap()
	assert.Equal(t, "request-1", fields["requestID"])
	assert.Contains(t, fmt.Sprint(fields["err"]), "boom")
}

func TestPlugin_Observation_GetSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name     string
		id       *vaultcommon.SecretIdentifier
		maxIDLen int
		err      string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "invalid secret identifier: cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "invalid secret identifier: key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		lggr := logger.TestLogger(t)
		store := requests.NewStore[*vaulttypes.Request]()
		maxIDLen := 256
		if tc.maxIDLen > 0 {
			maxIDLen = tc.maxIDLen
		}
		r := &ReportingPlugin{
			lggr:    lggr,
			store:   store,
			metrics: newTestMetrics(t),
			cfg: makeReportingPluginConfig(
				t,
				10,
				nil,
				nil,
				1,
				maxIDLen/3,
				maxIDLen/3,
				maxIDLen/3,
				maxIDLen/3,
				10,
			),
			marshalBlob:   mockMarshalBlob,
			unmarshalBlob: mockUnmarshalBlob,
		}
		seqNr := uint64(1)
		rdr := &kv{
			m: make(map[string]response),
		}
		p := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{
					Id:             tc.id,
					EncryptionKeys: []string{"foo"},
				},
			},
		}
		anyp, err := anypb.New(p)
		require.NoError(t, err)
		err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
			[]*vaultcommon.StoredPendingQueueItem{
				{Id: "request-1", Item: anyp},
			},
		)
		require.NoError(t, err)
		bf := &blobber{}
		data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		err = proto.Unmarshal(data, obs)
		require.NoError(t, err)

		assert.Len(t, obs.Observations, 1)
		o := obs.Observations[0]

		assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
		assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

		batchResp := o.GetGetSecretsResponse()
		assert.Len(t, p.Requests, 1)
		assert.Len(t, p.Requests, len(batchResp.Responses))

		assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_GetSecretsRequest_FillsInNamespace(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	ciphertext, err := tdh2easy.Encrypt(pk, plaintext)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	createdID := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), createdID, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
		WorkflowOwner: "owner",
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	bf := &blobber{}
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(batchResp.Responses[0].Id, createdID))
}

func TestPlugin_Observation_GetSecretsRequest_OrgIdLabelAcceptedWhenEnabled(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	cfg := makeReportingPluginConfig(t, 10, pk, shares[0], 1, 1024, 100, 100, 100, 10)
	cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)

	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		cfg:           cfg,
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

	encrypted, err := vaultutils.EncryptSecretWithOrgID("my-secret-value", pk, orgID)
	require.NoError(t, err)
	ciphertextBytes, err := hex.DecodeString(encrypted)
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{hex.EncodeToString(pubK[:])},
			},
		},
		OrgId: orgID,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetGetSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	assert.Empty(t, batchResp.Responses[0].GetError())
}

func TestPlugin_Observation_GetSecretsRequest_OrgIdLabelRejectedWhenDisabled(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	cfg := makeReportingPluginConfig(t, 10, pk, shares[0], 1, 1024, 100, 100, 100, 10)

	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		cfg:           cfg,
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

	encrypted, err := vaultutils.EncryptSecretWithOrgID("my-secret-value", pk, orgID)
	require.NoError(t, err)
	ciphertextBytes, err := hex.DecodeString(encrypted)
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{hex.EncodeToString(pubK[:])},
			},
		},
		OrgId: orgID,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetGetSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	assert.Contains(t, batchResp.Responses[0].GetError(), "failed to handle get secret request")
}

func TestPlugin_Observation_GetSecretsRequest_SecretDoesNotExist(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	bf := &blobber{}
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "key does not exist")
}

func TestPlugin_Observation_GetSecretsRequest_SecretExistsButIsIncorrect(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("invalid-ciphertext"),
	})
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	// The error returned is user-friendly
	assert.Contains(t, resp.GetError(), "failed to handle get secret request")

	// Inspect logs to get true source of error
	logs := observed.FilterMessage("failed to observe get secret request item")
	assert.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	errString := fields["error"]
	assert.Contains(t, errString, "failed to unmarshal ciphertext")
}

func TestPlugin_Observation_GetSecretsRequest_PublicKeyIsInvalid(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	ciphertext, err := tdh2easy.Encrypt(pk, plaintext)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
		WorkflowOwner: "owner",
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to convert public key to bytes")
}

func TestPlugin_Observation_GetSecretsRequest_SecretLabelIsInvalid(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	var label [32]byte
	ownerAddress := common.HexToAddress("0x0001020304050607080900010203040506070809")
	copy(label[12:], ownerAddress.Bytes()) // left-pad with 12 zero
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to handle get secret request")
}

func TestPlugin_Observation_GetSecretsRequest_Success(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	var label [32]byte
	ownerAddress := common.HexToAddress(owner)
	copy(label[12:], ownerAddress.Bytes()) // left-pad with 12 zero
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, privK, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
		WorkflowOwner: owner,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Empty(t, resp.GetError())

	assert.Equal(t, hex.EncodeToString(ciphertextBytes), resp.GetData().EncryptedValue)

	assert.Len(t, resp.GetData().EncryptedDecryptionKeyShares, 1)
	shareString := resp.GetData().EncryptedDecryptionKeyShares[0].Shares[0]

	share, err := hex.DecodeString(shareString)
	require.NoError(t, err)
	msg, ok := box.OpenAnonymous(nil, share, pubK, privK)
	assert.True(t, ok)

	ds := &tdh2easy.DecryptionShare{}
	err = ds.Unmarshal(msg)
	require.NoError(t, err)

	ct := &tdh2easy.Ciphertext{}
	ctb, err := hex.DecodeString(resp.GetData().EncryptedValue)
	require.NoError(t, err)
	err = ct.UnmarshalVerify(ctb, pk)
	require.NoError(t, err)

	gotSecret, err := tdh2easy.Aggregate(ct, []*tdh2easy.DecryptionShare{ds}, 3)
	require.NoError(t, err)

	assert.Equal(t, plaintext, gotSecret)
}

func TestPlugin_Observation_CreateSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name     string
		id       *vaultcommon.SecretIdentifier
		maxIDLen int
		err      string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "invalid secret identifier: cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "invalid secret identifier: key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		lggr := logger.TestLogger(t)
		store := requests.NewStore[*vaulttypes.Request]()
		maxIDLen := 256
		if tc.maxIDLen > 0 {
			maxIDLen = tc.maxIDLen
		}
		r := &ReportingPlugin{
			lggr:          lggr,
			store:         store,
			metrics:       newTestMetrics(t),
			marshalBlob:   mockMarshalBlob,
			unmarshalBlob: mockUnmarshalBlob,
			cfg: makeReportingPluginConfig(
				t,
				10,
				nil,
				nil,
				1,
				1024,
				maxIDLen/3,
				maxIDLen/3,
				maxIDLen/3,
				10,
			),
		}

		seqNr := uint64(1)
		rdr := &kv{
			m: make(map[string]response),
		}
		p := &vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id:             tc.id,
					EncryptedValue: "foo",
				},
			},
		}
		anyp, err := anypb.New(p)
		require.NoError(t, err)
		err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
			[]*vaultcommon.StoredPendingQueueItem{
				{Id: "request-1", Item: anyp},
			},
		)
		require.NoError(t, err)
		bf := &blobber{}
		data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		err = proto.Unmarshal(data, obs)
		require.NoError(t, err)

		assert.Len(t, obs.Observations, 1)
		o := obs.Observations[0]

		assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)
		assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

		batchResp := o.GetCreateSecretsResponse()
		assert.Len(t, p.EncryptedSecrets, 1)
		assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

		assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_CreateSecretsRequest_DisallowsDuplicateRequests(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: "foo",
			},
			{
				Id:             id,
				EncryptedValue: "bla",
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 2)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "duplicate request for secret identifier")

	assert.True(t, proto.Equal(p.EncryptedSecrets[1].Id, batchResp.Responses[1].Id))
	resp = batchResp.Responses[1]
	assert.Contains(t, resp.GetError(), "duplicate request for secret identifier")
}

func TestPlugin_StateTransition_CreateSecretsRequest_CorrectlyTracksLimits(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	id1 := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	req1 := &vaultcommon.CreateSecretsRequest{
		RequestId: "req1",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id1,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp1 := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id1,
				Success: false,
			},
		},
	}

	id2 := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret2",
	}
	req2 := &vaultcommon.CreateSecretsRequest{
		RequestId: "req2",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id2,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp2 := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id2,
				Success: false,
			},
		},
	}

	obs := marshalObservations(t, observation{id1, req1, resp1}, observation{id2, req2, resp2})

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: obs},
			{Observer: 1, Observation: obs},
			{Observer: 2, Observation: obs},
		},
		rdr,
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 2)

	o1 := os.Outcomes[0]
	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o1.RequestType)
	assert.Len(t, o1.GetCreateSecretsResponse().Responses, 1)
	r1 := o1.GetCreateSecretsResponse().Responses[0]
	assert.True(t, r1.Success)

	o2 := os.Outcomes[1]
	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o2.RequestType)
	assert.Len(t, o2.GetCreateSecretsResponse().Responses, 1)
	r2 := o2.GetCreateSecretsResponse().Responses[0]
	assert.False(t, r2.Success)
	assert.Contains(t, r2.GetError(), "owner has reached maximum number of secrets")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: "foo",
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "invalid hex encoding for ciphertext")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_TooLong(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			10,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ciphertext := []byte("a quick brown fox jumps over the lazy dog")
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertext),
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "ciphertext size exceeds maximum allowed size: 10b")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_EncryptedWithWrongPublicKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	// Wrong key
	_, wrongPublicKey, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	// Right key
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ct, err := tdh2easy.Encrypt(wrongPublicKey, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "failed to verify ciphertext")
}

func TestPlugin_Observation_CreateSecretsRequest_SecretLabelIsInvalid(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	owner := "0x1234567890abcdef1234567890abcdef12345678"
	wrongOwner := "0x0001020304050607080900010203040506070809"

	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "secret",
	}

	var wrongLabel [32]byte
	wrongOwnerAddr := common.HexToAddress(wrongOwner)
	copy(wrongLabel[12:], wrongOwnerAddr.Bytes())
	ct, err := tdh2easy.EncryptWithLabel(pk, []byte("my secret value"), wrongLabel)
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)

	batchResp := o.GetCreateSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "failed to verify ciphertext")
}

func TestPlugin_Observation_UpdateSecretsRequest_SecretLabelIsInvalid(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	owner := "0x1234567890abcdef1234567890abcdef12345678"
	wrongOwner := "0x0001020304050607080900010203040506070809"

	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "secret",
	}

	var wrongLabel [32]byte
	wrongOwnerAddr := common.HexToAddress(wrongOwner)
	copy(wrongLabel[12:], wrongOwnerAddr.Bytes())
	ct, err := tdh2easy.EncryptWithLabel(pk, []byte("my secret value"), wrongLabel)
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_UPDATE_SECRETS, o.RequestType)

	batchResp := o.GetUpdateSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "failed to verify ciphertext")
}

func TestPlugin_StateTransition_CreateSecretsRequest_TooManySecretsForOwner(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	kvstore := newTestWriteStore(t, rdr)
	err = kvstore.WriteMetadata(t.Context(), id.Owner, &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{
			{
				Owner:     "owner",
				Namespace: "main",
				Key:       "secret2",
			},
		},
	})
	require.NoError(t, err)

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	req := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
			},
		},
	}
	data, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{
				Observer:    0,
				Observation: marshalObservations(t, observation{id, req, resp}),
			},
		},
		rdr,
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(data, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)
	o := os.Outcomes[0]

	assert.Len(t, o.GetCreateSecretsResponse().Responses, 1)
	assert.Contains(t, o.GetCreateSecretsResponse().Responses[0].Error, "owner has reached maximum number of secrets")
}

func TestPlugin_StateTransition_CreateSecretsRequest_SecretExistsForKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	kvstore := newTestWriteStore(t, rdr)
	err = kvstore.WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("some-ciphertext"),
	})
	require.NoError(t, err)

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	req := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
			},
		},
	}
	data, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{
				Observer:    0,
				Observation: marshalObservations(t, observation{id, req, resp}),
			},
		},
		rdr,
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(data, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)
	o := os.Outcomes[0]

	assert.Len(t, o.GetCreateSecretsResponse().Responses, 1)
	assert.Contains(t, o.GetCreateSecretsResponse().Responses[0].Error, "key already exists")
}

func TestPlugin_Observation_CreateSecretsRequest_Success(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
		WorkflowOwner: "owner",
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Empty(t, resp.GetError())
}

func TestPlugin_Observation_CreateSecretsRequest_OrgIdLabelAcceptedWhenEnabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	cfg := makeReportingPluginConfig(t, 10, pk, shares[0], 1, 1024, 100, 100, 100, 10)
	cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)

	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg:           cfg,
	}

	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "secret",
	}

	encrypted, err := vaultutils.EncryptSecretWithOrgID("my secret value", pk, orgID)
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: id, EncryptedValue: encrypted},
		},
		OrgId: orgID,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetCreateSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	assert.Empty(t, batchResp.Responses[0].GetError())
}

func TestPlugin_Observation_CreateSecretsRequest_OrgIdLabelRejectedWhenDisabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	cfg := makeReportingPluginConfig(t, 10, pk, shares[0], 1, 1024, 100, 100, 100, 10)

	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg:           cfg,
	}

	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "secret",
	}

	encrypted, err := vaultutils.EncryptSecretWithOrgID("my secret value", pk, orgID)
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: id, EncryptedValue: encrypted},
		},
		OrgId: orgID,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetCreateSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	assert.Contains(t, batchResp.Responses[0].GetError(), "does not match any of the provided owner labels")
}

func makeEncryptedShares(t *testing.T, ciphertext *tdh2easy.Ciphertext, privateShare *tdh2easy.PrivateShare, keys []string) []*vaultcommon.EncryptedShares {
	t.Helper()
	share, err := tdh2easy.Decrypt(ciphertext, privateShare)
	require.NoError(t, err)
	shareBytes, err := share.Marshal()
	require.NoError(t, err)

	result := make([]*vaultcommon.EncryptedShares, len(keys))
	for i, pk := range keys {
		pkBytes, err := hex.DecodeString(pk)
		require.NoError(t, err)
		pubKey := [32]byte(pkBytes)
		encrypted, err := box.SealAnonymous(nil, shareBytes, &pubKey, rand.Reader)
		require.NoError(t, err)
		result[i] = &vaultcommon.EncryptedShares{
			EncryptionKey: pk,
			Shares:        []string{base64.StdEncoding.EncodeToString(encrypted)},
		}
	}
	return result
}

func makeGetSecretsObservations(
	t *testing.T,
	numRequests int,
	owner string,
	namespace string,
	encryptionKeys []string,
	encryptedValue string,
	ciphertext *tdh2easy.Ciphertext,
	privateShare *tdh2easy.PrivateShare,
) []byte {
	t.Helper()
	obs := make([]observation, 0, numRequests)
	for i := range numRequests {
		maxKey := fmt.Sprintf("%s%d", strings.Repeat("c", 64-1), i)

		id := &vaultcommon.SecretIdentifier{
			Owner:     owner,
			Namespace: namespace,
			Key:       maxKey,
		}
		req := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{
					Id:             id,
					EncryptionKeys: encryptionKeys,
				},
			},
		}
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{
				{
					Id: id,
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue:               encryptedValue,
							EncryptedDecryptionKeyShares: makeEncryptedShares(t, ciphertext, privateShare, encryptionKeys),
						},
					},
				},
			},
		}
		obs = append(obs, observation{id, req, resp})
	}
	return marshalObservations(t, obs...)
}

type observation struct {
	id   *vaultcommon.SecretIdentifier
	req  proto.Message
	resp proto.Message
}

func marshalObservations(t *testing.T, observations ...observation) []byte {
	obs := &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{},
	}
	for _, ob := range observations {
		o := &vaultcommon.Observation{
			Id: vaulttypes.KeyFor(ob.id),
		}
		switch tr := ob.req.(type) {
		case *vaultcommon.GetSecretsRequest:
			o.RequestType = vaultcommon.RequestType_GET_SECRETS
			o.Request = &vaultcommon.Observation_GetSecretsRequest{
				GetSecretsRequest: tr,
			}
		case *vaultcommon.CreateSecretsRequest:
			o.RequestType = vaultcommon.RequestType_CREATE_SECRETS
			o.Request = &vaultcommon.Observation_CreateSecretsRequest{
				CreateSecretsRequest: tr,
			}
		case *vaultcommon.UpdateSecretsRequest:
			o.RequestType = vaultcommon.RequestType_UPDATE_SECRETS
			o.Request = &vaultcommon.Observation_UpdateSecretsRequest{
				UpdateSecretsRequest: tr,
			}
		case *vaultcommon.DeleteSecretsRequest:
			o.RequestType = vaultcommon.RequestType_DELETE_SECRETS
			o.Request = &vaultcommon.Observation_DeleteSecretsRequest{
				DeleteSecretsRequest: tr,
			}
		case *vaultcommon.ListSecretIdentifiersRequest:
			o.RequestType = vaultcommon.RequestType_DELETE_SECRETS
			o.Request = &vaultcommon.Observation_ListSecretIdentifiersRequest{
				ListSecretIdentifiersRequest: tr,
			}
		}

		switch tr := ob.resp.(type) {
		case *vaultcommon.GetSecretsResponse:
			o.Response = &vaultcommon.Observation_GetSecretsResponse{
				GetSecretsResponse: tr,
			}
		case *vaultcommon.CreateSecretsResponse:
			o.Response = &vaultcommon.Observation_CreateSecretsResponse{
				CreateSecretsResponse: tr,
			}
		case *vaultcommon.UpdateSecretsResponse:
			o.Response = &vaultcommon.Observation_UpdateSecretsResponse{
				UpdateSecretsResponse: tr,
			}
		case *vaultcommon.DeleteSecretsResponse:
			o.Response = &vaultcommon.Observation_DeleteSecretsResponse{
				DeleteSecretsResponse: tr,
			}
		case *vaultcommon.ListSecretIdentifiersResponse:
			o.RequestType = vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS
			o.Response = &vaultcommon.Observation_ListSecretIdentifiersResponse{
				ListSecretIdentifiersResponse: tr,
			}
		}

		obs.Observations = append(obs.Observations, o)
	}

	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return b
}

func TestPlugin_StateTransition_InsufficientObservations(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id1 := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id1,
			},
		},
	}
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id1,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}

	obs1b := marshalObservations(t, observation{id1, req, resp})

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obs1b)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Empty(t, os.Outcomes, 0)

	assert.Equal(t, 1, observed.FilterMessage("insufficient observations found for id").Len())
}

func TestPlugin_StateTransition_GetSecretsRequest_ResponseSizeWithinLimit(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(4, 10)
	require.NoError(t, err)

	numObservers := 10
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 10,
			F: 3,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			100,
			2000,
			64,
			64,
			64,
			10,
		),
	}

	maxOwner := strings.Repeat("a", 64)
	maxNamespace := strings.Repeat("b", 64)

	numEncryptionKeys := 10
	encryptionKeys := make([]string, numEncryptionKeys)
	for i := range numEncryptionKeys {
		pubK, _, err2 := box.GenerateKey(rand.Reader)
		require.NoError(t, err2)
		encryptionKeys[i] = hex.EncodeToString(pubK[:])
	}

	plaintext := make([]byte, 1)
	_, err = rand.Read(plaintext)
	require.NoError(t, err)
	var label [32]byte
	copy(label[:], maxOwner[:32])
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)
	require.LessOrEqual(t, len(ciphertextBytes), 2000)
	encryptedValue := hex.EncodeToString(ciphertextBytes)

	// Create 10 observations from different observers, each with a distinct decryption share.
	aos := make([]types.AttributedObservation, numObservers)
	for i := range numObservers {
		aos[i] = types.AttributedObservation{
			Observer:    commontypes.OracleID(i), //nolint:gosec // G115 range is well within uint8 bounds
			Observation: types.Observation(makeGetSecretsObservations(t, 10, maxOwner, maxNamespace, encryptionKeys, encryptedValue, ciphertext, shares[i])),
		}
	}

	kvStore := &kv{m: make(map[string]response)}
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		aos, kvStore, nil)
	require.NoError(t, err)

	t.Logf("StateTransition response size: %d bytes (%.2f KB)", len(reportPrecursor), float64(len(reportPrecursor))/1024.0)
	maxResponseSize := 512 * 1024
	assert.LessOrEqual(t, len(reportPrecursor), maxResponseSize,
		"StateTransition response size %d exceeds 512KB limit", len(reportPrecursor))
}

func TestPlugin_ValidateObservations_InvalidObservations(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id1 := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id1,
			},
		},
	}
	resp := &vaultcommon.CreateSecretsResponse{}

	// Request and response don't match
	obsb := marshalObservations(t, observation{id1, req, resp})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	require.ErrorContains(t, err, "GetSecrets observation must have both request and response")

	// Invalid observation -- data can't be unmarshaled
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation([]byte("hello world"))},
		kv,
		nil,
	)

	require.ErrorContains(t, err, "failed to unmarshal observations")

	// Invalid observation -- a single observation set has observations for multiple request ids
	correctResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id1,
			},
		},
	}
	obsb = marshalObservations(t, observation{id1, req, correctResp}, observation{id1, req, correctResp})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	assert.ErrorContains(t, err, "invalid observation: a single observation cannot contain duplicate observations for the same request id")
}

func TestPlugin_ValidateObservations_IncludesAllItemsInPendingQueue(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	d := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-1",
		Ids:       []*vaultcommon.SecretIdentifier{id},
	}
	anyd, err := anypb.New(d)
	require.NoError(t, err)
	id2 := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret2",
	}
	g := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id2,
			},
		},
	}
	anyg, err := anypb.New(g)
	require.NoError(t, err)
	err = newTestWriteStore(t, kv).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: vaulttypes.KeyFor(id), Item: anyd},
			{Id: vaulttypes.KeyFor(id2), Item: anyg},
		},
	)
	require.NoError(t, err)

	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id2,
			},
		},
	}

	obsb := marshalObservations(t, observation{id, g, resp})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	require.ErrorContains(t, err, "number of observations doesn't match number of pending requests")

	resp2 := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id2,
				Success: true,
			},
		},
	}
	obsb = marshalObservations(t, observation{id, g, resp}, observation{id2, d, resp2})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	require.NoError(t, err)
}

func TestPlugin_ValidateObservations_DisallowsDuplicateBlobHandles(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
		marshalBlob:   mockMarshalBlob,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	obs := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0: 1},
			{0: 2},
		},
	}
	obsb, err := proto.Marshal(obs)
	require.NoError(t, err)

	bf := &blobber{
		blobs: [][]byte{
			{0: 1},
			{0: 1},
		},
	}
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		bf,
	)
	require.ErrorContains(t, err, "duplicate item found in pending queue item observation")
}

func TestPlugin_StateTransition_ShasDontMatch(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp1 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}
	resp2 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "something else",
				},
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp1}, observation{id, req, resp2}, observation{id, req, resp1})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Empty(t, os.Outcomes)

	assert.Equal(t, 1, observed.FilterMessage("insufficient observations found for id").Len())
}

func TestPlugin_StateTransition_AggregatesValidationErrors(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetGetSecretsRequest()))
	assert.True(t, proto.Equal(resp, o.GetGetSecretsResponse()))

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_GetSecretsRequest_CombinesShares(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp1 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-1"},
							},
						},
					},
				},
			},
		},
	}
	resp2 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-2"},
							},
						},
					},
				},
			},
		},
	}
	resp3 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-3"},
							},
						},
					},
				},
			},
		},
	}

	obsb1 := marshalObservations(t, observation{id, req, resp1})
	obsb2 := marshalObservations(t, observation{id, req, resp2})
	obsb3 := marshalObservations(t, observation{id, req, resp3})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb1)},
			{Observer: 1, Observation: types.Observation(obsb2)},
			{Observer: 2, Observation: types.Observation(obsb3)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetGetSecretsRequest()))

	expectedResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-1", "encrypted-share-2", "encrypted-share-3"},
							},
						},
					},
				},
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetGetSecretsResponse()), o.GetGetSecretsResponse())

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_CreateSecretsRequest_WritesSecrets(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	value := []byte("encrypted-value")
	enc := hex.EncodeToString(value)
	req := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: enc,
			},
		},
	}
	resp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetCreateSecretsRequest()))

	expectedResp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetCreateSecretsResponse()), o.GetCreateSecretsResponse())

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)

	assert.Equal(t, ss.EncryptedSecret, []byte("encrypted-value"))

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_CreateSecretsRequest_UsesWorkflowOwnerMetadataWhenGateEnabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	cfg := makeReportingPluginConfig(
		t,
		10,
		nil,
		nil,
		1,
		1024,
		100,
		100,
		100,
		10,
	)
	cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg:     cfg,
	}

	const (
		orgID         = "org-create"
		workflowOwner = "0x2222222222222222222222222222222222222222"
	)

	kv := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, kv).WriteSecret(t.Context(), &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "legacy",
	}, &vaultcommon.StoredSecret{EncryptedSecret: []byte("legacy-value")}))
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "new-secret",
	}
	req := &vaultcommon.CreateSecretsRequest{
		RequestId: "request-id",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString([]byte("encrypted-value")),
			},
		},
		OrgId:         orgID,
		WorkflowOwner: workflowOwner,
	}
	resp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		},
		kv,
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetCreateSecretsRequest()), o.GetCreateSecretsRequest())
	require.Len(t, o.GetCreateSecretsResponse().Responses, 1)
	assert.False(t, o.GetCreateSecretsResponse().Responses[0].Success)
	assert.Contains(t, o.GetCreateSecretsResponse().Responses[0].Error, "has reached maximum number of secrets")

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	assert.Nil(t, ss)
}

func TestPlugin_Reports(t *testing.T) {
	value := "encrypted-value"
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: value,
			},
		},
	}
	resp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	expectedOutcome1 := &vaultcommon.Outcome{
		Id:          vaulttypes.KeyFor(id),
		RequestType: vaultcommon.RequestType_CREATE_SECRETS,
		Request: &vaultcommon.Outcome_CreateSecretsRequest{
			CreateSecretsRequest: req,
		},
		Response: &vaultcommon.Outcome_CreateSecretsResponse{
			CreateSecretsResponse: resp,
		},
	}

	id2 := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret2",
	}
	req2 := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id2,
			},
		},
	}
	resp2 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id:     id2,
				Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{EncryptedValue: value}},
			},
		},
	}
	expectedOutcome2 := &vaultcommon.Outcome{
		Id:          vaulttypes.KeyFor(id2),
		RequestType: vaultcommon.RequestType_GET_SECRETS,
		Request: &vaultcommon.Outcome_GetSecretsRequest{
			GetSecretsRequest: req2,
		},
		Response: &vaultcommon.Outcome_GetSecretsResponse{
			GetSecretsResponse: resp2,
		},
	}
	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{
			expectedOutcome1,
			expectedOutcome2,
		},
	}

	osb, err := proto.Marshal(os)
	require.NoError(t, err)

	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	rs, err := r.Reports(t.Context(), uint64(1), osb)
	require.NoError(t, err)

	assert.Len(t, rs, 2)

	o1 := rs[0]
	info1, err := extractReportInfo(o1.ReportWithInfo)
	require.NoError(t, err)

	assert.True(t, proto.Equal(&vaultcommon.ReportInfo{
		Id:          vaulttypes.KeyFor(id),
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_JSON,
		RequestType: vaultcommon.RequestType_CREATE_SECRETS,
	}, info1))

	expectedBytes, err := vaultutils.ToCanonicalJSON(resp)
	require.NoError(t, err)
	assert.Equal(t, expectedBytes, []byte(o1.ReportWithInfo.Report))

	o2 := rs[1]
	info2, err := extractReportInfo(o2.ReportWithInfo)
	require.NoError(t, err)
	assert.True(t, proto.Equal(&vaultcommon.ReportInfo{
		Id:          vaulttypes.KeyFor(id2),
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_PROTOBUF,
		RequestType: vaultcommon.RequestType_GET_SECRETS,
	}, info2))

	o2r := &vaultcommon.GetSecretsResponse{}
	err = proto.Unmarshal(o2.ReportWithInfo.Report, o2r)
	require.NoError(t, err)
	assert.True(t, proto.Equal(resp2, o2r))
}

func TestPlugin_Observation_UpdateSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name     string
		id       *vaultcommon.SecretIdentifier
		maxIDLen int
		err      string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "invalid secret identifier: cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "invalid secret identifier: key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		lggr := logger.TestLogger(t)
		store := requests.NewStore[*vaulttypes.Request]()
		maxIDLen := 256
		if tc.maxIDLen > 0 {
			maxIDLen = tc.maxIDLen
		}
		r := &ReportingPlugin{
			lggr:          lggr,
			store:         store,
			metrics:       newTestMetrics(t),
			marshalBlob:   mockMarshalBlob,
			unmarshalBlob: mockUnmarshalBlob,
			cfg: makeReportingPluginConfig(
				t,
				10,
				nil,
				nil,
				1,
				1024,
				maxIDLen/3,
				maxIDLen/3,
				maxIDLen/3,
				10,
			),
		}

		seqNr := uint64(1)
		rdr := &kv{
			m: make(map[string]response),
		}
		p := &vaultcommon.UpdateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id:             tc.id,
					EncryptedValue: "foo",
				},
			},
		}
		anyp, err := anypb.New(p)
		require.NoError(t, err)
		err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
			[]*vaultcommon.StoredPendingQueueItem{
				{Id: "request-1", Item: anyp},
			},
		)
		require.NoError(t, err)
		bf := &blobber{}
		data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		err = proto.Unmarshal(data, obs)
		require.NoError(t, err)

		assert.Len(t, obs.Observations, 1)
		o := obs.Observations[0]

		assert.Equal(t, vaultcommon.RequestType_UPDATE_SECRETS, o.RequestType)
		assert.True(t, proto.Equal(o.GetUpdateSecretsRequest(), p))

		batchResp := o.GetUpdateSecretsResponse()
		assert.Len(t, p.EncryptedSecrets, 1)
		assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

		assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_UpdateSecretsRequest_DisallowsDuplicateRequests(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	p := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: "foo",
			},
			{
				Id:             id,
				EncryptedValue: "bla",
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_UPDATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetUpdateSecretsRequest(), p))

	batchResp := o.GetUpdateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 2)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "duplicate request for secret identifier")

	assert.True(t, proto.Equal(p.EncryptedSecrets[1].Id, batchResp.Responses[1].Id))
	resp = batchResp.Responses[1]
	assert.Contains(t, resp.GetError(), "duplicate request for secret identifier")
}

func TestPlugin_Observation_UpdateSecretsRequest_InvalidCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	p := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: "foo",
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_UPDATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetUpdateSecretsRequest(), p))

	batchResp := o.GetUpdateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "invalid hex encoding for ciphertext")
}

func TestPlugin_Observation_UpdateSecretsRequest_InvalidCiphertext_TooLong(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			10,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ciphertext := []byte("a quick brown fox jumps over the lazy dog")
	p := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertext),
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_UPDATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetUpdateSecretsRequest(), p))

	batchResp := o.GetUpdateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "ciphertext size exceeds maximum allowed size: 10b")
}

func TestPlugin_Observation_UpdateSecretsRequest_InvalidCiphertext_EncryptedWithWrongPublicKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	// Wrong key
	_, wrongPublicKey, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	// Right key
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ct, err := tdh2easy.Encrypt(wrongPublicKey, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString(ciphertextBytes),
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_UPDATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetUpdateSecretsRequest(), p))

	batchResp := o.GetUpdateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "failed to verify ciphertext")
}

func TestPlugin_StateTransition_UpdateSecretsRequest_SecretDoesntExist(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	value := []byte("encrypted-value")
	enc := hex.EncodeToString(value)
	req := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: enc,
			},
		},
	}
	resp := &vaultcommon.UpdateSecretsResponse{
		Responses: []*vaultcommon.UpdateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetUpdateSecretsRequest()))

	expectedResp := &vaultcommon.UpdateSecretsResponse{
		Responses: []*vaultcommon.UpdateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "could not write update to key value store: key does not exist",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetUpdateSecretsResponse()), o.GetUpdateSecretsResponse())

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_UpdateSecretsRequest_WritesSecrets(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	secret, err := proto.Marshal(&vaultcommon.StoredSecret{
		EncryptedSecret: []byte("old-encrypted-value"),
	})
	require.NoError(t, err)
	metadata, err := proto.Marshal(&vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{id},
	})
	require.NoError(t, err)
	kv := &kv{
		m: map[string]response{
			keyPrefix + vaulttypes.KeyFor(id): {
				data: secret,
			},
			metadataPrefix + "owner": {
				data: metadata,
			},
		},
	}
	rs := newTestReadStore(t, kv)

	value := []byte("encrypted-value")
	enc := hex.EncodeToString(value)
	req := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: enc,
			},
		},
	}
	resp := &vaultcommon.UpdateSecretsResponse{
		Responses: []*vaultcommon.UpdateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	seqNr := uint64(1)
	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetUpdateSecretsRequest()))

	expectedResp := &vaultcommon.UpdateSecretsResponse{
		Responses: []*vaultcommon.UpdateSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetUpdateSecretsResponse()), o.GetUpdateSecretsResponse())

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.NotNil(t, ss)

	assert.Equal(t, ss.EncryptedSecret, []byte("encrypted-value"))

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_UpdateSecretsRequest_MigratesWorkflowOwnerSecretWhenGateEnabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	cfg := makeReportingPluginConfig(
		t,
		10,
		nil,
		nil,
		5,
		1024,
		100,
		100,
		100,
		10,
	)
	cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg:     cfg,
	}

	const (
		orgID         = "org-update"
		workflowOwner = "0x3333333333333333333333333333333333333333"
	)

	kv := &kv{m: make(map[string]response)}
	legacyID := &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "secret",
	}
	require.NoError(t, newTestWriteStore(t, kv).WriteSecret(t.Context(), legacyID, &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("old-encrypted-value"),
	}))
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.UpdateSecretsRequest{
		RequestId: "request-id",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: hex.EncodeToString([]byte("encrypted-value")),
			},
		},
		OrgId:         orgID,
		WorkflowOwner: workflowOwner,
	}
	resp := &vaultcommon.UpdateSecretsResponse{
		Responses: []*vaultcommon.UpdateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		},
		kv,
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetUpdateSecretsRequest()), o.GetUpdateSecretsRequest())
	require.Len(t, o.GetUpdateSecretsResponse().Responses, 1)
	assert.True(t, o.GetUpdateSecretsResponse().Responses[0].Success)

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.NotNil(t, ss)
	assert.Equal(t, []byte("encrypted-value"), ss.EncryptedSecret)

	legacySecret, err := rs.GetSecret(t.Context(), legacyID)
	require.NoError(t, err)
	assert.Nil(t, legacySecret)
}

func TestPlugin_Reports_UpdateSecretsRequest(t *testing.T) {
	value := "encrypted-value"
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: value,
			},
		},
	}
	resp := &vaultcommon.UpdateSecretsResponse{
		Responses: []*vaultcommon.UpdateSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	expectedOutcome := &vaultcommon.Outcome{
		Id:          vaulttypes.KeyFor(id),
		RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
		Request: &vaultcommon.Outcome_UpdateSecretsRequest{
			UpdateSecretsRequest: req,
		},
		Response: &vaultcommon.Outcome_UpdateSecretsResponse{
			UpdateSecretsResponse: resp,
		},
	}

	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{
			expectedOutcome,
		},
	}

	osb, err := proto.Marshal(os)
	require.NoError(t, err)

	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	rs, err := r.Reports(t.Context(), uint64(1), osb)
	require.NoError(t, err)

	assert.Len(t, rs, 1)

	o := rs[0]
	info1, err := extractReportInfo(o.ReportWithInfo)
	require.NoError(t, err)

	assert.True(t, proto.Equal(&vaultcommon.ReportInfo{
		Id:          vaulttypes.KeyFor(id),
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_JSON,
		RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
	}, info1))

	expectedBytes, err := vaultutils.ToCanonicalJSON(resp)
	require.NoError(t, err)
	assert.Equal(t, expectedBytes, []byte(o.ReportWithInfo.Report))
}

func TestPlugin_Observation_DeleteSecrets(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "foo",
		Namespace: "main",
		Key:       "item4",
	}
	md := &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{
			id,
		},
	}
	mdb, err := proto.Marshal(md)
	require.NoError(t, err)

	ss := &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("encrypted-value"),
	}
	ssb, err := proto.Marshal(ss)
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{
			metadataPrefix + "foo": response{
				data: mdb,
			},
			keyPrefix + vaulttypes.KeyFor(id): response{
				data: ssb,
			},
		},
	}
	p := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-id",
		Ids: []*vaultcommon.SecretIdentifier{
			id,
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_DELETE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetDeleteSecretsRequest(), p))

	resp := o.GetDeleteSecretsResponse()
	assert.Len(t, resp.Responses, 1)
	assert.True(t, proto.Equal(id, resp.Responses[0].Id))
	assert.False(t, resp.Responses[0].Success, resp.Responses[0].GetError()) // false because it hasn't actually been deleted yet.
	assert.Empty(t, resp.Responses[0].GetError())
}

func TestPlugin_Observation_DeleteSecrets_IdDoesntExist(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{},
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "foo",
		Namespace: "main",
		Key:       "item4",
	}
	p := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-id",
		Ids: []*vaultcommon.SecretIdentifier{
			id,
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_DELETE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetDeleteSecretsRequest(), p))

	resp := o.GetDeleteSecretsResponse()
	assert.Len(t, resp.Responses, 1)
	assert.True(t, proto.Equal(id, resp.Responses[0].Id))
	assert.False(t, resp.Responses[0].Success, resp.Responses[0].GetError())
	assert.Contains(t, resp.Responses[0].GetError(), "key does not exist")
}

func TestPlugin_Observation_DeleteSecrets_InvalidRequestDuplicateIds(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{},
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "foo",
		Namespace: "main",
		Key:       "item4",
	}
	p := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-id",
		Ids: []*vaultcommon.SecretIdentifier{
			id,
			id,
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_DELETE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetDeleteSecretsRequest(), p))

	resp := o.GetDeleteSecretsResponse()
	assert.Len(t, resp.Responses, 2)
	assert.True(t, proto.Equal(id, resp.Responses[0].Id))
	assert.False(t, resp.Responses[0].Success, resp.Responses[0].GetError())
	assert.Contains(t, resp.Responses[0].GetError(), "duplicate request for secret identifier")

	assert.True(t, proto.Equal(id, resp.Responses[1].Id))
	assert.False(t, resp.Responses[1].Success, resp.Responses[1].GetError())
	assert.Contains(t, resp.Responses[1].GetError(), "duplicate request for secret identifier")
}

func TestPlugin_StateTransition_DeleteSecretsRequest(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "foo",
		Namespace: "main",
		Key:       "item4",
	}
	md := &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{
			id,
		},
	}
	mdb, err := proto.Marshal(md)
	require.NoError(t, err)

	ss := &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("encrypted-value"),
	}
	ssb, err := proto.Marshal(ss)
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{
			metadataPrefix + "foo": response{
				data: mdb,
			},
			keyPrefix + vaulttypes.KeyFor(id): response{
				data: ssb,
			},
		},
	}
	rs := newTestReadStore(t, rdr)

	req := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-id",
		Ids:       []*vaultcommon.SecretIdentifier{id},
	}
	resp := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, rdr, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetDeleteSecretsRequest()), o.GetDeleteSecretsRequest())
	expectedResp := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetDeleteSecretsResponse()))

	ss, err = rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_DeleteSecretsRequest_DeletesWorkflowOwnerSecretWhenGateEnabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	cfg := makeReportingPluginConfig(
		t,
		10,
		nil,
		nil,
		5,
		1024,
		100,
		100,
		100,
		10,
	)
	cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg:     cfg,
	}

	const (
		orgID         = "org-delete"
		workflowOwner = "0x4444444444444444444444444444444444444444"
	)

	kv := &kv{m: make(map[string]response)}
	legacyID := &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "item4",
	}
	require.NoError(t, newTestWriteStore(t, kv).WriteSecret(t.Context(), legacyID, &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("encrypted-value"),
	}))
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: "main",
		Key:       "item4",
	}
	req := &vaultcommon.DeleteSecretsRequest{
		RequestId:     "request-id",
		Ids:           []*vaultcommon.SecretIdentifier{id},
		OrgId:         orgID,
		WorkflowOwner: workflowOwner,
	}
	resp := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		},
		kv,
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetDeleteSecretsRequest()), o.GetDeleteSecretsRequest())
	require.Len(t, o.GetDeleteSecretsResponse().Responses, 1)
	assert.True(t, o.GetDeleteSecretsResponse().Responses[0].Success)

	ss, err := rs.GetSecret(t.Context(), legacyID)
	require.NoError(t, err)
	assert.Nil(t, ss)
}

func TestPlugin_StateTransition_DeleteSecretsRequest_SecretDoesNotExist(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "foo",
		Namespace: "main",
		Key:       "item4",
	}
	md := &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{},
	}
	mdb, err := proto.Marshal(md)
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{
			metadataPrefix + "foo": response{
				data: mdb,
			},
		},
	}
	rs := newTestReadStore(t, rdr)

	req := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-id",
		Ids:       []*vaultcommon.SecretIdentifier{id},
	}
	resp := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, rdr, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetDeleteSecretsRequest()), o.GetDeleteSecretsRequest())
	expectedResp := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "failed to handle delete secret request",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetDeleteSecretsResponse()), o.GetDeleteSecretsResponse())

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_Reports_DeleteSecretsRequest(t *testing.T) {
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-id",
		Ids:       []*vaultcommon.SecretIdentifier{id},
	}
	resp := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	expectedOutcome := &vaultcommon.Outcome{
		Id:          vaulttypes.KeyFor(id),
		RequestType: vaultcommon.RequestType_DELETE_SECRETS,
		Request: &vaultcommon.Outcome_DeleteSecretsRequest{
			DeleteSecretsRequest: req,
		},
		Response: &vaultcommon.Outcome_DeleteSecretsResponse{
			DeleteSecretsResponse: resp,
		},
	}

	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{
			expectedOutcome,
		},
	}

	osb, err := proto.Marshal(os)
	require.NoError(t, err)

	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	rs, err := r.Reports(t.Context(), uint64(1), osb)
	require.NoError(t, err)

	assert.Len(t, rs, 1)

	o := rs[0]
	info1, err := extractReportInfo(o.ReportWithInfo)
	require.NoError(t, err)

	assert.True(t, proto.Equal(&vaultcommon.ReportInfo{
		Id:          vaulttypes.KeyFor(id),
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_JSON,
		RequestType: vaultcommon.RequestType_DELETE_SECRETS,
	}, info1))

	expectedBytes, err := vaultutils.ToCanonicalJSON(resp)
	require.NoError(t, err)
	assert.Equal(t, expectedBytes, []byte(o.ReportWithInfo.Report))
}

func TestPlugin_Observation_ListSecretIdentifiers_OwnerRequired(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	p := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-id",
		Owner:     "",
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS, o.RequestType)
	assert.True(t, proto.Equal(o.GetListSecretIdentifiersRequest(), p))

	resp := o.GetListSecretIdentifiersResponse()
	assert.Empty(t, resp.Identifiers)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.GetError(), "owner cannot be empty")
}

func TestPlugin_Observation_ListSecretIdentifiers_NoNamespaceProvided(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	md := &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{
			{
				Owner:     "foo",
				Namespace: "main",
				Key:       "item4",
			},
			{
				Owner:     "foo",
				Namespace: "secondary",
				Key:       "item2",
			},
			{
				Owner:     "foo",
				Namespace: "main",
				Key:       "item3",
			},
		},
	}
	mdb, err := proto.Marshal(md)
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{
			metadataPrefix + "foo": response{
				data: mdb,
			},
		},
	}
	p := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-id",
		Owner:     "foo",
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS, o.RequestType)
	assert.True(t, proto.Equal(o.GetListSecretIdentifiersRequest(), p))

	resp := o.GetListSecretIdentifiersResponse()
	expectedIdentifiers := []*vaultcommon.SecretIdentifier{
		{
			Owner:     "foo",
			Namespace: "main",
			Key:       "item3",
		},
		{
			Owner:     "foo",
			Namespace: "main",
			Key:       "item4",
		},
		{
			Owner:     "foo",
			Namespace: "secondary",
			Key:       "item2",
		},
	}
	for i, id := range resp.Identifiers {
		assert.True(t, proto.Equal(expectedIdentifiers[i], id))
	}
	assert.Len(t, resp.Identifiers, 3)
	assert.True(t, resp.Success)
	assert.Empty(t, resp.GetError())
}

func TestPlugin_Observation_ListSecretIdentifiers_FilterByNamespace(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	md := &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{
			{
				Owner:     "foo",
				Namespace: "main",
				Key:       "item4",
			},
			{
				Owner:     "foo",
				Namespace: "secondary",
				Key:       "item2",
			},
			{
				Owner:     "foo",
				Namespace: "main",
				Key:       "item3",
			},
		},
	}
	mdb, err := proto.Marshal(md)
	require.NoError(t, err)

	seqNr := uint64(1)
	rdr := &kv{
		m: map[string]response{
			metadataPrefix + "foo": response{
				data: mdb,
			},
		},
	}
	p := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-id",
		Owner:     "foo",
		Namespace: "main",
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS, o.RequestType)
	assert.True(t, proto.Equal(o.GetListSecretIdentifiersRequest(), p))

	resp := o.GetListSecretIdentifiersResponse()
	expectedIdentifiers := []*vaultcommon.SecretIdentifier{
		{
			Owner:     "foo",
			Namespace: "main",
			Key:       "item3",
		},
		{
			Owner:     "foo",
			Namespace: "main",
			Key:       "item4",
		},
	}
	for i, id := range resp.Identifiers {
		assert.True(t, proto.Equal(expectedIdentifiers[i], id))
	}
	assert.Len(t, resp.Identifiers, 2)
	assert.True(t, resp.Success)
	assert.Empty(t, resp.GetError())
}

func TestPlugin_Observation_ListSecretIdentifiers_FallsBackToWorkflowOwnerWhenGateEnabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	cfg := makeReportingPluginConfig(
		t,
		10,
		nil,
		nil,
		5,
		1024,
		30,
		30,
		30,
		10,
	)
	cfg.OrgIDAsSecretOwnerEnabled = limits.NewGateLimiter(true)
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg:           cfg,
	}

	const (
		orgID         = "org-list"
		workflowOwner = "0x1111111111111111111111111111111111111111"
	)

	rdr := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, rdr).WriteSecret(t.Context(), &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "legacy-item",
	}, &vaultcommon.StoredSecret{EncryptedSecret: []byte("legacy-data")}))

	p := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-id",
		Owner:         orgID,
		OrgId:         orgID,
		WorkflowOwner: workflowOwner,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	))

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	require.Len(t, obs.Observations, 1)

	resp := obs.Observations[0].GetListSecretIdentifiersResponse()
	require.True(t, resp.Success, resp.GetError())
	require.Len(t, resp.Identifiers, 1)
	assert.Equal(t, orgID, resp.Identifiers[0].Owner)
	assert.Equal(t, "legacy-item", resp.Identifiers[0].Key)
}

func TestPlugin_Observation_ListSecretIdentifiers_DoesNotFallbackWhenGateDisabled(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	r := &ReportingPlugin{
		lggr:          lggr,
		store:         store,
		metrics:       newTestMetrics(t),
		marshalBlob:   mockMarshalBlob,
		unmarshalBlob: mockUnmarshalBlob,
		cfg: makeReportingPluginConfig(
			t,
			10,
			nil,
			nil,
			5,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	const (
		orgID         = "org-list"
		workflowOwner = "0x1111111111111111111111111111111111111111"
	)

	rdr := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, rdr).WriteSecret(t.Context(), &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "legacy-item",
	}, &vaultcommon.StoredSecret{EncryptedSecret: []byte("legacy-data")}))

	p := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-id",
		Owner:         orgID,
		OrgId:         orgID,
		WorkflowOwner: workflowOwner,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	))

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	require.Len(t, obs.Observations, 1)

	resp := obs.Observations[0].GetListSecretIdentifiersResponse()
	require.True(t, resp.Success, resp.GetError())
	assert.Empty(t, resp.Identifiers)
}

func TestPlugin_Reports_ListSecretIdentifiersRequest(t *testing.T) {
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-id",
		Owner:     "owner",
	}
	resp := &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: []*vaultcommon.SecretIdentifier{
			id,
		},
	}
	expectedOutcome := &vaultcommon.Outcome{
		Id:          vaulttypes.KeyFor(id),
		RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
		Request: &vaultcommon.Outcome_ListSecretIdentifiersRequest{
			ListSecretIdentifiersRequest: req,
		},
		Response: &vaultcommon.Outcome_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: resp,
		},
	}

	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{
			expectedOutcome,
		},
	}

	osb, err := proto.Marshal(os)
	require.NoError(t, err)

	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	rs, err := r.Reports(t.Context(), uint64(1), osb)
	require.NoError(t, err)

	assert.Len(t, rs, 1)

	o := rs[0]
	info1, err := extractReportInfo(o.ReportWithInfo)
	require.NoError(t, err)

	assert.True(t, proto.Equal(&vaultcommon.ReportInfo{
		Id:          vaulttypes.KeyFor(id),
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_JSON,
		RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
	}, info1))

	expectedBytes, err := vaultutils.ToCanonicalJSON(resp)
	require.NoError(t, err)
	assert.Equal(t, expectedBytes, []byte(o.ReportWithInfo.Report))
}

func TestPlugin_StateTransition_ListSecretIdentifiers(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	resp := &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: []*vaultcommon.SecretIdentifier{id},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetListSecretIdentifiersRequest()))

	assert.True(t, proto.Equal(resp, o.GetListSecretIdentifiersResponse()))

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func protoMarshal(t *testing.T, msg proto.Message) []byte {
	t.Helper()
	b, err := proto.Marshal(msg)
	require.NoError(t, err)
	return b
}

func mockUnmarshalBlob(data []byte) (ocr3_1types.BlobHandle, error) {
	return ocr3_1types.BlobHandle{}, nil
}

func mockMarshalBlob(ocr3_1types.BlobHandle) ([]byte, error) {
	return []byte{}, nil
}

type callbackBlobFetcher struct {
	fn func(payload []byte) error
}

func (f *callbackBlobFetcher) BroadcastBlob(_ context.Context, payload []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	if err := f.fn(payload); err != nil {
		return ocr3_1types.BlobHandle{}, err
	}
	return ocr3_1types.BlobHandle{}, nil
}

func (f *callbackBlobFetcher) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	panic("FetchBlob should not be called in broadcastBlobPayloads tests")
}

type ctxCallbackBlobFetcher struct {
	fn func(ctx context.Context, payload []byte) error
}

func (f *ctxCallbackBlobFetcher) BroadcastBlob(ctx context.Context, payload []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	if err := f.fn(ctx, payload); err != nil {
		return ocr3_1types.BlobHandle{}, err
	}
	return ocr3_1types.BlobHandle{}, nil
}

func (f *ctxCallbackBlobFetcher) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	panic("FetchBlob should not be called in broadcastBlobPayloads tests")
}

func TestPlugin_StateTransition_StoresPendingQueue(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			5,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	req := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	req1, err := anypb.New(req)
	require.NoError(t, err)

	req = &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "production",
		RequestId: "request-id2",
	}
	req2, err := anypb.New(req)
	require.NoError(t, err)

	req = &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "staging",
		RequestId: "request-id3",
	}
	req3, err := anypb.New(req)
	require.NoError(t, err)

	req = &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "testnet",
		RequestId: "request-id4",
	}
	req4, err := anypb.New(req)
	require.NoError(t, err)

	bf := &blobber{
		blobs: [][]byte{
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item1",
				Item: req1,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item2",
				Item: req2,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item3",
				Item: req3,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item2",
				Item: req2,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item1",
				Item: req1,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item3",
				Item: req3,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "item4",
				Item: req4,
			}),
		},
	}

	r.unmarshalBlob = bf.unmarshalBlob

	o1 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0}, // maps to item 0 in the blobs
			{1}, // maps to item 1 in the blobs
			{2}, // maps to item 2 in the blobs
		},
	}
	o1b, err := proto.Marshal(o1)
	require.NoError(t, err)

	o2 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{3}, // maps to item 3 in the blobs
		},
	}
	o2b, err := proto.Marshal(o2)
	require.NoError(t, err)

	o3 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{4}, // maps to item 4 in the blobs
			{5}, // maps to item 5 in the blobs
			{6}, // maps to item 6 in the blobs
		},
	}
	o3b, err := proto.Marshal(o3)
	require.NoError(t, err)

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: o1b},
			{Observer: 1, Observation: o2b},
			{Observer: 2, Observation: o3b},
		},
		rdr,
		bf,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Empty(t, os.Outcomes)

	pq, err := newTestReadStore(t, rdr).GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Len(t, pq, 3)

	ids := []string{}
	for _, item := range pq {
		ids = append(ids, item.Id)
	}

	assert.ElementsMatch(t, []string{"item1", "item2", "item3"}, ids)
}

func TestPlugin_StateTransition_StoresPendingQueue_LimitedToBatchSize(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	req1 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	areq1, err := anypb.New(req1)
	require.NoError(t, err)

	req2 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "production",
		RequestId: "request-id2",
	}
	areq2, err := anypb.New(req2)
	require.NoError(t, err)

	req3 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "staging",
		RequestId: "request-id3",
	}
	areq3, err := anypb.New(req3)
	require.NoError(t, err)

	req4 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "testnet",
		RequestId: "request-id4",
	}
	areq4, err := anypb.New(req4)
	require.NoError(t, err)

	bf := &blobber{
		blobs: [][]byte{
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id",
				Item: areq1,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id2",
				Item: areq2,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id3",
				Item: areq3,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id4",
				Item: areq4,
			}),
		},
	}

	r.unmarshalBlob = bf.unmarshalBlob

	o1 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0}, // maps to item 0 in the blobs
			{1}, // maps to item 1 in the blobs
			{2}, // maps to item 2 in the blobs
			{3}, // maps to item 3 in the blobs
		},
	}
	o1b, err := proto.Marshal(o1)
	require.NoError(t, err)

	o2 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0}, // maps to item 0 in the blobs
			{1}, // maps to item 1 in the blobs
			{2}, // maps to item 2 in the blobs
			{3}, // maps to item 3 in the blobs
		},
	}
	o2b, err := proto.Marshal(o2)
	require.NoError(t, err)

	o3 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0}, // maps to item 0 in the blobs
			{1}, // maps to item 1 in the blobs
			{2}, // maps to item 2 in the blobs
			{3}, // maps to item 3 in the blobs
		},
	}
	o3b, err := proto.Marshal(o3)
	require.NoError(t, err)

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: o1b},
			{Observer: 1, Observation: o2b},
			{Observer: 2, Observation: o3b},
		},
		rdr,
		bf,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Empty(t, os.Outcomes)

	pq, err := newTestReadStore(t, rdr).GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Len(t, pq, 1)

	ids := []string{}
	for _, item := range pq {
		ids = append(ids, item.Id)
	}

	// Batch size is 1, so only one item should be stored.
	assert.ElementsMatch(t, []string{"request-id"}, ids)
}

func TestPlugin_StateTransition_StoresPendingQueue_DoesntDoubleCountObservationsFromOneNode(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	req1 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	areq1, err := anypb.New(req1)
	require.NoError(t, err)

	bf := &blobber{
		blobs: [][]byte{
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id",
				Item: areq1,
			}),
		},
	}

	r.unmarshalBlob = bf.unmarshalBlob

	o1 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0}, // maps to item 0 in the blobs
			{0}, // maps to item 0 in the blobs (duplicate)
			{0}, // maps to item 0 in the blobs (duplicate)
		},
	}
	o1b, err := proto.Marshal(o1)
	require.NoError(t, err)

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: o1b},
		},
		rdr,
		bf,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Empty(t, os.Outcomes)

	pq, err := newTestReadStore(t, rdr).GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Empty(t, pq, 0)

	ids := []string{}
	for _, item := range pq {
		ids = append(ids, item.Id)
	}

	// 1 oracle submitted duplicates, so skipping
	assert.ElementsMatch(t, []string{}, ids)
}

func TestPlugin_ValidateObservation_RejectsIfMoreThan2xBatchSize(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	req1 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	areq1, err := anypb.New(req1)
	require.NoError(t, err)

	o1 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{}, // maps to item 0 in the blobs
			{}, // maps to item 1 in the blobs
			{}, // maps to item 2 in the blobs
			{}, // maps to item 3 in the blobs
		},
	}

	o1b, err := proto.Marshal(o1)
	require.NoError(t, err)

	bf := &blobber{
		blobs: [][]byte{
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id",
				Item: areq1,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id",
				Item: areq1,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id",
				Item: areq1,
			}),
			protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
				Id:   "request-id",
				Item: areq1,
			}),
		},
	}

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.ErrorContains(t, err, "invalid observation: too many pending queue items provided, have 4, want max 2")
}

// TestPlugin_ValidateObservation_AcceptsFullPendingQueueObservation verifies that an observation
// with exactly 2*batchSize pending queue items (the maximum Observation can produce) is accepted.
func TestPlugin_ValidateObservation_AcceptsFullPendingQueueObservation(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	batchSize := 1 // MaxBatchSize=1, so 2*batchSize=2 is the intended max pending queue items
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			batchSize,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	req1 := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	areq1, err := anypb.New(req1)
	require.NoError(t, err)

	// Build an observation with exactly 2*batchSize = 2 pending queue items.
	// This is the maximum that Observation() can produce.
	numItems := 2 * batchSize
	pendingQueueItems := make([][]byte, numItems)
	blobs := make([][]byte, numItems)
	for i := range numItems {
		pendingQueueItems[i] = []byte{}
		blobs[i] = protoMarshal(t, &vaultcommon.StoredPendingQueueItem{
			Id:   fmt.Sprintf("request-id-%d", i),
			Item: areq1,
		})
	}

	o1 := &vaultcommon.Observations{
		PendingQueueItems: pendingQueueItems,
	}

	o1b, err := proto.Marshal(o1)
	require.NoError(t, err)

	bf := &blobber{
		blobs: blobs,
	}

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.NoError(t, err)
}

func TestPlugin_ValidateObservation_GetSecretsRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ek, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(ek[:])
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-1", "encrypted-share-2"},
							},
						},
					},
				},
			},
		},
	}
	anyp, err := anypb.New(req)
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	bf := &blobber{}

	o1 := &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: req,
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b := protoMarshal(t, o1)

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.ErrorContains(t, err, "invalid observation: observation must have exactly 1 share per encryption key")

	resp = &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "foo",
				},
			},
		},
	}

	o1 = &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: req,
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b = protoMarshal(t, o1)

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.NoError(t, err)

	resp = &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "foo",
				},
			},
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "foo",
				},
			},
		},
	}

	o1 = &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: req,
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b = protoMarshal(t, o1)

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.ErrorContains(t, err, "invalid observation: GetSecrets request and response must have the same number of items")

	resp = &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue:               "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{},
					},
				},
			},
		},
	}

	o1 = &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: req,
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b = protoMarshal(t, o1)

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.ErrorContains(t, err, "invalid observation: observation must contain a share per encryption key provided")

	resp = &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								Shares: []string{strings.Repeat("1", 1000)},
							},
						},
					},
				},
			},
		},
	}

	o1 = &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: req,
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b = protoMarshal(t, o1)

	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.ErrorContains(t, err, "invalid observation: share provided exceeds maximum size allowed")
}

func TestPlugin_ValidateObservation_PanicsOnEmptyShares(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	rdr := &kv{m: make(map[string]response)}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ek, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(ek[:])

	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	// Malicious observation: EncryptedShares with an empty Shares slice.
	// This triggers an index-out-of-bounds panic at ds.Shares[0] in validateObservation.
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: pks,
								Shares:        []string{}, // empty — triggers panic
							},
						},
					},
				},
			},
		},
	}

	anyp, err := anypb.New(req)
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	bf := &blobber{}

	o1 := &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: req,
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b := protoMarshal(t, o1)

	// This should return an error, not panic.
	require.NotPanics(t, func() {
		err = r.ValidateObservation(
			t.Context(),
			seqNr,
			types.AttributedQuery{},
			types.AttributedObservation{
				Observer: 0, Observation: o1b,
			},
			rdr,
			bf,
		)
	})
	require.Error(t, err)
}

func TestPlugin_ValidateObservation_NilSecretIdentifier(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	bf := &blobber{}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}

	tests := []struct {
		name string
		obs  *vaultcommon.Observation
	}{
		{
			name: "GetSecrets request with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: &vaultcommon.GetSecretsRequest{
						Requests: []*vaultcommon.SecretRequest{
							{Id: nil},
						},
					},
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: &vaultcommon.GetSecretsResponse{
						Responses: []*vaultcommon.SecretResponse{
							{Id: id, Result: &vaultcommon.SecretResponse_Error{Error: "err"}},
						},
					},
				},
			},
		},
		{
			name: "GetSecrets response with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Request: &vaultcommon.Observation_GetSecretsRequest{
					GetSecretsRequest: &vaultcommon.GetSecretsRequest{
						Requests: []*vaultcommon.SecretRequest{
							{Id: id},
						},
					},
				},
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: &vaultcommon.GetSecretsResponse{
						Responses: []*vaultcommon.SecretResponse{
							{Id: nil, Result: &vaultcommon.SecretResponse_Error{Error: "err"}},
						},
					},
				},
			},
		},
		{
			name: "CreateSecrets with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_CREATE_SECRETS,
				Request: &vaultcommon.Observation_CreateSecretsRequest{
					CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: nil, EncryptedValue: "deadbeef"},
						},
					},
				},
				Response: &vaultcommon.Observation_CreateSecretsResponse{
					CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
						Responses: []*vaultcommon.CreateSecretResponse{
							{Id: id},
						},
					},
				},
			},
		},
		{
			name: "UpdateSecrets with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
				Request: &vaultcommon.Observation_UpdateSecretsRequest{
					UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: nil, EncryptedValue: "deadbeef"},
						},
					},
				},
				Response: &vaultcommon.Observation_UpdateSecretsResponse{
					UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
						Responses: []*vaultcommon.UpdateSecretResponse{
							{Id: id},
						},
					},
				},
			},
		},
		{
			name: "CreateSecrets response with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_CREATE_SECRETS,
				Request: &vaultcommon.Observation_CreateSecretsRequest{
					CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: "deadbeef"},
						},
					},
				},
				Response: &vaultcommon.Observation_CreateSecretsResponse{
					CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
						Responses: []*vaultcommon.CreateSecretResponse{
							{Id: nil},
						},
					},
				},
			},
		},
		{
			name: "UpdateSecrets response with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
				Request: &vaultcommon.Observation_UpdateSecretsRequest{
					UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: "deadbeef"},
						},
					},
				},
				Response: &vaultcommon.Observation_UpdateSecretsResponse{
					UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
						Responses: []*vaultcommon.UpdateSecretResponse{
							{Id: nil},
						},
					},
				},
			},
		},
		{
			name: "DeleteSecrets response with nil Id",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_DELETE_SECRETS,
				Request: &vaultcommon.Observation_DeleteSecretsRequest{
					DeleteSecretsRequest: &vaultcommon.DeleteSecretsRequest{
						Ids: []*vaultcommon.SecretIdentifier{id},
					},
				},
				Response: &vaultcommon.Observation_DeleteSecretsResponse{
					DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
						Responses: []*vaultcommon.DeleteSecretResponse{
							{Id: nil},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rdr := &kv{m: make(map[string]response)}

			anyp, err := anypb.New(tc.obs.GetGetSecretsRequest())
			if anyp == nil {
				// For non-GetSecrets types, use the appropriate request
				switch tc.obs.RequestType {
				case vaultcommon.RequestType_CREATE_SECRETS:
					anyp, err = anypb.New(tc.obs.GetCreateSecretsRequest())
				case vaultcommon.RequestType_UPDATE_SECRETS:
					anyp, err = anypb.New(tc.obs.GetUpdateSecretsRequest())
				case vaultcommon.RequestType_DELETE_SECRETS:
					anyp, err = anypb.New(tc.obs.GetDeleteSecretsRequest())
				default:
					t.FailNow()
				}
			}
			require.NoError(t, err)

			err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
				[]*vaultcommon.StoredPendingQueueItem{
					{Id: "request-1", Item: anyp},
				},
			)
			require.NoError(t, err)

			o := &vaultcommon.Observations{
				Observations: []*vaultcommon.Observation{tc.obs},
			}
			ob := protoMarshal(t, o)

			require.NotPanics(t, func() {
				err = r.ValidateObservation(
					t.Context(),
					seqNr,
					types.AttributedQuery{},
					types.AttributedObservation{
						Observer: 0, Observation: ob,
					},
					rdr,
					bf,
				)
			})
			require.Error(t, err, "expected an error for nil secret identifier, not a panic")
		})
	}
}

func TestPlugin_ValidateObservation_CiphertextSize(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	// maxCipherTextLengthBytes = 10 bytes, so any ciphertext > 10 decoded bytes should be rejected
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			1,
			pk,
			shares[0],
			1,
			10,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	seqNr := uint64(1)
	bf := &blobber{}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}

	// 11 bytes hex-encoded = 22 hex chars, exceeds the 10-byte limit
	oversizedCiphertext := hex.EncodeToString(make([]byte, 11))
	validCiphertext := hex.EncodeToString(make([]byte, 5))

	tests := []struct {
		name      string
		obs       *vaultcommon.Observation
		errSubstr string
	}{
		{
			name: "CreateSecrets with oversized ciphertext",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_CREATE_SECRETS,
				Request: &vaultcommon.Observation_CreateSecretsRequest{
					CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: oversizedCiphertext},
						},
					},
				},
				Response: &vaultcommon.Observation_CreateSecretsResponse{
					CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
						Responses: []*vaultcommon.CreateSecretResponse{
							{Id: id},
						},
					},
				},
			},
			errSubstr: "ciphertext size exceeds maximum allowed size",
		},
		{
			name: "UpdateSecrets with oversized ciphertext",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
				Request: &vaultcommon.Observation_UpdateSecretsRequest{
					UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: oversizedCiphertext},
						},
					},
				},
				Response: &vaultcommon.Observation_UpdateSecretsResponse{
					UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
						Responses: []*vaultcommon.UpdateSecretResponse{
							{Id: id},
						},
					},
				},
			},
			errSubstr: "ciphertext size exceeds maximum allowed size",
		},
		{
			name: "CreateSecrets with invalid hex ciphertext",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_CREATE_SECRETS,
				Request: &vaultcommon.Observation_CreateSecretsRequest{
					CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: "not-valid-hex!"},
						},
					},
				},
				Response: &vaultcommon.Observation_CreateSecretsResponse{
					CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
						Responses: []*vaultcommon.CreateSecretResponse{
							{Id: id},
						},
					},
				},
			},
			errSubstr: "invalid hex encoding for ciphertext",
		},
		{
			name: "UpdateSecrets with invalid hex ciphertext",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
				Request: &vaultcommon.Observation_UpdateSecretsRequest{
					UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: "not-valid-hex!"},
						},
					},
				},
				Response: &vaultcommon.Observation_UpdateSecretsResponse{
					UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
						Responses: []*vaultcommon.UpdateSecretResponse{
							{Id: id},
						},
					},
				},
			},
			errSubstr: "invalid hex encoding for ciphertext",
		},
		{
			name: "CreateSecrets with valid ciphertext passes",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_CREATE_SECRETS,
				Request: &vaultcommon.Observation_CreateSecretsRequest{
					CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: validCiphertext},
						},
					},
				},
				Response: &vaultcommon.Observation_CreateSecretsResponse{
					CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
						Responses: []*vaultcommon.CreateSecretResponse{
							{Id: id},
						},
					},
				},
			},
		},
		{
			name: "UpdateSecrets with valid ciphertext passes",
			obs: &vaultcommon.Observation{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
				Request: &vaultcommon.Observation_UpdateSecretsRequest{
					UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{Id: id, EncryptedValue: validCiphertext},
						},
					},
				},
				Response: &vaultcommon.Observation_UpdateSecretsResponse{
					UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
						Responses: []*vaultcommon.UpdateSecretResponse{
							{Id: id},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rdr := &kv{m: make(map[string]response)}

			var anyp *anypb.Any
			switch tc.obs.RequestType {
			case vaultcommon.RequestType_CREATE_SECRETS:
				anyp, err = anypb.New(tc.obs.GetCreateSecretsRequest())
			case vaultcommon.RequestType_UPDATE_SECRETS:
				anyp, err = anypb.New(tc.obs.GetUpdateSecretsRequest())
			default:
				t.FailNow()
			}
			require.NoError(t, err)

			err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
				[]*vaultcommon.StoredPendingQueueItem{
					{Id: "request-1", Item: anyp},
				},
			)
			require.NoError(t, err)

			o := &vaultcommon.Observations{
				Observations: []*vaultcommon.Observation{tc.obs},
			}
			ob := protoMarshal(t, o)

			err = r.ValidateObservation(
				t.Context(),
				seqNr,
				types.AttributedQuery{},
				types.AttributedObservation{
					Observer: 0, Observation: ob,
				},
				rdr,
				bf,
			)

			if tc.errSubstr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlugin_StateTransition_PendingQueueEnabled_NewQuora_NotGetRequest(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-id",
	}
	resp := &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: []*vaultcommon.SecretIdentifier{id},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetListSecretIdentifiersRequest()))
	assert.True(t, proto.Equal(resp, o.GetListSecretIdentifiersResponse()))

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_PendingQueueEnabled_GetRequest(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:   store,
		metrics: newTestMetrics(t),
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			1,
			1024,
			100,
			100,
			100,
			10,
		),
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ek, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(ek[:])
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetGetSecretsRequest()))
	assert.True(t, proto.Equal(resp, o.GetGetSecretsResponse()))

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_MaxShareSize(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	owner := "0x0001020304050607080900010203040506070809"
	ownerAddress := common.HexToAddress(owner)
	var label [32]byte
	copy(label[12:], ownerAddress.Bytes()) // left-pad with 12 zero

	recipientPub, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	expectedSize := cresettings.Default.VaultShareSizeLimit.DefaultValue
	for i := range 10 {
		plaintext := make([]byte, i*1024/9) // 0 to ~1kb

		ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
		require.NoError(t, err)

		ctb, err := ciphertext.Marshal()
		require.NoError(t, err)

		share, err := generatePlaintextShare(pk, shares[0], ctb, owner, "")
		require.NoError(t, err)

		eds, err := share.encryptWithKey(hex.EncodeToString(recipientPub[:]))
		require.NoError(t, err)

		assert.GreaterOrEqual(t, expectedSize, len(eds), "share size should be constant regardless of plaintext size (plaintext=%d bytes)", len(plaintext))
	}
}

func makeObservation(t *testing.T, reqType vaultcommon.RequestType, count int) *vaultcommon.Observation {
	ids := make([]*vaultcommon.SecretIdentifier, count)
	for i := range ids {
		ids[i] = &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret" + string(rune('0'+i))}
	}

	switch reqType {
	case vaultcommon.RequestType_GET_SECRETS:
		reqs := make([]*vaultcommon.SecretRequest, count)
		resps := make([]*vaultcommon.SecretResponse, count)
		for i, id := range ids {
			reqs[i] = &vaultcommon.SecretRequest{Id: id}
			resps[i] = &vaultcommon.SecretResponse{Id: id, Result: &vaultcommon.SecretResponse_Error{Error: "err"}}
		}
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: reqType,
			Request:     &vaultcommon.Observation_GetSecretsRequest{GetSecretsRequest: &vaultcommon.GetSecretsRequest{Requests: reqs}},
			Response:    &vaultcommon.Observation_GetSecretsResponse{GetSecretsResponse: &vaultcommon.GetSecretsResponse{Responses: resps}},
		}
	case vaultcommon.RequestType_CREATE_SECRETS:
		secrets := make([]*vaultcommon.EncryptedSecret, count)
		resps := make([]*vaultcommon.CreateSecretResponse, count)
		for i, id := range ids {
			secrets[i] = &vaultcommon.EncryptedSecret{Id: id, EncryptedValue: "deadbeef"}
			resps[i] = &vaultcommon.CreateSecretResponse{Id: id}
		}
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: reqType,
			Request:     &vaultcommon.Observation_CreateSecretsRequest{CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{EncryptedSecrets: secrets}},
			Response:    &vaultcommon.Observation_CreateSecretsResponse{CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{Responses: resps}},
		}
	case vaultcommon.RequestType_UPDATE_SECRETS:
		secrets := make([]*vaultcommon.EncryptedSecret, count)
		resps := make([]*vaultcommon.UpdateSecretResponse, count)
		for i, id := range ids {
			secrets[i] = &vaultcommon.EncryptedSecret{Id: id, EncryptedValue: "deadbeef"}
			resps[i] = &vaultcommon.UpdateSecretResponse{Id: id}
		}
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: reqType,
			Request:     &vaultcommon.Observation_UpdateSecretsRequest{UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{EncryptedSecrets: secrets}},
			Response:    &vaultcommon.Observation_UpdateSecretsResponse{UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{Responses: resps}},
		}
	case vaultcommon.RequestType_DELETE_SECRETS:
		resps := make([]*vaultcommon.DeleteSecretResponse, count)
		for i, id := range ids {
			resps[i] = &vaultcommon.DeleteSecretResponse{Id: id}
		}
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: reqType,
			Request:     &vaultcommon.Observation_DeleteSecretsRequest{DeleteSecretsRequest: &vaultcommon.DeleteSecretsRequest{Ids: ids}},
			Response:    &vaultcommon.Observation_DeleteSecretsResponse{DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{Responses: resps}},
		}
	default:
		t.Fatalf("unsupported request type: %s", reqType)
		return nil
	}
}

func TestPlugin_ValidateObservation_RequestBatchLimit(t *testing.T) {
	maxRequestBatchSize := 2

	tests := []struct {
		name      string
		reqType   vaultcommon.RequestType
		batchSize int
		wantErr   string
	}{
		{
			name:      "GetSecrets exceeding batch limit",
			reqType:   vaultcommon.RequestType_GET_SECRETS,
			batchSize: maxRequestBatchSize + 1,
			wantErr:   "max batch size exceeded for request",
		},
		{
			name:      "CreateSecrets exceeding batch limit",
			reqType:   vaultcommon.RequestType_CREATE_SECRETS,
			batchSize: maxRequestBatchSize + 1,
			wantErr:   "max batch size exceeded for request",
		},
		{
			name:      "UpdateSecrets exceeding batch limit",
			reqType:   vaultcommon.RequestType_UPDATE_SECRETS,
			batchSize: maxRequestBatchSize + 1,
			wantErr:   "max batch size exceeded for request",
		},
		{
			name:      "DeleteSecrets exceeding batch limit",
			reqType:   vaultcommon.RequestType_DELETE_SECRETS,
			batchSize: maxRequestBatchSize + 1,
			wantErr:   "max batch size exceeded for request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			store := requests.NewStore[*vaulttypes.Request]()
			_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
			require.NoError(t, err)
			r := &ReportingPlugin{
				lggr:    lggr,
				store:   store,
				metrics: newTestMetrics(t),
				onchainCfg: ocr3types.ReportingPluginConfig{
					N: 4,
					F: 1,
				},
				cfg: makeReportingPluginConfig(
					t,
					10,
					pk,
					shares[0],
					1,
					1024,
					30,
					30,
					30,
					maxRequestBatchSize,
				),
				unmarshalBlob: mockUnmarshalBlob,
			}
			rdr := &kv{m: make(map[string]response)}

			obs := &vaultcommon.Observations{
				Observations: []*vaultcommon.Observation{
					makeObservation(t, tc.reqType, tc.batchSize),
				},
			}
			ob := protoMarshal(t, obs)

			err = r.ValidateObservation(
				t.Context(),
				1,
				types.AttributedQuery{},
				types.AttributedObservation{Observer: 0, Observation: ob},
				rdr,
				&blobber{},
			)

			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlugin_ValidateObservation_ListSecretIdentifiersExceedsMaxSecretsPerOwner(t *testing.T) {
	maxSecretsPerOwner := 3

	lggr := logger.TestLogger(t)
	store := requests.NewStore[*vaulttypes.Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:    lggr,
		store:   store,
		metrics: newTestMetrics(t),
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		cfg: makeReportingPluginConfig(
			t,
			10,
			pk,
			shares[0],
			maxSecretsPerOwner,
			1024,
			30,
			30,
			30,
			10,
		),
		unmarshalBlob: mockUnmarshalBlob,
	}

	listReq := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "owner",
		Namespace: "main",
		RequestId: "request-1",
	}

	identifiers := make([]*vaultcommon.SecretIdentifier, maxSecretsPerOwner+1)
	for i := range identifiers {
		identifiers[i] = &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: fmt.Sprintf("secret%d", i)}
	}

	observation := &vaultcommon.Observation{
		Id:          "request-1",
		RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
		Request: &vaultcommon.Observation_ListSecretIdentifiersRequest{
			ListSecretIdentifiersRequest: listReq,
		},
		Response: &vaultcommon.Observation_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: &vaultcommon.ListSecretIdentifiersResponse{
				Identifiers: identifiers,
				Success:     true,
			},
		},
	}

	rdr := &kv{m: make(map[string]response)}
	anyReq, err := anypb.New(listReq)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyReq},
		},
	)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{observation},
	}
	ob := protoMarshal(t, obs)

	err = r.ValidateObservation(
		t.Context(),
		1,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: ob},
		rdr,
		&blobber{},
	)

	require.ErrorContains(t, err, "ListSecretIdentifiers response exceeds maximum number of secrets per owner")
}

func TestUserFacingError(t *testing.T) {
	t.Run("returns error message for userError", func(t *testing.T) {
		err := newUserError("key does not exist")
		assert.Equal(t, "key does not exist", userFacingError(err, "fallback"))
	})

	t.Run("returns fallback for non-userError", func(t *testing.T) {
		err := errors.New("internal failure")
		assert.Equal(t, "fallback msg", userFacingError(err, "fallback msg"))
	})

	t.Run("returns wrapped error message for wrapped userError", func(t *testing.T) {
		err := fmt.Errorf("context: %w", newUserError("bad key"))
		assert.Equal(t, "context: bad key", userFacingError(err, "fallback"))
	})
}

func TestLogUserErrorAware(t *testing.T) {
	t.Run("logs at debug level for userError", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := newUserError("key does not exist")

		logUserErrorAware(lggr, "failed to observe request", err, "id", "req-1")

		debugLogs := observed.FilterLevelExact(zapcore.DebugLevel)
		errorLogs := observed.FilterLevelExact(zapcore.ErrorLevel)
		assert.Equal(t, 1, debugLogs.Len())
		assert.Equal(t, 0, errorLogs.Len())

		entry := debugLogs.All()[0]
		assert.Equal(t, "failed to observe request", entry.Message)
		fields := entry.ContextMap()
		assert.Equal(t, "req-1", fields["id"])
		assert.Contains(t, fmt.Sprint(fields["error"]), "key does not exist")
	})

	t.Run("logs at error level for internal error", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := errors.New("database connection lost")

		logUserErrorAware(lggr, "failed to observe request", err, "id", "req-2")

		debugLogs := observed.FilterLevelExact(zapcore.DebugLevel)
		errorLogs := observed.FilterLevelExact(zapcore.ErrorLevel)
		assert.Equal(t, 0, debugLogs.Len())
		assert.Equal(t, 1, errorLogs.Len())

		entry := errorLogs.All()[0]
		assert.Equal(t, "failed to observe request", entry.Message)
		fields := entry.ContextMap()
		assert.Equal(t, "req-2", fields["id"])
		assert.Contains(t, fmt.Sprint(fields["error"]), "database connection lost")
	})

	t.Run("logs at debug level for wrapped userError", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := fmt.Errorf("validation: %w", newUserError("bad input"))

		logUserErrorAware(lggr, "request failed", err, "op", "create")

		debugLogs := observed.FilterLevelExact(zapcore.DebugLevel)
		errorLogs := observed.FilterLevelExact(zapcore.ErrorLevel)
		assert.Equal(t, 1, debugLogs.Len())
		assert.Equal(t, 0, errorLogs.Len())
	})

	t.Run("includes all key-value pairs in log entry", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := errors.New("internal error")

		logUserErrorAware(lggr, "op failed", err, "id", "req-3", "requestID", "abc-123")

		entry := observed.FilterLevelExact(zapcore.ErrorLevel).All()[0]
		fields := entry.ContextMap()
		assert.Equal(t, "req-3", fields["id"])
		assert.Equal(t, "abc-123", fields["requestID"])
		assert.Contains(t, fmt.Sprint(fields["error"]), "internal error")
	})
}

func TestPlugin_broadcastBlobPayloads(t *testing.T) {
	t.Run("empty payloads returns empty slice", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		fetcher := &callbackBlobFetcher{fn: func([]byte) error { return nil }}
		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, nil, nil)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("all payloads broadcast successfully", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		fetcher := &callbackBlobFetcher{fn: func([]byte) error { return nil }}
		payloads := [][]byte{[]byte("p1"), []byte("p2"), []byte("p3")}
		ids := []string{"req-1", "req-2", "req-3"}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Len(t, result, 3)
		for _, item := range result {
			assert.Equal(t, []byte("handle"), item)
		}
	})

	t.Run("failed broadcast is skipped and logged", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		fetcher := &callbackBlobFetcher{fn: func(payload []byte) error {
			if string(payload) == "p2" {
				return errors.New("broadcast error")
			}
			return nil
		}}

		payloads := [][]byte{[]byte("p1"), []byte("p2"), []byte("p3")}
		ids := []string{"req-1", "req-2", "req-3"}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 5, payloads, ids)
		require.NoError(t, err)
		assert.Len(t, result, 2)

		warnLogs := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		assert.Equal(t, 1, warnLogs.Len())
		fields := warnLogs.All()[0].ContextMap()
		assert.Equal(t, "req-2", fields["requestID"])
		assert.Equal(t, uint64(5), fields["seqNr"])
		assert.Contains(t, fmt.Sprint(fields["err"]), "broadcast error")
	})

	t.Run("all broadcasts fail returns empty slice", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		fetcher := &errorBlobBroadcastFetcher{err: errors.New("network down")}
		payloads := [][]byte{[]byte("p1"), []byte("p2")}
		ids := []string{"req-1", "req-2"}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Empty(t, result)

		warnLogs := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		assert.Equal(t, 2, warnLogs.Len())
	})

	t.Run("marshal blob failure skips item and logs warning", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return nil, errors.New("marshal error")
			},
		}

		fetcher := &callbackBlobFetcher{fn: func([]byte) error { return nil }}
		payloads := [][]byte{[]byte("p1"), []byte("p2")}
		ids := []string{"req-1", "req-2"}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Empty(t, result)

		warnLogs := observed.FilterMessage("failed to marshal blob handle, skipping")
		assert.Equal(t, 2, warnLogs.Len())
	})

	t.Run("mix of broadcast and marshal failures", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)

		marshalCallCount := atomic.Int32{}
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				n := marshalCallCount.Add(1)
				if n == 1 {
					return nil, errors.New("marshal error")
				}
				return []byte("handle"), nil
			},
		}

		fetcher := &callbackBlobFetcher{fn: func(payload []byte) error {
			if string(payload) == "p1" {
				return errors.New("broadcast error")
			}
			return nil
		}}

		payloads := [][]byte{[]byte("p1"), []byte("p2"), []byte("p3")}
		ids := []string{"req-1", "req-2", "req-3"}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)

		broadcastWarns := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		marshalWarns := observed.FilterMessage("failed to marshal blob handle, skipping")
		assert.Equal(t, 1, broadcastWarns.Len())
		assert.Equal(t, 1, marshalWarns.Len())
		assert.Len(t, result, 1)
	})

	t.Run("context cancellation propagates error", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		fetcher := &callbackBlobFetcher{fn: func([]byte) error {
			return ctx.Err()
		}}

		payloads := [][]byte{[]byte("p1"), []byte("p2")}
		ids := []string{"req-1", "req-2"}

		result, err := r.broadcastBlobPayloads(ctx, fetcher, 1, payloads, ids)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("context deadline exceeded propagates error", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		ctx, cancel := context.WithTimeout(t.Context(), 0)
		defer cancel()
		<-ctx.Done()

		fetcher := &callbackBlobFetcher{fn: func([]byte) error {
			return ctx.Err()
		}}

		payloads := [][]byte{[]byte("p1")}
		ids := []string{"req-1"}

		result, err := r.broadcastBlobPayloads(ctx, fetcher, 1, payloads, ids)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("slow broadcast hits per-call timeout and is skipped", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		r := &ReportingPlugin{
			lggr:    lggr,
			metrics: newTestMetrics(t),
			marshalBlob: func(ocr3_1types.BlobHandle) ([]byte, error) {
				return []byte("handle"), nil
			},
		}

		fetcher := &ctxCallbackBlobFetcher{fn: func(ctx context.Context, payload []byte) error {
			if string(payload) == "slow" {
				<-ctx.Done()
				return ctx.Err()
			}
			return nil
		}}

		payloads := [][]byte{[]byte("fast"), []byte("slow")}
		ids := []string{"req-fast", "req-slow"}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Len(t, result, 1)

		warnLogs := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		assert.Equal(t, 1, warnLogs.Len())
		fields := warnLogs.All()[0].ContextMap()
		assert.Equal(t, "req-slow", fields["requestID"])
	})
}
