package estimatorconfig

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// FeeEstimatorConfigProvider implements abstract storage for the DataAvailability settings in onRamp dynamic Config.
// It's implemented to transfer DA config from different entities offRamp, onRamp, commitStore without injecting the
// strong dependency between modules. ConfigProvider fetch ccip.OnRampReader object reads and returns only relevant
// fields for the daGasEstimator from the encapsulated onRampReader.
type FeeEstimatorConfigProvider interface {
	SetOnRampReader(reader ccip.OnRampReader)
	AddGasPriceInterceptor(GasPriceInterceptor)
	ModifyGasPriceComponents(ctx context.Context, execGasPrice, daGasPrice *big.Int) (modExecGasPrice, modDAGasPrice *big.Int, err error)
	GetDataAvailabilityConfig(ctx context.Context) (destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps int64, err error)
}

type GasPriceInterceptor interface {
	ModifyGasPriceComponents(ctx context.Context, execGasPrice, daGasPrice *big.Int) (modExecGasPrice, modDAGasPrice *big.Int, err error)
}

type FeeEstimatorConfigService struct {
	onRampReader         ccip.OnRampReader
	gasPriceInterceptors []GasPriceInterceptor
}

func NewFeeEstimatorConfigService() *FeeEstimatorConfigService {
	return &FeeEstimatorConfigService{}
}

// SetOnRampReader Sets the onRamp reader instance.
// must be called once for each instance.
func (c *FeeEstimatorConfigService) SetOnRampReader(reader ccip.OnRampReader) {
	c.onRampReader = reader
}

// GetDataAvailabilityConfig Returns dynamic config data availability parameters.
// GetDynamicConfig should be cached in the onRamp reader to avoid unnecessary on-chain calls
func (c *FeeEstimatorConfigService) GetDataAvailabilityConfig(ctx context.Context) (destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps int64, err error) {
	if c.onRampReader == nil {
		return 0, 0, 0, nil
	}

	cfg, err := c.onRampReader.GetDynamicConfig(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	return int64(cfg.DestDataAvailabilityOverheadGas),
		int64(cfg.DestGasPerDataAvailabilityByte),
		int64(cfg.DestDataAvailabilityMultiplierBps),
		err
}

// AddGasPriceInterceptor adds price interceptors that can modify gas price.
func (c *FeeEstimatorConfigService) AddGasPriceInterceptor(gpi GasPriceInterceptor) {
	if gpi != nil {
		c.gasPriceInterceptors = append(c.gasPriceInterceptors, gpi)
	}
}

// ModifyGasPriceComponents applies gasPrice interceptors and returns modified gasPrice.
func (c *FeeEstimatorConfigService) ModifyGasPriceComponents(ctx context.Context, gasPrice, daGasPrice *big.Int) (*big.Int, *big.Int, error) {
	if len(c.gasPriceInterceptors) == 0 {
		return gasPrice, daGasPrice, nil
	}

	// values are mutable, it is necessary to copy the values to protect the arguments from modification.
	cpGasPrice := new(big.Int).Set(gasPrice)
	cpDAGasPrice := new(big.Int).Set(daGasPrice)

	var err error
	for _, interceptor := range c.gasPriceInterceptors {
		if cpGasPrice, cpDAGasPrice, err = interceptor.ModifyGasPriceComponents(ctx, cpGasPrice, cpDAGasPrice); err != nil {
			return nil, nil, err
		}
	}

	return cpGasPrice, cpDAGasPrice, nil
}
