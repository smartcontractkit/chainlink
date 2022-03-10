package monitoring

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
	"github.com/stretchr/testify/require"
)

func TestRDDSource(t *testing.T) {
	t.Run("should filter out dead feeds", func(t *testing.T) {
		data, err := os.ReadFile("./fixtures/feeds.json")
		require.NoError(t, err)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-type", "application/json")
			_, err := w.Write(data)
			require.NoError(t, err)
		}))
		defer srv.Close()
		source := NewRDDSource(srv.URL, fakeFeedsParser, newNullLogger(), []string{})
		rawFeeds, err := source.Fetch(context.Background())
		require.NoError(t, err)
		feeds, ok := rawFeeds.([]FeedConfig)
		require.True(t, ok)
		require.Len(t, feeds, 4)
		for _, feed := range feeds {
			require.NotEqual(t, "dead", feed.GetContractStatus())
		}
	})
	t.Run("should filter out feeds that show up in FEEDS_IGNORE_IDS", func(t *testing.T) {
		// Configure fake weiwatchers server.
		data, err := os.ReadFile("./fixtures/feeds.json")
		require.NoError(t, err)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-type", "application/json")
			_, err := w.Write(data)
			require.NoError(t, err)
		}))
		defer srv.Close()
		// Build a source.
		os.Setenv("FEEDS_IGNORE_IDS", "HW3ipKzeeduJq6f1NqRCw4doknMeWkfrM4WxobtG3o5c, HW3ipKzeeduJq6f1NqRCw4doknMeWkfrM4WxobtG3o5d")
		defer os.Unsetenv("FEEDS_IGNORE_IDS")
		cfg, _ := config.Parse() // NOTE: purposefully ignoring config validation errors.
		source := NewRDDSource(srv.URL, fakeFeedsParser, newNullLogger(), cfg.Feeds.IgnoreIDs)
		// Fetch feeds from fake RDD.
		rawFeeds, err := source.Fetch(context.Background())
		require.NoError(t, err)
		feeds, ok := rawFeeds.([]FeedConfig)
		require.True(t, ok)
		require.Len(t, feeds, 2)
		for _, feed := range feeds {
			require.NotEqual(t, "dead", feed.GetContractStatus())
			require.NotContains(t, cfg.Feeds.IgnoreIDs, feed.GetID())
		}
	})
}
