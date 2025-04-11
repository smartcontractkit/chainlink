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
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

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
	keystonepor "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/por"
	keystonesecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
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
	BlockchainOutput                    *BlockchainOutput
	DonTopology                         *keystonetypes.DonTopology
	NodeOutput                          []*keystonetypes.WrappedNodeOutput
}

type SetupInput struct {
	ExtraAllowedPorts          []int
	CapabilitiesAwareNodeSets  []*keystonetypes.CapabilitiesAwareNodeSet
	CapabilityFactoryFunctions []func([]cretypes.CapabilityFlag) []keystone_changeset.DONCapabilityWithConfig
	JobSpecFactoryFunctions    []cretypes.JobSpecFactoryFn
	BlockchainsInput           blockchain.Input
	JdInput                    jd.Input
	InfraInput                 libtypes.InfraInput
	CustomBinariesPaths        map[cretypes.CapabilityFlag]string
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

	blockchainsInput := BlockchainsInput{
		blockchainInput: &input.BlockchainsInput,
		infraInput:      &input.InfraInput,
		nixShell:        nixShell,
	}

	blockchainsOutput, bcOutErr := CreateBlockchains(singeFileLogger, testLogger, blockchainsInput)
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to create blockchains")
	}

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	// but first, we need to create deployment.Environment that will contain only chain information in order to deploy contracts with the CLD
	chainsConfig := []devenv.ChainConfig{
		{
			ChainID:   blockchainsOutput.SethClient.Cfg.Network.ChainID,
			ChainName: blockchainsOutput.SethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(blockchainsOutput.BlockchainOutput.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: blockchainsOutput.BlockchainOutput.Nodes[0].ExternalWSUrl,
				Internal: blockchainsOutput.BlockchainOutput.Nodes[0].InternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: blockchainsOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
				Internal: blockchainsOutput.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			}},
			DeployerKey: blockchainsOutput.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
		},
	}

	chains, chainsErr := devenv.NewChains(singeFileLogger, chainsConfig)
	if chainsErr != nil {
		return nil, pkgerrors.Wrap(chainsErr, "failed to create chains")
	}

	chainsOnlyCld := &deployment.Environment{
		Logger:            singeFileLogger,
		Chains:            chains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return ctx
		},
	}

	keystoneContractsInput := &keystonetypes.KeystoneContractsInput{
		ChainSelector: blockchainsOutput.ChainSelector,
		CldEnv:        chainsOnlyCld,
	}
	keystoneContractsOutput, keyContrErr := libcontracts.DeployKeystone(testLogger, keystoneContractsInput)
	if keyContrErr != nil {
		return nil, pkgerrors.Wrap(keyContrErr, "failed to deploy keystone contracts")
	}

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topology, topoErr := libdon.BuildTopology(input.CapabilitiesAwareNodeSets, *blockchainsInput.infraInput)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}

	// Generate EVM and P2P keys, which are needed to prepare the node configs
	// That way we can pass them final configs and do away with restarting the nodes
	var keys *keystonetypes.GenerateKeysOutput
	chainIDInt, chainErr := strconv.Atoi(blockchainsOutput.BlockchainOutput.ChainID)
	if chainErr != nil {
		return nil, pkgerrors.Wrap(chainErr, "failed to convert chain ID to int")
	}

	generateKeysInput := &keystonetypes.GenerateKeysInput{
		GenerateEVMKeysForChainIDs: []int{chainIDInt},
		GenerateP2PKeys:            true,
		Topology:                   topology,
		Password:                   "", // since the test runs on private ephemeral blockchain we don't use real keys and do not care a lot about the password
	}
	keys, keysErr := libdon.GenereteKeys(generateKeysInput)
	if keysErr != nil {
		return nil, pkgerrors.Wrap(keysErr, "failed to generate keys")
	}

	topology, addKeysErr := libdon.AddKeysToTopology(topology, keys)
	if addKeysErr != nil {
		return nil, pkgerrors.Wrap(addKeysErr, "failed to add keys to topology")
	}

	// Configure Workflow Registry contract
	workflowRegistryInput := &keystonetypes.WorkflowRegistryInput{
		ChainSelector:  blockchainsOutput.ChainSelector,
		CldEnv:         chainsOnlyCld,
		AllowedDonIDs:  []uint32{topology.WorkflowDONID},
		WorkflowOwners: []common.Address{blockchainsOutput.SethClient.MustGetRootKeyAddress()},
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
		config, configErr := keystoneporconfig.GenerateConfigs(
			keystonetypes.GeneratePoRConfigsInput{
				DonMetadata:                 donMetadata,
				BlockchainOutput:            blockchainsOutput.BlockchainOutput,
				DonID:                       donMetadata.ID,
				Flags:                       donMetadata.Flags,
				PeeringData:                 peeringData,
				CapabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
				WorkflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
				ForwarderAddress:            keystoneContractsOutput.ForwarderAddress,
				GatewayConnectorOutput:      topology.GatewayConnectorOutput,
			},
		)
		if configErr != nil {
			return nil, pkgerrors.Wrap(configErr, "failed to generate config")
		}

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
			input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
			input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.TestSecretsOverrides = secrets[j]
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

	nodeOutput := make([]*keystonetypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))
	for _, nodeSetInput := range input.CapabilitiesAwareNodeSets {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, blockchainsOutput.BlockchainOutput)
		if nodesetErr != nil {
			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}

		nodeOutput = append(nodeOutput, &keystonetypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &keystonetypes.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutput:  blockchainsOutput.BlockchainOutput,
		SethClient:        blockchainsOutput.SethClient,
		NodeSetOutput:     nodeOutput,
		ExistingAddresses: chainsOnlyCld.ExistingAddresses,
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
			_, fundingErr := libfunding.SendFunds(zerolog.Logger{}, blockchainsOutput.SethClient, libtypes.FundsToSend{
				ToAddress:  common.HexToAddress(node.AccountAddr[blockchainsOutput.SethClient.Cfg.Network.ChainID]),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: blockchainsOutput.SethClient.MustGetRootPrivateKey(),
			})
			if fundingErr != nil {
				return nil, pkgerrors.Wrapf(fundingErr, "failed to fund node %s", node.AccountAddr[blockchainsOutput.SethClient.Cfg.Network.ChainID])
			}
		}
	}

	donToJobSpecs := make(keystonetypes.DonsToJobSpecs)

	for _, jobSpecGeneratingFn := range input.JobSpecFactoryFunctions {
		singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cretypes.JobSpecFactoryInput{
			CldEnvironment:          fullCldOutput.Environment,
			BlockchainOutput:        blockchainsOutput.BlockchainOutput,
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

	for idx, nodeSetOut := range nodeOutput {
		if flags.HasFlag(input.CapabilitiesAwareNodeSets[idx].DONTypes, cretypes.GatewayDON) || flags.HasFlag(input.CapabilitiesAwareNodeSets[idx].DONTypes, cretypes.CapabilitiesDON) {
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
		ChainSelector: blockchainsOutput.ChainSelector,
		CldEnv:        fullCldOutput.Environment,
		Topology:      topology,
	}

	keystoneErr := libcontracts.ConfigureKeystone(configureKeystoneInput, input.CapabilityFactoryFunctions)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	return &SetupOutput{
		KeystoneContractsOutput:             keystoneContractsOutput,
		WorkflowRegistryConfigurationOutput: workflowRegistryInput.Out, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    blockchainsOutput,
		DonTopology:                         fullCldOutput.DonTopology,
		NodeOutput:                          nodeOutput,
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
	cldLogger logger.Logger,
	testLogger zerolog.Logger,
	input BlockchainsInput,
) (*BlockchainOutput, error) {
	if input.blockchainInput == nil {
		return nil, pkgerrors.New("blockchain input is nil")
	}

	if input.infraInput.InfraType == libtypes.CRIB {
		if input.nixShell == nil {
			return nil, pkgerrors.New("nix shell is nil")
		}

		deployCribBlockchainInput := &keystonetypes.DeployCribBlockchainInput{
			BlockchainInput: input.blockchainInput,
			NixShell:        input.nixShell,
			CribConfigsDir:  cribConfigsDir,
		}

		var blockchainErr error
		input.blockchainInput.Out, blockchainErr = crib.DeployBlockchain(deployCribBlockchainInput)
		if blockchainErr != nil {
			return nil, pkgerrors.Wrap(blockchainErr, "failed to deploy blockchain")
		}
	}

	// Create a new blockchain network and Seth client to interact with it
	blockchainOutput, err := blockchain.NewBlockchainNetwork(input.blockchainInput)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create blockchain network")
	}

	pkey := os.Getenv("PRIVATE_KEY")
	if pkey == "" {
		return nil, pkgerrors.New("PRIVATE_KEY env var must be set")
	}

	err = keystonepor.WaitForRPCEndpoint(testLogger, blockchainOutput.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "RPC endpoint not available")
	}

	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(blockchainOutput.Nodes[0].ExternalWSUrl).
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

	return &BlockchainOutput{
		ChainSelector:      chainSelector,
		ChainID:            sethClient.Cfg.Network.ChainID,
		BlockchainOutput:   blockchainOutput,
		SethClient:         sethClient,
		DeployerPrivateKey: pkey,
		c:                  blockchainOutput,
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

func mergeJobSpecSlices(from, to keystonetypes.DonsToJobSpecs) {
	for fromDonID, fromJobSpecs := range from {
		if _, ok := to[fromDonID]; !ok {
			to[fromDonID] = make([]*jobv1.ProposeJobRequest, 0)
		}
		to[fromDonID] = append(to[fromDonID], fromJobSpecs...)
	}
}
