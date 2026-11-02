package events

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/durableemitter"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

// Emitter emits workflow platform events to beholder and optional durable sinks.
type Emitter struct {
	beholder beholder.Emitter
	durable  func(ctx context.Context, body []byte, attrKVs ...any) error
}

// EmitterOption configures an Emitter.
type EmitterOption func(*Emitter)

// WithDurableEmit sets the durable emit func. Nil keeps the default.
func WithDurableEmit(fn func(ctx context.Context, body []byte, attrKVs ...any) error) EmitterOption {
	return func(e *Emitter) {
		if fn != nil {
			e.durable = fn
		}
	}
}

// WithNoDurableEmit disables durable emission.
func WithNoDurableEmit() EmitterOption {
	return WithDurableEmit(func(context.Context, []byte, ...any) error { return nil })
}

// NewEmitter returns an Emitter writing to e. Nil defaults to beholder.GlobalEmitter().
func NewEmitter(e beholder.Emitter, opts ...EmitterOption) *Emitter {
	if e == nil {
		e = beholder.GlobalEmitter()
	}
	em := &Emitter{beholder: e, durable: durableemitter.GlobalEmit}
	for _, opt := range opts {
		opt(em)
	}
	return em
}

// DefaultEmitter returns an Emitter configured with global sinks.
func DefaultEmitter() *Emitter {
	return NewEmitter(nil)
}

func (e *Emitter) EmitWorkflowStatusChangedEventV2(
	ctx context.Context,
	labels map[string]string,
	head *types.Head,
	status string,
	binaryURL string,
	configURL string,
	eventErr error,
) error {
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

	return e.emitProtoMessage(ctx, v2Event)
}

func (e *Emitter) EmitWorkflowActivationAbandonedV2(
	ctx context.Context,
	labels map[string]string,
	binaryURL string,
	configURL string,
	reason eventsv2.ActivationAbandonReason,
	activationErr error,
	retryCount int32,
) error {
	var errorMessage string
	if activationErr != nil {
		errorMessage = activationErr.Error()
	}

	event := &eventsv2.WorkflowActivationAbandoned{
		CreInfo:      buildCREMetadataV2(labels),
		Workflow:     buildWorkflowV2(labels, binaryURL, configURL),
		Timestamp:    uint64(time.Now().Unix()), //nolint:gosec // G115: unix timestamp is non-negative
		ErrorMessage: errorMessage,
		Reason:       reason,
		RetryCount:   retryCount,
	}

	return e.emitProtoMessage(ctx, event)
}

func (e *Emitter) EmitExecutionStartedEvent(
	ctx context.Context,
	labels map[string]string,
	triggerEventID string,
	executionID string,
) error {
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	v2Event := &eventsv2.WorkflowExecutionStarted{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: executionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		TriggerID:           triggerEventID,
	}

	return e.emitProtoMessage(ctx, v2Event)
}

// ExecutionProfileStepInput describes timing for a single capability step invocation.
type ExecutionProfileStepInput struct {
	StepID       string
	StartTime    time.Time
	EndTime      time.Time // zero means use the execution end time
	CapabilityID string
	HasError     bool
}

func buildExecutionProfile(
	workflowID, executionID string,
	startTime, endTime time.Time,
	status string,
	steps []ExecutionProfileStepInput,
) *eventsv2.ExecutionProfile {
	stepProfiles := make([]*eventsv2.ExecutionProfileStep, 0, len(steps))
	for _, step := range steps {
		stepEnd := step.EndTime
		if stepEnd.IsZero() {
			stepEnd = endTime
		}
		stepProfiles = append(stepProfiles, &eventsv2.ExecutionProfileStep{
			StepID:       step.StepID,
			StartTime:    step.StartTime.UTC().Format(time.RFC3339Nano),
			EndTime:      stepEnd.UTC().Format(time.RFC3339Nano),
			CapabilityID: step.CapabilityID,
			HasError:     step.HasError,
		})
	}

	return &eventsv2.ExecutionProfile{
		WorkflowID:          workflowID,
		WorkflowExecutionID: executionID,
		StartTime:           startTime.UTC().Format(time.RFC3339Nano),
		EndTime:             endTime.UTC().Format(time.RFC3339Nano),
		Status:              status,
		Steps:               stepProfiles,
	}
}

func (e *Emitter) EmitExecutionProfile(
	ctx context.Context,
	workflowID, executionID string,
	startTime, endTime time.Time,
	status string,
	steps []ExecutionProfileStepInput,
) (*eventsv2.ExecutionProfile, error) {
	profile := buildExecutionProfile(workflowID, executionID, startTime, endTime, status, steps)
	return profile, e.emitProtoMessage(ctx, profile)
}

// ErrorClassification attributes a failed workflow execution to its root cause:
// the user's workflow (code, config, or a returned error) vs the platform
// (timeouts, unavailable capabilities, internal engine errors).
type ErrorClassification int

const (
	// ErrorClassificationUnspecified means the cause was not determined (e.g. the
	// legacy v1 engine, which does not classify). Emitted as UNSPECIFIED.
	ErrorClassificationUnspecified ErrorClassification = iota
	// ErrorClassificationUser attributes the failure to the user's workflow.
	ErrorClassificationUser
	// ErrorClassificationSystem attributes the failure to the platform.
	ErrorClassificationSystem
)

// ClassifyError resolves an execution error to a user/system classification.
// A caperrors.Error carrying an explicit origin wins (so user-origin errors that
// propagate up from a capability or the guest are attributed correctly); other
// errors fall back to the caller-supplied classification.
func ClassifyError(execErr error, fallback ErrorClassification) ErrorClassification {
	if execErr == nil {
		return ErrorClassificationUnspecified
	}
	if capErr, ok := errors.AsType[caperrors.Error](execErr); ok {
		switch capErr.Origin() {
		case caperrors.OriginUser:
			return ErrorClassificationUser
		case caperrors.OriginSystem:
			return ErrorClassificationSystem
		}
	}
	return fallback
}

// classifiedExecutionStatus folds the terminal status and the error
// classification into the v2 ClassifiedExecutionStatus enum.
func classifiedExecutionStatus(status eventsv2.ExecutionStatus, errClass ErrorClassification) eventsv2.ClassifiedExecutionStatus {
	switch status {
	case eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED:
		return eventsv2.ClassifiedExecutionStatus_CLASSIFIED_EXECUTION_STATUS_SUCCEEDED
	case eventsv2.ExecutionStatus_EXECUTION_STATUS_FAILED:
		switch errClass {
		case ErrorClassificationUser:
			return eventsv2.ClassifiedExecutionStatus_CLASSIFIED_EXECUTION_STATUS_USER_ERROR
		case ErrorClassificationSystem:
			return eventsv2.ClassifiedExecutionStatus_CLASSIFIED_EXECUTION_STATUS_SYSTEM_ERROR
		default:
			return eventsv2.ClassifiedExecutionStatus_CLASSIFIED_EXECUTION_STATUS_UNSPECIFIED
		}
	default:
		return eventsv2.ClassifiedExecutionStatus_CLASSIFIED_EXECUTION_STATUS_UNSPECIFIED
	}
}

func (e *Emitter) EmitExecutionFinishedEvent(ctx context.Context, labels map[string]string, status, executionID string, execErr error, errClass ErrorClassification, lggr logger.Logger) error {
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

	var errMsg string
	if execErr != nil {
		errMsg = execErr.Error()
	}

	v2Event := &eventsv2.WorkflowExecutionFinished{
		CreInfo:             creInfo,
		Workflow:            workflowKey,
		WorkflowExecutionID: executionID,
		Timestamp:           time.Now().Format(time.RFC3339),
		Status:              executionStatus,
		Error:               errMsg,
		ClassifiedStatus:    classifiedExecutionStatus(executionStatus, errClass),
	}

	return e.emitProtoMessage(ctx, v2Event)
}

func (e *Emitter) EmitCapabilityStartedEvent(ctx context.Context, labels map[string]string, executionID, capabilityID, stepRef, method string) error {
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// Convert stepRef string to int32; non-numeric refs map to -1
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
		Method:              method,
	}

	return e.emitProtoMessage(ctx, v2Event)
}

func (e *Emitter) EmitTriggerExecutionStarted(ctx context.Context, labels map[string]string, triggerID, workflowExecutionID string) error {
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

	return e.emitProtoMessage(ctx, v2Event)
}

func (e *Emitter) EmitCapabilityFinishedEvent(ctx context.Context, labels map[string]string, executionID, capabilityID, stepRef, status, method string, capErr error) error {
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// Convert stepRef string to int32; non-numeric refs map to -1
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
		Method:              method,
	}

	return e.emitProtoMessage(ctx, v2Event)
}

func (e *Emitter) EmitMeteringReport(ctx context.Context, labels map[string]string, rpt *events.MeteringReport) error {
	rpt.Metadata = buildWorkflowMetadata(labels, labels[platform.KeyWorkflowExecutionID])

	return e.emitProtoMessage(ctx, rpt)
}

func (e *Emitter) EmitUserLogs(ctx context.Context, labels map[string]string, logLines []*events.LogLine, executionID string) error {
	creInfo := buildCREMetadataV2(labels)
	workflowKey := buildWorkflowKeyV2(labels)

	// One event per log line
	var multiErr error
	for _, logLine := range logLines {
		v2Event := &eventsv2.WorkflowUserLog{
			CreInfo:             creInfo,
			Workflow:            workflowKey,
			WorkflowExecutionID: executionID,
			Timestamp:           logLine.NodeTimestamp,
			Msg:                 logLine.Message,
			Labels:              make(map[string]string), // Empty for now
		}

		if err := e.emitProtoMessage(ctx, v2Event); err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	return multiErr
}

func (e *Emitter) EmitUserMetric(ctx context.Context, labels map[string]string, metric *eventsv2.WorkflowUserMetric) error {
	metric.CreInfo = buildCREMetadataV2(labels)
	metric.Workflow = buildWorkflowKeyV2(labels)
	metric.Timestamp = time.Now().Format(time.RFC3339)

	return e.emitProtoMessage(ctx, metric)
}

func (e *Emitter) emitProtoMessage(ctx context.Context, msg proto.Message) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Determine the schema and entity based on the message type
	// entity must be prefixed with the proto package name
	var schema, entity string
	switch msg.(type) {
	case *events.MeteringReport:
		schema = MeteringReportSchema
		entity = fmt.Sprintf("%s.%s", ProtoPkg, MeteringReportEntity)
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
	case *eventsv2.WorkflowUserMetric:
		schema = SchemaUserMetricV2
		entity = "workflows.v2." + WorkflowUserMetric
	case *eventsv2.WorkflowActivated:
		schema = SchemaWorkflowActivatedV2
		entity = "workflows.v2." + WorkflowActivated
	case *eventsv2.WorkflowActivationAbandoned:
		schema = SchemaWorkflowActivationAbandonedV2
		entity = "workflows.v2." + WorkflowActivationAbandoned
	case *eventsv2.WorkflowPaused:
		schema = SchemaWorkflowPausedV2
		entity = "workflows.v2." + WorkflowPaused
	case *eventsv2.WorkflowDeleted:
		schema = SchemaWorkflowDeletedV2
		entity = "workflows.v2." + WorkflowDeleted
	case *eventsv2.ExecutionProfile:
		schema = SchemaWorkflowExecutionProfileV2
		entity = "workflows.v2." + WorkflowExecutionProfile
	default:
		return fmt.Errorf("unknown message type: %T", msg)
	}

	return e.emitRawMessage(ctx, b, schema, entity)
}

func (e *Emitter) emitRawMessage(ctx context.Context, body []byte, schema, entity string) error {
	if err := e.beholder.Emit(ctx, body,
		"beholder_data_schema", schema, // required
		"beholder_domain", "platform", // required
		"beholder_entity", entity); err != nil { // required
		return err
	}

	err := e.durable(ctx, body, "source", "platform", "type", entity)
	if err != nil && !errors.Is(err, durableemitter.ErrNotInitialized) {
		return err
	}
	return nil
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

	m.OrgID = kvs[platform.KeyOrganizationID]

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

	if id, err := strconv.ParseInt(kvs[platform.KeyWorkflowDonID], 10, 32); err == nil {
		m.WorkflowDonID = int32(id)
	}
	if idx, err := strconv.ParseUint(kvs[platform.KeyShardIndex], 10, 32); err == nil {
		m.ShardIndex = uint32(idx)
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
	m.WorkflowSource = kvs[platform.KeyWorkflowSource]

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
