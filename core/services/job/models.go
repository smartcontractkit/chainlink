package job

import (
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/smartcontractkit/chainlink/core/store/models"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	null "gopkg.in/guregu/null.v4"
)

const (
	DirectRequest     Type = "directrequest"
	FluxMonitor       Type = "fluxmonitor"
	OffchainReporting Type = "offchainreporting"
)

type IDEmbed struct {
	ID int32 `json:"-" toml:"-"                 gorm:"primary_key"`
}

func (id IDEmbed) GetID() string {
	return fmt.Sprintf("%v", id.ID)
}

func (id *IDEmbed) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	id.ID = int32(ID)
	return nil
}

type SpecDB struct {
	IDEmbed
	OffchainreportingOracleSpecID *int32                       `json:"-"`
	OffchainreportingOracleSpec   *OffchainReportingOracleSpec `json:"offChainReportingOracleSpec" gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
	DirectRequestSpecID           *int32                       `json:"-"`
	DirectRequestSpec             *DirectRequestSpec           `json:"DirectRequestSpec" gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
	FluxMonitorSpecID             *int32                       `json:"-"`
	FluxMonitorSpec               *FluxMonitorSpec             `json:"fluxMonitorSpec" gorm:"save_association:true;association_autoupdate:true;association_autocreate:true"`
	PipelineSpecID                int32                        `json:"-"`
	PipelineSpec                  *PipelineSpec                `json:"pipelineSpec"`
	JobSpecErrors                 []SpecError                  `json:"errors" gorm:"foreignKey:JobID"`
	Type                          Type                         `json:"type"`
	SchemaVersion                 uint32                       `json:"schemaVersion"`
	Name                          null.String                  `json:"name"`
	MaxTaskDuration               models.Interval              `json:"maxTaskDuration"`
	Pipeline                      pipeline.TaskDAG             `json:"-" toml:"observationSource"`
}

func (SpecDB) TableName() string {
	return "jobs"
}

type SpecError struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	JobID       int32     `json:"-"`
	Description string    `json:"description"`
	Occurrences uint      `json:"occurrences"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (SpecError) TableName() string {
	return "job_spec_errors_v2"
}

type PipelineRun struct {
	ID int64 `json:"-" gorm:"primary_key"`
}

func (pr PipelineRun) GetID() string {
	return fmt.Sprintf("%v", pr.ID)
}

func (pr *PipelineRun) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	pr.ID = int64(ID)
	return nil
}

type PipelineSpec struct {
	IDEmbed
	DotDagSource string    `json:"dotDagSource"`
	CreatedAt    time.Time `json:"-"`
}

// TODO: remove pointers when upgrading to gormv2
// which has https://github.com/go-gorm/gorm/issues/2748 fixed.
type OffchainReportingOracleSpec struct {
	IDEmbed
	ContractAddress                        models.EIP55Address  `json:"contractAddress" toml:"contractAddress"`
	P2PPeerID                              *models.PeerID       `json:"p2pPeerID" toml:"p2pPeerID" gorm:"column:p2p_peer_id;default:null"`
	P2PBootstrapPeers                      pq.StringArray       `json:"p2pBootstrapPeers" toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:text[]"`
	IsBootstrapPeer                        bool                 `json:"isBootstrapPeer" toml:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                *models.Sha256Hash   `json:"keyBundleID" toml:"keyBundleID"                 gorm:"type:bytea"`
	TransmitterAddress                     *models.EIP55Address `json:"transmitterAddress" toml:"transmitterAddress"`
	ObservationTimeout                     models.Interval      `json:"observationTimeout" toml:"observationTimeout" gorm:"type:bigint;default:null"`
	BlockchainTimeout                      models.Interval      `json:"blockchainTimeout" toml:"blockchainTimeout" gorm:"type:bigint;default:null"`
	ContractConfigTrackerSubscribeInterval models.Interval      `json:"contractConfigTrackerSubscribeInterval" toml:"contractConfigTrackerSubscribeInterval" gorm:"default:null"`
	ContractConfigTrackerPollInterval      models.Interval      `json:"contractConfigTrackerPollInterval" toml:"contractConfigTrackerPollInterval" gorm:"type:bigint;default:null"`
	ContractConfigConfirmations            uint16               `json:"contractConfigConfirmations" toml:"contractConfigConfirmations" gorm:"default:null"`
	CreatedAt                              time.Time            `json:"createdAt" toml:"-"`
	UpdatedAt                              time.Time            `json:"updatedAt" toml:"-"`
}

func (s OffchainReportingOracleSpec) GetID() string {
	return fmt.Sprintf("%v", s.ID)
}

func (s *OffchainReportingOracleSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

func (s *OffchainReportingOracleSpec) BeforeCreate() error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *OffchainReportingOracleSpec) BeforeSave() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (OffchainReportingOracleSpec) TableName() string {
	return "offchainreporting_oracle_specs"
}

type DirectRequestSpec struct {
	IDEmbed
	ContractAddress models.EIP55Address `json:"contractAddress" toml:"contractAddress"`
	// OnChainJobSpecID is the sha256 of the TOML that created this job spec
	OnChainJobSpecID gethCommon.Hash
	CreatedAt        time.Time `json:"createdAt" toml:"-"`
	UpdatedAt        time.Time `json:"updatedAt" toml:"-"`
}

func (DirectRequestSpec) TableName() string {
	return "direct_request_specs"
}

type FluxMonitorSpec struct {
	IDEmbed
	ContractAddress models.EIP55Address `json:"contractAddress" toml:"contractAddress"`
	Precision       int32               `json:"precision,omitempty" gorm:"type:smallint"`
	Threshold       float32             `json:"threshold,omitempty"`
	// AbsoluteThreshold is the maximum absolute change allowed in a fluxmonitored
	// value before a new round should be kicked off, so that the current value
	// can be reported on-chain.
	AbsoluteThreshold float32       `json:"absoluteThreshold" gorm:"type:float;not null"`
	PollTimerPeriod   time.Duration `json:"pollTimerPeriod,omitempty" gorm:"type:jsonb"`
	PollTimerDisabled bool          `json:"pollTimerDisabled,omitempty" gorm:"type:jsonb"`
	IdleTimerPeriod   time.Duration `json:"idleTimerPeriod,omitempty" gorm:"type:jsonb"`
	IdleTimerDisabled bool          `json:"idleTimerDisabled,omitempty" gorm:"type:jsonb"`
	CreatedAt         time.Time     `json:"createdAt" toml:"-"`
	UpdatedAt         time.Time     `json:"updatedAt" toml:"-"`
}
