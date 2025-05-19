package aggregation

import (
	"testing"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var (
	triggerEvent1 = map[string]any{"event": "triggerEvent1"}
	triggerEvent2 = map[string]any{"event": "triggerEvent2"}
)

func TestDefaultModeAggregator_Aggregate(t *testing.T) {
	val, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	capResponse1 := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			Outputs: val,
		},
		Err: nil,
	}
	marshaled1, err := pb.MarshalTriggerResponse(capResponse1)
	require.NoError(t, err)

	val2, err := values.NewMap(triggerEvent2)
	require.NoError(t, err)
	capResponse2 := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			Outputs: val2,
		},
		Err: nil,
	}
	marshaled2, err := pb.MarshalTriggerResponse(capResponse2)
	require.NoError(t, err)

	agg := NewDefaultModeAggregator(2)
	_, err = agg.Aggregate("", [][]byte{marshaled1})
	require.Error(t, err)

	_, err = agg.Aggregate("", [][]byte{marshaled1, marshaled2})
	require.Error(t, err)

	res, err := agg.Aggregate("", [][]byte{marshaled1, marshaled2, marshaled1})
	require.NoError(t, err)
	require.Equal(t, res, capResponse1)
}
