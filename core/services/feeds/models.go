package feeds

import (
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

// We only support OCR and FM for the feeds manager
const (
	JobTypeFluxMonitor       = "fluxmonitor"
	JobTypeOffchainReporting = "offchainreporting"
)

type FeedsManager struct {
	ID        int64
	Name      string
	URI       string
	PublicKey crypto.PublicKey
	JobTypes  pq.StringArray `gorm:"type:text[]"`
	Network   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (FeedsManager) TableName() string {
	return "feeds_managers"
}

// JobProposalStatus are the status codes that define the stage of a proposal
type JobProposalStatus string

const (
	JobProposalStatusPending  JobProposalStatus = "pending"
	JobProposalStatusApproved JobProposalStatus = "approved"
	JobProposalStatusRejected JobProposalStatus = "rejected"
)

type JobProposal struct {
	ID             int64
	Spec           string
	Status         JobProposalStatus
	JobID          uuid.NullUUID
	FeedsManagerID int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
