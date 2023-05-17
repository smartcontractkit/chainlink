package test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ median.DataSource = (*staticDataSource)(nil)

type staticDataSource struct {
	value *big.Int
}

func StaticDataSource() median.DataSource { return &staticDataSource{value} }

func StaticJuelsPerFeeCoinDataSource() median.DataSource { return &staticDataSource{juelsPerFeeCoin} }

func (s *staticDataSource) Observe(ctx context.Context, timestamp types.ReportTimestamp) (*big.Int, error) {
	if timestamp != reportContext.ReportTimestamp {
		return nil, fmt.Errorf("expected %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	return s.value, nil
}
