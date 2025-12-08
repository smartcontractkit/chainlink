//go:build dev

package devobservability

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
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

	// Store event (either with execution context or as orphan)
	workflowID, executionID, hasContext := GetExecutionContext(ctx)

	// If context is missing, try to extract from protobuf message using event type
	// if !hasContext && len(payload) > 0 && eventType != "" {
	// 	extractedWorkflowID, extractedExecutionID := extractIDsFromPayload(payload, eventType)
	// 	if extractedWorkflowID != "" {
	// 		workflowID = extractedWorkflowID
	// 	}
	// 	if extractedExecutionID != "" {
	// 		executionID = extractedExecutionID
	// 	}
	// 	if workflowID != "" && executionID != "" {
	// 		hasContext = true
	// 		fmt.Printf("[DevObservability] Extracted IDs from payload: workflowID=%s, executionID=%s\n", workflowID, executionID)
	// 	}
	// }

	fmt.Printf("[DevObservability] Context check: workflowID=%s, executionID=%s, hasContext=%v, eventType=%s\n", workflowID, executionID, hasContext, eventType)

	if !hasContext && strings.Contains(eventType, "BaseMessage") {
		msg := &commonevents.BaseMessage{}
		if err := proto.Unmarshal(payload, msg); err == nil {
			fmt.Printf("[DevObservability] BaseMessage without context: %s\n", msg.GetMsg())
		}
	}

	if eventType != "" {
		if store, ok := globalStore.(*devStore); ok {
			if hasContext && workflowID != "" && executionID != "" {
				fmt.Printf("[DevObservability] Storing event with context: workflow=%s, execution=%s, type=%s\n", workflowID, executionID, eventType)
				store.storeRawEvent(ctx, workflowID, executionID, eventType, payload)
			} else if workflowID != "" {
				fmt.Printf("[DevObservability] Storing workflow-level event: workflow=%s, type=%s\n", workflowID, eventType)
				store.storeWorkflowEvent(ctx, workflowID, eventType, payload)
			} else {
				fmt.Printf("[DevObservability] Storing orphan event (no context): type=%s\n", eventType)
				store.storeOrphanEvent(ctx, eventType, payload)
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

// extractIDsFromPayload extracts workflowID and executionID from a protobuf message using the known event type
func extractIDsFromPayload(payload []byte, eventType string) (workflowID, executionID string) {
	var msg proto.Message

	// Unmarshal to the correct type based on eventType
	switch eventType {
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

	// V1 events
	case "workflows.v1.WorkflowExecutionStarted":
		msg = &events.WorkflowExecutionStarted{}
	case "workflows.v1.WorkflowExecutionFinished":
		msg = &events.WorkflowExecutionFinished{}
	case "workflows.v1.CapabilityExecutionStarted":
		msg = &events.CapabilityExecutionStarted{}
	case "workflows.v1.CapabilityExecutionFinished":
		msg = &events.CapabilityExecutionFinished{}
	case "workflows.v1.UserLogs":
		msg = &events.UserLogs{}
	case "workflows.v1.WorkflowStatusChanged":
		msg = &events.WorkflowStatusChanged{}
	case "workflows.v1.MeteringReport":
		msg = &events.MeteringReport{}

	// BaseMessage
	case "BaseMessage", "common.BaseMessage":
		msg = &commonevents.BaseMessage{}

	default:
		fmt.Printf("[DevObservability] Unknown event type for extraction: %s\n", eventType)
		return "", ""
	}

	if err := proto.Unmarshal(payload, msg); err != nil {
		fmt.Printf("[DevObservability] Failed to unmarshal %s: %v\n", eventType, err)
		return "", ""
	}

	// Extract IDs based on event type family
	if len(eventType) >= 12 && eventType[:12] == "workflows.v2" {
		// V2 events: workflow.workflowID (nested) and workflowExecutionID (top-level)
		reflectMsg := msg.ProtoReflect()
		fields := reflectMsg.Descriptor().Fields()

		if workflowField := fields.ByName("workflow"); workflowField != nil && reflectMsg.Has(workflowField) {
			workflowMsg := reflectMsg.Get(workflowField).Message()
			workflowMsgFields := workflowMsg.Descriptor().Fields()

			if wfIDField := workflowMsgFields.ByName("workflowID"); wfIDField != nil && workflowMsg.Has(wfIDField) {
				workflowID = workflowMsg.Get(wfIDField).String()
			}
		}

		if execIDField := fields.ByName("workflowExecutionID"); execIDField != nil && reflectMsg.Has(execIDField) {
			executionID = reflectMsg.Get(execIDField).String()
		}

		fmt.Printf("[DevObservability] Extracted from V2 event %s: workflowID=%s, executionID=%s\n", eventType, workflowID, executionID)
	} else if len(eventType) >= 12 && eventType[:12] == "workflows.v1" {
		// V1 events: m.workflow_id and m.workflow_execution_id
		reflectMsg := msg.ProtoReflect()
		fields := reflectMsg.Descriptor().Fields()

		if mField := fields.ByName("m"); mField != nil && reflectMsg.Has(mField) {
			metadata := reflectMsg.Get(mField).Message()
			metadataFields := metadata.Descriptor().Fields()

			if wfIDField := metadataFields.ByName("workflow_id"); wfIDField != nil && metadata.Has(wfIDField) {
				if wfIDField.Kind() == protoreflect.BytesKind {
					workflowID = fmt.Sprintf("%x", metadata.Get(wfIDField).Bytes())
				} else {
					workflowID = metadata.Get(wfIDField).String()
				}
			}

			if execIDField := metadataFields.ByName("workflow_execution_id"); execIDField != nil && metadata.Has(execIDField) {
				if execIDField.Kind() == protoreflect.BytesKind {
					executionID = fmt.Sprintf("%x", metadata.Get(execIDField).Bytes())
				} else {
					executionID = metadata.Get(execIDField).String()
				}
			}
		}

		fmt.Printf("[DevObservability] Extracted from V1 event %s: workflowID=%s, executionID=%s\n", eventType, workflowID, executionID)
	} else {
		// BaseMessage: labels.workflow_id and labels.workflow_execution_id
		reflectMsg := msg.ProtoReflect()
		fields := reflectMsg.Descriptor().Fields()

		if labelsField := fields.ByName("labels"); labelsField != nil && reflectMsg.Has(labelsField) {
			labelsMap := reflectMsg.Get(labelsField).Map()

			if wfKey := labelsMap.Get(protoreflect.MapKey(protoreflect.ValueOfString("workflow_id"))); wfKey.IsValid() {
				workflowID = wfKey.String()
			}

			if execKey := labelsMap.Get(protoreflect.MapKey(protoreflect.ValueOfString("workflow_execution_id"))); execKey.IsValid() {
				executionID = execKey.String()
			}
		}

		fmt.Printf("[DevObservability] Extracted from BaseMessage: workflowID=%s, executionID=%s\n", workflowID, executionID)
	}

	return workflowID, executionID
}
