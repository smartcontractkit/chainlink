package shardownership

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	wfmon "github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

type ShardResolver interface {
	ResolveShard(ctx context.Context, workflowID string, ownerHex string) (donID uint32, found bool, err error)
	ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error)
}

type AllShardsResolver interface {
	ResolveAllShards(ctx context.Context, workflowID string, ownerHex string) ([]uint32, bool, error)
}

type OrgResolver interface {
	Get(ctx context.Context, owner string) (string, error)
}

func globalShardAssignmentMetrics(lggr logger.Logger) *wfmon.ShardAssignmentMetrics {
	m, err := wfmon.GlobalShardAssignmentMetrics()
	if err != nil {
		lggr.Warnw("Failed to register shard assignment metrics; continuing without assignment instrumentation", "err", err)
	}
	return m
}

func resolverOutcome(found bool, err error) string {
	switch {
	case err != nil:
		return wfmon.ShardResolverOutcomeError
	case found:
		return wfmon.ShardResolverOutcomeFound
	default:
		return wfmon.ShardResolverOutcomeNotFound
	}
}

type ringOCRShardResolver struct {
	client  shardorchestrator.ClientInterface
	lggr    logger.Logger
	metrics *wfmon.ShardAssignmentMetrics
}

func NewRingOCRShardResolver(client shardorchestrator.ClientInterface, lggr logger.Logger) ShardResolver {
	return &ringOCRShardResolver{client: client, lggr: logger.Named(lggr, "RingOCRShardResolver"), metrics: globalShardAssignmentMetrics(lggr)}
}

func (r *ringOCRShardResolver) ResolveShard(ctx context.Context, workflowID string, _ string) (uint32, bool, error) {
	start := time.Now()
	shard, found, err := r.resolveShard(ctx, workflowID)
	r.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverRing, resolverOutcome(found, err), "", time.Since(start))
	return shard, found, err
}

func (r *ringOCRShardResolver) resolveShard(ctx context.Context, workflowID string) (uint32, bool, error) {
	if r.client == nil {
		return 0, false, nil
	}
	resp, err := r.client.GetWorkflowShardMapping(ctx, []string{workflowID})
	if err != nil {
		return 0, false, err
	}
	shard, ok := resp.Mappings[workflowID]
	return shard, ok, nil
}

func (r *ringOCRShardResolver) ResolveShards(ctx context.Context, workflowIDs []string, _ []string) (map[string]uint32, error) {
	start := time.Now()
	mappings, err := r.resolveShards(ctx, workflowIDs)
	if err != nil {
		r.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverRing, wfmon.ShardResolverOutcomeError, "", time.Since(start))
		return nil, err
	}
	for _, wfID := range workflowIDs {
		outcome := wfmon.ShardResolverOutcomeNotFound
		if _, ok := mappings[wfID]; ok {
			outcome = wfmon.ShardResolverOutcomeFound
		}
		r.metrics.AddResolverRequest(ctx, wfmon.ShardResolverRing, outcome, "")
	}
	r.metrics.RecordResolverDuration(ctx, wfmon.ShardResolverRing, time.Since(start))
	return mappings, nil
}

func (r *ringOCRShardResolver) resolveShards(ctx context.Context, workflowIDs []string) (map[string]uint32, error) {
	if r.client == nil {
		return map[string]uint32{}, nil
	}
	if len(workflowIDs) == 0 {
		return map[string]uint32{}, nil
	}
	resp, err := r.client.GetWorkflowShardMapping(ctx, workflowIDs)
	if err != nil {
		return nil, fmt.Errorf("shard mapping unavailable: %w", err)
	}
	return resp.Mappings, nil
}

func (r *ringOCRShardResolver) GetRoutingResponse(ctx context.Context, workflowIDs []string) (*ringpb.GetWorkflowShardMappingResponse, error) {
	if r.client == nil {
		return nil, nil
	}
	return r.client.GetWorkflowShardMapping(ctx, workflowIDs)
}

type manualShardResolver struct {
	settings         *loop.AtomicSettings
	orgResolver      OrgResolver
	shardIndexMapper *ShardIndexMapper
	lggr             logger.Logger
	metrics          *wfmon.ShardAssignmentMetrics

	cfgMu   sync.Mutex
	lastCfg *cresettings.ShardAssignmentConfig
}

func NewManualShardResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, shardIndexMapper *ShardIndexMapper, lggr logger.Logger) ShardResolver {
	return &manualShardResolver{settings: settings, orgResolver: orgResolver, shardIndexMapper: shardIndexMapper, lggr: logger.Named(lggr, "ManualShardResolver"), metrics: globalShardAssignmentMetrics(lggr)}
}

func (m *manualShardResolver) ResolveShard(ctx context.Context, _ string, ownerHex string) (uint32, bool, error) {
	start := time.Now()
	donID, found, source, err := m.resolveShard(ctx, ownerHex)
	m.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverManual, resolverOutcome(found, err), source, time.Since(start))
	return donID, found, err
}

func (m *manualShardResolver) resolveShard(ctx context.Context, ownerHex string) (uint32, bool, string, error) {
	cfg, err := m.loadObservedConfig(ctx)
	if err != nil || cfg == nil {
		return 0, false, "", err
	}
	shardIndex, found, source, err := resolveManual(ctx, cfg, ownerHex, m.orgResolver)
	if err != nil || !found {
		return 0, found, source, err
	}
	return m.toDonID(ctx, shardIndex), true, source, nil
}

func (m *manualShardResolver) ResolveAllShards(ctx context.Context, _ string, ownerHex string) ([]uint32, bool, error) {
	start := time.Now()
	donIDs, found, source, err := m.resolveAllShards(ctx, ownerHex)
	m.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverManual, resolverOutcome(found, err), source, time.Since(start))
	return donIDs, found, err
}

func (m *manualShardResolver) resolveAllShards(ctx context.Context, ownerHex string) ([]uint32, bool, string, error) {
	cfg, err := m.loadObservedConfig(ctx)
	if err != nil || cfg == nil {
		return nil, false, "", err
	}
	shardIndices, found, source, err := resolveAllManual(ctx, cfg, ownerHex, m.orgResolver)
	if err != nil || !found {
		return nil, found, source, err
	}
	return m.toDonIDs(ctx, shardIndices), true, source, nil
}

func (m *manualShardResolver) ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	start := time.Now()
	result, err := m.resolveShards(ctx, workflowIDs, ownerHexes)
	if err != nil {
		m.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverManual, wfmon.ShardResolverOutcomeError, "", time.Since(start))
		return nil, err
	}
	m.metrics.RecordResolverDuration(ctx, wfmon.ShardResolverManual, time.Since(start))
	return result, nil
}

func (m *manualShardResolver) resolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	cfg, err := m.loadObservedConfig(ctx)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return map[string]uint32{}, nil
	}
	result := make(map[string]uint32, len(workflowIDs))
	for i, wfID := range workflowIDs {
		if i >= len(ownerHexes) {
			break
		}
		shardIndex, found, source, err := resolveManual(ctx, cfg, ownerHexes[i], m.orgResolver)
		if err != nil {
			return nil, err
		}
		m.metrics.AddResolverRequest(ctx, wfmon.ShardResolverManual, resolverOutcome(found, nil), source)
		if found {
			result[wfID] = m.toDonID(ctx, shardIndex)
		}
	}
	return result, nil
}

// toDonID translates a configured shard index into the DON ID currently
// assigned that index. Falls back to returning shardIndex unchanged when no
// donIndex is wired, or the index is out of range for the current registry
// snapshot.
func (m *manualShardResolver) toDonID(ctx context.Context, shardIndex uint32) uint32 {
	if m.shardIndexMapper == nil {
		return shardIndex
	}
	// Prevent a boot-time race where a workflow is launched before
	// the registry syncer has delivered its first snapshot.
	// NOTE: consider waiting indefinitely or put a global lock earlier
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := m.shardIndexMapper.WaitReady(waitCtx); err != nil {
		m.lggr.Warnw("shard DON index not ready, resolving with a possibly incomplete registry view", "shardIndex", shardIndex, "err", err)
	}
	don := m.shardIndexMapper.DonByShardIndex(ctx, shardIndex)
	if don == nil {
		return shardIndex
	}
	return don.ID
}

func (m *manualShardResolver) toDonIDs(ctx context.Context, shardIndices []uint32) []uint32 {
	donIDs := make([]uint32, len(shardIndices))
	for i, idx := range shardIndices {
		donIDs[i] = m.toDonID(ctx, idx)
	}
	return donIDs
}

// loadObservedConfig loads and parses the current shard assignment settings,
// recording the settings hash and emitting one diff log per actual settings
// change. A parse error keeps the previously applied assignment live.
func (m *manualShardResolver) loadObservedConfig(ctx context.Context) (*cresettings.ShardAssignmentConfig, error) {
	var raw string
	if m.settings != nil {
		update, err := m.settings.Load()
		if err != nil {
			return nil, err
		}
		raw = update.Settings
	}
	if raw == "" {
		return nil, nil
	}
	cfg, err := cresettings.ParseShardAssignmentConfig(raw)
	if err != nil {
		if m.metrics.RecordSettingsParseError(ctx, settingsStringHash(raw)) {
			m.lggr.Warnw("Failed to parse shard assignment settings; keeping previous assignment live", "err", err)
		}
		m.cfgMu.Lock()
		defer m.cfgMu.Unlock()
		return m.lastCfg, nil
	}
	hash := shardAssignmentSettingsHash(cfg)
	m.cfgMu.Lock()
	defer m.cfgMu.Unlock()
	if m.metrics.RecordSettingsApplied(ctx, hash) {
		m.lggr.Infow("Shard assignment settings applied", assignmentDiffKVs(m.lastCfg, cfg, hash)...)
	}
	m.lastCfg = cfg
	return cfg, nil
}

func resolveManual(ctx context.Context, cfg *cresettings.ShardAssignmentConfig, ownerHex string, orgResolver OrgResolver) (uint32, bool, string, error) {
	owner := normalizeOwner(ownerHex)

	if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
		return shards[0], true, wfmon.ShardResolverSourcePerOwner, nil
	}

	if orgResolver != nil && len(cfg.PerOrgAssignment) > 0 {
		orgID, err := orgResolver.Get(ctx, ownerHex)
		if err == nil && orgID != "" {
			if shards, ok := cfg.PerOrgAssignment[orgID]; ok && len(shards) > 0 {
				return shards[0], true, wfmon.ShardResolverSourcePerOrg, nil
			}
		}
	}

	if cfg.HashedOwnerAssignment[owner] {
		return 0, false, wfmon.ShardResolverSourceHashed, nil
	}

	if cfg.HashedDefaultAssignment {
		return 0, false, wfmon.ShardResolverSourceHashed, nil
	}

	if len(cfg.StaticDefaultAssignment) > 0 {
		return cfg.StaticDefaultAssignment[0], true, wfmon.ShardResolverSourceStaticDefault, nil
	}

	return 0, false, wfmon.ShardResolverSourceStaticDefault, nil
}

func resolveAllManual(ctx context.Context, cfg *cresettings.ShardAssignmentConfig, ownerHex string, orgResolver OrgResolver) ([]uint32, bool, string, error) {
	owner := normalizeOwner(ownerHex)

	if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
		return shards, true, wfmon.ShardResolverSourcePerOwner, nil
	}

	if orgResolver != nil && len(cfg.PerOrgAssignment) > 0 {
		orgID, err := orgResolver.Get(ctx, ownerHex)
		if err == nil && orgID != "" {
			if shards, ok := cfg.PerOrgAssignment[orgID]; ok && len(shards) > 0 {
				return shards, true, wfmon.ShardResolverSourcePerOrg, nil
			}
		}
	}

	if cfg.HashedOwnerAssignment[owner] {
		return nil, false, wfmon.ShardResolverSourceHashed, nil
	}

	if cfg.HashedDefaultAssignment {
		return nil, false, wfmon.ShardResolverSourceHashed, nil
	}

	if len(cfg.StaticDefaultAssignment) > 0 {
		return cfg.StaticDefaultAssignment, true, wfmon.ShardResolverSourceStaticDefault, nil
	}

	return nil, false, wfmon.ShardResolverSourceStaticDefault, nil
}

func shardAssignmentSettingsHash(cfg *cresettings.ShardAssignmentConfig) uint32 {
	h := fnv.New32a()
	_, _ = fmt.Fprintf(h, "static=%v;disabled=%v;hashedDefault=%v", cfg.StaticDefaultAssignment, cfg.DisabledShards, cfg.HashedDefaultAssignment)
	for _, owner := range slices.Sorted(maps.Keys(cfg.PerOwnerAssignment)) {
		_, _ = fmt.Fprintf(h, ";owner=%s=%v", owner, cfg.PerOwnerAssignment[owner])
	}
	for _, org := range slices.Sorted(maps.Keys(cfg.PerOrgAssignment)) {
		_, _ = fmt.Fprintf(h, ";org=%s=%v", org, cfg.PerOrgAssignment[org])
	}
	for _, owner := range slices.Sorted(maps.Keys(cfg.HashedOwnerAssignment)) {
		_, _ = fmt.Fprintf(h, ";hashedOwner=%s", owner)
	}
	return h.Sum32()
}

func settingsStringHash(raw string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(raw))
	return h.Sum32()
}

// assignmentDiffKVs returns the log fields describing what actually changed
// between two settings versions: only the owners added, moved or removed plus
// counts, never the full owner map.
func assignmentDiffKVs(prev, curr *cresettings.ShardAssignmentConfig, hash uint32) []any {
	kvs := []any{
		"settingsHash", hash,
		"perOwnerCount", len(curr.PerOwnerAssignment),
		"perOrgCount", len(curr.PerOrgAssignment),
		"hashedOwnerCount", len(curr.HashedOwnerAssignment),
		"staticDefaultAssignment", curr.StaticDefaultAssignment,
		"hashedDefaultAssignment", curr.HashedDefaultAssignment,
	}
	if prev == nil {
		return append(kvs, "initial", true)
	}
	var added, moved, removed []string
	for owner, shards := range curr.PerOwnerAssignment {
		if old, ok := prev.PerOwnerAssignment[owner]; ok {
			if !slices.Equal(old, shards) {
				moved = append(moved, owner+"="+fmt.Sprint(old)+"->"+fmt.Sprint(shards))
			}
		} else {
			added = append(added, owner+"="+fmt.Sprint(shards))
		}
	}
	for owner := range prev.PerOwnerAssignment {
		if _, ok := curr.PerOwnerAssignment[owner]; !ok {
			removed = append(removed, owner)
		}
	}
	return append(kvs,
		"perOwnerAdded", added,
		"perOwnerMoved", moved,
		"perOwnerRemoved", removed,
		"perOrgAdded", mapKeyDiffCount(prev.PerOrgAssignment, curr.PerOrgAssignment),
		"perOrgRemoved", mapKeyDiffCount(curr.PerOrgAssignment, prev.PerOrgAssignment),
	)
}

func mapKeyDiffCount(from, to map[string][]uint32) int {
	n := 0
	for k := range to {
		if _, ok := from[k]; !ok {
			n++
		}
	}
	return n
}

type overrideShardResolver struct {
	manual      *manualShardResolver
	orgResolver OrgResolver
	ringOCR     ShardResolver
	lggr        logger.Logger
	metrics     *wfmon.ShardAssignmentMetrics
}

func NewOverrideShardResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, shardIndexMapper *ShardIndexMapper, ringOCR ShardResolver, lggr logger.Logger) ShardResolver {
	return &overrideShardResolver{
		manual:      &manualShardResolver{settings: settings, orgResolver: orgResolver, shardIndexMapper: shardIndexMapper, lggr: logger.Named(lggr, "ManualShardResolver"), metrics: globalShardAssignmentMetrics(lggr)},
		orgResolver: orgResolver,
		ringOCR:     ringOCR,
		lggr:        logger.Named(lggr, "OverrideShardResolver"),
		metrics:     globalShardAssignmentMetrics(lggr),
	}
}

func (o *overrideShardResolver) ResolveShard(ctx context.Context, workflowID string, ownerHex string) (uint32, bool, error) {
	start := time.Now()
	donID, found, err := o.resolveShard(ctx, workflowID, ownerHex)
	o.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverOverride, resolverOutcome(found, err), "", time.Since(start))
	return donID, found, err
}

func (o *overrideShardResolver) resolveShard(ctx context.Context, workflowID string, ownerHex string) (uint32, bool, error) {
	cfg, err := o.manual.loadObservedConfig(ctx)
	if err != nil {
		return 0, false, err
	}
	if cfg != nil {
		shardIndex, found, _, err := resolveManual(ctx, cfg, ownerHex, o.orgResolver)
		if err != nil {
			return 0, false, err
		}
		if found {
			return o.manual.toDonID(ctx, shardIndex), true, nil
		}
	}
	return o.ringOCR.ResolveShard(ctx, workflowID, ownerHex)
}

func (o *overrideShardResolver) ResolveAllShards(ctx context.Context, workflowID string, ownerHex string) ([]uint32, bool, error) {
	start := time.Now()
	donIDs, found, err := o.resolveAllShards(ctx, workflowID, ownerHex)
	o.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverOverride, resolverOutcome(found, err), "", time.Since(start))
	return donIDs, found, err
}

func (o *overrideShardResolver) resolveAllShards(ctx context.Context, workflowID string, ownerHex string) ([]uint32, bool, error) {
	cfg, cfgErr := o.manual.loadObservedConfig(ctx)
	if cfgErr != nil {
		return nil, false, cfgErr
	}
	if cfg != nil {
		shardIndices, found, _, resolveErr := resolveAllManual(ctx, cfg, ownerHex, o.orgResolver)
		if resolveErr != nil {
			return nil, false, resolveErr
		}
		if found {
			return o.manual.toDonIDs(ctx, shardIndices), true, nil
		}
	}
	shardID, found, err := o.ringOCR.ResolveShard(ctx, workflowID, ownerHex)
	if err != nil || !found {
		return nil, false, err
	}
	return []uint32{shardID}, true, nil
}

func (o *overrideShardResolver) ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	start := time.Now()
	result, err := o.resolveShards(ctx, workflowIDs, ownerHexes)
	if err != nil {
		o.metrics.RecordResolverRequest(ctx, wfmon.ShardResolverOverride, wfmon.ShardResolverOutcomeError, "", time.Since(start))
		return nil, err
	}
	o.metrics.RecordResolverDuration(ctx, wfmon.ShardResolverOverride, time.Since(start))
	return result, nil
}

func (o *overrideShardResolver) resolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	cfg, err := o.manual.loadObservedConfig(ctx)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return o.ringOCR.ResolveShards(ctx, workflowIDs, ownerHexes)
	}

	result := make(map[string]uint32, len(workflowIDs))
	var ringWorkflowIDs []string

	for i, wfID := range workflowIDs {
		if i >= len(ownerHexes) {
			break
		}
		shardIndex, found, _, err := resolveManual(ctx, cfg, ownerHexes[i], o.orgResolver)
		if err != nil {
			return nil, err
		}
		o.metrics.AddResolverRequest(ctx, wfmon.ShardResolverOverride, resolverOutcome(found, nil), "")
		if found {
			result[wfID] = o.manual.toDonID(ctx, shardIndex)
		} else {
			ringWorkflowIDs = append(ringWorkflowIDs, wfID)
		}
	}

	if len(ringWorkflowIDs) > 0 {
		ringResult, err := o.ringOCR.ResolveShards(ctx, ringWorkflowIDs, nil)
		if err != nil {
			return nil, err
		}
		maps.Copy(result, ringResult)
	}

	return result, nil
}

func (o *overrideShardResolver) GetRoutingResponse(ctx context.Context, workflowIDs []string) (*ringpb.GetWorkflowShardMappingResponse, error) {
	if rr, ok := o.ringOCR.(interface {
		GetRoutingResponse(ctx context.Context, workflowIDs []string) (*ringpb.GetWorkflowShardMappingResponse, error)
	}); ok {
		return rr.GetRoutingResponse(ctx, workflowIDs)
	}
	return nil, nil
}

func normalizeOwner(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	s = strings.ToLower(s)
	if _, err := hex.DecodeString(s); err != nil {
		return s
	}
	return s
}
