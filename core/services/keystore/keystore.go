// Package keystore manages private keys.
// All raw key byte access is limited to internal/ package APIs, so they are not exposed outside of this subtree.
// Additionally, packages in this subtree may not import the logger packages. Instead, a few select places accept Logf
// funcs to announce new public keys. No other logging is allowed.
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
		return zero, fmt.Errorf("failed to ensure %T key: %w", zero, err)
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
