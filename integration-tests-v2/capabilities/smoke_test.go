package capabilities_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA      *blockchain.Input `toml:"blockchain_a" validate:"required"`
	BlockchainB      *blockchain.Input `toml:"blockchain_b" validate:"required"`
	FakeDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	CLNodeOne        *clnode.Input     `toml:"clnode_1" validate:"required"`
	CLNodeTwo        *clnode.Input     `toml:"clnode_2" validate:"required"`
}

func TestMultiNodeMultiNetwork(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bcNodes1, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	dpout, err := fake.NewMockedDataProvider(in.FakeDataProvider)
	require.NoError(t, err)
	in.CLNodeOne.DataProviderURL = dpout.Urls[0]
	in.CLNodeTwo.DataProviderURL = dpout.Urls[0]

	net, err := clnode.NewNetworkCfgOneNetworkAllNodes(bcNodes1)
	require.NoError(t, err)

	in.CLNodeOne.Node.TestConfigOverrides = net
	in.CLNodeTwo.Node.TestConfigOverrides = net

	_, err = clnode.NewNode(in.CLNodeOne)
	require.NoError(t, err)

	_, err = clnode.NewNode(in.CLNodeTwo)
	require.NoError(t, err)

	fmt.Printf("Node %d: http://%s\n", 1, in.CLNodeOne.Out.Node.Url)
	fmt.Printf("Node %d: http://%s\n", 2, in.CLNodeTwo.Out.Node.Url)

	t.Run("test feature A1", func(t *testing.T) {
		client := resty.New()
		_, err := client.R().
			Get("http://localhost:9111/mock1")
		require.NoError(t, err)
	})
	t.Run("test feature A2", func(t *testing.T) {
		fmt.Println("Complex testing in progress...")
		fmt.Println("Complex testing in progress...")
		fmt.Println("Complex testing in progress... Done!")
	})
}
