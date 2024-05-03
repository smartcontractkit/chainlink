package capabilities

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// RemoteTrigger should be created by generated code, or referenced from a capabilities' public library to assure correct usage with types
type RemoteTrigger[O any] struct {
	RefName  string
	TypeName string
}

func (r RemoteTrigger[O]) Transform(a values.Value) (O, error) {
	return UnwrapValue[O](a)
}

func (r RemoteTrigger[O]) Ref() string {
	return r.RefName
}

func (r RemoteTrigger[O]) Type() string {
	return r.TypeName
}

var _ Trigger[any] = &RemoteTrigger[any]{}
