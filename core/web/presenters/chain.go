package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/chains"
)

// chainResource is a generic chain resource for embedding in a typed EntityNamer.
type chainResource[C chains.Config] struct {
	JAID
	Enabled   bool      `json:"enabled"`
	Config    C         `json:"config"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
