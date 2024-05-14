package capabilities

import (
	"errors"
)

// RemoteTarget should be created by generated code, or *Referenced from a capabilities' public library to assure correct usage with types
type RemoteTarget[I any] struct {
	RefName  string
	TypeName string
}

func (r *RemoteTarget[I]) Ref() string {
	return r.RefName
}

func (r *RemoteTarget[I]) Type() string {
	return r.TypeName
}

func (r *RemoteTarget[I]) Invoke(_ I) error {
	return errors.New("host should be invoking remote targets, this should be a placeholder")
}

var _ Target[any] = &RemoteTarget[any]{}
