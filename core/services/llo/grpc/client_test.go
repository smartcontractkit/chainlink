package grpc

import (
	"crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-data-streams/rpc"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_Client(t *testing.T) {
	t.Run("Transmit errors if not started", func(t *testing.T) {
		c := NewClient(ClientOpts{
			Logger:        logger.TestLogger(t),
			ClientPrivKey: ed25519.PrivateKey{},
			ServerPubKey:  ed25519.PublicKey{},
			ServerURL:     "example.com",
		})

		resp, err := c.Transmit(tests.Context(t), &rpc.TransmitRequest{})
		assert.Nil(t, resp)
		require.EqualError(t, err, "service is Unstarted, not started")
	})
}
