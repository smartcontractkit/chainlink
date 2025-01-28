package automation

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	ctf_concurrency "github.com/smartcontractkit/chainlink-testing-framework/lib/concurrency"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
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

	grafanaURL, err := config.GetGrafanaBaseURL()
	if err != nil {
		return "", err
	}

	dashboardURL, err := config.GetGrafanaDashboardURL()
	if err != nil {
		return "", err
	}

	formattedDashboardURL := fmt.Sprintf("%s%s?orgId=1&from=%s&to=%s&var-namespace=%s&var-number_of_nodes=%s", grafanaURL, dashboardURL, startingTime, endingTime, namespace, numberOfNodes)
	l.Info().Str("Dashboard", formattedDashboardURL).Msg("Dashboard URL")

	var notificationBlocks []slack.Block

	notificationBlocks = append(notificationBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks = append(notificationBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", namespace, false, false)))
	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())
	if *config.Pyroscope.Enabled {
		pyroscopeServer := *config.Pyroscope.ServerUrl
		pyroscopeEnvironment := *config.Pyroscope.Environment

		formattedPyroscopeURL := fmt.Sprintf("%s/?query=chainlink-node.cpu{Environment=\"%s\"}&from=%s&to=%s", pyroscopeServer, pyroscopeEnvironment, startingTime, endingTime)
		l.Info().Str("Pyroscope", formattedPyroscopeURL).Msg("Dashboard URL")
		notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("<%s|Pyroscope>",
				formattedPyroscopeURL), false, true), nil, nil))
	}
	notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("<%s|Test Dashboard> \nNotifying <@%s>",
			formattedDashboardURL, reportModel.SlackUserID), false, true), nil, nil))

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

type deployedContractData struct {
	consumerContract contracts.KeeperConsumer
	triggerContract  contracts.LogEmitter
	triggerAddress   common.Address
	loadConfig       aconfig.Load
}

func (d deployedContractData) GetResult() deployedContractData {
	return d
}

type task struct {
	deployTrigger bool
}

func deployConsumerAndTriggerContracts(l zerolog.Logger, loadConfig aconfig.Load, chainClient *seth.Client, multicallAddress common.Address, maxConcurrency int, automationDefaultLinkFunds *big.Int, linkToken contracts.LinkToken) (DeploymentData, error) {
	data := DeploymentData{}

	concurrency, err := actions.GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return DeploymentData{}, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		l.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	l.Debug().
		Int("Number of Upkeeps", *loadConfig.NumberOfUpkeeps).
		Int("Concurrency", concurrency).
		Msg("Deployment parallelisation info")

	tasks := []task{}
	for i := 0; i < *loadConfig.NumberOfUpkeeps; i++ {
		if *loadConfig.SharedTrigger {
			if i == 0 {
				tasks = append(tasks, task{deployTrigger: true})
			} else {
				tasks = append(tasks, task{deployTrigger: false})
			}
			continue
		}
		tasks = append(tasks, task{deployTrigger: true})
	}

	var deployContractFn = func(deployedCh chan deployedContractData, errorCh chan error, keyNum int, task task) {
		data := deployedContractData{}
		consumerContract, err := contracts.DeployAutomationSimpleLogTriggerConsumerFromKey(chainClient, *loadConfig.IsStreamsLookup, keyNum)
		if err != nil {
			errorCh <- errors.Wrapf(err, "Error deploying simple log trigger contract")
			return
		}

		data.consumerContract = consumerContract

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

		if !task.deployTrigger {
			deployedCh <- data
			return
		}

		triggerContract, err := contracts.DeployLogEmitterContractFromKey(l, chainClient, keyNum)
		if err != nil {
			errorCh <- errors.Wrapf(err, "Error deploying log emitter contract")
			return
		}

		data.triggerContract = triggerContract
		data.triggerAddress = triggerContract.Address()
		deployedCh <- data
	}

	executor := ctf_concurrency.NewConcurrentExecutor[deployedContractData, deployedContractData, task](l)
	results, err := executor.Execute(concurrency, tasks, deployContractFn)
	if err != nil {
		return DeploymentData{}, err
	}

	for _, result := range results {
		if result.GetResult().triggerContract != nil {
			data.TriggerContracts = append(data.TriggerContracts, result.GetResult().triggerContract)
			data.TriggerAddresses = append(data.TriggerAddresses, result.GetResult().triggerAddress)
		}
		data.ConsumerContracts = append(data.ConsumerContracts, result.GetResult().consumerContract)
		data.LoadConfigs = append(data.LoadConfigs, result.GetResult().loadConfig)
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

	sendErr := actions.SendLinkFundsToDeploymentAddresses(chainClient, concurrency, *loadConfig.NumberOfUpkeeps, *loadConfig.NumberOfUpkeeps/concurrency, multicallAddress, automationDefaultLinkFunds, linkToken)
	if sendErr != nil {
		return DeploymentData{}, sendErr
	}

	return data, nil
}
