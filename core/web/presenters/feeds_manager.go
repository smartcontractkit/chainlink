package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"gopkg.in/guregu/null.v4"
)

// FeedsManagerResource represents a Feeds Manager JSONAPI resource.
type FeedsManagerResource struct {
	JAID
	Name                   string           `json:"name"`
	URI                    string           `json:"uri"`
	PublicKey              crypto.PublicKey `json:"publicKey"`
	JobTypes               []string         `json:"jobTypes"`
	IsBootstrapPeer        bool             `json:"isBootstrapPeer"`
	BootstrapPeerMultiaddr null.String      `json:"bootstrapPeerMultiaddr"`
	IsConnectionActive     bool             `json:"isConnectionActive"`
	CreatedAt              time.Time        `json:"createdAt"`
}

// GetName implements the api2go EntityNamer interface
func (r FeedsManagerResource) GetName() string {
	return "feeds_managers"
}

// NewFeedsManagerResource constructs a new FeedsManagerResource.
func NewFeedsManagerResource(ms feeds.FeedsManager) *FeedsManagerResource {
	return &FeedsManagerResource{
		JAID:                   NewJAIDInt64(ms.ID),
		Name:                   ms.Name,
		URI:                    ms.URI,
		PublicKey:              ms.PublicKey,
		JobTypes:               ms.JobTypes,
		IsBootstrapPeer:        ms.IsOCRBootstrapPeer,
		BootstrapPeerMultiaddr: ms.OCRBootstrapPeerMultiaddr,
		IsConnectionActive:     ms.IsConnectionActive,
		CreatedAt:              ms.CreatedAt,
	}
}

// NewJobResources initializes a slice of JSONAPI feed manager resources
func NewFeedsManagerResources(mss []feeds.FeedsManager) []FeedsManagerResource {
	rs := []FeedsManagerResource{}

	for _, ms := range mss {
		rs = append(rs, *NewFeedsManagerResource(ms))
	}

	return rs
}
