package shardownership

import (
	"context"
	"encoding/hex"
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
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

type ringOCRShardResolver struct {
	client shardorchestrator.ClientInterface
	lggr   logger.Logger
}

func NewRingOCRShardResolver(client shardorchestrator.ClientInterface, lggr logger.Logger) ShardResolver {
	return &ringOCRShardResolver{client: client, lggr: logger.Named(lggr, "RingOCRShardResolver")}
}

func (r *ringOCRShardResolver) ResolveShard(ctx context.Context, workflowID string, _ string) (uint32, bool, error) {
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
}

func NewManualShardResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, shardIndexMapper *ShardIndexMapper, lggr logger.Logger) ShardResolver {
	return &manualShardResolver{settings: settings, orgResolver: orgResolver, shardIndexMapper: shardIndexMapper, lggr: logger.Named(lggr, "ManualShardResolver")}
}

func (m *manualShardResolver) ResolveShard(ctx context.Context, _ string, ownerHex string) (uint32, bool, error) {
	cfg, err := loadShardAssignmentConfig(m.settings)
	if err != nil || cfg == nil {
		return 0, false, err
	}
	shardIndex, found, err := resolveManual(ctx, cfg, ownerHex, m.orgResolver)
	if err != nil || !found {
		return 0, found, err
	}
	return m.toDonID(ctx, shardIndex), true, nil
}

func (m *manualShardResolver) ResolveAllShards(ctx context.Context, _ string, ownerHex string) ([]uint32, bool, error) {
	cfg, err := loadShardAssignmentConfig(m.settings)
	if err != nil || cfg == nil {
		return nil, false, err
	}
	shardIndices, found, err := resolveAllManual(ctx, cfg, ownerHex, m.orgResolver)
	if err != nil || !found {
		return nil, found, err
	}
	return m.toDonIDs(ctx, shardIndices), true, nil
}

func (m *manualShardResolver) ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	cfg, err := loadShardAssignmentConfig(m.settings)
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
		shardIndex, found, err := resolveManual(ctx, cfg, ownerHexes[i], m.orgResolver)
		if err != nil {
			return nil, err
		}
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

func loadShardAssignmentConfig(settings *loop.AtomicSettings) (*cresettings.ShardAssignmentConfig, error) {
	if settings == nil {
		return nil, nil
	}
	update, err := settings.Load()
	if err != nil || update.Settings == "" {
		return nil, err
	}
	return cresettings.ParseShardAssignmentConfig(update.Settings)
}

func resolveManual(ctx context.Context, cfg *cresettings.ShardAssignmentConfig, ownerHex string, orgResolver OrgResolver) (uint32, bool, error) {
	owner := normalizeOwner(ownerHex)

	if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
		return shards[0], true, nil
	}

	if orgResolver != nil && len(cfg.PerOrgAssignment) > 0 {
		orgID, err := orgResolver.Get(ctx, ownerHex)
		if err == nil && orgID != "" {
			if shards, ok := cfg.PerOrgAssignment[orgID]; ok && len(shards) > 0 {
				return shards[0], true, nil
			}
		}
	}

	if cfg.HashedOwnerAssignment[owner] {
		return 0, false, nil
	}

	if cfg.HashedDefaultAssignment {
		return 0, false, nil
	}

	if len(cfg.StaticDefaultAssignment) > 0 {
		return cfg.StaticDefaultAssignment[0], true, nil
	}

	return 0, false, nil
}

func resolveAllManual(ctx context.Context, cfg *cresettings.ShardAssignmentConfig, ownerHex string, orgResolver OrgResolver) ([]uint32, bool, error) {
	owner := normalizeOwner(ownerHex)

	if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
		return shards, true, nil
	}

	if orgResolver != nil && len(cfg.PerOrgAssignment) > 0 {
		orgID, err := orgResolver.Get(ctx, ownerHex)
		if err == nil && orgID != "" {
			if shards, ok := cfg.PerOrgAssignment[orgID]; ok && len(shards) > 0 {
				return shards, true, nil
			}
		}
	}

	if cfg.HashedOwnerAssignment[owner] {
		return nil, false, nil
	}

	if cfg.HashedDefaultAssignment {
		return nil, false, nil
	}

	if len(cfg.StaticDefaultAssignment) > 0 {
		return cfg.StaticDefaultAssignment, true, nil
	}

	return nil, false, nil
}

type overrideShardResolver struct {
	manual      *manualShardResolver
	orgResolver OrgResolver
	ringOCR     ShardResolver
	lggr        logger.Logger
}

func NewOverrideShardResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, shardIndexMapper *ShardIndexMapper, ringOCR ShardResolver, lggr logger.Logger) ShardResolver {
	return &overrideShardResolver{
		manual:      &manualShardResolver{settings: settings, orgResolver: orgResolver, shardIndexMapper: shardIndexMapper, lggr: logger.Named(lggr, "ManualShardResolver")},
		orgResolver: orgResolver,
		ringOCR:     ringOCR,
		lggr:        logger.Named(lggr, "OverrideShardResolver"),
	}
}

func (o *overrideShardResolver) ResolveShard(ctx context.Context, workflowID string, ownerHex string) (uint32, bool, error) {
	cfg, err := loadShardAssignmentConfig(o.manual.settings)
	if err != nil {
		return 0, false, err
	}
	if cfg != nil {
		shardIndex, found, err := resolveManual(ctx, cfg, ownerHex, o.orgResolver)
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
	cfg, cfgErr := loadShardAssignmentConfig(o.manual.settings)
	if cfgErr != nil {
		return nil, false, cfgErr
	}
	if cfg != nil {
		shardIndices, found, resolveErr := resolveAllManual(ctx, cfg, ownerHex, o.orgResolver)
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
	cfg, err := loadShardAssignmentConfig(o.manual.settings)
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
		shardIndex, found, err := resolveManual(ctx, cfg, ownerHexes[i], o.orgResolver)
		if err != nil {
			return nil, err
		}
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
