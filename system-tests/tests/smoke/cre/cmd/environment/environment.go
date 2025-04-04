package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/chainreader"
	crepor "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/por"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/webapi"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

var EnvironmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment commands",
	Long:  `Commands to manage the environment`,
}

func init() {
	EnvironmentCmd.AddCommand(startCmd)
	EnvironmentCmd.AddCommand(stopCmd)
	EnvironmentCmd.AddCommand(DeployAndSetupKeystoneConsumerCmd)

	// add flags to the command
	DeployAndSetupKeystoneConsumerCmd.Flags().StringVar(&rpcHTTPURL, "rpc-http-url", "", "RPC HTTP URL")
	DeployAndSetupKeystoneConsumerCmd.Flags().StringVar(&rpcWSURL, "rpc-ws-url", "", "RPC WS URL")
	DeployAndSetupKeystoneConsumerCmd.Flags().StringVar(&forwarderAddress, "forwarder-address", "", "Forwarder address")
	DeployAndSetupKeystoneConsumerCmd.Flags().StringVar(&workflowName, "workflow-name", "", "Workflow name")
	DeployAndSetupKeystoneConsumerCmd.Flags().Uint64Var(&chainID, "chain-id", 1337, "Chain ID")

	// add required flags
	flagErr := DeployAndSetupKeystoneConsumerCmd.MarkFlagRequired("rpc-http-url")
	if flagErr != nil {
		panic(flagErr)
	}
	flagErr = DeployAndSetupKeystoneConsumerCmd.MarkFlagRequired("rpc-ws-url")
	if flagErr != nil {
		panic(flagErr)
	}
	flagErr = DeployAndSetupKeystoneConsumerCmd.MarkFlagRequired("forwarder-address")
	if flagErr != nil {
		panic(flagErr)
	}
	flagErr = DeployAndSetupKeystoneConsumerCmd.MarkFlagRequired("workflow-name")
	if flagErr != nil {
		panic(flagErr)
	}
}

const manualCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove containers with 'ctf' label and delete their volumes`

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the environment",
	Long:  `Start the local CRE environment with all supported capabilities`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func() {
			p := recover()

			if p != nil {
				fmt.Println("Panicked when starting environment")
				if err, ok := p.(error); ok {
					fmt.Fprint(os.Stderr, errors.Wrap(err, "error:\n%s").Error())
				} else {
					fmt.Fprintf(os.Stderr, "panic: %v", p)
				}

				removeErr := framework.RemoveTestContainers()
				if removeErr != nil {
					fmt.Fprint(os.Stderr, errors.Wrap(removeErr, manualCleanupMsg).Error())
				}
			}
		}()

		if os.Getenv("CTF_CONFIGS") == "" {
			return errors.New("CTF_CONFIGS environment variable is not set. It should contain paths to TOML files with configurations")
		}
		fmt.Println("Starting the environment...")

		// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
		setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if setErr != nil {
			return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
		}

		output, err := startCLIEnvironment()
		if err != nil {
			removeErr := framework.RemoveTestContainers()
			if removeErr != nil {
				return errors.Wrap(removeErr, manualCleanupMsg)
			}

			return errors.Wrap(err, "failed to start environment")
		}

		// TODO print urls?
		_ = output

		fmt.Println("Environment started successfully")
		fmt.Println("To terminate execute: ctf d rm")

		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the environment",
	Long:  `Stops the local CRE environment (if it's not running, it just fallsthrough)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		removeErr := framework.RemoveTestContainers()
		if removeErr != nil {
			fmt.Fprint(os.Stderr, errors.Wrap(removeErr, manualCleanupMsg).Error())
		}

		fmt.Println("Environment stopped successfully")
		return nil
	},
}

type Config struct {
	Blockchain        *blockchain.Input       `toml:"blockchain" validate:"required"`
	NodeSets          []*ns.Input             `toml:"nodesets" validate:"required"`
	JD                *jd.Input               `toml:"jd" validate:"required"`
	Infra             *libtypes.InfraInput    `toml:"infra" validate:"required"`
	ExtraCapabilities ExtraCapabilitiesConfig `toml:"extra_capabilities"`
}

type ExtraCapabilitiesConfig struct {
	CronCapabilityBinaryPath  string `toml:"cron_capability_binary_path"`
	LogEventTriggerBinaryPath string `toml:"log_event_trigger_binary_path"`
	ReadContractBinaryPath    string `toml:"read_contract_capability_binary_path"`
}

func startCLIEnvironment() (*creenv.SetupOutput, error) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[Config](nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load test configuration: %w", err)
	}

	if len(in.NodeSets) != 3 {
		return nil, fmt.Errorf("expected 3 nodesets, got %d", len(in.NodeSets))
	}

	capabilitiesBinaryPaths := map[cretypes.CapabilityFlag]string{}

	// add support for more binaries if needed
	workflowDONCapabilities := []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WebAPITriggerCapability}
	if in.ExtraCapabilities.CronCapabilityBinaryPath != "" {
		workflowDONCapabilities = append(workflowDONCapabilities, cretypes.CronCapability)
		capabilitiesBinaryPaths[cretypes.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
	}

	capabiliitesDONCapabilities := []string{cretypes.WriteEVMCapability, cretypes.WebAPITargetCapability}
	if in.ExtraCapabilities.LogEventTriggerBinaryPath != "" {
		capabiliitesDONCapabilities = append(capabiliitesDONCapabilities, cretypes.LogTriggerCapability)
		capabilitiesBinaryPaths[cretypes.LogTriggerCapability] = in.ExtraCapabilities.LogEventTriggerBinaryPath
	}

	if in.ExtraCapabilities.ReadContractBinaryPath != "" {
		capabiliitesDONCapabilities = append(capabiliitesDONCapabilities, cretypes.ReadContractCapability)
		capabilitiesBinaryPaths[cretypes.ReadContractCapability] = in.ExtraCapabilities.ReadContractBinaryPath
	}

	capabilitiesAwareNodeSets := []*cretypes.CapabilitiesAwareNodeSet{
		{
			Input:              in.NodeSets[0],
			Capabilities:       workflowDONCapabilities,
			DONTypes:           []string{cretypes.WorkflowDON},
			BootstrapNodeIndex: 0,
		},
		{
			Input:              in.NodeSets[1],
			Capabilities:       capabiliitesDONCapabilities,
			DONTypes:           []string{cretypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
			BootstrapNodeIndex: 0,
		},
		{
			Input:              in.NodeSets[2],
			Capabilities:       []string{},
			DONTypes:           []string{cretypes.GatewayDON}, // <----- it's crucial to set the correct DON type
			BootstrapNodeIndex: -1,                            // <----- it's crucial to indicate there's no bootstrap node
			GatewayNodeIndex:   0,
		},
	}

	fmt.Println("DON topology:")
	for _, nodeSet := range capabilitiesAwareNodeSets {
		fmt.Printf("%s\n", strings.ToUpper(nodeSet.Input.Name))
		fmt.Printf("\tNode count: %d\n", len(nodeSet.Input.NodeSpecs))
		capabilitiesDesc := "none"
		if len(nodeSet.Capabilities) > 0 {
			capabilitiesDesc = strings.Join(nodeSet.Capabilities, ", ")
		}
		fmt.Printf("\tCapabilities: %s\n", capabilitiesDesc)
		fmt.Printf("\tDON Types: %s\n", strings.Join(nodeSet.DONTypes, ", "))
		fmt.Println()
	}

	chainIDInt, chainErr := strconv.Atoi(in.Blockchain.ChainID)
	if chainErr != nil {
		return nil, fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
	}

	// add support for more capabilities if needed
	capabilityFactoryFns := []cretypes.DONCapabilityWithConfigFactoryFn{
		crecontracts.DefaultCapabilityFactoryFn,
		crecontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))),
		crecontracts.ChainReaderCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)), "evm"), // for now support only evm
		crecontracts.WebAPICapabilityFactoryFn,
	}

	chainReaderJobSpecFactoryFn := chainreader.ChainReaderJobSpecFactoryFn(
		chainIDInt,
		"evm",
		// path within the container/pod
		filepath.Join("/home/capabilities", filepath.Base(in.ExtraCapabilities.LogEventTriggerBinaryPath)),
		filepath.Join("/home/capabilities", filepath.Base(in.ExtraCapabilities.ReadContractBinaryPath)),
	)

	porJobSpecFactoryFn := crepor.PoRJobSpecFactoryFn(
		filepath.Join("/home/capabilities", filepath.Base(in.ExtraCapabilities.CronCapabilityBinaryPath)),
		[]int{},
		[]string{},
		[]string{"0.0.0.0/0"}, // allow all IPs
	)

	// add support for more job spec factory functions if needed
	jobSpecFactoryFns := []cretypes.JobSpecFactoryFn{chainReaderJobSpecFactoryFn, webapi.WebAPIJobSpecFactoryFn, porJobSpecFactoryFn}

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:  capabilitiesAwareNodeSets,
		CapabilityFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:           *in.Blockchain,
		JdInput:                    *in.JD,
		InfraInput:                 *in.Infra,
		CustomBinariesPaths:        capabilitiesBinaryPaths,
		JobSpecFactoryFunctions:    jobSpecFactoryFns,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(context.Background(), testLogger, cldlogger.NewSingleFileLogger(nil), universalSetupInput)
	if setupErr != nil {
		return nil, fmt.Errorf("failed to setup test environment: %w", setupErr)
	}

	return universalSetupOutput, nil
}

var (
	rpcHTTPURL, rpcWSURL, forwarderAddress, workflowName string
	chainID                                              uint64
)

var DeployAndSetupKeystoneConsumerCmd = &cobra.Command{
	Use:   "deploy-keystone-consumer",
	Short: "Deploy and setup keystone consumer",
	Long:  `Deploy and setup keystone consumer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkey := os.Getenv("PRIVATE_KEY")
		if pkey == "" {
			return errors.New("PRIVATE_KEY environment variable is not set")
		}

		chainSelector, err := chainselectors.SelectorFromChainId(chainID)
		if err != nil {
			return errors.Wrapf(err, "failed to get chain selector for chain id %d", chainID)
		}

		sethClient, err := seth.NewClientBuilder().
			WithRpcUrl(rpcWSURL).
			WithPrivateKeys([]string{pkey}).
			// do not check if there's a pending nonce nor check node's health
			WithProtections(false, false, seth.MustMakeDuration(time.Second)).
			Build()
		if err != nil {
			return errors.Wrap(err, "failed to create seth client")
		}

		chainsConfig := []devenv.ChainConfig{
			{
				ChainID:   chainID,
				ChainName: "doesn't matter",
				ChainType: "evm",
				WSRPCs: []devenv.CribRPCs{{
					External: rpcWSURL,
					Internal: rpcWSURL,
				}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: rpcHTTPURL,
					Internal: rpcHTTPURL,
				}},
				DeployerKey: sethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
			},
		}

		singeFileLogger := cldlogger.NewSingleFileLogger(nil)

		chains, chainsErr := devenv.NewChains(singeFileLogger, chainsConfig)
		if chainsErr != nil {
			return errors.Wrap(chainsErr, "failed to create chains")
		}

		chainsOnlyCld := &deployment.Environment{
			Logger:            singeFileLogger,
			Chains:            chains,
			ExistingAddresses: deployment.NewMemoryAddressBook(),
			GetContext:        context.Background,
		}

		addr, err := DeployAndSetupKeystoneConsumer(
			framework.L,
			sethClient,
			common.HexToAddress(forwarderAddress),
			chainsOnlyCld,
			chainSelector,
			workflowName,
		)

		if err != nil {
			return errors.Wrap(err, "failed to deploy keystone consumer")
		}
		fmt.Printf("Deployed keystone consumer at %s\n", addr.Hex())

		return nil
	},
}

func DeployAndSetupKeystoneConsumer(
	testLogger zerolog.Logger,
	sethClient *seth.Client,
	forwarderAddress common.Address,
	cldEnv *deployment.Environment,
	chainSelector uint64,
	workflowName string,
) (*common.Address, error) {
	deployFeedConsumerInput := &cretypes.DeployFeedConsumerInput{
		ChainSelector: chainSelector,
		CldEnv:        cldEnv,
	}
	deployFeedsConsumerOutput, err := libcontracts.DeployFeedsConsumer(testLogger, deployFeedConsumerInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy feeds consumer")
	}

	configureFeedConsumerInput := &cretypes.ConfigureFeedConsumerInput{
		SethClient:            sethClient,
		FeedConsumerAddress:   deployFeedsConsumerOutput.FeedConsumerAddress,
		AllowedSenders:        []common.Address{forwarderAddress},
		AllowedWorkflowOwners: []common.Address{sethClient.MustGetRootKeyAddress()},
		AllowedWorkflowNames:  []string{workflowName},
	}
	_, err = libcontracts.ConfigureFeedsConsumer(testLogger, configureFeedConsumerInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure feeds consumer")
	}

	return &deployFeedsConsumerOutput.FeedConsumerAddress, nil
}
