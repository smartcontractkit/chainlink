package fees

import "context"

//go:generate mockery --name Estimator --output ./mocks/
type Estimator interface {
	Start(context.Context) error
	Close() error
	BaseComputeUnitPrice() uint64
}
