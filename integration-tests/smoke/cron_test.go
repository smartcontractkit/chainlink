package smoke

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCronBasic(t *testing.T) {
	t.Parallel()

	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithMockServer(1).
		WithCLNodes(1).
		Build()
	require.NoError(t, err)
	err = env.MockServer.Client.SetValuePath("/variable", 5)
	require.NoError(t, err, "Setting value path in mockserver shouldn't fail")

	bta := &client.BridgeTypeAttributes{
		Name:        fmt.Sprintf("variable-%s", uuid.NewString()),
		URL:         fmt.Sprintf("%s/variable", env.MockServer.InternalEndpoint),
		RequestData: "{}",
	}
	err = env.CLNodes[0].API.MustCreateBridge(bta)
	require.NoError(t, err, "Creating bridge in chainlink node shouldn't fail")

	job, err := env.CLNodes[0].API.MustCreateJob(&client.CronJobSpec{
		Schedule:          "CRON_TZ=UTC * * * * * *",
		ObservationSource: client.ObservationSourceSpecBridge(bta),
	})
	require.NoError(t, err, "Creating Cron Job in chainlink node shouldn't fail")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Reading Job run data shouldn't fail")

		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 5), "Expected number of job runs to be greater than 5, but got %d", len(jobRuns.Data))

		for _, jr := range jobRuns.Data {
			g.Expect(jr.Attributes.Errors).Should(gomega.Equal([]interface{}{nil}), "Job run %s shouldn't have errors", jr.ID)
		}
	}, "20m", "3s").Should(gomega.Succeed())
}
