package chipsink

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	workfloweventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
)

// EventData decodes known CHiP workflow event types to human-readable JSON maps.
// Unknown or undecodable events fall back to a minimal metadata+base64 representation.
func EventData(event *pb.CloudEvent) map[string]any {
	if event == nil {
		return map[string]any{}
	}

	msg := typedMessageForEventType(strings.TrimSpace(event.GetType()))
	if msg != nil {
		if protoData := event.GetProtoData(); protoData != nil && len(protoData.GetValue()) > 0 {
			if err := proto.Unmarshal(protoData.GetValue(), msg); err == nil {
				if asMap, ok := protoMessageAsMap(msg); ok {
					return asMap
				}
			}
		}
	}

	fallback := map[string]any{
		"id":          strings.TrimSpace(event.GetId()),
		"type":        strings.TrimSpace(event.GetType()),
		"source":      strings.TrimSpace(event.GetSource()),
		"specVersion": strings.TrimSpace(event.GetSpecVersion()),
	}
	if protoData := event.GetProtoData(); protoData != nil && len(protoData.GetValue()) > 0 {
		fallback["protoDataBase64"] = base64.StdEncoding.EncodeToString(protoData.GetValue())
	}
	if textData := strings.TrimSpace(event.GetTextData()); textData != "" {
		fallback["textData"] = textData
	}
	return fallback
}

func protoMessageAsMap(msg proto.Message) (map[string]any, bool) {
	dataBytes, err := (protojson.MarshalOptions{Multiline: false}).Marshal(msg)
	if err != nil {
		return nil, false
	}
	var out map[string]any
	if err := json.Unmarshal(dataBytes, &out); err != nil {
		return nil, false
	}
	return out, true
}

func typedMessageForEventType(eventType string) proto.Message {
	switch eventType {
	// workflows.v1 events
	case "workflows.v1.CapabilityExecutionFinished":
		return &workflowevents.CapabilityExecutionFinished{}
	case "workflows.v1.CapabilityExecutionStarted":
		return &workflowevents.CapabilityExecutionStarted{}
	case "workflows.v1.MeteringReport":
		return &workflowevents.MeteringReport{}
	case "workflows.v1.TransmissionsScheduledEvent":
		return &workflowevents.TransmissionsScheduledEvent{}
	case "workflows.v1.TransmitScheduleEvent":
		return &workflowevents.TransmitScheduleEvent{}
	case "workflows.v1.WorkflowExecutionFinished":
		return &workflowevents.WorkflowExecutionFinished{}
	case "workflows.v1.WorkflowExecutionStarted":
		return &workflowevents.WorkflowExecutionStarted{}
	case "workflows.v1.WorkflowStatusChanged":
		return &workflowevents.WorkflowStatusChanged{}
	case "workflows.v1.UserLogs":
		return &workflowevents.UserLogs{}

	// workflows.v2 events
	case "workflows.v2.CapabilityExecutionFinished":
		return &workfloweventsv2.CapabilityExecutionFinished{}
	case "workflows.v2.CapabilityExecutionStarted":
		return &workfloweventsv2.CapabilityExecutionStarted{}
	case "workflows.v2.TriggerExecutionStarted":
		return &workfloweventsv2.TriggerExecutionStarted{}
	case "workflows.v2.WorkflowActivated":
		return &workfloweventsv2.WorkflowActivated{}
	case "workflows.v2.WorkflowDeleted":
		return &workfloweventsv2.WorkflowDeleted{}
	case "workflows.v2.WorkflowDeployed":
		return &workfloweventsv2.WorkflowDeployed{}
	case "workflows.v2.WorkflowExecutionFinished":
		return &workfloweventsv2.WorkflowExecutionFinished{}
	case "workflows.v2.WorkflowExecutionStarted":
		return &workfloweventsv2.WorkflowExecutionStarted{}
	case "workflows.v2.WorkflowPaused":
		return &workfloweventsv2.WorkflowPaused{}
	case "workflows.v2.WorkflowUpdated":
		return &workfloweventsv2.WorkflowUpdated{}
	case "workflows.v2.WorkflowUserLog":
		return &workfloweventsv2.WorkflowUserLog{}

	case "BaseMessage":
		return &commonevents.BaseMessage{}
	default:
		return nil
	}
}
