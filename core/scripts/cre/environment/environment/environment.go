package environment

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/deploy"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/verify"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	logeventtriggercap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/logevent"
	readcontractcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/readcontract"
	webapicap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapi"
	writeevmcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writeevm"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	crelogevent "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/logevent"
	crereadcontract "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/readcontract"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/webapi"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

const manualCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove containers with 'ctf' label and delete their volumes`

var topologyFlag string
var waitOnErrorTimeoutFlag string
var extraAllowedGatewayPortsFlag []int
var withExampleFlag bool
var exampleWorkflowTimeoutFlag string
var withPluginsDockerImageFlag string

func init() {
	EnvironmentCmd.AddCommand(startCmd)
	EnvironmentCmd.AddCommand(stopCmd)
	EnvironmentCmd.AddCommand(deployAndVerifyExampleWorkflowCmd)

	startCmd.Flags().StringVarP(&topologyFlag, "topology", "t", "simplified", "Topology to use for the environment (simiplified or full)")
	startCmd.Flags().StringVarP(&waitOnErrorTimeoutFlag, "wait-on-error-timeout", "w", "", "Wait on error timeout (e.g. 10s, 1m, 1h)")
	startCmd.Flags().IntSliceVarP(&extraAllowedGatewayPortsFlag, "extra-allowed-gateway-ports", "e", []int{}, "Extra allowed ports for outgoing connections from the Gateway DON (e.g. 8080,8081)")
	startCmd.Flags().BoolVarP(&withExampleFlag, "with-example", "x", false, "Deploy and register example workflow")
	startCmd.Flags().StringVarP(&exampleWorkflowTimeoutFlag, "example-workflow-timeout", "u", "5m", "Time to wait until example workflow succeeds")
	startCmd.Flags().StringVarP(&withPluginsDockerImageFlag, "with-plugins-docker-image", "p", "", "Docker image to use (must have all capabilities included)")

	deployAndVerifyExampleWorkflowCmd.Flags().StringVarP(&rpcURL, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	deployAndVerifyExampleWorkflowCmd.Flags().Uint64VarP(&chainID, "chain-id", "c", 1337, "Chain ID")
}

var waitOnErrorTimeoutDurationFn = func() {
	if waitOnErrorTimeoutFlag != "" {
		waitOnErrorTimeoutDuration, err := time.ParseDuration(waitOnErrorTimeoutFlag)
		if err != nil {
			return
		}

		fmt.Printf("Waiting %s on error before cleanup\n", waitOnErrorTimeoutFlag)
		time.Sleep(waitOnErrorTimeoutDuration)
	}
}

var EnvironmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment commands",
	Long:  `Commands to manage the environment`,
}

const (
	TopologySimplified = "simplified"
	TopologyFull       = "full"
)

type Config struct {
	Blockchains       []*blockchain.Input     `toml:"blockchains" validate:"required"`
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

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the environment",
	Long:  `Start the local CRE environment with all supported capabilities`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// remove all containers before starting the environment, just in case
		_ = framework.RemoveTestContainers()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			sig := <-sigCh
			fmt.Printf("\nReceived signal: %s\n", sig)

			removeErr := framework.RemoveTestContainers()
			if removeErr != nil {
				fmt.Fprint(os.Stderr, removeErr, manualCleanupMsg)
			}

			os.Exit(1)
		}()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func() {
			p := recover()

			if p != nil {
				fmt.Println("Panicked when starting environment")
				if err, ok := p.(error); ok {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err)
					fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))
				} else {
					fmt.Fprintf(os.Stderr, "panic: %v\n", p)
					fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))
				}

				waitOnErrorTimeoutDurationFn()

				removeErr := framework.RemoveTestContainers()
				if removeErr != nil {
					fmt.Fprint(os.Stderr, errors.Wrap(removeErr, manualCleanupMsg).Error())
				}
			}
		}()

		if topologyFlag != TopologySimplified && topologyFlag != TopologyFull {
			return fmt.Errorf("invalid topology: %s. Valid topologies are: %s, %s", topologyFlag, TopologySimplified, TopologyFull)
		}

		printCRELogo()
		startTime := time.Now()

		if os.Getenv("CTF_CONFIGS") == "" {
			// use default config
			if topologyFlag == TopologySimplified {
				setErr := os.Setenv("CTF_CONFIGS", "configs/single-don.toml")
				if setErr != nil {
					return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
				}
			} else {
				setErr := os.Setenv("CTF_CONFIGS", "configs/workflow-capabilities-don.toml")
				if setErr != nil {
					return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
				}
			}
			fmt.Printf("Set CTF_CONFIGS environment variable to default value: %s\n", os.Getenv("CTF_CONFIGS"))
		}

		if os.Getenv("PRIVATE_KEY") == "" {
			setErr := os.Setenv("PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
			if setErr != nil {
				return fmt.Errorf("failed to set PRIVATE_KEY environment variable: %w", setErr)
			}
			fmt.Printf("Set PRIVATE_KEY environment variable to default value: %s\n", os.Getenv("PRIVATE_KEY"))
		}

		// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
		setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if setErr != nil {
			return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
		}

		cmdContext := cmd.Context()

		output, err := startCLIEnvironment(cmdContext, topologyFlag, withPluginsDockerImageFlag, withExampleFlag, extraAllowedGatewayPortsFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

			waitOnErrorTimeoutDurationFn()
			removeErr := framework.RemoveTestContainers()
			if removeErr != nil {
				return errors.Wrap(removeErr, manualCleanupMsg)
			}

			return errors.Wrap(err, "failed to start environment")
		}

		homeChainOut := output.BlockchainOutput[0]

		sErr := func() error {
			rpcs := map[uint64]string{}
			for _, bcOut := range output.BlockchainOutput {
				rpcs[bcOut.ChainSelector] = bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl
			}
			creCLISettingsFile, settingsErr := crecli.PrepareCRECLISettingsFile(
				crecli.CRECLIProfile,
				homeChainOut.SethClient.MustGetRootKeyAddress(),
				output.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				output.DonTopology.WorkflowDonID,
				homeChainOut.ChainSelector,
				rpcs,
			)

			if settingsErr != nil {
				return settingsErr
			}

			// Copy the file to current directory as cre.yaml
			currentDir, cErr := os.Getwd()
			if cErr != nil {
				return cErr
			}

			targetPath := filepath.Join(currentDir, "cre.yaml")
			input, err := os.ReadFile(creCLISettingsFile.Name())
			if err != nil {
				return err
			}
			err = os.WriteFile(targetPath, input, 0600)
			if err != nil {
				return err
			}

			fmt.Printf("CRE CLI settings file created: %s\n", targetPath)

			return nil
		}()

		if sErr != nil {
			fmt.Fprintf(os.Stderr, "failed to create CRE CLI settings file: %s. You need to create it manually.", sErr)
		}

		// TODO print urls?
		fmt.Print(libformat.PurpleText("\nEnvironment setup completed successfully in %.2f seconds\n\n", time.Since(startTime).Seconds()))
		fmt.Print("To terminate execute: ctf d rm\n\n")

		if withExampleFlag {
			fmt.Print(libformat.PurpleText("\nRegistering and verifying example workflow\n\n"))
			deployErr := deployAndVerifyExampleWorkflow(cmdContext, homeChainOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl, homeChainOut.ChainID)
			if deployErr != nil {
				fmt.Printf("Failed to deploy and verify example workflow: %s\n", deployErr)
			}

			fmt.Print(libformat.PurpleText("\nEnvironment setup completed successfully in %.2f seconds\n\n", time.Since(startTime).Seconds()))
			fmt.Print("To terminate execute: ctf d rm\n\n")
		}

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

var (
	rpcURL  string
	chainID uint64
)

var deployAndVerifyExampleWorkflowCmd = &cobra.Command{
	Use:   "deploy-verify-example",
	Short: "Deploys and verifies example (optionally)",
	Long:  `Deploys a simple Proof-of-Reserve workflow and, optionally, wait until it succeeds`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return deployAndVerifyExampleWorkflow(cmd.Context(), rpcURL, chainID)
	},
}

func startCLIEnvironment(cmdContext context.Context, topologyFlag string, withPluginsDockerImageFlag string, withExampleFlag bool, extraAllowedGatewayPorts []int) (*creenv.SetupOutput, error) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[Config](nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load test configuration: %w", err)
	}

	// make sure that either cron is enabled or withPluginsDockerImageFlag is set
	if withExampleFlag && (in.ExtraCapabilities.CronCapabilityBinaryPath == "" && withPluginsDockerImageFlag == "") {
		return nil, fmt.Errorf("either cron binary path must be set in TOML config (%s) or you must use Docker image with all capabilities included and passed via withPluginsDockerImageFlag", os.Getenv("CTF_CONFIGS"))
	}

	capabilitiesBinaryPaths := map[cretypes.CapabilityFlag]string{}
	var capabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet

	if topologyFlag == TopologySimplified {
		if len(in.NodeSets) != 1 {
			return nil, fmt.Errorf("expected 1 nodeset, got %d", len(in.NodeSets))
		}
		// add support for more binaries if needed
		workflowDONCapabilities := []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WriteEVMCapability, cretypes.WebAPITriggerCapability, cretypes.WebAPITargetCapability}
		if in.ExtraCapabilities.CronCapabilityBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.CronCapability)
			capabilitiesBinaryPaths[cretypes.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
		}

		if in.ExtraCapabilities.LogEventTriggerBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.LogTriggerCapability)
			capabilitiesBinaryPaths[cretypes.LogTriggerCapability] = in.ExtraCapabilities.LogEventTriggerBinaryPath
		}

		if in.ExtraCapabilities.ReadContractBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.ReadContractCapability)
			capabilitiesBinaryPaths[cretypes.ReadContractCapability] = in.ExtraCapabilities.ReadContractBinaryPath
		}

		capabilitiesAwareNodeSets = []*cretypes.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],
				Capabilities:       workflowDONCapabilities,
				DONTypes:           []string{cretypes.WorkflowDON, cretypes.GatewayDON},
				BootstrapNodeIndex: 0,
				GatewayNodeIndex:   0,
			},
		}
	} else {
		if len(in.NodeSets) != 3 {
			return nil, fmt.Errorf("expected 3 nodesets, got %d", len(in.NodeSets))
		}

		// add support for more binaries if needed
		workflowDONCapabilities := []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WebAPITriggerCapability}
		if in.ExtraCapabilities.CronCapabilityBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.CronCapability)
			capabilitiesBinaryPaths[cretypes.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
		}

		if in.ExtraCapabilities.LogEventTriggerBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.LogTriggerCapability)
			capabilitiesBinaryPaths[cretypes.LogTriggerCapability] = in.ExtraCapabilities.LogEventTriggerBinaryPath
		}

		capabiliitesDONCapabilities := []string{cretypes.WriteEVMCapability, cretypes.WebAPITargetCapability}
		if in.ExtraCapabilities.ReadContractBinaryPath != "" || withPluginsDockerImageFlag != "" {
			capabiliitesDONCapabilities = append(capabiliitesDONCapabilities, cretypes.ReadContractCapability)
			capabilitiesBinaryPaths[cretypes.ReadContractCapability] = in.ExtraCapabilities.ReadContractBinaryPath
		}

		capabilitiesAwareNodeSets = []*cretypes.CapabilitiesAwareNodeSet{
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
				BootstrapNodeIndex: -1,                                 // <----- it's crucial to indicate there's no bootstrap node
			},
			{
				Input:              in.NodeSets[2],
				Capabilities:       []string{},
				DONTypes:           []string{cretypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                            // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
	}

	// unset DockerFilePath and DockerContext as we cannot use them with existing images
	if withPluginsDockerImageFlag != "" {
		for setIdx := range capabilitiesAwareNodeSets {
			for nodeIdx := range capabilitiesAwareNodeSets[setIdx].Input.NodeSpecs {
				capabilitiesAwareNodeSets[setIdx].Input.NodeSpecs[nodeIdx].Node.Image = withPluginsDockerImageFlag
				capabilitiesAwareNodeSets[setIdx].Input.NodeSpecs[nodeIdx].Node.DockerContext = ""
				capabilitiesAwareNodeSets[setIdx].Input.NodeSpecs[nodeIdx].Node.DockerFilePath = ""
			}
		}
	}

	fmt.Print(libformat.PurpleText("DON topology:\n"))
	for _, nodeSet := range capabilitiesAwareNodeSets {
		fmt.Print(libformat.PurpleText("%s\n", strings.ToUpper(nodeSet.Input.Name)))
		fmt.Print(libformat.PurpleText("\tNode count: %d\n", len(nodeSet.Input.NodeSpecs)))
		capabilitiesDesc := "none"
		if len(nodeSet.Capabilities) > 0 {
			capabilitiesDesc = strings.Join(nodeSet.Capabilities, ", ")
		}
		fmt.Print(libformat.PurpleText("\tCapabilities: %s\n", capabilitiesDesc))
		fmt.Print(libformat.PurpleText("\tDON Types: %s\n\n", strings.Join(nodeSet.DONTypes, ", ")))
	}

	// add support for more capabilities if needed
	capabilityFactoryFns := []cretypes.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPITriggerCapabilityFactoryFn,
		webapicap.WebAPITargetCapabilityFactoryFn,
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
	}

	containerPath, pathErr := crecapabilities.DefaultContainerDirectory(in.Infra.InfraType)
	if pathErr != nil {
		return nil, fmt.Errorf("failed to get default container directory: %w", pathErr)
	}

	homeChainIDInt, chainErr := strconv.Atoi(in.Blockchains[0].ChainID)
	if chainErr != nil {
		return nil, fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
	}

	cronBinaryName := filepath.Base(in.ExtraCapabilities.CronCapabilityBinaryPath)
	if withPluginsDockerImageFlag != "" {
		cronBinaryName = "cron"
	}

	logEventTriggerBinaryName := filepath.Base(in.ExtraCapabilities.LogEventTriggerBinaryPath)
	if withPluginsDockerImageFlag != "" {
		logEventTriggerBinaryName = "log-event-trigger"
	}

	readContractBinaryName := filepath.Base(in.ExtraCapabilities.ReadContractBinaryPath)
	if withPluginsDockerImageFlag != "" {
		readContractBinaryName = "readcontract"
	}

	jobSpecFactoryFunctions := []cretypes.JobSpecFactoryFn{
		// add support for more job spec factory functions if needed
		webapi.WebAPITriggerJobSpecFactoryFn,
		webapi.WebAPITargetJobSpecFactoryFn,
		creconsensus.ConsensusJobSpecFactoryFn(libc.MustSafeUint64(int64(homeChainIDInt))),
		crecron.CronJobSpecFactoryFn(filepath.Join(containerPath, cronBinaryName)),
		cregateway.GatewayJobSpecFactoryFn(extraAllowedGatewayPorts, []string{}, []string{"0.0.0.0/0"}),
		crecompute.ComputeJobSpecFactoryFn,
	}

	for _, blockchain := range in.Blockchains {
		chainIDInt, chainErr := strconv.Atoi(blockchain.ChainID)
		if chainErr != nil {
			return nil, fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
		}
		capabilityFactoryFns = append(capabilityFactoryFns, writeevmcap.WriteEVMCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))))
		capabilityFactoryFns = append(capabilityFactoryFns, readcontractcap.ReadContractCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)), "evm"))
		capabilityFactoryFns = append(capabilityFactoryFns, logeventtriggercap.LogEventTriggerCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)), "evm"))

		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, crelogevent.LogEventTriggerJobSpecFactoryFn(
			chainIDInt,
			"evm",
			// path within the container/pod
			filepath.Join(containerPath, logEventTriggerBinaryName),
		))

		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, crereadcontract.ReadContractJobSpecFactoryFn(
			chainIDInt,
			"evm",
			// path within the container/pod
			filepath.Join(containerPath, readContractBinaryName),
		))
	}

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            capabilitiesAwareNodeSets,
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		JobSpecFactoryFunctions:              jobSpecFactoryFunctions,
		ConfigFactoryFunctions: []cretypes.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
	}

	if withPluginsDockerImageFlag == "" {
		universalSetupInput.CustomBinariesPaths = capabilitiesBinaryPaths
	}

	ctx, cancel := context.WithTimeout(cmdContext, 10*time.Minute)
	defer cancel()
	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(ctx, testLogger, cldlogger.NewSingleFileLogger(nil), universalSetupInput)
	if setupErr != nil {
		return nil, fmt.Errorf("failed to setup test environment: %w", setupErr)
	}

	return universalSetupOutput, nil
}

func deployAndVerifyExampleWorkflow(cmdContext context.Context, rpcURL string, chainID uint64) error {
	totalStart := time.Now()
	start := time.Now()
	fmt.Print(libformat.PurpleText("[Stage 1/3] Deploying Permissionless Feeds Consumer\n\n"))
	consumerContractAddress, consumerErr := deploy.DeployPermissionlessFeedsConsumer(rpcURL)
	if consumerErr != nil {
		return errors.Wrap(consumerErr, "failed to deploy Permissionless Feeds Consumer contract")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 1/3] Deployed Permissionless Feeds Consumer in %.2f seconds\n", time.Since(start).Seconds()))

	workflowData, workflowDataErr := readWorkflowData()
	if workflowDataErr != nil {
		return errors.Wrap(workflowDataErr, "failed to read workflow data")
	}

	start = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 2/3] Registering example Proof-of-Reserve workflow\n\n"))

	deployErr := deployExampleWorkflow(chainID, *workflowData)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy example workflow")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 2/3] Registered workflow in %.2f seconds\n", time.Since(start).Seconds()))
	fmt.Print(libformat.PurpleText("[Stage 3/3] Waiting for %s for workflow to execute successuly\n\n", exampleWorkflowTimeoutFlag))
	waitTime, waitTimeErr := time.ParseDuration(exampleWorkflowTimeoutFlag)
	if waitTimeErr != nil {
		return errors.Wrapf(waitTimeErr, "failed to parse %s to time.Duration", exampleWorkflowTimeoutFlag)
	}

	var pauseWorkflow = func() {
		fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow executed in %.2f seconds\n", time.Since(totalStart).Seconds()))
		start = time.Now()
		fmt.Print(libformat.PurpleText("\n[CLEANUP] Pausing example workflow\n\n"))
		pauseErr := pauseExampleWorkflow(chainID)
		if pauseErr != nil {
			fmt.Printf("Failed to pause example workflow: %s\nPlease pause it manually\n", pauseErr)
		}

		fmt.Print(libformat.PurpleText("\n[CLEANUP] Paused example workflow in %.2f seconds\n\n", time.Since(start).Seconds()))
	}
	defer pauseWorkflow()

	// ignore return as if verification failed it will print that info
	verifyErr := verify.ProofOfReserve(rpcURL, consumerContractAddress.Hex(), workflowData.FeedID, true, waitTime)
	if verifyErr != nil {
		fmt.Print(libformat.PurpleText("\n[Stage 3/3] Example workflow failed to execute successfully in %.2f seconds\n", time.Since(totalStart).Seconds()))
		return errors.Wrap(verifyErr, "failed to verify example workflow")
	}

	if isBlockscoutRunning(cmdContext) {
		fmt.Print(libformat.PurpleText("Open http://localhost/address/0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE?tab=internal_txns to check consumer contract's transaction history\n"))
	}

	return nil
}

var creCLI = "cre_v0.2.0_darwin_arm64"
var exampleWorkflowName = "exampleworkflow"

func prepareCLIInput(chainID uint64) (*cretypes.ManageWorkflowWithCRECLIInput, error) {
	if !isCRECLIIsAvailable() {
		if downloadErr := tryToDownloadCRECLI(); downloadErr != nil {
			return nil, errors.Wrapf(downloadErr, "failed to download %s", creCLI)
		}
	}

	if os.Getenv("CRE_GITHUB_API_TOKEN") == "" {
		// set fake token to satisfy CRE CLI
		_ = os.Setenv("CRE_GITHUB_API_TOKEN", "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")
	}

	chainSelector, chainSelectorErr := chainselectors.SelectorFromChainId(chainID)
	if chainSelectorErr != nil {
		return nil, errors.Wrapf(chainSelectorErr, "failed to find chain selector for chainID %d", chainID)
	}

	CRECLIAbsPath, CRECLIAbsPathErr := creCLIAbsPath()
	if CRECLIAbsPathErr != nil {
		return nil, errors.Wrapf(CRECLIAbsPathErr, "failed to get absolute path of the %s binary", creCLI)
	}

	deployerPrivateKey := os.Getenv("PRIVATE_KEY")
	if deployerPrivateKey == "" {
		deployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	}

	privateKey, pkErr := crypto.HexToECDSA(deployerPrivateKey)
	if pkErr != nil {
		return nil, errors.Wrap(pkErr, "failed to parse the private key")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	deployerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	cliSettingsFileName := "cre.yaml"
	if _, cliFileErr := os.Stat(cliSettingsFileName); os.IsNotExist(cliFileErr) {
		return nil, errors.Wrap(cliFileErr, "CRE CLI settings file not found")
	}

	cliSettingsFile, cliSettingsFilhErr := os.OpenFile(cliSettingsFileName, os.O_RDONLY, 0600)
	if cliSettingsFilhErr != nil {
		return nil, errors.Wrap(cliSettingsFilhErr, "failed to open the CRE CLI settings file")
	}

	return &cretypes.ManageWorkflowWithCRECLIInput{
		ChainSelector:            chainSelector,
		WorkflowDonID:            1,
		WorkflowOwnerAddress:     deployerAddress,
		CRECLIPrivateKey:         deployerPrivateKey,
		CRECLIAbsPath:            CRECLIAbsPath,
		CRESettingsFile:          cliSettingsFile,
		WorkflowName:             exampleWorkflowName,
		ShouldCompileNewWorkflow: false,
		CRECLIProfile:            "test",
	}, nil
}

func deployExampleWorkflow(chainID uint64, workflowData workflowData) error {
	registerWorkflowInput, registerWorkflowInputErr := prepareCLIInput(chainID)
	if registerWorkflowInputErr != nil {
		return errors.Wrap(registerWorkflowInputErr, "failed to prepare CLI input")
	}

	registerWorkflowInput.ExistingWorkflow = &cretypes.ExistingWorkflow{
		BinaryURL: workflowData.BinaryURL,
		ConfigURL: &workflowData.ConfigURL,
	}

	registerErr := creworkflow.RegisterWithCRECLI(*registerWorkflowInput)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow")
	}

	return nil
}

func pauseExampleWorkflow(chainID uint64) error {
	pauseWorkflowInput, pauseWorkflowInputErr := prepareCLIInput(chainID)
	if pauseWorkflowInputErr != nil {
		return errors.Wrap(pauseWorkflowInputErr, "failed to prepare CLI input")
	}

	pauseErr := creworkflow.PauseWithCRECLI(*pauseWorkflowInput)
	if pauseErr != nil {
		return errors.Wrap(pauseErr, "failed to pause workflow")
	}

	return nil
}

type workflowData struct {
	BinaryURL string `json:"binary_url"`
	ConfigURL string `json:"config_url"`
	FeedID    string `json:"feed_id"`
}

func readWorkflowData() (*workflowData, error) {
	wdFileContent, wdFileErr := os.ReadFile("./examples/workflows/proof-of-reserve/workflow_data.json")
	if wdFileErr != nil {
		return nil, errors.Wrap(wdFileErr, "failed to open workflow_data.json file")
	}

	wdData := &workflowData{}
	unmarshallErr := json.Unmarshal(wdFileContent, wdData)
	if unmarshallErr != nil {
		return nil, errors.Wrap(unmarshallErr, "failed to unmarshall workflow data")
	}

	return wdData, nil
}

func isCRECLIIsAvailable() bool {
	if _, statErr := os.Stat(creCLI); statErr == nil {
		return true
	}

	pathCmd := exec.Command("which", creCLI)
	if err := pathCmd.Run(); err == nil {
		return true
	}

	return false
}

func tryToDownloadCRECLI() error {
	start := time.Now()
	fmt.Print(libformat.PurpleText("\n[Stage 2a/3] Downloading CRE CLI\n"))
	commandArgs := []string{"release", "download", "v0.2.0", "--repo", "smartcontractkit/dev-platform", "--pattern", "*darwin_arm64*", "--skip-existing"}

	ghCmd := exec.Command("gh", commandArgs...) // #nosec G204
	ghCmd.Stdout = os.Stdout
	ghCmd.Stderr = os.Stderr
	if startErr := ghCmd.Start(); startErr != nil {
		return errors.Wrap(startErr, "failed to start gh cli command")
	}

	if waitErr := ghCmd.Wait(); waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for gh cli command")
	}

	archiveName := "cre_v0.2.0_darwin_arm64.tar.gz"
	tarArgs := []string{"-xf", archiveName}
	tarCmd := exec.Command("tar", tarArgs...)

	tarCmd.Stdout = os.Stdout
	tarCmd.Stderr = os.Stderr
	if startErr := tarCmd.Start(); startErr != nil {
		return errors.Wrap(startErr, "failed to start tar command")
	}

	if waitErr := tarCmd.Wait(); waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for tar command")
	}

	removeErr := os.Remove(archiveName)
	if removeErr != nil {
		fmt.Fprintf(os.Stderr, "failed to remove %s. Please remove it manually.\n", archiveName)
	}

	fmt.Print(libformat.PurpleText("[Stage 2a/3] CRE CLI downloaded in %.2f seconds\n\n", time.Since(start).Seconds()))

	return nil
}

func creCLIAbsPath() (string, error) {
	var CRECLIAbsPath string

	_, statErr := os.Stat(creCLI)
	if statErr != nil {
		path, lookErr := exec.LookPath(creCLI)
		if lookErr != nil {
			return "", errors.Wrap(lookErr, "failed to find absolute path of the CRE CLI binary in PATH")
		}
		CRECLIAbsPath = path
	} else {
		var CRECLIAbsPathErr error
		CRECLIAbsPath, CRECLIAbsPathErr = filepath.Abs(creCLI)
		if CRECLIAbsPathErr != nil {
			return "", errors.Wrap(CRECLIAbsPathErr, "failed to find absolute path of the CRE CLI binary in current directory")
		}
	}

	return CRECLIAbsPath, nil
}

func isBlockscoutRunning(cmdContext context.Context) bool {
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(cmdContext, 15*time.Second)
	defer cancel()
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return false
	}

	for _, container := range containers {
		if strings.Contains(strings.ToLower(container.Names[0]), "blockscout") {
			return true
		}
	}

	return false
}

func printCRELogo() {
	blue := "\033[38;5;33m"
	reset := "\033[0m"

	fmt.Println()
	fmt.Println(blue + "	db       .d88b.   .o88b.  .d8b.  db            .o88b. d8888b. d88888b" + reset)
	fmt.Println(blue + "	88      .8P  Y8. d8P  Y8 d8' `8b 88           d8P  Y8 88  `8D 88'" + reset)
	fmt.Println(blue + "	88      88    88 8P      88ooo88 88           8P      88oobY' 88ooooo" + reset)
	fmt.Println(blue + "	88      88    88 8b      88~~~88 88           8b      88`8b   88~~~~~" + reset)
	fmt.Println(blue + "	88booo. `8b  d8' Y8b  d8 88   88 88booo.      Y8b  d8 88 `88. 88." + reset)
	fmt.Println(blue + "	Y88888P  `Y88P'   `Y88P' YP   YP Y88888P       `Y88P' 88   YD Y88888P" + reset)
	fmt.Println()
}
