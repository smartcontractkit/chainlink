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
		name      string
		input     string
		wantCode  int
		id        string
		signature string
	}{
		{
			"basic",
			base.String(),
			200,
			"0x8de92e4d6a2b527f1e3c39022ee0f1c177a6d874224fe54f5c4c0d5bcaa57d50",
			"0xe86c1172d713c57b8df0d9f1ae91e4524cbb4e40719798ea48d3f777d2c893c000214b70376ddb661a8dc02fcadb911582f0e00e62eece65d3f9b8d0bfb6702201",
		},
		{"fails validation", base.Delete("payment").String(), 400, "", ""},
		{"invalid JSON", "{", 400, "", ""},
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

				createdSA := cltest.FindServiceAgreement(app.Store, responseSA.ID)
				assert.Equal(t, test.id, createdSA.ID)
				assert.Equal(t, test.signature, createdSA.Signature)
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
	sa := cltest.ServiceAgreementFromString(string(input))
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
