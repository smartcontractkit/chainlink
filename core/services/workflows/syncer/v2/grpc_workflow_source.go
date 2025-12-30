package v2

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/grpcsource"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/sources/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const (
	// GRPCWorkflowSourceName is the name used for logging and identification.
	GRPCWorkflowSourceName = "GRPCWorkflowSource"
)

// GRPCWorkflowSource implements WorkflowMetadataSource by fetching from a GRPC server.
// This enables external systems to provide workflow metadata to the chainlink node.
type GRPCWorkflowSource struct {
	lggr   logger.Logger
	client *grpcsource.Client
	name   string
	mu     sync.RWMutex
	ready  bool
}

// GRPCWorkflowSourceConfig holds configuration for creating a GRPCWorkflowSource.
type GRPCWorkflowSourceConfig struct {
	// URL is the GRPC server address (e.g., "localhost:50051")
	URL string
	// Name is a human-readable identifier for this source
	Name string
	// TLSEnabled determines whether to use TLS for the connection
	TLSEnabled bool
}

// NewGRPCWorkflowSource creates a new GRPC-based workflow source.
func NewGRPCWorkflowSource(lggr logger.Logger, cfg GRPCWorkflowSourceConfig) (*GRPCWorkflowSource, error) {
	if cfg.URL == "" {
		return nil, errors.New("GRPC URL is required")
	}

	sourceName := cfg.Name
	if sourceName == "" {
		sourceName = GRPCWorkflowSourceName
	}

	client, err := grpcsource.NewClient(cfg.URL, sourceName, cfg.TLSEnabled)
	if err != nil {
		return nil, err
	}

	return &GRPCWorkflowSource{
		lggr:   lggr.Named(sourceName),
		client: client,
		name:   sourceName,
		ready:  true,
	}, nil
}

// ListWorkflowMetadata fetches workflow metadata from the GRPC source.
func (g *GRPCWorkflowSource) ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.ready {
		return nil, nil, errors.New("GRPC source not ready")
	}

	workflows, head, err := g.client.ListWorkflowMetadata(ctx, don.ID, don.Families)
	if err != nil {
		g.lggr.Errorw("Failed to fetch workflows from GRPC source", "error", err)
		return nil, nil, err
	}

	var views []WorkflowMetadataView
	for _, wf := range workflows {
		view, err := g.toWorkflowMetadataView(wf)
		if err != nil {
			g.lggr.Warnw("Failed to parse workflow metadata, skipping",
				"workflowName", wf.GetWorkflowName(),
				"error", err)
			continue
		}
		views = append(views, view)
	}

	g.lggr.Debugw("Loaded workflows from GRPC source",
		"count", len(views),
		"donID", don.ID,
		"donFamilies", don.Families)

	return views, g.toCommonHead(head), nil
}

// Name returns the name of this source.
func (g *GRPCWorkflowSource) Name() string {
	return g.name
}

// Ready returns nil if the GRPC client is connected.
func (g *GRPCWorkflowSource) Ready() error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.ready {
		return errors.New("GRPC source not ready")
	}
	return nil
}

// Close closes the underlying GRPC connection.
func (g *GRPCWorkflowSource) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.ready = false
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}

// toWorkflowMetadataView converts a protobuf WorkflowMetadata to a WorkflowMetadataView.
func (g *GRPCWorkflowSource) toWorkflowMetadataView(wf *pb.WorkflowMetadata) (WorkflowMetadataView, error) {
	// Validate workflow ID length
	workflowIDBytes := wf.GetWorkflowId()
	if len(workflowIDBytes) != 32 {
		return WorkflowMetadataView{}, errors.New("workflow_id must be 32 bytes")
	}
	var workflowID types.WorkflowID
	copy(workflowID[:], workflowIDBytes)

	// Get owner bytes directly
	ownerBytes := wf.GetOwner()

	// Get attributes directly (already bytes in proto)
	attributes := wf.GetAttributes()

	return WorkflowMetadataView{
		WorkflowID:   workflowID,
		Owner:        ownerBytes,
		CreatedAt:    wf.GetCreatedAt(),
		Status:       uint8(wf.GetStatus()),
		WorkflowName: wf.GetWorkflowName(),
		BinaryURL:    wf.GetBinaryUrl(),
		ConfigURL:    wf.GetConfigUrl(),
		Tag:          wf.GetTag(),
		Attributes:   attributes,
		DonFamily:    wf.GetDonFamily(),
	}, nil
}

// toCommonHead converts a protobuf Head to a common.Head.
func (g *GRPCWorkflowSource) toCommonHead(head *pb.Head) *commontypes.Head {
	if head == nil {
		// Return a synthetic head if none provided
		return &commontypes.Head{
			Height:    strconv.FormatInt(time.Now().Unix(), 10),
			Hash:      []byte("grpc-source"),
			Timestamp: uint64(time.Now().Unix()),
		}
	}
	return &commontypes.Head{
		Height:    head.GetHeight(),
		Hash:      []byte(head.GetHash()),
		Timestamp: head.GetTimestamp(),
	}
}
