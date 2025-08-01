package environment

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
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

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/tracking"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	logeventtriggercap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/logevent"
	readcontractcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/readcontract"
	webapicap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapi"
	writeevmcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writeevm"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	crelogevent "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/logevent"
	crereadcontract "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/readcontract"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/webapi"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const manualCtfCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove containers with 'ctf' label and delete their volumes`
const manualBeholderCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove the 'chip-ingress' stack`

var (
	binDir string
)

func init() {
	EnvironmentCmd.AddCommand(startCmd())
	EnvironmentCmd.AddCommand(stopCmd)
	EnvironmentCmd.AddCommand(workflowCmds())
	EnvironmentCmd.AddCommand(beholderCmds())

	rootPath, rootPathErr := os.Getwd()
	if rootPathErr != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", rootPathErr)
		os.Exit(1)
	}
	binDir = filepath.Join(rootPath, "bin")
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		if err := os.Mkdir(binDir, 0755); err != nil {
			panic(fmt.Errorf("failed to create bin directory: %w", err))
		}
	}
}

func waitToCleanUp(d time.Duration) {
	fmt.Printf("Waiting %s before cleanup\n", d)
	time.Sleep(d)
}

var EnvironmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment commands",
	Long:  `Commands to manage the environment`,
}

const (
	TopologySimplified = "simplified"
	TopologyFull       = "full"

	WorkflowTriggerWebTrigger = "web-trigger"
	WorkflowTriggerCron       = "cron"
)

type Config struct {
	Blockchains       []*cre.WrappedBlockchainInput `toml:"blockchains" validate:"required"`
	NodeSets          []*ns.Input                   `toml:"nodesets" validate:"required"`
	JD                *jd.Input                     `toml:"jd" validate:"required"`
	Infra             *infra.Input                  `toml:"infra" validate:"required"`
	ExtraCapabilities ExtraCapabilitiesConfig       `toml:"extra_capabilities"`
	S3ProviderInput   *s3provider.Input             `toml:"s3provider"`
}

func (c Config) Validate() error {
	if c.JD.CSAEncryptionKey == "" {
		return errors.New("jd.csa_encryption_key must be provided")
	}
	return nil
}

type ExtraCapabilitiesConfig struct {
	CronCapabilityBinaryPath  string `toml:"cron_capability_binary_path"`
	LogEventTriggerBinaryPath string `toml:"log_event_trigger_binary_path"`
	ReadContractBinaryPath    string `toml:"read_contract_capability_binary_path"`
}

// DX tracking
var (
	dxTracker             tracking.Tracker
	provisioningStartTime time.Time
)

var StartCmdPreRunFunc = func(cmd *cobra.Command, args []string) {
	provisioningStartTime = time.Now()

	// ensure non-nil dxTracker by default
	var trackerErr error
	dxTracker, trackerErr = tracking.NewDxTracker()
	if trackerErr != nil {
		fmt.Fprintf(os.Stderr, "failed to create DX tracker: %s\n", trackerErr)
		dxTracker = &tracking.NoOpTracker{}
	}

	// remove all containers before starting the environment, just in case
	_ = framework.RemoveTestContainers()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("\nReceived signal: %s\n", sig)

		removeErr := framework.RemoveTestContainers()
		if removeErr != nil {
			fmt.Fprint(os.Stderr, removeErr, manualCtfCleanupMsg)
		}

		os.Exit(1)
	}()
}

var StartCmdRecoverHandlerFunc = func(p interface{}, cleanupWait time.Duration) {
	if p != nil {
		fmt.Println("Panicked when starting environment")

		var errText string
		if err, ok := p.(error); ok {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

			errText = strings.SplitN(err.Error(), "\n", 1)[0]
		} else {
			fmt.Fprintf(os.Stderr, "panic: %v\n", p)
			fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

			errText = strings.SplitN(fmt.Sprintf("%v", p), "\n", 1)[0]
		}

		tracingErr := dxTracker.Track("startup.result", map[string]any{
			"success":  false,
			"error":    errText,
			"panicked": true,
		})

		if tracingErr != nil {
			fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", tracingErr)
		}

		waitToCleanUp(cleanupWait)

		removeErr := framework.RemoveTestContainers()
		if removeErr != nil {
			fmt.Fprint(os.Stderr, errors.Wrap(removeErr, manualCtfCleanupMsg).Error())
		}
	}
}

var StartCmdGenerateSettingsFile = func(homeChainOut *cre.WrappedBlockchainOutput, output *creenv.SetupOutput) error {
	rpcs := map[uint64]string{}
	for _, bcOut := range output.BlockchainOutput {
		rpcs[bcOut.ChainSelector] = bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl
	}

	creCLISettingsFile, settingsErr := crecli.PrepareCRECLISettingsFile(
		crecli.CRECLIProfile,
		homeChainOut.SethClient.MustGetRootKeyAddress(),
		output.CldEnvironment.ExistingAddresses, //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
		output.DonTopology.WorkflowDonID,
		homeChainOut.ChainSelector,
		rpcs,
		output.S3ProviderOutput,
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

	fmt.Printf("CRE CLI settings file created: %s\n\n", targetPath)

	return nil
}

func startCmd() *cobra.Command {
	var (
		topology                 string
		extraAllowedGatewayPorts []int
		withExampleFlag          bool
		exampleWorkflowTrigger   string
		exampleWorkflowTimeout   time.Duration
		withPluginsDockerImage   string
		doSetup                  bool
		cleanupWait              time.Duration
		withBeholder             bool
		protoConfigs             []string
	)

	cmd := &cobra.Command{
		Use:              "start",
		Short:            "Start the environment",
		Long:             `Start the local CRE environment with all supported capabilities`,
		Aliases:          []string{"restart"},
		PersistentPreRun: StartCmdPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				p := recover()
				StartCmdRecoverHandlerFunc(p, cleanupWait)
			}()

			if doSetup {
				setupErr := RunSetup(cmd.Context(), SetupConfig{}, false, false)
				if setupErr != nil {
					return errors.Wrap(setupErr, "failed to run setup")
				}
			}

			if topology != TopologySimplified && topology != TopologyFull {
				return fmt.Errorf("invalid topology: %s. Valid topologies are: %s, %s", topology, TopologySimplified, TopologyFull)
			}

			PrintCRELogo()

			if err := defaultCtfConfigs(topology); err != nil {
				return errors.Wrap(err, "failed to set default CTF configs")
			}

			if os.Getenv("PRIVATE_KEY") == "" {
				setErr := os.Setenv("PRIVATE_KEY", blockchain.DefaultAnvilPrivateKey)
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
			// Load and validate test configuration
			in, err := framework.Load[Config](nil)
			if err != nil {
				return errors.Wrap(err, "failed to load test configuration")
			}
			if err := in.Validate(); err != nil {
				return errors.Wrap(err, "failed to validate test configuration")
			}

			output, startErr := StartCLIEnvironment(cmdContext, in, topology, exampleWorkflowTrigger, withPluginsDockerImage, withExampleFlag, extraAllowedGatewayPorts, nil, nil)
			if startErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", startErr)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				dxErr := trackStartup(false, hasBuiltDockerImage(in, withPluginsDockerImage), in.Infra.Type, ptr.Ptr(strings.SplitN(startErr.Error(), "\n", 1)[0]), ptr.Ptr(false))
				if dxErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", dxErr)
				}

				waitToCleanUp(cleanupWait)
				removeErr := framework.RemoveTestContainers()
				if removeErr != nil {
					return errors.Wrap(removeErr, manualCtfCleanupMsg)
				}

				return errors.Wrap(startErr, "failed to start environment")
			}

			homeChainOut := output.BlockchainOutput[0]

			sErr := StartCmdGenerateSettingsFile(homeChainOut, output)

			if sErr != nil {
				fmt.Fprintf(os.Stderr, "failed to create CRE CLI settings file: %s. You need to create it manually.", sErr)
			}

			dxErr := trackStartup(true, hasBuiltDockerImage(in, withPluginsDockerImage), output.InfraInput.Type, nil, nil)
			if dxErr != nil {
				fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", dxErr)
			}

			if withBeholder {
				startBeholderErr := startBeholder(
					cmdContext,
					cleanupWait,
					protoConfigs,
				)
				if startBeholderErr != nil {
					if !strings.Contains(startBeholderErr.Error(), protoRegistrationErrMsg) {
						beholderRemoveErr := framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
						if beholderRemoveErr != nil {
							fmt.Fprint(os.Stderr, errors.Wrap(beholderRemoveErr, manualBeholderCleanupMsg).Error())
						}
					}
					return errors.Wrap(startBeholderErr, "failed to start Beholder")
				}
			}

			if withExampleFlag {
				gatewayURL := fmt.Sprintf("%s://%s:%d%s", output.DonTopology.GatewayConnectorOutput.Incoming.Protocol, output.DonTopology.GatewayConnectorOutput.Incoming.Host, output.DonTopology.GatewayConnectorOutput.Incoming.ExternalPort, output.DonTopology.GatewayConnectorOutput.Incoming.Path)

				fmt.Print(libformat.PurpleText("\nRegistering and verifying example workflow\n\n"))

				wfRegAddr := libcontracts.MustFindAddressesForChain(
					output.CldEnvironment.ExistingAddresses, //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
					output.BlockchainOutput[0].ChainSelector,
					keystone_changeset.WorkflowRegistry.String())
				deployErr := deployAndVerifyExampleWorkflow(cmdContext, homeChainOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl, gatewayURL, exampleWorkflowTimeout, exampleWorkflowTrigger, wfRegAddr.Hex())
				if deployErr != nil {
					fmt.Printf("Failed to deploy and verify example workflow: %s\n", deployErr)
				}
			}
			fmt.Print(libformat.PurpleText("\nEnvironment setup completed successfully in %.2f seconds\n\n", time.Since(provisioningStartTime).Seconds()))
			fmt.Print("To terminate execute:`go run . env stop`\n\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&topology, "topology", "t", "simplified", "Topology to use for the environment (simplified or full)")
	cmd.Flags().DurationVarP(&cleanupWait, "wait-on-error-timeout", "w", 15*time.Second, "Wait on error timeout (e.g. 10s, 1m, 1h)")
	cmd.Flags().IntSliceVarP(&extraAllowedGatewayPorts, "extra-allowed-gateway-ports", "e", []int{}, "Extra allowed ports for outgoing connections from the Gateway DON (e.g. 8080,8081)")
	cmd.Flags().BoolVarP(&withExampleFlag, "with-example", "x", false, "Deploy and register example workflow")
	cmd.Flags().DurationVarP(&exampleWorkflowTimeout, "example-workflow-timeout", "u", 5*time.Minute, "Time to wait until example workflow succeeds")
	cmd.Flags().StringVarP(&withPluginsDockerImage, "with-plugins-docker-image", "p", "", "Docker image to use (must have all capabilities included)")
	cmd.Flags().StringVarP(&exampleWorkflowTrigger, "example-workflow-trigger", "y", "web-trigger", "Trigger for example workflow to deploy (web-trigger or cron)")
	cmd.Flags().BoolVarP(&withBeholder, "with-beholder", "b", false, "Deploy Beholder (Chip Ingress + Red Panda)")
	cmd.Flags().StringArrayVarP(&protoConfigs, "with-proto-configs", "c", []string{"./proto-configs/default.toml"}, "Protos configs to use (e.g. './proto-configs/config_one.toml,./proto-configs/config_two.toml')")
	cmd.Flags().BoolVarP(&doSetup, "auto-setup", "a", false, "Run setup before starting the environment")
	return cmd
}

func trackStartup(success, hasBuiltDockerImage bool, infraType string, errorMessage *string, panicked *bool) error {
	metadata := map[string]any{
		"success": success,
		"infra":   infraType,
	}

	if errorMessage != nil {
		metadata["error"] = *errorMessage
	}

	if panicked != nil {
		metadata["panicked"] = *panicked
	}

	dxStartupErr := dxTracker.Track("cre.local.startup.result", metadata)
	if dxStartupErr != nil {
		fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", dxStartupErr)
	}

	if success {
		dxTimeErr := dxTracker.Track("cre.local.startup.time", map[string]any{
			"duration_seconds":       time.Since(provisioningStartTime).Seconds(),
			"has_built_docker_image": hasBuiltDockerImage,
		})

		if dxTimeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to track startup time: %s\n", dxTimeErr)
		}
	}

	return nil
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the environment",
	Long:  `Stops the local CRE environment (if it's not running, it just fallsthrough)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		removeErr := framework.RemoveTestContainers()
		if removeErr != nil {
			return errors.Wrap(removeErr, "failed to remove environment containers. Please remove them manually")
		}

		fmt.Println("Environment stopped successfully")
		return nil
	},
}

func StartCLIEnvironment(
	cmdContext context.Context,
	in *Config,
	topologyFlag string,
	workflowTrigger,
	withPluginsDockerImageFlag string,
	withExampleFlag bool,
	extraAllowedGatewayPorts []int,
	extraBinaries map[string]string,
	extraJobFactoryFns []cre.JobSpecFactoryFn,
) (*creenv.SetupOutput, error) {
	testLogger := framework.L

	// make sure that either cron is enabled or withPluginsDockerImageFlag is set, but only if workflowTrigger is cron
	if withExampleFlag && workflowTrigger == WorkflowTriggerCron && (in.ExtraCapabilities.CronCapabilityBinaryPath == "" && withPluginsDockerImageFlag == "") {
		return nil, fmt.Errorf("either cron binary path must be set in TOML config (%s) or you must use Docker image with all capabilities included and passed via withPluginsDockerImageFlag", os.Getenv("CTF_CONFIGS"))
	}

	capabilitiesBinaryPaths := map[cre.CapabilityFlag]string{}
	var capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet

	if topologyFlag == TopologySimplified {
		if len(in.NodeSets) != 1 {
			return nil, fmt.Errorf("expected 1 nodeset, got %d", len(in.NodeSets))
		}
		// add support for more binaries if needed
		workflowDONCapabilities := []string{cre.OCR3Capability, cre.CustomComputeCapability, cre.WriteEVMCapability, cre.WebAPITriggerCapability, cre.WebAPITargetCapability}
		if in.ExtraCapabilities.CronCapabilityBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cre.CronCapability)
			capabilitiesBinaryPaths[cre.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
		}

		if in.ExtraCapabilities.LogEventTriggerBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cre.LogTriggerCapability)
			capabilitiesBinaryPaths[cre.LogTriggerCapability] = in.ExtraCapabilities.LogEventTriggerBinaryPath
		}

		if in.ExtraCapabilities.ReadContractBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cre.ReadContractCapability)
			capabilitiesBinaryPaths[cre.ReadContractCapability] = in.ExtraCapabilities.ReadContractBinaryPath
		}

		for capabilityName, binaryPath := range extraBinaries {
			if binaryPath != "" || withPluginsDockerImageFlag != "" {
				workflowDONCapabilities = append(workflowDONCapabilities, capabilityName)
				capabilitiesBinaryPaths[capabilityName] = binaryPath
			}
		}

		capabilitiesAwareNodeSets = []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],
				Capabilities:       workflowDONCapabilities,
				DONTypes:           []string{cre.WorkflowDON, cre.GatewayDON},
				BootstrapNodeIndex: 0,
				GatewayNodeIndex:   0,
			},
		}
	} else {
		if len(in.NodeSets) != 3 {
			return nil, fmt.Errorf("expected 3 nodesets, got %d", len(in.NodeSets))
		}

		// add support for more binaries if needed
		workflowDONCapabilities := []string{cre.OCR3Capability, cre.CustomComputeCapability, cre.WebAPITriggerCapability}
		if in.ExtraCapabilities.CronCapabilityBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cre.CronCapability)
			capabilitiesBinaryPaths[cre.CronCapability] = in.ExtraCapabilities.CronCapabilityBinaryPath
		}

		if in.ExtraCapabilities.LogEventTriggerBinaryPath != "" || withPluginsDockerImageFlag != "" {
			workflowDONCapabilities = append(workflowDONCapabilities, cre.LogTriggerCapability)
			capabilitiesBinaryPaths[cre.LogTriggerCapability] = in.ExtraCapabilities.LogEventTriggerBinaryPath
		}

		for capabilityName, binaryPath := range extraBinaries {
			if binaryPath != "" || withPluginsDockerImageFlag != "" {
				workflowDONCapabilities = append(workflowDONCapabilities, capabilityName)
				capabilitiesBinaryPaths[capabilityName] = binaryPath
			}
		}

		capabiliitesDONCapabilities := []string{cre.WriteEVMCapability, cre.WebAPITargetCapability}
		if in.ExtraCapabilities.ReadContractBinaryPath != "" || withPluginsDockerImageFlag != "" {
			capabiliitesDONCapabilities = append(capabiliitesDONCapabilities, cre.ReadContractCapability)
			capabilitiesBinaryPaths[cre.ReadContractCapability] = in.ExtraCapabilities.ReadContractBinaryPath
		}

		capabilitiesAwareNodeSets = []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],
				Capabilities:       workflowDONCapabilities,
				DONTypes:           []string{cre.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              in.NodeSets[1],
				Capabilities:       capabiliitesDONCapabilities,
				DONTypes:           []string{cre.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                            // <----- it's crucial to indicate there's no bootstrap node
			},
			{
				Input:              in.NodeSets[2],
				Capabilities:       []string{},
				DONTypes:           []string{cre.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                       // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
	}

	// unset DockerFilePath and DockerContext as we cannot use them with existing images
	if withPluginsDockerImageFlag != "" {
		for setIdx := range capabilitiesAwareNodeSets {
			for nodeIdx := range capabilitiesAwareNodeSets[setIdx].NodeSpecs {
				capabilitiesAwareNodeSets[setIdx].NodeSpecs[nodeIdx].Node.Image = withPluginsDockerImageFlag
				capabilitiesAwareNodeSets[setIdx].NodeSpecs[nodeIdx].Node.DockerContext = ""
				capabilitiesAwareNodeSets[setIdx].NodeSpecs[nodeIdx].Node.DockerFilePath = ""
			}
		}
	}

	fmt.Print(libformat.PurpleText("DON topology:\n"))
	for _, nodeSet := range capabilitiesAwareNodeSets {
		fmt.Print(libformat.PurpleText("%s\n", strings.ToUpper(nodeSet.Name)))
		fmt.Print(libformat.PurpleText("\tNode count: %d\n", len(nodeSet.NodeSpecs)))
		capabilitiesDesc := "none"
		if len(nodeSet.Capabilities) > 0 {
			capabilitiesDesc = strings.Join(nodeSet.Capabilities, ", ")
		}
		fmt.Print(libformat.PurpleText("\tCapabilities: %s\n", capabilitiesDesc))
		fmt.Print(libformat.PurpleText("\tDON Types: %s\n\n", strings.Join(nodeSet.DONTypes, ", ")))
	}

	// add support for more capabilities if needed
	capabilityFactoryFns := []cre.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPITriggerCapabilityFactoryFn,
		webapicap.WebAPITargetCapabilityFactoryFn,
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
	}

	containerPath, pathErr := crecapabilities.DefaultContainerDirectory(in.Infra.Type)
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

	jobSpecFactoryFunctions := []cre.JobSpecFactoryFn{
		// add support for more job spec factory functions if needed
		webapi.WebAPITriggerJobSpecFactoryFn,
		webapi.WebAPITargetJobSpecFactoryFn,
		creconsensus.ConsensusJobSpecFactoryFn(libc.MustSafeUint64(int64(homeChainIDInt))),
		crecron.CronJobSpecFactoryFn(filepath.Join(containerPath, cronBinaryName)),
		cregateway.GatewayJobSpecFactoryFn(extraAllowedGatewayPorts, []string{}, []string{"0.0.0.0/0"}),
		crecompute.ComputeJobSpecFactoryFn,
	}

	jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, extraJobFactoryFns...)

	for _, blockchain := range in.Blockchains {
		chainIDInt, chainErr := strconv.Atoi(blockchain.ChainID)
		if chainErr != nil {
			return nil, fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
		}

		if !blockchain.ReadOnly {
			capabilityFactoryFns = append(capabilityFactoryFns, writeevmcap.WriteEVMCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))))
		}
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

	if in.JD.CSAEncryptionKey == "" {
		// generate a new key
		key, keyErr := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if keyErr != nil {
			return nil, fmt.Errorf("failed to generate CSA encryption key: %w", keyErr)
		}
		in.JD.CSAEncryptionKey = hex.EncodeToString(crypto.FromECDSA(key)[:32])
		fmt.Printf("Generated new CSA encryption key for JD: %s\n", in.JD.CSAEncryptionKey)
	}
	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            capabilitiesAwareNodeSets,
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		JobSpecFactoryFunctions:              jobSpecFactoryFunctions,
		ConfigFactoryFunctions: []cre.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
		S3ProviderInput: in.S3ProviderInput,
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

func PrintCRELogo() {
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

func defaultCtfConfigs(topologyFlag string) error {
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

	return nil
}

func hasBuiltDockerImage(in *Config, withPluginsDockerImageFlag string) bool {
	if withPluginsDockerImageFlag != "" {
		return false
	}

	hasBuilt := false

	for _, nodeset := range in.NodeSets {
		for _, nodeSpec := range nodeset.NodeSpecs {
			if nodeSpec.Node != nil && nodeSpec.Node.DockerFilePath != "" {
				hasBuilt = true
				break
			}
		}
	}

	return hasBuilt
}

func oneLineErrorMessage(errOrPanic any) string {
	if err, ok := errOrPanic.(error); ok {
		return strings.SplitN(err.Error(), "\n", 1)[0]
	}

	return strings.SplitN(fmt.Sprintf("%v", errOrPanic), "\n", 1)[0]
}
