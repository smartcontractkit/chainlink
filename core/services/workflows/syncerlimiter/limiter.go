package syncerlimiter

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

const (
	defaultGlobal   = 200
	defaultPerOwner = 200
)

type Config struct {
	// Global defines the maximum global number of workflows that can run on the node
	// across all owners.
	Global int32 `json:"global"`

	// PerOwner defines the maximum number of workflows that an owner may run.
	PerOwner int32 `json:"perOwner"`

	// PerOwnerOverrides is a map of owner address to a workflow limit.  If the map does
	// not exist, or an address is not found, then the PerOwner limit is used.
	PerOwnerOverrides map[string]int32 `json:"overrides"`
}

type keyedOwnerSettings struct {
	key  string
	vals map[string]string
}

func (k keyedOwnerSettings) GetScoped(ctx context.Context, scope settings.Scope, key string) (value string, err error) {
	if k.key != key || scope != settings.ScopeOwner {
		return "", nil
	}
	return k.vals[contexts.CREValue(ctx).Owner], nil
}

func NewWorkflowLimits(lggr logger.Logger, cfg Config, lf limits.Factory) (limits.ResourceLimiter[int32], error) {
	lggr = logger.Named(lggr, "WorkflowLimiter")
	if cfg.Global == 0 {
		cfg.Global = defaultGlobal
	}
	if cfg.PerOwner == 0 {
		cfg.PerOwner = defaultPerOwner
	}
	cfg.PerOwnerOverrides = normalizeOverrides(cfg.PerOwnerOverrides)

	// TODO cresettings
	ownerLimit := settings.Int32(cfg.PerOwner)
	ownerLimit.Key = "PerOwner.WorkflowLimit"
	ownerLimit.Scope = settings.ScopeOwner
	perOwner := make(map[string]string, len(cfg.PerOwnerOverrides))
	for k, v := range cfg.PerOwnerOverrides {
		perOwner[k] = strconv.Itoa(int(v))
	}
	lf.Settings = keyedOwnerSettings{key: ownerLimit.Key, vals: perOwner}
	owner, err := limits.NewResourcePoolLimiter(lf, ownerLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create owner resource limiter: %w", err)
	}
	globalLimit := settings.Int32(cfg.Global)
	globalLimit.Key = "WorkflowLimit"
	globalLimit.Scope = settings.ScopeGlobal
	global, err := limits.NewResourcePoolLimiter(lf, globalLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create global resource limiter: %w", err)
	}

	lggr.Debugw("workflow limits set", "perOwner", cfg.PerOwner, "global", cfg.Global, "overrides", cfg.PerOwnerOverrides)

	return limits.MultiResourcePoolLimiter[int32]{owner, global}, nil
}

// normalizeOverrides ensures all incoming keys are normalized
func normalizeOverrides(in map[string]int32) map[string]int32 {
	out := make(map[string]int32)
	for k, v := range in {
		norm := normalizeOwner(k)
		out[norm] = v
	}
	return out
}

// normalizeOwner removes any 0x prefix
func normalizeOwner(k string) string {
	norm := k
	if strings.HasPrefix(k, "0x") {
		norm = norm[2:]
	}
	return strings.ToLower(norm)
}
