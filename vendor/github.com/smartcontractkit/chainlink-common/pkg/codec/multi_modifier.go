package codec

import (
	"reflect"
)

// MultiModifier is a Modifier that applies each element for the slice in-order (reverse order for TransformForOnChain).
type MultiModifier []Modifier

func (c MultiModifier) RetypeToOffChain(onChainType reflect.Type, itemType string) (reflect.Type, error) {
	return forEach(c, onChainType, itemType, Modifier.RetypeToOffChain)
}

func (c MultiModifier) TransformToOnChain(offChainValue any, itemType string) (any, error) {
	onChainValue := offChainValue
	for i := len(c) - 1; i >= 0; i-- {
		var err error
		if onChainValue, err = c[i].TransformToOnChain(onChainValue, itemType); err != nil {
			return nil, err
		}
	}

	return onChainValue, nil
}

func (c MultiModifier) TransformToOffChain(onChainValue any, itemType string) (any, error) {
	return forEach(c, onChainValue, itemType, Modifier.TransformToOffChain)
}

func forEach[T any](c MultiModifier, input T, itemType string, fn func(Modifier, T, string) (T, error)) (T, error) {
	output := input
	for _, m := range c {
		var err error
		if output, err = fn(m, output, itemType); err != nil {
			return output, err
		}
	}
	return output, nil
}
