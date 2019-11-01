package web_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	"chainlink/core/store/models"
	"chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var endAt = time.Now().AddDate(0, 10, 0).Round(time.Second).UTC()
var endAtISO8601 = endAt.Format(time.RFC3339)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	eth := cltest.MockEthOnStore(t, app.GetStore())
	eth.RegisterSubscription("logs")

	client := app.NewHTTPClient()
	base := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	base = strings.Replace(base, "2019-10-19T22:17:19Z", endAtISO8601, 1)
	tests := []struct {
		name     string
		input    string
		wantCode int
	}{
		{"success", base, http.StatusOK},
		{"fails validation", cltest.MustJSONDel(t, base, "payment"), http.StatusUnprocessableEntity},
		{"invalid JSON", "{", http.StatusUnprocessableEntity},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, cleanup := client.Post("/v2/service_agreements", bytes.NewBufferString(test.input))
			defer cleanup()

			cltest.AssertServerResponse(t, resp, test.wantCode)
			if test.wantCode == http.StatusOK {
				responseSA := models.ServiceAgreement{}

				err := cltest.ParseJSONAPIResponse(t, resp, &responseSA)
				require.NoError(t, err)
				assert.NotEqual(t, "", responseSA.ID)
				assert.NotEqual(t, "", responseSA.Signature.String())

				createdSA := cltest.FindServiceAgreement(t, app.Store, responseSA.ID)
				assert.NotEqual(t, "", createdSA.ID)
				assert.NotEqual(t, "", createdSA.Signature.String())
				assert.Equal(t, endAt, createdSA.Encumbrance.EndAt.Time)

				var jobids []*models.ID
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

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	eth := cltest.MockEthOnStore(t, app.GetStore())
	eth.RegisterSubscription("logs")

	client := app.NewHTTPClient()

	base := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	base = strings.Replace(base, "2019-10-19T22:17:19Z", endAtISO8601, 1)
	reader := bytes.NewBuffer([]byte(base))

	resp, cleanup := client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	response1 := models.ServiceAgreement{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &response1))

	reader = bytes.NewBuffer([]byte(base))
	resp, cleanup = client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	response2 := models.ServiceAgreement{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &response2))

	assert.Equal(t, response1.ID, response2.ID)
	assert.Equal(t, response1.JobSpec.ID, response2.JobSpec.ID)
	eth.EventuallyAllCalled(t)
}

func TestServiceAgreementsController_Show(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	client := app.NewHTTPClient()

	input := cltest.MustReadFile(t, "testdata/hello_world_agreement.json")
	sa, err := cltest.ServiceAgreementFromString(string(input))
	require.NoError(t, err)
	require.NoError(t, app.Store.CreateServiceAgreement(&sa))

	resp, cleanup := client.Get("/v2/service_agreements/" + sa.ID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	normalizedInput := cltest.NormalizedJSON(t, input)
	parsed := presenters.ServiceAgreement{}
	cltest.ParseJSONAPIResponse(t, resp, &parsed)
	assert.Equal(t, normalizedInput, parsed.RequestBody)
}
