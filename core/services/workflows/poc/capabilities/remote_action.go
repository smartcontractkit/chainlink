package capabilities

import (
	"errors"
)

// RemoteAction should be created by generated code, or referenced from a capabilities' public library to assure correct usage with types
type RemoteAction[I, O any] struct {
	RefName  string
	TypeName string
}

func (r *RemoteAction[I, O]) Ref() string {
	return r.RefName
}

func (r *RemoteAction[I, O]) Type() string {
	return r.TypeName
}

func (r *RemoteAction[I, O]) Invoke(_ I) (O, bool, error) {
	var o O
	return o, false, errors.New("host should be invoking remote actions, this should be a placeholder")
}

var _ Action[any, any] = &RemoteAction[any, any]{}
