package reorg

//revive:disable:dot-imports
import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	networkSettingsLoad = &blockchain.EVMNetwork{
		Name:      "geth",
		Simulated: false,
		ChainID:   1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   2 * time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

var _ = Describe("LogPoller chaos tests @log-poller", func() {
	var (
		err            error
		c              blockchain.EVMClient
		cd             contracts.ContractDeployer
		le             *contracts.EthereumLogEmitter
		chainlinkNodes []*client.Chainlink
		e              *environment.Environment

		// Node params
		reorgBlocks             = 20
		EVMFinalityDepth        = "50"
		EVMTrackerHistoryDepth  = "100"
		EthLogPollInterval      = "5s"
		EthLogBackfillBatchSize = "100"

		//  Generator params
		RPS       = 20
		LogsPerTx = 50
	)
	DescribeTable("LogPoller can sustain chaos and reorgs", func(
		isReorg bool,
		chaosFunc chaos.ManifestFunc,
		chaosProps *chaos.Props,
	) {
		By("Deploying the environment")
		e = environment.New(&environment.Config{
			NamespacePrefix: "logpoller",
			TTL:             24 * time.Hour,
		})
		err := e.
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddChart(blockscout.New(&blockscout.Props{
				WsURL:   "ws://geth-ethereum-geth:8546",
				HttpURL: "http://geth-ethereum-geth:8544",
			})).
			AddHelm(reorg.New(&reorg.Props{
				NetworkName: "geth",
				NetworkType: "geth-reorg",
				Values: map[string]interface{}{
					"geth": map[string]interface{}{
						"genesis": map[string]interface{}{
							"networkId": "1337",
						},
					},
				},
			})).
			Run()
		Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
		// related https://app.shortcut.com/chainlinklabs/story/38295/creating-an-evm-chain-via-cli-or-api-immediately-polling-the-nodes-and-returning-an-error
		// node must work and reconnect even if network is not working
		time.Sleep(30 * time.Second)
		err = e.AddHelm(chainlink.New(0, map[string]interface{}{
			"env": map[string]interface{}{
				"eth_url":                        "ws://geth-ethereum-geth:8546",
				"eth_http_url":                   "http://geth-ethereum-geth:8544",
				"eth_chain_id":                   "1337",
				"FEATURE_LOG_POLLER":             "true",
				"ETH_LOG_BACKFILL_BATCH_SIZE":    EthLogBackfillBatchSize,
				"ETH_LOG_POLL_INTERVAL":          EthLogPollInterval,
				"ETH_FINALITY_DEPTH":             EVMFinalityDepth,
				"ETH_HEAD_TRACKER_HISTORY_DEPTH": EVMTrackerHistoryDepth,
			},
			"chainlink": map[string]interface{}{
				"resources": map[string]interface{}{
					"limits": map[string]interface{}{
						"cpu":    "2000m",
						"memory": "2048Mi",
					},
				},
			},
			"db": map[string]interface{}{
				"stateful": true,
				"capacity": "20Gi",
				"resources": map[string]interface{}{
					"limits": map[string]interface{}{
						"cpu":    "2000m",
						"memory": "1024Mi",
					},
				},
			},
		})).Run()
		Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")

		By("Connecting to launched resources")
		c, err = blockchain.NewEVMClient(networkSettingsLoad, e)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		cd, err = contracts.NewContractDeployer(c)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
		chainlinkNodes, err = client.ConnectChainlinkNodes(e)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")

		By("Deploying contracts")
		le, err = cd.DeployLogEmitter()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying LogEmitter Contract shouldn't fail")
		c.ParallelTransactions(true)

		fromBlockNumber, err := c.LatestBlockNumber(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		gen := blockchain.NewChainLoadGenerator(&blockchain.ChainLoadGeneratorConfig{
			RPS:      RPS,
			Contract: le,
			SharedData: &contracts.EthereumLogEmitterSharedData{
				EventType:           contracts.EventTypeInt,
				EventsPerRequest:    LogsPerTx,
				ConfirmTransactions: true,
			},
		})
		err = gen.Run()
		Expect(err).ShouldNot(HaveOccurred())

		var cid string
		if isReorg {
			rc, err := NewReorgController(
				&ReorgConfig{
					FromPodLabel:            reorg.TXNodesAppLabel,
					ToPodLabel:              reorg.MinerNodesAppLabel,
					Network:                 c,
					Env:                     e,
					BlockConsensusThreshold: 3,
					Timeout:                 5 * time.Minute,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			rc.ReOrg(reorgBlocks)
			err = rc.WaitDepthReached()
			Expect(err).ShouldNot(HaveOccurred())
		} else {
			cid, err = e.Chaos.Run(chaosFunc(e.Cfg.Namespace, chaosProps))
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(60 * time.Second)
		}

		e.Chaos.Stop(cid)

		gen.Stop()
		db, err := ctfClient.ConnectDB(0, e)
		Expect(err).ShouldNot(HaveOccurred())
		v := NewDBVerifier(db, gen.GetTransactions(), LogsPerTx, fromBlockNumber)
		v.VerifyAllTransactionsStored()
	}, []TableEntry{
		Entry("must survive reorg", true, nil, nil),
		Entry("must survive pod fail",
			false,
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{"app": a.Str("chainlink-0")},
				DurationStr:    "1m",
			}),
		Entry("must survive geth fail",
			false,
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": a.Str("chainlink-0")},
				ToLabels:    &map[string]*string{"app": a.Str("geth-ethereum-geth")},
				DurationStr: "30s",
			}),
	})
	AfterEach(func() {
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, utils.ProjectRoot, chainlinkNodes, nil, c)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
