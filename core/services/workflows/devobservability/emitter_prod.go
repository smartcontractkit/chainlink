//go:build !dev

package devobservability

import (
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// In production, no wrapping - just return the original
func WrapEmitter(underlying beholder.Emitter) beholder.Emitter {
	return underlying
}
