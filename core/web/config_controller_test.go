package web_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"chainlink/core/internal/cltest"
	"chainlink/core/store/assets"
	"chainlink/core/store/orm"
	"chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	cwl := presenters.ConfigWhitelist{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &cwl))

	assert.Equal(t, orm.LogLevel{Level: -1}, cwl.LogLevel)
	assert.Contains(t, cwl.RootDir, "/tmp/chainlink_test/")
	assert.Equal(t, uint16(6688), cwl.Port)
	assert.Equal(t, uint16(6689), cwl.TLSPort)
	assert.Equal(t, "", cwl.TLSHost)
	assert.Contains(t, cwl.EthereumURL, "ws://127.0.0.1:")
	assert.Equal(t, big.NewInt(3), cwl.ChainID)
	assert.Contains(t, cwl.ClientNodeURL, "http://127.0.0.1:")
	assert.Equal(t, uint64(6), cwl.MinOutgoingConfirmations)
	assert.Equal(t, uint32(1), cwl.MinIncomingConfirmations)
	assert.Equal(t, uint64(3), cwl.EthGasBumpThreshold)
	assert.Equal(t, uint64(300), cwl.MinimumRequestExpiration)
	assert.Equal(t, big.NewInt(5000000000), cwl.EthGasBumpWei)
	assert.Equal(t, big.NewInt(20000000000), cwl.EthGasPriceDefault)
	assert.Equal(t, orm.NewConfig().LinkContractAddress(), cwl.LinkContractAddress)
	assert.Equal(t, assets.NewLink(100), cwl.MinimumContractPayment)
	assert.Equal(t, (*common.Address)(nil), cwl.OracleContractAddress)
	assert.Equal(t, time.Millisecond*500, cwl.DatabaseTimeout)
}
