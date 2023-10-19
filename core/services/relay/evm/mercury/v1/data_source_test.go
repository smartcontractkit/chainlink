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
	"github.com/stretchr/testify/mock"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaymercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"
	commonmocks "github.com/smartcontractkit/chainlink/v2/common/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	mercurymocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
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
	c evmclient.Client
	h httypes.HeadTracker
}

func (m *mockHeadTracker) Client() evmclient.Client         { return m.c }
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
	ds := &datasource{lggr: logger.TestLogger(t), orm: orm, codec: (reportcodecv1.ReportCodec{})}
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
	c := evmclimocks.NewClient(t)
	ht := &mockHeadTracker{
		c: c,
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
				t.Run("if current block num errored", func(t *testing.T) {
					h2 := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
					h2.On("LatestChain").Return((*evmtypes.Head)(nil))
					ht.h = h2
					c2 := evmclimocks.NewClient(t)
					c2.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("head retrieval failed"))
					ht.c = c2

					obs, err := ds.Observe(ctx, repts, true)
					assert.NoError(t, err)

					assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed. No initialBlockNumber was set, tried to use current block number to determine maxFinalizedBlockNumber but got error: head retrieval failed")
				})
			})
		})
	})

	ht.h = h
	ht.c = c

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
		t.Run("if head tracker returns nil, falls back to RPC method", func(t *testing.T) {
			t.Run("if call succeeds", func(t *testing.T) {
				h = commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
				h.On("LatestChain").Return((*evmtypes.Head)(nil))
				ht.h = h
				c.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head, nil).Once()

				obs, err := ds.Observe(ctx, repts, false)
				assert.NoError(t, err)

				assert.Equal(t, head.Number, obs.CurrentBlockNum.Val)
				assert.NoError(t, obs.CurrentBlockNum.Err)
				assert.Equal(t, fmt.Sprintf("%x", head.Hash), fmt.Sprintf("%x", obs.CurrentBlockHash.Val))
				assert.NoError(t, obs.CurrentBlockHash.Err)
				assert.Equal(t, uint64(head.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)
				assert.NoError(t, obs.CurrentBlockTimestamp.Err)

				h.AssertExpectations(t)
				c.AssertExpectations(t)
			})
			t.Run("if call fails, returns error for that observation", func(t *testing.T) {
				c = evmclimocks.NewClient(t)
				ht.c = c
				c.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("client call failed")).Once()

				obs, err := ds.Observe(ctx, repts, false)
				assert.NoError(t, err)

				assert.Zero(t, obs.CurrentBlockNum.Val)
				assert.EqualError(t, obs.CurrentBlockNum.Err, "client call failed")
				assert.Zero(t, obs.CurrentBlockHash.Val)
				assert.EqualError(t, obs.CurrentBlockHash.Err, "client call failed")
				assert.Zero(t, obs.CurrentBlockTimestamp.Val)
				assert.EqualError(t, obs.CurrentBlockTimestamp.Err, "client call failed")

				c.AssertExpectations(t)
			})
		})
	})
}

func TestMercury_SetCurrentBlock(t *testing.T) {
	lggr := logger.TestLogger(t)
	ds := datasource{
		lggr: lggr,
	}

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
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, h.Number, obs.CurrentBlockNum.Val)
		assert.Equal(t, h.Hash.Bytes(), obs.CurrentBlockHash.Val)
		assert.Equal(t, uint64(h.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)

		chainHeadTracker.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call succeeds", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("Client").Return(ethClient)
		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return((*evmtypes.Head)(nil))
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&h, nil)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercuryv1.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, h.Number, obs.CurrentBlockNum.Val)
		assert.Equal(t, h.Hash.Bytes(), obs.CurrentBlockHash.Val)
		assert.Equal(t, uint64(h.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)

		chainHeadTracker.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call fails", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := commonmocks.NewHeadTracker[*evmtypes.Head, common.Hash](t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("Client").Return(ethClient)
		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return((*evmtypes.Head)(nil))
		err := errors.New("foo")
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, err)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercuryv1.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, err, obs.CurrentBlockNum.Err)
		assert.Equal(t, err, obs.CurrentBlockHash.Err)
		assert.Equal(t, err, obs.CurrentBlockTimestamp.Err)

		chainHeadTracker.AssertExpectations(t)
		ethClient.AssertExpectations(t)
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
