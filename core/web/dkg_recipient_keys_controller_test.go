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

func setupDKGRecipientKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	err := app.GetKeyStore().DKGRecipient().EnsureKey(ctx)
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, app.Stop()) })
	return client, app.GetKeyStore()
}

func TestDKGRecipientKeysController_Index_HappyPath(t *testing.T) {
	client, keyStore := setupDKGRecipientKeysControllerTests(t)
	keys, err := keyStore.DKGRecipient().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	response, cleanup := client.Get("/v2/keys/dkgrecipient")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.DKGRecipientKeyResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyString(), resources[0].PublicKey)
}
