package syncerlimiter

import (
	"sync"
)

const (
	defaultGlobal   = 200
	defaultPerOwner = 200
)

type Limits struct {
	global   *int32
	perOwner map[string]*int32
	config   Config
	mu       sync.Mutex
}

type Config struct {
	Global   int32 `json:"global"`
	PerOwner int32 `json:"perOwner"`
}

func NewWorkflowLimits(config Config) (*Limits, error) {
	if config.Global <= 0 || config.PerOwner <= 0 {
		config.Global = defaultGlobal
		config.PerOwner = defaultPerOwner
	}

	return &Limits{
		global:   new(int32),
		perOwner: make(map[string]*int32),
		config:   config,
	}, nil
}

func (l *Limits) Allow(owner string) (ownerAllow bool, globalAllow bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	ownerLimiter, ok := l.perOwner[owner]
	if !ok {
		l.perOwner[owner] = new(int32)
		ownerLimiter = l.perOwner[owner]
	}

	if *ownerLimiter < l.config.PerOwner {
		ownerAllow = true
	}

	if *l.global < l.config.Global {
		globalAllow = true
	}

	if ownerAllow && globalAllow {
		*ownerLimiter++
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
