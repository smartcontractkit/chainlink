package events

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

func EmitWorkflowStatusChangedEvent(
	ctx context.Context,
	labels map[string]string,
	status string,
) error {
	metadata := buildWorkflowMetadata(labels)
	event := &events.WorkflowStatusChanged{
		M:      metadata,
		Status: status,
	}

	return emitProtoMessage(ctx, event)
}

func EmitExecutionStartedEvent(
	ctx context.Context,
	labels map[string]string,
	triggerEventID string,
	executionID string,
) error {
	labels[platform.KeyWorkflowExecutionID] = executionID
	metadata := buildWorkflowMetadata(labels)

	event := &events.WorkflowExecutionStarted{
		M:         metadata,
		Timestamp: time.Now().Format(time.RFC3339Nano),
		TriggerID: triggerEventID,
	}

	return emitProtoMessage(ctx, event)
}

func EmitExecutionFinishedEvent(ctx context.Context, labels map[string]string, status string, executionID string) error {
	labels[platform.KeyWorkflowExecutionID] = executionID
	metadata := buildWorkflowMetadata(labels)

	event := &events.WorkflowExecutionFinished{
		M:         metadata,
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Status:    status,
	}

	return emitProtoMessage(ctx, event)
}

func EmitCapabilityStartedEvent(ctx context.Context, labels map[string]string, executionID, capabilityID, stepRef string) error {
	labels[platform.KeyWorkflowExecutionID] = executionID
	metadata := buildWorkflowMetadata(labels)

	event := &events.CapabilityExecutionStarted{
		M:            metadata,
		Timestamp:    time.Now().Format(time.RFC3339Nano),
		CapabilityID: capabilityID,
		StepRef:      stepRef,
	}

	return emitProtoMessage(ctx, event)
}

func EmitCapabilityFinishedEvent(ctx context.Context, labels map[string]string, executionID, capabilityID, stepRef, status string) error {
	labels[platform.KeyWorkflowExecutionID] = executionID
	metadata := buildWorkflowMetadata(labels)

	event := &events.CapabilityExecutionFinished{
		M:            metadata,
		Timestamp:    time.Now().Format(time.RFC3339Nano),
		CapabilityID: capabilityID,
		StepRef:      stepRef,
		Status:       status,
	}

	return emitProtoMessage(ctx, event)
}

func EmitMeteringReport(ctx context.Context, labels map[string]string, rpt *events.MeteringReport) error {
	rpt.Metadata = buildWorkflowMetadata(labels)

	return emitProtoMessage(ctx, rpt)
}

// EmitProtoMessage marshals a proto.Message and emits it via beholder.
func emitProtoMessage(ctx context.Context, msg proto.Message) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Determine the schema and entity based on the message type
	// entity must be prefixed with the proto package name
	var schema, entity string
	switch msg.(type) {
	case *events.WorkflowExecutionStarted:
		schema = SchemaWorkflowStarted
		entity = fmt.Sprintf("%s.%s", ProtoPkg, WorkflowExecutionStarted)
	case *events.WorkflowExecutionFinished:
		schema = SchemaWorkflowFinished
		entity = fmt.Sprintf("%s.%s", ProtoPkg, WorkflowExecutionFinished)
	case *events.CapabilityExecutionStarted:
		schema = SchemaCapabilityStarted
		entity = fmt.Sprintf("%s.%s", ProtoPkg, CapabilityExecutionStarted)
	case *events.CapabilityExecutionFinished:
		schema = SchemaCapabilityFinished
		entity = fmt.Sprintf("%s.%s", ProtoPkg, CapabilityExecutionFinished)
	case *events.MeteringReport:
		schema = MeteringReportSchema
		entity = fmt.Sprintf("%s.%s", ProtoPkg, MeteringReportEntity)
	default:
		return fmt.Errorf("unknown message type: %T", msg)
	}

	return beholder.GetEmitter().Emit(ctx, b,
		"beholder_data_schema", schema, // required
		"beholder_domain", "platform", // required
		"beholder_entity", entity) // required
}

// buildWorkflowMetadata populates a WorkflowMetadata from kvs (map[string]string).
func buildWorkflowMetadata(kvs map[string]string) *events.WorkflowMetadata {
	m := &events.WorkflowMetadata{}

	m.WorkflowOwner = kvs[platform.KeyWorkflowOwner]
	m.WorkflowName = kvs[platform.KeyWorkflowName]
	m.Version = kvs[platform.KeyWorkflowVersion]
	m.WorkflowID = kvs[platform.KeyWorkflowID]
	m.WorkflowExecutionID = kvs[platform.KeyWorkflowExecutionID]

	if donIDStr, ok := kvs[platform.KeyDonID]; ok {
		if id, err := strconv.ParseInt(donIDStr, 10, 32); err == nil {
			m.DonID = int32(id)
		}
	}

	m.P2PID = kvs[platform.KeyP2PID]

	if donFStr, ok := kvs[platform.KeyDonF]; ok {
		if id, err := strconv.ParseInt(donFStr, 10, 32); err == nil {
			m.DonF = int32(id)
		}
	}
	if donNStr, ok := kvs[platform.KeyDonN]; ok {
		if id, err := strconv.ParseInt(donNStr, 10, 32); err == nil {
			m.DonN = int32(id)
		}
	}

	return m
}
