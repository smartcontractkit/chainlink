package mercury_v3 //nolint:revive

import (
	"context"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury"
)

type testDataSource struct {
	Obs Observation
}

func (ds testDataSource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (Observation, error) {
	return ds.Obs, nil
}

type testReportCodec struct {
	observationTimestamp uint32
	builtReport          ocrtypes.Report

	builtReportFields *ReportFields
	err               error
}

func (rc *testReportCodec) BuildReport(rf ReportFields) (ocrtypes.Report, error) {
	rc.builtReportFields = &rf

	return rc.builtReport, nil
}

func (rc testReportCodec) MaxReportLength(n int) (int, error) {
	return 123, nil
}

func (rc testReportCodec) ObservationTimestampFromReport(ocrtypes.Report) (uint32, error) {
	return rc.observationTimestamp, rc.err
}

func newTestReportPlugin(t *testing.T, codec *testReportCodec, ds *testDataSource) *reportingPlugin {
	offchainConfig := mercury.OffchainConfig{
		ExpirationWindow: 1,
		BaseUSDFee:       decimal.NewFromInt32(1),
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

func newValidProtos() []*MercuryObservationProto {
	return []*MercuryObservationProto{
		&MercuryObservationProto{
			Timestamp: 42,

			BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(123)),
			Bid:            mercury.MustEncodeValueInt192(big.NewInt(120)),
			Ask:            mercury.MustEncodeValueInt192(big.NewInt(130)),
			PricesValid:    true,

			MaxFinalizedTimestamp:      40,
			MaxFinalizedTimestampValid: true,

			LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.1e18)),
			LinkFeeValid:   true,
			NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.1e18)),
			NativeFeeValid: true,
		},
		&MercuryObservationProto{
			Timestamp: 45,

			BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(234)),
			Bid:            mercury.MustEncodeValueInt192(big.NewInt(230)),
			Ask:            mercury.MustEncodeValueInt192(big.NewInt(240)),
			PricesValid:    true,

			MaxFinalizedTimestamp:      40,
			MaxFinalizedTimestampValid: true,

			LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.2e18)),
			LinkFeeValid:   true,
			NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.2e18)),
			NativeFeeValid: true,
		},
		&MercuryObservationProto{
			Timestamp: 47,

			BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(345)),
			Bid:            mercury.MustEncodeValueInt192(big.NewInt(340)),
			Ask:            mercury.MustEncodeValueInt192(big.NewInt(350)),
			PricesValid:    true,

			MaxFinalizedTimestamp:      39,
			MaxFinalizedTimestampValid: true,

			LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.3e18)),
			LinkFeeValid:   true,
			NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.3e18)),
			NativeFeeValid: true,
		},
		&MercuryObservationProto{
			Timestamp: 39,

			BenchmarkPrice: mercury.MustEncodeValueInt192(big.NewInt(456)),
			Bid:            mercury.MustEncodeValueInt192(big.NewInt(450)),
			Ask:            mercury.MustEncodeValueInt192(big.NewInt(460)),
			PricesValid:    true,

			MaxFinalizedTimestamp:      39,
			MaxFinalizedTimestampValid: true,

			LinkFee:        mercury.MustEncodeValueInt192(big.NewInt(1.4e18)),
			LinkFeeValid:   true,
			NativeFee:      mercury.MustEncodeValueInt192(big.NewInt(2.4e18)),
			NativeFeeValid: true,
		},
	}
}

func newValidAos(t *testing.T, protos ...*MercuryObservationProto) (aos []ocrtypes.AttributedObservation) {
	if len(protos) == 0 {
		protos = newValidProtos()
	}
	aos = make([]ocrtypes.AttributedObservation, len(protos))
	for i := range aos {
		marshalledObs, err := proto.Marshal(protos[i])
		require.NoError(t, err)
		aos[i] = ocrtypes.AttributedObservation{
			Observation: marshalledObs,
			Observer:    commontypes.OracleID(i),
		}
	}
	return
}

func Test_Plugin_Report(t *testing.T) {
	dataSource := &testDataSource{}
	codec := &testReportCodec{
		builtReport: []byte{1, 2, 3, 4},
	}
	rp := newTestReportPlugin(t, codec, dataSource)
	repts := ocrtypes.ReportTimestamp{}

	t.Run("when previous report is nil", func(t *testing.T) {
		t.Run("errors if not enough attributed observations", func(t *testing.T) {
			_, _, err := rp.Report(repts, nil, newValidAos(t)[0:1])
			assert.EqualError(t, err, "only received 1 valid attributed observations, but need at least f+1 (2)")
		})

		t.Run("errors if too many maxFinalizedTimestamp observations are invalid", func(t *testing.T) {
			ps := newValidProtos()
			ps[0].MaxFinalizedTimestampValid = false
			ps[1].MaxFinalizedTimestampValid = false
			ps[2].MaxFinalizedTimestampValid = false
			aos := newValidAos(t, ps...)

			should, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.False(t, should)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid maxFinalizedTimestamp (got: 1/4)")
		})
		t.Run("errors if maxFinalizedTimestamp is too large", func(t *testing.T) {
			ps := newValidProtos()
			ps[0].MaxFinalizedTimestamp = math.MaxUint32
			ps[1].MaxFinalizedTimestamp = math.MaxUint32
			ps[2].MaxFinalizedTimestamp = math.MaxUint32
			ps[3].MaxFinalizedTimestamp = math.MaxUint32
			aos := newValidAos(t, ps...)

			should, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.False(t, should)
			assert.EqualError(t, err, "maxFinalizedTimestamp is too large, got: 4294967295")
		})

		t.Run("succeeds and generates validFromTimestamp from maxFinalizedTimestamp when maxFinalizedTimestamp is positive", func(t *testing.T) {
			aos := newValidAos(t)

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, ReportFields{
				ValidFromTimestamp: 41, // consensus maxFinalizedTimestamp is 40, so validFrom should be 40+1
				Timestamp:          45,
				NativeFee:          big.NewInt(2300000000000000000), // 2.3e18
				LinkFee:            big.NewInt(1300000000000000000), // 1.3e18
				ExpiresAt:          46,
				BenchmarkPrice:     big.NewInt(345),
				Bid:                big.NewInt(340),
				Ask:                big.NewInt(350),
			}, *codec.builtReportFields)
		})
		t.Run("succeeds and generates validFromTimestamp from maxFinalizedTimestamp when maxFinalizedTimestamp is zero", func(t *testing.T) {
			protos := newValidProtos()
			for i := range protos {
				protos[i].MaxFinalizedTimestamp = 0
			}
			aos := newValidAos(t, protos...)

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, ReportFields{
				ValidFromTimestamp: 1,
				Timestamp:          45,
				NativeFee:          big.NewInt(2300000000000000000), // 2.3e18
				LinkFee:            big.NewInt(1300000000000000000), // 1.3e18
				ExpiresAt:          46,
				BenchmarkPrice:     big.NewInt(345),
				Bid:                big.NewInt(340),
				Ask:                big.NewInt(350),
			}, *codec.builtReportFields)
		})
		t.Run("succeeds and generates validFromTimestamp from maxFinalizedTimestamp when maxFinalizedTimestamp is -1 (missing feed)", func(t *testing.T) {
			protos := newValidProtos()
			for i := range protos {
				protos[i].MaxFinalizedTimestamp = -1
			}
			aos := newValidAos(t, protos...)

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)
			assert.True(t, should)
			assert.NoError(t, err)
			assert.Equal(t, codec.builtReport, report)
			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, ReportFields{
				ValidFromTimestamp: 45, // in case of missing feed, ValidFromTimestamp=Timestamp for first report
				Timestamp:          45,
				NativeFee:          big.NewInt(2300000000000000000), // 2.3e18
				LinkFee:            big.NewInt(1300000000000000000), // 1.3e18
				ExpiresAt:          46,
				BenchmarkPrice:     big.NewInt(345),
				Bid:                big.NewInt(340),
				Ask:                big.NewInt(350),
			}, *codec.builtReportFields)
		})

		t.Run("succeeds, ignoring unparseable attributed observation", func(t *testing.T) {
			aos := newValidAos(t)
			aos[0] = newUnparseableAttributedObservation()

			should, report, err := rp.Report(repts, nil, aos)
			require.NoError(t, err)

			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, ReportFields{
				ValidFromTimestamp: 40, // consensus maxFinalizedTimestamp is 39, so validFrom should be 39+1
				Timestamp:          45,
				NativeFee:          big.NewInt(2300000000000000000), // 2.3e18
				LinkFee:            big.NewInt(1300000000000000000), // 1.3e18
				ExpiresAt:          46,
				BenchmarkPrice:     big.NewInt(345),
				Bid:                big.NewInt(340),
				Ask:                big.NewInt(350),
			}, *codec.builtReportFields)
		})
	})

	t.Run("when previous report is present", func(t *testing.T) {
		*codec = testReportCodec{
			observationTimestamp: uint32(rand.Int31n(math.MaxInt16)),
			builtReport:          []byte{1, 2, 3, 4},
		}
		previousReport := ocrtypes.Report{}

		t.Run("succeeds and uses timestamp from previous report if valid", func(t *testing.T) {
			protos := newValidProtos()
			ts := codec.observationTimestamp + 1
			for i := range protos {
				protos[i].Timestamp = ts
			}
			aos := newValidAos(t, protos...)

			should, report, err := rp.Report(repts, previousReport, aos)
			require.NoError(t, err)

			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, ReportFields{
				ValidFromTimestamp: codec.observationTimestamp + 1, // previous observation timestamp +1 second
				Timestamp:          ts,
				NativeFee:          big.NewInt(2300000000000000000), // 2.3e18
				LinkFee:            big.NewInt(1300000000000000000), // 1.3e18
				ExpiresAt:          ts + 1,
				BenchmarkPrice:     big.NewInt(345),
				Bid:                big.NewInt(340),
				Ask:                big.NewInt(350),
			}, *codec.builtReportFields)
		})
		t.Run("errors if cannot extract timestamp from previous report", func(t *testing.T) {
			codec.err = errors.New("something exploded trying to extract timestamp")
			aos := newValidAos(t)

			should, _, err := rp.Report(ocrtypes.ReportTimestamp{}, previousReport, aos)
			assert.False(t, should)
			assert.EqualError(t, err, "something exploded trying to extract timestamp")
		})
		t.Run("does not report if observationTimestamp < validFromTimestamp", func(t *testing.T) {
			codec.observationTimestamp = 43
			codec.err = nil

			protos := newValidProtos()
			for i := range protos {
				protos[i].Timestamp = 42
			}
			aos := newValidAos(t, protos...)

			should, _, err := rp.Report(ocrtypes.ReportTimestamp{}, previousReport, aos)
			assert.False(t, should)
			assert.NoError(t, err)
		})
		t.Run("uses 0 values for link/native if they are invalid", func(t *testing.T) {
			codec.observationTimestamp = 42
			codec.err = nil

			protos := newValidProtos()
			for i := range protos {
				protos[i].LinkFeeValid = false
				protos[i].NativeFeeValid = false
			}
			aos := newValidAos(t, protos...)

			should, report, err := rp.Report(ocrtypes.ReportTimestamp{}, previousReport, aos)
			assert.True(t, should)
			assert.NoError(t, err)

			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, "0", codec.builtReportFields.LinkFee.String())
			assert.Equal(t, "0", codec.builtReportFields.NativeFee.String())
		})
	})

	t.Run("buildReport failures", func(t *testing.T) {
		t.Run("Report errors when the report is too large", func(t *testing.T) {
			aos := newValidAos(t)
			codec.builtReport = make([]byte, 1<<16)

			_, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)

			assert.EqualError(t, err, "report with len 65536 violates MaxReportLength limit set by ReportCodec (123)")
		})

		t.Run("Report errors when the report length is 0", func(t *testing.T) {
			aos := newValidAos(t)
			codec.builtReport = []byte{}
			_, _, err := rp.Report(ocrtypes.ReportTimestamp{}, nil, aos)

			assert.EqualError(t, err, "report may not have zero length (invariant violation)")
		})
	})
}

func Test_Plugin_validateReport(t *testing.T) {
	dataSource := &testDataSource{}
	codec := &testReportCodec{}
	rp := newTestReportPlugin(t, codec, dataSource)

	t.Run("valid reports", func(t *testing.T) {
		rf := ReportFields{
			ValidFromTimestamp: 42,
			Timestamp:          43,
			NativeFee:          big.NewInt(100),
			LinkFee:            big.NewInt(50),
			ExpiresAt:          44,
			BenchmarkPrice:     big.NewInt(150),
			Bid:                big.NewInt(140),
			Ask:                big.NewInt(160),
		}
		err := rp.validateReport(rf)
		require.NoError(t, err)

		rf = ReportFields{
			ValidFromTimestamp: 42,
			Timestamp:          42,
			NativeFee:          big.NewInt(0),
			LinkFee:            big.NewInt(0),
			ExpiresAt:          42,
			BenchmarkPrice:     big.NewInt(1),
			Bid:                big.NewInt(1),
			Ask:                big.NewInt(1),
		}
		err = rp.validateReport(rf)
		require.NoError(t, err)
	})
	t.Run("fails validation", func(t *testing.T) {
		rf := ReportFields{
			ValidFromTimestamp: 44, // later than timestamp not allowed
			Timestamp:          43,
			NativeFee:          big.NewInt(-1),     // negative value not allowed
			LinkFee:            big.NewInt(-1),     // negative value not allowed
			ExpiresAt:          42,                 // before timestamp
			BenchmarkPrice:     big.NewInt(150000), // exceeds max
			Bid:                big.NewInt(150000), // exceeds max
			Ask:                big.NewInt(150000), // exceeds max
		}
		err := rp.validateReport(rf)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "median benchmark price (Value: 150000) is outside of allowable range (Min: 1, Max: 1000)")
		assert.Contains(t, err.Error(), "median bid (Value: 150000) is outside of allowable range (Min: 1, Max: 1000)")
		assert.Contains(t, err.Error(), "median ask (Value: 150000) is outside of allowable range (Min: 1, Max: 1000)")
		assert.Contains(t, err.Error(), "median link fee (Value: -1) is outside of allowable range (Min: 0, Max: 3138550867693340381917894711603833208051177722232017256447)")
		assert.Contains(t, err.Error(), "median native fee (Value: -1) is outside of allowable range (Min: 0, Max: 3138550867693340381917894711603833208051177722232017256447)")
		assert.Contains(t, err.Error(), "observationTimestamp (Value: 43) must be >= validFromTimestamp (Value: 44)")
		assert.Contains(t, err.Error(), "expiresAt (Value: 42) must be ahead of observation timestamp (Value: 43)")
	})

	t.Run("zero values", func(t *testing.T) {
		rf := ReportFields{}
		err := rp.validateReport(rf)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "median benchmark price: got nil value")
		assert.Contains(t, err.Error(), "median bid: got nil value")
		assert.Contains(t, err.Error(), "median ask: got nil value")
		assert.Contains(t, err.Error(), "median native fee: got nil value")
		assert.Contains(t, err.Error(), "median link fee: got nil value")
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
			Bid:                        make([]byte, 24),
			Ask:                        make([]byte, 24),
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
		require.Equal(t, 11, numFields-3)

		b, err := proto.Marshal(&obs)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(b), maxObservationLength)
	})

	t.Run("all observations succeeded", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
			Bid: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
			Ask: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{
				Val: rand.Int63(),
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

		fee := mercury.CalculateFee(obs.LinkPrice.Val, decimal.NewFromInt32(1))
		assert.Equal(t, fee, mustDecodeBigInt(p.LinkFee))
		assert.True(t, p.LinkFeeValid)

		fee = mercury.CalculateFee(obs.NativePrice.Val, decimal.NewFromInt32(1))
		assert.Equal(t, fee, mustDecodeBigInt(p.NativeFee))
		assert.True(t, p.NativeFeeValid)
	})

	t.Run("negative link/native prices set fee to max int192", func(t *testing.T) {
		obs := Observation{
			LinkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(-1),
			},
			NativePrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(-1),
			},
		}
		dataSource.Obs = obs

		parsedObs, err := rp.Observation(context.Background(), ocrtypes.ReportTimestamp{}, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(parsedObs, &p))

		assert.Equal(t, mercury.MaxInt192, mustDecodeBigInt(p.LinkFee))
		assert.True(t, p.LinkFeeValid)
		assert.Equal(t, mercury.MaxInt192, mustDecodeBigInt(p.NativeFee))
		assert.True(t, p.NativeFeeValid)
	})

	t.Run("some observations failed", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: big.NewInt(rand.Int63()),
				Err: errors.New("bechmarkPrice error"),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{
				Val: rand.Int63(),
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

		fee := mercury.CalculateFee(obs.NativePrice.Val, decimal.NewFromInt32(1))
		assert.Equal(t, fee, mustDecodeBigInt(p.NativeFee))
		assert.True(t, p.NativeFeeValid)
	})

	t.Run("all observations failed", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Err: errors.New("bechmarkPrice error"),
			},
			Bid: mercury.ObsResult[*big.Int]{
				Err: errors.New("bid error"),
			},
			Ask: mercury.ObsResult[*big.Int]{
				Err: errors.New("ask error"),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{
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
		assert.Zero(t, p.Bid)
		assert.Zero(t, p.Ask)
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
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{
				Val: rand.Int63(),
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
	})

	t.Run("encoding fails on all observations", func(t *testing.T) {
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			Bid: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			Ask: mercury.ObsResult[*big.Int]{
				Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			},
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{
				Val: rand.Int63(),
			},
			// encoding never fails on calculated fees
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
		assert.Zero(t, p.Bid)
		assert.Zero(t, p.Ask)
		assert.False(t, p.PricesValid)
	})
}

func newUnparseableAttributedObservation() ocrtypes.AttributedObservation {
	return ocrtypes.AttributedObservation{
		Observation: []byte{1, 2},
		Observer:    commontypes.OracleID(42),
	}
}
