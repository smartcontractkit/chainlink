package core

import "context"

type KeyValueStore interface {
	Store(ctx context.Context, key string, val []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}
