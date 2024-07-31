package resolver

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

type PluginsResolver struct {
	plugins feeds.Plugins
}

// Commit returns the status of the commit plugin.
func (r PluginsResolver) Commit() bool {
	return r.plugins.Commit
}

// Execute returns the status of the execute plugin.
func (r PluginsResolver) Execute() bool {
	return r.plugins.Execute
}

// Median returns the status of the median plugin.
func (r PluginsResolver) Median() bool {
	return r.plugins.Median
}

// Mercury returns the status of the mercury plugin.
func (r PluginsResolver) Mercury() bool {
	return r.plugins.Mercury
}

// LiquidityManager returns the the status of the liquidity manager plugin.
func (r PluginsResolver) Rebalancer() bool {
	return r.plugins.LiquidityManager
}
