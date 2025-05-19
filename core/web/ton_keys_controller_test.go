package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/stretchr/testify/require"
)

func TestTONKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupTONKeysControllerTests(t)
	keys, _ := keyStore.TON().GetAll()

	response, cleanup := client.Get("/v2/keys/ton")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.TONKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))

	require.Equal(t, keys[0].ID(), resources[0].ID)
	require.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestTONKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/ton", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.TON().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.TONKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	require.NoError(t, err)

	require.Equal(t, keys[0].ID(), resource.ID)
	require.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)
	require.Equal(t, keys[0].AddressBase64(), resource.AddressBase64)
	require.Equal(t, keys[0].RawAddress(), resource.RawAddress)

	_, err = keyStore.TON().Get(resource.ID)
	require.NoError(t, err)
}

func TestTONKeysController_Delete_NonExistentTONKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupTONKeysControllerTests(t)

	nonExistentTONKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/ton/" + nonExistentTONKeyID)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestTONKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	client, keyStore := setupTONKeysControllerTests(t)

	keys, _ := keyStore.TON().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.TON().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/ton/" + key.ID())
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Error(t, utils.JustError(keyStore.TON().Get(key.ID())))

	keys, _ = keyStore.TON().GetAll()
	require.Len(t, keys, initialLength)
}

func setupTONKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.KeyStore.TON().Add(ctx, cltest.DefaultTONKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
