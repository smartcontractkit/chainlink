package testsetups

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	integrationactions "github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts/laneconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

var (
	GethResourceProfile = map[string]interface{}{
		"requests": map[string]interface{}{
			"cpu":    "4",
			"memory": "6Gi",
		},
		"limits": map[string]interface{}{
			"cpu":    "4",
			"memory": "6Gi",
		},
	}
)

type NetworkPair struct {
	NetworkA     blockchain.EVMNetwork
	NetworkB     blockchain.EVMNetwork
	ChainClientA blockchain.EVMClient
	ChainClientB blockchain.EVMClient
}

type CCIPTestConfig struct {
	Test                *testing.T
	EnvInput            *testconfig.Common
	TestGroupInput      *testconfig.CCIPTestConfig
	ContractsInput      *testconfig.CCIPContractConfig
	AllNetworks         map[string]blockchain.EVMNetwork
	SelectedNetworks    []blockchain.EVMNetwork
	NetworkPairs        []NetworkPair
	GethResourceProfile map[string]interface{}
}

func (p *CCIPTestConfig) AddPairToNetworkList(networkA, networkB blockchain.EVMNetwork) {
	if p.AllNetworks == nil {
		p.AllNetworks = make(map[string]blockchain.EVMNetwork)
	}
	firstOfPairs := []blockchain.EVMNetwork{networkA}
	secondOfPairs := []blockchain.EVMNetwork{networkB}
	// if no of lanes per pair is greater than 1, copy common contracts from the same network
	// if no of lanes per pair is more than 1, the networks are added into the inputs.AllNetworks with a suffix of -<lane number>
	// for example, if no of lanes per pair is 2, and the network pairs are called "testnetA", "testnetB",
	//	the network will be added as "testnetA-1", testnetA-2","testnetB-1", testnetB-2"
	// to deploy 4 lanes between same network pair "testnetA", "testnetB".
	// lanes - testnetA-1<->testnetB-1, testnetA-1<-->testnetB-2 , testnetA-2<--> testnetB-1, testnetA-2<--> testnetB-2
	if p.TestGroupInput.NoOfRoutersPerPair > 1 {
		firstOfPairs[0].Name = fmt.Sprintf("%s-%d", firstOfPairs[0].Name, 1)
		secondOfPairs[0].Name = fmt.Sprintf("%s-%d", secondOfPairs[0].Name, 1)
		for i := 1; i < p.TestGroupInput.NoOfRoutersPerPair; i++ {
			netsA := networkA
			netsA.Name = fmt.Sprintf("%s-%d", netsA.Name, i+1)
			netsB := networkB
			netsB.Name = fmt.Sprintf("%s-%d", netsB.Name, i+1)
			firstOfPairs = append(firstOfPairs, netsA)
			secondOfPairs = append(secondOfPairs, netsB)
		}
	}

	for i := range firstOfPairs {
		p.AllNetworks[firstOfPairs[i].Name] = firstOfPairs[i]
		p.AllNetworks[secondOfPairs[i].Name] = secondOfPairs[i]
		p.NetworkPairs = append(p.NetworkPairs, NetworkPair{
			NetworkA: firstOfPairs[i],
			NetworkB: secondOfPairs[i],
		})
	}
}

func (p *CCIPTestConfig) SetNetworkPairs(lggr zerolog.Logger) error {
	var allError error

	// if network pairs are provided, then use them
	if p.TestGroupInput.NetworkPairs != nil {
		networkPairs := p.TestGroupInput.NetworkPairs
		networkByChainID := make(map[int64]blockchain.EVMNetwork)
		for _, pair := range networkPairs {
			networkNames := strings.Split(pair, ",")
			if len(networkNames) != 2 {
				allError = multierr.Append(allError, fmt.Errorf("invalid network pair"))
			}
			nets := networks.SetNetworks(networkNames)
			if _, ok := networkByChainID[nets[0].ChainID]; !ok {
				networkByChainID[nets[0].ChainID] = nets[0]
			}
			if _, ok := networkByChainID[nets[1].ChainID]; !ok {
				networkByChainID[nets[1].ChainID] = nets[1]
			}
			p.AddPairToNetworkList(nets[0], nets[1])
		}

		for _, net := range networkByChainID {
			p.SelectedNetworks = append(p.SelectedNetworks, net)
		}
		return allError
	}

	p.SelectedNetworks = p.EnvInput.EVMNetworks()
	if p.TestGroupInput.NoOfNetworks == 0 {
		p.TestGroupInput.NoOfNetworks = len(p.SelectedNetworks)
	}
	// TODO remove this when CTF network timeout is fixed
	for i := range p.SelectedNetworks {
		p.SelectedNetworks[i].Timeout = blockchain.JSONStrDuration{
			Duration: 3 * time.Minute,
		}
	}
	simulated := p.SelectedNetworks[0].Simulated
	for i := 1; i < len(p.SelectedNetworks); i++ {
		if p.SelectedNetworks[i].Simulated != simulated {
			lggr.Fatal().Msg("networks must be of the same type either simulated or real")
		}
	}

	// if the networks are not simulated use the first p.NoOfNetworks networks from the selected networks
	if !simulated && len(p.SelectedNetworks) != p.TestGroupInput.NoOfNetworks {
		if len(p.SelectedNetworks) < p.TestGroupInput.NoOfNetworks {
			allError = multierr.Append(allError, fmt.Errorf("not enough networks provided"))
		} else {
			p.SelectedNetworks = p.SelectedNetworks[:p.TestGroupInput.NoOfNetworks]
		}
	}
	// If provided networks is lesser than the required number of networks
	// and the provided networks are simulated network, create replicas of the provided networks with
	// different chain ids
	if len(p.SelectedNetworks) < p.TestGroupInput.NoOfNetworks {
		if simulated {
			actualNoOfNetworks := len(p.SelectedNetworks)
			n := p.SelectedNetworks[0]
			var chainIDs []int64
			for _, id := range chainselectors.TestChainIds() {
				if id == 2337 {
					continue
				}
				chainIDs = append(chainIDs, int64(id))
			}
			for i := 0; i < p.TestGroupInput.NoOfNetworks-actualNoOfNetworks; i++ {
				chainID := chainIDs[i]
				p.SelectedNetworks = append(p.SelectedNetworks, blockchain.EVMNetwork{
					Name:                      fmt.Sprintf("simulated-non-dev%d", len(p.SelectedNetworks)+1),
					ChainID:                   chainID,
					Simulated:                 true,
					PrivateKeys:               []string{networks.AdditionalSimulatedPvtKeys[i]},
					ChainlinkTransactionLimit: n.ChainlinkTransactionLimit,
					Timeout:                   n.Timeout,
					MinimumConfirmations:      n.MinimumConfirmations,
					GasEstimationBuffer:       n.GasEstimationBuffer + 1000,
					ClientImplementation:      n.ClientImplementation,
					DefaultGasLimit:           n.DefaultGasLimit,
				})
			}
		}
	}

	if p.TestGroupInput.NoOfNetworks > 2 {
		p.FormNetworkPairCombinations()
	} else {
		p.AddPairToNetworkList(p.SelectedNetworks[0], p.SelectedNetworks[1])
	}

	// if the number of lanes is lesser than the number of network pairs, choose a random subset of network pairs
	if p.TestGroupInput.MaxNoOfLanes > 0 && p.TestGroupInput.MaxNoOfLanes < len(p.NetworkPairs) {
		rand.Shuffle(len(p.NetworkPairs), func(i, j int) {
			p.NetworkPairs[i], p.NetworkPairs[j] = p.NetworkPairs[j], p.NetworkPairs[i]
		})
		p.NetworkPairs = p.NetworkPairs[:p.TestGroupInput.MaxNoOfLanes]
	}

	for _, n := range p.NetworkPairs {
		lggr.Info().Str("NetworkA", n.NetworkA.Name).Str("NetworkB", n.NetworkB.Name).Msg("Network Pairs")
	}
	lggr.Info().Int("Pairs", len(p.NetworkPairs)).Msg("No Of Lanes")

	return allError
}

func (p *CCIPTestConfig) FormNetworkPairCombinations() {
	for i := 0; i < p.TestGroupInput.NoOfNetworks; i++ {
		for j := i + 1; j < p.TestGroupInput.NoOfNetworks; j++ {
			p.AddPairToNetworkList(p.SelectedNetworks[i], p.SelectedNetworks[j])
		}
	}
}

func NewCCIPTestConfig(t *testing.T, lggr zerolog.Logger, tType string) *CCIPTestConfig {
	var allError error
	if allError != nil {
		t.Fatal(allError)
	}
	ccipCfg := testconfig.GlobalTestConfig().CCIP.Env
	contractCfg := testconfig.GlobalTestConfig().CCIP.Deployments
	groupCfg, exists := testconfig.GlobalTestConfig().CCIP.Groups[tType]
	if !exists {
		t.Fatalf("group config for %s does not exist", tType)
	}

	p := &CCIPTestConfig{
		Test:                t,
		EnvInput:            ccipCfg,
		ContractsInput:      contractCfg,
		TestGroupInput:      groupCfg,
		GethResourceProfile: GethResourceProfile,
	}

	allError = multierr.Append(allError, p.SetNetworkPairs(lggr))

	if allError != nil {
		t.Fatal(allError)
	}

	return p
}

type BiDirectionalLaneConfig struct {
	NetworkA    blockchain.EVMNetwork
	NetworkB    blockchain.EVMNetwork
	ForwardLane *actions.CCIPLane
	ReverseLane *actions.CCIPLane
}

type CCIPTestSetUpOutputs struct {
	Cfg                      *CCIPTestConfig
	LaneContractsByNetwork   sync.Map
	laneMutex                *sync.Mutex
	Lanes                    []*BiDirectionalLaneConfig
	CommonContractsByNetwork sync.Map
	Reporter                 *testreporters.CCIPTestReporter
	LaneConfigFile           string
	LaneConfig               *laneconfig.Lanes
	TearDown                 func() error
	Env                      *actions.CCIPTestEnv
	Balance                  *actions.BalanceSheet
	BootstrapAdded           *atomic.Bool
	JobAddGrp                *errgroup.Group
}

func (o *CCIPTestSetUpOutputs) AddToLanes(lane *BiDirectionalLaneConfig) {
	o.laneMutex.Lock()
	defer o.laneMutex.Unlock()
	o.Lanes = append(o.Lanes, lane)
}

func (o *CCIPTestSetUpOutputs) ReadLanes() []*BiDirectionalLaneConfig {
	o.laneMutex.Lock()
	defer o.laneMutex.Unlock()
	return o.Lanes
}

func (o *CCIPTestSetUpOutputs) DeployChainContracts(
	chainClient blockchain.EVMClient,
	networkCfg blockchain.EVMNetwork,
	noOfTokens int,
	tokenDeployerFns []blockchain.ContractDeployer,
	lggr zerolog.Logger,
) error {
	var k8Env *environment.Environment
	ccipEnv := o.Env
	if ccipEnv != nil {
		k8Env = ccipEnv.K8Env
	}
	if k8Env != nil && chainClient.NetworkSimulated() {
		networkCfg.URLs = k8Env.URLs[chainClient.GetNetworkConfig().Name]
	}

	chain, err := blockchain.ConcurrentEVMClient(networkCfg, k8Env, chainClient, lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkCfg.Name, err))
	}

	chain.ParallelTransactions(true)
	defer chain.Close()
	ccipCommon, err := actions.DefaultCCIPModule(lggr, chain, pointer.GetBool(o.Cfg.TestGroupInput.ExistingDeployment), pointer.GetBool(o.Cfg.TestGroupInput.MulticallInOneTx))
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create ccip common module for %s: %w", networkCfg.Name, err))
	}

	cfg := o.LaneConfig.ReadLaneConfig(networkCfg.Name)

	err = ccipCommon.DeployContracts(noOfTokens, tokenDeployerFns, cfg)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy common ccip contracts for %s: %w", networkCfg.Name, err))
	}
	o.LaneContractsByNetwork.Store(networkCfg.Name, cfg)
	o.CommonContractsByNetwork.Store(networkCfg.Name, ccipCommon)
	return nil
}

func (o *CCIPTestSetUpOutputs) AddLanesForNetworkPair(
	lggr zerolog.Logger,
	networkA, networkB blockchain.EVMNetwork,
	chainClientA, chainClientB blockchain.EVMClient,
	transferAmounts []*big.Int,
	numOfCommitNodes int,
	commitAndExecOnSameDON, bidirectional bool,
) error {
	var allErrors atomic.Error
	t := o.Cfg.Test
	var k8Env *environment.Environment
	ccipEnv := o.Env
	namespace := o.Cfg.TestGroupInput.ExistingEnv
	if ccipEnv != nil {
		k8Env = ccipEnv.K8Env
		if k8Env != nil {
			namespace = k8Env.Cfg.Namespace
		}
	}
	configureCLNode := !pointer.GetBool(o.Cfg.TestGroupInput.ExistingDeployment)
	setUpFuncs, ctx := errgroup.WithContext(context.Background())

	// Use new set of clients(sourceChainClient,destChainClient)
	// with new header subscriptions(otherwise transactions
	// on one lane will keep on waiting for transactions on other lane for the same network)
	// Currently for simulated network clients(from same network) created with NewEVMClient does not sync nonce
	// ConcurrentEVMClient is a work-around for that.
	sourceChainClientA2B, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA, lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkA.Name, err))
	}

	sourceChainClientA2B.ParallelTransactions(true)

	destChainClientA2B, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB, lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkB.Name, err))
	}
	destChainClientA2B.ParallelTransactions(true)

	ccipLaneA2B := &actions.CCIPLane{
		Test:              t,
		TestEnv:           ccipEnv,
		SourceChain:       sourceChainClientA2B,
		DestChain:         destChainClientA2B,
		SourceNetworkName: actions.NetworkName(networkA.Name),
		DestNetworkName:   actions.NetworkName(networkB.Name),
		ValidationTimeout: o.Cfg.TestGroupInput.PhaseTimeout.Duration(),
		SentReqs:          make(map[common.Hash][]actions.CCIPRequest),
		TotalFee:          big.NewInt(0),
		Balance:           o.Balance,
		Context:           ctx,
	}
	contractsA, ok := o.LaneContractsByNetwork.Load(networkA.Name)
	if !ok {
		return errors.WithStack(fmt.Errorf("failed to load lane contracts for %s", networkA.Name))
	}
	srcCfg := contractsA.(*laneconfig.LaneConfig)
	ccipLaneA2B.SrcNetworkLaneCfg = srcCfg
	contractsB, ok := o.LaneContractsByNetwork.Load(networkB.Name)
	if !ok {
		return errors.WithStack(fmt.Errorf("failed to load lane contracts for %s", networkB.Name))
	}
	destCfg := contractsB.(*laneconfig.LaneConfig)
	ccipLaneA2B.DstNetworkLaneCfg = destCfg

	ccipLaneA2B.Logger = lggr.With().Str("env", namespace).Str("Lane",
		fmt.Sprintf("%s-->%s", ccipLaneA2B.SourceNetworkName, ccipLaneA2B.DestNetworkName)).Logger()
	ccipLaneA2B.Reports = o.Reporter.AddNewLane(fmt.Sprintf("%s To %s",
		networkA.Name, networkB.Name), ccipLaneA2B.Logger)

	bidirectionalLane := &BiDirectionalLaneConfig{
		NetworkA:    networkA,
		NetworkB:    networkB,
		ForwardLane: ccipLaneA2B,
	}

	var ccipLaneB2A *actions.CCIPLane

	if bidirectional {
		sourceChainClientB2A, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB, lggr)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkB.Name, err))
		}
		sourceChainClientB2A.ParallelTransactions(true)

		destChainClientB2A, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA, lggr)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkA.Name, err))
		}
		destChainClientB2A.ParallelTransactions(true)

		ccipLaneB2A = &actions.CCIPLane{
			Test:              t,
			TestEnv:           ccipEnv,
			SourceNetworkName: actions.NetworkName(networkB.Name),
			DestNetworkName:   actions.NetworkName(networkA.Name),
			SourceChain:       sourceChainClientB2A,
			DestChain:         destChainClientB2A,
			ValidationTimeout: o.Cfg.TestGroupInput.PhaseTimeout.Duration(),
			Balance:           o.Balance,
			SentReqs:          make(map[common.Hash][]actions.CCIPRequest),
			TotalFee:          big.NewInt(0),
			Context:           ctx,
			SrcNetworkLaneCfg: ccipLaneA2B.DstNetworkLaneCfg,
			DstNetworkLaneCfg: ccipLaneA2B.SrcNetworkLaneCfg,
		}
		ccipLaneB2A.Logger = lggr.With().Str("env", namespace).Str("Lane",
			fmt.Sprintf("%s-->%s", ccipLaneB2A.SourceNetworkName, ccipLaneB2A.DestNetworkName)).Logger()
		ccipLaneB2A.Reports = o.Reporter.AddNewLane(
			fmt.Sprintf("%s To %s", networkB.Name, networkA.Name), ccipLaneB2A.Logger)
		bidirectionalLane.ReverseLane = ccipLaneB2A
	}
	o.AddToLanes(bidirectionalLane)

	c1, ok := o.CommonContractsByNetwork.Load(networkA.Name)
	var networkACmn *actions.CCIPCommon
	if ok {
		networkACmn = c1.(*actions.CCIPCommon)
	}
	if networkACmn == nil {
		return errors.WithStack(fmt.Errorf("chain contracts for network %s not found", networkA.Name))
	}
	c2, ok := o.CommonContractsByNetwork.Load(networkB.Name)
	var networkBCmn *actions.CCIPCommon
	if ok {
		networkBCmn = c2.(*actions.CCIPCommon)
	}
	if networkBCmn == nil {
		return errors.WithStack(fmt.Errorf("chain contracts for network %s not found", networkB.Name))
	}

	setUpFuncs.Go(func() error {
		lggr.Info().Msgf("Setting up lane %s to %s", networkA.Name, networkB.Name)
		srcConfig, destConfig, err := ccipLaneA2B.DeployNewCCIPLane(numOfCommitNodes, commitAndExecOnSameDON, networkACmn, networkBCmn,
			transferAmounts, o.BootstrapAdded, configureCLNode, o.JobAddGrp)
		if err != nil {
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("deploying lane %s to %s; err - %w", networkA.Name, networkB.Name, errors.WithStack(err))))
			return err
		}
		err = o.LaneConfig.WriteLaneConfig(networkA.Name, srcConfig)
		if err != nil {
			lggr.Error().Err(err).Msgf("error deploying lane %s to %s", networkA.Name, networkB.Name)
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkA.Name, errors.WithStack(err))))
			return err
		}
		err = o.LaneConfig.WriteLaneConfig(networkB.Name, destConfig)
		if err != nil {
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkB.Name, errors.WithStack(err))))
			return err
		}
		lggr.Info().Msgf("done setting up lane %s to %s", networkA.Name, networkB.Name)
		return nil
	})

	setUpFuncs.Go(func() error {
		if bidirectional {
			lggr.Info().Msgf("Setting up lane %s to %s", networkB.Name, networkA.Name)
			srcConfig, destConfig, err := ccipLaneB2A.DeployNewCCIPLane(numOfCommitNodes, commitAndExecOnSameDON, networkBCmn, networkACmn,
				transferAmounts, o.BootstrapAdded, configureCLNode, o.JobAddGrp)
			if err != nil {
				lggr.Error().Err(err).Msgf("error deploying lane %s to %s", networkB.Name, networkA.Name)
				allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("deploying lane %s to %s; err -  %w", networkB.Name, networkA.Name, err)))
				return err
			}

			err = o.LaneConfig.WriteLaneConfig(networkB.Name, srcConfig)
			if err != nil {
				allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkA.Name, err)))
				return err
			}
			err = o.LaneConfig.WriteLaneConfig(networkA.Name, destConfig)
			if err != nil {
				allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkB.Name, err)))
				return err
			}
			lggr.Info().Msgf("done setting up lane %s to %s", networkB.Name, networkA.Name)
			return nil
		}
		return nil
	})

	errs := make(chan error, 1)
	go func() {
		errs <- setUpFuncs.Wait()
	}()

	// wait for either context to get cancelled or all the error-groups to finish execution
	for {
		select {
		case err := <-errs:
			// check if there has been any error while waiting for the error groups
			// to finish execution
			return err
		case <-ctx.Done():
			lggr.Print(ctx.Err())
			return allErrors.Load()
		}
	}
}

func (o *CCIPTestSetUpOutputs) StartEventWatchers() {
	for _, lane := range o.ReadLanes() {
		err := lane.ForwardLane.StartEventWatchers()
		require.NoError(o.Cfg.Test, err)
		if lane.ReverseLane != nil {
			err = lane.ReverseLane.StartEventWatchers()
			require.NoError(o.Cfg.Test, err)
		}
	}
}

func (o *CCIPTestSetUpOutputs) WaitForPriceUpdates(ctx context.Context) {
	t := o.Cfg.Test
	priceUpdateGrp, _ := errgroup.WithContext(ctx)
	priceUpdateTracker := sync.Map{}
	for _, lanes := range o.ReadLanes() {
		lanes := lanes
		waitForUpdate := func(lane actions.CCIPLane) error {
			if id, ok := priceUpdateTracker.Load(lane.Source.Common.PriceRegistry.Address()); ok &&
				id.(uint64) == lane.Source.DestinationChainId {
				return nil
			}
			priceUpdateTracker.Store(lane.Source.Common.PriceRegistry.Address(), lane.Source.DestinationChainId)
			lane.Logger.Info().
				Str("source_chain", lane.Source.Common.ChainClient.GetNetworkName()).
				Uint64("dest_chain", lane.Source.DestinationChainId).
				Str("price_registry", lane.Source.Common.PriceRegistry.Address()).
				Msgf("Waiting for price update")
			err := lane.Source.Common.WatchForPriceUpdates()
			if err != nil {
				return err
			}
			defer func() {
				lane.Logger.Info().
					Str("source_chain", lane.Source.Common.ChainClient.GetNetworkName()).
					Uint64("dest_chain", lane.Source.DestinationChainId).
					Str("price_registry", lane.Source.Common.PriceRegistry.Address()).
					Msg("Stopping price update watch")
				lane.Source.Common.StopWatchingPriceUpdates()
			}()
			err = lane.Source.Common.WaitForPriceUpdates(
				lane.Logger,
				30*time.Minute,
				lane.Source.DestinationChainId,
			)
			if err != nil {
				return errors.Wrapf(err, "waiting for price update failed on lane %s-->%s", lane.SourceNetworkName, lane.DestNetworkName)
			}
			return nil
		}

		priceUpdateGrp.Go(func() error {
			return waitForUpdate(*lanes.ForwardLane)
		})
		if lanes.ReverseLane != nil {
			priceUpdateGrp.Go(func() error {
				return waitForUpdate(*lanes.ReverseLane)
			})
		}
	}

	require.NoError(t, priceUpdateGrp.Wait())
}

// CCIPDefaultTestSetUp sets up the environment for CCIP tests
// if configureCLNode is set as false, it assumes:
// 1. contracts are already deployed on live networks
// 2. CL nodes are set up and configured with existing contracts
// 3. No k8 env deployment is needed
// It reuses already deployed contracts from the addresses provided in ../contracts/ccip/laneconfig/contracts.json
//
// If bidirectional is true it sets up two-way lanes between NetworkA and NetworkB. Same CL nodes are used for both the lanes.
// If bidirectional is false only one way lane is set up.
//
// Returns -
// 1. CCIPLane for NetworkA --> NetworkB
// 2. If bidirectional is true, CCIPLane for NetworkB --> NetworkA
// 3. If configureCLNode is true, the tearDown func to call when environment needs to be destroyed
func CCIPDefaultTestSetUp(
	t *testing.T,
	lggr zerolog.Logger,
	envName string,
	tokenDeployerFns []blockchain.ContractDeployer,
	inputs *CCIPTestConfig,
) *CCIPTestSetUpOutputs {
	var (
		ccipEnv *actions.CCIPTestEnv
		k8Env   *environment.Environment
		err     error
		chains  []blockchain.EVMClient
	)
	filename := fmt.Sprintf("./tmp_%s.json", strings.ReplaceAll(t.Name(), "/", "_"))
	inputs.Test = t
	var transferAmounts []*big.Int
	if inputs.TestGroupInput.MsgType == actions.TokenTransfer {
		for i := 0; i < inputs.TestGroupInput.NoOfTokensInMsg; i++ {
			transferAmounts = append(transferAmounts, big.NewInt(inputs.TestGroupInput.AmountPerToken))
		}
	}
	setUpArgs := &CCIPTestSetUpOutputs{
		Cfg:            inputs,
		Reporter:       testreporters.NewCCIPTestReporter(t, lggr),
		LaneConfigFile: filename,
		Balance:        actions.NewBalanceSheet(),
		BootstrapAdded: atomic.NewBool(false),
		JobAddGrp:      &errgroup.Group{},
		laneMutex:      &sync.Mutex{},
	}

	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	configureCLNode := !pointer.GetBool(inputs.TestGroupInput.ExistingDeployment)
	var deployCL func() error
	var local *test_env.CLClusterTestEnv
	envConfig := &environment.Config{
		NamespacePrefix: envName,
		Test:            t,
	}
	if inputs.EnvInput.TTL != nil {
		envConfig.TTL = inputs.EnvInput.TTL.Duration()
	}
	if inputs.TestGroupInput.TestDuration != nil {
		approxDur := inputs.TestGroupInput.TestDuration.Duration() + 3*time.Hour
		if envConfig.TTL < approxDur {
			envConfig.TTL = approxDur
		}
	}

	if configureCLNode {
		if pointer.GetBool(inputs.TestGroupInput.LocalCluster) {
			local, deployCL = DeployLocalCluster(t, inputs)
			ccipEnv = &actions.CCIPTestEnv{
				LocalCluster: local,
			}
		} else {
			lggr.Info().Msg("Deploying test environment")
			// deploy the env if configureCLNode is true
			k8Env = DeployEnvironments(t, envConfig, inputs)
			ccipEnv = &actions.CCIPTestEnv{K8Env: k8Env}
		}

		ccipEnv.CLNodeWithKeyReady, _ = errgroup.WithContext(parent)
		setUpArgs.Env = ccipEnv
		if ccipEnv.K8Env != nil && ccipEnv.K8Env.WillUseRemoteRunner() {
			return setUpArgs
		}
	} else {
		// if configureCLNode is false, use a placeholder env to create remote runner
		if value, set := os.LookupEnv(config.EnvVarJobImage); set && value != "" {
			k8Env = environment.New(envConfig)
			err = k8Env.Run()
			require.NoErrorf(t, err, "error creating environment remote runner")
			setUpArgs.Env = &actions.CCIPTestEnv{K8Env: k8Env}
			if k8Env.WillUseRemoteRunner() {
				return setUpArgs
			}
		}
	}

	_, err = os.Stat(setUpArgs.LaneConfigFile)
	if err == nil {
		// remove the existing lane config file
		err = os.Remove(setUpArgs.LaneConfigFile)
		require.NoError(t, err, "error while removing existing lane config file - %s", setUpArgs.LaneConfigFile)
	}

	setUpArgs.LaneConfig, err = laneconfig.ReadLanesFromExistingDeployment(setUpArgs.Cfg.ContractsInput.ContractsData())
	require.NoError(t, err)

	if setUpArgs.LaneConfig == nil {
		setUpArgs.LaneConfig = &laneconfig.Lanes{LaneConfigs: make(map[string]*laneconfig.LaneConfig)}
	}

	chainByChainID := make(map[int64]blockchain.EVMClient)
	if pointer.GetBool(inputs.TestGroupInput.LocalCluster) {
		require.NotNil(t, ccipEnv.LocalCluster, "Local cluster shouldn't be nil")
		for _, n := range ccipEnv.LocalCluster.PrivateChain {
			primaryNode := n.GetPrimaryNode()
			require.NotNil(t, primaryNode, "Primary node is nil in PrivateChain interface")
			chainByChainID[primaryNode.GetEVMClient().GetChainID().Int64()] = primaryNode.GetEVMClient()
			chains = append(chains, primaryNode.GetEVMClient())
		}
	} else {
		for _, n := range inputs.SelectedNetworks {
			if _, ok := chainByChainID[n.ChainID]; ok {
				continue
			}
			var ec blockchain.EVMClient
			if k8Env == nil {
				ec, err = blockchain.ConnectEVMClient(n, lggr)
			} else {
				ec, err = blockchain.NewEVMClient(n, k8Env, lggr)
			}
			require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
			chains = append(chains, ec)
			chainByChainID[n.ChainID] = ec
		}
	}
	printStats := func() {
		for k := range setUpArgs.Reporter.LaneStats {
			setUpArgs.Reporter.LaneStats[k].Finalize(k)
		}
	}
	t.Cleanup(func() {
		if configureCLNode {
			if ccipEnv.LocalCluster != nil {
				err := ccipEnv.LocalCluster.Terminate()
				require.NoError(t, err, "Local cluster termination shouldn't fail")
				for k := range setUpArgs.Reporter.LaneStats {
					setUpArgs.Reporter.LaneStats[k].Finalize(k)
				}
				return
			}
			if pointer.GetBool(inputs.TestGroupInput.KeepEnvAlive) {
				printStats()
				return
			}
			lggr.Info().Msg("Tearing down the environment")
			err = integrationactions.TeardownSuite(t, ccipEnv.K8Env, ccipEnv.CLNodes, setUpArgs.Reporter,
				zapcore.ErrorLevel, chains...)
			require.NoError(t, err, "Environment teardown shouldn't fail")
		} else {
			//just print
			printStats()
		}
	})

	if configureCLNode {
		ccipEnv.CLNodeWithKeyReady.Go(func() error {
			if ccipEnv.LocalCluster != nil {
				err = deployCL()
				if err != nil {
					return err
				}
			}
			return ccipEnv.SetUpNodesAndKeys(big.NewFloat(inputs.TestGroupInput.NodeFunding), chains, lggr)
		})
	}

	// if no of lanes per pair is greater than 1, copy common contracts from the same network
	// if no of lanes per pair is more than 1, the networks are added into the inputs.AllNetworks with a suffix of -<lane number>
	// for example, if no of lanes per pair is 2, and the network pairs are called "testnetA", "testnetB",
	//	the network will be added as "testnetA-1", testnetA-2","testnetB-1", testnetB-2"
	// to deploy 2 lanes between same network pair "testnetA", "testnetB".
	// In the following the common contracts will be copied from "testnetA" to "testnetA-1" and "testnetA-2" and
	// from "testnetB" to "testnetB-1" and "testnetB-2"
	for n := range inputs.AllNetworks {
		if setUpArgs.Cfg.TestGroupInput.NoOfRoutersPerPair > 1 {
			regex := regexp.MustCompile(`-(\d+)$`)
			networkNameToReadCfg := regex.ReplaceAllString(n, "")
			reuse := pointer.GetBool(inputs.TestGroupInput.ReuseContracts)
			// if reuse contracts is true, copy common contracts from the same network except the router contract
			setUpArgs.LaneConfig.CopyCommonContracts(
				networkNameToReadCfg, n,
				reuse, inputs.TestGroupInput.MsgType == actions.TokenTransfer)
		}
	}

	// deploy all chain specific common contracts
	chainAddGrp, _ := errgroup.WithContext(parent)
	lggr.Info().Msg("Deploying common contracts")
	for _, net := range inputs.AllNetworks {
		chain := chainByChainID[net.ChainID]
		net := net
		net.HTTPURLs = chain.GetNetworkConfig().HTTPURLs
		net.URLs = chain.GetNetworkConfig().URLs
		chainAddGrp.Go(func() error {
			return setUpArgs.DeployChainContracts(chain, net, inputs.TestGroupInput.NoOfTokensPerChain, tokenDeployerFns, lggr)
		})
	}
	require.NoError(t, chainAddGrp.Wait(), "Deploying common contracts shouldn't fail")

	// deploy all lane specific contracts
	lggr.Info().Msg("Deploying chain specific contracts")
	laneAddGrp, _ := errgroup.WithContext(parent)
	for i, n := range inputs.NetworkPairs {
		i := i
		n := n
		var ok bool
		inputs.NetworkPairs[i].ChainClientA, ok = chainByChainID[n.NetworkA.ChainID]
		require.True(t, ok, "Chain client for chainID %d not found", n.NetworkA.ChainID)
		inputs.NetworkPairs[i].ChainClientB, ok = chainByChainID[n.NetworkB.ChainID]
		require.True(t, ok, "Chain client for chainID %d not found", n.NetworkB.ChainID)

		n.NetworkA.HTTPURLs = inputs.NetworkPairs[i].ChainClientA.GetNetworkConfig().HTTPURLs
		n.NetworkA.URLs = inputs.NetworkPairs[i].ChainClientA.GetNetworkConfig().URLs
		n.NetworkB.HTTPURLs = inputs.NetworkPairs[i].ChainClientB.GetNetworkConfig().HTTPURLs
		n.NetworkB.URLs = inputs.NetworkPairs[i].ChainClientB.GetNetworkConfig().URLs

		laneAddGrp.Go(func() error {
			return setUpArgs.AddLanesForNetworkPair(
				lggr, n.NetworkA, n.NetworkB,
				chainByChainID[n.NetworkA.ChainID], chainByChainID[n.NetworkB.ChainID], transferAmounts,
				inputs.TestGroupInput.NoOfCommitNodes,
				pointer.GetBool(inputs.TestGroupInput.CommitAndExecuteOnSameDON),
				pointer.GetBool(inputs.TestGroupInput.BiDirectionalLane),
			)
		})
	}
	require.NoError(t, laneAddGrp.Wait())
	err = laneconfig.WriteLanesToJSON(setUpArgs.LaneConfigFile, setUpArgs.LaneConfig)
	require.NoError(t, err)
	require.Equal(t, len(setUpArgs.Lanes), len(inputs.NetworkPairs),
		"Number of bi-directional lanes should be equal to number of network pairs")

	if configureCLNode {
		// wait for all jobs to get created
		lggr.Info().Msg("Waiting for jobs to be created")
		require.NoError(t, setUpArgs.JobAddGrp.Wait(), "Creating jobs shouldn't fail")
		// wait for price updates to be available and start event watchers
		setUpArgs.WaitForPriceUpdates(parent)
	}

	// start event watchers for all lanes
	setUpArgs.StartEventWatchers()

	setUpArgs.TearDown = func() error {
		var errs error
		for _, lanes := range setUpArgs.Lanes {
			// if existing deployment is true, don't attempt to pay ccip fees
			err := lanes.ForwardLane.CleanUp(configureCLNode)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			if lanes.ReverseLane != nil {
				// if existing deployment is true, don't attempt to pay ccip fees
				err := lanes.ReverseLane.CleanUp(configureCLNode)
				if err != nil {
					errs = multierr.Append(errs, err)
				}
			}
		}
		return errs
	}
	lggr.Info().Msg("Test setup completed")
	return setUpArgs
}
