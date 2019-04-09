package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExternalInitiator(t *testing.T) {
	eia := NewExternalInitiatorAuthentication()
	assert.Len(t, eia.AccessKey, 32)
	assert.Len(t, eia.Secret, 64)

	ei, err := NewExternalInitiator(eia)
	assert.NoError(t, err)
	assert.NotEqual(t, ei.HashedSecret, eia.Secret)
	assert.Equal(t, ei.AccessKey, eia.AccessKey)
}
