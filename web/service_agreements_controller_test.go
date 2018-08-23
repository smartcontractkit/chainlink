package web_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/web"
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

				body := cltest.ParseResponseBody(resp)
				err := web.ParseJSONAPIResponse(body, &responseSA)
				assert.NoError(t, err)
				cltest.AssertValidHash(t, 32, responseSA.ID)
				cltest.AssertValidHash(t, 65, responseSA.Signature)

				createdSA := cltest.FindServiceAgreement(app.Store, responseSA.ID)
				cltest.AssertValidHash(t, 32, createdSA.ID)
				cltest.AssertValidHash(t, 65, createdSA.Signature)
			}
		})
	}
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
