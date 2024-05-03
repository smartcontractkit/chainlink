package workflow

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type mergeRunner2[I, II, O any] struct {
	mergeRunnerBase
	fn func(I I, II II) (O, error)
}

func (m mergeRunner2[I, II, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner2[any, any, any]{}

type mergeRunner3[I, II, III, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III) (O, error)
}

func (m mergeRunner3[I, II, III, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner3[any, any, any, any]{}

type mergeRunner4[I, II, III, IIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII) (O, error)
}

func (m mergeRunner4[I, II, III, IIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner4[any, any, any, any, any]{}

type mergeRunner5[I, II, III, IIII, IIIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII, IIIII IIIII) (O, error)
}

func (m mergeRunner5[I, II, III, IIII, IIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
	v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4, v5)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner5[any, any, any, any, any, any]{}

type mergeRunner6[I, II, III, IIII, IIIII, IIIIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII, IIIII IIIII, IIIIII IIIIII) (O, error)
}

func (m mergeRunner6[I, II, III, IIII, IIIII, IIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
	v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
	v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4, v5, v6)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner6[any, any, any, any, any, any, any]{}

type mergeRunner7[I, II, III, IIII, IIIII, IIIIII, IIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII, IIIII IIIII, IIIIII IIIIII, IIIIIII IIIIIII) (O, error)
}

func (m mergeRunner7[I, II, III, IIII, IIIII, IIIIII, IIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
	v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
	v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
	v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4, v5, v6, v7)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner7[any, any, any, any, any, any, any, any]{}

type mergeRunner8[I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII, IIIII IIIII, IIIIII IIIIII, IIIIIII IIIIIII, IIIIIIII IIIIIIII) (O, error)
}

func (m mergeRunner8[I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
	v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
	v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
	v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}
	v8, err := capabilities.UnwrapValue[IIIIIIII](ls[8-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4, v5, v6, v7, v8)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner8[any, any, any, any, any, any, any, any, any]{}

type mergeRunner9[I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, IIIIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII, IIIII IIIII, IIIIII IIIIII, IIIIIII IIIIIII, IIIIIIII IIIIIIII, IIIIIIIII IIIIIIIII) (O, error)
}

func (m mergeRunner9[I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, IIIIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
	v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
	v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
	v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}
	v8, err := capabilities.UnwrapValue[IIIIIIII](ls[8-1])
	if err != nil {
		return nil, false, err
	}
	v9, err := capabilities.UnwrapValue[IIIIIIIII](ls[9-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4, v5, v6, v7, v8, v9)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner9[any, any, any, any, any, any, any, any, any, any]{}

type mergeRunner10[I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, IIIIIIIII, IIIIIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I I, II II, III III, IIII IIII, IIIII IIIII, IIIIII IIIIII, IIIIIII IIIIIII, IIIIIIII IIIIIIII, IIIIIIIII IIIIIIIII, IIIIIIIIII IIIIIIIIII) (O, error)
}

func (m mergeRunner10[I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, IIIIIIIII, IIIIIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
	v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
	v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
	v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
	v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
	v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
	v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
	v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}
	v8, err := capabilities.UnwrapValue[IIIIIIII](ls[8-1])
	if err != nil {
		return nil, false, err
	}
	v9, err := capabilities.UnwrapValue[IIIIIIIII](ls[9-1])
	if err != nil {
		return nil, false, err
	}
	v10, err := capabilities.UnwrapValue[IIIIIIIIII](ls[10-1])
	if err != nil {
		return nil, false, err
	}

	merged, err := m.fn(v1, v2, v3, v4, v5, v6, v7, v8, v9, v10)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner10[any, any, any, any, any, any, any, any, any, any, any]{}
