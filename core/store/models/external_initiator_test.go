package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestNewExternalInitiator(t *testing.T) {
	eia := models.NewExternalInitiatorAuthentication()
	assert.Len(t, eia.AccessKey, 32)
	assert.Len(t, eia.Secret, 64)

	eir := &models.ExternalInitiatorRequest{
		Name: "bitcoin",
		URL:  cltest.WebURL(t, "http://localhost:8888"),
	}
	ei, err := models.NewExternalInitiator(eia, eir)
	assert.NoError(t, err)
	assert.NotEqual(t, ei.HashedSecret, eia.Secret)
	assert.Equal(t, ei.AccessKey, eia.AccessKey)
}
