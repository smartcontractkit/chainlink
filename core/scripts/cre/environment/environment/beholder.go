package environment

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

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

// getProtoSchemaSetFromGoMod uses `go list` to extract the version/commit ref
// from the github.com/smartcontractkit/chainlink-protos/workflows/go dependency.
// It returns a ProtoSchemaSet with hardcoded values matching default.toml config.
func getProtoSchemaSetFromGoMod(ctx context.Context) ([]chipingressset.ProtoSchemaSet, error) {
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

	// Return ProtoSchemaSet with hardcoded values from default.toml
	protoSchemaSet := chipingressset.ProtoSchemaSet{
		URI:           "https://github.com/smartcontractkit/chainlink-protos",
		Ref:           commitRef,
		Folders:       []string{"workflows"},
		SubjectPrefix: "cre-",
		ExcludeFiles:  []string{},
	}

	return []chipingressset.ProtoSchemaSet{protoSchemaSet}, nil
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

			startBeholderErr = startBeholder(cmd.Context(), timeout)
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

	return cmd
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

var protoRegistrationErrMsg = "proto registration failed"

func startBeholder(cmdContext context.Context, cleanupWait time.Duration) (startupErr error) {
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

	out, startErr := chipingressset.NewWithContext(cmdContext, in.ChipIngress)
	if startErr != nil {
		return errors.Wrap(startErr, "failed to create Chip Ingress set")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Started Chip Ingress stack in %.2f seconds", stageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Registering protos")))

	protoSchemaSets, getProtoErr := getProtoSchemaSetFromGoMod(cmdContext)
	if getProtoErr != nil {
		return errors.Wrap(getProtoErr, "failed to get proto schema set from go.mod")
	}

	registerProtosErr := parseConfigsAndRegisterProtos(cmdContext, protoSchemaSets, out.RedPanda.SchemaRegistryExternalURL)
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

func parseConfigsAndRegisterProtos(ctx context.Context, protoSchemaSets []chipingressset.ProtoSchemaSet, schemaRegistryExternalURL string) error {
	if len(protoSchemaSets) == 0 {
		framework.L.Warn().Msg("no proto configs provided, skipping proto registration")

		return nil
	}

	for _, protoSchemaSet := range protoSchemaSets {
		framework.L.Info().Msgf("Registering and fetching proto from %s", protoSchemaSet.URI)
		framework.L.Info().Msgf("Proto schema set config: %+v", protoSchemaSet)
	}

	reposErr := chipingressset.DefaultRegisterAndFetchProtos(
		ctx,
		nil, // GH client will be created dynamically, if needed
		protoSchemaSets,
		schemaRegistryExternalURL,
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
		schemaURL string
	)
	cmd := &cobra.Command{
		Use:              "register-protos",
		Short:            "Fetch and register protos",
		Long:             `Fetch and register protos`,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use default values if not provided
			if schemaURL == "" {
				schemaURL = "http://localhost:" + chipingressset.DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT
			}

			protoSchemaSets, getProtoErr := getProtoSchemaSetFromGoMod(cmd.Context())
			if getProtoErr != nil {
				return errors.Wrap(getProtoErr, "failed to get proto schema set from go.mod")
			}

			return parseConfigsAndRegisterProtos(cmd.Context(), protoSchemaSets, schemaURL)
		},
	}
	cmd.Flags().StringVarP(&schemaURL, "red-panda-schema-registry-url", "r", "http://localhost:"+chipingressset.DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT, "Red Panda Schema Registry URL")

	return cmd
}
