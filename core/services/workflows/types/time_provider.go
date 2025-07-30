package types

import (
	"context"
	"time"
)

type LocalTimeProvider struct{}

func (t *LocalTimeProvider) GetNodeTime() time.Time {
	return time.Now()
}

func (t *LocalTimeProvider) GetDONTime(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}
