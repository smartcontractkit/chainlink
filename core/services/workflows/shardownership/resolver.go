package shardownership

import (
	"context"
	"encoding/hex"
	"fmt"
	"maps"
	"strings"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"

	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

type ShardResolver interface {
	ResolveShard(ctx context.Context, workflowID string, ownerHex string) (shardID uint32, found bool, err error)
	ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error)
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
	config *atomic.Pointer[cresettings.ShardAssignmentConfig]
	lggr   logger.Logger
}

func NewManualShardResolver(config *atomic.Pointer[cresettings.ShardAssignmentConfig], lggr logger.Logger) ShardResolver {
	return &manualShardResolver{config: config, lggr: logger.Named(lggr, "ManualShardResolver")}
}

func (m *manualShardResolver) ResolveShard(_ context.Context, _ string, ownerHex string) (uint32, bool, error) {
	cfg := m.config.Load()
	if cfg == nil {
		return 0, false, nil
	}
	return resolveManual(cfg, ownerHex)
}

func (m *manualShardResolver) ResolveShards(_ context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	cfg := m.config.Load()
	if cfg == nil {
		return map[string]uint32{}, nil
	}
	result := make(map[string]uint32, len(workflowIDs))
	for i, wfID := range workflowIDs {
		if i >= len(ownerHexes) {
			break
		}
		shardID, found, err := resolveManual(cfg, ownerHexes[i])
		if err != nil {
			return nil, err
		}
		if found {
			result[wfID] = shardID
		}
	}
	return result, nil
}

func resolveManual(cfg *cresettings.ShardAssignmentConfig, ownerHex string) (uint32, bool, error) {
	owner := normalizeOwner(ownerHex)

	if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
		return shards[0], true, nil
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
	manual  *manualShardResolver
	ringOCR ShardResolver
	lggr    logger.Logger
}

func NewOverrideShardResolver(config *atomic.Pointer[cresettings.ShardAssignmentConfig], ringOCR ShardResolver, lggr logger.Logger) ShardResolver {
	return &overrideShardResolver{
		manual:  &manualShardResolver{config: config, lggr: logger.Named(lggr, "ManualShardResolver")},
		ringOCR: ringOCR,
		lggr:    logger.Named(lggr, "OverrideShardResolver"),
	}
}

func (o *overrideShardResolver) ResolveShard(ctx context.Context, workflowID string, ownerHex string) (uint32, bool, error) {
	cfg := o.manual.config.Load()
	if cfg != nil {
		owner := normalizeOwner(ownerHex)
		if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
			return shards[0], true, nil
		}
		if cfg.HashedOwnerAssignment[owner] {
			return o.ringOCR.ResolveShard(ctx, workflowID, ownerHex)
		}
		if cfg.HashedDefaultAssignment {
			return o.ringOCR.ResolveShard(ctx, workflowID, ownerHex)
		}
		if len(cfg.StaticDefaultAssignment) > 0 {
			return cfg.StaticDefaultAssignment[0], true, nil
		}
	}
	return o.ringOCR.ResolveShard(ctx, workflowID, ownerHex)
}

func (o *overrideShardResolver) ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	cfg := o.manual.config.Load()
	if cfg == nil {
		return o.ringOCR.ResolveShards(ctx, workflowIDs, ownerHexes)
	}

	result := make(map[string]uint32, len(workflowIDs))
	var ringWorkflowIDs []string

	for i, wfID := range workflowIDs {
		if i >= len(ownerHexes) {
			break
		}
		owner := normalizeOwner(ownerHexes[i])
		if shards, ok := cfg.PerOwnerAssignment[owner]; ok && len(shards) > 0 {
			result[wfID] = shards[0]
			continue
		}
		if cfg.HashedOwnerAssignment[owner] || cfg.HashedDefaultAssignment {
			ringWorkflowIDs = append(ringWorkflowIDs, wfID)
			continue
		}
		if len(cfg.StaticDefaultAssignment) > 0 {
			result[wfID] = cfg.StaticDefaultAssignment[0]
			continue
		}
		ringWorkflowIDs = append(ringWorkflowIDs, wfID)
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
