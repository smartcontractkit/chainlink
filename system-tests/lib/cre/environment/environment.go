package environment

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"

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
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	donconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	gateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/worker"
)

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cre.WorkflowRegistryOutput
	CreEnvironment                      *cre.Environment
	Dons                                *cre.Dons
	NodeOutput                          []*cre.WrappedNodeOutput
	S3ProviderOutput                    *s3provider.Output
	GatewayConnectors                   *cre.GatewayConnectors
}

type SetupInput struct {
	NodeSets               []*cre.NodeSet
	BlockchainsInput       []*blockchain.Input
	JdInput                *jd.Input
	Provider               infra.Provider
	ContractVersions       map[string]string
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
}

func (s *SetupInput) Validate() error {
	if s == nil {
		return pkgerrors.New("input is nil")
	}

	if len(s.NodeSets) == 0 {
		return pkgerrors.New("at least one nodeSet is required")
	}

	if len(s.BlockchainsInput) == 0 {
		return pkgerrors.New("at least one blockchain is required")
	}

	if s.JdInput == nil {
		return pkgerrors.New("jd input is nil")
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

	if input.Provider.Type == infra.CRIB {
		cribErr := crib.Bootstrap(ctx, input.Provider)
		if cribErr != nil {
			return nil, pkgerrors.Wrap(cribErr, "failed to bootstrap CRIB")
		}
	}

	s3Output, s3Err := workflow.StartS3(testLogger, input.S3ProviderInput, input.StageGen)
	if s3Err != nil {
		return nil, pkgerrors.Wrap(s3Err, "failed to start S3 provider")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting %d blockchain(s)", len(input.BlockchainsInput))))

	deployedBlockchains, startErr := blockchains.Start(
		ctx,
		testLogger,
		singleFileLogger,
		input.BlockchainsInput,
		input.BlockchainDeployers,
	)
	if startErr != nil {
		return nil, pkgerrors.Wrap(startErr, "failed to start blockchains")
	}

	creEnvironment := &cre.Environment{
		Blockchains:           deployedBlockchains.Outputs,
		ContractVersions:      input.ContractVersions,
		Provider:              input.Provider,
		CapabilityConfigs:     input.CapabilityConfigs,
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

	topology, tErr := cre.NewTopology(input.NodeSets, creEnvironment.Provider)
	if tErr != nil {
		return nil, pkgerrors.Wrap(tErr, "failed to create topology")
	}

	updatedNodeSets, topoErr := donconfig.PrepareNodeTOMLs(
		topology,
		creEnvironment,
		input.NodeSets,
		input.Capabilities,
		input.ConfigFactoryFunctions,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Applying Features before environment startup")))
	var donsCapabilities = make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
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
			}
			testLogger.Info().Msgf("PreEnvStartup for feature %s executed successfully", feature.Flag())
		}
	}
	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Applied Features in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	queue := worker.New(ctx, 10)
	defer queue.StopAndWait() // Ensure cleanup on any exit path

	jdStartedFuture := queue.SubmitAny(func(ctx context.Context) (any, error) {
		// TODO: pass context after we update the CTF to accept context, when creating new JD instance
		jdOutput, startJDErr := StartJD(ctx, testLogger, *input.JdInput, input.Provider)
		if startJDErr != nil {
			return nil, pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}
		return jdOutput, nil
	})

	donsStartedFuture := queue.SubmitAny(func(ctx context.Context) (any, error) {
		nodeSetOutput, startDonsErr := StartDONs(ctx, testLogger, topology, input.Provider, deployedBlockchains.RegistryChain().CtfOutput(), input.CapabilityConfigs, input.CopyCapabilityBinaries, updatedNodeSets)
		if startDonsErr != nil {
			return nil, pkgerrors.Wrap(startDonsErr, "failed to start DONs")
		}

		return nodeSetOutput, nil
	})

	// Await both futures to ensure proper cleanup even if one fails
	startedJD, jdStartErr := worker.AwaitAs[*StartedJD](ctx, jdStartedFuture)
	startedDONs, donStartErr := worker.AwaitAs[*StartedDONs](ctx, donsStartedFuture)

	// Check errors after both awaits complete
	// If both failed, prefer the non-context-cancelled error as it's likely the root cause
	if jdStartErr != nil && donStartErr != nil {
		// If one is context.Canceled, it was likely caused by the other task's error
		if pkgerrors.Is(jdStartErr, context.Canceled) && !pkgerrors.Is(donStartErr, context.Canceled) {
			return nil, pkgerrors.Wrap(donStartErr, "failed to start DONs")
		}
		if pkgerrors.Is(donStartErr, context.Canceled) && !pkgerrors.Is(jdStartErr, context.Canceled) {
			return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
		}
		// Both real errors
		return nil, pkgerrors.Wrap(errors.Join(fmt.Errorf("JD failed to start: %w", jdStartErr), fmt.Errorf("DONs failed to start: %w", donStartErr)), "failed to start Job Distributor AND Dons")
	}
	if jdStartErr != nil {
		return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
	}
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
	wfRegVersion := *semver.MustParse(input.ContractVersions[keystone_changeset.WorkflowRegistry.String()])
	workflowRegistryConfigurationOutput, wfErr := workflow.ConfigureWorkflowRegistry(
		ctx,
		testLogger,
		singleFileLogger,
		&cre.WorkflowRegistryInput{
			ContractAddress: common.HexToAddress(crecontracts.MustGetAddressFromDataStore(deployKeystoneContractsOutput.Env.DataStore, deployedBlockchains.RegistryChain().ChainSelector(), keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")),
			ContractVersion: cldf.TypeAndVersion{Version: wfRegVersion},
			ChainSelector:   deployedBlockchains.RegistryChain().ChainSelector(),
			CldEnv:          deployKeystoneContractsOutput.Env,
			AllowedDonIDs:   []uint64{topology.WorkflowDONID},
			WorkflowOwners:  []common.Address{deployedBlockchains.RegistryChain().(*evm.Blockchain).SethClient.MustGetRootKeyAddress()}, // registry chain is always EVM
		},
	)
	if wfErr != nil {
		return nil, pkgerrors.Wrap(wfErr, "failed to configure workflow registry")
	}

	wfFiltersFuture := queue.SubmitErr(func(ctx context.Context) error {
		// we currently have no way of checking if filters were registered, when code runs in CRIB
		// as we don't have a way to get its database connection string
		if input.Provider.Type == infra.CRIB {
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
	})

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

	if err := worker.AwaitErr(ctx, wfFiltersFuture); err != nil {
		return nil, pkgerrors.Wrap(err, "failed while waiting for workflow registry filters registration")
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
		S3ProviderOutput:                    s3Output,
		GatewayConnectors:                   topology.GatewayConnectors,
	}, nil
}

func appendOutputsToInput(input *SetupInput, nodeSetOutput []*cre.WrappedNodeOutput, blockchains []blockchains.Blockchain, jdOutput *jd.Output) {
	// append the nodeset output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	for idx, nsOut := range nodeSetOutput {
		input.NodeSets[idx].Out = nsOut.Output
	}

	for idx, blockchain := range blockchains {
		input.BlockchainsInput[idx].Out = blockchain.CtfOutput()
	}

	// append the jd output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	input.JdInput.Out = jdOutput
}

func newCldfEnvironment(ctx context.Context, singleFileLogger logger.Logger, cldfBlockchains cldf_chain.BlockChains) *cldf.Environment {
	memoryDatastore := datastore.NewMemoryDataStore()
	allChainsCLDEnvironment := &cldf.Environment{
		Name:              cre.EnvironmentName,
		Logger:            singleFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         memoryDatastore.Seal(),
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
