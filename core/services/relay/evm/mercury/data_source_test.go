package mercury

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	htmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/mocks"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	mercurymocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ relaymercury.Fetcher = &mockFetcher{}

type mockFetcher struct {
	num *int64
	err error
}

func (m *mockFetcher) FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error) {
	return m.num, m.err
}

var _ Runner = &mockRunner{}

type mockRunner struct {
	trrs pipeline.TaskRunResults
	err  error
}

func (m *mockRunner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run pipeline.Run, trrs pipeline.TaskRunResults, err error) {
	return pipeline.Run{ID: 42}, m.trrs, m.err
}

var _ pipeline.Task = &mockTask{}

type mockTask struct {
	result pipeline.Result
}

func (m *mockTask) Type() pipeline.TaskType { return "MockTask" }
func (m *mockTask) ID() int                 { return 0 }
func (m *mockTask) DotID() string           { return "" }
func (m *mockTask) Run(ctx context.Context, lggr logger.Logger, vars pipeline.Vars, inputs []pipeline.Result) (pipeline.Result, pipeline.RunInfo) {
	return m.result, pipeline.RunInfo{}
}
func (m *mockTask) Base() *pipeline.BaseTask           { return nil }
func (m *mockTask) Outputs() []pipeline.Task           { return nil }
func (m *mockTask) Inputs() []pipeline.TaskDependency  { return nil }
func (m *mockTask) OutputIndex() int32                 { return 0 }
func (m *mockTask) TaskTimeout() (time.Duration, bool) { return 0, false }
func (m *mockTask) TaskRetries() uint32                { return 0 }
func (m *mockTask) TaskMinBackoff() time.Duration      { return 0 }
func (m *mockTask) TaskMaxBackoff() time.Duration      { return 0 }

var _ ChainHeadTracker = &mockHeadTracker{}

type mockHeadTracker struct {
	c evmclient.Client
	h httypes.HeadTracker
}

func (m *mockHeadTracker) Client() evmclient.Client         { return m.c }
func (m *mockHeadTracker) HeadTracker() httypes.HeadTracker { return m.h }

func TestMercury_Observe(t *testing.T) {
	ds := &datasource{lggr: logger.TestLogger(t)}
	ctx := testutils.Context(t)
	repts := ocrtypes.ReportTimestamp{}

	fetcher := &mockFetcher{}
	ds.fetcher = fetcher

	trrs := []pipeline.TaskRunResult{
		pipeline.TaskRunResult{
			// benchmark price
			Result: pipeline.Result{Value: "122.345"},
			Task:   &mockTask{},
		},
		pipeline.TaskRunResult{
			// bid
			Result: pipeline.Result{Value: "121.993"},
			Task:   &mockTask{},
		},
		pipeline.TaskRunResult{
			// ask
			Result: pipeline.Result{Value: "123.111"},
			Task:   &mockTask{},
		},
	}

	runner := &mockRunner{
		trrs: trrs,
	}
	ds.pipelineRunner = runner

	spec := pipeline.Spec{}
	ds.spec = spec

	h := htmocks.NewHeadTracker(t)
	c := evmtest.NewEthClientMock(t)
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
				h2 := htmocks.NewHeadTracker(t)
				h2.On("LatestChain").Return(nil)
				ht.h = h2
				c2 := evmtest.NewEthClientMock(t)
				c2.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("head retrieval failed"))
				ht.c = c2

				obs, err := ds.Observe(ctx, repts, true)
				assert.NoError(t, err)

				assert.EqualError(t, obs.MaxFinalizedBlockNumber.Err, "FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed. No initialBlockNumber was set, tried to use current block number to determine maxFinalizedBlockNumber but got error: head retrieval failed")
			})
		})
	})

	ht.h = h
	ht.c = c

	t.Run("when fetchMaxFinalizedBlockNum=false", func(t *testing.T) {
		t.Run("when run execution fails, returns error", func(t *testing.T) {
			t.Cleanup(func() {
				runner.err = nil
			})
			runner.err = errors.New("run execution failed")

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

			ch := make(chan pipeline.Run, 1)
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
				h = htmocks.NewHeadTracker(t)
				h.On("LatestChain").Return(nil)
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
				c = evmtest.NewEthClientMock(t)
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
		headTracker := htmocks.NewHeadTracker(t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("HeadTracker").Return(headTracker)
		headTracker.On("LatestChain").Return(&h, nil)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercury.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, h.Number, obs.CurrentBlockNum.Val)
		assert.Equal(t, h.Hash.Bytes(), obs.CurrentBlockHash.Val)
		assert.Equal(t, uint64(h.Timestamp.Unix()), obs.CurrentBlockTimestamp.Val)

		chainHeadTracker.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call succeeds", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := htmocks.NewHeadTracker(t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("Client").Return(ethClient)
		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return(nil)
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&h, nil)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercury.Observation{}
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
		headTracker := htmocks.NewHeadTracker(t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("Client").Return(ethClient)
		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return(nil)
		err := errors.New("foo")
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, err)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercury.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, err, obs.CurrentBlockNum.Err)
		assert.Equal(t, err, obs.CurrentBlockHash.Err)
		assert.Equal(t, err, obs.CurrentBlockTimestamp.Err)

		chainHeadTracker.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})
}
