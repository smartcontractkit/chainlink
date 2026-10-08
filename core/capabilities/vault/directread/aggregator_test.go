package directread

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const n, f = 4, 1 // quorum 3

var (
	idA = &vaultcommon.SecretIdentifier{Owner: "0xowner", Namespace: "main", Key: "a"}
	idB = &vaultcommon.SecretIdentifier{Owner: "0xowner", Namespace: "main", Key: "b"}
)

func newReq(ids ...*vaultcommon.SecretIdentifier) *vaultcommon.GetSecretsRequest {
	req := &vaultcommon.GetSecretsRequest{GetSecretsDirectly: true}
	for _, id := range ids {
		req.Requests = append(req.Requests, &vaultcommon.SecretRequest{Id: id, EncryptionKeys: []string{"k1", "k2"}})
	}
	return req
}

func data(id *vaultcommon.SecretIdentifier, ct string, node byte) *vaultcommon.SecretResponse {
	return &vaultcommon.SecretResponse{Id: id, Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
		EncryptedValue: ct,
		EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
			{EncryptionKey: "k1", BinaryShares: [][]byte{{1, node}}},
			{EncryptionKey: "k2", BinaryShares: [][]byte{{2, node}}},
		},
	}}}
}

func itemErr(id *vaultcommon.SecretIdentifier, msg string) *vaultcommon.SecretResponse {
	return &vaultcommon.SecretResponse{Id: id, Result: &vaultcommon.SecretResponse_Error{Error: msg}}
}

func capResp(t *testing.T, rs ...*vaultcommon.SecretResponse) commoncap.CapabilityResponse {
	p, err := anypb.New(&vaultcommon.GetSecretsResponse{Responses: rs})
	require.NoError(t, err)
	return commoncap.CapabilityResponse{Payload: p}
}

func capRespWithKey(t *testing.T, publicKey string, rs ...*vaultcommon.SecretResponse) commoncap.CapabilityResponse {
	p, err := anypb.New(&vaultcommon.GetSecretsResponse{Responses: rs, RawVaultPublicKey: publicKey})
	require.NoError(t, err)
	return commoncap.CapabilityResponse{Payload: p}
}

func decode(t *testing.T, resp *commoncap.CapabilityResponse) *vaultcommon.GetSecretsResponse {
	require.NotNil(t, resp)
	out := &vaultcommon.GetSecretsResponse{}
	require.NoError(t, resp.Payload.UnmarshalTo(out))
	return out
}

func TestAggregator_SucceedsAtQuorumAndMergesShares(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA, idB), n, f)

	for node := range byte(2) {
		resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node), itemErr(idB, "key does not exist")))
		require.NoError(t, err)
		require.Nil(t, resp, "undecided before 2F+1 matching responses")
	}
	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 2), itemErr(idB, "key does not exist")))
	require.NoError(t, err)

	out := decode(t, resp)
	require.Len(t, out.Responses, 2)
	shares := out.Responses[0].GetData().GetEncryptedDecryptionKeyShares()
	require.Len(t, shares, 2, "one entry per encryption key")
	assert.Equal(t, "k1", shares[0].EncryptionKey)
	assert.Equal(t, [][]byte{{1, 0}, {1, 1}, {1, 2}}, shares[0].BinaryShares)
	assert.Equal(t, [][]byte{{2, 0}, {2, 1}, {2, 2}}, shares[1].BinaryShares)
	assert.Equal(t, "ctA", out.Responses[0].GetData().GetEncryptedValue())
	assert.Equal(t, "key does not exist", out.Responses[1].GetError())
}

func TestAggregator_ToleratesFFaultyResponses(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), n, f)

	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "forged", 9)))
	require.NoError(t, err)
	require.Nil(t, resp)
	for node := range byte(3) {
		resp, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node)))
		require.NoError(t, err)
	}
	out := decode(t, resp)
	assert.Equal(t, "ctA", out.Responses[0].GetData().GetEncryptedValue())
	assert.Len(t, out.Responses[0].GetData().GetEncryptedDecryptionKeyShares()[0].BinaryShares, 3, "the faulty node's share is excluded")
}

func TestAggregator_SplitVersionsIsSkew(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), n, f)

	_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v1", 0)))
	require.NoError(t, err)
	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v2", 1)))
	require.NoError(t, err)
	require.Nil(t, resp, "2 pending replies could still make either version reach 3")

	_, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v2", 2)))
	require.NoError(t, err, "v2 can still reach quorum")

	_, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v1", 3)))
	require.ErrorIs(t, err, vaultcommon.ErrSecretVersionSkew)
	assert.True(t, vaultcommon.IsSecretVersionSkew(err))
}

func TestAggregator_ReturnsVaultPublicKey(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), n, f)

	var resp *commoncap.CapabilityResponse
	for node := range byte(3) {
		var err error
		resp, err = a.OnResponse(p2ptypes.PeerID{}, capRespWithKey(t, "pk1", data(idA, "ctA", node)))
		require.NoError(t, err)
	}
	assert.Equal(t, "pk1", decode(t, resp).RawVaultPublicKey)
}

func TestAggregator_DifferentVaultPublicKeysAreSkew(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), n, f)

	for node, pk := range []string{"old", "old", "new"} {
		resp, err := a.OnResponse(p2ptypes.PeerID{}, capRespWithKey(t, pk, data(idA, "ctA", byte(node))))
		require.NoError(t, err)
		require.Nil(t, resp)
	}
	_, err := a.OnResponse(p2ptypes.PeerID{}, capRespWithKey(t, "new", data(idA, "ctA", 3)))
	require.ErrorIs(t, err, vaultcommon.ErrSecretVersionSkew, "same ciphertext from different DKG instances must not combine")
}

func TestAggregator_ItemErrorsIgnoreVaultPublicKey(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), n, f)

	var resp *commoncap.CapabilityResponse
	for _, pk := range []string{"old", "new", "new"} {
		var err error
		resp, err = a.OnResponse(p2ptypes.PeerID{}, capRespWithKey(t, pk, itemErr(idA, "key does not exist")))
		require.NoError(t, err)
	}
	out := decode(t, resp)
	assert.Equal(t, "key does not exist", out.Responses[0].GetError())
	assert.Empty(t, out.RawVaultPublicKey)
}

func TestAggregator_SkewFailsEarly(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), 7, 2) // quorum 5

	for i, ct := range []string{"v1", "v2", "v1"} {
		_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, ct, byte(i))))
		require.NoError(t, err)
	}
	_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v2", 3)))
	require.NoError(t, err)
	// Best is now 2 + 2 pending < 5.
	_, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v3", 4)))
	require.ErrorIs(t, err, vaultcommon.ErrSecretVersionSkew)
}

func TestAggregator_IdentifierIsPartOfTheVote(t *testing.T) {
	t.Parallel()
	a := NewAggregator(newReq(idA), n, f)
	other := &vaultcommon.SecretIdentifier{Owner: "0xowner", Namespace: "main", Key: "other"}

	_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(other, "ctA", 0)))
	require.NoError(t, err)
	_, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 1)))
	require.NoError(t, err)
	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 2)))
	require.NoError(t, err)
	require.Nil(t, resp)
	resp, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 3)))
	require.NoError(t, err)
	assert.Equal(t, idA.Key, decode(t, resp).Responses[0].GetId().GetKey())
}

func TestAggregator_PeerErrors(t *testing.T) {
	t.Parallel()
	t.Run("identical request errors are returned", func(t *testing.T) {
		t.Parallel()
		a := NewAggregator(newReq(idA), n, f)
		msg := caperrors.NewPublicUserError(assert.AnError, caperrors.InvalidArgument).SerializeToString()

		_, err := a.OnError(p2ptypes.PeerID{}, msg)
		require.NoError(t, err)
		_, err = a.OnError(p2ptypes.PeerID{}, msg)
		require.Error(t, err, "success is no longer possible")
		var capErr caperrors.Error
		require.ErrorAs(t, err, &capErr)
		assert.Equal(t, caperrors.InvalidArgument, capErr.Code())
	})

	t.Run("mixed failures without disagreement", func(t *testing.T) {
		t.Parallel()
		a := NewAggregator(newReq(idA), n, f)

		_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 0)))
		require.NoError(t, err)
		_, err = a.OnError(p2ptypes.PeerID{}, "vault direct read path is not ready")
		require.NoError(t, err)
		_, err = a.OnError(p2ptypes.PeerID{}, "some other error")
		require.ErrorIs(t, err, ErrNotEnoughMatchingResponses)
		assert.False(t, vaultcommon.IsSecretVersionSkew(err))
	})

	t.Run("malformed responses count as errors", func(t *testing.T) {
		t.Parallel()
		a := NewAggregator(newReq(idA, idB), n, f)

		_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 0))) // missing idB
		require.NoError(t, err)
		_, err = a.OnResponse(p2ptypes.PeerID{}, commoncap.CapabilityResponse{})
		require.ErrorIs(t, err, ErrNotEnoughMatchingResponses)
	})
}

func TestAggregatorFactory(t *testing.T) {
	t.Parallel()
	factory := NewAggregatorFactory()
	don := commoncap.DON{Members: make([]p2ptypes.PeerID, n), F: f}

	direct, err := anypb.New(newReq(idA))
	require.NoError(t, err)
	assert.NotNil(t, factory(commoncap.CapabilityRequest{Method: vaulttypes.MethodSecretsGet, Payload: direct}, don))

	viaOCR, err := anypb.New(&vaultcommon.GetSecretsRequest{Requests: newReq(idA).Requests})
	require.NoError(t, err)
	assert.Nil(t, factory(commoncap.CapabilityRequest{Method: vaulttypes.MethodSecretsGet, Payload: viaOCR}, don))

	assert.Nil(t, factory(commoncap.CapabilityRequest{Method: "other", Payload: direct}, don))
	assert.Nil(t, factory(commoncap.CapabilityRequest{Method: vaulttypes.MethodSecretsGet}, don))
}
