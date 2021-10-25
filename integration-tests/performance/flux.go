//go:build performance

package performance

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"golang.org/x/sync/errgroup"
)

// FluxTestOptions contains the parameters for the performance test to be executed
type FluxTestOptions struct {
	TestOptions
	RequiredSubmissions      int
	RestartDelayRounds       int
	NodePollTimePeriod       time.Duration
	MeasureLatenciesPerRound bool
}

// FluxTest is the implementation of Test that will configure and execute a performance test
// of FluxAggregator contracts & jobs
type FluxTest struct {
	TestOptions     FluxTestOptions
	ContractOptions contracts.FluxAggregatorOptions
	Environment     environment.Environment
	Blockchain      client.BlockchainClient
	Wallets         client.BlockchainWallets
	Deployer        contracts.ContractDeployer
	Prometheus      *client.Prometheus

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	contractInstances []contracts.FluxAggregator
	adapter           environment.ExternalAdapter
	submissionCount   int

	testResults PerfRoundTestResults
	jobMap      FluxJobMap

	headerTimestampCache map[uint64]time.Time
}

// NewFluxTest returns an instantiated instance of FluxTest
func NewFluxTest(
	testOptions FluxTestOptions,
	contractOptions contracts.FluxAggregatorOptions,
	env environment.Environment,
	blockchain client.BlockchainClient,
	wallets client.BlockchainWallets,
	deployer contracts.ContractDeployer,
	prom *client.Prometheus,
) Test {
	return &FluxTest{
		TestOptions:          testOptions,
		ContractOptions:      contractOptions,
		Environment:          env,
		Blockchain:           blockchain,
		Wallets:              wallets,
		Deployer:             deployer,
		Prometheus:           prom,
		testResults:          NewPerfTestResults(),
		jobMap:               FluxJobMap{},
		headerTimestampCache: map[uint64]time.Time{},
	}
}

// Setup will deploy all the contracts and create all the Chainlink jobs for the performance test
func (f *FluxTest) Setup() error {
	chainlinkClients, err := environment.GetChainlinkClients(f.Environment)
	if err != nil {
		return err
	}
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkClients)
	if err != nil {
		return err
	}
	adapter, err := environment.GetExternalAdapter(f.Environment)
	if err != nil {
		return err
	}
	f.chainlinkClients = chainlinkClients
	f.nodeAddresses = nodeAddresses
	f.adapter = adapter

	return f.deployContracts()
}

// Run will start the performance test by creating all the jobs within all Chainlink nodes, subscribing to events
// and ensuring all responses are received.
func (f *FluxTest) Run() error {
	g := errgroup.Group{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chainlinkMap, err := f.mapChainlinkByETHAddress()
	if err != nil {
		return err
	}
	g.Go(func() error {
		return f.watchSubmissions(ctx, chainlinkMap)
	})

	if err := f.createChainlinkJobs(); err != nil {
		return err
	}
	for i := 1; i <= f.TestOptions.NumberOfRounds; i++ {
		if err := f.adapter.SetVariable(i); err != nil {
			return err
		}
		if err := f.waitForAllContractRounds(big.NewInt(int64(i))); err != nil {
			return err
		}
	}
	if err := f.waitForAllSubmissions(cancel); err != nil {
		return err
	}
	return g.Wait()
}

// RecordValues will query all of the latencies of the FluxAggregator rounds and then record them within the
// test runner
func (f *FluxTest) RecordValues(b ginkgo.Benchmarker) error {
	// for each submission we can have 2 runs,
	// one is triggered by poll timer, another when node sees log with another node submission
	actions.SetChainlinkAPIPageSize(f.chainlinkClients, int(f.TestOptions.NumberOfRounds)*2)
	if err := f.setResultStartTimes(); err != nil {
		return err
	}
	return f.calculateLatencies(b)
}

func (f *FluxTest) deployContracts() error {
	contractChan := make(chan contracts.FluxAggregator, f.TestOptions.NumberOfContracts)
	g := errgroup.Group{}

	for i := 0; i < f.TestOptions.NumberOfContracts; i++ {
		g.Go(func() error {
			return f.deployContract(contractChan)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	close(contractChan)
	for contract := range contractChan {
		f.contractInstances = append(f.contractInstances, contract)
	}
	return f.Blockchain.WaitForEvents()
}

func (f *FluxTest) deployContract(contractChan chan<- contracts.FluxAggregator) error {
	fluxInstance, err := f.Deployer.DeployFluxAggregatorContract(f.Wallets.Default(), f.ContractOptions)
	if err != nil {
		return err
	}
	err = fluxInstance.Fund(f.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = fluxInstance.UpdateAvailableFunds(context.Background(), f.Wallets.Default())
	if err != nil {
		return err
	}

	if err := fluxInstance.SetOracles(
		f.Wallets.Default(),
		contracts.FluxAggregatorSetOraclesOptions{
			AddList:            f.nodeAddresses,
			RemoveList:         []common.Address{},
			AdminList:          f.nodeAddresses,
			MinSubmissions:     uint32(f.TestOptions.RequiredSubmissions),
			MaxSubmissions:     uint32(f.TestOptions.RequiredSubmissions),
			RestartDelayRounds: uint32(f.TestOptions.RestartDelayRounds),
		},
	); err != nil {
		return err
	}
	contractChan <- fluxInstance
	return nil
}

func (f *FluxTest) createChainlinkJobs() error {
	jobsChan := make(chan FluxJobMap, len(f.chainlinkClients)*len(f.contractInstances))

	for _, contract := range f.contractInstances {
		contract := contract
		g := errgroup.Group{}
		for _, node := range f.chainlinkClients {
			node := node
			g.Go(func() error {
				return f.createChainlinkJob(contract, node, jobsChan)
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}
	close(jobsChan)

	for jobMap := range jobsChan {
		for contractAddr, m := range jobMap {
			if _, ok := f.jobMap[contractAddr]; !ok {
				f.jobMap[contractAddr] = map[client.Chainlink]string{}
			}
			for k, v := range m {
				f.jobMap[contractAddr][k] = v
			}
		}
	}
	return nil
}

func (f *FluxTest) createChainlinkJob(
	contract contracts.FluxAggregator,
	chainlink client.Chainlink,
	jobsChan chan<- FluxJobMap,
) error {
	job, err := chainlink.CreateJob(&client.FluxMonitorJobSpec{
		Name:              contract.Address(),
		ContractAddress:   contract.Address(),
		PollTimerPeriod:   f.TestOptions.NodePollTimePeriod,
		IdleTimerDisabled: true,
		PollTimerDisabled: false,
		ObservationSource: client.ObservationSourceSpecHTTP(fmt.Sprintf("%s/variable", f.adapter.ClusterURL())),
	})
	if err != nil {
		return err
	}
	jobsChan <- FluxJobMap{contract: map[client.Chainlink]string{chainlink: job.Data.ID}}
	return nil
}

func (f *FluxTest) waitForAllContractRounds(roundID *big.Int) error {
	for _, contract := range f.contractInstances {
		contract := contract
		roundConfirmer := contracts.NewFluxAggregatorRoundConfirmer(contract, roundID, time.Minute*5)
		f.Blockchain.AddHeaderEventSubscription(contract.Address(), roundConfirmer)
	}
	return f.Blockchain.WaitForEvents()
}

func (f *FluxTest) setResultStartTimes() error {
	g := errgroup.Group{}
	for contract := range f.jobMap {
		contract := contract
		g.Go(func() error {
			return f.setResultStartTimeByContract(contract)
		})
	}
	return g.Wait()
}

func (f *FluxTest) setResultStartTimeByContract(contract contracts.FluxAggregator) error {
	for _, chainlink := range f.chainlinkClients {
		chainlink := chainlink

		startTimes, err := f.getJobStartTimes(chainlink, f.jobMap[contract][chainlink])
		if err != nil {
			return err
		}
		log.Debug().Str("Contract", contract.Address()).Interface("Start times", startTimes).Msg("Earliest start times")
		for roundID, startTime := range startTimes {
			result := f.testResults.Get(roundID, contract, chainlink)
			result.StartTime = startTime
		}
	}
	return nil
}

// takeEarliestRunsByAnswer takes earliest node job run for an answer
func (f *FluxTest) takeEarliestRunsByAnswer(runs []client.RunsResponseData) []client.RunsResponseData {
	deduplicated := make([]client.RunsResponseData, 0)
	sort.SliceStable(runs, func(i, j int) bool {
		return runs[i].Attributes.CreatedAt.Before(runs[j].Attributes.CreatedAt)
	})
	seen := make(map[int]bool)
	for _, data := range runs {
		answer := data.Attributes.Inputs.Parse
		if answer > f.TestOptions.NumberOfRounds {
			continue
		}
		if _, ok := seen[answer]; !ok {
			deduplicated = append(deduplicated, data)
		}
		seen[answer] = true
	}
	return deduplicated
}

func (f *FluxTest) getJobStartTimes(chainlink client.Chainlink, jobID string) (map[int64]time.Time, error) {
	jobRuns, err := chainlink.ReadRunsByJob(jobID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Str("Node", chainlink.URL()).
		Int("Runs", len(jobRuns.Data)).
		Interface("Data", jobRuns.Data).
		Msg("Total runs")
	earliestRuns := f.takeEarliestRunsByAnswer(jobRuns.Data)
	var runsStartTimes []time.Time
	for _, run := range earliestRuns {
		runsStartTimes = append(runsStartTimes, run.Attributes.CreatedAt)
	}

	// Place into a map to preserve order
	startTimesMap := map[int64]time.Time{}
	for i, startTime := range runsStartTimes {
		startTimesMap[int64(i+1)] = startTime
	}
	return startTimesMap, nil
}

func (f *FluxTest) watchSubmissions(ctx context.Context, chainlinkMap map[string]client.Chainlink) error {
	eventChan := make(chan *contracts.SubmissionEvent)
	g := errgroup.Group{}

	for _, contract := range f.contractInstances {
		contract := contract
		g.Go(func() error {
			return contract.WatchSubmissionReceived(ctx, eventChan)
		})
	}

	for {
		select {
		case event := <-eventChan:
			if event.Round > uint32(f.TestOptions.NumberOfRounds) {
				continue
			}
			log.Debug().
				Str("Contract Address", event.Contract.String()).
				Uint32("Round ID", event.Round).
				Str("Oracle", event.Oracle.String()).
				Msg("Received Submission")
			f.submissionCount++
			if err := f.setResultEndTime(chainlinkMap, event); err != nil {
				return err
			}
		case <-ctx.Done():
			return g.Wait()
		}
	}
}

func (f *FluxTest) mapChainlinkByETHAddress() (map[string]client.Chainlink, error) {
	chainlinkMap := map[string]client.Chainlink{}
	for _, chainlink := range f.chainlinkClients {
		primaryETHAddress, err := chainlink.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		chainlinkMap[strings.ToLower(primaryETHAddress)] = chainlink
	}
	return chainlinkMap, nil
}

func (f *FluxTest) setResultEndTime(
	chainlinkMap map[string]client.Chainlink,
	submission *contracts.SubmissionEvent,
) error {
	var contract contracts.FluxAggregator
	for _, c := range f.contractInstances {
		if c.Address() == submission.Contract.String() {
			contract = c
			break
		}
	}
	if contract == nil {
		return fmt.Errorf("contract with address of %s isn't stored within the test", submission.Contract.String())
	}

	if _, ok := f.headerTimestampCache[submission.BlockNumber]; !ok {
		blockTime, err := f.Blockchain.HeaderTimestampByNumber(
			context.Background(),
			big.NewInt(int64(submission.BlockNumber)),
		)
		if err != nil {
			return err
		}
		loc, _ := time.LoadLocation("UTC")
		f.headerTimestampCache[submission.BlockNumber] = time.Unix(int64(blockTime), 0).In(loc)
	}

	roundID := int64(submission.Round)
	testResult := f.testResults.Get(
		roundID,
		contract,
		chainlinkMap[strings.ToLower(submission.Oracle.String())],
	)
	testResult.EndTime = f.headerTimestampCache[submission.BlockNumber]
	log.Debug().Str("Contract", contract.Address()).Str("Oracle", submission.Oracle.Hex()).Time("End time", testResult.EndTime).Send()
	return nil
}

func (f *FluxTest) waitForAllSubmissions(ctxCancel context.CancelFunc) error {
	defer ctxCancel()

	ticker := time.NewTicker(time.Millisecond * 100)
	timeout := time.NewTimer(time.Second * 5)

	expectedSubmissionCount := len(f.contractInstances) * len(f.chainlinkClients) * int(f.TestOptions.NumberOfRounds)

	for {
		select {
		case <-ticker.C:
			if f.submissionCount >= expectedSubmissionCount {
				return nil
			}
		case <-timeout.C:
			return fmt.Errorf(
				"timeout while waiting for all submissions to be received, %d < %d",
				f.submissionCount,
				expectedSubmissionCount,
			)
		}
	}
}

func (f *FluxTest) calculateLatencies(b ginkgo.Benchmarker) error {
	var latencies []time.Duration

	for roundID, testResults := range f.testResults.GetAll() {
		if roundID > int64(f.TestOptions.NumberOfRounds) {
			continue
		}
		log.Info().Int64("Round ID", roundID).Msg("Calculating latencies for round")

		for contract, contractResults := range testResults {
			for node, nodeResults := range contractResults {
				if nodeResults.StartTime.IsZero() {
					log.Warn().
						Int64("Round ID", roundID).
						Str("Contract", contract.Address()).
						Str("Node", node.URL()).
						Msg("Start time zero")
				}
				if nodeResults.EndTime.IsZero() {
					log.Warn().
						Int64("Round ID", roundID).
						Str("Contract", contract.Address()).
						Str("Node", node.URL()).
						Msg("End time zero")
				}
				latency := nodeResults.EndTime.Sub(nodeResults.StartTime)
				if latency.Seconds() < 0 {
					log.Warn().
						Time("Start", nodeResults.StartTime).
						Time("End", nodeResults.EndTime).
						Int64("Round ID", roundID).
						Str("Contract", contract.Address()).
						Str("Node", node.URL()).
						Msg("Latency below zero")
				}
				latencies = append(latencies, latency)
			}
		}
		if f.TestOptions.MeasureLatenciesPerRound {
			if err := recordResults(
				b,
				fmt.Sprintf("Round_%d_Submission_Latency", roundID),
				latencies,
			); err != nil {
				return err
			}
		}
		if err := recordResults(b, "Submission_Latency", latencies); err != nil {
			return err
		}
	}
	return nil
}

// FluxJobMap is a custom map type that holds the record of jobs by the contract instance and the chainlink node
type FluxJobMap map[contracts.FluxAggregator]map[client.Chainlink]string
