package features

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient/gql/client"
	de "github.com/smartcontractkit/chainlink/devenv"
)

func TestMultipleJobDistributors(t *testing.T) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	node := in.NodeSets[0].Out.CLNodes[0].Node

	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	c, err := client.NewWithContext(t.Context(), node.ExternalURL, client.Credentials{
		Email:    node.APIAuthUser,
		Password: node.APIAuthPassword,
	})
	require.NoError(t, err)

	_, err = c.CreateJobDistributor(t.Context(), client.JobDistributorInput{
		Name:      "job-distributor-1",
		Uri:       "http://job-distributor-1:8080",
		PublicKey: "54227538d9352e0a24550a80ab6a7af6e4f1ffbb8a604e913cbb81c484a7f97d",
	})
	require.NoError(t, err)

	_, err = c.CreateJobDistributor(t.Context(), client.JobDistributorInput{
		Name:      "job-distributor-2",
		Uri:       "http://job-distributor-2:8080",
		PublicKey: "37346b7ea98af21e1309847e00f772826ac3689fe990b1920d01efc58ad2f250",
	})
	require.NoError(t, err)

	distributors, err := c.ListJobDistributors(t.Context())
	require.NoError(t, err)
	require.Len(t, distributors.FeedsManagers.Results, 2, "There should be 2 job distributors")

	require.Equal(t, "job-distributor-1", distributors.FeedsManagers.Results[0].Name)
	require.Equal(t, "job-distributor-2", distributors.FeedsManagers.Results[1].Name)
}
