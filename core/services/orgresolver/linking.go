package orgresolver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	log "github.com/smartcontractkit/chainlink-common/pkg/logger"
	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type Config struct {
	TTL             time.Duration
	RefreshStrategy RefreshStrategy
	URL             string
	TLSEnabled      bool
}

// Allows the caller control
// but in another version, we don't allow this
// never seen this before, unsure why. seems like
// it could be beneficial to allow a client to force its way through
// when needed
type CallerConfig struct {
	allowedStale time.Duration
}

type OrgResolver interface {
	services.ServiceCtx
	Get(ctx context.Context, owner string) (string, error)
}

type orgAndRefresh struct {
	orgID       string
	lastRefresh time.Time
}

type RefreshStrategy int

const (
	RefreshStrategyIndividual RefreshStrategy = iota // Smooth individual calls
	RefreshStrategyBulk                              // Single bulk call
)

// orgResolver is intended to work like so:
// on the first call to Get, the owner DNE in ownersMap, and cacheMiss is called.
// cacheMiss will use the configured linkingclient, which in production will be a network call
// Over time, owner addresses will collect in ownersMap.
// The cacheWorker, every tick duration of ttl, will refresh it the cache for the owners collected.
//

// we have cache strategies to flexibly handle the traffic to the linking service
// idea is to avoid overwhelming the linking service with spikes

// current issues are:
// when node is booted up, it will result in a direct call to the linking service for each owner
// not ideal, TODO: warming strategy
type orgResolver struct {
	ttl             time.Duration
	ownersMap       map[string]*orgAndRefresh
	refreshStrategy RefreshStrategy

	URL        string
	TLSEnabled bool

	client linkingclient.LinkingServiceClient
	clock  clockwork.Clock

	mu     sync.Mutex
	logger log.Logger

	staleEntryCount int // track refresh failures for health reporting

	shutdown chan struct{}
}

// NewOrgResolver creates a new org resolver with the specified configuration
func NewOrgResolver(cfg Config, client linkingclient.LinkingServiceClient, clock clockwork.Clock, logger log.Logger) *orgResolver {
	return &orgResolver{
		ttl:             cfg.TTL,
		refreshStrategy: cfg.RefreshStrategy,
		URL:             cfg.URL,
		TLSEnabled:      cfg.TLSEnabled,
		client:          client,
		clock:           clock,
		logger:          logger,
		ownersMap:       make(map[string]*orgAndRefresh),
		shutdown:        make(chan struct{}),
	}
}

func (o *orgResolver) Get(ctx context.Context, owner string) (string, error) {
	o.mu.Lock()

	entry, exists := o.ownersMap[owner]
	if !exists {
		o.mu.Unlock()
		return o.cacheMiss(ctx, owner)
	}

	// Return cached value (if linking service is unavailable, can be stale)
	defer o.mu.Unlock()
	return entry.orgID, nil
}

func (o *orgResolver) cacheMiss(ctx context.Context, owner string) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	req := linkingclient.GetOrganizationFromWorkflowOwnerRequest{
		WorkflowOwner:           owner,
		WorkflowRegistryAddress: "TODO",
		ChainSelector:           0,
	}

	// Set timeout for cache miss calls to prevent blocking Get() calls
	callTimeout := o.ttl / 4 // 25% of TTL
	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	resp, err := o.client.GetOrganizationFromWorkflowOwner(callCtx, &req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch organization from workflow owner: %w", err)
	}

	o.ownersMap[owner] = &orgAndRefresh{
		orgID:       resp.OrganizationId,
		lastRefresh: o.clock.Now(),
	}

	return resp.OrganizationId, nil
}

func (o *orgResolver) Start(_ context.Context) {
	go func() {
		o.cacheWorker()
	}()
}

func (o *orgResolver) cacheWorker() {
	timer := o.clock.NewTimer(o.ttl)
	for {
		select {
		case <-timer.Chan():
			timer.Stop()
			switch o.refreshStrategy {
			case RefreshStrategyBulk:
				o.bulkRefreshCache()
			case RefreshStrategyIndividual:
				o.individualRefreshCache()
			}
		case <-o.shutdown:
			timer.Stop()
			return
		}
	}
}

func (o *orgResolver) HealthReport() map[string]error {
	o.mu.Lock()
	failureCount := o.staleEntryCount // now tracks refresh failures
	totalCount := len(o.ownersMap)
	o.mu.Unlock()

	if failureCount > 0 {
		return map[string]error{
			"OrgResolver": fmt.Errorf("%d/%d cache entries failed to refresh due to linking service issues", failureCount, totalCount),
		}
	}
	return map[string]error{"OrgResolver": nil}
}

// bulkRefreshCache uses a single bulk endpoint (no smoothing needed)
func (o *orgResolver) bulkRefreshCache() {
	now := o.clock.Now()

	// Collect all owners to refresh
	o.mu.Lock()
	owners := make([]string, 0, len(o.ownersMap))
	for owner := range o.ownersMap {
		owners = append(owners, owner)
	}
	o.mu.Unlock()

	if len(owners) == 0 {
		return
	}

	// TODO: Replace with actual bulk endpoint call
	// When implementing, use timeout of o.ttl / 2 (50% of TTL) for the bulk call:
	// callTimeout := o.ttl / 2
	// callCtx, cancel := context.WithTimeout(context.Background(), callTimeout)
	// defer cancel()
	// resp, err := o.client.GetOrganizationsFromWorkflowOwners(callCtx, &bulkReq{WorkflowOwners: owners, ...})
	o.logger.Debugf("Bulk refresh not yet implemented - would refresh %d owners in single call", len(owners))

	// Placeholder: mark all as successfully refreshed
	o.mu.Lock()
	for owner := range o.ownersMap {
		if entry := o.ownersMap[owner]; entry != nil {
			entry.lastRefresh = now
		}
	}
	o.staleEntryCount = 0 // no failures in bulk refresh
	o.mu.Unlock()
}

// individualRefreshCache makes individual calls with smoothing for current single-owner endpoint
func (o *orgResolver) individualRefreshCache() {
	now := o.clock.Now()

	// Collect all owners to refresh
	o.mu.Lock()
	owners := make([]string, 0, len(o.ownersMap))
	for owner := range o.ownersMap {
		owners = append(owners, owner)
	}
	o.mu.Unlock()

	if len(owners) == 0 {
		return
	}

	// Spread calls over time to smooth traffic
	// Use 60% of TTL as spread window (e.g., 6 seconds for 10s TTL)
	spreadWindow := time.Duration(float64(o.ttl) * 0.6)
	callInterval := spreadWindow / time.Duration(len(owners))
	if callInterval < 10*time.Millisecond {
		callInterval = 10 * time.Millisecond // minimum interval
	}

	updates := make(map[string]*orgAndRefresh)
	failureCount := 0

	// Set timeout for all calls to ensure they complete well before next refresh cycle
	callTimeout := o.ttl / 4 // 25% of TTL

	for i, owner := range owners {
		// Add delay between calls to spread them over time
		if i > 0 {
			time.Sleep(callInterval)
		}

		req := linkingclient.GetOrganizationFromWorkflowOwnerRequest{
			WorkflowOwner:           owner,
			WorkflowRegistryAddress: "TODO",
			ChainSelector:           0,
		}

		callCtx, cancel := context.WithTimeout(context.Background(), callTimeout)
		resp, err := o.client.GetOrganizationFromWorkflowOwner(callCtx, &req)
		cancel() // immediately cancel after call

		if err != nil {
			o.logger.Warnf("failed to refresh organization for owner %s: %v", owner, err)
			failureCount++
			continue
		}
		updates[owner] = &orgAndRefresh{
			orgID:       resp.OrganizationId,
			lastRefresh: now,
		}
	}

	// Apply successful updates
	o.mu.Lock()
	for owner, entry := range updates {
		o.ownersMap[owner] = entry
	}
	o.staleEntryCount = failureCount
	o.mu.Unlock()

	o.logger.Debugf("Individual refresh: %d/%d cache entries over %v, %d failures", len(updates), len(owners), spreadWindow, failureCount)
}

func (o *orgResolver) Close() error {
	close(o.shutdown)
	return nil
}

func (o *orgResolver) Name() string {
	return "OrgResolver"
}

func (o *orgResolver) Ready() error {
	return nil
}
