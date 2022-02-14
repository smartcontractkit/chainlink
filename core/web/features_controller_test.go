package web_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FeaturesController_List(t *testing.T) {
	t.Parallel()

	_, client := setupFeaturesControllerTest(t)

	resp, cleanup := client.Get("/v2/features")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resources := []presenters.FeatureResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resources)
	require.NoError(t, err)
	require.Len(t, resources, 2)

	assert.Equal(t, "csa", resources[0].ID)
	assert.True(t, resources[0].Enabled)

	assert.Equal(t, "feeds_manager", resources[1].ID)
	assert.False(t, resources[1].Enabled)
}

func setupFeaturesControllerTest(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner) {
	os.Setenv("FEATURE_UI_CSA_KEYS", "true")

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient()

	return app, client
}
