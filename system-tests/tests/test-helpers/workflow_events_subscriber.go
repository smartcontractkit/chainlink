package helpers

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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
	F             int           // Number of worker nodes to subscribe to
	Logger        zerolog.Logger
	BufferSize    int           // Message channel buffer size (default: 200)
	RetryAttempts int           // Number of retries on API failure (default: 3)
	RetryDelay    time.Duration // Delay between retries (default: 1s)
}

// StartWorkflowEventsSubscriber creates and starts a workflow events subscriber.
// It spawns goroutines that poll F randomly-selected worker nodes for new events.
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

	// Validate F parameter
	if config.F <= 0 {
		return nil, nil, nil, fmt.Errorf("F must be greater than 0, got %d", config.F)
	}
	if config.F > len(workerNodes) {
		return nil, nil, nil, fmt.Errorf("F (%d) cannot be greater than number of worker nodes (%d)", config.F, len(workerNodes))
	}

	// Select F random worker nodes
	selectedNodes := selectRandomNodes(workerNodes, config.F)

	config.Logger.Info().
		Int("total_workers", len(workerNodes)).
		Int("selected", len(selectedNodes)).
		Strs("nodes", getNodeNames(selectedNodes)).
		Msg("Selected worker nodes for event subscription")

	// Create child context from parent
	subCtx, cancel := context.WithCancel(ctx)
	messageChan := make(chan WorkflowEventMessage, config.BufferSize)

	sub := &WorkflowEventsSubscriber{
		workflowID:    config.WorkflowID,
		don:           config.Don,
		workerNodes:   selectedNodes,
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
	for _, node := range selectedNodes {
		sub.nodeSequences[node.Name] = 0
	}

	// Start polling goroutines (one per selected node)
	for _, node := range selectedNodes {
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
type MatcherFunc func(msg proto.Message) bool

// AssertWorkflowEventMatched waits for F nodes to emit events matching the criteria.
// This is a BLOCKING function that returns when:
// - F unique nodes have emitted matching events, OR
// - The timeout is reached, OR
// - The subscriber context is cancelled (e.g., due to failures)
// The cancel function will be called when assertion completes to stop the subscriber.
// Returns error if conditions not met within timeout.
func AssertWorkflowEventMatched(
	ctx context.Context,
	cancel context.CancelFunc,
	f int,
	messageChan <-chan WorkflowEventMessage,
	matcher MatcherFunc,
	timeout time.Duration,
	logger zerolog.Logger,
) error {
	// Ensure subscriber is stopped when assertion completes
	defer cancel()

	matchedNodes := make(map[string]bool) // Track which nodes have matching events
	deadline := time.After(timeout)

	logger.Info().
		Int("required_nodes", f).
		Dur("timeout", timeout).
		Msg("Waiting for workflow event matches across nodes")

	for {
		select {
		case <-ctx.Done():
			// Subscriber context cancelled (likely due to failures)
			return fmt.Errorf(
				"subscriber context cancelled: only %d/%d nodes emitted matching events: %w",
				len(matchedNodes),
				f,
				ctx.Err(),
			)

		case <-deadline:
			// Timeout reached
			return fmt.Errorf(
				"timeout after %v: only %d/%d nodes emitted matching events",
				timeout,
				len(matchedNodes),
				f,
			)

		case msg, ok := <-messageChan:
			if !ok {
				// Channel closed (all pollers stopped)
				return fmt.Errorf(
					"message channel closed: only %d/%d nodes emitted matching events",
					len(matchedNodes),
					f,
				)
			}

			if matcher(msg.Event) {
				if !matchedNodes[msg.NodeName] {
					matchedNodes[msg.NodeName] = true

					logger.Info().
						Str("node", msg.NodeName).
						Str("type", msg.Type).
						Int64("sequence", msg.Sequence).
						Int("total_matched_nodes", len(matchedNodes)).
						Msg("Found matching event from node")

					if len(matchedNodes) >= f {
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

// Helper functions

// selectRandomNodes selects up to f random nodes from the list
func selectRandomNodes(nodes []*cre.Node, f int) []*cre.Node {
	if f >= len(nodes) {
		return nodes
	}

	// Fisher-Yates shuffle and take first f
	shuffled := make([]*cre.Node, len(nodes))
	copy(shuffled, nodes)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:f]
}

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
