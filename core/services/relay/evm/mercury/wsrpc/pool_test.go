package wsrpc

import (
	"context"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var _ Client = &mockClient{}

type mockClient struct {
	started   bool
	closed    bool
	rawClient pb.MercuryClient
}

func (c *mockClient) Transmit(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
	return
}
func (c *mockClient) LatestReport(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
	return
}
func (c *mockClient) Start(context.Context) error {
	c.started = true
	return nil
}
func (c *mockClient) Close() error {
	c.closed = true
	return nil
}
func (c *mockClient) Name() string                   { return "mock client" }
func (c *mockClient) Ready() error                   { return nil }
func (c *mockClient) HealthReport() map[string]error { return nil }
func (c *mockClient) ServerURL() string              { return "mock client url" }
func (c *mockClient) RawClient() pb.MercuryClient    { return c.rawClient }

func newMockClient(lggr logger.Logger) *mockClient {
	return &mockClient{}
}

func Test_Pool(t *testing.T) {
	lggr := logger.Sugared(logger.Test(t)).Named("PoolTestLogger")

	ctx := testutils.Context(t)

	t.Run("Checkout", func(t *testing.T) {
		p := newPool(lggr)
		p.cacheSet = &mockCacheSet{}

		t.Run("checks out one started client", func(t *testing.T) {
			clientPrivKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(rand.Int63()))
			serverPubKey := utils.NewHash().Bytes()
			serverURL := "example.com:443/ws"

			client := newMockClient(lggr)
			p.newClient = func(opts ClientOpts) Client {
				assert.Equal(t, clientPrivKey, opts.ClientPrivKey)
				assert.Equal(t, serverPubKey, opts.ServerPubKey)
				assert.Equal(t, serverURL, opts.ServerURL)
				return client
			}

			c, err := p.Checkout(ctx, clientPrivKey, serverPubKey, serverURL)
			require.NoError(t, err)

			assert.True(t, client.started)

			require.IsType(t, &clientCheckout{}, c)

			conn := c.(*clientCheckout).connection
			require.Equal(t, conn.Client, client)

			assert.Len(t, conn.checkouts, 1)
			assert.Same(t, lggr, conn.lggr)
			assert.Equal(t, clientPrivKey, conn.clientPrivKey)
			assert.Equal(t, serverPubKey, conn.serverPubKey)
			assert.Equal(t, serverURL, conn.serverURL)
			assert.Same(t, p, conn.pool)

			t.Run("checks in the clientCheckout when Close is called", func(t *testing.T) {
				err := c.Close()
				require.NoError(t, err)

				assert.Len(t, conn.checkouts, 0)
				require.IsType(t, nil, conn.Client)
				assert.Nil(t, conn.Client)
				assert.True(t, client.closed)
			})
		})

		t.Run("checks out multiple started clients and only closes if all of the clients for a given pk/server pair are checked back in", func(t *testing.T) {
			clientPrivKeys := []csakey.KeyV2{
				csakey.MustNewV2XXXTestingOnly(big.NewInt(rand.Int63())),
				csakey.MustNewV2XXXTestingOnly(big.NewInt(rand.Int63())),
			}
			serverPubKey := utils.NewHash().Bytes()
			serverURLs := []string{
				"example.com:443/ws",
				"example.invalid:8000/ws",
			}

			p.newClient = func(opts ClientOpts) Client {
				return newMockClient(opts.Logger)
			}

			// conn 1
			c1 := mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])
			c2 := mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])
			c3 := mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])
			assert.Len(t, p.connections, 1)
			assert.Len(t, p.connections[serverURLs[0]], 1)
			assert.Len(t, p.connections[serverURLs[1]], 0)

			// conn 2
			c4 := mustCheckout(t, p, clientPrivKeys[1], serverPubKey, serverURLs[0])
			assert.Len(t, p.connections, 1)
			assert.Len(t, p.connections[serverURLs[0]], 2)
			assert.Len(t, p.connections[serverURLs[1]], 0)

			// conn 3
			c5 := mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[1])
			c6 := mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[1])
			assert.Len(t, p.connections, 2)
			assert.Len(t, p.connections[serverURLs[0]], 2)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			conn1 := c1.(*clientCheckout).connection
			assert.Same(t, conn1, c2.(*clientCheckout).connection)
			assert.Same(t, conn1, c3.(*clientCheckout).connection)
			assert.Len(t, conn1.checkouts, 3)
			assert.True(t, conn1.Client.(*mockClient).started)

			conn2 := c4.(*clientCheckout).connection
			assert.NotEqual(t, conn1, conn2)
			assert.Len(t, conn2.checkouts, 1)
			assert.True(t, conn2.Client.(*mockClient).started)

			conn3 := c5.(*clientCheckout).connection
			assert.NotEqual(t, conn1, conn3)
			assert.NotEqual(t, conn2, conn3)
			assert.Same(t, conn3, c6.(*clientCheckout).connection)
			assert.Len(t, conn3.checkouts, 2)
			assert.True(t, conn3.Client.(*mockClient).started)

			require.NoError(t, c1.Close())
			assert.Len(t, conn1.checkouts, 2)
			assert.NotNil(t, conn1.Client)
			assert.Len(t, p.connections, 2)
			assert.Len(t, p.connections[serverURLs[0]], 2)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			require.NoError(t, c2.Close())
			assert.Len(t, conn1.checkouts, 1)
			assert.NotNil(t, conn1.Client)
			assert.Len(t, p.connections, 2)
			assert.Len(t, p.connections[serverURLs[0]], 2)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			require.NoError(t, c3.Close())
			assert.Len(t, conn1.checkouts, 0)
			assert.Nil(t, conn1.Client)
			assert.Len(t, p.connections, 2)
			assert.Len(t, p.connections[serverURLs[0]], 1)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			c7 := mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])
			// Not the same one, since previously all checkouts were checked in, the original connection was deleted from the map and a new one created
			assert.NotSame(t, conn1, c7.(*clientCheckout).connection)
			assert.Len(t, conn1.checkouts, 0) // actually, conn1 has already been removed from the map and will be garbage collected
			conn4 := c7.(*clientCheckout).connection
			assert.Len(t, conn4.checkouts, 1)
			assert.NotNil(t, conn4.Client)
			assert.Len(t, p.connections, 2)
			assert.Len(t, p.connections[serverURLs[0]], 2)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			require.NoError(t, c7.Close())
			assert.Len(t, p.connections, 2)
			assert.Len(t, p.connections[serverURLs[0]], 1)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			require.NoError(t, c4.Close())
			assert.Len(t, p.connections, 1)
			assert.Len(t, p.connections[serverURLs[0]], 0)
			assert.Len(t, p.connections[serverURLs[1]], 1)

			require.NoError(t, c5.Close())
			require.NoError(t, c6.Close())
			assert.Len(t, p.connections, 0)

			require.NoError(t, p.Close())
		})
	})

	p := newPool(lggr)
	p.cacheSet = &mockCacheSet{}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "PoolTestLogger", p.Name())
	})
	t.Run("Start", func(t *testing.T) {
		require.NoError(t, p.Start(ctx))
		assert.Nil(t, p.Ready())
		assert.Nil(t, p.HealthReport()["PoolTestLogger"])
	})
	t.Run("Close force closes all connections", func(t *testing.T) {
		clientPrivKeys := []csakey.KeyV2{
			csakey.MustNewV2XXXTestingOnly(big.NewInt(rand.Int63())),
			csakey.MustNewV2XXXTestingOnly(big.NewInt(rand.Int63())),
		}
		serverPubKey := utils.NewHash().Bytes()
		serverURLs := []string{
			"example.com:443/ws",
			"example.invalid:8000/ws",
		}

		var clients []*mockClient
		p.newClient = func(opts ClientOpts) Client {
			c := newMockClient(opts.Logger)
			clients = append(clients, c)
			return c
		}

		// conn 1
		mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])
		mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])
		mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[0])

		// conn 2
		mustCheckout(t, p, clientPrivKeys[1], serverPubKey, serverURLs[0])

		// conn 3
		mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[1])
		mustCheckout(t, p, clientPrivKeys[0], serverPubKey, serverURLs[1])

		for _, c := range clients {
			assert.True(t, c.started)
			assert.False(t, c.closed)
		}

		require.NoError(t, p.Close())
		assert.EqualError(t, p.Ready(), "pool is closed")
		assert.EqualError(t, p.HealthReport()["PoolTestLogger"], "pool is closed")

		for _, c := range clients {
			assert.True(t, c.closed)
		}
	})
}

func mustCheckout(t *testing.T, p *pool, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) Client {
	c, err := p.Checkout(testutils.Context(t), clientPrivKey, serverPubKey, serverURL)
	require.NoError(t, err)
	return c
}
