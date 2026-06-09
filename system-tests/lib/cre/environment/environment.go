package environment

import (
	"context"
	"fmt"
	"maps"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ctfchiprouter "github.com/smartcontractkit/chainlink-testing-framework/framework/components/chiprouter"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	donconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	gateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cre.WorkflowRegistryOutput
	CreEnvironment                      *cre.Environment
	Dons                                *cre.Dons
	NodeOutput                          []*cre.NodeSetOutput
	GatewayConnectors                   *cre.GatewayConnectors
}

type SetupInput struct {
	NodeSets               []*cre.NodeSet
	Blockchains            []*config.Blockchain
	JdInput                *config.JobDistributor
	ChipRouterInput        *ctfchiprouter.Input
	Provider               infra.Provider
	OCR3Config             *keystone_changeset.OracleConfig
	DONTimeConfig          *keystone_changeset.OracleConfig
	VaultOCR3Config        *keystone_changeset.OracleConfig
	S3ProviderInput        *s3provider.Input
	CapabilityConfigs      cre.CapabilityConfigs
	Capabilities           []cre.InstallableCapability //nolint:staticcheck //SA1019 - We can't remove until other repos are updated
	Features               cre.Features
	GatewayWhitelistConfig gateway.WhitelistConfig
	BlockchainDeployers    map[blockchain.ChainFamily]blockchains.Deployer
	ContractVersions       map[cre.ContractType]*semver.Version

	// allow to pass custom transformers for extensibility
	ConfigFactoryFunctions               []cre.NodeConfigTransformerFn
	JobSpecFactoryFunctions              []cre.JobSpecFn
	CapabilitiesContractFactoryFunctions []cre.CapabilityRegistryConfigFn

	StageGen *stagegen.StageGen

	// Optional hook executed after local dependencies are started (including JD),
	// and right before DON containers are started.
	PreDONsStartHook func(ctx context.Context) error
}

func (s *SetupInput) Validate() error {
	if s == nil {
		return pkgerrors.New("input is nil")
	}

	if len(s.NodeSets) == 0 {
		return pkgerrors.New("at least one nodeSet is required")
	}

	if len(s.Blockchains) == 0 {
		return pkgerrors.New("at least one blockchain is required")
	}

	if s.JdInput == nil {
		return pkgerrors.New("jd input is nil")
	}
	if err := s.JdInput.Validate(); err != nil {
		return pkgerrors.Wrap(err, "jd input validation failed")
	}

	return nil
}

func SetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input *SetupInput,
	relativePathToRepoRoot string,
) (*SetupOutput, error) {
	if input == nil {
		return nil, pkgerrors.New("input is nil")
	}

	if err := input.Validate(); err != nil {
		return nil, pkgerrors.Wrap(err, "input validation failed")
	}
	execPlan, err := buildPlacementPlan(input.Blockchains, input.JdInput, input.NodeSets)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "invalid component placement")
	}

	remoteRuntime, err := resolveRemoteRuntimeForSetup(testLogger, execPlan)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to resolve remote runtime settings")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting Chip Router")))
	_, err = ctfchiprouter.NewWithContext(ctx, input.ChipRouterInput)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to start chip router")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Chip Router started in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting %d blockchain(s)", len(input.Blockchains))))

	testLogger.Info().Msg("using persistent relay supervisor for mixed component relays")

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting %d blockchain(s)", len(input.Blockchains))))

	deployedBlockchains, startErr := startBlockchains(
		ctx,
		testLogger,
		singleFileLogger,
		input.Blockchains,
		input.BlockchainDeployers,
		remoteRuntime,
		execPlan.NodeSetPlacement.HasLocalTargets,
	)
	if startErr != nil {
		return nil, pkgerrors.Wrap(startErr, "failed to start blockchains")
	}

	creEnvironment := &cre.Environment{
		Blockchains:           deployedBlockchains.Outputs,
		Provider:              input.Provider,
		RegistryChainSelector: deployedBlockchains.RegistryChain().ChainSelector(),
		ContractVersions:      input.ContractVersions,
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Blockchains started in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Deploying Workflow and Capability Registry contracts")))

	deployKeystoneContractsOutput, deployErr := crecontracts.DeployKeystoneContracts(
		ctx,
		testLogger,
		singleFileLogger,
		crecontracts.DeployKeystoneContractsInput{
			CldfEnvironment: newCldfEnvironment(ctx, singleFileLogger, deployedBlockchains.CldfBlockChains),
			CtfBlockchains:  deployedBlockchains.Outputs,
		},
	)
	if deployErr != nil {
		return nil, pkgerrors.Wrap(deployErr, "failed to deploy Keystone contracts")
	}
	creEnvironment.CldfEnvironment = deployKeystoneContractsOutput.Env

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Workflow and Capability Registry contracts deployed in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Preparing DONs configuration")))

	topology, tErr := cre.NewTopology(input.NodeSets, creEnvironment.Provider, input.CapabilityConfigs)
	if tErr != nil {
		return nil, pkgerrors.Wrap(tErr, "failed to create topology")
	}
	remoteHostIP := ""
	if remoteRuntime != nil {
		remoteHostIP = remoteRuntime.RemoteHostIP
	}

	updatedNodeSets, topoErr := donconfig.PrepareNodeTOMLs(
		ctx,
		topology,
		creEnvironment,
		input.NodeSets,
		input.Blockchains,
		donconfig.PrepareNodeTOMLsOptions{
			RemoteHostIP: remoteHostIP,
		},
		input.Capabilities,
		input.ConfigFactoryFunctions,
		input.ChipRouterInput.Out.InternalGRPCURL,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Applying Features before environment startup")))
	var donsCapabilities = make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	var capabilityToOCR3Config = make(map[string]*ocr3.OracleConfig)
	capabilityToExtraSignerFamilies := make(map[string][]string)
	for _, feature := range input.Features.List() {
		for _, donMetadata := range topology.DonsMetadataWithFlag(feature.Flag()) {
			testLogger.Info().Msgf("Executing PreEnvStartup for feature %s for don '%s'", feature.Flag(), donMetadata.Name)
			output, preErr := feature.PreEnvStartup(
				ctx,
				testLogger,
				donMetadata,
				topology,
				creEnvironment,
			)
			if preErr != nil {
				return nil, fmt.Errorf("failed to execute PreEnvStartup for feature %s: %w", feature.Flag(), preErr)
			}
			if output != nil {
				if donsCapabilities[donMetadata.ID] == nil {
					donsCapabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
				}
				donsCapabilities[donMetadata.ID] = append(donsCapabilities[donMetadata.ID], output.DONCapabilityWithConfig...)
				maps.Copy(capabilityToOCR3Config, output.CapabilityToOCR3Config)
				for capability, families := range output.CapabilityToExtraSignerFamilies {
					capabilityToExtraSignerFamilies[capability] = append([]string(nil), families...)
				}
			}
			testLogger.Info().Msgf("PreEnvStartup for feature %s executed successfully", feature.Flag())
		}
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Applied Features in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting Job Distributor and DONs")))

	startedJD, jdStartErr := StartJD(ctx, testLogger, input.JdInput, input.Provider, remoteRuntime)
	if jdStartErr != nil {
		return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
	}
	if input.PreDONsStartHook != nil {
		if err := input.PreDONsStartHook(ctx); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to execute pre-DON startup hook")
		}
	}

	startedDONs, donStartErr := StartDONs(ctx, testLogger, topology, input.Provider, deployedBlockchains.RegistryChain().CtfOutput(), input.CapabilityConfigs, updatedNodeSets, remoteRuntime)
	if donStartErr != nil {
		return nil, pkgerrors.Wrap(donStartErr, "failed to start DONs")
	}
	dons := cre.NewDons(startedDONs.DONs(), topology.GatewayConnectors)
	deployKeystoneContractsOutput.Env.Offchain = startedJD.Client

	linkDonsToJDInput := &cre.LinkDonsToJDInput{
		Blockchains:     deployedBlockchains.Outputs,
		CldfEnvironment: deployKeystoneContractsOutput.Env,
		Topology:        topology,
		Dons:            dons,
		JDPlacement:     string(input.JdInput.Placement),
		JDInternalWSRPC: startedJD.JDOutput.InternalWSRPCUrl,
		JDExternalWSRPC: startedJD.JDOutput.ExternalWSRPCUrl,
	}

	cldErr := cre.LinkToJobDistributor(ctx, linkDonsToJDInput)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to link DONs to Job Distributor")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs and Job Distributor started and linked in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Creating Jobs with Job Distributor")))

	gJobErr := gateway.CreateJobs(ctx, creEnvironment, dons, topology.GatewayServiceConfigs, input.GatewayWhitelistConfig)
	if gJobErr != nil {
		return nil, pkgerrors.Wrap(gJobErr, "failed to create gateway jobs with Job Distributor")
	}

	// Deprecated: use Features instead. Support for InstallableCapability will be removed in the future.
	jobSpecFactoryFunctions := make([]cre.JobSpecFn, 0)
	for _, capability := range input.Capabilities {
		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, capability.JobSpecFn())
	}

	// allow to pass custom job spec factories for extensibility
	jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, input.JobSpecFactoryFunctions...)

	createJobsDeps := CreateJobsWithJdOpDeps{
		Logger:                        testLogger,
		SingleFileLogger:              singleFileLogger,
		RegistryChainBlockchainOutput: deployedBlockchains.RegistryChain().CtfOutput(),
		JobSpecFactoryFunctions:       jobSpecFactoryFunctions,
		CreEnvironment:                creEnvironment,
		Dons:                          dons,
		NodeSets:                      input.NodeSets,
		Capabilities:                  input.Capabilities,
	}
	_, createJobsErr := operations.ExecuteOperation(deployKeystoneContractsOutput.Env.OperationsBundle, CreateJobsWithJdOp, createJobsDeps, CreateJobsWithJdOpInput{})
	if createJobsErr != nil {
		return nil, pkgerrors.Wrap(createJobsErr, "failed to create jobs with Job Distributor")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Jobs created in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Funding Chainlink nodes")))

	fundingPerChainFamilyForEachNode := map[string]uint64{
		chainselectors.FamilyEVM:    10000000000000000, // 0.01 ETH
		chainselectors.FamilySolana: 50_000_000_000,    // 50 SOL
		chainselectors.FamilyTron:   100_000_000,       // 100 TRX in SUN
		chainselectors.FamilyAptos:  1_000_000_000_000, // 1,000 APT (octas) for local devnet sender accounts
	}

	fErr := FundNodes(
		ctx,
		testLogger,
		dons,
		deployedBlockchains.Outputs,
		fundingPerChainFamilyForEachNode,
	)
	if fErr != nil {
		return nil, pkgerrors.Wrap(fErr, "failed to fund chainlink nodes")
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Chainlink nodes funded in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Configuring Workflow and Capability Registry contracts")))

	// Configure Capabilities Registry first so we can resolve actual contract DON IDs
	// before wiring the workflow registry. Some downstream changesets read DON info
	// from CapReg state rather than the pre-contract topology shape.
	capRegInput := cre.ConfigureCapabilityRegistryInput{
		ChainSelector: deployedBlockchains.RegistryChain().ChainSelector(),
		CldEnv:        creEnvironment.CldfEnvironment,
		Blockchains:   deployedBlockchains.Outputs,
		Topology:      topology,
		Provider:      input.Provider,
		CapabilitiesRegistryAddress: new(crecontracts.MustGetAddressFromMemoryDataStore(
			deployKeystoneContractsOutput.MemoryDataStore,
			deployedBlockchains.RegistryChain().ChainSelector(),
			keystone_changeset.CapabilitiesRegistry.String(),
			crecontracts.V2Version,
			""),
		),
		NodeSets:                        input.NodeSets,
		DONCapabilityWithConfigs:        make(map[uint64][]keystone_changeset.DONCapabilityWithConfig),
		CapabilityToOCR3Config:          capabilityToOCR3Config,
		CapabilityToExtraSignerFamilies: capabilityToExtraSignerFamilies,
	}

	for _, capability := range input.Capabilities {
		configFn := capability.CapabilityRegistryV2ConfigFn()
		capRegInput.CapabilityRegistryConfigFns = append(capRegInput.CapabilityRegistryConfigFns, configFn)
	}
	capRegInput.CapabilityRegistryConfigFns = append(capRegInput.CapabilityRegistryConfigFns, input.CapabilitiesContractFactoryFunctions...)
	maps.Copy(capRegInput.DONCapabilityWithConfigs, donsCapabilities)

	capReg, capRegErr := crecontracts.ConfigureCapabilityRegistry(ctx, capRegInput)
	if capRegErr != nil {
		return nil, pkgerrors.Wrap(capRegErr, "failed to configure Capability Registry contracts")
	}

	// Resolve actual contract donIDs and apply to topology, dons, and NodeSets
	if err := crecontracts.ResolveAndApplyContractDonIDs(capReg, dons, topology, input.NodeSets); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to resolve and apply contract donIDs")
	}

	workflowRegistryConfigurationOutput, wfErr := workflow.ConfigureWorkflowRegistry(
		ctx,
		testLogger,
		singleFileLogger,
		&cre.WorkflowRegistryInput{
			ContractAddress: common.HexToAddress(crecontracts.MustGetAddressFromDataStore(deployKeystoneContractsOutput.Env.DataStore, deployedBlockchains.RegistryChain().ChainSelector(), keystone_changeset.WorkflowRegistry.String(), crecontracts.V2Version, "")),
			ContractVersion: cldf.TypeAndVersion{Version: *crecontracts.V2Version},
			ChainSelector:   deployedBlockchains.RegistryChain().ChainSelector(),
			CldEnv:          deployKeystoneContractsOutput.Env,
			AllowedDonIDs:   topology.WorkflowDONIDs,
			WorkflowOwners:  []common.Address{deployedBlockchains.RegistryChain().(*evm.Blockchain).SethClient.MustGetRootKeyAddress()}, // registry chain is always EVM
		},
	)
	if wfErr != nil {
		return nil, pkgerrors.Wrap(wfErr, "failed to configure workflow registry")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Workflow and Capability Registry contracts configured in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Applying Features after environment startup")))

	for _, feature := range input.Features.List() {
		for _, don := range dons.DonsWithFlag(feature.Flag()) {
			testLogger.Info().Msgf("Executing PostEnvStartup for feature %s for don '%s'", feature.Flag(), don.Name)
			if pErr := feature.PostEnvStartup(
				ctx,
				testLogger,
				don,
				dons,
				creEnvironment,
			); pErr != nil {
				return nil, fmt.Errorf("failed to execute PostEnvStartup for feature %s: %w", feature.Flag(), pErr)
			}
			testLogger.Info().Msgf("PostEnvStartup for feature %s executed successfully", feature.Flag())
		}
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Features applied in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	// Sharding setup moved AFTER PostEnvStartup to ensure OCR3 configs work properly
	if topology.DonsMetadata.ShardingEnabled() {
		fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Setting up Sharding")))
		err := sharding.SetupSharding(ctx, sharding.SetupShardingInput{
			Logger:   testLogger,
			CreEnv:   creEnvironment,
			Topology: topology,
			Dons:     dons,
		})
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to setup Sharding")
		}
		fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Sharding setup in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	}

	appendOutputsToInput(input, startedDONs.NodeOutputs(), deployedBlockchains.Outputs, startedJD.JDOutput)

	if err := workflowRegistryConfigurationOutput.Store(config.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to store workflow registry configuration output")
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryConfigurationOutput, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		Dons:                                dons,
		NodeOutput:                          startedDONs.NodeOutputs(),
		CreEnvironment:                      creEnvironment,
		GatewayConnectors:                   topology.GatewayConnectors,
	}, nil
}

func appendOutputsToInput(input *SetupInput, nodeSetOutput []*cre.NodeSetOutput, blockchains []blockchains.Blockchain, jdOutput *jd.Output) {
	// append the nodeset output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	for idx, nsOut := range nodeSetOutput {
		input.NodeSets[idx].Out = nsOut.Output
	}

	for idx, deployedBlockchain := range blockchains {
		if idx < len(input.Blockchains) && input.Blockchains[idx] != nil {
			input.Blockchains[idx].Out = deployedBlockchain.CtfOutput()
		}
	}

	// append the jd output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	input.JdInput.Out = jdOutput
}

func resolveRemoteRuntimeForSetup(
	testLogger zerolog.Logger,
	execPlan *placementPlan,
) (*remoteclient.Runtime, error) {
	if execPlan == nil || !execPlan.HasRemoteComponents {
		return nil, nil
	}
	runtimeInput, err := resolveRemoteRuntimeInput()
	if err != nil {
		return nil, err
	}
	return remoteclient.ResolveRuntimeWithInput(testLogger, runtimeInput)
}

func resolveRemoteRuntimeInput() (remoteclient.RuntimeInput, error) {
	input := remoteclient.RuntimeInput{
		AgentBaseURL: strings.TrimSpace(os.Getenv(remoteclient.EnvRemoteAgentURL)),
		RemoteHostIP: strings.TrimSpace(os.Getenv(runtimecfg.EnvRemoteHostIP)),
	}
	if configuredPort := strings.TrimSpace(os.Getenv(remoteclient.EnvRemoteAgentPort)); configuredPort != "" {
		parsedPort, err := strconv.Atoi(configuredPort)
		if err != nil || parsedPort <= 0 || parsedPort > 65535 {
			return remoteclient.RuntimeInput{}, fmt.Errorf("invalid %s: %q", remoteclient.EnvRemoteAgentPort, configuredPort)
		}
		input.AgentPort = parsedPort
	}
	return input, nil
}

func verifyRemoteToLocalBootstrapReachability(ctx context.Context, lggr zerolog.Logger, topology *cre.Topology) error {
	if topology == nil {
		return nil
	}
	hasRemoteDONs := false
	hasLocalBootstrap := false
	for _, don := range topology.DonsMetadata.List() {
		if don == nil || don.MustNodeSet() == nil {
			continue
		}
		placement := strings.TrimSpace(don.MustNodeSet().Placement)
		if placement == string(config.PlacementRemote) {
			hasRemoteDONs = true
		}
		if placement == string(config.PlacementLocal) {
			for _, node := range don.NodesMetadata {
				if node != nil && node.HasRole(cre.BootstrapNode) {
					hasLocalBootstrap = true
					break
				}
			}
		}
	}
	if !hasRemoteDONs || !hasLocalBootstrap {
		return nil
	}
	if !runtimecfg.IsDirectMode() {
		return nil
	}

	ec2HostIP, err := runtimecfg.DirectHostIP()
	if err != nil {
		return fmt.Errorf("resolve direct EC2 host ip: %w", err)
	}
	remoteRelayAddr := net.JoinHostPort(ec2HostIP, strconv.Itoa(cre.OCRPeeringPort))
	if err := waitForTCPReachable(ctx, remoteRelayAddr, 6*time.Second); err != nil {
		return fmt.Errorf("remote relay listener for bootstrap peering is not reachable at %s: %w", remoteRelayAddr, err)
	}
	lggr.Info().Str("remoteRelay", remoteRelayAddr).Msg("verified remote->local bootstrap relay listener reachability")
	return nil
}

func waitForTCPReachable(ctx context.Context, addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for {
		dialer := net.Dialer{Timeout: 600 * time.Millisecond}
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		lastErr = err
		if time.Now().After(deadline) {
			return lastErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func newCldfEnvironment(ctx context.Context, singleFileLogger logger.Logger, cldfBlockchains cldf_chain.BlockChains) *cldf.Environment {
	allChainsCLDEnvironment := &cldf.Environment{
		Name:              cre.EnvironmentName,
		Logger:            singleFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(), // can't set it to nil, because some changesets save addresses both to the address book and datastore
		DataStore:         datastore.NewMemoryDataStore().Seal(),
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: cldfBlockchains,
		OCRSecrets:  focr.XXXGenerateTestOCRSecrets(),
		OperationsBundle: operations.NewBundle(
			func() context.Context { return ctx },
			singleFileLogger, operations.NewMemoryReporter()),
	}

	return allChainsCLDEnvironment
}
