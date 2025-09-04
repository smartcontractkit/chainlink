package cre

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func loadBeholderStackCache(beholderConfigPath string) (*config.ChipIngressConfig, error) {
	originalCtfConfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", originalCtfConfigs)
		if setErr != nil {
			framework.L.Error().Err(setErr).Msg("failed to restore CTF_CONFIGS env var")
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", beholderConfigPath)
	if setErr != nil {
		return nil, errors.Wrap(setErr, "failed to set CTF_CONFIGS environment variable")
	}

	return framework.Load[config.ChipIngressConfig](nil)
}

func startBeholderStackIfIsNotRunning(stateFile, environmentDir string) error {
	split := strings.Split(stateFile, ",")
	if _, err := os.Stat(split[0]); os.IsNotExist(err) {
		framework.L.Info().Msg("Beholder has not been found. Starting Beholder...")
		cmd := exec.Command("go", "run", ".", "env", "beholder", "start")
		cmd.Dir = environmentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			return errors.Wrap(cmdErr, "failed to start Beholder")
		}
	}
	framework.L.Info().Msg("Beholder is running.")
	return nil
}

func subscribeToBeholderMessages(ctx context.Context, t *testing.T, testLogger zerolog.Logger,
	testEnv *TestEnvironment,
	messageTypes map[string]func() proto.Message) (<-chan proto.Message, <-chan error) {
	bErr := startBeholderStackIfIsNotRunning(testEnv.TestConfig.BeholderConfigPath, testEnv.TestConfig.EnvironmentDirPath)
	require.NoError(t, bErr, "failed to start Beholder")

	chipConfig, chipErr := loadBeholderStackCache(testEnv.TestConfig.BeholderConfigPath)
	require.NoError(t, chipErr, "failed to load chip ingress cache")
	require.NotNil(t, chipConfig.ChipIngress.Output.RedPanda.KafkaExternalURL, "kafka external url is not set in the cache")
	require.NotEmpty(t, chipConfig.Kafka.Topics, "kafka topics are not set in the cache")

	kafkaErrChan := make(chan error, 1)
	messageChan := make(chan proto.Message, 10)

	// Start listening for messages in the background
	go func() {
		listenForKafkaMessages(ctx, testLogger, chipConfig.ChipIngress.Output.RedPanda.KafkaExternalURL, chipConfig.Kafka.Topics[0], messageTypes, messageChan, kafkaErrChan)
	}()

	return messageChan, kafkaErrChan
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
				var kafkaErr kafka.Error
				if errors.As(err, &kafkaErr) && kafkaErr.Code() == kafka.ErrTimedOut {
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
		if header.Key == expectedHeader {
			return string(header.Value), nil
		}
	}
	return "", fmt.Errorf("%s not found in headers", expectedHeader)
}
