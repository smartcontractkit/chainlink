package environment

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

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
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	donconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	gateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const envUsePersistentRelaySupervisor = "CRE_USE_PERSISTENT_RELAY_SUPERVISOR"

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cre.WorkflowRegistryOutput
	CreEnvironment                      *cre.Environment
	Dons                                *cre.Dons
	NodeOutput                          []*cre.NodeSetOutput
	S3ProviderOutput                    *s3provider.Output
	GatewayConnectors                   *cre.GatewayConnectors

	tunnelManager tunnel.Manager
	relayManager  *componentRelayManager
	closeOnce     sync.Once
	closeErr      error
}

func (s *SetupOutput) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}
	manager := s.tunnelManager
	if manager == nil {
		manager = tunnel.NewNoopManager()
	}

	s.closeOnce.Do(func() {
		if s.relayManager != nil {
			_ = s.relayManager.Close(ctx)
		}
		s.closeErr = manager.Stop(ctx)
	})

	return s.closeErr
}

func (s *SetupOutput) TunnelBindings() []tunnel.TunnelBinding {
	if s == nil || s.tunnelManager == nil {
		return []tunnel.TunnelBinding{}
	}
	return s.tunnelManager.Snapshot()
}

type SetupInput struct {
	NodeSets               []*cre.NodeSet
	Blockchains            []*config.Blockchain
	JdInput                *config.JobDistributor
	Provider               infra.Provider
	ContractVersions       map[cre.ContractType]*semver.Version
	WithV2Registries       bool
	OCR3Config             *keystone_changeset.OracleConfig
	DONTimeConfig          *keystone_changeset.OracleConfig
	VaultOCR3Config        *keystone_changeset.OracleConfig
	S3ProviderInput        *s3provider.Input
	CapabilityConfigs      cre.CapabilityConfigs
	CopyCapabilityBinaries bool // if true, copy capability binaries to the containers (if false, we assume that the plugins image already has them)
	Capabilities           []cre.InstallableCapability
	Features               cre.Features
	GatewayWhitelistConfig gateway.WhitelistConfig
	BlockchainDeployers    map[blockchain.ChainFamily]blockchains.Deployer

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

	//TODO: remove these checks in December 2025, when everyone has migrated
	if val := os.Getenv("E2E_JD_IMAGE"); val != "" {
		return nil, errors.New("E2E_JD_IMAGE and E2E_JD_VERSION are deprecated, please use CTF_JD_IMAGE instead to specify the Job Distributor image with tag")
	}

	if val := os.Getenv("E2E_TEST_CHAINLINK_IMAGE"); val != "" {
		return nil, errors.New("E2E_TEST_CHAINLINK_IMAGE and E2E_TEST_CHAINLINK_VERSION are deprecated, please use CTF_CHAINLINK_IMAGE instead to specify the Chainlink Node image with tag")
	}

	if err := input.Validate(); err != nil {
		return nil, pkgerrors.Wrap(err, "input validation failed")
	}
	nodeSetPlacement, err := summarizeNodeSetPlacement(input.NodeSets)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "nodeset placement validation failed")
	}
	if err := validateUnsupportedPlacements(input.Blockchains, nodeSetPlacement); err != nil {
		return nil, pkgerrors.Wrap(err, "invalid component placement")
	}

	s3Output, s3Err := workflow.StartS3(testLogger, input.S3ProviderInput, input.StageGen)
	if s3Err != nil {
		return nil, pkgerrors.Wrap(s3Err, "failed to start S3 provider")
	}

	tunnelManager, tmErr := newEC2TunnelManager(testLogger)
	if tmErr != nil {
		return nil, pkgerrors.Wrap(tmErr, "failed to initialize tunnel manager")
	}
	var relayManager *componentRelayManager
	if !strings.EqualFold(strings.TrimSpace(os.Getenv(envUsePersistentRelaySupervisor)), "true") {
		rm, rmErr := newComponentRelayManager(testLogger)
		if rmErr != nil && nodeSetPlacement.HasRemoteTargets {
			return nil, pkgerrors.Wrap(rmErr, "failed to initialize relay manager")
		}
		relayManager = rm
	} else {
		testLogger.Info().Msg("persistent relay supervisor enabled; skipping in-process relay manager")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting %d blockchain(s)", len(input.Blockchains))))

	deployedBlockchains, startErr := startBlockchainsWithTargets(
		ctx,
		testLogger,
		input.Blockchains,
		input.BlockchainDeployers,
		tunnelManager,
		nodeSetPlacement.HasLocalTargets,
	)
	if startErr != nil {
		return nil, pkgerrors.Wrap(startErr, "failed to start blockchains")
	}
	cleanupTunnelsOnError := true
	cleanupRelaysOnError := true
	defer func() {
		if cleanupTunnelsOnError {
			_ = tunnelManager.Stop(ctx)
		}
		if cleanupRelaysOnError && relayManager != nil {
			_ = relayManager.Close(ctx)
		}
	}()

	creEnvironment := &cre.Environment{
		Blockchains:           deployedBlockchains.Outputs,
		ContractVersions:      input.ContractVersions,
		Provider:              input.Provider,
		RegistryChainSelector: deployedBlockchains.RegistryChain().ChainSelector(),
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Blockchains started in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Deploying Workflow and Capability Registry contracts")))

	deployKeystoneContractsOutput, deployErr := crecontracts.DeployKeystoneContracts(
		ctx,
		testLogger,
		singleFileLogger,
		crecontracts.DeployKeystoneContractsInput{
			CldfEnvironment:  newCldfEnvironment(ctx, singleFileLogger, deployedBlockchains.CldfBlockChains),
			CtfBlockchains:   deployedBlockchains.Outputs,
			ContractVersions: input.ContractVersions,
			WithV2Registries: input.WithV2Registries,
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
	blockchainTargetBySelector := blockchainTargetsBySelector(input.Blockchains, deployedBlockchains.Outputs)

	updatedNodeSets, topoErr := donconfig.PrepareNodeTOMLs(
		ctx,
		topology,
		creEnvironment,
		input.NodeSets,
		blockchainTargetBySelector,
		input.Capabilities,
		input.ConfigFactoryFunctions,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	if nodeSetPlacement.HasRemoteTargets && relayManager != nil {
		if err := ensureMixedRelaysForLocalBlockchains(ctx, relayManager, input.Blockchains, deployedBlockchains.Outputs); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to ensure mixed relays for local blockchains")
		}
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Applying Features before environment startup")))
	var donsCapabilities = make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	var capabilityToOCR3Config = make(map[string]*ocr3.OracleConfig)
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
			}
			testLogger.Info().Msgf("PreEnvStartup for feature %s executed successfully", feature.Flag())
		}
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Applied Features in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	// Start JD first when we need to expose local JD endpoints to remote nodesets.
	requireJDRelayBootstrap := nodeSetPlacement.HasRemoteTargets && input.JdInput != nil && input.JdInput.Target == config.TargetLocal
	startedJD, jdStartErr := StartJD(ctx, testLogger, input.JdInput, input.Provider, tunnelManager, nodeSetPlacement.HasLocalTargets)
	if jdStartErr != nil {
		return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
	}
	if requireJDRelayBootstrap && relayManager != nil {
		if err := ensureMixedRelaysForLocalJD(ctx, relayManager, startedJD.JDOutput); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to ensure mixed relays for local JD")
		}
	}
	if input.PreDONsStartHook != nil {
		if err := input.PreDONsStartHook(ctx); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to execute pre-DON startup hook")
		}
	}

	startedDONs, donStartErr := StartDONs(ctx, testLogger, topology, input.Provider, deployedBlockchains.RegistryChain().CtfOutput(), input.CapabilityConfigs, input.CopyCapabilityBinaries, updatedNodeSets, tunnelManager)
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
	}

	cldErr := cre.LinkToJobDistributor(ctx, linkDonsToJDInput)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to link DONs to Job Distributor")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs and Job Distributor started and linked in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Creating Jobs with Job Distributor")))

	gJobErr := gateway.CreateJobs(ctx, creEnvironment, dons, topology.GatewayConfigs, input.GatewayWhitelistConfig)
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
	wfRegVersion := input.ContractVersions[keystone_changeset.WorkflowRegistry.String()]
	workflowRegistryConfigurationOutput, wfErr := workflow.ConfigureWorkflowRegistry(
		ctx,
		testLogger,
		singleFileLogger,
		&cre.WorkflowRegistryInput{
			ContractAddress: common.HexToAddress(crecontracts.MustGetAddressFromDataStore(deployKeystoneContractsOutput.Env.DataStore, deployedBlockchains.RegistryChain().ChainSelector(), keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")),
			ContractVersion: cldf.TypeAndVersion{Version: *wfRegVersion},
			ChainSelector:   deployedBlockchains.RegistryChain().ChainSelector(),
			CldEnv:          deployKeystoneContractsOutput.Env,
			AllowedDonIDs:   topology.WorkflowDONIDs,
			WorkflowOwners:  []common.Address{deployedBlockchains.RegistryChain().(*evm.Blockchain).SethClient.MustGetRootKeyAddress()}, // registry chain is always EVM
		},
	)
	if wfErr != nil {
		return nil, pkgerrors.Wrap(wfErr, "failed to configure workflow registry")
	}

	waitForWorkflowFilters := func(ctx context.Context) error {
		// we currently have no way of checking if filters were registered in Kubernetes mode
		// as we don't have a way to get its database connection string
		if !input.Provider.IsDocker() {
			return nil
		}

		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Waiting for Workflow Registry filters registration\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Finished waiting for Workflow Registry filters registration\n\n"))

		// this operation can always safely run in the background, since it doesn't change on-chain state, it only reads data from databases
		switch wfRegVersion.Major() {
		case 2:
			// There are no filters registered with the V2 WF Registry Syncer
			return nil
		default:
			return workflow.WaitForAllNodesToHaveExpectedFiltersRegistered(ctx, singleFileLogger, testLogger, deployedBlockchains.RegistryChain().ChainID(), dons, updatedNodeSets)
		}
	}

	capRegInput := cre.ConfigureCapabilityRegistryInput{
		ChainSelector: deployedBlockchains.RegistryChain().ChainSelector(),
		CldEnv:        creEnvironment.CldfEnvironment,
		Blockchains:   deployedBlockchains.Outputs,
		Topology:      topology,
		CapabilitiesRegistryAddress: ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(
			deployKeystoneContractsOutput.MemoryDataStore,
			deployedBlockchains.RegistryChain().ChainSelector(),
			keystone_changeset.CapabilitiesRegistry.String(),
			input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
			""),
		),
		NodeSets:                 input.NodeSets,
		WithV2Registries:         input.WithV2Registries,
		DONCapabilityWithConfigs: make(map[uint64][]keystone_changeset.DONCapabilityWithConfig),
		CapabilityToOCR3Config:   capabilityToOCR3Config,
	}

	for _, capability := range input.Capabilities {
		configFn := capability.CapabilityRegistryV1ConfigFn()
		capRegInput.CapabilityRegistryConfigFns = append(capRegInput.CapabilityRegistryConfigFns, configFn)
	}
	capRegInput.CapabilityRegistryConfigFns = append(capRegInput.CapabilityRegistryConfigFns, input.CapabilitiesContractFactoryFunctions...)
	maps.Copy(capRegInput.DONCapabilityWithConfigs, donsCapabilities)

	_, capRegErr := crecontracts.ConfigureCapabilityRegistry(capRegInput)
	if capRegErr != nil {
		return nil, pkgerrors.Wrap(capRegErr, "failed to configure Capability Registry contracts")
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

	if err := waitForWorkflowFilters(ctx); err != nil {
		return nil, pkgerrors.Wrap(err, "failed while waiting for workflow registry filters registration")
	}

	appendOutputsToInput(input, startedDONs.NodeOutputs(), deployedBlockchains.Outputs, startedJD.JDOutput)

	if err := workflowRegistryConfigurationOutput.Store(config.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to store workflow registry configuration output")
	}

	cleanupTunnelsOnError = false
	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryConfigurationOutput, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		Dons:                                dons,
		NodeOutput:                          startedDONs.NodeOutputs(),
		CreEnvironment:                      creEnvironment,
		S3ProviderOutput:                    s3Output,
		GatewayConnectors:                   topology.GatewayConnectors,
		tunnelManager:                       tunnelManager,
		relayManager:                        relayManager,
	}, nil
}

func ensureMixedRelaysForLocalBlockchains(
	ctx context.Context,
	relayManager *componentRelayManager,
	configuredBlockchains []*config.Blockchain,
	deployedBlockchains []blockchains.Blockchain,
) error {
	attempted := 0
	for idx, configured := range configuredBlockchains {
		if configured == nil || configured.Target != config.TargetLocal {
			continue
		}
		if idx >= len(deployedBlockchains) || deployedBlockchains[idx] == nil {
			continue
		}
		for nodeIdx, node := range deployedBlockchains[idx].CtfOutput().Nodes {
			if node == nil {
				continue
			}
			if p, ok := extractEndpointPort(node.ExternalHTTPUrl); ok {
				attempted++
				if err := relayManager.EnsurePort(ctx, fmt.Sprintf("blockchain-http-%d-%d", idx, nodeIdx), p); err != nil {
					return err
				}
			}
			if p, ok := extractEndpointPort(node.ExternalWSUrl); ok {
				attempted++
				if err := relayManager.EnsurePort(ctx, fmt.Sprintf("blockchain-ws-%d-%d", idx, nodeIdx), p); err != nil {
					return err
				}
			}
		}
	}
	if attempted == 0 {
		relayManager.lggr.Warn().Msg("no local blockchain relay ports were detected; mixed remote nodesets may not reach local blockchains")
	}
	return nil
}

func ensureMixedRelaysForLocalJD(ctx context.Context, relayManager *componentRelayManager, jdOutput *jd.Output) error {
	if jdOutput == nil {
		return nil
	}
	attempted := 0
	if p, ok := extractEndpointPort(jdOutput.ExternalGRPCUrl); ok {
		attempted++
		if err := relayManager.EnsurePort(ctx, "jd-grpc", p); err != nil {
			return err
		}
	}
	if p, ok := extractEndpointPort(jdOutput.ExternalWSRPCUrl); ok {
		attempted++
		if err := relayManager.EnsurePort(ctx, "jd-wsrpc", p); err != nil {
			return err
		}
	}
	if attempted == 0 {
		relayManager.lggr.Warn().Msg("no local JD relay ports were detected; mixed remote nodesets may not reach local JD")
	}
	return nil
}

func extractEndpointPort(raw string) (int, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil || parsed.Port() == "" {
			return 0, false
		}
		port, convErr := strconv.Atoi(parsed.Port())
		if convErr != nil || port <= 0 || port > 65535 {
			return 0, false
		}
		return port, true
	}
	_, portRaw, err := net.SplitHostPort(trimmed)
	if err != nil {
		return 0, false
	}
	port, convErr := strconv.Atoi(portRaw)
	if convErr != nil || port <= 0 || port > 65535 {
		return 0, false
	}
	return port, true
}

func blockchainTargetsBySelector(configured []*config.Blockchain, deployed []blockchains.Blockchain) map[uint64]string {
	bySelector := make(map[uint64]string, len(deployed))
	for idx, blockchainCfg := range configured {
		if blockchainCfg == nil {
			continue
		}
		if idx >= len(deployed) || deployed[idx] == nil {
			continue
		}
		selector := deployed[idx].ChainSelector()
		bySelector[selector] = string(blockchainCfg.Target)
	}
	return bySelector
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

type nodeSetPlacementSummary struct {
	HasLocalTargets  bool
	HasRemoteTargets bool
}

func summarizeNodeSetPlacement(nodeSets []*cre.NodeSet) (*nodeSetPlacementSummary, error) {
	summary := &nodeSetPlacementSummary{}
	for _, nodeSet := range nodeSets {
		if nodeSet == nil {
			continue
		}
		configTarget := strings.TrimSpace(nodeSet.Target)
		if configTarget == "" || configTarget == string(config.TargetLocal) {
			summary.HasLocalTargets = true
			continue
		}
		if configTarget == string(config.TargetRemote) {
			summary.HasRemoteTargets = true
			continue
		}
		return nil, fmt.Errorf("invalid nodeset target: %s", nodeSet.Target)
	}

	// Mixed local and remote nodeset targets need per-DON node-facing URL config selection.
	// Current PrepareNodeTOMLs builds one node-facing URL shape, so keep this unsupported for now.
	if summary.HasLocalTargets && summary.HasRemoteTargets {
		return nil, errors.New("mixed nodeset targets are not supported yet; set all nodesets target=local or all target=remote")
	}
	return summary, nil
}

func validateUnsupportedPlacements(
	configuredBlockchains []*config.Blockchain,
	nodeSetPlacement *nodeSetPlacementSummary,
) error {
	if nodeSetPlacement == nil || !nodeSetPlacement.HasRemoteTargets {
		return nil
	}
	for _, bc := range configuredBlockchains {
		if bc == nil {
			continue
		}
		if bc.Target == config.TargetLocal {
			return errors.New(
				"remote nodesets with local blockchains are not supported in this PoC. " +
					"Set all blockchains to target=remote, or run nodesets with target=local so nodes stay colocated with local blockchains",
			)
		}
	}
	return nil
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
