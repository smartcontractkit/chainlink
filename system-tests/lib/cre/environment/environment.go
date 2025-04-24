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
	"time"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	libcaps "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdevenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	keystoneporconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/por"
	cresecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	keystonesecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
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
	KeystoneContractsOutput             *keystonetypes.KeystoneContractsOutput
	WorkflowRegistryConfigurationOutput *keystonetypes.WorkflowRegistryOutput
	CldEnvironment                      *deployment.Environment
	BlockchainOutput                    []*BlockchainOutput
	DonTopology                         *keystonetypes.DonTopology
	NodeOutput                          []*keystonetypes.WrappedNodeOutput
}

type SetupInput struct {
	ExtraAllowedPorts                    []int
	CapabilitiesAwareNodeSets            []*keystonetypes.CapabilitiesAwareNodeSet
	CapabilitiesContractFactoryFunctions []func([]cretypes.CapabilityFlag) []keystone_changeset.DONCapabilityWithConfig
	JobSpecFactoryFunctions              []cretypes.JobSpecFactoryFn
	BlockchainsInput                     []*blockchain.Input
	JdInput                              jd.Input
	InfraInput                           libtypes.InfraInput
	CustomBinariesPaths                  map[cretypes.CapabilityFlag]string
}

func SetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singeFileLogger *cldlogger.SingleFileLogger,
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
		startNixShellInput := &keystonetypes.StartNixShellInput{
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

	bi := make([]BlockchainsInput, 0)
	for _, bcInfraInput := range input.BlockchainsInput {
		bi = append(bi, BlockchainsInput{
			blockchainInput: bcInfraInput,
			infraInput:      &input.InfraInput,
			nixShell:        nixShell,
		})
	}
	homeChainInput := bi[0]

	blockchainsOutput, bcOutErr := CreateBlockchains(testLogger, bi)
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to create blockchains")
	}

	homeChainOutput := blockchainsOutput[0]

	// Home chain, where we deploy the Keystone

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	// but first, we need to create deployment.Environment that will contain only chain information in order to deploy contracts with the CLD
	homeChainConfigs := []devenv.ChainConfig{
		{
			ChainID:   homeChainOutput.SethClient.Cfg.Network.ChainID,
			ChainName: homeChainOutput.SethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(homeChainOutput.BlockchainOutput.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: homeChainOutput.BlockchainOutput.Nodes[0].ExternalWSUrl,
				Internal: homeChainOutput.BlockchainOutput.Nodes[0].InternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: homeChainOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
				Internal: homeChainOutput.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			}},
			DeployerKey: homeChainOutput.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
		},
	}

	chains, chainsErr := devenv.NewChains(singeFileLogger, homeChainConfigs)
	if chainsErr != nil {
		return nil, pkgerrors.Wrap(chainsErr, "failed to create chains")
	}

	homeChainCLDEnvironment := &deployment.Environment{
		Logger:            singeFileLogger,
		Chains:            chains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return ctx
		},
	}

	keystoneContractsInput := &keystonetypes.KeystoneContractsInput{
		ChainSelector: homeChainOutput.ChainSelector,
		CldEnv:        homeChainCLDEnvironment,
	}
	keystoneContractsOutput, keyContrErr := libcontracts.DeployKeystone(testLogger, keystoneContractsInput)
	if keyContrErr != nil {
		return nil, pkgerrors.Wrap(keyContrErr, "failed to deploy keystone contracts")
	}

	// Target chains where we deploy Forwarders and other contracts
	// we need another CLD environment to deploy forwarders for all the chains except the home chain
	targetChainsConfigs := make([]devenv.ChainConfig, 0)
	for idx, bcOut := range blockchainsOutput {
		if idx == 0 {
			// skip home chain
			continue
		}
		targetChainsConfigs = append(targetChainsConfigs, devenv.ChainConfig{
			ChainID:   bcOut.SethClient.Cfg.Network.ChainID,
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

	targetChains, chainsErr := devenv.NewChains(singeFileLogger, targetChainsConfigs)
	if chainsErr != nil {
		return nil, pkgerrors.Wrap(chainsErr, "failed to create chains")
	}

	targetChainsCLDEnvironment := &deployment.Environment{
		Logger:            singeFileLogger,
		Chains:            targetChains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return ctx
		},
	}

	for idx, chain := range blockchainsOutput {
		if idx == 0 {
			// skip home chain
			continue
		}
		_, err := libcontracts.DeployKeystoneForwarder(testLogger, targetChainsCLDEnvironment, chain.ChainSelector)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to deploy Keystone Forwarder contract")
		}
	}

	// merge all the addresses for home and target chains
	mergeErr := homeChainCLDEnvironment.ExistingAddresses.Merge(targetChainsCLDEnvironment.ExistingAddresses)
	if mergeErr != nil {
		return nil, pkgerrors.Wrap(mergeErr, "failed to merge existing addresses")
	}

	mergedAddressBook := homeChainCLDEnvironment.ExistingAddresses

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topology, topoErr := libdon.BuildTopology(input.CapabilitiesAwareNodeSets, *homeChainInput.infraInput, homeChainOutput.ChainSelector)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}

	// get chainIDs, they'll be used for identifying ETH keys and Forwarder addresses
	chainIDs := make([]int, 0)
	for _, bcOut := range blockchainsOutput {
		chainIDs = append(chainIDs, int(bcOut.ChainID))
	}

	// Generate EVM and P2P keys or read them from the config
	// That way we can pass them final configs and do away with restarting the nodes
	var keys *keystonetypes.GenerateKeysOutput

	keysOutput, keysOutputErr := cresecrets.KeysOutputFromConfig(input.CapabilitiesAwareNodeSets)
	if keysOutputErr != nil {
		return nil, pkgerrors.Wrap(keysOutputErr, "failed to generate keys output")
	}

	generateKeysInput := &keystonetypes.GenerateKeysInput{
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

	// Configure Workflow Registry contract
	workflowRegistryInput := &keystonetypes.WorkflowRegistryInput{
		ChainSelector:  homeChainOutput.ChainSelector,
		CldEnv:         homeChainCLDEnvironment,
		AllowedDonIDs:  []uint32{topology.WorkflowDONID},
		WorkflowOwners: []common.Address{homeChainOutput.SethClient.MustGetRootKeyAddress()},
	}

	_, workflowErr := libcontracts.ConfigureWorkflowRegistry(testLogger, workflowRegistryInput)
	if workflowErr != nil {
		return nil, pkgerrors.Wrap(workflowErr, "failed to configure workflow registry")
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

		bcOuts := make([]*blockchain.Output, 0)
		for _, bcOut := range blockchainsOutput {
			bcOuts = append(bcOuts, bcOut.BlockchainOutput)
		}

		// generate configs only if they are not provided
		if configsFound == 0 {
			config, configErr := keystoneporconfig.GenerateConfigs(
				keystonetypes.GeneratePoRConfigsInput{
					DonMetadata:            donMetadata,
					BlockchainOutput:       bcOuts,
					Flags:                  donMetadata.Flags,
					PeeringData:            peeringData,
					AddressBook:            mergedAddressBook,
					GatewayConnectorOutput: topology.GatewayConnectorOutput,
					HomeChainSelector:      topology.HomeChainSelector,
				},
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
			secretsInput := &keystonetypes.GenerateSecretsInput{
				DonMetadata: donMetadata,
			}

			if evmKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
				secretsInput.EVMKeys = evmKeys
			}

			if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
				secretsInput.P2PKeys = p2pKeys
			}

			// EVM and P2P keys will be provided to nodes as secrets
			secrets, secretsErr := keystonesecrets.GenerateSecrets(
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

	if input.InfraInput.InfraType == libtypes.CRIB {
		testLogger.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &keystonetypes.DeployCribDonsInput{
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

		deployCribJdInput := &keystonetypes.DeployCribJdInput{
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

	jdOutput, jdErr := CreateJobDistributor(&input.JdInput)
	if jdErr != nil {
		jdErr = fmt.Errorf("failed to start JD container for image %s: %w", input.JdInput.Image, jdErr)

		// useful end user messages
		if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
			jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
		}
		return nil, jdErr
	}

	nodeSetOutput := make([]*keystonetypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))
	for _, nodeSetInput := range input.CapabilitiesAwareNodeSets {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, homeChainOutput.BlockchainOutput)
		if nodesetErr != nil {
			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}

		nodeSetOutput = append(nodeSetOutput, &keystonetypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	bcOuts := make([]*blockchain.Output, 0)
	sethClients := make([]*seth.Client, 0)
	for _, bcOut := range blockchainsOutput {
		bcOuts = append(bcOuts, bcOut.BlockchainOutput)
		sethClients = append(sethClients, bcOut.SethClient)
	}

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &keystonetypes.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutputs: bcOuts,
		SethClients:       sethClients,
		NodeSetOutput:     nodeSetOutput,
		ExistingAddresses: homeChainCLDEnvironment.ExistingAddresses,
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

	fullCldOutput, cldErr := libdevenv.BuildFullCLDEnvironment(singeFileLogger, fullCldInput, creds)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to build full CLD environment")
	}

	// Fund the nodes
	for _, metaDon := range fullCldOutput.DonTopology.DonsWithMetadata {
		for _, node := range metaDon.DON.Nodes {
			for _, bcOut := range blockchainsOutput {
				_, fundingErr := libfunding.SendFunds(zerolog.Logger{}, bcOut.SethClient, libtypes.FundsToSend{
					ToAddress: common.HexToAddress(
						node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]),
					Amount:     big.NewInt(5000000000000000000),
					PrivateKey: bcOut.SethClient.MustGetRootPrivateKey(),
				})
				if fundingErr != nil {
					return nil, pkgerrors.Wrapf(fundingErr, "failed to fund node %s",
						node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)])
				}
			}
		}
	}

	donToJobSpecs := make(keystonetypes.DonsToJobSpecs)

	for _, jobSpecGeneratingFn := range input.JobSpecFactoryFunctions {
		singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cretypes.JobSpecFactoryInput{
			CldEnvironment:          fullCldOutput.Environment,
			BlockchainOutput:        homeChainOutput.BlockchainOutput,
			DonTopology:             fullCldOutput.DonTopology,
			KeystoneContractsOutput: keystoneContractsOutput,
		})
		if jobSpecsErr != nil {
			return nil, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
		}
		mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
	}

	createJobsInput := keystonetypes.CreateJobsInput{
		CldEnv:        fullCldOutput.Environment,
		DonTopology:   fullCldOutput.DonTopology,
		DonToJobSpecs: donToJobSpecs,
	}

	jobsErr := libdon.CreateJobs(testLogger, createJobsInput)
	if jobsErr != nil {
		return nil, pkgerrors.Wrap(jobsErr, "failed to create jobs")
	}

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	testLogger.Info().Msg("Waiting for ConfigWatcher health check")

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
	testLogger.Info().Msg("Proceeding to set OCR3 and Keystone configuration...")

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := keystonetypes.ConfigureKeystoneInput{
		ChainSelector: homeChainOutput.ChainSelector,
		CldEnv:        fullCldOutput.Environment,
		Topology:      topology,
	}

	keystoneErr := libcontracts.ConfigureKeystone(configureKeystoneInput, input.CapabilitiesContractFactoryFunctions)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	return &SetupOutput{
		KeystoneContractsOutput:             keystoneContractsOutput,
		WorkflowRegistryConfigurationOutput: workflowRegistryInput.Out, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    blockchainsOutput,
		DonTopology:                         fullCldOutput.DonTopology,
		NodeOutput:                          nodeSetOutput,
		CldEnvironment:                      fullCldOutput.Environment,
	}, nil
}

type BlockchainsInput struct {
	blockchainInput *blockchain.Input
	infraInput      *libtypes.InfraInput
	nixShell        *libnix.Shell
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
	input []BlockchainsInput,
) ([]*BlockchainOutput, error) {
	if len(input) == 0 {
		return nil, pkgerrors.New("blockchain input is nil")
	}
	blockchainOutput := make([]*BlockchainOutput, 0)
	for _, bi := range input {
		var bcOut *blockchain.Output
		var bcErr error
		if bi.infraInput.InfraType == libtypes.CRIB {
			if bi.nixShell == nil {
				return nil, pkgerrors.New("nix shell is nil")
			}

			deployCribBlockchainInput := &keystonetypes.DeployCribBlockchainInput{
				BlockchainInput: bi.blockchainInput,
				NixShell:        bi.nixShell,
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
			bcOut, bcErr = blockchain.NewBlockchainNetwork(bi.blockchainInput)
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

func mergeJobSpecSlices(from, to keystonetypes.DonsToJobSpecs) {
	for fromDonID, fromJobSpecs := range from {
		if _, ok := to[fromDonID]; !ok {
			to[fromDonID] = make([]*jobv1.ProposeJobRequest, 0)
		}
		to[fromDonID] = append(to[fromDonID], fromJobSpecs...)
	}
}
