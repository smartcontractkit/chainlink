package helpers

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	// Channel buffer sizes
	defaultMessageBufferSize = 200
	defaultErrorBufferSize   = 100

	// Kafka timings
	chipIngressStackStartTimeout   = 2 * time.Minute  // timeout for starting Chip Ingress stack
	maxConsumerConnectivityTimeout = 60 * time.Second // max timeout before Kafka consumer reconnection
	kafkaSessionTimeoutMs          = 20000            // keep it high enough to let Chip Ingress stack messages incoming
	messageReadInterval            = 50 * time.Millisecond

	// CloudEvents binary format
	// protobufOffset represents the number of bytes to skip in CloudEvents binary format messages
	// before the protobuf payload begins. This is a CloudEvents specification detail where the
	// first 6 bytes contain CloudEvents metadata in binary content mode.
	protobufOffset = 6

	// CloudEvents header for message type routing
	ceTypeHeader = "ce_type"

	// Error messages
	errChipIngressStackOrConfigNil = "chip ingress stack or config is nil"
)

var chipIngressStackStartupMu sync.Mutex

type ChipIngressStack struct {
	cfg  *config.ChipIngressConfig
	lggr zerolog.Logger
}

// All fields are optional; sensible defaults are applied when nil or empty.
type ConsumerOptions struct {
	GroupID                string // The consumer group to ensure independent message consumption. Defaults to "chip-ingress-stack-consumer".
	Topic                  string // If empty, uses the first topic from config.
	MessageBuffer          int
	ErrorBuffer            int
	CommitSync             bool // Default: true.Enables synchronous commits, safer (guaranteed commit). Async is less safe (potential reprocessing).
	IsolationReadCommitted bool // Ensures only committed messages are read. Defaults to "false".
}

// NewChipIngressStack creates a ChipIngressStack instance, even if it's not already running.
func NewChipIngressStack(lggr zerolog.Logger, testConfig *configuration.TestConfig) (*ChipIngressStack, error) {
	chipIngressStackStartupMu.Lock()
	defer chipIngressStackStartupMu.Unlock()

	// we don't need to pass the GRPC port here, because it will be automatically assigned by Docker
	if err := startChipIngressStackIfNotRunning(testConfig.RelativePathToRepoRoot, testConfig.EnvironmentDirPath); err != nil {
		return nil, errors.Wrap(err, "Chip Ingress stack failed to start")
	}

	chipConfig, err := loadChipIngressStackCache(testConfig.RelativePathToRepoRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load chip ingress stack cache")
	}

	return &ChipIngressStack{cfg: chipConfig, lggr: lggr}, nil
}

// startChipIngressStackIfNotRunning starts the Chip Ingress stack if it's not already running.
func startChipIngressStackIfNotRunning(relativePathToRepoRoot, environmentDir string) error {
	if config.ChipIngressStateFileExists(relativePathToRepoRoot) {
		framework.L.Info().Msg("No need to start Chip Ingress stack - it is already running")
		return nil
	}

	framework.L.Info().Dur("timeout", chipIngressStackStartTimeout).Msg("Chip Ingress stack state file not found. Starting Chip Ingress stack...")
	ctx, cancel := context.WithTimeout(context.Background(), chipIngressStackStartTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".", "env", "chip-ingress-stack", "start")
	cmd.Dir = environmentDir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.Wrap(err, "timeout starting Chip Ingress stack")
		}
		return errors.Wrap(err, "failed to start Chip Ingress stack")
	}

	framework.L.Info().Msg("Chip Ingress stack started successfully")
	return nil
}

func StopChipIngressStack(relativePathToRepoRoot, environmentDir string) error {
	if !config.ChipIngressStateFileExists(relativePathToRepoRoot) {
		framework.L.Info().Msg("No need to stop Chip Ingress stack - it is not running")
		return nil
	}

	framework.L.Info().Dur("timeout", chipIngressStackStartTimeout).Msg("Chip Ingress stack state file found. Stopping Chip Ingress stack...")
	ctx, cancel := context.WithTimeout(context.Background(), chipIngressStackStartTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".", "env", "chip-ingress-stack", "stop")
	cmd.Dir = environmentDir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.Wrap(err, "timeout stopping Chip Ingress stack")
		}
		return errors.Wrap(err, "failed to stop Chip Ingress stack")
	}

	framework.L.Info().Msg("Chip Ingress stack stopped successfully")
	return nil
}

// loadChipIngressStackCache loads and validates the Chip Ingress stack configuration.
func loadChipIngressStackCache(relativePathToRepoRoot string) (*config.ChipIngressConfig, error) {
	c := &config.ChipIngressConfig{}
	if err := c.Load(config.MustChipIngressStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return nil, errors.Wrap(err, "load cache")
	}

	if c.ChipIngress.Output.RedPanda.KafkaExternalURL == "" {
		return nil, errors.New("kafka external url not set in cache")
	}

	if len(c.Kafka.Topics) == 0 {
		return nil, errors.New("kafka topics not set in cache")
	}

	return c, nil
}

/*
SubscribeToChipIngressStackMessages starts a Kafka consumer and returns message/error channels.

1. Tests Kafka broker connectivity before starting the listener (FATAL - fails fast if not accessible)
2. Validates Chip Ingress stack heartbeat to ensure it's alive and healthy (FATAL - fails fast if not detected)
3. Validates topic existence and accessibility during subscription
4. Verifies topic metadata and partition availability
5. Coordinates consumer readiness to prevent race conditions with producers

Parameters:
  - ctx: Context for lifecycle management
  - messageTypes: Map of CloudEvents ce_type to protobuf factory functions

Returns:
  - Message channel (closed when consumer stops)
  - Error channel (buffered, reports fatal errors)
*/
func (b *ChipIngressStack) SubscribeToChipIngressStackMessages(ctx context.Context, messageTypes map[string]func() proto.Message,
) (<-chan proto.Message, <-chan error) {
	// If the Chip Ingress stack is not initialized, return an error channel
	if b == nil || b.cfg == nil {
		errCh := make(chan error, 1)
		errCh <- errors.New(errChipIngressStackOrConfigNil)
		close(errCh)
		return nil, errCh
	}

	// Create options internally with unique group ID (to enable tests parallelization)
	opts := &ConsumerOptions{
		GroupID:                fmt.Sprintf("chip-ingress-stack-consumer-%d", time.Now().UnixNano()),
		Topic:                  b.cfg.Kafka.Topics[0],
		MessageBuffer:          defaultMessageBufferSize,
		ErrorBuffer:            defaultErrorBufferSize,
		CommitSync:             true, // guaranteed Kafka acknowledgment
		IsolationReadCommitted: false,
	}

	msgCh := make(chan proto.Message, opts.MessageBuffer)
	errCh := make(chan error, opts.ErrorBuffer)
	readyCh := make(chan struct{}, 1)

	// Pre-flight validation: Kafka connectivity (fatal - fail early if Kafka is not accessible)
	b.lggr.Debug().Msg("Performing Kafka connectivity validation...")
	if err := b.validateConsumerConnectivity(ctx); err != nil {
		b.lggr.Error().Err(err).Msg("Kafka connectivity validation failed")
		errCh <- errors.Wrap(err, "kafka connectivity validation failed")
		close(errCh)
		close(msgCh)
		return msgCh, errCh
	}

	// Pre-flight validation: Chip Ingress stack heartbeat (fatal - fail early if stack is not healthy)
	b.lggr.Debug().Msg("Performing Chip Ingress stack heartbeat validation...")
	if err := b.validateChipIngressStackHeartbeat(ctx); err != nil {
		b.lggr.Error().Err(err).Msg("Chip Ingress stack heartbeat validation failed")
		errCh <- errors.Wrap(err, "chip ingress stack heartbeat validation failed")
		close(errCh)
		close(msgCh)
		return msgCh, errCh
	}

	// Start consumer in background  and wait for consumer readiness to coordinate with producers/workflows
	go b.consume(ctx, messageTypes, opts, msgCh, errCh, readyCh)
	select {
	case <-readyCh:
		b.lggr.Info().Msg("Kafka consumer ready and subscribed - safe to start workflow execution")
	case <-time.After(maxConsumerConnectivityTimeout): // Increased timeout for CI environments
		select {
		case errCh <- errors.New("timeout waiting for Kafka consumer readiness"):
		default:
		}
		b.lggr.Error().Msg("Timeout waiting for Kafka consumer readiness - check broker connectivity")
	case <-ctx.Done():
		b.lggr.Info().Msg("Context cancelled while waiting for Kafka consumer readiness")
	}

	return msgCh, errCh
}

// validateKafkaConnectivity explicitly validates Kafka broker connectivity.
func (b *ChipIngressStack) validateConsumerConnectivity(ctx context.Context) error {
	vctx, cancel := context.WithTimeout(ctx, maxConsumerConnectivityTimeout)
	defer cancel()

	consumer, err := b.createValidationConsumer(vctx, "validation")
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := consumer.Close(); closeErr != nil {
			b.lggr.Warn().Err(closeErr).Msg("Failed to close validation consumer")
		}
	}()

	topic := b.cfg.Kafka.Topics[0]
	if _, err := b.validateTopicMetadata(consumer, topic); err != nil {
		return err
	}

	b.lggr.Info().
		Str("broker", b.cfg.ChipIngress.Output.RedPanda.KafkaExternalURL).
		Str("topic", topic).
		Msg("Kafka connectivity validation successful")
	return nil
}

// validateChipIngressStackHeartbeat validates that the Chip Ingress stack is alive and sending heartbeat messages.
// Retries up to 3 times with a fixed 5-second delay between attempts.
func (b *ChipIngressStack) validateChipIngressStackHeartbeat(ctx context.Context) error {
	const (
		maxRetries = 3
		retryDelay = 5 * time.Second
	)

	b.lggr.Info().
		Int("max_retries", maxRetries).
		Dur("retry_delay", retryDelay).
		Dur("max_timeout", maxConsumerConnectivityTimeout).
		Int("session_timeout_ms", kafkaSessionTimeoutMs).
		Msg("Starting Chip Ingress stack heartbeat validation...")

	return retry.Do(
		func() error {
			return b.validateChipIngressStackHeartbeatOnce(ctx)
		},
		retry.Context(ctx),
		retry.Attempts(maxRetries),
		retry.Delay(retryDelay),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			b.lggr.Warn().
				Err(err).
				Uint("attempt", n+1).
				Uint("max_retries", maxRetries).
				Dur("retry_delay", retryDelay).
				Msg("Chip Ingress stack heartbeat validation attempt failed, retrying...")
		}),
	)
}

// validateChipIngressStackHeartbeatOnce performs a single heartbeat validation attempt.
func (b *ChipIngressStack) validateChipIngressStackHeartbeatOnce(ctx context.Context) error {
	hctx, cancel := context.WithTimeout(ctx, maxConsumerConnectivityTimeout)
	defer cancel()

	consumer, err := b.createValidationConsumer(hctx, "heartbeat-validation")
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := consumer.Close(); closeErr != nil {
			b.lggr.Warn().Err(closeErr).Msg("Failed to close heartbeat validation consumer")
		}
	}()
	b.lggr.Info().Msg("Created consumer for heartbeat validation")

	// Use blocking ReadMessage with timeout instead of ticker pattern
	for {
		select {
		case <-hctx.Done():
			return errors.New("timeout waiting for Chip Ingress stack heartbeat")
		default:
		}

		msg, err := consumer.ReadMessage(messageReadInterval)
		if err != nil {
			// Benign timeout - no messages available yet
			var kerr kafka.Error
			if errors.As(err, &kerr) && kerr.Code() == kafka.ErrTimedOut {
				continue
			}
			b.lggr.Error().Int("session_timeout_ms", kafkaSessionTimeoutMs).Err(err).Msg("Failed to read message during heartbeat validation (consider increasing the session timeout)")
			return errors.Wrap(err, "failed to read message during heartbeat validation")
		}

		// Check if this is a BaseMessage
		ceType, ok := getHeaderValue(ceTypeHeader, msg)
		if !ok || ceType != "BaseMessage" {
			continue
		}

		// Validate message length for CloudEvents binary format
		if len(msg.Value) <= protobufOffset {
			continue
		}

		// Unmarshal BaseMessage
		baseMsg := &commonevents.BaseMessage{}
		if err := proto.Unmarshal(msg.Value[protobufOffset:], baseMsg); err != nil {
			b.lggr.Debug().Err(err).Msg("Failed to unmarshal BaseMessage during heartbeat validation")
			continue
		}

		// Check if this is a heartbeat message
		if !isHeartbeatMessage(baseMsg) {
			b.lggr.Debug().
				Str("msg", baseMsg.Msg).
				Msg("Received BaseMessage but not a heartbeat; continuing to listen")
			continue
		}

		// Found heartbeat!
		b.lggr.Info().
			Str("msg", baseMsg.Msg).
			Interface("labels", baseMsg.Labels).
			Msg("Chip Ingress stack heartbeat detected successfully")
		return nil
	}
}

// createValidationConsumer creates a temporary Kafka consumer for validation purposes.
func (b *ChipIngressStack) createValidationConsumer(ctx context.Context, groupIDPrefix string) (*kafka.Consumer, error) {
	if b == nil || b.cfg == nil {
		return nil, errors.New(errChipIngressStackOrConfigNil)
	}

	groupID := fmt.Sprintf("%s-%d", groupIDPrefix, time.Now().UnixNano())
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":  b.cfg.ChipIngress.Output.RedPanda.KafkaExternalURL,
		"group.id":           groupID,
		"auto.offset.reset":  "latest",
		"session.timeout.ms": kafkaSessionTimeoutMs,
	}

	consumer, err := b.createAndSubscribeConsumer(cfg, b.cfg.Kafka.Topics[0])
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// isHeartbeatMessage checks if a BaseMessage is a Chip Ingress stack heartbeat.
// Heartbeat format: msg="heartbeat" and labels.system="Application"
func isHeartbeatMessage(msg *commonevents.BaseMessage) bool {
	if msg == nil {
		return false
	}

	if msg.Msg != "heartbeat" {
		return false
	}

	if msg.Labels == nil {
		return false
	}

	systemLabel, exists := msg.Labels["system"]
	if !exists {
		return false
	}

	// Case-insensitive comparison for robustness
	return strings.EqualFold(systemLabel, "Application")
}

// consume runs the Kafka consumer loop with offset management and automatic reconnection.
func (b *ChipIngressStack) consume(
	ctx context.Context,
	messageTypes map[string]func() proto.Message,
	opts *ConsumerOptions,
	out chan proto.Message,
	errCh chan<- error,
	readyCh chan<- struct{},
) {
	defer close(out)

	// Exponential backoff
	backoff := 2 * time.Second
	maxBackoffTimeout := 30 * time.Second
	backoffFactor := 2.0
	attempt := 0
	// Reconnection loop
	for {
		select {
		case <-ctx.Done():
			b.lggr.Info().Msg("Context cancelled; exiting Kafka consumer loop")
			return
		default:
			// Continue to connection attempt
		}

		err := b.consumeWithReconnect(ctx, messageTypes, opts, out, errCh, readyCh)
		if err == nil {
			// Clean exit (context cancelled)
			return
		}

		// Calculate backoff with jitter
		attempt++
		jitter := time.Duration(rand.Float64() * float64(backoff) * 0.1) // 10% jitter
		sleepDuration := min(backoff+jitter, maxBackoffTimeout)

		b.lggr.Warn().
			Dur("backoff", sleepDuration).
			Int("attempt", attempt).
			Err(err).
			Msg("Reconnecting Kafka consumer with exponential backoff...")

		select {
		case <-ctx.Done():
			b.lggr.Info().Msg("Context cancelled while attempting to reconnect Kafka consumer")
			return
		case <-time.After(sleepDuration):
			b.lggr.Info().Int("attempt", attempt).Msg("Attempting to reconnect Kafka consumer...")
			// Increase backoff for next iteration
			backoff = min(time.Duration(float64(backoff)*backoffFactor), maxBackoffTimeout)
		}
	}
}

// consumeWithReconnect runs a single consumer session with UserLogs timeout tracking.
func (b *ChipIngressStack) consumeWithReconnect(
	ctx context.Context,
	messageTypes map[string]func() proto.Message,
	opts *ConsumerOptions,
	out chan proto.Message,
	errCh chan<- error,
	readyCh chan<- struct{},
) error {
	// The isolation level determines which messages the Kafka consumer is allowed to read:
	// - [used by default] "read_uncommitted": The consumer can read all messages.
	// - [Chip Ingress stack does not use Kafka transactions] "read_committed": The consumer will only read messages from transactions that have been successfully committed, ensuring no uncommitted or aborted data is delivered.
	// This setting is important for applications that require strong data consistency and want to avoid processing uncommitted or potentially rolled-back messages.
	isolationLevel := "read_uncommitted"
	if opts.IsolationReadCommitted { // false by default
		isolationLevel = "read_committed"
	}

	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        b.cfg.ChipIngress.Output.RedPanda.KafkaExternalURL,
		"group.id":                 opts.GroupID,
		"auto.offset.reset":        "latest", // Only process new messages by default
		"session.timeout.ms":       kafkaSessionTimeoutMs,
		"enable.auto.commit":       false, // Manual commit for safety. Prevents premature commit.
		"enable.auto.offset.store": false, // Explicit commit control
		"isolation.level":          isolationLevel,
	}

	consumer, err := b.createAndSubscribeConsumer(cfg, opts.Topic)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := consumer.Close(); closeErr != nil {
			b.lggr.Warn().Err(closeErr).Msg("Failed to close Kafka consumer")
		}
	}()
	b.lggr.Info().Msg("Kafka consumer created successfully")

	// Verify and log subscription details
	if err := b.logSubscriptionInfo(consumer, opts, errCh); err != nil {
		return err
	}

	// Verify topic accessibility and log consumer ready
	if err := b.validateConsumerReadiness(consumer, opts, errCh); err != nil {
		return err
	}

	// This code signals (in a non-blocking way) that the Kafka consumer is ready to receive messages.
	// It attempts to send an empty struct to the readyCh channel, but if the channel is full, it does nothing.
	select {
	case readyCh <- struct{}{}:
	default:
		b.lggr.Warn().Msg("Kafka consumer readiness already signaled")
	}

	interestedTypes := getMessageTypeKeys(messageTypes)
	b.lggr.Debug().Strs("interested_types", interestedTypes).Msg("Starting message listening loop")

	// Recreate the timer on each UserLogs message to avoid timer reset race conditions
	timeoutTimer := time.NewTimer(maxConsumerConnectivityTimeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			b.lggr.Info().Msg("Context cancelled; exiting consumer loop")
			return nil

		case <-timeoutTimer.C:
			// No UserLogs received within the timeout period
			b.lggr.Warn().
				Dur("timeout", maxConsumerConnectivityTimeout).
				Msg("No UserLogs received within timeout period, triggering Kafka consumer reconnection")
			return errors.New("no UserLogs received within timeout period")

		default:
			// Use blocking ReadMessage with short timeout
			msg, err := consumer.ReadMessage(messageReadInterval)
			if err != nil {
				// Benign timeout - no messages available
				var kerr kafka.Error
				if errors.As(err, &kerr) && kerr.Code() == kafka.ErrTimedOut {
					continue
				}
				logError(b.lggr, errCh, errors.Wrap(err, "failed to read message"))
				return err
			}

			b.lggr.Debug().
				Str("key", string(msg.Key)).
				Int("value_length", len(msg.Value)).
				Int32("partition", msg.TopicPartition.Partition).
				Int64("offset", int64(msg.TopicPartition.Offset)).
				Time("timestamp", msg.Timestamp).
				Msg("Received Kafka message")

			// Extract and validate ce_type header
			ceType, ok := getHeaderValue(ceTypeHeader, msg)
			if !ok {
				b.lggr.Debug().
					Int64("offset", int64(msg.TopicPartition.Offset)).
					Msg("Skipping message without ce_type header")
				continue
			}
			b.lggr.Debug().
				Str("ce_type", ceType).
				Int64("offset", int64(msg.TopicPartition.Offset)).
				Int32("partition", msg.TopicPartition.Partition).
				Msg("Message type determined")

			// Check if we're interested in this message type
			factory, interested := messageTypes[ceType]
			if !interested {
				b.lggr.Debug().
					Str("ce_type", ceType).
					Int64("offset", int64(msg.TopicPartition.Offset)).
					Strs("interested_types", interestedTypes).
					Msg("Skipping other (uninterested) message type")
				continue
			}

			// Validate message length for CloudEvents binary format
			if len(msg.Value) <= protobufOffset {
				b.lggr.Debug().
					Int("len", len(msg.Value)).
					Int("offset", protobufOffset).
					Msg("Message too short for protobuf payload; skipping")
				continue
			}

			// Create and unmarshal protobuf message
			pm := factory()
			if pm == nil {
				b.lggr.Warn().Str("ce_type", ceType).Msg("Factory returned nil; skipping")
				continue
			}

			if err := proto.Unmarshal(msg.Value[protobufOffset:], pm); err != nil {
				b.lggr.Error().Err(err).Str("ce_type", ceType).Msg("Failed to unmarshal protobuf; skipping")
				continue
			}

			// Reset timeout if we received a UserLogs message
			if _, isUserLogs := pm.(*workflowevents.UserLogs); isUserLogs {
				// Go recommendation: don't reuse timers, create new ones to avoid race conditions
				if !timeoutTimer.Stop() {
					// Timer already fired, drain it to prevent blocking
					<-timeoutTimer.C
				}
				timeoutTimer = time.NewTimer(maxConsumerConnectivityTimeout)
				b.lggr.Info().
					Int64("offset", int64(msg.TopicPartition.Offset)).
					Int32("partition", msg.TopicPartition.Partition).
					Dur("timeout", maxConsumerConnectivityTimeout).
					Msg("UserLogs received - reconnection timeout reset")
			}

			// Send to output channel (blocking to prevent message loss)
			select {
			case out <- pm:
				// Commit offset after successful delivery
				if err := b.commitMessage(consumer, msg, opts.CommitSync); err != nil {
					logError(b.lggr, errCh, err)
					return err
				}
				b.lggr.Debug().
					Str("ce_type", ceType).
					Int64("offset", int64(msg.TopicPartition.Offset)).
					Int32("partition", msg.TopicPartition.Partition).
					Msg("Successfully processed and committed message")

			case <-ctx.Done():
				b.lggr.Info().Msg("Context cancelled while delivering message")
				return nil
			}
		}
	}
}

// commitMessage tells Kafka that the message processed and offset committed
// without commits, if the consumer crashes and restarts, it would re-read all messages from the beginning
//
// Two commit modes:
// 1. Synchronous - SLOWER but SAFER: StoreMessage + Commit. Guarantees the offset is persisted before continuing.
// 2. Asynchronous - FASTER but less safe: CommitMessage, don't wait for confirmation from Kafka, offset commit happens in the background
//
// Default: false.Re-reading some messages on crash/restart is acceptable
func (b *ChipIngressStack) commitMessage(consumer *kafka.Consumer, msg *kafka.Message, syncCommit bool) error {
	if syncCommit {
		// Synchronous: Store offset first, then commit synchronously
		if _, err := consumer.StoreMessage(msg); err != nil {
			return errors.Wrap(err, "store offset")
		}
		if _, err := consumer.Commit(); err != nil {
			return errors.Wrap(err, "commit sync")
		}
	} else {
		// Asynchronous: One-step fire-and-forget (stores + commits in one call)
		if _, err := consumer.CommitMessage(msg); err != nil {
			return errors.Wrap(err, "commit message async")
		}
	}

	return nil
}

// getHeaderValue extracts a header value from a Kafka message.
func getHeaderValue(key string, msg *kafka.Message) (string, bool) {
	for _, h := range msg.Headers {
		if h.Key == key {
			return string(h.Value), true
		}
	}
	return "", false
}

// getMessageTypeKeys returns the keys from the message types map for logging.
func getMessageTypeKeys(m map[string]func() proto.Message) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// createAndSubscribeConsumer creates a Kafka consumer and subscribes to a topic.
func (b *ChipIngressStack) createAndSubscribeConsumer(cfg *kafka.ConfigMap, topic string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		b.lggr.Error().Err(err).Msg("failed to create Kafka consumer")
		return nil, errors.Wrap(err, "failed to create Kafka consumer")
	}

	// Use SubscribeTopics for future multi-topic support
	if err := consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		if closeErr := consumer.Close(); closeErr != nil {
			b.lggr.Warn().Err(closeErr).Msg("Failed to close consumer after subscription failure")
		}
		b.lggr.Error().Err(err).Str("topic", topic).Msg("failed to subscribe to topic")
		return nil, errors.Wrapf(err, "failed to subscribe to topic %q", topic)
	}

	return consumer, nil
}

// logSubscriptionInfo fetches and logs subscription and partition assignment details.
func (b *ChipIngressStack) logSubscriptionInfo(consumer *kafka.Consumer, opts *ConsumerOptions, errCh chan<- error) error {
	// Verify subscription by fetching from consumer
	subscription, subErr := consumer.Subscription()
	if subErr != nil {
		logError(b.lggr, errCh, errors.Wrap(subErr, "failed to get subscription info"))
		return subErr
	}

	// Get partition assignment (may be empty initially, will be assigned after first poll)
	assignment, assignErr := consumer.Assignment()
	if assignErr != nil {
		b.lggr.Debug().Err(assignErr).Msg("Could not get partition assignment yet (will be assigned after first poll)")
	}

	logEvent := b.lggr.Info().
		Strs("subscribed_topics", subscription).
		Str("group_id", opts.GroupID)

	if len(assignment) > 0 {
		partitions := getPartitionFromAssignment(assignment)
		logEvent.Ints("assigned_partitions", partitions)
	}

	logEvent.Msg("Kafka consumer subscribed successfully")
	return nil
}

// validateConsumerReadiness verifies topic accessibility and logs consumer ready status.
func (b *ChipIngressStack) validateConsumerReadiness(consumer *kafka.Consumer, opts *ConsumerOptions, errCh chan<- error) error {
	// Get topic metadata to verify accessibility
	md, err := b.validateTopicMetadata(consumer, opts.Topic)
	if err != nil {
		logError(b.lggr, errCh, err)
		return err
	}

	// Log consumer ready with partition count
	b.logConsumerReady(consumer, opts, len(md.Topics[opts.Topic].Partitions))
	return nil
}

// validateTopicMetadata fetches topic metadata and validates accessibility.
func (b *ChipIngressStack) validateTopicMetadata(consumer *kafka.Consumer, topic string) (*kafka.Metadata, error) {
	md, err := consumer.GetMetadata(&topic, false, int(maxConsumerConnectivityTimeout/time.Millisecond))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get metadata")
	}

	if md == nil {
		return nil, errors.New("metadata is nil")
	}

	// Safely check if topic exists in metadata
	topicMd, exists := md.Topics[topic]
	if !exists {
		return nil, errors.Errorf("topic %q not found in metadata", topic)
	}

	// Validate topic error code and partitions
	if topicMd.Error.Code() != kafka.ErrNoError {
		return nil, errors.Errorf("topic %q has error: %v", topic, topicMd.Error)
	}

	if len(topicMd.Partitions) == 0 {
		return nil, errors.Errorf("topic %q has no partitions", topic)
	}

	return md, nil
}

// logConsumerReady logs consumer ready status with subscription and partition details.
func (b *ChipIngressStack) logConsumerReady(consumer *kafka.Consumer, opts *ConsumerOptions, totalPartitions int) {
	// Get updated partition assignment after metadata verification
	subscription, _ := consumer.Subscription()
	assignment, _ := consumer.Assignment()

	readyLog := b.lggr.Info().
		Strs("subscribed_topics", subscription).
		Str("group_id", opts.GroupID).
		Int("total_partitions", totalPartitions)

	if len(assignment) > 0 {
		partitions := getPartitionFromAssignment(assignment)
		readyLog.Ints("assigned_partitions", partitions)
	}

	readyLog.Msg("Consumer ready")
}

// getPartitionFromAssignment extracts partition numbers from TopicPartition slice.
func getPartitionFromAssignment(assignment []kafka.TopicPartition) []int {
	partitions := make([]int, len(assignment))
	for i, tp := range assignment {
		partitions[i] = int(tp.Partition)
	}
	return partitions
}

// logError logs an error and attempts to send it to the error channel.
// If the error channel is full (i.e., the send would block), it silently skips sending
// to avoid blocking the caller. This is useful in goroutines where you want to report
// errors but not risk deadlock if the channel is not being drained.
func logError(l zerolog.Logger, errCh chan<- error, err error) {
	l.Error().Err(err).Msg("Kafka consumer error")
	select {
	case errCh <- err:
		// Error sent to channel
	default:
		// Channel full, skip sending to avoid blocking
	}
}
