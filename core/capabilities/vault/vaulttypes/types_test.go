package vaulttypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGatewaySecretsMethod(t *testing.T) {
	t.Parallel()

	for _, method := range GatewaySecretsMethods {
		assert.True(t, IsGatewaySecretsMethod(method), method)
	}
	assert.False(t, IsGatewaySecretsMethod(MethodPublicKeyGet))
	assert.False(t, IsGatewaySecretsMethod(MethodSecretsGet))
	assert.False(t, IsGatewaySecretsMethod("vault.unsupported"))
}
