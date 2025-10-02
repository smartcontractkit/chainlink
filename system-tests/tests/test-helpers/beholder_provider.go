package helpers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

// Constants for configuration
const (
	// Channel buffer sizes
	messageChannelBufferSize = 20
	errorChannelBufferSize   = 1
	channelFullRetryTimeout  = 100 * time.Millisecond

	// Kafka configuration
	kafkaSessionTimeoutMs = 10000
	kafkaReadTimeoutMs    = 0 // Non-blocking read

	// Timing configuration
	messageReadInterval = 50 * time.Millisecond

	// CloudEvents protobuf offset
	protobufOffset = 6

	// Expected CloudEvents header
	ceTypeHeader = "ce_type"
)

type Beholder struct {
	cfg  *config.ChipIngressConfig
	lggr zerolog.Logger
}

func NewBeholder(lggr zerolog.Logger, relativePathToRepoRoot, environmentDir string) (*Beholder, error) {
	err := startBeholderStackIfIsNotRunning(relativePathToRepoRoot, environmentDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure beholder stack is running")
	}

	chipConfig, err := loadBeholderStackCache(relativePathToRepoRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load beholder stack cache")
	}
	return &Beholder{cfg: chipConfig, lggr: lggr}, nil
}

func loadBeholderStackCache(relativePathToRepoRoot string) (*config.ChipIngressConfig, error) {
	c := &config.ChipIngressConfig{}
	if loadErr := c.Load(config.MustChipIngressStateFileAbsPath(relativePathToRepoRoot)); loadErr != nil {
		return nil, errors.Wrap(loadErr, "failed to load beholder stack cache")
	}
	if c.ChipIngress.Output.RedPanda.KafkaExternalURL == "" {
		return nil, errors.New("kafka external url is not set in the cache")
	}

	if len(c.Kafka.Topics) == 0 {
		return nil, errors.New("kafka topics are not set in the cache")
	}

	return c, nil
}

func startBeholderStackIfIsNotRunning(relativePathToRepoRoot, environmentDir string) error {
	if !config.ChipIngressStateFileExists(relativePathToRepoRoot) {
		framework.L.Info().Str("state file", config.MustChipIngressStateFileAbsPath(relativePathToRepoRoot)).Msg("Beholder state file was not found. Starting Beholder...")
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

func (b *Beholder) SubscribeToBeholderMessages(
	ctx context.Context,
	messageTypes map[string]func() proto.Message,
) (<-chan proto.Message, <-chan error) {
	kafkaErrChan := make(chan error, errorChannelBufferSize)
	messageChan := make(chan proto.Message, messageChannelBufferSize)
	readyChan := make(chan bool, 1)

	// Start listening for messages in the background
	go func() {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				b.lggr.Error().Interface("panic", r).Msg("Panic in Kafka listener goroutine")
				select {
				case kafkaErrChan <- errors.Errorf("panic in listener: %v", r):
				default:
				}
			}
		}()

		kafkaURL := b.cfg.ChipIngress.Output.RedPanda.KafkaExternalURL
		topic := b.cfg.Kafka.Topics[0]
		listenForKafkaMessages(ctx, b.lggr, kafkaURL, topic, messageTypes, messageChan, kafkaErrChan, readyChan)
	}()

	// Wait for consumer to be ready before returning channels
	// This ensures proper coordination between consumer readiness and workflow execution
	select {
	case <-readyChan:
		b.lggr.Info().Msg("Kafka consumer is ready and subscribed - safe to start workflow execution")
	case <-time.After(15 * time.Second): // Increased timeout for CI environments
		select {
		case kafkaErrChan <- errors.New("timeout waiting for consumer to be ready"):
		default:
		}
		b.lggr.Error().Msg("Timeout waiting for Kafka consumer to be ready - check broker connectivity")
	case <-ctx.Done():
		b.lggr.Info().Msg("Context cancelled while waiting for consumer readiness")
	}

	return messageChan, kafkaErrChan
}

// Helper function to get map keys for logging
func getMapKeys(m map[string]func() proto.Message) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func listenForKafkaMessages(
	ctx context.Context,
	logger zerolog.Logger,
	brokerAddress string,
	topic string,
	messageTypes map[string]func() proto.Message, // ce_type -> protobuf factory function
	messageChan chan proto.Message, // channel to send deserialized messages
	errChan chan<- error,
	readyChan chan<- bool,
) {
	logger.Info().Str("broker", brokerAddress).Str("topic", topic).Msg("Starting Kafka listener with readiness signaling")

	// Ensure channel is closed when function exits to prevent goroutine leaks
	defer func() {
		close(messageChan)
		logger.Info().Msg("Listener message channel closed")
	}()

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":  brokerAddress,
		"group.id":           fmt.Sprintf("workshop-listener-%d", time.Now().Unix()), // Unique group per listener
		"auto.offset.reset":  "latest",
		"session.timeout.ms": kafkaSessionTimeoutMs,
		"enable.auto.commit": true,             // Commit messages after processing
		"isolation.level":    "read_committed", // Only read committed messages
	}

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to create consumer")
		return
	}
	defer consumer.Close()
	logger.Info().Msg("Kafka consumer created successfully")

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to subscribe to topic "+topic)
		return
	}

	logger.Info().Str("topic", topic).Msg("Subscribed to topic (consuming from latest offset)")

	// Record start time AFTER consumer is ready to avoid race condition
	startTime := time.Now()
	logger.Info().Time("start_time", startTime).Msg("Consumer ready - will process messages from this point forward")

	// Signal that consumer is ready - this is the key improvement for coordination
	select {
	case readyChan <- true:
		logger.Info().Msg("Signaled consumer readiness - workflow execution can now begin safely")
	default:
		logger.Debug().Msg("Ready channel already signaled or closed")
	}

	ticker := time.NewTicker(messageReadInterval)
	defer ticker.Stop()

	interestedTypes := getMapKeys(messageTypes)
	logger.Debug().Strs("interested_types", interestedTypes).Msg("Starting message listening loop")

	// Start consuming messages]
	for {
		select {
		case <-ctx.Done():
			logger.Warn().Msg("Context cancelled, stopping Kafka listener")
			return
		case <-ticker.C:
			msg, err := consumer.ReadMessage(kafkaReadTimeoutMs) // Non-blocking read
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

			// More lenient timestamp filtering - only skip very old messages (30+ seconds)
			msgTime := msg.Timestamp
			const oldMessageThreshold = 30 * time.Second
			if !msgTime.IsZero() && msgTime.Before(startTime.Add(-oldMessageThreshold)) {
				logger.Debug().
					Time("msg_time", msgTime).
					Time("start_time", startTime).
					Dur("old_message_threshold", oldMessageThreshold).
					Msg("Skipping old messages")
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

			ceType, err := getValueFromHeader(ceTypeHeader, msg)
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
				// Channel is full - try with a brief timeout instead of dropping immediately
				select {
				case messageChan <- message:
					logger.Warn().Msg("Message sent to channel after brief delay (channel was full)")
				case <-time.After(channelFullRetryTimeout):
					logger.Error().Msg("Message channel full for too long, dropping message")
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func getValueFromHeader(expectedHeader string, msg *kafka.Message) (string, error) {
	for _, header := range msg.Headers {
		if header.Key == expectedHeader {
			return string(header.Value), nil
		}
	}
	return "", fmt.Errorf("%s not found in headers", expectedHeader)
}
