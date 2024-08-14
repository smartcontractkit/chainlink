package fixtures

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func Test_MockTriggerMessageDistribution(t *testing.T) {
	ctx := tests.Context(t)
	sink := NewReportsSink()
	servicetest.Run(t, sink)

	trigger1 := sink.GetNewTrigger(t)
	trigger2 := sink.GetNewTrigger(t)
	trigger3 := sink.GetNewTrigger(t)

	trigger1Ch, err := trigger1.RegisterTrigger(ctx, capabilities.CapabilityRequest{})
	require.NoError(t, err)
	trigger2Ch, err := trigger2.RegisterTrigger(ctx, capabilities.CapabilityRequest{})
	require.NoError(t, err)

	sink.SendReports([]*datastreams.FeedReport{
		{FeedID: "DEADBEEF11"},
		{FeedID: "DEADBEEF22"},
	})

	tr1 := <-trigger1Ch
	assert.Equal(t, "DEADBEEF11", getFeedID(t, tr1, 0))
	assert.Equal(t, "DEADBEEF22", getFeedID(t, tr1, 1))

	tr2 := <-trigger2Ch
	assert.Equal(t, "DEADBEEF11", getFeedID(t, tr2, 0))
	assert.Equal(t, "DEADBEEF22", getFeedID(t, tr2, 1))

	// Late registered trigger should still receive all messages sent to the sink
	trigger3Ch, err := trigger3.RegisterTrigger(ctx, capabilities.CapabilityRequest{})
	require.NoError(t, err)

	tr3 := <-trigger3Ch
	assert.Equal(t, "DEADBEEF11", getFeedID(t, tr3, 0))
	assert.Equal(t, "DEADBEEF22", getFeedID(t, tr3, 1))

	sink.SendReports([]*datastreams.FeedReport{
		{FeedID: "DEADBEEF33"},
		{FeedID: "DEADBEEF44"},
	})

	tr1 = <-trigger1Ch
	assert.Equal(t, "DEADBEEF33", getFeedID(t, tr1, 0))
	assert.Equal(t, "DEADBEEF44", getFeedID(t, tr1, 1))

	tr2 = <-trigger2Ch
	assert.Equal(t, "DEADBEEF33", getFeedID(t, tr2, 0))
	assert.Equal(t, "DEADBEEF44", getFeedID(t, tr2, 1))

	tr3 = <-trigger3Ch
	assert.Equal(t, "DEADBEEF33", getFeedID(t, tr3, 0))
	assert.Equal(t, "DEADBEEF44", getFeedID(t, tr3, 1))

}

func getFeedID(t *testing.T, t1R capabilities.CapabilityResponse, feedIdx int) string {
	reports := t1R.Value.Underlying["Payload"].(*values.List)
	report := reports.Underlying[feedIdx]
	report1Map := map[string]any{}
	err := report.UnwrapTo(&report1Map)
	require.NoError(t, err)
	feedID := report1Map["FeedID"]
	return feedID.(string)
}
