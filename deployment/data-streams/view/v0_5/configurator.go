package v0_5

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
)

// ConfigState represents a single configuration state
type ConfigState struct {
	ConfigCount           uint64   `json:"configCount"`
	EncodedOffchainConfig string   `json:"encodedOffchainConfig"`
	EncodedOnchainConfig  string   `json:"encodedOnchainConfig"`
	F                     uint8    `json:"f"`
	OffchainConfigVersion uint64   `json:"offchainConfigVersion"`
	OffchainTransmitters  []string `json:"offchainTransmitters"`
	Signers               []string `json:"signers"`
	IsGreenProduction     bool     `json:"isGreenProduction"`
	BlockNumber           uint32   `json:"blockNumber"`
}

// ConfiguratorView represents a simplified view of the contract state
type ConfiguratorView struct {
	// Maps: configId -> configDigest -> ConfigState
	Configs        map[string]map[string]*ConfigState `json:"configs"`
	TypeAndVersion string                             `json:"typeAndVersion,omitempty"`
	Owner          common.Address                     `json:"owner,omitempty"`
}

// ConfiguratorView implements the ContractView interface
var _ interfaces.ContractView = (*ConfiguratorView)(nil)

// NewContractView creates a new empty ConfiguratorView
func NewContractView() *ConfiguratorView {
	return &ConfiguratorView{
		Configs: make(map[string]map[string]*ConfigState),
	}
}

type ConfiguratorViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// ConfiguratorContract defines the minimal interface needed for ConfiguratorViewGenerator
type ConfiguratorContract interface {
	// Call methods
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Owner(opts *bind.CallOpts) (common.Address, error)

	// Filter methods
	FilterProductionConfigSet(opts *bind.FilterOpts, configID [][32]byte) (evm.LogIterator[configurator.ConfiguratorProductionConfigSet], error)
	FilterStagingConfigSet(opts *bind.FilterOpts, configID [][32]byte) (evm.LogIterator[configurator.ConfiguratorStagingConfigSet], error)
	FilterPromoteStagingConfig(opts *bind.FilterOpts, configID [][32]byte, retiredConfigDigest [][32]byte) (evm.LogIterator[configurator.ConfiguratorPromoteStagingConfig], error)
}

type ConfiguratorViewGenerator struct {
	contract ConfiguratorContract
}

func NewConfiguratorViewGenerator(contract ConfiguratorContract) *ConfiguratorViewGenerator {
	return &ConfiguratorViewGenerator{
		contract: contract,
	}
}

// ConfiguratorViewGenerator implements ContractViewGenerator
var _ interfaces.ContractViewGenerator[ConfiguratorViewParams, ConfiguratorView] = (*ConfiguratorViewGenerator)(nil)

// Generate scans builds a view of the contract state from logs and calls
func (b *ConfiguratorViewGenerator) Generate(ctx context.Context, chainParams ConfiguratorViewParams) (ConfiguratorView, error) {
	view := NewContractView()

	// Define the filter options
	filterOpts := &bind.FilterOpts{
		Start:   chainParams.FromBlock,
		End:     chainParams.ToBlock,
		Context: ctx,
	}

	// Process both production and staging config events
	if err := b.processConfigEvents(filterOpts, view); err != nil {
		return ConfiguratorView{}, err
	}

	// Process all promote staging config events
	if err := b.processPromotions(filterOpts, view); err != nil {
		return ConfiguratorView{}, err
	}

	owner, err := b.contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return ConfiguratorView{}, fmt.Errorf("failed to get contract owner: %w", err)
	}

	view.Owner = owner

	return *view, nil
}

// SerializeView serializes the ConfiguratorView to JSON
func (v ConfiguratorView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// processConfigEvents processes both ProductionConfigSet and StagingConfigSet events
func (b *ConfiguratorViewGenerator) processConfigEvents(opts *bind.FilterOpts, view *ConfiguratorView) error {
	// Process production configs
	prodIter, err := b.contract.FilterProductionConfigSet(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter production config events: %w", err)
	}
	defer prodIter.Close()

	for prodIter.Next() {
		event := prodIter.GetEvent()
		b.processConfigEvent(event.ConfigId, event.ConfigDigest, event, view)
	}

	// Process staging configs
	stagingIter, err := b.contract.FilterStagingConfigSet(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter staging config events: %w", err)
	}
	defer stagingIter.Close()

	for stagingIter.Next() {
		event := stagingIter.GetEvent()
		b.processConfigEvent(event.ConfigId, event.ConfigDigest, event, view)
	}

	return nil
}

// processConfigEvent processes either a ProductionConfigSet or StagingConfigSet event
func (b *ConfiguratorViewGenerator) processConfigEvent(configID [32]byte, configDigest [32]byte, event interface{}, view *ConfiguratorView) {
	var (
		configCount           uint64
		offchainConfig        []byte
		onchainConfig         []byte
		f                     uint8
		offchainConfigVersion uint64
		offchainTransmitters  [][32]byte
		signers               [][]byte
		isGreenProduction     bool
		blockNumber           uint32
	)

	switch e := event.(type) {
	case *configurator.ConfiguratorProductionConfigSet:
		configCount = e.ConfigCount
		offchainConfig = e.OffchainConfig
		onchainConfig = e.OnchainConfig
		f = e.F
		offchainConfigVersion = e.OffchainConfigVersion
		offchainTransmitters = e.OffchainTransmitters
		signers = e.Signers
		isGreenProduction = e.IsGreenProduction
		blockNumber = e.PreviousConfigBlockNumber
	case *configurator.ConfiguratorStagingConfigSet:
		configCount = e.ConfigCount
		offchainConfig = e.OffchainConfig
		onchainConfig = e.OnchainConfig
		f = e.F
		offchainConfigVersion = e.OffchainConfigVersion
		offchainTransmitters = e.OffchainTransmitters
		signers = e.Signers
		isGreenProduction = e.IsGreenProduction
		blockNumber = e.PreviousConfigBlockNumber
	default:
		panic("unknown event type")
	}

	configIDHex := dsutil.HexEncodeBytes32(configID)
	configDigestHex := dsutil.HexEncodeBytes32(configDigest)

	if _, ok := view.Configs[configIDHex]; !ok {
		view.Configs[configIDHex] = make(map[string]*ConfigState)
	}

	// Convert types to readable hex strings
	signersHex := make([]string, len(signers))
	for i, signer := range signers {
		signersHex[i] = dsutil.HexEncodeBytes(signer)
	}

	// Convert types to readable hex strings
	transmittersHex := make([]string, len(offchainTransmitters))
	for i, transmitter := range offchainTransmitters {
		transmittersHex[i] = dsutil.HexEncodeBytes32(transmitter)
	}

	view.Configs[configIDHex][configDigestHex] = &ConfigState{
		ConfigCount:           configCount,
		EncodedOffchainConfig: dsutil.HexEncodeBytes(offchainConfig),
		EncodedOnchainConfig:  dsutil.HexEncodeBytes(onchainConfig),
		F:                     f,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainTransmitters:  transmittersHex,
		Signers:               signersHex,
		IsGreenProduction:     isGreenProduction,
		BlockNumber:           blockNumber,
	}
}

// processPromotions processes all PromoteStagingConfig events
func (b *ConfiguratorViewGenerator) processPromotions(opts *bind.FilterOpts, view *ConfiguratorView) error {
	iter, err := b.contract.FilterPromoteStagingConfig(opts, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to filter promote staging config events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.GetEvent()
		configIDHex := dsutil.HexEncodeBytes32(event.ConfigId)

		// Skip if this configId doesn't exist
		if configs, ok := view.Configs[configIDHex]; ok {
			// Update isGreenProduction for all configs with this configId
			for digest, config := range configs {
				config.IsGreenProduction = event.IsGreenProduction
				configs[digest] = config
			}
		}
	}

	return nil
}

func (v ConfiguratorView) GetConfigState(configID, configDigest string) (*ConfigState, error) {
	configs, ok := v.Configs[configID]
	if !ok {
		return nil, fmt.Errorf("configID %s not found", configID)
	}
	configState, ok := configs[configDigest]
	if !ok {
		return nil, fmt.Errorf("configDigest %s not found for configID %s", configDigest, configID)
	}
	return configState, nil
}
