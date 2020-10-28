package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeysController_Index_HappyPath(t *testing.T) {
	client, OCRKeyStore, cleanup := setupP2PKeysControllerTests(t)
	defer cleanup()

	p2pKeys := []p2pkey.EncryptedP2PKey{}

	keys, _ := OCRKeyStore.FindEncryptedP2PKeys()

	response, cleanup := client.Get("/v2/p2p_keys")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &p2pKeys)
	assert.NoError(t, err)

	require.Len(t, p2pKeys, len(keys))

	assert.Equal(t, keys[0].ID, p2pKeys[0].ID)
	assert.Equal(t, keys[0].PubKey, p2pKeys[0].PubKey)
	assert.Equal(t, keys[0].PeerID, p2pKeys[0].PeerID)
}

func setupP2PKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, *offchainreporting.KeyStore, func()) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	OCRKeyStore := app.GetStore().OCRKeyStore
	return client, OCRKeyStore, cleanup
}
