package evm

import (
	"context"
	"fmt"
	"strings"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// key is contract-readName
type contractBindings map[string]readBinding

func (b contractBindings) GetReadBinding(key string) (readBinding, error) {
	rb, rbExists := b[key]
	if !rbExists {
		return nil, fmt.Errorf("%w: no readbinding by key %s", commontypes.ErrInvalidType, key)
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

		address, err := validateEthereumAddress(bc.Address)
		if err != nil {
			return err
		}

		if err = rb.Bind(ctx, address); err != nil {
			return err
		}
	}
	return nil
}

func (b contractBindings) UnBind(ctx context.Context, boundContracts []commontypes.BoundContract) error {
	for _, bc := range boundContracts {
		rb, rbsExist := b[bc.Name]
		if rbsExist {
			address, err := validateEthereumAddress(bc.Address)
			if err != nil {
				return err
			}
			if err := rb.UnBind(ctx, address); err != nil {
				return err
			}
			delete(b, bc.Name)
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
