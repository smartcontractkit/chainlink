//go:build dev

package devobservability

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
)

// recordingEmitter wraps the beholder emitter to record all emissions
type recordingEmitter struct {
	underlying beholder.Emitter
}

// WrapEmitter wraps a beholder emitter with dev recording functionality
func WrapEmitter(underlying beholder.Emitter) beholder.Emitter {
	fmt.Printf("[DevObservability] Wrapping emitter (underlying type: %T)\n", underlying)
	wrapped := &recordingEmitter{underlying: underlying}
	fmt.Printf("[DevObservability] Created wrapped emitter: %p\n", wrapped)
	return wrapped
}

func (e *recordingEmitter) Emit(ctx context.Context, payload []byte, kvs ...any) error {
	fmt.Printf("[DevObservability] Emit called with %d kvs\n", len(kvs))

	// Extract event type from beholder metadata
	var eventType string
	for i := 0; i < len(kvs); i += 2 {
		if i+1 < len(kvs) {
			if key, ok := kvs[i].(string); ok {
				fmt.Printf("[DevObservability] KV pair: %s = %v\n", key, kvs[i+1])
				if key == "beholder_entity" {
					if val, ok := kvs[i+1].(string); ok {
						eventType = val
						fmt.Printf("[DevObservability] Found event type: %s\n", eventType)
						break
					}
				}
			}
		}
	}

	// If beholder_entity is missing, try to extract from protobuf payload
	if eventType == "" && len(payload) > 0 {
		fmt.Printf("[DevObservability] beholder_entity missing, attempting protobuf reflection\n")
		eventType = extractEventTypeFromProto(payload)
		if eventType != "" {
			fmt.Printf("[DevObservability] Extracted event type via reflection: %s\n", eventType)
		}
	}

	// Store event using workflow ID from CRE context
	creContext := contexts.CREValue(ctx)
	workflowID := creContext.Workflow

	fmt.Printf("[DevObservability] Context check: workflowID=%s, eventType=%s\n", workflowID, eventType)

	if eventType != "" {
		if store, ok := globalStore.(*devStore); ok {
			if workflowID != "" {
				fmt.Printf("[DevObservability] Storing workflow event: workflow=%s, type=%s\n", workflowID, eventType)
				store.storeWorkflowEvent(workflowID, eventType, payload)
			} else {
				fmt.Printf("[DevObservability] Storing orphan event (no workflow context): type=%s\n", eventType)
				store.storeOrphanEvent(eventType, payload)
			}
		}
	} else {
		fmt.Printf("[DevObservability] SKIPPED: No event type found\n")
	}

	// Forward to underlying emitter
	return e.underlying.Emit(ctx, payload, kvs...)
}

func (e *recordingEmitter) Close() error {
	return e.underlying.Close()
}

// extractEventTypeFromProto attempts to determine the event type from a protobuf payload
// using protobuf reflection to discover the message type dynamically
func extractEventTypeFromProto(payload []byte) string {
	// Iterate through all registered message types in the global registry
	var foundType string
	protoregistry.GlobalTypes.RangeMessages(func(mt protoreflect.MessageType) bool {
		// Create a new instance of this message type
		msg := mt.New().Interface()

		// Try to unmarshal the payload into this message type
		if err := proto.Unmarshal(payload, msg); err == nil {
			// Successfully unmarshaled - check if it's a workflow event
			fullName := string(mt.Descriptor().FullName())

			// Only consider workflow event types (workflows.v1.*, workflows.v2.*, or BaseMessage)
			if (len(fullName) > 10 && (fullName[:12] == "workflows.v1" || fullName[:12] == "workflows.v2")) || fullName == "common.BaseMessage" {
				fmt.Printf("[DevObservability] Proto reflection matched: %s\n", fullName)
				foundType = fullName
				return false // Stop iteration
			}
		}
		return true // Continue iteration
	})

	if foundType == "" {
		fmt.Printf("[DevObservability] Failed to match protobuf payload to any known event type\n")
	}

	return foundType
}
