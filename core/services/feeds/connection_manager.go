package feeds

import (
	"crypto"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc/connectivity"

	"github.com/smartcontractkit/wsrpc"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	pb "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
)

type ConnectionsManager interface {
	Connect(opts ConnectOpts)
	Disconnect(id int64) error
	Close()
	GetClient(id int64) (pb.FeedsManagerClient, error)
	IsConnected(id int64) bool
}

// connectionsManager manages the rpc connections to Feeds Manager services
type connectionsManager struct {
	mu       sync.Mutex
	wgClosed sync.WaitGroup

	connections map[int64]*connection
	lggr        logger.Logger
}

type connection struct {
	stopCh services.StopChan

	connected bool
	client    pb.FeedsManagerClient
}

func newConnectionsManager(lggr logger.Logger) *connectionsManager {
	return &connectionsManager{
		mu:          sync.Mutex{},
		connections: map[int64]*connection{},
		lggr:        lggr,
	}
}

// ConnectOpts defines the required options to connect to an FMS server
type ConnectOpts struct {
	FeedsManagerID int64

	// URI is the URI of the feeds manager
	URI string

	CSASigner crypto.Signer

	// Pubkey defines the Feeds Manager Service's public key
	Pubkey []byte

	// Handlers defines the wsrpc Handlers
	Handlers pb.NodeServiceServer

	// OnConnect defines a callback for when the dial succeeds
	OnConnect func(pb.FeedsManagerClient)
}

// Connects to a feeds manager
//
// Connection to FMS is handled in a goroutine because the Dial will block
// until it can establish a connection. This is important during startup because
// we do not want to block other services from starting.
//
// Eventually when FMS does come back up, wsrpc will establish the connection
// without any interaction on behalf of the node operator.
func (mgr *connectionsManager) Connect(opts ConnectOpts) {
	conn := &connection{
		stopCh:    make(chan struct{}),
		connected: false,
	}

	mgr.wgClosed.Add(1)

	mgr.mu.Lock()
	mgr.connections[opts.FeedsManagerID] = conn
	mgr.mu.Unlock()

	go recovery.WrapRecover(mgr.lggr, func() {
		ctx, cancel := conn.stopCh.NewCtx()
		defer cancel()
		defer mgr.wgClosed.Done()

		mgr.lggr.Infow("Connecting to Feeds Manager...", "feedsManagerID", opts.FeedsManagerID)

		clientConn, err := wsrpc.DialWithContext(ctx, opts.URI,
			wsrpc.WithTransportSigner(opts.CSASigner, opts.Pubkey),
			wsrpc.WithBlock(),
			wsrpc.WithLogger(mgr.lggr),
		)
		if err != nil {
			// We only want to log if there was an error that did not occur
			// from a context cancel.
			if ctx.Err() == nil {
				mgr.lggr.Warnf("Error connecting to Feeds Manager server: %v", err)
			} else {
				mgr.lggr.Infof("Closing wsrpc websocket connection: %v", err)
			}

			return
		}
		defer func() {
			cerr := clientConn.Close()
			if cerr != nil {
				mgr.lggr.Warnf("Error closing wsrpc client connection: %v", cerr)
			}
		}()

		mgr.lggr.Infow("Connected to Feeds Manager", "feedsManagerID", opts.FeedsManagerID)

		// Initialize a new wsrpc client to make RPC calls
		mgr.mu.Lock()
		conn.connected = true
		conn.client = pb.NewFeedsManagerClient(clientConn)
		mgr.connections[opts.FeedsManagerID] = conn
		mgr.mu.Unlock()

		// Initialize RPC call handlers on the client connection
		pb.RegisterNodeServiceServer(clientConn, opts.Handlers)

		if opts.OnConnect != nil {
			opts.OnConnect(conn.client)
		}

		// Detect changes in connection status
		go func() {
			for {
				s := clientConn.GetState()

				clientConn.WaitForStateChange(ctx, s)

				s = clientConn.GetState()

				// Exit the goroutine if we shutdown the connection
				if s == connectivity.Shutdown {
					break
				}

				mgr.mu.Lock()
				conn.connected = s == connectivity.Ready
				mgr.mu.Unlock()
			}
		}()

		// Wait for close
		<-ctx.Done()
	})
}

// Disconnect closes a single connection
func (mgr *connectionsManager) Disconnect(id int64) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	conn, ok := mgr.connections[id]
	if !ok {
		return errors.New("feeds manager is not connected")
	}

	close(conn.stopCh)
	delete(mgr.connections, id)

	mgr.lggr.Infow("Disconnected Feeds Manager", "feedsManagerID", id)

	return nil
}

// Close closes all connections
func (mgr *connectionsManager) Close() {
	mgr.mu.Lock()
	for _, conn := range mgr.connections {
		close(conn.stopCh)
	}

	mgr.mu.Unlock()

	mgr.wgClosed.Wait()
}

// GetClient returns a single client by id
func (mgr *connectionsManager) GetClient(id int64) (pb.FeedsManagerClient, error) {
	mgr.mu.Lock()
	conn, ok := mgr.connections[id]
	mgr.mu.Unlock()
	if !ok || !conn.connected {
		return nil, errors.New("feeds manager is not connected")
	}

	return conn.client, nil
}

// IsConnected returns true if the connection to a feeds manager is active
func (mgr *connectionsManager) IsConnected(id int64) bool {
	mgr.mu.Lock()
	conn, ok := mgr.connections[id]
	mgr.mu.Unlock()
	if !ok {
		return false
	}

	return conn.connected
}
