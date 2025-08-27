package environment

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/jmoiron/sqlx"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/scylladb/go-reflectx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	ks_sol_seq "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence"
	ks_sol_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdevenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"
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
	InfraInput                          infra.Input
	S3ProviderOutput                    *s3provider.Output
}

type SetupInput struct {
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
	BlockchainsInput          []blockchain.Input
	JdInput                   jd.Input
	InfraInput                infra.Input
	ContractVersions          map[string]string
	OCR3Config                *keystone_changeset.OracleConfig
	DONTimeConfig             *keystone_changeset.OracleConfig
	VaultOCR3Config           *keystone_changeset.OracleConfig
	S3ProviderInput           *s3provider.Input
	CapabilityConfigs         cre.CapabilityConfigs
	CopyCapabilityBinaries    bool // if true, copy capability binaries to the containers (if false, we assume that the plugins image already has them)
	Capabilities              []cre.InstallableCapability

	// Deprecated: use Capabilities []cre.InstallableCapability instead
	ConfigFactoryFunctions []cre.NodeConfigFn
	// Deprecated: use Capabilities []cre.InstallableCapability instead
	JobSpecFactoryFunctions []cre.JobSpecFn
	// Deprecated: use Capabilities []cre.InstallableCapability instead
	CapabilitiesContractFactoryFunctions []cre.CapabilityRegistryConfigFn
}

type backgroundStageResult struct {
	err            error
	successMessage string
	panicValue     any
	panicStack     []byte
}

func mustGetAddress(dataStore datastore.MutableDataStore, chainSel uint64, contractType string, version string, qualifier string) string {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s (qualifier=%s) address for chain %d: %s", contractType, qualifier, chainSel, err.Error()))
	}
	return addrRef.Address
}

func SetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input SetupInput,
) (*SetupOutput, error) {
	topologyErr := libdon.ValidateTopology(input.CapabilitiesAwareNodeSets, input.InfraInput)
	if topologyErr != nil {
		return nil, pkgerrors.Wrap(topologyErr, "failed to validate topology")
	}

	// Shell is only required, when using CRIB, because we want to run commands in the same "nix develop" context
	// We need to have this reference in the outer scope, because subsequent functions will need it
	var nixShell *libnix.Shell
	if input.InfraInput.Type == infra.CRIB {
		startNixShellInput := &cre.StartNixShellInput{
			InfraInput:     &input.InfraInput,
			CribConfigsDir: cribConfigsDir,
			PurgeNamespace: true,
		}

		var nixErr error
		nixShell, nixErr = crib.StartNixShell(startNixShellInput)
		if nixErr != nil {
			return nil, pkgerrors.Wrap(nixErr, "failed to start nix shell")
		}
		// In CRIB v2 we no longer rely on devspace to create a namespace so we need to do it before deploying
		err := crib.Bootstrap(&input.InfraInput)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to create namespace")
		}
	}

	defer func() {
		if nixShell != nil {
			_ = nixShell.Close()
		}
	}()

	stageGen := NewStageGen(7, "STAGE")

	var s3ProviderOutput *s3provider.Output
	if input.S3ProviderInput != nil {
		stageGen = NewStageGen(8, "STAGE")
		fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting MinIO")))
		var s3ProviderErr error
		s3ProviderOutput, s3ProviderErr = s3provider.NewMinioFactory().NewFrom(input.S3ProviderInput)
		if s3ProviderErr != nil {
			return nil, pkgerrors.Wrap(s3ProviderErr, "minio provider creation failed")
		}
		testLogger.Debug().Msgf("S3Provider.Output value: %#v", s3ProviderOutput)
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("MinIO started in %.2f seconds", stageGen.Elapsed().Seconds())))
	}

	bi := BlockchainsInput{
		infra:    &input.InfraInput,
		nixShell: nixShell,
	}
	bi.blockchainsInput = append(bi.blockchainsInput, input.BlockchainsInput...)

	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting %d blockchain(s)", len(bi.blockchainsInput))))

	startBlockchainsOutput, bcOutErr := StartBlockchains(BlockchainLoggers{
		lggr:       testLogger,
		singleFile: singleFileLogger,
	}, bi)
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to start blockchains")
	}

	blockchainOutputs := startBlockchainsOutput.BlockChainOutputs
	homeChainOutput := blockchainOutputs[0]
	blockChains := startBlockchainsOutput.BlockChains

	memoryDatastore := datastore.NewMemoryDataStore()
	allChainsCLDEnvironment := &cldf.Environment{
		Logger:            singleFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         memoryDatastore.Seal(),
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: cldf_chain.NewBlockChains(blockChains),
	}
	allChainsCLDEnvironment.OperationsBundle = operations.NewBundle(allChainsCLDEnvironment.GetContext, singleFileLogger, operations.NewMemoryReporter())

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Blockchains started in %.2f seconds", stageGen.Elapsed().Seconds())))

	// DEPLOY CONTRACTS
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Deploying Keystone contracts")))

	evmForwardersSelectors := make([]uint64, 0)
	solForwardersSelectors := make([]uint64, 0)
	for _, bcOut := range blockchainOutputs {
		for _, donMetadata := range input.CapabilitiesAwareNodeSets {
			if slices.Contains(evmForwardersSelectors, bcOut.ChainSelector) {
				continue
			}
			if flags.RequiresForwarderContract(donMetadata.ComputedCapabilities, bcOut.ChainID) {
				evmForwardersSelectors = append(evmForwardersSelectors, bcOut.ChainSelector)
			}
		}
		if bcOut.SolChain != nil {
			// we expect that if we have solana, solana write capability is enabled
			solForwardersSelectors = append(solForwardersSelectors, bcOut.SolChain.ChainSelector)
			continue
		}
	}

	var allNodeFlags []string
	for i := range input.CapabilitiesAwareNodeSets {
		nodeFlags, err := flags.NodeSetFlags(input.CapabilitiesAwareNodeSets[i])
		if err != nil {
			continue
		}
		allNodeFlags = append(allNodeFlags, nodeFlags...)
	}
	vaultOCR3AddrFlag := flags.HasFlag(allNodeFlags, cre.VaultCapability)
	evmOCR3AddrFlag := flags.HasFlagForAnyChain(allNodeFlags, cre.EVMCapability)
	consensusV2AddrFlag := flags.HasFlag(allNodeFlags, cre.ConsensusCapabilityV2)

	chainsWithEVMCapability := make(map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector)
	for _, chain := range blockchainOutputs {
		for _, donMetadata := range input.CapabilitiesAwareNodeSets {
			if flags.HasFlagForChain(donMetadata.ComputedCapabilities, cre.EVMCapability, chain.ChainID) {
				if chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] != 0 {
					continue
				}
				chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] = ks_contracts_op.Selector(chain.ChainSelector)
			}
		}
	}

	// use CLD to deploy the necessary contracts
	homeChainSelector := homeChainOutput.ChainSelector
	deployKeystoneReport, err := operations.ExecuteSequence(
		allChainsCLDEnvironment.OperationsBundle,
		ks_contracts_op.DeployKeystoneContractsSequence,
		ks_contracts_op.DeployKeystoneContractsSequenceDeps{
			Env: allChainsCLDEnvironment,
		},
		ks_contracts_op.DeployKeystoneContractsSequenceInput{
			RegistryChainSelector: homeChainSelector,
			ForwardersSelectors:   evmForwardersSelectors,
			DeployVaultOCR3:       vaultOCR3AddrFlag,
			DeployEVMOCR3:         evmOCR3AddrFlag,
			EVMChainIDs:           chainsWithEVMCapability,
			DeployConsensusOCR3:   consensusV2AddrFlag,
		},
	)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to deploy Keystone contracts")
	}

	if err = allChainsCLDEnvironment.ExistingAddresses.Merge(deployKeystoneReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, pkgerrors.Wrap(err, "failed to merge address book with Keystone contracts addresses")
	}

	if err = memoryDatastore.Merge(deployKeystoneReport.Output.Datastore); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}

	for _, sel := range solForwardersSelectors {
		out, err := operations.ExecuteSequence(
			allChainsCLDEnvironment.OperationsBundle,
			ks_sol_seq.DeployForwarderSeq,
			ks_sol_op.Deps{
				Env:       *allChainsCLDEnvironment,
				Chain:     allChainsCLDEnvironment.BlockChains.SolanaChains()[sel],
				Datastore: memoryDatastore.Seal(),
			},
			ks_sol_seq.DeployForwarderSeqInput{
				ChainSel:    sel,
				ProgramName: deployment.KeystoneForwarderProgramName,
			},
		)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to deploy sol forwarder")
		}

		err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
			Address:       out.Output.ProgramID.String(),
			ChainSelector: sel,
			Version:       semver.MustParse(input.ContractVersions[ks_sol.ForwarderContract.String()]),
			Qualifier:     ks_sol.DefaultForwarderQualifier,
			Type:          ks_sol.ForwarderContract,
		})
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to add address to the datastore for Solana Forwarder contract")
		}

		err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
			Address:       out.Output.State.String(),
			ChainSelector: sel,
			Version:       semver.MustParse(input.ContractVersions[ks_sol.ForwarderState.String()]),
			Qualifier:     ks_sol.DefaultForwarderQualifier,
			Type:          ks_sol.ForwarderState,
		})
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to add address to the datastore for Solana Forwarder state")
		}

		testLogger.Info().Msgf("Deployed Forwarder contract on Solana chain chain %d programID: %s state: %s", sel, out.Output.ProgramID.String(), out.Output.State.String())
	}
	allChainsCLDEnvironment.DataStore = memoryDatastore.Seal()

	ocr3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "capability_ocr3")
	testLogger.Info().Msgf("Deployed OCR3 contract on chain %d at %s", homeChainSelector, ocr3Addr)

	donTimeAddr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "DONTime")
	testLogger.Info().Msgf("Deployed DON Time contract on chain %d at %s", homeChainSelector, donTimeAddr)

	wfRegAddr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Workflow Registry contract on chain %d at %s", homeChainSelector, wfRegAddr)

	capRegAddr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Capabilities Registry contract on chain %d at %s", homeChainSelector, capRegAddr)

	var vaultOCR3CommonAddr common.Address
	if vaultOCR3AddrFlag {
		vaultOCR3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "capability_vault")
		testLogger.Info().Msgf("Deployed Vault OCR3 contract on chain %d at %s", homeChainSelector, vaultOCR3Addr)
		vaultOCR3CommonAddr = common.HexToAddress(vaultOCR3Addr)
	}

	evmOCR3CommonAddresses := make(map[uint64]common.Address)
	if evmOCR3AddrFlag {
		for chainID, selector := range chainsWithEVMCapability {
			qualifier := ks_contracts_op.GetCapabilityContractIdentifier(uint64(chainID))
			evmOCR3Addr := mustGetAddress(memoryDatastore, uint64(selector), keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], qualifier)
			testLogger.Info().Msgf("Deployed EVM OCR3 contract on chainID: %d, selector: %d, at: %s", chainID, selector, evmOCR3Addr)
			evmOCR3CommonAddresses[uint64(selector)] = common.HexToAddress(evmOCR3Addr)
		}
	}
	var consensusV2OCR3CommonAddr common.Address
	if consensusV2AddrFlag {
		consensusV2OCR3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "capability_consensus")
		testLogger.Info().Msgf("Deployed Consensus V2 OCR3 contract on chain %d at %s", homeChainSelector, consensusV2OCR3Addr)
		consensusV2OCR3CommonAddr = common.HexToAddress(consensusV2OCR3Addr)
	}

	for _, forwarderSelector := range evmForwardersSelectors {
		forwarderAddr := mustGetAddress(memoryDatastore, forwarderSelector, keystone_changeset.KeystoneForwarder.String(), input.ContractVersions[keystone_changeset.KeystoneForwarder.String()], "")
		testLogger.Info().Msgf("Deployed Forwarder contract on chain %d at %s", forwarderSelector, forwarderAddr)
	}

	for _, forwarderSelector := range solForwardersSelectors {
		forwarderAddr := mustGetAddress(memoryDatastore, forwarderSelector, ks_sol.ForwarderContract.String(), input.ContractVersions[ks_sol.ForwarderContract.String()], ks_sol.DefaultForwarderQualifier)
		forwarderStateAddr := mustGetAddress(memoryDatastore, forwarderSelector, ks_sol.ForwarderState.String(), input.ContractVersions[ks_sol.ForwarderState.String()], ks_sol.DefaultForwarderQualifier)
		testLogger.Info().Msgf("Deployed Forwarder contract on Solana chain %d at %s state %s", forwarderSelector, forwarderAddr, forwarderStateAddr)
	}

	// get chainIDs, they'll be used for identifying ETH keys and Forwarder addresses
	// and also for creating the CLD environment
	evmChainIDs := make([]int, 0)
	bcOuts := make(map[uint64]*cre.WrappedBlockchainOutput)
	sethClients := make(map[uint64]*seth.Client)
	solClients := make(map[uint64]*solrpc.Client)
	solChainIDs := make([]string, 0)
	for _, bcOut := range blockchainOutputs {
		if bcOut.SolChain != nil {
			sel := bcOut.SolChain.ChainSelector
			bcOuts[sel] = bcOut
			solClients[sel] = bcOut.SolClient
			bcOuts[sel].ChainSelector = sel
			bcOuts[sel].SolChain = bcOut.SolChain
			bcOuts[sel].SolChain.ArtifactsDir = bcOut.SolChain.ArtifactsDir
			solChainIDs = append(solChainIDs, bcOut.SolChain.ChainID)
			continue
		}
		bcOuts[bcOut.ChainSelector] = bcOut
		evmChainIDs = append(evmChainIDs, libc.MustSafeInt(bcOut.ChainID))
		sethClients[bcOut.ChainSelector] = bcOut.SethClient
	}

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topology, updatedNodeSets, topoErr := BuildTopology(
		homeChainSelector,
		input.CapabilitiesAwareNodeSets,
		input.InfraInput,
		evmChainIDs,
		solChainIDs,
		bcOuts,
		allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		allChainsCLDEnvironment.DataStore,
		input.Capabilities,
		input.CapabilityConfigs,
		input.CopyCapabilityBinaries,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", stageGen.Elapsed().Seconds())))

	// start 3 tasks in the background
	backgroundStagesCount := 3
	backgroundStagesWaitGroup := &sync.WaitGroup{}
	backgroundStagesCh := make(chan backgroundStageResult, backgroundStagesCount)
	backgroundStagesWaitGroup.Add(1)

	// configure workflow registry contract in the background, so that we can continue with the next stage
	var workflowRegistryInput *cre.WorkflowRegistryInput
	var startTime time.Time
	go func() {
		defer backgroundStagesWaitGroup.Done()
		defer func() {
			if p := recover(); p != nil {
				backgroundStagesCh <- backgroundStageResult{
					panicValue: p,
					panicStack: debug.Stack(),
				}
			}
		}()
		startTime = time.Now()
		fmt.Print(libformat.PurpleText("---> [BACKGROUND 1/3] Configuring Workflow Registry contract\n"))

		allAddresses, addrErr := allChainsCLDEnvironment.ExistingAddresses.Addresses() //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
		if addrErr != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(addrErr, "failed to get addresses from address book")}
			return
		}

		chainsWithContracts := make(map[uint64]bool)
		for chainSelector, addresses := range allAddresses {
			chainsWithContracts[chainSelector] = len(addresses) > 0
		}

		addresses, addrErr1 := allChainsCLDEnvironment.DataStore.Addresses().Fetch()
		if addrErr1 != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(addrErr1, "failed to get addresses from datastore")}
			return
		}

		for _, addr := range addresses {
			chainsWithContracts[addr.ChainSelector] = true
		}

		nonEmptyBlockchains := make(map[uint64]cldf_chain.BlockChain, 0)
		for chainSelector, chain := range allChainsCLDEnvironment.BlockChains.EVMChains() {
			if chainsWithContracts[chain.Selector] {
				nonEmptyBlockchains[chainSelector] = chain
			}
		}
		for chainSelector, chain := range allChainsCLDEnvironment.BlockChains.SolanaChains() {
			if chainsWithContracts[chain.Selector] {
				nonEmptyBlockchains[chainSelector] = chain
			}
		}

		nonEmptyChainsCLDEnvironment := &cldf.Environment{
			Logger:            singleFileLogger,
			ExistingAddresses: allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
			GetContext: func() context.Context {
				return ctx
			},
			DataStore:   allChainsCLDEnvironment.DataStore,
			BlockChains: cldf_chain.NewBlockChains(nonEmptyBlockchains),
		}
		nonEmptyChainsCLDEnvironment.OperationsBundle = operations.NewBundle(nonEmptyChainsCLDEnvironment.GetContext, singleFileLogger, operations.NewMemoryReporter())

		// Configure Workflow Registry contract
		workflowRegistryInput = &cre.WorkflowRegistryInput{
			ContractAddress: common.HexToAddress(wfRegAddr),
			ChainSelector:   homeChainOutput.ChainSelector,
			// TODO, here we might need to pass new environment that doesn't have chains that do not have forwarders deployed
			CldEnv:         nonEmptyChainsCLDEnvironment,
			AllowedDonIDs:  []uint64{topology.WorkflowDONID},
			WorkflowOwners: []common.Address{homeChainOutput.SethClient.MustGetRootKeyAddress()},
		}

		_, workflowErr := libcontracts.ConfigureWorkflowRegistry(testLogger, workflowRegistryInput)
		if workflowErr != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(workflowErr, "failed to configure workflow registry"), successMessage: libformat.PurpleText("\n<--- [BACKGROUND 1/3] Workflow Registry configured in %.2f seconds\n", time.Since(startTime).Seconds())}
			return
		}

		backgroundStagesCh <- backgroundStageResult{successMessage: libformat.PurpleText("\n<--- [BACKGROUND 1/3] Workflow Registry configured in %.2f seconds\n", time.Since(startTime).Seconds())}
	}()

	// JOB DISTRIBUTOR + JOBS (creation and distribution to nodes)
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting Job Distributor, DONs and creating Jobs with Job Distributor")))

	jdOutput, nodeSetOutput, jobsSeqErr := SetupJobs(
		testLogger,
		input.JdInput,
		nixShell,
		homeChainOutput.BlockchainOutput,
		topology,
		input.InfraInput,
		updatedNodeSets,
	)
	if jobsSeqErr != nil {
		return nil, pkgerrors.Wrap(jobsSeqErr, "failed to setup jobs")
	}

	// append the nodeset output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	for idx, nsOut := range nodeSetOutput {
		input.CapabilitiesAwareNodeSets[idx].Out = nsOut.Output
	}

	for idx, bcOut := range blockchainOutputs {
		input.BlockchainsInput[idx].Out = bcOut.BlockchainOutput
	}

	// append the jd output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	input.JdInput.Out = jdOutput

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &cre.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutputs: bcOuts,
		SethClients:       sethClients,
		SolClients:        solClients,
		NodeSetOutput:     nodeSetOutput,
		ExistingAddresses: allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		Datastore:         allChainsCLDEnvironment.DataStore,
		Topology:          topology,
		OperationsBundle:  allChainsCLDEnvironment.OperationsBundle,
	}

	fullCldOutput, cldErr := libdevenv.BuildFullCLDEnvironment(ctx, singleFileLogger, fullCldInput, insecure.NewCredentials())
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to build full CLD environment")
	}

	createJobsInput := CreateJobsWithJdOpInput{}

	jobSpecFactoryFunctions := make([]cre.JobSpecFn, 0)
	for _, capability := range input.Capabilities {
		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, capability.JobSpecFn())
	}

	// Deprecated, use Capabilities instead
	jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, input.JobSpecFactoryFunctions...)

	createJobsDeps := CreateJobsWithJdOpDeps{
		Logger:                    testLogger,
		SingleFileLogger:          singleFileLogger,
		HomeChainBlockchainOutput: homeChainOutput.BlockchainOutput,
		AddressBook:               allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		JobSpecFactoryFunctions:   jobSpecFactoryFunctions,
		FullCLDEnvOutput:          fullCldOutput,
		CapabilitiesAwareNodeSets: input.CapabilitiesAwareNodeSets,
		InfraInput:                &input.InfraInput,
		CapabilitiesConfigs:       input.CapabilityConfigs,
		Capabilities:              input.Capabilities,
	}
	_, createJobsErr := operations.ExecuteOperation(allChainsCLDEnvironment.OperationsBundle, CreateJobsWithJdOp, createJobsDeps, createJobsInput)
	if createJobsErr != nil {
		return nil, pkgerrors.Wrap(createJobsErr, "failed to create jobs with Job Distributor")
	}

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Jobs created in %.2f seconds", stageGen.Elapsed().Seconds())))

	// Fund nodes in the background, so that we can continue with the next stage
	backgroundStagesWaitGroup.Add(1)
	go func() {
		defer backgroundStagesWaitGroup.Done()
		defer func() {
			if p := recover(); p != nil {
				backgroundStagesCh <- backgroundStageResult{
					panicValue: p,
					panicStack: debug.Stack(),
				}
			}
		}()

		startTime = time.Now()
		fmt.Print(libformat.PurpleText("---> [BACKGROUND 2/3] Funding Chainlink nodes\n\n"))

		_, fundErr := operations.ExecuteOperation(fullCldOutput.Environment.OperationsBundle, FundCLNodesOp, FundCLNodesOpDeps{
			Env:               fullCldOutput.Environment,
			BlockchainOutputs: blockchainOutputs,
			DonTopology:       fullCldOutput.DonTopology,
		}, FundCLNodesOpInput{FundAmount: 5000000000000000000})
		if fundErr != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(fundErr, "failed to fund CL nodes")}
			return
		}

		backgroundStagesCh <- backgroundStageResult{successMessage: libformat.PurpleText("\n<--- [BACKGROUND 2/3] Chainlink nodes funded in %.2f seconds\033[0m\n", time.Since(startTime).Seconds())}
	}()

	startTime = time.Now()
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Waiting for Log Poller to start tracking OCR3 contract")))

	for idx, nodeSetOut := range nodeSetOutput {
		if !flags.HasFlag(updatedNodeSets[idx].ComputedCapabilities, cre.ConsensusCapability) || !flags.HasFlag(updatedNodeSets[idx].ComputedCapabilities, cre.VaultCapability) {
			continue
		}
		nsClients, cErr := clclient.New(nodeSetOut.CLNodes)
		if cErr != nil {
			return nil, pkgerrors.Wrap(cErr, "failed to create node set clients")
		}
		eg := &errgroup.Group{}
		for _, c := range nsClients {
			eg.Go(func() error {
				return c.WaitHealthy(".*ConfigWatcher", "passing", 100)
			})
		}
		if waitErr := eg.Wait(); waitErr != nil {
			return nil, pkgerrors.Wrap(waitErr, "failed to wait for ConfigWatcher health check")
		}
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Log Poller started in %.2f seconds", stageGen.Elapsed().Seconds())))

	// wait for log poller filters to be registered in the background, because we don't need it them at this stage yet
	backgroundStagesWaitGroup.Add(1)
	go func() {
		defer backgroundStagesWaitGroup.Done()
		defer func() {
			if p := recover(); p != nil {
				backgroundStagesCh <- backgroundStageResult{
					panicValue: p,
					panicStack: debug.Stack(),
				}
			}
		}()

		if input.InfraInput.Type != infra.CRIB {
			hasGateway := false
			for _, don := range fullCldOutput.DonTopology.DonsWithMetadata {
				if flags.HasFlag(don.Flags, cre.GatewayDON) {
					hasGateway = true
					break
				}
			}

			if hasGateway {
				startTime = time.Now()
				fmt.Print(libformat.PurpleText("---> [BACKGROUND 3/3] Waiting for all nodes to have expected LogPoller filters registered\n\n"))

				testLogger.Info().Msg("Waiting for all nodes to have expected LogPoller filters registered...")
				lpErr := waitForAllNodesToHaveExpectedFiltersRegistered(singleFileLogger, testLogger, homeChainOutput.ChainID, *fullCldOutput.DonTopology, updatedNodeSets)
				if lpErr != nil {
					backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(lpErr, "failed to wait for all nodes to have expected LogPoller filters registered")}
					return
				}
				backgroundStagesCh <- backgroundStageResult{successMessage: libformat.PurpleText("\n<--- [BACKGROUND 3/3] Waiting for all nodes to have expected LogPoller filters registered finished in %.2f seconds\n\n", time.Since(startTime).Seconds())}
			}
		}
	}()

	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Configuring OCR3 and Keystone contracts")))

	// Configure the Forwarder, OCR3 and Capabilities contracts
	ocr3CommonAddr := common.HexToAddress(ocr3Addr)
	donTimeCommonAddr := common.HexToAddress(donTimeAddr)
	capRegCommonAddr := common.HexToAddress(capRegAddr)

	configureKeystoneInput := cre.ConfigureKeystoneInput{
		ChainSelector:               homeChainSelector,
		CldEnv:                      fullCldOutput.Environment,
		Topology:                    topology,
		CapabilitiesRegistryAddress: &capRegCommonAddr,
		OCR3Address:                 &ocr3CommonAddr,
		DONTimeAddress:              &donTimeCommonAddr,
		VaultOCR3Address:            &vaultOCR3CommonAddr,
		EVMOCR3Addresses:            &evmOCR3CommonAddresses,
		ConsensusV2OCR3Address:      &consensusV2OCR3CommonAddr,
		NodeSets:                    input.CapabilitiesAwareNodeSets,
	}

	if input.OCR3Config != nil {
		configureKeystoneInput.OCR3Config = *input.OCR3Config
	} else {
		ocr3Config, ocr3ConfigErr := libcontracts.DefaultOCR3Config(topology)
		if ocr3ConfigErr != nil {
			return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
		}
		configureKeystoneInput.OCR3Config = *ocr3Config
	}

	if input.DONTimeConfig != nil {
		configureKeystoneInput.DONTimeConfig = *input.DONTimeConfig
	} else {
		donTimeConfig, donTimeConfigErr := libcontracts.DefaultOCR3Config(topology)
		donTimeConfig.DeltaRoundMillis = 0 // Fastest rounds possible
		if donTimeConfigErr != nil {
			return nil, pkgerrors.Wrap(donTimeConfigErr, "failed to generate default DON Time config")
		}
		configureKeystoneInput.DONTimeConfig = *donTimeConfig
	}

	ocr3Config, ocr3ConfigErr := libcontracts.DefaultOCR3Config(topology)
	if ocr3ConfigErr != nil {
		return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
	}
	configureKeystoneInput.VaultOCR3Config = *ocr3Config

	defaultOcr3Config, defaultOcr3ConfigErr := libcontracts.DefaultOCR3Config(topology)
	if defaultOcr3ConfigErr != nil {
		return nil, pkgerrors.Wrap(defaultOcr3ConfigErr, "failed to generate default OCR3 config for EVM")
	}
	configureKeystoneInput.EVMOCR3Config = *defaultOcr3Config
	configureKeystoneInput.ConsensusV2OCR3Config = *defaultOcr3Config

	capabilitiesContractFactoryFunctions := make([]cre.CapabilityRegistryConfigFn, 0)
	for _, capability := range input.Capabilities {
		capabilitiesContractFactoryFunctions = append(capabilitiesContractFactoryFunctions, capability.CapabilityRegistryV1ConfigFn())
	}

	// Deprecated, use Capabilities instead
	capabilitiesContractFactoryFunctions = append(capabilitiesContractFactoryFunctions, input.CapabilitiesContractFactoryFunctions...)

	keystoneErr := libcontracts.ConfigureKeystone(configureKeystoneInput, capabilitiesContractFactoryFunctions)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("OCR3 and Keystone contracts configured in %.2f seconds", stageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Writing bootstrapping data into disk (address book, data store, etc...)")))

	artifactPath, artifactErr := DumpArtifact(
		memoryDatastore.AddressRefStore,
		allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		*jdOutput,
		*fullCldOutput.DonTopology,
		fullCldOutput.Environment.Offchain,
		capabilitiesContractFactoryFunctions,
		input.CapabilitiesAwareNodeSets,
	)
	if artifactErr != nil {
		testLogger.Error().Err(artifactErr).Msg("failed to generate artifact")
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Failed to write bootstrapping data into disk in %.2f seconds", stageGen.Elapsed().Seconds())))
	} else {
		testLogger.Info().Msgf("Environment artifact saved to %s", artifactPath)
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Wrote bootstrapping data into disk in %.2f seconds", stageGen.Elapsed().Seconds())))
	}

	// block on background stages
	backgroundStagesWaitGroup.Wait()
	close(backgroundStagesCh)

	for result := range backgroundStagesCh {
		if result.err != nil {
			return nil, pkgerrors.Wrap(result.err, "background stage failed")
		}
		if result.panicValue != nil {
			// Print the original stack trace from the background goroutine
			if result.panicStack != nil {
				fmt.Fprintf(os.Stderr, "Original panic stack trace from background goroutine:\n%s\n", result.panicStack)
			}
			panic(result.panicValue)
		}
		fmt.Print(result.successMessage)
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryInput.Out, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    blockchainOutputs,
		DonTopology:                         fullCldOutput.DonTopology,
		NodeOutput:                          nodeSetOutput,
		CldEnvironment:                      fullCldOutput.Environment,
		S3ProviderOutput:                    s3ProviderOutput,
	}, nil
}

func CreateJobDistributor(input *jd.Input) (*jd.Output, error) {
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(E2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(E2eJobDistributorVersionEnvVarName)
		input.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}

	jdOutput, err := jd.NewJD(input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create new job distributor")
	}

	return jdOutput, nil
}

func mergeJobSpecSlices(from, to cre.DonsToJobSpecs) {
	for fromDonID, fromJobSpecs := range from {
		if _, ok := to[fromDonID]; !ok {
			to[fromDonID] = make([]*jobv1.ProposeJobRequest, 0)
		}
		to[fromDonID] = append(to[fromDonID], fromJobSpecs...)
	}
}

type ConcurrentNonceMap struct {
	mu             sync.Mutex
	nonceByChainID map[uint64]uint64
}

func NewConcurrentNonceMap(ctx context.Context, blockchainOutputs []*cre.WrappedBlockchainOutput) (*ConcurrentNonceMap, error) {
	nonceByChainID := make(map[uint64]uint64)
	for _, bcOut := range blockchainOutputs {
		if bcOut.BlockchainOutput.Family == chainselectors.FamilyEVM {
			var err error
			ctxWithTimeout, cancel := context.WithTimeout(ctx, bcOut.SethClient.Cfg.Network.TxnTimeout.Duration())
			nonceByChainID[bcOut.ChainID], err = bcOut.SethClient.Client.PendingNonceAt(ctxWithTimeout, bcOut.SethClient.MustGetRootKeyAddress())
			cancel()
			if err != nil {
				cancel()
				return nil, pkgerrors.Wrapf(err, "failed to get nonce for chain %d", bcOut.ChainID)
			}
		}
	}
	return &ConcurrentNonceMap{nonceByChainID: nonceByChainID}, nil
}

func (c *ConcurrentNonceMap) Decrement(chainID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nonceByChainID[chainID]--
}

func (c *ConcurrentNonceMap) Increment(chainID uint64) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nonceByChainID[chainID]++
	return c.nonceByChainID[chainID]
}

// must match nubmer of events we track in core/services/workflows/syncer/handler.go
const NumberOfTrackedWorkflowRegistryEvents = 6

// waitForAllNodesToHaveExpectedFiltersRegistered manually checks if all WorkflowRegistry filters used by the LogPoller are registered for all nodes. We want to see if this will help with the flakiness.
func waitForAllNodesToHaveExpectedFiltersRegistered(singeFileLogger logger.Logger, testLogger zerolog.Logger, homeChainID uint64, donTopology cre.DonTopology, nodeSetInput []*cre.CapabilitiesAwareNodeSet) error {
	for donIdx, don := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(don.Flags, cre.WorkflowDON) {
			continue
		}

		workderNodes, workersErr := crenode.FindManyWithLabel(don.NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode}, crenode.EqualLabels)
		if workersErr != nil {
			return pkgerrors.Wrap(workersErr, "failed to find worker nodes")
		}

		results := make(map[int]bool)
		ticker := 5 * time.Second
		timeout := 2 * time.Minute

	INNER_LOOP:
		for {
			select {
			case <-time.After(timeout):
				return fmt.Errorf("timed out, when waiting for %.2f seconds, waiting for all nodes to have expected filters registered", timeout.Seconds())
			case <-time.Tick(ticker):
				if len(results) == len(workderNodes) {
					testLogger.Info().Msgf("All %d nodes in DON %d have expected filters registered", len(workderNodes), don.ID)
					break INNER_LOOP
				}

				for _, workerNode := range workderNodes {
					nodeIndex, nodeIndexErr := crenode.FindLabelValue(workerNode, crenode.IndexKey)
					if nodeIndexErr != nil {
						return pkgerrors.Wrap(nodeIndexErr, "failed to find node index")
					}

					nodeIndexInt, nodeIdxErr := strconv.Atoi(nodeIndex)
					if nodeIdxErr != nil {
						return pkgerrors.Wrap(nodeIdxErr, "failed to convert node index to int")
					}

					if _, ok := results[nodeIndexInt]; ok {
						continue
					}

					testLogger.Info().Msgf("Checking if all WorkflowRegistry filters are registered for worker node %d", nodeIndexInt)
					allFilters, filtersErr := getAllFilters(context.Background(), singeFileLogger, big.NewInt(libc.MustSafeInt64(homeChainID)), nodeIndexInt, nodeSetInput[donIdx].DbInput.Port)
					if filtersErr != nil {
						return pkgerrors.Wrap(filtersErr, "failed to get filters")
					}

					for _, filter := range allFilters {
						if strings.Contains(filter.Name, "WorkflowRegistry") {
							if len(filter.EventSigs) == NumberOfTrackedWorkflowRegistryEvents {
								testLogger.Debug().Msgf("Found all WorkflowRegistry filters for node %d", nodeIndexInt)
								results[nodeIndexInt] = true
								continue
							}

							testLogger.Debug().Msgf("Found only %d WorkflowRegistry filters for node %d", len(filter.EventSigs), nodeIndexInt)
						}
					}
				}

				// return if we have results for all nodes, don't wait for next tick
				if len(results) == len(workderNodes) {
					testLogger.Info().Msgf("All %d nodes in DON %d have expected filters registered", len(workderNodes), don.ID)
					break INNER_LOOP
				}
			}
		}
	}

	return nil
}

func NewORM(logger logger.Logger, chainID *big.Int, nodeIndex, externalPort int) (logpoller.ORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", externalPort, postgres.User, postgres.Password, fmt.Sprintf("db_%d", nodeIndex))
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return logpoller.NewORM(chainID, db, logger), db, nil
}

func getAllFilters(ctx context.Context, logger logger.Logger, chainID *big.Int, nodeIndex, externalPort int) (map[string]logpoller.Filter, error) {
	orm, db, err := NewORM(logger, chainID, nodeIndex, externalPort)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	return orm.LoadFilters(ctx)
}
