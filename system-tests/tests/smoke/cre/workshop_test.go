package cre

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"

	common_events "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflow_events "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

func Test_V2_Workflow_Workshop(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping workshop test in CI")
	}

	testLogger := framework.L

	/*
		TEST SETUP:
		- set required env vars
		- load test config
		- start DON
		- deploy contracts
		- create jobs
		- register workflow
	*/

	// set required env vars
	setPkErr := os.Setenv("PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80") // not a secret, it's a known developer key used by Anvil
	require.NoError(t, setPkErr, "failed to set PRIVATE_KEY")

	setCtfConfigsErr := os.Setenv("CTF_CONFIGS", "workshop_test.toml")
	require.NoError(t, setCtfConfigsErr, "failed to set CTF_CONFIGS")

	// load test config
	in, err := framework.Load[V2WorkflowTestConfig](t)
	require.NoError(t, err, "couldn't load test config")

	// setup test environment
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
	require.NoError(t, pathErr, "failed to get default container directory")

	chainIDInt, err := strconv.Atoi(in.Blockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets: []*types.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],
				Capabilities:       []string{types.CronCapability, types.OCR3Capability, types.CustomComputeCapability, types.WriteEVMCapability},
				DONTypes:           []string{types.WorkflowDON, types.GatewayDON},
				BootstrapNodeIndex: 0,
			},
		},
		CapabilitiesContractFactoryFunctions: []func([]string) []keystone_changeset.DONCapabilityWithConfig{
			croncap.CronCapabilityFactoryFn,
			consensuscap.OCR3CapabilityFactoryFn,
		},
		BlockchainsInput: []*types.WrappedBlockchainInput{in.Blockchain},
		JdInput:          *in.JD,
		InfraInput:       *in.Infra,
		CustomBinariesPaths: map[string]string{
			types.CronCapability: in.Dependencies.CronBinaryPath,
		},
		JobSpecFactoryFunctions: []types.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			crecron.CronJobSpecFactoryFn(filepath.Join(containerPath, filepath.Base(in.Dependencies.CronBinaryPath))),
			cregateway.GatewayJobSpecFactoryFn([]int{}, []string{}, []string{"0.0.0.0/0"}),
		},
		ConfigFactoryFunctions: []types.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
		CustomAnvilMiner: in.CustomAnvilMiner,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	// compile and upload workflow
	containerTargetDir := "/home/chainlink/workflows"
	testLogger.Info().Msg("Proceeding to register test workflow...")
	workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(universalSetupOutput.CldEnvironment.ExistingAddresses, universalSetupOutput.BlockchainOutput[0].ChainSelector, keystone_changeset.WorkflowRegistry.String()) //nolint:staticcheck // won't migrate now
	require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainID)

	// TODO: add code that will compile your workflow
	testWorkflowPath := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/this_is_not_my_workflow.go"
	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(testWorkflowPath, "test-workflow")
	require.NoError(t, compileErr, "failed to compile workflow '%s'", testWorkflowPath)

	copyErr := creworkflow.CopyWorkflowToDockerContainers(compressedWorkflowWasmPath, "workflow-node", containerTargetDir)
	require.NoError(t, copyErr, "failed to copy workflow to docker containers")

	t.Cleanup(func() {
		_ = os.Remove(compressedWorkflowWasmPath)
	})

	// register workflow
	registerErr := creworkflow.RegisterWithContract(
		t.Context(),
		universalSetupOutput.BlockchainOutput[0].SethClient,
		workflowRegistryAddress,
		universalSetupOutput.DonTopology.WorkflowDonID,
		"test-workflow",
		"file://"+compressedWorkflowWasmPath,
		nil, // no config URL
		nil, // no secrets URL
		&containerTargetDir,
	)
	require.NoError(t, registerErr, "failed to register workflow")

	/*
		TEST EXECUTION:
		- consume Kafka messages in loop until workflow execution is detected
		- or until we detect workflow engine initialization failure
	*/

	listenerCtx, cancelListener := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(func() {
		cancelListener()
	})

	kafkaErrChan := make(chan error, 1)
	messageChan := make(chan proto.Message, 10)

	// We are interested in UserLogs (successful execution)
	// or BaseMessage with specific error message (engine initialization failure)
	messageTypes := map[string]func() proto.Message{
		"workflows.v1.UserLogs": func() proto.Message {
			return &workflow_events.UserLogs{}
		},
		"BaseMessage": func() proto.Message {
			return &common_events.BaseMessage{}
		},
	}

	// Start listening for messages in the background
	go func() {
		listenForKafkaMessages(listenerCtx, testLogger, "localhost:19092", "cre", messageTypes, messageChan, kafkaErrChan)
	}()

	// TODO: add a variable named expectedUserLog which contains the message you added to your workflow
	expectedUserLog := "This is not my message"

	foundExpectedLog := make(chan bool, 1) // Channel to signal when expected log is found
	foundErrorLog := make(chan bool, 1)    // Channel to signal when engine initialization failure is detected
	receivedUserLogs := 0
	// Start message processor goroutine
	go func() {
		for {
			select {
			case <-listenerCtx.Done():
				return
			case msg := <-messageChan:
				// Process received messages
				switch typedMsg := msg.(type) {
				case *common_events.BaseMessage:
					if strings.Contains(typedMsg.Msg, "Workflow Engine initialization failed") {
						foundErrorLog <- true
					}
				case *workflow_events.UserLogs:
					testLogger.Info().
						Msg("🎉 Received UserLogs message in test")
					receivedUserLogs++

					for _, logLine := range typedMsg.LogLines {
						if strings.Contains(logLine.Message, expectedUserLog) {
							testLogger.Info().
								Str("expected_log", expectedUserLog).
								Str("found_message", strings.TrimSpace(logLine.Message)).
								Msg("🎯 Found expected user log message!")

							select {
							case foundExpectedLog <- true:
							default: // Channel might already have a value
							}
							return // Exit the processor goroutine
						} else {
							testLogger.Warn().
								Str("expected_log", expectedUserLog).
								Str("found_message", strings.TrimSpace(logLine.Message)).
								Msg("Received UserLogs message, but it does not match expected log")
						}
					}
				default:
					// ignore other message types
				}
			}
		}
	}()

	timeout := 2 * time.Minute

	testLogger.Info().
		Str("expected_log", expectedUserLog).
		Dur("timeout", timeout).
		Msg("Waiting for expected user log message or timeout")

	// Wait for either the expected log to be found, or engine initialization failure to be detected, or timeout (2 minutes)
	select {
	case <-foundExpectedLog:
		testLogger.Info().
			Str("expected_log", expectedUserLog).
			Msg("✅ Test completed successfully - found expected user log message!")
		return
	case <-foundErrorLog:
		require.Fail(t, "Test completed with error - found engine initialization failure message!")
	case <-time.After(timeout):
		testLogger.Error().Msg("Timed out waiting for expected user log message")
		if receivedUserLogs > 0 {
			testLogger.Warn().Int("received_user_logs", receivedUserLogs).Msg("Received some UserLogs messages, but none matched expected log")
		} else {
			testLogger.Warn().Msg("Did not receive any UserLogs messages")
		}
		require.Failf(t, "Timed out waiting for expected user log message", "Expected user log message: '%s' not found after %s", expectedUserLog, timeout.String())
	case err := <-kafkaErrChan:
		testLogger.Error().Err(err).Msg("Kafka listener encountered an error during execution")
		require.Fail(t, "Kafka listener failed: %v", err)
	}

	testLogger.Info().Msg("Workshop test completed")
}

func listenForKafkaMessages(
	ctx context.Context,
	logger zerolog.Logger,
	brokerAddress string,
	topic string,
	messageTypes map[string]func() proto.Message, // ce_type -> protobuf factory function
	messageChan chan<- proto.Message, // channel to send deserialized messages
	errChan chan<- error,
) {
	logger.Info().Str("broker", brokerAddress).Str("topic", topic).Msg("Starting Kafka listener")
	startTime := time.Now()
	logger.Debug().Time("start_time", startTime).Msg("Listener start time - will process messages from this point forward")

	// Configure Kafka consumer to start from latest messages (test start time)
	config := &kafka.ConfigMap{
		"bootstrap.servers":  brokerAddress,
		"group.id":           fmt.Sprintf("workshop-listener-%d", startTime.Unix()), // Unique group per listener
		"auto.offset.reset":  "latest",                                              // Start from latest messages, not earliest
		"session.timeout.ms": 10000,
		"enable.auto.commit": true,             // Commit messages after processing
		"isolation.level":    "read_committed", // Only read committed messages
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to create consumer")
		return
	}
	defer consumer.Close()

	logger.Debug().Msg("Kafka consumer created successfully")

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to subscribe to topic "+topic)
		return
	}

	logger.Info().Str("topic", topic).Msg("Subscribed to topic (consuming from latest offset)")

	ticker := time.NewTicker(100 * time.Millisecond) // Check every 100ms instead of tight loop
	defer ticker.Stop()

	interestedTypes := getMapKeys(messageTypes)
	logger.Info().Strs("interested_types", interestedTypes).Msg("Starting message listening loop")

	// Start consuming messages
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Context cancelled, stopping Kafka listener")
			return
		case <-ticker.C:
			msg, err := consumer.ReadMessage(0) // Non-blocking read
			if err != nil {
				// Check if it's just a timeout (no messages available)
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					// Don't log timeouts as they're expected
					continue
				}
				logger.Error().Err(err).Msg("Consumer error")
				errChan <- errors.Wrap(err, "failed to consume message")
				return
			}

			// Check message timestamp to ensure it's from current listener session
			msgTime := msg.Timestamp
			if !msgTime.IsZero() && msgTime.Before(startTime) {
				logger.Debug().
					Time("msg_time", msgTime).
					Time("start_time", startTime).
					Msg("Skipping old message from before listener start")
				continue
			}

			logger.Debug().
				Str("key", string(msg.Key)).
				Int("value_length", len(msg.Value)).
				Str("topic", *msg.TopicPartition.Topic).
				Int32("partition", msg.TopicPartition.Partition).
				Int64("offset", int64(msg.TopicPartition.Offset)).
				Time("timestamp", msgTime).
				Msg("Received new message")

			ceType, err := getValueFromHeader("ce_type", msg)
			if err != nil {
				logger.Debug().Err(err).Msg("Failed to get ce_type, skipping")
				continue
			}

			logger.Debug().Str("ce_type", ceType).Msg("Message type determined")

			// Check if we're interested in this message type
			factory, interested := messageTypes[ceType]
			if !interested {
				logger.Debug().
					Str("ce_type", ceType).
					Strs("interested_types", interestedTypes).
					Msg("Skipping message type (not in interested types)")
				continue
			}

			// CloudEvents with ce_datacontenttype: application/protobuf
			// The protobuf data starts at offset 6 (after 6-byte binary header)
			const protobufOffset = 6

			if len(msg.Value) <= protobufOffset {
				logger.Debug().
					Int("message_length", len(msg.Value)).
					Int("required_offset", protobufOffset).
					Msg("Message too short for binary-wrapped protobuf")
				continue
			}

			protobufData := msg.Value[protobufOffset:]
			message := factory() // Create new instance using factory function passed in messageTypes map

			err = proto.Unmarshal(protobufData, message)
			if err != nil {
				logger.Error().
					Err(err).
					Int("protobuf_offset", protobufOffset).
					Str("ce_type", ceType).
					Msg("Failed to deserialize protobuf")
				continue
			}

			// Successfully processed the message! Send it back through channel
			logger.Debug().Str("ce_type", ceType).Msg("Successfully deserialized message, sending to channel")

			select {
			case messageChan <- message:
				logger.Debug().Msg("Message sent to channel successfully")
			case <-ctx.Done():
				logger.Info().Msg("Context cancelled while sending message")
				return
			default:
				logger.Warn().Msg("Message channel full, dropping message")
			}
		}
	}
}

// Helper function to get map keys for logging
func getMapKeys(m map[string]func() proto.Message) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getValueFromHeader(expectedHeader string, msg *kafka.Message) (string, error) {
	for _, header := range msg.Headers {
		if string(header.Key) == expectedHeader {
			return string(header.Value), nil
		}
	}
	return "", fmt.Errorf("%s not found in headers", expectedHeader)
}

type WorkflowV2DependenciesConfig struct {
	CronBinaryPath string `toml:"cron_capability_binary_path" validate:"required"`
}

type V2WorkflowTestConfig struct {
	Blockchain       *types.WrappedBlockchainInput `toml:"blockchain" validate:"required"`
	NodeSets         []*ns.Input                   `toml:"nodesets" validate:"required"`
	JD               *jd.Input                     `toml:"jd" validate:"required"`
	Infra            *libtypes.InfraInput          `toml:"infra" validate:"required"`
	Dependencies     *WorkflowV2DependenciesConfig `toml:"dependencies" validate:"required"`
	CustomAnvilMiner *types.CustomAnvilMiner       `toml:"custom_anvil_miner"`
}
