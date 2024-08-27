package handler

import (
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/config"
)

// BaseHandler is the common handler with a common logic
type BaseHandler struct {
	cfg *config.Config
}

// NewBaseHandler is the constructor of baseHandler
func NewBaseHandler(cfg *config.Config) *BaseHandler {
	return &BaseHandler{
		cfg: cfg,
	}
}
