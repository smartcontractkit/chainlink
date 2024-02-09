package v3_test

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercury_v3_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

type StaticDataSource struct{}

var _ mercury_v3_types.DataSource = StaticDataSource{}

func (StaticDataSource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (mercury_v3_types.Observation, error) {
	return Fixtures.Observation, nil
}
