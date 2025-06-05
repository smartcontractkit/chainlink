package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	solBinary "github.com/gagliardetto/binary"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangeSetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/cciptesthelpertypes"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

const (
	// NOTE: these are test values, real production values are configured in CLD.
	DefaultGasPriceDeviationPPB   = 1000
	DefaultDAGasPriceDeviationPPB = 1000
)

type EnvType string

const (
	Memory      EnvType = "in-memory"
	Docker      EnvType = "docker"
	ENVTESTTYPE         = "CCIP_V16_TEST_ENV"
)

type LogMessageToIgnore struct {
	Msg    string
	Reason string
	Level  zapcore.Level
}

type TestConfigs struct {
	Type      EnvType // set by env var CCIP_V16_TEST_ENV, defaults to Memory
	CreateJob bool
	// TODO: This should be CreateContracts so the booleans make sense?
	CreateJobAndContracts      bool
	PrerequisiteDeploymentOnly bool
	V1_5Cfg                    changeset.V1_5DeploymentConfig
	Chains                     int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	SolChains                  int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	AptosChains                int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	ChainIDs                   []uint64 // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	NumOfUsersPerChain         int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	Nodes                      int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	Bootstraps                 int      // only used in memory mode, for docker mode, this is determined by the integration-test config toml input
	IsUSDC                     bool
	IsTokenPoolFactory         bool
	IsUSDCAttestationMissing   bool
	IsMultiCall3               bool
	IsStaticLink               bool
	OCRConfigOverride          func(v1_6.CCIPOCRParams) v1_6.CCIPOCRParams
	RMNEnabled                 bool
	NumOfRMNNodes              int
	RMNConfDepth               int
	LinkPrice                  *big.Int
	WethPrice                  *big.Int
	BlockTime                  time.Duration
	// Test env related configs

	// LogMessagesToIgnore are log messages emitted by the chainlink node that cause
	// the test to auto-fail if they were logged.
	// In some tests we don't want this to happen where a failure is expected, e.g
	// we are purposely re-orging beyond finality.
	LogMessagesToIgnore []LogMessageToIgnore

	// ExtraConfigTomls contains the filenames of additional toml files to be loaded
	// to potentially override default configs.
	ExtraConfigTomls []string

	// CLNodeConfigOpts are the config options to be passed to the chainlink node.
	// Only used in memory mode.
	CLNodeConfigOpts []memory.ConfigOpt

	// RoleDONTopology is the chain-node topology of the role DON.
	// Only used in memory mode.
	RoleDONTopology cciptesthelpertypes.RoleDONTopology
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
		LinkPrice:             shared.MockLinkPrice,
		WethPrice:             shared.MockWethPrice,
		CreateJobAndContracts: true,
		BlockTime:             2 * time.Second,
	}
}

type TestOps func(testCfg *TestConfigs)

func WithLogMessagesToIgnore(logMessages []LogMessageToIgnore) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.LogMessagesToIgnore = logMessages
	}
}

func WithExtraConfigTomls(extraTomls []string) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.ExtraConfigTomls = extraTomls
	}
}

func WithCLNodeConfigOpts(opts ...memory.ConfigOpt) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.CLNodeConfigOpts = opts
	}
}

func WithBlockTime(blockTime time.Duration) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.BlockTime = blockTime
	}
}

func WithRMNConfDepth(depth int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.RMNConfDepth = depth
	}
}

func WithMultiCall3() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsMultiCall3 = true
	}
}

func WithStaticLink() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsStaticLink = true
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

func WithOCRConfigOverride(override func(v1_6.CCIPOCRParams) v1_6.CCIPOCRParams) TestOps {
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

func WithTokenPoolFactory() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsTokenPoolFactory = true
	}
}

func WithNumOfChains(numChains int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Chains = numChains
	}
}

func WithSolChains(numChains int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.SolChains = numChains
	}
}

func WithAptosChains(numChains int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.AptosChains = numChains
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

func WithRoleDONTopology(topology cciptesthelpertypes.RoleDONTopology) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.RoleDONTopology = topology
	}
}

type TestEnvironment interface {
	SetupJobs(t *testing.T)
	DeleteJobs(ctx context.Context, jobIDs map[string][]string) error
	StartNodes(t *testing.T, crConfig deployment.CapabilityRegistryConfig)
	StartChains(t *testing.T)
	TestConfigs() *TestConfigs
	DeployedEnvironment() DeployedEnv
	UpdateDeployedEnvironment(env DeployedEnv)
	MockUSDCAttestationServer(t *testing.T, isUSDCAttestationMissing bool) string
}

type DeployedEnv struct {
	Env                    cldf.Environment
	HomeChainSel           uint64
	FeedChainSel           uint64
	ReplayBlocks           map[uint64]uint64
	Users                  map[uint64][]*bind.TransactOpts
	RmnEnabledSourceChains map[uint64]bool
}

func (d *DeployedEnv) SetupJobs(t *testing.T) {
	_, err := commonchangeset.Apply(t, d.Env,
		commonchangeset.Configure(cldf.CreateLegacyChangeSet(v1_6.CCIPCapabilityJobspecChangeset), nil))
	require.NoError(t, err)
	ReplayLogs(t, d.Env.Offchain, d.ReplayBlocks)
}

type MemoryEnvironment struct {
	DeployedEnv
	nodes       map[string]memory.Node
	TestConfig  *TestConfigs
	Chains      map[uint64]cldf_evm.Chain
	SolChains   map[uint64]cldf_solana.Chain
	AptosChains map[uint64]cldf_aptos.Chain
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
	var chains map[uint64]cldf_evm.Chain
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
	m.SolChains = memory.NewMemoryChainsSol(t, tc.SolChains)
	aptosChains := memory.NewMemoryChainsAptos(t, tc.AptosChains)
	// if we have Aptos chains, we need to set the Aptos chain selectors on the wrapper
	// environment, so we have to convert it back to the concrete type. This needs to be refactored
	m.AptosChains = cldf_chain.NewBlockChainsFromSlice(aptosChains).AptosChains()

	blockChains := map[uint64]cldf_chain.BlockChain{}
	for selector, ch := range m.Chains {
		blockChains[selector] = ch
	}
	for selector, ch := range m.SolChains {
		blockChains[selector] = ch
	}
	for _, ch := range aptosChains {
		blockChains[ch.ChainSelector()] = ch
	}

	env := cldf.Environment{
		BlockChains: cldf_chain.NewBlockChains(blockChains),
	}
	homeChainSel, feedSel := allocateCCIPChainSelectors(chains)
	replayBlocks, err := LatestBlocksByChain(ctx, env)
	require.NoError(t, err)
	m.DeployedEnv = DeployedEnv{
		Env:          env,
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
	c := memory.NewNodesConfig{
		LogLevel:       zapcore.InfoLevel,
		Chains:         m.Chains,
		SolChains:      m.SolChains,
		AptosChains:    m.AptosChains,
		NumNodes:       tc.Nodes,
		NumBootstraps:  tc.Bootstraps,
		RegistryConfig: crConfig,
		CustomDBSetup:  nil,
	}
	nodes := memory.NewNodes(t, c, tc.CLNodeConfigOpts...)
	ctx := testcontext.Get(t)
	lggr := logger.Test(t)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}
	m.nodes = nodes
	m.Env = memory.NewMemoryEnvironmentFromChainsNodes(func() context.Context { return ctx }, lggr, m.Chains, m.SolChains, m.AptosChains, nodes)
}

func (m *MemoryEnvironment) DeleteJobs(ctx context.Context, jobIDs map[string][]string) error {
	for id, node := range m.nodes {
		if jobsToDelete, ok := jobIDs[id]; ok {
			for _, jobToDelete := range jobsToDelete {
				// delete job
				jobID, err := strconv.ParseInt(jobToDelete, 10, 32)
				if err != nil {
					return err
				}
				err = node.App.DeleteJob(ctx, int32(jobID))
				if err != nil {
					return fmt.Errorf("failed to delete job %s: %w", jobToDelete, err)
				}
			}
		}
	}
	return nil
}

func (m *MemoryEnvironment) MockUSDCAttestationServer(t *testing.T, isUSDCAttestationMissing bool) string {
	server := mockAttestationResponse(isUSDCAttestationMissing)
	endpoint := server.URL
	t.Cleanup(func() {
		server.Close()
	})
	return endpoint
}

// mineBlocks forces the simulated backend to produce a new block every X seconds
// NOTE: based on implementation in cltest/simulated_backend.go
func mineBlocks(backend *memory.Backend, blockTime time.Duration) (stopMining func()) {
	timer := time.NewTicker(blockTime)
	chStop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-timer.C:
				backend.Commit()
			case <-chStop:
				return
			}
		}
	}()
	return func() {
		close(chStop)
		timer.Stop()
		<-done
	}
}

func (m *MemoryEnvironment) MineBlocks(t *testing.T, blockTime time.Duration) {
	for _, chain := range m.Chains {
		if backend, ok := chain.Client.(*memory.Backend); ok {
			stopMining := mineBlocks(backend, blockTime)
			t.Cleanup(stopMining)
		}
	}
}

// NewMemoryEnvironment creates an in-memory environment based on the testconfig requested
func NewMemoryEnvironment(t *testing.T, opts ...TestOps) (DeployedEnv, TestEnvironment) {
	testCfg := DefaultTestConfigs()
	for _, opt := range opts {
		opt(testCfg)
	}
	testCfg.MustSetEnvTypeOrDefault(t)
	require.NoError(t, testCfg.Validate(), "invalid test config")
	env := &MemoryEnvironment{
		TestConfig: testCfg,
	}
	var dEnv DeployedEnv
	switch {
	case testCfg.PrerequisiteDeploymentOnly:
		dEnv = NewEnvironmentWithPrerequisitesContracts(t, env)
	case testCfg.CreateJobAndContracts:
		dEnv = NewEnvironmentWithJobsAndContracts(t, env)
	case testCfg.CreateJob:
		dEnv = NewEnvironmentWithJobs(t, env)
	default:
		dEnv = NewEnvironment(t, env)
	}
	env.UpdateDeployedEnvironment(dEnv)
	if testCfg.BlockTime > 0 {
		env.MineBlocks(t, testCfg.BlockTime)
	}
	return dEnv, env
}

func NewEnvironmentWithPrerequisitesContracts(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	var err error
	tc := tEnv.TestConfigs()
	e := NewEnvironment(t, tEnv)
	evmChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	solChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))

	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, c := range e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		mcmsCfg[c] = proposalutils.SingleGroupTimelockConfigV2(t)
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range evmChains {
		var opts []changeset.PrerequisiteOpt
		if tc != nil {
			if tc.IsTokenPoolFactory {
				opts = append(opts, changeset.WithTokenPoolFactoryEnabled())
			}
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
	deployLinkApp := commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
		evmChains,
	)

	if tc.IsStaticLink {
		deployLinkApp = commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployStaticLinkToken),
			evmChains,
		)
	}
	e.Env, err = commonchangeset.Apply(t, e.Env, deployLinkApp, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
		changeset.DeployPrerequisiteConfig{
			Configs: prereqCfg,
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		mcmsCfg,
	))
	require.NoError(t, err)
	if len(solChains) > 0 {
		solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()
		deploySolanaLinkApp := commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken),
			commonchangeset.DeploySolanaLinkTokenConfig{
				ChainSelector: solChains[0],
				TokenPrivKey:  solLinkTokenPrivKey,
				TokenDecimals: 9,
			},
		)
		e.Env, err = commonchangeset.Apply(t, e.Env,
			deploySolanaLinkApp,
		)
		require.NoError(t, err)
	}
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
	require.NotEmpty(t, dEnv.Env.BlockChains.EVMChains())
	ab := cldf.NewMemoryAddressBook()
	crConfig := DeployTestContracts(t, lggr, ab, dEnv.HomeChainSel, dEnv.FeedChainSel, dEnv.Env.BlockChains.EVMChains(), tc.LinkPrice, tc.WethPrice)
	tEnv.StartNodes(t, crConfig)
	dEnv = tEnv.DeployedEnvironment()
	dEnv.Env.ExistingAddresses = ab
	return dEnv
}

func NewEnvironmentWithJobsAndContracts(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	var err error
	e := NewEnvironmentWithPrerequisitesContracts(t, tEnv)
	evmChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	solChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	aptosChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))
	//nolint:gocritic // we need to segregate EVM and Solana chains
	allChains := append(evmChains, solChains...)
	allChains = append(allChains, aptosChains...)
	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)

	for _, c := range e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		mcmsCfg[c] = proposalutils.SingleGroupTimelockConfig(t)
	}

	tEnv.UpdateDeployedEnvironment(e)

	e = AddCCIPContractsToEnvironment(t, allChains, tEnv, false)
	// now we update RMNProxy to point to RMNRemote
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset),
			v1_6.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: evmChains,
			},
		),
	)
	require.NoError(t, err)

	// load the state again to get the latest addresses
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	err = state.ValidatePostDeploymentState(e.Env)
	require.NoError(t, err)
	return e
}
func DeployChainContractsToSolChainCS(e DeployedEnv, solChainSelector uint64, preload bool, buildSolConfig *ccipChangeSetSolana.BuildSolanaConfig) ([]commonchangeset.ConfiguredChangeSet, error) {
	var mcmsCfg *commontypes.MCMSWithTimelockConfigV2
	if preload {
		// Pre load default programs
		err := SavePreloadedSolAddresses(e.Env, solChainSelector)
		if err != nil {
			return nil, err
		}
	} else {
		mcmsCfg = &commontypes.MCMSWithTimelockConfigV2{
			Proposer: mcmstypes.Config{
				Quorum:  1,
				Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000001")},
			},
			Canceller: mcmstypes.Config{
				Quorum:  1,
				Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000002")},
			},
			Bypasser: mcmstypes.Config{
				Quorum:  1,
				Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000002")},
			},
			TimelockMinDelay: big.NewInt(1),
		}
	}
	state, err := stateview.LoadOnchainState(e.Env)
	if err != nil {
		return nil, err
	}
	value := [28]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 51, 51, 74, 153, 67, 41, 73, 55, 39, 96, 0, 0}
	return []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.DeployChainContractsChangeset),
			ccipChangeSetSolana.DeployChainContractsConfig{
				HomeChainSelector: e.HomeChainSel,
				ChainSelector:     solChainSelector,
				ContractParamsPerChain: ccipChangeSetSolana.ChainContractParams{
					FeeQuoterParams: ccipChangeSetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{
							Lo: 15532559262904483840, Hi: 10, Endianness: nil,
						},
						BillingConfig: []solFeeQuoter.BillingTokenConfig{
							{
								Enabled: true,
								Mint:    state.SolChains[solChainSelector].LinkToken,
								UsdPerToken: solFeeQuoter.TimestampedPackedU224{
									Value:     value,
									Timestamp: time.Now().Unix(),
								},
								PremiumMultiplierWeiPerEth: 9e17,
							},
							{
								Enabled: true,
								Mint:    state.SolChains[solChainSelector].WSOL,
								UsdPerToken: solFeeQuoter.TimestampedPackedU224{
									Value:     value,
									Timestamp: time.Now().Unix(),
								},
								PremiumMultiplierWeiPerEth: 1e18,
							},
						},
					},
					OffRampParams: ccipChangeSetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
				BuildConfig:            buildSolConfig,
				MCMSWithTimelockConfig: mcmsCfg,
			},
		)}, nil
}

func AddCCIPContractsToEnvironment(t *testing.T, allChains []uint64, tEnv TestEnvironment, mcmsEnabled bool) DeployedEnv {
	tc := tEnv.TestConfigs()
	e := tEnv.DeployedEnvironment()
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	var apps []commonchangeset.ConfiguredChangeSet
	evmContractParams := make(map[uint64]ccipseq.ChainContractParams)

	evmChains := []uint64{}
	for _, chain := range allChains {
		if _, ok := e.Env.BlockChains.EVMChains()[chain]; ok {
			evmChains = append(evmChains, chain)
		}
	}

	solChains := []uint64{}
	for _, chain := range allChains {
		if _, ok := e.Env.BlockChains.SolanaChains()[chain]; ok {
			solChains = append(solChains, chain)
		}
	}

	aptosChains := []uint64{}
	for _, chain := range allChains {
		if _, ok := e.Env.BlockChains.AptosChains()[chain]; ok {
			aptosChains = append(aptosChains, chain)
		}
	}

	for _, chain := range evmChains {
		evmContractParams[chain] = ccipseq.ChainContractParams{
			FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
			OffRampParams:   ccipops.DefaultOffRampParams(),
		}
	}

	apps = append(apps, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				NodeOperators:    NewTestNodeOperator(e.Env.BlockChains.EVMChains()[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			ccipseq.DeployChainContractsConfig{
				HomeChainSelector:      e.HomeChainSel,
				ContractParamsPerChain: evmContractParams,
			},
		),
	}...)
	if len(solChains) != 0 {
		solCs, err := DeployChainContractsToSolChainCS(e, solChains[0], true, nil)
		require.NoError(t, err)
		apps = append(apps, solCs...)
	}
	e.Env, _, err = commonchangeset.ApplyChangesets(t, e.Env, apps)
	require.NoError(t, err)

	if len(aptosChains) != 0 {
		// Currently only one aptos chain is supported in test environment
		aptosCs := DeployChainContractsToAptosCS(t, e, aptosChains[0])
		e.Env, _, err = commonchangeset.ApplyChangesets(t, e.Env, []commonchangeset.ConfiguredChangeSet{aptosCs})
		require.NoError(t, err)
	}

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	// Assert link present
	if tc.IsStaticLink {
		require.NotNil(t, state.MustGetEVMChainState(e.FeedChainSel).StaticLinkToken)
	} else {
		require.NotNil(t, state.MustGetEVMChainState(e.FeedChainSel).LinkToken)
	}
	require.NotNil(t, state.MustGetEVMChainState(e.FeedChainSel).Weth9)

	tokenConfig := shared.NewTestTokenConfig(state.MustGetEVMChainState(e.FeedChainSel).USDFeeds)
	var tokenDataProviders []pluginconfig.TokenDataObserverConfig
	if tc.IsUSDC {
		endpoint := tEnv.MockUSDCAttestationServer(t, tc.IsUSDCAttestationMissing)
		cctpContracts := make(map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig)
		for _, usdcChain := range evmChains {
			require.NotNil(t, state.MustGetEVMChainState(usdcChain).MockUSDCTokenMessenger)
			require.NotNil(t, state.MustGetEVMChainState(usdcChain).MockUSDCTransmitter)
			require.NotNil(t, state.MustGetEVMChainState(usdcChain).USDCTokenPools[deployment.Version1_5_1])
			cctpContracts[cciptypes.ChainSelector(usdcChain)] = pluginconfig.USDCCCTPTokenConfig{
				SourcePoolAddress:            state.MustGetEVMChainState(usdcChain).USDCTokenPools[deployment.Version1_5_1].Address().String(),
				SourceMessageTransmitterAddr: state.MustGetEVMChainState(usdcChain).MockUSDCTransmitter.Address().String(),
			}
		}
		tokenDataProviders = append(tokenDataProviders, pluginconfig.TokenDataObserverConfig{
			Type:    pluginconfig.USDCCCTPHandlerType,
			Version: "1.0",
			USDCCCTPObserverConfig: &pluginconfig.USDCCCTPObserverConfig{
				AttestationConfig: pluginconfig.AttestationConfig{
					AttestationAPI:         endpoint,
					AttestationAPITimeout:  commonconfig.MustNewDuration(time.Second),
					AttestationAPIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
				},
				Tokens: cctpContracts,
			}})
	}

	nodeInfo, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// generate the chainToNodeMapping if we have a topology provided.
	var chainToNodeMapping map[cciptypes.ChainSelector][][32]byte
	if tc.Type == Memory && tc.RoleDONTopology != nil {
		allSelectors := make([]cciptypes.ChainSelector, 0, len(evmChains)+len(solChains))
		for _, chain := range evmChains {
			// don't include the home chain, its supported by all nodes.
			if chain == e.HomeChainSel {
				continue
			}
			allSelectors = append(allSelectors, cciptypes.ChainSelector(chain))
		}
		for _, chain := range solChains {
			allSelectors = append(allSelectors, cciptypes.ChainSelector(chain))
		}
		chainToNodeMapping, err = tc.RoleDONTopology.ChainToNodeMapping(
			nodeInfo.NonBootstraps().PeerIDs(),
			allSelectors,
			cciptypes.ChainSelector(e.HomeChainSel),
		)
		require.NoError(t, err)
	}

	// Build the CCIPHome chain configs.
	chainConfigs := make(map[uint64]v1_6.ChainConfig)
	commitOCRConfigs := make(map[uint64]v1_6.CCIPOCRParams)
	execOCRConfigs := make(map[uint64]v1_6.CCIPOCRParams)
	for _, chain := range evmChains {
		var linkTokenAddr common.Address
		if tc.IsStaticLink {
			linkTokenAddr = state.MustGetEVMChainState(chain).StaticLinkToken.Address()
		} else {
			linkTokenAddr = state.MustGetEVMChainState(chain).LinkToken.Address()
		}
		ocrOverride := func(ocrParams v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			if tc.OCRConfigOverride != nil {
				tc.OCRConfigOverride(ocrParams)
			}
			if tc.RMNEnabled {
				if ocrParams.CommitOffChainConfig != nil {
					ocrParams.CommitOffChainConfig.RMNEnabled = true
				}
			} else {
				if ocrParams.CommitOffChainConfig != nil {
					ocrParams.CommitOffChainConfig.RMNEnabled = false
				}
			}
			return ocrParams
		}
		commitOCRConfigs[chain] = v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, e.FeedChainSel, tokenConfig.GetTokenInfo(e.Env.Logger, linkTokenAddr, state.MustGetEVMChainState(chain).Weth9.Address()), ocrOverride)
		execOCRConfigs[chain] = v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, tokenDataProviders, ocrOverride)

		var readers [][32]byte
		if chainToNodeMapping != nil {
			_, ok := chainToNodeMapping[cciptypes.ChainSelector(chain)]
			require.True(t, ok, "chain %d not found in chainToNodeMapping", chain)
			readers = chainToNodeMapping[cciptypes.ChainSelector(chain)]
			t.Logf("setting readers for chain %d to %v due to topology %v", chain, readers, chainToNodeMapping)
		} else {
			t.Logf("setting readers for chain %d to %v due to no topology", chain, nodeInfo.NonBootstraps().PeerIDs())
			readers = nodeInfo.NonBootstraps().PeerIDs()
		}
		chainConfigs[chain] = v1_6.ChainConfig{
			Readers: readers,
			// #nosec G115 - Overflow is not a concern in this test scenario
			FChain: uint8(len(readers) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:      cciptypes.BigInt{Int: big.NewInt(DefaultGasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(DefaultDAGasPriceDeviationPPB)},
				OptimisticConfirmations:   globals.OptimisticConfirmations,
				ChainFeeDeviationDisabled: false,
			},
		}
	}

	for _, chain := range solChains {
		// TODO: this is a workaround for tokenConfig.GetTokenInfo
		tokenInfo := map[cciptypes.UnknownEncodedAddress]pluginconfig.TokenInfo{}
		tokenInfo[cciptypes.UnknownEncodedAddress(state.SolChains[chain].LinkToken.String())] = tokenConfig.TokenSymbolToInfo[shared.LinkSymbol]
		// TODO: point this to proper SOL feed, apparently 0 signified SOL
		tokenInfo[cciptypes.UnknownEncodedAddress(solanago.SolMint.String())] = tokenConfig.TokenSymbolToInfo[shared.WethSymbol]

		ocrOverride := tc.OCRConfigOverride
		commitOCRConfigs[chain] = v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, e.FeedChainSel, tokenInfo, ocrOverride)
		execOCRConfigs[chain] = v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, tokenDataProviders, ocrOverride)

		var readers [][32]byte
		if chainToNodeMapping != nil {
			_, ok := chainToNodeMapping[cciptypes.ChainSelector(chain)]
			require.True(t, ok, "chain %d not found in chainToNodeMapping", chain)
			readers = chainToNodeMapping[cciptypes.ChainSelector(chain)]
			t.Logf("setting readers for chain %d to %v due to topology %v", chain, readers, chainToNodeMapping)
		} else {
			readers = nodeInfo.NonBootstraps().PeerIDs()
		}
		chainConfigs[chain] = v1_6.ChainConfig{
			Readers: readers,
			// #nosec G115 - Overflow is not a concern in this test scenario
			FChain: uint8(len(readers) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:      cciptypes.BigInt{Int: big.NewInt(DefaultGasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(DefaultDAGasPriceDeviationPPB)},
				OptimisticConfirmations:   globals.OptimisticConfirmations,
				ChainFeeDeviationDisabled: true,
			},
		}
	}

	for _, chain := range aptosChains {
		// TODO(aptos): update this for token transfers
		tokenInfo := map[cciptypes.UnknownEncodedAddress]pluginconfig.TokenInfo{}
		linkTokenAddress := state.AptosChains[chain].LinkTokenAddress
		tokenInfo[cciptypes.UnknownEncodedAddress(linkTokenAddress.String())] = tokenConfig.TokenSymbolToInfo[shared.LinkSymbol]

		ocrOverride := tc.OCRConfigOverride
		commitOCRConfigs[chain] = v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, e.FeedChainSel, tokenInfo, ocrOverride)
		execOCRConfigs[chain] = v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, tokenDataProviders, ocrOverride)
		chainConfigs[chain] = v1_6.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			// #nosec G115 - Overflow is not a concern in this test scenario
			FChain: uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(DefaultGasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(DefaultDAGasPriceDeviationPPB)},
				OptimisticConfirmations: globals.OptimisticConfirmations,
			},
		}
	}

	// Apply second set of changesets to configure the CCIP contracts.
	var mcmsConfig *proposalutils.TimelockConfig
	if mcmsEnabled {
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 0,
		}
	}
	apps = []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			// Add the chain configs for the new chains.
			cldf.CreateLegacyChangeSet(v1_6.UpdateChainConfigChangeset),
			v1_6.UpdateChainConfigConfig{
				HomeChainSelector: e.HomeChainSel,
				RemoteChainAdds:   chainConfigs,
				MCMS:              mcmsConfig,
			},
		),
		commonchangeset.Configure(
			// Add the DONs and candidate commit OCR instances for the chain.
			cldf.CreateLegacyChangeSet(v1_6.AddDonAndSetCandidateChangeset),
			v1_6.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: e.HomeChainSel,
					// TODO: we dont know what this means for solana
					FeedChainSelector: e.FeedChainSel,
					MCMS:              mcmsConfig,
				},
				PluginInfo: v1_6.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: commitOCRConfigs,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		),
		commonchangeset.Configure(
			// Add the exec OCR instances for the new chains.
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			v1_6.SetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: e.HomeChainSel,
					// TODO: we dont know what this means for solana
					FeedChainSelector: e.FeedChainSel,
					MCMS:              mcmsConfig,
				},
				PluginInfo: []v1_6.SetCandidatePluginInfo{
					{
						OCRConfigPerRemoteChainSelector: execOCRConfigs,
						PluginType:                      types.PluginTypeCCIPExec,
					},
				},
			},
		),
		commonchangeset.Configure(
			// Promote everything
			cldf.CreateLegacyChangeSet(v1_6.PromoteCandidateChangeset),
			v1_6.PromoteCandidateChangesetConfig{
				HomeChainSelector: e.HomeChainSel,
				PluginInfo: []v1_6.PromoteCandidatePluginInfo{
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
		),
		commonchangeset.Configure(
			// Enable the OCR config on the remote chains.
			cldf.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset),
			v1_6.SetOCR3OffRampConfig{
				HomeChainSel:       e.HomeChainSel,
				RemoteChainSels:    evmChains,
				CCIPHomeConfigType: globals.ConfigTypeActive,
			},
		),
		commonchangeset.Configure(
			// Enable the OCR config on the remote chains.
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.SetOCR3ConfigSolana),
			v1_6.SetOCR3OffRampConfig{
				HomeChainSel:       e.HomeChainSel,
				RemoteChainSels:    solChains,
				CCIPHomeConfigType: globals.ConfigTypeActive,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.CCIPCapabilityJobspecChangeset),
			nil, // Changeset ignores any config
		),
	}
	e.Env, _, err = commonchangeset.ApplyChangesets(t, e.Env, apps)
	require.NoError(t, err)

	ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	require.NotNil(t, state.MustGetEVMChainState(e.HomeChainSel).CapabilityRegistry)
	require.NotNil(t, state.MustGetEVMChainState(e.HomeChainSel).CCIPHome)
	require.NotNil(t, state.MustGetEVMChainState(e.HomeChainSel).RMNHome)
	for _, chain := range evmChains {
		if tc.IsStaticLink {
			require.NotNil(t, state.MustGetEVMChainState(chain).StaticLinkToken)
		} else {
			require.NotNil(t, state.MustGetEVMChainState(chain).LinkToken)
		}
		require.NotNil(t, state.MustGetEVMChainState(chain).Weth9)
		require.NotNil(t, state.MustGetEVMChainState(chain).TokenAdminRegistry)
		require.NotEmpty(t, state.MustGetEVMChainState(chain).RegistryModules1_6)
		require.NotNil(t, state.MustGetEVMChainState(chain).Router)
		require.NotNil(t, state.MustGetEVMChainState(chain).RMNRemote)
		require.NotNil(t, state.MustGetEVMChainState(chain).TestRouter)
		require.NotNil(t, state.MustGetEVMChainState(chain).NonceManager)
		require.NotNil(t, state.MustGetEVMChainState(chain).FeeQuoter)
		require.NotNil(t, state.MustGetEVMChainState(chain).OffRamp)
		require.NotNil(t, state.MustGetEVMChainState(chain).OnRamp)
	}
	err = ValidateSolanaState(e.Env, solChains)
	require.NoError(t, err)

	tEnv.UpdateDeployedEnvironment(e)
	return e
}

// NewEnvironmentWithJobs creates a new CCIP environment
// with home chain contracts, fee tokens, feeds, nodes and jobs set up.
func NewEnvironmentWithJobs(t *testing.T, tEnv TestEnvironment) DeployedEnv {
	e := NewEnvironment(t, tEnv)
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	// add home chain contracts, otherwise the job approval logic in chainlink fails silently
	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				NodeOperators:    NewTestNodeOperator(e.Env.BlockChains.EVMChains()[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			}))
	require.NoError(t, err)
	e.SetupJobs(t)
	return e
}
