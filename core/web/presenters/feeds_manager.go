package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

// FeedsManagerResource represents a Feeds Manager JSONAPI resource.
type FeedsManagerResource struct {
	JAID
	Name      string          `json:"name"`
	URI       string          `json:"uri"`
	PublicKey feeds.PublicKey `json:"publicKey"`
	JobTypes  []string        `json:"jobTypes"`
	Network   string          `json:"network"`
	CreatedAt time.Time       `json:"createdAt"`
}

// GetName implements the api2go EntityNamer interface
func (r FeedsManagerResource) GetName() string {
	return "feeds_managers"
}

// NewFeedsManagerResource constructs a new FeedsManagerResource.
func NewFeedsManagerResource(ms feeds.ManagerService) *FeedsManagerResource {
	return &FeedsManagerResource{
		JAID:      NewJAIDInt32(ms.ID),
		Name:      ms.Name,
		URI:       ms.URI,
		PublicKey: ms.PublicKey,
		JobTypes:  ms.JobTypes,
		Network:   ms.Network,
		CreatedAt: ms.CreatedAt,
	}
}

// NewJobResources initializes a slice of JSONAPI feed manager resources
func NewFeedsManagerResources(mss []feeds.ManagerService) []FeedsManagerResource {
	rs := []FeedsManagerResource{}

	for _, ms := range mss {
		rs = append(rs, *NewFeedsManagerResource(ms))
	}

	return rs
}
