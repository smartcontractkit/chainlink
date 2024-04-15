package automationv2_1

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/seth"

	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	aconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
)

func extraBlockWithText(text string) slack.Block {
	return slack.NewSectionBlock(slack.NewTextBlockObject(
		"mrkdwn", text, false, false), nil, nil)
}

func sendSlackNotification(header string, l zerolog.Logger, config *tc.TestConfig, namespace string, numberOfNodes,
	startingTime string, endingTime string, extraBlocks []slack.Block, msgOption slack.MsgOption) (string, error) {
	slackClient := slack.New(reportModel.SlackAPIKey)

	headerText := ":chainlink-keepers: Automation Load Test " + header

	grafanaUrl, err := config.GetGrafanaBaseURL()
	if err != nil {
		return "", err
	}

	dashboardUrl, err := config.GetGrafanaDashboardURL()
	if err != nil {
		return "", err
	}

	formattedDashboardUrl := fmt.Sprintf("%s%s?orgId=1&from=%s&to=%s&var-namespace=%s&var-number_of_nodes=%s", grafanaUrl, dashboardUrl, startingTime, endingTime, namespace, numberOfNodes)
	l.Info().Str("Dashboard", formattedDashboardUrl).Msg("Dashboard URL")

	var notificationBlocks []slack.Block

	notificationBlocks = append(notificationBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks = append(notificationBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", namespace, false, false)))
	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())
	if *config.Pyroscope.Enabled {
		pyroscopeServer := *config.Pyroscope.ServerUrl
		pyroscopeEnvironment := *config.Pyroscope.Environment

		formattedPyroscopeUrl := fmt.Sprintf("%s/?query=chainlink-node.cpu{Environment=\"%s\"}&from=%s&to=%s", pyroscopeServer, pyroscopeEnvironment, startingTime, endingTime)
		l.Info().Str("Pyroscope", formattedPyroscopeUrl).Msg("Dashboard URL")
		notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("<%s|Pyroscope>",
				formattedPyroscopeUrl), false, true), nil, nil))
	}
	notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("<%s|Test Dashboard> \nNotifying <@%s>",
			formattedDashboardUrl, reportModel.SlackUserID), false, true), nil, nil))

	if len(extraBlocks) > 0 {
		notificationBlocks = append(notificationBlocks, extraBlocks...)
	}

	ts, err := reportModel.SendSlackMessage(slackClient, slack.MsgOptionBlocks(notificationBlocks...), msgOption)
	l.Info().Str("ts", ts).Msg("Sent Slack Message")
	return ts, err
}

type DeploymentData struct {
	ConsumerContracts []contracts.KeeperConsumer
	TriggerContracts  []contracts.LogEmitter
	TriggerAddresses  []common.Address
	LoadConfigs       []aconfig.Load
}

func deployConsumerAndTriggerContracts(l zerolog.Logger, loadConfig aconfig.Load, chainClient *seth.Client, concurrency int, multicallAddress common.Address, automationDefaultLinkFunds *big.Int, linkToken contracts.LinkToken) (DeploymentData, error) {
	data := DeploymentData{}

	type deployedContractData struct {
		consumerContract contracts.KeeperConsumer
		triggerContract  contracts.LogEmitter
		triggerAddress   common.Address
		loadConfig       aconfig.Load
		err              error
	}

	l.Debug().
		Int("Number of Upkeeps", *loadConfig.NumberOfUpkeeps).
		Int("Concurrency", concurrency).
		Msg("Deployment parallelisation info")

	atomicCounter := atomic.Uint64{}
	var deplymentErr error
	deployedContractCh := make(chan deployedContractData, concurrency)
	stopCh := make(chan struct{})

	var deployContractFn = func(deployedCh chan deployedContractData, keyNum int) {
		data := deployedContractData{}
		consumerContract, err := contracts.DeployAutomationSimpleLogTriggerConsumerFromKey(chainClient, *loadConfig.IsStreamsLookup, keyNum)
		if err != nil {
			data.err = err
			deployedCh <- data
			return
		}

		data.consumerContract = consumerContract

		atomicCounter.Add(1)

		l.Debug().
			Str("Contract Address", consumerContract.Address()).
			Int("Number", int(atomicCounter.Load())).
			Int("Out Of", *loadConfig.NumberOfUpkeeps).
			Int("Key Number", keyNum).
			Msg("Deployed Automation Log Trigger Consumer Contract")

		loadCfg := aconfig.Load{
			NumberOfEvents:                loadConfig.NumberOfEvents,
			NumberOfSpamMatchingEvents:    loadConfig.NumberOfSpamMatchingEvents,
			NumberOfSpamNonMatchingEvents: loadConfig.NumberOfSpamNonMatchingEvents,
			CheckBurnAmount:               loadConfig.CheckBurnAmount,
			PerformBurnAmount:             loadConfig.PerformBurnAmount,
			UpkeepGasLimit:                loadConfig.UpkeepGasLimit,
			SharedTrigger:                 loadConfig.SharedTrigger,
			Feeds:                         []string{},
		}

		if *loadConfig.IsStreamsLookup {
			loadCfg.Feeds = loadConfig.Feeds
		}

		data.loadConfig = loadCfg

		// deploy only 1 trigger contract if it's a shared trigger
		if *loadConfig.SharedTrigger && atomicCounter.Load() > 1 {
			deployedCh <- data
			return
		}

		triggerContract, err := contracts.DeployLogEmitterContractFromKey(l, chainClient, keyNum)
		if err != nil {
			data.err = err
			deployedCh <- data
			return
		}

		data.triggerContract = triggerContract
		data.triggerAddress = triggerContract.Address()
		l.Debug().
			Str("Contract Address", triggerContract.Address().Hex()).
			Int("Number", int(atomicCounter.Load())).
			Int("Out Of", *loadConfig.NumberOfUpkeeps).
			Int("Key Number", keyNum).
			Msg("Deployed Automation Log Trigger Emitter Contract")

		deployedCh <- data
	}

	var wgProcess sync.WaitGroup
	for i := 0; i < *loadConfig.NumberOfUpkeeps; i++ {
		wgProcess.Add(1)
	}

	go func() {
		defer l.Debug().Msg("Finished listening to results of deploying consumer/trigger contracts")
		for contractData := range deployedContractCh {
			if contractData.err != nil {
				l.Error().Err(contractData.err).Msg("Error deploying customer/trigger contract")
				deplymentErr = contractData.err
				close(stopCh)
				return
			}
			if contractData.triggerContract != nil {
				data.TriggerContracts = append(data.TriggerContracts, contractData.triggerContract)
				data.TriggerAddresses = append(data.TriggerAddresses, contractData.triggerAddress)
			}
			data.ConsumerContracts = append(data.ConsumerContracts, contractData.consumerContract)
			data.LoadConfigs = append(data.LoadConfigs, contractData.loadConfig)
			wgProcess.Done()
		}
	}()

	operationsPerClient := *loadConfig.NumberOfUpkeeps / concurrency
	extraOperations := *loadConfig.NumberOfUpkeeps % concurrency

	for clientNum := 1; clientNum <= concurrency; clientNum++ {
		go func(key int) {
			numTasks := operationsPerClient
			if key <= extraOperations {
				numTasks++
			}

			if numTasks == 0 {
				return
			}

			l.Debug().
				Int("Key Number", key).
				Int("Number of Tasks", numTasks).
				Msg("Started deploying consumer/trigger contracts")

			for i := 0; i < numTasks; i++ {
				select {
				case <-stopCh:
					return
				default:
					deployContractFn(deployedContractCh, key)
					l.Trace().
						Int("Key Number", key).
						Msgf("Finished consumer/trigger contract deployment task %d/%d", (i + 1), numTasks)
				}
			}

			l.Debug().
				Int("Key Number", key).
				Msg("Finished deploying consumer/trigger contracts")
		}(clientNum)
	}

	wgProcess.Wait()
	close(deployedContractCh)

	if deplymentErr != nil {
		return DeploymentData{}, deplymentErr
	}

	// if there's more than 1 upkeep and it's a shared trigger, then we should use only the first address in triggerAddresses
	// as triggerAddresses array
	if *loadConfig.SharedTrigger {
		if len(data.TriggerAddresses) == 0 {
			return DeploymentData{}, errors.New("No trigger addresses found")
		}
		triggerAddress := data.TriggerAddresses[0]
		data.TriggerAddresses = make([]common.Address, 0)
		for i := 0; i < *loadConfig.NumberOfUpkeeps; i++ {
			data.TriggerAddresses = append(data.TriggerAddresses, triggerAddress)
		}
	}

	sendErr := actions_seth.SendLinkFundsToDeploymentAddresses(chainClient, concurrency, *loadConfig.NumberOfUpkeeps, operationsPerClient, multicallAddress, automationDefaultLinkFunds, linkToken)
	if sendErr != nil {
		return DeploymentData{}, sendErr
	}

	return data, nil
}
