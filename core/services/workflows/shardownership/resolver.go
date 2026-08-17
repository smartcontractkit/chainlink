package shardownership

import (
	"context"
	"encoding/hex"
	"fmt"
	"maps"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

type ShardResolver interface {
	ResolveShard(ctx context.Context, workflowID string, ownerHex string) (shardID uint32, found bool, err error)
	ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error)
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
	settings    *loop.AtomicSettings
	orgResolver OrgResolver
	lggr        logger.Logger
}

func NewManualShardResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, lggr logger.Logger) ShardResolver {
	return &manualShardResolver{settings: settings, orgResolver: orgResolver, lggr: logger.Named(lggr, "ManualShardResolver")}
}

func (m *manualShardResolver) ResolveShard(ctx context.Context, _ string, ownerHex string) (uint32, bool, error) {
	cfg, err := loadShardAssignmentConfig(m.settings)
	if err != nil || cfg == nil {
		return 0, false, err
	}
	return resolveManual(ctx, cfg, ownerHex, m.orgResolver)
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
		shardID, found, err := resolveManual(ctx, cfg, ownerHexes[i], m.orgResolver)
		if err != nil {
			return nil, err
		}
		if found {
			result[wfID] = shardID
		}
	}
	return result, nil
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

type overrideShardResolver struct {
	manual      *manualShardResolver
	orgResolver OrgResolver
	ringOCR     ShardResolver
	lggr        logger.Logger
}

func NewOverrideShardResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, ringOCR ShardResolver, lggr logger.Logger) ShardResolver {
	return &overrideShardResolver{
		manual:      &manualShardResolver{settings: settings, orgResolver: orgResolver, lggr: logger.Named(lggr, "ManualShardResolver")},
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
		shardID, found, err := resolveManual(ctx, cfg, ownerHex, o.orgResolver)
		if err != nil {
			return 0, false, err
		}
		if found {
			return shardID, true, nil
		}
	}
	return o.ringOCR.ResolveShard(ctx, workflowID, ownerHex)
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
		shardID, found, err := resolveManual(ctx, cfg, ownerHexes[i], o.orgResolver)
		if err != nil {
			return nil, err
		}
		if found {
			result[wfID] = shardID
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
