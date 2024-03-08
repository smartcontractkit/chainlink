package wsrpc

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/wsrpc/credentials"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/cache"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//var _ Client = &clientCheckout{}

type clientCheckout struct {
	connection // inherit all methods from client, with override on Start/Close
}

func (cco *clientCheckout) Start(_ context.Context) error {
	return nil
}

func (cco *clientCheckout) Close() error {
	cco.connection.checkin(cco)
	return nil
}

type baseConnection struct {
	lggr          logger.Logger
	clientPrivKey csakey.KeyV2
	serverPubKey  []byte
	serverURL     string

	checkouts []*clientCheckout // reference count, if this goes to zero the connection should be closed and *client nilified

	mu sync.Mutex
}

type wsrpcConnection struct {
	*baseConnection
	pool   *wsrpcPool
	Client *WsrpcClient
}

type grpcConnection struct {
	*baseConnection
	pool      *grpcPool
	Client    *GrpcClient
	tlsConfig *tlsConfig
}

// connection is a convenience interface to support both wsrpc and grpc connections
type connection interface {
	checkout(ctx context.Context) (*clientCheckout, error)
	checkin(checkinCco *clientCheckout)
	ensureStartedClient(ctx context.Context) error
	forceCloseAll() error
}

func (conn *wsrpcConnection) checkout(ctx context.Context) (cco *clientCheckout, err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if err = conn.ensureStartedClient(ctx); err != nil {
		return nil, err
	}
	cco = &clientCheckout{conn}
	conn.checkouts = append(conn.checkouts, cco)
	return cco, nil
}

func (conn *grpcConnection) checkout(ctx context.Context) (cco *clientCheckout, err error) {
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
func (conn *wsrpcConnection) ensureStartedClient(ctx context.Context) error {
	if len(conn.checkouts) == 0 {
		conn.Client = conn.pool.newClient(conn.lggr, conn.clientPrivKey, conn.serverPubKey, conn.serverURL, conn.pool.cacheSet)
		return conn.Client.Start(ctx)
	}
	return nil
}

// not thread-safe, access must be serialized
func (conn *grpcConnection) ensureStartedClient(ctx context.Context) error {
	if len(conn.checkouts) == 0 {
		conn.Client = conn.pool.newClient(conn.lggr, conn.clientPrivKey, conn.serverPubKey, conn.serverURL, conn.pool.cacheSet, conn.tlsConfig.CertFile)
		return conn.Client.Start(ctx)
	}
	return nil
}

func (conn *wsrpcConnection) checkin(checkinCco *clientCheckout) {
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

func (conn *grpcConnection) checkin(checkinCco *clientCheckout) {
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

func (conn *wsrpcConnection) forceCloseAll() (err error) {
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

func (conn *grpcConnection) forceCloseAll() (err error) {
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

// TlsConfig holds the TLS configuration for gRPC clients
type tlsConfig struct {
	CertFile *string
	Enabled  bool
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
	connections map[string]map[credentials.StaticSizedPublicKey]connection
	mu          sync.RWMutex

	closed bool
}

type grpcPool struct {
	*pool

	// embedding newClient makes testing/mocking easier
	newClient func(lggr logger.Logger, privKey csakey.KeyV2, serverPubKey []byte, serverURL string, cacheSet cache.GrpcCacheSet, tlsCertFile *string) *GrpcClient
	cacheSet  cache.GrpcCacheSet
}

type wsrpcPool struct {
	*pool

	// embedding newClient makes testing/mocking easier
	newClient func(lggr logger.Logger, privKey csakey.KeyV2, serverPubKey []byte, serverURL string, cacheSet cache.WsrpcCacheSet) *WsrpcClient
	cacheSet  cache.WsrpcCacheSet
}

func NewWsrpcPool(lggr logger.Logger, cacheCfg cache.Config) *wsrpcPool {
	return &wsrpcPool{
		&pool{
			lggr:        lggr.Named("Mercury.WSRPCPool"),
			connections: make(map[string]map[credentials.StaticSizedPublicKey]connection),
		},
		NewWsrpcClient,
		cache.NewWsrpcCacheSet(lggr, cacheCfg),
	}
}

func NewGrpcPool(lggr logger.Logger, cacheCfg cache.Config, tlsCfg tlsConfig) *grpcPool {
	return &grpcPool{
		&pool{
			lggr:        lggr.Named("Mercury.GRPCPool"),
			connections: make(map[string]map[credentials.StaticSizedPublicKey]connection),
		},
		NewGrpcClient,
		cache.NewGrpcCacheSet(lggr, cacheCfg),
	}
}

func (p *wsrpcPool) Checkout(ctx context.Context, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) (client *clientCheckout, err error) {
	clientPubKey := clientPrivKey.StaticSizedPublicKey()

	p.mu.Lock()

	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("pool is closed")
	}

	server, exists := p.connections[serverURL]
	if !exists {
		server = make(map[credentials.StaticSizedPublicKey]connection)
		p.connections[serverURL] = server
	}
	conn, exists := server[clientPubKey]
	if !exists {
		conn = p.newWsrpcConnection(p.lggr, clientPrivKey, serverPubKey, serverURL)
		server[clientPubKey] = conn
	}
	p.mu.Unlock()

	// checkout outside of pool lock since it might take non-trivial time
	// the clientCheckout will be checked in again when its Close() method is called
	// this also should avoid deadlocks between conn.mu and pool.mu
	return conn.checkout(ctx)
}

func (p *grpcPool) Checkout(ctx context.Context, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) (client *clientCheckout, err error) {
	clientPubKey := clientPrivKey.StaticSizedPublicKey()

	p.mu.Lock()

	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("pool is closed")
	}

	server, exists := p.connections[serverURL]
	if !exists {
		server = make(map[credentials.StaticSizedPublicKey]connection)
		p.connections[serverURL] = server
	}
	conn, exists := server[clientPubKey]
	if !exists {
		conn = p.newGrpcConnection(p.lggr, clientPrivKey, serverPubKey, serverURL)
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

func (p *wsrpcPool) newWsrpcConnection(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) *wsrpcConnection {
	return &wsrpcConnection{
		baseConnection: &baseConnection{
			lggr:          lggr,
			clientPrivKey: clientPrivKey,
			serverPubKey:  serverPubKey,
			serverURL:     serverURL,
		},
		pool: p,
	}
}

func (p *grpcPool) newGrpcConnection(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) *grpcConnection {
	return &grpcConnection{
		baseConnection: &baseConnection{
			lggr:          lggr,
			clientPrivKey: clientPrivKey,
			serverPubKey:  serverPubKey,
			serverURL:     serverURL,
		},
		pool: p,
	}
}

func (p *wsrpcPool) Start(ctx context.Context) error {
	return p.cacheSet.Start(ctx)
}

func (p *grpcPool) Start(ctx context.Context) error {
	return p.cacheSet.Start(ctx)
}

func (p *wsrpcPool) Close() (merr error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	for _, clientPubKeys := range p.connections {
		for _, conn := range clientPubKeys {
			merr = errors.Join(merr, conn.forceCloseAll())
		}
	}
	merr = errors.Join(merr, p.cacheSet.Close())
	return
}

func (p *grpcPool) Close() (merr error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	for _, clientPubKeys := range p.connections {
		for _, conn := range clientPubKeys {
			merr = errors.Join(merr, conn.forceCloseAll())
		}
	}
	merr = errors.Join(merr, p.cacheSet.Close())
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

func (p *wsrpcPool) HealthReport() map[string]error {
	hp := map[string]error{p.Name(): p.Ready()}
	maps.Copy(hp, p.cacheSet.HealthReport())
	return hp
}

func (p *grpcPool) HealthReport() map[string]error {
	hp := map[string]error{p.Name(): p.Ready()}
	maps.Copy(hp, p.cacheSet.HealthReport())
	return hp
}
