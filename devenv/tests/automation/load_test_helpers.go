package automation

import (
	"math"
	"math/big"
	"slices"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
	ctf_concurrency "github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"
)

// logTriggerCall pairs an address with its calldata to ensure they stay together when chunking
type logTriggerCall struct {
	address  string
	callData []byte
}

type DeploymentData struct {
	ConsumerContracts []contracts.KeeperConsumer
	TriggerContracts  []contracts.LogEmitter
	TriggerAddresses  []common.Address
	LoadConfigs       []Load
}

type deployedContractData struct {
	consumerContract contracts.KeeperConsumer
	triggerContract  contracts.LogEmitter
	triggerAddress   common.Address
	loadConfig       Load
}

func (d deployedContractData) GetResult() deployedContractData {
	return d
}

type task struct {
	deployTrigger bool
}

func deployConsumerAndTriggerContracts(l zerolog.Logger, tc loadtestcase, chainClient *seth.Client, multicallAddress common.Address, maxConcurrency int, automationDefaultLinkFunds *big.Int, linkToken contracts.LinkToken) (DeploymentData, error) {
	data := DeploymentData{}

	concurrency, err := automation.GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return DeploymentData{}, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		l.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	l.Debug().
		Int("Number of Upkeeps", tc.UpkeepCount).
		Int("Concurrency", concurrency).
		Msg("Deployment parallelisation info")

	tasks := []task{}
	for i := 0; i < tc.UpkeepCount; i++ {
		if tc.SharedTrigger {
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
		consumerContract, err := contracts.DeployAutomationSimpleLogTriggerConsumerFromKey(chainClient, tc.IsStreamsLookup, keyNum)
		if err != nil {
			errorCh <- errors.Wrapf(err, "Error deploying simple log trigger contract")
			return
		}

		data.consumerContract = consumerContract

		loadCfg := Load{
			NumberOfEvents:                tc.NumberOfEvents,
			NumberOfSpamMatchingEvents:    tc.NumberOfSpamMatchingEvents,
			NumberOfSpamNonMatchingEvents: tc.NumberOfSpamNonMatchingEvents,
			CheckBurnAmount:               tc.CheckBurnAmount,
			PerformBurnAmount:             tc.PerformBurnAmount,
			UpkeepGasLimit:                tc.UpkeepGasLimit,
			SharedTrigger:                 tc.SharedTrigger,
			Feeds:                         []string{},
		}

		if tc.IsStreamsLookup {
			loadCfg.Feeds = tc.Feeds
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
	if tc.SharedTrigger {
		if len(data.TriggerAddresses) == 0 {
			return DeploymentData{}, errors.New("No trigger addresses found")
		}
		triggerAddress := data.TriggerAddresses[0]
		data.TriggerAddresses = make([]common.Address, 0)
		for i := 0; i < tc.UpkeepCount; i++ {
			data.TriggerAddresses = append(data.TriggerAddresses, triggerAddress)
		}
	}

	sendErr := automation.SendLinkFundsToDeploymentAddresses(chainClient, concurrency, tc.UpkeepCount, tc.UpkeepCount/concurrency, multicallAddress, automationDefaultLinkFunds, linkToken)
	if sendErr != nil {
		return DeploymentData{}, sendErr
	}

	return data, nil
}

type LogTriggerConfig struct {
	Address                       string
	NumberOfEvents                int64
	NumberOfSpamMatchingEvents    int64
	NumberOfSpamNonMatchingEvents int64
}

type LogTriggerGun struct {
	calls            []logTriggerCall
	multiCallAddress string
	client           *seth.Client
	logger           zerolog.Logger
}

func generateCallData(int1 int64, int2 int64, count int64) []byte {
	abi, err := log_emitter.LogEmitterMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	data, err := abi.Pack("EmitLog4", big.NewInt(int1), big.NewInt(int2), big.NewInt(count))
	if err != nil {
		panic(err)
	}
	return data
}

func NewLogTriggerUser(
	logger zerolog.Logger,
	triggerConfigs []LogTriggerConfig,
	client *seth.Client,
	multicallAddress string,
) (*LogTriggerGun, error) {
	var calls []logTriggerCall

	// we need to sync nodes manually, because we are not using ephemeral addresses
	if err := client.NonceManager.UpdateNonces(); err != nil {
		return nil, err
	}

	for _, c := range triggerConfigs {
		if c.NumberOfEvents > 0 {
			d := generateCallData(1, 1, c.NumberOfEvents)
			calls = append(calls, logTriggerCall{address: c.Address, callData: d})
		}
		if c.NumberOfSpamMatchingEvents > 0 {
			d := generateCallData(1, 2, c.NumberOfSpamMatchingEvents)
			calls = append(calls, logTriggerCall{address: c.Address, callData: d})
		}
		if c.NumberOfSpamNonMatchingEvents > 0 {
			d := generateCallData(2, 2, c.NumberOfSpamNonMatchingEvents)
			calls = append(calls, logTriggerCall{address: c.Address, callData: d})
		}
	}

	return &LogTriggerGun{
		calls:            calls,
		logger:           logger,
		multiCallAddress: multicallAddress,
		client:           client,
	}, nil
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	var wg sync.WaitGroup

	// Chunk the paired calls to ensure addresses stay aligned with their calldata
	var dividedCalls [][]logTriggerCall
	chunkSize := 100
	for i := 0; i < len(m.calls); i += chunkSize {
		end := min(i+chunkSize, len(m.calls))
		dividedCalls = append(dividedCalls, m.calls[i:end])
	}

	resultCh := make(chan *wasp.Response, len(dividedCalls))

	for _, chunk := range dividedCalls {
		wg.Add(1)
		go func(chunk []logTriggerCall, m *LogTriggerGun) {
			defer wg.Done()

			// Convert to contracts.Call slice
			calls := make([]contracts.Call, len(chunk))
			for i, c := range chunk {
				calls[i] = contracts.Call{
					Target:       common.HexToAddress(c.address),
					AllowFailure: false,
					CallData:     c.callData,
				}
			}

			_, err := contracts.MultiCallLogTriggerLoadGen(m.client, m.multiCallAddress, calls)
			if err != nil {
				m.logger.Error().Err(err).Msg("Error calling MultiCallLogTriggerLoadGen")
				resultCh <- &wasp.Response{Error: err.Error(), Failed: true}
				return
			}

			resultCh <- &wasp.Response{}
		}(chunk, m)
	}

	wg.Wait()
	close(resultCh)

	r := &wasp.Response{}
	for result := range resultCh {
		if result.Failed {
			r.Failed = true
			if r.Error != "" {
				r.Error += "; " + result.Error
			} else {
				r.Error = result.Error
			}
		}
	}

	return r
}

// intListStats helper calculates some statistics on an int list: avg, median, 90pct, 99pct, max
//
//nolint:revive // we know what each int64 we return means
func IntListStats(in []int64) (float64, int64, int64, int64, int64) {
	length := len(in)
	if length == 0 {
		return 0, 0, 0, 0, 0
	}
	slices.Sort(in)
	var sum int64
	for _, num := range in {
		sum += num
	}
	return float64(sum) / float64(length), in[int(math.Floor(float64(length)*0.5))], in[int(math.Floor(float64(length)*0.9))], in[int(math.Floor(float64(length)*0.99))], in[length-1]
}
