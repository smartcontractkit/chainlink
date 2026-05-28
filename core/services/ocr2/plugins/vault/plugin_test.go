package vault

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
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
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func testRequestLifecycleTracker(t *testing.T, lggr logger.Logger) *vaultcap.RequestLifecycleTracker {
	t.Helper()
	lc, err := vaultcap.NewRequestLifecycleTracker(lggr)
	require.NoError(t, err)
	return lc
}

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

func TestPlugin_Observation_NothingInBatch(t *testing.T) {
	r := newTestReportingPlugin(t)

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

func TestPlugin_Observation_GetSecretsRequest_OmitsRequestWhenOptimizationsEnabled(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id}},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{{Id: "request-1", Item: anyp}},
	))

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	require.Len(t, obs.Observations, 1)
	assert.Nil(t, obs.Observations[0].GetGetSecretsRequest())
}

func TestPlugin_Observation_PendingQueueEnabled_EmptyPendingQueue(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store), withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
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

	// Two local GetSecrets requests share one batched blob; PendingQueueItems is one handle per blob.
	assert.Len(t, obs.PendingQueueItems, 1)

	assertPendingQueueItemsContain(t, bf.blobs, map[string]proto.Message{
		expectedID:  p,
		expectedID2: p,
	})

	assert.NotEmpty(t, obs.SortNonce)
}

func TestPlugin_Observation_PendingQueueEnabled_WithPendingQueueProvided(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store), withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
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

	// Two local GetSecrets requests share one batched blob; PendingQueueItems is one handle per blob.
	assert.Len(t, obs.PendingQueueItems, 1)

	assertPendingQueueItemsContain(t, bf.blobs, map[string]proto.Message{
		expectedID:  p,
		expectedID2: p,
	})

	assert.NotEmpty(t, obs.SortNonce)
}

func TestPlugin_Observation_PendingQueueEnabled_ItemBothInPendingQueueAndLocalQueue(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
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
	assert.True(t, proto.Equal(p, gotO.GetGetSecretsRequest()))

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

	remaining := make(map[string]proto.Message, len(expected))
	for id, payload := range expected {
		remaining[id] = payload
	}

	var total int
	for _, got := range gotItems {
		items, err := unmarshalPendingQueueBlob(got)
		require.NoError(t, err)
		total += len(items)
		for _, gotMsg := range items {
			expectedPayload, ok := remaining[gotMsg.Id]
			require.True(t, ok, "unexpected pending queue item id %q", gotMsg.Id)

			gotPayload, err := gotMsg.Item.UnmarshalNew()
			require.NoError(t, err)
			assert.True(t, proto.Equal(expectedPayload, gotPayload))

			delete(remaining, gotMsg.Id)
		}
	}

	require.Equal(t, len(expected), total, "item count across blob payloads")
	assert.Empty(t, remaining)
}

func TestPendingQueueBlobMarshalUnmarshal_legacyAndBatch(t *testing.T) {
	id := &vaultcommon.SecretIdentifier{Owner: "o", Namespace: "n", Key: "k"}
	req := &vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id}}}
	any1, err := anypb.New(req)
	require.NoError(t, err)
	s1 := &vaultcommon.StoredPendingQueueItem{Id: "a", Item: any1}
	any2, err := anypb.New(req)
	require.NoError(t, err)
	s2 := &vaultcommon.StoredPendingQueueItem{Id: "b", Item: any2}

	legacy, err := marshalPendingQueueBlobPayload([]*vaultcommon.StoredPendingQueueItem{s1})
	require.NoError(t, err)
	items, err := unmarshalPendingQueueBlob(legacy)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "a", items[0].Id)

	batchBytes, err := marshalPendingQueueBlobPayload([]*vaultcommon.StoredPendingQueueItem{s1, s2})
	require.NoError(t, err)
	items2, err := unmarshalPendingQueueBlob(batchBytes)
	require.NoError(t, err)
	require.Len(t, items2, 2)
	require.Equal(t, "a", items2[0].Id)
	require.Equal(t, "b", items2[1].Id)
}

func TestPrepareObservationPendingQueueBlobs_packsManySmallItemsInOneObservation(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store), withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{Id: id, EncryptionKeys: []string{pks}},
		},
	}
	const n = 25
	for i := range n {
		require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: fmt.Sprintf("req-%02d", i)}))
	}
	localQueueItems, err := store.All()
	require.NoError(t, err)
	slices.SortFunc(localQueueItems, func(a, b *vaulttypes.Request) int {
		return strings.Compare(a.ID(), b.ID())
	})

	maxBlobSz, err := r.cfg.MaxBlobPayloadBytes.Limit(t.Context())
	require.NoError(t, err)
	batchSize, err := r.cfg.MaxBatchSize.Limit(t.Context())
	require.NoError(t, err)

	pack, err := r.prepareObservationPendingQueueBlobs(t.Context(), 1, localQueueItems, map[string]bool{}, int(maxBlobSz), 2*batchSize)
	require.NoError(t, err)
	require.Equal(t, n, pack.packedItemCount)
	require.False(t, pack.truncated)
	require.LessOrEqual(t, len(pack.blobPayloads), 2*batchSize)

	var unpacked int
	for _, blob := range pack.blobPayloads {
		items, uerr := unmarshalPendingQueueBlob(blob)
		require.NoError(t, uerr)
		unpacked += len(items)
	}
	require.Equal(t, n, unpacked)
}

func TestPrepareObservationPendingQueueBlobs_flushesAndContinuesWhenBatchFull(t *testing.T) {
	// Items are larger than half maxBlobBytes so only one fits per blob.
	// The loop must flush the current batch and start a new one for each item.
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store), withVaultOptimizationsEnabled())

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "my_secret"}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{pks}}},
	}

	const n = 3
	for i := range n {
		require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: fmt.Sprintf("req-%02d", i)}))
	}
	localQueueItems, err := store.All()
	require.NoError(t, err)
	slices.SortFunc(localQueueItems, func(a, b *vaulttypes.Request) int {
		return strings.Compare(a.ID(), b.ID())
	})

	// Compute the wire size of a single item so we can set maxBlobBytes to just above it,
	// forcing each item into its own blob.
	anyMsg, err := anypb.New(p)
	require.NoError(t, err)
	singleItem := &vaultcommon.StoredPendingQueueItem{Id: "req-00", Item: anyMsg}
	singlePayload, err := marshalPendingQueueBlobPayload([]*vaultcommon.StoredPendingQueueItem{singleItem})
	require.NoError(t, err)
	maxBlobBytes := len(singlePayload) + 5 // fits exactly one item

	pack, err := r.prepareObservationPendingQueueBlobs(t.Context(), 1, localQueueItems, map[string]bool{}, maxBlobBytes, 10)
	require.NoError(t, err)
	require.Equal(t, n, pack.packedItemCount)
	require.False(t, pack.truncated)
	require.Len(t, pack.blobPayloads, n, "each item should be its own blob")

	var unpacked int
	for _, blob := range pack.blobPayloads {
		items, uerr := unmarshalPendingQueueBlob(blob)
		require.NoError(t, uerr)
		unpacked += len(items)
	}
	require.Equal(t, n, unpacked)
}

func TestPrepareObservationPendingQueueBlobs_truncatesWhenHandleCountExceeded(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store), withVaultOptimizationsEnabled())

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "my_secret"}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{pks}}},
	}

	const n = 10
	for i := range n {
		require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: fmt.Sprintf("req-%02d", i)}))
	}
	localQueueItems, err := store.All()
	require.NoError(t, err)
	slices.SortFunc(localQueueItems, func(a, b *vaulttypes.Request) int {
		return strings.Compare(a.ID(), b.ID())
	})

	anyMsg, err := anypb.New(p)
	require.NoError(t, err)
	singleItem := &vaultcommon.StoredPendingQueueItem{Id: "req-00", Item: anyMsg}
	singlePayload, err := marshalPendingQueueBlobPayload([]*vaultcommon.StoredPendingQueueItem{singleItem})
	require.NoError(t, err)
	// Force one item per blob by making maxBlobBytes just above single-item size.
	maxBlobBytes := len(singlePayload) + 5
	const maxBlobHandleCount = 3

	pack, err := r.prepareObservationPendingQueueBlobs(t.Context(), 1, localQueueItems, map[string]bool{}, maxBlobBytes, maxBlobHandleCount)
	require.NoError(t, err)
	require.True(t, pack.truncated)
	require.Len(t, pack.blobPayloads, maxBlobHandleCount)
	require.Equal(t, maxBlobHandleCount, pack.packedItemCount)
}

func TestPrepareObservationPendingQueueBlobs_errorWhenSingleItemTooLarge(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store), withVaultOptimizationsEnabled())

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "my_secret"}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{pks}}},
	}
	require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: "req-00"}))

	localQueueItems, err := store.All()
	require.NoError(t, err)

	_, err = r.prepareObservationPendingQueueBlobs(t.Context(), 1, localQueueItems, map[string]bool{}, 1, 10)
	require.Error(t, err)
	require.Contains(t, err.Error(), "single pending queue item exceeds max blob payload size")
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
	store := requests.NewStore[*vaulttypes.Request]()
	// Blob payload cap: one pending item fits in a blob, two batched items do not — two BroadcastBlob calls.
	// Blob cap between one item (155B) and two batched (316B) so two BroadcastBlob calls occur.
	r := newTestReportingPlugin(t, withStore(store), withMaxBlobPayloadBytes(310))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
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
	r := newTestReportingPlugin(t, withStore(store), withLggr(lggr))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
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
		name            string
		id              *vaultcommon.SecretIdentifier
		maxIDLen        int
		maxOwnerLen     int
		maxNamespaceLen int
		maxKeyLen       int
		err             string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "secret identifier cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner exceeds maximum length of 3b",
		},
		{
			name:            "namespace exceeds maximum length",
			maxNamespaceLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "namespace exceeds maximum length of 3b",
		},
		{
			name:      "key exceeds maximum length",
			maxKeyLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "key exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		ownerLen, namespaceLen, keyLen := 256, 256, 256
		if tc.maxIDLen > 0 {
			ownerLen = tc.maxIDLen / 3
			namespaceLen = tc.maxIDLen / 3
			keyLen = tc.maxIDLen / 3
		}
		if tc.maxOwnerLen > 0 {
			ownerLen = tc.maxOwnerLen
		}
		if tc.maxNamespaceLen > 0 {
			namespaceLen = tc.maxNamespaceLen
		}
		if tc.maxKeyLen > 0 {
			keyLen = tc.maxKeyLen
		}
		r := newTestReportingPlugin(t, withMaxIdentifierLengths(ownerLen, namespaceLen, keyLen))

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
		assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

		batchResp := o.GetGetSecretsResponse()
		assert.Len(t, p.Requests, 1)
		assert.Len(t, p.Requests, len(batchResp.Responses))

		assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_GetSecretsRequest_ResponseUsesCanonicalIdentifier(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(batchResp.Responses[0].Id, createdID))
}

func TestPlugin_Observation_GetSecretsRequest_WorkflowOwnerLabelAccepted(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	ownerAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	owner := ownerAddr.Hex()
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("my-secret-value", pk, ownerAddr)
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
	require.NotNil(t, batchResp.Responses[0].GetId())
	assert.Equal(t, owner, batchResp.Responses[0].GetId().GetOwner())
	assert.Empty(t, batchResp.Responses[0].GetError())
}

func TestPlugin_Observation_GetSecretsRequest_WrongOwnerLabelRejected(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	goodOwner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	wrongOwner := common.HexToAddress("0x2222222222222222222222222222222222222222")

	id := &vaultcommon.SecretIdentifier{
		Owner:     wrongOwner.Hex(),
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("my-secret-value", pk, goodOwner)
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
	r := newTestReportingPlugin(t)

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
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "key does not exist")
}

func TestPlugin_Observation_GetSecretsRequest_SecretExistsButIsIncorrect(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]))

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
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to convert public key to bytes")
}

func TestPlugin_Observation_GetSecretsRequest_SecretLabelIsInvalid(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to handle get secret request")
}

func TestPlugin_Observation_GetSecretsRequest_Success(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

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

func TestPlugin_Observation_GetSecretsRequest_BinarySharesWhenOptimizationsEnabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withVaultOptimizationsEnabled())

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
	copy(label[12:], ownerAddress.Bytes())
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

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	require.Len(t, obs.Observations, 1)
	resp := obs.Observations[0].GetGetSecretsResponse().Responses[0]
	require.Empty(t, resp.GetError())

	encShares := resp.GetData().EncryptedDecryptionKeyShares
	require.Len(t, encShares, 1)
	require.Empty(t, encShares[0].Shares)
	require.Len(t, encShares[0].BinaryShares, 1)

	msg, ok := box.OpenAnonymous(nil, encShares[0].BinaryShares[0], pubK, privK)
	require.True(t, ok)

	ds := &tdh2easy.DecryptionShare{}
	require.NoError(t, ds.Unmarshal(msg))

	ct := &tdh2easy.Ciphertext{}
	require.NoError(t, ct.UnmarshalVerify(ciphertextBytes, pk))
	gotSecret, err := tdh2easy.Aggregate(ct, []*tdh2easy.DecryptionShare{ds}, 3)
	require.NoError(t, err)
	require.Equal(t, plaintext, gotSecret)

	err = r.ValidateObservation(
		t.Context(),
		1,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: data},
		rdr,
		&blobber{},
	)
	require.NoError(t, err)
}

func TestPlugin_Observation_CreateSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name            string
		id              *vaultcommon.SecretIdentifier
		maxIDLen        int
		maxOwnerLen     int
		maxNamespaceLen int
		maxKeyLen       int
		err             string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "secret identifier cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner exceeds maximum length of 3b",
		},
		{
			name:            "namespace exceeds maximum length",
			maxNamespaceLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "namespace exceeds maximum length of 3b",
		},
		{
			name:      "key exceeds maximum length",
			maxKeyLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "key exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		ownerLen, namespaceLen, keyLen := 256, 256, 256
		if tc.maxIDLen > 0 {
			ownerLen = tc.maxIDLen / 3
			namespaceLen = tc.maxIDLen / 3
			keyLen = tc.maxIDLen / 3
		}
		if tc.maxOwnerLen > 0 {
			ownerLen = tc.maxOwnerLen
		}
		if tc.maxNamespaceLen > 0 {
			namespaceLen = tc.maxNamespaceLen
		}
		if tc.maxKeyLen > 0 {
			keyLen = tc.maxKeyLen
		}
		r := newTestReportingPlugin(t, withMaxIdentifierLengths(ownerLen, namespaceLen, keyLen))

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
	r := newTestReportingPlugin(t)

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withBatchSize(10), withMaxIdentifierLengths(30, 30, 30), withKeys(pk, shares[0]))

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

	// o1 := os.Outcomes[0]
	// assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o1.RequestType)
	// assert.Len(t, o1.GetCreateSecretsResponse().Responses, 1)
	// r1 := o1.GetCreateSecretsResponse().Responses[0]
	// assert.True(t, r1.Success)

	// o2 := os.Outcomes[1]
	// assert.Equal(t, vaultcommon.RequestType_CREATE_SECRETS, o2.RequestType)
	// assert.Len(t, o2.GetCreateSecretsResponse().Responses, 1)
	// r2 := o2.GetCreateSecretsResponse().Responses[0]
	// assert.False(t, r2.Success)
	// assert.Contains(t, r2.GetError(), "owner has reached maximum number of secrets")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext(t *testing.T) {
	r := newTestReportingPlugin(t)

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
	assert.Contains(t, resp.GetError(), "failed to decode encrypted value")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_TooLong(t *testing.T) {
	r := newTestReportingPlugin(t, withMaxCiphertextLengthBytes(10))

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
	// Wrong key
	_, wrongPublicKey, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	// Right key
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withBatchSize(10), withKeys(pk, shares[0]))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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

func TestPlugin_Observation_CreateSecretsRequest_WorkflowOwnerLabelAccepted(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	ownerAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	owner := ownerAddr.Hex()
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "secret",
	}

	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("my secret value", pk, ownerAddr)
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: id, EncryptedValue: encrypted},
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

func TestPlugin_Observation_CreateSecretsRequest_WrongOwnerLabelRejected(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	goodOwner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	wrongOwner := common.HexToAddress("0x2222222222222222222222222222222222222222")
	id := &vaultcommon.SecretIdentifier{
		Owner:     wrongOwner.Hex(),
		Namespace: "main",
		Key:       "secret",
	}

	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("my secret value", pk, goodOwner)
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	p := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: id, EncryptedValue: encrypted},
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

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetCreateSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	assert.Contains(t, batchResp.Responses[0].GetError(), "does not match workflow owner label")
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
			Shares:        []string{hex.EncodeToString(encrypted)},
		}
	}
	return result
}

func makeBinaryEncryptedShares(t *testing.T, ciphertext *tdh2easy.Ciphertext, privateShare *tdh2easy.PrivateShare, keys []string) []*vaultcommon.EncryptedShares {
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
			BinaryShares:  [][]byte{encrypted},
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

func makeGetSecretsBinaryObservations(
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
							EncryptedDecryptionKeyShares: makeBinaryEncryptedShares(t, ciphertext, privateShare, encryptionKeys),
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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(4, 10)
	require.NoError(t, err)

	numObservers := 10
	r := newTestReportingPlugin(
		t,
		withKeys(pk, shares[0]),
		withOnchainCfg(10, 3),
		withMaxCiphertextLengthBytes(2000),
		withMaxIdentifierLengths(64, 64, 64),
		withVaultOptimizationsEnabled(),
	)

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
			Observer:    commontypes.OracleID(i),
			Observation: types.Observation(makeGetSecretsBinaryObservations(t, 10, maxOwner, maxNamespace, encryptionKeys, encryptedValue, ciphertext, shares[i])),
		}
	}

	kvStore := &kv{m: make(map[string]response)}
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		aos, kvStore, nil)
	require.NoError(t, err)

	twoFPlusOne := 2*r.onchainCfg.F + 1
	osOut := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, osOut)
	require.NoError(t, err)
	for _, outcome := range osOut.Outcomes {
		for _, secretResp := range outcome.GetGetSecretsResponse().GetResponses() {
			data := secretResp.GetData()
			if data == nil {
				continue
			}
			for _, enc := range data.EncryptedDecryptionKeyShares {
				assert.Empty(t, enc.Shares)
				assert.Len(t, enc.BinaryShares, twoFPlusOne,
					"expected at most 2f+1 binary shares per encryption key, got %d (N=%d)", len(enc.BinaryShares), numObservers)
			}
		}
	}

	t.Logf("StateTransition response size: %d bytes (%.2f KB)", len(reportPrecursor), float64(len(reportPrecursor))/1024.0)
	maxResponseSize := 512 * 1024
	assert.LessOrEqual(t, len(reportPrecursor), maxResponseSize,
		"StateTransition response size %d exceeds 512KB limit", len(reportPrecursor))
}

func TestPlugin_ValidateObservations_InvalidObservations(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	require.ErrorContains(t, err, "GetSecrets observation must have a response")

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

func TestPlugin_ValidateObservations_RequiresObservedIDsInPendingQueue(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

	// Observation id not present in the pending queue must be rejected.
	unknown := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "unknown",
	}
	respUnknown := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: unknown,
			},
		},
	}
	obsUnknown := marshalObservations(t, observation{unknown, g, respUnknown})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsUnknown)},
		kv,
		nil,
	)
	require.ErrorContains(t, err, "is not present in the pending queue")

	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id2,
			},
		},
	}

	obsb := marshalObservations(t, observation{id2, g, resp})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	require.NoError(t, err)

	respDel := &vaultcommon.DeleteSecretsResponse{
		Responses: []*vaultcommon.DeleteSecretResponse{
			{Id: id, Success: false, Error: ""},
		},
	}
	obsb = marshalObservations(t, observation{id, d, respDel}, observation{id2, g, resp})
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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	req := &vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id}}}
	anyMsg, err := anypb.New(req)
	require.NoError(t, err)
	stored := &vaultcommon.StoredPendingQueueItem{Id: "request-1", Item: anyMsg}
	dupBytes, err := proto.Marshal(stored)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0},
			{1},
		},
	}
	obsb, err := proto.Marshal(obs)
	require.NoError(t, err)

	bf := &blobber{
		blobs: [][]byte{
			dupBytes,
			dupBytes,
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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	assert.True(t, proto.Equal(resp, o.GetGetSecretsResponse()))

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_GetSecretsRequest_CombinesShares(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

func TestPlugin_StateTransition_GetSecretsRequest_CombinesBinaryShares(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	share1 := []byte("encrypted-share-1")
	share2 := []byte("encrypted-share-2")
	share3 := []byte("encrypted-share-3")
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
								BinaryShares:  [][]byte{share1},
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
								BinaryShares:  [][]byte{share2},
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
								BinaryShares:  [][]byte{share3},
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

	require.Len(t, os.Outcomes, 1)
	o := os.Outcomes[0]

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
								BinaryShares:  [][]byte{share1, share2, share3},
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

func TestPlugin_StateTransition_GetSecretsRequest_CapsSharesAtTwoFPlusOne(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1), withVaultOptimizationsEnabled())

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
	makeResp := func(share string) *vaultcommon.GetSecretsResponse {
		return &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{
				{
					Id: id,
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue: "encrypted-value",
							EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
								{
									EncryptionKey: "my-encryption-key",
									Shares:        []string{share},
								},
							},
						},
					},
				},
			},
		}
	}

	obsb1 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-1")})
	obsb2 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-2")})
	obsb3 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-3")})
	obsb4 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-4")})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb1)},
			{Observer: 1, Observation: types.Observation(obsb2)},
			{Observer: 2, Observation: types.Observation(obsb3)},
			{Observer: 3, Observation: types.Observation(obsb4)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]

	twoFPlusOne := 2*r.onchainCfg.F + 1
	got := o.GetGetSecretsResponse().Responses[0].GetData().EncryptedDecryptionKeyShares[0].Shares
	assert.Len(t, got, twoFPlusOne)
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

func TestPlugin_StateTransition_GetSecretsRequest_DoesNotCapSharesWhenOptimizationsDisabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id}},
	}
	makeResp := func(share string) *vaultcommon.GetSecretsResponse {
		return &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{
				{
					Id: id,
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue: "encrypted-value",
							EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
								{EncryptionKey: "my-encryption-key", Shares: []string{share}},
							},
						},
					},
				},
			},
		}
	}

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-1")}))},
			{Observer: 1, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-2")}))},
			{Observer: 2, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-3")}))},
			{Observer: 3, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-4")}))},
		},
		&kv{m: make(map[string]response)},
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)

	got := os.Outcomes[0].GetGetSecretsResponse().Responses[0].GetData().EncryptedDecryptionKeyShares[0].Shares
	assert.Len(t, got, 4)
	require.NotNil(t, os.Outcomes[0].GetGetSecretsRequest())
}

func TestPlugin_StateTransition_GetSecretsRequest_OmitsOutcomeRequestWhenOptimizationsEnabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1), withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	req := &vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id}}}
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{EncryptionKey: "my-encryption-key", Shares: []string{"share-1"}},
						},
					},
				},
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
		&kv{m: make(map[string]response)},
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)
	assert.Nil(t, os.Outcomes[0].Request)
}

func TestPlugin_PrepareLegacyObservationPendingQueueBlobs_oneBlobPerItem(t *testing.T) {
	store := requests.NewStore[*vaulttypes.Request]()
	r := newTestReportingPlugin(t, withStore(store))

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "my_secret"}
	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{pks}}},
	}
	const n = 5
	for i := range n {
		require.NoError(t, store.Add(&vaulttypes.Request{Payload: p, IDVal: fmt.Sprintf("req-%02d", i)}))
	}
	localQueueItems, err := store.All()
	require.NoError(t, err)

	batchSize, err := r.cfg.MaxBatchSize.Limit(t.Context())
	require.NoError(t, err)

	pack, err := r.prepareLegacyObservationPendingQueueBlobs(t.Context(), 1, localQueueItems, map[string]bool{}, 2*batchSize)
	require.NoError(t, err)
	require.Equal(t, n, pack.packedItemCount)
	require.Len(t, pack.blobPayloads, n)
	for _, blob := range pack.blobPayloads {
		items, uerr := unmarshalPendingQueueBlob(blob)
		require.NoError(t, uerr)
		require.Len(t, items, 1)
	}
}

func TestPlugin_StateTransition_CreateSecretsRequest_WritesSecrets(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

func TestPlugin_StateTransition_CreateSecretsRequest_PerOwnerLimitEnforcedWhenAtCapacity(t *testing.T) {
	r := newTestReportingPlugin(t, withMaxSecretsPerOwner(1), withOnchainCfg(4, 1))

	const owner = "0x2222222222222222222222222222222222222222"

	kv := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, kv).WriteSecret(t.Context(), &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "existing",
	}, &vaultcommon.StoredSecret{EncryptedSecret: []byte("legacy-value")}))
	rs := newTestReadStore(t, kv)

	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
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
	require.Len(t, o.GetCreateSecretsResponse().Responses, 1)
	assert.False(t, o.GetCreateSecretsResponse().Responses[0].Success)
	assert.Contains(t, o.GetCreateSecretsResponse().Responses[0].Error, "has reached maximum number of secrets")

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	assert.Nil(t, ss)
}

func TestPlugin_StateTransition_CreateSecrets_ResponseOwnerMatchesStoredIdentifier(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(
		t,
		withLggr(lggr),
		withKeys(pk, shares[0]),
		withBatchSize(10),
		withMaxSecretsPerOwner(5),
		withOnchainCfg(4, 1),
	)

	const workflowOwner = "0x5555555555555555555555555555555555555555"

	requestID := &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "secret",
	}

	value := []byte("encrypted-value")
	req := &vaultcommon.CreateSecretsRequest{
		RequestId: "request-id",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id:             requestID,
				EncryptedValue: hex.EncodeToString(value),
			},
		},
	}
	resp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      requestID,
				Success: false,
				Error:   "",
			},
		},
	}

	kv := &kv{m: make(map[string]response)}
	rs := newTestReadStore(t, kv)
	obsb := marshalObservations(t, observation{requestID, req, resp})
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
	expectedResp := &vaultcommon.CreateSecretsResponse{
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id:      requestID,
				Success: true,
				Error:   "",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetCreateSecretsResponse()), o.GetCreateSecretsResponse())

	ss, err := rs.GetSecret(t.Context(), requestID)
	require.NoError(t, err)
	assert.Equal(t, []byte("encrypted-value"), ss.EncryptedSecret)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
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

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
		name            string
		id              *vaultcommon.SecretIdentifier
		maxIDLen        int
		maxOwnerLen     int
		maxNamespaceLen int
		maxKeyLen       int
		err             string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "secret identifier cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner exceeds maximum length of 3b",
		},
		{
			name:            "namespace exceeds maximum length",
			maxNamespaceLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "a",
				Key:       "b",
				Namespace: "world",
			},
			err: "namespace exceeds maximum length of 3b",
		},
		{
			name:      "key exceeds maximum length",
			maxKeyLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "a",
				Namespace: "b",
				Key:       "hello",
			},
			err: "key exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		ownerLen, namespaceLen, keyLen := 256, 256, 256
		if tc.maxIDLen > 0 {
			ownerLen = tc.maxIDLen / 3
			namespaceLen = tc.maxIDLen / 3
			keyLen = tc.maxIDLen / 3
		}
		if tc.maxOwnerLen > 0 {
			ownerLen = tc.maxOwnerLen
		}
		if tc.maxNamespaceLen > 0 {
			namespaceLen = tc.maxNamespaceLen
		}
		if tc.maxKeyLen > 0 {
			keyLen = tc.maxKeyLen
		}
		r := newTestReportingPlugin(t, withMaxIdentifierLengths(ownerLen, namespaceLen, keyLen))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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
	r := newTestReportingPlugin(t)

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
	assert.Contains(t, resp.GetError(), "failed to decode encrypted value")
}

func TestPlugin_Observation_UpdateSecretsRequest_InvalidCiphertext_TooLong(t *testing.T) {
	r := newTestReportingPlugin(t, withMaxCiphertextLengthBytes(10))

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
	// Wrong key
	_, wrongPublicKey, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	// Right key
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

func TestPlugin_StateTransition_DeleteSecretsRequest_SecretDoesNotExist(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30))

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

func TestPlugin_Observation_ListSecretIdentifiers_ListsSecretsForRequestOwner(t *testing.T) {
	r := newTestReportingPlugin(
		t,
		withMaxSecretsPerOwner(5),
		withMaxIdentifierLengths(30, 30, 30),
	)

	const workflowOwner = "0x1111111111111111111111111111111111111111"

	rdr := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, rdr).WriteSecret(t.Context(), &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "legacy-item",
	}, &vaultcommon.StoredSecret{EncryptedSecret: []byte("legacy-data")}))

	p := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-id",
		Owner:     workflowOwner,
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
	assert.Equal(t, workflowOwner, resp.Identifiers[0].Owner)
	assert.Equal(t, "legacy-item", resp.Identifiers[0].Key)
}

func TestPlugin_Observation_ListSecretIdentifiers_DoesNotFallbackWhenGateDisabled(t *testing.T) {
	r := newTestReportingPlugin(t, withMaxSecretsPerOwner(5), withMaxIdentifierLengths(30, 30, 30))

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
		RequestId: "request-id",
		Owner:     orgID,
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

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(
		t,
		withMaxSecretsPerOwner(5),
		withMaxIdentifierLengths(30, 30, 30),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

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

func TestPlugin_StateTransition_StoresPendingQueue_AllConsensusItems(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

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

	stored := []*vaultcommon.StoredPendingQueueItem{
		{Id: "request-id", Item: areq1},
		{Id: "request-id2", Item: areq2},
		{Id: "request-id3", Item: areq3},
		{Id: "request-id4", Item: areq4},
	}
	salt := []byte{}
	slices.SortFunc(stored, func(i, j *vaultcommon.StoredPendingQueueItem) int {
		return bytes.Compare(sortKey(i.Id, salt), sortKey(j.Id, salt))
	})

	bf := &blobber{
		blobs: [][]byte{
			protoMarshal(t, stored[0]),
			protoMarshal(t, stored[1]),
			protoMarshal(t, stored[2]),
			protoMarshal(t, stored[3]),
		},
	}

	o1 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0},
			{1},
			{2},
			{3},
		},
	}
	o1b, err := proto.Marshal(o1)
	require.NoError(t, err)

	o2 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0},
			{1},
			{2},
			{3},
		},
	}
	o2b, err := proto.Marshal(o2)
	require.NoError(t, err)

	o3 := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{
			{0},
			{1},
			{2},
			{3},
		},
	}
	o3b, err := proto.Marshal(o3)
	require.NoError(t, err)

	runStateTransition := func(t *testing.T, r *ReportingPlugin) []*vaultcommon.StoredPendingQueueItem {
		t.Helper()
		ctx := t.Context()
		seqNr := uint64(1)
		rdr := &kv{m: make(map[string]response)}
		r.unmarshalBlob = bf.unmarshalBlob

		reportPrecursor, err := r.StateTransition(
			ctx,
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

		pq, err := newTestReadStore(t, rdr).GetPendingQueue(ctx)
		require.NoError(t, err)
		return pq
	}

	sortedIDs := func(items []*vaultcommon.StoredPendingQueueItem) []string {
		ids := make([]string, len(items))
		for i, it := range items {
			ids[i] = it.Id
		}
		slices.Sort(ids)
		return ids
	}

	t.Run("optimizations enabled writes all F+1 consensus items", func(t *testing.T) {
		r := newTestReportingPlugin(
			t,
			withBatchSize(1),
			withMaxIdentifierLengths(30, 30, 30),
			withKeys(pk, shares[0]),
			withOnchainCfg(4, 1),
			withVaultOptimizationsEnabled(),
		)
		pq := runStateTransition(t, r)
		require.Len(t, pq, len(stored))
		require.Equal(t, sortedIDs(stored), sortedIDs(pq), "all F+1 consensus pending queue items are written without MaxBatchSize truncation")
	})

	t.Run("optimizations disabled applies MaxBatchSize", func(t *testing.T) {
		r := newTestReportingPlugin(
			t,
			withBatchSize(1),
			withMaxIdentifierLengths(30, 30, 30),
			withKeys(pk, shares[0]),
			withOnchainCfg(4, 1),
		)
		pq := runStateTransition(t, r)
		require.Len(t, pq, 1)
		require.Equal(t, stored[0].Id, pq[0].Id, "MaxBatchSize keeps the first item after deterministic sort")
	})
}

func TestPlugin_StateTransition_OutcomesStoppedByPrecursorWireSize(t *testing.T) {
	ctx := context.Background()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	buildListObs := func(obsID string) *vaultcommon.Observation {
		req := &vaultcommon.ListSecretIdentifiersRequest{Owner: "owner", Namespace: "main", RequestId: obsID}
		resp := &vaultcommon.ListSecretIdentifiersResponse{Success: true, Identifiers: []*vaultcommon.SecretIdentifier{}, Error: ""}
		return &vaultcommon.Observation{
			Id:          obsID,
			RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
			Request: &vaultcommon.Observation_ListSecretIdentifiersRequest{
				ListSecretIdentifiersRequest: req,
			},
			Response: &vaultcommon.Observation_ListSecretIdentifiersResponse{
				ListSecretIdentifiersResponse: resp,
			},
		}
	}
	ws := NewWriteStore(&kv{m: make(map[string]response)}, newTestMetrics(t))

	rPrec := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1), withVaultOptimizationsEnabled())

	out1 := &vaultcommon.Outcome{Id: "list-1", RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS}
	rPrec.stateTransitionListSecretIdentifiers(ctx, ws, []*vaultcommon.Observation{buildListObs("list-1")}, out1)
	sz1 := proto.Size(&vaultcommon.Outcomes{Outcomes: []*vaultcommon.Outcome{out1}})

	out2 := &vaultcommon.Outcome{Id: "list-2", RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS}
	rPrec.stateTransitionListSecretIdentifiers(ctx, ws, []*vaultcommon.Observation{buildListObs("list-2")}, out2)
	szBoth := proto.Size(&vaultcommon.Outcomes{Outcomes: []*vaultcommon.Outcome{out1, out2}})
	require.Greater(t, szBoth, sz1)

	rPrec.maxReportsPlusPrecursorBytes = sz1

	kv := &kv{m: make(map[string]response)}

	obsPack := &vaultcommon.Observations{Observations: []*vaultcommon.Observation{buildListObs("list-1"), buildListObs("list-2")}}
	obsb, err := proto.Marshal(obsPack)
	require.NoError(t, err)

	aos := []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(obsb)},
		{Observer: 1, Observation: types.Observation(obsb)},
	}
	prec, err := rPrec.StateTransition(ctx, 1, types.AttributedQuery{}, aos, kv, nil)
	require.NoError(t, err)
	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(prec, os))
	require.Len(t, os.Outcomes, 1)
	require.Equal(t, "list-1", os.Outcomes[0].Id)
}

func TestPlugin_StateTransition_OutcomesNotStoppedByPrecursorWireSizeWhenOptimizationsDisabled(t *testing.T) {
	ctx := context.Background()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	buildListObs := func(obsID string) *vaultcommon.Observation {
		req := &vaultcommon.ListSecretIdentifiersRequest{Owner: "owner", Namespace: "main", RequestId: obsID}
		resp := &vaultcommon.ListSecretIdentifiersResponse{Success: true, Identifiers: []*vaultcommon.SecretIdentifier{}, Error: ""}
		return &vaultcommon.Observation{
			Id:          obsID,
			RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
			Request: &vaultcommon.Observation_ListSecretIdentifiersRequest{
				ListSecretIdentifiersRequest: req,
			},
			Response: &vaultcommon.Observation_ListSecretIdentifiersResponse{
				ListSecretIdentifiersResponse: resp,
			},
		}
	}

	out1 := &vaultcommon.Outcome{Id: "list-1", RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS}
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))
	ws := NewWriteStore(&kv{m: make(map[string]response)}, newTestMetrics(t))
	r.stateTransitionListSecretIdentifiers(ctx, ws, []*vaultcommon.Observation{buildListObs("list-1")}, out1)
	sz1 := proto.Size(&vaultcommon.Outcomes{Outcomes: []*vaultcommon.Outcome{out1}})

	rPrec := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1), withMaxReportsPlusPrecursorBytes(sz1))

	obsPack := &vaultcommon.Observations{Observations: []*vaultcommon.Observation{buildListObs("list-1"), buildListObs("list-2")}}
	obsb, err := proto.Marshal(obsPack)
	require.NoError(t, err)

	aos := []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(obsb)},
		{Observer: 1, Observation: types.Observation(obsb)},
	}
	prec, err := rPrec.StateTransition(ctx, 1, types.AttributedQuery{}, aos, &kv{m: make(map[string]response)}, nil)
	require.NoError(t, err)
	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(prec, os))
	require.Len(t, os.Outcomes, 2)
}

func TestPlugin_ValidateObservation_WireSizeExceedsMax(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1), withMaxObservationBytes(20))

	obs := &vaultcommon.Observations{
		SortNonce:    bytes.Repeat([]byte{'a'}, 500),
		Observations: []*vaultcommon.Observation{},
	}
	obsb, err := proto.Marshal(obs)
	require.NoError(t, err)
	require.Greater(t, len(obsb), r.maxObservationBytes)

	err = r.ValidateObservation(
		t.Context(),
		1,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		&kv{m: make(map[string]response)},
		nil,
	)
	require.ErrorContains(t, err, "wire size")
}

func TestPlugin_StateTransition_StoresPendingQueue_DoesntDoubleCountObservationsFromOneNode(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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

// TestPlugin_ValidateObservation_AcceptsFullPendingQueueObservation verifies that an observation
// with exactly 2*batchSize pending queue items (the maximum Observation can produce) is accepted.
func TestPlugin_ValidateObservation_AcceptsFullPendingQueueObservation(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	batchSize := 1 // MaxBatchSize=1, so 2*batchSize=2 is the intended max pending queue items
	r := newTestReportingPlugin(t, withBatchSize(batchSize), withMaxIdentifierLengths(30, 30, 30), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								BinaryShares:  [][]byte{[]byte("encrypted-binary-share")},
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
	require.NoError(t, err)

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

	// More responses than requests in the original GetSecretsRequest is now valid because
	// request content is no longer included in the observation — validation is response-only.
	resp = &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: "foo",
				},
			},
			{
				Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret2"},
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

	// An empty EncryptedShares list is now valid because the per-key cross-check against
	// the request's EncryptionKeys is no longer performed.
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

	// Batch limit is now enforced on the response count, not the request count.
	rBatchLimited := newTestReportingPlugin(t, withMaxRequestBatchSize(1), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	resp = &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{Id: id, Result: &vaultcommon.SecretResponse_Error{Error: "err"}},
			{Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret2"}, Result: &vaultcommon.SecretResponse_Error{Error: "err"}},
		},
	}

	o1 = &vaultcommon.Observations{
		Observations: []*vaultcommon.Observation{
			{
				Id:          "request-1",
				RequestType: vaultcommon.RequestType_GET_SECRETS,
				Response: &vaultcommon.Observation_GetSecretsResponse{
					GetSecretsResponse: resp,
				},
			},
		},
	}
	o1b = protoMarshal(t, o1)

	err = rBatchLimited.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{
			Observer: 0, Observation: o1b,
		},
		rdr,
		bf,
	)
	require.ErrorContains(t, err, "invalid observation: max batch size exceeded for request")
}

func TestPlugin_ValidateObservation_GetSecrets_MismatchedResponseOwnerRejected(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(
		t,
		withMaxRequestBatchSize(10),
		withMaxIdentifierLengths(30, 30, 30),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	workflowOwner := "workflowowner"
	orgID := "org_abc123def456ghi"
	secretID := &vaultcommon.SecretIdentifier{
		Owner:     workflowOwner,
		Namespace: "main",
		Key:       "secret",
	}
	responseID := &vaultcommon.SecretIdentifier{
		Owner:     orgID,
		Namespace: secretID.Namespace,
		Key:       secretID.Key,
	}

	obs := &vaultcommon.Observation{
		Id:          "request-1",
		RequestType: vaultcommon.RequestType_GET_SECRETS,
		Request: &vaultcommon.Observation_GetSecretsRequest{
			GetSecretsRequest: &vaultcommon.GetSecretsRequest{
				Requests: []*vaultcommon.SecretRequest{
					{Id: secretID},
				},
			},
		},
		Response: &vaultcommon.Observation_GetSecretsResponse{
			GetSecretsResponse: &vaultcommon.GetSecretsResponse{
				Responses: []*vaultcommon.SecretResponse{
					{
						Id: responseID,
						Result: &vaultcommon.SecretResponse_Error{
							Error: "not found",
						},
					},
				},
			},
		},
	}

	require.NoError(t, r.validateObservation(t.Context(), obs))
}

func TestPlugin_ValidateObservation_PanicsOnEmptyShares(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

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
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withMaxIdentifierLengths(30, 30, 30), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	seqNr := uint64(1)
	bf := &blobber{}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}

	tests := []struct {
		name        string
		obs         *vaultcommon.Observation
		expectError bool
	}{
		{
			// The request is no longer included in GetSecrets observations, so a nil Id
			// in the request field does not cause a validation error.
			name: "GetSecrets request with nil Id passes",
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
			expectError: false,
		},
		{
			name:        "GetSecrets response with nil Id",
			expectError: true,
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
			name:        "CreateSecrets with nil Id",
			expectError: true,
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
			name:        "UpdateSecrets with nil Id",
			expectError: true,
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
			name:        "CreateSecrets response with nil Id",
			expectError: true,
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
			name:        "UpdateSecrets response with nil Id",
			expectError: true,
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
			name:        "DeleteSecrets response with nil Id",
			expectError: true,
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
			if tc.expectError {
				require.Error(t, err, "expected an error for nil secret identifier, not a panic")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlugin_ValidateObservation_CiphertextSize(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	// maxCipherTextLengthBytes = 10 bytes, so any ciphertext > 10 decoded bytes should be rejected
	r := newTestReportingPlugin(
		t,
		withMaxCiphertextLengthBytes(10),
		withMaxIdentifierLengths(30, 30, 30),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

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
			errSubstr: "failed to decode encrypted value",
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
			errSubstr: "failed to decode encrypted value",
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

func TestPlugin_ValidateObservation_SecretIdentifierValidation(t *testing.T) {
	validID := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	validCiphertext := hex.EncodeToString(make([]byte, 5))

	type testCase struct {
		name      string
		obs       *vaultcommon.Observation
		errSubstr string
	}

	makeCreateSecretsObs := func(id *vaultcommon.SecretIdentifier, ciphertext string) *vaultcommon.Observation {
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: vaultcommon.RequestType_CREATE_SECRETS,
			Request: &vaultcommon.Observation_CreateSecretsRequest{
				CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{{Id: id, EncryptedValue: ciphertext}},
				},
			},
			Response: &vaultcommon.Observation_CreateSecretsResponse{
				CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
					Responses: []*vaultcommon.CreateSecretResponse{{Id: validID}},
				},
			},
		}
	}

	makeUpdateSecretsObs := func(id *vaultcommon.SecretIdentifier, ciphertext string) *vaultcommon.Observation {
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
			Request: &vaultcommon.Observation_UpdateSecretsRequest{
				UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{{Id: id, EncryptedValue: ciphertext}},
				},
			},
			Response: &vaultcommon.Observation_UpdateSecretsResponse{
				UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
					Responses: []*vaultcommon.UpdateSecretResponse{{Id: validID}},
				},
			},
		}
	}

	makeDeleteSecretsObs := func(id *vaultcommon.SecretIdentifier) *vaultcommon.Observation {
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: vaultcommon.RequestType_DELETE_SECRETS,
			Request: &vaultcommon.Observation_DeleteSecretsRequest{
				DeleteSecretsRequest: &vaultcommon.DeleteSecretsRequest{
					Ids: []*vaultcommon.SecretIdentifier{id},
				},
			},
			Response: &vaultcommon.Observation_DeleteSecretsResponse{
				DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
					Responses: []*vaultcommon.DeleteSecretResponse{{Id: validID}},
				},
			},
		}
	}

	makeListObs := func(owner, namespace string) *vaultcommon.Observation {
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
			Request: &vaultcommon.Observation_ListSecretIdentifiersRequest{
				ListSecretIdentifiersRequest: &vaultcommon.ListSecretIdentifiersRequest{
					RequestId: "request-1",
					Owner:     owner,
					Namespace: namespace,
				},
			},
			Response: &vaultcommon.Observation_ListSecretIdentifiersResponse{
				ListSecretIdentifiersResponse: &vaultcommon.ListSecretIdentifiersResponse{Success: true},
			},
		}
	}

	// makeGetSecretsObs puts the identifier only in the request; the response always uses validID.
	// Since GetSecrets observations no longer include the request, identifier validation on the
	// request side is not performed — all these cases pass regardless of the request ID.
	makeGetSecretsObs := func(id *vaultcommon.SecretIdentifier) *vaultcommon.Observation {
		req := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{"key"}}},
		}
		return &vaultcommon.Observation{
			Id:          "request-1",
			RequestType: vaultcommon.RequestType_GET_SECRETS,
			Request:     &vaultcommon.Observation_GetSecretsRequest{GetSecretsRequest: req},
			Response: &vaultcommon.Observation_GetSecretsResponse{
				GetSecretsResponse: &vaultcommon.GetSecretsResponse{
					Responses: []*vaultcommon.SecretResponse{
						{Id: validID, Result: &vaultcommon.SecretResponse_Error{Error: "err"}},
					},
				},
			},
		}
	}

	tests := []testCase{
		// --- GetSecrets ---
		// Request identifiers are not validated for GetSecrets (request is not included in the
		// observation). All cases pass regardless of the identifier in the request.
		{
			name:      "GetSecrets valid identifier passes",
			obs:       makeGetSecretsObs(validID),
			errSubstr: "",
		},
		{
			name:      "GetSecrets empty key in request passes (not validated)",
			obs:       makeGetSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: ""}),
			errSubstr: "",
		},
		{
			name:      "GetSecrets empty owner in request passes (not validated)",
			obs:       makeGetSecretsObs(&vaultcommon.SecretIdentifier{Owner: "", Namespace: "main", Key: "secret"}),
			errSubstr: "",
		},
		{
			name:      "GetSecrets empty namespace in request passes (not validated)",
			obs:       makeGetSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "", Key: "secret"}),
			errSubstr: "",
		},
		{
			name:      "GetSecrets owner too long in request passes (not validated)",
			obs:       makeGetSecretsObs(&vaultcommon.SecretIdentifier{Owner: "toolongowner", Namespace: "main", Key: "secret"}),
			errSubstr: "",
		},
		{
			name:      "GetSecrets namespace too long in request passes (not validated)",
			obs:       makeGetSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "toolongnamespace", Key: "secret"}),
			errSubstr: "",
		},
		{
			name:      "GetSecrets key too long in request passes (not validated)",
			obs:       makeGetSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "toolongkey123"}),
			errSubstr: "",
		},
		// --- CreateSecrets ---
		{
			name:      "CreateSecrets valid identifier passes",
			obs:       makeCreateSecretsObs(validID, validCiphertext),
			errSubstr: "",
		},
		{
			name:      "CreateSecrets empty key rejected",
			obs:       makeCreateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: ""}, validCiphertext),
			errSubstr: "key cannot be empty",
		},
		{
			name:      "CreateSecrets empty owner rejected",
			obs:       makeCreateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "", Namespace: "main", Key: "secret"}, validCiphertext),
			errSubstr: "owner cannot be empty",
		},
		{
			name:      "CreateSecrets empty namespace rejected",
			obs:       makeCreateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "", Key: "secret"}, validCiphertext),
			errSubstr: "namespace cannot be empty",
		},
		{
			name:      "CreateSecrets owner too long rejected",
			obs:       makeCreateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "toolongowner", Namespace: "main", Key: "secret"}, validCiphertext),
			errSubstr: "owner exceeds maximum length",
		},
		{
			name:      "CreateSecrets key too long rejected",
			obs:       makeCreateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "toolongkey123"}, validCiphertext),
			errSubstr: "key exceeds maximum length",
		},
		// --- UpdateSecrets ---
		{
			name:      "UpdateSecrets valid identifier passes",
			obs:       makeUpdateSecretsObs(validID, validCiphertext),
			errSubstr: "",
		},
		{
			name:      "UpdateSecrets empty key rejected",
			obs:       makeUpdateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: ""}, validCiphertext),
			errSubstr: "key cannot be empty",
		},
		{
			name:      "UpdateSecrets empty owner rejected",
			obs:       makeUpdateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "", Namespace: "main", Key: "secret"}, validCiphertext),
			errSubstr: "owner cannot be empty",
		},
		{
			name:      "UpdateSecrets empty namespace rejected",
			obs:       makeUpdateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "", Key: "secret"}, validCiphertext),
			errSubstr: "namespace cannot be empty",
		},
		{
			name:      "UpdateSecrets namespace too long rejected",
			obs:       makeUpdateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "toolongnamespace", Key: "secret"}, validCiphertext),
			errSubstr: "namespace exceeds maximum length",
		},
		{
			name:      "UpdateSecrets key too long rejected",
			obs:       makeUpdateSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "toolongkey123"}, validCiphertext),
			errSubstr: "key exceeds maximum length",
		},
		// --- DeleteSecrets ---
		{
			name:      "DeleteSecrets valid identifier passes",
			obs:       makeDeleteSecretsObs(validID),
			errSubstr: "",
		},
		{
			name:      "DeleteSecrets empty key rejected",
			obs:       makeDeleteSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: ""}),
			errSubstr: "key cannot be empty",
		},
		{
			name:      "DeleteSecrets empty owner rejected",
			obs:       makeDeleteSecretsObs(&vaultcommon.SecretIdentifier{Owner: "", Namespace: "main", Key: "secret"}),
			errSubstr: "owner cannot be empty",
		},
		{
			name:      "DeleteSecrets empty namespace rejected",
			obs:       makeDeleteSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "", Key: "secret"}),
			errSubstr: "namespace cannot be empty",
		},
		{
			name:      "DeleteSecrets owner too long rejected",
			obs:       makeDeleteSecretsObs(&vaultcommon.SecretIdentifier{Owner: "toolongowner", Namespace: "main", Key: "secret"}),
			errSubstr: "owner exceeds maximum length",
		},
		{
			name:      "DeleteSecrets key too long rejected",
			obs:       makeDeleteSecretsObs(&vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "toolongkey123"}),
			errSubstr: "key exceeds maximum length",
		},
		// --- ListSecretIdentifiers ---
		{
			name:      "ListSecretIdentifiers valid owner and namespace passes",
			obs:       makeListObs("owner", "main"),
			errSubstr: "",
		},
		{
			name:      "ListSecretIdentifiers empty owner rejected",
			obs:       makeListObs("", "main"),
			errSubstr: "key cannot be empty",
		},
		{
			name:      "ListSecretIdentifiers empty namespace rejected",
			obs:       makeListObs("owner", ""),
			errSubstr: "namespace cannot be empty",
		},
		{
			name:      "ListSecretIdentifiers owner too long rejected",
			obs:       makeListObs("toolongowner", "main"),
			errSubstr: "owner exceeds maximum length",
		},
		{
			name:      "ListSecretIdentifiers namespace too long rejected",
			obs:       makeListObs("owner", "toolongnamespace"),
			errSubstr: "namespace exceeds maximum length",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Use small limits (10 bytes) to trigger length errors on identifiers above.
			r := newTestReportingPlugin(
				t,
				withMaxIdentifierLengths(10, 10, 10),
				withOnchainCfg(4, 1),
			)

			rdr := &kv{m: make(map[string]response)}

			var anyp *anypb.Any
			var err error
			switch tc.obs.RequestType {
			case vaultcommon.RequestType_GET_SECRETS:
				anyp, err = anypb.New(tc.obs.GetGetSecretsRequest())
			case vaultcommon.RequestType_CREATE_SECRETS:
				anyp, err = anypb.New(tc.obs.GetCreateSecretsRequest())
			case vaultcommon.RequestType_UPDATE_SECRETS:
				anyp, err = anypb.New(tc.obs.GetUpdateSecretsRequest())
			case vaultcommon.RequestType_DELETE_SECRETS:
				anyp, err = anypb.New(tc.obs.GetDeleteSecretsRequest())
			case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
				anyp, err = anypb.New(tc.obs.GetListSecretIdentifiersRequest())
			default:
				t.FailNow()
			}
			require.NoError(t, err)

			err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
				[]*vaultcommon.StoredPendingQueueItem{{Id: "request-1", Item: anyp}},
			)
			require.NoError(t, err)

			ob := protoMarshal(t, &vaultcommon.Observations{
				Observations: []*vaultcommon.Observation{tc.obs},
			})

			err = r.ValidateObservation(
				t.Context(),
				1,
				types.AttributedQuery{},
				types.AttributedObservation{Observer: 0, Observation: ob},
				rdr,
				&blobber{},
			)

			if tc.errSubstr != "" {
				require.ErrorContains(t, err, tc.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlugin_StateTransition_PendingQueueEnabled_NewQuora_NotGetRequest(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(
		t,
		withLggr(lggr),
		withBatchSize(10),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

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
	assert.True(t, proto.Equal(resp, o.GetListSecretIdentifiersResponse()))

	ss, err := rs.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, ss)

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_PendingQueueEnabled_GetRequest(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(
		t,
		withLggr(lggr),
		withBatchSize(10),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

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

		share, err := generatePlaintextShare(pk, shares[0], ctb, owner)
		require.NoError(t, err)

		eds, err := share.encryptWithKeyBinary(hex.EncodeToString(recipientPub[:]))
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
			_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
			require.NoError(t, err)
			r := newTestReportingPlugin(
				t,
				withBatchSize(10),
				withMaxRequestBatchSize(maxRequestBatchSize),
				withMaxIdentifierLengths(30, 30, 30),
				withKeys(pk, shares[0]),
				withOnchainCfg(4, 1),
			)
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

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(
		t,
		withBatchSize(10),
		withMaxSecretsPerOwner(maxSecretsPerOwner),
		withMaxIdentifierLengths(30, 30, 30),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

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
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withMarshalBlob(marshalBlobOverride))

		fetcher := &callbackBlobFetcher{fn: func([]byte) error { return nil }}
		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, nil, nil)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("all payloads broadcast successfully", func(t *testing.T) {
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withMarshalBlob(marshalBlobOverride))

		fetcher := &callbackBlobFetcher{fn: func([]byte) error { return nil }}
		payloads := [][]byte{[]byte("p1"), []byte("p2"), []byte("p3")}
		ids := [][]string{{"req-1"}, {"req-2"}, {"req-3"}}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Len(t, result, 3)
		for _, item := range result {
			assert.Equal(t, []byte("handle"), item)
		}
	})

	t.Run("does not exceed max concurrent broadcasts", func(t *testing.T) {
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withMarshalBlob(marshalBlobOverride))

		payloads := make([][]byte, maxConcurrentBlobBroadcasts*2+1)
		ids := make([][]string, len(payloads))
		for i := range payloads {
			payloads[i] = []byte(fmt.Sprintf("payload-%d", i))
			ids[i] = []string{fmt.Sprintf("req-%d", i)}
		}

		var active atomic.Int32
		var maxActive atomic.Int32
		started := make(chan struct{}, len(payloads))
		release := make(chan struct{})
		released := atomic.Bool{}
		releaseBroadcasts := func() {
			if released.CompareAndSwap(false, true) {
				close(release)
			}
		}
		defer releaseBroadcasts()

		fetcher := &ctxCallbackBlobFetcher{fn: func(ctx context.Context, _ []byte) error {
			current := active.Add(1)
			defer active.Add(-1)

			for {
				maxSeen := maxActive.Load()
				if current <= maxSeen || maxActive.CompareAndSwap(maxSeen, current) {
					break
				}
			}

			started <- struct{}{}
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}}

		type broadcastResult struct {
			payloads [][]byte
			err      error
		}
		done := make(chan broadcastResult, 1)
		go func() {
			result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
			done <- broadcastResult{payloads: result, err: err}
		}()

		for i := 0; i < maxConcurrentBlobBroadcasts; i++ {
			select {
			case <-started:
			case <-time.After(time.Second):
				t.Fatalf("timed out waiting for broadcast %d to start", i+1)
			}
		}

		assert.Never(t, func() bool {
			return maxActive.Load() > int32(maxConcurrentBlobBroadcasts)
		}, 100*time.Millisecond, 10*time.Millisecond)

		releaseBroadcasts()

		select {
		case result := <-done:
			require.NoError(t, result.err)
			assert.Len(t, result.payloads, len(payloads))
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for broadcasts to complete")
		}
		assert.LessOrEqual(t, maxActive.Load(), int32(maxConcurrentBlobBroadcasts))
	})

	t.Run("failed broadcast is skipped and logged", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withLggr(lggr), withMarshalBlob(marshalBlobOverride))

		fetcher := &callbackBlobFetcher{fn: func(payload []byte) error {
			if string(payload) == "p2" {
				return errors.New("broadcast error")
			}
			return nil
		}}

		payloads := [][]byte{[]byte("p1"), []byte("p2"), []byte("p3")}
		ids := [][]string{{"req-1"}, {"req-2"}, {"req-3"}}

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
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withLggr(lggr), withMarshalBlob(marshalBlobOverride))

		fetcher := &errorBlobBroadcastFetcher{err: errors.New("network down")}
		payloads := [][]byte{[]byte("p1"), []byte("p2")}
		ids := [][]string{{"req-1"}, {"req-2"}}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Empty(t, result)

		warnLogs := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		assert.Equal(t, 2, warnLogs.Len())
	})

	t.Run("marshal blob failure skips item and logs warning", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return nil, errors.New("marshal error")
		}
		r := newTestReportingPlugin(t, withLggr(lggr), withMarshalBlob(marshalBlobOverride))

		fetcher := &callbackBlobFetcher{fn: func([]byte) error { return nil }}
		payloads := [][]byte{[]byte("p1"), []byte("p2")}
		ids := [][]string{{"req-1"}, {"req-2"}}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Empty(t, result)

		warnLogs := observed.FilterMessage("failed to marshal blob handle, skipping")
		assert.Equal(t, 2, warnLogs.Len())
	})

	t.Run("mix of broadcast and marshal failures", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)

		marshalCallCount := atomic.Int32{}
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			n := marshalCallCount.Add(1)
			if n == 1 {
				return nil, errors.New("marshal error")
			}
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withLggr(lggr), withMarshalBlob(marshalBlobOverride))

		fetcher := &callbackBlobFetcher{fn: func(payload []byte) error {
			if string(payload) == "p1" {
				return errors.New("broadcast error")
			}
			return nil
		}}

		payloads := [][]byte{[]byte("p1"), []byte("p2"), []byte("p3")}
		ids := [][]string{{"req-1"}, {"req-2"}, {"req-3"}}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)

		broadcastWarns := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		marshalWarns := observed.FilterMessage("failed to marshal blob handle, skipping")
		assert.Equal(t, 1, broadcastWarns.Len())
		assert.Equal(t, 1, marshalWarns.Len())
		assert.Len(t, result, 1)
	})

	t.Run("context cancellation propagates error", func(t *testing.T) {
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withMarshalBlob(marshalBlobOverride))

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		fetcher := &callbackBlobFetcher{fn: func([]byte) error {
			return ctx.Err()
		}}

		payloads := [][]byte{[]byte("p1"), []byte("p2")}
		ids := [][]string{{"req-1"}, {"req-2"}}

		result, err := r.broadcastBlobPayloads(ctx, fetcher, 1, payloads, ids)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("context deadline exceeded propagates error", func(t *testing.T) {
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withMarshalBlob(marshalBlobOverride))

		ctx, cancel := context.WithTimeout(t.Context(), 0)
		defer cancel()
		<-ctx.Done()

		fetcher := &callbackBlobFetcher{fn: func([]byte) error {
			return ctx.Err()
		}}

		payloads := [][]byte{[]byte("p1")}
		ids := [][]string{{"req-1"}}

		result, err := r.broadcastBlobPayloads(ctx, fetcher, 1, payloads, ids)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("slow broadcast hits per-call timeout and is skipped", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		marshalBlobOverride := func(ocr3_1types.BlobHandle) ([]byte, error) {
			return []byte("handle"), nil
		}
		r := newTestReportingPlugin(t, withLggr(lggr), withMarshalBlob(marshalBlobOverride))

		fetcher := &ctxCallbackBlobFetcher{fn: func(ctx context.Context, payload []byte) error {
			if string(payload) == "slow" {
				<-ctx.Done()
				return ctx.Err()
			}
			return nil
		}}

		payloads := [][]byte{[]byte("fast"), []byte("slow")}
		ids := [][]string{{"req-fast"}, {"req-slow"}}

		result, err := r.broadcastBlobPayloads(t.Context(), fetcher, 1, payloads, ids)
		require.NoError(t, err)
		assert.Len(t, result, 1)

		warnLogs := observed.FilterMessage("failed to broadcast pending queue item as blob, skipping")
		assert.Equal(t, 1, warnLogs.Len())
		fields := warnLogs.All()[0].ContextMap()
		assert.Equal(t, "req-slow", fields["requestID"])
	})
}

func TestProperty_broadcastBlobPayloads_MaxSizePayloadsWithinBlobLimit(t *testing.T) {
	maxRequestBatchSize := cresettings.Default.VaultRequestBatchSizeLimit.DefaultValue
	maxCiphertextBytes := cresettings.Default.VaultCiphertextSizeLimit.DefaultValue
	maxIDKeySize := cresettings.Default.VaultIdentifierKeySizeLimit.DefaultValue
	maxIDOwnerSize := cresettings.Default.VaultIdentifierOwnerSizeLimit.DefaultValue
	maxIDNamespaceSize := cresettings.Default.VaultIdentifierNamespaceSizeLimit.DefaultValue
	maxSecretsPerReq := vaulttypes.MaxBatchSize
	maxBlobPayloadBytes := cresettings.Default.VaultMaxBlobPayloadSizeLimit.DefaultValue

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	encKey := hex.EncodeToString(pubK[:])

	maxIDKeyField := strings.Repeat("a", int(maxIDKeySize))
	maxIDOwnerField := strings.Repeat("a", int(maxIDOwnerSize))
	maxIDNamespaceField := strings.Repeat("a", int(maxIDNamespaceSize))
	maxCiphertext := strings.Repeat("a", int(maxCiphertextBytes))

	maxIdentifier := func() *vaultcommon.SecretIdentifier {
		return &vaultcommon.SecretIdentifier{
			Owner:     maxIDOwnerField,
			Namespace: maxIDNamespaceField,
			Key:       maxIDKeyField,
		}
	}

	buildMaxEncryptedSecrets := func() []*vaultcommon.EncryptedSecret {
		secs := make([]*vaultcommon.EncryptedSecret, maxRequestBatchSize)
		for i := range secs {
			secs[i] = &vaultcommon.EncryptedSecret{
				Id:             maxIdentifier(),
				EncryptedValue: maxCiphertext,
			}
		}
		return secs
	}

	buildMaxSecretRequests := func() []*vaultcommon.SecretRequest {
		reqs := make([]*vaultcommon.SecretRequest, maxSecretsPerReq)
		encKeys := make([]string, 10)
		for i := range encKeys {
			encKeys[i] = encKey
		}
		for i := range reqs {
			reqs[i] = &vaultcommon.SecretRequest{
				Id:             maxIdentifier(),
				EncryptionKeys: encKeys,
			}
		}
		return reqs
	}

	buildMaxIdentifiers := func() []*vaultcommon.SecretIdentifier {
		ids := make([]*vaultcommon.SecretIdentifier, maxSecretsPerReq)
		for i := range ids {
			ids[i] = maxIdentifier()
		}
		return ids
	}

	requestTypes := []struct {
		name    string
		payload proto.Message
	}{
		{
			name: "GetSecretsRequest",
			payload: &vaultcommon.GetSecretsRequest{
				Requests: buildMaxSecretRequests(),
			},
		},
		{
			name: "CreateSecretsRequest",
			payload: &vaultcommon.CreateSecretsRequest{
				RequestId:        "req",
				EncryptedSecrets: buildMaxEncryptedSecrets(),
			}},
		{
			name: "UpdateSecretsRequest",
			payload: &vaultcommon.UpdateSecretsRequest{
				RequestId:        "req",
				EncryptedSecrets: buildMaxEncryptedSecrets(),
			},
		},
		{
			name: "DeleteSecretsRequest",
			payload: &vaultcommon.DeleteSecretsRequest{
				RequestId: "req",
				Ids:       buildMaxIdentifiers(),
			},
		},
		{
			name: "ListSecretIdentifiersRequest",
			payload: &vaultcommon.ListSecretIdentifiersRequest{
				RequestId: "req",
				Owner:     maxIDOwnerField,
				Namespace: maxIDNamespaceField,
			},
		},
	}

	for _, rt := range requestTypes {
		t.Run(rt.name, func(t *testing.T) {
			anyMsg, err := anypb.New(rt.payload)
			require.NoError(t, err)

			item := &vaultcommon.StoredPendingQueueItem{
				Id:   "req",
				Item: anyMsg,
			}
			itemBytes := protoMarshal(t, item)

			assert.LessOrEqualf(t, len(itemBytes), maxBlobPayloadBytes,
				"marshaled %s StoredPendingQueueItem (%d bytes) exceeds VaultMaxBlobPayloadSizeLimit (%d bytes)",
				rt.name, len(itemBytes), maxBlobPayloadBytes)
		})
	}
}
