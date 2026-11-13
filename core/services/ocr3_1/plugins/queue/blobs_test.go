package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

func Test_extractEvent(t *testing.T) {
	payload, err := anypb.New(timestamppb.New(time.Now()))
	require.NoError(t, err)
	event := &pb.CoordinatedEvent{
		WorkflowID:   "want",
		TriggerCapID: "capID",
		TriggerIndex: 1,
		ObservedAt:   &timestamppb.Timestamp{Seconds: 10, Nanos: 42},
		Deadline:     timestamppb.New(time.Now()),
		Event: &pb.TriggerResponse{
			Event: &pb.TriggerEvent{
				TriggerType: "type",
				Id:          "id",
				Outputs:     nil,
				Payload:     payload,
			},
			Error: "asdf",
		},
	}
	otherPayload, err := anypb.New(&pb.CapabilityResponse{})
	require.NoError(t, err)
	others := &pb.CoordinatedEvent{
		WorkflowID:   "foo",
		TriggerCapID: "bar",
		TriggerIndex: 99,
		ObservedAt:   &timestamppb.Timestamp{Nanos: 99},
		Deadline:     &timestamppb.Timestamp{Seconds: 1000},
		Event: &pb.TriggerResponse{
			Event: &pb.TriggerEvent{
				TriggerType: "baz",
				Id:          "trigger-id",
				Outputs:     nil,
				Payload:     otherPayload,
			},
			Error: "none",
		},
	}
	for _, tt := range []struct {
		name     string
		len, idx int
	}{
		{"first", 4, 0},
		{"last", 4, 3},
		{"middle", 4, 2},
		{"bounds", 4, 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var events []*pb.CoordinatedEvent
			for i := range tt.len {
				if i == tt.idx {
					events = append(events, event)
				} else {
					events = append(events, others)
				}
			}
			b, err := proto.MarshalOptions{Deterministic: true}.Marshal(&pb.CoordinatedEventList{Events: events})
			require.NoError(t, err)
			got, err := extractEvent(b, tt.idx)
			if tt.idx < tt.len {
				require.True(t, proto.Equal(event, got))
			} else {
				require.Error(t, err)
			}
		})
	}
}
