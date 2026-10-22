package directread

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const n, f = 4, 1 // quorum 3

// shareSizeLimit matches the limit the launcher wires into the aggregator factory.
var shareSizeLimit = func() limits.BoundLimiter[commonconfig.Size] {
	l, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, cresettings.Default.VaultShareSizeLimit)
	if err != nil {
		panic(err)
	}
	return l
}()

func newAggregator(req *vaultcommon.GetSecretsRequest, n, f int) *Aggregator {
	return NewAggregator(req, n, f, shareSizeLimit)
}

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
	a := newAggregator(newReq(idA, idB), n, f)

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
	a := newAggregator(newReq(idA), n, f)

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
	a := newAggregator(newReq(idA), n, f)

	_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v1", 0)))
	require.NoError(t, err)
	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v2", 1)))
	require.NoError(t, err)
	require.Nil(t, resp, "2 pending replies could still make either version reach 3")

	_, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v2", 2)))
	require.NoError(t, err, "v2 can still reach quorum")

	_, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "v1", 3)))
	require.ErrorIs(t, err, vaultcommon.ErrSecretVersionSkew)
}

func TestAggregator_ReturnsVaultPublicKey(t *testing.T) {
	t.Parallel()
	a := newAggregator(newReq(idA), n, f)

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
	a := newAggregator(newReq(idA), n, f)

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
	a := newAggregator(newReq(idA), n, f)

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
	a := newAggregator(newReq(idA), 7, 2) // quorum 5

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
	a := newAggregator(newReq(idA), n, f)
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

func TestAggregator_MergedResponseIsLabelledWithRequestedIdentifier(t *testing.T) {
	t.Parallel()
	a := newAggregator(newReq(idA), n, f)

	// A faulty node answers first with a blank namespace. KeyFor normalizes it to
	// "main", so its answer joins the honest vote bucket and counts toward quorum.
	blankNS := &vaultcommon.SecretIdentifier{Owner: "0xowner", Namespace: "", Key: "a"}
	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(blankNS, "ctA", 9)))
	require.NoError(t, err)
	require.Nil(t, resp)

	for node := range byte(3) {
		resp, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node)))
		require.NoError(t, err)
	}

	out := decode(t, resp)
	require.Len(t, out.Responses, 1)
	assert.True(t, proto.Equal(idA, out.Responses[0].GetId()), "merged response must carry the requested identifier, not the first responder's")
}

func TestAggregator_RejectsResponseWithDifferentIdentifier(t *testing.T) {
	t.Parallel()
	a := newAggregator(newReq(idA), n, f)

	otherNS := &vaultcommon.SecretIdentifier{Owner: "0xowner", Namespace: "other", Key: "a"}
	_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(otherNS, "ctA", 0)))
	require.NoError(t, err)

	var resp *commoncap.CapabilityResponse
	for node := range byte(3) {
		var err2 error
		resp, err2 = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node)))
		require.NoError(t, err2)
	}

	out := decode(t, resp)
	require.Len(t, out.Responses, 1)
	assert.True(t, proto.Equal(idA, out.Responses[0].GetId()))
	assert.Equal(t, "ctA", out.Responses[0].GetData().GetEncryptedValue())
}

func TestAggregator_PeerErrors(t *testing.T) {
	t.Parallel()
	t.Run("identical request errors are returned", func(t *testing.T) {
		t.Parallel()
		a := newAggregator(newReq(idA), n, f)
		msg := caperrors.NewPublicUserError(assert.AnError, caperrors.InvalidArgument).SerializeToString()

		_, err := a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, msg)
		require.NoError(t, err)
		_, err = a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, msg)
		require.Error(t, err, "success is no longer possible")
		var capErr caperrors.Error
		require.ErrorAs(t, err, &capErr)
		assert.Equal(t, caperrors.InvalidArgument, capErr.Code())
	})

	t.Run("mixed failures without disagreement", func(t *testing.T) {
		t.Parallel()
		a := newAggregator(newReq(idA), n, f)

		_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 0)))
		require.NoError(t, err)
		_, err = a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "vault direct read path is not ready")
		require.NoError(t, err)
		_, err = a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "some other error")
		require.ErrorIs(t, err, ErrNotEnoughMatchingResponses)
		assert.NotErrorIs(t, err, vaultcommon.ErrSecretVersionSkew)
	})

	t.Run("malformed responses count as errors", func(t *testing.T) {
		t.Parallel()
		a := newAggregator(newReq(idA, idB), n, f)

		_, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", 0))) // missing idB
		require.NoError(t, err)
		_, err = a.OnResponse(p2ptypes.PeerID{}, commoncap.CapabilityResponse{})
		require.ErrorIs(t, err, ErrNotEnoughMatchingResponses)
	})
}

func TestAggregatorFactory(t *testing.T) {
	t.Parallel()
	factory := NewAggregatorFactory(shareSizeLimit)
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

// Padded share entries are rejected; the honest nodes still reach quorum.
func TestAggregator_RejectsPaddedShareEntries(t *testing.T) {
	t.Parallel()
	padded := func(id *vaultcommon.SecretIdentifier, ct string, entry *vaultcommon.EncryptedShares) *vaultcommon.SecretResponse {
		return &vaultcommon.SecretResponse{Id: id, Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
			EncryptedValue: ct,
			EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
				entry,
				{EncryptionKey: "k2", BinaryShares: [][]byte{{2, 9}}},
			},
		}}}
	}

	tests := map[string]*vaultcommon.EncryptedShares{
		"multiple binary shares":   {EncryptionKey: "k1", BinaryShares: [][]byte{{1, 9}, {1, 9}, {1, 9}}},
		"binary and legacy shares": {EncryptionKey: "k1", BinaryShares: [][]byte{{1, 9}}, Shares: []string{"01"}},
		"multiple legacy shares":   {EncryptionKey: "k1", Shares: []string{"01", "02", "03"}},
		"no share at all":          {EncryptionKey: "k1"},
	}
	for name, entry := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := newAggregator(newReq(idA), n, f)

			resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, padded(idA, "ctA", entry)))
			require.NoError(t, err, "the padded response counts as an error reply, not a failure")
			require.Nil(t, resp, "the padded response joins no vote bucket")

			// The honest nodes still reach quorum without the padded shares.
			var out *commoncap.CapabilityResponse
			for node := range byte(3) {
				out, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node)))
				require.NoError(t, err)
			}
			res := decode(t, out)
			assert.Equal(t, "ctA", res.Responses[0].GetData().GetEncryptedValue())
			shares := res.Responses[0].GetData().GetEncryptedDecryptionKeyShares()
			require.Len(t, shares, 2, "one entry per encryption key")
			assert.Equal(t, [][]byte{{1, 0}, {1, 1}, {1, 2}}, shares[0].BinaryShares, "only the honest nodes' shares are merged")
		})
	}
}

func TestAggregator_RejectsOversizedShares(t *testing.T) {
	t.Parallel()
	a := newAggregator(newReq(idA), n, f)

	oversized := &vaultcommon.SecretResponse{Id: idA, Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
		EncryptedValue: "ctA",
		EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
			// One share above the VaultShareSizeLimit default of 600 bytes.
			{EncryptionKey: "k1", BinaryShares: [][]byte{make([]byte, 601)}},
			{EncryptionKey: "k2", BinaryShares: [][]byte{{2, 9}}},
		},
	}}}

	resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, oversized))
	require.NoError(t, err)
	require.Nil(t, resp, "the oversized response counts as an error reply, not a failure")

	var out *commoncap.CapabilityResponse
	for node := range byte(3) {
		out, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node)))
		require.NoError(t, err)
	}
	res := decode(t, out)
	assert.Equal(t, "ctA", res.Responses[0].GetData().GetEncryptedValue())
	assert.Equal(t, [][]byte{{1, 0}, {1, 1}, {1, 2}}, res.Responses[0].GetData().GetEncryptedDecryptionKeyShares()[0].BinaryShares)
}

func TestAggregator_RejectsInvalidShareKeys(t *testing.T) {
	t.Parallel()
	tests := map[string][]*vaultcommon.EncryptedShares{
		"unexpected key": {
			{EncryptionKey: "k1", BinaryShares: [][]byte{{1, 9}}},
			{EncryptionKey: "k3", BinaryShares: [][]byte{{2, 9}}},
		},
		"duplicate key": {
			{EncryptionKey: "k1", BinaryShares: [][]byte{{1, 9}}},
			{EncryptionKey: "k1", BinaryShares: [][]byte{{2, 9}}},
		},
		"missing requested key": {
			{EncryptionKey: "k1", BinaryShares: [][]byte{{1, 9}}},
		},
		"nil entry": {
			{EncryptionKey: "k1", BinaryShares: [][]byte{{1, 9}}},
			nil,
		},
	}
	for name, entries := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := newAggregator(newReq(idA), n, f)

			bad := &vaultcommon.SecretResponse{Id: idA, Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
				EncryptedValue:               "ctA",
				EncryptedDecryptionKeyShares: entries,
			}}}

			resp, err := a.OnResponse(p2ptypes.PeerID{}, capResp(t, bad))
			require.NoError(t, err)
			require.Nil(t, resp, "the invalid response counts as an error reply, not a failure")

			var out *commoncap.CapabilityResponse
			for node := range byte(3) {
				out, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, data(idA, "ctA", node)))
				require.NoError(t, err)
			}
			assert.Equal(t, "ctA", decode(t, out).Responses[0].GetData().GetEncryptedValue())
		})
	}
}

func TestAggregator_AcceptsLegacyHexShares(t *testing.T) {
	t.Parallel()
	a := newAggregator(newReq(idA), n, f)

	legacy := func(id *vaultcommon.SecretIdentifier, ct string, node string) *vaultcommon.SecretResponse {
		return &vaultcommon.SecretResponse{Id: id, Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
			EncryptedValue: ct,
			EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
				{EncryptionKey: "k1", Shares: []string{node}},
				{EncryptionKey: "k2", Shares: []string{node}},
			},
		}}}
	}

	var resp *commoncap.CapabilityResponse
	var err error
	for _, node := range []string{"aa", "bb", "cc"} {
		resp, err = a.OnResponse(p2ptypes.PeerID{}, capResp(t, legacy(idA, "ctA", node)))
		require.NoError(t, err)
	}
	out := decode(t, resp)
	shares := out.Responses[0].GetData().GetEncryptedDecryptionKeyShares()
	require.Len(t, shares, 2)
	assert.Equal(t, []string{"aa", "bb", "cc"}, shares[0].Shares, "legacy hex shares are merged one per node")
}
