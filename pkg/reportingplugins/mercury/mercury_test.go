package mercury

import (
	"context"
	"math"
	"math/big"
	"math/rand"
	reflect "reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

type testReportCodec struct {
	currentBlock          int64
	currentBlockErr       error
	builtReport           ocrtypes.Report
	buildReportShouldFail bool

	lastBuildReportPaos              []ParsedAttributedObservation
	lastBuildReportF                 int
	lastBuildReportValidFromBlockNum int64
}

func (trc *testReportCodec) reset() {
	trc.currentBlockErr = nil
	trc.buildReportShouldFail = false
	trc.lastBuildReportPaos = nil
	trc.lastBuildReportF = 0
	trc.lastBuildReportValidFromBlockNum = 0
}

func (trc *testReportCodec) BuildReport(paos []ParsedAttributedObservation, f int, validFromBlockNum int64) (ocrtypes.Report, error) {
	if trc.buildReportShouldFail {
		return nil, errors.New("buildReportShouldFail=true")
	}
	trc.lastBuildReportPaos = paos
	trc.lastBuildReportF = f
	trc.lastBuildReportValidFromBlockNum = validFromBlockNum
	return trc.builtReport, nil
}

func (trc *testReportCodec) MaxReportLength(n int) (int, error) {
	return 8*32 + // feed ID
			32 + // timestamp
			192 + // benchmarkPrice
			192 + // bid
			192 + // ask
			64 + //currentBlockNum
			8*32 + // currentBlockHash
			64, // validFromBlockNum
		nil
}

func (trc *testReportCodec) CurrentBlockNumFromReport(types.Report) (int64, error) {
	return trc.currentBlock, trc.currentBlockErr
}

func newReportingPlugin(t *testing.T, codec *testReportCodec) *reportingPlugin {
	maxReportLength, err := codec.MaxReportLength(4)
	require.NoError(t, err)
	return &reportingPlugin{
		f:               1,
		onchainConfig:   OnchainConfig{Min: big.NewInt(0), Max: big.NewInt(1000)},
		logger:          logger.Test(t),
		reportCodec:     codec,
		maxReportLength: maxReportLength,
	}
}

func Test_ReportingPlugin_shouldReport(t *testing.T) {
	rp := newReportingPlugin(t, &testReportCodec{})
	repts := types.ReportTimestamp{}
	paos := NewValidParsedAttributedObservations()

	t.Run("reports if all reports have currentBlockNum > validFromBlockNum", func(t *testing.T) {
		for i := range paos {
			paos[i].CurrentBlockNum = 500
		}
		shouldReport, err := rp.shouldReport(499, repts, paos)
		require.NoError(t, err)

		assert.True(t, shouldReport)
	})
	t.Run("reporta if all reports have currentBlockNum == validFromBlockNum", func(t *testing.T) {
		for i := range paos {
			paos[i].CurrentBlockNum = 500
		}
		shouldReport, err := rp.shouldReport(500, repts, paos)
		require.NoError(t, err)

		assert.True(t, shouldReport)
	})
	t.Run("does not report if all reports have currentBlockNum < validFromBlockNum", func(t *testing.T) {
		paos := NewValidParsedAttributedObservations()
		for i := range paos {
			paos[i].CurrentBlockNum = 499
		}
		shouldReport, err := rp.shouldReport(500, repts, paos)
		require.NoError(t, err)

		assert.False(t, shouldReport)
	})
	t.Run("returns error if it cannot come to consensus about currentBlockNum", func(t *testing.T) {
		paos := NewValidParsedAttributedObservations()
		for i := range paos {
			paos[i].CurrentBlockNum = 500 + int64(i)
		}
		shouldReport, err := rp.shouldReport(499, repts, paos)
		require.NoError(t, err)

		assert.False(t, shouldReport)
	})
}

var _ DataSource = &mockDataSource{}

type mockDataSource struct{ obs Observation }

func (m mockDataSource) Observe(context.Context, ocrtypes.ReportTimestamp, bool) (Observation, error) {
	return m.obs, nil
}

func randBigInt() *big.Int {
	return big.NewInt(rand.Int63())
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func mustDecodeBigInt(b []byte) *big.Int {
	n, err := DecodeValueInt192(b)
	if err != nil {
		panic(err)
	}
	return n
}

func Test_Plugin_Observation(t *testing.T) {
	ctx := context.Background()
	repts := ocrtypes.ReportTimestamp{}
	codec := &testReportCodec{
		currentBlock: int64(rand.Int31()),
		builtReport:  []byte{},
	}

	rp := newReportingPlugin(t, codec)

	t.Run("with previous report", func(t *testing.T) {
		// content of previousReport is irrelevant, the only thing that matters
		// for this test is that it's not nil
		previousReport := ocrtypes.Report{}

		t.Run("when all observations are successful", func(t *testing.T) {
			obs := Observation{
				BenchmarkPrice: ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Bid: ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Ask: ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				CurrentBlockNum: ObsResult[int64]{
					Val: rand.Int63(),
				},
				CurrentBlockHash: ObsResult[[]byte]{
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: ObsResult[uint64]{
					Val: rand.Uint64(),
				},
			}

			rp.dataSource = mockDataSource{obs}

			pbObs, err := rp.Observation(ctx, repts, previousReport)
			require.NoError(t, err)

			var p MercuryObservationProto
			require.NoError(t, proto.Unmarshal(pbObs, &p))

			assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
			assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
			assert.Equal(t, obs.Bid.Val, mustDecodeBigInt(p.Bid))
			assert.Equal(t, obs.Ask.Val, mustDecodeBigInt(p.Ask))
			assert.Equal(t, obs.CurrentBlockNum.Val, p.CurrentBlockNum)
			assert.Equal(t, obs.CurrentBlockHash.Val, p.CurrentBlockHash)
			assert.Equal(t, obs.CurrentBlockTimestamp.Val, p.CurrentBlockTimestamp)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.Zero(t, p.MaxFinalizedBlockNumber)

			assert.True(t, p.PricesValid)
			assert.True(t, p.CurrentBlockValid)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.False(t, p.MaxFinalizedBlockNumberValid)
		})

		t.Run("when all observations have failed", func(t *testing.T) {
			obs := Observation{
				// Vals should be ignored, this is asserted with .Zero below
				BenchmarkPrice: ObsResult[*big.Int]{
					Val: randBigInt(),
					Err: errors.New("benchmarkPrice exploded"),
				},
				Bid: ObsResult[*big.Int]{
					Val: randBigInt(),
					Err: errors.New("bid exploded"),
				},
				Ask: ObsResult[*big.Int]{
					Val: randBigInt(),
					Err: errors.New("ask exploded"),
				},
				CurrentBlockNum: ObsResult[int64]{
					Err: errors.New("currentBlockNum exploded"),
					Val: rand.Int63(),
				},
				CurrentBlockHash: ObsResult[[]byte]{
					Err: errors.New("currentBlockHash exploded"),
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: ObsResult[uint64]{
					Err: errors.New("currentBlockTimestamp exploded"),
					Val: rand.Uint64(),
				},
			}
			rp.dataSource = mockDataSource{obs}

			pbObs, err := rp.Observation(ctx, repts, previousReport)
			require.NoError(t, err)

			var p MercuryObservationProto
			require.NoError(t, proto.Unmarshal(pbObs, &p))

			assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
			assert.Zero(t, p.BenchmarkPrice)
			assert.Zero(t, p.Bid)
			assert.Zero(t, p.Ask)
			assert.Zero(t, p.CurrentBlockNum)
			assert.Zero(t, p.CurrentBlockHash)
			assert.Zero(t, p.CurrentBlockTimestamp)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.Zero(t, p.MaxFinalizedBlockNumber)

			assert.False(t, p.PricesValid)
			assert.False(t, p.CurrentBlockValid)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.False(t, p.MaxFinalizedBlockNumberValid)
		})

		t.Run("when some observations have failed", func(t *testing.T) {
			obs := Observation{
				BenchmarkPrice: ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Bid: ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Ask: ObsResult[*big.Int]{
					Err: errors.New("ask exploded"),
				},
				CurrentBlockNum: ObsResult[int64]{
					Err: errors.New("currentBlockNum exploded"),
				},
				CurrentBlockHash: ObsResult[[]byte]{
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: ObsResult[uint64]{
					Val: rand.Uint64(),
				},
			}
			rp.dataSource = mockDataSource{obs}

			pbObs, err := rp.Observation(ctx, repts, previousReport)
			require.NoError(t, err)

			var p MercuryObservationProto
			require.NoError(t, proto.Unmarshal(pbObs, &p))

			assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
			assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
			assert.Equal(t, obs.Bid.Val, mustDecodeBigInt(p.Bid))
			assert.Zero(t, p.Ask)
			assert.Zero(t, p.CurrentBlockNum)
			assert.Equal(t, obs.CurrentBlockHash.Val, p.CurrentBlockHash)
			assert.Equal(t, obs.CurrentBlockTimestamp.Val, p.CurrentBlockTimestamp)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.Zero(t, p.MaxFinalizedBlockNumber)

			assert.False(t, p.PricesValid)
			assert.False(t, p.CurrentBlockValid)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.False(t, p.MaxFinalizedBlockNumberValid)
		})
	})

	t.Run("without previous report, includes maxFinalizedBlockNumber observation", func(t *testing.T) {
		currentBlockNum := int64(rand.Int31())
		obs := Observation{
			BenchmarkPrice: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Bid: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Ask: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			CurrentBlockNum: ObsResult[int64]{
				Val: currentBlockNum,
			},
			CurrentBlockHash: ObsResult[[]byte]{
				Val: randBytes(32),
			},
			CurrentBlockTimestamp: ObsResult[uint64]{
				Val: rand.Uint64(),
			},
			MaxFinalizedBlockNumber: ObsResult[int64]{
				Val: currentBlockNum - 42,
			},
		}
		rp.dataSource = mockDataSource{obs}

		pbObs, err := rp.Observation(ctx, repts, nil)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(pbObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
		assert.Equal(t, obs.Bid.Val, mustDecodeBigInt(p.Bid))
		assert.Equal(t, obs.Ask.Val, mustDecodeBigInt(p.Ask))
		assert.Equal(t, obs.CurrentBlockNum.Val, p.CurrentBlockNum)
		assert.Equal(t, obs.CurrentBlockHash.Val, p.CurrentBlockHash)
		assert.Equal(t, obs.CurrentBlockTimestamp.Val, p.CurrentBlockTimestamp)
		assert.Equal(t, obs.MaxFinalizedBlockNumber.Val, p.MaxFinalizedBlockNumber)

		assert.True(t, p.PricesValid)
		assert.True(t, p.CurrentBlockValid)
		assert.True(t, p.MaxFinalizedBlockNumberValid)
	})

}

func newAttributedObservation(t *testing.T, p *MercuryObservationProto) ocrtypes.AttributedObservation {
	marshalledObs, err := proto.Marshal(p)
	require.NoError(t, err)
	return ocrtypes.AttributedObservation{
		Observation: ocrtypes.Observation(marshalledObs),
		Observer:    commontypes.OracleID(42),
	}
}

func newUnparseableAttributedObservation() ocrtypes.AttributedObservation {
	return ocrtypes.AttributedObservation{
		Observation: []byte{1, 2},
		Observer:    commontypes.OracleID(42),
	}
}

var blockHash = randBytes(32)

func newValidMercuryObservationProto() *MercuryObservationProto {
	return &MercuryObservationProto{
		Timestamp:                    42,
		BenchmarkPrice:               MustEncodeValueInt192(big.NewInt(43)),
		Bid:                          MustEncodeValueInt192(big.NewInt(44)),
		Ask:                          MustEncodeValueInt192(big.NewInt(45)),
		PricesValid:                  true,
		CurrentBlockNum:              49,
		CurrentBlockHash:             blockHash,
		CurrentBlockTimestamp:        46,
		CurrentBlockValid:            true,
		MaxFinalizedBlockNumber:      47,
		MaxFinalizedBlockNumberValid: true,
	}
}

func newInvalidMercuryObservationProto() *MercuryObservationProto {
	return &MercuryObservationProto{
		PricesValid:                  false,
		CurrentBlockValid:            false,
		MaxFinalizedBlockNumberValid: false,
	}
}

func Test_Plugin_parseAttributedObservation(t *testing.T) {
	t.Run("with all valid values", func(t *testing.T) {
		obs := newValidMercuryObservationProto()
		ao := newAttributedObservation(t, obs)

		pao, err := parseAttributedObservation(ao)
		assert.NoError(t, err)

		assert.Equal(t,
			ParsedAttributedObservation{
				Timestamp:                    0x2a,
				Observer:                     0x2a,
				BenchmarkPrice:               big.NewInt(43),
				Bid:                          big.NewInt(44),
				Ask:                          big.NewInt(45),
				PricesValid:                  true,
				CurrentBlockNum:              49,
				CurrentBlockHash:             obs.CurrentBlockHash,
				CurrentBlockTimestamp:        46,
				CurrentBlockValid:            true,
				MaxFinalizedBlockNumber:      47,
				MaxFinalizedBlockNumberValid: true,
			},
			pao,
		)
	})

	t.Run("with all invalid values", func(t *testing.T) {
		obs := newInvalidMercuryObservationProto()
		ao := newAttributedObservation(t, obs)

		pao, err := parseAttributedObservation(ao)
		assert.NoError(t, err)

		assert.Equal(t,
			ParsedAttributedObservation{
				Observer:                     0x2a,
				PricesValid:                  false,
				CurrentBlockValid:            false,
				MaxFinalizedBlockNumberValid: false,
			},
			pao,
		)
	})

	t.Run("with unparseable values", func(t *testing.T) {
		t.Run("ao cannot be unmarshalled", func(t *testing.T) {
			ao := newUnparseableAttributedObservation()

			_, err := parseAttributedObservation(ao)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "attributed observation cannot be unmarshaled")
		})
		t.Run("bad benchmark price", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.BenchmarkPrice = randBytes(16)
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, "benchmarkPrice cannot be converted to big.Int: expected b to have length 24, but got length 16")
		})
		t.Run("bad bid", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.Bid = []byte{1}
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, "bid cannot be converted to big.Int: expected b to have length 24, but got length 1")
		})
		t.Run("bad ask", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.Ask = []byte{1}
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, "ask cannot be converted to big.Int: expected b to have length 24, but got length 1")
		})
		t.Run("bad block hash", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.CurrentBlockHash = []byte{1}
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, "wrong len for hash: 1 (expected: 32)")
		})
		t.Run("negative block number", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.CurrentBlockNum = -1
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, "negative block number: -1")
		})
	})
}

func Test_Plugin_Report(t *testing.T) {
	repts := types.ReportTimestamp{}

	t.Run("when previous report is nil", func(t *testing.T) {
		codec := &testReportCodec{
			currentBlock: int64(rand.Int31()),
			builtReport:  []byte{1, 2, 3, 4},
		}
		rp := newReportingPlugin(t, codec)

		t.Run("errors if not enough attributed observations", func(t *testing.T) {
			_, _, err := rp.Report(repts, nil, []types.AttributedObservation{})
			assert.EqualError(t, err, "only received 0 valid attributed observations, but need at least f+1 (2)")
		})
		t.Run("succeeds, ignoring unparseable attributed observations", func(t *testing.T) {
			aos := []types.AttributedObservation{
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newUnparseableAttributedObservation(),
			}
			should, report, err := rp.Report(repts, nil, aos)

			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
			assert.NoError(t, err)
		})
		t.Run("succeeds and generates validFromBlockNum from maxFinalizedBlockNumber", func(t *testing.T) {
			codec.reset()

			aos := []types.AttributedObservation{
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
			}
			should, report, err := rp.Report(repts, nil, aos)

			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
			assert.NoError(t, err)

			assert.Equal(t, 48, int(codec.lastBuildReportValidFromBlockNum))
			assert.Len(t, codec.lastBuildReportPaos, 4)
			assert.Equal(t, 1, codec.lastBuildReportF)
		})
		t.Run("errors if cannot get consensus maxFinalizedBlockNumber", func(t *testing.T) {
			obs := []*MercuryObservationProto{
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
			}
			for i := range obs {
				obs[i].MaxFinalizedBlockNumber = int64(i)
			}
			aos := []types.AttributedObservation{
				newAttributedObservation(t, obs[0]),
				newAttributedObservation(t, obs[1]),
				newAttributedObservation(t, obs[2]),
				newAttributedObservation(t, obs[3]),
			}
			_, _, err := rp.Report(repts, nil, aos)

			assert.EqualError(t, err, "no valid maxFinalizedBlockNumber with at least f+1 votes (got counts: map[0:1 1:1 2:1 3:1])")
		})
		t.Run("returns false if shouldReport returns false", func(t *testing.T) {
			obs := []*MercuryObservationProto{
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
			}
			for i := range obs {
				obs[i].BenchmarkPrice = MustEncodeValueInt192(big.NewInt(-1)) // benchmark price below min of 0, cannot report
			}
			aos := []types.AttributedObservation{
				newAttributedObservation(t, obs[0]),
				newAttributedObservation(t, obs[1]),
				newAttributedObservation(t, obs[2]),
				newAttributedObservation(t, obs[3]),
			}
			should, report, err := rp.Report(repts, nil, aos)

			assert.False(t, should)
			assert.Nil(t, report)
			assert.NoError(t, err)
		})
		t.Run("BuildReport failures", func(t *testing.T) {
			t.Run("errors if BuildReport returns error", func(t *testing.T) {
				codec.buildReportShouldFail = true
				defer codec.reset()

				aos := []types.AttributedObservation{
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newUnparseableAttributedObservation(),
				}
				_, _, err := rp.Report(repts, nil, aos)

				assert.EqualError(t, err, "buildReportShouldFail=true")
			})
			t.Run("errors if BuildReport returns a report that is too long", func(t *testing.T) {
				codec.builtReport = randBytes(9999)
				aos := []types.AttributedObservation{
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newUnparseableAttributedObservation(),
				}
				_, _, err := rp.Report(repts, nil, aos)

				assert.EqualError(t, err, "report with len 9999 violates MaxReportLength limit set by ReportCodec (1248)")
			})
			t.Run("errors if BuildReport returns a report that is too short", func(t *testing.T) {
				codec.builtReport = []byte{}
				aos := []types.AttributedObservation{
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newAttributedObservation(t, newValidMercuryObservationProto()),
					newUnparseableAttributedObservation(),
				}
				_, _, err := rp.Report(repts, nil, aos)

				assert.EqualError(t, err, "report may not have zero length (invariant violation)")
			})
		})
	})
	t.Run("when previous report is present", func(t *testing.T) {
		codec := &testReportCodec{
			currentBlock: int64(rand.Int31()),
			builtReport:  []byte{1, 2, 3, 4},
		}
		rp := newReportingPlugin(t, codec)
		previousReport := types.Report{}

		t.Run("succeeds and uses block number in previous report if valid", func(t *testing.T) {
			currentBlock := int64(32)
			codec.currentBlock = currentBlock

			aos := []types.AttributedObservation{
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
			}
			should, report, err := rp.Report(repts, previousReport, aos)

			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
			assert.NoError(t, err)

			// current block of previous report + 1 is the validFromBlockNum of current report
			assert.Equal(t, 33, int(codec.lastBuildReportValidFromBlockNum))
			assert.Len(t, codec.lastBuildReportPaos, 4)
			assert.Equal(t, 1, codec.lastBuildReportF)
		})
		t.Run("errors if cannot extract block number from previous report", func(t *testing.T) {
			codec.currentBlockErr = errors.New("test error current block fail")

			aos := []types.AttributedObservation{
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
			}
			should, _, err := rp.Report(repts, previousReport, aos)

			assert.False(t, should)
			assert.EqualError(t, err, "test error current block fail")
		})
	})
}

func Test_MaxObservationLength(t *testing.T) {
	t.Run("maximally sized pbuf does not exceed maxObservationLength", func(t *testing.T) {
		maxInt192Bytes := make([]byte, 24)
		for i := 0; i < 24; i++ {
			maxInt192Bytes[i] = 255
		}
		maxHash := make([]byte, 32)
		for i := 0; i < 32; i++ {
			maxHash[i] = 255
		}
		obs := MercuryObservationProto{
			Timestamp:                    math.MaxUint32,
			BenchmarkPrice:               maxInt192Bytes,
			Bid:                          maxInt192Bytes,
			Ask:                          maxInt192Bytes,
			PricesValid:                  true,
			CurrentBlockNum:              math.MaxInt64,
			CurrentBlockHash:             maxHash,
			CurrentBlockTimestamp:        math.MaxUint64,
			CurrentBlockValid:            true,
			MaxFinalizedBlockNumber:      math.MaxInt64,
			MaxFinalizedBlockNumberValid: true,
		}
		// This assertion is here to force this test to fail if a new field is
		// added to the protobuf. In this case, you must add the max value of
		// the field to the MercuryObservationProto in the test and only after
		// that increment the count below
		numFields := reflect.TypeOf(obs).NumField() //nolint:all
		// 3 fields internal to pbuf struct
		require.Equal(t, 11, numFields-3)

		// the actual test
		b, err := proto.Marshal(&obs)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(b), maxObservationLength)
	})
}
