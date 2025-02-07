package syncerlimiter

import (
	"sync"
)

const (
	defaultGlobalEngineCountLimit   = 50
	defaultPerOwnerEngineCountLimit = 5
)

type WorkflowSyncerLimiter struct {
	global   *int32
	perOwner map[string]*int32
	config   Config
	mu       sync.Mutex
}

type Config struct {
	Global   int32 `json:"global"`
	PerOwner int32 `json:"perOwner"`
}

func NewWorkflowSyncerLimiter(config Config) (*WorkflowSyncerLimiter, error) {
	if config.Global <= 0 || config.PerOwner <= 0 {
		config.Global = defaultGlobalEngineCountLimit
		config.PerOwner = defaultPerOwnerEngineCountLimit
	}

	return &WorkflowSyncerLimiter{
		global:   new(int32),
		perOwner: make(map[string]*int32),
		config:   config,
	}, nil
}

func (wsl *WorkflowSyncerLimiter) Allow(owner string) (ownerAllow bool, globalAllow bool) {
	wsl.mu.Lock()
	defer wsl.mu.Unlock()
	ownerLimiter, ok := wsl.perOwner[owner]
	if !ok {
		wsl.perOwner[owner] = new(int32)
		ownerLimiter = wsl.perOwner[owner]
	}

	if *ownerLimiter < wsl.config.PerOwner {
		ownerAllow = true
	}

	if *wsl.global < wsl.config.Global {
		globalAllow = true
	}

	if ownerAllow && globalAllow {
		*ownerLimiter++
		*wsl.global++
	}

	return ownerAllow, globalAllow
}

func (wsl *WorkflowSyncerLimiter) Decrement(owner string) {
	wsl.mu.Lock()
	defer wsl.mu.Unlock()
	ownerLimiter, ok := wsl.perOwner[owner]
	if !ok || *ownerLimiter <= 0 {
		return
	}

	*ownerLimiter--
	*wsl.global--
}
