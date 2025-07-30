package utils

import (
	"context"
	"time"

	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2" //nolint:revive // required alias for v2 module
)

type LocalTimeProvider struct{}

var _ v2.TimeProvider = &LocalTimeProvider{}

func (t *LocalTimeProvider) GetNodeTime() time.Time {
	return time.Now()
}

func (t *LocalTimeProvider) GetDONTime(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}
