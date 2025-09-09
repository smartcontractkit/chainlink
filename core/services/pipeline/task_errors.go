package pipeline

import (
	"errors"
	"fmt"
)

var ErrUnsupportedInLOOPPMode = fmt.Errorf("legacy task not available in LOOP Plugin mode: %w", errors.ErrUnsupported)
