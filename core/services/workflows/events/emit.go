package events

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

func EmitWorkflowStatusChangedEvent(
	ctx context.Context,
	labels map[string]string,
	status string,
) error {
	metadata := buildWorkflowMetadata(labels, "")
	event := &events.WorkflowStatusChanged{
		M:      metadata,
		Status: status,
	}

	return emitProtoMessage(ctx, event)
}

func EmitWorkflowStatusChangedEventV2(
	ctx context.Context,
	labels map[string]string,
	head *types.Head,
	status string,
	binaryURL string,
	configURL string,
	eventErr error,
) error {
	// Emit v1 event
	var multiErr error
	if err := EmitWorkflowStatusChangedEvent(ctx, labels, status); err != nil {
		multiErr = errors.Join(multiErr, err)
	}

	// Prepare v2 event data
	creInfo := buildCREMetadataV2(labels)
	workflow := buildWorkflowV2(labels, binaryURL, configURL)
	txInfo := &eventsv2.TransactionInfo{
		ChainSelector: labels[platform.WorkflowRegistryChainSelector],
		TxHash:        hex.EncodeToString(head.Hash),
	}

	var v2Event proto.Message
	var errorMessage string
	if eventErr != nil {
		errorMessage = eventErr.Error()
	}

	switch status {
	case WorkflowActivated:
		v2Event = &eventsv2.WorkflowActivated{
			CreInfo:      creInfo,
			Workflow:     workflow,
			TxInfo:       txInfo,
			Timestamp:    time.Now().Format(time.RFC3339),
			ErrorMessage: errorMessage,
		}

	case WorkflowPaused:
		v2Event = &eventsv2.WorkflowPaused{
			CreInfo:      creInfo,
			Workflow:     workflow,
			TxInfo:       txInfo,
			Timestamp:    time.Now().Format(time.RFC3339),
			ErrorMessage: errorMessage,
		}

	case WorkflowDeleted:
		v2Event = &eventsv2.WorkflowDeleted{
			CreInfo:      creInfo,
			Workflow:     workflow,
			TxInfo:       txInfo,
			Timestamp:    time.Now().Format(time.RFC3339),
			ErrorMessage: errorMessage,
		}
	}

	// Emit v2 event
	if err := emitProtoMessage(ctx, v2Event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}

	return multiErr
}

func EmitExecutionStartedEvent(
	ctx context.Context,
	labels map[string]string,
	triggerEventID string,
	executionID string,
) error {
	metadata := buildWorkflowMetadata(labels, executionID)

	event := &events.WorkflowExecutionStarted{
		M:         metadata,
		Timestamp: time.Now().Format(time.RFC3339Nano),
		TriggerID: triggerEventID,
	}

	// Also emit v2 event
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	v2Event := &eventsv2.WorkflowExecutionStarted{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: executionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		TriggerID:           triggerEventID,
	}

	// Emit both v1 and v2 events
	var multiErr error
	if err := emitProtoMessage(ctx, event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	if err := emitProtoMessage(ctx, v2Event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	return multiErr
}

func EmitExecutionFinishedEvent(ctx context.Context, labels map[string]string, status string, executionID string, lggr logger.Logger) error {
	metadata := buildWorkflowMetadata(labels, executionID)

	event := &events.WorkflowExecutionFinished{
		M:         metadata,
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Status:    status,
	}

	// Also emit v2 event
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// Convert status string to v2 ExecutionStatus enum
	var executionStatus eventsv2.ExecutionStatus
	switch status {
	case "completed", "completed_early_exit": // there are enums in workflows/store, but we shouldn't import that here
		executionStatus = eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED
	case "errored", "timeout":
		executionStatus = eventsv2.ExecutionStatus_EXECUTION_STATUS_FAILED
	default:
		executionStatus = eventsv2.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED
	}

	v2Event := &eventsv2.WorkflowExecutionFinished{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: executionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		Status:              executionStatus,
	}

	// Emit both v1 and v2 events
	var multiErr error
	if err := emitProtoMessage(ctx, event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	if err := emitProtoMessage(ctx, v2Event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	return multiErr
}

func EmitCapabilityStartedEvent(ctx context.Context, labels map[string]string, executionID, capabilityID, stepRef string) error {
	metadata := buildWorkflowMetadata(labels, executionID)

	event := &events.CapabilityExecutionStarted{
		M:            metadata,
		Timestamp:    time.Now().Format(time.RFC3339Nano),
		CapabilityID: capabilityID,
		StepRef:      stepRef,
	}

	// Also emit v2 event
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// Convert stepRef string to int32
	// V1 engine has arbitrary string stepRefs, v2 engine has monotonically increasing integers
	// We will support both v1 and v2 events for the short term, so need to handle both cases
	stepRefInt, err := strconv.ParseInt(stepRef, 10, 32)
	if err != nil {
		stepRefInt = -1
	}

	v2Event := &eventsv2.CapabilityExecutionStarted{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: executionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		CapabilityID:        capabilityID,
		StepRef:             int32(stepRefInt),
	}

	// Emit both v1 and v2 events
	var multiErr error
	if err := emitProtoMessage(ctx, event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	if err := emitProtoMessage(ctx, v2Event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	return multiErr
}

func EmitTriggerExecutionStarted(ctx context.Context, labels map[string]string, triggerID, workflowExecutionID string) error {
	// Emit v2 event
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	v2Event := &eventsv2.TriggerExecutionStarted{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: workflowExecutionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		TriggerID:           triggerID,
	}

	return emitProtoMessage(ctx, v2Event)
}

func EmitCapabilityFinishedEvent(ctx context.Context, labels map[string]string, executionID, capabilityID, stepRef, status string, capErr error) error {
	metadata := buildWorkflowMetadata(labels, executionID)

	event := &events.CapabilityExecutionFinished{
		M:            metadata,
		Timestamp:    time.Now().Format(time.RFC3339Nano),
		CapabilityID: capabilityID,
		StepRef:      stepRef,
		Status:       status,
	}

	// Also emit v2 event
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// Convert stepRef string to int32
	// V1 engine has arbitrary string stepRefs, v2 engine has monotonically increasing integers
	// We will support both v1 and v2 events for the short term, so need to handle both cases
	stepRefInt, err := strconv.ParseInt(stepRef, 10, 32)
	if err != nil {
		stepRefInt = -1
	}

	// Convert status string to v2 ExecutionStatus enum
	var executionStatus eventsv2.ExecutionStatus
	switch status {
	case "completed", "completed_early_exit":
		executionStatus = eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED
	case "errored", "timeout":
		executionStatus = eventsv2.ExecutionStatus_EXECUTION_STATUS_FAILED
	default:
		executionStatus = eventsv2.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED
	}

	var errMsg string
	if capErr != nil {
		errMsg = capErr.Error()
	}

	v2Event := &eventsv2.CapabilityExecutionFinished{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: executionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		CapabilityID:        capabilityID,
		StepRef:             int32(stepRefInt),
		Status:              executionStatus,
		Error:               errMsg,
	}

	// Emit both v1 and v2 events
	var multiErr error
	if err := emitProtoMessage(ctx, event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	if err := emitProtoMessage(ctx, v2Event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}
	return multiErr
}

func EmitMeteringReport(ctx context.Context, labels map[string]string, rpt *events.MeteringReport) error {
	rpt.Metadata = buildWorkflowMetadata(labels, labels[platform.KeyWorkflowExecutionID])

	return emitProtoMessage(ctx, rpt)
}

func EmitUserLogs(ctx context.Context, labels map[string]string, logLines []*events.LogLine, executionID string) error {
	metadata := buildWorkflowMetadata(labels, executionID)
	event := &events.UserLogs{
		M:        metadata,
		LogLines: logLines,
	}

	// Also emit v2 events - one per log line
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// Emit v1 event
	var multiErr error
	if err := emitProtoMessage(ctx, event); err != nil {
		multiErr = errors.Join(multiErr, err)
	}

	// Emit v2 events - one per log line
	for _, logLine := range logLines {
		v2Event := &eventsv2.WorkflowUserLog{
			CreInfo:             creInfo,
			Workflow:            workflowKey,
			WorkflowExecutionID: executionID,
			Timestamp:           logLine.NodeTimestamp,
			Msg:                 logLine.Message,
			Labels:              make(map[string]string), // Empty for now
		}

		if err := emitProtoMessage(ctx, v2Event); err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	return multiErr
}

// GenerateExecutionID generates a deterministic execution ID from workflowID and triggerEventID
// hash of (workflowID, triggerEventID)
func GenerateExecutionID(workflowID, triggerEventID string) (string, error) {
	s := sha256.New()
	_, err := s.Write([]byte(workflowID))
	if err != nil {
		return "", err
	}

	_, err = s.Write([]byte(triggerEventID))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
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
	case *events.WorkflowStatusChanged:
		schema = SchemaWorkflowStatusChanged
		entity = fmt.Sprintf("%s.%s", ProtoPkg, WorkflowStatusChanged)
	case *events.UserLogs:
		schema = SchemaUserLogs
		entity = fmt.Sprintf("%s.%s", ProtoPkg, UserLogs)
	// V2 event types
	case *eventsv2.WorkflowExecutionStarted:
		schema = SchemaWorkflowStartedV2
		entity = "workflows.v2." + WorkflowExecutionStarted
	case *eventsv2.WorkflowExecutionFinished:
		schema = SchemaWorkflowFinishedV2
		entity = "workflows.v2." + WorkflowExecutionFinished
	case *eventsv2.CapabilityExecutionStarted:
		schema = SchemaCapabilityStartedV2
		entity = "workflows.v2." + CapabilityExecutionStarted
	case *eventsv2.CapabilityExecutionFinished:
		schema = SchemaCapabilityFinishedV2
		entity = "workflows.v2." + CapabilityExecutionFinished
	case *eventsv2.TriggerExecutionStarted:
		schema = SchemaTriggerStartedV2
		entity = "workflows.v2." + TriggerExecutionStarted
	case *eventsv2.WorkflowUserLog:
		schema = SchemaUserLogsV2
		entity = "workflows.v2." + WorkflowUserLog
	case *eventsv2.WorkflowActivated:
		schema = SchemaWorkflowActivatedV2
		entity = "workflows.v2." + WorkflowActivated
	case *eventsv2.WorkflowPaused:
		schema = SchemaWorkflowPausedV2
		entity = "workflows.v2." + WorkflowPaused
	case *eventsv2.WorkflowDeleted:
		schema = SchemaWorkflowDeletedV2
		entity = "workflows.v2." + WorkflowDeleted
	default:
		return fmt.Errorf("unknown message type: %T", msg)
	}

	return beholder.GetEmitter().Emit(ctx, b,
		"beholder_data_schema", schema, // required
		"beholder_domain", "platform", // required
		"beholder_entity", entity) // required
}

// buildWorkflowMetadata populates a WorkflowMetadata from kvs (map[string]string).
func buildWorkflowMetadata(kvs map[string]string, workflowExecutionID string) *events.WorkflowMetadata {
	m := &events.WorkflowMetadata{}

	m.WorkflowOwner = kvs[platform.KeyWorkflowOwner]
	m.WorkflowName = kvs[platform.KeyWorkflowName]
	m.Version = kvs[platform.KeyWorkflowVersion]
	m.WorkflowID = kvs[platform.KeyWorkflowID]
	m.WorkflowExecutionID = workflowExecutionID

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

// buildCREMetadataV2 populates a CREInfo from kvs (map[string]string).
func buildCREMetadataV2(kvs map[string]string) *eventsv2.CreInfo {
	m := &eventsv2.CreInfo{}

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

	m.WorkflowRegistryAddress = kvs[platform.WorkflowRegistryAddress]
	m.WorkflowRegistryVersion = kvs[platform.WorkflowRegistryVersion]
	m.WorkflowRegistryChain = kvs[platform.WorkflowRegistryChainSelector]
	m.EngineVersion = kvs[platform.EngineVersion]
	m.CapabilitiesRegistryVersion = kvs[platform.CapabilitiesRegistryVersion]
	m.DonVersion = kvs[platform.DonVersion]

	return m
}

// buildWorkflowKeyV2 populates a WorkflowKey from kvs (map[string]string).
func buildWorkflowKeyV2(kvs map[string]string) *eventsv2.WorkflowKey {
	w := &eventsv2.WorkflowKey{}

	w.WorkflowOwner = kvs[platform.KeyWorkflowOwner]
	w.WorkflowName = kvs[platform.KeyWorkflowName]
	w.WorkflowID = kvs[platform.KeyWorkflowID]
	w.OrganizationID = kvs[platform.KeyOrganizationID]

	return w
}

func buildWorkflowV2(kvs map[string]string, binaryURL, configURL string) *eventsv2.Workflow {
	w := &eventsv2.Workflow{}

	w.WorkflowKey = buildWorkflowKeyV2(kvs)
	w.Version = kvs[platform.KeyWorkflowVersion]
	w.BinaryURL = binaryURL
	w.ConfigURL = configURL

	return w
}
