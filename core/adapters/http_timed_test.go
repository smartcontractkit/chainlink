// +build probabilistic_tests

package adapters_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPGet_TimeoutAllowsRetries(t *testing.T) {
	store := leanStore()
	store.Config.Set("DEFAULT_HTTP_TIMEOUT", "80ms")
	store.Config.Set("MAX_HTTP_ATTEMPTS", "2")

	attempts := make(chan struct{}, 2)
	timeoutOnce := sync.Once{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Greater(t, len(b), 0)
		attempts <- struct{}{}
		timeoutOnce.Do(func() { time.Sleep(100 * time.Millisecond) })
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	hga := adapters.HTTPPost{
		URL:                            cltest.WebURL(t, server.URL),
		AllowUnrestrictedNetworkAccess: true,
	}

	input := cltest.NewRunInputWithResult("inputValue")
	result := hga.Perform(input, store, nil)
	require.NoError(t, result.Error())

	for i := 0; i < 2; i++ {
		select {
		case <-attempts:
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for attempt %v", i)
		}
	}
}
