package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWorkflowKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	err := app.GetKeyStore().Workflow().EnsureKey(ctx)
	require.NoError(t, err)

	return client, app.GetKeyStore()
}

func TestWorkflowKeysController_Index_HappyPath(t *testing.T) {
	client, keyStore := setupWorkflowKeysControllerTests(t)
	keys, err := keyStore.Workflow().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	response, cleanup := client.Get("/v2/keys/workflow")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.WorkflowKeyResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyString(), resources[0].PublicKey)
}
