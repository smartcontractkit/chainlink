package fees

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var _ Estimator = &fixedPriceEstimator{}

type fixedPriceEstimator struct {
	cfg config.Config
}

func NewFixedPriceEstimator(cfg config.Config) (Estimator, error) {
	if cfg.DefaultComputeUnitPrice() < cfg.MinComputeUnitPrice() || cfg.DefaultComputeUnitPrice() > cfg.MaxComputeUnitPrice() {
		return nil, fmt.Errorf("default price (%d) is not within the min (%d) and max (%d) price bounds", cfg.DefaultComputeUnitPrice(), cfg.MinComputeUnitPrice(), cfg.MaxComputeUnitPrice())
	}

	return &fixedPriceEstimator{
		cfg: cfg,
	}, nil
}

func (est *fixedPriceEstimator) Start(ctx context.Context) error {
	return nil
}

func (est *fixedPriceEstimator) Close() error {
	return nil
}

func (est *fixedPriceEstimator) BaseComputeUnitPrice() uint64 {
	return est.cfg.DefaultComputeUnitPrice()
}
