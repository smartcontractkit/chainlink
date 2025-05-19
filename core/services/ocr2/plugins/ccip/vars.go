package ccip

import (
	"github.com/pkg/errors"
)

const (
	MaxQueryLength             = 0       // empty for both plugins
	MaxObservationLength       = 250_000 // plugins's Observation should make sure to cap to this limit
	CommitPluginLabel          = "commit"
	ExecPluginLabel            = "exec"
	DefaultSourceFinalityDepth = uint32(2)
	DefaultDestFinalityDepth   = uint32(2)
)

var ErrChainIsNotHealthy = errors.New("lane processing is stopped because of healthcheck failure, please see crit logs")
