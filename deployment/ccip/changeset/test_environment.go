package changeset

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

type EnvType string

const (
	Memory      EnvType = "in-memory"
	Docker      EnvType = "docker"
	ENVTESTTYPE         = "CCIP_V16_TEST_ENV"
)

type TestConfigs struct {
	Type      EnvType // set by env var CCIP_V16_TEST_ENV, defaults to Memory
	CreateJob bool
	// TODO: This should be CreateContracts so the booleans make sense?
	CreateJobAndContracts    bool
	Chains                   int // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	NumOfUsersPerChain       int // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	Nodes                    int // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	Bootstraps               int // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	IsUSDC                   bool
	IsUSDCAttestationMissing bool
	IsMultiCall3             bool
	OCRConfigOverride        func(CCIPOCRParams) CCIPOCRParams
	RMNEnabled               bool
	NumOfRMNNodes            int
	LinkPrice                *big.Int
	WethPrice                *big.Int
}

func (tc *TestConfigs) Validate() error {
	if tc.Chains < 2 {
		return fmt.Errorf("chains must be at least 2")
	}
	if tc.Nodes < 4 {
		return fmt.Errorf("nodes must be at least 4")
	}
	if tc.Bootstraps < 1 {
		return fmt.Errorf("bootstraps must be at least 1")
	}
	if tc.Type == Memory && tc.RMNEnabled {
		return fmt.Errorf("cannot run RMN tests in memory mode")
	}
	return nil
}

func (tc *TestConfigs) MustSetEnvTypeOrDefault(t *testing.T) {
	envType := os.Getenv(ENVTESTTYPE)
	if envType == "" || envType == string(Memory) {
		tc.Type = Memory
	} else if envType == string(Docker) {
		tc.Type = Docker
	} else {
		t.Fatalf("env var CCIP_V16_TEST_ENV must be either %s or %s, defaults to %s if unset, got: %s", Memory, Docker, Memory, envType)
	}
}

func DefaultTestConfigs() *TestConfigs {
	return &TestConfigs{
		Chains:                2,
		NumOfUsersPerChain:    1,
		Nodes:                 4,
		Bootstraps:            1,
		LinkPrice:             MockLinkPrice,
		WethPrice:             MockWethPrice,
		CreateJobAndContracts: true,
	}
}

type TestOps func(testCfg *TestConfigs)

func WithMultiCall3() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsMultiCall3 = true
	}
}

func WithJobsOnly() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.CreateJobAndContracts = false
		testCfg.CreateJob = true
	}
}

func WithNoJobsAndContracts() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.CreateJobAndContracts = false
		testCfg.CreateJob = false
	}
}

func WithRMNEnabled(numOfNode int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.RMNEnabled = true
		testCfg.NumOfRMNNodes = numOfNode
	}
}

func WithOCRConfigOverride(override func(CCIPOCRParams) CCIPOCRParams) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.OCRConfigOverride = override
	}
}

func WithUSDCAttestationMissing() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsUSDCAttestationMissing = true
	}
}

func WithUSDC() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsUSDC = true
	}
}

func WithChains(numChains int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Chains = numChains
	}
}

func WithUsersPerChain(numUsers int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.NumOfUsersPerChain = numUsers
	}
}

func WithNodes(numNodes int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Nodes = numNodes
	}
}

func WithBootstraps(numBootstraps int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Bootstraps = numBootstraps
	}
}

type TestEnvironment interface {
	SetupJobs(t *testing.T)
	StartNodes(t *testing.T, tc *TestConfigs, crConfig deployment.CapabilityRegistryConfig)
	StartChains(t *testing.T, tc *TestConfigs)
	DeployedEnvironment() DeployedEnv
	MockUSDCAttestationServer(t *testing.T, isUSDCAttestationMissing bool) string
}

type DeployedEnv struct {
	Env          deployment.Environment
	HomeChainSel uint64
	FeedChainSel uint64
	ReplayBlocks map[uint64]uint64
	Users        map[uint64][]*bind.TransactOpts
}

func (d *DeployedEnv) TimelockContracts(t *testing.T) map[uint64]*proposalutils.TimelockExecutionContracts {
	timelocks := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	state, err := LoadOnchainState(d.Env)
	require.NoError(t, err)
	for chain, chainState := range state.Chains {
		timelocks[chain] = &proposalutils.TimelockExecutionContracts{
			Timelock:  chainState.Timelock,
			CallProxy: chainState.CallProxy,
		}
	}
	return timelocks
}

func (d *DeployedEnv) SetupJobs(t *testing.T) {
	ctx := testcontext.Get(t)
	out, err := CCIPCapabilityJobspec(d.Env, struct{}{})
	require.NoError(t, err)
	for nodeID, jobs := range out.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := d.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}
	// Wait for plugins to register filters?
	// TODO: Investigate how to avoid.
	time.Sleep(30 * time.Second)
	ReplayLogs(t, d.Env.Offchain, d.ReplayBlocks)
}

type MemoryEnvironment struct {
	DeployedEnv
	chains map[uint64]deployment.Chain
}

func (m *MemoryEnvironment) DeployedEnvironment() DeployedEnv {
	return m.DeployedEnv
}

func (m *MemoryEnvironment) StartChains(t *testing.T, tc *TestConfigs) {
	ctx := testcontext.Get(t)
	chains, users := memory.NewMemoryChains(t, tc.Chains, tc.NumOfUsersPerChain)
	m.chains = chains
	homeChainSel, feedSel := allocateCCIPChainSelectors(chains)
	replayBlocks, err := LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)
	m.DeployedEnv = DeployedEnv{
		Env: deployment.Environment{
			Chains: m.chains,
		},
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
		Users:        users,
	}
}

func (m *MemoryEnvironment) StartNodes(t *testing.T, tc *TestConfigs, crConfig deployment.CapabilityRegistryConfig) {
	require.NotNil(t, m.chains, "start chains first, chains are empty")
	require.NotNil(t, m.DeployedEnv, "start chains and initiate deployed env first before starting nodes")
	nodes := memory.NewNodes(t, zapcore.InfoLevel, m.chains, tc.Nodes, tc.Bootstraps, crConfig)
	ctx := testcontext.Get(t)
	lggr := logger.Test(t)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}
	m.DeployedEnv.Env = memory.NewMemoryEnvironmentFromChainsNodes(func() context.Context { return ctx }, lggr, m.chains, nodes)
}

func (m *MemoryEnvironment) MockUSDCAttestationServer(t *testing.T, isUSDCAttestationMissing bool) string {
	server := mockAttestationResponse(isUSDCAttestationMissing)
	endpoint := server.URL
	t.Cleanup(func() {
		server.Close()
	})
	return endpoint
}

// NewMemoryEnvironment creates an in-memory environment based on the testconfig requested
func NewMemoryEnvironment(t *testing.T, opts ...TestOps) DeployedEnv {
	testCfg := DefaultTestConfigs()
	for _, opt := range opts {
		opt(testCfg)
	}
	require.NoError(t, testCfg.Validate(), "invalid test config")
	env := &MemoryEnvironment{}
	if testCfg.CreateJobAndContracts {
		return NewEnvironmentWithJobsAndContracts(t, testCfg, env)
	}
	if testCfg.CreateJob {
		return NewEnvironmentWithJobs(t, testCfg, env)
	}
	return NewEnvironment(t, testCfg, env)
}

func NewEnvironment(t *testing.T, tc *TestConfigs, tEnv TestEnvironment) DeployedEnv {
	lggr := logger.Test(t)
	tEnv.StartChains(t, tc)
	dEnv := tEnv.DeployedEnvironment()
	require.NotEmpty(t, dEnv.FeedChainSel)
	require.NotEmpty(t, dEnv.HomeChainSel)
	require.NotEmpty(t, dEnv.Env.Chains)
	ab := deployment.NewMemoryAddressBook()
	crConfig := DeployTestContracts(t, lggr, ab, dEnv.HomeChainSel, dEnv.FeedChainSel, dEnv.Env.Chains, tc.LinkPrice, tc.WethPrice)
	tEnv.StartNodes(t, tc, crConfig)
	dEnv = tEnv.DeployedEnvironment()
	// TODO: Should use ApplyChangesets here.
	envNodes, err := deployment.NodeInfo(dEnv.Env.NodeIDs, dEnv.Env.Offchain)
	require.NoError(t, err)
	dEnv.Env.ExistingAddresses = ab
	_, err = deployHomeChain(lggr, dEnv.Env, dEnv.Env.ExistingAddresses, dEnv.Env.Chains[dEnv.HomeChainSel],
		NewTestRMNStaticConfig(),
		NewTestRMNDynamicConfig(),
		NewTestNodeOperator(dEnv.Env.Chains[dEnv.HomeChainSel].DeployerKey.From),
		map[string][][32]byte{
			"NodeOperator": envNodes.NonBootstraps().PeerIDs(),
		},
	)
	require.NoError(t, err)

	return dEnv
}

func NewEnvironmentWithJobsAndContracts(t *testing.T, tc *TestConfigs, tEnv TestEnvironment) DeployedEnv {
	var err error
	e := NewEnvironment(t, tc, tEnv)
	allChains := e.Env.AllChainSelectors()
	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)

	for _, c := range e.Env.AllChainSelectors() {
		mcmsCfg[c] = proposalutils.SingleGroupTimelockConfig(t)
	}
	var (
		usdcChains   []uint64
		isMulticall3 bool
	)
	if tc != nil {
		if tc.IsUSDC {
			usdcChains = allChains
		}
		isMulticall3 = tc.IsMultiCall3
	}
	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    allChains,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployPrerequisites),
			Config: DeployPrerequisiteConfig{
				ChainSelectors: allChains,
				Opts: []PrerequisiteOpt{
					WithUSDCChains(usdcChains),
					WithMulticall3(isMulticall3),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    allChains,
				HomeChainSelector: e.HomeChainSel,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(SetRMNRemoteOnRMNProxy),
			Config: SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: allChains,
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)
	// Assert USDC set up as expected.
	for _, chain := range usdcChains {
		require.NotNil(t, state.Chains[chain].MockUSDCTokenMessenger)
		require.NotNil(t, state.Chains[chain].MockUSDCTransmitter)
		require.NotNil(t, state.Chains[chain].USDCTokenPool)
	}
	// Assert link present
	require.NotNil(t, state.Chains[e.FeedChainSel].LinkToken)
	require.NotNil(t, state.Chains[e.FeedChainSel].Weth9)

	tokenConfig := NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	var tokenDataProviders []pluginconfig.TokenDataObserverConfig
	if len(usdcChains) > 0 {
		endpoint := tEnv.MockUSDCAttestationServer(t, tc.IsUSDCAttestationMissing)
		cctpContracts := make(map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig)
		for _, usdcChain := range usdcChains {
			cctpContracts[cciptypes.ChainSelector(usdcChain)] = pluginconfig.USDCCCTPTokenConfig{
				SourcePoolAddress:            state.Chains[usdcChain].USDCTokenPool.Address().String(),
				SourceMessageTransmitterAddr: state.Chains[usdcChain].MockUSDCTransmitter.Address().String(),
			}
		}
		tokenDataProviders = append(tokenDataProviders, pluginconfig.TokenDataObserverConfig{
			Type:    pluginconfig.USDCCCTPHandlerType,
			Version: "1.0",
			USDCCCTPObserverConfig: &pluginconfig.USDCCCTPObserverConfig{
				Tokens:                 cctpContracts,
				AttestationAPI:         endpoint,
				AttestationAPITimeout:  commonconfig.MustNewDuration(time.Second),
				AttestationAPIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
			}})
	}
	// Build the per chain config.
	ocrConfigs := make(map[uint64]CCIPOCRParams)
	chainConfigs := make(map[uint64]ChainConfig)
	timelockContractsPerChain := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	nodeInfo, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	for _, chain := range allChains {
		timelockContractsPerChain[chain] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain].Timelock,
			CallProxy: state.Chains[chain].CallProxy,
		}
		tokenInfo := tokenConfig.GetTokenInfo(e.Env.Logger, state.Chains[chain].LinkToken, state.Chains[chain].Weth9)
		ocrParams := DefaultOCRParams(e.FeedChainSel, tokenInfo, tokenDataProviders)
		if tc.OCRConfigOverride != nil {
			ocrParams = tc.OCRConfigOverride(ocrParams)
		}
		ocrConfigs[chain] = ocrParams
		chainConfigs[chain] = ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			FChain:  uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(internal.GasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(internal.DAGasPriceDeviationPPB)},
				OptimisticConfirmations: internal.OptimisticConfirmations,
			},
		}
	}
	// Deploy second set of changesets to deploy and configure the CCIP contracts.
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, timelockContractsPerChain, []commonchangeset.ChangesetApplication{
		{
			// Add the chain configs for the new chains.
			Changeset: commonchangeset.WrapChangeSet(UpdateChainConfig),
			Config: UpdateChainConfigConfig{
				HomeChainSelector: e.HomeChainSel,
				RemoteChainAdds:   chainConfigs,
			},
		},
		{
			// Add the DONs and candidate commit OCR instances for the chain.
			Changeset: commonchangeset.WrapChangeSet(AddDonAndSetCandidateChangeset),
			Config: AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase{
					HomeChainSelector:               e.HomeChainSel,
					FeedChainSelector:               e.FeedChainSel,
					OCRConfigPerRemoteChainSelector: ocrConfigs,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		},
		{
			// Add the exec OCR instances for the new chains.
			Changeset: commonchangeset.WrapChangeSet(SetCandidateChangeset),
			Config: SetCandidateChangesetConfig{
				SetCandidateConfigBase{
					HomeChainSelector:               e.HomeChainSel,
					FeedChainSelector:               e.FeedChainSel,
					OCRConfigPerRemoteChainSelector: ocrConfigs,
					PluginType:                      types.PluginTypeCCIPExec,
				},
			},
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
			Config: PromoteAllCandidatesChangesetConfig{
				HomeChainSelector:    e.HomeChainSel,
				RemoteChainSelectors: allChains,
				PluginType:           types.PluginTypeCCIPCommit,
			},
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
			Config: PromoteAllCandidatesChangesetConfig{
				HomeChainSelector:    e.HomeChainSel,
				RemoteChainSelectors: allChains,
				PluginType:           types.PluginTypeCCIPExec,
			},
		},
		{
			// Enable the OCR config on the remote chains.
			Changeset: commonchangeset.WrapChangeSet(SetOCR3OffRamp),
			Config: SetOCR3OffRampConfig{
				HomeChainSel:    e.HomeChainSel,
				RemoteChainSels: allChains,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(CCIPCapabilityJobspec),
		},
	})
	require.NoError(t, err)

	ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[e.HomeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[e.HomeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[e.HomeChainSel].RMNHome)
	for _, chain := range allChains {
		require.NotNil(t, state.Chains[chain].LinkToken)
		require.NotNil(t, state.Chains[chain].Weth9)
		require.NotNil(t, state.Chains[chain].TokenAdminRegistry)
		require.NotNil(t, state.Chains[chain].RegistryModule)
		require.NotNil(t, state.Chains[chain].Router)
		require.NotNil(t, state.Chains[chain].RMNRemote)
		require.NotNil(t, state.Chains[chain].TestRouter)
		require.NotNil(t, state.Chains[chain].NonceManager)
		require.NotNil(t, state.Chains[chain].FeeQuoter)
		require.NotNil(t, state.Chains[chain].OffRamp)
		require.NotNil(t, state.Chains[chain].OnRamp)
	}
	return e
}

// NewEnvironmentWithJobs creates a new CCIP environment
// with capreg, fee tokens, feeds, nodes and jobs set up.
func NewEnvironmentWithJobs(t *testing.T, tc *TestConfigs, tEnv TestEnvironment) DeployedEnv {
	e := NewEnvironment(t, tc, tEnv)
	e.SetupJobs(t)
	return e
}
