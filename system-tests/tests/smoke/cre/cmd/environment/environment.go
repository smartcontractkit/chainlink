package environment

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	chainreadercap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/chainreader"
	chainwritercap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/chainwriter"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	webapicap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapi"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/chainreader"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/webapi"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

var EnvironmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment commands",
	Long:  `Commands to manage the environment`,
}

func init() {
	EnvironmentCmd.AddCommand(startCmd)
	EnvironmentCmd.AddCommand(stopCmd)

	startCmd.Flags().StringVarP(&topologyFlag, "topology", "t", "simplified", "Topology to use for the environment (simiplified or full)")
	startCmd.Flags().StringVarP(&waitOnErrorTimeoutFlag, "wait-on-error-timeout", "w", "", "Wait on error timeout (e.g. 10s, 1m, 1h)")
	startCmd.Flags().IntSliceVarP(&extraAllowedPortsFlag, "extra-allowed-ports", "e", []int{}, "Extra allowed ports (e.g. 8080,8081)")
}

const manualCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove containers with 'ctf' label and delete their volumes`

var topologyFlag string
var waitOnErrorTimeoutFlag string
var extraAllowedPortsFlag []int

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
					fmt.Fprint(os.Stderr, errors.Wrap(err, "error:\n%s").Error())
				} else {
					fmt.Fprintf(os.Stderr, "panic: %v", p)
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

		fmt.Println("Starting the environment...")

		// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
		setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if setErr != nil {
			return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
		}

		output, err := startCLIEnvironment(topologyFlag, extraAllowedPortsFlag)
		if err != nil {
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
				homeChainOut.SethClient.MustGetRootKeyAddress(),
				output.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				output.DonTopology.WorkflowDonID,
				homeChainOut.ChainSelector,
				rpcs,
			)

			if settingsErr != nil {
				return settingsErr
			}

			// Copy the file to current directory as cre.settings.yaml
			currentDir, cErr := os.Getwd()
			if cErr != nil {
				return cErr
			}

			targetPath := filepath.Join(currentDir, "cre.settings.yaml")
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

		fmt.Println()
		fmt.Println("Environment started successfully")
		fmt.Println()
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

func startCLIEnvironment(topologyFlag string, extraAllowedPorts []int) (*creenv.SetupOutput, error) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[Config](nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load test configuration: %w", err)
	}

	capabilitiesBinaryPaths := map[cretypes.CapabilityFlag]string{}
	var capabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet

	if topologyFlag == TopologySimplified {
		if len(in.NodeSets) != 1 {
			return nil, fmt.Errorf("expected 1 nodeset, got %d", len(in.NodeSets))
		}
		// add support for more binaries if needed
		workflowDONCapabilities := []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WebAPITriggerCapability, cretypes.WriteEVMCapability, cretypes.WebAPITargetCapability}
		if in.ExtraCapabilities.CronCapabilityBinaryPath != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.CronCapability)
			capabilitiesBinaryPaths[cretypes.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
		}

		if in.ExtraCapabilities.LogEventTriggerBinaryPath != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cretypes.LogTriggerCapability)
			capabilitiesBinaryPaths[cretypes.LogTriggerCapability] = in.ExtraCapabilities.LogEventTriggerBinaryPath
		}

		if in.ExtraCapabilities.ReadContractBinaryPath != "" {
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

	// add support for more capabilities if needed
	capabilityFactoryFns := []cretypes.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPICapabilityFactoryFn,
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

	jobSpecFactoryFunctions := []cretypes.JobSpecFactoryFn{
		// add support for more job spec factory functions if needed
		webapi.WebAPIJobSpecFactoryFn,
		creconsensus.ConsensusJobSpecFactoryFn(libc.MustSafeUint64(int64(homeChainIDInt))),
		crecron.CronJobSpecFactoryFn(filepath.Join(containerPath, filepath.Base(in.ExtraCapabilities.CronCapabilityBinaryPath))),
		cregateway.GatewayJobSpecFactoryFn([]int{}, []string{}, []string{"0.0.0.0/0"}),
		crecompute.ComputeJobSpecFactoryFn,
	}

	for _, blockchain := range in.Blockchains {
		chainIDInt, chainErr := strconv.Atoi(blockchain.ChainID)
		if chainErr != nil {
			return nil, fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
		}
		capabilityFactoryFns = append(capabilityFactoryFns, chainwritercap.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))))
		capabilityFactoryFns = append(capabilityFactoryFns, chainreadercap.ChainReaderCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)), "evm"))

		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, chainreader.ChainReaderJobSpecFactoryFn(
			chainIDInt,
			"evm",
			// path within the container/pod
			filepath.Join(containerPath, filepath.Base(in.ExtraCapabilities.LogEventTriggerBinaryPath)),
			filepath.Join(containerPath, filepath.Base(in.ExtraCapabilities.ReadContractBinaryPath)),
		))
	}

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            capabilitiesAwareNodeSets,
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  capabilitiesBinaryPaths,
		JobSpecFactoryFunctions:              jobSpecFactoryFunctions,
		ConfigFactoryFunctions: []cretypes.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(context.Background(), testLogger, cldlogger.NewSingleFileLogger(nil), universalSetupInput)
	if setupErr != nil {
		return nil, fmt.Errorf("failed to setup test environment: %w", setupErr)
	}

	return universalSetupOutput, nil
}
