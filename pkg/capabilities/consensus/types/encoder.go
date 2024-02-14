package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Encoder interface {
	Encode(ctx context.Context, input values.Map) ([]byte, error)
}
