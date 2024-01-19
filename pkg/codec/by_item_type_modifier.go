package codec

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// NewByItemTypeModifier returns a Modifier that uses modByItemType to determine which Modifier to use for a given itemType.
func NewByItemTypeModifier(modByItemType map[string]Modifier) (Modifier, error) {
	if modByItemType == nil {
		modByItemType = map[string]Modifier{}
	}

	return &byItemTypeModifier{
		modByitemType: modByItemType,
	}, nil
}

type byItemTypeModifier struct {
	modByitemType map[string]Modifier
}

func (b *byItemTypeModifier) RetypeToOffChain(onChainType reflect.Type, itemType string) (reflect.Type, error) {
	mod, ok := b.modByitemType[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find modifier for %s", types.ErrInvalidType, itemType)
	}

	return mod.RetypeToOffChain(onChainType, itemType)
}

func (b *byItemTypeModifier) TransformToOnChain(offChainValue any, itemType string) (any, error) {
	return b.transform(offChainValue, itemType, Modifier.TransformToOnChain)
}

func (b *byItemTypeModifier) TransformToOffChain(onChainValue any, itemType string) (any, error) {
	return b.transform(onChainValue, itemType, Modifier.TransformToOffChain)
}

func (b *byItemTypeModifier) transform(
	val any, itemType string, transform func(Modifier, any, string) (any, error)) (any, error) {
	mod, ok := b.modByitemType[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find modifier for %s", types.ErrInvalidType, itemType)
	}

	return transform(mod, val, itemType)
}

var _ Modifier = &byItemTypeModifier{}
