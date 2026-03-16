package automation

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
)

const (
	StartupWaitTime = 30 * time.Second
	StopWaitTime    = 60 * time.Second
)

func TestLoad(t *testing.T) {
	testCases := []loadtestcase{
		{
			Testcase: Testcase{
				RegistryVersion:   contracts.RegistryVersion_2_1,
				Name:              "registry_2_1",
				UpkeepCount:       5,
				TestKeyFundingEth: 100,
				UpkeepFundingLink: 1_000_000,
			},
			Load: Load{
				DurationSec:                   10800, // 3h
				NumberOfEvents:                1,
				NumberOfSpamMatchingEvents:    1,
				NumberOfSpamNonMatchingEvents: 0,
				CheckBurnAmount:               big.NewInt(0),
				PerformBurnAmount:             big.NewInt(0),
				UpkeepGasLimit:                1000000,
				SharedTrigger:                 false,
				IsStreamsLookup:               false,
				Feeds:                         []string{"0x000200"},
			},
		},
	}

	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			start := time.Now()

			l := framework.L
			l.Info().Msg("Running test " + tc.Name + " with registry version " + tc.RegistryVersion.String())

			tcStr, mErr := toml.Marshal(tc)
			require.NoError(t, mErr, "failed to marshall test case")
			fmt.Println("------ TEST CONFIGURATION ------")
			fmt.Print(string(tcStr))
			fmt.Println("--------------------------------")

			// dangerous: takes a lot of time, if test runs for a long time
			// t.Cleanup(func() {
			// 	err := products.ScanLogs(l, products.DefaultSettings())
			// 	require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

			// 	_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
			// 	require.NoError(t, cErr)
			// })

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
			require.NoError(t, err)

			// used only to determine which config to use
			isMercuryV02 := strings.Contains(tc.Name, "mercury_v02")
			isMercuryV03 := strings.Contains(tc.Name, "mercury_v03")
			isMercury := isMercuryV02 || isMercuryV03

			var config *automation.Automation
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == tc.RegistryVersion {
					if !isMercury {
						config = candidate
						break
					}

					if isMercuryV02 && candidate.MercurySettings != nil && candidate.MercurySettings.Version == "v2" {
						config = candidate
						break
					}

					if isMercuryV03 && candidate.MercurySettings != nil && candidate.MercurySettings.Version == "v3" {
						config = candidate
						break
					}
				}
			}
			require.NotNil(t, config, "failed to find matching config with registry version %v; mercury v2: %v, mercury v3: %v", tc.RegistryVersion.String(), isMercuryV02, isMercuryV03)

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we ignore key at index 0, because it is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

			pks := []string{products.NetworkPrivateKey()}
			keysRequired := tc.UpkeepCount * 5
			if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != keysRequired {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)
				require.NoError(t, err, "Failed to create ETH client")

				newPks, err := products.FundNewAddresses(t.Context(), keysRequired, c, tc.TestKeyFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), keysRequired+1, "you must provide at least %d private keys", keysRequired+1)
			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			a, err := NewTest(chainClient, config)
			require.NoError(t, err, "Failed to create automation test")

			loadDuration := time.Duration(tc.DurationSec) * time.Second

			startTimeTestSetup := time.Now()
			l.Info().Str("START_TIME", startTimeTestSetup.String()).Msg("Test setup started")

			consumerContracts := make([]contracts.KeeperConsumer, 0)
			triggerContracts := make([]contracts.LogEmitter, 0)
			triggerAddresses := make([]common.Address, 0)

			convenienceABI, err := ac.AutomationCompatibleUtilsMetaData.GetAbi()
			require.NoError(t, err, "Error getting automation utils abi")
			emitterABI, err := log_emitter.LogEmitterMetaData.GetAbi()
			require.NoError(t, err, "Error getting log emitter abi")
			consumerABI, err := simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.GetAbi()
			require.NoError(t, err, "Error getting consumer abi")

			bytes0 := [32]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			}

			bytes1 := [32]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			}

			upkeepConfigs := make([]UpkeepConfig, 0)
			loadConfigs := make([]Load, 0)

			expectedTotalUpkeepCount := tc.UpkeepCount

			maxDeploymentConcurrency := 100

			multicallAddress := common.HexToAddress(config.DeployedContracts.MultiCall)
			deploymentData, err := deployConsumerAndTriggerContracts(l, tc, a.ChainClient, multicallAddress, maxDeploymentConcurrency, big.NewInt(0).Mul(big.NewInt(tc.UpkeepFundingLink), big.NewInt(1e18)), a.LinkToken)
			require.NoError(t, err, "Error deploying consumer and trigger contracts")

			consumerContracts = append(consumerContracts, deploymentData.ConsumerContracts...)
			triggerContracts = append(triggerContracts, deploymentData.TriggerContracts...)
			triggerAddresses = append(triggerAddresses, deploymentData.TriggerAddresses...)
			loadConfigs = append(loadConfigs, deploymentData.LoadConfigs...)

			require.Len(t, consumerContracts, expectedTotalUpkeepCount, "Incorrect number of consumer/trigger contracts deployed")

			for i, consumerContract := range consumerContracts {
				logTriggerConfigStruct := ac.IAutomationV21PlusCommonLogTriggerConfig{
					ContractAddress: triggerAddresses[i],
					FilterSelector:  1,
					Topic0:          emitterABI.Events["Log4"].ID,
					Topic1:          bytes1,
					Topic2:          bytes0,
					Topic3:          bytes0,
				}
				encodedLogTriggerConfig, err := convenienceABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
				require.NoError(t, err, "Error encoding log trigger config")
				l.Debug().
					Interface("logTriggerConfigStruct", logTriggerConfigStruct).
					Str("Encoded Log Trigger Config", hex.EncodeToString(encodedLogTriggerConfig)).Msg("Encoded Log Trigger Config")

				checkDataStruct := simple_log_upkeep_counter_wrapper.CheckData{
					CheckBurnAmount:   loadConfigs[i].CheckBurnAmount,
					PerformBurnAmount: loadConfigs[i].PerformBurnAmount,
					EventSig:          bytes1,
					Feeds:             loadConfigs[i].Feeds,
				}

				encodedCheckDataStruct, err := consumerABI.Methods["_checkDataConfig"].Inputs.Pack(&checkDataStruct)
				require.NoError(t, err, "Error encoding check data struct")
				l.Debug().
					Interface("checkDataStruct", checkDataStruct).
					Str("Encoded Check Data Struct", hex.EncodeToString(encodedCheckDataStruct)).Msg("Encoded Check Data Struct")

				upkeepConfig := UpkeepConfig{
					UpkeepName:     fmt.Sprintf("LogTriggerUpkeep-%d", i),
					EncryptedEmail: []byte("test@mail.com"),
					UpkeepContract: common.HexToAddress(consumerContract.Address()),
					GasLimit:       loadConfigs[i].UpkeepGasLimit,
					AdminAddress:   chainClient.MustGetRootKeyAddress(),
					TriggerType:    uint8(1),
					CheckData:      encodedCheckDataStruct,
					TriggerConfig:  encodedLogTriggerConfig,
					OffchainConfig: []byte(""),
					FundingAmount:  big.NewInt(0).Mul(big.NewInt(tc.UpkeepFundingLink), big.NewInt(1e18)),
				}
				l.Debug().Interface("Upkeep Config", upkeepConfig).Msg("Upkeep Config")
				upkeepConfigs = append(upkeepConfigs, upkeepConfig)
			}

			require.Len(t, upkeepConfigs, expectedTotalUpkeepCount, "Incorrect number of upkeep configs created")
			registrationTxHashes, err := a.RegisterUpkeeps(upkeepConfigs, maxDeploymentConcurrency)
			require.NoError(t, err, "Error registering upkeeps")

			upkeepIDs, err := a.ConfirmUpkeepsRegistered(registrationTxHashes, maxDeploymentConcurrency)
			require.NoError(t, err, "Error confirming upkeeps registered")
			require.Len(t, upkeepIDs, expectedTotalUpkeepCount, "Incorrect number of upkeeps registered")

			l.Info().Msg("Successfully registered all Automation Upkeeps")
			l.Info().Interface("Upkeep IDs", upkeepIDs).Msg("Upkeeps Registered")
			l.Info().Str("STARTUP_WAIT_TIME", StartupWaitTime.String()).Msg("Waiting for plugin to start")
			time.Sleep(StartupWaitTime)

			startBlock, err := a.ChainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Error getting latest block number")

			p := wasp.NewProfile()

			configs := make([]LogTriggerConfig, 0)
			var numberOfEventsEmitted int64
			var numberOfEventsEmittedPerSec int64

			for i, triggerContract := range triggerContracts {
				c := LogTriggerConfig{
					Address:                       triggerContract.Address().String(),
					NumberOfEvents:                int64(loadConfigs[i].NumberOfEvents),
					NumberOfSpamMatchingEvents:    int64(loadConfigs[i].NumberOfSpamMatchingEvents),
					NumberOfSpamNonMatchingEvents: int64(loadConfigs[i].NumberOfSpamNonMatchingEvents),
				}
				numberOfEventsEmittedPerSec += int64(loadConfigs[i].NumberOfEvents)
				configs = append(configs, c)
			}

			endTimeTestSetup := time.Now()
			testSetupDuration := endTimeTestSetup.Sub(startTimeTestSetup)
			l.Info().
				Str("END_TIME", endTimeTestSetup.String()).
				Str("Duration", testSetupDuration.String()).
				Msg("Test setup ended")

			gun, gErr := NewLogTriggerUser(
				l,
				configs,
				a.ChainClient,
				multicallAddress.Hex(),
			)
			require.NoError(t, gErr, "failed to create LogTriggerUser WASP gun")

			g, err := wasp.NewGenerator(&wasp.Config{
				T:           t,
				LoadType:    wasp.RPS,
				GenName:     "log_trigger_gen",
				CallTimeout: time.Minute * 3,
				Schedule: wasp.Plain(
					1,
					loadDuration,
				),
				Gun:              gun,
				CallResultBufLen: 1000,
			})
			p.Add(g, err)

			startTimeTestEx := time.Now()
			l.Info().Str("START_TIME", startTimeTestEx.String()).Msg("Test execution started")

			l.Info().Msg("Starting load generators")
			_, err = p.Run(true)
			require.NoError(t, err, "Error running load generators")

			l.Info().Msg("Finished load generators")
			l.Info().Str("STOP_WAIT_TIME", StopWaitTime.String()).Msg("Waiting for upkeeps to be performed")
			time.Sleep(StopWaitTime)
			l.Info().Msg("Finished waiting 60s for upkeeps to be performed")
			endTimeTestEx := time.Now()
			testExDuration := endTimeTestEx.Sub(startTimeTestEx)
			l.Info().
				Str("END_TIME", endTimeTestEx.String()).
				Str("Duration", testExDuration.String()).
				Msg("Test execution ended")

			l.Info().Str("Duration", testExDuration.String()).Msg("Test Execution Duration")
			endBlock, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Error getting latest block number")
			l.Info().Uint64("Starting Block", startBlock).Uint64("Ending Block", endBlock).Msg("Test Block Range")

			startTimeTestReport := time.Now()
			l.Info().Str("START_TIME", startTimeTestReport.String()).Msg("Test reporting started")

			for _, gen := range p.Generators {
				if len(gen.Errors()) != 0 {
					l.Error().Strs("Errors", gen.Errors()).Msg("Error in load gen")
					t.Fail()
				}
			}

			upkeepDelaysFast := make([][]int64, 0)
			upkeepDelaysRecovery := make([][]int64, 0)

			var batchSize uint64 = 500

			if endBlock-startBlock < batchSize {
				batchSize = endBlock - startBlock
			}

			for _, consumerContract := range consumerContracts {
				var (
					logs    []types.Log
					address = common.HexToAddress(consumerContract.Address())
					timeout = 5 * time.Second
				)
				for fromBlock := startBlock; fromBlock < endBlock; fromBlock += batchSize + 1 {
					filterQuery := geth.FilterQuery{
						Addresses: []common.Address{address},
						FromBlock: big.NewInt(0).SetUint64(fromBlock),
						ToBlock:   big.NewInt(0).SetUint64(fromBlock + batchSize),
						Topics:    [][]common.Hash{{consumerABI.Events["PerformingUpkeep"].ID}},
					}
					err = errors.New("initial error") // to ensure our for loop runs at least once
					for err != nil {
						var logsInBatch []types.Log
						ctx2, cancel := context.WithTimeout(t.Context(), timeout)
						logsInBatch, err = a.ChainClient.Client.FilterLogs(ctx2, filterQuery)
						cancel()
						if err != nil {
							l.Error().Err(err).
								Interface("FilterQuery", filterQuery).
								Str("Contract Address", consumerContract.Address()).
								Str("Timeout", timeout.String()).
								Msg("Error getting consumer contract logs")
							timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
							continue
						}
						l.Debug().
							Interface("FilterQuery", filterQuery).
							Str("Contract Address", consumerContract.Address()).
							Str("Timeout", timeout.String()).
							Int("Number of Logs", len(logsInBatch)).
							Msg("Collected consumer contract logs")
						logs = append(logs, logsInBatch...)
					}
				}

				if len(logs) > 0 {
					delayFast := make([]int64, 0)
					delayRecovery := make([]int64, 0)
					for _, log := range logs {
						eventDetails, err := consumerABI.EventByID(log.Topics[0])
						require.NoError(t, err, "Error getting event details")
						consumer, err := simple_log_upkeep_counter_wrapper.NewSimpleLogUpkeepCounter(
							address, a.ChainClient.Client,
						)
						require.NoError(t, err, "Error getting consumer contract")
						if eventDetails.Name == "PerformingUpkeep" {
							parsedLog, err := consumer.ParsePerformingUpkeep(log)
							require.NoError(t, err, "Error parsing log")
							if parsedLog.IsRecovered {
								delayRecovery = append(delayRecovery, parsedLog.TimeToPerform.Int64())
							} else {
								delayFast = append(delayFast, parsedLog.TimeToPerform.Int64())
							}
						}
					}
					upkeepDelaysFast = append(upkeepDelaysFast, delayFast)
					upkeepDelaysRecovery = append(upkeepDelaysRecovery, delayRecovery)
				}
			}

			for _, triggerContract := range triggerContracts {
				var (
					logs    []types.Log
					address = triggerContract.Address()
					timeout = 5 * time.Second
				)
				for fromBlock := startBlock; fromBlock < endBlock; fromBlock += batchSize + 1 {
					filterQuery := geth.FilterQuery{
						Addresses: []common.Address{address},
						FromBlock: big.NewInt(0).SetUint64(fromBlock),
						ToBlock:   big.NewInt(0).SetUint64(fromBlock + batchSize),
						Topics:    [][]common.Hash{{emitterABI.Events["Log4"].ID}, {bytes1}, {bytes1}},
					}
					err = errors.New("initial error") // to ensure our for loop runs at least once
					for err != nil {
						var logsInBatch []types.Log
						ctx2, cancel := context.WithTimeout(t.Context(), timeout)
						logsInBatch, err = chainClient.Client.FilterLogs(ctx2, filterQuery)
						cancel()
						if err != nil {
							l.Error().Err(err).
								Interface("FilterQuery", filterQuery).
								Str("Contract Address", address.Hex()).
								Str("Timeout", timeout.String()).
								Msg("Error getting trigger contract logs")
							timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
							continue
						}
						l.Debug().
							Interface("FilterQuery", filterQuery).
							Str("Contract Address", address.Hex()).
							Str("Timeout", timeout.String()).
							Int("Number of Logs", len(logsInBatch)).
							Msg("Collected trigger contract logs")
						logs = append(logs, logsInBatch...)
					}
				}
				numberOfEventsEmitted += int64(len(logs))
			}

			l.Info().Int64("Number of Events Emitted", numberOfEventsEmitted).Msg("Number of Events Emitted")

			l.Info().
				Interface("Upkeep Delays Fast", upkeepDelaysFast).
				Interface("Upkeep Delays Recovered", upkeepDelaysRecovery).
				Msg("Upkeep Delays")

			var allUpkeepDelays []int64
			var allUpkeepDelaysFast []int64
			var allUpkeepDelaysRecovery []int64

			for _, upkeepDelay := range upkeepDelaysFast {
				allUpkeepDelays = append(allUpkeepDelays, upkeepDelay...)
				allUpkeepDelaysFast = append(allUpkeepDelaysFast, upkeepDelay...)
			}

			for _, upkeepDelay := range upkeepDelaysRecovery {
				allUpkeepDelays = append(allUpkeepDelays, upkeepDelay...)
				allUpkeepDelaysRecovery = append(allUpkeepDelaysRecovery, upkeepDelay...)
			}

			avgF, medianF, ninetyPctF, ninetyNinePctF, maximumF := IntListStats(allUpkeepDelaysFast)
			avgR, medianR, ninetyPctR, ninetyNinePctR, maximumR := IntListStats(allUpkeepDelaysRecovery)
			eventsMissed := (numberOfEventsEmitted) - int64(len(allUpkeepDelays))
			percentMissed := float64(eventsMissed) / float64(numberOfEventsEmitted) * 100
			l.Info().
				Float64("Average", avgF).Int64("Median", medianF).
				Int64("90th Percentile", ninetyPctF).Int64("99th Percentile", ninetyNinePctF).
				Int64("Max", maximumF).Msg("Upkeep Delays Fast Execution in seconds")
			l.Info().
				Float64("Average", avgR).Int64("Median", medianR).
				Int64("90th Percentile", ninetyPctR).Int64("99th Percentile", ninetyNinePctR).
				Int64("Max", maximumR).Msg("Upkeep Delays Recovery Execution in seconds")
			l.Info().
				Int("Total Perform Count", len(allUpkeepDelays)).
				Int("Perform Count Fast Execution", len(allUpkeepDelaysFast)).
				Int("Perform Count Recovery Execution", len(allUpkeepDelaysRecovery)).
				Int64("Total Events Emitted", numberOfEventsEmitted).
				Int64("Total Events Missed", eventsMissed).
				Float64("Percent Missed", percentMissed).
				Msg("Test completed")

			testReportFormat := `Upkeep Delays in seconds - Fast Execution
Average: %f
Median: %d
90th Percentile: %d
99th Percentile: %d
Max: %d

Upkeep Delays in seconds - Recovery Execution
Average: %f
Median: %d
90th Percentile: %d
99th Percentile: %d
Max: %d

Total Perform Count: %d
Perform Count Fast Execution: %d
Perform Count Recovery Execution: %d
Total Expected Log Triggering Events: %d
Total Log Triggering Events Emitted: %d
Total Events Missed: %d
Percent Missed: %f
Test Duration: %s`

			endTimeTestReport := time.Now()
			testReDuration := endTimeTestReport.Sub(startTimeTestReport)
			l.Info().
				Str("END_TIME", endTimeTestReport.String()).
				Str("Duration", testReDuration.String()).
				Msg("Test reporting ended")

			numberOfExpectedEvents := numberOfEventsEmittedPerSec * int64(loadDuration.Seconds())
			testReport := fmt.Sprintf(testReportFormat, avgF, medianF, ninetyPctF, ninetyNinePctF, maximumF,
				avgR, medianR, ninetyPctR, ninetyNinePctR, maximumR, len(allUpkeepDelays), len(allUpkeepDelaysFast),
				len(allUpkeepDelaysRecovery), numberOfExpectedEvents, numberOfEventsEmitted, eventsMissed, percentMissed, testExDuration.String())
			l.Info().Str("Test Report", testReport).Msg("Test Report prepared")

			// it might happen that the number of events emitted is less than expected due to the fact that sometimes Anvil gets stuck
			// we cannot get any private key without a pending transaction and we need to drop all pending transactions
			diff := numberOfExpectedEvents - numberOfEventsEmitted
			maxDiff := int64(float64(numberOfExpectedEvents) * 0.05)
			if diff-maxDiff > 0 {
				l.Error().
					Int64("number of events emitted", numberOfEventsEmitted).
					Int64("number of expected events", numberOfExpectedEvents).
					Int64("Difference", diff).
					Int64("Max Difference", maxDiff).
					Msg("Number of events emitted is less than expected")
				t.FailNow()
			}

			require.LessOrEqual(t, percentMissed, 10.0, "Too many events were missed")

			leaks, err := leak.NewCLNodesLeakDetector(leak.NewResourceLeakChecker())
			require.NoError(t, err)
			errs := leaks.Check(&leak.CLNodesCheck{
				// since the test is stable we assert absolute values
				// no more than 30% CPU and 200Mb (last 5m)
				ComparisonMode:  leak.ComparisonModeAbsolute,
				NumNodes:        in.NodeSets[0].Nodes,
				Start:           start,
				End:             time.Now(),
				WarmUpDuration:  30 * time.Minute,
				CPUThreshold:    30.0,
				MemoryThreshold: 240.0, // max observed so far was 227 mb, adding a buffer to be safe
			})
			require.NoError(t, errs)
		})
	}
}
