package v2_test

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
)

type StaticDataSource struct{}

var _ mercury_v2_types.DataSource = StaticDataSource{}

func (StaticDataSource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (mercury_v2_types.Observation, error) {
	return Fixtures.Observation, nil
}
