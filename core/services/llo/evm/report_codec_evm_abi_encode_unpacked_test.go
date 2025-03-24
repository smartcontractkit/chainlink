package evm

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

// AllTrue returns false if any element in the array is false.
func AllTrue(arr []bool) bool {
	for _, v := range arr {
		if !v {
			return false
		}
	}
	return true
}

func TestReportFormatEVMABIEncodeOpts_Decode_Encode_properties(t *testing.T) {
	properties := gopter.NewProperties(nil)

	runTest := func(opts ReportFormatEVMABIEncodeOpts) bool {
		encoded, err := opts.Encode()
		require.NoError(t, err)

		decoded := ReportFormatEVMABIEncodeOpts{}
		err = decoded.Decode(encoded)
		require.NoError(t, err)

		return decoded.BaseUSDFee.Equal(opts.BaseUSDFee) && decoded.ExpirationWindow == opts.ExpirationWindow && decoded.FeedID == opts.FeedID && assert.Equal(t, opts.ABI, decoded.ABI)
	}
	properties.Property("Encodes values", prop.ForAll(
		runTest,
		gen.StrictStruct(reflect.TypeOf(&ReportFormatEVMABIEncodeOpts{}), map[string]gopter.Gen{
			"BaseUSDFee":       genBaseUSDFee(),
			"ExpirationWindow": genExpirationWindow(),
			"FeedID":           genFeedID(),
			"ABI":              genABI(),
		})))

	properties.TestingRun(t)
}

func genABI() gopter.Gen {
	return gen.SliceOf(genABIEncoder())
}

func genABIEncoder() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&ABIEncoder{}), map[string]gopter.Gen{
		"encoders": gen.SliceOf(genSingleABIEncoder()),
	})
}

func genSingleABIEncoder() gopter.Gen {
	return gen.StrictStruct(reflect.TypeOf(&singleABIEncoder{}), map[string]gopter.Gen{
		"Multiplier": genMultiplier(),
		"Type":       gen.AnyString(),
	})
}

func TestReportCodecEVMABIEncodeUnpacked_Encode_properties(t *testing.T) {
	ctx := tests.Context(t)
	codec := ReportCodecEVMABIEncodeUnpacked{}

	properties := gopter.NewProperties(nil)

	t.Run("DEX-based asset schema example", func(t *testing.T) {
		expectedDEXBasedAssetSchema := abi.Arguments([]abi.Argument{
			{Name: "feedId", Type: mustNewABIType("bytes32")},
			{Name: "validFromTimestamp", Type: mustNewABIType("uint32")},
			{Name: "observationsTimestamp", Type: mustNewABIType("uint32")},
			{Name: "nativeFee", Type: mustNewABIType("uint192")},
			{Name: "linkFee", Type: mustNewABIType("uint192")},
			{Name: "expiresAt", Type: mustNewABIType("uint32")},
			{Name: "price", Type: mustNewABIType("int192")},
			{Name: "baseMarketDepth", Type: mustNewABIType("int192")},
			{Name: "quoteMarketDepth", Type: mustNewABIType("int192")},
		})
		runTest := func(sampleFeedID common.Hash, sampleObservationTimestampNanoseconds, sampleValidAfterNanoseconds uint64, sampleExpirationWindow uint32, priceMultiplier, marketDepthMultiplier *ubig.Big, sampleBaseUSDFee, sampleLinkBenchmarkPrice, sampleNativeBenchmarkPrice, sampleDexBasedAssetPrice, sampleBaseMarketDepth, sampleQuoteMarketDepth decimal.Decimal) bool {
			report := llo.Report{
				ConfigDigest:                    types.ConfigDigest{0x01},
				SeqNr:                           0x02,
				ChannelID:                       llotypes.ChannelID(0x03),
				ValidAfterNanoseconds:           sampleValidAfterNanoseconds,
				ObservationTimestampNanoseconds: sampleObservationTimestampNanoseconds,
				Values: []llo.StreamValue{
					&llo.Quote{Bid: decimal.NewFromFloat(6.1), Benchmark: sampleLinkBenchmarkPrice, Ask: decimal.NewFromFloat(8.2332)},  // Link price
					&llo.Quote{Bid: decimal.NewFromFloat(9.4), Benchmark: sampleNativeBenchmarkPrice, Ask: decimal.NewFromFloat(11.33)}, // Native price
					llo.ToDecimal(sampleDexBasedAssetPrice), // DEX-based asset price
					llo.ToDecimal(sampleBaseMarketDepth),    // Base market depth
					llo.ToDecimal(sampleQuoteMarketDepth),   // Quote market depth
				},
				Specimen: false,
			}

			opts := ReportFormatEVMABIEncodeOpts{
				BaseUSDFee:       sampleBaseUSDFee,
				ExpirationWindow: sampleExpirationWindow,
				FeedID:           sampleFeedID,
				ABI: []ABIEncoder{
					// benchmark price
					newSingleABIEncoder("int192", priceMultiplier),
					// base market depth
					newSingleABIEncoder("int192", marketDepthMultiplier),
					// quote market depth
					newSingleABIEncoder("int192", marketDepthMultiplier),
				},
			}
			serializedOpts, err := opts.Encode()
			require.NoError(t, err)

			cd := llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
				Streams: []llotypes.Stream{
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorQuote,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: serializedOpts,
			}

			encoded, err := codec.Encode(ctx, report, cd)
			require.NoError(t, err)

			values, err := expectedDEXBasedAssetSchema.Unpack(encoded)
			require.NoError(t, err)

			require.Len(t, values, len(expectedDEXBasedAssetSchema))

			expectedLinkFee := CalculateFee(sampleLinkBenchmarkPrice, sampleBaseUSDFee)
			expectedNativeFee := CalculateFee(sampleNativeBenchmarkPrice, sampleBaseUSDFee)

			// doesn't crash if values are nil
			for i := range report.Values {
				report.Values[i] = nil
			}
			_, err = codec.Encode(ctx, report, cd)
			require.Error(t, err)

			return AllTrue([]bool{
				assert.Equal(t, sampleFeedID, (common.Hash)(values[0].([32]byte))),                                                                  //nolint:testifylint // false positive // feedId
				assert.Equal(t, uint32(sampleValidAfterNanoseconds/1e9)+1, values[1].(uint32)),                                                      //nolint:gosec // G115 // validFromTimestamp
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9), values[2].(uint32)),                                              //nolint:gosec // G115 // observationsTimestamp
				assert.Equal(t, expectedLinkFee.String(), values[3].(*big.Int).String()),                                                            // linkFee
				assert.Equal(t, expectedNativeFee.String(), values[4].(*big.Int).String()),                                                          // nativeFee
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9)+sampleExpirationWindow, values[5].(uint32)),                       //nolint:gosec // G115 generator ensures it wont overflow
				assert.Equal(t, sampleDexBasedAssetPrice.Mul(decimal.NewFromBigInt(priceMultiplier.ToInt(), 0)).BigInt(), values[6].(*big.Int)),     // price
				assert.Equal(t, sampleBaseMarketDepth.Mul(decimal.NewFromBigInt(marketDepthMultiplier.ToInt(), 0)).BigInt(), values[7].(*big.Int)),  // baseMarketDepth
				assert.Equal(t, sampleQuoteMarketDepth.Mul(decimal.NewFromBigInt(marketDepthMultiplier.ToInt(), 0)).BigInt(), values[8].(*big.Int)), // quoteMarketDepth
			})
		}

		properties.Property("Encodes values", prop.ForAll(
			runTest,
			genFeedID(),
			genObservationTimestampNanoseconds(),
			genValidAfterNanoseconds(),
			genExpirationWindow(),
			genPriceMultiplier(),
			genMarketDepthMultiplier(),
			genBaseUSDFee(),
			genLinkBenchmarkPrice(),
			genNativeBenchmarkPrice(),
			genBenchmarkPrice(),
			genBaseMarketDepth(),
			genQuoteMarketDepth(),
		))

		properties.TestingRun(t)
	})

	t.Run("Market status schema", func(t *testing.T) {
		expectedRWASchema := abi.Arguments([]abi.Argument{
			{Name: "feedId", Type: mustNewABIType("bytes32")},
			{Name: "validFromTimestamp", Type: mustNewABIType("uint32")},
			{Name: "observationsTimestamp", Type: mustNewABIType("uint32")},
			{Name: "nativeFee", Type: mustNewABIType("uint192")},
			{Name: "linkFee", Type: mustNewABIType("uint192")},
			{Name: "expiresAt", Type: mustNewABIType("uint32")},
			{Name: "marketStatus", Type: mustNewABIType("uint32")},
		})

		runTest := func(sampleFeedID common.Hash, sampleObservationTimestampNanoseconds, sampleValidAfterNanoseconds uint64, sampleExpirationWindow uint32, sampleBaseUSDFee, sampleLinkBenchmarkPrice, sampleNativeBenchmarkPrice, sampleMarketStatus decimal.Decimal) bool {
			report := llo.Report{
				ConfigDigest:                    types.ConfigDigest{0x01},
				SeqNr:                           0x02,
				ChannelID:                       llotypes.ChannelID(0x03),
				ValidAfterNanoseconds:           sampleValidAfterNanoseconds,
				ObservationTimestampNanoseconds: sampleObservationTimestampNanoseconds,
				Values: []llo.StreamValue{
					&llo.Quote{Bid: decimal.NewFromFloat(6.1), Benchmark: sampleLinkBenchmarkPrice, Ask: decimal.NewFromFloat(8.2332)},  // Link price
					&llo.Quote{Bid: decimal.NewFromFloat(9.4), Benchmark: sampleNativeBenchmarkPrice, Ask: decimal.NewFromFloat(11.33)}, // Native price
					llo.ToDecimal(sampleMarketStatus), // DEX-based asset price
				},
				Specimen: false,
			}

			opts := ReportFormatEVMABIEncodeOpts{
				BaseUSDFee:       sampleBaseUSDFee,
				ExpirationWindow: sampleExpirationWindow,
				FeedID:           sampleFeedID,
				ABI: []ABIEncoder{
					// market status
					newSingleABIEncoder("uint32", nil),
				},
			}
			serializedOpts, err := opts.Encode()
			require.NoError(t, err)

			cd := llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
				Streams: []llotypes.Stream{
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorQuote,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: serializedOpts,
			}

			encoded, err := codec.Encode(ctx, report, cd)
			require.NoError(t, err)

			values, err := expectedRWASchema.Unpack(encoded)
			require.NoError(t, err)

			require.Len(t, values, len(expectedRWASchema))

			expectedLinkFee := CalculateFee(sampleLinkBenchmarkPrice, sampleBaseUSDFee)
			expectedNativeFee := CalculateFee(sampleNativeBenchmarkPrice, sampleBaseUSDFee)

			// doesn't crash if values are nil
			for i := range report.Values {
				report.Values[i] = nil
			}
			_, err = codec.Encode(ctx, report, cd)
			require.Error(t, err)

			return AllTrue([]bool{
				assert.Equal(t, sampleFeedID, (common.Hash)(values[0].([32]byte))),                                            //nolint:testifylint // false positive // feedId
				assert.Equal(t, uint32(sampleValidAfterNanoseconds/1e9)+1, values[1].(uint32)),                                //nolint:gosec // G115 // validFromTimestamp
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9), values[2].(uint32)),                        //nolint:gosec // G115 //  observationsTimestamp
				assert.Equal(t, expectedLinkFee.String(), values[3].(*big.Int).String()),                                      // linkFee
				assert.Equal(t, expectedNativeFee.String(), values[4].(*big.Int).String()),                                    // nativeFee
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9)+sampleExpirationWindow, values[5].(uint32)), //nolint:gosec // G115 // expiresAt
				assert.Equal(t, uint32(sampleMarketStatus.BigInt().Int64()), values[6].(uint32)),                              //nolint:gosec // G115 // market status
			})
		}

		properties.Property("Encodes values", prop.ForAll(
			runTest,
			genFeedID(),
			genObservationTimestampNanoseconds(),
			genValidAfterNanoseconds(),
			genExpirationWindow(),
			genBaseUSDFee(),
			genLinkBenchmarkPrice(),
			genNativeBenchmarkPrice(),
			genMarketStatus(),
		))

		properties.TestingRun(t)
	})

	t.Run("benchmark price schema example", func(t *testing.T) {
		expectedDEXBasedAssetSchema := abi.Arguments([]abi.Argument{
			{Name: "feedId", Type: mustNewABIType("bytes32")},
			{Name: "validFromTimestamp", Type: mustNewABIType("uint32")},
			{Name: "observationsTimestamp", Type: mustNewABIType("uint32")},
			{Name: "nativeFee", Type: mustNewABIType("uint192")},
			{Name: "linkFee", Type: mustNewABIType("uint192")},
			{Name: "expiresAt", Type: mustNewABIType("uint32")},
			{Name: "price", Type: mustNewABIType("int192")},
		})
		runTest := func(sampleFeedID common.Hash, sampleObservationTimestampNanoseconds, sampleValidAfterNanoseconds uint64, sampleExpirationWindow uint32, priceMultiplier, marketDepthMultiplier *ubig.Big, sampleBaseUSDFee, sampleLinkBenchmarkPrice, sampleNativeBenchmarkPrice, sampleBenchmarkPrice decimal.Decimal) bool {
			report := llo.Report{
				ConfigDigest:                    types.ConfigDigest{0x01},
				SeqNr:                           0x02,
				ChannelID:                       llotypes.ChannelID(0x03),
				ValidAfterNanoseconds:           sampleValidAfterNanoseconds,
				ObservationTimestampNanoseconds: sampleObservationTimestampNanoseconds,
				Values: []llo.StreamValue{
					&llo.Quote{Bid: decimal.NewFromFloat(6.1), Benchmark: sampleLinkBenchmarkPrice, Ask: decimal.NewFromFloat(8.2332)},  // Link price
					&llo.Quote{Bid: decimal.NewFromFloat(9.4), Benchmark: sampleNativeBenchmarkPrice, Ask: decimal.NewFromFloat(11.33)}, // Native price
					llo.ToDecimal(sampleBenchmarkPrice), // price
				},
				Specimen: false,
			}

			opts := ReportFormatEVMABIEncodeOpts{
				BaseUSDFee:       sampleBaseUSDFee,
				ExpirationWindow: sampleExpirationWindow,
				FeedID:           sampleFeedID,
				ABI: []ABIEncoder{
					// benchmark price
					newSingleABIEncoder("int192", priceMultiplier),
				},
			}
			serializedOpts, err := opts.Encode()
			require.NoError(t, err)

			cd := llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
				Streams: []llotypes.Stream{
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorQuote,
					},
				},
				Opts: serializedOpts,
			}

			encoded, err := codec.Encode(ctx, report, cd)
			require.NoError(t, err)

			values, err := expectedDEXBasedAssetSchema.Unpack(encoded)
			require.NoError(t, err)

			require.Len(t, values, len(expectedDEXBasedAssetSchema))

			expectedLinkFee := CalculateFee(sampleLinkBenchmarkPrice, sampleBaseUSDFee)
			expectedNativeFee := CalculateFee(sampleNativeBenchmarkPrice, sampleBaseUSDFee)

			// doesn't crash if values are nil
			for i := range report.Values {
				report.Values[i] = nil
			}
			_, err = codec.Encode(ctx, report, cd)
			require.Error(t, err)

			return AllTrue([]bool{
				assert.Equal(t, sampleFeedID, (common.Hash)(values[0].([32]byte))),                                                          //nolint:testifylint // false positive // feedId
				assert.Equal(t, uint32(sampleValidAfterNanoseconds/1e9)+1, values[1].(uint32)),                                              //nolint:gosec // G115 // validFromTimestamp
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9), values[2].(uint32)),                                      //nolint:gosec // G115 //  observationsTimestamp
				assert.Equal(t, expectedLinkFee.String(), values[3].(*big.Int).String()),                                                    // linkFee
				assert.Equal(t, expectedNativeFee.String(), values[4].(*big.Int).String()),                                                  // nativeFee
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9)+sampleExpirationWindow, values[5].(uint32)),               //nolint:gosec // G115 // expiresAt
				assert.Equal(t, sampleBenchmarkPrice.Mul(decimal.NewFromBigInt(priceMultiplier.ToInt(), 0)).BigInt(), values[6].(*big.Int)), // price
			})
		}

		properties.Property("Encodes values", prop.ForAll(
			runTest,
			genFeedID(),
			genObservationTimestampNanoseconds(),
			genValidAfterNanoseconds(),
			genExpirationWindow(),
			genPriceMultiplier(),
			genMarketDepthMultiplier(),
			genBaseUSDFee(),
			genLinkBenchmarkPrice(),
			genNativeBenchmarkPrice(),
			genBenchmarkPrice(),
		))

		properties.TestingRun(t)
	})

	t.Run("funding rate schema example", func(t *testing.T) {
		expectedFundingRateSchema := abi.Arguments([]abi.Argument{
			{Name: "feedId", Type: mustNewABIType("bytes32")},
			{Name: "validFromTimestamp", Type: mustNewABIType("uint32")},
			{Name: "observationsTimestamp", Type: mustNewABIType("uint32")},
			{Name: "nativeFee", Type: mustNewABIType("uint192")},
			{Name: "linkFee", Type: mustNewABIType("uint192")},
			{Name: "expiresAt", Type: mustNewABIType("uint32")},
			{Name: "binanceFundingRate", Type: mustNewABIType("int192")},
			{Name: "binanceFundingTime", Type: mustNewABIType("uint32")},
			{Name: "binanceFundingIntervalHours", Type: mustNewABIType("uint32")},
			{Name: "deribitFundingRate", Type: mustNewABIType("int192")},
			{Name: "deribitFundingTime", Type: mustNewABIType("uint32")},
			{Name: "deribitFundingIntervalHours", Type: mustNewABIType("uint32")},
		})

		runTest := func(sampleFeedID common.Hash, sampleObservationTimestampNanoseconds, sampleValidAfterNanoseconds uint64, sampleExpirationWindow uint32, sampleBaseUSDFee, sampleLinkBenchmarkPrice, sampleNativeBenchmarkPrice, sampleBinanceFundingRate, sampleBinanceFundingTime, sampleBinanceFundingIntervalHours, sampleDeribitFundingRate, sampleDeribitFundingTime, sampleDeribitFundingIntervalHours decimal.Decimal) bool {
			report := llo.Report{
				ConfigDigest:                    types.ConfigDigest{0x01},
				SeqNr:                           0x02,
				ChannelID:                       llotypes.ChannelID(0x03),
				ValidAfterNanoseconds:           sampleValidAfterNanoseconds,
				ObservationTimestampNanoseconds: sampleObservationTimestampNanoseconds,
				Values: []llo.StreamValue{
					&llo.Quote{Bid: decimal.NewFromFloat(6.1), Benchmark: sampleLinkBenchmarkPrice, Ask: decimal.NewFromFloat(8.2332)},  // Link price
					&llo.Quote{Bid: decimal.NewFromFloat(9.4), Benchmark: sampleNativeBenchmarkPrice, Ask: decimal.NewFromFloat(11.33)}, // Native price
					llo.ToDecimal(sampleBinanceFundingRate),          // Binance funding rate
					llo.ToDecimal(sampleBinanceFundingTime),          // Binance funding time
					llo.ToDecimal(sampleBinanceFundingIntervalHours), // Binance funding interval hours
					llo.ToDecimal(sampleDeribitFundingRate),          // Deribit funding rate
					llo.ToDecimal(sampleDeribitFundingTime),          // Deribit funding time
					llo.ToDecimal(sampleDeribitFundingIntervalHours), // Deribit funding interval hours
				},
				Specimen: false,
			}

			opts := ReportFormatEVMABIEncodeOpts{
				BaseUSDFee:       sampleBaseUSDFee,
				ExpirationWindow: sampleExpirationWindow,
				FeedID:           sampleFeedID,
				ABI: []ABIEncoder{
					newSingleABIEncoder("int192", nil),
					newSingleABIEncoder("uint32", nil),
					newSingleABIEncoder("uint32", nil),
					newSingleABIEncoder("int192", nil),
					newSingleABIEncoder("uint32", nil),
					newSingleABIEncoder("uint32", nil),
				},
			}
			serializedOpts, err := opts.Encode()
			require.NoError(t, err)

			cd := llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
				Streams: []llotypes.Stream{
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorQuote,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: serializedOpts,
			}

			encoded, err := codec.Encode(ctx, report, cd)
			require.NoError(t, err)

			values, err := expectedFundingRateSchema.Unpack(encoded)
			require.NoError(t, err)

			require.Len(t, values, len(expectedFundingRateSchema))

			expectedLinkFee := CalculateFee(sampleLinkBenchmarkPrice, sampleBaseUSDFee)
			expectedNativeFee := CalculateFee(sampleNativeBenchmarkPrice, sampleBaseUSDFee)

			// doesn't crash if values are nil
			for i := range report.Values {
				report.Values[i] = nil
			}
			_, err = codec.Encode(ctx, report, cd)
			require.Error(t, err)

			return AllTrue([]bool{
				assert.Equal(t, sampleFeedID, (common.Hash)(values[0].([32]byte))),                                            //nolint:testifylint // false positive // feedId
				assert.Equal(t, uint32(sampleValidAfterNanoseconds/1e9)+1, values[1].(uint32)),                                //nolint:gosec // G115 // validFromTimestamp
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9), values[2].(uint32)),                        //nolint:gosec // G115 //  observationsTimestamp
				assert.Equal(t, expectedLinkFee.String(), values[3].(*big.Int).String()),                                      // linkFee
				assert.Equal(t, expectedNativeFee.String(), values[4].(*big.Int).String()),                                    // nativeFee
				assert.Equal(t, uint32(sampleObservationTimestampNanoseconds/1e9)+sampleExpirationWindow, values[5].(uint32)), //nolint:gosec // G115 // expiresAt
				assert.Equal(t, sampleBinanceFundingRate.BigInt().String(), values[6].(*big.Int).String()),                    // binanceFundingRate
				assert.Equal(t, uint32(sampleBinanceFundingTime.BigInt().Int64()), values[7].(uint32)),                        //nolint:gosec // G115 // binanceFundingTime
				assert.Equal(t, uint32(sampleBinanceFundingIntervalHours.BigInt().Int64()), values[8].(uint32)),               //nolint:gosec // G115 // binanceFundingIntervalHours
				assert.Equal(t, sampleDeribitFundingRate.BigInt().String(), values[9].(*big.Int).String()),                    // deribitFundingRate
				assert.Equal(t, uint32(sampleDeribitFundingTime.BigInt().Int64()), values[10].(uint32)),                       //nolint:gosec // G115 // deribitFundingTime
				assert.Equal(t, uint32(sampleDeribitFundingIntervalHours.BigInt().Int64()), values[11].(uint32)),              //nolint:gosec // G115 // deribitFundingIntervalHours
			})
		}

		properties.Property("Encodes values", prop.ForAll(
			runTest,
			genFeedID(),
			genObservationTimestampNanoseconds(),
			genValidAfterNanoseconds(),
			genExpirationWindow(),
			genBaseUSDFee(),
			genLinkBenchmarkPrice(),
			genNativeBenchmarkPrice(),
			genFundingRate(),
			genFundingTime(),
			genFundingIntervalHours(),
			genFundingRate(),
			genFundingTime(),
			genFundingIntervalHours(),
		))

		properties.TestingRun(t)
	})
}

func TestReportCodecEVMABIEncodeUnpacked_Encode(t *testing.T) {
	t.Run("ABI and values length mismatch error", func(t *testing.T) {
		report := llo.Report{
			ConfigDigest:                    types.ConfigDigest{0x01},
			SeqNr:                           0x02,
			ChannelID:                       llotypes.ChannelID(0x03),
			ValidAfterNanoseconds:           0x04,
			ObservationTimestampNanoseconds: 0x05,
			Values: []llo.StreamValue{
				&llo.Quote{Bid: decimal.NewFromFloat(6.1), Benchmark: decimal.NewFromFloat(7.4), Ask: decimal.NewFromFloat(8.2332)},
				&llo.Quote{Bid: decimal.NewFromFloat(9.4), Benchmark: decimal.NewFromFloat(10.0), Ask: decimal.NewFromFloat(11.33)},
				llo.ToDecimal(decimal.NewFromFloat(100)),
				llo.ToDecimal(decimal.NewFromFloat(101)),
				llo.ToDecimal(decimal.NewFromFloat(102)),
			},
			Specimen: false,
		}

		opts := ReportFormatEVMABIEncodeOpts{
			ABI: []ABIEncoder{},
		}
		serializedOpts, err := opts.Encode()
		require.NoError(t, err)
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Streams: []llotypes.Stream{
				{
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					Aggregator: llotypes.AggregatorQuote,
				},
				{
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					Aggregator: llotypes.AggregatorMedian,
				},
			},
			Opts: serializedOpts,
		}

		codec := ReportCodecEVMABIEncodeUnpacked{}
		_, err = codec.Encode(tests.Context(t), report, cd)
		require.EqualError(t, err, "failed to build payload; ABI and values length mismatch; ABI: 0, Values: 3")

		report.Values = []llo.StreamValue{}
		_, err = codec.Encode(tests.Context(t), report, cd)
		require.EqualError(t, err, "ReportCodecEVMABIEncodeUnpacked requires at least 2 values (NativePrice, LinkPrice, ...); got report.Values: []")
	})
}

func genFeedID() gopter.Gen {
	return func(p *gopter.GenParameters) *gopter.GenResult {
		var feedID common.Hash
		p.Rng.Read(feedID[:])
		return gopter.NewGenResult(feedID, gopter.NoShrinker)
	}
}

func genObservationTimestampNanoseconds() gopter.Gen {
	return genTimestampThatFitsUint32Seconds()
}

func genValidAfterNanoseconds() gopter.Gen {
	return genTimestampThatFitsUint32Seconds()
}

func genTimestampThatFitsUint32Seconds() gopter.Gen {
	return gen.UInt32().Map(func(i uint32) uint64 {
		return uint64(i) * 1e9
	})
}

func genExpirationWindow() gopter.Gen {
	return gen.UInt32()
}

func genPriceMultiplier() gopter.Gen {
	return genMultiplier()
}

func genMarketDepthMultiplier() gopter.Gen {
	return genMultiplier()
}

func genMultiplier() gopter.Gen {
	return gen.UInt32().Map(func(i uint32) *ubig.Big {
		return ubig.NewI(int64(i))
	})
}

func genDecimal() gopter.Gen {
	return gen.Float32Range(-2e32, 2e32).Map(decimal.NewFromFloat32)
}

func genBaseUSDFee() gopter.Gen {
	return genDecimal()
}

func genLinkBenchmarkPrice() gopter.Gen {
	return genDecimal()
}

func genNativeBenchmarkPrice() gopter.Gen {
	return genDecimal()
}

func genBenchmarkPrice() gopter.Gen {
	return genDecimal()
}

func genBaseMarketDepth() gopter.Gen {
	return genDecimal()
}

func genQuoteMarketDepth() gopter.Gen {
	return genDecimal()
}

func genMarketStatus() gopter.Gen {
	return gen.UInt32().Map(func(i uint32) decimal.Decimal {
		return decimal.NewFromInt(int64(i))
	})
}

func genFundingRate() gopter.Gen {
	return genDecimal()
}

func genFundingTime() gopter.Gen {
	// Unix epochs
	return gen.UInt32Range(1500000000, 2000000000).Map(func(i uint32) decimal.Decimal {
		return decimal.NewFromInt(int64(i))
	})
}

func genFundingIntervalHours() gopter.Gen {
	return gen.UInt32().Map(func(i uint32) decimal.Decimal {
		return decimal.NewFromInt(int64(i))
	})
}

func mustNewABIType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}

func TestReportCodecEVMABIEncodeUnpacked_Verify(t *testing.T) {
	c := ReportCodecEVMABIEncodeUnpacked{}
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
		assert.EqualError(t, err, "invalid Opts, got: \"\\\"invalid\\\"\"; json: cannot unmarshal string into Go value of type evm.ReportFormatEVMABIEncodeOpts")
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
			Opts:         []byte(`{"feedID":"0x0000000000000000000000000000000000000000000000000000000000000000"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "feedID must not be zero")
	})
	t.Run("missing feedID", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "feedID must not be zero")
	})
	t.Run("not enough streams", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Streams: []llotypes.Stream{
				{StreamID: 1},
				{StreamID: 2},
			},
			Opts: []byte(`{"ABI":[{"type":"int192"}],"feedID":"0x1111111111111111111111111111111111111111111111111111111111111111"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected at least 3 streams; got: 2")
	})
	t.Run("ABI length does not match streams length", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Streams: []llotypes.Stream{
				{StreamID: 1},
				{StreamID: 2},
				{StreamID: 3},
				{StreamID: 4},
			},
			Opts: []byte(`{"ABI":[{"type":"int192"}],"feedID":"0x1111111111111111111111111111111111111111111111111111111111111111"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ABI length mismatch; expected: 2, got: 1")
	})
	t.Run("invalid feedID", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"baseUSDFee":"1","feedID":"0x"}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.Error(t, err)
		assert.EqualError(t, err, "invalid Opts, got: \"{\\\"baseUSDFee\\\":\\\"1\\\",\\\"feedID\\\":\\\"0x\\\"}\"; hex string has length 0, want 64 for common.Hash")
	})
	t.Run("valid", func(t *testing.T) {
		cd := llotypes.ChannelDefinition{
			Streams: []llotypes.Stream{
				{StreamID: 1},
				{StreamID: 2},
				{StreamID: 3},
			},
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
			Opts:         []byte(`{"baseUSDFee":"1","feedID":"0x1111111111111111111111111111111111111111111111111111111111111111","ABI":[{"streamID":1,"type":"int192"}]}`),
		}
		err := c.Verify(tests.Context(t), cd)
		require.NoError(t, err)
	})
}

func padLeft32Byte(str string) string {
	if len(str) >= 64 {
		return str
	}
	padding := strings.Repeat("0", 64-len(str))
	return padding + str
}

func newSingleABIEncoder(t string, m *ubig.Big) ABIEncoder {
	return ABIEncoder{
		[]singleABIEncoder{{
			Type:       t,
			Multiplier: m,
		}},
	}
}
