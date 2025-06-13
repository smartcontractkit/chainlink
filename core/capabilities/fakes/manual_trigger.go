package fakes

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type ManualTriggerCapability interface {
	services.Service
	ManualTrigger(ctx context.Context) error
}
