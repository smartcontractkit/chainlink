package smoke

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	graphqlClient "github.com/smartcontractkit/chainlink/integration-tests/web/sdk/client"
)

func TestRegisteringMultipleJobDistributor(t *testing.T) {
	t.Parallel()

	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Smoke"}, "job_distributor")
	require.NoError(t, err, "Error getting config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&config).
		WithTestInstance(t).
		WithStandardCleanup().
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	ctx := context.Background()
	_, err = env.ClCluster.Nodes[0].GraphqlAPI.CreateJobDistributor(ctx, graphqlClient.JobDistributorInput{
		Name:      "job-distributor-1",
		Uri:       "http://job-distributor-1:8080",
		PublicKey: "54227538d9352e0a24550a80ab6a7af6e4f1ffbb8a604e913cbb81c484a7f97d",
	})
	require.NoError(t, err, "Creating first job distributor in chainlink node shouldn't fail")

	_, err = env.ClCluster.Nodes[0].GraphqlAPI.CreateJobDistributor(ctx, graphqlClient.JobDistributorInput{
		Name:      "job-distributor-2",
		Uri:       "http://job-distributor-2:8080",
		PublicKey: "37346b7ea98af21e1309847e00f772826ac3689fe990b1920d01efc58ad2f250",
	})
	require.NoError(t, err, "Creating second job distributor in chainlink node shouldn't fail")

	distributors, err := env.ClCluster.Nodes[0].GraphqlAPI.ListJobDistributors(ctx)
	require.NoError(t, err, "Listing job distributors in chainlink node shouldn't fail")
	require.Len(t, distributors.FeedsManagers.Results, 2, "There should be 2 job distributors")

	assert.Equal(t, "job-distributor-1", distributors.FeedsManagers.Results[0].Name)
	assert.Equal(t, "job-distributor-2", distributors.FeedsManagers.Results[1].Name)
}
