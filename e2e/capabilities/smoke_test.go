package capabilities_test

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	burn_mint_erc677 "github.com/smartcontractkit/chainlink/e2e/capabilities/components/gethwrappers"
	"github.com/smartcontractkit/chainlink/e2e/capabilities/components/onchain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type Config struct {
	FundingETH         float64           `toml:"funds_eth" default:"30.0"`
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Contracts          *onchain.Input    `toml:"contracts" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestDON(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)
	pkey := os.Getenv("PRIVATE_KEY")

	// deploy docker test environment
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	for i, n := range out.CLNodes {
		fmt.Printf("Node %d --> %s\n", i, n.Node.HostURL)
		fmt.Printf("Node P2P %d --> %s\n", i, n.Node.HostP2PURL)
	}

	// deploy product contracts
	in.Contracts.URL = bc.Nodes[0].HostWSUrl
	contracts, err := onchain.NewProductOnChainDeployment(in.Contracts)
	require.NoError(t, err)

	// connect clients
	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err)

	c, err := clclient.NewCLCDefaultlients(out.CLNodes, framework.L)
	require.NoError(t, err)
	err = onchain.FundNodes(sc, c, pkey, in.FundingETH)
	require.NoError(t, err)

	t.Run("smoke test", func(t *testing.T) {
		// interact with contracts
		i, err := burn_mint_erc677.NewBurnMintERC677(contracts.Addresses[0], sc.Client)
		require.NoError(t, err)
		balance, err := i.BalanceOf(sc.NewCallOpts(), contracts.Addresses[0])
		require.NoError(t, err)
		fmt.Println(balance)

		// create jobs using deployed contracts data, this is just an example
		r, _, err := c[0].CreateJobRaw(`
		type            = "cron"
		schemaVersion   = 1
		schedule        = "CRON_TZ=UTC */10 * * * * *" # every 10 secs
		observationSource   = """
		   // data source 2
		   fetch         [type=http method=GET url="https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD"];
		   parse       [type=jsonparse path="RAW,ETH,USD,PRICE"];
		   multiply    [type="multiply" input="$(parse)" times=100]
		   encode_tx   [type="ethabiencode"
		                abi="submit(uint256 value)"
		                data="{ \\"value\\": $(multiply) }"]
		   submit_tx   [type="ethtx" to="0x859AAa51961284C94d970B47E82b8771942F1980" data="$(encode_tx)"]
		
		   fetch -> parse -> multiply -> encode_tx -> submit_tx
		"""`)
		require.NoError(t, err)
		require.Equal(t, len(r.Errors), 0)
		fmt.Printf("error: %v", err)
		fmt.Println(r)
	})
}
