package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

func generateTestKeys(t *testing.T) (*tdh2easy.PublicKey, []*tdh2easy.PrivateShare) {
	t.Helper()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	return pk, shares
}

func encryptWithEthAddressLabel(t *testing.T, pk *tdh2easy.PublicKey, owner string) string {
	t.Helper()
	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("test-secret", pk, common.HexToAddress(owner))
	require.NoError(t, err)
	return encrypted
}

func encryptWithOrgIDLabel(t *testing.T, pk *tdh2easy.PublicKey, orgID string) string {
	t.Helper()
	encrypted, err := vaultutils.EncryptSecretWithOrgID("test-secret", pk, orgID)
	require.NoError(t, err)
	return encrypted
}

func TestWorkflowOwnerToLabel(t *testing.T) {
	t.Run("ethereum address with 0x prefix", func(t *testing.T) {
		addr := "0x0001020304050607080900010203040506070809"
		label := vaultutils.WorkflowOwnerToLabel(addr)

		var expected [32]byte
		copy(expected[12:], common.HexToAddress(addr).Bytes())
		assert.Equal(t, expected, label)
	})

	t.Run("ethereum address without 0x prefix", func(t *testing.T) {
		addr := "0001020304050607080900010203040506070809"
		label := vaultutils.WorkflowOwnerToLabel(addr)

		var expected [32]byte
		copy(expected[12:], common.HexToAddress(addr).Bytes())
		assert.Equal(t, expected, label)
	})

	t.Run("checksummed ethereum address", func(t *testing.T) {
		addr := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		label := vaultutils.WorkflowOwnerToLabel(addr)

		var expected [32]byte
		copy(expected[12:], common.HexToAddress(addr).Bytes())
		assert.Equal(t, expected, label)
	})
}

func TestOrgIDToLabel(t *testing.T) {
	t.Run("org_id produces SHA256 label", func(t *testing.T) {
		orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
		label := vaultutils.OrgIDToLabel(orgID)

		expected := sha256.Sum256([]byte(orgID))
		assert.Equal(t, expected, label)
	})

	t.Run("short string", func(t *testing.T) {
		orgID := "my-org-id"
		label := vaultutils.OrgIDToLabel(orgID)

		expected := sha256.Sum256([]byte(orgID))
		assert.Equal(t, expected, label)
	})
}

func TestEnsureRightLabelOnSecret_WorkflowOwnerOnly(t *testing.T) {
	pk, _ := generateTestKeys(t)
	owner := "0x0001020304050607080900010203040506070809"
	secret := encryptWithEthAddressLabel(t, pk, owner)

	err := EnsureRightLabelOnSecret(pk, secret, owner, "")
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_OrgIDOnly(t *testing.T) {
	pk, _ := generateTestKeys(t)
	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	secret := encryptWithOrgIDLabel(t, pk, orgID)

	err := EnsureRightLabelOnSecret(pk, secret, "", orgID)
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_DualMatchesWorkflowOwner(t *testing.T) {
	pk, _ := generateTestKeys(t)
	ethAddr := "0x0001020304050607080900010203040506070809"
	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	secret := encryptWithEthAddressLabel(t, pk, ethAddr)

	err := EnsureRightLabelOnSecret(pk, secret, ethAddr, orgID)
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_DualMatchesOrgID(t *testing.T) {
	pk, _ := generateTestKeys(t)
	ethAddr := "0x0001020304050607080900010203040506070809"
	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	secret := encryptWithOrgIDLabel(t, pk, orgID)

	err := EnsureRightLabelOnSecret(pk, secret, ethAddr, orgID)
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_NeitherMatches(t *testing.T) {
	pk, _ := generateTestKeys(t)
	ethAddr := "0x0001020304050607080900010203040506070809"
	wrongAddr := "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	wrongOrgID := "org_wrong"
	secret := encryptWithEthAddressLabel(t, pk, ethAddr)
	expectedWorkflowOwnerLabelBytes := vaultutils.WorkflowOwnerToLabel(wrongAddr)
	expectedOrgIDLabelBytes := vaultutils.OrgIDToLabel(wrongOrgID)
	expectedWorkflowOwnerLabel := hex.EncodeToString(expectedWorkflowOwnerLabelBytes[:])
	expectedOrgIDLabel := hex.EncodeToString(expectedOrgIDLabelBytes[:])

	err := EnsureRightLabelOnSecret(pk, secret, wrongAddr, wrongOrgID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match any of the provided owner labels")
	assert.Contains(t, err.Error(), "expectedLabels=["+expectedWorkflowOwnerLabel+", "+expectedOrgIDLabel+"]")
}

func TestEnsureRightLabelOnSecret_BothEmpty(t *testing.T) {
	pk, _ := generateTestKeys(t)
	ethAddr := "0x0001020304050607080900010203040506070809"
	secret := encryptWithEthAddressLabel(t, pk, ethAddr)

	err := EnsureRightLabelOnSecret(pk, secret, "", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match any of the provided owner labels")
	assert.Contains(t, err.Error(), "expectedLabels=[]")
}

func TestEnsureRightLabelOnSecret_NilPublicKey(t *testing.T) {
	pk, _ := generateTestKeys(t)
	ethAddr := "0x0001020304050607080900010203040506070809"
	secret := encryptWithEthAddressLabel(t, pk, ethAddr)

	err := EnsureRightLabelOnSecret(nil, secret, ethAddr, "")
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_InvalidHexSecret(t *testing.T) {
	pk, _ := generateTestKeys(t)

	err := EnsureRightLabelOnSecret(pk, "not-valid-hex!", "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode encrypted value")
}

func TestEnsureRightLabelOnSecret_InvalidCiphertext(t *testing.T) {
	pk, _ := generateTestKeys(t)

	err := EnsureRightLabelOnSecret(pk, hex.EncodeToString([]byte("garbage")), "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to verify encrypted value")
}

func TestEnsureRightLabelOnSecret_WrongPublicKey(t *testing.T) {
	pk, _ := generateTestKeys(t)
	wrongPK, _ := generateTestKeys(t)
	ethAddr := "0x0001020304050607080900010203040506070809"
	secret := encryptWithEthAddressLabel(t, pk, ethAddr)

	err := EnsureRightLabelOnSecret(wrongPK, secret, ethAddr, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to verify encrypted value")
}

func TestEnsureRightLabelOnSecret_BackwardCompatSingleOwner(t *testing.T) {
	pk, _ := generateTestKeys(t)
	owner := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	secret := encryptWithEthAddressLabel(t, pk, owner)

	err := EnsureRightLabelOnSecret(pk, secret, owner, "")
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_LegacySecretReadViaNewFlow(t *testing.T) {
	pk, _ := generateTestKeys(t)
	workflowOwner := "0x0001020304050607080900010203040506070809"
	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"

	secret := encryptWithEthAddressLabel(t, pk, workflowOwner)
	err := EnsureRightLabelOnSecret(pk, secret, workflowOwner, orgID)
	assert.NoError(t, err)
}

func TestEnsureRightLabelOnSecret_NewSecretReadViaNewFlow(t *testing.T) {
	pk, _ := generateTestKeys(t)
	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	workflowOwner := "0x0001020304050607080900010203040506070809"

	secret := encryptWithOrgIDLabel(t, pk, orgID)
	err := EnsureRightLabelOnSecret(pk, secret, workflowOwner, orgID)
	assert.NoError(t, err)
}

func TestRequestValidator_CiphertextSizeLimit(t *testing.T) {
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter[pkgconfig.Size](10*pkgconfig.Byte),
	)

	id := &vaultcommon.SecretIdentifier{
		Key:       "key",
		Namespace: "namespace",
		Owner:     "0x1111111111111111111111111111111111111111",
	}

	tests := []struct {
		name      string
		call      func(*testing.T, *RequestValidator, string) error
		value     string
		errSubstr string
	}{
		{
			name: "create accepts ciphertext at the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateCreateSecretsRequest(nil, &vaultcommon.CreateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				})
			},
			value: hex.EncodeToString(make([]byte, 10)),
		},
		{
			name: "create rejects ciphertext above the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateCreateSecretsRequest(nil, &vaultcommon.CreateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				})
			},
			value:     hex.EncodeToString(make([]byte, 11)),
			errSubstr: "ciphertext size exceeds maximum allowed size",
		},
		{
			name: "update accepts ciphertext at the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateUpdateSecretsRequest(nil, &vaultcommon.UpdateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				})
			},
			value: hex.EncodeToString(make([]byte, 10)),
		},
		{
			name: "update rejects ciphertext above the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateUpdateSecretsRequest(nil, &vaultcommon.UpdateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				})
			},
			value:     hex.EncodeToString(make([]byte, 11)),
			errSubstr: "ciphertext size exceeds maximum allowed size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call(t, validator, tt.value)
			if tt.errSubstr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.ErrorContains(t, err, tt.errSubstr)
		})
	}
}
