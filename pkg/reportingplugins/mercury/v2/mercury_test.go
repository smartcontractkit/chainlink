package mercury_v2

import (
	"context"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
)

type testDataSource struct {
	Obs Observation
}

func (ds testDataSource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (Observation, error) {
	return ds.Obs, nil
}

type testReportCodec struct {
	observationTimestamp uint32
	validFromTimestamp   uint32
	builtReport          ocrtypes.Report
}

func (rc *testReportCodec) BuildReport(paos []ParsedAttributedObservation, f int, validFromTimestamp uint32, expiresAt uint32) (ocrtypes.Report, error) {
	rc.validFromTimestamp = validFromTimestamp
	return rc.builtReport, nil
}

func (rc testReportCodec) MaxReportLength(n int) (int, error) {
	return 456, nil
}

func (rc testReportCodec) ObservationTimestampFromReport(ocrtypes.Report) (uint32, error) {
	return rc.observationTimestamp, nil
}

func newAttributedObservation(t *testing.T, p *MercuryObservationProto) ocrtypes.AttributedObservation {
	marshalledObs, err := proto.Marshal(p)
	require.NoError(t, err)
	return ocrtypes.AttributedObservation{
		Observation: ocrtypes.Observation(marshalledObs),
		Observer:    commontypes.OracleID(42),
	}
}

func newTestReportPlugin(t *testing.T, codec *testReportCodec, ds *testDataSource) *reportingPlugin {
	offchainConfig := mercury.OffchainConfig{
		ExpirationWindow: 1,
		BaseUSDFeeCents:  100,
	}
	onchainConfig := mercury.OnchainConfig{
		Min: big.NewInt(1),
		Max: big.NewInt(1000),
	}
	maxReportLength, _ := codec.MaxReportLength(4)
	return &reportingPlugin{
		offchainConfig:           offchainConfig,
		onchainConfig:            onchainConfig,
		dataSource:               ds,
		logger:                   logger.Test(t),
		reportCodec:              codec,
		configDigest:             ocrtypes.ConfigDigest{},
		f:                        1,
		latestAcceptedEpochRound: mercury.EpochRound{},
		latestAcceptedMedian:     big.NewInt(0),
		maxReportLength:          maxReportLength,
	}
}

func Test_Plugin_Report(t *testing.T) {
	dataSource := &testDataSource{}
	codec := &testReportCodec{}
	rp := newTestReportPlugin(t, codec, dataSource)

	t.Run("when previous report is not nil", func(t *testing.T) {
		previousReport := ocrtypes.Report{}

		t.Run("reports if more than f+1 observations are valid", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}
			codec.observationTimestamp = 11
			codec.builtReport = ocrtypes.Report{1, 2, 3}

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, previousReport, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			assert.Equal(t, codec.validFromTimestamp, codec.observationTimestamp)
		})

		t.Run("reports if no f+1 maxFinalizedTimestamp observations available", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					PricesValid: false,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}
			codec.observationTimestamp = 22
			codec.builtReport = ocrtypes.Report{2, 3, 4}

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, previousReport, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			assert.Equal(t, codec.validFromTimestamp, codec.observationTimestamp)
		})

		t.Run("errors when less than f+1 valid observations available", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
			}

			_, _, err := rp.Report(ocrtypes.ReportTimestamp{}, previousReport, aos)
			assert.EqualError(t, err, "only received 1 valid attributed observations, but need at least f+1 (2)")
		})
	})

	t.Run("when previous report is nil", func(t *testing.T) {

		t.Run("reports if more than f+1 observations are valid", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      120,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}
			codec.builtReport = ocrtypes.Report{1, 2, 3}

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			assert.Equal(t, int(codec.validFromTimestamp), 40)
		})

		t.Run("errors when less than f+1 maxFinalizedTimestamp observations available", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: false,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}

			should, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.False(t, should)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid maxFinalizedTimestamp (got: 1/4)")
		})

		t.Run("errors when cannot come to consensus on MaxFinalizedTimestamp", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      41,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      42,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      43,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}

			should, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.False(t, should)
			assert.EqualError(t, err, "no valid maxFinalizedTimestamp with at least f+1 votes (got counts: map[40:1 41:1 42:1 43:1])")
		})

		t.Run("maxFinalizedTimestamp equals to observationTimestamp when consensus on MaxFinalizedTimestamp = 0", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 55,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      0,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 55,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      0,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 55,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      0,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 55,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      43,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}
			codec.builtReport = ocrtypes.Report{7, 8, 9}

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			assert.Equal(t, int(codec.validFromTimestamp), 55)
		})
	})

	t.Run("checkBenchmarkPrice", func(t *testing.T) {
		t.Run("checkBenchmarkPrice errors when fewer than f+1 observations have valid price", func(t *testing.T) {
			paos := []ParsedAttributedObservation{
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(123),
					PricesValid:    false,
				},
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(456),
					PricesValid:    false,
				},
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(789),
					PricesValid:    false,
				},
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(456),
					PricesValid:    true,
				},
			}

			err := rp.checkBenchmarkPrice(paos)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price (got: 1/4)")
		})

		t.Run("checkBenchmarkPrice errors when consensus benchmark price is outside of allowable range", func(t *testing.T) {
			paos := []ParsedAttributedObservation{
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(123123),
					PricesValid:    true,
				},
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(456456),
					PricesValid:    true,
				},
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(789789),
					PricesValid:    true,
				},
				parsedAttributedObservation{
					BenchmarkPrice: big.NewInt(123890),
					PricesValid:    true,
				},
			}

			err := rp.checkBenchmarkPrice(paos)
			assert.EqualError(t, err, "median benchmark price 456456 is outside of allowable range (Min: 1, Max: 1000)")
		})
	})

	t.Run("checkValidFromTimestamp errors when observationTimestamp < validFromTimestamp", func(t *testing.T) {
		err := rp.checkValidFromTimestamp(123, 456)
		assert.EqualError(t, err, "observationTimestamp (123) must be >= validFromTimestamp (456)")
	})

	t.Run("checkExpiresAt errors when expiresAt overflows", func(t *testing.T) {
		err := rp.checkExpiresAt(math.MaxUint32, math.MaxUint32)
		assert.EqualError(t, err, "timestamp 4294967295 + expiration window 4294967295 overflows uint32")
	})

	t.Run("buildReport failures", func(t *testing.T) {
		t.Run("Report errors when the report is too large", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}
			codec.builtReport = make([]byte, 1<<16)

			_, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)

			assert.EqualError(t, err, "report with len 65536 violates MaxReportLength limit set by ReportCodec (456)")
		})

		t.Run("Report errors when the report length is 0", func(t *testing.T) {
			aos := []ocrtypes.AttributedObservation{
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 42,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 45,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 47,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
					NativeFeeValid: true,
				}),
				newAttributedObservation(t, &MercuryObservationProto{
					Timestamp: 39,

					BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
					PricesValid:    true,

					MaxFinalizedTimestamp:      40,
					MaxFinalizedTimestampValid: true,

					LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
					LinkFeeValid:   true,
					NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
					NativeFeeValid: true,
				}),
			}
			codec.builtReport = []byte{}
			_, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)

			assert.EqualError(t, err, "report may not have zero length (invariant violation)")
		})
	})
}

func mustDecodeBigInt(b []byte) *big.Int {
	n, err := mercury.DecodeValueInt192(b)
	if err != nil {
		panic(err)
	}
	return n
}

func Test_Plugin_Observation(t *testing.T) {
	dataSource := &testDataSource{}
	codec := &testReportCodec{}
	rp := newTestReportPlugin(t, codec, dataSource)
	t.Run("Observation protobuf doesn't exceed maxObservationLength", func(t *testing.T) {
		obs := MercuryObservationProto{
			Timestamp:                  math.MaxUint32,
			BenchmarkPrice:             make([]byte, 24),
			PricesValid:                true,
			MaxFinalizedTimestamp:      math.MaxUint32,
			MaxFinalizedTimestampValid: true,
			LinkFee:                    make([]byte, 24),
			LinkFeeValid:               true,
			NativeFee:                  make([]byte, 24),
			NativeFeeValid:             true,
		}
		// This assertion is here to force this test to fail if a new field is
		// added to the protobuf. In this case, you must add the max value of
		// the field to the MercuryObservationProto in the test and only after
		// that increment the count below
		numFields := reflect.TypeOf(obs).NumField() //nolint:all
		// 3 fields internal to pbuf struct
		require.Equal(t, 9, numFields-3)

		b, err := proto.Marshal(&obs)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(b), maxObservationLength)
	})

	t.Run("all observations succeeded", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[uint32]{
				Val: rand.Uint32(),
			},
			LinkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
			NativePrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
		}
		dataSource.Obs = obs

		parsedObs, err := rp.Observation(context.Background(), ocrtypes.ReportTimestamp{}, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(parsedObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
		assert.True(t, p.PricesValid)
		assert.Equal(t, obs.MaxFinalizedTimestamp.Val, p.MaxFinalizedTimestamp)
		assert.True(t, p.MaxFinalizedTimestampValid)
		assert.Equal(t, mercury.CalculateFee(obs.LinkPrice.Val, 100), mustDecodeBigInt(p.LinkFee))
		assert.True(t, p.LinkFeeValid)
		assert.Equal(t, mercury.CalculateFee(obs.NativePrice.Val, 100), mustDecodeBigInt(p.NativeFee))
		assert.True(t, p.NativeFeeValid)
	})

	t.Run("some observations failed", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
				Err: errors.New("bechmarkPrice error"),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[uint32]{
				Val: rand.Uint32(),
				Err: errors.New("maxFinalizedTimestamp error"),
			},
			LinkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
				Err: errors.New("linkPrice error"),
			},
			NativePrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
		}

		dataSource.Obs = obs

		parsedObs, err := rp.Observation(context.Background(), ocrtypes.ReportTimestamp{}, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(parsedObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Zero(t, p.BenchmarkPrice)
		assert.False(t, p.PricesValid)
		assert.Zero(t, p.MaxFinalizedTimestamp)
		assert.False(t, p.MaxFinalizedTimestampValid)
		assert.Zero(t, p.LinkFee)
		assert.False(t, p.LinkFeeValid)
		assert.Equal(t, mercury.CalculateFee(obs.NativePrice.Val, 100), mustDecodeBigInt(p.NativeFee))
		assert.True(t, p.NativeFeeValid)
	})

	t.Run("all observations failed", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Err: errors.New("bechmarkPrice error"),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[uint32]{
				Err: errors.New("maxFinalizedTimestamp error"),
			},
			LinkPrice: mercury.ObsResult[*big.Int]{
				Err: errors.New("linkPrice error"),
			},
			NativePrice: mercury.ObsResult[*big.Int]{
				Err: errors.New("nativePrice error"),
			},
		}

		dataSource.Obs = obs

		parsedObs, err := rp.Observation(context.Background(), ocrtypes.ReportTimestamp{}, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(parsedObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Zero(t, p.BenchmarkPrice)
		assert.False(t, p.PricesValid)
		assert.Zero(t, p.MaxFinalizedTimestamp)
		assert.False(t, p.MaxFinalizedTimestampValid)
		assert.Zero(t, p.LinkFee)
		assert.False(t, p.LinkFeeValid)
		assert.Zero(t, p.NativeFee)
		assert.False(t, p.NativeFeeValid)
	})

	t.Run("encoding fails on some observations", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[uint32]{
				Val: rand.Uint32(),
			},
			LinkPrice: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			NativePrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
		}

		dataSource.Obs = obs

		parsedObs, err := rp.Observation(context.Background(), ocrtypes.ReportTimestamp{}, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(parsedObs, &p))

		assert.Zero(t, p.BenchmarkPrice)
		assert.False(t, p.PricesValid)
		assert.Zero(t, p.LinkFee)
		assert.False(t, p.LinkFeeValid)
	})

	t.Run("encoding fails on all observations", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[uint32]{
				Val: rand.Uint32(),
			},
			LinkPrice: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			NativePrice: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
		}

		dataSource.Obs = obs

		parsedObs, err := rp.Observation(context.Background(), ocrtypes.ReportTimestamp{}, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(parsedObs, &p))

		assert.Zero(t, p.BenchmarkPrice)
		assert.False(t, p.PricesValid)
		assert.Zero(t, p.LinkFee)
		assert.False(t, p.LinkFeeValid)
		assert.Zero(t, p.NativeFee)
		assert.False(t, p.NativeFeeValid)
	})
}
