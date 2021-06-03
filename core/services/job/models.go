package job

import (
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

const (
	Cron              Type = "cron"
	DirectRequest     Type = "directrequest"
	FluxMonitor       Type = "fluxmonitor"
	OffchainReporting Type = "offchainreporting"
	Keeper            Type = "keeper"
	VRF               Type = "vrf"
	Webhook           Type = "webhook"
)

type Job struct {
	ID                            int32 `toml:"-" gorm:"primary_key"`
	OffchainreportingOracleSpecID *int32
	OffchainreportingOracleSpec   *OffchainReportingOracleSpec
	CronSpecID                    *int32
	CronSpec                      *CronSpec
	DirectRequestSpecID           *int32
	DirectRequestSpec             *DirectRequestSpec
	FluxMonitorSpecID             *int32
	FluxMonitorSpec               *FluxMonitorSpec
	KeeperSpecID                  *int32
	KeeperSpec                    *KeeperSpec
	VRFSpecID                     *int32
	VRFSpec                       *VRFSpec
	WebhookSpecId                 *int32
	WebhookSpec                   *WebhookSpec
	PipelineSpecID                int32
	PipelineSpec                  *pipeline.Spec
	JobSpecErrors                 []SpecError `gorm:"foreignKey:JobID"`
	Type                          Type
	SchemaVersion                 uint32
	Name                          null.String
	MaxTaskDuration               models.Interval
	Pipeline                      pipeline.TaskDAG `toml:"observationSource" gorm:"-"`
}

func (Job) TableName() string {
	return "jobs"
}

// SetID takes the id as a string and attempts to convert it to an int32. If
// it succeeds, it will set it as the id on the job
func (job *Job) SetID(value string) error {
	id, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	job.ID = int32(id)
	return nil
}

type SpecError struct {
	ID          int64 `gorm:"primary_key"`
	JobID       int32
	Description string
	Occurrences uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

// TODO: remove pointers when upgrading to gormv2
// which has https://github.com/go-gorm/gorm/issues/2748 fixed.
type OffchainReportingOracleSpec struct {
	ID                                     int32                `toml:"-" gorm:"primary_key"`
	ContractAddress                        models.EIP55Address  `toml:"contractAddress"`
	P2PPeerID                              *models.PeerID       `toml:"p2pPeerID" gorm:"column:p2p_peer_id;default:null"`
	P2PBootstrapPeers                      pq.StringArray       `toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:text[]"`
	IsBootstrapPeer                        bool                 `toml:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                *models.Sha256Hash   `toml:"keyBundleID" gorm:"type:bytea"`
	TransmitterAddress                     *models.EIP55Address `toml:"transmitterAddress"`
	ObservationTimeout                     models.Interval      `toml:"observationTimeout" gorm:"type:bigint;default:null"`
	BlockchainTimeout                      models.Interval      `toml:"blockchainTimeout" gorm:"type:bigint;default:null"`
	ContractConfigTrackerSubscribeInterval models.Interval      `toml:"contractConfigTrackerSubscribeInterval" gorm:"default:null"`
	ContractConfigTrackerPollInterval      models.Interval      `toml:"contractConfigTrackerPollInterval" gorm:"type:bigint;default:null"`
	ContractConfigConfirmations            uint16               `toml:"contractConfigConfirmations"`
	CreatedAt                              time.Time            `toml:"-"`
	UpdatedAt                              time.Time            `toml:"-"`
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

func (s *OffchainReportingOracleSpec) BeforeCreate(db *gorm.DB) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *OffchainReportingOracleSpec) BeforeSave(db *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func (OffchainReportingOracleSpec) TableName() string {
	return "offchainreporting_oracle_specs"
}

type WebhookSpec struct {
	ID               int32        `toml:"-" gorm:"primary_key"`
	OnChainJobSpecID models.JobID `toml:"jobID"`
	CreatedAt        time.Time    `json:"createdAt" toml:"-"`
	UpdatedAt        time.Time    `json:"updatedAt" toml:"-"`
}

func (w WebhookSpec) GetID() string {
	return fmt.Sprintf("%v", w.ID)
}

func (w *WebhookSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	w.ID = int32(ID)
	return nil
}

func (w *WebhookSpec) BeforeCreate(db *gorm.DB) error {
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	return nil
}

func (w *WebhookSpec) BeforeSave(db *gorm.DB) error {
	w.UpdatedAt = time.Now()
	return nil
}

func (WebhookSpec) TableName() string {
	return "webhook_specs"
}

type DirectRequestSpec struct {
	ID                       int32               `toml:"-" gorm:"primary_key"`
	ContractAddress          models.EIP55Address `toml:"contractAddress"`
	OnChainJobSpecID         models.JobID        `toml:"jobID"`
	MinIncomingConfirmations clnull.Uint32       `toml:"minIncomingConfirmations"`
	CreatedAt                time.Time           `toml:"-"`
	UpdatedAt                time.Time           `toml:"-"`
}

func (DirectRequestSpec) TableName() string {
	return "direct_request_specs"
}

type CronSpec struct {
	ID           int32     `toml:"-" gorm:"primary_key"`
	CronSchedule string    `toml:"schedule"`
	CreatedAt    time.Time `toml:"-"`
	UpdatedAt    time.Time `toml:"-"`
}

func (s CronSpec) GetID() string {
	return fmt.Sprintf("%v", s.ID)
}

func (s *CronSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

func (s *CronSpec) BeforeCreate(db *gorm.DB) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *CronSpec) BeforeSave(db *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func (CronSpec) TableName() string {
	return "cron_specs"
}

type FluxMonitorSpec struct {
	ID              int32               `toml:"-" gorm:"primary_key"`
	ContractAddress models.EIP55Address `toml:"contractAddress"`
	Precision       int32               `gorm:"type:smallint"`
	Threshold       float32             `toml:"threshold,float"`
	// AbsoluteThreshold is the maximum absolute change allowed in a fluxmonitored
	// value before a new round should be kicked off, so that the current value
	// can be reported on-chain.
	AbsoluteThreshold float32       `toml:"absoluteThreshold,float" gorm:"type:float;not null"`
	PollTimerPeriod   time.Duration `gorm:"type:jsonb"`
	PollTimerDisabled bool          `gorm:"type:jsonb"`
	IdleTimerPeriod   time.Duration `gorm:"type:jsonb"`
	IdleTimerDisabled bool          `gorm:"type:jsonb"`
	MinPayment        *assets.Link
	CreatedAt         time.Time `toml:"-"`
	UpdatedAt         time.Time `toml:"-"`
}

type KeeperSpec struct {
	ID              int32               `toml:"-" gorm:"primary_key"`
	ContractAddress models.EIP55Address `toml:"contractAddress"`
	FromAddress     models.EIP55Address `toml:"fromAddress"`
	CreatedAt       time.Time           `toml:"-"`
	UpdatedAt       time.Time           `toml:"-"`
}

type VRFSpec struct {
	ID                 int32
	CoordinatorAddress models.EIP55Address `toml:"coordinatorAddress"`
	PublicKey          secp256k1.PublicKey `toml:"publicKey"`
	Confirmations      uint32              `toml:"confirmations"`
	CreatedAt          time.Time           `toml:"-"`
	UpdatedAt          time.Time           `toml:"-"`
}
