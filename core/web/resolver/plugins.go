package resolver

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

type PluginsResolver struct {
	plugins feeds.Plugins
}

// Commit returns the the status of the commit plugin.
func (r PluginsResolver) Commit() bool {
	return r.plugins.Commit
}

// Execute returns the the status of the execute plugin.
func (r PluginsResolver) Execute() bool {
	return r.plugins.Execute
}

// Median returns the the status of the median plugin.
func (r PluginsResolver) Median() bool {
	return r.plugins.Median
}

// Mercury returns the the status of the mercury plugin.
func (r PluginsResolver) Mercury() bool {
	return r.plugins.Mercury
}
