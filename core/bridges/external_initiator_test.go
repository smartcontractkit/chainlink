package bridges_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
)

func TestNewExternalInitiator(t *testing.T) {
	eia := auth.NewToken()
	assert.Len(t, eia.AccessKey, 32)
	assert.Len(t, eia.Secret, 64)

	url := cltest.WebURL(t, "http://localhost:8888")
	eir := &bridges.ExternalInitiatorRequest{
		Name: "bitcoin",
		URL:  &url,
	}
	ei, err := bridges.NewExternalInitiator(eia, eir)
	assert.NoError(t, err)
	assert.NotEqual(t, ei.HashedSecret, eia.Secret)
	assert.Equal(t, ei.AccessKey, eia.AccessKey)
}
