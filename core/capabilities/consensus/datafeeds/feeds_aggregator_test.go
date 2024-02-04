package datafeeds_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/consensus/datafeeds"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/consensus/types"
)

func TestDataFeedsAggregator(t *testing.T) {
	agg, err := datafeeds.NewDataFeedsAggregator(values.Map{}, nil, nil)
	require.NoError(t, err)

	_, err = agg.Aggregate(&types.AggregationOutcome{}, map[commontypes.OracleID][]values.Value{})
	require.NoError(t, err)
}
