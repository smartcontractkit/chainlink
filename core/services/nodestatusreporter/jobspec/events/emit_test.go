package events_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"

	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/jobspec/events"
)

func TestEmitJobSpecEvent_RoundTrip(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	emitter := beholder.GetEmitter()

	event := &events.JobSpecEvent{
		ExternalJobId:   "test-job-id",
		Name:            "test-job",
		JobType:         "offchainreporting2",
		EmissionTrigger: events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT,
	}

	err := events.EmitJobSpecEvent(context.Background(), emitter, event)
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	msg := msgs[0]
	require.Equal(t, "test-job-id", msg.Attrs["partitionkey"])
	require.Equal(t, events.SchemaJobSpec, msg.Attrs["beholder_data_schema"])
	require.Equal(t, events.BeholderDomain, msg.Attrs["beholder_domain"])
	require.NotContains(t, msg.Attrs, "source")

	var decoded events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msg.Body, &decoded))
	require.Equal(t, "test-job-id", decoded.ExternalJobId)
	require.Equal(t, "test-job", decoded.Name)
	require.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT, decoded.EmissionTrigger)
	require.NotEmpty(t, decoded.Timestamp)
}

func TestEmitJobSpecEvent_SetsTimestampIfEmpty(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	emitter := beholder.GetEmitter()

	event := &events.JobSpecEvent{
		ExternalJobId: "ts-test",
	}
	require.Empty(t, event.Timestamp)

	err := events.EmitJobSpecEvent(context.Background(), emitter, event)
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var decoded events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &decoded))
	require.NotEmpty(t, decoded.Timestamp)
}
