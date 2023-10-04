package types

import (
	"context"
)

// Deprecated: use services.Service
type Service interface {
	Name() string
	Start(context.Context) error
	Close() error
	Ready() error
	HealthReport() map[string]error
}
