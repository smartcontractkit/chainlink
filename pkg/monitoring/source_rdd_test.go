package monitoring

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/monitoring/config"
)

func TestRDDSource(t *testing.T) {
	t.Run("should filter out dead feeds", func(t *testing.T) {
		srv := serveJSON(t, "./fixtures/feeds.json")
		defer srv.Close()
		source := NewRDDSource(srv.URL, fakeFeedsParser, []string{}, "no-nodes", fakeNodesParser, newNullLogger()).(*rddSource)
		feeds, err := source.fetchFeeds(context.Background())
		require.NoError(t, err)
		require.Len(t, feeds, 4)
		for _, feed := range feeds {
			require.NotEqual(t, "dead", feed.GetContractStatus())
		}
	})
	t.Run("should filter out feeds that show up in FEEDS_IGNORE_IDS", func(t *testing.T) {
		srv := serveJSON(t, "./fixtures/feeds.json")
		defer srv.Close()
		// Build a source.
		os.Setenv("FEEDS_IGNORE_IDS", "HW3ipKzeeduJq6f1NqRCw4doknMeWkfrM4WxobtG3o5c, HW3ipKzeeduJq6f1NqRCw4doknMeWkfrM4WxobtG3o5d")
		defer os.Unsetenv("FEEDS_IGNORE_IDS")
		cfg, _ := config.Parse() // NOTE: purposefully ignoring config validation errors.
		source := NewRDDSource(srv.URL, fakeFeedsParser, cfg.Feeds.IgnoreIDs, "no-nodes", fakeNodesParser, newNullLogger()).(*rddSource)
		// Fetch feeds from fake RDD.
		feeds, err := source.fetchFeeds(context.Background())
		require.NoError(t, err)
		require.Len(t, feeds, 2)
		for _, feed := range feeds {
			require.NotEqual(t, "dead", feed.GetContractStatus())
			require.NotContains(t, cfg.Feeds.IgnoreIDs, feed.GetID())
		}
	})
	t.Run("should fetch feeds and nodes data", func(t *testing.T) {
		feedsSrv := serveJSON(t, "./fixtures/feeds.json")
		defer feedsSrv.Close()
		nodesSrv := serveJSON(t, "./fixtures/nodes.json")
		defer nodesSrv.Close()

		cfg := config.Config{}
		source := NewRDDSource(
			feedsSrv.URL, fakeFeedsParser, cfg.Feeds.IgnoreIDs,
			nodesSrv.URL, fakeNodesParser,
			newNullLogger(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		rawData, err := source.Fetch(ctx)
		require.NoError(t, err)
		data, ok := rawData.(RDDData)
		require.True(t, ok)
		require.Len(t, data.Feeds, 4)
		require.Len(t, data.Nodes, 2)
	})
}

// Helpers

func serveJSON(t *testing.T, path string) *httptest.Server {
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-type", "application/json")
		_, err := w.Write(data)
		require.NoError(t, err)
	}))
}
