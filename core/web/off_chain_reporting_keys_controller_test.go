package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffChainReportingKeysController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	ocrKeys := []ocrkey.EncryptedKeyBundle{}

	keys, _ := app.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()

	resp, cleanup := client.Get("/v2/off_chain_reporting_keys")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &ocrKeys)
	assert.NoError(t, err)

	require.Len(t, ocrKeys, len(keys))

	assert.Equal(t, keys[0].ID, ocrKeys[0].ID)
	assert.Equal(t, keys[0].OnChainSigningAddress, ocrKeys[0].OnChainSigningAddress)
	assert.Equal(t, keys[0].OffChainPublicKey, ocrKeys[0].OffChainPublicKey)
	assert.Equal(t, keys[0].ConfigPublicKey, ocrKeys[0].ConfigPublicKey)
}
