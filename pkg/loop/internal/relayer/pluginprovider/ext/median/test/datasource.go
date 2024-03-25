package median_test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

var _ median.DataSource = (*staticDataSource)(nil)
var _ testtypes.Evaluator[median.DataSource] = (*staticDataSource)(nil)

type staticDataSourceConfig struct {
	ReportContext types.ReportContext
	Value         *big.Int
}

type staticDataSource struct {
	staticDataSourceConfig
}

var (
	DataSource = staticDataSource{
		staticDataSourceConfig{
			ReportContext: reportContext,
			Value:         value,
		},
	}

	JuelsPerFeeCoinDataSource = staticDataSource{
		staticDataSourceConfig{
			ReportContext: reportContext,
			Value:         juelsPerFeeCoin,
		},
	}
)

func (s staticDataSource) Observe(ctx context.Context, timestamp types.ReportTimestamp) (*big.Int, error) {
	if timestamp != s.ReportContext.ReportTimestamp {
		return nil, fmt.Errorf("expected %v but got %v", s.ReportContext.ReportTimestamp, timestamp)
	}
	return s.Value, nil
}

func (s staticDataSource) Evaluate(ctx context.Context, ds median.DataSource) error {
	gotVal, err := ds.Observe(ctx, s.ReportContext.ReportTimestamp)
	if err != nil {
		return fmt.Errorf("failed to observe dataSource: %w", err)
	}
	if gotVal.Cmp(s.Value) != 0 {
		return fmt.Errorf("expected Value %s but got %s", value, gotVal)
	}
	return nil
}
