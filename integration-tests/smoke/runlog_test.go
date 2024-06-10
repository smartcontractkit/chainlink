package smoke

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestRunLogBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.RunLog)
	require.NoError(t, err, "Error getting config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	sethClient, err := env.GetSethClient(network.ChainID)
	require.NoError(t, err, "Error getting seth client")

	lt, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	oracle, err := contracts.DeployOracle(sethClient, lt.Address())
	require.NoError(t, err, "Deploying Oracle Contract shouldn't fail")
	consumer, err := contracts.DeployAPIConsumer(sethClient, lt.Address())
	require.NoError(t, err, "Deploying Consumer Contract shouldn't fail")

	err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
	require.NoError(t, err, "Transferring %d to consumer contract shouldn't fail", big.NewInt(2e18))

	err = env.MockAdapter.SetAdapterBasedIntValuePath("/variable", []string{http.MethodPost}, 5)
	require.NoError(t, err, "Setting mock adapter value path shouldn't fail")

	jobUUID := uuid.New()

	bta := client.BridgeTypeAttributes{
		Name: fmt.Sprintf("five-%s", jobUUID.String()),
		URL:  fmt.Sprintf("%s/variable", env.MockAdapter.InternalEndpoint),
	}
	err = env.ClCluster.Nodes[0].API.MustCreateBridge(&bta)
	require.NoError(t, err, "Creating bridge shouldn't fail")

	os := &client.DirectRequestTxPipelineSpec{
		BridgeTypeAttributes: bta,
		DataPath:             "data,result",
	}
	ost, err := os.String()
	require.NoError(t, err, "Building observation source spec shouldn't fail")

	_, err = env.ClCluster.Nodes[0].API.MustCreateJob(&client.DirectRequestJobSpec{
		Name:                     fmt.Sprintf("direct-request-%s", uuid.NewString()),
		MinIncomingConfirmations: "1",
		ContractAddress:          oracle.Address(),
		EVMChainID:               fmt.Sprint(sethClient.ChainID),
		ExternalJobID:            jobUUID.String(),
		ObservationSource:        ost,
	})
	require.NoError(t, err, "Creating direct_request job shouldn't fail")

	jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
	var jobID [32]byte
	copy(jobID[:], jobUUIDReplaces)
	err = consumer.CreateRequestTo(
		oracle.Address(),
		jobID,
		big.NewInt(1e18),
		fmt.Sprintf("%s/variable", env.MockAdapter.InternalEndpoint),
		"data,result",
		big.NewInt(100),
	)
	require.NoError(t, err, "Calling oracle contract shouldn't fail")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		d, err := consumer.Data(testcontext.Get(t))
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Getting data from consumer contract shouldn't fail")
		g.Expect(d).ShouldNot(gomega.BeNil(), "Expected the initial on chain data to be nil")
		l.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
		g.Expect(d.Int64()).Should(gomega.BeNumerically("==", 5), "Expected the on-chain data to be 5, but found %d", d.Int64())
	}, "2m", "1s").Should(gomega.Succeed())
}
