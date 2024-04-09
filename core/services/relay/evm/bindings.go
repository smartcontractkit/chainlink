package evm

import (
	"context"
	"fmt"
	"strings"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// key is contractName.readName
type contractBindings map[string]readBinding

func (b contractBindings) GetReadBinding(key string) (readBinding, error) {
	tokens := strings.Split(key, ".")
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid key: %s, key should look like: contractName.readName", key)
	}

	rb, rbExists := b[key]
	if !rbExists {
		return nil, fmt.Errorf("%w: no readbinding by key: %s", commontypes.ErrInvalidType, key)
	}

	return rb, nil
}

func (b contractBindings) AddReadBinding(key string, reader readBinding) {
	_, rbsExists := b[key]
	if !rbsExists {
		return
	}
	b[key] = reader
}

func (b contractBindings) Bind(ctx context.Context, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		rb, rbsExist := b[bc.Name]
		if !rbsExist {
			return fmt.Errorf("%w: no contract named %s", commontypes.ErrInvalidConfig, bc.Name)
		}

		if err := rb.Bind(ctx, bc); err != nil {
			return err
		}
	}
	return nil
}

func (b contractBindings) ForEach(ctx context.Context, fn func(readBinding, context.Context) error) error {
	for _, rb := range b {
		if err := fn(rb, ctx); err != nil {
			return err
		}
	}
	return nil
}

func formatKey(str ...string) string {
	return strings.Join(str, ".")
}
