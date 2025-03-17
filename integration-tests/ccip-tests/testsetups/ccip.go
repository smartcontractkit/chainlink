package testsetups

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/blockchain_client_wrapper"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	integrationactions "github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts/laneconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	ccipconfig "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	testutils "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/utils"

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
	// to set default values through test config use sync.once
	setContractVersion sync.Once
	setOCRParams       sync.Once
	setConfigOverrides sync.Once
)

type NetworkPair struct {
	NetworkA     blockchain.EVMNetwork
	NetworkB     blockchain.EVMNetwork
	ChainClientA blockchain.EVMClient
	ChainClientB blockchain.EVMClient
}

// LeaderLane is to hold the details of leader lane source and destination network
type LeaderLane struct {
	source string
	dest   string
}

type CCIPTestConfig struct {
	Test                *testing.T
	EnvInput            *testconfig.Common
	TestGroupInput      *testconfig.CCIPTestGroupConfig
	VersionInput        map[contracts.Name]contracts.Version
	ContractsInput      *testconfig.CCIPContractConfig
	AllNetworks         map[string]blockchain.EVMNetwork
	SelectedNetworks    []blockchain.EVMNetwork
	NetworkPairs        []NetworkPair
	LeaderLanes         []LeaderLane
	GethResourceProfile map[string]interface{}
}

func (c *CCIPTestConfig) useExistingDeployment() bool {
	return pointer.GetBool(c.TestGroupInput.ExistingDeployment)
}

func (c *CCIPTestConfig) useSeparateTokenDeployer() bool {
	return contracts.NeedTokenAdminRegistry() &&
		!pointer.GetBool(c.TestGroupInput.TokenConfig.CCIPOwnerTokens) &&
		!c.useExistingDeployment()
}

func (c *CCIPTestConfig) MultiCallEnabled() bool {
	return pointer.GetBool(c.TestGroupInput.MulticallInOneTx)
}

func (c *CCIPTestConfig) localCluster() bool {
	return pointer.GetBool(c.TestGroupInput.LocalCluster)
}

func (c *CCIPTestConfig) ExistingCLCluster() bool {
	return c.EnvInput.ExistingCLCluster != nil
}

func (c *CCIPTestConfig) CLClusterNeedsUpgrade() bool {
	if c.EnvInput.NewCLCluster == nil {
		return false
	}
	if c.EnvInput.NewCLCluster.Common != nil && c.EnvInput.NewCLCluster.Common.ChainlinkUpgradeImage != nil {
		return true
	}
	for _, node := range c.EnvInput.NewCLCluster.Nodes {
		if node.ChainlinkUpgradeImage != nil {
			return true
		}
	}
	return false
}

func (c *CCIPTestConfig) AddPairToNetworkList(networkA, networkB blockchain.EVMNetwork) {
	if c.AllNetworks == nil {
		c.AllNetworks = make(map[string]blockchain.EVMNetwork)
	}
	firstOfPairs := []blockchain.EVMNetwork{networkA}
	secondOfPairs := []blockchain.EVMNetwork{networkB}
	// if no of lanes per pair is greater than 1, copy common contracts from the same network
	// if no of lanes per pair is more than 1, the networks are added into the inputs.AllNetworks with a suffix of -<lane number>
	// for example, if no of lanes per pair is 2, and the network pairs are called "testnetA", "testnetB",
	//	the network will be added as "testnetA-1", testnetA-2","testnetB-1", testnetB-2"
	// to deploy 4 lanes between same network pair "testnetA", "testnetB".
	// lanes - testnetA-1<->testnetB-1, testnetA-1<-->testnetB-2 , testnetA-2<--> testnetB-1, testnetA-2<--> testnetB-2
	if c.TestGroupInput.NoOfRoutersPerPair > 1 {
		firstOfPairs[0].Name = fmt.Sprintf("%s-%d", firstOfPairs[0].Name, 1)
		secondOfPairs[0].Name = fmt.Sprintf("%s-%d", secondOfPairs[0].Name, 1)
		for i := 1; i < c.TestGroupInput.NoOfRoutersPerPair; i++ {
			netsA := networkA
			netsA.Name = fmt.Sprintf("%s-%d", netsA.Name, i+1)
			netsB := networkB
			netsB.Name = fmt.Sprintf("%s-%d", netsB.Name, i+1)
			firstOfPairs = append(firstOfPairs, netsA)
			secondOfPairs = append(secondOfPairs, netsB)
		}
	}

	for i := range firstOfPairs {
		c.AllNetworks[firstOfPairs[i].Name] = firstOfPairs[i]
		c.AllNetworks[secondOfPairs[i].Name] = secondOfPairs[i]
		c.NetworkPairs = append(c.NetworkPairs, NetworkPair{
			NetworkA: firstOfPairs[i],
			NetworkB: secondOfPairs[i],
		})
	}
}

func (c *CCIPTestConfig) SetNetworkPairs(lggr zerolog.Logger) error {
	var allError error
	var err error
	var inputNetworks []string
	c.SelectedNetworks, inputNetworks, err = c.EnvInput.EVMNetworks()
	if err != nil {
		allError = multierr.Append(allError, fmt.Errorf("failed to get networks: %w", err))
		return allError
	}

	networkByChainName := make(map[string]blockchain.EVMNetwork)
	for i, net := range c.SelectedNetworks {
		networkByChainName[inputNetworks[i]] = net
	}
	// if network pairs are provided, then use them
	if c.TestGroupInput.NetworkPairs != nil {
		networkPairs := c.TestGroupInput.NetworkPairs

		for _, pair := range networkPairs {
			networkNames := strings.Split(pair, ",")
			if len(networkNames) != 2 {
				allError = multierr.Append(allError, fmt.Errorf("invalid network pair"))
			}
			// check if the network names are valid
			network1, ok := networkByChainName[networkNames[0]]
			if !ok {
				allError = multierr.Append(allError, fmt.Errorf("network %s not found in network config", networkNames[0]))
			}
			network2, ok := networkByChainName[networkNames[1]]
			if !ok {
				allError = multierr.Append(allError, fmt.Errorf("network %s not found in network config", networkNames[1]))
			}
			c.AddPairToNetworkList(network1, network2)
		}
		lggr.Info().Int("Pairs", len(c.NetworkPairs)).Msg("No Of Lanes")
		return allError
	}

	if c.TestGroupInput.NoOfNetworks == 0 {
		c.TestGroupInput.NoOfNetworks = len(c.SelectedNetworks)
	}
	// TODO remove this when CTF network timeout is fixed
	for i := range c.SelectedNetworks {
		c.SelectedNetworks[i].Timeout = blockchain.StrDuration{
			Duration: 3 * time.Minute,
		}
	}
	simulated := c.SelectedNetworks[0].Simulated
	for i := 1; i < len(c.SelectedNetworks); i++ {
		if c.SelectedNetworks[i].Simulated != simulated {
			lggr.Fatal().Msg("networks must be of the same type either simulated or real")
		}
	}

	// if the networks are not simulated use the first p.NoOfNetworks networks from the selected networks
	if !simulated && len(c.SelectedNetworks) != c.TestGroupInput.NoOfNetworks {
		if len(c.SelectedNetworks) < c.TestGroupInput.NoOfNetworks {
			allError = multierr.Append(allError, fmt.Errorf("not enough networks provided"))
		} else {
			c.SelectedNetworks = c.SelectedNetworks[:c.TestGroupInput.NoOfNetworks]
		}
	}
	// If provided networks is lesser than the required number of networks
	// and the provided networks are simulated network, create replicas of the provided networks with
	// different chain ids
	if simulated && len(c.SelectedNetworks) < c.TestGroupInput.NoOfNetworks {
		actualNoOfNetworks := len(c.SelectedNetworks)
		n := c.SelectedNetworks[0]
		var chainIDs []int64
		existingChainIDs := make(map[uint64]struct{})
		for _, net := range c.SelectedNetworks {
			if net.ChainID < 0 {
				return fmt.Errorf("negative chain ID: %d", net.ChainID)
			}
			existingChainIDs[uint64(net.ChainID)] = struct{}{}
		}
		for _, id := range chainselectors.TestChainIds() {
			// if the chain id already exists in the already provided selected networks, skip it
			if _, exists := existingChainIDs[id]; exists {
				continue
			}
			if id > math.MaxInt64 {
				return fmt.Errorf("chain ID overflows int64: %d", id)
			}
			chainIDs = append(chainIDs, int64(id))
		}
		for i := 0; i < c.TestGroupInput.NoOfNetworks-actualNoOfNetworks; i++ {
			chainID := chainIDs[i]
			// if i is greater than the number of simulated pvt keys, rotate the keys
			if i > len(networks.AdditionalSimulatedPvtKeys)-1 {
				networks.AdditionalSimulatedPvtKeys = append(networks.AdditionalSimulatedPvtKeys, networks.AdditionalSimulatedPvtKeys...)
			}
			name := fmt.Sprintf("private-chain-%d", len(c.SelectedNetworks)+1)
			c.SelectedNetworks = append(c.SelectedNetworks, blockchain.EVMNetwork{
				Name:      name,
				ChainID:   chainID,
				Simulated: true,
				PrivateKeys: []string{
					networks.AdditionalSimulatedPvtKeys[i],
					"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // second key for token deployments
				},
				ChainlinkTransactionLimit: n.ChainlinkTransactionLimit,
				Timeout:                   n.Timeout,
				MinimumConfirmations:      n.MinimumConfirmations,
				GasEstimationBuffer:       n.GasEstimationBuffer + 1000,
				ClientImplementation:      n.ClientImplementation,
				DefaultGasLimit:           n.DefaultGasLimit,
				FinalityDepth:             n.FinalityDepth,
				SupportsEIP1559:           true,
			})
			if existing, ok := c.EnvInput.Network.AnvilConfigs[strings.ToUpper(n.Name)]; c.EnvInput.Network.AnvilConfigs != nil && ok {
				c.EnvInput.Network.AnvilConfigs[strings.ToUpper(name)] = existing
			}

			chainConfig := &ctf_config.EthereumChainConfig{}
			err := chainConfig.Default()
			if err != nil {
				allError = multierr.Append(allError, fmt.Errorf("failed to get default chain config: %w", err))
			} else {
				chainConfig.ChainID = int(chainID)
				eth1 := ctf_config_types.EthereumVersion_Eth1
				geth := ctf_config_types.ExecutionLayer_Geth

				c.EnvInput.PrivateEthereumNetworks[fmt.Sprint(chainID)] = &ctf_config.EthereumNetworkConfig{
					EthereumVersion:     &eth1,
					ExecutionLayer:      &geth,
					EthereumChainConfig: chainConfig,
				}
			}
		}
	}

	if c.TestGroupInput.NoOfNetworks > 2 {
		c.FormNetworkPairCombinations()
	} else {
		c.AddPairToNetworkList(c.SelectedNetworks[0], c.SelectedNetworks[1])
	}

	// if the number of lanes is lesser than the number of network pairs, choose first c.TestGroupInput.MaxNoOfLanes pairs
	if c.TestGroupInput.MaxNoOfLanes > 0 && c.TestGroupInput.MaxNoOfLanes < len(c.NetworkPairs) {
		var newNetworkPairs []NetworkPair
		denselyConnectedNetworks := make(map[string]struct{})
		// if densely connected networks are provided, choose all the network pairs containing the networks mentioned in the list for DenselyConnectedNetworkChainIds
		if len(c.TestGroupInput.DenselyConnectedNetworkChainIDs) > 0 {
			for _, n := range c.TestGroupInput.DenselyConnectedNetworkChainIDs {
				denselyConnectedNetworks[n] = struct{}{}
			}
			for _, pair := range c.NetworkPairs {
				if _, exists := denselyConnectedNetworks[strconv.FormatInt(pair.NetworkA.ChainID, 10)]; exists {
					newNetworkPairs = append(newNetworkPairs, pair)
				}
			}
		}
		// shuffle the network pairs, we want to randomly distribute the network pairs among all available networks
		rand.Shuffle(len(c.NetworkPairs), func(i, j int) {
			c.NetworkPairs[i], c.NetworkPairs[j] = c.NetworkPairs[j], c.NetworkPairs[i]
		})
		// now add the remaining network pairs by skipping the already covered networks
		// and adding the remaining pair from the shuffled list
		i := len(newNetworkPairs)
		j := 0
		for i < c.TestGroupInput.MaxNoOfLanes {
			pair := c.NetworkPairs[j]
			// if the network is already covered, skip it
			if _, exists := denselyConnectedNetworks[strconv.FormatInt(pair.NetworkA.ChainID, 10)]; !exists {
				newNetworkPairs = append(newNetworkPairs, pair)
				i++
			}
			j++
		}
		c.NetworkPairs = newNetworkPairs
	}

	// setting leader lane details to network pairs if it is enabled and only in simulated environments
	if !pointer.GetBool(c.TestGroupInput.ExistingDeployment) {
		c.defineLeaderLanes(lggr)
	}
	for _, n := range c.NetworkPairs {
		lggr.Info().
			Str("NetworkA", fmt.Sprintf("%s-%d", n.NetworkA.Name, n.NetworkA.ChainID)).
			Str("NetworkB", fmt.Sprintf("%s-%d", n.NetworkB.Name, n.NetworkB.ChainID)).
			Msg("Network Pairs")
	}
	for _, lane := range c.LeaderLanes {
		lggr.Info().
			Str("Source", lane.source).
			Str("Destination", lane.dest).
			Msg("Leader Lane: ")
	}
	lggr.Info().Int("Pairs", len(c.NetworkPairs)).Msg("No Of Lanes")

	return allError
}

// defineLeaderLanes goes over the available network pairs and define one leader lane per destination
func (c *CCIPTestConfig) defineLeaderLanes(lggr zerolog.Logger) {
	if !isLeaderLaneFeatureEnabled(&lggr) {
		return
	}
	// the way we are defining leader lane is simply by tagging the destination as key along with the first source network
	// as value to the map.
	leaderLanes := make(map[string]string)
	for _, n := range c.NetworkPairs {
		if _, ok := leaderLanes[n.NetworkB.Name]; !ok {
			leaderLanes[n.NetworkB.Name] = n.NetworkA.Name
		}
		if pointer.GetBool(c.TestGroupInput.BiDirectionalLane) {
			if _, ok := leaderLanes[n.NetworkA.Name]; !ok {
				leaderLanes[n.NetworkA.Name] = n.NetworkB.Name
			}
		}
	}
	for k, v := range leaderLanes {
		c.LeaderLanes = append(c.LeaderLanes, LeaderLane{
			source: v,
			dest:   k,
		})
	}
}

// isPriceReportingDisabled checks the given lane is leader lane and return boolean accordingly
func (c *CCIPTestConfig) isPriceReportingDisabled(lggr *zerolog.Logger, source, dest string) bool {
	for _, leader := range c.LeaderLanes {
		if leader.source == source && leader.dest == dest {
			lggr.Debug().
				Str("Source", source).
				Str("Destination", dest).
				Msg("Non-leader lane")
			return true
		}
	}
	return false
}

func isLeaderLaneFeatureEnabled(lggr *zerolog.Logger) bool {
	if err := contracts.MatchContractVersionsOrAbove(map[contracts.Name]contracts.Version{
		contracts.OffRampContract: contracts.V1_2_0,
		contracts.OnRampContract:  contracts.V1_2_0,
	}); err != nil {
		lggr.Info().Str("Required contract version", contracts.V1_2_0.String()).Msg("Leader lane feature is not enabled")
		return false
	}
	return true
}

func (c *CCIPTestConfig) FormNetworkPairCombinations() {
	for i := 0; i < c.TestGroupInput.NoOfNetworks; i++ {
		for j := i + 1; j < c.TestGroupInput.NoOfNetworks; j++ {
			c.AddPairToNetworkList(c.SelectedNetworks[i], c.SelectedNetworks[j])
		}
	}
}

func (c *CCIPTestConfig) SetContractVersion() error {
	if c.VersionInput == nil {
		return nil
	}
	for contractName, version := range c.VersionInput {
		err := contracts.CheckVersionSupported(contractName, version)
		if err != nil {
			return err
		}
		contracts.VersionMap[contractName] = version
	}
	return nil
}

func (c *CCIPTestConfig) SetOCRParams() error {
	if c.TestGroupInput.OffRampConfig != nil {
		if c.TestGroupInput.OffRampConfig.InflightExpiry != nil &&
			c.TestGroupInput.OffRampConfig.InflightExpiry.Duration() > 0 {
			actions.InflightExpiryExec = c.TestGroupInput.OffRampConfig.InflightExpiry.Duration()
		}
		if pointer.GetUint32(c.TestGroupInput.OffRampConfig.BatchGasLimit) > 0 {
			actions.BatchGasLimit = pointer.GetUint32(c.TestGroupInput.OffRampConfig.BatchGasLimit)
		}
		if pointer.GetUint32(c.TestGroupInput.OffRampConfig.MaxDataBytes) > 0 {
			actions.MaxDataBytes = pointer.GetUint32(c.TestGroupInput.OffRampConfig.MaxDataBytes)
		}
		if c.TestGroupInput.OffRampConfig.RootSnooze != nil &&
			c.TestGroupInput.OffRampConfig.RootSnooze.Duration() > 0 {
			actions.RootSnoozeTime = c.TestGroupInput.OffRampConfig.RootSnooze.Duration()
		}
	}
	if c.TestGroupInput.CommitInflightExpiry != nil && c.TestGroupInput.CommitInflightExpiry.Duration() > 0 {
		actions.InflightExpiryCommit = c.TestGroupInput.CommitInflightExpiry.Duration()
	}
	return nil
}

// TestConfigOverrideOption is a function that modifies the test config and overrides any values passed in by test files
// This is useful for setting up test specific configurations.
// The return should be a short, explanatory string that describes the change made by the override.
// This is logged at the beginning of the test run.
type TestConfigOverrideOption func(*CCIPTestConfig) string

// UseCCIPOwnerTokens defines whether all tokens are deployed by the same address as the CCIP owner
func UseCCIPOwnerTokens(yes bool) TestConfigOverrideOption {
	return func(c *CCIPTestConfig) string {
		c.TestGroupInput.TokenConfig.CCIPOwnerTokens = pointer.ToBool(yes)
		return fmt.Sprintf("CCIPOwnerTokens set to %t", yes)
	}
}

// WithTokensPerChain sets the number of tokens to deploy on each chain
func WithTokensPerChain(count int) TestConfigOverrideOption {
	return func(c *CCIPTestConfig) string {
		c.TestGroupInput.TokenConfig.NoOfTokensPerChain = pointer.ToInt(count)
		return fmt.Sprintf("NoOfTokensPerChain set to %d", count)
	}
}

// WithMsgDetails sets the message details for the test
func WithMsgDetails(details *testconfig.MsgDetails) TestConfigOverrideOption {
	return func(c *CCIPTestConfig) string {
		c.TestGroupInput.MsgDetails = details
		return "Message set"
	}
}

// WithNoTokensPerMessage sets how many tokens can be sent in a single message
func WithNoTokensPerMessage(noOfTokens int) TestConfigOverrideOption {
	return func(c *CCIPTestConfig) string {
		c.TestGroupInput.MsgDetails.NoOfTokens = pointer.ToInt(noOfTokens)
		return fmt.Sprintf("MsgDetails.NoOfTokens set to %d", noOfTokens)
	}
}

// NewCCIPTestConfig reads the CCIP test config from TOML files, applies any overrides, and configures the test environment
func NewCCIPTestConfig(t *testing.T, lggr zerolog.Logger, tType string, overrides ...TestConfigOverrideOption) *CCIPTestConfig {
	testCfg := ccipconfig.GlobalTestConfig()
	groupCfg, exists := testCfg.CCIP.Groups[tType]
	if !exists {
		t.Fatalf("group config for %s does not exist", tType)
	}
	if tType == ccipconfig.Load {
		if testCfg.CCIP.Env.Logging == nil || testCfg.CCIP.Env.Logging.Loki == nil {
			t.Fatal("loki config is required to be set for load test")
		}
		if testCfg.CCIP.Env.Logging == nil || testCfg.CCIP.Env.Logging.Grafana == nil {
			t.Fatal("grafana config is required for load test")
		}
	}
	if pointer.GetBool(groupCfg.KeepEnvAlive) {
		err := os.Setenv(config.EnvVarKeepEnvironments, "ALWAYS")
		if err != nil {
			t.Fatal(err)
		}
	}
	ccipTestConfig := &CCIPTestConfig{
		Test:                t,
		EnvInput:            testCfg.CCIP.Env,
		ContractsInput:      testCfg.CCIP.Deployments,
		VersionInput:        testCfg.CCIP.ContractVersions,
		TestGroupInput:      groupCfg,
		GethResourceProfile: GethResourceProfile,
	}
	setContractVersion.Do(func() {
		err := ccipTestConfig.SetContractVersion()
		if err != nil {
			t.Fatal(err)
		}
	})
	setOCRParams.Do(func() {
		err := ccipTestConfig.SetOCRParams()
		if err != nil {
			t.Fatal(err)
		}
	})
	setConfigOverrides.Do(func() {
		overrideMessages := []string{}
		for _, override := range overrides {
			if override != nil {
				overrideMessages = append(overrideMessages, override(ccipTestConfig))
			}
		}
		if len(overrideMessages) > 0 {
			lggr.Debug().Int("Overrides", len(overrideMessages)).Msg("Test Specific Config Overrides Applied")
			for _, msg := range overrideMessages {
				lggr.Debug().Msg(msg)
			}
		}
	})
	err := ccipTestConfig.SetNetworkPairs(lggr)
	if err != nil {
		t.Fatal(err)
	}

	return ccipTestConfig
}

type BiDirectionalLaneConfig struct {
	NetworkA    blockchain.EVMNetwork
	NetworkB    blockchain.EVMNetwork
	ForwardLane *actions.CCIPLane
	ReverseLane *actions.CCIPLane
}

type CCIPTestSetUpOutputs struct {
	SetUpContext           context.Context
	Cfg                    *CCIPTestConfig
	LaneContractsByNetwork *sync.Map
	laneMutex              *sync.Mutex
	Lanes                  []*BiDirectionalLaneConfig
	Reporter               *testreporters.CCIPTestReporter
	LaneConfigFile         string
	LaneConfig             *laneconfig.Lanes
	TearDown               func() error
	Env                    *actions.CCIPTestEnv
	Balance                *actions.BalanceSheet
	BootstrapAdded         *atomic.Bool
	JobAddGrp              *errgroup.Group
	SC                     *sentinel.SentinelCoordinator
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
	lggr *zerolog.Logger,
	chainClient blockchain.EVMClient,
	networkCfg blockchain.EVMNetwork,
	noOfTokens int,
	tokenDeployerFns []blockchain.ContractDeployer,
) error {
	var k8Env *environment.Environment
	ccipEnv := o.Env
	if ccipEnv != nil {
		k8Env = ccipEnv.K8Env
	}
	if k8Env != nil && chainClient.NetworkSimulated() {
		networkCfg.URLs = k8Env.URLs[chainClient.GetNetworkConfig().Name]
	}

	mainChainClient, err := blockchain.ConcurrentEVMClient(networkCfg, k8Env, chainClient, *lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkCfg.Name, err))
	}

	mainChainClient.ParallelTransactions(true)
	defer mainChainClient.Close()
	ccipCommon, err := actions.DefaultCCIPModule(
		lggr, o.Cfg.TestGroupInput, mainChainClient,
	)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create ccip common module for %s: %w", networkCfg.Name, err))
	}

	cfg := o.LaneConfig.ReadLaneConfig(networkCfg.Name)

	err = ccipCommon.DeployContracts(noOfTokens, tokenDeployerFns, cfg)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy common ccip contracts for %s: %w", networkCfg.Name, err))
	}
	ccipCommon.WriteLaneConfig(cfg)
	o.LaneContractsByNetwork.Store(networkCfg.Name, cfg)

	return nil
}

func (o *CCIPTestSetUpOutputs) SetupDynamicTokenPriceUpdates() error {
	interval := o.Cfg.TestGroupInput.TokenConfig.DynamicPriceUpdateInterval.Duration()
	covered := make(map[string]struct{})
	for _, lanes := range o.ReadLanes() {
		lane := lanes.ForwardLane
		if _, exists := covered[lane.SourceNetworkName]; !exists {
			covered[lane.SourceNetworkName] = struct{}{}
			err := lane.Source.Common.UpdateTokenPricesAtRegularInterval(lane.Context, lane.Logger, interval, o.LaneConfig.ReadLaneConfig(lane.SourceNetworkName))
			if err != nil {
				return err
			}
		}
		if _, exists := covered[lane.DestNetworkName]; !exists {
			covered[lane.DestNetworkName] = struct{}{}
			err := lane.Dest.Common.UpdateTokenPricesAtRegularInterval(lane.Context, lane.Logger, interval, o.LaneConfig.ReadLaneConfig(lane.DestNetworkName))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *CCIPTestSetUpOutputs) AddLanesForNetworkPair(
	lggr *zerolog.Logger,
	networkA, networkB blockchain.EVMNetwork,
	chainClientA, chainClientB blockchain.EVMClient,
) error {
	var (
		t         = o.Cfg.Test
		allErrors atomic.Error
		k8Env     *environment.Environment
		ccipEnv   = o.Env
		namespace = ""
	)

	if o.Cfg.TestGroupInput.LoadProfile != nil {
		namespace = o.Cfg.TestGroupInput.LoadProfile.TestRunName
	}
	bidirectional := pointer.GetBool(o.Cfg.TestGroupInput.BiDirectionalLane)
	if ccipEnv != nil {
		k8Env = ccipEnv.K8Env
		if k8Env != nil {
			namespace = k8Env.Cfg.Namespace
		}
	}

	setUpFuncs, ctx := errgroup.WithContext(testcontext.Get(t))

	// Use new set of clients(sourceChainClient,destChainClient)
	// with new header subscriptions(otherwise transactions
	// on one lane will keep on waiting for transactions on other lane for the same network)
	// Currently for simulated network clients(from same network) created with NewEVMClient does not sync nonce
	// ConcurrentEVMClient is a work-around for that.
	sourceChainClientA2B, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA, *lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkA.Name, err))
	}

	sourceChainClientA2B.ParallelTransactions(true)

	destChainClientA2B, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB, *lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkB.Name, err))
	}
	destChainClientA2B.ParallelTransactions(true)

	ccipLaneA2B := &actions.CCIPLane{
		Test:              t,
		SourceChain:       sourceChainClientA2B,
		DestChain:         destChainClientA2B,
		SourceNetworkName: actions.NetworkName(networkA.Name),
		DestNetworkName:   actions.NetworkName(networkB.Name),
		ValidationTimeout: o.Cfg.TestGroupInput.PhaseTimeout.Duration(),
		SentReqs:          make(map[common.Hash][]actions.CCIPRequest),
		TotalFee:          big.NewInt(0),
		Balance:           o.Balance,
		Context:           testcontext.Get(t),
	}
	// if it non leader lane, disable the price reporting
	ccipLaneA2B.PriceReportingDisabled = len(o.Cfg.LeaderLanes) > 0 &&
		!o.Cfg.isPriceReportingDisabled(lggr, ccipLaneA2B.SourceNetworkName, ccipLaneA2B.DestNetworkName)

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

	a2blogger := lggr.With().Str("env", namespace).Str("Lane",
		fmt.Sprintf("%s-->%s", ccipLaneA2B.SourceNetworkName, ccipLaneA2B.DestNetworkName)).Logger()
	ccipLaneA2B.Logger = &a2blogger
	ccipLaneA2B.Reports = o.Reporter.AddNewLane(fmt.Sprintf("%s To %s",
		networkA.Name, networkB.Name), ccipLaneA2B.Logger)

	bidirectionalLane := &BiDirectionalLaneConfig{
		NetworkA:    networkA,
		NetworkB:    networkB,
		ForwardLane: ccipLaneA2B,
	}

	var ccipLaneB2A *actions.CCIPLane

	if bidirectional {
		sourceChainClientB2A, err := blockchain.ConcurrentEVMClient(networkB, k8Env, chainClientB, *lggr)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkB.Name, err))
		}
		sourceChainClientB2A.ParallelTransactions(true)

		destChainClientB2A, err := blockchain.ConcurrentEVMClient(networkA, k8Env, chainClientA, *lggr)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkA.Name, err))
		}
		destChainClientB2A.ParallelTransactions(true)

		ccipLaneB2A = &actions.CCIPLane{
			Test:              t,
			SourceNetworkName: actions.NetworkName(networkB.Name),
			DestNetworkName:   actions.NetworkName(networkA.Name),
			SourceChain:       sourceChainClientB2A,
			DestChain:         destChainClientB2A,
			ValidationTimeout: o.Cfg.TestGroupInput.PhaseTimeout.Duration(),
			Balance:           o.Balance,
			SentReqs:          make(map[common.Hash][]actions.CCIPRequest),
			TotalFee:          big.NewInt(0),
			Context:           testcontext.Get(t),
			SrcNetworkLaneCfg: ccipLaneA2B.DstNetworkLaneCfg,
			DstNetworkLaneCfg: ccipLaneA2B.SrcNetworkLaneCfg,
		}
		// if it non leader lane, disable the price reporting
		ccipLaneB2A.PriceReportingDisabled = len(o.Cfg.LeaderLanes) > 0 &&
			!o.Cfg.isPriceReportingDisabled(lggr, ccipLaneB2A.SourceNetworkName, ccipLaneB2A.DestNetworkName)
		b2aLogger := lggr.With().Str("env", namespace).Str("Lane",
			fmt.Sprintf("%s-->%s", ccipLaneB2A.SourceNetworkName, ccipLaneB2A.DestNetworkName)).Logger()
		ccipLaneB2A.Logger = &b2aLogger
		ccipLaneB2A.Reports = o.Reporter.AddNewLane(
			fmt.Sprintf("%s To %s", networkB.Name, networkA.Name), ccipLaneB2A.Logger)
		bidirectionalLane.ReverseLane = ccipLaneB2A
	}
	o.AddToLanes(bidirectionalLane)

	setUpFuncs.Go(func() error {
		lggr.Info().Msgf("Setting up lane %s to %s", networkA.Name, networkB.Name)
		err := ccipLaneA2B.DeployNewCCIPLane(
			o.SetUpContext, o.Env,
			o.Cfg.TestGroupInput, o.BootstrapAdded, o.JobAddGrp,
		)
		if err != nil {
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("deploying lane %s to %s; err - %w", networkA.Name, networkB.Name, errors.WithStack(err))))
			return err
		}
		err = o.LaneConfig.WriteLaneConfig(networkA.Name, ccipLaneA2B.SrcNetworkLaneCfg)
		if err != nil {
			lggr.Error().Err(err).Msgf("error deploying lane %s to %s", networkA.Name, networkB.Name)
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkA.Name, errors.WithStack(err))))
			return err
		}
		err = o.LaneConfig.WriteLaneConfig(networkB.Name, ccipLaneA2B.DstNetworkLaneCfg)
		if err != nil {
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkB.Name, errors.WithStack(err))))
			return err
		}

		// we need to set the remote chains on the pool after the lane is deployed
		// it's sufficient to do this only for the forward lane, as the destination pools will also be updated with source pool updates
		// The reverse lane will have the same pools as the forward lane but in reverse order of source and destination
		err = ccipLaneA2B.SetRemoteChainsOnPool()
		if err != nil {
			allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("error setting remote chains; err - %w", errors.WithStack(err))))
			return err
		}
		lggr.Info().Msgf("done setting up lane %s to %s", networkA.Name, networkB.Name)
		if o.Cfg.TestGroupInput.LoadProfile != nil && pointer.GetBool(o.Cfg.TestGroupInput.LoadProfile.OptimizeSpace) {
			// This is to optimize memory space for load tests with high number of networks, lanes, tokens
			ccipLaneA2B.OptimizeStorage()
		}
		return nil
	})

	setUpFuncs.Go(func() error {
		if bidirectional {
			lggr.Info().Msgf("Setting up lane %s to %s", networkB.Name, networkA.Name)
			err := ccipLaneB2A.DeployNewCCIPLane(
				o.SetUpContext, o.Env,
				o.Cfg.TestGroupInput, o.BootstrapAdded, o.JobAddGrp,
			)
			if err != nil {
				lggr.Error().Err(err).Msgf("error deploying lane %s to %s", networkB.Name, networkA.Name)
				allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("deploying lane %s to %s; err -  %w", networkB.Name, networkA.Name, errors.WithStack(err))))
				return err
			}

			err = o.LaneConfig.WriteLaneConfig(networkB.Name, ccipLaneB2A.SrcNetworkLaneCfg)
			if err != nil {
				allErrors.Store(multierr.Append(allErrors.Load(), fmt.Errorf("writing lane config for %s; err - %w", networkA.Name, errors.WithStack(err))))
				return err
			}
			err = o.LaneConfig.WriteLaneConfig(networkA.Name, ccipLaneB2A.DstNetworkLaneCfg)
			if err != nil {
				allErrors.Store(
					multierr.Append(
						allErrors.Load(),
						fmt.Errorf("writing lane config for %s; err - %w", networkB.Name, errors.WithStack(err)),
					),
				)
				return err
			}
			lggr.Info().Msgf("done setting up lane %s to %s", networkB.Name, networkA.Name)
			if o.Cfg.TestGroupInput.LoadProfile != nil && pointer.GetBool(o.Cfg.TestGroupInput.LoadProfile.OptimizeSpace) {
				// This is to optimize memory space for load tests with high number of networks, lanes, tokens
				ccipLaneB2A.OptimizeStorage()
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

func (o *CCIPTestSetUpOutputs) StartEventWatchersPolling() {
	for _, lane := range o.ReadLanes() {
		err := lane.ForwardLane.StartEventWatchersPolling(o.SC)
		require.NoError(o.Cfg.Test, err)
		if lane.ReverseLane != nil {
			err = lane.ReverseLane.StartEventWatchersPolling(o.SC)
			require.NoError(o.Cfg.Test, err)
		}
	}
}

func (o *CCIPTestSetUpOutputs) WaitForPriceUpdates() {
	t := o.Cfg.Test
	priceUpdateGrp, _ := errgroup.WithContext(o.SetUpContext)
	for _, lanes := range o.ReadLanes() {
		lanes := lanes
		forwardLane := lanes.ForwardLane
		reverseLane := lanes.ReverseLane
		waitForUpdate := func(lane *actions.CCIPLane) error {
			defer func() {
				lane.Logger.Info().
					Str("source_chain", lane.Source.Common.ChainClient.GetNetworkName()).
					Uint64("dest_chain", lane.Source.DestinationChainId).
					Str("price_registry", lane.Source.Common.PriceRegistry.Address()).
					Msg("Stopping price update watch")

			}()
			var allTokens []common.Address
			for _, token := range lane.Source.Common.BridgeTokens {
				allTokens = append(allTokens, token.ContractAddress)
			}
			allTokens = append(allTokens, lane.Source.Common.FeeToken.EthAddress)
			lane.Logger.Info().
				Str("source_chain", lane.Source.Common.ChainClient.GetNetworkName()).
				Uint64("dest_chain", lane.Source.DestinationChainId).
				Str("price_registry", lane.Source.Common.PriceRegistry.Address()).
				Msgf("Waiting for price update")
			err := lane.Source.Common.WaitForPriceUpdates(
				o.SetUpContext, lane.Logger,
				o.Cfg.TestGroupInput.TokenConfig.TimeoutForPriceUpdate.Duration(),
				lane.Source.DestinationChainId,
				allTokens,
			)
			if err != nil {
				return errors.Wrapf(err, "waiting for price update failed on lane %s-->%s", lane.SourceNetworkName, lane.DestNetworkName)
			}
			return nil
		}

		priceUpdateGrp.Go(func() error {
			return waitForUpdate(forwardLane)
		})
		if lanes.ReverseLane != nil {
			priceUpdateGrp.Go(func() error {
				return waitForUpdate(reverseLane)
			})
		}
	}

	require.NoError(t, priceUpdateGrp.Wait())
}

// CheckGasUpdateTransaction checks the gas update transactions count
func (o *CCIPTestSetUpOutputs) CheckGasUpdateTransaction(lggr *zerolog.Logger) error {
	transactionsBySource := make(map[string]string)
	destToSourcesList := make(map[string][]string)
	// create a map to hold the unique destination with list of sources
	for _, n := range o.Cfg.NetworkPairs {
		destToSourcesList[n.NetworkB.Name] = append(destToSourcesList[n.NetworkB.Name], n.NetworkA.Name)
		if pointer.GetBool(o.Cfg.TestGroupInput.BiDirectionalLane) {
			destToSourcesList[n.NetworkA.Name] = append(destToSourcesList[n.NetworkA.Name], n.NetworkB.Name)
		}
	}
	lggr.Debug().Interface("list", destToSourcesList).Msg("Dest to Source")
	// a function to read the gas update events and create a map with unique source and store the tx hash
	filterGasUpdateEventTxBySource := func(lane *actions.CCIPLane) error {
		for _, g := range lane.Source.Common.GasUpdateEvents {
			if g.Value == nil {
				return fmt.Errorf("gas update value should not be nil in tx %s", g.Tx)
			}
			if _, ok := transactionsBySource[g.Source]; !ok {
				transactionsBySource[g.Source] = g.Tx
			}
		}
		return nil
	}

	for _, lane := range o.ReadLanes() {
		if err := filterGasUpdateEventTxBySource(lane.ForwardLane); err != nil {
			return fmt.Errorf("error in filtering gas update transactions in the lane source: %s and destination: %s, error: %w",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName, err)
		}
		if lane.ReverseLane != nil {
			if err := filterGasUpdateEventTxBySource(lane.ReverseLane); err != nil {
				return fmt.Errorf("error in filtering gas update transactions in the lane source: %s and destination: %s, error: %w",
					lane.ReverseLane.SourceNetworkName, lane.ReverseLane.DestNetworkName, err)
			}
		}
	}

	lggr.Debug().Interface("Tx hashes by source", transactionsBySource).Msg("Checked Gas Update Transactions by Source")
	// when leader lane setup is enabled, number of unique transaction from the source
	// should match the number of leader lanes defined.
	if len(transactionsBySource) != len(o.Cfg.LeaderLanes) {
		lggr.Error().
			Int("Tx hashes expected", len(o.Cfg.LeaderLanes)).
			Int("Tx hashes received", len(transactionsBySource)).
			Int("Leader lanes count", len(o.Cfg.LeaderLanes)).
			Msg("Checked Gas Update transactions count doesn't match")
		return fmt.Errorf("checked Gas Update transactions count doesn't match")
	}
	lggr.Debug().
		Int("Tx hashes by source", len(transactionsBySource)).
		Int("Leader lanes count", len(o.Cfg.LeaderLanes)).
		Msg("Checked Gas Update transactions count matches")

	return nil
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
	lggr *zerolog.Logger,
	envName string,
	tokenDeployerFns []blockchain.ContractDeployer,
	testConfig *CCIPTestConfig,
) *CCIPTestSetUpOutputs {
	var err error
	reportPath := "tmp_laneconfig"
	filepath := fmt.Sprintf("./%s/tmp_%s.json", reportPath, strings.ReplaceAll(t.Name(), "/", "_"))
	reportFile := testutils.FileNameFromPath(filepath)
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	setUpArgs := &CCIPTestSetUpOutputs{
		SetUpContext:           parent,
		Cfg:                    testConfig,
		Reporter:               testreporters.NewCCIPTestReporter(t, lggr),
		LaneConfigFile:         filepath,
		LaneContractsByNetwork: &sync.Map{},
		Balance:                actions.NewBalanceSheet(),
		BootstrapAdded:         atomic.NewBool(false),
		JobAddGrp:              &errgroup.Group{},
		laneMutex:              &sync.Mutex{},
	}

	contractsData, err := setUpArgs.Cfg.ContractsInput.ContractsData()
	require.NoError(t, err, "error reading existing lane config")

	chainClientByChainID := setUpArgs.CreateEnvironment(lggr, envName, reportPath)
	// if test is run in remote runner, register a clean-up to copy the laneconfig file
	if value, set := os.LookupEnv(config.EnvVarJobImage); set && value != "" &&
		(setUpArgs.Env != nil && setUpArgs.Env.K8Env != nil) &&
		pointer.GetBool(setUpArgs.Cfg.TestGroupInput.StoreLaneConfig) {
		t.Cleanup(func() {
			path := fmt.Sprintf("reports/%s/%s", reportPath, reportFile)
			dir, err := os.Getwd()
			require.NoError(t, err)
			destPath := fmt.Sprintf("%s/%s", dir, reportFile)
			lggr.Info().Str("srcPath", path).Str("dstPath", destPath).Msg("copying lane config")
			err = setUpArgs.Env.K8Env.CopyFromPod("app=runner-data",
				"remote-test-runner-data-files", path, destPath)
			require.NoError(t, err, "error getting lane config")
		})
	}
	if setUpArgs.Env != nil {
		ccipEnv := setUpArgs.Env
		if ccipEnv.K8Env != nil && ccipEnv.K8Env.WillUseRemoteRunner() {
			return setUpArgs
		}
	}

	setUpArgs.LaneConfig, err = laneconfig.ReadLanesFromExistingDeployment(contractsData)
	require.NoError(t, err)

	if setUpArgs.LaneConfig == nil {
		setUpArgs.LaneConfig = &laneconfig.Lanes{LaneConfigs: make(map[string]*laneconfig.LaneConfig)}
	}
	laneCfgFile, err := os.Stat(setUpArgs.LaneConfigFile)
	if err == nil && laneCfgFile.Size() > 0 {
		// remove the existing lane config file
		err = os.Remove(setUpArgs.LaneConfigFile)
		require.NoError(t, err, "error while removing existing lane config file - %s", setUpArgs.LaneConfigFile)
	}

	configureCLNode := !testConfig.useExistingDeployment()

	// if no of lanes per pair is greater than 1, copy common contracts from the same network
	// if no of lanes per pair is more than 1, the networks are added into the testConfig.AllNetworks with a suffix of -<lane number>
	// for example, if no of lanes per pair is 2, and the network pairs are called "testnetA", "testnetB",
	//	the network will be added as "testnetA-1", testnetA-2","testnetB-1", testnetB-2"
	// to deploy 2 lanes between same network pair "testnetA", "testnetB".
	// In the following the common contracts will be copied from "testnetA" to "testnetA-1" and "testnetA-2" and
	// from "testnetB" to "testnetB-1" and "testnetB-2"
	for n := range testConfig.AllNetworks {
		if setUpArgs.Cfg.TestGroupInput.NoOfRoutersPerPair > 1 {
			regex := regexp.MustCompile(`-(\d+)$`)
			networkNameToReadCfg := regex.ReplaceAllString(n, "")
			reuse := pointer.GetBool(testConfig.TestGroupInput.ReuseContracts)
			// if reuse contracts is true, copy common contracts from the same network except the router contract
			setUpArgs.LaneConfig.CopyCommonContracts(
				networkNameToReadCfg, n,
				reuse, testConfig.TestGroupInput.MsgDetails.IsTokenTransfer(),
			)
		}
	}

	// deploy all chain specific common contracts
	chainAddGrp, _ := errgroup.WithContext(setUpArgs.SetUpContext)
	lggr.Info().Msg("Deploying common contracts")

	// If we have a token admin registry, we need to use a separate to deploy our test tokens from so that the tokens
	// are not owned by the same account that owns the other CCIP contracts. This emulates self-serve token setups where
	// the token owner is different from the CCIP contract owner.
	if testConfig.useSeparateTokenDeployer() {
		for _, net := range testConfig.AllNetworks {
			chainClient := chainClientByChainID[net.ChainID]
			require.NotNil(t, chainClient, "Chain client not found for chainID %d", net.ChainID)
			require.GreaterOrEqual(t, len(chainClient.GetWallets()), 2, "The test is using a TokenAdminRegistry, and has CCIPOwnerTokens set to 'false'. The test needs a second wallet to deploy token contracts from. Please add a second wallet to the 'evm_clients' config option.")
			tokenDeployerWallet := chainClient.GetWallets()[1]
			// TODO: This is a total guess at how much funds we need to deploy the tokens. This could be way off, especially on live chains.
			// There aren't a lot of good ways to estimate this though. See CCIP-2471.
			recommendedTokenBalance := new(big.Int).Mul(big.NewInt(5e18), big.NewInt(int64(pointer.GetInt(testConfig.TestGroupInput.TokenConfig.NoOfTokensPerChain))))
			currentTokenBalance, err := chainClient.BalanceAt(context.Background(), common.HexToAddress(tokenDeployerWallet.Address()))
			require.NoError(t, err)
			if currentTokenBalance.Cmp(recommendedTokenBalance) < 0 {
				lggr.Warn().
					Str("Token Deployer Address", tokenDeployerWallet.Address()).
					Uint64("Current Balance", currentTokenBalance.Uint64()).
					Uint64("Recommended Balance", recommendedTokenBalance.Uint64()).
					Msg("Token Deployer wallet may be underfunded. Please ensure it has enough funds to deploy the tokens.")
			}
		}
	}

	for _, net := range testConfig.AllNetworks {
		chainClient := chainClientByChainID[net.ChainID]
		net := net
		net.HTTPURLs = chainClient.GetNetworkConfig().HTTPURLs
		net.URLs = chainClient.GetNetworkConfig().URLs
		chainAddGrp.Go(func() error {
			return setUpArgs.DeployChainContracts(
				lggr, chainClient, net,
				pointer.GetInt(testConfig.TestGroupInput.TokenConfig.NoOfTokensPerChain),
				tokenDeployerFns,
			)
		})
	}
	require.NoError(t, chainAddGrp.Wait(), "Deploying common contracts shouldn't fail")

	// set up mock server for price pipeline and usdc attestation if not using existing deployment
	if !pointer.GetBool(setUpArgs.Cfg.TestGroupInput.ExistingDeployment) {
		var killgrave *ctftestenv.Killgrave
		if setUpArgs.Env.LocalCluster != nil {
			killgrave = setUpArgs.Env.LocalCluster.MockAdapter
		}
		if setUpArgs.Cfg.TestGroupInput.TokenConfig.IsPipelineSpec() {
			// set up mock server for price pipeline. need to set it once for all the lanes as the price pipeline path uses
			// regex to match the path for all tokens across all lanes
			actions.SetMockserverWithTokenPriceValue(killgrave, setUpArgs.Env.MockServer)
		}
		if pointer.GetBool(setUpArgs.Cfg.TestGroupInput.USDCMockDeployment) {
			// if it's a new USDC deployment, set up mock server for attestation,
			// we need to set it only once for all the lanes as the attestation path uses regex to match the path for
			// all messages across all lanes
			err = actions.SetMockServerWithUSDCAttestation(killgrave, setUpArgs.Env.MockServer, false)
			require.NoError(t, err, "failed to set up mock server for USDC attestation")
		}
		if pointer.GetBool(setUpArgs.Cfg.TestGroupInput.LBTCMockDeployment) {
			// if it's a new LBTC deployment, set up mock server for attestation,
			// we need to set it only once for all the lanes as the attestation path uses regex to match the path for
			// all messages across all lanes
			err = actions.SetMockServerWithLBTCAttestation(killgrave, setUpArgs.Env.MockServer)
			require.NoError(t, err, "failed to set up mock server for LBTC attestation")
		}
	}
	// deploy all lane specific contracts
	lggr.Info().Msg("Deploying lane specific contracts")
	laneAddGrp, _ := errgroup.WithContext(setUpArgs.SetUpContext)
	// for memory management set a batch size for active lane deployment group
	laneAddGrp.SetLimit(200)
	for _, networkPair := range testConfig.NetworkPairs {
		n := networkPair
		var ok bool
		n.ChainClientA, ok = chainClientByChainID[n.NetworkA.ChainID]
		require.True(t, ok, "Chain client for chainID %d not found", n.NetworkA.ChainID)
		n.ChainClientB, ok = chainClientByChainID[n.NetworkB.ChainID]
		require.True(t, ok, "Chain client for chainID %d not found", n.NetworkB.ChainID)

		n.NetworkA.HTTPURLs = n.ChainClientA.GetNetworkConfig().HTTPURLs
		n.NetworkA.URLs = n.ChainClientA.GetNetworkConfig().URLs
		n.NetworkB.HTTPURLs = n.ChainClientB.GetNetworkConfig().HTTPURLs
		n.NetworkB.URLs = n.ChainClientB.GetNetworkConfig().URLs

		laneAddGrp.Go(func() error {
			return setUpArgs.AddLanesForNetworkPair(
				lggr, n.NetworkA, n.NetworkB,
				chainClientByChainID[n.NetworkA.ChainID], chainClientByChainID[n.NetworkB.ChainID],
			)
		})
	}
	require.NoError(t, laneAddGrp.Wait())
	err = laneconfig.WriteLanesToJSON(setUpArgs.LaneConfigFile, setUpArgs.LaneConfig)
	require.NoError(t, err)

	require.Equal(t, len(setUpArgs.Lanes), len(testConfig.NetworkPairs),
		"Number of bi-directional lanes should be equal to number of network pairs")
	// only required for env set up
	setUpArgs.LaneContractsByNetwork = nil

	if configureCLNode {
		// wait for all jobs to get created
		lggr.Info().Msg("Waiting for jobs to be created")
		require.NoError(t, setUpArgs.JobAddGrp.Wait(), "Creating jobs shouldn't fail")
		// wait for price updates to be available
		setUpArgs.WaitForPriceUpdates()
		if isLeaderLaneFeatureEnabled(lggr) && !pointer.GetBool(setUpArgs.Cfg.TestGroupInput.ExistingDeployment) {
			require.NoError(t, setUpArgs.CheckGasUpdateTransaction(lggr), "gas update transaction check shouldn't fail")
		}
		// if dynamic price update is required
		if setUpArgs.Cfg.TestGroupInput.TokenConfig.IsDynamicPriceUpdate() {
			require.NoError(t, setUpArgs.SetupDynamicTokenPriceUpdates(), "setting up dynamic price update should not fail")
		}
	}

	// start event watchers for all lanes
	if useWebSocket(chainClientByChainID) {
		setUpArgs.StartEventWatchers()
	} else {
		setUpArgs.SC = sentinel.NewSentinelCoordinator(*lggr)
		err := setUpArgs.addChains()
		require.NoError(t, err, "error adding chain to Sentinel")
		setUpArgs.StartEventWatchersPolling()
		t.Cleanup(func() {
			if setUpArgs.SC != nil {
				lggr.Info().Msg("Closing Sentinel")
				setUpArgs.SC.Sentinel.Close()
			}
		})
	}

	// now that lane configs are already dumped to file, we can clean up the lane config map
	setUpArgs.LaneConfig = nil
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

func useWebSocket(chainClientByChainID map[int64]blockchain.EVMClient) bool {
	for _, c := range chainClientByChainID {
		if !c.GetEthClient().Client().SupportsSubscriptions() {
			return false
		}
	}
	return true
}

func (o *CCIPTestSetUpOutputs) addChains() error {
	for _, lane := range o.ReadLanes() {
		// Add both forward and reverse lanes
		err := o.addChainToSentinel(lane.ForwardLane.SourceChain)
		if err != nil {
			return err
		}
		err = o.addChainToSentinel(lane.ForwardLane.DestChain)
		if err != nil {
			return err
		}
	}
	return nil
}

// addChainToSentinel is a helper function to add a chain to Sentinel
func (o *CCIPTestSetUpOutputs) addChainToSentinel(chain blockchain.EVMClient) error {
	blockchainClient := blockchain_client_wrapper.NewGethClientWrapper(chain.GetEthClient())

	// Define the chain poller service configuration
	addChainConfig := sentinel.AddChainConfig{
		ChainID:          chain.GetChainID().Int64(),
		PollInterval:     30 * time.Second,
		BlockchainClient: blockchainClient,
	}

	// Add the chain to Sentinel
	err := o.SC.Sentinel.AddChain(addChainConfig)

	return err
}

// CreateEnvironment creates the environment for the test and registers the test clean-up function to tear down the set-up environment
// It returns the map of chainID to EVMClient
func (o *CCIPTestSetUpOutputs) CreateEnvironment(
	lggr *zerolog.Logger,
	envName string,
	reportPath string,
) map[int64]blockchain.EVMClient {
	var (
		testConfig = o.Cfg
		t          = o.Cfg.Test

		ccipEnv  *actions.CCIPTestEnv
		k8Env    *environment.Environment
		err      error
		chains   []blockchain.EVMClient
		local    *test_env.CLClusterTestEnv
		deployCL func() error
	)

	envConfig := createEnvironmentConfig(t, envName, testConfig, reportPath)

	configureCLNode := !testConfig.useExistingDeployment() || pointer.GetString(testConfig.EnvInput.EnvToConnect) != ""
	namespace := ""
	if testConfig.TestGroupInput.LoadProfile != nil {
		namespace = testConfig.TestGroupInput.LoadProfile.TestRunName
	}
	require.False(t, testConfig.localCluster() && testConfig.ExistingCLCluster(),
		"local cluster and existing cluster cannot be true at the same time")
	// if it's a new deployment, deploy the env
	// Or if EnvToConnect is given connect to that k8 environment
	if configureCLNode {
		if !testConfig.ExistingCLCluster() {
			// if it's a local cluster, deploy the local cluster in docker
			if testConfig.localCluster() {
				local, deployCL = DeployLocalCluster(t, testConfig)
				ccipEnv = &actions.CCIPTestEnv{
					LocalCluster: local,
				}
				namespace = "local-docker-deployment"
			} else {
				// Otherwise, deploy the k8s env
				lggr.Info().Msg("Deploying test environment")
				// deploy the env if configureCLNode is true
				k8Env = DeployEnvironments(t, envConfig, testConfig)
				ccipEnv = &actions.CCIPTestEnv{K8Env: k8Env}
				namespace = ccipEnv.K8Env.Cfg.Namespace
			}
		} else {
			// if there is already a cluster, use the existing cluster to connect to the nodes
			ccipEnv = &actions.CCIPTestEnv{}
			mockserverURL := pointer.GetString(testConfig.EnvInput.Mockserver)
			require.NotEmpty(t, mockserverURL, "mockserver URL cannot be nil")
			ccipEnv.MockServer = ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
				LocalURL:   mockserverURL,
				ClusterURL: mockserverURL,
			})
		}
		ccipEnv.CLNodeWithKeyReady, _ = errgroup.WithContext(o.SetUpContext)
		o.Env = ccipEnv
		if ccipEnv.K8Env != nil && ccipEnv.K8Env.WillUseRemoteRunner() {
			return nil
		}
	} else {
		// if configureCLNode is false it means we don't need to deploy any additional pods,
		// use a placeholder env to create just the remote runner in it.
		if value, set := os.LookupEnv(config.EnvVarJobImage); set && value != "" {
			k8Env = environment.New(envConfig)
			err = k8Env.Run()
			require.NoErrorf(t, err, "error creating environment remote runner")
			o.Env = &actions.CCIPTestEnv{K8Env: k8Env}
			if k8Env.WillUseRemoteRunner() {
				return nil
			}
		}
	}
	if o.Cfg.TestGroupInput.LoadProfile != nil {
		o.Cfg.TestGroupInput.LoadProfile.SetTestRunName(namespace)
	}
	chainByChainID := make(map[int64]blockchain.EVMClient)
	if pointer.GetBool(testConfig.TestGroupInput.LocalCluster) {
		require.NotNil(t, ccipEnv.LocalCluster, "Local cluster shouldn't be nil")
		for _, n := range ccipEnv.LocalCluster.EVMNetworks {
			if evmClient, err := blockchain.NewEVMClientFromNetwork(*n, *lggr); err == nil {
				chainByChainID[evmClient.GetChainID().Int64()] = evmClient
				chains = append(chains, evmClient)
			} else {
				lggr.Error().Err(err).Msgf("EVMClient for chainID %d not found", n.ChainID)
			}
		}
	} else {
		for _, n := range testConfig.SelectedNetworks {
			if _, ok := chainByChainID[n.ChainID]; ok {
				continue
			}
			var ec blockchain.EVMClient
			if k8Env == nil {
				ec, err = blockchain.ConnectEVMClient(n, *lggr)
			} else {
				log.Info().Interface("urls", k8Env.URLs).Msg("URLs")
				ec, err = blockchain.NewEVMClient(n, k8Env, *lggr)
			}
			require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
			chains = append(chains, ec)
			chainByChainID[n.ChainID] = ec
		}
	}
	if configureCLNode {
		ccipEnv.CLNodeWithKeyReady.Go(func() error {
			var totalNodes int
			if !o.Cfg.ExistingCLCluster() {
				if ccipEnv.LocalCluster != nil {
					err = deployCL()
					if err != nil {
						return err
					}
				}
				err = ccipEnv.ConnectToDeployedNodes()
				if err != nil {
					return fmt.Errorf("error connecting to chainlink nodes: %w", err)
				}
				totalNodes = pointer.GetInt(testConfig.EnvInput.NewCLCluster.NoOfNodes)
			} else {
				totalNodes = pointer.GetInt(testConfig.EnvInput.ExistingCLCluster.NoOfNodes)
				err = ccipEnv.ConnectToExistingNodes(o.Cfg.EnvInput)
				if err != nil {
					return fmt.Errorf("error deploying and connecting to chainlink nodes: %w", err)
				}
			}
			err = ccipEnv.SetUpNodeKeysAndFund(lggr, big.NewFloat(testConfig.TestGroupInput.NodeFunding), chains)
			if err != nil {
				return fmt.Errorf("error setting up nodes and keys %w", err)
			}
			// first node is the bootstrapper
			ccipEnv.CommitNodeStartIndex = 1
			ccipEnv.ExecNodeStartIndex = 1
			ccipEnv.NumOfCommitNodes = testConfig.TestGroupInput.NoOfCommitNodes
			ccipEnv.NumOfExecNodes = ccipEnv.NumOfCommitNodes
			if !pointer.GetBool(testConfig.TestGroupInput.CommitAndExecuteOnSameDON) {
				if len(ccipEnv.CLNodesWithKeys) < 11 {
					return fmt.Errorf("not enough CL nodes for separate commit and execution nodes")
				}
				if testConfig.TestGroupInput.NoOfCommitNodes >= totalNodes {
					return fmt.Errorf("number of commit nodes can not be greater than total number of nodes in DON")
				}
				// first two nodes are reserved for bootstrap commit and bootstrap exec
				ccipEnv.CommitNodeStartIndex = 2
				ccipEnv.ExecNodeStartIndex = 2 + testConfig.TestGroupInput.NoOfCommitNodes
				ccipEnv.NumOfExecNodes = totalNodes - (2 + testConfig.TestGroupInput.NoOfCommitNodes)
				if ccipEnv.NumOfExecNodes < 4 {
					return fmt.Errorf("insufficient number of exec nodes")
				}
			}
			ccipEnv.NumOfAllowedFaultyExec = (ccipEnv.NumOfExecNodes - 1) / 3
			ccipEnv.NumOfAllowedFaultyCommit = (ccipEnv.NumOfCommitNodes - 1) / 3
			return nil
		})
	}

	t.Cleanup(func() {
		if configureCLNode {
			if ccipEnv.LocalCluster != nil {
				err := ccipEnv.LocalCluster.Terminate()
				require.NoError(t, err, "Local cluster termination shouldn't fail")
				require.NoError(t, o.Reporter.SendReport(t, namespace, false), "Aggregating and sending report shouldn't fail")
				return
			}
			if pointer.GetBool(testConfig.TestGroupInput.KeepEnvAlive) || testConfig.ExistingCLCluster() {
				require.NoError(t, o.Reporter.SendReport(t, namespace, true), "Aggregating and sending report shouldn't fail")
				return
			}
			lggr.Info().Msg("Tearing down the environment")
			err = integrationactions.TeardownSuite(t, nil, ccipEnv.K8Env, ccipEnv.CLNodes, o.Reporter, zapcore.DPanicLevel, o.Cfg.EnvInput)
			require.NoError(t, err, "Environment teardown shouldn't fail")
		} else {
			//just send the report
			require.NoError(t, o.Reporter.SendReport(t, namespace, true), "Aggregating and sending report shouldn't fail")
		}
	})
	return chainByChainID
}

func createEnvironmentConfig(t *testing.T, envName string, testConfig *CCIPTestConfig, reportPath string) *environment.Config {
	testType := testConfig.TestGroupInput.Type
	nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels(string(tc.CCIP), testType)
	require.NoError(t, err, "Error creating required chain.link labels for namespace")

	workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels(string(tc.CCIP), testType)
	require.NoError(t, err, "Error creating required chain.link labels for workloads and pods")

	envConfig := &environment.Config{
		NamespacePrefix: envName,
		Test:            t,
		Labels:          nsLabels,
		WorkloadLabels:  workloadPodLabels,
		PodLabels:       workloadPodLabels,
		//	PreventPodEviction: true, //TODO: enable this once we have a way to handle pod eviction
	}
	if pointer.GetBool(testConfig.TestGroupInput.StoreLaneConfig) {
		envConfig.ReportPath = reportPath
	}
	// if there is already existing namespace, no need to update any manifest there, we just connect to it
	existingEnv := pointer.GetString(testConfig.EnvInput.EnvToConnect)
	if existingEnv != "" {
		envConfig.Namespace = existingEnv
		envConfig.NamespacePrefix = ""
		envConfig.SkipManifestUpdate = true
		envConfig.RunnerName = fmt.Sprintf("%s-%s", environment.REMOTE_RUNNER_NAME, uuid.NewString()[0:5])
	}
	if testConfig.EnvInput.TTL != nil {
		envConfig.TTL = testConfig.EnvInput.TTL.Duration()
	}
	if testConfig.TestGroupInput.LoadProfile != nil && testConfig.TestGroupInput.LoadProfile.TestDuration != nil {
		approxDur := testConfig.TestGroupInput.LoadProfile.TestDuration.Duration() + 3*time.Hour
		if envConfig.TTL < approxDur {
			envConfig.TTL = approxDur
		}
	}
	return envConfig
}
