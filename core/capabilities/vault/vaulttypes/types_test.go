package vaulttypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeNamespace(t *testing.T) {
	assert.Equal(t, DefaultNamespace, NormalizeNamespace(""))
	assert.Equal(t, "custom", NormalizeNamespace("custom"))
}

func TestIsUserSecretsMethod(t *testing.T) {
	for _, method := range UserSecretsMethods {
		assert.True(t, IsUserSecretsMethod(method), method)
	}
	assert.False(t, IsUserSecretsMethod(MethodPublicKeyGet))
	assert.False(t, IsUserSecretsMethod(MethodSecretsGet))
	assert.False(t, IsUserSecretsMethod("vault.unsupported"))
}
