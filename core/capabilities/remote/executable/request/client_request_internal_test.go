package request

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientRequest_Expired_aggregationGrace(t *testing.T) {
	t.Parallel()

	requestTimeout := 100 * time.Millisecond
	t.Run("not expired before requestTimeout plus grace", func(t *testing.T) {
		t.Parallel()
		c := &ClientRequest{
			createdAt:      time.Now().Add(-requestTimeout - time.Millisecond), // less than defaultResponseAggregationGrace
			requestTimeout: requestTimeout,
		}
		require.False(t, c.Expired())
	})

	t.Run("expired after requestTimeout plus grace", func(t *testing.T) {
		t.Parallel()
		c := &ClientRequest{
			createdAt:      time.Now().Add(-requestTimeout - defaultResponseAggregationGrace - time.Millisecond),
			requestTimeout: requestTimeout,
		}
		require.True(t, c.Expired())
	})
}
