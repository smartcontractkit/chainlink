package web_test

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	eth := cltest.MockEthOnStore(app.GetStore())
	eth.RegisterSubscription("logs")

	client := app.NewHTTPClient()
	base := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))

	tests := []struct {
		name     string
		input    string
		wantCode int
	}{
		{"success", base, 200},
		{"fails validation", cltest.MustJSONDel(t, base, "payment"), 422},
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
				assert.Equal(t, time.Unix(1571523439, 0).UTC(), createdSA.Encumbrance.EndAt.Time)

				var jobids []string
				for _, j := range app.JobSubscriber.Jobs() {
					jobids = append(jobids, j.ID)
				}
				assert.Contains(t, jobids, createdSA.JobSpec.ID)
				eth.EventuallyAllCalled(t)
			}
		})
	}
}

func TestServiceAgreementsController_Create_isIdempotent(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	eth := cltest.MockEthOnStore(app.GetStore())
	eth.RegisterSubscription("logs")

	client := app.NewHTTPClient()

	reader := bytes.NewBuffer(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	resp, cleanup := client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	response1 := models.ServiceAgreement{}
	assert.NoError(t, cltest.ParseJSONAPIResponse(resp, &response1))

	reader = bytes.NewBuffer(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	resp, cleanup = client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	response2 := models.ServiceAgreement{}
	assert.NoError(t, cltest.ParseJSONAPIResponse(resp, &response2))

	assert.Equal(t, response1.ID, response2.ID)
	assert.Equal(t, response1.JobSpec.ID, response2.JobSpec.ID)
	eth.EventuallyAllCalled(t)
}

func TestServiceAgreementsController_Show(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	input := cltest.MustReadFile(t, "testdata/hello_world_agreement.json")
	sa, err := cltest.ServiceAgreementFromString(string(input))
	require.NoError(t, err)
	require.NoError(t, app.Store.CreateServiceAgreement(&sa))

	resp, cleanup := client.Get("/v2/service_agreements/" + sa.ID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	normalizedInput := cltest.NormalizedJSON(input)
	saBody := cltest.JSONFromBytes(t, b).Get("data").Get("attributes")
	assert.Equal(t, normalizedInput, saBody.String())
}
