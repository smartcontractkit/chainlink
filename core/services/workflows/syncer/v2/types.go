package v2

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type SyncStrategy string

const (
	SyncStrategyReconciliation = "reconciliation"
	defaultSyncStrategy        = SyncStrategyReconciliation
)

const (
	WorkflowStatusActive uint8 = iota
	WorkflowStatusPaused
)

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
