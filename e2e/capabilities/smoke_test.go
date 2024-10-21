package capabilities_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/don"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink/integration-tests/v2/components"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	BlockchainB        *blockchain.Input `toml:"blockchain_b" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	DONInput           *don.Input        `toml:"don" validate:"required"`
}

func TestDON(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)

	out, err := don.NewBasicDON(in.DONInput, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	for i, n := range out.Nodes {
		fmt.Printf("Node %d --> %s\n", i, n.Node.HostURL)
	}

	t.Run("can access mockserver", func(t *testing.T) {
		// on the host, locally
		client := resty.New()
		_, err := client.R().
			Get(fmt.Sprintf("%s/mock1", dp.BaseURLHost))
		require.NoError(t, err)
		// inside docker
		err = components.NewDockerFakeTester(fmt.Sprintf("%s/mock1", dp.BaseURLDocker))
		require.NoError(t, err)
	})
	t.Run("smoke test", func(t *testing.T) {
		c, err := clclient.NewCLCDefaultlients(out.Nodes, framework.L)
		require.NoError(t, err)
		for _, cl := range c {
			r, _, err := cl.Health()
			require.NoError(t, err)
			framework.L.Info().Any("Response", r).Msg("Response is...")
		}
	})
	t.Run("load test", func(t *testing.T) {

	})
	t.Run("chaos test", func(t *testing.T) {

	})
}
