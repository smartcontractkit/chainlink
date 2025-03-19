package keystore

import (
	"context"
	"fmt"
)

type getDefault[K any] interface {
	EnsureKey(context.Context) error
	GetAll() ([]K, error)
}

func GetDefault[K any, KS getDefault[K]](ctx context.Context, ks KS) (K, error) {
	var zero K
	if err := ks.EnsureKey(ctx); err != nil {
		return zero, fmt.Errorf("failed to ensure %T key", zero)
	}
	keys, err := ks.GetAll()
	if err != nil {
		return zero, err
	}
	if len(keys) < 1 {
		return zero, fmt.Errorf("no %T keys available", zero)
	}
	return keys[0], nil
}
