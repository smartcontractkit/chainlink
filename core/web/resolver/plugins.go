package resolver

import (
	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

type PluginsResolver struct {
	plugins feeds.Plugins
}

// Median returns the the status of the median plugin.
func (r PluginsResolver) Median() bool {
	return r.plugins.Median
}
