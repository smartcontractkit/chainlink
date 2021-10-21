package performance

import (
	"context"
	"fmt"
	"math/big"
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

// KeeperJobMap is a custom map type that holds the record of jobs by the contract instance and the chainlink node
type KeeperJobMap map[contracts.KeeperConsumer]map[client.Chainlink]string

// KeeperTestOptions contains the parameters for the Keeper performance/soak test to be executed
type KeeperTestOptions struct {
	TestOptions
	PaymentPremiumPPB     uint32
	RegistryCheckGasLimit uint32
	BlockCountPerTurn     *big.Int
	StalenessSeconds      *big.Int
	GasCeilingMultiplier  uint16
	RoundTimeout          time.Duration
	TestDuration          time.Duration
}

// KeeperTest is the implementation of Test that will configure and execute soak test
// of Keeper contracts & jobs
type KeeperTest struct {
	TestOptions KeeperTestOptions
	Environment environment.Environment
	Blockchain  client.BlockchainClient
	Wallets     client.BlockchainWallets
	Deployer    contracts.ContractDeployer
	// common contracts
	Link      contracts.LinkToken
	GasFeed   contracts.MockGasFeed
	LinkFeed  contracts.MockETHLINKFeed
	Registry  contracts.KeeperRegistry
	Registrar contracts.UpkeepRegistrar

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	contractInstances []contracts.KeeperConsumer
	adapter           environment.ExternalAdapter

	jobMap KeeperJobMap
}

// NewKeeperTest creates new Keeper performance/soak test
func NewKeeperTest(
	testOptions KeeperTestOptions,
	env environment.Environment,
	blockchain client.BlockchainClient,
	wallets client.BlockchainWallets,
	deployer contracts.ContractDeployer,
	adapter environment.ExternalAdapter,
	link contracts.LinkToken,
) Test {
	return &KeeperTest{
		TestOptions: testOptions,
		Environment: env,
		Blockchain:  blockchain,
		Wallets:     wallets,
		Deployer:    deployer,
		adapter:     adapter,
		jobMap:      KeeperJobMap{},
		Link:        link,
	}
}

// RecordValues records Keeper metrics
func (f *KeeperTest) RecordValues(b ginkgo.Benchmarker) error {
	// TODO: collect metrics
	return nil
}

// Setup setups Keeper performance/soak test
func (f *KeeperTest) Setup() error {
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
	if err := f.deployMockFeeds(); err != nil {
		return err
	}
	if err := f.deployRegistry(); err != nil {
		return err
	}
	if err := f.deployRegistrar(); err != nil {
		return err
	}
	if err := f.deployConsumers(); err != nil {
		return err
	}
	if err := f.setKeepers(); err != nil {
		return err
	}
	return f.fundRegistrarBalances()
}

func (f *KeeperTest) deployConsumerContract(c chan<- contracts.KeeperConsumer) error {
	consumer, err := f.Deployer.DeployKeeperConsumer(f.Wallets.Default(), big.NewInt(5))
	if err != nil {
		return err
	}
	if err = consumer.Fund(f.Wallets.Default(), big.NewFloat(0), big.NewFloat(1)); err != nil {
		return err
	}
	c <- consumer
	return nil
}

// deployMockFeeds deploys mock ETH/Link & Gas price feeds
func (f *KeeperTest) deployMockFeeds() error {
	var err error
	if f.LinkFeed, err = f.Deployer.DeployMockETHLINKFeed(f.Wallets.Default(), big.NewInt(2e18)); err != nil {
		return err
	}
	if f.GasFeed, err = f.Deployer.DeployMockGasFeed(f.Wallets.Default(), big.NewInt(2e11)); err != nil {
		return err
	}
	return f.Blockchain.WaitForEvents()
}

// deployRegistrar deploys contract for registering upkeeps
func (f *KeeperTest) deployRegistrar() error {
	var err error
	f.Registrar, err = f.Deployer.DeployUpkeepRegistrationRequests(
		f.Wallets.Default(),
		f.Link.Address(),
		big.NewInt(0),
	)
	if err != nil {
		return err
	}
	if err = f.Registry.SetRegistrar(f.Wallets.Default(), f.Registrar.Address()); err != nil {
		return err
	}
	if err = f.Registrar.SetRegistrarConfig(
		f.Wallets.Default(),
		true,
		uint32(999),
		uint16(999),
		f.Registry.Address(),
		big.NewInt(0),
	); err != nil {
		return err
	}
	return f.Blockchain.WaitForEvents()
}

// deployRegistry deploys keeper registry
func (f *KeeperTest) deployRegistry() error {
	var err error
	f.Registry, err = f.Deployer.DeployKeeperRegistry(
		f.Wallets.Default(),
		&contracts.KeeperRegistryOpts{
			LinkAddr:             f.Link.Address(),
			ETHFeedAddr:          f.LinkFeed.Address(),
			GasFeedAddr:          f.GasFeed.Address(),
			PaymentPremiumPPB:    f.TestOptions.PaymentPremiumPPB,
			BlockCountPerTurn:    f.TestOptions.BlockCountPerTurn,
			CheckGasLimit:        f.TestOptions.RegistryCheckGasLimit,
			StalenessSeconds:     f.TestOptions.StalenessSeconds,
			GasCeilingMultiplier: f.TestOptions.GasCeilingMultiplier,
			FallbackGasPrice:     big.NewInt(2e11),
			FallbackLinkPrice:    big.NewInt(2e18),
		},
	)
	if err != nil {
		return err
	}
	return f.Registry.Fund(f.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
}

// setKeepers sets keepers, all keepers are "payees" too
func (f *KeeperTest) setKeepers() error {
	nodeAddresses := make([]string, 0)
	for _, node := range f.chainlinkClients {
		nodeAddr, err := node.PrimaryEthAddress()
		if err != nil {
			return err
		}
		nodeAddresses = append(nodeAddresses, nodeAddr)
	}
	if err := f.Registry.SetKeepers(f.Wallets.Default(), nodeAddresses, nodeAddresses); err != nil {
		return err
	}
	return f.Blockchain.WaitForEvents()
}

//fundRegistrarBalances funds registrar balances so payees can be charged
func (f *KeeperTest) fundRegistrarBalances() error {
	g := errgroup.Group{}
	for _, consumer := range f.contractInstances {
		consumer := consumer
		g.Go(func() error {
			req, err := f.Registrar.EncodeRegisterRequest(
				fmt.Sprintf("upkeep_perf_%s", uuid.NewV4().String()),
				[]byte("0x1234"),
				consumer.Address(),
				f.TestOptions.RegistryCheckGasLimit,
				f.Wallets.Default().Address(),
				[]byte("0x"),
				big.NewInt(9e18),
				0,
			)
			if err != nil {
				return err
			}
			return f.Link.TransferAndCall(f.Wallets.Default(), f.Registrar.Address(), big.NewInt(9e18), req)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return f.Blockchain.WaitForEvents()
}

// deployConsumers deploys consumers on which upkeeps will be performed
func (f *KeeperTest) deployConsumers() error {
	contractChan := make(chan contracts.KeeperConsumer, f.TestOptions.NumberOfContracts)
	g := errgroup.Group{}

	for i := 0; i < f.TestOptions.NumberOfContracts; i++ {
		g.Go(func() error {
			return f.deployConsumerContract(contractChan)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	if err := f.Blockchain.WaitForEvents(); err != nil {
		return err
	}
	close(contractChan)
	for contract := range contractChan {
		f.contractInstances = append(f.contractInstances, contract)
	}
	return nil
}

// Run runs Keeper performance/soak test
func (f *KeeperTest) Run() error {
	if err := f.createKeeperJobs(); err != nil {
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
			requiredUpkeeps := i * len(f.chainlinkClients)
			log.Info().Int("RoundID", i).Int("Required upkeeps", requiredUpkeeps).Msg("New round")
			if err := f.waitForUpkeeps(requiredUpkeeps); err != nil {
				return err
			}
			i++
		}
	}
}

// waitForUpkeeps await that every consumer upkeep was called and counter inside contract is incremented for that round
func (f *KeeperTest) waitForUpkeeps(upkeeps int) error {
	for _, consumer := range f.contractInstances {
		upkeepRound := contracts.NewKeeperConsumerRoundConfirmer(consumer, upkeeps, f.TestOptions.RoundTimeout)
		f.Blockchain.AddHeaderEventSubscription(consumer.Address(), upkeepRound)
	}
	return f.Blockchain.WaitForEvents()
}

// createKeeperJobs creates keeper jobs
func (f *KeeperTest) createKeeperJobs() error {
	for _, node := range f.chainlinkClients {
		nodeAddr, err := node.PrimaryEthAddress()
		if err != nil {
			return err
		}
		_, err = node.CreateJob(&client.KeeperJobSpec{
			Name:            "keeper",
			ContractAddress: f.Registry.Address(),
			FromAddress:     nodeAddr,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
