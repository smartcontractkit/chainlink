package presenters

import (
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

// JobProposalResource represents a job proposal JSONAPI resource.
type JobProposalResource struct {
	JAID
	Spec           string                  `json:"spec"`
	Status         feeds.JobProposalStatus `json:"status"`
	ExternalJobID  *string                 `json:"external_job_id"`
	FeedsManagerID string                  `json:"feeds_manager_id"`
	Multiaddrs     []string                `json:"multiaddrs"`
	CreatedAt      time.Time               `json:"createdAt"`
}

// GetName implements the api2go EntityNamer interface
func (r JobProposalResource) GetName() string {
	return "job_proposals"
}

// JobProposalResource constructs a new JobProposalResource.
func NewJobProposalResource(jp feeds.JobProposal) *JobProposalResource {
	res := &JobProposalResource{
		JAID:           NewJAIDInt64(jp.ID),
		Status:         jp.Status,
		Spec:           jp.Spec,
		FeedsManagerID: strconv.FormatInt(jp.FeedsManagerID, 10),
		Multiaddrs:     jp.Multiaddrs,
		CreatedAt:      jp.CreatedAt,
	}

	if jp.ExternalJobID.Valid {
		uuid := jp.ExternalJobID.UUID.String()
		res.ExternalJobID = &uuid
	}

	return res
}

// NewJobProposalResources initializes a slice of JSONAPI job proposal resources
func NewJobProposalResources(jps []feeds.JobProposal) []JobProposalResource {
	rs := []JobProposalResource{}

	for _, jp := range jps {
		rs = append(rs, *NewJobProposalResource(jp))
	}

	return rs
}
