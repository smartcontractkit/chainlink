package wsrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	mocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

func Test_Client_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	req := &pb.TransmitRequest{}

	t.Run("sends on reset channel after MaxConsecutiveTransmitFailures timed out transmits", func(t *testing.T) {
		calls := 0
		transmitErr := context.DeadlineExceeded
		wsrpcClient := &mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
				calls++
				return nil, transmitErr
			},
		}
		conn := &mocks.MockConn{
			Ready: true,
		}
		c := newClient(lggr, csakey.KeyV2{}, nil, "")
		c.conn = conn
		c.client = wsrpcClient
		require.NoError(t, c.StartOnce("Mock WSRPC Client", func() error { return nil }))
		for i := 1; i < MaxConsecutiveTransmitFailures; i++ {
			_, err := c.Transmit(ctx, req)
			require.EqualError(t, err, "context deadline exceeded")
		}
		assert.Equal(t, 4, calls)
		select {
		case <-c.chResetTransport:
			t.Fatal("unexpected send on chResetTransport")
		default:
		}
		_, err := c.Transmit(ctx, req)
		require.EqualError(t, err, "context deadline exceeded")
		assert.Equal(t, 5, calls)
		select {
		case <-c.chResetTransport:
		default:
			t.Fatal("expected send on chResetTransport")
		}

		t.Run("successful transmit resets the counter", func(t *testing.T) {
			transmitErr = nil
			// working transmit to reset counter
			_, err = c.Transmit(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, 6, calls)
			assert.Equal(t, 0, int(c.consecutiveTimeoutCnt.Load()))
		})

		t.Run("doesn't block in case channel is full", func(t *testing.T) {
			transmitErr = context.DeadlineExceeded
			c.chResetTransport = nil // simulate full channel
			for i := 0; i < MaxConsecutiveTransmitFailures; i++ {
				_, err := c.Transmit(ctx, req)
				require.EqualError(t, err, "context deadline exceeded")
			}
		})
	})
}
