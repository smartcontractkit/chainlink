package environment

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	retry "github.com/avast/retry-go/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

const DefaultBeholderConfigFile = "configs/chip-ingress.toml"

// moduleInfo represents the JSON output from `go list -m -json`
type moduleInfo struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
}

// getSchemaSetFromGoMod uses `go list` to extract the version/commit ref
// from the github.com/smartcontractkit/chainlink-protos/workflows/go dependency.
// It returns a SchemaSet with hardcoded values matching default.toml config.
func getSchemaSetFromGoMod(ctx context.Context) ([]chipingressset.SchemaSet, error) {
	const targetModule = "github.com/smartcontractkit/chainlink-protos/workflows/go"

	// Get the absolute path to the repository root (where go.mod is located)
	repoRoot, err := filepath.Abs(relativePathToRepoRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get absolute path to repository root")
	}

	// Use `go list -m -json` to get module information
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-json", targetModule)
	cmd.Dir = repoRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run 'go list -m -json %s'", targetModule)
	}

	// Parse JSON output
	var modInfo moduleInfo
	if err := json.Unmarshal(output, &modInfo); err != nil {
		return nil, errors.Wrap(err, "failed to parse go list JSON output")
	}

	if modInfo.Version == "" {
		return nil, errors.Errorf("no version found for module %s", targetModule)
	}

	// Extract commit ref from version string
	// Support various formats:
	// 1. v1.2.1 -> use as-is
	// 2. v0.0.0-20211026045750-20ab5afb07e3 -> extract short hash (20ab5afb07e3)
	// 3. 2a35b54f48ae06be4cc81c768dc9cc9e92249571 -> full commit hash, use as-is
	// 4. v0.0.0-YYYYMMDDHHMMSS-SHORTHASH -> extract short hash
	commitRef := extractCommitRef(modInfo.Version)

	framework.L.Info().Msgf("Extracted commit ref for %s: %s (from version: %s)", targetModule, commitRef, modInfo.Version)

	// Return SchemaSet with hardcoded values from default.toml
	schemaSet := chipingressset.SchemaSet{
		URI:        "https://github.com/smartcontractkit/chainlink-protos",
		Ref:        commitRef,
		SchemaDir:  "workflows",
		ConfigFile: "chip-cre.json", // file with mappings of protobufs to subjects, together with references
	}

	return []chipingressset.SchemaSet{schemaSet}, nil
}

// extractCommitRef extracts a commit reference from various version formats
func extractCommitRef(version string) string {
	// If it looks like a full commit hash (40 hex characters, no dashes or dots)
	if len(version) == 40 && isHexString(version) {
		return version
	}

	// If version contains hyphens, it might be pseudo-version format:
	// v0.0.0-YYYYMMDDHHMMSS-SHORTHASH or v1.2.3-0.YYYYMMDDHHMMSS-SHORTHASH
	if strings.Contains(version, "-") {
		parts := strings.Split(version, "-")
		// The last part should be the short hash
		if len(parts) >= 2 {
			lastPart := parts[len(parts)-1]
			// Verify it looks like a hash (12 hex characters typically)
			if len(lastPart) >= 7 && isHexString(lastPart) {
				return lastPart
			}
		}
	}

	// Otherwise, use the version as-is (e.g., v1.2.1, v0.10.0)
	return version
}

// isHexString checks if a string contains only hexadecimal characters
func isHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

// getComposeFileFromGoMod extracts the version of chainlink-testing-framework/framework/components/dockercompose
// from go.mod and returns the URL to the docker-compose.yml file for that version.
// It caches the file locally to avoid re-downloading.
func getComposeFileFromGoMod(ctx context.Context) (string, error) {
	const targetModule = "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose"

	// Get the absolute path to the core/scripts directory (where go.mod is located for this package)
	scriptsDir, err := filepath.Abs("../../")
	if err != nil {
		return "", errors.Wrap(err, "failed to get absolute path to scripts directory")
	}

	// Use `go list -m -json` to get module information
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-json", targetModule)
	cmd.Dir = scriptsDir

	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "failed to run 'go list -m -json %s'", targetModule)
	}

	// Parse JSON output
	var modInfo moduleInfo
	if unmarshalErr := json.Unmarshal(output, &modInfo); unmarshalErr != nil {
		return "", errors.Wrap(unmarshalErr, "failed to parse go list JSON output")
	}

	if modInfo.Version == "" {
		return "", errors.Errorf("no version found for module %s", targetModule)
	}

	// Determine the GitHub ref to use
	version := modInfo.Version
	var githubRef string
	var cacheKey string

	// Check if it's a pseudo-version (format: v0.1.19-0.20260130101725-678aa4ae7ce6)
	if strings.Contains(version, "-0.") && strings.Count(version, "-") >= 2 {
		// Extract commit hash from pseudo-version
		parts := strings.Split(version, "-")
		commitHash := parts[len(parts)-1]
		githubRef = commitHash // Use commit hash directly
		cacheKey = commitHash  // Use commit hash for cache
		framework.L.Info().Msgf("Detected pseudo-version: %s, using commit: %s", version, commitHash)
	} else {
		// It's a proper version tag
		githubRef = "framework/components/dockercompose/" + version
		cacheKey = version
		framework.L.Info().Msgf("Detected version tag: %s", version)
	}

	// Check if file is already cached locally
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}
	cacheDir := filepath.Join(homeDir, ".local", "share", "chip_ingress_set")
	cachedFile := filepath.Join(cacheDir, fmt.Sprintf("docker-compose-%s.yml", cacheKey))

	if _, statErr := os.Stat(cachedFile); statErr == nil {
		framework.L.Info().Msgf("Using cached compose file: %s", cachedFile)
		return "file://" + cachedFile, nil
	}

	// Download and cache the file
	url := fmt.Sprintf("https://raw.githubusercontent.com/smartcontractkit/chainlink-testing-framework/%s/framework/components/dockercompose/chip_ingress_set/docker-compose.yml", githubRef)
	framework.L.Info().Msgf("Downloading compose file from: %s", url)

	// Create cache directory
	if mkdirErr := os.MkdirAll(cacheDir, 0o755); mkdirErr != nil {
		return "", errors.Wrap(mkdirErr, "failed to create cache directory")
	}

	// Download file with retries to withstand transient GitHub/network issues
	var respBody []byte
	downloadErr := retry.Do(
		func() error {
			req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if reqErr != nil {
				return errors.Wrapf(reqErr, "failed to create HTTP request for %s", url)
			}

			resp, httpErr := http.DefaultClient.Do(req)
			if httpErr != nil {
				return errors.Wrapf(httpErr, "failed to download compose file from %s", url)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return errors.Errorf("failed to download compose file: HTTP %d", resp.StatusCode)
			}

			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				return errors.Wrap(readErr, "failed to read compose file contents")
			}
			respBody = bodyBytes
			return nil
		},
		retry.Context(ctx),
		retry.Delay(500*time.Millisecond),
		retry.Attempts(5),
		retry.DelayType(retry.BackOffDelay),
	)
	if downloadErr != nil {
		return "", errors.Wrap(downloadErr, "failed to download compose file")
	}

	// Save to cache
	if writeErr := os.WriteFile(cachedFile, respBody, 0o644); writeErr != nil { //nolint: gosec // it's fine for permissions to be a bit wider
		return "", errors.Wrap(writeErr, "failed to write compose file to cache")
	}

	framework.L.Info().Msgf("Cached compose file at: %s", cachedFile)
	return "file://" + cachedFile, nil
}

func beholderCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "beholder",
		Short: "Beholder commands",
		Long:  `Commands to manage the Beholder stack`,
	}

	cmd.AddCommand(startBeholderCmd())
	cmd.AddCommand(stopBeholderCmd)
	cmd.AddCommand(createKafkaTopicsCmd())
	cmd.AddCommand(fetchAndRegisterProtosCmd())

	return cmd
}

func startBeholderCmd() *cobra.Command {
	var (
		timeout time.Duration
		port    int
	)
	cmd := &cobra.Command{
		Use:              "start",
		Short:            "Start the Beholder",
		Long:             `Start the Beholder`,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			initDxTracker()
			var startBeholderErr error

			defer func() {
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
			}()

			// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
			setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			if setErr != nil {
				return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
			}

			startBeholderErr = startBeholder(cmd.Context(), timeout, port)
			if startBeholderErr != nil {
				// remove the stack if the error is not related to proto registration
				if !strings.Contains(startBeholderErr.Error(), protoRegistrationErrMsg) {
					waitToCleanUp(timeout)
					beholderRemoveErr := framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
					if beholderRemoveErr != nil {
						fmt.Fprint(os.Stderr, errors.Wrap(beholderRemoveErr, manualBeholderCleanupMsg).Error())
					}
				}
				return errors.Wrap(startBeholderErr, "failed to start Beholder")
			}

			return nil
		},
	}

	cmd.Flags().DurationVarP(&timeout, "wait-on-error-timeout", "w", 15*time.Second, "Time to wait before removing Docker containers if environment fails to start (e.g. 10s, 1m, 1h)")
	cmd.Flags().IntVarP(&port, "grpc-port", "g", mustStringToInt(chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT), "GRPC port for Chip Ingress")

	return cmd
}

func mustStringToInt(in string) int {
	out, err := strconv.Atoi(in)
	if err != nil {
		panic(fmt.Errorf("failed to parse default ChIP Ingress port: %w", err))
	}

	return out
}

var stopBeholderCmd = &cobra.Command{
	Use:              "stop",
	Short:            "Stop the Beholder",
	Long:             "Stop the Beholder",
	PersistentPreRun: globalPreRunFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopBeholder()
	},
}

func stopBeholder() error {
	setErr := os.Setenv("CTF_CONFIGS", DefaultBeholderConfigFile)
	if setErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
	}

	removeCacheErr := removeBeholderStateFiles(relativePathToRepoRoot)
	if removeCacheErr != nil {
		framework.L.Warn().Msgf("failed to remove cache files: %s\n", removeCacheErr)
	}

	return framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
}

func removeBeholderStateFiles(relativePathToRepoRoot string) error {
	path := filepath.Join(relativePathToRepoRoot, envconfig.StateDirname, envconfig.ChipIngressStateFilename)
	absPath, absErr := filepath.Abs(path)
	if absErr != nil {
		return errors.Wrap(absErr, "error getting absolute path for chip ingress state file")
	}

	return os.Remove(absPath)
}

func isPortAvailable(addr string) bool {
	lc := net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return false // already in use or permission denied
	}
	_ = l.Close()
	return true
}

var protoRegistrationErrMsg = "proto registration failed"

func startBeholder(cmdContext context.Context, cleanupWait time.Duration, port int) (startupErr error) {
	// just in case, remove the stack if it exists
	_ = framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)

	defer func() {
		p := recover()

		if p != nil {
			fmt.Println("Panicked when starting Beholder")

			if err, ok := p.(error); ok {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				startupErr = err
			} else {
				fmt.Fprintf(os.Stderr, "panic: %v\n", p)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				startupErr = fmt.Errorf("panic: %v", p)
			}

			time.Sleep(cleanupWait)

			beholderRemoveErr := framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
			if beholderRemoveErr != nil {
				fmt.Fprint(os.Stderr, errors.Wrap(beholderRemoveErr, manualBeholderCleanupMsg).Error())
			}

			os.Exit(1)
		}
	}()

	stageGen := stagegen.NewStageGen(3, "STAGE")
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting Chip Ingress stack")))

	if !isPortAvailable(":" + strconv.Itoa(port)) {
		return fmt.Errorf(`port %d is already in use. Most probably an instance of ChIP Test Sink is already running.
If you want to use both together start ChIP Ingress on a different port with '--grpc-port' flag
and make sure that the sink is pointing to correct upstream endpoint ('localhost:<grpc-port>' in most cases)`, port)
	}

	// set both internal and external (host) ChIP Ingress GRPC port to the same value
	if err := os.Setenv(chipingressset.ChipIngressGRPCHostPortEnvVar, strconv.Itoa(port)); err != nil {
		return fmt.Errorf("failed to set %s environment variable: %w", chipingressset.ChipIngressGRPCHostPortEnvVar, err)
	}

	if err := os.Setenv(chipingressset.ChipIngressGRPCPortEnvVar, strconv.Itoa(port)); err != nil {
		return fmt.Errorf("failed to set %s environment variable: %w", chipingressset.ChipIngressGRPCPortEnvVar, err)
	}

	// we want to restore previous configs, because Beholder might be started within the context of a different command,
	// which is also using CTF_CONFIGS environment variable to load or later store configs
	previousCTFConfig := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", previousCTFConfig)
		if setErr != nil {
			framework.L.Warn().Msgf("failed to restore previous CTF_CONFIGS environment variable: %s", setErr)
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", DefaultBeholderConfigFile)
	if setErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
	}

	// Load and validate test configuration
	in, err := framework.Load[envconfig.ChipIngressConfig](nil)
	if err != nil {
		return errors.Wrap(err, "failed to load test configuration")
	}

	// Auto-detect compose file from go.mod if not specified
	if in.ChipIngress != nil && in.ChipIngress.ComposeFile == "" {
		composeFile, composeErr := getComposeFileFromGoMod(cmdContext)
		if composeErr != nil {
			return errors.Wrap(composeErr, "failed to get compose file from go.mod")
		}
		in.ChipIngress.ComposeFile = composeFile
	}

	out, startErr := chipingressset.NewWithContext(cmdContext, in.ChipIngress)
	if startErr != nil {
		return errors.Wrap(startErr, "failed to create Chip Ingress set")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Started Chip Ingress stack in %.2f seconds", stageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Registering protos")))

	schemaSets, setErr := getSchemaSetFromGoMod(cmdContext)
	if setErr != nil {
		return errors.Wrap(setErr, "failed to get chainlink-proto version from go.mod")
	}

	registerProtosErr := parseConfigsAndRegisterProtos(cmdContext, schemaSets, out.ChipIngress)
	if registerProtosErr != nil {
		return errors.Wrap(registerProtosErr, "failed to register protos")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Registered protos in %.2f seconds", stageGen.Elapsed().Seconds())))

	fmt.Println()
	framework.L.Info().Msgf("Red Panda Console URL: %s", out.RedPanda.ConsoleExternalURL)

	topicsErr := chipingressset.CreateTopics(cmdContext, out.RedPanda.KafkaExternalURL, in.Kafka.Topics)
	if topicsErr != nil {
		return errors.Wrap(topicsErr, "failed to create topics")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Created topics in %.2f seconds", stageGen.Elapsed().Seconds())))

	for _, topic := range in.Kafka.Topics {
		framework.L.Info().Msgf("Topic URL: %s", fmt.Sprintf("%s/topics/%s", out.RedPanda.ConsoleExternalURL, topic))
	}
	fmt.Println()
	fmt.Println("To exclude a flood of heartbeat messages it is recommended that you register a JS filter with following code: `return value.msg !== 'heartbeat';`")
	fmt.Println()
	fmt.Print("To terminate Beholder stack execute: `go run . env beholder stop`\n\n")

	return in.Store(envconfig.MustChipIngressStateFileAbsPath(relativePathToRepoRoot))
}

func parseConfigsAndRegisterProtos(ctx context.Context, schemaSets []chipingressset.SchemaSet, chipIngressOutput *chipingressset.ChipIngressOutput) error {
	if len(schemaSets) == 0 {
		framework.L.Warn().Msg("no proto configs provided, skipping proto registration")

		return nil
	}

	for _, protoSchemaSet := range schemaSets {
		framework.L.Info().Msgf("Registering and fetching proto from %s", protoSchemaSet.URI)
		framework.L.Info().Msgf("Proto schema set config: %+v", protoSchemaSet)
	}

	reposErr := chipingressset.FetchAndRegisterProtos(
		ctx,
		nil, // GH client will be created dynamically, if needed
		chipIngressOutput,
		schemaSets,
	)
	if reposErr != nil {
		return errors.Wrap(reposErr, protoRegistrationErrMsg+"failed to fetch and register protos")
	}
	return nil
}

func createKafkaTopicsCmd() *cobra.Command {
	var (
		url    string
		topics []string
		purge  bool
	)
	cmd := &cobra.Command{
		Use:              "create-topics",
		Short:            "Create Kafka topics",
		Long:             `Create Kafka topics (with or without removing existing topics)`,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				return errors.New("red-panda-kafka-url cannot be empty")
			}

			if len(topics) == 0 {
				return errors.New("kafka topics list cannot be empty")
			}

			if purge {
				topicsErr := chipingressset.DeleteAllTopics(cmd.Context(), url)
				if topicsErr != nil {
					return errors.Wrap(topicsErr, "failed to remove topics")
				}
			}

			topicsErr := chipingressset.CreateTopics(cmd.Context(), url, topics)
			if topicsErr != nil {
				return errors.Wrap(topicsErr, "failed to create topics")
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&url, "red-panda-kafka-url", "k", "localhost:"+chipingressset.DEFAULT_RED_PANDA_KAFKA_PORT, "Red Panda Kafka URL")
	cmd.Flags().StringArrayVarP(&topics, "topics", "t", []string{}, "Kafka topics to create (e.g. 'topic1,topic2')")
	cmd.Flags().BoolVarP(&purge, "purge-topics", "p", false, "Remove existing Kafka topics")
	_ = cmd.MarkFlagRequired("topics")
	_ = cmd.MarkFlagRequired("red-panda-kafka-url")

	return cmd
}

func fetchAndRegisterProtosCmd() *cobra.Command {
	var (
		chipIngressGRPCURL string
	)
	cmd := &cobra.Command{
		Use:              "register-protos",
		Short:            "Fetch and register protos",
		Long:             `Fetch and register protos`,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaSets, setErr := getSchemaSetFromGoMod(cmd.Context())
			if setErr != nil {
				return errors.Wrap(setErr, "failed to get proto schema set from go.mod")
			}

			return parseConfigsAndRegisterProtos(cmd.Context(), schemaSets, &chipingressset.ChipIngressOutput{
				GRPCExternalURL: chipIngressGRPCURL,
			})
		},
	}
	cmd.Flags().StringVarP(&chipIngressGRPCURL, "chip-ingress-grpc-url", "h", "localhost:"+chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT, "Chip Ingress GRPC URL")
	return cmd
}
