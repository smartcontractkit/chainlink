package v2

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/sources"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type SyncStrategy string

const (
	SyncStrategyReconciliation = "reconciliation"
	defaultSyncStrategy        = SyncStrategyReconciliation
)

// Workflow status values. These match the on-chain contract status values
// (0=Active, 1=Paused) to avoid any translation errors.
const (
	WorkflowStatusActive uint8 = 0
	WorkflowStatusPaused uint8 = 1
)

// ContractStatusToInternal converts on-chain contract status to internal representation.
// Contract and internal use the same values (0=Active, 1=Paused).
func ContractStatusToInternal(s uint8) uint8 {
	switch s {
	case WorkflowStatusActive, WorkflowStatusPaused:
		return s
	default:
		// Unknown status defaults to paused
		return WorkflowStatusPaused
	}
}

// FileStatusToInternal converts file source status to internal representation.
// File format uses same values as contract (0=Active, 1=Paused).
func FileStatusToInternal(s uint8) uint8 {
	return ContractStatusToInternal(s)
}

// GRPCStatusToInternal converts proto WorkflowStatus enum to internal representation.
// Proto uses: UNSPECIFIED=0, ACTIVE=1, PAUSED=2
// Internal uses: Active=0, Paused=1
func GRPCStatusToInternal(s pb.WorkflowStatus, lggr logger.Logger) uint8 {
	switch s {
	case pb.WorkflowStatus_WORKFLOW_STATUS_ACTIVE:
		return WorkflowStatusActive
	case pb.WorkflowStatus_WORKFLOW_STATUS_PAUSED:
		return WorkflowStatusPaused
	case pb.WorkflowStatus_WORKFLOW_STATUS_UNSPECIFIED:
		lggr.Warn("Received WORKFLOW_STATUS_UNSPECIFIED from proto, treating as paused")
		return WorkflowStatusPaused
	default:
		lggr.Warnw("Unknown proto status, treating as paused", "status", s)
		return WorkflowStatusPaused
	}
}

type Head struct {
	Hash      string
	Height    string
	Timestamp uint64
}

type Config struct {
	QueryCount   uint64
	SyncStrategy SyncStrategy
}

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error)

type GetActiveAllowlistedRequestsReverseParams struct {
	EndIndex   *big.Int
	StartIndex *big.Int
}

type WorkflowMetadataView struct {
	WorkflowID   types.WorkflowID
	Owner        []byte
	CreatedAt    uint64
	Status       uint8
	WorkflowName string
	BinaryURL    string
	ConfigURL    string
	Tag          string
	Attributes   []byte
	DonFamily    string
	// Source identifies where this workflow metadata came from.
	// Format varies by source type:
	//   - Onchain contract: "contract:{chain_selector}:{contract_address}"
	//   - GRPC source:      "grpc:{source_name}:v1"
	//   - File source:      "file:{source_name}:v1"
	Source string
}

type GetWorkflowListByDONParams struct {
	DonFamily string
	Start     *big.Int
	Limit     *big.Int
}

type GetWorkflowListByDONReturnVal struct {
	WorkflowMetadataList []WorkflowMetadataView
}

type WorkflowRegistryEventName string

var (
	// A WorkflowRegistered event represents when a workflow is registered
	WorkflowRegistered WorkflowRegistryEventName = "WorkflowRegistered"
	// A WorkflowActivated event represents when a workflow is activated
	WorkflowActivated WorkflowRegistryEventName = "WorkflowActivated"
	// A WorkflowPaused event represents when a workflow is paused
	WorkflowPaused WorkflowRegistryEventName = "WorkflowPaused"
	// A WorkflowDeleted event represents when a workflow is deleted
	WorkflowDeleted WorkflowRegistryEventName = "WorkflowDeleted"
)

type Event struct {
	Name WorkflowRegistryEventName
	Data any
	Head Head
	Info string // additional human-readable metadata
}

// NOTE: The following types differ from gethwrappers in that they are chain agnostic definitions (owners are represented as bytes / workflow IDs might be more than bytes32)

type WorkflowRegisteredEvent struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	CreatedAt     uint64
	Status        uint8
	WorkflowName  string
	WorkflowTag   string
	BinaryURL     string
	ConfigURL     string
	Tag           string
	Attributes    []byte
	Source        string // source that provided this workflow metadata
}

type WorkflowActivatedEvent struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	CreatedAt     uint64
	Status        uint8
	WorkflowName  string
	WorkflowTag   string
	BinaryURL     string
	ConfigURL     string
	Tag           string
	Attributes    []byte
	Source        string // source that provided this workflow metadata
}

type WorkflowPausedEvent struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	CreatedAt     uint64
	Status        uint8
	WorkflowName  string
	WorkflowTag   string
	BinaryURL     string
	ConfigURL     string
	Tag           string
	Attributes    []byte
	Source        string
}

type WorkflowDeletedEvent struct {
	WorkflowID types.WorkflowID
	Source     string
}

// WorkflowMetadataSource is an interface for fetching workflow metadata from various sources.
// This abstraction allows the workflow registry syncer to aggregate workflows from multiple
// sources (e.g., on-chain contract, file-based, API-based) while treating them uniformly.
type WorkflowMetadataSource interface {
	// ListWorkflowMetadata returns all workflow metadata for the given DON.
	ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error)

	// Name returns a human-readable name for this source.
	Name() string

	// SourceIdentifier returns the source identifier used in WorkflowMetadataView.Source.
	// This identifier is used in engine registry lookups and to differeniate between wf registries in workflow events.
	SourceIdentifier() string

	// Ready returns nil if the source is ready to be queried.
	Ready() error
}
