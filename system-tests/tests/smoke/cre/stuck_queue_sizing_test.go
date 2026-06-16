package cre

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const (
	stuckQueueRecoveryMaxObservationBytes = 4 * 1024
	stuckQueueRecoveryMaxCiphertextBytes  = 8 * 1024
)

func marshalStuckQueueCreateObservationWireSize(t *testing.T, enc string, owner common.Address, requestID, secretID string) int {
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

func pickStuckQueueWedgePlaintextSize(t *testing.T, pk *tdh2easy.PublicKey, owner common.Address) int {
	t.Helper()
	for n := 500; n <= 4000; n += 50 {
		enc, err := vaultutils.EncryptSecretWithWorkflowOwner(strings.Repeat("x", n), pk, owner)
		require.NoError(t, err)
		raw, err := hex.DecodeString(enc)
		require.NoError(t, err)
		if len(raw) >= stuckQueueRecoveryMaxCiphertextBytes {
			continue
		}
		wire := marshalStuckQueueCreateObservationWireSize(t, enc, owner, "req-stuck-wedge", "stucksecret")
		if wire > stuckQueueRecoveryMaxObservationBytes {
			return n
		}
	}
	t.Fatalf("no plaintext size found where %d < raw ciphertext < %d and observation wire > %d",
		0, stuckQueueRecoveryMaxCiphertextBytes, stuckQueueRecoveryMaxObservationBytes)
	return 0
}

func TestStuckQueueWedgeCiphertextSizing(t *testing.T) {
	t.Parallel()

	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	n := pickStuckQueueWedgePlaintextSize(t, pk, owner)
	enc, err := vaultutils.EncryptSecretWithWorkflowOwner(strings.Repeat("x", n), pk, owner)
	require.NoError(t, err)
	raw, err := hex.DecodeString(enc)
	require.NoError(t, err)
	wire := marshalStuckQueueCreateObservationWireSize(t, enc, owner, "req-stuck-wedge", "stucksecret")
	require.Less(t, len(raw), stuckQueueRecoveryMaxCiphertextBytes)
	require.Greater(t, wire, stuckQueueRecoveryMaxObservationBytes)
	t.Logf("stuck queue wedge plaintext=%d rawCipher=%d observationWire=%d", n, len(raw), wire)
}
