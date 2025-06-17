package job

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/tomlutils"
)

const (
	BlockHeaderFeeder       Type = (Type)(pipeline.BlockHeaderFeederJobType)
	BlockhashStore          Type = (Type)(pipeline.BlockhashStoreJobType)
	Bootstrap               Type = (Type)(pipeline.BootstrapJobType)
	Cron                    Type = (Type)(pipeline.CronJobType)
	CCIP                    Type = (Type)(pipeline.CCIPJobType)
	DirectRequest           Type = (Type)(pipeline.DirectRequestJobType)
	FluxMonitor             Type = (Type)(pipeline.FluxMonitorJobType)
	Gateway                 Type = (Type)(pipeline.GatewayJobType)
	Keeper                  Type = (Type)(pipeline.KeeperJobType)
	LegacyGasStationServer  Type = (Type)(pipeline.LegacyGasStationServerJobType)
	LegacyGasStationSidecar Type = (Type)(pipeline.LegacyGasStationSidecarJobType)
	OffchainReporting       Type = (Type)(pipeline.OffchainReportingJobType)
	OffchainReporting2      Type = (Type)(pipeline.OffchainReporting2JobType)
	Stream                  Type = (Type)(pipeline.StreamJobType)
	VRF                     Type = (Type)(pipeline.VRFJobType)
	Webhook                 Type = (Type)(pipeline.WebhookJobType)
	Workflow                Type = (Type)(pipeline.WorkflowJobType)
	StandardCapabilities    Type = (Type)(pipeline.StandardCapabilitiesJobType)
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
		BlockHeaderFeeder:       false,
		BlockhashStore:          false,
		Bootstrap:               false,
		Cron:                    true,
		CCIP:                    false,
		DirectRequest:           true,
		FluxMonitor:             true,
		Gateway:                 false,
		Keeper:                  false, // observationSource is injected in the upkeep executor
		LegacyGasStationServer:  false,
		LegacyGasStationSidecar: false,
		OffchainReporting2:      false, // bootstrap jobs do not require it
		OffchainReporting:       false, // bootstrap jobs do not require it
		Stream:                  true,
		VRF:                     true,
		Webhook:                 true,
		Workflow:                false,
		StandardCapabilities:    false,
	}
	supportsAsync = map[Type]bool{
		BlockHeaderFeeder:       false,
		BlockhashStore:          false,
		Bootstrap:               false,
		Cron:                    true,
		CCIP:                    false,
		DirectRequest:           true,
		FluxMonitor:             false,
		Gateway:                 false,
		Keeper:                  true,
		LegacyGasStationServer:  false,
		LegacyGasStationSidecar: false,
		OffchainReporting2:      false,
		OffchainReporting:       false,
		Stream:                  true,
		VRF:                     true,
		Webhook:                 true,
		Workflow:                false,
		StandardCapabilities:    false,
	}
	schemaVersions = map[Type]uint32{
		BlockHeaderFeeder:       1,
		BlockhashStore:          1,
		Bootstrap:               1,
		Cron:                    1,
		CCIP:                    1,
		DirectRequest:           1,
		FluxMonitor:             1,
		Gateway:                 1,
		Keeper:                  1,
		LegacyGasStationServer:  1,
		LegacyGasStationSidecar: 1,
		OffchainReporting2:      1,
		OffchainReporting:       1,
		Stream:                  1,
		VRF:                     1,
		Webhook:                 1,
		Workflow:                1,
		StandardCapabilities:    1,
	}
)

type Job struct {
	ID                            int32     `toml:"-"`
	ExternalJobID                 uuid.UUID `toml:"externalJobID"`
	StreamID                      *uint32   `toml:"streamID"`
	OCROracleSpecID               *int32
	OCROracleSpec                 *OCROracleSpec
	OCR2OracleSpecID              *int32
	OCR2OracleSpec                *OCR2OracleSpec
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
	BlockhashStoreSpecID          *int32
	BlockhashStoreSpec            *BlockhashStoreSpec
	BlockHeaderFeederSpecID       *int32
	BlockHeaderFeederSpec         *BlockHeaderFeederSpec
	BALSpecID                     *int32
	LegacyGasStationServerSpecID  *int32
	LegacyGasStationServerSpec    *LegacyGasStationServerSpec
	LegacyGasStationSidecarSpecID *int32
	LegacyGasStationSidecarSpec   *LegacyGasStationSidecarSpec
	BootstrapSpec                 *BootstrapSpec
	BootstrapSpecID               *int32
	GatewaySpec                   *GatewaySpec
	GatewaySpecID                 *int32
	EALSpec                       *EALSpec
	EALSpecID                     *int32
	LiquidityBalancerSpec         *LiquidityBalancerSpec
	LiquidityBalancerSpecID       *int32
	PipelineSpecID                int32 // This is deprecated in favor of the `job_pipeline_specs` table relationship
	PipelineSpec                  *pipeline.Spec
	WorkflowSpecID                *int32
	WorkflowSpec                  *WorkflowSpec
	StandardCapabilitiesSpecID    *int32
	StandardCapabilitiesSpec      *StandardCapabilitiesSpec
	CCIPSpecID                    *int32
	CCIPSpec                      *CCIPSpec
	CCIPBootstrapSpecID           *int32
	JobSpecErrors                 []SpecError
	Type                          Type          `toml:"type"`
	SchemaVersion                 uint32        `toml:"schemaVersion"`
	GasLimit                      clnull.Uint32 `toml:"gasLimit"`
	ForwardingAllowed             bool          `toml:"forwardingAllowed"`
	Name                          null.String   `toml:"name"`
	MaxTaskDuration               models.Interval
	Pipeline                      pipeline.Pipeline `toml:"observationSource"`
	CreatedAt                     time.Time
}

func ExternalJobIDEncodeStringToTopic(id uuid.UUID) common.Hash {
	return common.BytesToHash([]byte(strings.Replace(id.String(), "-", "", 4)))
}

func ExternalJobIDEncodeBytesToTopic(id uuid.UUID) common.Hash {
	return common.BytesToHash(common.RightPadBytes(id[:], utils.EVMWordByteLen))
}

// ExternalIDEncodeStringToTopic encodes the external job ID (UUID) into a log topic (32 bytes)
// by taking the string representation of the UUID, removing the dashes
// so that its 32 characters long and then encoding those characters to bytes.
func (j Job) ExternalIDEncodeStringToTopic() common.Hash {
	return ExternalJobIDEncodeStringToTopic(j.ExternalJobID)
}

// ExternalIDEncodeBytesToTopic encodes the external job ID (UUID) into a log topic (32 bytes)
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

type PipelineSpec struct {
	JobID          int32 `json:"-"`
	PipelineSpecID int32 `json:"-"`
	IsPrimary      bool  `json:"is_primary"`
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
	ID         int64 `json:"-"`
	PruningKey int64 `json:"-"`
}

func (pr PipelineRun) GetID() string {
	return strconv.FormatInt(pr.ID, 10)
}

func (pr *PipelineRun) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	pr.ID = ID
	return nil
}

// OCROracleSpec defines the job spec for OCR jobs.
type OCROracleSpec struct {
	ID                                     int32                  `toml:"-"`
	ContractAddress                        evmtypes.EIP55Address  `toml:"contractAddress"`
	P2PV2Bootstrappers                     pq.StringArray         `toml:"p2pv2Bootstrappers" db:"p2pv2_bootstrappers"`
	IsBootstrapPeer                        bool                   `toml:"isBootstrapPeer"`
	EncryptedOCRKeyBundleID                *models.Sha256Hash     `toml:"keyBundleID"`
	TransmitterAddress                     *evmtypes.EIP55Address `toml:"transmitterAddress"`
	ObservationTimeout                     models.Interval        `toml:"observationTimeout"`
	BlockchainTimeout                      models.Interval        `toml:"blockchainTimeout"`
	ContractConfigTrackerSubscribeInterval models.Interval        `toml:"contractConfigTrackerSubscribeInterval"`
	ContractConfigTrackerPollInterval      models.Interval        `toml:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations            uint16                 `toml:"contractConfigConfirmations"`
	EVMChainID                             *big.Big               `toml:"evmChainID" db:"evm_chain_id"`
	DatabaseTimeout                        *models.Interval       `toml:"databaseTimeout"`
	ObservationGracePeriod                 *models.Interval       `toml:"observationGracePeriod"`
	ContractTransmitterTransmitTimeout     *models.Interval       `toml:"contractTransmitterTransmitTimeout"`
	CaptureEATelemetry                     bool                   `toml:"captureEATelemetry"`
	CreatedAt                              time.Time              `toml:"-"`
	UpdatedAt                              time.Time              `toml:"-"`
}

// GetID is a getter function that returns the ID of the spec.
func (s OCROracleSpec) GetID() string {
	return strconv.Itoa(int(s.ID))
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

// JSONConfig is a map for config properties which are encoded as JSON in the database by implementing
// sql.Scanner and driver.Valuer.
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

func (r JSONConfig) MercuryCredentialName() (string, error) {
	url, ok := r["mercuryCredentialName"]
	if !ok {
		return "", nil
	}
	name, ok := url.(string)
	if !ok {
		return "", fmt.Errorf("expected string mercuryCredentialName but got: %T", url)
	}
	return name, nil
}

func (r JSONConfig) ApplyDefaultsOCR2(cfg ocr2Config) {
	_, ok := r["defaultTransactionQueueDepth"]
	if !ok {
		r["defaultTransactionQueueDepth"] = cfg.DefaultTransactionQueueDepth()
	}
	_, ok = r["simulateTransactions"]
	if !ok {
		r["simulateTransactions"] = cfg.SimulateTransactions()
	}
}

type ocr2Config interface {
	DefaultTransactionQueueDepth() uint32
	SimulateTransactions() bool
}

var ForwardersSupportedPlugins = []types.OCR2PluginType{types.Median, types.OCR2Keeper, types.Functions}

// OCR2OracleSpec defines the job spec for OCR2 jobs.
// Relay config is chain specific config for a relay (chain adapter).
type OCR2OracleSpec struct {
	ID         int32        `toml:"-"`
	ContractID string       `toml:"contractID"`
	FeedID     *common.Hash `toml:"feedID"`
	// Network
	Relay string `toml:"relay"`
	// TODO BCF-2442 implement ChainID as top level parameter rathe than buried in RelayConfig.
	ChainID                           string               `toml:"chainID"`
	RelayConfig                       JSONConfig           `toml:"relayConfig"`
	P2PV2Bootstrappers                pq.StringArray       `toml:"p2pv2Bootstrappers"`
	OCRKeyBundleID                    null.String          `toml:"ocrKeyBundleID"`
	MonitoringEndpoint                null.String          `toml:"monitoringEndpoint"`
	TransmitterID                     null.String          `toml:"transmitterID"`
	BlockchainTimeout                 models.Interval      `toml:"blockchainTimeout"`
	ContractConfigTrackerPollInterval models.Interval      `toml:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations       uint16               `toml:"contractConfigConfirmations"`
	OnchainSigningStrategy            JSONConfig           `toml:"onchainSigningStrategy"`
	PluginConfig                      JSONConfig           `toml:"pluginConfig"`
	PluginType                        types.OCR2PluginType `toml:"pluginType"`
	CreatedAt                         time.Time            `toml:"-"`
	UpdatedAt                         time.Time            `toml:"-"`
	CaptureEATelemetry                bool                 `toml:"captureEATelemetry"`
	CaptureAutomationCustomTelemetry  bool                 `toml:"captureAutomationCustomTelemetry"`
	// AllowNoBootstrappers is a flag that allows the job to start without any bootstrappers
	// This is useful for testing and deployments where the node is not configured to conduct consensus (i.e. f = 0 and n = 1).
	AllowNoBootstrappers bool `toml:"allowNoBootstrappers"`
}

func validateRelayID(id types.RelayID) error {
	// only the EVM has specific requirements
	if id.Network == relay.NetworkEVM {
		_, err := toml.ChainIDInt64(id.ChainID)
		if err != nil {
			return fmt.Errorf("invalid EVM chain id %s: %w", id.ChainID, err)
		}
	}
	return nil
}

func (s *OCR2OracleSpec) RelayID() (types.RelayID, error) {
	cid, err := s.getChainID()
	if err != nil {
		return types.RelayID{}, err
	}
	rid := types.NewRelayID(s.Relay, cid)
	err = validateRelayID(rid)
	if err != nil {
		return types.RelayID{}, err
	}
	return rid, nil
}

func (s *OCR2OracleSpec) getChainID() (string, error) {
	if s.ChainID != "" {
		return s.ChainID, nil
	}
	// backward compatible job spec
	return s.getChainIdFromRelayConfig()
}

func (s *OCR2OracleSpec) getChainIdFromRelayConfig() (string, error) {
	v, exists := s.RelayConfig["chainID"]
	if !exists {
		return "", errors.New("chainID does not exist")
	}
	switch t := v.(type) {
	case string:
		return t, nil
	case int, int64, int32:
		return fmt.Sprintf("%d", v), nil
	case float64:
		// backward compatibility with JSONConfig.EVMChainID
		i := int64(t)
		return strconv.FormatInt(i, 10), nil

	default:
		return "", fmt.Errorf("unable to parse chainID: unexpected type %T", t)
	}
}

// GetID is a getter function that returns the ID of the spec.
func (s OCR2OracleSpec) GetID() string {
	return strconv.Itoa(int(s.ID))
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
	return strconv.Itoa(int(w.ID))
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
	ID                       int32                    `toml:"-"`
	ContractAddress          evmtypes.EIP55Address    `toml:"contractAddress"`
	MinIncomingConfirmations clnull.Uint32            `toml:"minIncomingConfirmations"`
	Requesters               models.AddressCollection `toml:"requesters"`
	MinContractPayment       *commonassets.Link       `toml:"minContractPaymentLinkJuels"`
	EVMChainID               *big.Big                 `toml:"evmChainID"`
	CreatedAt                time.Time                `toml:"-"`
	UpdatedAt                time.Time                `toml:"-"`
}

type CronSpec struct {
	ID           int32     `toml:"-"`
	CronSchedule string    `toml:"schedule"`
	EVMChainID   *big.Big  `toml:"evmChainID"`
	CreatedAt    time.Time `toml:"-"`
	UpdatedAt    time.Time `toml:"-"`
}

func (s CronSpec) GetID() string {
	return strconv.Itoa(int(s.ID))
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
	ID              int32                 `toml:"-"`
	ContractAddress evmtypes.EIP55Address `toml:"contractAddress"`
	Threshold       tomlutils.Float32     `toml:"threshold,float"`
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
	MinPayment          *commonassets.Link
	EVMChainID          *big.Big  `toml:"evmChainID"`
	CreatedAt           time.Time `toml:"-"`
	UpdatedAt           time.Time `toml:"-"`
}

type KeeperSpec struct {
	ID                       int32                 `toml:"-"`
	ContractAddress          evmtypes.EIP55Address `toml:"contractAddress"`
	MinIncomingConfirmations *uint32               `toml:"minIncomingConfirmations"`
	FromAddress              evmtypes.EIP55Address `toml:"fromAddress"`
	EVMChainID               *big.Big              `toml:"evmChainID"`
	CreatedAt                time.Time             `toml:"-"`
	UpdatedAt                time.Time             `toml:"-"`
}

type VRFSpec struct {
	ID int32

	// BatchCoordinatorAddress is the address of the batch vrf coordinator to use.
	// This is required if batchFulfillmentEnabled is set to true in the job spec.
	BatchCoordinatorAddress *evmtypes.EIP55Address `toml:"batchCoordinatorAddress"`
	// BatchFulfillmentEnabled indicates to the vrf job to use the batch vrf coordinator
	// for fulfilling requests. If set to true, batchCoordinatorAddress must be set in
	// the job spec.
	BatchFulfillmentEnabled bool `toml:"batchFulfillmentEnabled"`
	// CustomRevertsPipelineEnabled indicates to the vrf job to run the
	// custom reverted txns pipeline along with VRF listener
	CustomRevertsPipelineEnabled bool `toml:"customRevertsPipelineEnabled"`
	// BatchFulfillmentGasMultiplier is used to determine the final gas estimate for the batch
	// fulfillment.
	BatchFulfillmentGasMultiplier tomlutils.Float64 `toml:"batchFulfillmentGasMultiplier"`

	// VRFOwnerAddress is the address of the VRFOwner address to use.
	//
	// V2 only.
	VRFOwnerAddress *evmtypes.EIP55Address `toml:"vrfOwnerAddress"`

	CoordinatorAddress       evmtypes.EIP55Address   `toml:"coordinatorAddress"`
	PublicKey                secp256k1.PublicKey     `toml:"publicKey"`
	MinIncomingConfirmations uint32                  `toml:"minIncomingConfirmations"`
	EVMChainID               *big.Big                `toml:"evmChainID"`
	FromAddresses            []evmtypes.EIP55Address `toml:"fromAddresses"`
	PollPeriod               time.Duration           `toml:"pollPeriod"`          // For v2 jobs
	RequestedConfsDelay      int64                   `toml:"requestedConfsDelay"` // For v2 jobs. Optional, defaults to 0 if not provided.
	RequestTimeout           time.Duration           `toml:"requestTimeout"`      // Optional, defaults to 24hr if not provided.

	// GasLanePrice specifies the gas lane price for this VRF job.
	// If the specified keys in FromAddresses do not have the provided gas price the job
	// will not start.
	//
	// Optional, for v2 jobs only.
	GasLanePrice *assets.Wei `toml:"gasLanePrice" db:"gas_lane_price"`

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
	CoordinatorV1Address *evmtypes.EIP55Address `toml:"coordinatorV1Address"`

	// CoordinatorV2Address is the VRF V2 coordinator to watch for unfulfilled requests. If empty,
	// no V2 coordinator will be watched.
	CoordinatorV2Address *evmtypes.EIP55Address `toml:"coordinatorV2Address"`

	// CoordinatorV2PlusAddress is the VRF V2Plus coordinator to watch for unfulfilled requests. If empty,
	// no V2Plus coordinator will be watched.
	CoordinatorV2PlusAddress *evmtypes.EIP55Address `toml:"coordinatorV2PlusAddress"`

	// LookbackBlocks defines the maximum age of blocks whose hashes should be stored.
	LookbackBlocks int32 `toml:"lookbackBlocks"`

	// WaitBlocks defines the minimum age of blocks whose hashes should be stored.
	WaitBlocks int32 `toml:"waitBlocks"`

	// HeartbeatPeriodTime defines the number of seconds by which we "heartbeat store"
	// a blockhash into the blockhash store contract.
	// This is so that we always have a blockhash to anchor to in the event we need to do a
	// backwards mode on the contract.
	HeartbeatPeriod time.Duration `toml:"heartbeatPeriod"`

	// BlockhashStoreAddress is the address of the BlockhashStore contract to store blockhashes
	// into.
	BlockhashStoreAddress evmtypes.EIP55Address `toml:"blockhashStoreAddress"`

	// BatchBlockhashStoreAddress is the address of the trusted BlockhashStore contract to store blockhashes
	TrustedBlockhashStoreAddress *evmtypes.EIP55Address `toml:"trustedBlockhashStoreAddress"`

	// BatchBlockhashStoreBatchSize is the number of blockhashes to store in a single batch
	TrustedBlockhashStoreBatchSize int32 `toml:"trustedBlockhashStoreBatchSize"`

	// PollPeriod defines how often recent blocks should be scanned for blockhash storage.
	PollPeriod time.Duration `toml:"pollPeriod"`

	// RunTimeout defines the timeout for a single run of the blockhash store feeder.
	RunTimeout time.Duration `toml:"runTimeout"`

	// EVMChainID defines the chain ID for monitoring and storing of blockhashes.
	EVMChainID *big.Big `toml:"evmChainID"`

	// FromAddress is the sender address that should be used to store blockhashes.
	FromAddresses []evmtypes.EIP55Address `toml:"fromAddresses"`

	// CreatedAt is the time this job was created.
	CreatedAt time.Time `toml:"-"`

	// UpdatedAt is the time this job was last updated.
	UpdatedAt time.Time `toml:"-"`
}

// BlockHeaderFeederSpec defines the job spec for the blockhash store feeder.
type BlockHeaderFeederSpec struct {
	ID int32

	// CoordinatorV1Address is the VRF V1 coordinator to watch for unfulfilled requests. If empty,
	// no V1 coordinator will be watched.
	CoordinatorV1Address *evmtypes.EIP55Address `toml:"coordinatorV1Address"`

	// CoordinatorV2Address is the VRF V2 coordinator to watch for unfulfilled requests. If empty,
	// no V2 coordinator will be watched.
	CoordinatorV2Address *evmtypes.EIP55Address `toml:"coordinatorV2Address"`

	// CoordinatorV2PlusAddress is the VRF V2Plus coordinator to watch for unfulfilled requests. If empty,
	// no V2Plus coordinator will be watched.
	CoordinatorV2PlusAddress *evmtypes.EIP55Address `toml:"coordinatorV2PlusAddress"`

	// LookbackBlocks defines the maximum age of blocks whose hashes should be stored.
	LookbackBlocks int32 `toml:"lookbackBlocks"`

	// WaitBlocks defines the minimum age of blocks whose hashes should be stored.
	WaitBlocks int32 `toml:"waitBlocks"`

	// BlockhashStoreAddress is the address of the BlockhashStore contract to store blockhashes
	// into.
	BlockhashStoreAddress evmtypes.EIP55Address `toml:"blockhashStoreAddress"`

	// BatchBlockhashStoreAddress is the address of the BatchBlockhashStore contract to store blockhashes
	// into.
	BatchBlockhashStoreAddress evmtypes.EIP55Address `toml:"batchBlockhashStoreAddress"`

	// PollPeriod defines how often recent blocks should be scanned for blockhash storage.
	PollPeriod time.Duration `toml:"pollPeriod"`

	// RunTimeout defines the timeout for a single run of the blockhash store feeder.
	RunTimeout time.Duration `toml:"runTimeout"`

	// EVMChainID defines the chain ID for monitoring and storing of blockhashes.
	EVMChainID *big.Big `toml:"evmChainID"`

	// FromAddress is the sender address that should be used to store blockhashes.
	FromAddresses []evmtypes.EIP55Address `toml:"fromAddresses"`

	// GetBlockHashesBatchSize is the RPC call batch size for retrieving blockhashes
	GetBlockhashesBatchSize uint16 `toml:"getBlockhashesBatchSize"`

	// StoreBlockhashesBatchSize is the RPC call batch size for storing blockhashes
	StoreBlockhashesBatchSize uint16 `toml:"storeBlockhashesBatchSize"`

	// CreatedAt is the time this job was created.
	CreatedAt time.Time `toml:"-"`

	// UpdatedAt is the time this job was last updated.
	UpdatedAt time.Time `toml:"-"`
}

// LegacyGasStationServerSpec defines the job spec for the legacy gas station server.
type LegacyGasStationServerSpec struct {
	ID int32

	// ForwarderAddress is the address of EIP2771 forwarder that verifies signature
	// and forwards requests to target contracts
	ForwarderAddress evmtypes.EIP55Address `toml:"forwarderAddress"`

	// EVMChainID defines the chain ID from which the meta-transaction request originates.
	EVMChainID *big.Big `toml:"evmChainID"`

	// CCIPChainSelector is the CCIP chain selector that corresponds to EVMChainID param.
	// This selector is equivalent to (source) chainID specified in SendTransaction request
	CCIPChainSelector *big.Big `toml:"ccipChainSelector"`

	// FromAddress is the sender address that should be used to send meta-transactions
	FromAddresses []evmtypes.EIP55Address `toml:"fromAddresses"`

	// CreatedAt is the time this job was created.
	CreatedAt time.Time `toml:"-"`

	// UpdatedAt is the time this job was last updated.
	UpdatedAt time.Time `toml:"-"`
}

// LegacyGasStationSidecarSpec defines the job spec for the legacy gas station sidecar.
type LegacyGasStationSidecarSpec struct {
	ID int32

	// ForwarderAddress is the address of EIP2771 forwarder that verifies signature
	// and forwards requests to target contracts
	ForwarderAddress evmtypes.EIP55Address `toml:"forwarderAddress"`

	// OffRampAddress is the address of CCIP OffRamp for the given chainID
	OffRampAddress evmtypes.EIP55Address `toml:"offRampAddress"`

	// LookbackBlocks defines the maximum number of blocks to search for on-chain events.
	LookbackBlocks int32 `toml:"lookbackBlocks"`

	// PollPeriod defines how frequently legacy gas station sidecar runs.
	PollPeriod time.Duration `toml:"pollPeriod"`

	// RunTimeout defines the timeout for a single run of the legacy gas station sidecar.
	RunTimeout time.Duration `toml:"runTimeout"`

	// EVMChainID defines the chain ID for the on-chain events tracked by sidecar
	EVMChainID *big.Big `toml:"evmChainID"`

	// CCIPChainSelector is the CCIP chain selector that corresponds to EVMChainID param
	CCIPChainSelector *big.Big `toml:"ccipChainSelector"`

	// CreatedAt is the time this job was created.
	CreatedAt time.Time `toml:"-"`

	// UpdatedAt is the time this job was last updated.
	UpdatedAt time.Time `toml:"-"`
}

// BootstrapSpec defines the spec to handles the node communication setup process.
type BootstrapSpec struct {
	ID                                int32        `toml:"-"`
	ContractID                        string       `toml:"contractID"`
	FeedID                            *common.Hash `toml:"feedID"`
	Relay                             string       `toml:"relay"` // RelayID.Network
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
		P2PV2Bootstrappers:                pq.StringArray{},
	}
}

type GatewaySpec struct {
	ID            int32      `toml:"-"`
	GatewayConfig JSONConfig `toml:"gatewayConfig"`
	CreatedAt     time.Time  `toml:"-"`
	UpdatedAt     time.Time  `toml:"-"`
}

func (s GatewaySpec) GetID() string {
	return strconv.Itoa(int(s.ID))
}

func (s *GatewaySpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

// AuthGatewayID returns AuthGatewayId or empty string, if not found or it's not a string
func (s *GatewaySpec) AuthGatewayID() string {
	// not using config.GatewayConfig directly to avoid import cycle
	if nsc, ok := s.GatewayConfig["ConnectionManagerConfig"]; ok {
		if nscMap, ok := nsc.(map[string]interface{}); ok {
			if authGatewayID, ok := nscMap["AuthGatewayId"]; ok {
				if authGatewayIDStr, ok := authGatewayID.(string); ok {
					return authGatewayIDStr
				}
			}
		}
	}

	return ""
}

// EALSpec defines the job spec for the gas station.
type EALSpec struct {
	ID int32

	// ForwarderAddress is the address of EIP2771 forwarder that verifies signature
	// and forwards requests to target contracts
	ForwarderAddress evmtypes.EIP55Address `toml:"forwarderAddress"`

	// EVMChainID defines the chain ID from which the meta-transaction request originates.
	EVMChainID *big.Big `toml:"evmChainID"`

	// FromAddress is the sender address that should be used to send meta-transactions
	FromAddresses []evmtypes.EIP55Address `toml:"fromAddresses"`

	// LookbackBlocks defines the maximum age of blocks to lookback in status tracker
	LookbackBlocks int32 `toml:"lookbackBlocks"`

	// PollPeriod defines how frequently EAL status tracker runs
	PollPeriod time.Duration `toml:"pollPeriod"`

	// RunTimeout defines the timeout for a single run of EAL status tracker
	RunTimeout time.Duration `toml:"runTimeout"`

	// CreatedAt is the time this job was created.
	CreatedAt time.Time `toml:"-"`

	// UpdatedAt is the time this job was last updated.
	UpdatedAt time.Time `toml:"-"`
}

type LiquidityBalancerSpec struct {
	ID int32

	LiquidityBalancerConfig string `toml:"liquidityBalancerConfig" db:"liquidity_balancer_config"`
}

type WorkflowSpecType string

const (
	YamlSpec        WorkflowSpecType = "yaml"
	WASMFile        WorkflowSpecType = "wasm_file"
	DefaultSpecType                  = ""
)

type WorkflowSpecStatus string

const (
	WorkflowSpecStatusActive  WorkflowSpecStatus = "active"
	WorkflowSpecStatusPaused  WorkflowSpecStatus = "paused"
	WorkflowSpecStatusDefault WorkflowSpecStatus = ""
)

type WorkflowSpec struct {
	ID       int32  `toml:"-"`
	Workflow string `toml:"workflow"`           // the raw representation of the workflow
	Config   string `toml:"config" db:"config"` // the raw representation of the config
	// fields derived from the yaml spec, used for indexing the database
	// note: i tried to make these private, but translating them to the database seems to require them to be public
	WorkflowID    string             `toml:"-" db:"workflow_id"` // Derived. Do not modify. the CID of the workflow.
	WorkflowOwner string             `toml:"workflow_owner" db:"workflow_owner"`
	WorkflowName  string             `toml:"workflow_name" db:"workflow_name"`
	Status        WorkflowSpecStatus `db:"status"`
	BinaryURL     string             `db:"binary_url"`
	ConfigURL     string             `db:"config_url"`
	SecretsID     sql.NullInt64      `db:"secrets_id"`
	CreatedAt     time.Time          `toml:"-"`
	UpdatedAt     time.Time          `toml:"-"`
	SpecType      WorkflowSpecType   `toml:"spec_type" db:"spec_type"`
	sdkWorkflow   *sdk.WorkflowSpec
	rawSpec       []byte
	config        []byte
}

var (
	ErrInvalidWorkflowID       = errors.New("invalid workflow id")
	ErrInvalidWorkflowYAMLSpec = errors.New("invalid workflow yaml spec")
)

const (
	workflowIDLen = 64 // sha256 hash
)

// Validate checks the workflow spec for correctness
func (w *WorkflowSpec) Validate(ctx context.Context) error {
	s, err := w.SDKSpec(ctx)
	if err != nil {
		return err
	}

	// For yaml-based workflow specs, use the owner & name fields defined there.
	// For wasm workflows, use the `workflow_name` & `workflow_owner` fields directly from the job spec.
	if s.Owner+s.Name != "" {
		w.WorkflowOwner = strings.TrimPrefix(s.Owner, "0x") // the json schema validation ensures it is a hex string with 0x prefix, but the database does not store the prefix
		w.WorkflowName = s.Name
	} else {
		w.WorkflowOwner = strings.TrimPrefix(w.WorkflowOwner, "0x")
	}

	if len(w.WorkflowID) != workflowIDLen {
		return fmt.Errorf("%w: incorrect length for id %s: expected %d, got %d", ErrInvalidWorkflowID, w.WorkflowID, workflowIDLen, len(w.WorkflowID))
	}

	return nil
}

func (w *WorkflowSpec) SDKSpec(ctx context.Context) (sdk.WorkflowSpec, error) {
	if w.sdkWorkflow != nil {
		return *w.sdkWorkflow, nil
	}

	workflowSpecFactory, ok := workflowSpecFactories[w.SpecType]
	if !ok {
		return sdk.WorkflowSpec{}, fmt.Errorf("unknown spec type %s", w.SpecType)
	}
	spec, rawSpec, cid, err := workflowSpecFactory.Spec(ctx, w.Workflow, w.Config)
	if err != nil {
		return sdk.WorkflowSpec{}, fmt.Errorf("spec factory failed: %w", err)
	}
	w.sdkWorkflow = &spec
	w.rawSpec = rawSpec
	w.WorkflowID = cid
	return spec, nil
}

func (w *WorkflowSpec) RawSpec(ctx context.Context) ([]byte, error) {
	if w.rawSpec != nil {
		return w.rawSpec, nil
	}

	workflowSpecFactory, ok := workflowSpecFactories[w.SpecType]
	if !ok {
		return nil, fmt.Errorf("unknown spec type %s", w.SpecType)
	}

	rs, err := workflowSpecFactory.RawSpec(ctx, w.Workflow, w.Config)
	if err != nil {
		return nil, err
	}

	w.rawSpec = rs
	return rs, nil
}

func (w *WorkflowSpec) GetConfig(ctx context.Context) ([]byte, error) {
	if w.config != nil {
		return w.config, nil
	}

	workflowSpecFactory, ok := workflowSpecFactories[w.SpecType]
	if !ok {
		return nil, fmt.Errorf("unknown spec type %s", w.SpecType)
	}

	rs, err := workflowSpecFactory.Config(ctx, w.Config)
	if err != nil {
		return nil, err
	}

	w.config = rs
	return rs, nil
}

type OracleFactoryConfig struct {
	Enabled            bool     `toml:"enabled"`
	BootstrapPeers     []string `toml:"bootstrap_peers"`      // e.g.,["12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690"]
	OCRContractAddress string   `toml:"ocr_contract_address"` // e.g., 0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6
	ChainID            string   `toml:"chain_id"`             // e.g., "31337"
	Network            string   `toml:"network"`              // e.g., "evm"
}

// Value returns this instance serialized for database storage.
func (ofc OracleFactoryConfig) Value() (driver.Value, error) {
	return json.Marshal(ofc)
}

// Scan reads the database value and returns an instance.
func (ofc *OracleFactoryConfig) Scan(value interface{}) error {
	if value == nil {
		return nil // field is nullable
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("expected bytes got %T", value)
	}
	return json.Unmarshal(b, &ofc)
}

type StandardCapabilitiesSpec struct {
	ID            int32
	CreatedAt     time.Time           `toml:"-"`
	UpdatedAt     time.Time           `toml:"-"`
	Command       string              `toml:"command" db:"command"`
	Config        string              `toml:"config" db:"config"`
	OracleFactory OracleFactoryConfig `toml:"oracle_factory" db:"oracle_factory"`
}

func (w *StandardCapabilitiesSpec) GetID() string {
	return strconv.Itoa(int(w.ID))
}

func (w *StandardCapabilitiesSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	w.ID = int32(ID)
	return nil
}

type CCIPSpec struct {
	ID        int32
	CreatedAt time.Time `toml:"-"`
	UpdatedAt time.Time `toml:"-"`

	// P2PV2Bootstrappers is a list of "peer_id@ip_address:port" strings that are used to
	// identify the bootstrap nodes of the P2P network.
	// These bootstrappers will be used to bootstrap all CCIP DONs.
	P2PV2Bootstrappers pq.StringArray `toml:"p2pV2Bootstrappers" db:"p2pv2_bootstrappers"`

	// CapabilityVersion is the semantic version of the CCIP capability.
	// This capability version must exist in the onchain capability registry.
	CapabilityVersion string `toml:"capabilityVersion" db:"capability_version"`

	// CapabilityLabelledName is the labelled name of the CCIP capability.
	// Corresponds to the labelled name of the capability in the onchain capability registry.
	CapabilityLabelledName string `toml:"capabilityLabelledName" db:"capability_labelled_name"`

	// OCRKeyBundleIDs is a mapping from chain type to OCR key bundle ID.
	// These are explicitly specified here so that we don't run into strange errors auto-detecting
	// the valid bundle, since nops can create as many bundles as they want.
	// This will most likely never change for a particular CCIP capability version,
	// since new chain families will likely require a new capability version.
	// {"evm": "evm_key_bundle_id", "solana": "solana_key_bundle_id", ... }
	OCRKeyBundleIDs JSONConfig `toml:"ocrKeyBundleIDs" db:"ocr_key_bundle_ids"`

	// RelayConfigs consists of relay specific configuration.
	// Chain reader configurations are stored here, and are defined on a chain family basis, e.g
	// we will have one chain reader config for EVM, one for solana, starknet, etc.
	// Chain writer configurations are also stored here, and are also defined on a chain family basis,
	// e.g we will have one chain writer config for EVM, one for solana, starknet, etc.
	// See tests for examples of relay configs in TOML.
	// { "evm": {"chainReader": {...}, "chainWriter": {...}}, "solana": {...}, ... }
	// see core/services/relay/evm/types/types.go for EVM configs.
	RelayConfigs JSONConfig `toml:"relayConfigs" db:"relay_configs"`

	// P2PKeyID is the ID of the P2P key of the node.
	// This must be present in the capability registry otherwise the job will not start correctly.
	P2PKeyID string `toml:"p2pKeyID" db:"p2p_key_id"`

	// PluginConfig contains plugin-specific config, like token price pipelines
	// and RMN network info for offchain blessing.
	PluginConfig JSONConfig `toml:"pluginConfig"`
}
