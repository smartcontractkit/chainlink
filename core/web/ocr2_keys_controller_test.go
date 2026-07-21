package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestOCR2KeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()
	client, OCRKeyStore := setupOCR2KeysControllerTests(t)

	keys, _ := OCRKeyStore.GetAll()

	response, cleanup := client.Get("/v2/keys/ocr2")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.OCR2KeysBundleResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))
	assert.Equal(t, keys[0].ID(), resources[0].ID)
}

func TestOCR2KeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name      string
		chainType corekeys.ChainType
	}{
		{"EVM keys", "evm"},
		{"Solana Keys", "solana"},
	} {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			client, OCRKeyStore := setupOCR2KeysControllerTests(tt)

			keys, _ := OCRKeyStore.GetAll()
			initialLength := len(keys)

			response, cleanup := client.Post(fmt.Sprintf("/v2/keys/ocr2/%s", test.chainType), nil)
			tt.Cleanup(cleanup)
			cltest.AssertServerResponse(tt, response, http.StatusOK)

			keys, _ = OCRKeyStore.GetAll()
			require.Len(tt, keys, initialLength+1)

			resource := presenters.OCR2KeysBundleResource{}
			err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(tt, response), &resource)
			require.NoError(tt, err)

			ids := make([]string, 0, len(keys))
			for _, key := range keys {
				ids = append(ids, key.ID())
			}
			require.Contains(tt, ids, resource.ID)

			_, err = OCRKeyStore.Get(resource.ID)
			assert.NoError(tt, err)
		})
	}
}

func TestOCR2KeysController_Delete_NonExistentOCRKeyID(t *testing.T) {
	t.Parallel()
	client, _ := setupOCR2KeysControllerTests(t)

	nonExistentOCRKeyID := "eb81f4a35033ac8dd68b9d33a039a713d6fd639af6852b81f47ffeda1c95de54"
	response, cleanup := client.Delete("/v2/keys/ocr2/" + nonExistentOCRKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestOCR2KeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	client, OCRKeyStore := setupOCR2KeysControllerTests(t)

	keys, _ := OCRKeyStore.GetAll()
	initialLength := len(keys)
	key, _ := OCRKeyStore.Create(ctx, "evm")

	response, cleanup := client.Delete("/v2/keys/ocr2/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Error(t, utils.JustError(OCRKeyStore.Get(key.ID())))

	keys, _ = OCRKeyStore.GetAll()
	assert.Len(t, keys, initialLength)
}

func setupOCR2KeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.OCR2) {
	ctx := t.Context()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(t.Context()))
	client := app.NewHTTPClient(nil)

	_, err := app.KeyStore.OCR2().Create(ctx, "evm")
	require.NoError(t, err)

	return client, app.GetKeyStore().OCR2()
}
