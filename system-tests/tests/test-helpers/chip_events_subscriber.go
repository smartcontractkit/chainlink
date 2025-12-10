package helpers

// ChipEventsSubscriber provides efficient event polling from a standalone chip test sink server.
//
// Key features:
// - Events are automatically enriched with workflow ID and node P2P ID as they arrive
// - WorkflowID and NodeP2PID fields are populated in each ChipEventMessage
// - No need to re-parse events - metadata is extracted once during polling
// - Simple channel-based consumption pattern
//
// Architecture:
// - Single polling goroutine fetches new events from the chip server
// - Each event is enriched with extracted metadata (workflow ID, node P2P ID)
// - Events are sent to a channel for streaming consumption
// - Clients can group events by WorkflowID/NodeP2PID as needed using the populated fields

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

const (
	defaultEventBufferSize = 200
	defaultPollInterval    = 2 * time.Second
	defaultRetryAttempts   = 3
	defaultRetryDelay      = 1 * time.Second
	defaultEventLimit      = 100
	defaultBatchSize       = 20
)

// ChipEventMessage wraps an event from the chip test sink with metadata
type ChipEventMessage struct {
	Sequence uint64         `json:"sequence"`
	Domain   string         `json:"domain"`
	Entity   string         `json:"entity"`
	Schema   string         `json:"schema"`
	Body     string         `json:"body"` // base64 encoded
	Attrs    map[string]any `json:"attrs"`
	// Extracted metadata (populated by subscriber if available)
	WorkflowID string `json:"workflowId,omitempty"` // Extracted workflow ID
	NodeP2PID  string `json:"nodeP2pId,omitempty"`  // Extracted node P2P ID
}

// ChipEventsSubscriber polls events from the standalone chip test sink server
type ChipEventsSubscriber struct {
	chipServerURL string // URL of the chip test sink HTTP API (e.g., "http://localhost:8080")
	pollInterval  time.Duration
	retryAttempts int
	retryDelay    time.Duration
	eventLimit    int // Maximum events to fetch per poll
	batchSize     int

	// Sequence tracking
	lastSequence uint64
	sequenceMu   sync.RWMutex

	logger      zerolog.Logger
	messageChan chan ChipEventMessage
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	httpClient  *http.Client
}

// ChipEventsSubscriberConfig configures the subscriber
type ChipEventsSubscriberConfig struct {
	ChipServerURL string        // URL of the chip test sink HTTP API (e.g., "http://localhost:8080")
	PollInterval  time.Duration // How often to poll for events
	Logger        zerolog.Logger
	BufferSize    int           // Message channel buffer size (default: 200)
	RetryAttempts int           // Number of retries on API failure (default: 3)
	RetryDelay    time.Duration // Delay between retries (default: 1s)
	EventLimit    int           // Maximum number of events to fetch per poll (default: 100)
	BatchSize     int           // Number of events to fetch per batch
}

// getCurrentSequence fetches the latest sequence number from the chip server
// Returns 0 if no events exist yet, or the highest sequence number
func getCurrentSequence(ctx context.Context, chipServerURL string, httpClient *http.Client) (uint64, error) {
	url := fmt.Sprintf("%s/sequence", chipServerURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Sequence uint64 `json:"sequence"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	return result.Sequence, nil
}

// StartChipEventsSubscriber creates and starts a chip events subscriber.
// It spawns a goroutine that polls the chip test sink server for new events.
// Events are automatically enriched with workflow ID and node P2P ID as they arrive.
//
// Usage example:
//
//	msgChan, ctx, cancel, err := StartChipEventsSubscriber(parentCtx, ChipEventsSubscriberConfig{
//	    ChipServerURL: "http://localhost:8080",
//	    PollInterval:  5 * time.Second,
//	    Logger:        logger,
//	})
//	if err != nil {
//	    return err
//	}
//	defer cancel()
//
//	for event := range msgChan {
//	    // WorkflowID and NodeP2PID are already populated
//	    fmt.Printf("Workflow: %s, Node: %s\n", event.WorkflowID, event.NodeP2PID)
//	}
//
// Returns: message channel, context, cancel function, error
func StartChipEventsSubscriber(ctx context.Context, config ChipEventsSubscriberConfig) (
	<-chan ChipEventMessage,
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
	if config.EventLimit == 0 {
		config.EventLimit = defaultEventLimit
	}
	if config.BatchSize == 0 {
		config.BatchSize = defaultBatchSize
	}

	if config.ChipServerURL == "" {
		return nil, nil, nil, fmt.Errorf("ChipServerURL is required")
	}

	config.Logger.Info().
		Str("chip_server", config.ChipServerURL).
		Dur("poll_interval", config.PollInterval).
		Msg("Starting chip events subscriber")

	// Create child context from parent
	subCtx, cancel := context.WithCancel(ctx)
	messageChan := make(chan ChipEventMessage, config.BufferSize)

	// Get current sequence to avoid reading old events from previous test runs
	httpClient := &http.Client{Timeout: 30 * time.Second}
	currentSequence, err := getCurrentSequence(subCtx, config.ChipServerURL, httpClient)
	if err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("failed to get current sequence: %w", err)
	}

	config.Logger.Info().
		Uint64("starting_sequence", currentSequence).
		Msg("Subscriber will start from current sequence to avoid stale events")

	sub := &ChipEventsSubscriber{
		chipServerURL: config.ChipServerURL,
		pollInterval:  config.PollInterval,
		retryAttempts: config.RetryAttempts,
		retryDelay:    config.RetryDelay,
		eventLimit:    config.EventLimit,
		batchSize:     config.BatchSize,
		lastSequence:  currentSequence,
		logger:        config.Logger,
		messageChan:   messageChan,
		ctx:           subCtx,
		cancel:        cancel,
		httpClient:    httpClient,
	}

	// Start polling goroutine
	sub.wg.Add(1)
	go sub.pollLoop()

	// Goroutine to close channel when poller exits
	go func() {
		sub.wg.Wait()
		close(messageChan)
		config.Logger.Debug().Msg("Polling goroutine stopped, message channel closed")
	}()

	return messageChan, subCtx, cancel, nil
}

// pollLoop continuously polls the chip test sink server for new events
func (s *ChipEventsSubscriber) pollLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	consecutiveFailures := 0

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Debug().Msg("Stopping event polling (context cancelled)")
			return

		case <-ticker.C:
			err := s.fetchAndProcessEvents()
			if err != nil {
				consecutiveFailures++
				s.logger.Warn().
					Err(err).
					Int("consecutive_failures", consecutiveFailures).
					Msg("Failed to fetch events")

				// If we've exceeded max retries, cancel the context and exit
				if consecutiveFailures >= s.retryAttempts {
					s.logger.Error().
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

// fetchAndProcessEvents calls the chip test sink API and processes new events
func (s *ChipEventsSubscriber) fetchAndProcessEvents() error {
	// Get current sequence
	s.sequenceMu.RLock()
	minSequence := s.lastSequence
	s.sequenceMu.RUnlock()

	// Build URL with sequence filter and limit
	url := fmt.Sprintf("%s/events?sequence=%d&limit=%d", s.chipServerURL, minSequence, s.eventLimit)

	var events []ChipEventMessage
	err := retry.Do(
		func() error {
			req, err := http.NewRequestWithContext(s.ctx, http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("create request: %w", err)
			}

			resp, err := s.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("http request: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
			}

			if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
				return fmt.Errorf("decode response: %w", err)
			}

			return nil
		},
		retry.Attempts(uint(s.retryAttempts)),
		retry.Delay(s.retryDelay),
		retry.OnRetry(func(n uint, err error) {
			s.logger.Debug().
				Uint("attempt", n+1).
				Err(err).
				Msg("Retrying fetch events")
		}),
	)

	if err != nil {
		return fmt.Errorf("fetch events failed after %d attempts: %w", s.retryAttempts, err)
	}

	if len(events) == 0 {
		return nil
	}

	// Process events in parallel batches for better throughput
	batchSize := s.batchSize
	numBatches := (len(events) + batchSize - 1) / batchSize
	processedEvents := make([]ChipEventMessage, 0, len(events))
	var processMu sync.Mutex
	var wg sync.WaitGroup

	for i := range numBatches {
		start := i * batchSize
		end := min(start+batchSize, len(events))
		batch := events[start:end]

		wg.Add(1)
		go func(eventBatch []ChipEventMessage) {
			defer wg.Done()

			localProcessed := make([]ChipEventMessage, 0, len(eventBatch))
			for _, event := range eventBatch {
				// Check if we already processed this sequence
				s.sequenceMu.RLock()
				lastSeq := s.lastSequence
				s.sequenceMu.RUnlock()

				if event.Sequence <= lastSeq {
					continue
				}

				// Extract workflow ID and node P2P ID and populate them in the event
				metadata, wfErr := extractMetadata(event)
				if wfErr != nil {
					s.logger.Debug().
						Err(wfErr).
						Str("entity", event.Entity).
						Uint64("sequence", event.Sequence).
						Msg("Failed to extract metadata")
				}
				
				// Populate fields if metadata was successfully extracted
				if metadata != nil {
					if metadata.workflowID != "" {
						event.WorkflowID = metadata.workflowID
					}
					if metadata.p2pID != "" {
						event.NodeP2PID = metadata.p2pID
					}
				}

				localProcessed = append(localProcessed, event)
			}

			// Append to global processed events (thread-safe)
			processMu.Lock()
			processedEvents = append(processedEvents, localProcessed...)
			processMu.Unlock()
		}(batch)
	}

	// Wait for all batches to complete
	wg.Wait()

	// Send all processed events to channel in sequence order
	for _, event := range processedEvents {
		select {
		case s.messageChan <- event:
			s.logger.Debug().
				Uint64("sequence", event.Sequence).
				Str("entity", event.Entity).
				Str("domain", event.Domain).
				Str("workflow_id", event.WorkflowID).
				Str("node_p2p_id", event.NodeP2PID).
				Msg("Event sent to channel")

			// Update last sequence
			s.sequenceMu.Lock()
			if event.Sequence > s.lastSequence {
				s.lastSequence = event.Sequence
			}
			s.sequenceMu.Unlock()

		case <-s.ctx.Done():
			return nil

		default:
			s.logger.Warn().
				Uint64("sequence", event.Sequence).
				Str("entity", event.Entity).
				Msg("Channel full, dropping event")
		}
	}

	return nil
}

type eventMetadata struct {
	workflowID, p2pID string
}

// extractWorkflowID extracts the workflow ID from an event
func extractMetadata(event ChipEventMessage) (*eventMetadata, error) {
	switch event.Entity {
	case "workflows.v1.UserLogs":
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
		n, err := base64.StdEncoding.Decode(decoded, []byte(event.Body))
		if err != nil {
			return nil, fmt.Errorf("failed to decode entity body from base64: %w", err)
		}
		decoded = decoded[:n]
		msg := &workflowevents.UserLogs{}
		err = proto.Unmarshal(decoded, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal entity with type '%s' to %T: %w", event.Entity, msg, err)
		}

		if msg.M != nil {
			return &eventMetadata{
				workflowID: msg.M.WorkflowID,
				p2pID:      msg.M.P2PID,
			}, nil

		}
		return nil, fmt.Errorf("UserLogs message has nil M field")
	case "BaseMessage":
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
		n, err := base64.StdEncoding.Decode(decoded, []byte(event.Body))
		if err != nil {
			return nil, fmt.Errorf("failed to decode entity body from base64: %w", err)
		}
		decoded = decoded[:n]
		msg := &commonevents.BaseMessage{}
		err = proto.Unmarshal(decoded, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal entity with type '%s' to %T: %w", event.Entity, msg, err)
		}

		if msg.Labels != nil {
			return &eventMetadata{
				workflowID: msg.Labels["workflowID"],
				p2pID:      msg.Labels["p2pID"],
			}, nil
		}
		return nil, fmt.Errorf("BaseMessage has nil Labels field")
	}

	return nil, nil
}

// LogChipEvents logs chip events as they arrive on the message channel.
// This function will block until the context is cancelled or the channel is closed.
func LogChipEvents(
	ctx context.Context,
	messageChan <-chan ChipEventMessage,
	logger zerolog.Logger,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-messageChan:
			if !ok {
				return
			}

			// Decode and log based on event type
			switch event.Entity {
			case "workflows.v1.UserLogs":
				decoded := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
				n, err := base64.StdEncoding.Decode(decoded, []byte(event.Body))
				if err != nil {
					logger.Warn().Err(err).Msg("Failed to decode UserLogs event")
					continue
				}
				decoded = decoded[:n]
				msg := &workflowevents.UserLogs{}
				if err := proto.Unmarshal(decoded, msg); err != nil {
					logger.Warn().Err(err).Msg("Failed to unmarshal UserLogs event")
					continue
				}

				for _, line := range msg.LogLines {
					logger.Info().
						Str("node_p2p_id", event.NodeP2PID).
						Str("workflow_id", event.WorkflowID).
						Uint64("sequence", event.Sequence).
						Msgf("[UserLog] %s", line.Message)
				}

			case "BaseMessage":
				decoded := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
				n, err := base64.StdEncoding.Decode(decoded, []byte(event.Body))
				if err != nil {
					logger.Warn().Err(err).Msg("Failed to decode BaseMessage event")
					continue
				}
				decoded = decoded[:n]
				msg := &commonevents.BaseMessage{}
				if err := proto.Unmarshal(decoded, msg); err != nil {
					logger.Warn().Err(err).Msg("Failed to unmarshal BaseMessage event")
					continue
				}

				if strings.Contains(msg.Msg, "heartbeat") || strings.Contains(msg.Msg, "metering report") {
					continue
				}

				logger.Info().
					Str("node_p2p_id", event.NodeP2PID).
					Str("workflow_id", event.WorkflowID).
					Msgf("[BaseMessage]#%d %s", event.Sequence, msg.Msg)

			default:
				logger.Debug().
					Str("entity", event.Entity).
					Str("domain", event.Domain).
					Str("node_p2p_id", event.NodeP2PID).
					Str("workflow_id", event.WorkflowID).
					Uint64("sequence", event.Sequence).
					Msg("Chip event received")
			}
		}
	}
}

// MatcherFunc is a function that checks if an event matches criteria
type ChipEventMatcherFunc func(event ChipEventMessage) (bool, error)

// AssertChipEventMatchedByNodes waits for N unique nodes to emit matching events for a specific workflow.
// Each node is identified by the NodeP2PID field in the event.
// This is a BLOCKING function that returns when:
// - N unique nodes have emitted matching events, OR
// - The timeout is reached, OR
// - The subscriber context is cancelled
func AssertChipEventMatchedByNodes(
	ctx context.Context,
	workflowID string,
	n int,
	messageChan <-chan ChipEventMessage,
	matcher ChipEventMatcherFunc,
	timeout time.Duration,
	logger zerolog.Logger,
) error {
	matchedNodes := make(map[string]bool) // Track which nodes (by P2P ID) have matching events
	deadline := time.After(timeout)

	logger.Info().
		Str("workflow_id", workflowID).
		Int("required_nodes", n).
		Dur("timeout", timeout).
		Msg("Waiting for event matches across nodes")

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"subscriber context cancelled: only %d/%d nodes emitted matching events: %w",
				len(matchedNodes),
				n,
				ctx.Err(),
			)

		case <-deadline:
			return fmt.Errorf(
				"timeout after %v: only %d/%d nodes emitted matching events",
				timeout,
				len(matchedNodes),
				n,
			)

		case event, ok := <-messageChan:
			if !ok {
				return fmt.Errorf(
					"message channel closed: only %d/%d nodes emitted matching events",
					len(matchedNodes),
					n,
				)
			}

			// Skip events without workflow ID - we can't attribute them to any workflow
			if event.WorkflowID == "" {
				logger.Debug().
					Str("entity", event.Entity).
					Uint64("sequence", event.Sequence).
					Msg("Event has no workflow ID, skipping")
				continue
			}

			// Filter by workflow ID if provided
			if workflowID != "" && event.WorkflowID != workflowID {
				continue
			}

			// Use the NodeP2PID field that was already extracted by the subscriber
			if event.NodeP2PID == "" {
				logger.Debug().
					Str("entity", event.Entity).
					Uint64("sequence", event.Sequence).
					Msg("Event has no node P2P ID, skipping")
				continue
			}

			if matched, err := matcher(event); err != nil {
				logger.Error().
					Err(err).
					Str("entity", event.Entity).
					Uint64("sequence", event.Sequence).
					Msg("Matcher function error")
				return fmt.Errorf("matcher function errored: %w", err)
			} else if matched {
				if !matchedNodes[event.NodeP2PID] {
					matchedNodes[event.NodeP2PID] = true

					logger.Info().
						Str("workflow_id", event.WorkflowID).
						Str("node_p2p_id", event.NodeP2PID).
						Str("entity", event.Entity).
						Uint64("sequence", event.Sequence).
						Int("total_matched_nodes", len(matchedNodes)).
						Msg("Found matching event from node")

					if len(matchedNodes) >= n {
						logger.Info().
							Int("matched_nodes", len(matchedNodes)).
							Msg("Required number of nodes have matching events")
						return nil
					}
				}
			}
		}
	}
}

// FanOutChipEvents distributes events from a source channel to multiple consumer channels
// Useful for having multiple goroutines process the same event stream
func FanOutChipEvents(
	ctx context.Context,
	source <-chan ChipEventMessage,
	numConsumers int,
) []<-chan ChipEventMessage {
	destinations := make([]chan ChipEventMessage, numConsumers)
	outputs := make([]<-chan ChipEventMessage, numConsumers)

	// Use larger buffers for destination channels to handle bursts
	for i := range numConsumers {
		destinations[i] = make(chan ChipEventMessage, 500)
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
				return
			case msg, ok := <-source:
				if !ok {
					return
				}
				// Broadcast to all consumers (non-blocking - drop if slow)
				for _, dest := range destinations {
					select {
					case dest <- msg:
						// Successfully sent
					case <-ctx.Done():
						return
					default:
						// Channel full, drop message for this consumer
						// This prevents slow consumers from blocking the fan-out
					}
				}
			}
		}
	}()

	return outputs
}

// GetUserLogMatcherFn returns a matcher function that checks if a UserLogs event contains the expected message.
// It also checks for workflow engine initialization failures and returns an error if found.
func GetUserLogMatcherFn(expectedMessage string) ChipEventMatcherFunc {
	return func(event ChipEventMessage) (bool, error) {
		// Check for UserLogs events
		if event.Entity == "workflows.v1.UserLogs" {
			decoded := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
			n, err := base64.StdEncoding.Decode(decoded, []byte(event.Body))
			if err != nil {
				return false, nil // Ignore decode errors
			}
			decoded = decoded[:n]

			msg := &workflowevents.UserLogs{}
			if err := proto.Unmarshal(decoded, msg); err != nil {
				return false, nil // Ignore unmarshal errors
			}

			for _, line := range msg.LogLines {
				if len(line.Message) > 0 && strings.Contains(line.Message, expectedMessage) {
					return true, nil
				}
			}
		}

		// Check for BaseMessage events that might indicate failures
		if event.Entity == "BaseMessage" {
			decoded := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
			n, err := base64.StdEncoding.Decode(decoded, []byte(event.Body))
			if err != nil {
				return false, nil // Ignore decode errors
			}
			decoded = decoded[:n]

			msg := &commonevents.BaseMessage{}
			if err := proto.Unmarshal(decoded, msg); err != nil {
				return false, nil // Ignore unmarshal errors
			}

			if strings.Contains(msg.Msg, "Workflow Engine initialization failed") {
				return false, fmt.Errorf("found engine initialization failure message: %s", msg.Msg)
			}
		}

		return false, nil
	}
}
