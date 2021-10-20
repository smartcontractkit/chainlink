package web_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"
)

func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	cp := presenters.ConfigPrinter{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &cp))

	assert.Contains(t, cp.RootDir, "/tmp/chainlink_test/")
	assert.Equal(t, uint16(6688), cp.Port)
	assert.Equal(t, uint16(6689), cp.TLSPort)
	assert.Equal(t, "", cp.TLSHost)
	assert.Contains(t, cp.EthereumURL, "ws://localhost:8546")
	assert.Equal(t, big.NewInt(eth.NullClientChainID), cp.ChainID)
	assert.Contains(t, cp.ClientNodeURL, "http://127.0.0.1:")
	assert.Equal(t, uint64(6), cp.MinRequiredOutgoingConfirmations)
	assert.Equal(t, uint32(1), cp.MinIncomingConfirmations)
	assert.Equal(t, uint64(3), cp.EthGasBumpThreshold)
	assert.Equal(t, uint64(300), cp.MinimumRequestExpiration)
	assert.Equal(t, big.NewInt(5000000000), cp.EthGasBumpWei)
	assert.Equal(t, big.NewInt(20000000000), cp.EthGasPriceDefault)
	assert.Equal(t, config.NewConfig().BlockBackfillDepth(), cp.BlockBackfillDepth)
	assert.Equal(t, assets.NewLink(100), cp.MinimumContractPayment)
	assert.Equal(t, common.Address{}, cp.OperatorContractAddress)
	assert.Equal(t, time.Second*5, cp.DatabaseTimeout.Duration())
}
