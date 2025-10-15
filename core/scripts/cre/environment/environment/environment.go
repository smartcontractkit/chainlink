package environment

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
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
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/sets"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	gateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	billingplatformservice "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/billing_platform_service"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/tracking"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const (
	manualCtfCleanupMsg      = `unexpected startup error. this may have stranded resources. please manually remove containers with 'ctf' label and delete their volumes`
	manualBeholderCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove the 'chip-ingress' stack`
	manualBillingCleanupMsg  = `unexpected startup error. this may have stranded resources. please manually remove the 'billing-platform-service' stack`
)

var (
	binDir string

	defaultCapabilitiesConfigFile = "configs/capability_defaults.toml"
)

// DX tracking
var (
	dxTracker             tracking.Tracker
	provisioningStartTime time.Time
)

const (
	TopologyWorkflow                    = "workflow"
	TopologyWorkflowGateway             = "workflow-gateway"
	TopologyWorkflowGatewayCapabilities = "workflow-gateway-capabilities"
	TopologyMock                        = "mock"

	WorkflowTriggerWebTrigger = "web-trigger"
	WorkflowTriggerCron       = "cron"
)

var EnvironmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment commands",
	Long:  `Commands to manage the environment`,
}

func init() {
	EnvironmentCmd.AddCommand(startCmd())
	EnvironmentCmd.AddCommand(stopCmd())
	EnvironmentCmd.AddCommand(workflowCmds())
	EnvironmentCmd.AddCommand(beholderCmds())
	EnvironmentCmd.AddCommand(swapCmds())
	EnvironmentCmd.AddCommand(stateCmd())
	EnvironmentCmd.AddCommand(billingCmds())

	rootPath, rootPathErr := os.Getwd()
	if rootPathErr != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", rootPathErr)
		os.Exit(1)
	}
	binDir = filepath.Join(rootPath, "bin")
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		if err := os.Mkdir(binDir, 0o755); err != nil {
			panic(fmt.Errorf("failed to create bin directory: %w", err))
		}
	}
}

func waitToCleanUp(d time.Duration) {
	fmt.Printf("Waiting %s before cleanup\n", d)
	time.Sleep(d)
}

var StartCmdPreRunFunc = func(cmd *cobra.Command, args []string) {
	globalPreRunFunc(cmd, args)
	provisioningStartTime = time.Now()

	// ensure non-nil dxTracker by default
	initDxTracker()

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

var StartCmdRecoverHandlerFunc = func(p any, cleanupWait time.Duration) {
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

		tracingErr := dxTracker.Track(MetricStartupResult, map[string]any{
			"success":  false,
			"error":    errText,
			"panicked": true,
		})

		if tracingErr != nil {
			fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", tracingErr)
		}

		waitToCleanUp(cleanupWait)
		_, saveErr := framework.SaveContainerLogs("./logs")
		if saveErr != nil {
			fmt.Fprintf(os.Stderr, "failed to save container logs: %s\n", saveErr)
		}

		removeErr := framework.RemoveTestContainers()
		if removeErr != nil {
			fmt.Fprint(os.Stderr, errors.Wrap(removeErr, manualCtfCleanupMsg).Error())
		}

		// signal that the environment failed to start
		os.Exit(1)
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
	err = os.WriteFile(targetPath, input, 0o600)
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
		withContractsVersion     string
		doSetup                  bool
		cleanupWait              time.Duration
		withBeholder             bool
		withDashboards           bool
		withBilling              bool
		protoConfigs             []string
		setupConfig              SetupConfig
	)

	cmd := &cobra.Command{
		Use:              "start",
		Short:            "Start the environment",
		Long:             `Start the local CRE environment with all supported capabilities`,
		Aliases:          []string{"restart"},
		PersistentPreRun: StartCmdPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				StartCmdRecoverHandlerFunc(recover(), cleanupWait)
			}()

			if doSetup {
				setupErr := RunSetup(cmd.Context(), SetupConfig{ConfigPath: DefaultSetupConfigPath}, true, false, withBilling)
				if setupErr != nil {
					return errors.Wrap(setupErr, "failed to run setup")
				}
			}

			if topology != TopologyWorkflow && topology != TopologyWorkflowGatewayCapabilities && topology != TopologyWorkflowGateway && topology != TopologyMock {
				framework.L.Warn().Msgf("'%s' is an unknown topology. Using whatever configuration was passed in CTF_CONFIGs", topology)
			}

			PrintCRELogo()

			if err := defaultCtfConfigs(topology); err != nil {
				return errors.Wrap(err, "failed to set default CTF configs")
			}

			if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
				return errors.Wrap(pkErr, "failed to set default private key")
			}

			cleanUpErr := envconfig.RemoveAllEnvironmentStateDir(relativePathToRepoRoot)
			if cleanUpErr != nil {
				return errors.Wrap(cleanUpErr, "failed to clean up environment state files")
			}

			// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
			setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			if setErr != nil {
				return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
			}

			cmdContext := cmd.Context()
			in := &envconfig.Config{}

			if err := in.Load(os.Getenv("CTF_CONFIGS")); err != nil {
				return errors.Wrap(err, "failed to load environment configuration")
			}

			if err := ensureDockerIsRunning(cmdContext); err != nil {
				return err
			}

			// This will not work with remote images that require authentication, but it will catch early most of the issues with missing env setup
			if err := ensureDockerImagesExist(cmdContext, framework.L, in, withPluginsDockerImage); err != nil {
				return err
			}

			withV2Registries := withContractsVersion == "v2"
			envDependencies := cre.NewEnvironmentDependencies(
				flags.NewDefaultCapabilityFlagsProvider(),
				cre.NewContractVersionsProvider(envconfig.DefaultContractSet(withV2Registries)),
				cre.NewCLIFlagsProvider(withV2Registries),
			)

			if err := in.Validate(envDependencies); err != nil {
				return errors.Wrap(err, "failed to validate test configuration")
			}

			homeChainIDInt, chainErr := strconv.Atoi(in.Blockchains[0].ChainID)
			if chainErr != nil {
				return fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
			}

			defaultCapabilities, defaultCapabilitiesErr := sets.NewDefaultSet(libc.MustSafeUint64FromInt(homeChainIDInt))
			if defaultCapabilitiesErr != nil {
				return errors.Wrap(defaultCapabilitiesErr, "failed to create default capabilities")
			}

			if err := validateWorkflowTriggerAndCapabilities(in, withExampleFlag, exampleWorkflowTrigger, withPluginsDockerImage); err != nil {
				return errors.Wrap(err, "either cron binary path must be set in TOML config (%s) or you must use Docker image with all capabilities included and passed via withPluginsDockerImageFlag")
			}

			extraJobSpecFunctions := []cre.JobSpecFn{
				// temporary solution until we figure out where that jobspec should live. Gateway is not a capability, it's more of a role
				// but we don't have a good expression of that abstraction yet
				gateway.JobSpec(append(extraAllowedGatewayPorts, in.Fake.Port), []string{}, []string{"0.0.0.0/0"}),
			}

			output, startErr := StartCLIEnvironment(cmdContext, relativePathToRepoRoot, in, topology, withPluginsDockerImage, defaultCapabilities, extraJobSpecFunctions, envDependencies)
			if startErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", startErr)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				dxErr := trackStartup(false, hasBuiltDockerImage(in, withPluginsDockerImage), in.Infra.Type, ptr.Ptr(strings.SplitN(startErr.Error(), "\n", 1)[0]), ptr.Ptr(false))
				if dxErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", dxErr)
				}

				waitToCleanUp(cleanupWait)
				_, saveErr := framework.SaveContainerLogs("./logs")
				if saveErr != nil {
					fmt.Fprintf(os.Stderr, "failed to save container logs: %s\n", saveErr)
				}

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

				metaData := map[string]any{}
				if startBeholderErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(startBeholderErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track(MetricBeholderStart, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track beholder start: %s\n", trackingErr)
				}

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

			if withDashboards {
				err := setupDashboards(setupConfig)
				if err != nil {
					return errors.Wrap(err, "failed to setup dashboards")
				}
			}

			if withBilling {
				startBillingErr := startBilling(
					cmdContext,
					cleanupWait,
					output,
				)

				metaData := map[string]any{}
				if startBillingErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(startBillingErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track(MetricBillingStart, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track billing start: %s\n", trackingErr)
				}

				if startBillingErr != nil {
					if !strings.Contains(startBillingErr.Error(), protoRegistrationErrMsg) {
						billingRemoveErr := framework.RemoveTestStack(billingplatformservice.DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
						if billingRemoveErr != nil {
							fmt.Fprint(os.Stderr, errors.Wrap(billingRemoveErr, manualBillingCleanupMsg).Error())
						}
					}

					return errors.Wrap(startBillingErr, "failed to start Billing Platform Service")
				}
			}

			if withExampleFlag {
				if output.DonTopology.GatewayConnectorOutput == nil || len(output.DonTopology.GatewayConnectorOutput.Configurations) == 0 {
					return errors.New("no gateway connector configurations found")
				}

				// use first gateway for example workflow
				gatewayURL := fmt.Sprintf("%s://%s:%d%s", output.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol, output.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host, output.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort, output.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)

				fmt.Print(libformat.PurpleText("\nRegistering and verifying example workflow\n\n"))

				wfRegAddr := libcontracts.MustFindAddressesForChain(
					output.CldEnvironment.ExistingAddresses, //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
					output.BlockchainOutput[0].ChainSelector,
					keystone_changeset.WorkflowRegistry.String())

				var workflowDonID uint32
				for idx, don := range output.DonTopology.Dons.List() {
					if don.HasFlag(cre.WorkflowDON) {
						workflowDonID = libc.MustSafeUint32(idx + 1)
						break
					}
				}

				if workflowDonID == 0 {
					return errors.New("no workflow DON found")
				}

				deployErr := deployAndVerifyExampleWorkflow(cmdContext, homeChainOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl, gatewayURL, output.DonTopology.GatewayConnectorOutput.Configurations[0].Dons[0].ID, workflowDonID, exampleWorkflowTimeout, exampleWorkflowTrigger, wfRegAddr.Hex())
				if deployErr != nil {
					fmt.Printf("Failed to deploy and verify example workflow: %s\n", deployErr)
				}
			}
			fmt.Print(libformat.PurpleText("\nEnvironment setup completed successfully in %.2f seconds\n\n", time.Since(provisioningStartTime).Seconds()))
			fmt.Print("To terminate execute:`go run . env stop`\n\n")

			storeErr := in.Store(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
			if storeErr != nil {
				return errors.Wrap(storeErr, "failed to store local CRE state")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&topology, "topology", "t", TopologyWorkflow, "Topology to use for the environment (workflow, workflow-gateway, workflow-gateway-capabilities)")
	cmd.Flags().DurationVarP(&cleanupWait, "wait-on-error-timeout", "w", 15*time.Second, "Wait on error timeout (e.g. 10s, 1m, 1h)")
	cmd.Flags().IntSliceVarP(&extraAllowedGatewayPorts, "extra-allowed-gateway-ports", "e", []int{}, "Extra allowed ports for outgoing connections from the Gateway DON (e.g. 8080,8081)")
	cmd.Flags().BoolVarP(&withExampleFlag, "with-example", "x", false, "Deploy and register example workflow")
	cmd.Flags().DurationVarP(&exampleWorkflowTimeout, "example-workflow-timeout", "u", 5*time.Minute, "Time to wait until example workflow succeeds")
	cmd.Flags().StringVarP(&withPluginsDockerImage, "with-plugins-docker-image", "p", "", "Docker image to use (must have all capabilities included)")
	cmd.Flags().StringVarP(&exampleWorkflowTrigger, "example-workflow-trigger", "y", "web-trigger", "Trigger for example workflow to deploy (web-trigger or cron)")
	cmd.Flags().BoolVarP(&withBeholder, "with-beholder", "b", false, "Deploy Beholder (Chip Ingress + Red Panda)")
	cmd.Flags().BoolVarP(&withDashboards, "with-dashboards", "d", false, "Deploy Observability Stack and Grafana Dashboards")
	cmd.Flags().BoolVar(&withBilling, "with-billing", false, "Deploy Billing Platform Service")
	cmd.Flags().StringArrayVarP(&protoConfigs, "with-proto-configs", "c", []string{"./proto-configs/default.toml"}, "Protos configs to use (e.g. './proto-configs/config_one.toml,./proto-configs/config_two.toml')")
	cmd.Flags().BoolVarP(&doSetup, "auto-setup", "a", false, "Run setup before starting the environment")
	cmd.Flags().StringVar(&withContractsVersion, "with-contracts-version", "v1", "Version of workflow and capabilities registry contracts to use (v1 or v2)")
	cmd.Flags().StringVarP(&setupConfig.ConfigPath, "config", "s", DefaultSetupConfigPath, "Path to the TOML configuration file")
	return cmd
}

func setupDashboards(setupCfg SetupConfig) error {
	cfg, cfgErr := readConfig(setupCfg.ConfigPath)
	if cfgErr != nil {
		return errors.Wrap(cfgErr, "failed to read config")
	}

	// Run the `ctf obs up -f` command from the ./bin directory
	ctfCmd := exec.Command("./bin/ctf", "obs", "up", "-f")

	obsOutput, err := ctfCmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return errors.Wrap(err, "failed to start ctf observability stack: "+string(obsOutput))
	}

	fmt.Print(libformat.PurpleText("\nObservabilty stack setup completed successfully\n"))

	// Wait for grafana at localhost:3000 to be available
	fmt.Print(libformat.PurpleText("\nWaiting for Grafana to be available at http://localhost:3000\n"))
	grafanaContacted := false
	for range 30 {
		time.Sleep(1 * time.Second)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:3000", nil)
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		grafanaContacted = true
		break
	}

	if !grafanaContacted {
		return errors.New("timed out waiting for Grafana to be available at http://localhost:3000")
	}

	// Check the file exists before trying to run the script
	scriptPath := filepath.Join(cfg.Observability.TargetPath, "deploy-cre-local.sh")
	if _, err = os.Stat(scriptPath); os.IsNotExist(err) {
		return errors.New("deploy-cre-local.sh script does not exist, ensure the setup command has been run")
	}

	deployDashboardsCmd := exec.Command("./deploy-cre-local.sh")
	deployDashboardsCmd.Dir = cfg.Observability.TargetPath
	deployOutput, err := deployDashboardsCmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return errors.Wrap(err, "failed to deploy dashboards: "+string(deployOutput))
	}

	fmt.Print(libformat.PurpleText("\nDashboards successfully deployed\n"))
	return nil
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

	dxStartupErr := dxTracker.Track(MetricStartupResult, metadata)
	if dxStartupErr != nil {
		fmt.Fprintf(os.Stderr, "failed to track startup: %s\n", dxStartupErr)
	}

	if success {
		dxTimeErr := dxTracker.Track(MetricStartupTime, map[string]any{
			"duration_seconds":       time.Since(provisioningStartTime).Seconds(),
			"has_built_docker_image": hasBuiltDockerImage,
		})

		if dxTimeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to track startup time: %s\n", dxTimeErr)
		}
	}

	return nil
}

func stopCmd() *cobra.Command {
	var allFlag bool
	cmd := &cobra.Command{
		Use:              "stop",
		Short:            "Stops the environment",
		Long:             `Stops the local CRE environment (if it's not running, it just fallsthrough)`,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			removeErr := framework.RemoveTestContainers()
			if removeErr != nil {
				return errors.Wrap(removeErr, "failed to remove environment containers. Please remove them manually")
			}

			if allFlag {
				stopBeholderErr := stopBeholder()
				if stopBeholderErr != nil {
					return errors.Wrap(stopBeholderErr, "failed to stop beholder")
				}

				stopBillingErr := stopBilling()
				if stopBillingErr != nil {
					return errors.Wrap(stopBillingErr, "failed to stop billing")
				}

				removeCacheErr := envconfig.RemoveAllEnvironmentStateDir(relativePathToRepoRoot)
				if removeCacheErr != nil {
					framework.L.Warn().Msgf("failed to remove local CRE state files: %s", removeCacheErr)
				}
			} else {
				creStateFile := envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)
				cErr := os.Remove(creStateFile)
				if cErr != nil {
					framework.L.Warn().Msgf("failed to remove local CRE state file: %s", cErr)
				} else {
					framework.L.Info().Msgf("removed local CRE state file: %s", creStateFile)
				}

				envArtifactFile := creenv.MustEnvArtifactAbsPath(relativePathToRepoRoot)
				eErr := os.Remove(envArtifactFile)
				if eErr != nil {
					framework.L.Warn().Msgf("failed to remove local CRE environment artifact file: %s", eErr)
				} else {
					framework.L.Info().Msgf("removed local CRE environment artifact file: %s", envArtifactFile)
				}
			}

			fmt.Println("Environment stopped successfully")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Remove all environment state files")

	return cmd
}

func StartCLIEnvironment(
	cmdContext context.Context,
	relativePathToRepoRoot string,
	in *envconfig.Config,
	topologyFlag string,
	withPluginsDockerImageFlag string,
	capabilities []cre.InstallableCapability,
	extraJobSpecFunctions []cre.JobSpecFn,
	env cre.CLIEnvironmentDependencies,
) (*creenv.SetupOutput, error) {
	testLogger := framework.L

	// unset DockerFilePath and DockerContext as we cannot use them with existing images
	if withPluginsDockerImageFlag != "" {
		for setIdx := range in.NodeSets {
			for nodeIdx := range in.NodeSets[setIdx].NodeSpecs {
				in.NodeSets[setIdx].NodeSpecs[nodeIdx].Node.Image = withPluginsDockerImageFlag
				in.NodeSets[setIdx].NodeSpecs[nodeIdx].Node.DockerContext = ""
				in.NodeSets[setIdx].NodeSpecs[nodeIdx].Node.DockerFilePath = ""
			}
		}
	}

	fmt.Print(libformat.PurpleText("DON topology:\n"))
	for _, nodeSet := range in.NodeSets {
		fmt.Print(libformat.PurpleText("%s\n", strings.ToUpper(nodeSet.Name)))
		fmt.Print(libformat.PurpleText("\tNode count: %d\n", len(nodeSet.NodeSpecs)))
		capabilitiesDesc := "none"
		if len(nodeSet.Capabilities) > 0 {
			capabilitiesDesc = strings.Join(nodeSet.Capabilities, ", ")
		}
		fmt.Print(libformat.PurpleText("\tGlobal capabilities: %s\n", capabilitiesDesc))
		chainCapabilitiesDesc := "none"
		if len(nodeSet.ChainCapabilities) > 0 {
			chainCapList := []string{}
			for capabilityName, chainCapability := range nodeSet.ChainCapabilities {
				for _, chainID := range chainCapability.EnabledChains {
					chainCapList = append(chainCapList, fmt.Sprintf("%s-%d", capabilityName, chainID))
				}
			}
			chainCapabilitiesDesc = strings.Join(chainCapList, ", ")
		}
		fmt.Print(libformat.PurpleText("\tChain capabilities: %s\n", chainCapabilitiesDesc))
		fmt.Print(libformat.PurpleText("\tDON Types: %s\n\n", strings.Join(nodeSet.DONTypes, ", ")))
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
	universalSetupInput := &creenv.SetupInput{
		CapabilitiesAwareNodeSets: in.NodeSets,
		BlockchainsInput:          in.Blockchains,
		ContractVersions:          env.ContractVersions(),
		WithV2Registries:          env.WithV2Registries(),
		JdInput:                   in.JD,
		Provider:                  *in.Infra,
		S3ProviderInput:           in.S3ProviderInput,
		CapabilityConfigs:         in.CapabilityConfigs,
		CopyCapabilityBinaries:    withPluginsDockerImageFlag == "", // do not copy any binaries to the containers, if we are using plugins image (they already have them)
		Capabilities:              capabilities,
		JobSpecFactoryFunctions:   extraJobSpecFunctions,
		StageGen:                  initLocalCREStageGen(in),
	}

	ctx, cancel := context.WithTimeout(cmdContext, 10*time.Minute)
	defer cancel()
	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(ctx, testLogger, cldlogger.NewSingleFileLogger(nil), universalSetupInput, relativePathToRepoRoot)
	if setupErr != nil {
		return nil, fmt.Errorf("failed to setup test environment: %w", setupErr)
	}

	capabilitiesContractFactoryFunctions := []cre.CapabilityRegistryConfigFn{}
	for _, cap := range capabilities {
		capabilitiesContractFactoryFunctions = append(capabilitiesContractFactoryFunctions, cap.CapabilityRegistryV1ConfigFn())
	}

	artifactPath, artifactErr := creenv.DumpArtifact(
		creenv.MustEnvArtifactAbsPath(relativePathToRepoRoot),
		universalSetupOutput.CldEnvironment.DataStore.Addresses(),
		universalSetupOutput.CldEnvironment.ExistingAddresses,
		*in.JD.Out,
		*universalSetupOutput.DonTopology,
		universalSetupOutput.CldEnvironment.Offchain,
		capabilitiesContractFactoryFunctions,
		in.NodeSets,
	)

	if artifactErr != nil {
		testLogger.Error().Err(artifactErr).Msg("failed to generate env artifact")
	} else {
		testLogger.Info().Msgf("Environment artifact saved to %s", artifactPath)
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
		var setErr error
		// use default configs for each
		switch topologyFlag {
		case TopologyWorkflow:
			setErr = os.Setenv("CTF_CONFIGS", "configs/workflow-don.toml")
		case TopologyWorkflowGateway:
			setErr = os.Setenv("CTF_CONFIGS", "configs/workflow-gateway-don.toml")
		case TopologyWorkflowGatewayCapabilities:
			setErr = os.Setenv("CTF_CONFIGS", "configs/workflow-gateway-capabilities-don.toml")
		case TopologyMock:
			setErr = os.Setenv("CTF_CONFIGS", "configs/workflow-gateway-mock-don.toml")
		default:
			return fmt.Errorf("unknown topology: %s. Please use a known one or indicate which TOML config to use via CTF_CONFIGS environment variable", topologyFlag)
		}

		if setErr != nil {
			return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
		}

		fmt.Printf("Set CTF_CONFIGS environment variable to default value: %s\n", os.Getenv("CTF_CONFIGS"))
	}

	// set the defaults before the configs, so that they can be overridden by the configs
	defaultsSetErr := os.Setenv("CTF_CONFIGS", defaultCapabilitiesConfigFile+","+os.Getenv("CTF_CONFIGS"))
	if defaultsSetErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", defaultsSetErr)
	}

	return nil
}

func hasBuiltDockerImage(in *envconfig.Config, withPluginsDockerImageFlag string) bool {
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

func initDxTracker() {
	if dxTracker != nil {
		return
	}

	var trackerErr error
	dxTracker, trackerErr = tracking.NewDxTracker(GetDXGitHubVariableName, GetDXProductName)
	if trackerErr != nil {
		fmt.Fprintf(os.Stderr, "failed to create DX tracker: %s\n", trackerErr)
		dxTracker = &tracking.NoOpTracker{}
	}
}

func validateWorkflowTriggerAndCapabilities(in *envconfig.Config, withExampleFlag bool, workflowTrigger, withPluginsDockerImageFlag string) error {
	if withExampleFlag && workflowTrigger == WorkflowTriggerCron {
		// assume it has cron binary if we are using plugins image
		if withPluginsDockerImageFlag != "" {
			return nil
		}

		// otherwise, make sure we have cron binary path set in TOML config
		if in.CapabilityConfigs == nil {
			return errors.New("capability configs is not set in TOML config")
		}

		cronCapConfig, ok := in.CapabilityConfigs[cre.CronCapability]
		if !ok {
			return errors.New("cron capability config is not set in TOML config")
		}

		if cronCapConfig.BinaryPath == "" {
			return errors.New("cron binary path must be set in TOML config")
		}

		return nil
	}

	return nil
}

func ensureDockerIsRunning(ctx context.Context) error {
	dockerClient, dockerClientErr := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return errors.Wrap(dockerClientErr, "failed to create Docker client")
	}

	_, pingErr := dockerClient.Ping(ctx)
	if pingErr != nil {
		return errors.Wrap(pingErr, "docker is not running. Please start Docker and try again")
	}
	return nil
}

func ensureDockerImagesExist(ctx context.Context, logger zerolog.Logger, in *envconfig.Config, withPluginsDockerImageFlag string) error {
	// skip this check in CI, as we inject images at runtime and this check would fail
	if os.Getenv("CI") == "true" {
		return nil
	}

	if withPluginsDockerImageFlag != "" {
		if err := ensureDockerImageExists(ctx, logger, withPluginsDockerImageFlag); err != nil {
			return errors.Wrapf(err, "Plugins image '%s' not found. Make sure it exists locally", withPluginsDockerImageFlag)
		}
	}

	if in.JD != nil {
		if err := ensureDockerImageExists(ctx, logger, in.JD.Image); err != nil {
			return errors.Wrapf(err, "Job Distributor image '%s' not found. Make sure it exists locally or run 'go run . env setup' to pull it and other dependencies that also might be missing", in.JD.Image)
		}
	}

	for _, nodeSet := range in.NodeSets {
		for _, nodeSpec := range nodeSet.NodeSpecs {
			if nodeSpec.Node != nil && nodeSpec.Node.Image != "" {
				if err := ensureDockerImageExists(ctx, logger, nodeSpec.Node.Image); err != nil {
					return errors.Wrapf(err, "Node image '%s' not found. Make sure it exists locally", nodeSpec.Node.Image)
				}
			}
		}
	}

	return nil
}

// ensureDockerImageExists checks if the image exists locally, if not, it pulls it
// it returns nil if the image exists locally or was pulled successfully
// it returns an error if the image does not exist locally and pulling fails
// it doesn't handle registries that require authentication
func ensureDockerImageExists(ctx context.Context, logger zerolog.Logger, imageName string) error {
	dockerClient, dErr := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if dErr != nil {
		return errors.Wrap(dErr, "failed to create Docker client")
	}

	logger.Debug().Msgf("Checking if image '%s' exists locally", imageName)

	_, err := dockerClient.ImageInspect(ctx, imageName)
	if err != nil {
		logger.Debug().Msgf("Image '%s' not found locally, trying to pull it", imageName)

		ioRead, pullErr := dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
		if pullErr != nil {
			return fmt.Errorf("image '%s' not found locally and pulling failed", imageName)
		}
		defer ioRead.Close()

		logger.Debug().Msgf("Image '%s' pulled successfully", imageName)

		return nil
	}

	return nil
}

func globalPreRunFunc(cmd *cobra.Command, args []string) {
	if err := cobra.NoArgs(cmd, args); err != nil {
		fmt.Fprint(os.Stderr, libformat.RedText("\n%v. Please use flags to parameterise the command\n\n", err))
		os.Exit(1)
	}
}

func stateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "state",
		Short:            "State commands",
		Long:             `Commands to manage and view the state of the environment`,
		PersistentPreRun: globalPreRunFunc,
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(purgeStateCmd())
	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all state files in the environment",
		Long:  `List all state files in the environment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			envStateFilenames, sErr := allEnvironmentStateFiles()
			if sErr != nil {
				return errors.Wrap(sErr, "failed to get environment state files")
			}

			fmt.Println()
			fmt.Println("Environment state files:")
			for _, file := range envStateFilenames {
				fmt.Println("- " + file)
			}
			if len(envStateFilenames) == 0 {
				fmt.Println("- no state files found")
			}

			fmt.Println()
			fmt.Println("Cache folders:")
			cacheDirs, cErr := allCacheFolders()
			if cErr != nil {
				return errors.Wrap(cErr, "failed to get cache folders")
			}
			for _, dir := range cacheDirs {
				fmt.Println("- " + dir)
			}

			fmt.Println()

			return nil
		},
	}

	return cmd
}

func purgeStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge all state and cache files in the environment",
		Long:  `Purge all state and cache files in the environment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print()
			remErr := envconfig.RemoveAllEnvironmentStateDir(relativePathToRepoRoot)
			if remErr != nil {
				return errors.Wrap(remErr, "failed to remove environment state files")
			}

			fmt.Println()
			fmt.Println("Removing cache folders:")
			cacheDirs, cErr := allCacheFolders()
			if cErr != nil {
				return errors.Wrap(cErr, "failed to get cache folders")
			}

			for _, dir := range cacheDirs {
				rErr := os.RemoveAll(dir)
				if rErr != nil {
					fmt.Fprintf(os.Stderr, "failed to remove cache folder %s: %s\n", dir, rErr)
				} else {
					fmt.Printf("\n- removed cache folder: %s\n", dir)
				}
			}

			if len(cacheDirs) == 0 {
				fmt.Println("- no cache folders found")
			}

			fmt.Println()
			fmt.Println("Purge completed successfully")
			fmt.Println()

			return nil
		},
	}

	return cmd
}

func allCacheFolders() ([]string, error) {
	// TODO get this path from Beholder in the CTF
	knownCacheDirRoots := []string{"~/.local/share/beholder"}

	cacheDirs := []string{}
	for _, root := range knownCacheDirRoots {
		rootPath := strings.ReplaceAll(root, "~", os.Getenv("HOME"))
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, errors.Wrapf(err, "failed to read cache directory root: %s", rootPath)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				cacheDirs = append(cacheDirs, filepath.Join(rootPath, entry.Name()))
			}
		}
	}

	return cacheDirs, nil
}

func allEnvironmentStateFiles() ([]string, error) {
	stateDirAbs, absErr := filepath.Abs(filepath.Join(relativePathToRepoRoot, envconfig.StateDirname))
	if absErr != nil {
		return nil, errors.Wrap(absErr, "failed to get absolute path for state directory")
	}

	if _, statErr := os.Stat(stateDirAbs); os.IsNotExist(statErr) {
		return nil, nil
	}

	files, err := os.ReadDir(stateDirAbs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read state directory")
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no state files found in %s", stateDirAbs)
	}

	stateFiles := []string{}
	for _, file := range files {
		if !file.IsDir() {
			stateFiles = append(stateFiles, filepath.Join(stateDirAbs, file.Name()))
		}
	}

	return stateFiles, nil
}

func initLocalCREStageGen(in *envconfig.Config) *stagegen.StageGen {
	stages := 9
	if in.S3ProviderInput != nil {
		stages++
	}

	return stagegen.NewStageGen(stages, "STAGE")
}
