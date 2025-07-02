package environment

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v72/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

var (
	withBeholderFlag              bool
	protoConfigsFlag              []string
	redPandaKafkaURLFlag          string
	redPandaSchemaRegistryURLFlag string
	kafkaCreateTopicsFlag         []string
	kafkaRemoveTopicsFlag         bool
)

type ChipIngressConfig struct {
	ChipIngress *chipingressset.Input `toml:"chip_ingress"`
	Kafka       *KafkaConfig          `toml:"kafka"`
}

type KafkaConfig struct {
	Topics []string `toml:"topics"`
}

var startBeholderCmd = &cobra.Command{
	Use:   "start-beholder",
	Short: "Start the Beholder",
	Long:  `Start the Beholder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if topologyFlag != TopologySimplified && topologyFlag != TopologyFull {
			return fmt.Errorf("invalid topology: %s. Valid topologies are: %s, %s", topologyFlag, TopologySimplified, TopologyFull)
		}

		// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
		setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if setErr != nil {
			return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
		}

		dockerNetworks, dockerNetworksErr := getCtfDockerNetworks()
		if dockerNetworksErr != nil {
			return errors.Wrap(dockerNetworksErr, "failed to get CTF Docker networks")
		}

		startBeholderErr := startBeholder(cmd.Context(), protoConfigsFlag, dockerNetworks)
		if startBeholderErr != nil {
			// remove the stack if the error is not related to proto registration
			if !strings.Contains(startBeholderErr.Error(), protoRegistrationErrMsg) {
				WaitOnErrorTimeoutDurationFn(waitOnErrorTimeoutFlag)
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

var stopBeholderCmd = &cobra.Command{
	Use:   "stop-beholder",
	Short: "Stop the Beholder",
	Long:  `Stop the Beholder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
	},
}

var protoRegistrationErrMsg = "proto registration failed"

func startBeholder(cmdContext context.Context, protoConfigsFlag []string, dockerNetworks []string) (startupErr error) {
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

			WaitOnErrorTimeoutDurationFn(waitOnErrorTimeoutFlag)

			beholderRemoveErr := framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
			if beholderRemoveErr != nil {
				fmt.Fprint(os.Stderr, errors.Wrap(beholderRemoveErr, manualBeholderCleanupMsg).Error())
			}
		}
	}()

	stageGen := creenv.NewStageGen(3, "STAGE")
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting Chip Ingress stack")))

	setErr := os.Setenv("CTF_CONFIGS", "configs/chip-ingress.toml")
	if setErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
	}

	// Load and validate test configuration
	in, err := framework.Load[ChipIngressConfig](nil)
	if err != nil {
		return errors.Wrap(err, "failed to load test configuration")
	}

	// connect to existing network if provided, that should only be used, when chip-ingress is started for an already running environment
	if len(dockerNetworks) > 0 {
		in.ChipIngress.ExtraDockerNetworks = append(in.ChipIngress.ExtraDockerNetworks, dockerNetworks...)
	}

	out, startErr := chipingressset.New(in.ChipIngress)
	if startErr != nil {
		return errors.Wrap(startErr, "failed to create Chip Ingress set")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Started Chip Ingress stack in %.2f seconds", stageGen.Elapsed().Seconds())))
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Registering protos")))

	registerProtosErr := parseConfigsAndRegisterProtos(cmdContext, protoConfigsFlag, out.RedPanda.SchemaRegistryExternalURL)
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
	fmt.Print("To terminate Beholder stack execute: `go run . env stop-beholder`\n\n")

	return nil
}

func parseConfigsAndRegisterProtos(ctx context.Context, protoConfigsFlag []string, schemaRegistryExternalURL string) error {
	var protoSchemaSets []chipingressset.ProtoSchemaSet
	for _, protoConfig := range protoConfigsFlag {
		file, fileErr := os.ReadFile(protoConfig)
		if fileErr != nil {
			return errors.Wrap(fileErr, protoRegistrationErrMsg+"failed to read proto config file: "+protoConfig)
		}

		type wrappedProtoSchemaSets struct {
			ProtoSchemaSets []chipingressset.ProtoSchemaSet `toml:"proto_schema_sets"`
		}

		var schemaSets wrappedProtoSchemaSets
		if err := toml.Unmarshal(file, &schemaSets); err != nil {
			return errors.Wrap(err, protoRegistrationErrMsg+"failed to unmarshal proto config file: "+protoConfig)
		}

		protoSchemaSets = append(protoSchemaSets, schemaSets.ProtoSchemaSets...)
	}

	if len(protoSchemaSets) == 0 {
		framework.L.Warn().Msg("no proto configs provided, skipping proto registration")

		return nil
	}

	for _, protoSchemaSet := range protoSchemaSets {
		framework.L.Info().Msgf("Registering and fetching proto from %s", protoSchemaSet.Repository)
		framework.L.Info().Msgf("Proto schema set config: %+v", protoSchemaSet)
	}

	var client *github.Client
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		framework.L.Warn().Msg("GITHUB_TOKEN is not set, using unauthenticated GitHub client. This may cause rate limiting issues when downloading proto files")
		client = github.NewClient(nil)
	}

	reposErr := chipingressset.DefaultRegisterAndFetchProtos(ctx, client, protoSchemaSets, schemaRegistryExternalURL)
	if reposErr != nil {
		return errors.Wrap(reposErr, protoRegistrationErrMsg+"failed to fetch and register protos")
	}
	return nil
}

var createKafkaTopicsCmd = &cobra.Command{
	Use:   "create-kafka-topics",
	Short: "Create Kafka topics",
	Long:  `Create Kafka topics (with or without removing existing topics)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if redPandaKafkaURLFlag == "" {
			return errors.New("red-panda-kafka-url cannot be empty")
		}

		if len(kafkaCreateTopicsFlag) == 0 {
			return errors.New("kafka topics list cannot be empty")
		}

		if kafkaRemoveTopicsFlag {
			topicsErr := chipingressset.DeleteAllTopics(cmd.Context(), redPandaKafkaURLFlag)
			if topicsErr != nil {
				return errors.Wrap(topicsErr, "failed to remove topics")
			}
		}

		topicsErr := chipingressset.CreateTopics(cmd.Context(), redPandaKafkaURLFlag, kafkaCreateTopicsFlag)
		if topicsErr != nil {
			return errors.Wrap(topicsErr, "failed to create topics")
		}

		return nil
	},
}

var fetchAndRegisterProtosCmd = &cobra.Command{
	Use:   "fetch-and-register-protos",
	Short: "Fetch and register protos",
	Long:  `Fetch and register protos`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if redPandaSchemaRegistryURLFlag == "" {
			return errors.New("red-panda-schema-registry-url cannot be empty")
		}

		if len(protoConfigsFlag) == 0 {
			framework.L.Warn().Msg("no proto configs provided, skipping proto registration")

			return nil
		}

		return parseConfigsAndRegisterProtos(cmd.Context(), protoConfigsFlag, redPandaSchemaRegistryURLFlag)
	},
}
