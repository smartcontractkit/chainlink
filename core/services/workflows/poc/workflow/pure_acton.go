package workflow

import "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"

func AddTransform[I, O any](builder *Builder[I], ref string, fn func(I) (O, error)) (*Builder[O], error) {
	return AddTransformAndFilter[I, O](builder, ref, func(i I) (O, bool, error) {
		o, err := fn(i)
		return o, err == nil, err
	})
}

func AddFilter[T any](builder *Builder[T], ref string, fn func(T) bool) (*Builder[T], error) {
	return AddTransformAndFilter[T, T](builder, ref, func(i T) (T, bool, error) {
		return i, fn(i), nil
	})
}

func AddTransformAndFilter[I, O any](builder *Builder[I], ref string, fn func(I) (O, bool, error)) (*Builder[O], error) {
	// Note there's probably a way to be more efficient than to have the channel, just call the function but it's fine for now.
	return AddStep[I, O](builder, &pureAction[I, O]{ref: ref, fn: fn})
}

type pureAction[I, O any] struct {
	ref string
	fn  func(I) (O, bool, error)
}

func (p pureAction[I, O]) Type() string {
	return capabilities.LocalCodeActionCapability
}

func (p pureAction[I, O]) Ref() string {
	return p.ref
}

func (p pureAction[I, O]) Invoke(input I) (O, bool, error) {
	return p.fn(input)
}
