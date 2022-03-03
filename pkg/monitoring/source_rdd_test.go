package monitoring

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
		source := NewRDDSource(srv.URL, fakeFeedsParser, newNullLogger())
		rawFeeds, err := source.Fetch(context.Background())
		require.NoError(t, err)
		feeds, ok := rawFeeds.([]FeedConfig)
		require.True(t, ok)
		require.Len(t, feeds, 4)
		for _, feed := range feeds {
			require.NotEqual(t, "dead", feed.GetContractStatus())
		}
	})
}
