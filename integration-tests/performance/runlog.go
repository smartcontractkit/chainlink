//go:build performance

package performance

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"golang.org/x/sync/errgroup"
)

// ConsumerOraclePair consumer and oracle pair
type ConsumerOraclePair struct {
	consumer contracts.APIConsumer
	oracle   contracts.Oracle
	jobUUID  string
}

// RunlogTestOptions contains the parameters for the Runlog soak test to be executed
type RunlogTestOptions struct {
	TestOptions
	AdapterValue int
}

// RunlogTest is the implementation of Test that will configure and execute soak test
// of Runlog contracts & jobs
type RunlogTest struct {
	TestOptions RunlogTestOptions
	Environment environment.Environment
	Blockchain  client.BlockchainClient
	Wallets     client.BlockchainWallets
	Deployer    contracts.ContractDeployer
	Link        contracts.LinkToken

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	contractInstances []*ConsumerOraclePair
	adapter           environment.ExternalAdapter

	testResults *PerfRequestIDTestResults
	jobMap      ContractsNodesJobsMap
}

// NewRunlogTest creates new Runlog performance/soak test
func NewRunlogTest(
	testOptions RunlogTestOptions,
	env environment.Environment,
	link contracts.LinkToken,
	blockchain client.BlockchainClient,
	wallets client.BlockchainWallets,
	deployer contracts.ContractDeployer,
	adapter environment.ExternalAdapter,
) Test {
	return &RunlogTest{
		TestOptions: testOptions,
		Environment: env,
		Link:        link,
		Blockchain:  blockchain,
		Wallets:     wallets,
		Deployer:    deployer,
		adapter:     adapter,
		testResults: NewPerfRequestIDTestResults(),
		jobMap:      ContractsNodesJobsMap{},
	}
}

// RecordValues records Runlog metrics
func (f *RunlogTest) RecordValues(b ginkgo.Benchmarker) error {
	// can't estimate perf metrics in soak mode
	if f.TestOptions.NumberOfRounds == 0 {
		return nil
	}
	actions.SetChainlinkAPIPageSize(f.chainlinkClients, f.TestOptions.NumberOfRounds*f.TestOptions.NumberOfContracts)
	if err := f.testResults.setResultStartTimes(f.chainlinkClients, f.jobMap); err != nil {
		return err
	}
	return f.testResults.calculateLatencies(b)
}

// Setup setups Runlog performance/soak test
func (f *RunlogTest) Setup() error {
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

func (f *RunlogTest) deployContract(c chan<- *ConsumerOraclePair) error {
	oracle, err := f.Deployer.DeployOracle(f.Wallets.Default(), f.Link.Address())
	if err != nil {
		return err
	}
	if err = oracle.SetFulfillmentPermission(f.Wallets.Default(), f.nodeAddresses[0].Hex(), true); err != nil {
		return err
	}
	consumer, err := f.Deployer.DeployAPIConsumer(f.Wallets.Default(), f.Link.Address())
	if err != nil {
		return err
	}
	err = consumer.Fund(f.Wallets.Default(), nil, big.NewFloat(20000))
	if err != nil {
		return err
	}
	c <- &ConsumerOraclePair{consumer: consumer, oracle: oracle}
	return nil
}

func (f *RunlogTest) deployContracts() error {
	contractChan := make(chan *ConsumerOraclePair, f.TestOptions.NumberOfContracts)
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

func (f *RunlogTest) requestData() error {
	g := errgroup.Group{}
	for _, p := range f.contractInstances {
		p := p
		g.Go(func() error {
			jobUUIDReplaces := strings.Replace(p.jobUUID, "-", "", 4)
			var jobID [32]byte
			copy(jobID[:], jobUUIDReplaces)
			return p.consumer.CreateRequestTo(
				f.Wallets.Default(),
				p.oracle.Address(),
				jobID,
				big.NewInt(1e18),
				fmt.Sprintf("%s/five", f.adapter.ClusterURL()),
				"data,result",
				big.NewInt(100),
			)
		})
	}
	return g.Wait()
}

// Run runs Runlog performance/soak test
func (f *RunlogTest) Run() error {
	if err := f.createChainlinkJobs(); err != nil {
		return err
	}
	var ctx context.Context
	var testCtxCancel context.CancelFunc
	if f.TestOptions.TestDuration.Seconds() > 0 {
		ctx, testCtxCancel = context.WithTimeout(context.Background(), f.TestOptions.TestDuration)
	} else {
		ctx, testCtxCancel = context.WithCancel(context.Background())
	}
	defer testCtxCancel()
	cancelPerfEvents := f.watchPerfEvents()
	currentRound := 0
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Test finished")
			time.Sleep(f.TestOptions.GracefulStopDuration)
			cancelPerfEvents()
			return nil
		default:
			log.Info().Int("RoundID", currentRound).Msg("New round")
			if err := f.requestData(); err != nil {
				return err
			}
			if err := f.waitRoundEnd(currentRound + 1); err != nil {
				return err
			}
			if f.TestOptions.NumberOfRounds != 0 && currentRound >= f.TestOptions.NumberOfRounds {
				log.Info().Msg("Final round is reached")
				testCtxCancel()
			}
			currentRound++
		}
	}
}

func (f *RunlogTest) watchPerfEvents() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan *contracts.PerfEvent)
		g := errgroup.Group{}
		for _, p := range f.contractInstances {
			p := p
			g.Go(func() error {
				return p.consumer.WatchPerfEvents(context.Background(), ch)
			})
		}
		for {
			select {
			case event := <-ch:
				rqID := common.Bytes2Hex(event.RequestID[:])
				r := f.testResults.Get(rqID)
				loc, _ := time.LoadLocation("UTC")
				r.EndTime = time.Unix(event.BlockTimestamp.Int64(), 0).In(loc)
				log.Debug().
					Int64("Round", event.Round.Int64()).
					Str("RequestID", rqID).
					Time("EndTime", r.EndTime).
					Msg("Perf event received")
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel
}

func (f *RunlogTest) waitRoundEnd(roundID int) error {
	for _, p := range f.contractInstances {
		rc := contracts.NewRunlogRoundConfirmer(p.consumer, big.NewInt(int64(roundID)), f.TestOptions.RoundTimeout)
		f.Blockchain.AddHeaderEventSubscription(p.consumer.Address(), rc)
	}
	return f.Blockchain.WaitForEvents()
}

func (f *RunlogTest) createChainlinkJobs() error {
	jobsChan := make(chan ContractsNodesJobsMap, len(f.contractInstances))
	g := NewLimitErrGroup(30)

	bta := client.BridgeTypeAttributes{
		Name: "five",
		URL:  fmt.Sprintf("%s/five", f.adapter.ClusterURL()),
	}
	if err := f.chainlinkClients[0].CreateBridge(&bta); err != nil {
		return err
	}
	os := &client.DirectRequestTxPipelineSpec{
		BridgeTypeAttributes: bta,
		DataPath:             "data,result",
	}
	ost, err := os.String()
	if err != nil {
		return err
	}

	for _, p := range f.contractInstances {
		p := p
		g.Go(func() error {
			jobUUID := uuid.NewV4()
			p.jobUUID = jobUUID.String()
			job, err := f.chainlinkClients[0].CreateJob(&client.DirectRequestJobSpec{
				Name:              fmt.Sprintf("direct_request_%s", p.jobUUID),
				ContractAddress:   p.oracle.Address(),
				ExternalJobID:     jobUUID.String(),
				ObservationSource: ost,
			})
			if err != nil {
				return err
			}
			jobsChan <- ContractsNodesJobsMap{p.consumer: map[client.Chainlink]NodeData{f.chainlinkClients[0]: RunlogNodeData{JobID: job.Data.ID}}}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	close(jobsChan)
	f.jobMap.FromJobsChan(jobsChan)
	return nil
}
