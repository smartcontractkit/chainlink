package capabilities

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Trigger[O any] interface {
	Base
	Transform(a values.Value) (O, error)
}
