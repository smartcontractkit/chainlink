package fees

import (
	"context"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var _ Estimator = &fixedPriceEstimator{}

type fixedPriceEstimator struct {
	price uint64
}

func NewFixedPriceEstimator(cfg config.Config) (Estimator, error) {
	return &fixedPriceEstimator{
		price: cfg.DefaultComputeBudgetPrice(),
	}, nil
}

func (est *fixedPriceEstimator) Start(ctx context.Context) error {
	return nil
}

func (est *fixedPriceEstimator) Close() error {
	return nil
}

func (est *fixedPriceEstimator) GetComputeUnitPrice() (uint64, error) {
	return est.price, nil
}