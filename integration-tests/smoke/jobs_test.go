package smoke

import (
	"fmt"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
)

type jobTest struct {
	jobType string
	jobSpec client.JobSpec
}

var memoTask = `memo [type="memo" value="21"]`

var jobsToTest = []jobTest{
	{
		jobType: job.Cron,
		jobSpec: &client.CronJobSpec{
			Schedule:          "CRON_TZ=UTC * * * * * *",
			ObservationSource: memoTask,
		},
	},
	{
		jobType: job.DirectRequest,
		jobSpec: &client.DirectRequestJobSpec{
			Name:                     "-",
			ContractAddress:          "",
			ExternalJobID:            "",
			MinIncomingConfirmations: "",
			ObservationSource:        memoTask,
		},
	},
	{
		jobType: job.FluxMonitor,
		jobSpec: &client.FluxMonitorJobSpec{
			Name:              "",
			ContractAddress:   "",
			Precision:         0,
			Threshold:         0,
			AbsoluteThreshold: 0,
			IdleTimerPeriod:   0,
			IdleTimerDisabled: false,
			PollTimerPeriod:   0,
			PollTimerDisabled: false,
			MaxTaskDuration:   0,
			ObservationSource: "",
		},
	},
	{
		jobType: job.OffchainReporting,
		jobSpec: &client.OCRTaskJobSpec{
			Name:                     "",
			BlockChainTimeout:        0,
			ContractConfirmations:    0,
			TrackerPollInterval:      0,
			TrackerSubscribeInterval: 0,
			ForwardingAllowed:        false,
			ContractAddress:          "",
			P2PBootstrapPeers:        nil,
			IsBootstrapPeer:          false,
			P2PPeerID:                "",
			KeyBundleID:              "",
			MonitoringEndpoint:       "",
			TransmitterAddress:       "",
			ObservationSource:        "",
		},
	},
	{
		jobType: job.OffchainReporting2,
		jobSpec: &client.OCR2TaskJobSpec{
			Name:              "",
			JobType:           "",
			MaxTaskDuration:   "",
			ForwardingAllowed: false,
			OCR2OracleSpec:    job.OCR2OracleSpec{},
			ObservationSource: "",
		},
	},
	{
		jobType: job.Keeper,
		jobSpec: &client.KeeperJobSpec{
			Name:                     "",
			ContractAddress:          "",
			FromAddress:              "",
			MinIncomingConfirmations: 0,
		},
	},
	{
		jobType: job.VRF,
		jobSpec: &client.VRFJobSpec{
			Name:                     "",
			CoordinatorAddress:       "",
			PublicKey:                "",
			ExternalJobID:            "",
			ObservationSource:        "",
			MinIncomingConfirmations: 0,
		},
	},
	{
		jobType: job.Webhook,
		jobSpec: &client.WebhookJobSpec{
			Name:              "",
			Initiator:         "",
			InitiatorSpec:     "",
			ObservationSource: "",
		},
	},
	{
		jobType: job.BlockhashStore,
		jobSpec: &client.BlockhashStoreJobSpec{
			Name:                  "",
			CoordinatorV2Address:  "",
			WaitBlocks:            0,
			LookbackBlocks:        0,
			BlockhashStoreAddress: "",
			PollPeriod:            "",
			RunTimeout:            "",
			EVMChainID:            "",
		},
	},
	//{
	//	jobType: job.BlockHeaderFeeder,
	//	jobSpec: nil,
	//},
	{
		jobType: job.Bootstrap,
		jobSpec: &client.OCRBootstrapJobSpec{
			Name:                     "",
			BlockChainTimeout:        0,
			ContractConfirmations:    0,
			TrackerPollInterval:      0,
			TrackerSubscribeInterval: 0,
			ContractAddress:          "",
			IsBootstrapPeer:          false,
			P2PPeerID:                "",
		},
	},
	//{
	//	jobType: job.Gateway,
	//	jobSpec: nil,
	//},
}

func TestJobSwitching(t *testing.T) {
	t.Parallel()
	testEnvironment := setupJobsTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	mockServer, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Creating mockserver client shouldn't fail")
	chainlinkNode := chainlinkNodes[0]
	err = mockServer.SetValuePath("/variable", 5)
	require.NoError(t, err, "Setting value path in mockserver shouldn't fail")
	// Register cleanup for any test
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, nil)
		require.NoError(t, err, "Error tearing down environment")
	})

	gom := gomega.NewGomegaWithT(t)
	for _, jb := range jobsToTest {
		createAndWaitForJobToRun(jb, t, gom, chainlinkNode)
		createAndWaitForJobToRun(jb, t, gom, chainlinkNode)
	}

}

func createAndWaitForJobToRun(jb jobTest, t *testing.T, gom *gomega.WithT, chainlinkNode *client.Chainlink) {
	job, err := chainlinkNode.MustCreateJob(jb.jobSpec)
	require.NoError(t, err, "Cannot create job for:", jb.jobType)

	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := chainlinkNode.MustReadRunsByJob(job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Reading Job run data shouldn't fail")

		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 3), "Expected number of job runs to be greater than 3, but got %d", len(jobRuns.Data))
	}, "2m", "1s").Should(gomega.Succeed())
	chainlinkNode.MustDeleteJob(job.Data.ID)
}

func setupJobsTest(t *testing.T) (testEnvironment *environment.Environment) {
	network := networks.SimulatedEVM
	evmConfig := ethereum.New(nil)
	if !network.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	cd, err := chainlink.NewDeployment(1, map[string]interface{}{
		"toml": client.AddNetworksConfig("", network),
	})
	require.NoError(t, err, "Error creating chainlink deployment")
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-jobs-test-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelmCharts(cd)
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment
}
