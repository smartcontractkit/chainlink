package environment

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
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
}

// update e2e-test.ysml with new command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the environment",
	Long:  `Start the local CRE environment with all supported capabilities`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("CTF_CONFIGS") == "" {
			return fmt.Errorf("CTF_CONFIGS environment variable is not set. It should contain paths to TOML files with configurations")
		}
		fmt.Println("Starting the environment...")

		// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
		setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if setErr != nil {
			return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
		}

		output, err := startCLIEnvironment()
		if err != nil {
			return err
		}

		//TODO print urls?
		_ = output

		fmt.Println("Environment started successfully")
		fmt.Println("To terminate execute: ctf d rm")

		return nil
	},
}

type EnvironmentConfig struct {
	Blockchain        *blockchain.Input       `toml:"blockchain" validate:"required"`
	NodeSets          []*ns.Input             `toml:"nodesets" validate:"required"`
	JD                *jd.Input               `toml:"jd" validate:"required"`
	Infra             *libtypes.InfraInput    `toml:"infra" validate:"required"`
	ExtraCapabilities ExtraCapabilitiesConfig `toml:"extra_capabilities" validate:"required"`
}

type ExtraCapabilitiesConfig struct {
	CronCapabilityBinaryPath  string `toml:"cron_capability_binary_path"`
	LogEventTriggerBinaryPath string `toml:"log_event_trigger_binary_path"`
	ReadContractBinaryPath    string `toml:"read_contract_capability_binary_path"`
}

func startCLIEnvironment() (*creenv.SetupOutput, error) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[EnvironmentConfig](nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load test configuration: %w", err)
	}

	if len(in.NodeSets) != 3 {
		return nil, fmt.Errorf("expected 3 nodesets, got %d", len(in.NodeSets))
	}

	capabilitiesBinaryPaths := map[cretypes.CapabilityFlag]string{}

	workflowDONCapabilities := []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WebAPITriggerCapability}
	if in.ExtraCapabilities.CronCapabilityBinaryPath != "" {
		workflowDONCapabilities = append(workflowDONCapabilities, cretypes.CronCapability)
		capabilitiesBinaryPaths[cretypes.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
	}

	capabiliitesDONCapabilities := []string{cretypes.WriteEVMCapability, cretypes.ReadContractCapability, cretypes.WebAPITargetCapability}
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
			GatewayNodeIndex:   0,
			BootstrapNodeIndex: -1, // <----- it's crucial to indicate there's no bootstrap node
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

	capabilityFactoryFns := []cretypes.DONCapabilityWithConfigFactoryFn{
		crecontracts.DefaultCapabilityFactoryFn,
		crecontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))),
		crecontracts.ChainReaderCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)), "evm"), // for now support only evm
		crecontracts.WebAPICapabilityFactoryFn,
	}

	chainReaderJobSpecFactoryFn := chainreader.ChainReaderJobSpecFactoryFn(
		chainIDInt,
		"evm",
		in.ExtraCapabilities.LogEventTriggerBinaryPath,
		in.ExtraCapabilities.ReadContractBinaryPath,
	)

	porJobSpecFactoryFn := crepor.PoRJobSpecFactoryFn(
		in.ExtraCapabilities.CronCapabilityBinaryPath,
		[]int{},
		[]string{},
		[]string{"0.0.0.0/0"}, // allow all IPs
	)

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:  capabilitiesAwareNodeSets,
		CapabilityFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:           *in.Blockchain,
		JdInput:                    *in.JD,
		InfraInput:                 *in.Infra,
		CustomBinariesPaths:        capabilitiesBinaryPaths,
		JobSpecFactoryFunctions:    []cretypes.JobSpecFactoryFn{chainReaderJobSpecFactoryFn, webapi.WebAPIJobSpecFactoryFn, porJobSpecFactoryFn},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testLogger, cldlogger.NewSingleFileLogger(nil), context.Background(), universalSetupInput)
	if setupErr != nil {
		return nil, fmt.Errorf("failed to setup test environment: %w", setupErr)
	}

	return universalSetupOutput, nil
}
