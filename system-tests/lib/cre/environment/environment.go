package environment

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	donconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/worker"
)

const (
	GithubReadTokenEnvVarName          = "GITHUB_READ_TOKEN"
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
	cribConfigsDir                     = "crib-configs"
)

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cre.WorkflowRegistryOutput
	CldEnvironment                      *cldf.Environment
	BlockchainOutput                    []*cre.WrappedBlockchainOutput
	DonTopology                         *cre.DonTopology
	NodeOutput                          []*cre.WrappedNodeOutput
	InfraInput                          infra.Provider
	S3ProviderOutput                    *s3provider.Output
}

type SetupInput struct {
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
	BlockchainsInput          []blockchain.Input
	JdInput                   *jd.Input
	Provider                  infra.Provider
	ContractVersions          map[string]string
	WithV2Registries          bool
	OCR3Config                *keystone_changeset.OracleConfig
	DONTimeConfig             *keystone_changeset.OracleConfig
	VaultOCR3Config           *keystone_changeset.OracleConfig
	S3ProviderInput           *s3provider.Input
	CapabilityConfigs         cre.CapabilityConfigs
	CopyCapabilityBinaries    bool // if true, copy capability binaries to the containers (if false, we assume that the plugins image already has them)
	Capabilities              []cre.InstallableCapability

	// Deprecated: use Capabilities []cre.InstallableCapability instead
	ConfigFactoryFunctions []cre.NodeConfigTransformerFn
	// Deprecated: use Capabilities []cre.InstallableCapability instead
	JobSpecFactoryFunctions []cre.JobSpecFn
	// Deprecated: use Capabilities []cre.InstallableCapability instead
	CapabilitiesContractFactoryFunctions []cre.CapabilityRegistryConfigFn

	StageGen *stagegen.StageGen
}

func (s *SetupInput) Validate() error {
	if s == nil {
		return pkgerrors.New("input is nil")
	}

	if len(s.CapabilitiesAwareNodeSets) == 0 {
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

	if err := input.Validate(); err != nil {
		return nil, pkgerrors.Wrap(err, "input validation failed")
	}

	if input.Provider.Type == infra.CRIB {
		cribErr := crib.Bootstrap(input.Provider)
		if cribErr != nil {
			return nil, pkgerrors.Wrap(cribErr, "failed to bootstrap CRIB")
		}
	}

	s3Output, s3Err := workflow.StartS3(testLogger, input.S3ProviderInput, input.StageGen)
	if s3Err != nil {
		return nil, pkgerrors.Wrap(s3Err, "failed to start S3 provider")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting %d blockchain(s)", len(input.BlockchainsInput))))

	startBlockchainsOutput, bcOutErr := StartBlockchains(BlockchainLoggers{
		lggr:       testLogger,
		singleFile: singleFileLogger,
	}, BlockchainsInput{
		infra:            input.Provider,
		blockchainsInput: input.BlockchainsInput,
	})
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to start blockchains")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Blockchains started in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Deploying Keystone contracts")))

	deployKeystoneContractsOutput, deployErr := crecontracts.DeployKeystoneContracts(
		ctx,
		testLogger,
		singleFileLogger,
		crecontracts.DeployKeystoneContractsInput{
			CldfEnvironment:           newCldfEnvironment(ctx, singleFileLogger, startBlockchainsOutput.BlockChains),
			CtfBlockchains:            startBlockchainsOutput.BlockChainOutputs,
			ContractVersions:          input.ContractVersions,
			WithV2Registries:          input.WithV2Registries,
			CapabilitiesAwareNodeSets: input.CapabilitiesAwareNodeSets,
		},
	)
	if deployErr != nil {
		return nil, pkgerrors.Wrap(deployErr, "failed to deploy Keystone contracts")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Keystone contracts deployed in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Preparing DONs configuration")))

	topology, updatedNodeSets, topoErr := donconfig.PrepareNodeTOMLs(
		startBlockchainsOutput.RegistryChain().ChainSelector,
		input.CapabilitiesAwareNodeSets,
		input.Provider,
		startBlockchainsOutput.BlockChainOutputs,
		deployKeystoneContractsOutput.Env.ExistingAddresses, //nolint:staticcheck // won't migrate now
		deployKeystoneContractsOutput.Env.DataStore,
		input.Capabilities,
		input.CapabilityConfigs,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting Job Distributor and DONs and linking them to JD")))

	queue := worker.New(10)
	jdStartedFuture := queue.SubmitAny(func() (any, error) {
		jdOutput, startJDErr := StartJD(testLogger, *input.JdInput, input.Provider)
		if startJDErr != nil {
			return nil, pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}
		return jdOutput, nil
	})

	donsStartedFuture := queue.SubmitAny(func() (any, error) {
		nodeSetOutput, startDonsErr := StartDONs(ctx, testLogger, topology, input.Provider, startBlockchainsOutput.RegistryChain().BlockchainOutput, input.CapabilityConfigs, input.CopyCapabilityBinaries, updatedNodeSets)
		if startDonsErr != nil {
			return nil, pkgerrors.Wrap(startDonsErr, "failed to start DONs")
		}

		return nodeSetOutput, nil
	})

	// First wait for JD to start, because it will be faster than DONs
	startedJD, jdStartErr := worker.AwaitAs[*StartedJD](ctx, jdStartedFuture)
	if jdStartErr != nil {
		return nil, pkgerrors.Wrap(jdStartErr, "failed to start Job Distributor")
	}

	startedDONs, donStartErr := worker.AwaitAs[*StartedDONs](ctx, donsStartedFuture)
	if donStartErr != nil {
		return nil, pkgerrors.Wrap(donStartErr, "failed to start DONs")
	}

	linkDonsToJDInput := &cre.LinkDonsToJDInput{
		JDClient:          startedJD.Client,
		BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
		CldfEnvironment:   deployKeystoneContractsOutput.Env,
		Topology:          topology,
		DONs:              startedDONs.DONs(),
	}

	cldfEnvironment, cldErr := cre.LinkToJobDistributor(ctx, linkDonsToJDInput)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to link DONs to Job Distributor")
	}
	creEnvironment := &cre.Environment{
		CldfEnvironment: cldfEnvironment,
		DonTopology:     cre.NewDonTopology(startBlockchainsOutput.RegistryChain().ChainSelector, topology, cre.NewDons(startedDONs.DONs())),
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("DONs and Job Distributor started and linked in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Creating Jobs with Job Distributor")))

	jobSpecFactoryFunctions := make([]cre.JobSpecFn, 0)
	for _, capability := range input.Capabilities {
		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, capability.JobSpecFn())
	}

	jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, input.JobSpecFactoryFunctions...) // Deprecated, use Capabilities instead

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the OCR3 contracts.
	createJobsDeps := CreateJobsWithJdOpDeps{
		Logger:                    testLogger,
		SingleFileLogger:          singleFileLogger,
		HomeChainBlockchainOutput: startBlockchainsOutput.RegistryChain().BlockchainOutput,
		JobSpecFactoryFunctions:   jobSpecFactoryFunctions,
		CreEnvironment:            creEnvironment,
		CapabilitiesAwareNodeSets: input.CapabilitiesAwareNodeSets,
		InfraInput:                input.Provider,
		CapabilitiesConfigs:       input.CapabilityConfigs,
		Capabilities:              input.Capabilities,
	}
	_, createJobsErr := operations.ExecuteOperation(deployKeystoneContractsOutput.Env.OperationsBundle, CreateJobsWithJdOp, createJobsDeps, CreateJobsWithJdOpInput{})
	if createJobsErr != nil {
		return nil, pkgerrors.Wrap(createJobsErr, "failed to create jobs with Job Distributor")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Jobs created in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Preparing Chainlink Node funding")))

	// This operation cannot execute in the background, because it uses master private key and we want to avoid nonce issues
	// Once we have generated and funded new private keys for each chain, we can execute fanning out of funds to nodes in the background
	preFundingOutput, prefundErr := operations.ExecuteOperation(creEnvironment.CldfEnvironment.OperationsBundle, PrepareCLNodesFundingOp, PrepareFundCLNodesOpDeps{
		TestLogger:        testLogger,
		Env:               creEnvironment.CldfEnvironment,
		BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
		DonTopology:       creEnvironment.DonTopology,
	}, PrepareFundCLNodesOpInput{FundingPerChainFamilyForEachNode: map[string]uint64{
		chainselectors.FamilyEVM:    10000000000000000, // 0.01 ETH
		chainselectors.FamilySolana: 50_000_000_000,    // 50 SOL
		chainselectors.FamilyTron:   100_000_000,       // 100 TRX in SUN
	}})
	if prefundErr != nil {
		return nil, pkgerrors.Wrap(prefundErr, "failed to prepare funding of CL nodes")
	}

	fundNodesFuture := queue.SubmitErr(func() error {
		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Funding Chainlink nodes\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Finished Funding Chainlink nodes\n\n"))

		_, fundErr := operations.ExecuteOperation(creEnvironment.CldfEnvironment.OperationsBundle, FundCLNodesOp, FundCLNodesOpDeps{
			TestLogger:        testLogger,
			Env:               creEnvironment.CldfEnvironment,
			BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
			DonTopology:       creEnvironment.DonTopology,
		}, FundCLNodesOpInput{
			FundingAmountPerChainFamily: preFundingOutput.Output.FundingPerChainFamilyForEachNode,
			PrivateKeyPerChainFamily:    preFundingOutput.Output.PrivateKeysPerChainFamily,
		})

		return fundErr
	})

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Chainlink Node funding prepared in %.2f seconds", input.StageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Waiting for Log Poller to start tracking OCR3 contract")))

	// Wait for Log Poller to be up and running. If it misses the ConfigSet event, OCR protocol will not start.
	// TODO: we might want to add similar checks for other OCR3 contracts.
	if err := waitForLogPollerToBeHealthy(updatedNodeSets, startedDONs.NodeOutputs()); err != nil {
		return nil, pkgerrors.Wrap(err, "failed while waiting for Log Poller to become healthy")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Log Poller started in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Starting Workflow Registry Contract configuration")))

	wfRegVersion := *semver.MustParse(input.ContractVersions[keystone_changeset.WorkflowRegistry.String()])
	workflowRegistryConfigurationOutput, wfErr := workflow.ConfigureWorkflowRegistry(
		ctx,
		testLogger,
		singleFileLogger,
		&cre.WorkflowRegistryInput{
			ContractAddress: common.HexToAddress(crecontracts.MustGetAddressFromDataStore(deployKeystoneContractsOutput.Env.DataStore, startBlockchainsOutput.RegistryChain().ChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")),
			ContractVersion: cldf.TypeAndVersion{Version: wfRegVersion},
			ChainSelector:   startBlockchainsOutput.RegistryChain().ChainSelector,
			CldEnv:          deployKeystoneContractsOutput.Env,
			AllowedDonIDs:   []uint64{topology.WorkflowDONID},
			WorkflowOwners:  []common.Address{startBlockchainsOutput.RegistryChain().SethClient.MustGetRootKeyAddress()},
		},
	)

	if wfErr != nil {
		return nil, pkgerrors.Wrap(wfErr, "failed to configure workflow registry")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("Workflow Registry Contract configured in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	wfFiltersFuture := queue.SubmitErr(func() error {
		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Waiting for Workflow Registry filters registration\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Finished waiting for Workflow Registry filters registration\n\n"))

		// this operation can always safely run in the background, since it doesn't change on-chain state, it only reads data from databases
		switch wfRegVersion.Major() {
		case 2:
			// There are no filters registered with the V2 WF Registry Syncer
			return nil
		default:
			return workflow.WaitForWorkflowRegistryFiltersRegistration(testLogger, singleFileLogger, input.Provider.Type, startBlockchainsOutput.RegistryChain().ChainID, creEnvironment.DonTopology, updatedNodeSets)
		}
	})

	fmt.Print(libformat.PurpleText("%s", input.StageGen.Wrap("Configuring OCR3 and Keystone contracts")))

	configureKeystoneInput, ksErr := prepareKeystoneConfigurationInput(*input, startBlockchainsOutput.RegistryChain().ChainSelector, topology, updatedNodeSets, creEnvironment.CldfEnvironment, deployKeystoneContractsOutput, startBlockchainsOutput)
	if ksErr != nil {
		return nil, pkgerrors.Wrap(ksErr, "failed to prepare keystone configuration input")
	}

	keystoneErr := crecontracts.ConfigureKeystone(*configureKeystoneInput)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	fmt.Print(libformat.PurpleText("%s", input.StageGen.WrapAndNext("OCR3 and Keystone contracts configured in %.2f seconds", input.StageGen.Elapsed().Seconds())))

	if err := worker.AwaitErr(ctx, fundNodesFuture); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to fund chainlink nodes")
	}

	if err := worker.AwaitErr(ctx, wfFiltersFuture); err != nil {
		return nil, pkgerrors.Wrap(err, "failed while waiting for workflow registry filters registration")
	}
	queue.StopAndWait()

	appendOutputsToInput(input, startedDONs.NodeOutputs(), startBlockchainsOutput, startedJD.JDOutput)

	if err := workflowRegistryConfigurationOutput.Store(config.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to store workflow registry configuration output")
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryConfigurationOutput, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    startBlockchainsOutput.BlockChainOutputs,
		DonTopology:                         creEnvironment.DonTopology,
		NodeOutput:                          startedDONs.NodeOutputs(),
		CldEnvironment:                      creEnvironment.CldfEnvironment,
		S3ProviderOutput:                    s3Output,
	}, nil
}

func evmOCR3AddressesFromDataStore(blockchains []*cre.WrappedBlockchainOutput, nodeSets []*cre.CapabilitiesAwareNodeSet, ds *datastore.MemoryDataStore, homeChainSelector uint64) map[uint64]common.Address {
	chainsWithEVMCapability := crecontracts.ChainsWithEVMCapability(blockchains, nodeSets)
	evmOCR3CommonAddresses := make(map[uint64]common.Address)
	for chainID := range chainsWithEVMCapability {
		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
		// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
		// because that violates a DB constraint for offchain reporting jobs
		evmOCR3Addr := crecontracts.MustGetAddressFromMemoryDataStore(ds, homeChainSelector, keystone_changeset.OCR3Capability.String(), "1.0.0", qualifier)
		evmOCR3CommonAddresses[homeChainSelector] = evmOCR3Addr
	}

	return evmOCR3CommonAddresses
}

func mergeJobSpecSlices(from, to cre.DonsToJobSpecs) {
	for fromDonID, fromJobSpecs := range from {
		if _, ok := to[fromDonID]; !ok {
			to[fromDonID] = make([]*jobv1.ProposeJobRequest, 0)
		}
		to[fromDonID] = append(to[fromDonID], fromJobSpecs...)
	}
}

func prepareKeystoneConfigurationInput(input SetupInput, homeChainSelector uint64, topology *cre.Topology, updatedNodeSets []*cre.CapabilitiesAwareNodeSet, cldEnvironment *cldf.Environment, deployKeystoneContractsOutput *crecontracts.DeployKeystoneContractsOutput, startBlockchainsOutput StartBlockchainsOutput) (*cre.ConfigureKeystoneInput, error) {
	configureKeystoneInput := cre.ConfigureKeystoneInput{
		ChainSelector:               homeChainSelector,
		CldEnv:                      cldEnvironment,
		BlockchainOutputs:           startBlockchainsOutput.BlockChainOutputs,
		Topology:                    topology,
		CapabilitiesRegistryAddress: ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")),
		OCR3Address:                 ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.OCR3ContractQualifier)),
		DONTimeAddress:              ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.DONTimeContractQualifier)),
		VaultOCR3Address:            crecontracts.MightGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.VaultOCR3ContractQualifier+"_plugin"),
		DKGOCR3Address:              crecontracts.MightGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.VaultOCR3ContractQualifier+"_dkg"),
		EVMOCR3Addresses:            evmOCR3AddressesFromDataStore(startBlockchainsOutput.BlockChainOutputs, updatedNodeSets, deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector),
		ConsensusV2OCR3Address:      crecontracts.MightGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.ConsensusV2ContractQualifier),
		NodeSets:                    input.CapabilitiesAwareNodeSets,
		WithV2Registries:            input.WithV2Registries,
	}

	if input.OCR3Config != nil {
		configureKeystoneInput.OCR3Config = *input.OCR3Config
	} else {
		ocr3Config, ocr3ConfigErr := crecontracts.DefaultOCR3Config(topology)
		if ocr3ConfigErr != nil {
			return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
		}
		configureKeystoneInput.OCR3Config = *ocr3Config
	}

	if input.DONTimeConfig != nil {
		configureKeystoneInput.DONTimeConfig = *input.DONTimeConfig
	} else {
		donTimeConfig, donTimeConfigErr := crecontracts.DefaultOCR3Config(topology)
		donTimeConfig.DeltaRoundMillis = 0 // Fastest rounds possible
		if donTimeConfigErr != nil {
			return nil, pkgerrors.Wrap(donTimeConfigErr, "failed to generate default DON Time config")
		}
		configureKeystoneInput.DONTimeConfig = *donTimeConfig
	}

	if configureKeystoneInput.VaultOCR3Address != nil && configureKeystoneInput.VaultOCR3Address.Cmp(common.Address{}) != 0 {
		ocr3Config, ocr3ConfigErr := crecontracts.DefaultOCR3Config(topology)
		if ocr3ConfigErr != nil {
			return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
		}
		configureKeystoneInput.VaultOCR3Config = *ocr3Config

		dkgReportingPluginConfig, err := crecontracts.DKGReportingPluginConfig(topology, input.CapabilitiesAwareNodeSets)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to generate DKG reporting plugin config")
		}
		configureKeystoneInput.DKGReportingPluginConfig = dkgReportingPluginConfig
		configureKeystoneInput.DKGOCR3Config = *ocr3Config
	}

	chainOCR3Config, chainOCR3ConfigErr := crecontracts.DefaultChainCapabilityOCR3Config(topology)
	if chainOCR3ConfigErr != nil {
		return nil, pkgerrors.Wrap(chainOCR3ConfigErr, "failed to generate default Chain OCR3 config")
	}

	configureKeystoneInput.EVMOCR3Config = *chainOCR3Config

	defaultOcr3Config, defaultOcr3ConfigErr := crecontracts.DefaultOCR3Config(topology)
	if defaultOcr3ConfigErr != nil {
		return nil, pkgerrors.Wrap(defaultOcr3ConfigErr, "failed to generate default OCR3 config for EVM")
	}
	configureKeystoneInput.ConsensusV2OCR3Config = *defaultOcr3Config

	for _, capability := range input.Capabilities {
		configFn := capability.CapabilityRegistryV1ConfigFn()
		configureKeystoneInput.CapabilityRegistryConfigFns = append(configureKeystoneInput.CapabilityRegistryConfigFns, configFn)
	}

	// Deprecated, use Capabilities instead
	configureKeystoneInput.CapabilityRegistryConfigFns = append(configureKeystoneInput.CapabilityRegistryConfigFns, input.CapabilitiesContractFactoryFunctions...)

	return &configureKeystoneInput, nil
}

func waitForLogPollerToBeHealthy(nodeSetInput []*cre.CapabilitiesAwareNodeSet, nodeSetOutput []*cre.WrappedNodeOutput) error {
	for idx, nodeSetOut := range nodeSetOutput {
		if !flags.HasFlag(nodeSetInput[idx].ComputedCapabilities, cre.ConsensusCapability) || !flags.HasFlag(nodeSetInput[idx].ComputedCapabilities, cre.VaultCapability) {
			continue
		}
		nsClients, cErr := clclient.New(nodeSetOut.CLNodes)
		if cErr != nil {
			return pkgerrors.Wrap(cErr, "failed to create node set clients")
		}
		eg := &errgroup.Group{}
		for _, c := range nsClients {
			eg.Go(func() error {
				return c.WaitHealthy(".*ConfigWatcher", "passing", 100)
			})
		}
		if waitErr := eg.Wait(); waitErr != nil {
			return pkgerrors.Wrap(waitErr, "failed to wait for ConfigWatcher health check")
		}
	}

	return nil
}

func appendOutputsToInput(input *SetupInput, nodeSetOutput []*cre.WrappedNodeOutput, startBlockchainsOutput StartBlockchainsOutput, jdOutput *jd.Output) {
	// append the nodeset output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	for idx, nsOut := range nodeSetOutput {
		input.CapabilitiesAwareNodeSets[idx].Out = nsOut.Output
	}

	for idx, bcOut := range startBlockchainsOutput.BlockChainOutputs {
		input.BlockchainsInput[idx].Out = bcOut.BlockchainOutput
	}

	// append the jd output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	input.JdInput.Out = jdOutput
}

func newCldfEnvironment(ctx context.Context, singleFileLogger logger.Logger, cldfBlockchains map[uint64]cldf_chain.BlockChain) *cldf.Environment {
	memoryDatastore := datastore.NewMemoryDataStore()
	allChainsCLDEnvironment := &cldf.Environment{
		Name:              "local CRE",
		Logger:            singleFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         memoryDatastore.Seal(),
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: cldf_chain.NewBlockChains(cldfBlockchains),
		OCRSecrets:  focr.XXXGenerateTestOCRSecrets(),
	}
	allChainsCLDEnvironment.OperationsBundle = operations.NewBundle(allChainsCLDEnvironment.GetContext, singleFileLogger, operations.NewMemoryReporter())

	return allChainsCLDEnvironment
}
