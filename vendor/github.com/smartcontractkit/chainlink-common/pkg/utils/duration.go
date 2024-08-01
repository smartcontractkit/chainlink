package utils

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Deprecated: use config.Duration
type Duration = config.Duration

// Deprecated: use config.NewDuration
func NewDuration(d time.Duration) (Duration, error) { return config.NewDuration(d) }

// Deprecated: use config.MustNewDuration
func MustNewDuration(d time.Duration) *Duration { return config.MustNewDuration(d) }
