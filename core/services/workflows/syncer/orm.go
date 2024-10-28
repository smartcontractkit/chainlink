package syncer

import (
	"context"
	"errors"
)

// TODO(mstreet3): Define the methods on the Local Registry
type LocalRegistry any

// TODO(mstreet3): Define a type T for the ORM
type ORM[T any] interface {
	// Persists the workflow metadata as rows into workflow_specs and
	// workflow_artifacts tables along with the fetched content
	//
	// TODO(mstreet3): Should we limit the number of rows of state kept?
	// 				   Should we keep raw state and hash separately?
	AddLocalRegistry(ctx context.Context, localRegistry T) error

	// Fetches representation of state.  Most likely some map of DON id to
	// workflow metadata entries
	LatestLocalRegistry(ctx context.Context) (*T, error)
}

type WorkflowRegistryDS = ORM[LocalRegistry]

var _ WorkflowRegistryDS = (*orm)(nil)

func NewUnimplementedDS() *orm {
	return &orm{}
}

type orm struct{}

func (o *orm) AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
	return errors.New("not implemented")
}

func (o *orm) LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error) {
	return nil, errors.New("not implemented")
}
