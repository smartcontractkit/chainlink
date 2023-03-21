package presenters

import (
	"github.com/smartcontractkit/chainlink/core/chains"
)

type ChainResource[C chains.Config] struct {
	JAID
	Enabled bool `json:"enabled"`
	Config  C    `json:"config"`
}

func (r ChainResource[C]) GetConfig() any  { return r.Config }
func (r ChainResource[C]) IsEnabled() bool { return r.Enabled }
