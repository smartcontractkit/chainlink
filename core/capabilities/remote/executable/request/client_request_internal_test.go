package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
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

func Test_getPeerIDToTransmissionDelays(t *testing.T) {
	t.Parallel()

	peer1 := [32]byte(fmt.Appendf(nil, "%-32s", "one"))
	peer2 := [32]byte(fmt.Appendf(nil, "%-32s", "two"))
	peer3 := [32]byte(fmt.Appendf(nil, "%-32s", "three"))
	peer4 := [32]byte(fmt.Appendf(nil, "%-32s", "four"))

	ids := []p2ptypes.PeerID{peer1, peer2, peer3, peer4}

	peerIDToDelay := getPeerIDToTransmissionDelays(ids)

	assert.Len(t, peerIDToDelay, len(ids))
	for _, id := range ids {
		assert.Equal(t, time.Duration(0), peerIDToDelay[id])
	}
}
