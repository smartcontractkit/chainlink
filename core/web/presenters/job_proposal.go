package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

// JobProposalResource represents a job proposal JSONAPI resource.
type JobProposalResource struct {
	JAID
	Spec      string                  `json:"spec"`
	Status    feeds.JobProposalStatus `json:"status"`
	CreatedAt time.Time               `json:"createdAt"`
}

// GetName implements the api2go EntityNamer interface
func (r JobProposalResource) GetName() string {
	return "job_proposals"
}

// JobProposalResource constructs a new JobProposalResource.
func NewJobProposalResource(jp feeds.JobProposal) *JobProposalResource {
	return &JobProposalResource{
		JAID:      NewJAIDInt64(jp.ID),
		Status:    jp.Status,
		Spec:      jp.Spec,
		CreatedAt: jp.CreatedAt,
	}
}

// NewJobProposalResources initializes a slice of JSONAPI job proposal resources
func NewJobProposalResources(jps []feeds.JobProposal) []JobProposalResource {
	rs := []JobProposalResource{}

	for _, jp := range jps {
		rs = append(rs, *NewJobProposalResource(jp))
	}

	return rs
}
