package securemint

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// MultiChainSecureMintServices creates secure mint services that can handle multiple chains
// similar to how CCIP handles multiple chains with separate readers and writers per chain
type MultiChainSecureMintServices struct {
	services.StateMachine

	logger logger.Logger

	// Map of chain selectors to their respective contract readers
	contractReaders map[por.ChainSelector]sm_plugin.ContractReader

	// Map of chain selectors to their respective contract transmitters
	contractTransmitters map[por.ChainSelector]ocr3types.ContractTransmitter[por.ChainSelector]

	// Map of chain selectors to their respective chain writers
	chainWriters map[por.ChainSelector]types.ContractWriter

	// Map of relay IDs to relayers for different chains
	relayers map[types.RelayID]loop.Relayer

	// Configuration for each chain
	chainConfigs map[por.ChainSelector]*ChainConfig

	mu sync.RWMutex
}

// ChainConfig holds configuration specific to a chain
type ChainConfig struct {
	ChainSelector   por.ChainSelector
	RelayID         types.RelayID
	ContractAddress string
	GasLimit        uint64
	MaxGasPrice     string
	FromAddress     string
}

// NewMultiChainSecureMintServices creates a new multi-chain secure mint services instance
func NewMultiChainSecureMintServices(
	logger logger.Logger,
	relayers map[types.RelayID]loop.Relayer,
	chainConfigs map[por.ChainSelector]*ChainConfig,
) *MultiChainSecureMintServices {
	return &MultiChainSecureMintServices{
		logger:               logger.Named("MultiChainSecureMint"),
		contractReaders:      make(map[por.ChainSelector]sm_plugin.ContractReader),
		contractTransmitters: make(map[por.ChainSelector]ocr3types.ContractTransmitter[por.ChainSelector]),
		chainWriters:         make(map[por.ChainSelector]types.ContractWriter),
		relayers:             relayers,
		chainConfigs:         chainConfigs,
	}
}

// Start initializes all chain-specific services
func (m *MultiChainSecureMintServices) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.StateMachine.StartOnce("MultiChainSecureMintServices", func() error {
		return m.initializeChainServices(ctx)
	}); err != nil {
		return err
	}

	return nil
}

// Close shuts down all chain-specific services
func (m *MultiChainSecureMintServices) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.StateMachine.StopOnce("MultiChainSecureMintServices", func() error {
		var errs []error

		// Close all chain writers
		for chainSelector, chainWriter := range m.chainWriters {
			if err := chainWriter.Close(); err != nil {
				m.logger.Errorw("Failed to close chain writer", "chainSelector", chainSelector, "error", err)
				errs = append(errs, fmt.Errorf("failed to close chain writer for chain %d: %w", chainSelector, err))
			}
		}

		// Note: Contract readers don't have a Close method, so we skip them

		if len(errs) > 0 {
			return fmt.Errorf("errors closing services: %v", errs)
		}

		return nil
	})
}

// initializeChainServices creates contract readers, writers, and transmitters for each chain
func (m *MultiChainSecureMintServices) initializeChainServices(ctx context.Context) error {
	for chainSelector, chainConfig := range m.chainConfigs {
		relayer, exists := m.relayers[chainConfig.RelayID]
		if !exists {
			return fmt.Errorf("relayer not found for chain selector %d with relay ID %s", chainSelector, chainConfig.RelayID)
		}

		// Create chain writer for this chain
		chainWriter, err := m.createChainWriter(ctx, relayer, chainConfig)
		if err != nil {
			return fmt.Errorf("failed to create chain writer for chain %d: %w", chainSelector, err)
		}
		m.chainWriters[chainSelector] = chainWriter

		// Create contract reader for this chain
		contractReader, err := m.createContractReader(ctx, relayer, chainConfig)
		if err != nil {
			return fmt.Errorf("failed to create contract reader for chain %d: %w", chainSelector, err)
		}
		m.contractReaders[chainSelector] = contractReader

		// Create contract transmitter for this chain
		contractTransmitter, err := m.createContractTransmitter(ctx, relayer, chainConfig, chainWriter)
		if err != nil {
			return fmt.Errorf("failed to create contract transmitter for chain %d: %w", chainSelector, err)
		}
		m.contractTransmitters[chainSelector] = contractTransmitter

		m.logger.Infow("Initialized services for chain",
			"chainSelector", chainSelector,
			"relayID", chainConfig.RelayID,
			"contractAddress", chainConfig.ContractAddress)
	}

	return nil
}

// createChainWriter creates a chain writer for a specific chain
func (m *MultiChainSecureMintServices) createChainWriter(ctx context.Context, relayer loop.Relayer, chainConfig *ChainConfig) (types.ContractWriter, error) {
	// Create chain writer config similar to CCIP's approach
	chainWriterConfig := map[string]interface{}{
		"contracts": map[string]interface{}{
			"SecureMint": map[string]interface{}{
				"contractABI": `[{"inputs":[],"name":"transmit","outputs":[],"stateMutability":"nonpayable","type":"function"}]`,
				"configs": map[string]interface{}{
					"transmit": map[string]interface{}{
						"chainSpecificName": "transmit",
						"fromAddress":       chainConfig.FromAddress,
						"gasLimit":          chainConfig.GasLimit,
					},
				},
			},
		},
		"maxGasPrice": chainConfig.MaxGasPrice,
	}

	configBytes, err := json.Marshal(chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain writer config: %w", err)
	}

	chainWriter, err := relayer.NewContractWriter(ctx, configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer: %w", err)
	}

	if err := chainWriter.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start chain writer: %w", err)
	}

	return chainWriter, nil
}

// createContractReader creates a contract reader for a specific chain
func (m *MultiChainSecureMintServices) createContractReader(ctx context.Context, relayer loop.Relayer, chainConfig *ChainConfig) (sm_plugin.ContractReader, error) {
	// Create chain reader config
	chainReaderConfig := map[string]interface{}{
		"contracts": map[string]interface{}{
			"SecureMint": map[string]interface{}{
				"contractABI": `[{"inputs":[],"name":"latestTransmittedDetails","outputs":[{"components":[{"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"internalType":"uint64","name":"seqNr","type":"uint64"},{"internalType":"uint64","name":"latestTimestamp","type":"uint64"}],"internalType":"struct TransmittedReportDetails","name":"","type":"tuple"}],"stateMutability":"view","type":"function"}]`,
			},
		},
	}

	configBytes, err := json.Marshal(chainReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	chainReader, err := relayer.NewContractReader(ctx, configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain reader: %w", err)
	}

	// Bind the contract
	err = chainReader.Bind(ctx, []types.BoundContract{
		{
			Address: chainConfig.ContractAddress,
			Name:    "SecureMint",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to bind contract reader: %w", err)
	}

	if err := chainReader.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start chain reader: %w", err)
	}

	// Wrap the chain reader in a secure mint contract reader interface
	return &secureMintContractReader{
		chainReader:   chainReader,
		chainSelector: chainConfig.ChainSelector,
		logger:        m.logger.Named(fmt.Sprintf("ChainReader-%d", chainConfig.ChainSelector)),
	}, nil
}

// createContractTransmitter creates a contract transmitter for a specific chain
func (m *MultiChainSecureMintServices) createContractTransmitter(ctx context.Context, relayer loop.Relayer, chainConfig *ChainConfig, chainWriter types.ContractWriter) (ocr3types.ContractTransmitter[por.ChainSelector], error) {
	// Create a real contract transmitter that uses the chain writer
	return &secureMintContractTransmitter{
		logger:          m.logger.Named(fmt.Sprintf("Transmitter-%d", chainConfig.ChainSelector)),
		chainWriter:     chainWriter,
		fromAccount:     ocrtypes.Account(chainConfig.FromAddress),
		contractAddress: chainConfig.ContractAddress,
		chainSelector:   chainConfig.ChainSelector,
	}, nil
}

// GetContractReader returns the contract reader for a specific chain
func (m *MultiChainSecureMintServices) GetContractReader(chainSelector por.ChainSelector) (sm_plugin.ContractReader, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	reader, exists := m.contractReaders[chainSelector]
	if !exists {
		return nil, fmt.Errorf("contract reader not found for chain selector %d", chainSelector)
	}
	return reader, nil
}

// GetContractTransmitter returns the contract transmitter for a specific chain
func (m *MultiChainSecureMintServices) GetContractTransmitter(chainSelector por.ChainSelector) (ocr3types.ContractTransmitter[por.ChainSelector], error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	transmitter, exists := m.contractTransmitters[chainSelector]
	if !exists {
		return nil, fmt.Errorf("contract transmitter not found for chain selector %d", chainSelector)
	}
	return transmitter, nil
}

// GetChainWriter returns the chain writer for a specific chain
func (m *MultiChainSecureMintServices) GetChainWriter(chainSelector por.ChainSelector) (types.ContractWriter, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	writer, exists := m.chainWriters[chainSelector]
	if !exists {
		return nil, fmt.Errorf("chain writer not found for chain selector %d", chainSelector)
	}
	return writer, nil
}

// ListSupportedChains returns all supported chain selectors
func (m *MultiChainSecureMintServices) ListSupportedChains() []por.ChainSelector {
	m.mu.RLock()
	defer m.mu.RUnlock()

	chains := make([]por.ChainSelector, 0, len(m.chainConfigs))
	for chainSelector := range m.chainConfigs {
		chains = append(chains, chainSelector)
	}
	return chains
}

// HealthReport returns health status for all chains
func (m *MultiChainSecureMintServices) HealthReport() map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	report := make(map[string]error)

	for chainSelector, chainWriter := range m.chainWriters {
		report[fmt.Sprintf("chain_writer_%d", chainSelector)] = chainWriter.HealthReport()[chainWriter.Name()]
	}

	// Note: Contract readers don't implement ServiceCtx, so we skip them

	return report
}

// Name returns the service name
func (m *MultiChainSecureMintServices) Name() string {
	return "MultiChainSecureMintServices"
}
