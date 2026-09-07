package cre

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const (
	stallPurgeMaxObservationBytes = 4 * 1024
	stallPurgeMaxCiphertextBytes  = 8 * 1024
)

func marshalStallPurgeCreateObservationWireSize(t *testing.T, enc string, owner common.Address, requestID, secretID string) int {
	t.Helper()
	req := &vaultcommon.CreateSecretsRequest{
		RequestId: requestID,
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{{
			Id:             &vaultcommon.SecretIdentifier{Owner: owner.Hex(), Namespace: "main", Key: secretID},
			EncryptedValue: enc,
		}},
	}
	obs := &vaultcommon.Observations{
		SortNonce: make([]byte, 32),
		Observations: []*vaultcommon.Observation{{
			Id:          requestID,
			RequestType: vaultcommon.RequestType_CREATE_SECRETS,
			Request: &vaultcommon.Observation_CreateSecretsRequest{
				CreateSecretsRequest: req,
			},
			Response: &vaultcommon.Observation_CreateSecretsResponse{
				CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
					Responses: []*vaultcommon.CreateSecretResponse{{
						Id:      req.EncryptedSecrets[0].Id,
						Success: false,
					}},
				},
			},
		}},
	}
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(obs)
	require.NoError(t, err)
	return len(b)
}

func encryptVaultSecretForOwner(t *testing.T, plaintext string, pk *tdh2easy.PublicKey, owner common.Address) (string, error) {
	t.Helper()
	return vaultutils.EncryptSecretWithWorkflowOwner(plaintext, pk, owner)
}

func encryptVaultSecretOfByteSize(t *testing.T, pk *tdh2easy.PublicKey, owner common.Address, plaintextBytes int) string {
	t.Helper()
	enc, err := vaultutils.EncryptSecretWithWorkflowOwner(strings.Repeat("x", plaintextBytes), pk, owner)
	require.NoError(t, err)
	return enc
}

func pickStallPurgeWedgePlaintextSize(t *testing.T, pk *tdh2easy.PublicKey, owner common.Address) int {
	t.Helper()
	for n := 500; n <= 4000; n += 50 {
		enc, err := vaultutils.EncryptSecretWithWorkflowOwner(strings.Repeat("x", n), pk, owner)
		require.NoError(t, err)
		raw, err := hex.DecodeString(enc)
		require.NoError(t, err)
		if len(raw) >= stallPurgeMaxCiphertextBytes {
			continue
		}
		wire := marshalStallPurgeCreateObservationWireSize(t, enc, owner, "req-stall-wedge", "stalledsecret")
		if wire > stallPurgeMaxObservationBytes {
			return n
		}
	}
	t.Fatalf("no plaintext size found with raw ciphertext below %d and observation wire above %d",
		stallPurgeMaxCiphertextBytes, stallPurgeMaxObservationBytes)
	return 0
}

func TestPendingQueueStallWedgeCiphertextSizing(t *testing.T) {
	t.Parallel()

	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	n := pickStallPurgeWedgePlaintextSize(t, pk, owner)
	enc, err := vaultutils.EncryptSecretWithWorkflowOwner(strings.Repeat("x", n), pk, owner)
	require.NoError(t, err)
	raw, err := hex.DecodeString(enc)
	require.NoError(t, err)
	wire := marshalStallPurgeCreateObservationWireSize(t, enc, owner, "req-stall-wedge", "stalledsecret")
	require.Less(t, len(raw), stallPurgeMaxCiphertextBytes)
	require.Greater(t, wire, stallPurgeMaxObservationBytes)
}
