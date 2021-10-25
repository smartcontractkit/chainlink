//go:build performance

package performance

import (
	"context"
	"fmt"
	"math/big"
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

// OCRJobMap is a custom map type that holds the record of jobs by the contract instance and the chainlink node
type OCRJobMap map[contracts.OffchainAggregator]map[client.Chainlink]string

// OCRTestOptions contains the parameters for the OCR soak test to be executed
type OCRTestOptions struct {
	TestOptions
	RoundTimeout time.Duration
	AdapterValue int
	TestDuration time.Duration
}

// OCRTest is the implementation of Test that will configure and execute soak test
// of OCR contracts & jobs
type OCRTest struct {
	TestOptions     OCRTestOptions
	ContractOptions contracts.OffchainOptions
	Environment     environment.Environment
	Blockchain      client.BlockchainClient
	Wallets         client.BlockchainWallets
	Deployer        contracts.ContractDeployer

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	contractInstances []contracts.OffchainAggregator
	adapter           environment.ExternalAdapter

	jobMap OCRJobMap
}

// NewOCRTest creates new OCR performance/soak test
func NewOCRTest(
	testOptions OCRTestOptions,
	contractOptions contracts.OffchainOptions,
	env environment.Environment,
	blockchain client.BlockchainClient,
	wallets client.BlockchainWallets,
	deployer contracts.ContractDeployer,
	adapter environment.ExternalAdapter,
) Test {
	return &OCRTest{
		TestOptions:     testOptions,
		ContractOptions: contractOptions,
		Environment:     env,
		Blockchain:      blockchain,
		Wallets:         wallets,
		Deployer:        deployer,
		adapter:         adapter,
		jobMap:          OCRJobMap{},
	}
}

// RecordValues records OCR metrics
func (f *OCRTest) RecordValues(b ginkgo.Benchmarker) error {
	// TODO: collect metrics
	return nil
}

// Setup setups OCR performance/soak test
func (f *OCRTest) Setup() error {
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

func (f *OCRTest) deployContract(c chan<- contracts.OffchainAggregator) error {
	ocrInstance, err := f.Deployer.DeployOffChainAggregator(f.Wallets.Default(), f.ContractOptions)
	if err != nil {
		return err
	}
	err = ocrInstance.Fund(f.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
	if err != nil {
		return err
	}
	if err = ocrInstance.SetConfig(
		f.Wallets.Default(),
		f.chainlinkClients,
		contracts.DefaultOffChainAggregatorConfig(len(f.chainlinkClients)),
	); err != nil {
		return err
	}
	if err = ocrInstance.Fund(f.Wallets.Default(), nil, big.NewFloat(2)); err != nil {
		return err
	}
	c <- ocrInstance
	return nil
}

func (f *OCRTest) deployContracts() error {
	contractChan := make(chan contracts.OffchainAggregator, f.TestOptions.NumberOfContracts)
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

// changeAdapterValue changes adapter value to trigger new round
func (f *OCRTest) changeAdapterValue(roundID int) (int, error) {
	var val int
	if roundID%2 == 0 {
		val = f.TestOptions.AdapterValue * 5
	} else {
		val = f.TestOptions.AdapterValue
	}
	if err := f.adapter.SetVariable(val); err != nil {
		return 0, err
	}
	return val, nil
}

// Run runs OCR performance/soak test
func (f *OCRTest) Run() error {
	if err := f.createChainlinkJobs(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), f.TestOptions.TestDuration)
	defer cancel()
	i := 1
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Test finished")
			return nil
		default:
			log.Info().Int("RoundID", i).Msg("New round")
			val, err := f.changeAdapterValue(i)
			if err != nil {
				return err
			}
			if err := f.waitRoundEnd(i); err != nil {
				return err
			}
			if err := f.checkAllRounds(val); err != nil {
				return err
			}
			i++
		}
	}
}

func (f *OCRTest) waitRoundEnd(roundID int) error {
	for _, ci := range f.contractInstances {
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ci, big.NewInt(int64(roundID)), f.TestOptions.RoundTimeout)
		f.Blockchain.AddHeaderEventSubscription(ci.Address(), ocrRound)
	}
	return f.Blockchain.WaitForEvents()
}

func (f *OCRTest) checkAllRounds(val int) error {
	g := errgroup.Group{}
	log.Info().Msg("Asserting results")
	for _, ci := range f.contractInstances {
		ci := ci
		g.Go(func() error {
			answer, err := ci.GetLatestAnswer(context.Background())
			if err != nil {
				return err
			}
			log.Debug().
				Str("Contract", ci.Address()).
				Int64("Answer", answer.Int64()).
				Msg("Round answer")
			if answer.Int64() != int64(val) {
				return fmt.Errorf("round answer value is different: %d", answer.Int64())
			}
			return nil
		})
	}
	return g.Wait()
}

func (f *OCRTest) createChainlinkJobs() error {
	jobsChan := make(chan OCRJobMap, len(f.chainlinkClients)*len(f.contractInstances))

	bridgeAttrs := make([]client.BridgeTypeAttributes, 0)
	for _, n := range f.chainlinkClients {
		bta := client.BridgeTypeAttributes{
			Name: "variable",
			URL:  fmt.Sprintf("%s/variable", f.adapter.ClusterURL()),
		}
		bridgeAttrs = append(bridgeAttrs, bta)
		if err := n.CreateBridge(&bta); err != nil {
			return err
		}
	}

	for _, contract := range f.contractInstances {
		contract := contract
		g := errgroup.Group{}
		g.Go(func() error {
			return f.createChainlinkJobsPerContract(contract, bridgeAttrs, jobsChan)
		})
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

func (f *OCRTest) createChainlinkJobsPerContract(
	contract contracts.OffchainAggregator,
	bridgesAttrs []client.BridgeTypeAttributes,
	jobsChan chan<- OCRJobMap,
) error {
	// Initialize bootstrap node
	bootstrapNode := f.chainlinkClients[0]
	bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
	if err != nil {
		return err
	}
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
	bootstrapSpec := &client.OCRBootstrapJobSpec{
		ContractAddress: contract.Address(),
		P2PPeerID:       bootstrapP2PId,
		IsBootstrapPeer: true,
	}
	bootstrapJob, err := bootstrapNode.CreateJob(bootstrapSpec)
	if err != nil {
		return err
	}

	jobsChan <- OCRJobMap{contract: map[client.Chainlink]string{bootstrapNode: bootstrapJob.Data.ID}}

	// Send OCR job to other nodes
	g := errgroup.Group{}
	for index := 1; index < len(f.chainlinkClients); index++ {
		index := index
		g.Go(func() error {
			nodeP2PIds, err := f.chainlinkClients[index].ReadP2PKeys()
			if err != nil {
				return err
			}
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := f.chainlinkClients[index].PrimaryEthAddress()
			if err != nil {
				return err
			}
			nodeOCRKeys, err := f.chainlinkClients[index].ReadOCRKeys()
			if err != nil {
				return err
			}
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    contract.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bridgesAttrs[index]),
			}
			jobID, err := f.chainlinkClients[index].CreateJob(ocrSpec)
			if err != nil {
				return err
			}
			jobsChan <- OCRJobMap{contract: map[client.Chainlink]string{f.chainlinkClients[index]: jobID.Data.ID}}
			return nil
		})
	}
	return g.Wait()
}
