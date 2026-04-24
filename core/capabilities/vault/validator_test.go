package vault

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

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
