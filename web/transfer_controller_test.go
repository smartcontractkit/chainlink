package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestTransfersController_CreateSuccess(t *testing.T) {
	config, _ := cltest.NewConfig()
	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		Amount:             assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 200)
}
