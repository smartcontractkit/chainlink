package job

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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

//revive:disable:redefines-builtin-id
type Type string

func (t Type) String() string {
	return string(t)
}

func (t Type) RequiresPipelineSpec() bool {
	return requiresPipelineSpec[t]
}

func (t Type) SupportsAsync() bool {
	return supportsAsync[t]
}

func (t Type) SchemaVersion() uint32 {
	return schemaVersions[t]
}

var (
	requiresPipelineSpec = map[Type]bool{
		Cron:              true,
		DirectRequest:     true,
		FluxMonitor:       true,
		OffchainReporting: false, // bootstrap jobs do not require it
		Keeper:            true,
		VRF:               true,
		Webhook:           true,
	}
	supportsAsync = map[Type]bool{
		Cron:              true,
		DirectRequest:     true,
		FluxMonitor:       false,
		OffchainReporting: false,
		Keeper:            true,
		VRF:               true,
		Webhook:           true,
	}
	schemaVersions = map[Type]uint32{
		Cron:              1,
		DirectRequest:     1,
		FluxMonitor:       1,
		OffchainReporting: 1,
		Keeper:            2,
		VRF:               1,
		Webhook:           1,
	}
)

type Job struct {
	ID                            int32     `toml:"-" gorm:"primary_key"`
	ExternalJobID                 uuid.UUID `toml:"externalJobID"`
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
	WebhookSpecID                 *int32
	WebhookSpec                   *WebhookSpec
	PipelineSpecID                int32
	PipelineSpec                  *pipeline.Spec
	JobSpecErrors                 []SpecError `gorm:"foreignKey:JobID"`
	Type                          Type
	SchemaVersion                 uint32
	Name                          null.String
	MaxTaskDuration               models.Interval
	Pipeline                      pipeline.Pipeline `toml:"observationSource" gorm:"-"`
}

func ExternalJobIDEncodeStringToTopic(id uuid.UUID) common.Hash {
	return common.BytesToHash([]byte(strings.Replace(id.String(), "-", "", 4)))
}

func ExternalJobIDEncodeBytesToTopic(id uuid.UUID) common.Hash {
	return common.BytesToHash(common.RightPadBytes(id.Bytes(), utils.EVMWordByteLen))
}

// The external job ID (UUID) can be encoded into a log topic (32 bytes)
// by taking the string representation of the UUID, removing the dashes
// so that its 32 characters long and then encoding those characters to bytes.
func (j Job) ExternalIDEncodeStringToTopic() common.Hash {
	return ExternalJobIDEncodeStringToTopic(j.ExternalJobID)
}

// The external job ID (UUID) can also be encoded into a log topic (32 bytes)
// by taking the 16 bytes undelying the UUID and right padding it.
func (j Job) ExternalIDEncodeBytesToTopic() common.Hash {
	return ExternalJobIDEncodeBytesToTopic(j.ExternalJobID)
}

func (j Job) TableName() string {
	return "jobs"
}

// SetID takes the id as a string and attempts to convert it to an int32. If
// it succeeds, it will set it as the id on the job
func (j *Job) SetID(value string) error {
	id, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	j.ID = int32(id)
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
	return "job_spec_errors"
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
	ContractAddress                        ethkey.EIP55Address  `toml:"contractAddress"`
	P2PPeerID                              *p2pkey.PeerID       `toml:"p2pPeerID" gorm:"column:p2p_peer_id;default:null"`
	P2PBootstrapPeers                      pq.StringArray       `toml:"p2pBootstrapPeers" gorm:"column:p2p_bootstrap_peers;type:text[]"`
	IsBootstrapPeer                        bool                 `toml:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                *models.Sha256Hash   `toml:"keyBundleID" gorm:"type:bytea"`
	TransmitterAddress                     *ethkey.EIP55Address `toml:"transmitterAddress"`
	ObservationTimeout                     models.Interval      `toml:"observationTimeout" gorm:"type:bigint;default:null"`
	BlockchainTimeout                      models.Interval      `toml:"blockchainTimeout" gorm:"type:bigint;default:null"`
	ContractConfigTrackerSubscribeInterval models.Interval      `toml:"contractConfigTrackerSubscribeInterval" gorm:"default:null"`
	ContractConfigTrackerPollInterval      models.Interval      `toml:"contractConfigTrackerPollInterval" gorm:"type:bigint;default:null"`
	ContractConfigConfirmations            uint16               `toml:"contractConfigConfirmations"`
	EVMChainID                             *utils.Big           `toml:"evmChainID" gorm:"column:evm_chain_id"`
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

type ExternalInitiatorWebhookSpec struct {
	ExternalInitiatorID int64
	ExternalInitiator   bridges.ExternalInitiator `gorm:"foreignkey:ExternalInitiatorID;->"`
	WebhookSpecID       int32
	WebhookSpec         WebhookSpec `gorm:"foreignkey:WebhookSpecID;->"`
	Spec                models.JSON
}

type WebhookSpec struct {
	ID                            int32 `toml:"-" gorm:"primary_key"`
	ExternalInitiatorWebhookSpecs []ExternalInitiatorWebhookSpec
	CreatedAt                     time.Time `json:"createdAt" toml:"-"`
	UpdatedAt                     time.Time `json:"updatedAt" toml:"-"`
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

func (WebhookSpec) TableName() string {
	return "webhook_specs"
}

type DirectRequestSpec struct {
	ID                       int32                    `toml:"-" gorm:"primary_key"`
	ContractAddress          ethkey.EIP55Address      `toml:"contractAddress"`
	MinIncomingConfirmations clnull.Uint32            `toml:"minIncomingConfirmations"`
	Requesters               models.AddressCollection `toml:"requesters"`
	MinContractPayment       *assets.Link             `toml:"minContractPaymentLinkJuels"`
	EVMChainID               *utils.Big               `toml:"evmChainID" gorm:"column:evm_chain_id"`
	CreatedAt                time.Time                `toml:"-"`
	UpdatedAt                time.Time                `toml:"-"`
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

// Need to also try integer thresholds until
// https://github.com/pelletier/go-toml/issues/571 is addressed.
// The UI's TOML.stringify({"threshold": 1.0}) (https://github.com/iarna/iarna-toml)
// will return "threshold = 1" since ts/js doesn't know the
// difference between 1.0 and 1, so we need to address it on the backend.
type FluxMonitorSpecIntThreshold struct {
	ContractAddress     ethkey.EIP55Address `toml:"contractAddress"`
	Threshold           int                 `toml:"threshold"`
	AbsoluteThreshold   int                 `toml:"absoluteThreshold"`
	PollTimerPeriod     time.Duration
	PollTimerDisabled   bool
	IdleTimerPeriod     time.Duration
	IdleTimerDisabled   bool
	DrumbeatSchedule    string
	DrumbeatRandomDelay time.Duration
	DrumbeatEnabled     bool
	MinPayment          *assets.Link
}

type FluxMonitorSpec struct {
	ID              int32               `toml:"-" gorm:"primary_key"`
	ContractAddress ethkey.EIP55Address `toml:"contractAddress"`
	Threshold       float32             `toml:"threshold,float"`
	// AbsoluteThreshold is the maximum absolute change allowed in a fluxmonitored
	// value before a new round should be kicked off, so that the current value
	// can be reported on-chain.
	AbsoluteThreshold   float32 `toml:"absoluteThreshold,float" gorm:"type:float;not null"`
	PollTimerPeriod     time.Duration
	PollTimerDisabled   bool
	IdleTimerPeriod     time.Duration
	IdleTimerDisabled   bool
	DrumbeatSchedule    string
	DrumbeatRandomDelay time.Duration
	DrumbeatEnabled     bool
	MinPayment          *assets.Link
	EVMChainID          *utils.Big `toml:"evmChainID" gorm:"column:evm_chain_id"`
	CreatedAt           time.Time  `toml:"-"`
	UpdatedAt           time.Time  `toml:"-"`
}

type KeeperSpec struct {
	ID              int32               `toml:"-" gorm:"primary_key"`
	ContractAddress ethkey.EIP55Address `toml:"contractAddress"`
	FromAddress     ethkey.EIP55Address `toml:"fromAddress"`
	EVMChainID      *utils.Big          `toml:"evmChainID" gorm:"column:evm_chain_id"`
	CreatedAt       time.Time           `toml:"-"`
	UpdatedAt       time.Time           `toml:"-"`
}

type VRFSpec struct {
	ID                 int32
	CoordinatorAddress ethkey.EIP55Address  `toml:"coordinatorAddress"`
	PublicKey          secp256k1.PublicKey  `toml:"publicKey"`
	Confirmations      uint32               `toml:"confirmations"`
	EVMChainID         *utils.Big           `toml:"evmChainID" gorm:"column:evm_chain_id"`
	FromAddress        *ethkey.EIP55Address `toml:"fromAddress"`
	PollPeriod         *time.Duration       `toml:"pollPeriod"` // For v2 jobs
	CreatedAt          time.Time            `toml:"-"`
	UpdatedAt          time.Time            `toml:"-"`
}
