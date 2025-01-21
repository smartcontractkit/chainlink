package testhelpers

import (
	"context"
	"errors"
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

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
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
	CreateJobAndContracts      bool
	PrerequisiteDeploymentOnly bool
	V1_5Cfg                    changeset.V1_5DeploymentConfig
	Chains                     int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	ChainIDs                   []uint64 // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	NumOfUsersPerChain         int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	Nodes                      int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	Bootstraps                 int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	IsUSDC                     bool
	IsUSDCAttestationMissing   bool
	IsMultiCall3               bool
	OCRConfigOverride          func(*changeset.CCIPOCRParams)
	RMNEnabled                 bool
	NumOfRMNNodes              int
	LinkPrice                  *big.Int
	WethPrice                  *big.Int
}

func (tc *TestConfigs) Validate() error {
	if tc.Chains < 2 {
		return errors.New("chains must be at least 2")
	}
	if tc.Nodes < 4 {
		return errors.New("nodes must be at least 4")
	}
	if tc.Bootstraps < 1 {
		return errors.New("bootstraps must be at least 1")
	}
	if tc.Type == Memory && tc.RMNEnabled {
		return errors.New("cannot run RMN tests in memory mode")
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
		LinkPrice:             changeset.MockLinkPrice,
		WethPrice:             changeset.MockWethPrice,
		CreateJobAndContracts: true,
	}
}

type TestOps func(testCfg *TestConfigs)

func WithMultiCall3() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsMultiCall3 = true
	}
}

func WithPrerequisiteDeploymentOnly(v1_5Cfg *changeset.V1_5DeploymentConfig) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.PrerequisiteDeploymentOnly = true
		if v1_5Cfg != nil {
			testCfg.V1_5Cfg = *v1_5Cfg
		}
	}
}

func WithChainIDs(chainIDs []uint64) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.ChainIDs = chainIDs
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

func WithOCRConfigOverride(override func(*changeset.CCIPOCRParams)) TestOps {
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

func WithNumOfChains(numChains int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Chains = numChains
	}
}

func WithNumOfUsersPerChain(numUsers int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.NumOfUsersPerChain = numUsers
	}
}

func WithNumOfNodes(numNodes int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Nodes = numNodes
	}
}

func WithNumOfBootstrapNodes(numBootstraps int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Bootstraps = numBootstraps
	}
}

type TestEnvironment interface {
	SetupJobs(t *testing.T)
	StartNodes(t *testing.T, crConfig deployment.CapabilityRegistryConfig)
	StartChains(t *testing.T)
	TestConfigs() *TestConfigs
	DeployedEnvironment() DeployedEnv
	UpdateDeployedEnvironment(env DeployedEnv)
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
	state, err := changeset.LoadOnchainState(d.Env)
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
	out, err := changeset.CCIPCapabilityJobspecChangeset(d.Env, struct{}{})
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
	TestConfig *TestConfigs
	Chains     map[uint64]deployment.Chain
}

func (m *MemoryEnvironment) TestConfigs() *TestConfigs {
	return m.TestConfig
}

func (m *MemoryEnvironment) DeployedEnvironment() DeployedEnv {
	return m.DeployedEnv
}

func (m *MemoryEnvironment) UpdateDeployedEnvironment(env DeployedEnv) {
	m.DeployedEnv = env
}

func (m *MemoryEnvironment) StartChains(t *testing.T) {
	ctx := testcontext.Get(t)
	tc := m.TestConfig
	var chains map[uint64]deployment.Chain
	var users map[uint64][]*bind.TransactOpts
	if len(tc.ChainIDs) > 0 {
		chains, users = memory.NewMemoryChainsWithChainIDs(t, tc.ChainIDs, tc.NumOfUsersPerChain)
		if tc.Chains > len(tc.ChainIDs) {
			additionalChains, additionalUsers := memory.NewMemoryChains(t, tc.Chains-len(tc.ChainIDs), tc.NumOfUsersPerChain)
			for k, v := range additionalChains {
				chains[k] = v
			}
			for k, v := range additionalUsers {
				users[k] = v
			}
		}
	} else {
		chains, users = memory.NewMemoryChains(t, tc.Chains, tc.NumOfUsersPerChain)
	}
	m.Chains = chains
	homeChainSel, feedSel := allocateCCIPChainSelectors(chains)
	replayBlocks, err := LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)
	m.DeployedEnv = DeployedEnv{
		Env: deployment.Environment{
			Chains: m.Chains,
		},
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
		Users:        users,
	}
}

func (m *MemoryEnvironment) StartNodes(t *testing.T, crConfig deployment.CapabilityRegistryConfig) {
	require.NotNil(t, m.Chains, "start chains first, chains are empty")
	require.NotNil(t, m.DeployedEnv, "start chains and initiate deployed env first before starting nodes")
	tc := m.TestConfig
	nodes := memory.NewNodes(t, zapcore.InfoLevel, m.Chains, tc.Nodes, tc.Bootstraps, crConfig)
	ctx := testcontext.Get(t)
	lggr := logger.Test(t)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}
	m.DeployedEnv.Env = memory.NewMemoryEnvironmentFromChainsNodes(func() context.Context { return ctx }, lggr, m.Chains, nodes)
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
func NewMemoryEnvironment(t *testing.T, opts ...TestOps) (DeployedEnv, TestEnvironment) {
	testCfg := DefaultTestConfigs()
	for _, opt := range opts {
		opt(testCfg)
	}
	require.NoError(t, testCfg.Validate(), "invalid test config")
	env := &MemoryEnvironment{
		TestConfig: testCfg,
	}
	if testCfg.PrerequisiteDeploymentOnly {
		dEnv := NewEnvironmentWithPrerequisitesContracts(t, env)
		env.UpdateDeployedEnvironment(dEnv)
		return dEnv, env
	}
	if testCfg.CreateJobAndContracts {
		dEnv := NewEnvironmentWithJobsAndContracts(t, env)
		env.UpdateDeployedEnvironment(dEnv)
		return dEnv, env
	}
	if testCfg.CreateJob {
		dEnv := NewEnvironmentWithJobs(t, env)
		env.UpdateDeployedEnvironment(dEnv)
		return dEnv, env
	}
	dEnv := NewEnvironment(t, env)
	env.UpdateDeployedEnvironment(dEnv)
	return dEnv, env
}

func NewEnvironmentWithPrerequisitesContracts(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	var err error
	tc := tEnv.TestConfigs()
	e := NewEnvironment(t, tEnv)
	allChains := e.Env.AllChainSelectors()

	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, c := range e.Env.AllChainSelectors() {
		mcmsCfg[c] = proposalutils.SingleGroupTimelockConfig(t)
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range allChains {
		var opts []changeset.PrerequisiteOpt
		if tc != nil {
			if tc.IsUSDC {
				opts = append(opts, changeset.WithUSDCEnabled())
			}
			if tc.IsMultiCall3 {
				opts = append(opts, changeset.WithMultiCall3Enabled())
			}
		}
		if tc.V1_5Cfg != (changeset.V1_5DeploymentConfig{}) {
			opts = append(opts, changeset.WithLegacyDeploymentEnabled(tc.V1_5Cfg))
		}
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
			Opts:          opts,
		})
	}

	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    allChains,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)
	return e
}

func NewEnvironment(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	lggr := logger.Test(t)
	tc := tEnv.TestConfigs()
	tEnv.StartChains(t)
	dEnv := tEnv.DeployedEnvironment()
	require.NotEmpty(t, dEnv.FeedChainSel)
	require.NotEmpty(t, dEnv.HomeChainSel)
	require.NotEmpty(t, dEnv.Env.Chains)
	ab := deployment.NewMemoryAddressBook()
	crConfig := DeployTestContracts(t, lggr, ab, dEnv.HomeChainSel, dEnv.FeedChainSel, dEnv.Env.Chains, tc.LinkPrice, tc.WethPrice)
	tEnv.StartNodes(t, crConfig)
	dEnv = tEnv.DeployedEnvironment()
	dEnv.Env.ExistingAddresses = ab
	return dEnv
}

func NewEnvironmentWithJobsAndContracts(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	var err error
	tc := tEnv.TestConfigs()
	e := NewEnvironment(t, tEnv)
	allChains := e.Env.AllChainSelectors()
	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)

	for _, c := range e.Env.AllChainSelectors() {
		mcmsCfg[c] = proposalutils.SingleGroupTimelockConfig(t)
	}

	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range allChains {
		var opts []changeset.PrerequisiteOpt
		if tc != nil {
			if tc.IsUSDC {
				opts = append(opts, changeset.WithUSDCEnabled())
			}
			if tc.IsMultiCall3 {
				opts = append(opts, changeset.WithMultiCall3Enabled())
			}
		}
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
			Opts:          opts,
		})
	}
	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    allChains,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)
	e = AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tEnv, true, true, false)
	// now we update RMNProxy to point to RMNRemote
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetRMNRemoteOnRMNProxyChangeset),
			Config: changeset.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: allChains,
			},
		},
	})
	require.NoError(t, err)
	return e
}

func AddCCIPContractsToEnvironment(t *testing.T, allChains []uint64, tEnv TestEnvironment, deployJobs, deployHomeChain, mcmsEnabled bool) DeployedEnv {
	tc := tEnv.TestConfigs()
	e := tEnv.DeployedEnvironment()
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	var apps []commonchangeset.ChangesetApplication
	if deployHomeChain {
		apps = append(apps, commonchangeset.ChangesetApplication{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				NodeOperators:    NewTestNodeOperator(e.Env.Chains[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			},
		})
	}
	apps = append(apps, commonchangeset.ChangesetApplication{
		Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
		Config: changeset.DeployChainContractsConfig{
			ChainSelectors:    allChains,
			HomeChainSelector: e.HomeChainSel,
		},
	})
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, apps)
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	// Assert link present
	require.NotNil(t, state.Chains[e.FeedChainSel].LinkToken)
	require.NotNil(t, state.Chains[e.FeedChainSel].Weth9)

	tokenConfig := changeset.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	var tokenDataProviders []pluginconfig.TokenDataObserverConfig
	if tc.IsUSDC {
		endpoint := tEnv.MockUSDCAttestationServer(t, tc.IsUSDCAttestationMissing)
		cctpContracts := make(map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig)
		for _, usdcChain := range allChains {
			require.NotNil(t, state.Chains[usdcChain].MockUSDCTokenMessenger)
			require.NotNil(t, state.Chains[usdcChain].MockUSDCTransmitter)
			require.NotNil(t, state.Chains[usdcChain].USDCTokenPool)
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
	ocrConfigs := make(map[uint64]changeset.CCIPOCRParams)
	chainConfigs := make(map[uint64]changeset.ChainConfig)
	timelockContractsPerChain := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	nodeInfo, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	for _, chain := range allChains {
		timelockContractsPerChain[chain] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain].Timelock,
			CallProxy: state.Chains[chain].CallProxy,
		}
		tokenInfo := tokenConfig.GetTokenInfo(e.Env.Logger, state.Chains[chain].LinkToken, state.Chains[chain].Weth9)
		ocrOverride := tc.OCRConfigOverride
		if tc.RMNEnabled {
			ocrOverride = func(ocrParams *changeset.CCIPOCRParams) {
				if tc.OCRConfigOverride != nil {
					tc.OCRConfigOverride(ocrParams)
				}
				ocrParams.CommitOffChainConfig.RMNEnabled = true
			}
		}
		ocrParams := changeset.DeriveCCIPOCRParams(
			changeset.WithDefaultCommitOffChainConfig(e.FeedChainSel, tokenInfo),
			changeset.WithDefaultExecuteOffChainConfig(tokenDataProviders),
			changeset.WithOCRParamOverride(ocrOverride),
		)
		ocrConfigs[chain] = ocrParams
		chainConfigs[chain] = changeset.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			FChain:  uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(internal.GasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(internal.DAGasPriceDeviationPPB)},
				OptimisticConfirmations: internal.OptimisticConfirmations,
			},
		}
	}
	timelockContractsPerChain[e.HomeChainSel] = &proposalutils.TimelockExecutionContracts{
		Timelock:  state.Chains[e.HomeChainSel].Timelock,
		CallProxy: state.Chains[e.HomeChainSel].CallProxy,
	}
	// Apply second set of changesets to configure the CCIP contracts.
	var mcmsConfig *changeset.MCMSConfig
	if mcmsEnabled {
		mcmsConfig = &changeset.MCMSConfig{
			MinDelay: 0,
		}
	}
	apps = []commonchangeset.ChangesetApplication{
		{
			// Add the chain configs for the new chains.
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateChainConfigChangeset),
			Config: changeset.UpdateChainConfigConfig{
				HomeChainSelector: e.HomeChainSel,
				RemoteChainAdds:   chainConfigs,
				MCMS:              mcmsConfig,
			},
		},
		{
			// Add the DONs and candidate commit OCR instances for the chain.
			Changeset: commonchangeset.WrapChangeSet(changeset.AddDonAndSetCandidateChangeset),
			Config: changeset.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: e.HomeChainSel,
					FeedChainSelector: e.FeedChainSel,
					MCMS:              mcmsConfig,
				},
				PluginInfo: changeset.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: ocrConfigs,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		},
		{
			// Add the exec OCR instances for the new chains.
			Changeset: commonchangeset.WrapChangeSet(changeset.SetCandidateChangeset),
			Config: changeset.SetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: e.HomeChainSel,
					FeedChainSelector: e.FeedChainSel,
					MCMS:              mcmsConfig,
				},
				PluginInfo: []changeset.SetCandidatePluginInfo{
					{
						OCRConfigPerRemoteChainSelector: ocrConfigs,
						PluginType:                      types.PluginTypeCCIPExec,
					},
				},
			},
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(changeset.PromoteCandidateChangeset),
			Config: changeset.PromoteCandidateChangesetConfig{
				HomeChainSelector: e.HomeChainSel,
				PluginInfo: []changeset.PromoteCandidatePluginInfo{
					{
						PluginType:           types.PluginTypeCCIPCommit,
						RemoteChainSelectors: allChains,
					},
					{
						PluginType:           types.PluginTypeCCIPExec,
						RemoteChainSelectors: allChains,
					},
				},
				MCMS: mcmsConfig,
			},
		},
		{
			// Enable the OCR config on the remote chains.
			Changeset: commonchangeset.WrapChangeSet(changeset.SetOCR3OffRampChangeset),
			Config: changeset.SetOCR3OffRampConfig{
				HomeChainSel:    e.HomeChainSel,
				RemoteChainSels: allChains,
			},
		},
	}
	if deployJobs {
		apps = append(apps, commonchangeset.ChangesetApplication{
			Changeset: commonchangeset.WrapChangeSet(changeset.CCIPCapabilityJobspecChangeset),
		})
	}
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, timelockContractsPerChain, apps)
	require.NoError(t, err)

	ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

	state, err = changeset.LoadOnchainState(e.Env)
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
	tEnv.UpdateDeployedEnvironment(e)
	return e
}

// NewEnvironmentWithJobs creates a new CCIP environment
// with capreg, fee tokens, feeds, nodes and jobs set up.
func NewEnvironmentWithJobs(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	e := NewEnvironment(t, tEnv)
	e.SetupJobs(t)
	return e
}
