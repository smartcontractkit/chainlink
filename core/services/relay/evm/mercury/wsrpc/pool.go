package wsrpc

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ Client = &clientCheckout{}

type clientCheckout struct {
	*connection // inherit all methods from client, with override on Start/Close
}

func (cco *clientCheckout) Start(_ context.Context) error {
	return nil
}

func (cco *clientCheckout) Close() error {
	cco.connection.checkin(cco)
	return nil
}

type connection struct {
	// Client will be nil when checkouts is empty, if len(checkouts) > 0 then it is expected to be a non-nil, started client
	Client

	lggr          logger.Logger
	clientPrivKey csakey.KeyV2
	serverPubKey  []byte
	serverURL     string

	pool *pool

	checkouts []*clientCheckout // reference count, if this goes to zero the connection should be closed and *client nilified

	mu sync.Mutex
}

func (conn *connection) checkout(ctx context.Context) (cco *clientCheckout, err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if err = conn.ensureStartedClient(ctx); err != nil {
		return nil, err
	}
	cco = &clientCheckout{conn}
	conn.checkouts = append(conn.checkouts, cco)
	return cco, nil
}

// not thread-safe, access must be serialized
func (conn *connection) ensureStartedClient(ctx context.Context) error {
	if len(conn.checkouts) == 0 {
		conn.Client = conn.pool.newClient(conn.lggr, conn.clientPrivKey, conn.serverPubKey, conn.serverURL)
		return conn.Client.Start(ctx)
	}
	return nil
}

func (conn *connection) checkin(checkinCco *clientCheckout) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	var removed bool
	for i, cco := range conn.checkouts {
		if cco == checkinCco {
			conn.checkouts = utils.DeleteUnstable(conn.checkouts, i)
			removed = true
			break
		}
	}
	if !removed {
		panic("tried to check in client that was never checked out")
	}
	if len(conn.checkouts) == 0 {
		if err := conn.Client.Close(); err != nil {
			// programming error if we hit this
			panic(err)
		}
		conn.Client = nil
		conn.pool.remove(conn.serverURL, conn.clientPrivKey.StaticSizedPublicKey())
	}
}

func (conn *connection) forceCloseAll() (err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.Client != nil {
		err = conn.Client.Close()
		if errors.Is(err, utils.ErrAlreadyStopped) {
			// ignore error if it has already been stopped; no problem
			err = nil
		}
		conn.Client = nil
		conn.checkouts = nil
	}
	return
}

type Pool interface {
	services.ServiceCtx
	// Checkout gets a wsrpc.Client for the given arguments
	// The same underlying client can be checked out multiple times, the pool
	// handles lifecycle management. The consumer can treat it as if it were
	// its own unique client.
	Checkout(ctx context.Context, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) (client Client, err error)
}

// WSRPC allows only one connection per client key per server
type pool struct {
	lggr logger.Logger
	// server url => client public key => connection
	connections map[string]map[credentials.StaticSizedPublicKey]*connection

	// embedding newClient makes testing/mocking easier
	newClient func(lggr logger.Logger, privKey csakey.KeyV2, serverPubKey []byte, serverURL string) Client

	mu sync.RWMutex

	closed bool
}

func NewPool(lggr logger.Logger) Pool {
	p := newPool(lggr)
	p.newClient = NewClient
	return p
}

func newPool(lggr logger.Logger) *pool {
	return &pool{
		lggr:        lggr,
		connections: make(map[string]map[credentials.StaticSizedPublicKey]*connection),
	}
}

func (p *pool) Checkout(ctx context.Context, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) (client Client, err error) {
	clientPubKey := clientPrivKey.StaticSizedPublicKey()

	p.mu.Lock()

	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("pool is closed")
	}

	server, exists := p.connections[serverURL]
	if !exists {
		server = make(map[credentials.StaticSizedPublicKey]*connection)
		p.connections[serverURL] = server
	}
	conn, exists := server[clientPubKey]
	if !exists {
		conn = p.newConnection(p.lggr, clientPrivKey, serverPubKey, serverURL)
		server[clientPubKey] = conn
	}
	p.mu.Unlock()

	// checkout outside of pool lock since it might take non-trivial time
	// the clientCheckout will be checked in again when its Close() method is called
	// this also should avoid deadlocks between conn.mu and pool.mu
	return conn.checkout(ctx)
}

// remove performs garbage collection on the connections map after connections are no longer used
func (p *pool) remove(serverURL string, clientPubKey credentials.StaticSizedPublicKey) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.connections[serverURL], clientPubKey)
	if len(p.connections[serverURL]) == 0 {
		delete(p.connections, serverURL)
	}

}

func (p *pool) newConnection(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) *connection {
	return &connection{
		lggr:          lggr,
		clientPrivKey: clientPrivKey,
		serverPubKey:  serverPubKey,
		serverURL:     serverURL,
		pool:          p,
	}
}

func (p *pool) Start(_ context.Context) error { return nil }

func (p *pool) Close() (merr error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	for _, clientPubKeys := range p.connections {
		for _, conn := range clientPubKeys {
			merr = errors.Join(merr, conn.forceCloseAll())
		}
	}
	return
}

func (p *pool) Name() string {
	return p.lggr.Name()
}

func (p *pool) Ready() error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return errors.New("pool is closed")
	}
	return nil
}

func (p *pool) HealthReport() map[string]error {
	return map[string]error{
		p.Name(): p.Ready(),
	}
}
