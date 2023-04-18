package loop_test

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

var _ median.DataSource = (*staticDataSource)(nil)

type staticDataSource struct {
	value *big.Int
}

func (s *staticDataSource) Observe(ctx context.Context) (*big.Int, error) {
	return s.value, nil
}
