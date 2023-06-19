package txmgr

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
)

// Config encompasses config used by txmgr package
// Unless otherwise specified, these should support changing at runtime
//
//go:generate mockery --quiet --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	ChainType() coreconfig.ChainType
	EvmEIP1559DynamicFees() bool
	EvmFinalityDepth() uint32
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint32
	EvmGasLimitDefault() uint32
	EvmGasPriceDefault() *assets.Wei
	EvmGasTipCapMinimum() *assets.Wei
	EvmMaxGasPriceWei() *assets.Wei
	EvmMinGasPriceWei() *assets.Wei
	KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei
}

type ChainConfig interface {
	NonceAutoSync() bool
	RPCDefaultBatchSize() uint32
}

type DatabaseConfig interface {
	DefaultQueryTimeout() time.Duration
	LogSQL() bool
}

type ListenerConfig interface {
	FallbackPollInterval() time.Duration
}

type (
	EvmTxmConfig         txmgrtypes.TransactionManagerConfig
	EvmBroadcasterConfig txmgrtypes.BroadcasterConfig
	EvmConfirmerConfig   txmgrtypes.ConfirmerConfig
	EvmResenderConfig    txmgrtypes.ResenderChainConfig
	EvmReaperConfig      txmgrtypes.ReaperConfig
)

var _ EvmTxmConfig = (*evmTxmConfig)(nil)

type evmTxmConfig struct {
	Config
}

func NewEvmTxmConfig(c Config) *evmTxmConfig {
	return &evmTxmConfig{c}
}

func (c evmTxmConfig) IsL2() bool { return c.ChainType().IsL2() }

func (c evmTxmConfig) MaxFeePrice() string { return c.EvmMaxGasPriceWei().String() }

func (c evmTxmConfig) FeePriceDefault() string { return c.EvmGasPriceDefault().String() }

func (c evmTxmConfig) FeeBumpTxDepth() uint32 { return c.EvmGasBumpTxDepth() }

func (c evmTxmConfig) FeeLimitDefault() uint32 { return c.EvmGasLimitDefault() }

func (c evmTxmConfig) FeeBumpThreshold() uint64 { return c.EvmGasBumpThreshold() }

func (c evmTxmConfig) FinalityDepth() uint32 { return c.EvmFinalityDepth() }

func (c evmTxmConfig) FeeBumpPercent() uint16 { return c.EvmGasBumpPercent() }
