package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func TestRequestValidator_BatchSizeLimit(t *testing.T) {
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(2),
		limits.NewUpperBoundLimiter[pkgconfig.Size](1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
	)

	validValue := hex.EncodeToString(make([]byte, 10))

	makeSecrets := func(n int) []*vaultcommon.EncryptedSecret {
		secrets := make([]*vaultcommon.EncryptedSecret, n)
		for i := range secrets {
			secrets[i] = &vaultcommon.EncryptedSecret{
				Id: &vaultcommon.SecretIdentifier{
					Key:       fmt.Sprintf("key%d", i),
					Namespace: "namespace",
					Owner:     "0x1111111111111111111111111111111111111111",
				},
				EncryptedValue: validValue,
			}
		}
		return secrets
	}

	tests := []struct {
		name      string
		call      func(*testing.T, *RequestValidator) error
		errSubstr string
	}{
		{
			name: "create accepts batch at the limit",
			call: func(t *testing.T, v *RequestValidator) error {
				return v.ValidateCreateSecretsRequest(t.Context(), nil, &vaultcommon.CreateSecretsRequest{
					RequestId:        "request-id",
					EncryptedSecrets: makeSecrets(2),
				}, false)
			},
		},
		{
			name: "create rejects batch above the limit",
			call: func(t *testing.T, v *RequestValidator) error {
				return v.ValidateCreateSecretsRequest(t.Context(), nil, &vaultcommon.CreateSecretsRequest{
					RequestId:        "request-id",
					EncryptedSecrets: makeSecrets(3),
				}, false)
			},
			errSubstr: "request batch size exceeds maximum of 2",
		},
		{
			name: "update accepts batch at the limit",
			call: func(t *testing.T, v *RequestValidator) error {
				return v.ValidateUpdateSecretsRequest(t.Context(), nil, &vaultcommon.UpdateSecretsRequest{
					RequestId:        "request-id",
					EncryptedSecrets: makeSecrets(2),
				}, false)
			},
		},
		{
			name: "update rejects batch above the limit",
			call: func(t *testing.T, v *RequestValidator) error {
				return v.ValidateUpdateSecretsRequest(t.Context(), nil, &vaultcommon.UpdateSecretsRequest{
					RequestId:        "request-id",
					EncryptedSecrets: makeSecrets(3),
				}, false)
			},
			errSubstr: "request batch size exceeds maximum of 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call(t, validator)
			if tt.errSubstr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errSubstr)
		})
	}
}

func TestRequestValidator_CiphertextSizeLimit(t *testing.T) {
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter[pkgconfig.Size](10*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
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
				return validator.ValidateCreateSecretsRequest(t.Context(), nil, &vaultcommon.CreateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				}, false)
			},
			value: hex.EncodeToString(make([]byte, 10)),
		},
		{
			name: "create rejects ciphertext above the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateCreateSecretsRequest(t.Context(), nil, &vaultcommon.CreateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				}, false)
			},
			value:     hex.EncodeToString(make([]byte, 11)),
			errSubstr: "ciphertext size exceeds maximum allowed size",
		},
		{
			name: "update accepts ciphertext at the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateUpdateSecretsRequest(t.Context(), nil, &vaultcommon.UpdateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				}, false)
			},
			value: hex.EncodeToString(make([]byte, 10)),
		},
		{
			name: "update rejects ciphertext above the limit",
			call: func(t *testing.T, validator *RequestValidator, value string) error {
				return validator.ValidateUpdateSecretsRequest(t.Context(), nil, &vaultcommon.UpdateSecretsRequest{
					RequestId: "request-id",
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: id, EncryptedValue: value},
					},
				}, false)
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

func TestRequestValidator_ValidateCreateSecretsRequest_UsesRequestIdentityForOrgLabels(t *testing.T) {
	pk, _ := generateTestKeys(t)
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter[pkgconfig.Size](10*pkgconfig.KByte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
	)

	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	workflowOwner := "0x0001020304050607080900010203040506070809"
	encrypted := encryptWithOrgIDLabel(t, pk, orgID)

	err := validator.ValidateCreateSecretsRequest(t.Context(), pk, &vaultcommon.CreateSecretsRequest{
		RequestId:     "request-id",
		OrgId:         orgID,
		WorkflowOwner: workflowOwner,
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "key",
					Namespace: "namespace",
					Owner:     orgID,
				},
				EncryptedValue: encrypted,
			},
		},
	}, false)

	require.NoError(t, err)
}

func TestRequestValidator_ValidateCreateSecretsRequest_FallsBackToSecretOwnerForLegacyRequests(t *testing.T) {
	pk, _ := generateTestKeys(t)
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter[pkgconfig.Size](10*pkgconfig.KByte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
	)

	workflowOwner := "0x0001020304050607080900010203040506070809"
	encrypted := encryptWithEthAddressLabel(t, pk, workflowOwner)

	err := validator.ValidateCreateSecretsRequest(t.Context(), pk, &vaultcommon.CreateSecretsRequest{
		RequestId: "request-id",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "key",
					Namespace: "namespace",
					Owner:     workflowOwner,
				},
				EncryptedValue: encrypted,
			},
		},
	}, false)

	require.NoError(t, err)
}

func TestValidateSecretIdentifier(t *testing.T) {
	const (
		keyLimit   = 10 * pkgconfig.Byte
		ownerLimit = 10 * pkgconfig.Byte
		nsLimit    = 10 * pkgconfig.Byte
	)
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(100),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(keyLimit),
		limits.NewUpperBoundLimiter(ownerLimit),
		limits.NewUpperBoundLimiter(nsLimit),
	)

	tests := []struct {
		name      string
		key       string
		owner     string
		namespace string
		errSubstr string
	}{
		{
			name:      "valid identifier",
			key:       "mykey",
			owner:     "owner1",
			namespace: "main",
		},
		{
			name:      "empty namespace is rejected",
			key:       "mykey",
			owner:     "owner1",
			namespace: "",
			errSubstr: "namespace cannot be empty",
		},
		{
			name:      "empty key",
			key:       "",
			owner:     "owner1",
			namespace: "main",
			errSubstr: "key cannot be empty",
		},
		{
			name:      "empty owner",
			key:       "mykey",
			owner:     "",
			namespace: "main",
			errSubstr: "owner cannot be empty",
		},
		{
			name:      "invalid chars in key",
			key:       "key-invalid",
			owner:     "owner1",
			namespace: "main",
			errSubstr: "must only contain alphanumeric characters",
		},
		{
			name:      "invalid chars in namespace",
			key:       "mykey",
			owner:     "owner1",
			namespace: "bad.ns",
			errSubstr: "must only contain alphanumeric characters",
		},
		{
			name:      "invalid chars in owner",
			key:       "mykey",
			owner:     "bad-owner",
			namespace: "main",
			errSubstr: "must only contain alphanumeric characters",
		},
		{
			name:      "key at limit",
			key:       "tenbytekey", // exactly 10 bytes
			owner:     "owner1",
			namespace: "main",
		},
		{
			name:      "key exceeds limit",
			key:       "tenbytekey1", // 11 bytes
			owner:     "owner1",
			namespace: "main",
			errSubstr: "key exceeds maximum length",
		},
		{
			name:      "owner at limit",
			key:       "mykey",
			owner:     "owner12345", // exactly 10 bytes
			namespace: "main",
		},
		{
			name:      "owner exceeds limit",
			key:       "mykey",
			owner:     "owner123456", // 11 bytes
			namespace: "main",
			errSubstr: "owner exceeds maximum length",
		},
		{
			name:      "namespace at limit",
			key:       "mykey",
			owner:     "owner1",
			namespace: "tenbytekey", // exactly 10 bytes
		},
		{
			name:      "namespace exceeds limit",
			key:       "mykey",
			owner:     "owner1",
			namespace: "tenbytekey1", // 11 bytes
			errSubstr: "namespace exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSecretIdentifier(t.Context(), tt.key, tt.owner, tt.namespace)
			if tt.errSubstr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errSubstr)
		})
	}
}

func TestValidateSecretIdentifier_OwnerSpecificKeyLimit(t *testing.T) {
	const (
		defaultKeyLimit = 5 * pkgconfig.Byte
		privilegedOwner = "privilegedowner"
		privilegedLimit = 20 * pkgconfig.Byte
	)

	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(100),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		&ownerOverrideLimiter{
			defaultBound: defaultKeyLimit,
			overrides:    map[string]pkgconfig.Size{privilegedOwner: privilegedLimit},
		},
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	longKey := "averylongkeyname" // 16 bytes: exceeds default (5) but within privileged (20)

	// Regular owner cannot use the long key
	err := validator.ValidateSecretIdentifier(t.Context(), longKey, "owner1", "main")
	require.Error(t, err)
	require.ErrorContains(t, err, "key exceeds maximum length")

	// Privileged owner is allowed the same long key
	err = validator.ValidateSecretIdentifier(t.Context(), longKey, privilegedOwner, "main")
	require.NoError(t, err)
}

func TestRequestValidator_IdentifierLengths(t *testing.T) {
	const (
		keyLimit   = 5 * pkgconfig.Byte
		ownerLimit = 6 * pkgconfig.Byte
		nsLimit    = 4 * pkgconfig.Byte
	)
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(keyLimit),
		limits.NewUpperBoundLimiter(ownerLimit),
		limits.NewUpperBoundLimiter(nsLimit),
	)

	validValue := hex.EncodeToString(make([]byte, 10))

	makeRequest := func(key, owner, ns string) *vaultcommon.CreateSecretsRequest {
		return &vaultcommon.CreateSecretsRequest{
			RequestId: "req-id",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
						Key:       key,
						Owner:     owner,
						Namespace: ns,
					},
					EncryptedValue: validValue,
				},
			},
		}
	}

	tests := []struct {
		name      string
		key       string
		owner     string
		namespace string
		errSubstr string
	}{
		{
			name:      "all fields at limit",
			key:       "abcde",  // 5 bytes
			owner:     "owner1", // 6 bytes
			namespace: "main",   // 4 bytes
		},
		{
			name:      "key exceeds limit",
			key:       "abcdef", // 6 bytes
			owner:     "owner1",
			namespace: "main",
			errSubstr: "key exceeds maximum length",
		},
		{
			name:      "owner exceeds limit",
			key:       "mykey",
			owner:     "owner12", // 7 bytes
			namespace: "main",
			errSubstr: "owner exceeds maximum length",
		},
		{
			name:      "namespace exceeds limit",
			key:       "mykey",
			owner:     "owner1",
			namespace: "mains", // 5 bytes
			errSubstr: "namespace exceeds maximum length",
		},
		{
			name:      "invalid chars in key",
			key:       "key-1",
			owner:     "owner1",
			namespace: "main",
			errSubstr: "must only contain alphanumeric characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateSecretsRequest(t.Context(), nil, makeRequest(tt.key, tt.owner, tt.namespace), false)
			if tt.errSubstr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errSubstr)
		})
	}
}

func TestValidateGetSecretsRequest(t *testing.T) {
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	validID := func(key, owner, ns string) *vaultcommon.SecretIdentifier {
		return &vaultcommon.SecretIdentifier{Key: key, Owner: owner, Namespace: ns}
	}

	tests := []struct {
		name      string
		requests  []*vaultcommon.SecretRequest
		errSubstr string
	}{
		{
			name: "valid single request",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("mykey", "owner1", "main")},
			},
		},
		{
			name: "valid multiple requests",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "owner1", "main")},
				{Id: validID("key2", "owner2", "main")},
			},
		},
		{
			name:      "empty request list",
			requests:  []*vaultcommon.SecretRequest{},
			errSubstr: "no GetSecret request specified in request",
		},
		{
			name: "batch size at limit is accepted",
			requests: func() []*vaultcommon.SecretRequest {
				reqs := make([]*vaultcommon.SecretRequest, 9) // MaxBatchSize-1 = 9
				for i := range reqs {
					reqs[i] = &vaultcommon.SecretRequest{Id: validID(fmt.Sprintf("key%d", i), "owner1", "main")}
				}
				return reqs
			}(),
		},
		{
			name: "batch size equals MaxBatchSize is rejected",
			requests: func() []*vaultcommon.SecretRequest {
				reqs := make([]*vaultcommon.SecretRequest, 10) // MaxBatchSize = 10
				for i := range reqs {
					reqs[i] = &vaultcommon.SecretRequest{Id: validID(fmt.Sprintf("key%d", i), "owner1", "main")}
				}
				return reqs
			}(),
			errSubstr: "request batch size exceeds maximum of",
		},
		{
			name: "nil ID at index",
			requests: []*vaultcommon.SecretRequest{
				{Id: nil},
			},
			errSubstr: "secret ID must have id set at index",
		},
		{
			name: "empty key",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("", "owner1", "main")},
			},
			errSubstr: "secret ID must have key set at index",
		},
		{
			name: "key with invalid characters (hyphen) is rejected",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key-invalid", "owner1", "main")},
			},
			errSubstr: "must only contain alphanumeric characters",
		},
		{
			name: "key with invalid characters (slash) is rejected",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key/name", "owner1", "main")},
			},
			errSubstr: "must only contain alphanumeric characters",
		},
		{
			name: "invalid chars in owner",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "bad-owner", "main")},
			},
			errSubstr: "invalid secret identifier at index 0",
		},
		{
			name: "invalid chars in namespace",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "owner1", "bad.ns")},
			},
			errSubstr: "invalid secret identifier at index 0",
		},
		{
			name: "invalid identifier at second index",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "owner1", "main")},
				{Id: validID("key-bad", "owner1", "main")},
			},
			errSubstr: "invalid secret identifier at index 1",
		},
		{
			name: "empty owner at second index",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "owner1", "main")},
				{Id: validID("key2", "", "main")},
			},
			errSubstr: "invalid secret identifier at index 1",
		},
		{
			name: "nil id at second index",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "owner1", "main")},
				{Id: nil},
			},
			errSubstr: "secret ID must have id set at index 1",
		},
		{
			name: "empty key at second index",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("key1", "owner1", "main")},
				{Id: validID("", "owner1", "main")},
			},
			errSubstr: "secret ID must have key set at index 1",
		},
		{
			name: "empty namespace is rejected",
			requests: []*vaultcommon.SecretRequest{
				{Id: validID("mykey", "owner1", "")},
			},
			errSubstr: "namespace cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateGetSecretsRequest(t.Context(), &vaultcommon.GetSecretsRequest{Requests: tt.requests})
			if tt.errSubstr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errSubstr)
		})
	}
}

func TestValidateGetSecretsRequest_OwnerLengthPerBatchItem(t *testing.T) {
	const ownerLimit = 6 * pkgconfig.Byte
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(ownerLimit),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	// First item is within the owner length limit; second item exceeds it.
	err := validator.ValidateGetSecretsRequest(t.Context(), &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{Id: &vaultcommon.SecretIdentifier{Key: "mykey", Owner: "owner1", Namespace: "main"}},  // 6 bytes
			{Id: &vaultcommon.SecretIdentifier{Key: "mykey", Owner: "owner12", Namespace: "main"}}, // 7 bytes
		},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid secret identifier at index 1")
	require.ErrorContains(t, err, "owner exceeds maximum length")
}

func TestValidateGetSecretsRequest_KeyLengthPerBatchItem(t *testing.T) {
	const keyLimit = 5 * pkgconfig.Byte
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(keyLimit),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	// First item is within the key length limit; second item exceeds it.
	err := validator.ValidateGetSecretsRequest(t.Context(), &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{Id: &vaultcommon.SecretIdentifier{Key: "abcde", Owner: "owner1", Namespace: "main"}},  // 5 bytes
			{Id: &vaultcommon.SecretIdentifier{Key: "abcdef", Owner: "owner1", Namespace: "main"}}, // 6 bytes
		},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid secret identifier at index 1")
	require.ErrorContains(t, err, "key exceeds maximum length")
}

func TestValidateGetSecretsRequest_OwnerSpecificKeyLimit(t *testing.T) {
	const (
		defaultKeyLimit = 5 * pkgconfig.Byte
		privilegedOwner = "privilegedowner"
		privilegedLimit = 20 * pkgconfig.Byte
	)

	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		&ownerOverrideLimiter{
			defaultBound: defaultKeyLimit,
			overrides:    map[string]pkgconfig.Size{privilegedOwner: privilegedLimit},
		},
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	longKey := "averylongkeyname" // 16 bytes: exceeds default (5) but within privileged (20)

	makeRequest := func(key, owner string) *vaultcommon.GetSecretsRequest {
		return &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{Id: &vaultcommon.SecretIdentifier{Key: key, Owner: owner, Namespace: "main"}},
			},
		}
	}

	// Regular owner is rejected because the key exceeds their limit
	err := validator.ValidateGetSecretsRequest(t.Context(), makeRequest(longKey, "regularowner"))
	require.Error(t, err)
	require.ErrorContains(t, err, "key exceeds maximum length")

	// Privileged owner is accepted because the key is within their limit
	err = validator.ValidateGetSecretsRequest(t.Context(), makeRequest(longKey, privilegedOwner))
	require.NoError(t, err)
}

func TestValidateGetSecretsRequest_OwnerSpecificNamespaceLimit(t *testing.T) {
	const (
		defaultNsLimit  = 5 * pkgconfig.Byte
		privilegedOwner = "privilegedowner"
		privilegedLimit = 20 * pkgconfig.Byte
	)

	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter(1024*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		&ownerOverrideLimiter{
			defaultBound: defaultNsLimit,
			overrides:    map[string]pkgconfig.Size{privilegedOwner: privilegedLimit},
		},
	)

	longNamespace := "averylongnamespace" // 18 bytes: exceeds default (5) but within privileged (20)

	makeRequest := func(ns, owner string) *vaultcommon.GetSecretsRequest {
		return &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{Id: &vaultcommon.SecretIdentifier{Key: "mykey", Owner: owner, Namespace: ns}},
			},
		}
	}

	// Regular owner is rejected because the namespace exceeds their limit
	err := validator.ValidateGetSecretsRequest(t.Context(), makeRequest(longNamespace, "regularowner"))
	require.Error(t, err)
	require.ErrorContains(t, err, "namespace exceeds maximum length")

	// Privileged owner is accepted because the namespace is within their limit
	err = validator.ValidateGetSecretsRequest(t.Context(), makeRequest(longNamespace, privilegedOwner))
	require.NoError(t, err)
}

func TestRequestValidator_OwnerSpecificCiphertextLimit(t *testing.T) {
	const (
		defaultLimit    = 10 * pkgconfig.Byte
		privilegedOwner = "privilegedowner"
		privilegedLimit = 20 * pkgconfig.Byte
	)

	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(100),
		&ownerOverrideLimiter{
			defaultBound: defaultLimit,
			overrides:    map[string]pkgconfig.Size{privilegedOwner: privilegedLimit},
		},
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter(64*pkgconfig.Byte),
	)

	// 15 raw bytes: exceeds default (10) but within privileged (20)
	largeValue := hex.EncodeToString(make([]byte, 15))

	makeRequest := func(requestID, owner string) *vaultcommon.CreateSecretsRequest {
		return &vaultcommon.CreateSecretsRequest{
			RequestId: requestID,
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
						Key:       "mykey",
						Namespace: "main",
						Owner:     owner,
					},
					EncryptedValue: largeValue,
				},
			},
		}
	}

	// Regular owner is rejected because 15 bytes exceeds their 10-byte limit
	err := validator.ValidateCreateSecretsRequest(t.Context(), nil, makeRequest("req-1", "regularowner"), false)
	require.Error(t, err)
	require.ErrorContains(t, err, "ciphertext size exceeds maximum allowed size")

	// Privileged owner is accepted because 15 bytes is within their 20-byte limit
	err = validator.ValidateCreateSecretsRequest(t.Context(), nil, makeRequest("req-2", privilegedOwner), false)
	require.NoError(t, err)
}

func TestRequestValidator_ValidateCreateSecretsRequest_SkipsLabelValidationWithBool(t *testing.T) {
	pk, _ := generateTestKeys(t)
	validator := NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		limits.NewUpperBoundLimiter[pkgconfig.Size](10*pkgconfig.KByte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
	)

	orgID := "org_2xAbCdEfGhIjKlMnOpQrStUvWxYz"
	workflowOwner := "0x0001020304050607080900010203040506070809"
	encrypted := encryptWithOrgIDLabel(t, pk, orgID)
	request := &vaultcommon.CreateSecretsRequest{
		RequestId:     "request-id",
		WorkflowOwner: workflowOwner,
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "key",
					Namespace: "namespace",
					Owner:     workflowOwner,
				},
				EncryptedValue: encrypted,
			},
		},
	}

	err := validator.ValidateCreateSecretsRequest(t.Context(), pk, request, false)
	require.ErrorContains(t, err, "doesn't have owner as the label")

	err = validator.ValidateCreateSecretsRequest(t.Context(), pk, request, true)
	require.NoError(t, err)
}
