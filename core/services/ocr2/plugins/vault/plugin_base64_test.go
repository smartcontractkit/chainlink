package vault

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestHexEncodeGetSecretsResponseForReport(t *testing.T) {
	r := newTestReportingPlugin(t)
	raw := []byte{1, 2, 3, 4, 5}
	b64ev := base64.StdEncoding.EncodeToString(raw)
	b64share := base64.StdEncoding.EncodeToString([]byte{9, 9, 9})
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: &vaultcommon.SecretIdentifier{Key: "k", Namespace: "n", Owner: "o"},
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: b64ev,
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{EncryptionKey: "pk", Shares: []string{b64share}},
						},
					},
				},
			},
		},
	}
	got, err := r.hexEncodeGetSecretsResponseForReport(resp)
	require.NoError(t, err)
	require.Equal(t, b64ev, resp.Responses[0].GetData().EncryptedValue)
	require.Equal(t, b64share, resp.Responses[0].GetData().EncryptedDecryptionKeyShares[0].Shares[0])
	require.Equal(t, hex.EncodeToString(raw), got.Responses[0].GetData().EncryptedValue)
	require.Equal(t, hex.EncodeToString([]byte{9, 9, 9}), got.Responses[0].GetData().EncryptedDecryptionKeyShares[0].Shares[0])
}

func TestPlugin_Observation_GetSecrets_Base64EncodingEnabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withBase64EncodingEnabled())

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

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

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{Id: id, EncryptionKeys: []string{pks}},
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

	batchResp := obs.Observations[0].GetGetSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	resp := batchResp.Responses[0].GetData()
	require.NotNil(t, resp)

	_, err = base64.StdEncoding.DecodeString(resp.EncryptedValue)
	require.NoError(t, err, "EncryptedValue should be base64 when flag is enabled")
	require.Len(t, resp.EncryptedDecryptionKeyShares, 1)
	_, err = base64.StdEncoding.DecodeString(resp.EncryptedDecryptionKeyShares[0].Shares[0])
	require.NoError(t, err, "share should be base64 when flag is enabled")
}

func TestStateTransitionCreateSecretsRequest_DecodesBase64EncryptedValue(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withBase64EncodingEnabled())

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "k"}
	plaintext := []byte("secret")
	var label [32]byte
	copy(label[12:], common.HexToAddress(owner).Bytes())
	ct, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	rawCipher, err := ct.Marshal()
	require.NoError(t, err)
	b64 := base64.StdEncoding.EncodeToString(rawCipher)

	rdr := &kv{m: make(map[string]response)}
	store := newTestWriteStore(t, rdr)

	resp := &vaultcommon.CreateSecretResponse{Id: id, Success: false, Error: ""}
	out, err := r.stateTransitionCreateSecretsRequest(t.Context(), store, &vaultcommon.EncryptedSecret{
		Id: id, EncryptedValue: b64,
	}, resp, "")
	require.NoError(t, err)
	require.True(t, out.Success)

	stored, err := store.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.Equal(t, rawCipher, stored.EncryptedSecret)
}

func TestStateTransitionUpdateSecretsRequest_DecodesBase64EncryptedValue(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withBase64EncodingEnabled())

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "k"}
	plaintext := []byte("original")
	var label [32]byte
	copy(label[12:], common.HexToAddress(owner).Bytes())
	ct, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	rawCipherOriginal, err := ct.Marshal()
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	store := newTestWriteStore(t, rdr)
	require.NoError(t, store.WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{EncryptedSecret: rawCipherOriginal}))

	plaintext2 := []byte("updated")
	ct2, err := tdh2easy.EncryptWithLabel(pk, plaintext2, label)
	require.NoError(t, err)
	rawCipherNew, err := ct2.Marshal()
	require.NoError(t, err)
	b64 := base64.StdEncoding.EncodeToString(rawCipherNew)

	resp := &vaultcommon.UpdateSecretResponse{Id: id, Success: false, Error: ""}
	out, err := r.stateTransitionUpdateSecretsRequest(t.Context(), store, &vaultcommon.EncryptedSecret{
		Id: id, EncryptedValue: b64,
	}, resp, "")
	require.NoError(t, err)
	require.True(t, out.Success)

	stored, err := store.GetSecret(t.Context(), id)
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.Equal(t, rawCipherNew, stored.EncryptedSecret)
}

func TestPlugin_Reports_GetSecrets_Base64Outcome_NormalizesToHexInReport(t *testing.T) {
	rawCipher := []byte{1, 2, 3, 10, 11}
	rawShare1 := []byte{9}
	rawShare2 := []byte{8, 7}
	rawCipher2 := []byte{5, 6}

	id := &vaultcommon.SecretIdentifier{Owner: "o", Namespace: "main", Key: "a"}
	id2 := &vaultcommon.SecretIdentifier{Owner: "o", Namespace: "main", Key: "b"}
	id3 := &vaultcommon.SecretIdentifier{Owner: "o", Namespace: "main", Key: "c"}

	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: base64.StdEncoding.EncodeToString(rawCipher),
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "pk1",
								Shares: []string{
									base64.StdEncoding.EncodeToString(rawShare1),
									base64.StdEncoding.EncodeToString(rawShare2),
								},
							},
						},
					},
				},
			},
			{
				Id: id2,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: base64.StdEncoding.EncodeToString(rawCipher2),
					},
				},
			},
			{
				Id:     id3,
				Result: &vaultcommon.SecretResponse_Error{Error: "boom"},
			},
		},
	}

	req := &vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id}}}
	out := &vaultcommon.Outcome{
		Id:          vaulttypes.KeyFor(id),
		RequestType: vaultcommon.RequestType_GET_SECRETS,
		Request:     &vaultcommon.Outcome_GetSecretsRequest{GetSecretsRequest: req},
		Response:    &vaultcommon.Outcome_GetSecretsResponse{GetSecretsResponse: resp},
	}
	os := &vaultcommon.Outcomes{Outcomes: []*vaultcommon.Outcome{out}}
	osb, err := proto.Marshal(os)
	require.NoError(t, err)

	rp := newTestReportingPlugin(t, withBase64EncodingEnabled())
	reports, err := rp.Reports(t.Context(), 1, osb)
	require.NoError(t, err)
	require.Len(t, reports, 1)

	info, err := extractReportInfo(reports[0].ReportWithInfo)
	require.NoError(t, err)
	require.Equal(t, vaultcommon.ReportFormat_REPORT_FORMAT_PROTOBUF, info.Format)
	require.Equal(t, vaultcommon.RequestType_GET_SECRETS, info.RequestType)

	got := &vaultcommon.GetSecretsResponse{}
	require.NoError(t, proto.Unmarshal(reports[0].ReportWithInfo.Report, got))
	require.Len(t, got.Responses, 3)

	d0 := got.Responses[0].GetData()
	require.NotNil(t, d0)
	dec0, err := hex.DecodeString(d0.EncryptedValue)
	require.NoError(t, err)
	require.Equal(t, rawCipher, dec0)
	require.Len(t, d0.EncryptedDecryptionKeyShares, 1)
	sh := d0.EncryptedDecryptionKeyShares[0].Shares
	require.Len(t, sh, 2)
	s0, err := hex.DecodeString(sh[0])
	require.NoError(t, err)
	require.Equal(t, rawShare1, s0)
	s1, err := hex.DecodeString(sh[1])
	require.NoError(t, err)
	require.Equal(t, rawShare2, s1)

	d1 := got.Responses[1].GetData()
	require.NotNil(t, d1)
	dec1, err := hex.DecodeString(d1.EncryptedValue)
	require.NoError(t, err)
	require.Equal(t, rawCipher2, dec1)

	require.Equal(t, "boom", got.Responses[2].GetError())
}
