package common

import (
	"context"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// GetChainReaderParams is a struct that contains the parameters for GetChainReader.
type GetChainReaderParams struct {
	Lggr          logger.Logger
	Relayer       loop.Relayer
	ChainID       string
	DestChainID   string
	HomeChainID   string
	Ofc           OffChainConfig
	ChainSelector cciptypes.ChainSelector
}

// GetChainWriterParams is a struct that contains the parameters for GetChainWriter.
type GetChainWriterParams struct {
	ChainID               string
	Relayer               loop.Relayer
	Transmitters          map[types.RelayID][]string
	ExecBatchGasLimit     uint64
	ChainFamily           string
	OfframpProgramAddress []byte
}

// GetChainReaderWriter is an interface that defines the methods to get a ContractReader and a ContractWriter.
type GetChainReaderWriter interface {
	GetChainReader(ctx context.Context, params GetChainReaderParams) (types.ContractReader, error)
	GetChainWriter(ctx context.Context, params GetChainWriterParams) (types.ContractWriter, error)
}

// PluginConfig is a struct that contains the configuration for a plugin.
type PluginConfig struct {
	CommitPluginCodec   cciptypes.CommitPluginCodec
	ExecutePluginCodec  cciptypes.ExecutePluginCodec
	MessageHasher       func(lggr logger.Logger) cciptypes.MessageHasher
	TokenDataEncoder    cciptypes.TokenDataEncoder
	GasEstimateProvider cciptypes.EstimateProvider
	RMNCrypto           func(lggr logger.Logger) cciptypes.RMNCrypto
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn    string
	GetChainReaderWriter GetChainReaderWriter
}

// OffChainPluginConfig is an interface that defines the method to create a PluginConfig.
type OffChainPluginConfig interface {
	InitializePluginConfig() PluginConfig
}

// PluginConfigFactory is a factory for creating PluginConfig instances.
type PluginConfigFactory struct {
	EVMPluginConfig    OffChainPluginConfig
	SolanaPluginConfig OffChainPluginConfig
}

// NewPluginConfigFactory is a constructor for PluginConfigFactory.
func NewPluginConfigFactory(evmPluginConfig, solanaPluginConfig OffChainPluginConfig) *PluginConfigFactory {
	return &PluginConfigFactory{
		EVMPluginConfig:    evmPluginConfig,
		SolanaPluginConfig: solanaPluginConfig,
	}
}

// CreatePluginConfig creates a PluginConfig instance based on the chain family.
func (f *PluginConfigFactory) CreatePluginConfig(chainFamily string) (PluginConfig, error) {
	switch chainFamily {
	case chainsel.FamilyEVM:
		return f.EVMPluginConfig.InitializePluginConfig(), nil
	case chainsel.FamilySolana:
		return f.SolanaPluginConfig.InitializePluginConfig(), nil
	default:
		return PluginConfig{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
	}
}
