package test

import (
	"context"
)

type KeyValueStore struct {
}

func (t KeyValueStore) Store(ctx context.Context, key string, val []byte) error {
	return nil
}

func (t KeyValueStore) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
