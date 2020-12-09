package web_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var endAt = time.Now().AddDate(0, 10, 0).Round(time.Second).UTC()
var endAtISO8601 = endAt.Format(time.RFC3339)

func TestServiceAgreementsController_Show(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

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
