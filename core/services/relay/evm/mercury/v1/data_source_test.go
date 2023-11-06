package mercury_v1

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaymercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"
	commonmocks "github.com/smartcontractkit/chainlink/v2/common/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	mercurymocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reportcodecv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ relaymercury.MercuryServerFetcher = &mockFetcher{}

type mockFetcher struct {
	num *int64
	err error
}

func (m *mockFetcher) FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error) {
	return m.num, m.err
}

func (m *mockFetcher) LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error) {
	return nil, nil
}

func (m *mockFetcher) LatestTimestamp(context.Context) (int64, error) {
	return 0, nil
}

var _ types.ChainHeadTracker = &mockHeadTracker{}

type mockHeadTracker struct {
	h httypes.HeadTracker
}

func (m *mockHeadTracker) HeadTracker() httypes.HeadTracker { return m.h }

type mockORM struct {
	report []byte
	err    error
}

func (m *mockORM) LatestReport(ctx context.Context, feedID [32]byte, qopts ...pg.QOpt) (report []byte, err error) {
	return m.report, m.err
}

func TestMercury_Observe(t *testing.T) {
	orm := &mockORM{}
	lggr := logger.TestLogger(t)
	ds := NewDataSource(orm, nil, job.Job{}, pipeline.Spec{}, lggr, nil, nil, nil, nil, nil, mercuryutils.FeedID{})
	ctx := testutils.Context(t)
	repts := ocrtypes.ReportTimestamp{}

	fetcher := &mockFetcher{}
	ds.fetcher = fetcher

	trrs := []pipeline.TaskRunResult{
		{
			// benchmark price
			Result: pipeline.Result{Value: "122.345"},
			Task:   &mercurymocks.MockTask{},
		},
		{
			// bid
			Result: pipeline.Result{Value: "121.993"},
			Task:   &mercurymocks.MockTask{},
		},
		{
			// ask
			Result: pipeline.Result{Value: "123.111"},
			Task:   &mercurymocks.MockTask{},
		},
	}

	runner := &mercurymocks.MockRunner{
		Trrs: trrs,
	}
	ds.pipelineRunner = runner

	spec := pipeline.Spec{}
	ds.spec = spec

	h := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
	ht := &mockHeadTracker{
		h: h,
	}
	ds.chainHeadTracker = ht

	head := &evmtypes.Head{
		Number:    int64(rand.Int31()),
		Hash:      utils.NewHash(),
		Timestamp: time.Now(),
	}
	h.On("LatestChain").Return(head)

	t.Run("when fetchMaxFinalizedBlockNum=true", func(t *testing.T) {
		t.Run("with latest report in database", func(t *testing.T) {
			orm.report = buildSampleV1Report()
			orm.err = nil

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.NoError(t, obs.MaxFinalizedBlockNumber.Err)
			assert.Equal(t, int64(143), obs.MaxFinalizedBlockNumber.Val)
		})
		t.Run("if querying latest report fails", func(t *testing.T) {
			orm.report = nil
			orm.err = errors.New("something exploded")

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "something exploded")
			assert.Zero(t, obs.MaxFinalizedBlockNumber.Val)
		})
		t.Run("if decoding latest report fails", func(t *testing.T) {
			orm.report = []byte{1, 2, 3}
			orm.err = nil

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")
			assert.Zero(t, obs.MaxFinalizedBlockNumber.Val)
		})

		orm.report = nil
		orm.err = nil

		t.Run("without latest report in database", func(t *testing.T) {
			t.Run("if FetchInitialMaxFinalizedBlockNumber returns error", func(t *testing.T) {
				fetcher.err = errors.New("mock fetcher error")

				obs, err := ds.Observe(ctx, repts, true)
				assert.NoError(t, err)

				assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "mock fetcher error")
				assert.Zero(t, obs.MaxFinalizedBlockNumber.Val)
			})
			t.Run("if FetchInitialMaxFinalizedBlockNumber succeeds", func(t *testing.T) {
				fetcher.err = nil
				var num int64 = 32
				fetcher.num = &num

				obs, err := ds.Observe(ctx, repts, true)
				assert.NoError(t, err)

				assert.NoError(t, obs.MaxFinalizedBlockNumber.Err)
				assert.Equal(t, int64(32), obs.MaxFinalizedBlockNumber.Val)
			})
			t.Run("if FetchInitialMaxFinalizedBlockNumber returns nil (new feed) and initialBlockNumber is set", func(t *testing.T) {
				var initialBlockNumber int64 = 50
				ds.initialBlockNumber = &initialBlockNumber
				fetcher.err = nil
				fetcher.num = nil

				obs, err := ds.Observe(ctx, repts, true)
				assert.NoError(t, err)

				assert.NoError(t, obs.MaxFinalizedBlockNumber.Err)
				assert.Equal(t, int64(49), obs.MaxFinalizedBlockNumber.Val)
			})
			t.Run("if FetchInitialMaxFinalizedBlockNumber returns nil (new feed) and initialBlockNumber is not set", func(t *testing.T) {
				ds.initialBlockNumber = nil
				t.Run("if current block num is valid", func(t *testing.T) {
					fetcher.err = nil
					fetcher.num = nil

					obs, err := ds.Observe(ctx, repts, true)
					assert.NoError(t, err)

					assert.NoError(t, obs.MaxFinalizedBlockNumber.Err)
					assert.Equal(t, head.Number-1, obs.MaxFinalizedBlockNumber.Val)
				})
				t.Run("if no current block available", func(t *testing.T) {
					h2 := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
					h2.On("LatestChain").Return((*evmtypes.Head)(nil))
					ht.h = h2

					obs, err := ds.Observe(ctx, repts, true)
					assert.NoError(t, err)

					assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed. No initialBlockNumber was set, tried to use current block number to determine maxFinalizedBlockNumber but got error: no blocks available")
				})
			})
		})
	})

	ht.h = h

	t.Run("when fetchMaxFinalizedBlockNum=false", func(t *testing.T) {
		t.Run("when run execution fails, returns error", func(t *testing.T) {
			t.Cleanup(func() {
				runner.Err = nil
			})
			runner.Err = errors.New("run execution failed")

			_, err := ds.Observe(ctx, repts, false)
			assert.EqualError(t, err, "Observe failed while executing run: error executing run for spec ID 0: run execution failed")
		})
		t.Run("makes observation using pipeline, when all tasks succeed", func(t *testing.T) {
			obs, err := ds.Observe(ctx, repts, false)
			assert.NoError(t, err)

			assert.Equal(t, big.NewInt(122), obs.BenchmarkPrice.Val)
			assert.NoError(t, obs.BenchmarkPrice.Err)
			assert.Equal(t, big.NewInt(121), obs.Bid.Val)
			assert.NoError(t, obs.Bid.Err)
			assert.Equal(t, big.NewInt(123), obs.Ask.Val)
			assert.NoError(t, obs.Ask.Err)
			assert.Equal(t, head.Number, obs.CurrentBlockNum.Val)
			assert.NoError(t, obs.CurrentBlockNum.Err)
			assert.Equal(t, fmt.Sprintf("%x", head.Hash), fmt.Sprintf("%x", obs.CurrentBlockHash.Val))
			assert.NoError(t, obs.CurrentBlockHash.Err)
			assert.Equal(t, uint64(head.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)
			assert.NoError(t, obs.CurrentBlockTimestamp.Err)

			assert.Zero(t, obs.MaxFinalizedBlockNumber.Val)
			assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "fetchMaxFinalizedBlockNum=false")
		})
		t.Run("makes observation using pipeline, with erroring tasks", func(t *testing.T) {
			for i := range trrs {
				trrs[i].Result.Error = fmt.Errorf("task error %d", i)
			}

			obs, err := ds.Observe(ctx, repts, false)
			assert.NoError(t, err)

			assert.Zero(t, obs.BenchmarkPrice.Val)
			assert.EqualError(t, obs.BenchmarkPrice.Err, "task error 0")
			assert.Zero(t, obs.Bid.Val)
			assert.EqualError(t, obs.Bid.Err, "task error 1")
			assert.Zero(t, obs.Ask.Val)
			assert.EqualError(t, obs.Ask.Err, "task error 2")
			assert.Equal(t, head.Number, obs.CurrentBlockNum.Val)
			assert.NoError(t, obs.CurrentBlockNum.Err)
			assert.Equal(t, fmt.Sprintf("%x", head.Hash), fmt.Sprintf("%x", obs.CurrentBlockHash.Val))
			assert.NoError(t, obs.CurrentBlockHash.Err)
			assert.Equal(t, uint64(head.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)
			assert.NoError(t, obs.CurrentBlockTimestamp.Err)

			assert.Zero(t, obs.MaxFinalizedBlockNumber.Val)
			assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "fetchMaxFinalizedBlockNum=false")
		})
		t.Run("makes partial observation using pipeline, if only some results have errored", func(t *testing.T) {
			trrs[0].Result.Error = fmt.Errorf("task failed")
			trrs[1].Result.Value = "33"
			trrs[1].Result.Error = nil
			trrs[2].Result.Value = nil
			trrs[2].Result.Error = fmt.Errorf("task failed")

			obs, err := ds.Observe(ctx, repts, false)
			assert.NoError(t, err)

			assert.Zero(t, obs.BenchmarkPrice.Val)
			assert.EqualError(t, obs.BenchmarkPrice.Err, "task failed")
			assert.Equal(t, big.NewInt(33), obs.Bid.Val)
			assert.NoError(t, obs.Bid.Err)
			assert.Zero(t, obs.Ask.Val)
			assert.EqualError(t, obs.Ask.Err, "task failed")
		})
		t.Run("returns error if at least one result is unparseable", func(t *testing.T) {
			trrs[0].Result.Error = fmt.Errorf("task failed")
			trrs[1].Result.Value = "foo"
			trrs[1].Result.Error = nil
			trrs[2].Result.Value = "123456"
			trrs[2].Result.Error = nil

			_, err := ds.Observe(ctx, repts, false)
			assert.EqualError(t, err, "Observe failed while parsing run results: failed to parse Bid: can't convert foo to decimal")
		})
		t.Run("sends run to runResults channel", func(t *testing.T) {
			for i := range trrs {
				trrs[i].Result.Value = "123"
				trrs[i].Result.Error = nil
			}
			ch := make(chan *pipeline.Run, 1)

			ds.runResults = ch

			_, err := ds.Observe(ctx, repts, false)
			assert.NoError(t, err)

			select {
			case run := <-ch:
				assert.Equal(t, int64(42), run.ID)
			default:
				t.Fatal("expected run on channel")
			}
		})
	})

	t.Run("LatestBlocks is populated correctly", func(t *testing.T) {
		t.Run("when chain length is zero", func(t *testing.T) {
			ht2 := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
			ht2.On("LatestChain").Return((*evmtypes.Head)(nil))
			ht.h = ht2

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.Len(t, obs.LatestBlocks, 0)

			ht2.AssertExpectations(t)
		})
		t.Run("when chain is too short", func(t *testing.T) {
			h4 := &evmtypes.Head{
				Number: 4,
				Parent: nil,
			}
			h5 := &evmtypes.Head{
				Number: 5,
				Parent: h4,
			}
			h6 := &evmtypes.Head{
				Number: 6,
				Parent: h5,
			}

			ht2 := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
			ht2.On("LatestChain").Return(h6)
			ht.h = ht2

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.Len(t, obs.LatestBlocks, 3)
			assert.Equal(t, 6, int(obs.LatestBlocks[0].Num))
			assert.Equal(t, 5, int(obs.LatestBlocks[1].Num))
			assert.Equal(t, 4, int(obs.LatestBlocks[2].Num))

			ht2.AssertExpectations(t)
		})
		t.Run("when chain is long enough", func(t *testing.T) {
			h1 := &evmtypes.Head{
				Number: 1,
			}
			h2 := &evmtypes.Head{
				Number: 2,
				Parent: h1,
			}
			h3 := &evmtypes.Head{
				Number: 3,
				Parent: h2,
			}
			h4 := &evmtypes.Head{
				Number: 4,
				Parent: h3,
			}
			h5 := &evmtypes.Head{
				Number: 5,
				Parent: h4,
			}
			h6 := &evmtypes.Head{
				Number: 6,
				Parent: h5,
			}

			ht2 := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
			ht2.On("LatestChain").Return(h6)
			ht.h = ht2

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.Len(t, obs.LatestBlocks, 5)
			assert.Equal(t, 6, int(obs.LatestBlocks[0].Num))
			assert.Equal(t, 5, int(obs.LatestBlocks[1].Num))
			assert.Equal(t, 4, int(obs.LatestBlocks[2].Num))
			assert.Equal(t, 3, int(obs.LatestBlocks[3].Num))
			assert.Equal(t, 2, int(obs.LatestBlocks[4].Num))

			ht2.AssertExpectations(t)
		})
	})
}

func TestMercury_SetLatestBlocks(t *testing.T) {
	lggr := logger.TestLogger(t)
	ds := NewDataSource(nil, nil, job.Job{}, pipeline.Spec{}, lggr, nil, nil, nil, nil, nil, mercuryutils.FeedID{})

	h := evmtypes.Head{
		Number:           testutils.NewRandomPositiveInt64(),
		Hash:             utils.NewHash(),
		ParentHash:       utils.NewHash(),
		Timestamp:        time.Now(),
		BaseFeePerGas:    assets.NewWeiI(testutils.NewRandomPositiveInt64()),
		ReceiptsRoot:     utils.NewHash(),
		TransactionsRoot: utils.NewHash(),
		StateRoot:        utils.NewHash(),
	}

	t.Run("returns head from headtracker if present", func(t *testing.T) {
		headTracker := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("HeadTracker").Return(headTracker)
		headTracker.On("LatestChain").Return(&h, nil)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercuryv1.Observation{}
		ds.setLatestBlocks(context.Background(), &obs)

		assert.Equal(t, h.Number, obs.CurrentBlockNum.Val)
		assert.Equal(t, h.Hash.Bytes(), obs.CurrentBlockHash.Val)
		assert.Equal(t, uint64(h.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)

		assert.Len(t, obs.LatestBlocks, 1)

		chainHeadTracker.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head", func(t *testing.T) {
		headTracker := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return((*evmtypes.Head)(nil))

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercuryv1.Observation{}
		ds.setLatestBlocks(context.Background(), &obs)

		assert.Zero(t, obs.CurrentBlockNum.Val)
		assert.Zero(t, obs.CurrentBlockHash.Val)
		assert.Zero(t, obs.CurrentBlockTimestamp.Val)
		assert.EqualError(t, obs.CurrentBlockNum.Err, "no blocks available")
		assert.EqualError(t, obs.CurrentBlockHash.Err, "no blocks available")
		assert.EqualError(t, obs.CurrentBlockTimestamp.Err, "no blocks available")

		assert.Len(t, obs.LatestBlocks, 0)

		chainHeadTracker.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})
}

var sampleFeedID = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}

func buildSampleV1Report() []byte {
	feedID := sampleFeedID
	timestamp := uint32(42)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(143)
	currentBlockHash := utils.NewHash()
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(142)

	b, err := reportcodecv1.ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, currentBlockTimestamp, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}
