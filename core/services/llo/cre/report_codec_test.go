package cre

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/wsrpc/logger"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

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
		assert.Equal(t, &pb.Value_Uint64Value{Uint64Value: 34}, pbuf.Outputs.Fields["ObservationTimestampNanoseconds"].Value)
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
		assert.Equal(t, &pb.Value_Uint64Value{Uint64Value: 34}, pbuf.Outputs.Fields["ObservationTimestampNanoseconds"].Value)
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
	t.Run("Decode: With Opts SUCCESS", func(t *testing.T) {
		optBytes := []byte{123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 115, 34, 58, 91, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 48, 49, 48, 49, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 48, 49, 48, 50, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 49, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 50, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 51, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 52, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 53, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 54, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 55, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 56, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 48, 57, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 48, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 49, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 50, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 51, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 52, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 53, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 54, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 55, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 56, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 49, 57, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 50, 48, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 50, 49, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 50, 50, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 50, 51, 125, 44, 123, 34, 109, 117, 108, 116, 105, 112, 108, 105, 101, 114, 34, 58, 34, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 116, 114, 101, 97, 109, 73, 68, 34, 58, 49, 48, 50, 48, 48, 48, 49, 48, 50, 52, 125, 93, 125}

		opts := &ReportCodecCapabilityTriggerOpts{}
		err := opts.Decode(optBytes)

		require.NoError(t, err)

		// Verify the decoded opts structure contains expected multipliers and stream IDs
		require.Len(t, opts.Multipliers, 26)

		expectedMultiplier, err := decimal.NewFromString("1000000000000000000") // 10^18
		require.NoError(t, err)

		expectedStreamIDs := []uint32{
			1020000101, 1020000102, 1020001001, 1020001002, 1020001003, 1020001004,
			1020001005, 1020001006, 1020001007, 1020001008, 1020001009, 1020001010,
			1020001011, 1020001012, 1020001013, 1020001014, 1020001015, 1020001016,
			1020001017, 1020001018, 1020001019, 1020001020, 1020001021, 1020001022,
			1020001023, 1020001024,
		}

		for i, multiplier := range opts.Multipliers {
			assert.True(t, multiplier.Multiplier.Equal(expectedMultiplier), "Multiplier %d should be %s", i, expectedMultiplier.String())
			assert.Equal(t, expectedStreamIDs[i], multiplier.StreamID, "StreamID %d should be %d", i, expectedStreamIDs[i])
		}
	})
	t.Run("Verify: Without Opts SUCCESS", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		err := c.Verify(
			llotypes.ChannelDefinition{
				Streams: []llotypes.Stream{
					{StreamID: 1},
					{StreamID: 2},
				},
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
	t.Run("Verify: Multiplier isn't an integer FAIL", func(t *testing.T) {
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
