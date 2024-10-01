package smoke

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
)

func TestOnChainReadTaskCronBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Smoke"}, tc.Cron)
	if err != nil {
		t.Fatal(err)
	}

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	sethClient, err := seth_utils.GetChainClient(config, *evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	//Set up the link contract to read from
	linkContract, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Error deploying link token contract")

	job, err := env.ClCluster.Nodes[0].API.MustCreateJob(&client.CronJobSpec{
		Schedule: "CRON_TZ=UTC * * * * * *",
		ObservationSource: `
 name   [type="onchainread" contractAddress="` + linkContract.Address() + `" contractName="linkContract" methodName="Name"]
 total_supply   [type="onchainread" contractAddress="` + linkContract.Address() + `" contractName="linkContract" methodName="TotalSupply"]
 decimals   [type="onchainread" contractAddress="` + linkContract.Address() + `" contractName="linkContract" methodName="Decimals"]
 symbol   [type="onchainread" contractAddress="` + linkContract.Address() + `" contractName="linkContract" methodName="Symbol"]`,
		Relay: "evm",
		RelayConfig: `chainID = "1337"
[relayConfig.chainReader.contracts.linkContract]
contractABI = '''
[
     {
      "constant":true,
      "inputs":[],
      "name":"name",
      "outputs":[
         {
            "name":"",
            "type":"string"
         }
      ],
      "payable":false,
      "stateMutability":"view",
      "type":"function"
    },
    {
      "constant":true,
      "inputs":[],
      "name":"totalSupply",
      "outputs":[
         {
            "name":"",
            "type":"uint256"
         }
      ],
      "payable":false,
      "stateMutability":"view",
      "type":"function"
   },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [
            {
                "name": "",
                "type": "uint8"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }
]
'''
[relayConfig.chainReader.contracts.linkContract.configs]
Name = '''
{
  "chainSpecificName": "name"
}
'''
Symbol = '''
{
  "chainSpecificName": "symbol"
}
'''
TotalSupply = '''
{
  "chainSpecificName": "totalSupply"
}
'''
Decimals = '''
{
  "chainSpecificName": "decimals"
}
'''
`,
	})

	require.NoError(t, err, "Creating Cron Job in chainlink node shouldn't fail")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(job.Data.ID)
		if err != nil {
			l.Info().Err(err).Msg("error while waiting for job runs")
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Reading Job run data shouldn't fail")

		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 8), "Expected number of job runs to be greater than 8, but got %d", len(jobRuns.Data))

		for _, jr := range jobRuns.Data {
			g.Expect(jr.Attributes.Errors).Should(gomega.Equal([]interface{}{nil, nil, nil, nil}), "Job run %s shouldn't have errors.", jr.ID)

			g.Expect(len(jr.Attributes.Outputs)).Should(gomega.Equal(4), "Should have 4 outputs but has %d", len(jr.Attributes.TaskRuns))

			g.Expect(jr.Attributes.Outputs).Should(gomega.ContainElements("1000000000000000000000000000", "ChainLink Token", "18", "LINK"), "OnChainReadTask outputs should be [\"1000000000000000000000000000\", \"ChainLink Token\", \"18\", \"LINK\"] but is ", jr.Attributes.Outputs)
		}
	}, "30s", "5s").Should(gomega.Succeed())
}
