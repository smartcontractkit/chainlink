package smoke

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestCronBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithMockAdapter().
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	err = env.MockAdapter.SetAdapterBasedIntValuePath("/variable", []string{http.MethodGet, http.MethodPost}, 5)
	require.NoError(t, err, "Setting value path in mock adapter shouldn't fail")

	bta := &client.BridgeTypeAttributes{
		Name:        fmt.Sprintf("variable-%s", uuid.NewString()),
		URL:         fmt.Sprintf("%s/variable", env.MockAdapter.InternalEndpoint),
		RequestData: "{}",
	}
	err = env.ClCluster.Nodes[0].API.MustCreateBridge(bta)
	require.NoError(t, err, "Creating bridge in chainlink node shouldn't fail")

	job, err := env.ClCluster.Nodes[0].API.MustCreateJob(&client.CronJobSpec{
		Schedule:          "CRON_TZ=UTC * * * * * *",
		ObservationSource: client.ObservationSourceSpecBridge(bta),
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
			g.Expect(jr.Attributes.Errors).Should(gomega.Equal([]interface{}{nil}), "Job run %s shouldn't have errors", jr.ID)
		}
	}, "2m", "3s").Should(gomega.Succeed())
}
