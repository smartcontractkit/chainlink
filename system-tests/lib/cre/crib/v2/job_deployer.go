package cribv2

import (
	"context"
	"github.com/smartcontractkit/crib-sdk/crib"
)

// Component returns a new Anvil composite component.
func Component(props *crib.Props, opts ...PropOpt) crib.ComponentFunc {
	return func(ctx context.Context) (crib.Component, error) {
		// remote execution here
		// deploy jobs
		return
	}
}
