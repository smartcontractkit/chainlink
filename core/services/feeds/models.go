package feeds

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

const (
	JobTypeFluxMonitor        = "fluxmonitor"
	JobTypeOffchainReporting  = "ocr"
	JobTypeOffchainReporting2 = "ocr2"
)

type ChainType string

const (
	ChainTypeUnknown ChainType = "UNKNOWN"
	ChainTypeEVM     ChainType = "EVM"
)

func NewChainType(s string) (ChainType, error) {
	switch s {
	case "EVM":
		return ChainTypeEVM, nil
	default:
		return ChainTypeUnknown, errors.New("invalid chain type")
	}
}

// FeedsManager defines a registered Feeds Manager Service and the connection
// information.
type FeedsManager struct {
	ID                 int64
	Name               string
	URI                string
	PublicKey          crypto.PublicKey
	IsConnectionActive bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// ChainConfig defines the chain configuration for a Feeds Manager.
type ChainConfig struct {
	ID                int64
	FeedsManagerID    int64
	ChainID           string
	ChainType         ChainType
	AccountAddress    string
	AdminAddress      string
	FluxMonitorConfig FluxMonitorConfig
	OCR1Config        OCR1Config
	OCR2Config        OCR2Config
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// FluxMonitorConfig defines configuration for FluxMonitorJobs.
type FluxMonitorConfig struct {
	Enabled bool `json:"enabled"`
}

func (c FluxMonitorConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *FluxMonitorConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

// OCR1Config defines configuration for OCR1 Jobs.
type OCR1Config struct {
	Enabled     bool        `json:"enabled"`
	IsBootstrap bool        `json:"is_bootstrap"`
	Multiaddr   null.String `json:"multiaddr"`
	P2PPeerID   null.String `json:"p2p_peer_id"`
	KeyBundleID null.String `json:"key_bundle_id"`
}

func (c OCR1Config) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *OCR1Config) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

// OCR2Config defines configuration for OCR2 Jobs.
type OCR2Config struct {
	Enabled     bool        `json:"enabled"`
	IsBootstrap bool        `json:"is_bootstrap"`
	Multiaddr   null.String `json:"multiaddr"`
	P2PPeerID   null.String `json:"p2p_peer_id"`
	KeyBundleID null.String `json:"key_bundle_id"`
}

func (c OCR2Config) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *OCR2Config) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
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
	ID             int64
	RemoteUUID     uuid.UUID // RemoteUUID is the uuid of the proposal in FMS.
	Status         JobProposalStatus
	ExternalJobID  uuid.NullUUID // ExternalJobID is the external job id in the job spec.
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

// JobProposalSpec defines a versioned proposed spec for a JobProposal.
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
