package internal

import (
	"fmt"
	"math"
	"time"
)

func SafeDurationFromSeconds(s uint) (time.Duration, error) {
	if s > uint(math.MaxInt64/time.Second) {
		return 0, fmt.Errorf("int64 overflow: %d", s)
	}
	return time.Duration(s) * time.Second, nil
}
