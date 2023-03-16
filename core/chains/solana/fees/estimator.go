package fees

import "context"

type Estimator interface {
	Start(context.Context) error
	Close() error
	BaseComputeUnitPrice() uint64
}
