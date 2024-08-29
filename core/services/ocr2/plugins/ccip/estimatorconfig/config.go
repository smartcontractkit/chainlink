package estimatorconfig

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// FeeEstimatorConfigProvider implements abstract storage for the DataAvailability settings in onRamp dynamic Config.
// It's implemented to transfer DA config from different entities offRamp, onRamp, commitStore without injecting the
// strong dependency between modules. ConfigProvider fetch ccip.OnRampReader object reads and returns only relevant
// fields for the daGasEstimator from the encapsulated onRampReader.
type FeeEstimatorConfigProvider interface {
	SetOnRampReader(reader ccip.OnRampReader)
	GetDataAvailabilityConfig(ctx context.Context) (destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps int64, err error)
}

type FeeEstimatorConfigService struct {
	onRampReader ccip.OnRampReader
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
		return 0, 0, 0, errors.New("no OnRampReader has been configured")
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
