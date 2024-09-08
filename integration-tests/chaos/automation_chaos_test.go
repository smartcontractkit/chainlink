package chaos

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/automationv2"

	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	baseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`

	defaultAutomationSettings = map[string]interface{}{
		"replicas": 6,
		"toml":     "",
		"db": map[string]interface{}{
			"stateful": true,
			"capacity": "1Gi",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "250m",
					"memory": "256Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "250m",
					"memory": "256Mi",
				},
			},
		},
	}

	defaultEthereumSettings = ethereum.Props{
		Values: map[string]interface{}{
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "4000m",
					"memory": "4Gi",
				},
				"limits": map[string]interface{}{
					"cpu":    "4000m",
					"memory": "4Gi",
				},
			},
			"geth": map[string]interface{}{
				"blocktime": "1",
			},
		},
	}
)

func getDefaultAutomationSettings(config *tc.TestConfig) map[string]interface{} {
	defaultAutomationSettings["toml"] = networks.AddNetworksConfig(baseTOML, config.Pyroscope, networks.MustGetSelectedNetworkConfig(config.Network)[0])
	return defaultAutomationSettings
}

func getDefaultEthereumSettings(config *tc.TestConfig) *ethereum.Props {
	network := networks.MustGetSelectedNetworkConfig(config.Network)[0]
	defaultEthereumSettings.NetworkName = network.Name
	defaultEthereumSettings.Simulated = network.Simulated
	defaultEthereumSettings.WsURLs = network.URLs

	return &defaultEthereumSettings
}

type KeeperConsumerContracts int32

const (
	BasicCounter KeeperConsumerContracts = iota

	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
	numberOfUpkeeps       = 2
)

func TestAutomationChaos(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	registryVersions := map[string]eth_contracts.KeeperRegistryVersion{
		"registry_2_0": eth_contracts.RegistryVersion_2_0,
		"registry_2_1": eth_contracts.RegistryVersion_2_1,
		"registry_2_2": eth_contracts.RegistryVersion_2_2,
		"registry_2_3": eth_contracts.RegistryVersion_2_3,
	}

	for name, registryVersion := range registryVersions {
		rv := registryVersion
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config, err := tc.GetConfig([]string{"Chaos"}, tc.Automation)
			if err != nil {
				t.Fatal(err)
			}

			var overrideFn = func(_ interface{}, target interface{}) {
				ctf_config.MustConfigOverrideChainlinkVersion(config.GetChainlinkImageConfig(), target)
				ctf_config.MightConfigOverridePyroscopeKey(config.GetPyroscopeConfig(), target)
			}

			chainlinkCfg := chainlink.NewWithOverride(0, getDefaultAutomationSettings(&config), config.ChainlinkImage, overrideFn)

			testCases := map[string]struct {
				networkChart environment.ConnectedChart
				clChart      environment.ConnectedChart
				chaosFunc    chaos.ManifestFunc
				chaosProps   *chaos.Props
			}{
				// see ocr_chaos.test.go for comments
				PodChaosFailMinorityNodes: {
					ethereum.New(getDefaultEthereumSettings(&config)),
					chainlinkCfg,
					chaos.NewFailPods,
					&chaos.Props{
						LabelsSelector: &map[string]*string{ChaosGroupMinority: ptr.Ptr("1")},
						DurationStr:    "1m",
					},
				},
				PodChaosFailMajorityNodes: {
					ethereum.New(getDefaultEthereumSettings(&config)),
					chainlinkCfg,
					chaos.NewFailPods,
					&chaos.Props{
						LabelsSelector: &map[string]*string{ChaosGroupMajority: ptr.Ptr("1")},
						DurationStr:    "1m",
					},
				},
				PodChaosFailMajorityDB: {
					ethereum.New(getDefaultEthereumSettings(&config)),
					chainlinkCfg,
					chaos.NewFailPods,
					&chaos.Props{
						LabelsSelector: &map[string]*string{ChaosGroupMajority: ptr.Ptr("1")},
						DurationStr:    "1m",
						ContainerNames: &[]*string{ptr.Ptr("chainlink-db")},
					},
				},
				NetworkChaosFailMajorityNetwork: {
					ethereum.New(getDefaultEthereumSettings(&config)),
					chainlinkCfg,
					chaos.NewNetworkPartition,
					&chaos.Props{
						FromLabels:  &map[string]*string{ChaosGroupMajority: ptr.Ptr("1")},
						ToLabels:    &map[string]*string{ChaosGroupMinority: ptr.Ptr("1")},
						DurationStr: "1m",
					},
				},
				NetworkChaosFailBlockchainNode: {
					ethereum.New(getDefaultEthereumSettings(&config)),
					chainlinkCfg,
					chaos.NewNetworkPartition,
					&chaos.Props{
						FromLabels:  &map[string]*string{"app": ptr.Ptr("geth")},
						ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: ptr.Ptr("1")},
						DurationStr: "1m",
					},
				},
			}

			for name, testCase := range testCases {
				name := name
				testCase := testCase
				t.Run(fmt.Sprintf("Automation_%s", name), func(t *testing.T) {
					t.Parallel()
					network := networks.MustGetSelectedNetworkConfig(config.Network)[0] // Need a new copy of the network for each test

					testEnvironment := environment.
						New(&environment.Config{
							NamespacePrefix: fmt.Sprintf("chaos-automation-%s", name),
							TTL:             time.Hour * 1,
							Test:            t,
						}).
						AddHelm(testCase.networkChart).
						AddHelm(testCase.clChart)
					// TODO we need to update the image in CTF, the old one is not available anymore
					// deploy blockscout if running on simulated
					// AddHelm(testCase.clChart).
					// AddChart(blockscout.New(&blockscout.Props{
					// 	Name:    "geth-blockscout",
					// 	WsURL:   network.URL,
					// 	HttpURL: network.HTTPURLs[0],
					// })
					err = testEnvironment.Run()
					require.NoError(t, err, "Error setting up test environment")
					if testEnvironment.WillUseRemoteRunner() {
						return
					}

					err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, "instance=node-", 1, 2, ChaosGroupMinority)
					require.NoError(t, err)
					err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, "instance=node-", 3, 5, ChaosGroupMajority)
					require.NoError(t, err)
					err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, "instance=node-", 2, 5, ChaosGroupMajorityPlus)
					require.NoError(t, err)

					chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
					require.NoError(t, err, "Error connecting to Chainlink nodes")

					network = seth_utils.MustReplaceSimulatedNetworkUrlWithK8(l, network, *testEnvironment)

					chainClient, err := seth_utils.GetChainClientWithConfigFunction(&config, network, seth_utils.OneEphemeralKeysLiveTestnetAutoFixFn)
					require.NoError(t, err, "Error creating seth client")

					// Register cleanup for any test
					t.Cleanup(func() {
						err := actions.TeardownSuite(t, chainClient, testEnvironment, chainlinkNodes, nil, zapcore.PanicLevel, &config)
						require.NoError(t, err, "Error tearing down environment")
					})

					a := automationv2.NewAutomationTestK8s(l, chainClient, chainlinkNodes)
					a.SetMercuryCredentialName("cred1")
					conf := config.Automation.AutomationConfig
					a.RegistrySettings = contracts.KeeperRegistrySettings{
						PaymentPremiumPPB:    *conf.RegistrySettings.PaymentPremiumPPB,
						FlatFeeMicroLINK:     *conf.RegistrySettings.FlatFeeMicroLINK,
						CheckGasLimit:        *conf.RegistrySettings.CheckGasLimit,
						StalenessSeconds:     conf.RegistrySettings.StalenessSeconds,
						GasCeilingMultiplier: *conf.RegistrySettings.GasCeilingMultiplier,
						MaxPerformGas:        *conf.RegistrySettings.MaxPerformGas,
						MinUpkeepSpend:       conf.RegistrySettings.MinUpkeepSpend,
						FallbackGasPrice:     conf.RegistrySettings.FallbackGasPrice,
						FallbackLinkPrice:    conf.RegistrySettings.FallbackLinkPrice,
						MaxCheckDataSize:     *conf.RegistrySettings.MaxCheckDataSize,
						MaxPerformDataSize:   *conf.RegistrySettings.MaxPerformDataSize,
						MaxRevertDataSize:    *conf.RegistrySettings.MaxRevertDataSize,
						RegistryVersion:      rv,
					}
					a.RegistrarSettings = contracts.KeeperRegistrarSettings{
						AutoApproveConfigType: uint8(2),
						AutoApproveMaxAllowed: 1000,
						MinLinkJuels:          big.NewInt(0),
					}
					plCfg := config.GetAutomationConfig().AutomationConfig.PluginConfig
					a.PluginConfig = ocr2keepers30config.OffchainConfig{
						TargetProbability:    *plCfg.TargetProbability,
						TargetInRounds:       *plCfg.TargetInRounds,
						PerformLockoutWindow: *plCfg.PerformLockoutWindow,
						GasLimitPerReport:    *plCfg.GasLimitPerReport,
						GasOverheadPerUpkeep: *plCfg.GasOverheadPerUpkeep,
						MinConfirmations:     *plCfg.MinConfirmations,
						MaxUpkeepBatchSize:   *plCfg.MaxUpkeepBatchSize,
					}
					pubCfg := config.GetAutomationConfig().AutomationConfig.PublicConfig
					a.PublicConfig = ocr3.PublicConfig{
						DeltaProgress:                           *pubCfg.DeltaProgress,
						DeltaResend:                             *pubCfg.DeltaResend,
						DeltaInitial:                            *pubCfg.DeltaInitial,
						DeltaRound:                              *pubCfg.DeltaRound,
						DeltaGrace:                              *pubCfg.DeltaGrace,
						DeltaCertifiedCommitRequest:             *pubCfg.DeltaCertifiedCommitRequest,
						DeltaStage:                              *pubCfg.DeltaStage,
						RMax:                                    *pubCfg.RMax,
						MaxDurationQuery:                        *pubCfg.MaxDurationQuery,
						MaxDurationObservation:                  *pubCfg.MaxDurationObservation,
						MaxDurationShouldAcceptAttestedReport:   *pubCfg.MaxDurationShouldAcceptAttestedReport,
						MaxDurationShouldTransmitAcceptedReport: *pubCfg.MaxDurationShouldTransmitAcceptedReport,
						F:                                       *pubCfg.F,
					}

					a.SetupAutomationDeployment(t)

					err = actions.FundChainlinkNodesFromRootAddress(l, a.ChainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes[1:]), big.NewFloat(*config.Common.ChainlinkNodeFunding))
					require.NoError(t, err, "Error funding Chainlink nodes")

					var consumersLogTrigger, consumersConditional []contracts.KeeperConsumer
					var upkeepidsConditional, upkeepidsLogTrigger []*big.Int
					consumersConditional, upkeepidsConditional = actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, numberOfUpkeeps, big.NewInt(defaultLinkFunds), defaultUpkeepGasLimit, false, false, false, nil)
					consumers := consumersConditional
					upkeepIDs := upkeepidsConditional
					if rv >= eth_contracts.RegistryVersion_2_1 {
						consumersLogTrigger, upkeepidsLogTrigger = actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, numberOfUpkeeps, big.NewInt(defaultLinkFunds), defaultUpkeepGasLimit, true, false, false, nil)

						consumers = append(consumersConditional, consumersLogTrigger...)
						upkeepIDs = append(upkeepidsConditional, upkeepidsLogTrigger...)

						for _, c := range consumersLogTrigger {
							err = c.Start()
							require.NoError(t, err, "Error starting consumer")
						}
					}

					l.Info().Msg("Waiting for all upkeeps to be performed")

					gom := gomega.NewGomegaWithT(t)
					gom.Eventually(func(g gomega.Gomega) {
						// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(testcontext.Get(t))
							require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
							expect := 5
							l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
							g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
								"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
						}
					}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

					_, err = testEnvironment.Chaos.Run(testCase.chaosFunc(testEnvironment.Cfg.Namespace, testCase.chaosProps))
					require.NoError(t, err)

					if rv >= eth_contracts.RegistryVersion_2_1 {
						for _, c := range consumersLogTrigger {
							err = c.Start()
							require.NoError(t, err, "Error starting consumer")
						}
					}

					gom.Eventually(func(g gomega.Gomega) {
						// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(testcontext.Get(t))
							require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
							expect := 10
							l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
							g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
								"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
						}
					}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer
				})
			}

		})
	}
}
