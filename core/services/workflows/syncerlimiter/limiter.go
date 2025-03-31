package syncerlimiter

import (
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	defaultGlobal   = 200
	defaultPerOwner = 200
)

type Limits struct {
	// global tracks the count of all running workflows
	global *int32

	// perOwner tracks the count of all running workflows for an address
	perOwner map[string]*int32

	// perOwnerOverrides maps an address to a running workflow limit
	perOwnerOverrides map[string]int32
	config            Config
	lggr              logger.Logger
	mu                sync.Mutex
}

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

func NewWorkflowLimits(lggr logger.Logger, config Config) (*Limits, error) {
	lggr = lggr.Named("WorkflowLimiter")
	cfg := Config{
		Global:            config.Global,
		PerOwner:          config.PerOwner,
		PerOwnerOverrides: normalizeOverrides(config.PerOwnerOverrides),
	}

	if cfg.Global == 0 {
		cfg.Global = defaultGlobal
	}

	if cfg.PerOwner == 0 {
		cfg.PerOwner = defaultPerOwner
	}

	lggr.Debugw("workflow limits set", "perOwner", cfg.PerOwner, "global", cfg.Global, "overrides", cfg.PerOwnerOverrides)

	return &Limits{
		global:            new(int32),
		perOwner:          make(map[string]*int32),
		perOwnerOverrides: cfg.PerOwnerOverrides,
		config:            cfg,
		lggr:              lggr,
	}, nil
}

func (l *Limits) Allow(owner string) (ownerAllow bool, globalAllow bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	owner = normalizeOwner(owner)
	countForOwner, ok := l.perOwner[owner]
	if !ok {
		l.perOwner[owner] = new(int32)
		countForOwner = l.perOwner[owner]
	}

	limit := l.getPerOwnerLimit(owner)

	l.lggr.Debugw("determined owner limit", "owner", owner, "limit", limit, "countForOwner", countForOwner)

	if *countForOwner < limit {
		ownerAllow = true
	}

	if *l.global < l.config.Global {
		globalAllow = true
	}

	if ownerAllow && globalAllow {
		*countForOwner++
		*l.global++
	}

	l.lggr.Debugw("assesed if owner is allowed", "owner", owner, "ownerAllow", ownerAllow, "globalAllow", globalAllow)

	return ownerAllow, globalAllow
}

func (l *Limits) Decrement(owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ownerLimiter, ok := l.perOwner[normalizeOwner(owner)]
	if !ok || *ownerLimiter <= 0 {
		return
	}

	*ownerLimiter--
	*l.global--
}

// getPerOwnerLimit returns the default limit per owner if there are no overrides found
// for the given owner.
func (l *Limits) getPerOwnerLimit(owner string) int32 {
	if l.perOwnerOverrides == nil {
		return l.config.PerOwner
	}
	limit, found := l.perOwnerOverrides[normalizeOwner(owner)]
	if found {
		l.lggr.Debugw("overriding limit for owner", "owner", owner, "limit", limit)
		return limit
	}

	l.lggr.Debugw("did not find owner in overrides, returning default", "owner", owner, "limit", l.config.PerOwner)
	return l.config.PerOwner
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
