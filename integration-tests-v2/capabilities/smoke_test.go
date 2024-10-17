package capabilities_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/don"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA      *blockchain.Input `toml:"blockchain_a" validate:"required"`
	BlockchainB      *blockchain.Input `toml:"blockchain_b" validate:"required"`
	FakeDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	DONInput         *don.Input        `toml:"don" validate:"required"`
}

func TestMultiNodeMultiNetwork(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bcNodes1, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	dpout, err := fake.NewMockedDataProvider(in.FakeDataProvider)
	require.NoError(t, err)

	out, err := don.NewBasicDON(in.DONInput, bcNodes1, dpout.Urls[0])
	require.NoError(t, err)

	for i, n := range out.Nodes {
		fmt.Printf("Node %d --> http://%s\n", i, n.Node.Url)
	}

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
