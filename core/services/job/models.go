package job

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

const (
	Cron               Type = "cron"
	DirectRequest      Type = "directrequest"
	FluxMonitor        Type = "fluxmonitor"
	OffchainReporting  Type = "offchainreporting"
	OffchainReporting2 Type = "offchainreporting2"
	Keeper             Type = "keeper"
	VRF                Type = "vrf"
	Webhook            Type = "webhook"
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
		Cron:               true,
		DirectRequest:      true,
		FluxMonitor:        true,
		OffchainReporting:  false, // bootstrap jobs do not require it
		OffchainReporting2: false, // bootstrap jobs do not require it
		Keeper:             true,
		VRF:                true,
		Webhook:            true,
	}
	supportsAsync = map[Type]bool{
		Cron:               true,
		DirectRequest:      true,
		FluxMonitor:        false,
		OffchainReporting:  false,
		OffchainReporting2: false,
		Keeper:             true,
		VRF:                true,
		Webhook:            true,
	}
	schemaVersions = map[Type]uint32{
		Cron:               1,
		DirectRequest:      1,
		FluxMonitor:        1,
		OffchainReporting:  1,
		OffchainReporting2: 1,
		Keeper:             3,
		VRF:                1,
		Webhook:            1,
	}
)

type Job struct {
	ID                             int32     `toml:"-"`
	ExternalJobID                  uuid.UUID `toml:"externalJobID"`
	OffchainreportingOracleSpecID  *int32
	OffchainreportingOracleSpec    *OffchainReportingOracleSpec
	Offchainreporting2OracleSpecID *int32
	Offchainreporting2OracleSpec   *OffchainReporting2OracleSpec
	CronSpecID                     *int32
	CronSpec                       *CronSpec
	DirectRequestSpecID            *int32
	DirectRequestSpec              *DirectRequestSpec
	FluxMonitorSpecID              *int32
	FluxMonitorSpec                *FluxMonitorSpec
	KeeperSpecID                   *int32
	KeeperSpec                     *KeeperSpec
	VRFSpecID                      *int32
	VRFSpec                        *VRFSpec
	WebhookSpecID                  *int32
	WebhookSpec                    *WebhookSpec
	PipelineSpecID                 int32
	PipelineSpec                   *pipeline.Spec
	JobSpecErrors                  []SpecError
	Type                           Type
	SchemaVersion                  uint32
	Name                           null.String
	MaxTaskDuration                models.Interval
	Pipeline                       pipeline.Pipeline `toml:"observationSource"`
	CreatedAt                      time.Time
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
	ID          int64
	JobID       int32
	Description string
	Occurrences uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SpecError) TableName() string {
	return "job_spec_errors"
}

// SetID takes the id as a string and attempts to convert it to an int32. If
// it succeeds, it will set it as the id on the job
func (j *SpecError) SetID(value string) error {
	id, err := stringutils.ToInt64(value)
	if err != nil {
		return err
	}
	j.ID = id
	return nil
}

type PipelineRun struct {
	ID int64 `json:"-"`
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

type OffchainReportingOracleSpec struct {
	ID                                        int32               `toml:"-"`
	ContractAddress                           ethkey.EIP55Address `toml:"contractAddress"`
	P2PBootstrapPeers                         pq.StringArray      `toml:"p2pBootstrapPeers" db:"p2p_bootstrap_peers"`
	IsBootstrapPeer                           bool                `toml:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                   *models.Sha256Hash  `toml:"keyBundleID"`
	EncryptedOCRKeyBundleIDEnv                bool
	TransmitterAddress                        *ethkey.EIP55Address `toml:"transmitterAddress"`
	TransmitterAddressEnv                     bool
	ObservationTimeout                        models.Interval `toml:"observationTimeout"`
	ObservationTimeoutEnv                     bool
	BlockchainTimeout                         models.Interval `toml:"blockchainTimeout"`
	BlockchainTimeoutEnv                      bool
	ContractConfigTrackerSubscribeInterval    models.Interval `toml:"contractConfigTrackerSubscribeInterval"`
	ContractConfigTrackerSubscribeIntervalEnv bool
	ContractConfigTrackerPollInterval         models.Interval `toml:"contractConfigTrackerPollInterval"`
	ContractConfigTrackerPollIntervalEnv      bool
	ContractConfigConfirmations               uint16 `toml:"contractConfigConfirmations"`
	ContractConfigConfirmationsEnv            bool
	EVMChainID                                *utils.Big       `toml:"evmChainID" db:"evm_chain_id"`
	DatabaseTimeout                           *models.Interval `toml:"databaseTimeout"`
	DatabaseTimeoutEnv                        bool
	ObservationGracePeriod                    *models.Interval `toml:"observationGracePeriod"`
	ObservationGracePeriodEnv                 bool
	ContractTransmitterTransmitTimeout        *models.Interval `toml:"contractTransmitterTransmitTimeout"`
	ContractTransmitterTransmitTimeoutEnv     bool
	CreatedAt                                 time.Time `toml:"-"`
	UpdatedAt                                 time.Time `toml:"-"`
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

func (OffchainReportingOracleSpec) TableName() string {
	return "offchainreporting_oracle_specs"
}

type RelayConfig map[string]interface{}

func (r RelayConfig) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

func (r RelayConfig) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *RelayConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("expected bytes got %T", b)
	}
	return json.Unmarshal(b, &r)
}

// Relay config is chain specific config for a relay (chain adapter).
type OffchainReporting2OracleSpec struct {
	ID                                     int32              `toml:"-"`
	ContractID                             string             `toml:"contractID"`
	Relay                                  relaytypes.Network `toml:"relay"`
	RelayConfig                            RelayConfig        `toml:"relayConfig"`
	P2PBootstrapPeers                      pq.StringArray     `toml:"p2pBootstrapPeers"`
	IsBootstrapPeer                        bool               `toml:"isBootstrapPeer"`
	OCRKeyBundleID                         null.String        `toml:"ocrKeyBundleID"`
	MonitoringEndpoint                     null.String        `toml:"monitoringEndpoint"`
	TransmitterID                          null.String        `toml:"transmitterID"`
	BlockchainTimeout                      models.Interval    `toml:"blockchainTimeout"`
	ContractConfigTrackerSubscribeInterval models.Interval    `toml:"contractConfigTrackerSubscribeInterval"`
	ContractConfigTrackerPollInterval      models.Interval    `toml:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations            uint16             `toml:"contractConfigConfirmations"`
	JuelsPerFeeCoinPipeline                string             `toml:"juelsPerFeeCoinSource"`
	CreatedAt                              time.Time          `toml:"-"`
	UpdatedAt                              time.Time          `toml:"-"`
}

func (s OffchainReporting2OracleSpec) GetID() string {
	return fmt.Sprintf("%v", s.ID)
}

func (s *OffchainReporting2OracleSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

func (OffchainReporting2OracleSpec) TableName() string {
	return "offchainreporting2_oracle_specs"
}

type ExternalInitiatorWebhookSpec struct {
	ExternalInitiatorID int64
	ExternalInitiator   bridges.ExternalInitiator
	WebhookSpecID       int32
	WebhookSpec         WebhookSpec
	Spec                models.JSON
}

type WebhookSpec struct {
	ID                            int32 `toml:"-"`
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
	ID                          int32                    `toml:"-"`
	ContractAddress             ethkey.EIP55Address      `toml:"contractAddress"`
	MinIncomingConfirmations    clnull.Uint32            `toml:"minIncomingConfirmations"`
	MinIncomingConfirmationsEnv bool                     `toml:"minIncomingConfirmationsEnv"`
	Requesters                  models.AddressCollection `toml:"requesters"`
	MinContractPayment          *assets.Link             `toml:"minContractPaymentLinkJuels"`
	EVMChainID                  *utils.Big               `toml:"evmChainID"`
	CreatedAt                   time.Time                `toml:"-"`
	UpdatedAt                   time.Time                `toml:"-"`
}

func (DirectRequestSpec) TableName() string {
	return "direct_request_specs"
}

type CronSpec struct {
	ID           int32     `toml:"-"`
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
	EVMChainID          *utils.Big `toml:"evmChainID"`
}

type FluxMonitorSpec struct {
	ID              int32               `toml:"-"`
	ContractAddress ethkey.EIP55Address `toml:"contractAddress"`
	Threshold       float32             `toml:"threshold,float"`
	// AbsoluteThreshold is the maximum absolute change allowed in a fluxmonitored
	// value before a new round should be kicked off, so that the current value
	// can be reported on-chain.
	AbsoluteThreshold   float32 `toml:"absoluteThreshold,float"`
	PollTimerPeriod     time.Duration
	PollTimerDisabled   bool
	IdleTimerPeriod     time.Duration
	IdleTimerDisabled   bool
	DrumbeatSchedule    string
	DrumbeatRandomDelay time.Duration
	DrumbeatEnabled     bool
	MinPayment          *assets.Link
	EVMChainID          *utils.Big `toml:"evmChainID"`
	CreatedAt           time.Time  `toml:"-"`
	UpdatedAt           time.Time  `toml:"-"`
}

type KeeperSpec struct {
	ID                       int32               `toml:"-"`
	ContractAddress          ethkey.EIP55Address `toml:"contractAddress"`
	MinIncomingConfirmations *uint32             `toml:"minIncomingConfirmations"`
	FromAddress              ethkey.EIP55Address `toml:"fromAddress"`
	EVMChainID               *utils.Big          `toml:"evmChainID"`
	CreatedAt                time.Time           `toml:"-"`
	UpdatedAt                time.Time           `toml:"-"`
}

type VRFSpec struct {
	ID                       int32
	CoordinatorAddress       ethkey.EIP55Address  `toml:"coordinatorAddress"`
	PublicKey                secp256k1.PublicKey  `toml:"publicKey"`
	MinIncomingConfirmations uint32               `toml:"minIncomingConfirmations"`
	ConfirmationsEnv         bool                 `toml:"-"`
	EVMChainID               *utils.Big           `toml:"evmChainID"`
	FromAddress              *ethkey.EIP55Address `toml:"fromAddress"`
	PollPeriod               time.Duration        `toml:"pollPeriod"` // For v2 jobs
	PollPeriodEnv            bool
	RequestedConfsDelay      int64         `toml:"requestedConfsDelay"` // For v2 jobs. Optional, defaults to 0 if not provided.
	RequestTimeout           time.Duration `toml:"requestTimeout"`      // For v2 jobs. Optional, defaults to 24hr if not provided.
	CreatedAt                time.Time     `toml:"-"`
	UpdatedAt                time.Time     `toml:"-"`
}
