package syncerlimiter

import (
	"sync"
)

const (
	defaultGlobal   = 200
	defaultPerOwner = 200
)

type Limits struct {
	global            *int32
	perOwner          map[string]*int32
	perOwnerOverrides map[string]int32
	config            Config
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

func NewWorkflowLimits(config Config) (*Limits, error) {
	cfg := Config{
		Global:            config.Global,
		PerOwner:          config.PerOwner,
		PerOwnerOverrides: config.PerOwnerOverrides,
	}

	if cfg.Global == 0 {
		cfg.Global = defaultGlobal
	}

	if cfg.PerOwner == 0 {
		cfg.PerOwner = defaultPerOwner
	}

	return &Limits{
		global:            new(int32),
		perOwner:          make(map[string]*int32),
		perOwnerOverrides: cfg.PerOwnerOverrides,
		config:            cfg,
	}, nil
}

func (l *Limits) Allow(owner string) (ownerAllow bool, globalAllow bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	countForOwner, ok := l.perOwner[owner]
	if !ok {
		l.perOwner[owner] = new(int32)
		countForOwner = l.perOwner[owner]
	}

	if *countForOwner < l.getPerOwnerLimit(owner) {
		ownerAllow = true
	}

	if *l.global < l.config.Global {
		globalAllow = true
	}

	if ownerAllow && globalAllow {
		*countForOwner++
		*l.global++
	}

	return ownerAllow, globalAllow
}

func (l *Limits) Decrement(owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	ownerLimiter, ok := l.perOwner[owner]
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
	limit, found := l.perOwnerOverrides[owner]
	if found {
		return limit
	}
	return l.config.PerOwner
}
