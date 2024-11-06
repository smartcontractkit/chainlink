package syncer

import (
	"context"
	"errors"
)

type LocalRegistry any

type ORM[T any] interface {
	// Persists the workflow metadata as rows into workflow_specs and
	// workflow_artifacts tables along with the fetched content
	AddLocalRegistry(ctx context.Context, localRegistry T) error

	// Fetches representation of state.  Most likely some map of DON id to
	// workflow metadata entries
	LatestLocalRegistry(ctx context.Context) (*T, error)
}

type WorkflowRegistryDS = ORM[LocalRegistry]

var _ WorkflowRegistryDS = (*orm)(nil)

type orm struct{}

func NewUnimplementedDS() WorkflowRegistryDS {
	return &orm{}
}

func (o *orm) AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
	return errors.New("not implemented")
}

func (o *orm) LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error) {
	return nil, errors.New("not implemented")
}
