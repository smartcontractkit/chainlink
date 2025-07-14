package utils

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

type LocalTimeProvider struct{}

var _ v2.TimeProvider = &LocalTimeProvider{}

func (t *LocalTimeProvider) GetNodeTime() time.Time {
	return time.Now()
}

func (t *LocalTimeProvider) GetDONTime(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}
