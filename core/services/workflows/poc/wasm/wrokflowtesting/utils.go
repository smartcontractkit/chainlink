package wrokflowtesting

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type TargetMock[T any] struct {
	Seen []T
}

func (t *TargetMock[T]) AddTarget(ref string, registry *Registry) {
	registry.RegisterRemoteTarget(ref, func(value values.Value) error {
		unwrapped, err := capabilities.UnwrapValue[T](value)
		if err != nil {
			return err
		}
		t.Seen = append(t.Seen, unwrapped)
		return nil
	})
}
