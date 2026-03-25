package chainlink

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	evmtoml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
)

const evmChainConfigHealthName = "EVMChainConfigHealth"

// EVMChainConfigHealth reports readiness failure when any enabled EVM chain was skipped during startup
// (partial load). Implements [services.Service] for registration with the application health checker.
type EVMChainConfigHealth struct {
	mu sync.RWMutex

	expected []string
	loaded   map[string]struct{}
	skipped  []evmSkippedRecord
}

type evmSkippedRecord struct {
	chainID string
	reason  string
	err     error
}

// NewEVMChainConfigHealth tracks whether all enabled EVM chains from config were loaded successfully.
func NewEVMChainConfigHealth(expectedEnabledChainIDs []string) *EVMChainConfigHealth {
	cp := append([]string(nil), expectedEnabledChainIDs...)
	return &EVMChainConfigHealth{
		expected: cp,
		loaded:   make(map[string]struct{}),
	}
}

// RecordLoaded marks a chain ID as successfully initialized (relayer created).
func (h *EVMChainConfigHealth) RecordLoaded(chainID string) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.loaded[chainID] = struct{}{}
}

// RecordSkipped records a chain that could not be loaded; readiness will fail until config/binary is fixed.
func (h *EVMChainConfigHealth) RecordSkipped(chainID, reason string, err error) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.skipped = append(h.skipped, evmSkippedRecord{chainID: chainID, reason: reason, err: err})
}

// ExpectedCount returns the number of enabled EVM chains in configuration.
func (h *EVMChainConfigHealth) ExpectedCount() int {
	if h == nil {
		return 0
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.expected)
}

func (h *EVMChainConfigHealth) Start(context.Context) error { return nil }

func (h *EVMChainConfigHealth) Close() error { return nil }

func (h *EVMChainConfigHealth) Name() string { return evmChainConfigHealthName }

func (h *EVMChainConfigHealth) Ready() error {
	if h == nil {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.skipped) == 0 {
		return nil
	}
	var b strings.Builder
	for i, s := range h.skipped {
		if i > 0 {
			b.WriteString("; ")
		}
		fmt.Fprintf(&b, "chain %s (%s)", s.chainID, s.reason)
		if s.err != nil {
			fmt.Fprintf(&b, ": %v", s.err)
		}
	}
	return fmt.Errorf("evm chain configuration degraded: %d of %d enabled chains failed to load: %s",
		len(h.skipped), len(h.expected), b.String())
}

func (h *EVMChainConfigHealth) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Ready()}
}

var _ services.Service = (*EVMChainConfigHealth)(nil)

// enabledEVMChainIDStrings returns chain ID strings for enabled EVM configs.
func enabledEVMChainIDStrings(cfgs evmtoml.EVMConfigs) []string {
	var out []string
	for _, c := range cfgs {
		if c == nil || !c.IsEnabled() || c.ChainID == nil {
			continue
		}
		out = append(out, c.ChainID.String())
	}
	return out
}
