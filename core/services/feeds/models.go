package feeds

import (
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"gopkg.in/guregu/null.v4"
)

// We only support OCR and FM for the feeds manager
const (
	JobTypeFluxMonitor       = "fluxmonitor"
	JobTypeOffchainReporting = "ocr"
)

// FeedsManager contains feeds manager related fields
type FeedsManager struct {
	ID        int64
	Name      string
	URI       string
	PublicKey crypto.PublicKey
	JobTypes  pq.StringArray `gorm:"type:text[]"`

	// Determines whether the node will be used as a bootstrap peer. If this is
	// true, you must have both an OCRBootstrapAddr and OCRBootstrapPeerID.
	IsOCRBootstrapPeer bool

	// The libp2p multiaddress which the node operator will assign to this node
	// for bootstrap peer discovery.
	OCRBootstrapPeerMultiaddr null.String

	// IsConnectionActive is the indicator of connection activeness
	IsConnectionActive bool

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
	ID int64
	// RemoteUUID is the unique id of the proposal in FMS.
	RemoteUUID uuid.UUID
	Spec       string
	Status     JobProposalStatus
	// ExternalJobID is the external job id in the spec.
	ExternalJobID  uuid.NullUUID
	FeedsManagerID int64
	Multiaddrs     pq.StringArray `gorm:"type:text[]"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
