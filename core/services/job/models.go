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
	POR                     Type = (Type)(pipeline.PORJobType)
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
		POR:                     false,
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
		POR:                     false,
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
		POR:                     1,
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
	PORSpecID                     *int32
	PORSpec                       *PORSpec
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
	P2PV2Bootstrappers                     pq.StringArray         `toml:"p2pv2Bootstrappers" db:"p2v2_bootstrappers"`
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

type PORSpec struct {
	ID        int32     `toml:"-"`
	CreatedAt time.Time `toml:"-"`
	UpdatedAt time.Time `toml:"-"`
}

func (p PORSpec) GetID() string {
	return strconv.Itoa(int(p.ID))
}

func (p *PORSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	p.ID = int32(ID)
	return nil
}

type StandardCapabilitiesSpec struct {
	ID        int32     `toml:"-"`
	CreatedAt time.Time `toml:"-"`
	UpdatedAt time.Time `toml:"-"`
	Command   string    `toml:"command"`
	Config    string    `toml:"config"`
	
	// Oracle Factory configuration for standard capabilities
	OracleFactory StandardCapabilitiesOracleFactory `toml:"oracleFactory"`
}

type StandardCapabilitiesOracleFactory struct {
	Enabled        bool     `toml:"enabled"`
	BootstrapPeers []string `toml:"bootstrapPeers"`
}

func (s StandardCapabilitiesSpec) GetID() string {
	return strconv.Itoa(int(s.ID))
}

func (s *StandardCapabilitiesSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	s.ID = int32(ID)
	return nil
}

type CCIPSpec struct {
	ID                     int32                  `toml:"-"`
	CreatedAt              time.Time              `toml:"-"`
	UpdatedAt              time.Time              `toml:"-"`
	CapabilityVersion      string                 `toml:"capabilityVersion"`
	CapabilityLabelledName string                 `toml:"capabilityLabelledName"`
	OCRKeyBundleIDs        map[string]interface{} `toml:"ocrKeyBundleIDs"`
	P2PKeyID               string                 `toml:"p2pKeyID"`
	P2PV2Bootstrappers     []string               `toml:"p2pV2Bootstrappers"`
	RelayConfigs           map[string]interface{} `toml:"relayConfigs"`
}

func (c CCIPSpec) GetID() string {
	return strconv.Itoa(int(c.ID))
}

func (c *CCIPSpec) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	c.ID = int32(ID)
	return nil
}

// ...existing code...
