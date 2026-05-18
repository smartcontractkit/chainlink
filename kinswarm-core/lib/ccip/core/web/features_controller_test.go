package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FeaturesController_List(t *testing.T) {
	app := cltest.NewApplicationWithConfig(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		csa := true
		c.Feature.UICSAKeys = &csa
		c.Feature.MultiFeedsManagers = &csa
	}))
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Get("/v2/features")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resources := []presenters.FeatureResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resources)
	require.NoError(t, err)
	require.Len(t, resources, 3)

	assert.Equal(t, "csa", resources[0].ID)
	assert.True(t, resources[0].Enabled)

	assert.Equal(t, "feeds_manager", resources[1].ID)
	assert.True(t, resources[1].Enabled)

	assert.Equal(t, "multi_feeds_managers", resources[2].ID)
	assert.True(t, resources[2].Enabled)
}
