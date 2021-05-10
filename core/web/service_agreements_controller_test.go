package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAgreementsController_Show(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	j := models.NewJob()
	require.NoError(t, app.Store.CreateJob(&j))

	input := cltest.MustReadFile(t, "../testdata/jsonspecs/hello_world_agreement.json")
	sa, err := cltest.ServiceAgreementFromString(string(input))
	require.NoError(t, err)
	sa.JobSpecID = j.ID
	require.NoError(t, app.Store.CreateServiceAgreement(&sa))

	resp, cleanup := client.Get("/v2/service_agreements/" + sa.ID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	normalizedInput := cltest.NormalizedJSON(t, input)
	parsed := presenters.ServiceAgreement{}
	cltest.ParseJSONAPIResponse(t, resp, &parsed)
	assert.Equal(t, normalizedInput, parsed.RequestBody)
}
