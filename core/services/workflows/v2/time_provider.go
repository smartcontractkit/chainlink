package v2

import (
	"context"
	"time"
)

// TODO: https://smartcontract-it.atlassian.net/browse/CAPPL-903 implement the methods

type TimeProvider struct{}

func (tp *TimeProvider) GetNodeTime() time.Time {
	return time.Now().UTC()
}

func (tp *TimeProvider) GetDONTime(ctx context.Context) (time.Time, error) {
	return time.Now().UTC(), nil
}
