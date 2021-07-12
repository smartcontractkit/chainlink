package presenters

import (
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

// ExplorerStatus represents the connected server and status of the connection
// This is rendered as normal JSON (as opposed to a JSONAPI resource)
type ExplorerStatus struct {
	Status string `json:"status"`
	Url    string `json:"url"`
}

// NewExplorerStatus returns an initialized ExplorerStatus from the store
func NewExplorerStatus(statsPusher synchronization.StatsPusher) ExplorerStatus {
	url := statsPusher.GetURL()

	return ExplorerStatus{
		Status: string(statsPusher.GetStatus()),
		Url:    url.String(),
	}
}
