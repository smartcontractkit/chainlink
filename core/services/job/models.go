package job

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/utils/tomlutils"
)

const (
	Cron               Type = "cron"
	DirectRequest      Type = "directrequest"
	FluxMonitor        Type = "fluxmonitor"
	OffchainReporting  Type = "offchainreporting"
	OffchainReporting2 Type = "offchainreporting2"
	Keeper             Type = "keeper"
	VRF                Type = "vrf"
	BlockhashStore     Type = "blockhashstore"
	Webhook            Type = "webhook"
	Bootstrap          Type = "bootstrap"
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
		BlockhashStore:     false,
		Bootstrap:          false,
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
		BlockhashStore:     false,
		Bootstrap:          false,
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
		BlockhashStore:     1,
		Bootstrap:          1,
	}
)

type Job struct {
	ID                   int32     `toml:"-"`
	ExternalJobID        uuid.UUID `toml:"externalJobID"`
	OCROracleSpecID      *int32
	OCROracleSpec        *OCROracleSpec
	OCR2OracleSpecID     *int32
	OCR2OracleSpec       *OCR2OracleSpec
	CronSpecID           *int32
	CronSpec             *CronSpec
	DirectRequestSpecID  *int32
	DirectRequestSpec    *DirectRequestSpec
	FluxMonitorSpecID    *int32
	FluxMonitorSpec      *FluxMonitorSpec
	KeeperSpecID         *int32
	KeeperSpec           *KeeperSpec
	VRFSpecID            *int32
	VRFSpec              *VRFSpec
	WebhookSpecID        *int32
	WebhookSpec          *WebhookSpec
	BlockhashStoreSpecID *int32
	BlockhashStoreSpec   *BlockhashStoreSpec
	BootstrapSpec        *BootstrapSpec
	BootstrapSpecID      *int32
	PipelineSpecID       int32
	PipelineSpec         *pipeline.Spec
	JobSpecErrors        []SpecError
	Type                 Type
	SchemaVersion        uint32
	Name                 null.String
	MaxTaskDuration      models.Interval
	Pipeline             pipeline.Pipeline `toml:"observationSource"`
	CreatedAt            time.Time
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
// by taking the 16 bytes underlying the UUID and right padding it.
func (j Job) ExternalIDEncodeBytesToTopic() common.Hash {
	return ExternalJobIDEncodeBytesToTopic(j.ExternalJobID)
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

// OCROracleSpec defines the job spec for OCR jobs.
type OCROracleSpec struct {
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

// GetID is a getter function that returns the ID of the spec.
func (s OCROracleSpec) GetID() string {
	return fmt.Sprintf("%v", s.ID)
}

// SetID is a setter function that sets the ID of the spec.
func (s *OCROracleSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

// JSONConfig is a Go mapping for JSON based database properties.
type JSONConfig map[string]interface{}

// Bytes returns the raw bytes
func (r JSONConfig) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

// Value returns this instance serialized for database storage.
func (r JSONConfig) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan reads the database value and returns an instance.
func (r *JSONConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("expected bytes got %T", b)
	}
	return json.Unmarshal(b, &r)
}

// OCR2PluginType defines supported OCR2 plugin types.
type OCR2PluginType string

const (
	// Median refers to the median.Median type
	Median OCR2PluginType = "median"
)

// OCR2OracleSpec defines the job spec for OCR2 jobs.
// Relay config is chain specific config for a relay (chain adapter).
type OCR2OracleSpec struct {
	ID                                int32              `toml:"-"`
	ContractID                        string             `toml:"contractID"`
	Relay                             relaytypes.Network `toml:"relay"`
	RelayConfig                       JSONConfig         `toml:"relayConfig"`
	P2PBootstrapPeers                 pq.StringArray     `toml:"p2pBootstrapPeers"`
	OCRKeyBundleID                    null.String        `toml:"ocrKeyBundleID"`
	MonitoringEndpoint                null.String        `toml:"monitoringEndpoint"`
	TransmitterID                     null.String        `toml:"transmitterID"`
	BlockchainTimeout                 models.Interval    `toml:"blockchainTimeout"`
	ContractConfigTrackerPollInterval models.Interval    `toml:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations       uint16             `toml:"contractConfigConfirmations"`
	PluginConfig                      JSONConfig         `toml:"pluginConfig"`
	PluginType                        OCR2PluginType     `toml:"pluginType"`
	CreatedAt                         time.Time          `toml:"-"`
	UpdatedAt                         time.Time          `toml:"-"`
}

// GetID is a getter function that returns the ID of the spec.
func (s OCR2OracleSpec) GetID() string {
	return fmt.Sprintf("%v", s.ID)
}

// SetID is a setter function that sets the ID of the spec.
func (s *OCR2OracleSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
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

type FluxMonitorSpec struct {
	ID              int32               `toml:"-"`
	ContractAddress ethkey.EIP55Address `toml:"contractAddress"`
	Threshold       tomlutils.Float32   `toml:"threshold,float"`
	// AbsoluteThreshold is the maximum absolute change allowed in a fluxmonitored
	// value before a new round should be kicked off, so that the current value
	// can be reported on-chain.
	AbsoluteThreshold   tomlutils.Float32 `toml:"absoluteThreshold,float"`
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
	ID int32

	// BatchCoordinatorAddress is the address of the batch vrf coordinator to use.
	// This is required if batchFulfillmentEnabled is set to true in the job spec.
	BatchCoordinatorAddress *ethkey.EIP55Address `toml:"batchCoordinatorAddress"`
	// BatchFulfillmentEnabled indicates to the vrf job to use the batch vrf coordinator
	// for fulfilling requests. If set to true, batchCoordinatorAddress must be set in
	// the job spec.
	BatchFulfillmentEnabled bool `toml:"batchFulfillmentEnabled"`
	// BatchFulfillmentGasMultiplier is used to determine the final gas estimate for the batch
	// fulfillment.
	BatchFulfillmentGasMultiplier tomlutils.Float64 `toml:"batchFulfillmentGasMultiplier"`

	CoordinatorAddress       ethkey.EIP55Address   `toml:"coordinatorAddress"`
	PublicKey                secp256k1.PublicKey   `toml:"publicKey"`
	MinIncomingConfirmations uint32                `toml:"minIncomingConfirmations"`
	ConfirmationsEnv         bool                  `toml:"-"`
	EVMChainID               *utils.Big            `toml:"evmChainID"`
	FromAddresses            []ethkey.EIP55Address `toml:"fromAddresses"`
	PollPeriod               time.Duration         `toml:"pollPeriod"` // For v2 jobs
	PollPeriodEnv            bool
	RequestedConfsDelay      int64         `toml:"requestedConfsDelay"` // For v2 jobs. Optional, defaults to 0 if not provided.
	RequestTimeout           time.Duration `toml:"requestTimeout"`      // Optional, defaults to 24hr if not provided.

	// ChunkSize is the number of pending VRF V2 requests to process in parallel. Optional, defaults
	// to 20 if not provided.
	ChunkSize uint32 `toml:"chunkSize"`

	// BackoffInitialDelay is the amount of time to wait before retrying a failed request after the
	// first failure. V2 only.
	BackoffInitialDelay time.Duration `toml:"backoffInitialDelay"`

	// BackoffMaxDelay is the maximum amount of time to wait before retrying a failed request. V2
	// only.
	BackoffMaxDelay time.Duration `toml:"backoffMaxDelay"`

	CreatedAt time.Time `toml:"-"`
	UpdatedAt time.Time `toml:"-"`
}

// BlockhashStoreSpec defines the job spec for the blockhash store feeder.
type BlockhashStoreSpec struct {
	ID int32

	// CoordinatorV1Address is the VRF V1 coordinator to watch for unfulfilled requests. If empty,
	// no V1 coordinator will be watched.
	CoordinatorV1Address *ethkey.EIP55Address `toml:"coordinatorV1Address"`

	// CoordinatorV2Address is the VRF V2 coordinator to watch for unfulfilled requests. If empty,
	// no V2 coordinator will be watched.
	CoordinatorV2Address *ethkey.EIP55Address `toml:"coordinatorV2Address"`

	// WaitBlocks defines the number of blocks to wait before a hash is stored.
	WaitBlocks int32 `toml:"waitBlocks"`

	// LookbackBlocks defines the maximum age of blocks whose hashes should be stored.
	LookbackBlocks int32 `toml:"lookbackBlocks"`

	// BlockhashStoreAddress is the address of the BlockhashStore contract to store blockhashes
	// into.
	BlockhashStoreAddress ethkey.EIP55Address `toml:"blockhashStoreAddress"`

	// PollPeriod defines how often recent blocks should be scanned for blockhash storage.
	PollPeriod time.Duration `toml:"pollPeriod"`

	// RunTimeout defines the timeout for a single run of the blockhash store feeder.
	RunTimeout time.Duration `toml:"runTimeout"`

	// EVMChainID defines the chain ID for monitoring and storing of blockhashes.
	EVMChainID *utils.Big `toml:"evmChainID"`

	// FromAddress is the sender address that should be used to store blockhashes.
	FromAddress *ethkey.EIP55Address `toml:"fromAddress"`

	// CreatedAt is the time this job was created.
	CreatedAt time.Time `toml:"-"`

	// UpdatedAt is the time this job was last updated.
	UpdatedAt time.Time `toml:"-"`
}

// BootstrapSpec defines the spec to handles the node communication setup process.
type BootstrapSpec struct {
	ID                                int32              `toml:"-"`
	ContractID                        string             `toml:"contractID"`
	Relay                             relaytypes.Network `toml:"relay"`
	RelayConfig                       JSONConfig
	MonitoringEndpoint                null.String     `toml:"monitoringEndpoint"`
	BlockchainTimeout                 models.Interval `toml:"blockchainTimeout"`
	ContractConfigTrackerPollInterval models.Interval `toml:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations       uint16          `toml:"contractConfigConfirmations"`
	CreatedAt                         time.Time       `toml:"-"`
	UpdatedAt                         time.Time       `toml:"-"`
}

// AsOCR2Spec transforms the bootstrap spec into a generic OCR2 format to enable code sharing between specs.
func (s BootstrapSpec) AsOCR2Spec() OCR2OracleSpec {
	return OCR2OracleSpec{
		ID:                                s.ID,
		ContractID:                        s.ContractID,
		Relay:                             s.Relay,
		RelayConfig:                       s.RelayConfig,
		MonitoringEndpoint:                s.MonitoringEndpoint,
		BlockchainTimeout:                 s.BlockchainTimeout,
		ContractConfigTrackerPollInterval: s.ContractConfigTrackerPollInterval,
		ContractConfigConfirmations:       s.ContractConfigConfirmations,
		CreatedAt:                         s.CreatedAt,
		UpdatedAt:                         s.UpdatedAt,
	}
}
