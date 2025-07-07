package cre

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/wsrpc/logger"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReportCodec(t *testing.T) {
	t.Run("Encode: Without Opts SUCCESS", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		r := datastreamsllo.Report{
			ConfigDigest:                    types.ConfigDigest{1, 2, 3},
			SeqNr:                           32,
			ChannelID:                       llotypes.ChannelID(31),
			ValidAfterNanoseconds:           28,
			ObservationTimestampNanoseconds: 34,
			Values:                          []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36))},
			Specimen:                        false,
		}
		encoded, err := c.Encode(r, llotypes.ChannelDefinition{
			Streams: []llotypes.Stream{
				{StreamID: 1},
				{StreamID: 2},
			},
		})
		require.NoError(t, err)

		var pbuf capabilitiespb.OCRTriggerReport
		err = proto.Unmarshal(encoded, &pbuf)
		require.NoError(t, err)

		assert.Equal(t, "streams_1_34", pbuf.EventID)
		assert.Equal(t, uint64(34), pbuf.Timestamp)
		require.Len(t, pbuf.Outputs.Fields, 2)
		assert.Equal(t, &pb.Value_Int64Value{Int64Value: 34}, pbuf.Outputs.Fields["ObservationTimestampNanoseconds"].Value)
		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields, 2)

		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[0].Value.(*pb.Value_MapValue).MapValue.Fields, 2)
		decimalBytes := pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[0].Value.(*pb.Value_MapValue).MapValue.Fields["Decimal"].Value.(*pb.Value_BytesValue).BytesValue
		d := decimal.Decimal{}
		require.NoError(t, (&d).UnmarshalBinary(decimalBytes))
		assert.Equal(t, "35", d.String())
		assert.Equal(t, int64(1), pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[0].Value.(*pb.Value_MapValue).MapValue.Fields["StreamID"].Value.(*pb.Value_Int64Value).Int64Value)

		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[1].Value.(*pb.Value_MapValue).MapValue.Fields, 2)
		decimalBytes = pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[1].Value.(*pb.Value_MapValue).MapValue.Fields["Decimal"].Value.(*pb.Value_BytesValue).BytesValue
		d = decimal.Decimal{}
		require.NoError(t, (&d).UnmarshalBinary(decimalBytes))
		assert.Equal(t, "36", d.String())
		assert.Equal(t, int64(2), pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[1].Value.(*pb.Value_MapValue).MapValue.Fields["StreamID"].Value.(*pb.Value_Int64Value).Int64Value)
	})
	t.Run("Encode: With Opts SUCCESS", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		r := datastreamsllo.Report{
			ConfigDigest:                    types.ConfigDigest{1, 2, 3},
			SeqNr:                           32,
			ChannelID:                       llotypes.ChannelID(31),
			ValidAfterNanoseconds:           28,
			ObservationTimestampNanoseconds: 34,
			Values:                          []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36)), llo.ToDecimal(decimal.NewFromInt(37))},
			Specimen:                        false,
		}

		multiplier1, err := decimal.NewFromString("1")
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("1000000000000000000") // 10^18
		require.NoError(t, err)
		multiplier3, err := decimal.NewFromString("1000000") // 10^6
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 2},
				{Multiplier: multiplier3, StreamID: 3},
			},
		}).Encode()
		require.NoError(t, err)
		encoded, err := c.Encode(r, llotypes.ChannelDefinition{
			Streams: []llotypes.Stream{
				{StreamID: 1},
				{StreamID: 2},
				{StreamID: 3},
			},
			Opts: opts,
		})
		require.NoError(t, err)

		var pbuf capabilitiespb.OCRTriggerReport
		err = proto.Unmarshal(encoded, &pbuf)
		require.NoError(t, err)

		assert.Equal(t, "streams_1_34", pbuf.EventID)
		assert.Equal(t, uint64(34), pbuf.Timestamp)
		require.Len(t, pbuf.Outputs.Fields, 2)
		assert.Equal(t, &pb.Value_Int64Value{Int64Value: 34}, pbuf.Outputs.Fields["ObservationTimestampNanoseconds"].Value)
		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields, 3)

		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[0].Value.(*pb.Value_MapValue).MapValue.Fields, 2)
		decimalBytes := pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[0].Value.(*pb.Value_MapValue).MapValue.Fields["Decimal"].Value.(*pb.Value_BytesValue).BytesValue
		d := decimal.Decimal{}
		require.NoError(t, (&d).UnmarshalBinary(decimalBytes))
		assert.Equal(t, "35", d.String())
		assert.Equal(t, int64(1), pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[0].Value.(*pb.Value_MapValue).MapValue.Fields["StreamID"].Value.(*pb.Value_Int64Value).Int64Value)

		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[1].Value.(*pb.Value_MapValue).MapValue.Fields, 2)
		decimalBytes = pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[1].Value.(*pb.Value_MapValue).MapValue.Fields["Decimal"].Value.(*pb.Value_BytesValue).BytesValue
		d = decimal.Decimal{}
		require.NoError(t, (&d).UnmarshalBinary(decimalBytes))
		assert.Equal(t, "36000000000000000000", d.String())
		assert.Equal(t, int64(2), pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[1].Value.(*pb.Value_MapValue).MapValue.Fields["StreamID"].Value.(*pb.Value_Int64Value).Int64Value)

		require.Len(t, pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[2].Value.(*pb.Value_MapValue).MapValue.Fields, 2)
		decimalBytes = pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[2].Value.(*pb.Value_MapValue).MapValue.Fields["Decimal"].Value.(*pb.Value_BytesValue).BytesValue
		d = decimal.Decimal{}
		require.NoError(t, (&d).UnmarshalBinary(decimalBytes))
		assert.Equal(t, "37000000", d.String())
		assert.Equal(t, int64(3), pbuf.Outputs.Fields["Payload"].Value.(*pb.Value_ListValue).ListValue.Fields[2].Value.(*pb.Value_MapValue).MapValue.Fields["StreamID"].Value.(*pb.Value_Int64Value).Int64Value)
	})
	t.Run("Verify: Without Opts SUCCESS", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		multiplier1, err := decimal.NewFromString("1")
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("1000000000000000000") // 10^18
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 2},
			},
		}).Encode()
		require.NoError(t, err)
		err = c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 2},
				},
				Opts: opts,
			},
		)
		require.NoError(t, err)
	})
	t.Run("Verify: Misaligned Multiplier StreamIDs FAIL", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		multiplier1, err := decimal.NewFromString("1")
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("1000000000000000000") // 10^18
		require.NoError(t, err)
		multiplier3, err := decimal.NewFromString("1000000") // 10^6
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 3},
				{Multiplier: multiplier3, StreamID: 2},
			},
		}).Encode()
		require.NoError(t, err)
		err = c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 2},
					{StreamID: 3},
				},
				Opts: opts,
			},
		)
		require.EqualError(t, err, "LLO StreamID 2 mismatched with Multiplier StreamID 3")
	})
	t.Run("Encode: Multiplier isn't an integer FAIL", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		multiplier1, err := decimal.NewFromString("123.4567")
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("89.01234")
		require.NoError(t, err)
		multiplier3, err := decimal.NewFromString("1000000") // 10^6
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 2},
				{Multiplier: multiplier3, StreamID: 3},
			},
		}).Encode()
		require.NoError(t, err)
		err = c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 2},
					{StreamID: 3},
				},
				Opts: opts,
			},
		)
		require.EqualError(t, err, "multiplier for StreamID 1 must be an integer")
	})
	t.Run("Verify: Multiplier is zero FAIL", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		multiplier1, err := decimal.NewFromString("0")
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("0")
		require.NoError(t, err)
		multiplier3, err := decimal.NewFromString("1000000") // 10^6
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 2},
				{Multiplier: multiplier3, StreamID: 3},
			},
		}).Encode()
		require.NoError(t, err)
		err = c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 2},
					{StreamID: 3},
				},
				Opts: opts,
			},
		)
		require.EqualError(t, err, "multiplier for StreamID 1 can't be zero")
	})
	t.Run("Verify: Multiplier is negative FAIL", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		multiplier1, err := decimal.NewFromString("-1000000000000000000") // -10^18
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("-1")
		require.NoError(t, err)
		multiplier3, err := decimal.NewFromString("1000000") // 10^6
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 2},
				{Multiplier: multiplier3, StreamID: 3},
			},
		}).Encode()
		require.NoError(t, err)
		err = c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 2},
					{StreamID: 3},
				},
				Opts: opts,
			},
		)
		require.EqualError(t, err, "multiplier for StreamID 1 can't be negative")
	})
	t.Run("Verify: Multipliers length, StreamValues length mismatch FAIL", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		multiplier1, err := decimal.NewFromString("1000000000000000000") // 10^18
		require.NoError(t, err)
		multiplier2, err := decimal.NewFromString("1")
		require.NoError(t, err)
		multiplier3, err := decimal.NewFromString("1000000") // 10^6
		require.NoError(t, err)

		opts, err := (&ReportCodecCapabilityTriggerOpts{
			Multipliers: []ReportCodecCapabilityTriggerMultiplier{
				{Multiplier: multiplier1, StreamID: 1},
				{Multiplier: multiplier2, StreamID: 2},
				{Multiplier: multiplier3, StreamID: 3},
			},
		}).Encode()
		require.NoError(t, err)

		err = c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 3},
				},
				Opts: opts,
			},
		)
		require.EqualError(t, err, "multipliers length 3 != StreamValues length 2")
	})
}
