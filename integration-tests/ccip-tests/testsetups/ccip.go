package testsetups

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockserver_cfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	integrationactions "github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts/laneconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/types/config/node"
	ccipnode "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

const (
	Load                            string = "Load"
	Chaos                           string = "Chaos"
	Smoke                           string = "Smoke"
	DefaultTTLForLongTests                 = 5 * time.Hour
	DefaultNoOfNetworks             int    = 2
	DefaultLoadRPS                  int64  = 2
	DefaultLoadTimeOut                     = 30 * time.Minute
	DefaultPhaseTimeoutForLongTests        = 50 * time.Minute
	DefaultPhaseTimeout                    = 10 * time.Minute
	DefaultTestDuration                    = 10 * time.Minute
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
	DONResourceProfile = map[string]interface{}{
		"requests": map[string]interface{}{
			"cpu":    "4",
			"memory": "8Gi",
		},
		"limits": map[string]interface{}{
			"cpu":    "4",
			"memory": "8Gi",
		},
	}
	DONDBResourceProfile = map[string]interface{}{
		"stateful": true,
		"capacity": "10Gi",
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2",
				"memory": "4Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2",
				"memory": "4Gi",
			},
		},
	}
	NodeFundingForLoad = big.NewFloat(20)
	DefaultNodeFunding = big.NewFloat(1)
)

type NetworkPair struct {
	NetworkA     blockchain.EVMNetwork
	NetworkB     blockchain.EVMNetwork
	ChainClientA blockchain.EVMClient
	ChainClientB blockchain.EVMClient
}

type CCIPTestConfig struct {
	Test                    *testing.T
	EnvTTL                  time.Duration
	MsgType                 string
	PhaseTimeout            time.Duration
	TestDuration            time.Duration
	LocalCluster            bool
	ExistingDeployment      bool
	ExistingEnv             string
	ReuseContracts          bool
	SequentialLaneAddition  bool
	NodeFunding             *big.Float
	Load                    *CCIPLoadInput
	AllNetworks             []blockchain.EVMNetwork
	NetworkPairs            []NetworkPair
	NoOfNetworks            int
	GethResourceProfile     map[string]interface{}
	CLNodeResourceProfile   map[string]interface{}
	CLNodeDBResourceProfile map[string]interface{}
}

type CCIPLoadInput struct {
	RequestPerUnitTime         []int64
	LoadTimeOut                time.Duration
	TimeUnit                   time.Duration
	StepDuration               []time.Duration
	WaitBetweenChaosDuringLoad time.Duration
}

func (p *CCIPTestConfig) setLoadInputs() {
	var allError error
	p.Load = &CCIPLoadInput{
		RequestPerUnitTime:         []int64{DefaultLoadRPS},
		LoadTimeOut:                DefaultLoadTimeOut,
		TimeUnit:                   time.Second,
		StepDuration:               []time.Duration{p.TestDuration},
		WaitBetweenChaosDuringLoad: 1 * time.Minute,
	}

	timeUnit, _ := utils.GetEnv("CCIP_LOAD_TEST_RATEUNIT")
	if timeUnit != "" {
		d, err := time.ParseDuration(timeUnit)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.Load.TimeUnit = d
		}
	}

	schedule, _ := utils.GetEnv("CCIP_LOAD_TEST_STEP_DURATION")
	if schedule != "" {
		steps := strings.Split(schedule, ",")
		var durations []time.Duration
		for _, i := range steps {
			d, err := time.ParseDuration(i)
			if err != nil {
				allError = multierr.Append(allError, err)
			} else {
				durations = append(durations, d)
			}
		}
		p.Load.StepDuration = durations
	}

	inputRps, _ := utils.GetEnv("CCIP_LOAD_TEST_RATE")
	if inputRps != "" {
		var rpss []int64
		inputRpss := strings.Split(inputRps, ",")
		for _, r := range inputRpss {
			rps, err := strconv.ParseInt(r, 10, 64)
			if err != nil {
				allError = multierr.Append(allError, err)
			} else {
				rpss = append(rpss, rps)
			}
		}
		p.Load.RequestPerUnitTime = rpss
	}

	// if all phases take max time to complete, then the load test will run for 4 times the individual phase time
	// the goal of setting high timeout is to avoid load test failure due to load generator timeout
	// In case of failure the test should fail for individual phase timeout
	if p.PhaseTimeout.Seconds() > 0 {
		p.Load.LoadTimeOut = time.Duration(p.PhaseTimeout.Minutes()*4) * time.Minute
	}
	if len(p.Load.RequestPerUnitTime) != len(p.Load.StepDuration) {
		allError = multierr.Append(allError, fmt.Errorf(
			"the number of request per unit time %d and step duration %d should be equal",
			len(p.Load.RequestPerUnitTime), len(p.Load.StepDuration)))
	}
	if allError != nil {
		p.Test.Fatal(allError)
	}
}

func (p *CCIPTestConfig) SetNetworkPairs(t *testing.T, lggr zerolog.Logger) error {
	var allError error
	// if network pairs are provided, then use them
	// example usage - CCIP_NETWORK_PAIRS="networkA,networkB|networkC,networkD|networkA,networkC"
	lanes, _ := utils.GetEnv("CCIP_NETWORK_PAIRS")
	if lanes != "" {
		p.NetworkPairs = []NetworkPair{}
		networkPairs := strings.Split(lanes, "|")
		networkMap := make(map[int64]blockchain.EVMNetwork)
		for _, pair := range networkPairs {
			networkNames := strings.Split(pair, ",")
			if len(networkNames) != 2 {
				allError = multierr.Append(allError, fmt.Errorf("invalid network pair"))
			}
			nets := networks.SetNetworks(networkNames)
			p.NetworkPairs = append(p.NetworkPairs, NetworkPair{
				NetworkA: nets[0],
				NetworkB: nets[1],
			})
			if _, ok := networkMap[nets[0].ChainID]; !ok {
				networkMap[nets[0].ChainID] = nets[0]
			}
			if _, ok := networkMap[nets[1].ChainID]; !ok {
				networkMap[nets[1].ChainID] = nets[1]
			}
		}
		for _, net := range networkMap {
			p.AllNetworks = append(p.AllNetworks, net)
		}
		return allError
	}
	// if network pairs are not provided with CCIP_NETWORK_PAIRS, then form all possible network pair combination from SELECTED_NETWORKS
	if len(networks.SelectedNetworks) < 3 {
		lggr.Fatal().
			Interface("SELECTED_NETWORKS", networks.SelectedNetworks).
			Msg("Set source and destination network in index 1 & 2 of env variable SELECTED_NETWORKS")
	}
	p.AllNetworks = networks.SelectedNetworks[1:]
	p.NoOfNetworks = DefaultNoOfNetworks
	inputNoOfNetworks, _ := utils.GetEnv("CCIP_NO_OF_NETWORKS")
	if inputNoOfNetworks != "" {
		n, err := strconv.Atoi(inputNoOfNetworks)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.NoOfNetworks = n
		}
	}
	// skip the first index as it is generally set to Simulated EVM in dev mode to be used by other tests
	simulated := p.AllNetworks[0].Simulated
	for i := 1; i < len(p.AllNetworks); i++ {
		if p.AllNetworks[i].Simulated != simulated {
			t.Fatal("networks must be of the same type either simulated or real")
		}
	}
	// if the networks are not simulated use the first p.NoOfNetworks networks from the selected networks
	if !simulated && len(p.AllNetworks) != p.NoOfNetworks {
		if len(p.AllNetworks) < p.NoOfNetworks {
			allError = multierr.Append(allError, fmt.Errorf("not enough networks provided"))
		} else {
			p.AllNetworks = p.AllNetworks[:p.NoOfNetworks]
		}
	}
	// If provided networks is lesser than the required number of networks
	// and the provided networks are simulated network, create replicas of the provided networks with
	// different chain ids
	if len(p.AllNetworks) < p.NoOfNetworks {
		if simulated {
			actualNoOfNetworks := len(p.AllNetworks)
			n := p.AllNetworks[0]
			for i := 0; i < p.NoOfNetworks-actualNoOfNetworks; i++ {
				chainID := networks.AdditionalSimulatedChainIds[i]
				p.AllNetworks = append(p.AllNetworks, blockchain.EVMNetwork{
					Name:                      fmt.Sprintf("simulated-non-dev%d", len(p.AllNetworks)+1),
					ChainID:                   chainID,
					Simulated:                 true,
					PrivateKeys:               []string{networks.AdditionalSimulatedPvtKeys[i]},
					ChainlinkTransactionLimit: n.ChainlinkTransactionLimit,
					Timeout:                   n.Timeout,
					MinimumConfirmations:      n.MinimumConfirmations,
					GasEstimationBuffer:       n.GasEstimationBuffer,
					ClientImplementation:      n.ClientImplementation,
				})
			}
		}
	}
	lggr.Info().Interface("Networks", p.AllNetworks).Msg("Running tests with networks")
	if p.NoOfNetworks > 2 {
		p.FormNetworkPairCombinations()
	} else {
		p.NetworkPairs = []NetworkPair{
			{
				NetworkA: p.AllNetworks[0],
				NetworkB: p.AllNetworks[1],
			},
		}
	}

	for _, n := range p.NetworkPairs {
		lggr.Info().Str("NetworkA", n.NetworkA.Name).Str("NetworkB", n.NetworkB.Name).Msg("Network Pairs")
	}

	return allError
}

func (p *CCIPTestConfig) FormNetworkPairCombinations() {
	for i := 0; i < p.NoOfNetworks; i++ {
		for j := i + 1; j < p.NoOfNetworks; j++ {
			p.NetworkPairs = append(p.NetworkPairs, NetworkPair{
				NetworkA: p.AllNetworks[i],
				NetworkB: p.AllNetworks[j],
			})
		}
	}
}

// NewCCIPTestConfig collects all test related CCIPTestConfig from environment variables
func NewCCIPTestConfig(t *testing.T, lggr zerolog.Logger, tType string) *CCIPTestConfig {
	var allError error
	p := &CCIPTestConfig{
		Test:                t,
		MsgType:             actions.TokenTransfer,
		PhaseTimeout:        DefaultPhaseTimeout,
		TestDuration:        DefaultTestDuration,
		NodeFunding:         DefaultNodeFunding,
		GethResourceProfile: GethResourceProfile,
	}

	if tType != Smoke {
		p.CLNodeDBResourceProfile = DONDBResourceProfile
	}

	if tType == Load {
		p.EnvTTL = DefaultTTLForLongTests
		p.CLNodeResourceProfile = DONResourceProfile
		p.NodeFunding = NodeFundingForLoad
		p.PhaseTimeout = DefaultPhaseTimeoutForLongTests
	}

	allError = multierr.Append(allError, p.SetNetworkPairs(t, lggr))

	ttlDuration, _ := utils.GetEnv("CCIP_KEEP_ENV_TTL")
	if ttlDuration != "" {
		keepEnvFor, err := time.ParseDuration(ttlDuration)
		if err != nil {
			allError = multierr.Append(allError, fmt.Errorf("invalid KEEP_ENV_TTL %s", ttlDuration))
		} else {
			if keepEnvFor.Minutes() < 20 {
				allError = multierr.Append(allError, fmt.Errorf("invalid timeout %s - must be greater than 20m", keepEnvFor))
			} else {
				p.EnvTTL = keepEnvFor
			}
		}
	}

	phaseTimeOut, _ := utils.GetEnv("CCIP_PHASE_VALIDATION_TIMEOUT")
	if phaseTimeOut != "" {
		timeout, err := time.ParseDuration(phaseTimeOut)
		if err != nil {
			allError = multierr.Append(allError, fmt.Errorf("invalid PHASE_VALIDATION_TIMEOUT %s", phaseTimeOut))
		} else {
			if timeout.Minutes() < 1 || timeout.Minutes() > 50 {
				allError = multierr.Append(allError, fmt.Errorf("invalid timeout %s - must be between 1m and 50m", timeout))
			} else {
				p.PhaseTimeout = timeout
			}
		}
	}

	inputDuration, _ := utils.GetEnv("CCIP_TEST_DURATION")
	if inputDuration != "" {
		d, err := time.ParseDuration(inputDuration)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			if d.Minutes() < 1 {
				allError = multierr.Append(allError, fmt.Errorf("invalid duration %d - should be atleast 1m", d))
			} else {
				p.TestDuration = d
			}
		}
	}

	inputMsgType, _ := utils.GetEnv("CCIP_MSG_TYPE")

	if inputMsgType != "" {
		if inputMsgType != actions.DataOnlyTransfer && inputMsgType != actions.TokenTransfer {
			allError = multierr.Append(allError, fmt.Errorf("invalid msg type %s", inputMsgType))
		} else {
			p.MsgType = inputMsgType
		}
	}

	fundingAmountStr, _ := utils.GetEnv("CCIP_CHAINLINK_NODE_FUNDING")
	if fundingAmountStr != "" {
		fundingAmount, _ := big.NewFloat(0).SetString(fundingAmountStr)
		if fundingAmount == nil {
			allError = multierr.Append(allError, fmt.Errorf("invalid CCIP_CHAINLINK_NODE_FUNDING env variable value: %s", fundingAmountStr))
		} else {
			p.NodeFunding = fundingAmount
		}
	}

	local, _ := utils.GetEnv("CCIP_DEPLOY_ON_LOCAL")
	if local != "" {
		e, err := strconv.ParseBool(local)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.LocalCluster = e
		}
	}

	existing, _ := utils.GetEnv("CCIP_TESTS_ON_EXISTING_DEPLOYMENT")
	if existing != "" {
		e, err := strconv.ParseBool(existing)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.ExistingDeployment = e
		}
	}

	reuse, _ := utils.GetEnv("CCIP_REUSE_CONTRACTS")
	if reuse != "" {
		e, err := strconv.ParseBool(reuse)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.ReuseContracts = e
		}
	}
	if p.ExistingDeployment {
		envName, _ := utils.GetEnv("CCIP_EXISTING_ENV")
		if envName != "" {
			p.ExistingEnv = envName
		} else {
			p.ExistingEnv = fmt.Sprintf("Existing-Deployment-%s", uuid.NewString()[0:5])
		}
	}

	if allError != nil {
		t.Fatal(allError)
	}
	if tType == Load {
		p.setLoadInputs()
	}

	return p
}

type BiDirectionalLaneConfig struct {
	NetworkA     blockchain.EVMNetwork
	NetworkB     blockchain.EVMNetwork
	ForwardLane  *actions.CCIPLane
	ReverseLane  *actions.CCIPLane
	LaneDeployed bool
}

type CCIPTestSetUpOutputs struct {
	Cfg            *CCIPTestConfig
	Lanes          []*BiDirectionalLaneConfig
	Reporter       *testreporters.CCIPTestReporter
	LaneConfigFile string
	LaneConfig     *laneconfig.Lanes
	TearDown       func()
	Env            *actions.CCIPTestEnv
	Balance        *actions.BalanceSheet
}

func (o *CCIPTestSetUpOutputs) AddLanesForNetworkPair(
	lggr zerolog.Logger,
	networkA, networkB blockchain.EVMNetwork,
	chainClientA, chainClientB blockchain.EVMClient,
	transferAmounts []*big.Int,
	tokenDeployerFns []blockchain.ContractDeployer,
	numOfCommitNodes int,
	commitAndExecOnSameDON, bidirectional bool,
	newBootstrap bool,
) error {
	var allErrors error
	t := o.Cfg.Test
	var k8Env *environment.Environment
	ccipEnv := o.Env
	namespace := o.Cfg.ExistingEnv
	if ccipEnv != nil {
		k8Env = ccipEnv.K8Env
		if k8Env != nil {
			namespace = k8Env.Cfg.Namespace
		}
	}
	configureCLNode := !o.Cfg.ExistingDeployment
	setUpFuncs, ctx := errgroup.WithContext(context.Background())

	// Use new set of clients(sourceChainClient,destChainClient)
	// with new header subscriptions(otherwise transactions
	// on one lane will keep on waiting for transactions on other lane for the same network)
	// Currently for simulated network clients(from same network) created with NewEVMClient does not sync nonce
	// ConcurrentEVMClient is a work-around for that.
	sourceChainClientA2B, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	sourceChainClientA2B.ParallelTransactions(true)

	destChainClientA2B, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	destChainClientA2B.ParallelTransactions(true)

	ccipLaneA2B := &actions.CCIPLane{
		Test:              t,
		TestEnv:           ccipEnv,
		SourceChain:       sourceChainClientA2B,
		DestChain:         destChainClientA2B,
		SourceNetworkName: actions.NetworkName(networkA.Name),
		DestNetworkName:   actions.NetworkName(networkB.Name),
		ValidationTimeout: o.Cfg.PhaseTimeout,
		SentReqs:          make(map[int64]actions.CCIPRequest),
		TotalFee:          big.NewInt(0),
		Balance:           o.Balance,
		Context:           ctx,
		CommonContractsWg: &sync.WaitGroup{},
	}
	ccipLaneA2B.SrcNetworkLaneCfg, err = o.LaneConfig.ReadLaneConfig(networkA.Name)
	require.NoError(t, err, "Reading lane config shouldn't fail")
	ccipLaneA2B.DstNetworkLaneCfg, err = o.LaneConfig.ReadLaneConfig(networkB.Name)
	require.NoError(t, err, "Reading lane config shouldn't fail")

	ccipLaneA2B.Logger = lggr.With().Str("env", namespace).Str("Lane",
		fmt.Sprintf("%s-->%s", ccipLaneA2B.SourceNetworkName, ccipLaneA2B.DestNetworkName)).Logger()
	ccipLaneA2B.Reports = o.Reporter.AddNewLane(fmt.Sprintf("%d To %d",
		networkA.ChainID, networkB.ChainID), ccipLaneA2B.Logger)

	bidirectionalLane := &BiDirectionalLaneConfig{
		NetworkA:     networkA,
		NetworkB:     networkB,
		ForwardLane:  ccipLaneA2B,
		LaneDeployed: true,
	}

	var ccipLaneB2A *actions.CCIPLane

	if bidirectional {
		sourceChainClientB2A, err := blockchain.ConcurrentEVMClient(networkB, k8Env, destChainClientA2B)
		require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
		sourceChainClientB2A.ParallelTransactions(true)

		destChainClientB2A, err := blockchain.ConcurrentEVMClient(networkA, k8Env, sourceChainClientA2B)
		require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
		destChainClientB2A.ParallelTransactions(true)

		ccipLaneB2A = &actions.CCIPLane{
			Test:              t,
			TestEnv:           ccipEnv,
			SourceNetworkName: actions.NetworkName(networkB.Name),
			DestNetworkName:   actions.NetworkName(networkA.Name),
			SourceChain:       sourceChainClientB2A,
			DestChain:         destChainClientB2A,
			ValidationTimeout: o.Cfg.PhaseTimeout,
			Balance:           o.Balance,
			SentReqs:          make(map[int64]actions.CCIPRequest),
			TotalFee:          big.NewInt(0),
			Context:           ctx,
			CommonContractsWg: &sync.WaitGroup{},
			SrcNetworkLaneCfg: ccipLaneA2B.DstNetworkLaneCfg,
			DstNetworkLaneCfg: ccipLaneA2B.SrcNetworkLaneCfg,
		}
		ccipLaneB2A.Logger = lggr.With().Str("env", namespace).Str("Lane",
			fmt.Sprintf("%s-->%s", ccipLaneB2A.SourceNetworkName, ccipLaneB2A.DestNetworkName)).Logger()
		ccipLaneB2A.Reports = o.Reporter.AddNewLane(
			fmt.Sprintf("%d To %d", networkB.ChainID, networkA.ChainID), ccipLaneB2A.Logger)
		bidirectionalLane.ReverseLane = ccipLaneB2A
	}
	o.Lanes = append(o.Lanes, bidirectionalLane)

	ccipLaneA2B.CommonContractsWg.Add(1)
	setUpFuncs.Go(func() error {
		lggr.Info().Msgf("Setting up lane %s to %s", networkA.Name, networkB.Name)
		err := ccipLaneA2B.DeployNewCCIPLane(numOfCommitNodes, commitAndExecOnSameDON, nil, nil,
			transferAmounts, tokenDeployerFns, newBootstrap, configureCLNode, o.Cfg.ExistingDeployment)
		if err != nil {
			allErrors = multierr.Append(allErrors, fmt.Errorf("deploying lane %s to %s; err - %+v", networkA.Name, networkB.Name, err))
		}
		return err
	})

	if ccipLaneB2A != nil {
		ccipLaneB2A.CommonContractsWg.Add(1)
	}

	setUpFuncs.Go(func() error {
		if bidirectional {
			ccipLaneA2B.CommonContractsWg.Wait()
			srcCommon := ccipLaneA2B.Dest.Common.CopyAddresses(ccipLaneB2A.Context, ccipLaneB2A.SourceChain, o.Cfg.ExistingDeployment)
			destCommon := ccipLaneA2B.Source.Common.CopyAddresses(ccipLaneB2A.Context, ccipLaneB2A.DestChain, o.Cfg.ExistingDeployment)
			lggr.Info().Msgf("Setting up lane %s to %s", networkB.Name, networkA.Name)
			err := ccipLaneB2A.DeployNewCCIPLane(numOfCommitNodes, commitAndExecOnSameDON, srcCommon, destCommon,
				transferAmounts, tokenDeployerFns, false, configureCLNode, o.Cfg.ExistingDeployment)
			if err != nil {
				allErrors = multierr.Append(allErrors, fmt.Errorf("deploying lane %s to %s; err -  %+v", networkB.Name, networkA.Name, err))
			}
			return err
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
			return allErrors
		}
	}
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
	numOfCLNodes int64,
	transferAmounts []*big.Int,
	tokenDeployerFns []blockchain.ContractDeployer,
	numOfCommitNodes int, commitAndExecOnSameDON, bidirectional bool,
	inputs *CCIPTestConfig,
) *CCIPTestSetUpOutputs {
	var (
		ccipEnv *actions.CCIPTestEnv
		k8Env   *environment.Environment
		ctx     context.Context
		err     error
		chains  []blockchain.EVMClient
	)
	filename := fmt.Sprintf("./tmp_%s.json", strings.ReplaceAll(t.Name(), "/", "_"))
	inputs.Test = t
	setUpArgs := &CCIPTestSetUpOutputs{
		Cfg:            inputs,
		Reporter:       testreporters.NewCCIPTestReporter(t, lggr),
		LaneConfigFile: filename,
		Balance:        actions.NewBalanceSheet(),
	}
	_, err = os.Stat(setUpArgs.LaneConfigFile)
	if err == nil {
		// remove the existing lane config file
		err = os.Remove(setUpArgs.LaneConfigFile)
		require.NoError(t, err, "error while removing existing lane config file - %s", setUpArgs.LaneConfigFile)
	}
	if inputs.ExistingDeployment || inputs.ReuseContracts {
		setUpArgs.LaneConfig, err = laneconfig.ReadLanesFromExistingDeployment()
		require.NoError(t, err)
	} else {
		setUpArgs.LaneConfig, err = laneconfig.CreateDeploymentJSON(setUpArgs.LaneConfigFile)
		require.NoError(t, err)
		if setUpArgs.LaneConfig == nil {
			setUpArgs.LaneConfig = &laneconfig.Lanes{LaneConfigs: make(map[string]*laneconfig.LaneConfig)}
		}
	}

	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	configureCLNode := !inputs.ExistingDeployment
	var deployCL func() error
	var local *test_env.CLClusterTestEnv
	if configureCLNode {
		if inputs.LocalCluster {
			local, deployCL = DeployLocalCluster(t, numOfCLNodes, inputs.AllNetworks)
			ccipEnv = &actions.CCIPTestEnv{
				LocalCluster: local,
			}
		} else {
			clProps := make(map[string]interface{})
			clProps["replicas"] = strconv.FormatInt(numOfCLNodes, 10)
			clProps["db"] = inputs.CLNodeDBResourceProfile
			clProps["chainlink"] = map[string]interface{}{
				"resources": inputs.CLNodeResourceProfile,
			}

			// deploy the env if configureCLNode is true
			k8Env = DeployEnvironments(
				t,
				&environment.Config{
					TTL:             inputs.EnvTTL,
					NamespacePrefix: envName,
					Test:            t,
				}, clProps, inputs.GethResourceProfile, inputs.AllNetworks)
			ccipEnv = &actions.CCIPTestEnv{K8Env: k8Env}
		}

		ccipEnv.CLNodeWithKeyReady, ctx = errgroup.WithContext(parent)
		setUpArgs.Env = ccipEnv
		if ccipEnv.K8Env != nil && ccipEnv.K8Env.WillUseRemoteRunner() {
			return setUpArgs
		}
	} else {
		// if configureCLNode is false, use a placeholder env to create remote runner
		k8Env = environment.New(
			&environment.Config{
				TTL:             inputs.EnvTTL,
				NamespacePrefix: envName,
				Test:            t,
			})
		err = k8Env.Run()
		require.NoErrorf(t, err, "error creating environment remote runner")
		setUpArgs.Env = &actions.CCIPTestEnv{K8Env: k8Env}
		if k8Env.WillUseRemoteRunner() {
			return setUpArgs
		}
	}

	chainByChainID := make(map[int64]blockchain.EVMClient)
	if inputs.LocalCluster {
		require.NotNil(t, ccipEnv.LocalCluster, "Local cluster shouldn't be nil")
		for _, n := range ccipEnv.LocalCluster.PrivateChain {
			primaryNode := n.GetPrimaryNode()
			require.NotNil(t, primaryNode, "Primary node is nil in PrivateChain interface")
			chainByChainID[primaryNode.GetEVMClient().GetChainID().Int64()] = primaryNode.GetEVMClient()
			chains = append(chains, primaryNode.GetEVMClient())
		}
	} else {
		for _, network := range inputs.AllNetworks {
			ec, err := blockchain.NewEVMClient(network, k8Env)
			require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
			chains = append(chains, ec)
			chainByChainID[network.ChainID] = ec
		}
	}

	t.Cleanup(func() {
		if configureCLNode {
			lggr.Info().Msg("Tearing down the environment")
			if ccipEnv.LocalCluster != nil {
				err := ccipEnv.LocalCluster.Terminate()
				require.NoError(t, err, "Local cluster termination shouldn't fail")
				for k := range setUpArgs.Reporter.LaneStats {
					setUpArgs.Reporter.LaneStats[k].Finalize(k)
				}
				return
			}
			err = integrationactions.TeardownSuite(t, ccipEnv.K8Env, utils.ProjectRoot, ccipEnv.CLNodes, setUpArgs.Reporter,
				zapcore.ErrorLevel, chains...)
			require.NoError(t, err, "Environment teardown shouldn't fail")
		} else {
			//just print
			for k := range setUpArgs.Reporter.LaneStats {
				setUpArgs.Reporter.LaneStats[k].Finalize(k)
			}
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
			return ccipEnv.SetUpNodesAndKeys(ctx, inputs.NodeFunding, chains)
		})
	}
	if inputs.MsgType == actions.DataOnlyTransfer {
		transferAmounts = []*big.Int{}
	}
	for i, n := range inputs.NetworkPairs {
		newBootstrap := false
		if i == 0 {
			// create bootstrap job once
			newBootstrap = true
		}
		var ok bool
		inputs.NetworkPairs[i].ChainClientA, ok = chainByChainID[n.NetworkA.ChainID]
		require.True(t, ok, "Chain client for chainID %d not found", n.NetworkA.ChainID)
		inputs.NetworkPairs[i].ChainClientB, ok = chainByChainID[n.NetworkB.ChainID]
		require.True(t, ok, "Chain client for chainID %d not found", n.NetworkB.ChainID)

		n.NetworkA = *inputs.NetworkPairs[i].ChainClientA.GetNetworkConfig()
		n.NetworkB = *inputs.NetworkPairs[i].ChainClientB.GetNetworkConfig()

		// if sequential lane addition is true, continue after adding the first bidirectional lane
		// and add rest of the lanes later, after the previously added lane(s) starts getting requests
		// this is mainly used for testing the new lane addition functionality while other lanes are running
		if i > 0 && inputs.SequentialLaneAddition {
			continue
		}
		err = setUpArgs.AddLanesForNetworkPair(
			lggr, n.NetworkA, n.NetworkB,
			chainByChainID[n.NetworkA.ChainID], chainByChainID[n.NetworkB.ChainID],
			transferAmounts, tokenDeployerFns, numOfCommitNodes, commitAndExecOnSameDON,
			bidirectional, newBootstrap)
		require.NoError(t, err)
		err = laneconfig.WriteLanesToJSON(setUpArgs.LaneConfigFile, setUpArgs.LaneConfig)
		require.NoError(t, err)
	}

	setUpArgs.TearDown = func() {
		for _, lanes := range setUpArgs.Lanes {
			lanes.ForwardLane.CleanUp()
			if lanes.ReverseLane != nil {
				lanes.ReverseLane.CleanUp()
			}
		}
	}
	return setUpArgs
}

// CCIPExistingDeploymentTestSetUp is same as CCIPDefaultTestSetUp
// except it's called when
// 1. contracts are already deployed on live networks
// 2. CL nodes are set up and configured with existing contracts
// 3. No k8 env deployment is needed
// It reuses already deployed contracts from the addresses provided in ../contracts/ccip/laneconfig/contracts.json
// Returns -
// CCIPLane for NetworkA --> NetworkB
// CCIPLane for NetworkB --> NetworkA
func CCIPExistingDeploymentTestSetUp(
	t *testing.T,
	lggr zerolog.Logger,
	transferAmounts []*big.Int,
	bidirectional bool,
	input *CCIPTestConfig,
) *CCIPTestSetUpOutputs {
	return CCIPDefaultTestSetUp(t, lggr, "ccip-runner", 0, transferAmounts,
		nil, 0, false, bidirectional, input)
}

func DeployLocalCluster(
	t *testing.T,
	noOfCLNodes int64,
	networks []blockchain.EVMNetwork,
) (*test_env.CLClusterTestEnv, func() error) {
	env, err := test_env.NewCLTestEnvBuilder().
		WithPrivateGethChains(networks).
		WithMockServer(1).
		Build()
	require.NoError(t, err)
	// a func to start the CL nodes asynchronously
	deployCL := func() error {
		var nonDevGethNetworks []blockchain.EVMNetwork
		for i, n := range env.PrivateChain {
			primaryNode := n.GetPrimaryNode()
			require.NotNil(t, primaryNode, "Primary node is nil in PrivateChain interface")
			nonDevGethNetworks = append(nonDevGethNetworks, *n.GetNetworkConfig())
			nonDevGethNetworks[i].URLs = []string{primaryNode.GetInternalWsUrl()}
			nonDevGethNetworks[i].HTTPURLs = []string{primaryNode.GetInternalHttpUrl()}
		}
		if nonDevGethNetworks == nil {
			return errors.New("cannot create nodes with custom config without nonDevGethNetworks")
		}
		toml, err := node.NewConfigFromToml(ccipnode.CCIPTOML,
			node.WithPrivateEVMs(nonDevGethNetworks))
		if err != nil {
			return err
		}
		return env.StartClNodes(toml, int(noOfCLNodes))
	}
	return env, deployCL
}

// DeployEnvironments deploys K8 env for CCIP tests. For tests running on simulated geth it deploys -
// 1. two simulated geth network in non-dev mode
// 2. mockserver ( to set mock price feed details)
// 3. chainlink nodes
func DeployEnvironments(
	t *testing.T,
	envconfig *environment.Config,
	clProps map[string]interface{},
	gethResource map[string]interface{},
	networks []blockchain.EVMNetwork,
) *environment.Environment {
	testEnvironment := environment.New(envconfig)
	numOfTxNodes := 1
	for _, network := range networks {
		if network.Simulated {
			testEnvironment.
				AddHelm(reorg.New(&reorg.Props{
					NetworkName: network.Name,
					NetworkType: "simulated-geth-non-dev",
					Values: map[string]interface{}{
						"geth": map[string]interface{}{
							"genesis": map[string]interface{}{
								"networkId": fmt.Sprint(network.ChainID),
							},
							"tx": map[string]interface{}{
								"replicas":  strconv.Itoa(numOfTxNodes),
								"resources": gethResource,
							},
							"miner": map[string]interface{}{
								"replicas":  "0",
								"resources": gethResource,
							},
						},
						"bootnode": map[string]interface{}{
							"replicas": "1",
						},
					},
				}))
		}
	}

	err := testEnvironment.
		AddHelm(mockserver_cfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		Run()
	require.NoError(t, err)
	if testEnvironment.WillUseRemoteRunner() {
		return testEnvironment
	}
	urlFinder := func(network blockchain.EVMNetwork) ([]string, []string) {
		if !network.Simulated {
			return network.URLs, network.HTTPURLs
		}
		networkName := network.Name
		var internalWsURLs, internalHttpURLs []string
		for i := 0; i < numOfTxNodes; i++ {
			podName := fmt.Sprintf("%s-ethereum-geth:%d", networkName, i)
			txNodeInternalWs, err := testEnvironment.Fwd.FindPort(podName, "geth", "ws-rpc").As(client.RemoteConnection, client.WS)
			require.NoError(t, err, "Error finding WS ports")
			internalWsURLs = append(internalWsURLs, txNodeInternalWs)
			txNodeInternalHttp, err := testEnvironment.Fwd.FindPort(podName, "geth", "http-rpc").As(client.RemoteConnection, client.HTTP)
			require.NoError(t, err, "Error finding HTTP ports")
			internalHttpURLs = append(internalHttpURLs, txNodeInternalHttp)
		}
		return internalWsURLs, internalHttpURLs
	}
	var nets []blockchain.EVMNetwork
	for i := range networks {
		nets = append(nets, networks[i])
		nets[i].URLs, nets[i].HTTPURLs = urlFinder(networks[i])
		// skip adding blockscout for simplified deployments
		// uncomment the following to debug on-chain transactions
		/*
			testEnvironment.AddChart(blockscout.New(&blockscout.Props{
					Name:    fmt.Sprintf("%s-blockscout", networks[i].Name),
					WsURL:   networks[i].URLs[0],
					HttpURL: networks[i].HTTPURLs[0],
				}))
		*/
	}

	tomlCfg, err := node.NewConfigFromToml(ccipnode.CCIPTOML, ccipnode.WithPrivateEVMs(nets))
	tomlStr, err := tomlCfg.TOMLString()
	require.NoError(t, err)
	clProps["toml"] = tomlStr

	err = testEnvironment.
		AddHelm(chainlink.New(0, clProps)).
		Run()
	require.NoError(t, err)
	return testEnvironment
}
