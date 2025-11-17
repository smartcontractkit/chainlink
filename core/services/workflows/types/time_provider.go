package types

import (
	"time"
)

type LocalTimeProvider struct{}

func (t *LocalTimeProvider) GetNodeTime() time.Time {
	return time.Now()
}

func (t *LocalTimeProvider) GetDONTime() (time.Time, error) {
	return time.Now(), nil
}
