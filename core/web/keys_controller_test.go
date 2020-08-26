package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func TestKeysController_CreateSuccess(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	request := models.CreateKeyRequest{
		CurrentPassword: cltest.Password,
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/keys", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 201)

	ethMock.AllCalled()
}

func TestKeysController_InvalidPassword(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	request := models.CreateKeyRequest{
		CurrentPassword: "12345",
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/keys", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 401)

	ethMock.AllCalled()
}

func TestKeysController_JSONBindingError(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	resp, cleanup := client.Post("/v2/keys", bytes.NewBuffer([]byte(`{"current_password":12}`)))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 422)

	ethMock.AllCalled()
}
