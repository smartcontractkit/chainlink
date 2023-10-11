package testsetups

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
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
	integrationnodes "github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
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
			"cpu":    "2",
			"memory": "4Gi",
		},
		"limits": map[string]interface{}{
			"cpu":    "2",
			"memory": "4Gi",
		},
	}
	DONDBResourceProfile = map[string]interface{}{
		"image": map[string]interface{}{
			"image":   "postgres",
			"version": "13.12",
		},
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
	KeepEnvAlive            bool
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
	AllNetworks             map[string]blockchain.EVMNetwork
	SelectedNetworks        []blockchain.EVMNetwork
	NetworkPairs            []NetworkPair
	NoOfNetworks            int
	NoOfLanesPerPair        int
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
	if p.NoOfLanesPerPair > 1 {
		firstOfPairs[0].Name = fmt.Sprintf("%s-%d", firstOfPairs[0].Name, 1)
		secondOfPairs[0].Name = fmt.Sprintf("%s-%d", secondOfPairs[0].Name, 1)
		for i := 1; i < p.NoOfLanesPerPair; i++ {
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
	noRouter, _ := utils.GetEnv("CCIP_NO_OF_LANES_PER_PAIR")
	if noRouter != "" {
		n, err := strconv.Atoi(noRouter)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.NoOfLanesPerPair = n
		}
	}
	// if network pairs are provided, then use them
	// example usage - CCIP_NETWORK_PAIRS="networkA,networkB|networkC,networkD|networkA,networkC"
	lanes, _ := utils.GetEnv("CCIP_NETWORK_PAIRS")
	if lanes != "" {
		p.NetworkPairs = []NetworkPair{}
		networkPairs := strings.Split(lanes, "|")
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
	// if network pairs are not provided with CCIP_NETWORK_PAIRS, then form all possible network pair combination from SELECTED_NETWORKS
	if len(networks.SelectedNetworks) < 3 {
		lggr.Fatal().
			Interface("SELECTED_NETWORKS", networks.SelectedNetworks).
			Msg("Set source and destination network in index 1 & 2 of env variable SELECTED_NETWORKS")
	}

	// skip the first index as it is generally set to Simulated EVM in dev mode to be used by other tests
	p.SelectedNetworks = networks.SelectedNetworks[1:]
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

	// if the networks are not simulated use the first p.NoOfNetworks networks from the selected networks
	if !simulated && len(p.SelectedNetworks) != p.NoOfNetworks {
		if len(p.SelectedNetworks) < p.NoOfNetworks {
			allError = multierr.Append(allError, fmt.Errorf("not enough networks provided"))
		} else {
			p.SelectedNetworks = p.SelectedNetworks[:p.NoOfNetworks]
		}
	}
	// If provided networks is lesser than the required number of networks
	// and the provided networks are simulated network, create replicas of the provided networks with
	// different chain ids
	if len(p.SelectedNetworks) < p.NoOfNetworks {
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
			for i := 0; i < p.NoOfNetworks-actualNoOfNetworks; i++ {
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

	if p.NoOfNetworks > 2 {
		p.FormNetworkPairCombinations()
	} else {
		p.AddPairToNetworkList(p.SelectedNetworks[0], p.SelectedNetworks[1])
	}

	for _, n := range p.NetworkPairs {
		lggr.Info().Str("NetworkA", n.NetworkA.Name).Str("NetworkB", n.NetworkB.Name).Msg("Network Pairs")
	}
	lggr.Info().Int("Pairs", len(p.NetworkPairs)).Msg("No Of Lanes")

	return allError
}

func (p *CCIPTestConfig) FormNetworkPairCombinations() {
	for i := 0; i < p.NoOfNetworks; i++ {
		for j := i + 1; j < p.NoOfNetworks; j++ {
			p.AddPairToNetworkList(p.SelectedNetworks[i], p.SelectedNetworks[j])
		}
	}
}

func SetResourceProfile(defaultcpu, defaultmem, cpu, mem string) map[string]interface{} {
	if cpu == "" {
		cpu = defaultcpu
	}
	if mem == "" {
		mem = defaultmem
	}
	return map[string]interface{}{
		"requests": map[string]interface{}{
			"cpu":    cpu,
			"memory": mem,
		},
		"limits": map[string]interface{}{
			"cpu":    cpu,
			"memory": mem,
		},
	}
}

// NewCCIPTestConfig collects all test related CCIPTestConfig from environment variables
func NewCCIPTestConfig(t *testing.T, lggr zerolog.Logger, tType string) *CCIPTestConfig {
	var allError error
	nodeMem, _ := utils.GetEnv("CCIP_NODE_MEM")
	nodeCPU, _ := utils.GetEnv("CCIP_NODE_CPU")
	DONResourceProfile["resources"] = SetResourceProfile("2", "4Gi", nodeCPU, nodeMem)

	dbMem, _ := utils.GetEnv("CCIP_DB_MEM")
	dbCPU, _ := utils.GetEnv("CCIP_DB_CPU")
	DONDBResourceProfile["resources"] = SetResourceProfile("2", "4Gi", dbCPU, dbMem)

	dbArgs, _ := utils.GetEnv("CCIP_DB_ARGS")
	if dbArgs != "" {
		args := strings.Split(dbArgs, ",")
		var formattedArgs []string
		for _, arg := range args {
			formattedArgs = append(formattedArgs, "-c")
			formattedArgs = append(formattedArgs, arg)
		}
		DONDBResourceProfile["additionalArgs"] = formattedArgs
	}

	ccipTOML, _ := utils.GetEnv("CCIP_TOML_PATH")
	if ccipTOML != "" {
		tomlFile, err := os.Open(ccipTOML)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			defer tomlFile.Close()
			_, err := tomlFile.Read(node.CCIPTOML)
			if err != nil {
				allError = multierr.Append(allError, err)
			}
		}
	}

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

	allError = multierr.Append(allError, p.SetNetworkPairs(lggr))

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

	alive, _ := utils.GetEnv("CCIP_KEEP_ENV_ALIVE")
	if alive != "" {
		e, err := strconv.ParseBool(alive)
		if err != nil {
			allError = multierr.Append(allError, err)
		} else {
			p.KeepEnvAlive = e
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
		}
	}

	if allError != nil {
		t.Fatal(allError)
	}
	if tType == Load {
		p.setLoadInputs()
	}

	// if no of NetworkPairs are more than 3 , need to increase the db profile
	if len(p.NetworkPairs) > 10 {
		p.CLNodeDBResourceProfile = DONDBResourceProfile
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
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %v", networkCfg.Name, err))
	}

	chain.ParallelTransactions(true)
	defer chain.Close()
	ccipCommon, err := actions.DefaultCCIPModule(lggr, chain, o.Cfg.ExistingDeployment)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create ccip common module for %s: %v", networkCfg.Name, err))
	}

	cfg := o.LaneConfig.ReadLaneConfig(networkCfg.Name)

	err = ccipCommon.DeployContracts(noOfTokens, tokenDeployerFns, cfg)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy common ccip contracts for %s: %v", networkCfg.Name, err))
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
	sourceChainClientA2B, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA, lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %v", networkA.Name, err))
	}

	sourceChainClientA2B.ParallelTransactions(true)

	destChainClientA2B, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB, lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %v", networkB.Name, err))
	}
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
		NetworkA:     networkA,
		NetworkB:     networkB,
		ForwardLane:  ccipLaneA2B,
		LaneDeployed: true,
	}

	var ccipLaneB2A *actions.CCIPLane

	if bidirectional {
		sourceChainClientB2A, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB, lggr)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %v", networkB.Name, err))
		}
		sourceChainClientB2A.ParallelTransactions(true)

		destChainClientB2A, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA, lggr)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %v", networkA.Name, err))
		}
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
			allErrors = multierr.Append(allErrors, fmt.Errorf("deploying lane %s to %s; err - %+v", networkA.Name, networkB.Name, err))
			return err
		}
		err = o.LaneConfig.WriteLaneConfig(networkA.Name, srcConfig)
		if err != nil {
			allErrors = multierr.Append(allErrors, fmt.Errorf("writing lane config for %s; err - %+v", networkA.Name, err))
			return err
		}
		err = o.LaneConfig.WriteLaneConfig(networkB.Name, destConfig)
		if err != nil {
			allErrors = multierr.Append(allErrors, fmt.Errorf("writing lane config for %s; err - %+v", networkB.Name, err))
			return err
		}
		return nil
	})

	setUpFuncs.Go(func() error {
		if bidirectional {
			lggr.Info().Msgf("Setting up lane %s to %s", networkB.Name, networkA.Name)
			srcConfig, destConfig, err := ccipLaneB2A.DeployNewCCIPLane(numOfCommitNodes, commitAndExecOnSameDON, networkBCmn, networkACmn,
				transferAmounts, o.BootstrapAdded, configureCLNode, o.JobAddGrp)
			if err != nil {
				allErrors = multierr.Append(allErrors, fmt.Errorf("deploying lane %s to %s; err -  %+v", networkB.Name, networkA.Name, err))
				return err
			}

			err = o.LaneConfig.WriteLaneConfig(networkB.Name, srcConfig)
			if err != nil {
				allErrors = multierr.Append(allErrors, fmt.Errorf("writing lane config for %s; err - %+v", networkA.Name, err))
				return err
			}
			err = o.LaneConfig.WriteLaneConfig(networkA.Name, destConfig)
			if err != nil {
				allErrors = multierr.Append(allErrors, fmt.Errorf("writing lane config for %s; err - %+v", networkB.Name, err))
				return err
			}
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
			return allErrors
		}
	}
}

func (o *CCIPTestSetUpOutputs) StartEventWatchers() {
	for _, lane := range o.ReadLanes() {
		err := lane.ForwardLane.StartEventWatchers()
		require.NoError(o.Cfg.Test, err)
		err = lane.ReverseLane.StartEventWatchers()
		require.NoError(o.Cfg.Test, err)
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
		priceUpdateGrp.Go(func() error {
			return waitForUpdate(*lanes.ReverseLane)
		})
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
	numOfCLNodes int,
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
	if inputs.MsgType == actions.DataOnlyTransfer {
		transferAmounts = []*big.Int{}
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
	_, err = os.Stat(setUpArgs.LaneConfigFile)
	if err == nil {
		// remove the existing lane config file
		err = os.Remove(setUpArgs.LaneConfigFile)
		require.NoError(t, err, "error while removing existing lane config file - %s", setUpArgs.LaneConfigFile)
	}

	setUpArgs.LaneConfig, err = laneconfig.ReadLanesFromExistingDeployment()
	require.NoError(t, err)

	if setUpArgs.LaneConfig == nil {
		setUpArgs.LaneConfig = &laneconfig.Lanes{LaneConfigs: make(map[string]*laneconfig.LaneConfig)}
	}

	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	configureCLNode := !inputs.ExistingDeployment
	var deployCL func() error
	var local *test_env.CLClusterTestEnv
	if configureCLNode {
		if inputs.LocalCluster {
			local, deployCL = DeployLocalCluster(t, numOfCLNodes, inputs.SelectedNetworks)
			ccipEnv = &actions.CCIPTestEnv{
				LocalCluster: local,
			}
		} else {
			clProps := make(map[string]interface{})
			clProps["replicas"] = numOfCLNodes
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
				}, clProps, inputs.GethResourceProfile, inputs.SelectedNetworks)
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
		for _, n := range inputs.SelectedNetworks {
			if _, ok := chainByChainID[n.ChainID]; ok {
				continue
			}
			ec, err := blockchain.NewEVMClient(n, k8Env, lggr)
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
			if inputs.KeepEnvAlive {
				printStats()
				return
			}
			lggr.Info().Msg("Tearing down the environment")
			err = integrationactions.TeardownSuite(t, ccipEnv.K8Env, utils.ProjectRoot, ccipEnv.CLNodes, setUpArgs.Reporter,
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
			return ccipEnv.SetUpNodesAndKeys(ctx, inputs.NodeFunding, chains, lggr)
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
		if setUpArgs.Cfg.NoOfLanesPerPair > 1 {
			regex := regexp.MustCompile(`-(\d+)$`)
			networkNameToReadCfg := regex.ReplaceAllString(n, "")
			// if reuse contracts is true, copy common contracts from the same network except the router contract
			setUpArgs.LaneConfig.CopyCommonContracts(networkNameToReadCfg, n, inputs.ReuseContracts, inputs.MsgType == actions.TokenTransfer)
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
			return setUpArgs.DeployChainContracts(chain, net, len(transferAmounts), tokenDeployerFns, lggr)
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

		// if sequential lane addition is true, continue after adding the first bidirectional lane
		// and add rest of the lanes later, after the previously added lane(s) starts getting requests
		// this is mainly used for testing the new lane addition functionality while other lanes are running
		if i > 0 && inputs.SequentialLaneAddition {
			continue
		}
		laneAddGrp.Go(func() error {
			return setUpArgs.AddLanesForNetworkPair(
				lggr, n.NetworkA, n.NetworkB,
				chainByChainID[n.NetworkA.ChainID], chainByChainID[n.NetworkB.ChainID],
				transferAmounts, numOfCommitNodes, commitAndExecOnSameDON,
				bidirectional)
		})
	}
	require.NoError(t, laneAddGrp.Wait())
	err = laneconfig.WriteLanesToJSON(setUpArgs.LaneConfigFile, setUpArgs.LaneConfig)
	require.NoError(t, err)
	require.Equal(t, len(setUpArgs.Lanes), len(inputs.NetworkPairs),
		"Number of bi-directional lanes should be equal to number of network pairs")

	if !inputs.ExistingDeployment {
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
			err := lanes.ForwardLane.CleanUp(!setUpArgs.Cfg.ExistingDeployment)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			if lanes.ReverseLane != nil {
				// if existing deployment is true, don't attempt to pay ccip fees
				err := lanes.ReverseLane.CleanUp(!setUpArgs.Cfg.ExistingDeployment)
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
	noOfCLNodes int,
	networks []blockchain.EVMNetwork,
) (*test_env.CLClusterTestEnv, func() error) {
	env, err := test_env.NewCLTestEnvBuilder().
		WithPrivateGethChains(networks).
		Build()
	require.NoError(t, err)
	for _, n := range env.PrivateChain {
		primaryNode := n.GetPrimaryNode()
		require.NotNil(t, primaryNode, "Primary node is nil in PrivateChain interface")
		for i, networkCfg := range networks {
			if networkCfg.ChainID == n.GetNetworkConfig().ChainID {
				networks[i].URLs = []string{primaryNode.GetInternalWsUrl()}
				networks[i].HTTPURLs = []string{primaryNode.GetInternalHttpUrl()}
			}
		}
	}
	configOpts := []integrationnodes.NodeConfigOpt{
		node.WithPrivateEVMs(networks),
	}
	// a func to start the CL nodes asynchronously
	deployCL := func() error {
		toml, err := node.NewConfigFromToml(ccipnode.CCIPTOML, configOpts...)
		if err != nil {
			return err
		}
		return env.StartClNodes(toml, noOfCLNodes, "")
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
		if !network.Simulated {
			continue
		}
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
	err := testEnvironment.Run()
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

	tomlCfg, err := node.NewConfigFromToml(
		ccipnode.CCIPTOML,
		ccipnode.WithPrivateEVMs(nets),
	)
	tomlStr, err := tomlCfg.TOMLString()
	require.NoError(t, err)
	clProps["toml"] = tomlStr
	clProps["prometheus"] = true
	err = testEnvironment.
		AddHelm(chainlink.New(0, clProps)).
		Run()
	require.NoError(t, err)
	return testEnvironment
}
