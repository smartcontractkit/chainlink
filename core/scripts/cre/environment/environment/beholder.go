package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

const DefaultBeholderConfigFile = "configs/chip-ingress.toml"

// getProtoSchemaSetFromGoMod reads the root go.mod file and extracts the commit ref
// from the github.com/smartcontractkit/chainlink-protos/workflows/go dependency.
// It returns a ProtoSchemaSet with hardcoded values matching default.toml config.
func getProtoSchemaSetFromGoMod() ([]chipingressset.ProtoSchemaSet, error) {
	// Find the root go.mod file (6 levels up from this file's location)
	// This file is at: core/scripts/cre/environment/environment/beholder.go
	// Root is at: go.mod
	goModPath := filepath.Join(relativePathToRepoRoot, "go.mod")
	absGoModPath, err := filepath.Abs(goModPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get absolute path to go.mod")
	}

	// Read go.mod file
	goModContent, err := os.ReadFile(absGoModPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read go.mod file")
	}

	// Parse go.mod
	modFile, err := modfile.Parse("go.mod", goModContent, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse go.mod file")
	}

	// Find the chainlink-protos/workflows/go dependency
	const targetModule = "github.com/smartcontractkit/chainlink-protos/workflows/go"
	var commitRef string
	for _, req := range modFile.Require {
		if req.Mod.Path == targetModule {
			// Version format: v0.0.0-YYYYMMDDHHMMSS-SHORTHASH
			// Extract the short hash after the last hyphen
			parts := strings.Split(req.Mod.Version, "-")
			if len(parts) >= 3 {
				commitRef = parts[len(parts)-1]
			}
			break
		}
	}

	if commitRef == "" {
		return nil, errors.Errorf("failed to find %s dependency or extract commit ref from go.mod", targetModule)
	}

	framework.L.Info().Msgf("Extracted commit ref from go.mod: %s", commitRef)

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
	cmd.Flags().DurationVarP(&timeout, "wait-on-error-timeout", "w", 15*time.Second, "Wait on error timeout (e.g. 10s, 1m, 1h)")

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

	out, startErr := chipingressset.New(in.ChipIngress)
	if startErr != nil {
		return errors.Wrap(startErr, "failed to create Chip Ingress set")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Started Chip Ingress stack in %.2f seconds", stageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Registering protos")))

	protoSchemaSets, getProtoErr := getProtoSchemaSetFromGoMod()
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

			protoSchemaSets, getProtoErr := getProtoSchemaSetFromGoMod()
			if getProtoErr != nil {
				return errors.Wrap(getProtoErr, "failed to get proto schema set from go.mod")
			}

			return parseConfigsAndRegisterProtos(cmd.Context(), protoSchemaSets, schemaURL)
		},
	}
	cmd.Flags().StringVarP(&schemaURL, "red-panda-schema-registry-url", "r", "http://localhost:"+chipingressset.DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT, "Red Panda Schema Registry URL")
	return cmd
}
