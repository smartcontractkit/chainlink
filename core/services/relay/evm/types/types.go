package types

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"gopkg.in/guregu/null.v2"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ContractReaders map[string]ContractReader
type ChainReaderConfig struct {
	// ContractReaders key is contract name
	ContractReaders `json:"contractReaders"`
}

type CodecConfig struct {
	// ChainCodecConfigs is the type's name for the codec
	ChainCodecConfigs map[string]ChainCodedConfig `json:"chainCodecConfig"`
}

type ChainCodedConfig struct {
	TypeAbi         string `json:"typeAbi"`
	ModifierConfigs codec.ModifiersConfig
}

type ContractReader struct {
	ContractABI       string `json:"contractABI"`
	parsedContractABI *abi.ABI
	// key is genericName from config
	ReadingDefinitions map[string]ReadingDefinition `json:"readingDefinitions"`
}

func (cr ContractReaders) GetReadingDefinition(contractName, readName string) (ReadingDefinition, error) {
	contractReader, ok := cr[contractName]
	if ok {
		return ReadingDefinition{}, fmt.Errorf("chain reading not defined for contract %q", contractName)
	}

	readingDefinition, exists := contractReader.ReadingDefinitions[readName]
	if !exists {
		return ReadingDefinition{}, fmt.Errorf("chain reading not defined for:%s on contract:%s", readName, contractName)
	}

	return readingDefinition, nil
}

func (cr ContractReaders) GetContractABI(contractName string) (abi.ABI, error) {
	contractReader, ok := cr[contractName]
	if ok {
		return abi.ABI{}, fmt.Errorf("chain reading not defined for contract %q", contractName)
	}

	if contractReader.parsedContractABI == nil {
		contractABI, err := abi.JSON(strings.NewReader(contractReader.ContractABI))
		if err != nil {
			return abi.ABI{}, err
		}
		contractReader.parsedContractABI = &contractABI
	}

	return *contractReader.parsedContractABI, nil
}

func (cr ContractReaders) GetEventHash(contractName, eventName string) (common.Hash, error) {
	contractABI, err := cr.GetContractABI(contractName)
	if err != nil {
		return common.Hash{}, err
	}

	readingDefinition, err := cr.GetReadingDefinition(contractName, eventName)
	if err != nil {
		return common.Hash{}, err
	}

	event, ok := contractABI.Events[readingDefinition.ChainSpecificName]
	if !ok {
		return common.Hash{}, fmt.Errorf("event %q hash doesn't exist in contract %q abi", eventName, contractName)
	}

	return event.ID, nil
}

func (cr ContractReaders) GetReadingDefinitionContractAddress(contractName, readName string) (common.Address, error) {
	readingDefinition, err := cr.GetReadingDefinition(contractName, readName)
	if err != nil {
		return common.Address{}, err
	}

	return readingDefinition.ContractAddress, nil
}

func (cr ContractReaders) GetReadType(contractName, readName string) (ReadType, error) {
	contractABI, err := cr.GetContractABI(contractName)
	if err != nil {
		return 0, err
	}
	_, isMethod := contractABI.Methods[readName]
	_, isEvent := contractABI.Events[readName]
	if !isMethod && !isEvent {
		return 0, fmt.Errorf("contract %q doesn't have a method or an event named %q", contractName, readName)
	}
	if isMethod {
		return Method, nil
	}
	return Event, nil
}

type ReadingDefinition struct {
	ChainSpecificName   string `json:"chainSpecificName"` // chain specific contract method name or event type.
	CacheEnabled        bool   `json:"cacheEnabled"`
	ContractAddress     common.Address
	ReadType            `json:"readType"`
	InputModifications  codec.ModifiersConfig `json:"input_modifications"`
	OutputModifications codec.ModifiersConfig `json:"output_modifications"`
}

type ReadType int64

const (
	Method ReadType = 0
	Event  ReadType = 1
)

type RelayConfig struct {
	ChainID                *utils.Big         `json:"chainID"`
	FromBlock              uint64             `json:"fromBlock"`
	EffectiveTransmitterID null.String        `json:"effectiveTransmitterID"`
	ConfigContractAddress  *common.Address    `json:"configContractAddress"`
	ChainReader            *ChainReaderConfig `json:"chainReader"`
	Codec                  *CodecConfig       `json:"codec"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID *common.Hash `json:"feedID"`
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
