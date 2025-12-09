package helpers

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	defaultEventBufferSize = 200
	defaultPollInterval    = 5 * time.Second
	defaultRetryAttempts   = 3
	defaultRetryDelay      = 1 * time.Second
)

// WorkflowEventMessage wraps a deserialized protobuf event with metadata
type WorkflowEventMessage struct {
	NodeName  string        // Which node emitted this event
	Sequence  int64         // Event sequence number
	Type      string        // Event type (e.g., "workflows.v2.CapabilityExecutionStarted")
	Timestamp time.Time     // Event timestamp
	Event     proto.Message // Deserialized protobuf message
}

// WorkflowEventsSubscriber polls workflow events from multiple Chainlink worker nodes
type WorkflowEventsSubscriber struct {
	workflowID    string
	don           *cre.Don
	workerNodes   []*cre.Node
	pollInterval  time.Duration
	retryAttempts int
	retryDelay    time.Duration

	// Per-node sequence tracking
	nodeSequences map[string]int64 // map[nodeName]lastSequence
	sequenceMu    sync.RWMutex

	logger      zerolog.Logger
	messageChan chan WorkflowEventMessage
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// WorkflowEventsSubscriberConfig configures the subscriber
type WorkflowEventsSubscriberConfig struct {
	WorkflowID    string
	Don           *cre.Don
	PollInterval  time.Duration // How often to poll each node
	Logger        zerolog.Logger
	BufferSize    int           // Message channel buffer size (default: 200)
	RetryAttempts int           // Number of retries on API failure (default: 3)
	RetryDelay    time.Duration // Delay between retries (default: 1s)
}

// StartWorkflowEventsSubscriber creates and starts a workflow events subscriber.
// It spawns goroutines that poll all worker nodes for new events.
// Returns: message channel, context, cancel function, error
func StartWorkflowEventsSubscriber(ctx context.Context, config WorkflowEventsSubscriberConfig) (
	<-chan WorkflowEventMessage,
	context.Context,
	context.CancelFunc,
	error,
) {
	// Apply defaults
	if config.BufferSize == 0 {
		config.BufferSize = defaultEventBufferSize
	}
	if config.PollInterval == 0 {
		config.PollInterval = defaultPollInterval
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = defaultRetryAttempts
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = defaultRetryDelay
	}

	// Get worker nodes using Don.Workers()
	workerNodes, err := config.Don.Workers()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get worker nodes: %w", err)
	}

	config.Logger.Info().
		Int("total_workers", len(workerNodes)).
		Strs("nodes", getNodeNames(workerNodes)).
		Msg("Selected worker nodes for event subscription")

	// Create child context from parent
	subCtx, cancel := context.WithCancel(ctx)
	messageChan := make(chan WorkflowEventMessage, config.BufferSize)

	sub := &WorkflowEventsSubscriber{
		workflowID:    config.WorkflowID,
		don:           config.Don,
		workerNodes:   workerNodes,
		pollInterval:  config.PollInterval,
		retryAttempts: config.RetryAttempts,
		retryDelay:    config.RetryDelay,
		nodeSequences: make(map[string]int64),
		logger:        config.Logger,
		messageChan:   messageChan,
		ctx:           subCtx,
		cancel:        cancel,
	}

	// Initialize sequence tracking for each node
	for _, node := range workerNodes {
		sub.nodeSequences[node.Name] = 0
	}

	// Start polling goroutines (one per node)
	for _, node := range workerNodes {
		sub.wg.Add(1)
		go sub.pollNodeLoop(node)
	}

	// Goroutine to close channel when all pollers exit
	go func() {
		sub.wg.Wait()
		close(messageChan)
		config.Logger.Debug().Msg("All polling goroutines stopped, message channel closed")
	}()

	return messageChan, subCtx, cancel, nil
}

// pollNodeLoop continuously polls a single node for new events
func (s *WorkflowEventsSubscriber) pollNodeLoop(node *cre.Node) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	nodeLogger := s.logger.With().Str("node", node.Name).Logger()
	consecutiveFailures := 0

	for {
		select {
		case <-s.ctx.Done():
			nodeLogger.Debug().Msg("Stopping event polling (context cancelled)")
			return

		case <-ticker.C:
			err := s.fetchAndProcessEvents(node, nodeLogger)
			if err != nil {
				consecutiveFailures++
				nodeLogger.Warn().
					Err(err).
					Int("consecutive_failures", consecutiveFailures).
					Msg("Failed to fetch events")

				// If we've exceeded max retries, cancel the context and exit
				if consecutiveFailures >= s.retryAttempts {
					nodeLogger.Error().
						Int("max_retries", s.retryAttempts).
						Msg("Max consecutive failures reached, cancelling subscription")
					s.cancel()
					return
				}
			} else {
				// Reset failure counter on success
				consecutiveFailures = 0
			}
		}
	}
}

// fetchAndProcessEvents calls the API and processes new events
func (s *WorkflowEventsSubscriber) fetchAndProcessEvents(node *cre.Node, logger zerolog.Logger) error {
	// Get the last sequence for this specific node
	s.sequenceMu.RLock()
	minSequence := s.nodeSequences[node.Name] + 1
	s.sequenceMu.RUnlock()

	// Call ReadWorkflowEvents with retry logic
	var response *clclient.WorkflowDebugEvents
	err := retry.Do(
		func() error {
			var apiErr error
			response, _, apiErr = node.Clients.RestClient.ReadWorkflowEvents(s.workflowID, minSequence, 100)
			return apiErr
		},
		retry.Attempts(uint(s.retryAttempts)),
		retry.Delay(s.retryDelay),
		retry.OnRetry(func(n uint, err error) {
			logger.Debug().
				Uint("attempt", n+1).
				Err(err).
				Msg("Retrying ReadWorkflowEvents")
		}),
	)

	if err != nil {
		return fmt.Errorf("ReadWorkflowEvents failed after %d attempts: %w", s.retryAttempts, err)
	}

	// Extract events from response
	events := response.Data.Attributes.Events

	// Process each event
	for _, eventEntry := range events {
		// Deserialize protobuf based on event type
		protoMsg, err := deserializeWorkflowEvent(eventEntry.Type, eventEntry.Message)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("type", eventEntry.Type).
				Int64("sequence", eventEntry.Sequence).
				Msg("Failed to deserialize event")
			continue
		}

		msg := WorkflowEventMessage{
			NodeName:  node.Name,
			Sequence:  eventEntry.Sequence,
			Type:      eventEntry.Type,
			Timestamp: eventEntry.Timestamp,
			Event:     protoMsg,
		}

		// Send to channel (non-blocking with context check)
		select {
		case s.messageChan <- msg:
			logger.Debug().
				Str("type", eventEntry.Type).
				Int64("sequence", eventEntry.Sequence).
				Msg("Event sent to channel")

			// Update last sequence for this node
			s.sequenceMu.Lock()
			if eventEntry.Sequence > s.nodeSequences[node.Name] {
				s.nodeSequences[node.Name] = eventEntry.Sequence
			}
			s.sequenceMu.Unlock()

		case <-s.ctx.Done():
			return nil

		default:
			logger.Warn().
				Str("type", eventEntry.Type).
				Int64("sequence", eventEntry.Sequence).
				Msg("Channel full, dropping event")
		}
	}

	return nil
}

// MatcherFunc is a function that checks if an event matches criteria
type MatcherFunc func(msg proto.Message) (bool, error)

// AssertWorkflowEventMatched waits for N nodes to emit events matching the criteria.
// This is a BLOCKING function that returns when:
// - N unique nodes have emitted matching events, OR
// - The timeout is reached, OR
// - The subscriber context is cancelled (e.g., due to failures)
// Returns error if conditions not met within timeout.
func AssertWorkflowEventMatched(
	ctx context.Context,
	n int,
	messageChan <-chan WorkflowEventMessage,
	matcher MatcherFunc,
	timeout time.Duration,
	logger zerolog.Logger,
) error {

	matchedNodes := make(map[string]bool) // Track which nodes have matching events
	deadline := time.After(timeout)

	logger.Info().
		Int("required_nodes", n).
		Dur("timeout", timeout).
		Msg("Waiting for workflow event matches across nodes")

	for {
		select {
		case <-ctx.Done():
			// Subscriber context cancelled (likely due to failures)
			return fmt.Errorf(
				"subscriber context cancelled: only %d/%d nodes emitted matching events: %w",
				len(matchedNodes),
				n,
				ctx.Err(),
			)

		case <-deadline:
			// Timeout reached
			return fmt.Errorf(
				"timeout after %v: only %d/%d nodes emitted matching events",
				timeout,
				len(matchedNodes),
				n,
			)

		case msg, ok := <-messageChan:
			if !ok {
				// Channel closed (all pollers stopped)
				return fmt.Errorf(
					"message channel closed: only %d/%d nodes emitted matching events",
					len(matchedNodes),
					n,
				)
			}

			if ok, err := matcher(msg.Event); err != nil {
				logger.Error().
					Err(err).
					Str("type", msg.Type).
					Int64("sequence", msg.Sequence).
					Msg("Matcher function error")

				return fmt.Errorf("Matcher function errored: %w", err)
			} else if ok {
				if !matchedNodes[msg.NodeName] {
					matchedNodes[msg.NodeName] = true

					logger.Info().
						Str("node", msg.NodeName).
						Str("type", msg.Type).
						Int64("sequence", msg.Sequence).
						Int("total_matched_nodes", len(matchedNodes)).
						Msg("Found matching event from node")

					if len(matchedNodes) >= n {
						logger.Info().
							Int("matched_nodes", len(matchedNodes)).
							Strs("nodes", getMapKeys(matchedNodes)).
							Msg("Required number of nodes have matching events")
						return nil
					}
				}
			}
		}
	}
}

func LogWorkflowEvent(
	ctx context.Context,
	messageChan <-chan WorkflowEventMessage,
) {
	if os.Getenv("TEST_WORKFLOW_DEBUG_NO_LOGS") == "true" {
		return
	}

	for {
		select {
		case <-ctx.Done():
			// Subscriber context cancelled (likely due to failures)
			return
		case msg, ok := <-messageChan:
			if !ok {
				// Channel closed (all pollers stopped)
				return
			}

			// Map event type strings to protobuf message types
			switch msg.Type {
			// V1 events
			case "workflows.v1.WorkflowExecutionStarted":
				if asEvent, ok := msg.Event.(*workflowevents.WorkflowExecutionStarted); ok {
					fmt.Printf(" [%s] --------> WorkflowExecutionStarted: %s\n", msg.NodeName, asEvent.M.WorkflowID)
				}
			case "workflows.v1.WorkflowExecutionFinished":
				if asEvent, ok := msg.Event.(*workflowevents.WorkflowExecutionFinished); ok {
					fmt.Printf(" [%s] --------> WorkflowExecutionFinished: %s\n", msg.NodeName, asEvent.M.WorkflowID)
				}
			case "workflows.v1.WorkflowStatusChanged":
				if asEvent, ok := msg.Event.(*workflowevents.WorkflowStatusChanged); ok {
					fmt.Printf(" [%s] --------> WorkflowStatusChanged: %s\n", msg.NodeName, asEvent.Status) // weirdly WorkflowID here is empty!
				}
			case "workflows.v1.UserLogs":
				if asEvent, ok := msg.Event.(*workflowevents.UserLogs); ok {
					for _, line := range asEvent.LogLines {
						re := regexp.MustCompile(`msg=\"(.*?)\"`)
						result := re.FindStringSubmatch(line.Message)
						if len(result) > 1 {
							for _, match := range result[1:] {
								fmt.Printf(" [%s] --------> UserLogs: %s\n", msg.NodeName, match)
							}
						} else {
							fmt.Printf(" [%s] --------> UserLogs: %s\n", msg.NodeName, line.Message)
						}
					}
				}
			// Common events
			case "BaseMessage":
				if asEvent, ok := msg.Event.(*commonevents.BaseMessage); ok {
					fmt.Printf(" [%s] --------> BaseMessage: %s\n", msg.NodeName, asEvent.Msg)
				}
			}
		}
	}
}

// Helper functions

// getNodeNames extracts node names from a list of nodes
func getNodeNames(nodes []*cre.Node) []string {
	names := make([]string, len(nodes))
	for i, node := range nodes {
		names[i] = node.Name
	}
	return names
}

// getMapKeys returns the keys of a map as a slice
func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// deserializeWorkflowEvent deserializes event data based on type
// Complete list from core/web/presenters/debug_workflow_formatter.go
func deserializeWorkflowEvent(eventType string, data []byte) (proto.Message, error) {
	var msg proto.Message

	// Map event type strings to protobuf message types
	switch eventType {
	// V1 events
	case "workflows.v1.WorkflowExecutionStarted":
		msg = &workflowevents.WorkflowExecutionStarted{}
	case "workflows.v1.WorkflowExecutionFinished":
		msg = &workflowevents.WorkflowExecutionFinished{}
	case "workflows.v1.CapabilityExecutionStarted":
		msg = &workflowevents.CapabilityExecutionStarted{}
	case "workflows.v1.CapabilityExecutionFinished":
		msg = &workflowevents.CapabilityExecutionFinished{}
	case "workflows.v1.MeteringReport":
		msg = &workflowevents.MeteringReport{}
	case "workflows.v1.WorkflowStatusChanged":
		msg = &workflowevents.WorkflowStatusChanged{}
	case "workflows.v1.UserLogs":
		msg = &workflowevents.UserLogs{}
	case "workflows.v1.TransmissionsScheduledEvent":
		msg = &workflowevents.TransmissionsScheduledEvent{}

	// Common events
	case "BaseMessage":
		msg = &commonevents.BaseMessage{}

	// V2 events
	case "workflows.v2.WorkflowExecutionStarted":
		msg = &eventsv2.WorkflowExecutionStarted{}
	case "workflows.v2.WorkflowExecutionFinished":
		msg = &eventsv2.WorkflowExecutionFinished{}
	case "workflows.v2.CapabilityExecutionStarted":
		msg = &eventsv2.CapabilityExecutionStarted{}
	case "workflows.v2.CapabilityExecutionFinished":
		msg = &eventsv2.CapabilityExecutionFinished{}
	case "workflows.v2.TriggerExecutionStarted":
		msg = &eventsv2.TriggerExecutionStarted{}
	case "workflows.v2.WorkflowUserLog":
		msg = &eventsv2.WorkflowUserLog{}
	case "workflows.v2.WorkflowActivated":
		msg = &eventsv2.WorkflowActivated{}
	case "workflows.v2.WorkflowPaused":
		msg = &eventsv2.WorkflowPaused{}
	case "workflows.v2.WorkflowDeleted":
		msg = &eventsv2.WorkflowDeleted{}

	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", eventType, err)
	}

	return msg, nil
}

func GetUserLogMatcherFn(expectedMessage string) func(msg proto.Message) (bool, error) {
	return func(msg proto.Message) (bool, error) {
		if log, ok := msg.(*workflowevents.UserLogs); ok {
			for _, line := range log.LogLines {
				if strings.Contains(line.Message, expectedMessage) {
					return true, nil
					// } else {
					// fmt.Printf("Found a different user log: %s\n", line.Message)
				}
			}
		}
		if base, ok := msg.(*commonevents.BaseMessage); ok {
			if strings.Contains(base.Msg, "Workflow Engine initialization failed") {
				return false, errors.New("Workflow Engine initialization failed")
			}
		}
		return false, nil
	}
}

func GetStandardWorkflowEventsSubscriberConfig(testEnv *ttypes.TestEnvironment, workflowID string) WorkflowEventsSubscriberConfig {
	wfDon := testEnv.Dons.MustWorkflowDON()
	return WorkflowEventsSubscriberConfig{
		WorkflowID:   workflowID,
		Don:          wfDon,
		PollInterval: 5 * time.Second,
		Logger:       testEnv.Logger,
	}
}

func FanOutWorkflowEvents(
	ctx context.Context,
	source <-chan WorkflowEventMessage,
	numConsumers int,
) []<-chan WorkflowEventMessage {
	destinations := make([]chan WorkflowEventMessage, numConsumers)
	outputs := make([]<-chan WorkflowEventMessage, numConsumers)

	for i := range numConsumers {
		destinations[i] = make(chan WorkflowEventMessage, 100)
		outputs[i] = destinations[i]
	}

	go func() {
		defer func() {
			// Close all destination channels when done
			for _, dest := range destinations {
				close(dest)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				// Context cancelled, stop fan-out
				return
			case msg, ok := <-source:
				if !ok {
					// Source channel closed
					return
				}
				// Broadcast to all consumers
				for _, dest := range destinations {
					select {
					case dest <- msg:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return outputs
}
