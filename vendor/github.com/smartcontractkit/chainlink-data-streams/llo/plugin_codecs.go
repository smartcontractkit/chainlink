package llo

import (
	"errors"
	"fmt"
	"sort"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"golang.org/x/exp/maps"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"google.golang.org/protobuf/proto"
)

// NOTE: These codecs make a lot of allocations which will be hard on the
// garbage collector, this can probably be made more efficient

// OBSERVATION CODEC

var (
	_ ObservationCodec = (*protoObservationCodec)(nil)
)

type ObservationCodec interface {
	Encode(obs Observation) (types.Observation, error)
	Decode(encoded types.Observation) (obs Observation, err error)
}

type protoObservationCodec struct{}

func (c protoObservationCodec) Encode(obs Observation) (types.Observation, error) {
	dfns := channelDefinitionsToProtoObservation(obs.UpdateChannelDefinitions)

	streamValues := make(map[uint32]*LLOStreamValue, len(obs.StreamValues))
	for id, sv := range obs.StreamValues {
		if sv != nil {
			enc, err := sv.MarshalBinary()
			if errors.Is(err, ErrNilStreamValue) {
				// Ignore nil values
				continue
			} else if err != nil {
				return nil, err
			}
			streamValues[id] = &LLOStreamValue{
				Type:  sv.Type(),
				Value: enc,
			}
		}
	}

	pbuf := &LLOObservationProto{
		AttestedPredecessorRetirement: obs.AttestedPredecessorRetirement,
		ShouldRetire:                  obs.ShouldRetire,
		UnixTimestampNanoseconds:      obs.UnixTimestampNanoseconds,
		RemoveChannelIDs:              maps.Keys(obs.RemoveChannelIDs),
		UpdateChannelDefinitions:      dfns,
		StreamValues:                  streamValues,
	}

	return proto.Marshal(pbuf)
}

func channelDefinitionsToProtoObservation(in llotypes.ChannelDefinitions) (out map[uint32]*LLOChannelDefinitionProto) {
	if len(in) > 0 {
		out = make(map[uint32]*LLOChannelDefinitionProto, len(in))
		for id, d := range in {
			streams := make([]*LLOStreamDefinition, len(d.Streams))
			for i, strm := range d.Streams {
				streams[i] = &LLOStreamDefinition{
					StreamID:   strm.StreamID,
					Aggregator: uint32(strm.Aggregator),
				}
			}
			out[id] = &LLOChannelDefinitionProto{
				ReportFormat: uint32(d.ReportFormat),
				Streams:      streams,
				Opts:         d.Opts,
			}
		}
	}
	return
}

// TODO: Guard against untrusted inputs!
// MERC-3524
func (c protoObservationCodec) Decode(b types.Observation) (Observation, error) {
	pbuf := &LLOObservationProto{}
	err := proto.Unmarshal(b, pbuf)
	if err != nil {
		return Observation{}, fmt.Errorf("failed to decode observation: expected protobuf (got: 0x%x); %w", b, err)
	}
	var removeChannelIDs map[llotypes.ChannelID]struct{}
	if len(pbuf.RemoveChannelIDs) > 0 {
		removeChannelIDs = make(map[llotypes.ChannelID]struct{}, len(pbuf.RemoveChannelIDs))
		for _, id := range pbuf.RemoveChannelIDs {
			removeChannelIDs[id] = struct{}{}
		}
	}
	dfns := channelDefinitionsFromProtoObservation(pbuf.UpdateChannelDefinitions)
	var streamValues StreamValues
	if len(pbuf.StreamValues) > 0 {
		streamValues = make(StreamValues, len(pbuf.StreamValues))
		for id, enc := range pbuf.StreamValues {
			// StreamValues shouldn't have explicit nils, but for safety we
			// ought to handle it anyway
			sv, err := UnmarshalProtoStreamValue(enc)
			if err != nil {
				return Observation{}, err
			}
			streamValues[id] = sv
		}
	}
	obs := Observation{
		AttestedPredecessorRetirement: pbuf.AttestedPredecessorRetirement,
		ShouldRetire:                  pbuf.ShouldRetire,
		UnixTimestampNanoseconds:      pbuf.UnixTimestampNanoseconds,
		RemoveChannelIDs:              removeChannelIDs,
		UpdateChannelDefinitions:      dfns,
		StreamValues:                  streamValues,
	}
	return obs, nil
}

// TODO: Needs fuzz testing
// MERC-3524
func channelDefinitionsFromProtoObservation(channelDefinitions map[uint32]*LLOChannelDefinitionProto) llotypes.ChannelDefinitions {
	if len(channelDefinitions) == 0 {
		return nil
	}
	dfns := make(map[llotypes.ChannelID]llotypes.ChannelDefinition, len(channelDefinitions))
	for id, d := range channelDefinitions {
		streams := make([]llotypes.Stream, len(d.Streams))
		for i, strm := range d.Streams {
			streams[i] = llotypes.Stream{
				StreamID:   strm.StreamID,
				Aggregator: llotypes.Aggregator(strm.Aggregator),
			}
		}
		dfns[id] = llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormat(d.ReportFormat),
			Streams:      streams,
			Opts:         d.Opts,
		}
	}
	return dfns
}

// OUTCOME CODEC

var _ OutcomeCodec = (*protoOutcomeCodec)(nil)

type OutcomeCodec interface {
	Encode(outcome Outcome) (ocr3types.Outcome, error)
	Decode(encoded ocr3types.Outcome) (outcome Outcome, err error)
}

type protoOutcomeCodec struct{}

func (protoOutcomeCodec) Encode(outcome Outcome) (ocr3types.Outcome, error) {
	dfns := channelDefinitionsToProtoOutcome(outcome.ChannelDefinitions)

	streamAggregates, err := StreamAggregatesToProtoOutcome(outcome.StreamAggregates)
	if err != nil {
		return nil, err
	}

	validAfterSeconds := validAfterSecondsToProtoOutcome(outcome.ValidAfterSeconds)

	pbuf := &LLOOutcomeProto{
		LifeCycleStage:                   string(outcome.LifeCycleStage),
		ObservationsTimestampNanoseconds: outcome.ObservationsTimestampNanoseconds,
		ChannelDefinitions:               dfns,
		ValidAfterSeconds:                validAfterSeconds,
		StreamAggregates:                 streamAggregates,
	}

	// It's very important that Outcome serialization be deterministic across all nodes!
	// Should be reliable since we don't use maps
	return proto.MarshalOptions{Deterministic: true}.Marshal(pbuf)
}

func channelDefinitionsToProtoOutcome(in llotypes.ChannelDefinitions) (out []*LLOChannelIDAndDefinitionProto) {
	if len(in) > 0 {
		out = make([]*LLOChannelIDAndDefinitionProto, 0, len(in))
		for id, d := range in {
			streams := make([]*LLOStreamDefinition, len(d.Streams))
			for i, strm := range d.Streams {
				streams[i] = &LLOStreamDefinition{
					StreamID:   strm.StreamID,
					Aggregator: uint32(strm.Aggregator),
				}
			}
			out = append(out, &LLOChannelIDAndDefinitionProto{
				ChannelID: id,
				ChannelDefinition: &LLOChannelDefinitionProto{
					ReportFormat: uint32(d.ReportFormat),
					Streams:      streams,
					Opts:         d.Opts,
				},
			})
		}
		sort.Slice(out, func(i, j int) bool {
			return out[i].ChannelID < out[j].ChannelID
		})
	}
	return
}

// TODO: Needs thorough unit testing of all paths including nil handling
// MERC-3524
func StreamAggregatesToProtoOutcome(in StreamAggregates) (out []*LLOStreamAggregate, err error) {
	if len(in) > 0 {
		out = make([]*LLOStreamAggregate, 0, len(in))
		for sid, aggregates := range in {
			if aggregates == nil {
				return nil, fmt.Errorf("cannot marshal protobuf; nil value for stream ID: %d", sid)
			}
			for agg, v := range aggregates {
				if v == nil {
					return nil, fmt.Errorf("cannot marshal protobuf; nil value for stream ID: %d, aggregator: %v", sid, agg)
				}
				value, err := v.MarshalBinary()
				if err != nil {
					return nil, err
				}

				out = append(out, &LLOStreamAggregate{
					StreamID:    sid,
					StreamValue: &LLOStreamValue{Type: v.Type(), Value: value},
					Aggregator:  uint32(agg),
				})
			}
		}
		sort.Slice(out, func(i, j int) bool {
			if out[i].StreamID == out[j].StreamID {
				return out[i].Aggregator < out[j].Aggregator
			}
			return out[i].StreamID < out[j].StreamID
		})
	}
	return
}

func validAfterSecondsToProtoOutcome(in map[llotypes.ChannelID]uint32) (out []*LLOChannelIDAndValidAfterSecondsProto) {
	if len(in) > 0 {
		out = make([]*LLOChannelIDAndValidAfterSecondsProto, 0, len(in))
		for id, v := range in {
			out = append(out, &LLOChannelIDAndValidAfterSecondsProto{
				ChannelID:         id,
				ValidAfterSeconds: v,
			})
		}
		sort.Slice(out, func(i, j int) bool {
			return out[i].ChannelID < out[j].ChannelID
		})
	}
	return
}

// TODO: Guard against untrusted inputs!
// MERC-3524
func (protoOutcomeCodec) Decode(b ocr3types.Outcome) (outcome Outcome, err error) {
	pbuf := &LLOOutcomeProto{}
	err = proto.Unmarshal(b, pbuf)
	if err != nil {
		return Outcome{}, fmt.Errorf("failed to decode outcome: expected protobuf (got: 0x%x); %w", b, err)
	}
	dfns := channelDefinitionsFromProtoOutcome(pbuf.ChannelDefinitions)
	streamAggregates, err := streamAggregatesFromProtoOutcome(pbuf.StreamAggregates)
	if err != nil {
		return Outcome{}, err
	}
	validAfterSeconds := validAfterSecondsFromProtoOutcome(pbuf.ValidAfterSeconds)
	outcome = Outcome{
		LifeCycleStage:                   llotypes.LifeCycleStage(pbuf.LifeCycleStage),
		ObservationsTimestampNanoseconds: pbuf.ObservationsTimestampNanoseconds,
		ChannelDefinitions:               dfns,
		ValidAfterSeconds:                validAfterSeconds,
		StreamAggregates:                 streamAggregates,
	}
	return outcome, nil
}

// TODO: Needs fuzz testing
// MERC-3524
func channelDefinitionsFromProtoOutcome(in []*LLOChannelIDAndDefinitionProto) (out llotypes.ChannelDefinitions) {
	if len(in) > 0 {
		out = make(map[llotypes.ChannelID]llotypes.ChannelDefinition, len(in))
		for _, d := range in {
			streams := make([]llotypes.Stream, len(d.ChannelDefinition.Streams))
			for i, strm := range d.ChannelDefinition.Streams {
				streams[i] = llotypes.Stream{
					StreamID:   strm.StreamID,
					Aggregator: llotypes.Aggregator(strm.Aggregator),
				}
			}
			out[d.ChannelID] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormat(d.ChannelDefinition.ReportFormat),
				Streams:      streams,
				Opts:         d.ChannelDefinition.Opts,
			}
		}
	}
	return
}

// TODO: Needs fuzz testing
// MERC-3524
func streamAggregatesFromProtoOutcome(in []*LLOStreamAggregate) (out StreamAggregates, err error) {
	if len(in) > 0 {
		out = make(StreamAggregates, len(in))
		for _, enc := range in {
			var sv StreamValue
			sv, err = UnmarshalProtoStreamValue(enc.StreamValue)
			if err != nil {
				return
			}
			m, exists := out[enc.StreamID]
			if !exists {
				m = make(map[llotypes.Aggregator]StreamValue)
				out[enc.StreamID] = m
			}
			m[llotypes.Aggregator(enc.Aggregator)] = sv
		}
	}
	return
}

// TODO: Needs fuzz testing
// MERC-3524
func validAfterSecondsFromProtoOutcome(in []*LLOChannelIDAndValidAfterSecondsProto) (out map[llotypes.ChannelID]uint32) {
	if len(in) > 0 {
		out = make(map[llotypes.ChannelID]uint32, len(in))
		for _, v := range in {
			out[v.ChannelID] = v.ValidAfterSeconds
		}
	}
	return
}
