package v2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultMock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault/mock"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestParseVaultPublicKeyHex(t *testing.T) {
	t.Parallel()

	_, pk, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	full, err := pk.Marshal()
	require.NoError(t, err)

	t.Run("parses a full key", func(t *testing.T) {
		t.Parallel()
		got, perr := parseVaultPublicKeyHex(hex.EncodeToString(full))
		require.NoError(t, perr)
		require.NotNil(t, got)
	})

	t.Run("parses an encrypt-only key (no HArray)", func(t *testing.T) {
		t.Parallel()
		var m map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(full, &m))
		delete(m, "HArray")
		stripped, merr := json.Marshal(m)
		require.NoError(t, merr)
		got, perr := parseVaultPublicKeyHex(hex.EncodeToString(stripped))
		require.NoError(t, perr)
		require.NotNil(t, got)
	})

	t.Run("errors on invalid hex", func(t *testing.T) {
		t.Parallel()
		_, perr := parseVaultPublicKeyHex("nothex")
		require.Error(t, perr)
	})
}

// TestSecretsFetcher_PrefersRawVaultPublicKeyFromResponse proves that when the
// GetSecrets response carries RawVaultPublicKey, decryption uses it rather than
// the (deliberately different/wrong) key from the capability config.
func TestSecretsFetcher_PrefersRawVaultPublicKeyFromResponse(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	reg := registry.NewRegistry(lggr)
	peer := RandomUTF8BytesWord()
	workflowEncryptionKey := workflowkey.MustNewXXXTestingOnly(big.NewInt(1))
	workflowKeyBytes := workflowEncryptionKey.PublicKey()

	// The config key placed in the registry is a DIFFERENT key than the one that
	// actually produced the shares — decryption only succeeds if the code uses the
	// response's RawVaultPublicKey.
	_, configKey, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	configKeyBytes, err := configKey.Marshal()
	require.NoError(t, err)
	reg.SetRegistryMetadata(CreateLocalRegistryWith1Node(t, peer, workflowKeyBytes, configKeyBytes))

	rawSecret := "Raw Secret Value"
	f, n := 2, 3
	_, responseKey, privateShares, err := tdh2easy.GenerateKeys(f, n)
	require.NoError(t, err)
	responseKeyBytes, err := responseKey.Marshal()
	require.NoError(t, err)

	cipher, err := tdh2easy.Encrypt(responseKey, []byte(rawSecret))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)

	share0, err := tdh2easy.Decrypt(cipher, privateShares[0])
	require.NoError(t, err)
	share0Bytes, err := share0.Marshal()
	require.NoError(t, err)
	share1, err := tdh2easy.Decrypt(cipher, privateShares[1])
	require.NoError(t, err)
	share1Bytes, err := share1.Marshal()
	require.NoError(t, err)

	encShare0, err := workflowEncryptionKey.Encrypt(share0Bytes)
	require.NoError(t, err)
	encShare1, err := workflowEncryptionKey.Encrypt(share1Bytes)
	require.NoError(t, err)

	owner := "1234567890abcdef1234567890abcdef12345678"
	normalizedOwner, err := normalizeOwner(owner)
	require.NoError(t, err)

	mc := vaultMock.Vault{
		Fn: func(ctx context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
			return &vault.GetSecretsResponse{
				RawVaultPublicKey: hex.EncodeToString(responseKeyBytes),
				Responses: []*vault.SecretResponse{
					{
						Id: &vault.SecretIdentifier{Key: "R1", Namespace: "Bar", Owner: normalizedOwner},
						Result: &vault.SecretResponse_Data{
							Data: &vault.SecretData{
								EncryptedValue: hex.EncodeToString(cipherBytes),
								EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
									{
										BinaryShares:  [][]byte{encShare0, encShare1},
										EncryptionKey: hex.EncodeToString(workflowKeyBytes[:]),
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}
	require.NoError(t, reg.Add(t.Context(), mc))

	sf := NewSecretsFetcher(
		MetricsLabelerTest(t),
		reg,
		lggr,
		limits.WorkflowResourcePoolLimiter[int](5),
		limits.NewUpperBoundLimiter[int](5),
		nil,
		"",
		owner,
		"workflowName",
		"workflowID",
		"workflowExecID",
		time.Time{},
		workflowEncryptionKey,
		nil,
	)

	resp, err := sf.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
		Requests: []*sdkpb.SecretRequest{{Id: "R1", Namespace: "Bar"}},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Nil(t, resp[0].GetError(), "decryption should succeed using the response key, not the config key")
	assert.Equal(t, rawSecret, resp[0].GetSecret().Value)

	// Guard: the config key really is different, so success proves the response key was used.
	assert.NotEqual(t, hex.EncodeToString(configKeyBytes), hex.EncodeToString(responseKeyBytes))
}
