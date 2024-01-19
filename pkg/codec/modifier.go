package codec

import (
	"reflect"
)

// Modifier allows you to modify the off-chain type to be used on-chain, and vice-versa.
// A modifier is set up by retyping the on-chain type to a type used off-chain.
type Modifier interface {
	RetypeToOffChain(onChainType reflect.Type, itemType string) (reflect.Type, error)

	// TransformToOnChain transforms a type returned from AdjustForInput into the outputType.
	// You may also pass a pointer to the type returned by AdjustForInput to get a pointer to outputType.
	TransformToOnChain(offChainValue any, itemType string) (any, error)

	// TransformToOffChain is the reverse of TransformForOnChain input.
	// It is used to send back the object after it has been decoded
	TransformToOffChain(onChainValue any, itemType string) (any, error)
}
