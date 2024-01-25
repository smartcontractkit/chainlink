package types

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"gopkg.in/guregu/null.v2"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

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
	ModifierConfigs codec.ModifiersConfig `toml:"modifierConfigs,omitempty"`
}

type ChainContractReader struct {
	ContractABI string `json:"contractABI" toml:"contractABI"`
	// key is genericName from config
	Configs map[string]*ChainReaderDefinition `json:"configs" toml:"configs"`
}

type ChainReaderDefinition chainReaderDefinitionFields

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

	// EventInputFields allows you to choose which indexed fields are expected from the input
	EventInputFields []string `json:"eventInputFields,omitempty"`
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

type RelayConfig struct {
	ChainID                *big.Big           `json:"chainID"`
	FromBlock              uint64             `json:"fromBlock"`
	EffectiveTransmitterID null.String        `json:"effectiveTransmitterID"`
	ConfigContractAddress  *common.Address    `json:"configContractAddress"`
	ChainReader            *ChainReaderConfig `json:"chainReader"`
	Codec                  *CodecConfig       `json:"codec"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID *common.Hash `json:"feedID"`

	// Rebalancer specific
	// FromBlocks specifies the block numbers to replay from for each chain.
	FromBlocks map[string]int64 `json:"fromBlocks"`
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
	//TODO this should be done once and the error should be cached
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
	UpdateRoutes(activeCoordinator common.Address, proposedCoordinator common.Address) error
}

// A LogPoller wrapper that understands router proxy contracts
//
//go:generate mockery --quiet --name LogPollerWrapper --output ./mocks/ --case=underscore
type LogPollerWrapper interface {
	services.Service
	LatestEvents() ([]OracleRequest, []OracleResponse, error)

	// TODO (FUN-668): Remove from the LOOP interface and only use internally within the EVM relayer
	SubscribeToUpdates(name string, subscriber RouteUpdateSubscriber)
}
