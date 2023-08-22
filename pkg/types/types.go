package types

import (
	"context"
)

type Service interface {
	Name() string
	Start(context.Context) error
	Close() error
	Ready() error
	HealthReport() map[string]error
}
