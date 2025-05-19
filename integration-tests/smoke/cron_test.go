package smoke

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestCronBasic(t *testing.T) {
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
		WithMockAdapter().
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	route := &parrot.Route{
		Method:             parrot.MethodAny,
		Path:               "/variable",
		ResponseBody:       5,
		ResponseStatusCode: http.StatusOK,
	}
	err = env.MockAdapter.SetAdapterRoute(route)
	require.NoError(t, err, "Failed to set route in mock adapter")

	bta := &nodeclient.BridgeTypeAttributes{
		Name:        "variable-" + uuid.NewString(),
		URL:         env.MockAdapter.InternalEndpoint + "/variable",
		RequestData: "{}",
	}
	err = env.ClCluster.Nodes[0].API.MustCreateBridge(bta)
	require.NoError(t, err, "Creating bridge in chainlink node shouldn't fail")

	job, err := env.ClCluster.Nodes[0].API.MustCreateJob(&nodeclient.CronJobSpec{
		Schedule:          "CRON_TZ=UTC * * * * * *",
		ObservationSource: nodeclient.ObservationSourceSpecBridge(bta),
	})
	require.NoError(t, err, "Creating Cron Job in chainlink node shouldn't fail")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(job.Data.ID)
		if err != nil {
			l.Info().Err(err).Msg("error while waiting for job runs")
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Reading Job run data shouldn't fail")

		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 5), "Expected number of job runs to be greater than 5, but got %d", len(jobRuns.Data))

		for _, jr := range jobRuns.Data {
			g.Expect(jr.Attributes.Errors).Should(gomega.Equal([]any{nil}), "Job run %s shouldn't have errors", jr.ID)
		}
	}, "2m", "3s").Should(gomega.Succeed())
}

func TestCronJobReplacement(t *testing.T) {
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
		WithMockAdapter().
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	route := &parrot.Route{
		Method:             parrot.MethodAny,
		Path:               "/variable",
		ResponseBody:       5,
		ResponseStatusCode: http.StatusOK,
	}
	err = env.MockAdapter.SetAdapterRoute(route)
	require.NoError(t, err, "Failed to set route in mock adapter")

	bta := &nodeclient.BridgeTypeAttributes{
		Name:        "variable-" + uuid.NewString(),
		URL:         env.MockAdapter.InternalEndpoint + "/variable",
		RequestData: "{}",
	}
	err = env.ClCluster.Nodes[0].API.MustCreateBridge(bta)
	require.NoError(t, err, "Creating bridge in chainlink node shouldn't fail")

	// CRON job creation and replacement
	job, err := env.ClCluster.Nodes[0].API.MustCreateJob(&nodeclient.CronJobSpec{
		Schedule:          "CRON_TZ=UTC * * * * * *",
		ObservationSource: nodeclient.ObservationSourceSpecBridge(bta),
	})
	require.NoError(t, err, "Creating Cron Job in chainlink node shouldn't fail")

	gom := gomega.NewWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(job.Data.ID)
		if err != nil {
			l.Info().Err(err).Msg("error while waiting for job runs")
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Reading Job run data shouldn't fail")

		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 5), "Expected number of job runs to be greater than 5, but got %d", len(jobRuns.Data))

		for _, jr := range jobRuns.Data {
			g.Expect(jr.Attributes.Errors).Should(gomega.Equal([]interface{}{nil}), "Job run %s shouldn't have errors", jr.ID)
		}
	}, "3m", "3s").Should(gomega.Succeed())

	err = env.ClCluster.Nodes[0].API.MustDeleteJob(job.Data.ID)
	require.NoError(t, err)

	job, err = env.ClCluster.Nodes[0].API.MustCreateJob(&nodeclient.CronJobSpec{
		Schedule:          "CRON_TZ=UTC * * * * * *",
		ObservationSource: nodeclient.ObservationSourceSpecBridge(bta),
	})
	require.NoError(t, err, "Recreating Cron Job in chainlink node shouldn't fail")

	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(job.Data.ID)
		if err != nil {
			l.Info().Err(err).Msg("error while waiting for job runs")
		}
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Reading Job run data shouldn't fail")

		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 5), "Expected number of job runs to be greater than 5, but got %d", len(jobRuns.Data))

		for _, jr := range jobRuns.Data {
			g.Expect(jr.Attributes.Errors).Should(gomega.Equal([]interface{}{nil}), "Job run %s shouldn't have errors", jr.ID)
		}
	}, "3m", "3s").Should(gomega.Succeed())
}
