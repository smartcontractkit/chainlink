package mercury_v1

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury"
)

type testReportCodec struct {
	currentBlock          int64
	currentBlockErr       error
	builtReport           ocrtypes.Report
	buildReportShouldFail bool

	builtReportFields *ReportFields
}

func (trc *testReportCodec) reset() {
	trc.currentBlockErr = nil
	trc.buildReportShouldFail = false
	trc.builtReportFields = nil
}

func (trc *testReportCodec) BuildReport(rf ReportFields) (ocrtypes.Report, error) {
	if trc.buildReportShouldFail {
		return nil, errors.New("buildReportShouldFail=true")
	}
	trc.builtReportFields = &rf
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
		onchainConfig:   mercury.OnchainConfig{Min: big.NewInt(0), Max: big.NewInt(1000)},
		logger:          logger.Test(t),
		reportCodec:     codec,
		maxReportLength: maxReportLength,
	}
}

func newValidReportFields() ReportFields {
	return ReportFields{
		BenchmarkPrice:    big.NewInt(42),
		Bid:               big.NewInt(42),
		Ask:               big.NewInt(42),
		CurrentBlockNum:   42,
		ValidFromBlockNum: 42,
		CurrentBlockHash:  make([]byte, 32),
	}
}

func Test_ReportingPlugin_validateReport(t *testing.T) {
	rp := newReportingPlugin(t, &testReportCodec{})
	rf := newValidReportFields()

	t.Run("reports if currentBlockNum > validFromBlockNum", func(t *testing.T) {
		rf.CurrentBlockNum = 500
		rf.ValidFromBlockNum = 499
		err := rp.validateReport(rf)
		require.NoError(t, err)
	})
	t.Run("reports if currentBlockNum == validFromBlockNum", func(t *testing.T) {
		rf.CurrentBlockNum = 500
		rf.ValidFromBlockNum = 500
		err := rp.validateReport(rf)
		require.NoError(t, err)
	})
	t.Run("does not report if currentBlockNum < validFromBlockNum", func(t *testing.T) {
		rf.CurrentBlockNum = 499
		rf.ValidFromBlockNum = 500
		err := rp.validateReport(rf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validFromBlockNum (Value: 500) must be less than or equal to CurrentBlockNum (Value: 499)")
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
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func mustDecodeBigInt(b []byte) *big.Int {
	n, err := mercury.DecodeValueInt192(b)
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
				BenchmarkPrice: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Bid: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Ask: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				CurrentBlockNum: mercury.ObsResult[int64]{
					Val: rand.Int63(),
				},
				CurrentBlockHash: mercury.ObsResult[[]byte]{
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: mercury.ObsResult[uint64]{
					Val: rand.Uint64(),
				},
				LatestBlocks: []Block{
					Block{Num: rand.Int63(), Hash: string(randBytes(32)), Ts: rand.Uint64()},
					Block{Num: rand.Int63(), Hash: string(randBytes(32)), Ts: rand.Uint64()},
					Block{Num: rand.Int63(), Hash: string(randBytes(32)), Ts: rand.Uint64()},
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
			assert.Equal(t, len(obs.LatestBlocks), len(p.LatestBlocks))
			for i := range obs.LatestBlocks {
				assert.Equal(t, obs.LatestBlocks[i].Num, p.LatestBlocks[i].Num)
				assert.Equal(t, []byte(obs.LatestBlocks[i].Hash), p.LatestBlocks[i].Hash)
				assert.Equal(t, obs.LatestBlocks[i].Ts, p.LatestBlocks[i].Ts)
			}
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
				BenchmarkPrice: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
					Err: errors.New("benchmarkPrice exploded"),
				},
				Bid: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
					Err: errors.New("bid exploded"),
				},
				Ask: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
					Err: errors.New("ask exploded"),
				},
				CurrentBlockNum: mercury.ObsResult[int64]{
					Err: errors.New("currentBlockNum exploded"),
					Val: rand.Int63(),
				},
				CurrentBlockHash: mercury.ObsResult[[]byte]{
					Err: errors.New("currentBlockHash exploded"),
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: mercury.ObsResult[uint64]{
					Err: errors.New("currentBlockTimestamp exploded"),
					Val: rand.Uint64(),
				},
				LatestBlocks: ([]Block)(nil),
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
			assert.Len(t, p.LatestBlocks, 0)

			assert.False(t, p.PricesValid)
			assert.False(t, p.CurrentBlockValid)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.False(t, p.MaxFinalizedBlockNumberValid)
		})

		t.Run("when some observations have failed", func(t *testing.T) {
			obs := Observation{
				BenchmarkPrice: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Bid: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Ask: mercury.ObsResult[*big.Int]{
					Err: errors.New("ask exploded"),
				},
				CurrentBlockNum: mercury.ObsResult[int64]{
					Err: errors.New("currentBlockNum exploded"),
				},
				CurrentBlockHash: mercury.ObsResult[[]byte]{
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: mercury.ObsResult[uint64]{
					Val: rand.Uint64(),
				},
				LatestBlocks: []Block{
					Block{Num: rand.Int63(), Hash: string(randBytes(32)), Ts: rand.Uint64()},
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
			assert.Equal(t, len(obs.LatestBlocks), len(p.LatestBlocks))
			for i := range obs.LatestBlocks {
				assert.Equal(t, obs.LatestBlocks[i].Num, p.LatestBlocks[i].Num)
				assert.Equal(t, []byte(obs.LatestBlocks[i].Hash), p.LatestBlocks[i].Hash)
				assert.Equal(t, obs.LatestBlocks[i].Ts, p.LatestBlocks[i].Ts)
			}
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.Zero(t, p.MaxFinalizedBlockNumber)

			assert.False(t, p.PricesValid)
			assert.False(t, p.CurrentBlockValid)
			// since previousReport is not nil, maxFinalizedBlockNumber is skipped
			assert.False(t, p.MaxFinalizedBlockNumberValid)
		})

		t.Run("when encoding fails on some price observations", func(t *testing.T) {
			obs := Observation{
				BenchmarkPrice: mercury.ObsResult[*big.Int]{
					// too large to encode
					Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
				},
				Bid: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				Ask: mercury.ObsResult[*big.Int]{
					Val: randBigInt(),
				},
				CurrentBlockNum: mercury.ObsResult[int64]{
					Val: rand.Int63(),
				},
				CurrentBlockHash: mercury.ObsResult[[]byte]{
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: mercury.ObsResult[uint64]{
					Val: rand.Uint64(),
				},
			}
			rp.dataSource = mockDataSource{obs}

			pbObs, err := rp.Observation(ctx, repts, previousReport)
			require.NoError(t, err)

			var p MercuryObservationProto
			require.NoError(t, proto.Unmarshal(pbObs, &p))

			assert.False(t, p.PricesValid)
			assert.Zero(t, p.BenchmarkPrice)
			assert.NotZero(t, p.Bid)
			assert.NotZero(t, p.Ask)
		})
		t.Run("when encoding fails on all price observations", func(t *testing.T) {
			obs := Observation{
				BenchmarkPrice: mercury.ObsResult[*big.Int]{
					// too large to encode
					Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
				},
				Bid: mercury.ObsResult[*big.Int]{
					// too large to encode
					Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
				},
				Ask: mercury.ObsResult[*big.Int]{
					// too large to encode
					Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
				},
				CurrentBlockNum: mercury.ObsResult[int64]{
					Val: rand.Int63(),
				},
				CurrentBlockHash: mercury.ObsResult[[]byte]{
					Val: randBytes(32),
				},
				CurrentBlockTimestamp: mercury.ObsResult[uint64]{
					Val: rand.Uint64(),
				},
			}
			rp.dataSource = mockDataSource{obs}

			pbObs, err := rp.Observation(ctx, repts, previousReport)
			require.NoError(t, err)

			var p MercuryObservationProto
			require.NoError(t, proto.Unmarshal(pbObs, &p))

			assert.False(t, p.PricesValid)
			assert.Zero(t, p.BenchmarkPrice)
			assert.Zero(t, p.Bid)
			assert.Zero(t, p.Ask)
		})
	})

	t.Run("without previous report, includes maxFinalizedBlockNumber observation", func(t *testing.T) {
		currentBlockNum := int64(rand.Int31())
		obs := Observation{
			BenchmarkPrice: mercury.ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Bid: mercury.ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Ask: mercury.ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			CurrentBlockNum: mercury.ObsResult[int64]{
				Val: currentBlockNum,
			},
			CurrentBlockHash: mercury.ObsResult[[]byte]{
				Val: randBytes(32),
			},
			CurrentBlockTimestamp: mercury.ObsResult[uint64]{
				Val: rand.Uint64(),
			},
			MaxFinalizedBlockNumber: mercury.ObsResult[int64]{
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

func genRandHash(seed int64) []byte {
	r := rand.New(rand.NewSource(seed))

	b := make([]byte, 32)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func newValidMercuryObservationProto() *MercuryObservationProto {
	var latestBlocks []*BlockProto
	for i := 0; i < MaxAllowedBlocks; i++ {
		latestBlocks = append(latestBlocks, &BlockProto{Num: int64(49 - i), Hash: genRandHash(int64(i)), Ts: uint64(46 - i)})
	}

	return &MercuryObservationProto{
		Timestamp:                    42,
		BenchmarkPrice:               mercury.MustEncodeValueInt192(big.NewInt(43)),
		Bid:                          mercury.MustEncodeValueInt192(big.NewInt(44)),
		Ask:                          mercury.MustEncodeValueInt192(big.NewInt(45)),
		PricesValid:                  true,
		CurrentBlockNum:              latestBlocks[0].Num,
		CurrentBlockHash:             latestBlocks[0].Hash,
		CurrentBlockTimestamp:        latestBlocks[0].Ts,
		CurrentBlockValid:            true,
		LatestBlocks:                 latestBlocks,
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
	t.Run("with all valid values, and > 0 LatestBlocks", func(t *testing.T) {
		obs := newValidMercuryObservationProto()
		ao := newAttributedObservation(t, obs)
		expectedLatestBlocks := make([]Block, len(obs.LatestBlocks))
		for i, b := range obs.LatestBlocks {
			expectedLatestBlocks[i] = NewBlock(b.Num, b.Hash, b.Ts)
		}

		pao, err := parseAttributedObservation(ao)
		require.NoError(t, err)

		assert.Equal(t,
			parsedAttributedObservation{
				Timestamp:                    0x2a,
				Observer:                     0x2a,
				BenchmarkPrice:               big.NewInt(43),
				Bid:                          big.NewInt(44),
				Ask:                          big.NewInt(45),
				PricesValid:                  true,
				MaxFinalizedBlockNumber:      47,
				MaxFinalizedBlockNumberValid: true,
				LatestBlocks:                 expectedLatestBlocks,
			},
			pao,
		)
		t.Run("with 0 LatestBlocks", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.LatestBlocks = nil
			ao := newAttributedObservation(t, obs)

			pao, err := parseAttributedObservation(ao)
			require.NoError(t, err)

			assert.Equal(t,
				parsedAttributedObservation{
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
	})

	t.Run("with all invalid values", func(t *testing.T) {
		obs := newInvalidMercuryObservationProto()
		ao := newAttributedObservation(t, obs)

		pao, err := parseAttributedObservation(ao)
		assert.NoError(t, err)

		assert.Equal(t,
			parsedAttributedObservation{
				Observer:                     0x2a,
				PricesValid:                  false,
				CurrentBlockValid:            false,
				MaxFinalizedBlockNumberValid: false,
				LatestBlocks:                 ([]Block)(nil),
			},
			pao,
		)
	})

	t.Run("when LatestBlocks is valid", func(t *testing.T) {
		t.Run("sorts blocks if they are out of order", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			slices.Reverse(obs.LatestBlocks)

			ao := newAttributedObservation(t, obs)

			pao, err := parseAttributedObservation(ao)
			assert.NoError(t, err)

			assert.Len(t, pao.GetLatestBlocks(), MaxAllowedBlocks)
			assert.Equal(t, 49, int(pao.GetLatestBlocks()[0].Num))
			assert.Equal(t, 48, int(pao.GetLatestBlocks()[1].Num))
			assert.Equal(t, 47, int(pao.GetLatestBlocks()[2].Num))
			assert.Equal(t, 46, int(pao.GetLatestBlocks()[3].Num))
			assert.Equal(t, 45, int(pao.GetLatestBlocks()[4].Num))
		})
	})

	t.Run("when LatestBlocks is invalid", func(t *testing.T) {
		t.Run("contains duplicate block numbers", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.LatestBlocks = []*BlockProto{&BlockProto{Num: 32, Hash: randBytes(32)}, &BlockProto{Num: 32, Hash: randBytes(32)}}
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, "observation invalid for observer 42; got duplicate block number: 32")
		})
		t.Run("contains duplicate block hashes", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			h := randBytes(32)
			obs.LatestBlocks = []*BlockProto{&BlockProto{Num: 1, Hash: h}, &BlockProto{Num: 2, Hash: h}}
			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, fmt.Sprintf("observation invalid for observer 42; got duplicate block hash: 0x%x", h))
		})
		t.Run("contains too many blocks", func(t *testing.T) {
			obs := newValidMercuryObservationProto()
			obs.LatestBlocks = nil
			for i := 0; i < MaxAllowedBlocks+1; i++ {
				obs.LatestBlocks = append(obs.LatestBlocks, &BlockProto{Num: int64(i)})
			}

			ao := newAttributedObservation(t, obs)

			_, err := parseAttributedObservation(ao)
			assert.EqualError(t, err, fmt.Sprintf("LatestBlocks too large; got: %d, max: %d", MaxAllowedBlocks+1, MaxAllowedBlocks))
		})
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
			t.Run("CurrentBlockHash", func(t *testing.T) {
				obs := newValidMercuryObservationProto()
				obs.LatestBlocks = nil
				obs.CurrentBlockHash = []byte{1}
				ao := newAttributedObservation(t, obs)

				_, err := parseAttributedObservation(ao)
				assert.EqualError(t, err, "wrong len for hash: 1 (expected: 32)")
			})

			t.Run("LatestBlocks", func(t *testing.T) {
				obs := newValidMercuryObservationProto()
				obs.LatestBlocks[0].Hash = []byte{1}
				ao := newAttributedObservation(t, obs)

				_, err := parseAttributedObservation(ao)
				assert.EqualError(t, err, "wrong len for hash: 1 (expected: 32)")
			})
		})
		t.Run("negative block number", func(t *testing.T) {
			t.Run("CurrentBlockNum", func(t *testing.T) {
				obs := newValidMercuryObservationProto()
				obs.LatestBlocks = nil
				obs.CurrentBlockNum = -1
				ao := newAttributedObservation(t, obs)

				_, err := parseAttributedObservation(ao)
				assert.EqualError(t, err, "negative block number: -1")
			})
			t.Run("LatestBlocks", func(t *testing.T) {
				obs := newValidMercuryObservationProto()
				obs.LatestBlocks[0].Num = -1
				ao := newAttributedObservation(t, obs)

				_, err := parseAttributedObservation(ao)
				assert.EqualError(t, err, "negative block number: -1")
			})
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
			assert.EqualError(t, err, "got zero valid attributed observations")
		})
		t.Run("succeeds, ignoring unparseable attributed observations", func(t *testing.T) {
			aos := []types.AttributedObservation{
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newUnparseableAttributedObservation(),
			}
			should, report, err := rp.Report(repts, nil, aos)

			assert.NoError(t, err)
			assert.True(t, should)
			assert.Equal(t, codec.builtReport, report)
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

			require.NotNil(t, codec.builtReportFields)
			assert.Equal(t, 48, int(codec.builtReportFields.ValidFromBlockNum))
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

			assert.EqualError(t, err, "no valid maxFinalizedBlockNumber with at least f+1 votes (got counts: map[0:1 1:1 2:1 3:1], f=1)")
		})
		t.Run("errors if it cannot come to consensus about currentBlockNum", func(t *testing.T) {
			obs := []*MercuryObservationProto{
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
			}
			for i := range obs {
				obs[i].LatestBlocks = nil
				obs[i].CurrentBlockNum = int64(i)
			}
			aos := []types.AttributedObservation{
				newAttributedObservation(t, obs[0]),
				newAttributedObservation(t, obs[1]),
				newAttributedObservation(t, obs[2]),
				newAttributedObservation(t, obs[3]),
			}
			_, _, err := rp.Report(repts, nil, aos)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "GetConsensusCurrentBlock failed: cannot come to consensus on latest block number")
		})
		t.Run("errors if it cannot come to consensus on LatestBlocks", func(t *testing.T) {
			obs := []*MercuryObservationProto{
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
			}
			for i := range obs {
				for j := range obs[i].LatestBlocks {
					obs[i].LatestBlocks[j].Hash = randBytes(32)
				}
			}
			aos := []types.AttributedObservation{
				newAttributedObservation(t, obs[0]),
				newAttributedObservation(t, obs[1]),
				newAttributedObservation(t, obs[2]),
				newAttributedObservation(t, obs[3]),
			}
			_, _, err := rp.Report(repts, nil, aos)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "GetConsensusCurrentBlock failed: cannot come to consensus on latest block number")
		})
		t.Run("errors if price is invalid", func(t *testing.T) {
			obs := []*MercuryObservationProto{
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
				newValidMercuryObservationProto(),
			}
			for i := range obs {
				obs[i].BenchmarkPrice = mercury.MustEncodeValueInt192(big.NewInt(-1)) // benchmark price below min of 0, cannot report
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
			assert.EqualError(t, err, "median benchmark price (Value: -1) is outside of allowable range (Min: 0, Max: 1000)")
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

			require.NotNil(t, codec.builtReportFields)
			// current block of previous report + 1 is the validFromBlockNum of current report
			assert.Equal(t, 33, int(codec.builtReportFields.ValidFromBlockNum))
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
		t.Run("does not report if currentBlockNum < validFromBlockNum", func(t *testing.T) {
			codec.currentBlock = 49 // means that validFromBlockNum=50 which is > currentBlockNum of 49
			codec.currentBlockErr = nil

			aos := []types.AttributedObservation{
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
				newAttributedObservation(t, newValidMercuryObservationProto()),
			}
			should, _, err := rp.Report(repts, previousReport, aos)

			assert.False(t, should)
			assert.NoError(t, err)
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
		maxLatestBlocks := []*BlockProto{}
		for i := 0; i < MaxAllowedBlocks; i++ {
			maxLatestBlocks = append(maxLatestBlocks, &BlockProto{Num: math.MaxInt64, Hash: maxHash, Ts: math.MaxUint64})
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
			LatestBlocks:                 maxLatestBlocks,
		}
		// This assertion is here to force this test to fail if a new field is
		// added to the protobuf. In this case, you must add the max value of
		// the field to the MercuryObservationProto in the test and only after
		// that increment the count below
		numFields := reflect.TypeOf(obs).NumField() //nolint:all
		// 3 fields internal to pbuf struct
		require.Equal(t, 12, numFields-3)

		// the actual test
		b, err := proto.Marshal(&obs)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(b), maxObservationLength)
	})
}
