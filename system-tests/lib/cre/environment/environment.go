package environment

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/scylladb/go-reflectx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcaps "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdevenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	creconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	cresecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	libinfra "github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

const (
	cronCapabilityAssetFile            = "cron"
	GithubReadTokenEnvVarName          = "GITHUB_READ_TOKEN"
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
	cribConfigsDir                     = "crib-configs"
)

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cretypes.WorkflowRegistryOutput
	CldEnvironment                      *cldf.Environment
	BlockchainOutput                    []*BlockchainOutput
	DonTopology                         *cretypes.DonTopology
	NodeOutput                          []*cretypes.WrappedNodeOutput
}

type SetupInput struct {
	CapabilitiesAwareNodeSets            []*cretypes.CapabilitiesAwareNodeSet
	CapabilitiesContractFactoryFunctions []func([]cretypes.CapabilityFlag) []keystone_changeset.DONCapabilityWithConfig
	ConfigFactoryFunctions               []cretypes.ConfigFactoryFn
	JobSpecFactoryFunctions              []cretypes.JobSpecFactoryFn
	BlockchainsInput                     []*blockchain.Input
	JdInput                              jd.Input
	InfraInput                           libtypes.InfraInput
	CustomBinariesPaths                  map[cretypes.CapabilityFlag]string
	OCR3Config                           *keystone_changeset.OracleConfig
}

type backgroundStageResult struct {
	err            error
	successMessage string
}

func SetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singeFileLogger logger.Logger,
	input SetupInput,
) (*SetupOutput, error) {
	topologyErr := libdon.ValidateTopology(input.CapabilitiesAwareNodeSets, input.InfraInput)
	if topologyErr != nil {
		return nil, pkgerrors.Wrap(topologyErr, "failed to validate topology")
	}

	// Shell is only required, when using CRIB, because we want to run commands in the same "nix develop" context
	// We need to have this reference in the outer scope, because subsequent functions will need it
	var nixShell *libnix.Shell
	if input.InfraInput.InfraType == libtypes.CRIB {
		startNixShellInput := &cretypes.StartNixShellInput{
			InfraInput:     &input.InfraInput,
			CribConfigsDir: cribConfigsDir,
			PurgeNamespace: true,
		}

		var nixErr error
		nixShell, nixErr = crib.StartNixShell(startNixShellInput)
		if nixErr != nil {
			return nil, pkgerrors.Wrap(nixErr, "failed to start nix shell")
		}
	}

	defer func() {
		if nixShell != nil {
			_ = nixShell.Close()
		}
	}()

	bi := BlockchainsInput{
		infra:    &input.InfraInput,
		nixShell: nixShell,
	}
	bi.blockchainsInput = append(bi.blockchainsInput, input.BlockchainsInput...)

	startTime := time.Now()
	fmt.Print(libformat.PurpleText("\n[Stage 1/8] Starting %d blockchain(s)\n\n", len(bi.blockchainsInput)))

	blockchainsOutput, bcOutErr := CreateBlockchains(testLogger, bi)
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to create blockchains")
	}

	homeChainOutput := blockchainsOutput[0]
	chainsConfigs := []devenv.ChainConfig{}

	for _, bcOut := range blockchainsOutput {
		chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
			ChainID:   strconv.FormatUint(bcOut.SethClient.Cfg.Network.ChainID, 10),
			ChainName: bcOut.SethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(bcOut.BlockchainOutput.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: bcOut.BlockchainOutput.Nodes[0].ExternalWSUrl,
				Internal: bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
				Internal: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			}},
			DeployerKey: bcOut.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
		})
	}

	allChains, _, allChainsErr := devenv.NewChains(singeFileLogger, chainsConfigs)
	if allChainsErr != nil {
		return nil, pkgerrors.Wrap(allChainsErr, "failed to create chains")
	}

	blockChains := map[uint64]chain.BlockChain{}
	for selector, ch := range allChains {
		blockChains[selector] = ch
	}

	allChainsCLDEnvironment := &cldf.Environment{
		Logger:            singeFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: chain.NewBlockChains(blockChains),
	}

	fmt.Print(libformat.PurpleText("\n[Stage 1/8] Blockchains started in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 2/8 Deploying Keystone contracts\n\n"))

	// we could try to parallelise deployment of these contracts, but it's difficult, because there's no way to make chain.DeployerKey dynamic
	// in order to manually increment the nonce for each contract
	ocr3Output, ocr3Err := keystone_changeset.DeployOCR3V2(*allChainsCLDEnvironment, &keystone_changeset.DeployRequestV2{
		ChainSel: homeChainOutput.ChainSelector,
	})
	if ocr3Err != nil {
		return nil, pkgerrors.Wrap(ocr3Err, "failed to deploy OCR3 contract")
	}

	mergeErr := allChainsCLDEnvironment.ExistingAddresses.Merge(ocr3Output.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		return nil, pkgerrors.Wrap(mergeErr, "failed to merge address book")
	}
	testLogger.Info().Msgf("Deployed OCR3 contract on chain %d at %s", homeChainOutput.ChainSelector, libcontracts.MustFindAddressesForChain(allChainsCLDEnvironment.ExistingAddresses, homeChainOutput.ChainSelector, keystone_changeset.OCR3Capability.String())) //nolint:staticcheck // won't migrate now

	capabilitiesRegistryOutput, capabilitiesRegistryErr := keystone_changeset.DeployCapabilityRegistry(*allChainsCLDEnvironment, homeChainOutput.ChainSelector)
	if capabilitiesRegistryErr != nil {
		return nil, pkgerrors.Wrap(capabilitiesRegistryErr, "failed to deploy Capabilities Registry contract")
	}

	mergeErr = allChainsCLDEnvironment.ExistingAddresses.Merge(capabilitiesRegistryOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		return nil, pkgerrors.Wrap(mergeErr, "failed to merge address book")
	}
	testLogger.Info().Msgf("Deployed Capabilities Registry contract on chain %d at %s", homeChainOutput.ChainSelector, libcontracts.MustFindAddressesForChain(allChainsCLDEnvironment.ExistingAddresses, homeChainOutput.ChainSelector, keystone_changeset.CapabilitiesRegistry.String())) //nolint:staticcheck // won't migrate now

	workflowRegistryOutput, workflowRegistryErr := workflow_registry_changeset.Deploy(*allChainsCLDEnvironment, homeChainOutput.ChainSelector)
	if workflowRegistryErr != nil {
		return nil, pkgerrors.Wrap(workflowRegistryErr, "failed to deploy workflow registry contract")
	}

	mergeErr = allChainsCLDEnvironment.ExistingAddresses.Merge(workflowRegistryOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		return nil, pkgerrors.Wrap(mergeErr, "failed to merge address book")
	}
	testLogger.Info().Msgf("Deployed Workflow Registry contract on chain %d at %s", homeChainOutput.ChainSelector, libcontracts.MustFindAddressesForChain(allChainsCLDEnvironment.ExistingAddresses, homeChainOutput.ChainSelector, keystone_changeset.WorkflowRegistry.String())) //nolint:staticcheck // won't migrate now

	// Deploy forwarders for all chains
	contractErrGroup := &errgroup.Group{}
	for _, bcOut := range blockchainsOutput {
		contractErrGroup.Go(func() error {
			output, err := keystone_changeset.DeployForwarder(*allChainsCLDEnvironment, keystone_changeset.DeployForwarderRequest{
				ChainSelectors: []uint64{bcOut.ChainSelector},
			})
			if err != nil {
				return pkgerrors.Wrap(err, "failed to deploy forwarder contract")
			}

			mergeErr := allChainsCLDEnvironment.ExistingAddresses.Merge(output.AddressBook) //nolint:staticcheck // won't migrate now
			if mergeErr != nil {
				return pkgerrors.Wrap(mergeErr, "failed to merge address book")
			}
			testLogger.Info().Msgf("Deployed Forwarder contract on chain %d at %s", bcOut.ChainSelector, libcontracts.MustFindAddressesForChain(allChainsCLDEnvironment.ExistingAddresses, bcOut.ChainSelector, keystone_changeset.KeystoneForwarder.String())) //nolint:staticcheck // won't migrate now

			return nil
		})
	}

	if err := contractErrGroup.Wait(); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to deploy Keystone contracts")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 2/8] Contracts deployed in %.2f seconds\n", time.Since(startTime).Seconds()))

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topology, topoErr := libdon.BuildTopology(input.CapabilitiesAwareNodeSets, input.InfraInput, homeChainOutput.ChainSelector)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}

	// start 3 tasks in the background
	backgroundStagesCount := 3
	backgroundStagesWaitGroup := &sync.WaitGroup{}
	backgroundStagesCh := make(chan backgroundStageResult, backgroundStagesCount)
	backgroundStagesWaitGroup.Add(1)

	// configure workflow registry contract in the background, so that we can continue with the next stage
	var workflowRegistryInput *cretypes.WorkflowRegistryInput
	go func() {
		defer backgroundStagesWaitGroup.Done()
		startTime = time.Now()
		fmt.Print(libformat.PurpleText("---> [BACKGROUND 1/3] Configuring Workflow Registry contract\n"))

		// Configure Workflow Registry contract
		workflowRegistryInput = &cretypes.WorkflowRegistryInput{
			ChainSelector:  homeChainOutput.ChainSelector,
			CldEnv:         allChainsCLDEnvironment,
			AllowedDonIDs:  []uint32{topology.WorkflowDONID},
			WorkflowOwners: []common.Address{homeChainOutput.SethClient.MustGetRootKeyAddress()},
		}

		_, workflowErr := libcontracts.ConfigureWorkflowRegistry(testLogger, workflowRegistryInput)
		if workflowErr != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(workflowErr, "failed to configure workflow registry"), successMessage: libformat.PurpleText("\n<--- [BACKGROUND 1/3] Workflow Registry configured in %.2f seconds\n", time.Since(startTime).Seconds())}
			return
		}

		backgroundStagesCh <- backgroundStageResult{successMessage: libformat.PurpleText("\n<--- [BACKGROUND 1/3] Workflow Registry configured in %.2f seconds\n", time.Since(startTime).Seconds())}
	}()

	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 3/8] Preparing DON(s) configuration\n\n"))

	// Generate EVM and P2P keys or read them from the config
	// That way we can pass them final configs and do away with restarting the nodes
	var keys *cretypes.GenerateKeysOutput

	keysOutput, keysOutputErr := cresecrets.KeysOutputFromConfig(input.CapabilitiesAwareNodeSets)
	if keysOutputErr != nil {
		return nil, pkgerrors.Wrap(keysOutputErr, "failed to generate keys output")
	}

	// get chainIDs, they'll be used for identifying ETH keys and Forwarder addresses
	// and also for creating the CLD environment
	chainIDs := make([]int, 0)
	bcOuts := make(map[uint64]*blockchain.Output)
	sethClients := make(map[uint64]*seth.Client)
	for _, bcOut := range blockchainsOutput {
		chainIDs = append(chainIDs, libc.MustSafeInt(bcOut.ChainID))
		bcOuts[bcOut.ChainSelector] = bcOut.BlockchainOutput
		sethClients[bcOut.ChainSelector] = bcOut.SethClient
	}

	generateKeysInput := &cretypes.GenerateKeysInput{
		GenerateEVMKeysForChainIDs: chainIDs,
		GenerateP2PKeys:            true,
		Topology:                   topology,
		Password:                   "", // since the test runs on private ephemeral blockchain we don't use real keys and do not care a lot about the password
		Out:                        keysOutput,
	}
	keys, keysErr := cresecrets.GenereteKeys(generateKeysInput)
	if keysErr != nil {
		return nil, pkgerrors.Wrap(keysErr, "failed to generate keys")
	}

	topology, addKeysErr := cresecrets.AddKeysToTopology(topology, keys)
	if addKeysErr != nil {
		return nil, pkgerrors.Wrap(addKeysErr, "failed to add keys to topology")
	}

	peeringData, peeringErr := libdon.FindPeeringData(topology)
	if peeringErr != nil {
		return nil, pkgerrors.Wrap(peeringErr, "failed to find peering data")
	}

	for i, donMetadata := range topology.DonsMetadata {
		configsFound := 0
		secretsFound := 0
		for _, nodeSpec := range input.CapabilitiesAwareNodeSets[i].NodeSpecs {
			if nodeSpec.Node.TestConfigOverrides != "" {
				configsFound++
			}
			if nodeSpec.Node.TestSecretsOverrides != "" {
				secretsFound++
			}
		}
		if configsFound != 0 && configsFound != len(input.CapabilitiesAwareNodeSets[i].NodeSpecs) {
			return nil, fmt.Errorf("%d out of %d node specs have config overrides. Either provide overrides for all nodes or none at all", configsFound, len(input.CapabilitiesAwareNodeSets[i].NodeSpecs))
		}

		if secretsFound != 0 && secretsFound != len(input.CapabilitiesAwareNodeSets[i].NodeSpecs) {
			return nil, fmt.Errorf("%d out of %d node specs have secrets overrides. Either provide overrides for all nodes or none at all", secretsFound, len(input.CapabilitiesAwareNodeSets[i].NodeSpecs))
		}

		// Allow providing only secrets, because we can decode them and use them to generate configs
		// We can't allow providing only configs, because we can't replace secret-related values in the configs
		// If both are provided, we assume that the user knows what they are doing and we don't need to validate anything
		// And that configs match the secrets
		if configsFound > 0 && secretsFound == 0 {
			return nil, fmt.Errorf("nodese config overrides are provided for DON %d, but not secrets. You need to either provide both, only secrets or nothing at all", donMetadata.ID)
		}

		// generate configs only if they are not provided
		if configsFound == 0 {
			config, configErr := creconfig.Generate(
				cretypes.GenerateConfigsInput{
					DonMetadata:            donMetadata,
					BlockchainOutput:       bcOuts,
					Flags:                  donMetadata.Flags,
					PeeringData:            peeringData,
					AddressBook:            allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
					HomeChainSelector:      topology.HomeChainSelector,
					GatewayConnectorOutput: topology.GatewayConnectorOutput,
				},
				input.ConfigFactoryFunctions,
			)
			if configErr != nil {
				return nil, pkgerrors.Wrap(configErr, "failed to generate config")
			}

			for j := range donMetadata.NodesMetadata {
				input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
			}
		}

		// generate secrets only if they are not provided
		if secretsFound == 0 {
			secretsInput := &cretypes.GenerateSecretsInput{
				DonMetadata: donMetadata,
			}

			if evmKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
				secretsInput.EVMKeys = evmKeys
			}

			if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
				secretsInput.P2PKeys = p2pKeys
			}

			// EVM and P2P keys will be provided to nodes as secrets
			secrets, secretsErr := cresecrets.GenerateSecrets(
				secretsInput,
			)
			if secretsErr != nil {
				return nil, pkgerrors.Wrap(secretsErr, "failed to generate secrets")
			}

			for j := range donMetadata.NodesMetadata {
				input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.TestSecretsOverrides = secrets[j]
			}
		}

		var appendErr error
		input.CapabilitiesAwareNodeSets[i], appendErr = libcaps.AppendBinariesPathsNodeSpec(input.CapabilitiesAwareNodeSets[i], donMetadata, input.CustomBinariesPaths)
		if appendErr != nil {
			return nil, pkgerrors.Wrapf(appendErr, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
		}
	}

	// Deploy the DONs
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range input.CapabilitiesAwareNodeSets {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range input.CapabilitiesAwareNodeSets[i].NodeSpecs {
				input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.Image = image
			}
		}
	}

	fmt.Print(libformat.PurpleText("\n[Stage 3/8] DONs configuration prepared in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 4/8] Starting Job Distributor\n"))

	if input.InfraInput.InfraType == libtypes.CRIB {
		deployCribJdInput := &cretypes.DeployCribJdInput{
			JDInput:        &input.JdInput,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var jdErr error
		input.JdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
		}
	}

	jdAndDonsErrGroup := &errgroup.Group{}
	var jdOutput *jd.Output

	jdAndDonsErrGroup.Go(func() error {
		var jdErr error
		jdOutput, jdErr = CreateJobDistributor(&input.JdInput)
		if jdErr != nil {
			jdErr = fmt.Errorf("failed to start JD container for image %s: %w", input.JdInput.Image, jdErr)

			// useful end user messages
			if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
				jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
			}
			return jdErr
		}

		fmt.Print(libformat.PurpleText("\n[Stage 4/8] Job Distributor started in %.2f seconds\n", time.Since(startTime).Seconds()))

		return nil
	})

	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 5/8] Starting %d DON(s)\n\n", len(input.CapabilitiesAwareNodeSets)))

	if input.InfraInput.InfraType == libtypes.CRIB {
		testLogger.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &cretypes.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  input.CapabilitiesAwareNodeSets,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var devspaceErr error
		input.CapabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with devspace")
		}
	}

	nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))

	jdAndDonsErrGroup.Go(func() error {
		// TODO we could parallelise this as well in the future, but for single DON env this doesn't matter
		for _, nodeSetInput := range input.CapabilitiesAwareNodeSets {
			nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, homeChainOutput.BlockchainOutput)
			if nodesetErr != nil {
				return pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
			}

			nodeSetOutput = append(nodeSetOutput, &cretypes.WrappedNodeOutput{
				Output:       nodeset,
				NodeSetName:  nodeSetInput.Name,
				Capabilities: nodeSetInput.Capabilities,
			})
		}

		return nil
	})

	if jdAndDonErr := jdAndDonsErrGroup.Wait(); jdAndDonErr != nil {
		return nil, pkgerrors.Wrap(jdAndDonErr, "failed to start Job Distributor or DONs")
	}

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &cretypes.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutputs: bcOuts,
		SethClients:       sethClients,
		NodeSetOutput:     nodeSetOutput,
		ExistingAddresses: allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		Topology:          topology,
	}

	// We need to use TLS for CRIB, because it exposes HTTPS endpoints
	var creds credentials.TransportCredentials
	if input.InfraInput.InfraType == libtypes.CRIB {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	fullCldOutput, cldErr := libdevenv.BuildFullCLDEnvironment(ctx, singeFileLogger, fullCldInput, creds)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to build full CLD environment")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 5/8] DONs started in %.2f seconds\n", time.Since(startTime).Seconds()))

	// Fund nodes in the background, so that we can continue with the next stage
	backgroundStagesWaitGroup.Add(1)
	go func() {
		defer backgroundStagesWaitGroup.Done()

		startTime = time.Now()
		fmt.Print(libformat.PurpleText("---> [BACKGROUND 2/3] Funding Chainlink nodes\n\n"))

		// Fund the nodes
		concurrentNonceMap, concurrentNonceMapErr := NewConcurrentNonceMap(ctx, blockchainsOutput)
		if concurrentNonceMapErr != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(concurrentNonceMapErr, "failed to create concurrent nonce map")}
			return
		}

		// Decrement the nonce for each chain, because we will increment it in the next loop
		for _, bcOut := range blockchainsOutput {
			concurrentNonceMap.Decrement(bcOut.ChainID)
		}

		errGroup := &errgroup.Group{}
		for _, metaDon := range fullCldOutput.DonTopology.DonsWithMetadata {
			for _, bcOut := range blockchainsOutput {
				for _, node := range metaDon.DON.Nodes {
					errGroup.Go(func() error {
						nodeAddress := node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]
						if nodeAddress == "" {
							return nil
						}

						nonce := concurrentNonceMap.Increment(bcOut.ChainID)

						_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bcOut.SethClient, libtypes.FundsToSend{
							ToAddress:  common.HexToAddress(nodeAddress),
							Amount:     big.NewInt(5000000000000000000),
							PrivateKey: bcOut.SethClient.MustGetRootPrivateKey(),
							Nonce:      ptr.Ptr(nonce),
						})
						if fundingErr != nil {
							return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", nodeAddress)
						}
						return nil
					})
				}
			}
		}

		if err := errGroup.Wait(); err != nil {
			backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(err, "failed to fund nodes")}
			return
		}

		backgroundStagesCh <- backgroundStageResult{successMessage: libformat.PurpleText("\n<--- [BACKGROUND 2/3] Chainlink nodes funded in %.2f seconds\033[0m\n", time.Since(startTime).Seconds())}
	}()

	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 6/8] Creating jobs with Job Distributor\n\n"))

	donToJobSpecs := make(cretypes.DonsToJobSpecs)

	for _, jobSpecGeneratingFn := range input.JobSpecFactoryFunctions {
		singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cretypes.JobSpecFactoryInput{
			CldEnvironment:   fullCldOutput.Environment,
			BlockchainOutput: homeChainOutput.BlockchainOutput,
			DonTopology:      fullCldOutput.DonTopology,
			AddressBook:      allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		})
		if jobSpecsErr != nil {
			return nil, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
		}
		mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
	}

	createJobsInput := cretypes.CreateJobsInput{
		CldEnv:        fullCldOutput.Environment,
		DonTopology:   fullCldOutput.DonTopology,
		DonToJobSpecs: donToJobSpecs,
	}

	jobsErr := libdon.CreateJobs(ctx, testLogger, createJobsInput)
	if jobsErr != nil {
		return nil, pkgerrors.Wrap(jobsErr, "failed to create jobs")
	}

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	fmt.Print(libformat.PurpleText("\n[Stage 6/8] Jobs created in %.2f seconds\033[0m\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 7/8] Waiting for Log Poller to start tracking OCR3 contract\n\n"))

	for idx, nodeSetOut := range nodeSetOutput {
		if !flags.HasFlag(input.CapabilitiesAwareNodeSets[idx].Capabilities, cretypes.OCR3Capability) {
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
		if err := eg.Wait(); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to wait for ConfigWatcher health check")
		}
	}

	fmt.Print(libformat.PurpleText("\n[Stage 7/8] Log Poller started in %.2f seconds\n", time.Since(startTime).Seconds()))

	// wait for log poller filters to be registered in the background, because we don't need it them at this stage yet
	backgroundStagesWaitGroup.Add(1)
	go func() {
		defer backgroundStagesWaitGroup.Done()

		if input.InfraInput.InfraType != libtypes.CRIB {
			hasGateway := false
			for _, don := range fullCldOutput.DonTopology.DonsWithMetadata {
				if flags.HasFlag(don.Flags, cretypes.GatewayDON) {
					hasGateway = true
					break
				}
			}

			if hasGateway {
				startTime = time.Now()
				fmt.Print(libformat.PurpleText("---> [BACKGROUND 3/3] Waiting for all nodes to have expected LogPoller filters registered\n\n"))

				testLogger.Info().Msg("Waiting for all nodes to have expected LogPoller filters registered...")
				lpErr := waitForAllNodesToHaveExpectedFiltersRegistered(singeFileLogger, testLogger, homeChainOutput.ChainID, *fullCldOutput.DonTopology, input.CapabilitiesAwareNodeSets)
				if lpErr != nil {
					backgroundStagesCh <- backgroundStageResult{err: pkgerrors.Wrap(lpErr, "failed to wait for all nodes to have expected LogPoller filters registered")}
					return
				}
				backgroundStagesCh <- backgroundStageResult{successMessage: libformat.PurpleText("\n<--- [BACKGROUND 3/3] Waiting for all nodes to have expected LogPoller filters registered finished in %.2f seconds\n\n", time.Since(startTime).Seconds())}
			}
		}
	}()

	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 8/8] Configuring OCR3 and Keystone contracts\n\n"))

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := cretypes.ConfigureKeystoneInput{
		ChainSelector: homeChainOutput.ChainSelector,
		CldEnv:        fullCldOutput.Environment,
		Topology:      topology,
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

	keystoneErr := libcontracts.ConfigureKeystone(configureKeystoneInput, input.CapabilitiesContractFactoryFunctions)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 8/8] OCR3 and Keystone contracts configured in %.2f seconds\n", time.Since(startTime).Seconds()))

	// block on background stages
	backgroundStagesWaitGroup.Wait()
	close(backgroundStagesCh)

	for result := range backgroundStagesCh {
		if result.err != nil {
			return nil, pkgerrors.Wrap(result.err, "background stage failed")
		}
		fmt.Print(result.successMessage)
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryInput.Out, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    blockchainsOutput,
		DonTopology:                         fullCldOutput.DonTopology,
		NodeOutput:                          nodeSetOutput,
		CldEnvironment:                      fullCldOutput.Environment,
	}, nil
}

type BlockchainsInput struct {
	blockchainsInput []*blockchain.Input
	infra            *libtypes.InfraInput
	nixShell         *libnix.Shell
}

type BlockchainOutput struct {
	ChainSelector      uint64
	ChainID            uint64
	BlockchainOutput   *blockchain.Output
	SethClient         *seth.Client
	DeployerPrivateKey string

	// private data depending crib vs docker
	c *blockchain.Output // non-nil if running in docker
}

func CreateBlockchains(
	testLogger zerolog.Logger,
	input BlockchainsInput,
) ([]*BlockchainOutput, error) {
	if len(input.blockchainsInput) == 0 {
		return nil, pkgerrors.New("blockchain input is nil")
	}

	blockchainOutput := make([]*BlockchainOutput, 0)
	for _, bi := range input.blockchainsInput {
		var bcOut *blockchain.Output
		var bcErr error
		if input.infra.InfraType == libtypes.CRIB {
			if input.nixShell == nil {
				return nil, pkgerrors.New("nix shell is nil")
			}

			deployCribBlockchainInput := &cretypes.DeployCribBlockchainInput{
				BlockchainInput: bi,
				NixShell:        input.nixShell,
				CribConfigsDir:  cribConfigsDir,
			}
			bcOut, bcErr = crib.DeployBlockchain(deployCribBlockchainInput)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
			err := libinfra.WaitForRPCEndpoint(testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
			}
		} else {
			bcOut, bcErr = blockchain.NewBlockchainNetwork(bi)
			if bcErr != nil {
				return nil, pkgerrors.Wrap(bcErr, "failed to deploy blockchain")
			}
		}

		pkey := os.Getenv("PRIVATE_KEY")
		if pkey == "" {
			return nil, pkgerrors.New("PRIVATE_KEY env var must be set")
		}

		sethClient, err := seth.NewClientBuilder().
			WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
			WithPrivateKeys([]string{pkey}).
			// do not check if there's a pending nonce nor check node's health
			WithProtections(false, false, seth.MustMakeDuration(time.Second)).
			Build()
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to create seth client")
		}

		chainSelector, err := chainselectors.SelectorFromChainId(sethClient.Cfg.Network.ChainID)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %d", sethClient.Cfg.Network.ChainID)
		}
		chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
		}

		blockchainOutput = append(blockchainOutput, &BlockchainOutput{
			ChainSelector:      chainSelector,
			ChainID:            chainID,
			BlockchainOutput:   bcOut,
			SethClient:         sethClient,
			DeployerPrivateKey: pkey,
			c:                  bcOut,
		})
	}

	return blockchainOutput, nil
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

func mergeJobSpecSlices(from, to cretypes.DonsToJobSpecs) {
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

func NewConcurrentNonceMap(ctx context.Context, blockchainOutputs []*BlockchainOutput) (*ConcurrentNonceMap, error) {
	nonceByChainID := make(map[uint64]uint64)
	for _, bcOut := range blockchainOutputs {
		var err error
		ctxWithTimeout, cancel := context.WithTimeout(ctx, bcOut.SethClient.Cfg.Network.TxnTimeout.Duration())
		nonceByChainID[bcOut.ChainID], err = bcOut.SethClient.Client.PendingNonceAt(ctxWithTimeout, bcOut.SethClient.MustGetRootKeyAddress())
		cancel()
		if err != nil {
			cancel()
			return nil, pkgerrors.Wrapf(err, "failed to get nonce for chain %d", bcOut.ChainID)
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
func waitForAllNodesToHaveExpectedFiltersRegistered(singeFileLogger logger.Logger, testLogger zerolog.Logger, homeChainID uint64, donTopology cretypes.DonTopology, nodeSetInput []*cretypes.CapabilitiesAwareNodeSet) error {
	for donIdx, don := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(don.Flags, cretypes.WorkflowDON) {
			continue
		}

		workderNodes, workersErr := crenode.FindManyWithLabel(don.NodesMetadata, &cretypes.Label{Key: crenode.NodeTypeKey, Value: cretypes.WorkerNode}, crenode.EqualLabels)
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
