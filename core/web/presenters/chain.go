package presenters

import (
	"time"

	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/chains"
)

type ChainResource[C chains.Config] interface {
	IsEnabled() bool
	GetConfig() C
	jsonapi.EntityNamer
}

// chainResource is a generic chain resource for embedding in a ChainResource implementation.
type chainResource[C chains.Config] struct {
	JAID
	Enabled   bool      `json:"enabled"`
	Config    C         `json:"config"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (r chainResource[C]) GetConfig() C    { return r.Config }
func (r chainResource[C]) IsEnabled() bool { return r.Enabled }
