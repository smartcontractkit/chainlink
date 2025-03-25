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
	t.Run("Encode", func(t *testing.T) {
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
	t.Run("Verify", func(t *testing.T) {
		donID := uint32(1)
		c := NewReportCodecCapabilityTrigger(logger.Test(t), donID)

		err := c.Verify(llotypes.ChannelDefinition{})
		require.NoError(t, err)

		err = c.Verify(llotypes.ChannelDefinition{Opts: []byte{1, 2, 3}})
		require.EqualError(t, err, "capability trigger does not support channel definitions with options")
	})
}
