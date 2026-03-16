package fakes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	customhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestDirectHTTPAction_RequestHeaders(t *testing.T) {
	t.Run("MultiHeaders are sent in request", func(t *testing.T) {
		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		t.Cleanup(srv.Close)

		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:    srv.URL,
			Method: "GET",
			MultiHeaders: map[string]*customhttp.HeaderValues{
				"Authorization": {Values: []string{"Bearer test-token"}},
			},
		}
		metadata := commonCap.RequestMetadata{}

		result, err := action.SendRequest(context.Background(), metadata, input)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Bearer test-token", receivedAuth, "Authorization header should be sent")
	})

	t.Run("Headers (deprecated) are sent in request when MultiHeaders empty", func(t *testing.T) {
		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		t.Cleanup(srv.Close)

		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:     srv.URL,
			Method:  "GET",
			Headers: map[string]string{"Authorization": "Basic legacy-auth"},
		}
		metadata := commonCap.RequestMetadata{}

		result, err := action.SendRequest(context.Background(), metadata, input)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Basic legacy-auth", receivedAuth, "Authorization header should be sent via deprecated Headers")
	})
}

func TestDirectHTTPAction_ResponseHeadersAndMultiHeaders(t *testing.T) {
	t.Run("response has both Headers and MultiHeaders populated", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("Set-Cookie", "sessionid=abc123; Path=/")
			w.Header().Add("Set-Cookie", "csrf=xyz789; Path=/")
			w.Header().Add("X-Custom", "single-value")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		t.Cleanup(srv.Close)

		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:    srv.URL,
			Method: "GET",
		}
		metadata := commonCap.RequestMetadata{}

		result, err := action.SendRequest(context.Background(), metadata, input)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Response)

		resp := result.Response
		assert.Equal(t, uint32(200), resp.StatusCode)

		// Headers (comma-joined, backwards compat)
		require.NotNil(t, resp.Headers)                                    //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers, "Content-Type")                   //nolint:staticcheck // testing deprecated field
		assert.Equal(t, "application/json", resp.Headers["Content-Type"])  //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers, "Set-Cookie")                     //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers["Set-Cookie"], "sessionid=abc123") //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers["Set-Cookie"], "csrf=xyz789")      //nolint:staticcheck // testing deprecated field
		assert.Equal(t, "single-value", resp.Headers["X-Custom"])          //nolint:staticcheck // testing deprecated field

		// MultiHeaders (per-value slices)
		require.NotNil(t, resp.MultiHeaders)
		assert.Contains(t, resp.MultiHeaders, "Content-Type")
		assert.Equal(t, []string{"application/json"}, resp.MultiHeaders["Content-Type"].GetValues())

		setCookie := resp.MultiHeaders["Set-Cookie"]
		require.NotNil(t, setCookie)
		vals := setCookie.GetValues()
		require.Len(t, vals, 2)
		assert.Contains(t, vals, "sessionid=abc123; Path=/")
		assert.Contains(t, vals, "csrf=xyz789; Path=/")

		assert.Contains(t, resp.MultiHeaders, "X-Custom")
		assert.Equal(t, []string{"single-value"}, resp.MultiHeaders["X-Custom"].GetValues())
	})
}
