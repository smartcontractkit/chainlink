package directrequest

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/test_api_consumer_wrapper"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/directrequest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	ctx := t.Context()
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	productCfg, err := products.LoadOutput[directrequest.Configurator](outputFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	c, auth, _, err := products.ETHClient(
		ctx,
		in.Blockchains[0].Out.Nodes[0].ExternalWSUrl,
		productCfg.Config[0].GasSettings.FeeCapMultiplier,
		productCfg.Config[0].GasSettings.TipCapMultiplier,
	)
	require.NoError(t, err)

	consumer, err := test_api_consumer_wrapper.NewTestAPIConsumer(common.HexToAddress(productCfg.Config[0].Out.Consumer), c)
	require.NoError(t, err)

	var jobIDBytes [32]byte
	copy(jobIDBytes[:], []byte(productCfg.Config[0].Out.JobID))

	tx, err := consumer.CreateRequestTo(
		auth,
		common.HexToAddress(productCfg.Config[0].Out.Oracle),
		jobIDBytes,
		big.NewInt(1e18),
		in.FakeServer.Out.BaseURLDocker+"/direct_request_response",
		"data,result",
		big.NewInt(10),
	)
	require.NoError(t, err)
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	require.NoError(t, err)
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		d, err := consumer.Data(&bind.CallOpts{})
		assert.NoError(c, err)
		assert.Equal(c, int64(200), d.Int64())
	}, 2*time.Minute, 2*time.Second)
}
