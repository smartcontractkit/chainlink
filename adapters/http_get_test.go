package adapters_test

import (
	"net/url"
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

// Gives a non-URL and returns an error
func TestHttpGetNotAUrlError(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("NotAUrl")
	assert.Nil(t, err)

	httpGet := adapters.HttpGet{Endpoint: models.WebURL{u}}
	input := models.RunResult{}
	result := httpGet.Perform(input, nil)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}

// Gives a valid URL that does not respond and returns an error
func TestHttpGetResponseError(t *testing.T) {
	defer gock.Off()
	u, err := url.Parse(`https://example.com/api`)
	assert.Nil(t, err)

	gock.New(u.String()).
		Get("").
		Reply(400).
		JSON(`Invalid request`)

	httpGet := adapters.HttpGet{Endpoint: models.WebURL{u}}
	input := models.RunResult{}
	result := httpGet.Perform(input, nil)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}
