package web_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfigWithPrivateKey()
	app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
	defer cleanup()

	client := app.NewHTTPClient()
	base := cltest.EasyJSONFromFixture("../internal/fixtures/web/hello_world_agreement.json")

	tests := []struct {
		name     string
		input    string
		wantCode int
	}{
		{"success", base.String(), 200},
		{"fails validation", base.Delete("payment").String(), 422},
		{"invalid JSON", "{", 422},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, cleanup := client.Post("/v2/service_agreements", bytes.NewBufferString(test.input))
			defer cleanup()

			cltest.AssertServerResponse(t, resp, test.wantCode)
			if test.wantCode == 200 {
				responseSA := models.ServiceAgreement{}

				err := cltest.ParseJSONAPIResponse(resp, &responseSA)
				assert.NoError(t, err)
				assert.NotEqual(t, "", responseSA.ID)
				assert.NotEqual(t, "", responseSA.Signature.String())

				createdSA := cltest.FindServiceAgreement(app.Store, responseSA.ID)
				assert.NotEqual(t, "", createdSA.ID)
				assert.NotEqual(t, "", createdSA.Signature.String())
			}
		})
	}
}

func TestServiceAgreementsController_Create_isIdempotent(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfigWithPrivateKey()
	app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
	defer cleanup()

	client := app.NewHTTPClient()

	reader := bytes.NewBufferString(cltest.EasyJSONFromFixture("../internal/fixtures/web/hello_world_agreement.json").String())
	resp, cleanup := client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	response1 := models.ServiceAgreement{}
	assert.NoError(t, cltest.ParseJSONAPIResponse(resp, &response1))

	reader = bytes.NewBufferString(cltest.EasyJSONFromFixture("../internal/fixtures/web/hello_world_agreement.json").String())
	resp, cleanup = client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	response2 := models.ServiceAgreement{}
	assert.NoError(t, cltest.ParseJSONAPIResponse(resp, &response2))

	assert.Equal(t, response1.ID, response2.ID)
	assert.Equal(t, response1.JobSpec.ID, response2.JobSpec.ID)
}

func TestServiceAgreementsController_Show(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	input := cltest.LoadJSON("../internal/fixtures/web/hello_world_agreement.json")
	sa, err := cltest.ServiceAgreementFromString(string(input))
	assert.NoError(t, err)
	assert.NoError(t, app.Store.SaveServiceAgreement(&sa))

	resp, cleanup := client.Get("/v2/service_agreements/" + sa.ID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	normalizedInput := cltest.NormalizedJSON(input)
	saBody := cltest.JSONFromString(string(b)).Get("data").Get("attributes")
	assert.Equal(t, normalizedInput, saBody.String())
}
