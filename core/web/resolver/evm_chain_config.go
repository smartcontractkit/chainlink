package resolver

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

type ChainType string

const (
	ChainTypeArbitrum ChainType = "ARBITRUM"
	ChainTypeOptimism ChainType = "OPTIMISM"
	ChainTypeXDAI     ChainType = "XDAI"
)

func ToChainType(s string) (ChainType, error) {
	switch s {
	case "arbitrum":
		return ChainTypeArbitrum, nil
	case "optimism":
		return ChainTypeOptimism, nil
	case "xdai":
		return ChainTypeXDAI, nil
	default:
		return "", errors.New("invalid chain type")
	}
}

func FromChainType(ct ChainType) string {
	switch ct {
	case ChainTypeArbitrum:
		return "arbitrum"
	case ChainTypeOptimism:
		return "optimism"
	case ChainTypeXDAI:
		return "xdai"
	default:
		return strings.ToLower(string(ct))
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

func FromGasEstimatorMode(gsm GasEstimatorMode) string {
	switch gsm {
	case GasEstimatorModeBlockHistory:
		return "BlockHistory"
	case GasEstimatorModeFixedPrice:
		return "FixedPrice"
	case GasEstimatorModeOptimism:
		return "Optimism"
	case GasEstimatorModeOptimism2:
		return "Optimism2"
	default:
		return strings.ToLower(string(gsm))
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

// LinkContractAddress resolves the LINK contract address for the chain
func (r *ChainConfigResolver) LinkContractAddress() *string {
	if r.cfg.LinkContractAddress.Valid {
		addr := r.cfg.LinkContractAddress.String

		return &addr
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
	BlockHistoryEstimatorBlockDelay       *int32
	BlockHistoryEstimatorBlockHistorySize *int32
	EthTxReaperThreshold                  *string
	EthTxResendAfterThreshold             *string
	EvmEIP1559DynamicFees                 *bool
	EvmFinalityDepth                      *int32
	EvmGasBumpPercent                     *int32
	EvmGasBumpTxDepth                     *int32
	EvmGasBumpWei                         *string
	EvmGasLimitDefault                    *int32
	EvmGasLimitMultiplier                 *float64
	EvmGasPriceDefault                    *string
	EvmGasTipCapDefault                   *string
	EvmGasTipCapMinimum                   *string
	EvmHeadTrackerHistoryDepth            *int32
	EvmHeadTrackerMaxBufferSize           *int32
	EvmHeadTrackerSamplingInterval        *string
	EvmLogBackfillBatchSize               *int32
	EvmMaxGasPriceWei                     *string
	EvmNonceAutoSync                      *bool
	EvmRPCDefaultBatchSize                *int32
	FlagsContractAddress                  *string
	GasEstimatorMode                      *GasEstimatorMode
	ChainType                             *ChainType
	MinIncomingConfirmations              *int32
	MinRequiredOutgoingConfirmations      *int32
	MinimumContractPayment                *string
	OCRObservationTimeout                 *string
	LinkContractAddress                   *string
}

type KeySpecificChainConfigInput struct {
	Address string
	Config  ChainConfigInput
}

func ToChainConfig(input ChainConfigInput) (*types.ChainCfg, map[string]string) {
	cfg := types.ChainCfg{}
	inputErrs := map[string]string{}

	if input.BlockHistoryEstimatorBlockDelay != nil {
		cfg.BlockHistoryEstimatorBlockDelay = null.IntFrom(int64(*input.BlockHistoryEstimatorBlockDelay))
	}

	if input.BlockHistoryEstimatorBlockHistorySize != nil {
		cfg.BlockHistoryEstimatorBlockHistorySize = null.IntFrom(int64(*input.BlockHistoryEstimatorBlockHistorySize))
	}

	if input.EthTxReaperThreshold != nil {
		d, err := models.MakeDurationFromString(*input.EthTxReaperThreshold)
		if err != nil {
			inputErrs["EthTxReaperThreshold"] = "invalid value"
		} else {
			cfg.EthTxReaperThreshold = &d
		}
	}

	if input.EthTxResendAfterThreshold != nil {
		d, err := models.MakeDurationFromString(*input.EthTxResendAfterThreshold)
		if err != nil {
			inputErrs["EthTxResendAfterThreshold"] = "invalid value"
		} else {
			cfg.EthTxResendAfterThreshold = &d
		}
	}

	if input.EvmEIP1559DynamicFees != nil {
		cfg.EvmEIP1559DynamicFees = null.BoolFrom(*input.EvmEIP1559DynamicFees)
	}

	if input.EvmFinalityDepth != nil {
		cfg.EvmFinalityDepth = null.IntFrom(int64(*input.EvmFinalityDepth))
	}

	if input.EvmGasBumpPercent != nil {
		cfg.EvmGasBumpPercent = null.IntFrom(int64(*input.EvmGasBumpPercent))
	}

	if input.EvmGasBumpTxDepth != nil {
		cfg.EvmGasBumpTxDepth = null.IntFrom(int64(*input.EvmGasBumpTxDepth))
	}

	if input.EvmGasBumpWei != nil {
		val, err := stringutils.ToInt64(*input.EvmGasBumpWei)
		if err != nil {
			inputErrs["EvmGasBumpWei"] = "invalid value"
		} else {
			cfg.EvmGasBumpWei = utils.NewBigI(val)
		}
	}

	if input.EvmGasLimitDefault != nil {
		cfg.EvmGasLimitDefault = null.IntFrom(int64(*input.EvmGasLimitDefault))
	}

	if input.EvmGasLimitMultiplier != nil {
		cfg.EvmGasLimitMultiplier = null.FloatFrom(*input.EvmGasLimitMultiplier)
	}

	if input.EvmGasPriceDefault != nil {
		val, err := stringutils.ToInt64(*input.EvmGasPriceDefault)
		if err != nil {
			inputErrs["EvmGasPriceDefault"] = "invalid value"
		} else {
			cfg.EvmGasPriceDefault = utils.NewBigI(val)
		}
	}

	if input.EvmGasTipCapDefault != nil {
		val, err := stringutils.ToInt64(*input.EvmGasTipCapDefault)
		if err != nil {
			inputErrs["EvmGasTipCapDefault"] = "invalid value"
		} else {
			cfg.EvmGasTipCapDefault = utils.NewBigI(val)
		}
	}

	if input.EvmGasTipCapMinimum != nil {
		val, err := stringutils.ToInt64(*input.EvmGasTipCapMinimum)
		if err != nil {
			inputErrs["EvmGasTipCapMinimum"] = "invalid value"
		} else {
			cfg.EvmGasTipCapMinimum = utils.NewBigI(val)
		}
	}

	if input.EvmHeadTrackerHistoryDepth != nil {
		cfg.EvmHeadTrackerHistoryDepth = null.IntFrom(int64(*input.EvmHeadTrackerHistoryDepth))
	}

	if input.EvmHeadTrackerMaxBufferSize != nil {
		cfg.EvmHeadTrackerMaxBufferSize = null.IntFrom(int64(*input.EvmHeadTrackerMaxBufferSize))
	}

	if input.EvmHeadTrackerSamplingInterval != nil {
		d, err := models.MakeDurationFromString(*input.EvmHeadTrackerSamplingInterval)
		if err != nil {
			inputErrs["EvmHeadTrackerSamplingInterval"] = "invalid value"
		} else {
			cfg.EvmHeadTrackerSamplingInterval = &d
		}
	}

	if input.EvmLogBackfillBatchSize != nil {
		cfg.EvmLogBackfillBatchSize = null.IntFrom(int64(*input.EvmLogBackfillBatchSize))
	}

	if input.EvmMaxGasPriceWei != nil {
		val, err := stringutils.ToInt64(*input.EvmMaxGasPriceWei)
		if err != nil {
			inputErrs["EvmMaxGasPriceWei"] = "invalid value"
		} else {
			cfg.EvmMaxGasPriceWei = utils.NewBigI(val)
		}
	}

	if input.EvmNonceAutoSync != nil {
		cfg.EvmNonceAutoSync = null.BoolFrom(*input.EvmNonceAutoSync)
	}

	if input.EvmRPCDefaultBatchSize != nil {
		cfg.EvmRPCDefaultBatchSize = null.IntFrom(int64(*input.EvmRPCDefaultBatchSize))
	}

	if input.FlagsContractAddress != nil {
		cfg.FlagsContractAddress = null.StringFrom(*input.FlagsContractAddress)
	}

	if input.GasEstimatorMode != nil {
		cfg.GasEstimatorMode = null.StringFrom(FromGasEstimatorMode(*input.GasEstimatorMode))
	}

	if input.ChainType != nil {
		cfg.ChainType = null.StringFrom(FromChainType(*input.ChainType))
	}

	if input.MinIncomingConfirmations != nil {
		cfg.MinIncomingConfirmations = null.IntFrom(int64(*input.MinIncomingConfirmations))
	}

	if input.MinRequiredOutgoingConfirmations != nil {
		cfg.MinRequiredOutgoingConfirmations = null.IntFrom(int64(*input.MinRequiredOutgoingConfirmations))
	}

	if input.MinimumContractPayment != nil {
		val, err := stringutils.ToInt64(*input.MinimumContractPayment)
		if err != nil {
			inputErrs["MinimumContractPayment"] = "invalid value"
		} else {
			cfg.MinimumContractPayment = assets.NewLinkFromJuels(val)
		}
	}

	if input.OCRObservationTimeout != nil {
		d, err := models.MakeDurationFromString(*input.OCRObservationTimeout)
		if err != nil {
			inputErrs["MinimumContractPayment"] = "invalid value"
		} else {
			cfg.OCRObservationTimeout = &d
		}
	}

	if input.LinkContractAddress != nil {
		cfg.LinkContractAddress = null.StringFrom(*input.LinkContractAddress)
	}

	return &cfg, inputErrs
}
