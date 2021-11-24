package resolver

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

type ChainType string

const (
	ChainTypeArbitrum ChainType = "ARBITRUM"
	ChainTypeExChain  ChainType = "EXCHAIN"
	ChainTypeOptimism ChainType = "OPTIMISM"
	ChainTypeXDAI     ChainType = "XDAI"
)

func ToChainType(s string) (ChainType, error) {
	switch s {
	case "arbitrum":
		return ChainTypeArbitrum, nil
	case "exchain":
		return ChainTypeExChain, nil
	case "optimism":
		return ChainTypeOptimism, nil
	case "xdai":
		return ChainTypeXDAI, nil
	default:
		return "", errors.New("invalid chain type")
	}
}

type GasEstimatorMode string

const (
	GasEstimatorModeBlockHistory GasEstimatorMode = "BLOCK_HISTORY"
	GasEstimatorModeFixedPrice   GasEstimatorMode = "FIXED_PRICE"
	GasEstimatorModeOptimism     GasEstimatorMode = "OPTIMISM"
	GasEstimatorModeOptimism2    GasEstimatorMode = "OPTIMISM2"
)

func ToGasEstimatorMode(s string) (GasEstimatorMode, error) {
	switch s {
	case "BlockHistory":
		return GasEstimatorModeBlockHistory, nil
	case "FixedPrice":
		return GasEstimatorModeFixedPrice, nil
	case "Optimism":
		return GasEstimatorModeOptimism, nil
	case "Optimism2":
		return GasEstimatorModeOptimism2, nil
	default:
		return "", errors.New("invalid gas estimator mode")
	}
}

type ChainConfigResolver struct {
	cfg types.ChainCfg
}

func NewChainConfig(cfg types.ChainCfg) *ChainConfigResolver {
	return &ChainConfigResolver{cfg}
}

type KeySpecificChainConfigResolver struct {
	addr string
	cfg  types.ChainCfg
}

func NewKeySpecificChainConfig(address string, cfg types.ChainCfg) *KeySpecificChainConfigResolver {
	return &KeySpecificChainConfigResolver{
		cfg:  cfg,
		addr: address,
	}
}

func (r *ChainConfigResolver) BlockHistoryEstimatorBlockDelay() *int32 {
	if r.cfg.BlockHistoryEstimatorBlockDelay.Valid {
		val := r.cfg.BlockHistoryEstimatorBlockDelay.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) BlockHistoryEstimatorBlockHistorySize() *int32 {
	if r.cfg.BlockHistoryEstimatorBlockHistorySize.Valid {
		val := r.cfg.BlockHistoryEstimatorBlockHistorySize.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EthTxReaperThreshold() *string {
	if r.cfg.EthTxReaperThreshold != nil {
		threshold := r.cfg.EthTxReaperThreshold.Duration().String()

		return &threshold
	}

	return nil
}

func (r *ChainConfigResolver) EthTxResendAfterThreshold() *string {
	if r.cfg.EthTxResendAfterThreshold != nil {
		threshold := r.cfg.EthTxResendAfterThreshold.Duration().String()

		return &threshold
	}

	return nil
}

func (r *ChainConfigResolver) EvmEIP1559DynamicFees() *bool {
	if r.cfg.EvmEIP1559DynamicFees.Valid {
		return r.cfg.EvmEIP1559DynamicFees.Ptr()
	}

	return nil
}

func (r *ChainConfigResolver) EvmFinalityDepth() *int32 {
	if r.cfg.EvmFinalityDepth.Valid {
		val := r.cfg.EvmFinalityDepth.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasBumpPercent() *int32 {
	if r.cfg.EvmGasBumpPercent.Valid {
		val := r.cfg.EvmGasBumpPercent.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasBumpTxDepth() *int32 {
	if r.cfg.EvmGasBumpTxDepth.Valid {
		val := r.cfg.EvmGasBumpTxDepth.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasBumpWei() *string {
	if r.cfg.EvmGasBumpWei != nil {
		value := r.cfg.EvmGasBumpWei.String()

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasLimitDefault() *int32 {
	if r.cfg.EvmGasLimitDefault.Valid {
		val := r.cfg.EvmGasLimitDefault.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasLimitMultiplier() *float64 {
	if r.cfg.EvmGasLimitMultiplier.Valid {
		return r.cfg.EvmGasLimitMultiplier.Ptr()
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasPriceDefault() *string {
	if r.cfg.EvmGasPriceDefault != nil {
		value := r.cfg.EvmGasPriceDefault.String()

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasTipCapDefault() *string {
	if r.cfg.EvmGasTipCapDefault != nil {
		value := r.cfg.EvmGasTipCapDefault.String()

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) EvmGasTipCapMinimum() *string {
	if r.cfg.EvmGasTipCapMinimum != nil {
		value := r.cfg.EvmGasTipCapMinimum.String()

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) EvmHeadTrackerHistoryDepth() *int32 {
	if r.cfg.EvmHeadTrackerHistoryDepth.Valid {
		val := r.cfg.EvmHeadTrackerHistoryDepth.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmHeadTrackerMaxBufferSize() *int32 {
	if r.cfg.EvmHeadTrackerMaxBufferSize.Valid {
		val := r.cfg.EvmHeadTrackerMaxBufferSize.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmHeadTrackerSamplingInterval() *string {
	if r.cfg.EvmHeadTrackerSamplingInterval != nil {
		interval := r.cfg.EvmHeadTrackerSamplingInterval.Duration().String()

		return &interval
	}

	return nil
}

func (r *ChainConfigResolver) EvmLogBackfillBatchSize() *int32 {
	if r.cfg.EvmLogBackfillBatchSize.Valid {
		val := r.cfg.EvmLogBackfillBatchSize.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) EvmMaxGasPriceWei() *string {
	if r.cfg.EvmMaxGasPriceWei != nil {
		value := r.cfg.EvmMaxGasPriceWei.String()

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) EvmNonceAutoSync() *bool {
	if r.cfg.EvmNonceAutoSync.Valid {
		return r.cfg.EvmNonceAutoSync.Ptr()
	}

	return nil
}

func (r *ChainConfigResolver) EvmRPCDefaultBatchSize() *int32 {
	if r.cfg.EvmRPCDefaultBatchSize.Valid {
		val := r.cfg.EvmRPCDefaultBatchSize.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) FlagsContractAddress() *string {
	if r.cfg.FlagsContractAddress.Valid {
		value := r.cfg.FlagsContractAddress.String

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) GasEstimatorMode() *GasEstimatorMode {
	if r.cfg.GasEstimatorMode.Valid {
		value, err := ToGasEstimatorMode(r.cfg.GasEstimatorMode.String)
		if err != nil {
			return nil
		}

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) ChainType() *ChainType {
	if r.cfg.ChainType.Valid {
		value, err := ToChainType(r.cfg.ChainType.String)
		if err != nil {
			return nil
		}

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) MinIncomingConfirmations() *int32 {
	if r.cfg.MinIncomingConfirmations.Valid {
		val := r.cfg.MinIncomingConfirmations.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) MinRequiredOutgoingConfirmations() *int32 {
	if r.cfg.MinRequiredOutgoingConfirmations.Valid {
		val := r.cfg.MinRequiredOutgoingConfirmations.Int64
		intVal := int32(val)

		return &intVal
	}

	return nil
}

func (r *ChainConfigResolver) MinimumContractPayment() *string {
	if r.cfg.MinimumContractPayment != nil {
		value := r.cfg.MinimumContractPayment.String()

		return &value
	}

	return nil
}

func (r *ChainConfigResolver) OCRObservationTimeout() *string {
	if r.cfg.OCRObservationTimeout != nil {
		timeout := r.cfg.OCRObservationTimeout.Duration().String()

		return &timeout
	}

	return nil
}

func (r *ChainConfigResolver) KeySpecificConfigs() []*KeySpecificChainConfigResolver {
	var resolvers []*KeySpecificChainConfigResolver

	for addr, cfg := range r.cfg.KeySpecific {
		resolvers = append(resolvers, NewKeySpecificChainConfig(addr, cfg))
	}

	return resolvers
}

func (r *KeySpecificChainConfigResolver) Address() string {
	return r.addr
}

func (r *KeySpecificChainConfigResolver) Config() *ChainConfigResolver {
	return NewChainConfig(r.cfg)
}

type ChainConfigInput struct {
	BlockHistoryEstimatorBlockDelay       *int
	BlockHistoryEstimatorBlockHistorySize *int
	EthTxReaperThreshold                  *string
	EthTxResendAfterThreshold             *string
	EvmEIP1559DynamicFees                 *bool
	EvmFinalityDepth                      *int
	EvmGasBumpPercent                     *int
	EvmGasBumpTxDepth                     *int
	EvmGasBumpWei                         *string
	EvmGasLimitDefault                    *int
	EvmGasLimitMultiplier                 *float64
	EvmGasPriceDefault                    *string
	EvmGasTipCapDefault                   *string
	EvmGasTipCapMinimum                   *string
	EvmHeadTrackerHistoryDepth            *int
	EvmHeadTrackerMaxBufferSize           *int
	EvmHeadTrackerSamplingInterval        *string
	EvmLogBackfillBatchSize               *int
	EvmMaxGasPriceWei                     *string
	EvmNonceAutoSync                      *bool
	EvmRPCDefaultBatchSize                *int
	FlagsContractAddress                  *string
	GasEstimatorMode                      string
	ChainType                             string
	MinIncomingConfirmations              *int
	MinRequiredOutgoingConfirmations      *int
	MinimumContractPayment                *string
	OCRObservationTimeout                 *string
}

type KeySpecificChainConfigInput struct {
	Address string
	Config  ChainConfigInput
}
