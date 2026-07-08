package vault

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"

	libocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPlugin_ValidateObservation_GetSecrets_BogusShareLabelRejected(t *testing.T) {
	t.Parallel()
	_, vaultPub, vaultShares, err := tdh2easy.GenerateKeys(2, 4)
	require.NoError(t, err)

	encKey := strings.Repeat("ab", 32)
	id := &vaultcommon.SecretIdentifier{Owner: "52bc44d5378309ee2abf1539bf71de1b7d7be3b5", Namespace: "main", Key: "mysecret"}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{encKey}}},
	}

	esHex, err := vaultutils.EncryptSecretWithWorkflowOwner("secret", vaultPub, common.HexToAddress(id.Owner))
	require.NoError(t, err)
	esBytes, err := hex.DecodeString(esHex)
	require.NoError(t, err)

	honestResp := buildHonestGetSecretsResponse(t, id, esHex, vaultPub, vaultShares[0], []string{encKey}, esBytes, id.Owner)
	byzResp := proto.Clone(honestResp).(*vaultcommon.GetSecretsResponse)
	byzResp.Responses[0].GetData().EncryptedDecryptionKeyShares[0].EncryptionKey = strings.Repeat("ba", 32)

	rdr := &kv{m: make(map[string]response)}
	writeGetSecretsPendingQueueItem(t, rdr, vaulttypes.KeyFor(id), req)

	plugin := newTestReportingPlugin(t, withKeys(vaultPub, vaultShares[0]), withOnchainCfg(4, 1))
	byzObs := marshalObservations(t, observation{id, req, byzResp})

	err = plugin.ValidateObservation(
		context.Background(),
		1,
		libocrtypes.AttributedQuery{},
		libocrtypes.AttributedObservation{Observer: 0, Observation: libocrtypes.Observation(byzObs)},
		rdr,
		&blobber{},
	)
	require.ErrorContains(t, err, "unexpected encryption key in response")
}

func TestPlugin_ValidateObservation_GetSecrets_EmbeddedRequestMismatchRejected(t *testing.T) {
	t.Parallel()
	_, vaultPub, vaultShares, err := tdh2easy.GenerateKeys(2, 4)
	require.NoError(t, err)

	realKey := strings.Repeat("ab", 32)
	bogusKey := strings.Repeat("ba", 32)
	id := &vaultcommon.SecretIdentifier{Owner: "52bc44d5378309ee2abf1539bf71de1b7d7be3b5", Namespace: "main", Key: "mysecret"}
	pendingReq := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{realKey}}},
	}
	embeddedReq := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{bogusKey}}},
	}

	esHex, err := vaultutils.EncryptSecretWithWorkflowOwner("secret", vaultPub, common.HexToAddress(id.Owner))
	require.NoError(t, err)
	esBytes, err := hex.DecodeString(esHex)
	require.NoError(t, err)

	resp := buildHonestGetSecretsResponse(t, id, esHex, vaultPub, vaultShares[0], []string{bogusKey}, esBytes, id.Owner)
	rdr := &kv{m: make(map[string]response)}
	writeGetSecretsPendingQueueItem(t, rdr, vaulttypes.KeyFor(id), pendingReq)

	plugin := newTestReportingPlugin(t, withKeys(vaultPub, vaultShares[0]), withOnchainCfg(4, 1))
	obs := marshalObservations(t, observation{id, embeddedReq, resp})

	err = plugin.ValidateObservation(
		context.Background(),
		1,
		libocrtypes.AttributedQuery{},
		libocrtypes.AttributedObservation{Observer: 0, Observation: libocrtypes.Observation(obs)},
		rdr,
		&blobber{},
	)
	require.ErrorContains(t, err, "embedded GetSecrets request does not match pending queue request")
}

func TestPlugin_ShaForObservation_ShareAggregationIncludesPublicKeysFlag(t *testing.T) {
	t.Parallel()
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	realKey := strings.Repeat("ab", 32)
	bogusKey := strings.Repeat("ba", 32)
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{realKey}}},
	}
	honestResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: id,
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: "encrypted-value",
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
						EncryptionKey: realKey,
						Shares:        []string{"share"},
					}},
				},
			},
		}},
	}
	byzResp := proto.Clone(honestResp).(*vaultcommon.GetSecretsResponse)
	byzResp.Responses[0].GetData().EncryptedDecryptionKeyShares[0].EncryptionKey = bogusKey

	honestObs := &vaultcommon.Observation{
		Id:          vaulttypes.KeyFor(id),
		RequestType: vaultcommon.RequestType_GET_SECRETS,
		Request:     &vaultcommon.Observation_GetSecretsRequest{GetSecretsRequest: req},
		Response:    &vaultcommon.Observation_GetSecretsResponse{GetSecretsResponse: honestResp},
	}
	byzObs := proto.Clone(honestObs).(*vaultcommon.Observation)
	byzObs.Response = &vaultcommon.Observation_GetSecretsResponse{GetSecretsResponse: byzResp}

	ctx := context.Background()
	pluginOff := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	shaHonestOff, err := pluginOff.shaForObservation(ctx, honestObs)
	require.NoError(t, err)
	shaByzOff, err := pluginOff.shaForObservation(ctx, byzObs)
	require.NoError(t, err)
	require.Equal(t, shaHonestOff, shaByzOff)

	pluginOn := newTestReportingPlugin(t, withOnchainCfg(4, 1), withVaultGetSecretsShareAggregationIncludesPublicKeys())
	shaHonestOn, err := pluginOn.shaForObservation(ctx, honestObs)
	require.NoError(t, err)
	shaByzOn, err := pluginOn.shaForObservation(ctx, byzObs)
	require.NoError(t, err)
	require.NotEqual(t, shaHonestOn, shaByzOn)
}

func TestPlugin_ShaForObservation_ShareAggregationIncludesPublicKeys_PermutedEntryOrder(t *testing.T) {
	t.Parallel()
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	keyA := strings.Repeat("aa", 32)
	keyB := strings.Repeat("bb", 32)
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{keyA, keyB}}},
	}

	makeObs := func(order []string, shares []string) *vaultcommon.Observation {
		entries := make([]*vaultcommon.EncryptedShares, len(order))
		for i, key := range order {
			entries[i] = &vaultcommon.EncryptedShares{EncryptionKey: key, Shares: []string{shares[i]}}
		}
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue:               "encrypted-value",
						EncryptedDecryptionKeyShares: entries,
					},
				},
			}},
		}
		return &vaultcommon.Observation{
			Id:          vaulttypes.KeyFor(id),
			RequestType: vaultcommon.RequestType_GET_SECRETS,
			Request:     &vaultcommon.Observation_GetSecretsRequest{GetSecretsRequest: req},
			Response:    &vaultcommon.Observation_GetSecretsResponse{GetSecretsResponse: resp},
		}
	}

	ctx := context.Background()
	plugin := newTestReportingPlugin(t, withOnchainCfg(4, 1), withVaultGetSecretsShareAggregationIncludesPublicKeys())

	shaAB, err := plugin.shaForObservation(ctx, makeObs([]string{keyA, keyB}, []string{"share-a1", "share-b1"}))
	require.NoError(t, err)
	shaBA, err := plugin.shaForObservation(ctx, makeObs([]string{keyB, keyA}, []string{"share-b2", "share-a2"}))
	require.NoError(t, err)
	require.NotEqual(t, shaAB, shaBA)
}

func TestPlugin_ShaForObservation_ShareAggregationIncludesPublicKeys_DifferentShareBytesSameLabels(t *testing.T) {
	t.Parallel()
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	encKey := strings.Repeat("ab", 32)
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{encKey}}},
	}

	makeObs := func(share string) *vaultcommon.Observation {
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
							EncryptionKey: encKey,
							Shares:        []string{share},
						}},
					},
				},
			}},
		}
		return &vaultcommon.Observation{
			Id:          vaulttypes.KeyFor(id),
			RequestType: vaultcommon.RequestType_GET_SECRETS,
			Request:     &vaultcommon.Observation_GetSecretsRequest{GetSecretsRequest: req},
			Response:    &vaultcommon.Observation_GetSecretsResponse{GetSecretsResponse: resp},
		}
	}

	ctx := context.Background()
	plugin := newTestReportingPlugin(t, withOnchainCfg(4, 1), withVaultGetSecretsShareAggregationIncludesPublicKeys())

	sha1, err := plugin.shaForObservation(ctx, makeObs("share-from-node-1"))
	require.NoError(t, err)
	sha2, err := plugin.shaForObservation(ctx, makeObs("share-from-node-2"))
	require.NoError(t, err)
	require.Equal(t, sha1, sha2)
}

func TestPlugin_StateTransition_GetSecretsRequest_ShareAggregationIncludesPublicKeys_CombinesShares(t *testing.T) {
	t.Parallel()
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withLggr(lggr),
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
		withVaultGetSecretsShareAggregationIncludesPublicKeys(),
	)

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	keyA := strings.Repeat("aa", 32)
	keyB := strings.Repeat("bb", 32)
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{keyA, keyB}}},
	}

	makeResp := func(keyOrder []string, shareA, shareB string) *vaultcommon.GetSecretsResponse {
		entries := map[string]*vaultcommon.EncryptedShares{
			keyA: {EncryptionKey: keyA, Shares: []string{shareA}},
			keyB: {EncryptionKey: keyB, Shares: []string{shareB}},
		}
		ordered := make([]*vaultcommon.EncryptedShares, len(keyOrder))
		for i, key := range keyOrder {
			ordered[i] = entries[key]
		}
		return &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue:               "encrypted-value",
						EncryptedDecryptionKeyShares: ordered,
					},
				},
			}},
		}
	}

	obsb1 := marshalObservations(t, observation{id, req, makeResp([]string{keyA, keyB}, "share-a1", "share-b1")})
	obsb2 := marshalObservations(t, observation{id, req, makeResp([]string{keyA, keyB}, "share-a2", "share-b2")})
	obsb3 := marshalObservations(t, observation{id, req, makeResp([]string{keyA, keyB}, "share-a3", "share-b3")})

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		libocrtypes.AttributedQuery{},
		[]libocrtypes.AttributedObservation{
			{Observer: 0, Observation: libocrtypes.Observation(obsb1)},
			{Observer: 1, Observation: libocrtypes.Observation(obsb2)},
			{Observer: 2, Observation: libocrtypes.Observation(obsb3)},
		},
		&kv{m: make(map[string]response)},
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)

	expectedResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: id,
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: "encrypted-value",
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
						{EncryptionKey: keyA, Shares: []string{"share-a1", "share-a2", "share-a3"}},
						{EncryptionKey: keyB, Shares: []string{"share-b1", "share-b2", "share-b3"}},
					},
				},
			},
		}},
	}
	assert.True(t, proto.Equal(expectedResp, os.Outcomes[0].GetGetSecretsResponse()))
	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func buildHonestGetSecretsResponse(
	t *testing.T,
	id *vaultcommon.SecretIdentifier,
	esHex string,
	pub *tdh2easy.PublicKey,
	share *tdh2easy.PrivateShare,
	encKeys []string,
	esBytes []byte,
	owner string,
) *vaultcommon.GetSecretsResponse {
	t.Helper()
	entries := make([]*vaultcommon.EncryptedShares, len(encKeys))
	for j, encKey := range encKeys {
		sh, err := generatePlaintextShare(pub, share, esBytes, owner)
		require.NoError(t, err)
		enc, err := sh.encryptWithKeyBinary(encKey)
		require.NoError(t, err)
		entries[j] = &vaultcommon.EncryptedShares{
			EncryptionKey: encKey,
			Shares:        []string{hex.EncodeToString(enc)},
		}
	}
	return &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: id,
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue:               esHex,
					EncryptedDecryptionKeyShares: entries,
				},
			},
		}},
	}
}
