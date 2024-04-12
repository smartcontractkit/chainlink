package automationv2_1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/seth"

	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	aconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
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

	sendErr := sendLinkFundsToDepolymentAddresses(chainClient, concurrency, *loadConfig.NumberOfUpkeeps, operationsPerClient, multicallAddress, automationDefaultLinkFunds, linkToken)
	if sendErr != nil {
		return DeploymentData{}, sendErr
	}

	return data, nil
}

func sendLinkFundsToDepolymentAddresses(
	chainClient *seth.Client,
	concurrency,
	totalUpkeeps,
	operationsPerAddress int,
	multicallAddress common.Address,
	automationDefaultLinkFunds *big.Int,
	linkToken contracts.LinkToken,
) error {
	var generateCallData = func(receiver common.Address, amount *big.Int) ([]byte, error) {
		abi, err := link_token_interface.LinkTokenMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
		data, err := abi.Pack("transfer", receiver, amount)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if *chainClient.Cfg.EphemeralAddrs == 0 {
		return nil
	}

	toTransferToMultiCallContract := big.NewInt(0).Mul(automationDefaultLinkFunds, big.NewInt(int64(totalUpkeeps+concurrency)))
	toTransferPerClient := big.NewInt(0).Mul(automationDefaultLinkFunds, big.NewInt(int64(operationsPerAddress+1)))
	err := linkToken.Transfer(multicallAddress.Hex(), toTransferToMultiCallContract)
	if err != nil {
		return errors.Join(err, errors.New("Error transferring LINK to multicall contract"))
	}

	balance, err := linkToken.BalanceOf(context.Background(), multicallAddress.Hex())
	if err != nil {
		return errors.Join(err, errors.New("Error getting LINK balance of multicall contract"))
	}

	if toTransferToMultiCallContract.Cmp(balance) != 0 {
		return fmt.Errorf("Incorrect LINK balance of multicall contract. Expected: %s. Got: %s", toTransferToMultiCallContract.String(), balance.String())
	}

	// Transfer LINK to ephemeral keys
	multiCallData := make([][]byte, 0)
	for i := 1; i <= concurrency; i++ {
		data, err := generateCallData(chainClient.Addresses[i], toTransferPerClient)
		if err != nil {
			return errors.Join(err, errors.New("Error generating call data for LINK transfer"))
		}
		multiCallData = append(multiCallData, data)
	}

	var call []contracts.Call
	for _, d := range multiCallData {
		data := contracts.Call{Target: common.HexToAddress(linkToken.Address()), AllowFailure: false, CallData: d}
		call = append(call, data)
	}

	multiCallABI, err := abi.JSON(strings.NewReader(contracts.MultiCallABI))
	if err != nil {
		return errors.Join(err, errors.New("Error getting Multicall contract ABI"))
	}
	boundContract := bind.NewBoundContract(multicallAddress, multiCallABI, chainClient.Client, chainClient.Client, chainClient.Client)
	// call aggregate3 to group all msg call data and send them in a single transaction
	_, err = chainClient.Decode(boundContract.Transact(chainClient.NewTXOpts(), "aggregate3", call))
	if err != nil {
		return errors.Join(err, errors.New("Error calling Multicall contract"))
	}

	for i := 1; i <= concurrency; i++ {
		balance, err := linkToken.BalanceOf(context.Background(), chainClient.Addresses[i].Hex())
		if err != nil {
			return errors.Join(err, fmt.Errorf("Error getting LINK balance of ephemeral key %d", i))
		}
		if toTransferPerClient.Cmp(balance) != 0 {
			return fmt.Errorf("Incorrect LINK balance after transferring for ephemeral key %d. Expected: %s. Got: %s", i, toTransferPerClient.String(), balance.String())
		}
	}

	return nil
}
