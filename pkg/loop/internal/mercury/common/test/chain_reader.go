package mercury_common_test

import (
	"context"

	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

type StaticMercuryChainReader struct{}

var _ mercury_types.ChainReader = StaticMercuryChainReader{}

func (StaticMercuryChainReader) LatestHeads(ctx context.Context, n int) ([]mercury_types.Head, error) {
	return nil, nil
}
