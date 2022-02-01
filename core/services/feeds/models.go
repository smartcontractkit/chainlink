package feeds

import (
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/utils/crypto"
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
	JobTypes  pq.StringArray

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

// JobProposalStatus are the status codes that define the stage of a proposal
type JobProposalStatus string

const (
	JobProposalStatusPending   JobProposalStatus = "pending"
	JobProposalStatusApproved  JobProposalStatus = "approved"
	JobProposalStatusRejected  JobProposalStatus = "rejected"
	JobProposalStatusCancelled JobProposalStatus = "cancelled"
)

// JobProposal represents a proposal which has been sent by a Feeds Manager.
//
// A job proposal has multiple spec versions which are created each time
// the Feeds Manager sends a new proposal version.
type JobProposal struct {
	ID int64
	// RemoteUUID is the unique id of the proposal in FMS.
	RemoteUUID uuid.UUID
	Status     JobProposalStatus
	// ExternalJobID is the external job id in the spec.
	ExternalJobID  uuid.NullUUID
	FeedsManagerID int64
	Multiaddrs     pq.StringArray
	PendingUpdate  bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SpecStatus is the status of each proposed spec.
type SpecStatus string

const (
	// SpecStatusPending defines a spec status  which has been proposed by the
	// FMS.
	SpecStatusPending SpecStatus = "pending"
	// SpecStatusApproved defines a spec status which the node op has approved.
	// An approved spec is currently being run by the node.
	SpecStatusApproved SpecStatus = "approved"
	// SpecStatusRejected defines a spec status which was proposed, but was
	// rejected by the node op.
	SpecStatusRejected SpecStatus = "rejected"
	// SpecStatusCancelled defines a spec status which was previously approved,
	// but cancelled by the node op. A cancelled spec is not being run by the
	// node.
	SpecStatusCancelled SpecStatus = "cancelled"
)

// JobProposalSpec defines a spec version for a job proposal
type JobProposalSpec struct {
	ID              int64
	Definition      string
	Status          SpecStatus
	Version         int32
	JobProposalID   int64
	StatusUpdatedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CanEditDefinition checks if the spec definition can be edited.
func (s *JobProposalSpec) CanEditDefinition() bool {
	return s.Status == SpecStatusPending ||
		s.Status == SpecStatusCancelled
}
