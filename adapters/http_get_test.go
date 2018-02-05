package adapters_test

import (
	"net/url"
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHttpGetNotAUrlError(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("NotAUrl")
	assert.Nil(t, err)

	httpGet := adapters.HttpGet{URL: models.WebURL{u}}
	input := models.RunResult{}
	result := httpGet.Perform(input, nil)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}

func TestHttpGetResponseError(t *testing.T) {
	defer gock.Off()
	u, err := url.Parse(`https://example.com/api`)
	assert.Nil(t, err)

	gock.New(u.String()).
		Get("").
		Reply(400).
		JSON(`Invalid request`)

	httpGet := adapters.HttpGet{URL: models.WebURL{u}}
	input := models.RunResult{}
	result := httpGet.Perform(input, nil)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}
