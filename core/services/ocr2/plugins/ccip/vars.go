package ccip

import (
	"github.com/pkg/errors"
)

const (
	MaxQueryLength       = 0       // empty for both plugins
	MaxObservationLength = 250_000 // plugins's Observation should make sure to cap to this limit
	CommitPluginLabel    = "commit"
	ExecPluginLabel      = "exec"
)

var ErrCommitStoreIsDown = errors.New("commitStoreReader is down")
