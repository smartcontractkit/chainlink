package v1_test

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
)

type StaticDataSource struct{}

var _ mercury_v1_types.DataSource = StaticDataSource{}

func (StaticDataSource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, FetchMaxFinalizedBlockNum bool) (mercury_v1_types.Observation, error) {
	return Fixtures.Observation, nil
}
