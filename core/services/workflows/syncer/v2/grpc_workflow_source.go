package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	nodeauthjwt "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/grpcsource"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/sources"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const (
	// GRPCWorkflowSourceName is the name used for logging and identification.
	GRPCWorkflowSourceName = "GRPCWorkflowSource"

	// Default configuration values
	defaultPageSize       int64         = 1000
	defaultMaxRetries     int           = 2
	defaultRetryBaseDelay time.Duration = 100 * time.Millisecond
	defaultRetryMaxDelay  time.Duration = 5 * time.Second
)

// grpcClient is an interface for the GRPC client to enable testing.
type grpcClient interface {
	ListWorkflowMetadata(ctx context.Context, families []string, start, limit int64) ([]*pb.WorkflowMetadata, bool, error)
	Close() error
}

// GRPCWorkflowSource implements WorkflowMetadataSource by fetching from a GRPC server.
// This enables external systems to provide workflows for deployment.
type GRPCWorkflowSource struct {
	lggr           logger.Logger
	client         grpcClient
	name           string
	pageSize       int64
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
	mu             sync.RWMutex
	ready          bool
}

// GRPCWorkflowSourceConfig holds configuration for creating a GRPCWorkflowSource.
type GRPCWorkflowSourceConfig struct {
	// URL is the GRPC server address (e.g., "localhost:50051")
	URL string
	// Name is a human-readable identifier for this source
	Name string
	// TLSEnabled determines whether to use TLS for the connection
	TLSEnabled bool
	// JWTGenerator is the JWT generator for authentication (always enabled, matching billing/storage pattern)
	JWTGenerator nodeauthjwt.JWTGenerator
	// PageSize is the number of workflows to fetch per page (default: 1000)
	PageSize int64
	// MaxRetries is the maximum number of retry attempts for transient errors (default: 2)
	MaxRetries int
	// RetryBaseDelay is the base delay for exponential backoff (default: 100ms)
	RetryBaseDelay time.Duration
	// RetryMaxDelay is the maximum delay between retries (default: 5s)
	RetryMaxDelay time.Duration
}

// NewGRPCWorkflowSource creates a new GRPC-based workflow source.
// The Name field in cfg is required and must be unique across all workflow sources.
func NewGRPCWorkflowSource(lggr logger.Logger, cfg GRPCWorkflowSourceConfig) (*GRPCWorkflowSource, error) {
	if cfg.Name == "" {
		return nil, errors.New("source name is required")
	}
	if cfg.URL == "" {
		return nil, errors.New("GRPC URL is required")
	}

	// Build client options - JWT auth is always enabled
	clientOpts := []grpcsource.ClientOption{
		grpcsource.WithTLS(cfg.TLSEnabled),
	}
	if cfg.JWTGenerator != nil {
		clientOpts = append(clientOpts, grpcsource.WithJWTGenerator(cfg.JWTGenerator))
	}

	client, err := grpcsource.NewClient(cfg.URL, cfg.Name, clientOpts...)
	if err != nil {
		return nil, err
	}

	return newGRPCWorkflowSourceWithClient(lggr, client, cfg)
}

// NewGRPCWorkflowSourceWithClient creates a new GRPC-based workflow source with an injected client.
// This is useful for testing with mock clients.
func NewGRPCWorkflowSourceWithClient(lggr logger.Logger, client grpcClient, cfg GRPCWorkflowSourceConfig) (*GRPCWorkflowSource, error) {
	return newGRPCWorkflowSourceWithClient(lggr, client, cfg)
}

func newGRPCWorkflowSourceWithClient(lggr logger.Logger, client grpcClient, cfg GRPCWorkflowSourceConfig) (*GRPCWorkflowSource, error) {
	if cfg.Name == "" {
		return nil, errors.New("source name is required")
	}

	pageSize := cfg.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}

	retryBaseDelay := cfg.RetryBaseDelay
	if retryBaseDelay <= 0 {
		retryBaseDelay = defaultRetryBaseDelay
	}

	retryMaxDelay := cfg.RetryMaxDelay
	if retryMaxDelay <= 0 {
		retryMaxDelay = defaultRetryMaxDelay
	}

	return &GRPCWorkflowSource{
		lggr:           logger.Named(lggr, cfg.Name),
		client:         client,
		name:           cfg.Name,
		pageSize:       pageSize,
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
		ready:          true,
	}, nil
}

// ListWorkflowMetadata fetches workflow metadata from the GRPC source.
// Pagination is handled internally - this method fetches all pages and returns all workflows.
// Transient errors (Unavailable, ResourceExhausted) are retried with exponential backoff.
// Returns a synthetic head since GRPC sources don't have blockchain state.
func (g *GRPCWorkflowSource) ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	g.tryInitialize(ctx)

	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.ready {
		return nil, nil, errors.New("GRPC source not ready")
	}

	var allViews []WorkflowMetadataView
	var start int64

	// Fetch all pages
	for {
		workflows, hasMore, err := g.fetchPageWithRetry(ctx, don.Families, start)
		if err != nil {
			return nil, nil, err
		}

		// Convert workflows to views, skipping invalid ones
		for _, wf := range workflows {
			view, err := g.toWorkflowMetadataView(wf)
			if err != nil {
				g.lggr.Warnw("Failed to parse workflow metadata, skipping",
					"workflowName", wf.GetWorkflowName(),
					"error", err)
				continue
			}
			allViews = append(allViews, view)
		}

		// Check if we've fetched all pages
		if !hasMore {
			break
		}

		// Move to next page
		start += g.pageSize
	}

	g.lggr.Debugw("Loaded workflows from GRPC source",
		"count", len(allViews),
		"donID", don.ID,
		"donFamilies", don.Families)

	return allViews, g.syntheticHead(), nil
}

// fetchPageWithRetry fetches a single page with retry logic for transient errors.
func (g *GRPCWorkflowSource) fetchPageWithRetry(ctx context.Context, families []string, start int64) ([]*pb.WorkflowMetadata, bool, error) {
	var lastErr error

	for attempt := 0; attempt <= g.maxRetries; attempt++ {
		// Check context before making request
		if ctx.Err() != nil {
			return nil, false, ctx.Err()
		}

		workflows, hasMore, err := g.client.ListWorkflowMetadata(ctx, families, start, g.pageSize)
		if err == nil {
			return workflows, hasMore, nil
		}

		lastErr = err

		// Check if this is a retryable error
		if !g.isRetryableError(err) {
			g.lggr.Errorw("Non-retryable error from GRPC source",
				"error", err,
				"start", start,
				"pageSize", g.pageSize)
			return nil, false, err
		}

		// Log retry attempt
		g.lggr.Warnw("Retryable error from GRPC source",
			"error", err,
			"attempt", attempt+1,
			"maxRetries", g.maxRetries)

		// If we've exhausted retries, return the error
		if attempt >= g.maxRetries {
			g.lggr.Errorw("Max retries exceeded for GRPC request",
				"error", err,
				"maxRetries", g.maxRetries)
			return nil, false, fmt.Errorf("max retries exceeded: %w", err)
		}

		// Calculate backoff with jitter
		backoff := g.calculateBackoff(attempt + 1)

		// Wait for backoff or context cancellation
		select {
		case <-ctx.Done():
			return nil, false, ctx.Err()
		case <-time.After(backoff):
			g.lggr.Debugw("Retrying GRPC request",
				"attempt", attempt+1,
				"delay", backoff,
				"lastError", lastErr)
		}
	}

	return nil, false, lastErr
}

// isRetryableError determines if an error should be retried.
func (g *GRPCWorkflowSource) isRetryableError(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable, codes.ResourceExhausted:
		return true
	default:
		return false
	}
}

// calculateBackoff calculates the backoff duration for a given attempt with jitter.
func (g *GRPCWorkflowSource) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^(attempt-1)
	backoff := g.retryBaseDelay * time.Duration(1<<(attempt-1))

	// Apply jitter (0.5 to 1.5 multiplier) - math/rand/v2 is auto-seeded and concurrent-safe
	jitter := 0.5 + rand.Float64() //nolint:gosec // G404: weak random is fine for retry jitter
	backoff = time.Duration(float64(backoff) * jitter)

	// Cap at max delay
	if backoff > g.retryMaxDelay {
		backoff = g.retryMaxDelay
	}

	return backoff
}

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

// tryInitialize returns the current ready state (GRPC client initialized in constructor).
func (g *GRPCWorkflowSource) tryInitialize(_ context.Context) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.ready
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
		return WorkflowMetadataView{}, fmt.Errorf("workflow_id must be 32 bytes, got %d", len(workflowIDBytes))
	}
	var workflowID types.WorkflowID
	copy(workflowID[:], workflowIDBytes)

	// Parse owner from hex string (proto changed from bytes to string)
	ownerStr := wf.GetOwner()
	ownerBytes, err := hex.DecodeString(strings.TrimPrefix(ownerStr, "0x"))
	if err != nil {
		return WorkflowMetadataView{}, fmt.Errorf("invalid owner hex string: %w", err)
	}

	// Get attributes directly (already bytes in proto)
	attributes := wf.GetAttributes()

	// Map proto status enum to internal representation
	statusVal := GRPCStatusToInternal(wf.GetStatus(), g.lggr)

	return WorkflowMetadataView{
		WorkflowID:   workflowID,
		Owner:        ownerBytes,
		CreatedAt:    wf.GetCreatedAt(),
		Status:       statusVal,
		WorkflowName: wf.GetWorkflowName(),
		BinaryURL:    wf.GetBinaryUrl(),
		ConfigURL:    wf.GetConfigUrl(),
		Tag:          wf.GetTag(),
		Attributes:   attributes,
		DonFamily:    wf.GetDonFamily(),
		Source:       g.SourceIdentifier(),
	}, nil
}

// SourceIdentifier returns the source identifier used in WorkflowMetadataView.Source.
// Format: grpc:{source_name}:v1
func (g *GRPCWorkflowSource) SourceIdentifier() string {
	return fmt.Sprintf("grpc:%s:v1", g.name)
}

// syntheticHead returns a synthetic head for GRPC sources.
// GRPC sources don't have blockchain state, so we generate a synthetic head
// with the current timestamp for consistency with the WorkflowMetadataSource interface.
func (g *GRPCWorkflowSource) syntheticHead() *commontypes.Head {
	now := time.Now().Unix()
	var timestamp uint64
	if now >= 0 {
		timestamp = uint64(now)
	}
	return &commontypes.Head{
		Height:    strconv.FormatInt(now, 10),
		Hash:      []byte("grpc-source"),
		Timestamp: timestamp,
	}
}
