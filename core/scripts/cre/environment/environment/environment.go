package environment

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	ctypes "github.com/docker/docker/api/types/container"
	dc "github.com/docker/docker/client"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/sets"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/tracking"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const manualCtfCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove containers with 'ctf' label and delete their volumes`
const manualBeholderCleanupMsg = `unexpected startup error. this may have stranded resources. please manually remove the 'chip-ingress' stack`

var (
	binDir string

	defaultCapabilitiesConfigFile = "configs/capability_defaults.toml"
	defaultArtifactsPathFile      = "artifact_paths.json"
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

func init() {
	EnvironmentCmd.AddCommand(startCmd())
	EnvironmentCmd.AddCommand(stopCmd())
	EnvironmentCmd.AddCommand(workflowCmds())
	EnvironmentCmd.AddCommand(beholderCmds())
	EnvironmentCmd.AddCommand(swapCmds())

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

var StartCmdPreRunFunc = func(cmd *cobra.Command, args []string) {
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
				StartCmdRecoverHandlerFunc(recover(), cleanupWait)
			}()

			if doSetup {
				setupErr := RunSetup(cmd.Context(), SetupConfig{}, false, false)
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

			removeCacheErr := removeCacheFiles(removeAll)
			if removeCacheErr != nil {
				framework.L.Warn().Msgf("failed to remove cache files: %s\n", removeCacheErr)
			}

			if removeDirErr := os.RemoveAll("env_artifact"); removeDirErr != nil {
				framework.L.Warn().Msgf("failed to remove env_artifact folder: %s\n", removeDirErr)
			}

			// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
			setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			if setErr != nil {
				return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
			}

			cmdContext := cmd.Context()
			// Load and validate test configuration
			in, err := framework.Load[envconfig.Config](nil)
			if err != nil {
				return errors.Wrap(err, "failed to load test configuration")
			}

			// TODO since UnmarshalTOML is not supported by the TOML library we use :head_exploding:
			// we need to parse chain capabilities manually, but we need to handle it properly, maybe by adding hooks to Load()?
			for _, nodeSet := range in.NodeSets {
				if err := nodeSet.ParseChainCapabilities(); err != nil {
					return errors.Wrap(err, "failed to parse chain capabilities")
				}

				if err := nodeSet.ValidateChainCapabilities(in.Blockchains); err != nil {
					return errors.Wrap(err, "failed to validate chain capabilities")
				}
			}

			capabilityFlagsProvider := flags.NewDefaultCapabilityFlagsProvider()

			if err := in.Validate(capabilityFlagsProvider); err != nil {
				return errors.Wrap(err, "failed to validate test configuration")
			}

			homeChainIDInt, chainErr := strconv.Atoi(in.Blockchains[0].ChainID)
			if chainErr != nil {
				return fmt.Errorf("failed to convert chain ID to int: %w", chainErr)
			}

			defaultCapabilities, defaultCapabilitiesErr := sets.NewDefaultSet(libc.MustSafeUint64FromInt(homeChainIDInt), append(extraAllowedGatewayPorts, in.Fake.Port), []string{}, []string{"0.0.0.0/0"})
			if defaultCapabilitiesErr != nil {
				return errors.Wrap(defaultCapabilitiesErr, "failed to create default capabilities")
			}

			if err := validateWorkflowTriggerAndCapabilities(in, withExampleFlag, exampleWorkflowTrigger, withPluginsDockerImage); err != nil {
				return errors.Wrap(err, "either cron binary path must be set in TOML config (%s) or you must use Docker image with all capabilities included and passed via withPluginsDockerImageFlag")
			}

			output, startErr := StartCLIEnvironment(cmdContext, in, topology, withPluginsDockerImage, defaultCapabilities, capabilityFlagsProvider)
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
				deployErr := deployAndVerifyExampleWorkflow(cmdContext, homeChainOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl, gatewayURL, output.DonTopology.GatewayConnectorOutput.Configurations[0].Dons[0].ID, exampleWorkflowTimeout, exampleWorkflowTrigger, wfRegAddr.Hex())
				if deployErr != nil {
					fmt.Printf("Failed to deploy and verify example workflow: %s\n", deployErr)
				}
			}
			fmt.Print(libformat.PurpleText("\nEnvironment setup completed successfully in %.2f seconds\n\n", time.Since(provisioningStartTime).Seconds()))
			fmt.Print("To terminate execute:`go run . env stop`\n\n")

			if err := storeArtifacts(in); err != nil {
				return errors.Wrap(err, "failed to store artifacts")
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
	cmd.Flags().StringArrayVarP(&protoConfigs, "with-proto-configs", "c", []string{"./proto-configs/default.toml"}, "Protos configs to use (e.g. './proto-configs/config_one.toml,./proto-configs/config_two.toml')")
	cmd.Flags().BoolVarP(&doSetup, "auto-setup", "a", false, "Run setup before starting the environment")
	return cmd
}

// Store the config with cached output so subsequent runs can reuse the
// environment without full setup. Then persist absolute paths to the
// generated artifacts (env artifact JSON and the cached CTF config) in
// `artifact_paths.json`. System tests use these to reload environment
// state across runs (see `system-tests/tests/smoke/cre/capabilities_test.go`),
// where the cached config and env artifact are consumed to reconstruct
// the in-memory CLDF environment without re-provisioning.
//
// This makes local iteration and CI reruns faster and deterministic.
func storeArtifacts(in *envconfig.Config) error {
	if err := storeCTFConfigs(in); err != nil {
		return err
	}

	return saveArtifactPaths()
}

func storeCTFConfigs[ConfigType any](config *ConfigType) error {
	// hack, because CTF takes the first config file from the list to select the name of the cache file, we need to remove the default capabilities config file (which we added as the first one, so that other configs can override it)
	ctfConfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", ctfConfigs)
		if setErr != nil {
			framework.L.Warn().Msgf("failed to restore CTF_CONFIGS env var: %s", setErr)
		}
	}()

	splitConfigs := strings.Split(ctfConfigs, ",")
	if len(splitConfigs) > 1 {
		if strings.Contains(splitConfigs[0], defaultCapabilitiesConfigFile) {
			splitConfigs = splitConfigs[1:]
		}

		setErr := os.Setenv("CTF_CONFIGS", strings.Join(splitConfigs, ","))
		if setErr != nil {
			return errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
		}
	}

	cachedOutName := splitConfigs[0]
	if !strings.HasSuffix(splitConfigs[0], "-cache.toml") {
		cachedOutName = strings.ReplaceAll(splitConfigs[0], ".toml", "") + "-cache.toml"
	}

	if _, err := os.Stat(cachedOutName); err == nil {
		removeErr := os.Remove(cachedOutName)
		if removeErr != nil {
			return errors.Wrap(removeErr, "failed to remove cached config file")
		}
	}

	// just in case remove "-cache.toml" suffix, because if it's present in the env var name
	// config won't be cached as logic assumes it already exists
	setErr := os.Setenv("CTF_CONFIGS", strings.ReplaceAll(splitConfigs[0], "-cache.toml", ".toml"))
	if setErr != nil {
		return errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
	}

	storeErr := framework.Store(config)
	if storeErr != nil {
		return errors.Wrap(storeErr, "failed to store environment cached config")
	}

	return nil
}

type artifactPaths struct {
	EnvArtifact string `json:"env_artifact"`
	EnvConfig   string `json:"env_config"`
}

func saveArtifactPaths() error {
	artifactAbsPath, artifactAbsPathErr := filepath.Abs(filepath.Join(creenv.ArtifactDirName, creenv.ArtifactFileName))
	if artifactAbsPathErr != nil {
		return artifactAbsPathErr
	}

	// hack, because CTF takes the first config file from the list to select the name of the cache file, we need to remove the default capabilities config file (which we added as the first one, so that other configs can override it)
	ctfConfigs := os.Getenv("CTF_CONFIGS")
	// defer func() {
	// 	setErr := os.Setenv("CTF_CONFIGS", ctfConfigs)
	// 	if setErr != nil {
	// 		framework.L.Warn().Msgf("failed to restore CTF_CONFIGS env var: %s", setErr)
	// 	}
	// }()

	splitConfigs := strings.Split(ctfConfigs, ",")
	if len(splitConfigs) > 1 {
		if strings.Contains(splitConfigs[0], defaultCapabilitiesConfigFile) {
			splitConfigs = splitConfigs[1:]
		}

		// setErr := os.Setenv("CTF_CONFIGS", strings.Join(splitConfigs, ","))
		// if setErr != nil {
		// 	return errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
		// }
	}

	splitConfigs[0] = strings.ReplaceAll(splitConfigs[0], "-cache.toml", ".toml")

	ctfConfigsAbsPath, ctfConfigsAbsPathErr := filepath.Abs(splitConfigs[0])
	if ctfConfigsAbsPathErr != nil {
		return ctfConfigsAbsPathErr
	}

	ap := artifactPaths{
		EnvArtifact: artifactAbsPath,
		EnvConfig:   ctfConfigsAbsPath,
	}

	marshalled, mErr := json.Marshal(ap)
	if mErr != nil {
		return errors.Wrap(mErr, "failed to marshal artifact paths")
	}

	return os.WriteFile(defaultArtifactsPathFile, marshalled, 0600)
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

func stopCmd() *cobra.Command {
	var allFlag bool
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the environment",
		Long:  `Stops the local CRE environment (if it's not running, it just fallsthrough)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			removeErr := framework.RemoveTestContainers()
			if removeErr != nil {
				return errors.Wrap(removeErr, "failed to remove environment containers. Please remove them manually")
			}

			stopBeholderErr := stopBeholder()
			if stopBeholderErr != nil {
				return errors.Wrap(stopBeholderErr, "failed to stop beholder")
			}

			var shouldRemove shouldRemove
			if allFlag {
				shouldRemove = removeAll
			} else {
				shouldRemove = removeCurrentCtfConfigs
			}

			removeCacheErr := removeCacheFiles(shouldRemove)
			if removeCacheErr != nil {
				framework.L.Warn().Msgf("failed to remove cache files: %s\n", removeCacheErr)
			}

			if removeDirErr := os.RemoveAll("env_artifact"); removeDirErr != nil {
				framework.L.Warn().Msgf("failed to remove env_artifact folder: %s\n", removeDirErr)
			} else {
				framework.L.Debug().Msg("Removed env_artifact folder")
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
	in *envconfig.Config,
	topologyFlag string,
	withPluginsDockerImageFlag string,
	capabilities []cre.InstallableCapability,
	capabilityFlagsProvider cre.CapabilityFlagsProvider,
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
	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets: in.NodeSets,
		BlockchainsInput:          in.Blockchains,
		JdInput:                   *in.JD,
		InfraInput:                *in.Infra,
		S3ProviderInput:           in.S3ProviderInput,
		CapabilityConfigs:         in.CapabilityConfigs,
		CopyCapabilityBinaries:    withPluginsDockerImageFlag == "", // do not copy any binaries to the containers, if we are using plugins image (they already have them)
		Capabilities:              capabilities,
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
	dxTracker, trackerErr = tracking.NewDxTracker()
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

type shouldRemove = func(file string) bool

var removeAll = func(_ string) bool {
	return true
}

var removeCurrentCtfConfigs = func(file string) bool {
	ctfConfigs := os.Getenv("CTF_CONFIGS")
	if ctfConfigs != "" {
		for config := range strings.SplitSeq(ctfConfigs, ",") {
			if strings.Contains(file, strings.ReplaceAll(config, ".toml", "-cache.toml")) {
				return true
			}
		}

		return false
	}

	content, readErr := os.ReadFile(defaultArtifactsPathFile)
	if readErr != nil {
		return false
	}

	var paths artifactPaths
	if err := json.Unmarshal(content, &paths); err != nil {
		return false
	}

	if paths.EnvConfig == file {
		return true
	}

	return false
}

func removeCacheFiles(shouldRemove shouldRemove) error {
	framework.L.Info().Msg("Removing environment state files")

	cacheConfigPattern := "configs/*-cache.toml"
	cacheFiles, globErr := filepath.Glob(cacheConfigPattern)
	if globErr != nil {
		fmt.Fprintf(os.Stderr, "failed to find cache config files: %s\n", globErr)
	} else {
		for _, file := range cacheFiles {
			absFile, absFileErr := filepath.Abs(file)
			if absFileErr != nil {
				framework.L.Warn().Msgf("failed to get absolute path of cache config file %s: %s\n", file, absFileErr)
				continue
			}

			if shouldRemove(absFile) {
				if removeFileErr := os.Remove(file); removeFileErr != nil {
					framework.L.Warn().Msgf("failed to remove cache config file %s: %s\n", file, removeFileErr)
				} else {
					framework.L.Debug().Msgf("Removed cache config file: %s\n", file)
				}
			}
		}
	}

	return nil
}

func swapCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap",
		Short: "Swaps parts of the local CRE without restarting the environment",
	}

	cmd.AddCommand(capabilitySwapCmd())
	cmd.AddCommand(nodesSwapCmd())

	return cmd
}

func capabilitySwapCmd() *cobra.Command {
	var (
		capabilityFlag string
		binaryPath     string
	)

	cmd := &cobra.Command{
		Use:     "capability",
		Short:   "Swaps the capability binary of the Chainlink nodes in the environment",
		Long:    "Swaps the capability binary of the Chainlink nodes in the environment. Capability flag is used to find jobs with names containing the capability name, which are cancelled and approved, so that capability binary is reloaded. Only DONs that have the capability are impacted.",
		Aliases: []string{"c", "cap"},
		RunE: func(cmd *cobra.Command, args []string) error {
			supportedCapabilities := flags.NewSwappableCapabilityFlagsProvider()
			if !slices.Contains(supportedCapabilities.SupportedCapabilityFlags(), capabilityFlag) {
				return fmt.Errorf("capability %s cannot be hot-reloaded", capabilityFlag)
			}

			content, readErr := os.ReadFile(defaultArtifactsPathFile)
			if readErr != nil {
				return errors.Wrap(readErr, "failed to read artifact paths file")
			}

			var paths artifactPaths
			if err := json.Unmarshal(content, &paths); err != nil {
				return errors.Wrap(err, "failed to unmarshal artifact paths file")
			}

			setErr := os.Setenv("CTF_CONFIGS", strings.ReplaceAll(paths.EnvConfig, ".toml", "-cache.toml"))
			if setErr != nil {
				return errors.Wrap(setErr, "failed to set CTF_CONFIGS environment variable")
			}

			config, loadErr := framework.Load[envconfig.Config](nil)
			if loadErr != nil {
				return errors.Wrap(loadErr, "failed to load CTF config")
			}

			var envArtifact creenv.EnvArtifact
			artFile, artErr := os.ReadFile(paths.EnvArtifact)
			if artErr != nil {
				return errors.Wrap(artErr, "failed to read artifact file")
			}
			unmarshalErr := json.Unmarshal(artFile, &envArtifact)
			if unmarshalErr != nil {
				return errors.Wrap(unmarshalErr, "failed to unmarshal artifact file")
			}

			cldLogger := cldlogger.NewSingleFileLogger(nil)

			fullCldEnvOutput, _, loadErr := creenv.BuildFromSavedState(cmd.Context(), cldLogger, config, envArtifact)
			if loadErr != nil {
				return errors.Wrap(loadErr, "failed to load environment")
			}

			var nameContainsCapability = func(jobSpec string) bool {
				r := regexp.MustCompile(`name\s+=\s+"(.*)"`)
				matches := r.FindStringSubmatch(jobSpec)
				if len(matches) < 2 {
					return false
				}
				return strings.Contains(matches[1], capabilityFlag)
			}

			donIdxToNodeIdToJobIds := map[int]map[string][]string{}
			for idx, nodeSet := range fullCldEnvOutput.DonTopology.DonsWithMetadata {
				if !flags.HasFlagForAnyChain(nodeSet.Flags, capabilityFlag) {
					continue
				}
				donIdxToNodeIdToJobIds[idx] = map[string][]string{}
				for _, node := range nodeSet.DON.Nodes {
					framework.L.Info().Msgf("Cancelling all jobs for node %s", node.Name)
					jobIds, cancelErr := node.CancelAllJobs(cmd.Context(), nameContainsCapability)
					if cancelErr != nil {
						return errors.Wrap(cancelErr, "failed to cancel jobs")
					}
					framework.L.Info().Msgf("Cancelled %d jobs for node %s", len(jobIds), node.Name)
					donIdxToNodeIdToJobIds[idx][node.NodeID] = jobIds
				}
			}

			if len(donIdxToNodeIdToJobIds) == 0 {
				return fmt.Errorf("no nodes found with capability %s", capabilityFlag)
			}

			for donIdx, nodeIdToJobIds := range donIdxToNodeIdToJobIds {
				pattern := fullCldEnvOutput.DonTopology.DonsWithMetadata[donIdx].Name + "-node"
				capDir, dirErr := crecapabilities.DefaultContainerDirectory(config.Infra.Type)
				if dirErr != nil {
					return errors.Wrapf(dirErr, "failed to get default capabilities directory for infra type %s", config.Infra.Type)
				}

				copyErr := creworkflow.CopyArtifactToDockerContainers(binaryPath, pattern, capDir)
				if copyErr != nil {
					return errors.Wrapf(copyErr, "failed to copy %s capability binary to Docker containers with pattern %s", binaryPath, pattern)
				}

				for _, node := range fullCldEnvOutput.DonTopology.DonsWithMetadata[donIdx].DON.Nodes {
					jobIds, ok := nodeIdToJobIds[node.NodeID]
					if ok {
						framework.L.Info().Msgf("Approving %d jobs for node %s", len(jobIds), node.Name)
						approveErr := node.ApproveAllJobs(cmd.Context(), jobIds)
						if approveErr != nil {
							return errors.Wrap(approveErr, "failed to approve jobs")
						}
						framework.L.Info().Msgf("Approved %d jobs for node %s", len(jobIds), node.Name)
					}
				}
			}

			return storeArtifacts(config)
		},
	}

	cmd.Flags().StringVarP(&capabilityFlag, "name", "n", "", "The name of the capability to swap")
	cmd.Flags().StringVarP(&binaryPath, "binary", "b", "", "The binary path to swap")
	cmd.MarkFlagRequired("binary")
	cmd.MarkFlagRequired("name")

	return cmd
}

func nodesSwapCmd() *cobra.Command {
	var (
		forceFlag bool
	)

	cmd := &cobra.Command{
		Use:     "nodes",
		Short:   "Swaps the Docker images of the Chainlink nodes in the environment",
		Long:    "Swaps the Docker images of the Chainlink nodes in the environment. If environment is configured to build the Docker image, it will be rebuilt if any change is detected in the source code.",
		Aliases: []string{"n", "node"},
		RunE: func(cmd *cobra.Command, args []string) error {
			content, readErr := os.ReadFile(defaultArtifactsPathFile)
			if readErr != nil {
				return errors.Wrap(readErr, "failed to read artifact paths file")
			}

			var paths artifactPaths
			if err := json.Unmarshal(content, &paths); err != nil {
				return errors.Wrap(err, "failed to unmarshal artifact paths file")
			}

			setErr := os.Setenv("CTF_CONFIGS", strings.ReplaceAll(paths.EnvConfig, ".toml", "-cache.toml"))
			if setErr != nil {
				return errors.Wrap(setErr, "failed to set CTF_CONFIGS environment variable")
			}

			config, loadErr := framework.Load[envconfig.Config](nil)
			if loadErr != nil {
				return errors.Wrap(loadErr, "failed to load CTF config")
			}

			// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
			setErr = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			if setErr != nil {
				return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
			}

			for _, nodeSet := range config.NodeSets {
				framework.L.Info().Msgf("Removing Docker containers for DON %s", nodeSet.Name)
				containerIDs, containerIDsErr := findAllDockerContainerIDs(nodeSet.Name + "-node")
				if containerIDsErr != nil {
					return errors.Wrapf(containerIDsErr, "failed to find Docker containers for node set %s", nodeSet.Name)
				}

				for _, id := range containerIDs {
					framework.L.Debug().Msgf("Removing Docker container %s", id)
					dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
					if dockerClientErr != nil {
						return errors.Wrap(dockerClientErr, "failed to create Docker client")
					}

					if !forceFlag {
						stopErr := dockerClient.ContainerStop(context.Background(), id, ctypes.StopOptions{})
						if stopErr != nil {
							return errors.Wrapf(stopErr, "failed to stop Docker container %s", id)
						}
					}

					removeErr := dockerClient.ContainerRemove(context.Background(), id, ctypes.RemoveOptions{Force: forceFlag})
					if removeErr != nil {
						return errors.Wrapf(removeErr, "failed to remove Docker container %s", id)
					}
				}

				framework.L.Info().Msgf("Starting new Docker containers for DON %s", nodeSet.Name)
				nodeSet.Out = nil
				var nodesetErr error
				nodeSet.Out, nodesetErr = ns.NewSharedDBNodeSet(nodeSet.Input, config.Blockchains[0].Out)
				if nodesetErr != nil {
					time.Sleep(2 * time.Minute)
					return errors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSet.Name)
				}
			}

			return storeArtifacts(config)
		},
	}

	cmd.Flags().BoolVarP(&forceFlag, "force", "f", true, "Force removal of Docker containers")

	return cmd
}

func findAllDockerContainerIDs(pattern string) ([]string, error) {
	dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return nil, errors.Wrap(dockerClientErr, "failed to create Docker client")
	}

	containers, containersErr := dockerClient.ContainerList(context.Background(), ctypes.ListOptions{})
	if containersErr != nil {
		return nil, errors.Wrap(containersErr, "failed to list Docker containers")
	}

	containerIDs := []string{}
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, pattern) {
				containerIDs = append(containerIDs, container.ID)
			}
		}
	}

	return containerIDs, nil
}
