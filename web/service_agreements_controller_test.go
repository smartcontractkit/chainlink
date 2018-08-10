package web_test

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	sa := cltest.FixtureCreateServiceAgreementViaWeb(t, app, "../internal/fixtures/web/hello_world_agreement.json")
	assert.NotEqual(t, "", sa.ID)
	js := cltest.FindJob(app.Store, sa.JobSpecID)
	assert.Equal(t, "0x85820c5ec619a1f517ee6cfeff545ec0ca1a90206e1a38c47f016d4137e801dd", js.Digest)

	assert.Equal(t, big.NewInt(100), sa.Encumbrance.Payment)
	assert.Equal(t, big.NewInt(2), sa.Encumbrance.Expiration)
	assert.Equal(t, "0x1d121fcb2f850b0a49018f12b26c110d87ce99fc7835608c237a88d944558fe0", sa.ID)
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
	normalizedInput := cltest.NormalizedJSONString(input)
	saBody := cltest.JSONFromString(string(b)).Get("data").Get("attributes")
	assert.Equal(t, normalizedInput, saBody.String())
}
