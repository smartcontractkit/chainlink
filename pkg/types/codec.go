package types

import (
	"context"
)

type Encoder interface {
	Encode(ctx context.Context, item any, itemType string) ([]byte, error)
	// GetMaxEncodingSize returns the max size in bytes if n elements are supplied for all top level dynamically sized elements.
	// If no elements are dynamically sized, the returned value will be the same for all n.
	// If there are multiple levels of dynamically sized elements, or itemType cannot be found,
	// ErrInvalidType will be returned.
	GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error)
}

type Decoder interface {
	Decode(ctx context.Context, raw []byte, into any, itemType string) error
	// GetMaxDecodingSize returns the max size in bytes if n elements are supplied for all top level dynamically sized elements.
	// If no elements are dynamically sized, the returned value will be the same for all n.
	// If there are multiple levels of dynamically sized elements, or itemType cannot be found,
	// ErrInvalidType will be returned.
	GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error)
}

type Codec interface {
	Encoder
	Decoder
}

type TypeProvider interface {
	CreateType(itemType string, forEncoding bool) (any, error)
}

type RemoteCodec interface {
	Codec
	TypeProvider
}
