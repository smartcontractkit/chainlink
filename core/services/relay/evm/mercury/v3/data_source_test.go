package v3

import (
	"context"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	mercurytypes "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	relaymercuryv3 "github.com/smartcontractkit/chainlink-data-streams/mercury/v3"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	mercurymocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

var _ mercurytypes.ServerFetcher = &mockFetcher{}

type mockFetcher struct {
	ts             int64
	tsErr          error
	linkPrice      *big.Int
	linkPriceErr   error
	nativePrice    *big.Int
	nativePriceErr error
}

var feedId utils.FeedID = [32]byte{1}
var linkFeedId utils.FeedID = [32]byte{2}
var nativeFeedId utils.FeedID = [32]byte{3}

func (m *mockFetcher) FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error) {
	return nil, nil
}

func (m *mockFetcher) LatestPrice(ctx context.Context, fId [32]byte) (*big.Int, error) {
	if fId == linkFeedId {
		return m.linkPrice, m.linkPriceErr
	} else if fId == nativeFeedId {
		return m.nativePrice, m.nativePriceErr
	}
	return nil, nil
}

func (m *mockFetcher) LatestTimestamp(context.Context) (int64, error) {
	return m.ts, m.tsErr
}

type mockORM struct {
	report []byte
	err    error
}

func (m *mockORM) LatestReport(ctx context.Context, feedID [32]byte) (report []byte, err error) {
	return m.report, m.err
}

type mockSaver struct {
	r *pipeline.Run
}

func (ms *mockSaver) Save(r *pipeline.Run) {
	ms.r = r
}

func Test_Datasource(t *testing.T) {
	orm := &mockORM{}
	ds := &datasource{orm: orm, lggr: logger.TestLogger(t)}
	ctx := testutils.Context(t)
	repts := ocrtypes.ReportTimestamp{}

	fetcher := &mockFetcher{}
	ds.fetcher = fetcher

	saver := &mockSaver{}
	ds.saver = saver

	goodTrrs := []pipeline.TaskRunResult{
		{
			// bp
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

	ds.pipelineRunner = &mercurymocks.MockRunner{
		Trrs: goodTrrs,
	}

	spec := pipeline.Spec{}
	ds.spec = spec

	t.Run("when fetchMaxFinalizedTimestamp=true", func(t *testing.T) {
		t.Run("with latest report in database", func(t *testing.T) {
			orm.report = buildSampleV3Report()
			orm.err = nil

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.NoError(t, obs.MaxFinalizedTimestamp.Err)
			assert.Equal(t, int64(124), obs.MaxFinalizedTimestamp.Val)
		})
		t.Run("if querying latest report fails", func(t *testing.T) {
			orm.report = nil
			orm.err = errors.New("something exploded")

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.EqualError(t, obs.MaxFinalizedTimestamp.Err, "something exploded")
			assert.Zero(t, obs.MaxFinalizedTimestamp.Val)
		})
		t.Run("if codec fails to decode", func(t *testing.T) {
			orm.report = []byte{1, 2, 3}
			orm.err = nil

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.EqualError(t, obs.MaxFinalizedTimestamp.Err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")
			assert.Zero(t, obs.MaxFinalizedTimestamp.Val)
		})

		orm.report = nil
		orm.err = nil

		t.Run("if LatestTimestamp returns error", func(t *testing.T) {
			fetcher.tsErr = errors.New("some error")

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.EqualError(t, obs.MaxFinalizedTimestamp.Err, "some error")
			assert.Zero(t, obs.MaxFinalizedTimestamp.Val)
		})

		t.Run("if LatestTimestamp succeeds", func(t *testing.T) {
			fetcher.tsErr = nil
			fetcher.ts = 123

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.Equal(t, int64(123), obs.MaxFinalizedTimestamp.Val)
			assert.NoError(t, obs.MaxFinalizedTimestamp.Err)
		})

		t.Run("if LatestTimestamp succeeds but ts=0 (new feed)", func(t *testing.T) {
			fetcher.tsErr = nil
			fetcher.ts = 0

			obs, err := ds.Observe(ctx, repts, true)
			assert.NoError(t, err)

			assert.NoError(t, obs.MaxFinalizedTimestamp.Err)
			assert.Zero(t, obs.MaxFinalizedTimestamp.Val)
		})

		t.Run("when run execution succeeded", func(t *testing.T) {
			t.Run("when feedId=linkFeedID=nativeFeedId", func(t *testing.T) {
				t.Cleanup(func() {
					ds.feedID, ds.linkFeedID, ds.nativeFeedID = feedId, linkFeedId, nativeFeedId
				})

				ds.feedID, ds.linkFeedID, ds.nativeFeedID = feedId, feedId, feedId

				fetcher.ts = 123123
				fetcher.tsErr = nil

				obs, err := ds.Observe(ctx, repts, true)
				assert.NoError(t, err)

				assert.Equal(t, big.NewInt(122), obs.BenchmarkPrice.Val)
				assert.NoError(t, obs.BenchmarkPrice.Err)
				assert.Equal(t, big.NewInt(121), obs.Bid.Val)
				assert.NoError(t, obs.Bid.Err)
				assert.Equal(t, big.NewInt(123), obs.Ask.Val)
				assert.NoError(t, obs.Ask.Err)
				assert.Equal(t, int64(123123), obs.MaxFinalizedTimestamp.Val)
				assert.NoError(t, obs.MaxFinalizedTimestamp.Err)
				assert.Equal(t, big.NewInt(122), obs.LinkPrice.Val)
				assert.NoError(t, obs.LinkPrice.Err)
				assert.Equal(t, big.NewInt(122), obs.NativePrice.Val)
				assert.NoError(t, obs.NativePrice.Err)
			})
		})
	})

	t.Run("when fetchMaxFinalizedTimestamp=false", func(t *testing.T) {
		t.Run("when run execution fails, returns error", func(t *testing.T) {
			t.Cleanup(func() {
				ds.pipelineRunner = &mercurymocks.MockRunner{
					Trrs: goodTrrs,
					Err:  nil,
				}
			})

			ds.pipelineRunner = &mercurymocks.MockRunner{
				Trrs: goodTrrs,
				Err:  errors.New("run execution failed"),
			}

			_, err := ds.Observe(ctx, repts, false)
			assert.EqualError(t, err, "Observe failed while executing run: error executing run for spec ID 0: run execution failed")
		})

		t.Run("when parsing run results fails, return error", func(t *testing.T) {
			t.Cleanup(func() {
				runner := &mercurymocks.MockRunner{
					Trrs: goodTrrs,
					Err:  nil,
				}
				ds.pipelineRunner = runner
			})

			badTrrs := []pipeline.TaskRunResult{
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
					Result: pipeline.Result{Error: errors.New("some error with ask")},
					Task:   &mercurymocks.MockTask{},
				},
			}

			ds.pipelineRunner = &mercurymocks.MockRunner{
				Trrs: badTrrs,
				Err:  nil,
			}

			_, err := ds.Observe(ctx, repts, false)
			assert.EqualError(t, err, "Observe failed while parsing run results: some error with ask")
		})

		t.Run("when run execution succeeded", func(t *testing.T) {
			t.Run("when feedId=linkFeedID=nativeFeedId", func(t *testing.T) {
				t.Cleanup(func() {
					ds.feedID, ds.linkFeedID, ds.nativeFeedID = feedId, linkFeedId, nativeFeedId
				})

				var feedId utils.FeedID = [32]byte{1}
				ds.feedID, ds.linkFeedID, ds.nativeFeedID = feedId, feedId, feedId

				obs, err := ds.Observe(ctx, repts, false)
				assert.NoError(t, err)

				assert.Equal(t, big.NewInt(122), obs.BenchmarkPrice.Val)
				assert.NoError(t, obs.BenchmarkPrice.Err)
				assert.Equal(t, big.NewInt(121), obs.Bid.Val)
				assert.NoError(t, obs.Bid.Err)
				assert.Equal(t, big.NewInt(123), obs.Ask.Val)
				assert.NoError(t, obs.Ask.Err)
				assert.Equal(t, int64(0), obs.MaxFinalizedTimestamp.Val)
				assert.NoError(t, obs.MaxFinalizedTimestamp.Err)
				assert.Equal(t, big.NewInt(122), obs.LinkPrice.Val)
				assert.NoError(t, obs.LinkPrice.Err)
				assert.Equal(t, big.NewInt(122), obs.NativePrice.Val)
				assert.NoError(t, obs.NativePrice.Err)
			})

			t.Run("when fails to fetch linkPrice or nativePrice", func(t *testing.T) {
				t.Cleanup(func() {
					fetcher.linkPriceErr = nil
					fetcher.nativePriceErr = nil
				})

				fetcher.linkPriceErr = errors.New("some error fetching link price")
				fetcher.nativePriceErr = errors.New("some error fetching native price")

				obs, err := ds.Observe(ctx, repts, false)
				assert.NoError(t, err)

				assert.Nil(t, obs.LinkPrice.Val)
				assert.EqualError(t, obs.LinkPrice.Err, "some error fetching link price")
				assert.Nil(t, obs.NativePrice.Val)
				assert.EqualError(t, obs.NativePrice.Err, "some error fetching native price")
			})

			t.Run("when succeeds to fetch linkPrice or nativePrice but got nil (new feed)", func(t *testing.T) {
				obs, err := ds.Observe(ctx, repts, false)
				assert.NoError(t, err)

				assert.Equal(t, obs.LinkPrice.Val, relaymercuryv3.MissingPrice)
				assert.Nil(t, obs.LinkPrice.Err)
				assert.Equal(t, obs.NativePrice.Val, relaymercuryv3.MissingPrice)
				assert.Nil(t, obs.NativePrice.Err)
			})
		})
	})
}

var sampleFeedID = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}

func buildSampleV3Report() []byte {
	feedID := sampleFeedID
	timestamp := uint32(124)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	validFromTimestamp := uint32(123)
	expiresAt := uint32(456)
	linkFee := big.NewInt(3334455)
	nativeFee := big.NewInt(556677)

	b, err := reportcodecv3.ReportTypes.Pack(feedID, validFromTimestamp, timestamp, nativeFee, linkFee, expiresAt, bp, bid, ask)
	if err != nil {
		panic(err)
	}
	return b
}
