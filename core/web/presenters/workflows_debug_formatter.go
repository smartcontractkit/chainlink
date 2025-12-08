//go:build dev

package presenters

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/devobservability"
)

// FormattedEventEntry represents an event with decoded protobuf message
type FormattedEventEntry struct {
	Sequence    int64                  `json:"sequence"`
	Type        string                 `json:"type"`
	Timestamp   string                 `json:"timestamp"`
	MessageJSON map[string]interface{} `json:"messageJson,omitempty"`
	DecodeError string                 `json:"decodeError,omitempty"`
}

// WorkflowDebugOrphanEventsFormattedResource represents orphan events with decoded protobufs
type WorkflowDebugOrphanEventsFormattedResource struct {
	JAID
	Events []FormattedEventEntry `json:"events"`
}

// NewWorkflowDebugOrphanEventsFormattedResource creates a formatted orphan events resource with decoded protobufs
func NewWorkflowDebugOrphanEventsFormattedResource(events []devobservability.EventEntry) *WorkflowDebugOrphanEventsFormattedResource {
	formattedEvents := make([]FormattedEventEntry, len(events))

	for i, evt := range events {
		formatted := FormattedEventEntry{
			Sequence:  evt.Sequence,
			Type:      evt.Type,
			Timestamp: evt.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"),
			// Don't include base64 in decoded format - it's redundant
		}

		// Try to decode the protobuf message based on event type
		var protoMsg proto.Message
		switch evt.Type {
		case "workflows.v1.WorkflowExecutionStarted":
			protoMsg = &workflowevents.WorkflowExecutionStarted{}
		case "workflows.v1.WorkflowExecutionFinished":
			protoMsg = &workflowevents.WorkflowExecutionFinished{}
		case "workflows.v1.CapabilityExecutionStarted":
			protoMsg = &workflowevents.CapabilityExecutionStarted{}
		case "workflows.v1.CapabilityExecutionFinished":
			protoMsg = &workflowevents.CapabilityExecutionFinished{}
		case "workflows.v1.MeteringReport":
			protoMsg = &workflowevents.MeteringReport{}
		case "workflows.v1.WorkflowStatusChanged":
			protoMsg = &workflowevents.WorkflowStatusChanged{}
		case "workflows.v1.UserLogs":
			protoMsg = &workflowevents.UserLogs{}
		case "BaseMessage":
			protoMsg = &commonevents.BaseMessage{}
		case "workflows.v2.WorkflowExecutionStarted":
			protoMsg = &eventsv2.WorkflowExecutionStarted{}
		case "workflows.v2.WorkflowExecutionFinished":
			protoMsg = &eventsv2.WorkflowExecutionFinished{}
		case "workflows.v2.CapabilityExecutionStarted":
			protoMsg = &eventsv2.CapabilityExecutionStarted{}
		case "workflows.v2.CapabilityExecutionFinished":
			protoMsg = &eventsv2.CapabilityExecutionFinished{}
		case "workflows.v2.TriggerExecutionStarted":
			protoMsg = &eventsv2.TriggerExecutionStarted{}
		case "workflows.v2.WorkflowUserLog":
			protoMsg = &eventsv2.WorkflowUserLog{}
		case "workflows.v2.WorkflowActivated":
			protoMsg = &eventsv2.WorkflowActivated{}
		case "workflows.v2.WorkflowPaused":
			protoMsg = &eventsv2.WorkflowPaused{}
		case "workflows.v2.WorkflowDeleted":
			protoMsg = &eventsv2.WorkflowDeleted{}
		}

		if protoMsg != nil {
			if err := proto.Unmarshal(evt.Message, protoMsg); err != nil {
				formatted.DecodeError = fmt.Sprintf("Failed to unmarshal protobuf: %v", err)
			} else {
				// Convert protobuf to JSON
				jsonBytes, err := protojson.Marshal(protoMsg)
				if err != nil {
					formatted.DecodeError = fmt.Sprintf("Failed to convert to JSON: %v", err)
				} else {
					var jsonObj map[string]interface{}
					if err := json.Unmarshal(jsonBytes, &jsonObj); err != nil {
						formatted.DecodeError = fmt.Sprintf("Failed to parse JSON: %v", err)
					} else {
						formatted.MessageJSON = jsonObj
					}
				}
			}
		}

		formattedEvents[i] = formatted
	}

	return &WorkflowDebugOrphanEventsFormattedResource{
		JAID:   NewJAID("orphan_events_formatted"),
		Events: formattedEvents,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugOrphanEventsFormattedResource) GetName() string {
	return "orphan_events_formatted"
}

// WorkflowDebugWorkflowEventsFormattedResource represents workflow-level events with decoded protobuf messages
type WorkflowDebugWorkflowEventsFormattedResource struct {
	JAID
	WorkflowID string                `json:"workflowId"`
	Events     []FormattedEventEntry `json:"events"`
}

// NewWorkflowDebugWorkflowEventsFormattedResource creates a formatted workflow events resource with decoded protobufs
func NewWorkflowDebugWorkflowEventsFormattedResource(workflowID string, events []devobservability.EventEntry) *WorkflowDebugWorkflowEventsFormattedResource {
	formattedEvents := make([]FormattedEventEntry, len(events))

	for i, evt := range events {
		formatted := FormattedEventEntry{
			Sequence:  evt.Sequence,
			Type:      evt.Type,
			Timestamp: evt.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"),
		}

		// Try to decode the protobuf message based on event type
		var protoMsg proto.Message
		switch evt.Type {
		case "workflows.v1.WorkflowExecutionStarted":
			protoMsg = &workflowevents.WorkflowExecutionStarted{}
		case "workflows.v1.WorkflowExecutionFinished":
			protoMsg = &workflowevents.WorkflowExecutionFinished{}
		case "workflows.v1.CapabilityExecutionStarted":
			protoMsg = &workflowevents.CapabilityExecutionStarted{}
		case "workflows.v1.CapabilityExecutionFinished":
			protoMsg = &workflowevents.CapabilityExecutionFinished{}
		case "workflows.v1.MeteringReport":
			protoMsg = &workflowevents.MeteringReport{}
		case "workflows.v1.WorkflowStatusChanged":
			protoMsg = &workflowevents.WorkflowStatusChanged{}
		case "workflows.v1.UserLogs":
			protoMsg = &workflowevents.UserLogs{}
		case "BaseMessage":
			protoMsg = &commonevents.BaseMessage{}
		case "workflows.v2.WorkflowExecutionStarted":
			protoMsg = &eventsv2.WorkflowExecutionStarted{}
		case "workflows.v2.WorkflowExecutionFinished":
			protoMsg = &eventsv2.WorkflowExecutionFinished{}
		case "workflows.v2.CapabilityExecutionStarted":
			protoMsg = &eventsv2.CapabilityExecutionStarted{}
		case "workflows.v2.CapabilityExecutionFinished":
			protoMsg = &eventsv2.CapabilityExecutionFinished{}
		case "workflows.v2.TriggerExecutionStarted":
			protoMsg = &eventsv2.TriggerExecutionStarted{}
		case "workflows.v2.WorkflowUserLog":
			protoMsg = &eventsv2.WorkflowUserLog{}
		case "workflows.v2.WorkflowActivated":
			protoMsg = &eventsv2.WorkflowActivated{}
		case "workflows.v2.WorkflowPaused":
			protoMsg = &eventsv2.WorkflowPaused{}
		case "workflows.v2.WorkflowDeleted":
			protoMsg = &eventsv2.WorkflowDeleted{}
		}

		if protoMsg != nil {
			if err := proto.Unmarshal(evt.Message, protoMsg); err == nil {
				if jsonBytes, err := protojson.Marshal(protoMsg); err == nil {
					var jsonObj map[string]interface{}
					if err := json.Unmarshal(jsonBytes, &jsonObj); err == nil {
						formatted.MessageJSON = jsonObj
					}
				}
			}
		}

		formattedEvents[i] = formatted
	}

	return &WorkflowDebugWorkflowEventsFormattedResource{
		JAID:       NewJAID("workflow_events_formatted"),
		WorkflowID: workflowID,
		Events:     formattedEvents,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugWorkflowEventsFormattedResource) GetName() string {
	return "workflow_events_formatted"
}
