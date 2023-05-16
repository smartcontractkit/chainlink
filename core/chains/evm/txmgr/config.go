package txmgr

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// Config encompasses config used by txmgr package
// Unless otherwise specified, these should support changing at runtime
//
//go:generate mockery --quiet --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	gas.Config
	pg.QConfig
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint16
	EvmGasLimitDefault() uint32
	EvmMaxInFlightTransactions() uint32
	EvmMaxQueuedTransactions() uint64
	EvmNonceAutoSync() bool
	EvmUseForwarders() bool
	EvmRPCDefaultBatchSize() uint32
	KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei
	TriggerFallbackDBPollInterval() time.Duration
}

type (
	EvmTxmConfig         txmgrtypes.TxmConfig[*assets.Wei]
	EvmBroadcasterConfig txmgrtypes.BroadcasterConfig[*assets.Wei]
	EvmConfirmerConfig   txmgrtypes.ConfirmerConfig[*assets.Wei]
	EvmResenderConfig    txmgrtypes.ResenderConfig
	EvmReaperConfig      txmgrtypes.ReaperConfig
)

var _ EvmTxmConfig = (*evmTxmConfig)(nil)

type evmTxmConfig struct {
	Config
}

func NewEvmTxmConfig(c Config) *evmTxmConfig {
	return &evmTxmConfig{c}
}

func (c evmTxmConfig) SequenceAutoSync() bool { return c.EvmNonceAutoSync() }

func (c evmTxmConfig) UseForwarders() bool { return c.EvmUseForwarders() }

func (c evmTxmConfig) MaxQueuedTransactions() uint64 { return c.EvmMaxQueuedTransactions() }

func (c evmTxmConfig) MaxInFlightTransactions() uint32 { return c.EvmMaxInFlightTransactions() }

func (c evmTxmConfig) IsL2() bool { return c.ChainType().IsL2() }

func (c evmTxmConfig) MaxFeePrice() *assets.Wei { return c.EvmMaxGasPriceWei() }

func (c evmTxmConfig) FeePriceDefault() *assets.Wei { return c.EvmGasPriceDefault() }

func (c evmTxmConfig) RPCDefaultBatchSize() uint32 { return c.EvmRPCDefaultBatchSize() }

func (c evmTxmConfig) FeeBumpTxDepth() uint16 { return c.EvmGasBumpTxDepth() }

func (c evmTxmConfig) FeeLimitDefault() uint32 { return c.EvmGasLimitDefault() }

func (c evmTxmConfig) FeeBumpThreshold() uint64 { return c.EvmGasBumpThreshold() }

func (c evmTxmConfig) FinalityDepth() uint32 { return c.EvmFinalityDepth() }

func (c evmTxmConfig) FeeBumpPercent() uint16 { return c.EvmGasBumpPercent() }

func (c evmTxmConfig) TxResendAfterThreshold() time.Duration { return c.EthTxResendAfterThreshold() }

func (c evmTxmConfig) TxReaperInterval() time.Duration { return c.EthTxReaperInterval() }

func (c evmTxmConfig) TxReaperThreshold() time.Duration { return c.EthTxReaperThreshold() }
