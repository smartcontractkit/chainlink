package job

import (
	"context"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

//go:generate mockery --quiet --name ServiceCtx --output ./mocks/ --case=underscore

type Service interface {
	Start() error
	Close() error
}

// ServiceCtx is the same as Service, but Start method receives a context.
type ServiceCtx interface {
	Start(context.Context) error
	Close() error
}

type Config interface {
	DatabaseURL() url.URL
	TriggerFallbackDBPollInterval() time.Duration
	pg.QConfig
}

// ServiceAdapter is a helper introduced for transitioning from Service to ServiceCtx.
type ServiceAdapter interface {
	ServiceCtx
}

type adapter struct {
	service Service
}

// NewServiceAdapter creates an adapter instance for the given Service.
func NewServiceAdapter(service Service) ServiceCtx {
	return &adapter{
		service,
	}
}

// Start forwards the call to the underlying service.Start().
// Context is not used in this case.
func (a adapter) Start(context.Context) error {
	return a.service.Start()
}

// Close forwards the call to the underlying service.Close().
func (a adapter) Close() error {
	return a.service.Close()
}
