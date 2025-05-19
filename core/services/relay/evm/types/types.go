package types

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type ChainWriterConfig struct {
	Contracts   map[string]*ContractConfig
	MaxGasPrice *assets.Wei
}

type ContractConfig struct {
	ContractABI string `json:"contractABI" toml:"contractABI"`
	// key is genericName from config
	Configs map[string]*ChainWriterDefinition `json:"configs" toml:"configs"`
}

type ChainWriterDefinition struct {
	// chain specific contract method name or event type.
	ChainSpecificName  string                `json:"chainSpecificName"`
	Checker            string                `json:"checker"`
	FromAddress        common.Address        `json:"fromAddress"`
	GasLimit           uint64                `json:"gasLimit"` // TODO(archseer): what if this has to be configured per call?
	InputModifications codec.ModifiersConfig `json:"inputModifications,omitempty"`
}

type ChainReaderConfig struct {
	// Contracts key is contract name
	Contracts map[string]ChainContractReader `json:"contracts" toml:"contracts"`
}

type CodecConfig struct {
	// Configs key is the type's name for the codec
	Configs map[string]ChainCodecConfig `json:"configs" toml:"configs"`
}

type ChainCodecConfig struct {
	TypeABI         string                `json:"typeAbi" toml:"typeABI"`
	ModifierConfigs codec.ModifiersConfig `json:"modifierConfigs,omitempty" toml:"modifierConfigs,omitempty"`
}

type DataWordDetail struct {
	Name string `json:"name"`
	// Index is indexed from 0. Index should only be used as an override in specific edge case scenarios where the index can't be programmatically calculated, otherwise leave this as nil.
	Index *int `json:"index,omitempty"`
	// Type should follow the geth ABI types naming convention
	Type string `json:"type,omitempty"`
}

type ContractPollingFilter struct {
	GenericEventNames []string `json:"genericEventNames"`
	PollingFilter     `json:"pollingFilter"`
}

type PollingFilter struct {
	Topic2       evmtypes.HashArray `json:"topic2"`       // list of possible values for topic2
	Topic3       evmtypes.HashArray `json:"topic3"`       // list of possible values for topic3
	Topic4       evmtypes.HashArray `json:"topic4"`       // list of possible values for topic4
	Retention    models.Interval    `json:"retention"`    // maximum amount of time to retain logs
	MaxLogsKept  uint64             `json:"maxLogsKept"`  // maximum number of logs to retain ( 0 = unlimited )
	LogsPerBlock uint64             `json:"logsPerBlock"` // rate limit ( maximum # of logs per block, 0 = unlimited )
}

func (f *PollingFilter) ToLPFilter(eventSigs evmtypes.HashArray) logpoller.Filter {
	return logpoller.Filter{
		EventSigs:    eventSigs,
		Topic2:       f.Topic2,
		Topic3:       f.Topic3,
		Topic4:       f.Topic4,
		Retention:    f.Retention.Duration(),
		MaxLogsKept:  f.MaxLogsKept,
		LogsPerBlock: f.LogsPerBlock,
	}
}

type ChainContractReader struct {
	ContractABI           string `json:"contractABI" toml:"contractABI"`
	ContractPollingFilter `json:"contractPollingFilter,omitempty" toml:"contractPollingFilter,omitempty"`
	// key is genericName from config
	Configs map[string]*ChainReaderDefinition `json:"configs" toml:"configs"`
}

type ChainReaderDefinition chainReaderDefinitionFields

type EventDefinitions struct {
	// GenericTopicNames helps QueryingKeys not rely on EVM specific topic names. Key is chain specific name, value is generic name.
	// This helps us translate chain agnostic querying key "transfer-value" to EVM specific "evmTransferEvent-weiAmountTopic".
	GenericTopicNames map[string]string `json:"genericTopicNames,omitempty"`
	// GenericDataWordDetails key is generic name for evm log event data word that maps to chain details.
	// For e.g. first evm data word(32bytes) of USDC log event is value so the key can be called value.
	GenericDataWordDetails map[string]DataWordDetail `json:"genericDataWordDetails,omitempty"`
	// PollingFilter should be defined on a contract level in ContractPollingFilter,
	// unless event needs to override the contract level filter options.
	// This will create a separate log poller filter for this event.
	*PollingFilter `json:"pollingFilter,omitempty"`
}

// chainReaderDefinitionFields has the fields for ChainReaderDefinition but no methods.
// This is necessary because package json recognizes the text encoding methods used for TOML,
// and would infinitely recurse on itself.
type chainReaderDefinitionFields struct {
	CacheEnabled bool `json:"cacheEnabled,omitempty"`
	// chain specific contract method name or event type.
	ChainSpecificName   string                `json:"chainSpecificName"`
	ReadType            ReadType              `json:"readType,omitempty"`
	InputModifications  codec.ModifiersConfig `json:"inputModifications,omitempty"`
	OutputModifications codec.ModifiersConfig `json:"outputModifications,omitempty"`
	EventDefinitions    *EventDefinitions     `json:"eventDefinitions,omitempty" toml:"eventDefinitions,omitempty"`
	// ConfidenceConfirmations is a mapping between a ConfidenceLevel and the confirmations associated. Confidence levels
	// should be valid float values.
	ConfidenceConfirmations map[string]int `json:"confidenceConfirmations,omitempty"`
}

func (d *ChainReaderDefinition) HasPollingFilter() bool {
	return d.EventDefinitions != nil && d.EventDefinitions.PollingFilter != nil
}

func (d *ChainReaderDefinition) MarshalText() ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetIndent("", "  ")
	if err := e.Encode((*chainReaderDefinitionFields)(d)); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (d *ChainReaderDefinition) UnmarshalText(b []byte) error {
	return json.Unmarshal(b, (*chainReaderDefinitionFields)(d))
}

type ReadType int

const (
	Method ReadType = iota
	Event
)

func (r ReadType) String() string {
	switch r {
	case Method:
		return "method"
	case Event:
		return "event"
	}
	return fmt.Sprintf("ReadType(%d)", r)
}

func (r ReadType) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *ReadType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "method":
		*r = Method
		return nil
	case "event":
		*r = Event
		return nil
	}
	return fmt.Errorf("unrecognized ReadType: %s", string(text))
}

type LLOConfigMode string

const (
	LLOConfigModeMercury   LLOConfigMode = "mercury"
	LLOConfigModeBlueGreen LLOConfigMode = "bluegreen"
)

func (c LLOConfigMode) String() string {
	return string(c)
}

type DualTransmissionConfig struct {
	ContractAddress    common.Address      `json:"contractAddress" toml:"contractAddress"`
	TransmitterAddress common.Address      `json:"transmitterAddress" toml:"transmitterAddress"`
	Meta               map[string][]string `json:"meta" toml:"meta"`
}

type RelayConfig struct {
	ChainID                *big.Big           `json:"chainID"`
	FromBlock              uint64             `json:"fromBlock"`
	EffectiveTransmitterID null.String        `json:"effectiveTransmitterID"`
	ConfigContractAddress  *common.Address    `json:"configContractAddress"`
	ChainReader            *ChainReaderConfig `json:"chainReader"`
	Codec                  *CodecConfig       `json:"codec"`

	DefaultTransactionQueueDepth uint32 `json:"defaultTransactionQueueDepth"`
	SimulateTransactions         bool   `json:"simulateTransactions"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID                   *common.Hash `json:"feedID"`
	EnableTriggerCapability  bool         `json:"enableTriggerCapability"`
	TriggerCapabilityName    string       `json:"triggerCapabilityName"`
	TriggerCapabilityVersion string       `json:"triggerCapabilityVersion"`

	// LLO-specific
	LLODONID      uint32        `json:"lloDonID" toml:"lloDonID"`
	LLOConfigMode LLOConfigMode `json:"lloConfigMode" toml:"lloConfigMode"`

	// DualTransmission specific
	EnableDualTransmission bool                    `json:"enableDualTransmission" toml:"enableDualTransmission"`
	DualTransmissionConfig *DualTransmissionConfig `json:"dualTransmission" toml:"dualTransmission"`
}

var ErrBadRelayConfig = errors.New("bad relay config")

type RelayOpts struct {
	// TODO BCF-2508 -- should anyone ever get the raw config bytes that are embedded in args? if not,
	// make this private and wrap the arg fields with funcs on RelayOpts
	types.RelayArgs
	c *RelayConfig
}

func NewRelayOpts(args types.RelayArgs) *RelayOpts {
	return &RelayOpts{
		RelayArgs: args,
		c:         nil, // lazy initialization
	}
}

func (o *RelayOpts) RelayConfig() (RelayConfig, error) {
	var empty RelayConfig
	// TODO this should be done once and the error should be cached
	if o.c == nil {
		var c RelayConfig
		err := json.Unmarshal(o.RelayArgs.RelayConfig, &c)
		if err != nil {
			return empty, fmt.Errorf("%w: failed to deserialize relay config: %w", ErrBadRelayConfig, err)
		}
		o.c = &c
	}
	return *o.c, nil
}

type ConfigPoller interface {
	ocrtypes.ContractConfigTracker

	Start()
	Close() error
	Replay(ctx context.Context, fromBlock int64) error
}

// TODO(FUN-668): Migrate this fully into types.FunctionsProvider
type FunctionsProvider interface {
	types.FunctionsProvider
	LogPollerWrapper() LogPollerWrapper
}

type OracleRequest struct {
	RequestId           [32]byte
	RequestingContract  common.Address
	RequestInitiator    common.Address
	SubscriptionId      uint64
	SubscriptionOwner   common.Address
	Data                []byte
	DataVersion         uint16
	Flags               [32]byte
	CallbackGasLimit    uint64
	TxHash              common.Hash
	CoordinatorContract common.Address
	OnchainMetadata     []byte
}

type OracleResponse struct {
	RequestId [32]byte
}

type RouteUpdateSubscriber interface {
	UpdateRoutes(ctx context.Context, activeCoordinator common.Address, proposedCoordinator common.Address) error
}

// A LogPoller wrapper that understands router proxy contracts
type LogPollerWrapper interface {
	services.Service
	LatestEvents(ctx context.Context) ([]OracleRequest, []OracleResponse, error)

	// TODO (FUN-668): Remove from the LOOP interface and only use internally within the EVM relayer
	SubscribeToUpdates(ctx context.Context, name string, subscriber RouteUpdateSubscriber)
}

// ChainReaderConfigFromBytes function applies json decoding on provided bytes to return a
// valid ChainReaderConfig. If other unmarshaling or parsing techniques are required, place it
// here.
func ChainReaderConfigFromBytes(bts []byte) (ChainReaderConfig, error) {
	decoder := json.NewDecoder(bytes.NewBuffer(bts))

	decoder.UseNumber()

	var cfg ChainReaderConfig
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal chain reader config err: %s", err)
	}

	return cfg, nil
}
