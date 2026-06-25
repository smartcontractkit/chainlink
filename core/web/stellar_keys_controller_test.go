package web_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestStellarKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupStellarKeysControllerTests(t)
	keys, _ := keyStore.Stellar().GetAll()

	response, cleanup := client.Get("/v2/keys/stellar")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var resources []presenters.StellarKeyResource
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestStellarKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(t.Context()))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/stellar", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Stellar().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.StellarKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	require.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)

	_, err = keyStore.Stellar().Get(resource.ID)
	require.NoError(t, err)
}

func TestStellarKeysController_Delete_NonExistentStellarKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupStellarKeysControllerTests(t)

	nonExistentStellarKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/stellar/" + nonExistentStellarKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestStellarKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	client, keyStore := setupStellarKeysControllerTests(t)

	keys, _ := keyStore.Stellar().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Stellar().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/stellar/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Error(t, utils.JustError(keyStore.Stellar().Get(key.ID())))

	keys, _ = keyStore.Stellar().GetAll()
	assert.Len(t, keys, initialLength)
}

func setupStellarKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := t.Context()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	stellarKeyStore := app.GetKeyStore().Stellar()
	require.NotNil(t, stellarKeyStore)
	require.NoError(t, stellarKeyStore.Add(ctx, cltest.DefaultStellarKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
