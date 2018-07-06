package web_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/config")
	cltest.AssertServerResponse(t, resp, 200)

	cwl := presenters.ConfigWhitelist{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(resp), &cwl)
	assert.NoError(t, err)

	assert.Equal(t, store.LogLevel{Level: -1}, cwl.LogLevel)
	assert.Contains(t, cwl.RootDir, "/tmp/chainlink_test/")
	assert.Equal(t, "", cwl.Port)
	assert.Equal(t, "testusername", cwl.BasicAuthUsername)
	assert.Contains(t, cwl.EthereumURL, "ws://127.0.0.1:")
	assert.Equal(t, uint64(3), cwl.ChainID)
	assert.Contains(t, cwl.ClientNodeURL, "http://127.0.0.1:")
	assert.Equal(t, uint64(6), cwl.MinOutgoingConfirmations)
	assert.Equal(t, uint64(0), cwl.MinIncomingConfirmations)
	assert.Equal(t, uint64(3), cwl.EthGasBumpThreshold)
	assert.Equal(t, big.NewInt(5000000000), cwl.EthGasBumpWei)
	assert.Equal(t, big.NewInt(20000000000), cwl.EthGasPriceDefault)
	assert.Equal(t, "", cwl.LinkContractAddress)
	assert.Equal(t, big.NewInt(0), cwl.MinimumContractPayment)
	assert.Equal(t, (*common.Address)(nil), cwl.OracleContractAddress)
	assert.Equal(t, store.Duration{Duration: time.Millisecond * 500}, cwl.DatabasePollInterval)
}
