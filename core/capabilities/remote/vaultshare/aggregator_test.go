package vaultshare

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func makeVaultCapabilityResponse(t *testing.T, owner, encKey string, share []byte) []byte {
	t.Helper()
	vaultResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "secret"},
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: "abc123",
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
						EncryptionKey: encKey,
						BinaryShares:  [][]byte{share},
					}},
				},
			},
		}},
	}
	anyPayload, err := anypb.New(vaultResp)
	require.NoError(t, err)
	payload, err := pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Payload: anyPayload})
	require.NoError(t, err)
	return payload
}

func peer(id byte) p2ptypes.PeerID {
	var p p2ptypes.PeerID
	p[0] = id
	return p
}

func TestAggregator_MergesFPlusOneShares(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(3, 5)
	encKey := "enc-key"
	for i := byte(1); i <= 3; i++ {
		final, err := agg.OnResponse(peer(i), &types.MessageBody{
			Error:   types.Error_OK,
			Payload: makeVaultCapabilityResponse(t, "owner", encKey, []byte{i}),
		})
		require.NoError(t, err)
		if i < 3 {
			require.Nil(t, final)
			continue
		}
		require.NotNil(t, final)
		merged := &vaultcommon.GetSecretsResponse{}
		require.NoError(t, final.Payload.UnmarshalTo(merged))
		require.Len(t, merged.Responses, 1)
		shares := merged.Responses[0].GetData().EncryptedDecryptionKeyShares
		require.Len(t, shares, 1)
		require.Len(t, shares[0].BinaryShares, 3)
	}
}

func makeVaultCapabilityResponseWithValue(t *testing.T, owner, encKey, encryptedValue string, share []byte) []byte {
	t.Helper()
	vaultResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "secret"},
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: encryptedValue,
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
						EncryptionKey: encKey,
						BinaryShares:  [][]byte{share},
					}},
				},
			},
		}},
	}
	anyPayload, err := anypb.New(vaultResp)
	require.NoError(t, err)
	payload, err := pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Payload: anyPayload})
	require.NoError(t, err)
	return payload
}

func TestAggregator_ErrorsAfterMaxResponses(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(3, 4)
	encKey := "k"
	for i := byte(1); i <= 4; i++ {
		_, err := agg.OnResponse(peer(i), &types.MessageBody{
			Error:   types.Error_OK,
			Payload: makeVaultCapabilityResponseWithValue(t, "owner", encKey, string([]byte{'a', i}), []byte("x")),
		})
		if i < 4 {
			require.NoError(t, err)
			continue
		}
		require.ErrorIs(t, err, ErrFastPathReconstruction)
	}
}

func TestAggregator_ErrorQuorum(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 4)
	for i := byte(1); i <= 2; i++ {
		final, err := agg.OnResponse(peer(i), &types.MessageBody{
			Error:    types.Error_INTERNAL_ERROR,
			ErrorMsg: "key does not exist",
		})
		if i < 2 {
			require.NoError(t, err)
			require.Nil(t, final)
			continue
		}
		require.Error(t, err)
		require.Nil(t, final)
	}
}
