package vaultshare

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

func TestNewAggregatorFactory_GatedByVaultFastPathGetSecretsEnabled(t *testing.T) {
	// NOTE: not parallel because it mutates the global cresettings.DefaultGetter.

	t.Run("disabled returns nil", func(t *testing.T) {
		factory := NewAggregatorFactory(1)
		require.NotNil(t, factory, "outer factory should always be non-nil")

		// DefaultGetter is nil in unit tests, so the gate falls back to the default value which is false.
		agg := factory(t.Context(), commoncap.CapabilityRequest{})
		require.Nil(t, agg, "factory should return nil when VaultFastPathGetSecretsEnabled is disabled")
	})

	t.Run("enabled returns aggregator", func(t *testing.T) {
		getter, err := settings.NewJSONGetter([]byte(`{"global":{"VaultFastPathGetSecretsEnabled":"true"}}`))
		require.NoError(t, err)

		oldGetter := cresettings.DefaultGetter
		cresettings.DefaultGetter = getter
		t.Cleanup(func() { cresettings.DefaultGetter = oldGetter })

		factory := NewAggregatorFactory(1)
		require.NotNil(t, factory, "outer factory should always be non-nil")

		agg := factory(t.Context(), commoncap.CapabilityRequest{})
		require.NotNil(t, agg, "factory should return an aggregator when VaultFastPathGetSecretsEnabled is enabled")
	})
}

func TestAggregator_VaultLevelErrorConsensus(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 4) // threshold=2, max=4

	for i := byte(1); i <= 2; i++ {
		final, err := agg.OnResponse(peer(i), makeVaultErrorMessage(t, "owner", "key does not exist"))
		if i < 2 {
			require.NoError(t, err)
			require.Nil(t, final)
			continue
		}
		require.Error(t, err)
		require.Equal(t, "key does not exist", err.Error())
		require.Nil(t, final)
	}
}

func TestAggregator_OnTimeout_ReturnsReconstructionError(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(3, 5)
	final, err := agg.OnTimeout()
	require.ErrorIs(t, err, ErrFastPathReconstruction)
	require.Nil(t, final)
}

func TestAggregator_OnTimeout_BelowThresholdReturnsReconstructionError(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 4)
	_, err := agg.OnResponse(peer(1), makeVaultErrorMessage(t, "owner", "key does not exist"))
	require.NoError(t, err)
	final, err := agg.OnTimeout()
	require.ErrorIs(t, err, ErrFastPathReconstruction)
	require.Nil(t, final)
}

func TestAggregator_Merge_MultipleSecretsAndKeys(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 5)
	encKey1 := "enc-key-1"
	encKey2 := "enc-key-2"

	for i := byte(1); i <= 2; i++ {
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{
				{
					Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret1"},
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue: "enc1",
							EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
								EncryptionKey: encKey1,
								BinaryShares:  [][]byte{[]byte{i}},
							}},
						},
					},
				},
				{
					Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret2"},
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue: "enc2",
							EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
								EncryptionKey: encKey2,
								BinaryShares:  [][]byte{[]byte{i + 10}},
							}},
						},
					},
				},
			},
		}
		payload, err := makeVaultCapabilityResponsePayload(t, resp)
		require.NoError(t, err)
		final, err := agg.OnResponse(peer(i), &types.MessageBody{Error: types.Error_OK, Payload: payload})
		if i < 2 {
			require.NoError(t, err)
			require.Nil(t, final)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, final)
		merged := &vaultcommon.GetSecretsResponse{}
		require.NoError(t, final.Payload.UnmarshalTo(merged))
		require.Len(t, merged.Responses, 2)
		for _, sr := range merged.Responses {
			require.Len(t, sr.GetData().EncryptedDecryptionKeyShares, 1)
			require.Len(t, sr.GetData().EncryptedDecryptionKeyShares[0].BinaryShares, 2)
		}
	}
}

func TestAggregator_Merge_InconsistentEncryptedValue(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 4)
	encKey := "enc-key"

	for i, value := range []string{"abc", "xyz"} {
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{{
				Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"},
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: value,
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
							EncryptionKey: encKey,
							BinaryShares:  [][]byte{[]byte{byte(i + 1)}},
						}},
					},
				},
			}},
		}
		payload, err := makeVaultCapabilityResponsePayload(t, resp)
		require.NoError(t, err)
		_, err = agg.OnResponse(peer(byte(i+1)), &types.MessageBody{Error: types.Error_OK, Payload: payload})
		require.NoError(t, err)
	}

	// Two more peers with consistent but different values to hit maxResponses.
	for i := byte(3); i <= 4; i++ {
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{{
				Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"},
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: string([]byte{'a', i}),
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
							EncryptionKey: encKey,
							BinaryShares:  [][]byte{{i}},
						}},
					},
				},
			}},
		}
		payload, err := makeVaultCapabilityResponsePayload(t, resp)
		require.NoError(t, err)
		final, err := agg.OnResponse(peer(i), &types.MessageBody{Error: types.Error_OK, Payload: payload})
		if i < 4 {
			require.NoError(t, err)
			require.Nil(t, final)
			continue
		}
		require.ErrorIs(t, err, ErrFastPathReconstruction)
		require.Nil(t, final)
	}
}

func TestAggregator_Receives_NoMoreThanMaxResponses(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 2)
	encKey := "enc-key"

	// The first two responses should merge and return immediately because threshold=2 and max=2.
	for i := byte(1); i <= 2; i++ {
		resp := &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{{
				Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"},
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "abc",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
							EncryptionKey: encKey,
							BinaryShares:  [][]byte{{i}},
						}},
					},
				},
			}},
		}
		payload, err := makeVaultCapabilityResponsePayload(t, resp)
		require.NoError(t, err)
		final, err := agg.OnResponse(peer(i), &types.MessageBody{Error: types.Error_OK, Payload: payload})
		if i < 2 {
			require.NoError(t, err)
			require.Nil(t, final)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, final)
	}

	// A third response after the final has been sent does not panic and still returns the same
	// merged final (the aggregator recomputes it). The ClientRequest layer deduplicates the final
	// response; here we just verify graceful handling of extra responses.
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"},
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: "abc",
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
						EncryptionKey: encKey,
						BinaryShares:  [][]byte{{3}},
					}},
				},
			},
		}},
	}
	payload, err := makeVaultCapabilityResponsePayload(t, resp)
	require.NoError(t, err)
	final, err := agg.OnResponse(peer(3), &types.MessageBody{Error: types.Error_OK, Payload: payload})
	require.NoError(t, err)
	require.NotNil(t, final)
}

func TestAggregator_OnResponse_FinalPayloadMarshaledCorrectly(t *testing.T) {
	t.Parallel()

	agg := NewAggregator(2, 3)
	resp, err := agg.OnResponse(peer(1), &types.MessageBody{Error: types.Error_OK, Payload: makeVaultCapabilityResponse(t, "owner", "enc-key", []byte{1})})
	require.NoError(t, err)
	require.Nil(t, resp)

	resp, err = agg.OnResponse(peer(2), &types.MessageBody{Error: types.Error_OK, Payload: makeVaultCapabilityResponse(t, "owner", "enc-key", []byte{2})})
	require.NoError(t, err)
	require.NotNil(t, resp)

	merged := &vaultcommon.GetSecretsResponse{}
	require.NoError(t, resp.Payload.UnmarshalTo(merged))
	require.Len(t, merged.Responses, 1)
	require.Len(t, merged.Responses[0].GetData().EncryptedDecryptionKeyShares[0].BinaryShares, 2)
}

func makeVaultErrorMessage(t *testing.T, owner, errMsg string) *types.MessageBody {
	t.Helper()
	vaultResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "secret"},
			Result: &vaultcommon.SecretResponse_Error{
				Error: errMsg,
			},
		}},
	}
	anyPayload, err := anypb.New(vaultResp)
	require.NoError(t, err)
	payload, err := pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Payload: anyPayload})
	require.NoError(t, err)
	return &types.MessageBody{Error: types.Error_OK, Payload: payload}
}

func makeVaultCapabilityResponsePayload(t *testing.T, vaultResp *vaultcommon.GetSecretsResponse) ([]byte, error) {
	t.Helper()
	anyPayload, err := anypb.New(vaultResp)
	if err != nil {
		return nil, err
	}
	return pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Payload: anyPayload})
}

// Ensure types implement interfaces at compile time.
var _ request.ResponseAggregator = (*Aggregator)(nil)
