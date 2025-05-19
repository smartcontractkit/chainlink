package cache

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

const neverExpireTTL = 1000 * time.Hour // some massive value that will never expire during a test

func Test_Cache(t *testing.T) {
	lggr := logger.Test(t)
	client := &mockClient{}
	cfg := Config{}
	ctx := testutils.Context(t)

	req1 := &pb.LatestReportRequest{FeedId: []byte{1}}
	req2 := &pb.LatestReportRequest{FeedId: []byte{2}}
	req3 := &pb.LatestReportRequest{FeedId: []byte{3}}

	feedID1Hex := mercuryutils.BytesToFeedID(req1.FeedId).String()

	t.Run("errors with nil req", func(t *testing.T) {
		c := newMemCache(lggr, client, cfg)

		_, err := c.LatestReport(ctx, nil)
		assert.EqualError(t, err, "req must not be nil")
	})

	t.Run("with LatestReportTTL=0 does no caching", func(t *testing.T) {
		c := newMemCache(lggr, client, cfg)

		req := &pb.LatestReportRequest{}
		for i := 0; i < 5; i++ {
			client.resp = &pb.LatestReportResponse{Report: &pb.Report{Price: []byte(strconv.Itoa(i))}}

			resp, err := c.LatestReport(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, client.resp, resp)
		}

		client.resp = nil
		client.err = errors.New("something exploded")

		resp, err := c.LatestReport(ctx, req)
		assert.EqualError(t, err, "something exploded")
		assert.Nil(t, resp)
	})

	t.Run("caches repeated calls to LatestReport, keyed by request", func(t *testing.T) {
		cfg.LatestReportTTL = neverExpireTTL
		client.err = nil
		c := newMemCache(lggr, client, cfg)

		t.Run("if cache is unstarted, returns error", func(t *testing.T) {
			// starting the cache is required for state management if we
			// actually cache results, since fetches are initiated async and
			// need to be cleaned up properly on close
			_, err := c.LatestReport(ctx, &pb.LatestReportRequest{})
			assert.EqualError(t, err, "memCache must be started, but is: Unstarted")
		})

		err := c.StartOnce("test start", func() error { return nil })
		require.NoError(t, err)

		t.Run("returns cached value for key", func(t *testing.T) {
			var firstResp *pb.LatestReportResponse
			for i := 0; i < 5; i++ {
				client.resp = &pb.LatestReportResponse{Report: &pb.Report{Price: []byte(strconv.Itoa(i))}}
				if firstResp == nil {
					firstResp = client.resp
				}

				resp, err := c.LatestReport(ctx, req1)
				require.NoError(t, err)
				assert.Equal(t, firstResp, resp)
			}
		})

		t.Run("cache keys do not conflict", func(t *testing.T) {
			var firstResp1 *pb.LatestReportResponse
			for i := 5; i < 10; i++ {
				client.resp = &pb.LatestReportResponse{Report: &pb.Report{Price: []byte(strconv.Itoa(i))}}
				if firstResp1 == nil {
					firstResp1 = client.resp
				}

				resp, err := c.LatestReport(ctx, req2)
				require.NoError(t, err)
				assert.Equal(t, firstResp1, resp)
			}

			var firstResp2 *pb.LatestReportResponse
			for i := 10; i < 15; i++ {
				client.resp = &pb.LatestReportResponse{Report: &pb.Report{Price: []byte(strconv.Itoa(i))}}
				if firstResp2 == nil {
					firstResp2 = client.resp
				}

				resp, err := c.LatestReport(ctx, req3)
				require.NoError(t, err)
				assert.Equal(t, firstResp2, resp)
			}

			// req1 key still has same value
			resp, err := c.LatestReport(ctx, req1)
			require.NoError(t, err)
			assert.Equal(t, []byte(strconv.Itoa(0)), resp.Report.Price)

			// req2 key still has same value
			resp, err = c.LatestReport(ctx, req2)
			require.NoError(t, err)
			assert.Equal(t, []byte(strconv.Itoa(5)), resp.Report.Price)
		})

		t.Run("re-queries when a cache item has expired", func(t *testing.T) {
			vi, exists := c.cache.Load(feedID1Hex)
			require.True(t, exists)
			v := vi.(*cacheVal)
			v.expiresAt = time.Now().Add(-1 * time.Second)

			client.resp = &pb.LatestReportResponse{Report: &pb.Report{Price: []byte(strconv.Itoa(15))}}

			resp, err := c.LatestReport(ctx, req1)
			require.NoError(t, err)
			assert.Equal(t, client.resp, resp)

			// querying again yields the same cached item
			resp, err = c.LatestReport(ctx, req1)
			require.NoError(t, err)
			assert.Equal(t, client.resp, resp)
		})
	})

	t.Run("complete fetch", func(t *testing.T) {
		t.Run("does not change expiry if fetch returns error", func(t *testing.T) {
			expires := time.Now().Add(-1 * time.Second)
			v := &cacheVal{
				fetching:  true,
				fetchCh:   make(chan (struct{})),
				val:       nil,
				err:       nil,
				expiresAt: expires,
			}
			v.completeFetch(nil, errors.New("foo"), time.Now().Add(neverExpireTTL))
			assert.Equal(t, expires, v.expiresAt)

			v = &cacheVal{
				fetching:  true,
				fetchCh:   make(chan (struct{})),
				val:       nil,
				err:       nil,
				expiresAt: expires,
			}
			expires = time.Now().Add(neverExpireTTL)
			v.completeFetch(nil, nil, expires)
			assert.Equal(t, expires, v.expiresAt)
		})
	})

	t.Run("timeouts", func(t *testing.T) {
		c := newMemCache(lggr, client, cfg)
		// simulate fetch already executing in background
		v := &cacheVal{
			fetching:  true,
			fetchCh:   make(chan (struct{})),
			val:       nil,
			err:       nil,
			expiresAt: time.Now().Add(-1 * time.Second),
		}
		c.cache.Store(feedID1Hex, v)

		canceledCtx, cancel := context.WithCancel(testutils.Context(t))
		cancel()

		t.Run("returns context deadline exceeded error if fetch takes too long", func(t *testing.T) {
			_, err := c.LatestReport(canceledCtx, req1)
			require.Error(t, err)
			assert.True(t, errors.Is(err, context.Canceled))
			assert.EqualError(t, err, "context canceled")
		})
		t.Run("returns wrapped context deadline exceeded error if fetch has errored and is in the retry loop", func(t *testing.T) {
			v.err = errors.New("some background fetch error")

			_, err := c.LatestReport(canceledCtx, req1)
			require.Error(t, err)
			assert.True(t, errors.Is(err, context.Canceled))
			assert.EqualError(t, err, "some background fetch error\ncontext canceled")
		})
	})
}
