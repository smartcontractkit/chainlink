package adapters_test

import (
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHttpGetNotAUrlError(t *testing.T) {
	t.Parallel()
	httpGet := adapters.HttpGet{Endpoint: "NotAUrl"}
	input := models.RunResult{}
	result := httpGet.Perform(input)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}

func TestHttpGetResponseError(t *testing.T) {
	defer gock.Off()
	url := `https://example.com/api`

	gock.New(url).
		Get("").
		Reply(400).
		JSON(`Invalid request`)

	httpGet := adapters.HttpGet{Endpoint: url}
	input := models.RunResult{}
	result := httpGet.Perform(input)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}
