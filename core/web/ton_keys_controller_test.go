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

func TestTonKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupTonKeysControllerTests(t)
	keys, _ := keyStore.Ton().GetAll()

	response, cleanup := client.Get("/v2/keys/ton")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.TonKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))

	require.Equal(t, keys[0].ID(), resources[0].ID)
	require.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestTonKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/ton", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Ton().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.TonKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	require.NoError(t, err)

	require.Equal(t, keys[0].ID(), resource.ID)
	require.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)
	require.Equal(t, keys[0].UserFriendlyAddress(), resource.UserFriendlyAddress)
	require.Equal(t, keys[0].RawAddress(), resource.RawAddress)

	_, err = keyStore.Ton().Get(resource.ID)
	require.NoError(t, err)
}

func TestTonKeysController_Delete_NonExistentTonKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupTonKeysControllerTests(t)

	nonExistentTonKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/ton/" + nonExistentTonKeyID)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestTonKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	client, keyStore := setupTonKeysControllerTests(t)

	keys, _ := keyStore.Ton().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Ton().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/ton/" + key.ID())
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Error(t, utils.JustError(keyStore.Ton().Get(key.ID())))

	keys, _ = keyStore.Ton().GetAll()
	require.Len(t, keys, initialLength)
}

func setupTonKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.KeyStore.Ton().Add(ctx, cltest.DefaultTonKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
