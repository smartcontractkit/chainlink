package web_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	cp := presenters.ConfigPrinter{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &cp))

	assert.Equal(t, orm.LogLevel{Level: 0}, cp.LogLevel)
	assert.Contains(t, cp.RootDir, "/tmp/chainlink_test/")
	assert.Equal(t, uint16(6688), cp.Port)
	assert.Equal(t, uint16(6689), cp.TLSPort)
	assert.Equal(t, "", cp.TLSHost)
	assert.Contains(t, cp.EthereumURL, "ws://127.0.0.1:")
	assert.Equal(t, big.NewInt(3), cp.ChainID)
	assert.Contains(t, cp.ClientNodeURL, "http://127.0.0.1:")
	assert.Equal(t, uint64(6), cp.MinRequiredOutgoingConfirmations)
	assert.Equal(t, uint32(1), cp.MinIncomingConfirmations)
	assert.Equal(t, uint64(3), cp.EthGasBumpThreshold)
	assert.Equal(t, uint64(300), cp.MinimumRequestExpiration)
	assert.Equal(t, big.NewInt(5000000000), cp.EthGasBumpWei)
	assert.Equal(t, big.NewInt(20000000000), cp.EthGasPriceDefault)
	assert.Equal(t, orm.NewConfig().LinkContractAddress(), cp.LinkContractAddress)
	assert.Equal(t, orm.NewConfig().BlockBackfillDepth(), cp.BlockBackfillDepth)
	assert.Equal(t, assets.NewLink(100), cp.MinimumContractPayment)
	assert.Equal(t, (*common.Address)(nil), cp.OracleContractAddress)
	assert.Equal(t, time.Millisecond*500, cp.DatabaseTimeout.Duration())
}
