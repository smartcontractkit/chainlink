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

	GasPriceSubunitsDataSource = staticDataSource{
		staticDataSourceConfig{
			ReportContext: reportContext,
			Value:         gasPriceSubunits,
		},
	}
)

func (s staticDataSource) Observe(ctx context.Context, timestamp types.ReportTimestamp) (*big.Int, error) {
	if timestamp != s.ReportContext.ReportTimestamp {
		return nil, fmt.Errorf("expected %v but got %v", s.ReportContext.ReportTimestamp, timestamp)
	}
	return s.Value, nil
}

type CompareError struct {
	Got      *big.Int
	Expected *big.Int
}

func (e *CompareError) Error() string {
	return fmt.Sprintf("expected Value %s but got %s", e.Expected, e.Got)
}

func (e *CompareError) GotZero() bool {
	return e.Got.Uint64() == 0
}

func (s staticDataSource) Evaluate(ctx context.Context, ds median.DataSource) error {
	gotVal, err := ds.Observe(ctx, s.ReportContext.ReportTimestamp)
	if err != nil {
		return fmt.Errorf("failed to observe dataSource: %w", err)
	}
	if gotVal.Cmp(s.Value) != 0 {
		return &CompareError{Got: gotVal, Expected: s.Value}
	}
	return nil
}

// Only to be used for testing
type ZeroDataSource struct {
}

func (s *ZeroDataSource) Observe(ctx context.Context, _ types.ReportTimestamp) (*big.Int, error) {
	return big.NewInt(0), nil
}
