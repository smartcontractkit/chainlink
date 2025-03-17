package evm

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

func FuzzReportCodecPremiumLegacy_Decode(f *testing.F) {
	f.Add([]byte("not a protobuf"))
	f.Add([]byte{0x0a, 0x00})             // empty protobuf
	f.Add([]byte{0x0a, 0x02, 0x08, 0x01}) // invalid protobuf
	f.Add(([]byte)(nil))
	f.Add([]byte{})

	validReport := newValidPremiumLegacyReport()
	feedID := [32]uint8{0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	cd := llotypes.ChannelDefinition{Opts: llotypes.ChannelOpts(fmt.Sprintf(`{"baseUSDFee":"10.50","expirationWindow":60,"feedId":"0x%x","multiplier":10}`, feedID))}

	codec := ReportCodecPremiumLegacy{logger.NullLogger, 100002}

	validEncodedReport, err := codec.Encode(tests.Context(f), validReport, cd)
	require.NoError(f, err)
	f.Add(validEncodedReport)

	f.Fuzz(func(t *testing.T, data []byte) {
		codec.Decode(data) //nolint:errcheck // test that it doesn't panic, don't care about errors
	})
}

func newValidPremiumLegacyReport() llo.Report {
	return llo.Report{
		ConfigDigest:                    types.ConfigDigest{1, 2, 3},
		SeqNr:                           32,
		ChannelID:                       llotypes.ChannelID(31),
		ValidAfterNanoseconds:           28 * 1e9,
		ObservationTimestampNanoseconds: 34 * 1e9,
		Values:                          []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36)), &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}},
		Specimen:                        false,
	}
}

func Test_ReportCodecPremiumLegacy(t *testing.T) {
	rc := ReportCodecPremiumLegacy{logger.TestLogger(t), 2}

	feedID := [32]uint8{0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	cd := llotypes.ChannelDefinition{Opts: llotypes.ChannelOpts(fmt.Sprintf(`{"baseUSDFee":"10.50","expirationWindow":60,"feedId":"0x%x","multiplier":10}`, feedID))}

	t.Run("Encode errors if no values", func(t *testing.T) {
		ctx := tests.Context(t)
		_, err := rc.Encode(ctx, llo.Report{}, cd)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "ReportCodecPremiumLegacy cannot encode; got unusable report; ReportCodecPremiumLegacy requires exactly 3 values (NativePrice, LinkPrice, Quote{Bid, Mid, Ask}); got report.Values: []")
	})

	t.Run("does not encode specimen reports", func(t *testing.T) {
		ctx := tests.Context(t)
		report := newValidPremiumLegacyReport()
		report.Specimen = true

		_, err := rc.Encode(ctx, report, cd)
		require.Error(t, err)
		require.EqualError(t, err, "ReportCodecPremiumLegacy does not support encoding specimen reports")
	})

	t.Run("Encode constructs a report from observations", func(t *testing.T) {
		ctx := tests.Context(t)
		report := newValidPremiumLegacyReport()

		encoded, err := rc.Encode(ctx, report, cd)
		require.NoError(t, err)

		assert.Len(t, encoded, 288)
		assert.Equal(t, []byte{0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x29, 0xd0, 0x69, 0x18, 0x9e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0xc, 0x35, 0x49, 0xbb, 0x7d, 0x2a, 0xab, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x7c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x86}, encoded)

		decoded, err := reporttypes.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, feedID, decoded.FeedId)
		assert.Equal(t, uint32(34), decoded.ObservationsTimestamp)
		assert.Equal(t, big.NewInt(380), decoded.BenchmarkPrice)
		assert.Equal(t, big.NewInt(370), decoded.Bid)
		assert.Equal(t, big.NewInt(390), decoded.Ask)
		assert.Equal(t, uint32(29), decoded.ValidFromTimestamp)
		assert.Equal(t, uint32(94), decoded.ExpiresAt)
		assert.Equal(t, big.NewInt(291666666666666667), decoded.LinkFee)
		assert.Equal(t, big.NewInt(300000000000000000), decoded.NativeFee)

		t.Run("Decode decodes the report", func(t *testing.T) {
			decoded, err := rc.Decode(encoded)
			require.NoError(t, err)

			assert.Equal(t, &reporttypes.Report{
				FeedId:                [32]uint8{0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				ObservationsTimestamp: 0x22,
				BenchmarkPrice:        big.NewInt(380),
				Bid:                   big.NewInt(370),
				Ask:                   big.NewInt(390),
				ValidFromTimestamp:    0x1d,
				ExpiresAt:             uint32(94),
				LinkFee:               big.NewInt(291666666666666667),
				NativeFee:             big.NewInt(300000000000000000),
			}, decoded)
		})
	})

	t.Run("uses zero values if fees are missing", func(t *testing.T) {
		ctx := tests.Context(t)
		report := llo.Report{
			Values: []llo.StreamValue{nil, nil, &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}},
		}

		encoded, err := rc.Encode(ctx, report, cd)
		require.NoError(t, err)

		assert.Len(t, encoded, 288)
		assert.Equal(t, []byte{0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x7c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x86}, encoded)

		decoded, err := rc.Decode(encoded)
		require.NoError(t, err)

		assert.Equal(t, feedID, decoded.FeedId)
		assert.Equal(t, uint32(0), decoded.ObservationsTimestamp)
		assert.Equal(t, big.NewInt(380).String(), decoded.BenchmarkPrice.String())
		assert.Equal(t, big.NewInt(370).String(), decoded.Bid.String())
		assert.Equal(t, big.NewInt(390).String(), decoded.Ask.String())
		assert.Equal(t, uint32(1), decoded.ValidFromTimestamp)
		assert.Equal(t, uint32(60), decoded.ExpiresAt)
		assert.Equal(t, big.NewInt(0).String(), decoded.LinkFee.String())
		assert.Equal(t, big.NewInt(0).String(), decoded.NativeFee.String())
	})

	t.Run("Decode errors on invalid report", func(t *testing.T) {
		_, err := rc.Decode([]byte{1, 2, 3})
		require.EqualError(t, err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")

		longBad := make([]byte, 64)
		for i := 0; i < len(longBad); i++ {
			longBad[i] = byte(i)
		}
		_, err = rc.Decode(longBad)
		require.EqualError(t, err, "failed to decode report: abi: improperly encoded uint32 value")
	})
}

type UnhandledStreamValue struct{}

var _ llo.StreamValue = &UnhandledStreamValue{}

func (sv *UnhandledStreamValue) MarshalBinary() (data []byte, err error) { return }
func (sv *UnhandledStreamValue) UnmarshalBinary(data []byte) error       { return nil }
func (sv *UnhandledStreamValue) MarshalText() (text []byte, err error)   { return }
func (sv *UnhandledStreamValue) UnmarshalText(text []byte) error         { return nil }
func (sv *UnhandledStreamValue) Type() llo.LLOStreamValue_Type           { return 0 }

func Test_ExtractReportValues(t *testing.T) {
	t.Run("with wrong number of stream values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36))}}
		_, _, _, err := ExtractReportValues(report)
		require.EqualError(t, err, "ReportCodecPremiumLegacy requires exactly 3 values (NativePrice, LinkPrice, Quote{Bid, Mid, Ask}); got report.Values: [35 36]")
	})
	t.Run("with (nil, nil, nil) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{nil, nil, nil}}
		_, _, _, err := ExtractReportValues(report)

		require.EqualError(t, err, "ReportCodecPremiumLegacy expects third stream value to be of type *Quote; got: <nil>")
	})
	t.Run("with ((*llo.Quote)(nil), nil, (*llo.Quote)(nil)) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{(*llo.Quote)(nil), nil, (*llo.Quote)(nil)}}
		nativePrice, linkPrice, quote, err := ExtractReportValues(report)

		require.EqualError(t, err, "ReportCodecPremiumLegacy expects third stream value to be non-nil")
		assert.Equal(t, decimal.Zero, nativePrice)
		assert.Equal(t, decimal.Zero, linkPrice)
		assert.Nil(t, quote)
	})
	t.Run("with (*llo.Decimal, *llo.Decimal, *llo.Decimal) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36)), llo.ToDecimal(decimal.NewFromInt(37))}}
		_, _, _, err := ExtractReportValues(report)

		require.EqualError(t, err, "ReportCodecPremiumLegacy expects third stream value to be of type *Quote; got: *llo.Decimal")
	})
	t.Run("with ((*llo.Quote)(nil), nil, *llo.Quote) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{(*llo.Quote)(nil), nil, &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}}}
		nativePrice, linkPrice, quote, err := ExtractReportValues(report)

		require.NoError(t, err)
		assert.Equal(t, decimal.Zero, nativePrice)
		assert.Equal(t, decimal.Zero, linkPrice)
		assert.Equal(t, &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}, quote)
	})
	t.Run("with unrecognized types", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{&UnhandledStreamValue{}, &UnhandledStreamValue{}, &UnhandledStreamValue{}}}
		_, _, _, err := ExtractReportValues(report)

		require.EqualError(t, err, "ReportCodecPremiumLegacy failed to extract native price: expected *Decimal or *Quote; got: *evm.UnhandledStreamValue")

		report = llo.Report{Values: []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), &UnhandledStreamValue{}, &UnhandledStreamValue{}}}
		_, _, _, err = ExtractReportValues(report)

		require.EqualError(t, err, "ReportCodecPremiumLegacy failed to extract link price: expected *Decimal or *Quote; got: *evm.UnhandledStreamValue")

		report = llo.Report{Values: []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36)), &UnhandledStreamValue{}}}
		_, _, _, err = ExtractReportValues(report)

		require.EqualError(t, err, "ReportCodecPremiumLegacy expects third stream value to be of type *Quote; got: *evm.UnhandledStreamValue")
	})
	t.Run("with (*llo.Decimal, *llo.Decimal, *llo.Quote) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36)), &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}}}
		nativePrice, linkPrice, quote, err := ExtractReportValues(report)

		require.NoError(t, err)
		assert.Equal(t, decimal.NewFromInt(35), nativePrice)
		assert.Equal(t, decimal.NewFromInt(36), linkPrice)
		assert.Equal(t, &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}, quote)
	})
	t.Run("with (*llo.Quote, *llo.Quote, *llo.Quote) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{&llo.Quote{Bid: decimal.NewFromInt(35), Benchmark: decimal.NewFromInt(36), Ask: decimal.NewFromInt(37)}, &llo.Quote{Bid: decimal.NewFromInt(38), Benchmark: decimal.NewFromInt(39), Ask: decimal.NewFromInt(40)}, &llo.Quote{Bid: decimal.NewFromInt(41), Benchmark: decimal.NewFromInt(42), Ask: decimal.NewFromInt(43)}}}
		nativePrice, linkPrice, quote, err := ExtractReportValues(report)

		require.NoError(t, err)
		assert.Equal(t, decimal.NewFromInt(36), nativePrice)
		assert.Equal(t, decimal.NewFromInt(39), linkPrice)
		assert.Equal(t, &llo.Quote{Bid: decimal.NewFromInt(41), Benchmark: decimal.NewFromInt(42), Ask: decimal.NewFromInt(43)}, quote)
	})
	t.Run("with (nil, nil, *llo.Quote) values", func(t *testing.T) {
		report := llo.Report{Values: []llo.StreamValue{nil, nil, &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}}}
		nativePrice, linkPrice, quote, err := ExtractReportValues(report)

		require.NoError(t, err)
		assert.Equal(t, decimal.Zero, nativePrice)
		assert.Equal(t, decimal.Zero, linkPrice)
		assert.Equal(t, &llo.Quote{Bid: decimal.NewFromInt(37), Benchmark: decimal.NewFromInt(38), Ask: decimal.NewFromInt(39)}, quote)
	})
}

func Test_LLOExtraHash(t *testing.T) {
	donID := uint32(8)
	extraHash := LLOExtraHash(donID)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000800000001", extraHash.String())
}

func Test_ReportCodecPremiumLegacy_Verify(t *testing.T) {
	c := ReportCodecPremiumLegacy{}
	t.Run("unrecognized fields in opts", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"unknown":"field"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown field")
	})
	t.Run("invalid opts", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`"invalid"`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Opts, got: \"\\\"invalid\\\"\"; json: cannot unmarshal string into Go value of type evm.ReportFormatEVMPremiumLegacyOpts")
	})
	t.Run("negative BaseUSDFee", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"baseUSDFee":"-1"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "baseUSDFee must be non-negative")
	})
	t.Run("zero feedID", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"feedID":"0x"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Opts, got: \"{\\\"feedID\\\":\\\"0x\\\"}\"; hex string has length 0, want 64 for common.Hash")
	})
	t.Run("incorrect number of streams", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"baseUSDFee":"1","feedID":"0x1111111111111111111111111111111111111111111111111111111111111111"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ReportFormatEVMPremiumLegacy requires exactly 3 streams (NativePrice, LinkPrice, Quote); got: []")
	})
	t.Run("invalid feedID", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"baseUSDFee":"1","feedID":"foo"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Opts, got: \"{\\\"baseUSDFee\\\":\\\"1\\\",\\\"feedID\\\":\\\"foo\\\"}\"; json: cannot unmarshal hex string without 0x prefix into Go struct field ReportFormatEVMPremiumLegacyOpts.feedID of type common.Hash")
	})
	t.Run("valid", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			Streams: []llotypes.Stream{
				{
					StreamID:   1,
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					StreamID:   2,
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					StreamID:   3,
					Aggregator: llotypes.AggregatorMedian,
				},
			},
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"baseUSDFee":"1","feedID":"0x1111111111111111111111111111111111111111111111111111111111111111"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.NoError(t, err)
	})
}
