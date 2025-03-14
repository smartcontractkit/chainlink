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

type GetChainReaderParams struct {
	Ctx           context.Context
	Lggr          logger.Logger
	Relayer       loop.Relayer
	ChainID       string
	DestChainID   string
	HomeChainID   string
	Ofc           OffChainConfig
	ChainSelector cciptypes.ChainSelector
}

type GetChainWriterParams struct {
	Ctx                   context.Context
	ChainID               string
	Relayer               loop.Relayer
	Transmitters          map[types.RelayID][]string
	ExecBatchGasLimit     uint64
	ChainFamily           string
	OfframpProgramAddress []byte
}

type GetChainReaderWriter interface {
	GetChainReader(params GetChainReaderParams) (types.ContractReader, error)
	GetChainWriter(params GetChainWriterParams) (types.ContractWriter, error)
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
