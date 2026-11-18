package v2

import (
	"context"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
)

// minRegistryPageSize is the smallest page size adaptive paging will shrink to.
const minRegistryPageSize = int64(10)

// isGasError reports whether err looks like an eth_call failing because it needs more gas than the
// RPC node allows (e.g. "out of gas: gas required exceeds 16777216", geth's "gas required exceeds
// allowance"). Retrying the read with a smaller page can fix these. The match is intentionally
// broad: any error mentioning "gas" is treated as a gas-limit failure, since RPC providers word
// them differently. A non-gas error that happens to mention gas only causes a few smaller retries
// before the call fails as it would have anyway.
func isGasError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "gas")
}

// registryPageSizer decides the page size for paged WorkflowRegistry eth_calls.
//
// The page size comes from the cresettings limit WorkflowRegistryMaxResultsPerQuery, re-read on
// every call so operators can change it without a restart. On top of that it supports a runtime
// reduction: when the RPC node rejects a page for exceeding its gas cap, shrink halves the size and
// the lower value is kept so later ticks do not repeat the failing size. A reduction is dropped
// whenever the configured value changes, and the effective size never exceeds the configured one.
type registryPageSizer struct {
	mu         sync.Mutex
	getter     settings.Getter // nil falls back to the cresettings default
	configured int64           // configured value that `reduced` was derived from
	reduced    int64           // 0 means no reduction
}

func (p *registryPageSizer) setGetter(g settings.Getter) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.getter = g
}

// current returns the page size to use for the next call.
func (p *registryPageSizer) current(ctx context.Context, lggr logger.Logger) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	setting := cresettings.Default.WorkflowRegistryMaxResultsPerQuery
	v, err := setting.GetOrDefault(ctx, p.getter)
	if err != nil {
		lggr.Warnw("Failed to read WorkflowRegistryMaxResultsPerQuery, using default", "err", err)
	}
	if v < 1 {
		v = setting.DefaultValue
	}
	configured := int64(v)

	if configured != p.configured {
		p.configured, p.reduced = configured, 0
	}
	if p.reduced > 0 && p.reduced < configured {
		return p.reduced
	}
	return configured
}

// shrink halves failed (the page size that just hit the gas cap) and remembers the result.
// It returns false if failed is already at the minimum, i.e. there is nothing left to try.
func (p *registryPageSizer) shrink(failed int64) (int64, bool) {
	if failed <= minRegistryPageSize {
		return failed, false
	}
	next := max(failed/2, minRegistryPageSize)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.reduced = next
	return next, true
}
