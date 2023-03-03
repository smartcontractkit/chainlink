package wsrpc

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/internal/message"
	"github.com/smartcontractkit/wsrpc/internal/transport"
	"github.com/smartcontractkit/wsrpc/internal/wsrpcsync"
	"github.com/smartcontractkit/wsrpc/peer"
)

var ErrNotConnected = errors.New("client not connected")

// Server is a wsrpc server to both perform and serve RPC requests.
type Server struct {
	mu sync.RWMutex

	wssrv *http.Server

	opts serverOptions

	// Manages the open client connections
	connMgr *connectionsManager

	// Parameters for upgrading a websocket connection
	upgrader websocket.Upgrader
	// The RPC service definition
	service *serviceInfo

	// Contains all pending method call ids and the channel to respond to when
	// a result is received
	methodCalls map[string]chan<- *message.Response

	// Signals a quit event when the server wants to quit
	quit *wsrpcsync.Event
	// Signals a done event once the server has finished shutting down
	done *wsrpcsync.Event

	serveWG sync.WaitGroup
}

// NewServer initializes a new wsrpc server.
func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	s := &Server{
		opts: opts,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  opts.readBufferSize,
			WriteBufferSize: opts.writeBufferSize,
		},
		connMgr:     newConnectionsManager(),
		methodCalls: map[string]chan<- *message.Response{},
		quit:        wsrpcsync.NewEvent(),
		done:        wsrpcsync.NewEvent(),
		serveWG:     sync.WaitGroup{},
	}

	return s
}

// Serve accepts incoming connections on the listener lis, creating a new
// ServerTransport and service goroutine for each.
func (s *Server) Serve(lis net.Listener) {
	// Serve Healthcheck over HTTP
	if s.opts.healthcheckAddr != "" {
		hclis, err := net.Listen("tcp", s.opts.healthcheckAddr)
		if err != nil {
			panic(err)
		}
		hchandler := http.NewServeMux()
		hchandler.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		hcsrv := &http.Server{
			Handler: hchandler,
		}

		//nolint:errcheck
		go hcsrv.Serve(hclis)
		defer hcsrv.Close()
	}

	// Serve websockets over HTTPS
	wshandler := http.NewServeMux()
	wshandler.HandleFunc("/", s.wshandler)
	wssrv := &http.Server{
		TLSConfig: s.opts.creds.Config,
		Handler:   wshandler,
	}

	//nolint:errcheck
	go wssrv.ServeTLS(lis, "", "")
	defer wssrv.Close()

	s.wssrv = wssrv

	<-s.done.Done()
}

// wshandler upgrades the HTTP connection to a websocket connection and
// registers the connection's pub key for the client.
func (s *Server) wshandler(w http.ResponseWriter, r *http.Request) {
	// Do not establish a new connection if quit has already been fired
	if s.quit.HasFired() {
		return
	}

	pubKey, err := s.ensureSingleClientConnection(r.TLS.PeerCertificates[0])
	if err != nil {
		return
	}

	// Upgrade the websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// A signal channel to close down running go routines (i.e receiving
	// messages) and then ensures the handler to returns
	done := make(chan struct{})

	config := &transport.ServerConfig{}
	onClose := func() {
		// There is no connection manager when we are shutting down, so
		// we can ignore removing the connection.
		if s.connMgr != nil {
			s.connMgr.mu.Lock()
			s.connMgr.removeConnection(pubKey)
			s.connMgr.mu.Unlock()
		}
		s.serveWG.Done()
		close(done)
	}

	// Initialize the transport
	tr, err := transport.NewServerTransport(conn, config, onClose)
	if err != nil {
		return
	}

	// Register the transport against the public key
	s.connMgr.registerConnection(pubKey, tr)

	s.serveWG.Add(1)

	// Start the reader handler
	go s.handleRead(pubKey, done)

	select {
	case <-done:
		log.Println("[wsrpc] Connection dropped")
	case <-s.quit.Done():
		log.Println("[wsrpc] Connection closed due to shutdown")
	}
}

// sendMsg writes the message to the connection which matches the public key.
func (s *Server) sendMsg(pub [32]byte, msg []byte) error {
	// Find the transport matching the public key
	tr, err := s.connMgr.getTransport(pub)
	if err != nil {
		return err
	}

	return tr.Write(msg)
}

// handleRead listens to the transport read channel and passes the message to the
// readFn handler.
func (s *Server) handleRead(pubKey credentials.StaticSizedPublicKey, done <-chan struct{}) {
	tr, err := s.connMgr.getTransport(pubKey)
	if err != nil {
		return
	}

	for {
		select {
		case in := <-tr.Read():
			// Unmarshal the message
			msg := &message.Message{}
			if err := UnmarshalProtoMessage(in, msg); err != nil {
				continue
			}

			// Handle the message request or response
			switch ex := msg.Exchange.(type) {
			case *message.Message_Request:
				go s.handleMessageRequest(pubKey, ex.Request)
			case *message.Message_Response:
				go s.handleMessageResponse(ex.Response)
			default:
				log.Println("Invalid message type")
			}
		case <-done:
			return
		}
	}
}

// handleMessageRequest looks up the method matching the method name and calls
// the handler. The connection client's public key is injected into the context,
// so the handler is able to identify the caller.
func (s *Server) handleMessageRequest(pubKey credentials.StaticSizedPublicKey, r *message.Request) {
	methodName := r.GetMethod()
	if md, ok := s.service.methods[methodName]; ok {
		// Create a decoder function to unmarshal the message
		dec := func(v interface{}) error {
			return UnmarshalProtoMessage(r.GetPayload(), v)
		}

		// Inject the peer's public key into the context so the handler can use it
		ctx := peer.NewContext(context.Background(), &peer.Peer{PublicKey: pubKey})
		v, herr := md.Handler(s.service.serviceImpl, ctx, dec)

		msg, err := message.NewResponse(r.GetCallId(), v, herr)
		if err != nil {
			return
		}

		replyMsg, err := MarshalProtoMessage(msg)
		if err != nil {
			return
		}

		if err := s.sendMsg(pubKey, replyMsg); err != nil {
			log.Printf("error sending message: %s", err)
		}
	}
}

// handleMessageResponse finds the call which matches the method call id of the
// response and sends the payload to the call channel.
func (s *Server) handleMessageResponse(r *message.Response) {
	s.mu.Lock()
	defer s.mu.Unlock()

	callID := r.GetCallId()
	if call, ok := s.methodCalls[callID]; ok {
		call <- r

		s.removeMethodCall(callID) // Delete the call now that we have completed the request/response cycle
	}
}

// RegisterService registers a service and its implementation to the wsrpc
// server. This must be called before invoking Serve.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	s.register(sd, ss)
}

func (s *Server) register(sd *ServiceDesc, ss interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	info := &serviceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*MethodDesc),
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	s.service = info
}

// Invoke sends the RPC request on the connection which is connected with the
// public key and returns after response is received.
func (s *Server) Invoke(ctx context.Context, method string, args interface{}, reply interface{}) error {
	callID := uuid.NewString()
	msg, err := message.NewRequest(callID, method, args)
	if err != nil {
		return err
	}

	req, err := MarshalProtoMessage(msg)
	if err != nil {
		return err
	}

	s.mu.Lock()
	wait := s.registerMethodCall(callID)
	s.mu.Unlock()

	// Extract the public key from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return errors.New("could not extract public key")
	}
	pubKey := p.PublicKey

	if err = s.sendMsg(pubKey, req); err != nil {
		return err
	}

	// Wait for the response
	select {
	case msg := <-wait:
		// Handle error
		if msg.Error != "" {
			return errors.New(msg.Error)
		}

		// Unmarshal the payload into the reply
		if err := UnmarshalProtoMessage(msg.GetPayload(), reply); err != nil {
			return err
		}
	case <-ctx.Done():
		// Remove the call since we have timeout
		s.mu.Lock()
		s.removeMethodCall(callID)
		s.mu.Unlock()

		return fmt.Errorf("call timeout: %w", ctx.Err())
	}

	return nil
}

// UpdatePublicKeys updates the list of allowable public keys in the TLS config
// and drops the connections which match the deleted keys.
func (s *Server) UpdatePublicKeys(pubKeys []ed25519.PublicKey) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.opts.creds.PublicKeys.Replace(pubKeys)
	s.removeConnectionsToDeletedKeys(pubKeys)
}

// GetConnectionNotifyChan gets the connection notification channel.
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func (s *Server) GetConnectionNotifyChan() <-chan struct{} {
	return s.connMgr.getNotifyChan()
}

// GetConnectedPeerPublicKeys gets the public keys for all peers which are
// connected.
//
// Notice: This API is EXPERIMENTAL and may be changed or removed in a
// later release.
func (s *Server) GetConnectedPeerPublicKeys() []credentials.StaticSizedPublicKey {
	return s.connMgr.getConnectionPublicKeys()
}

// Stop stops the wsRPC server. It immediately closes all open
// connections and listeners.
func (s *Server) Stop() {
	s.quit.Fire()
	defer func() {
		s.done.Fire()
	}()

	s.mu.Lock()
	connMgr := s.connMgr
	s.connMgr = nil
	s.mu.Unlock()

	connMgr.close()

	// Wait for all the connections to close
	s.serveWG.Wait()
}

// When the list of allowable certs are updated, we need to refresh the existing
// connections as well and shutdown any client connections no longer allowed.
func (s *Server) removeConnectionsToDeletedKeys(pubKeys []ed25519.PublicKey) {
	pubKeysMap := make(map[credentials.StaticSizedPublicKey]bool)

	for _, pk := range pubKeys {
		pubKey, err := credentials.ToStaticallySizedPublicKey(pk)
		if err != nil {
			log.Print("[Server] error reading keys while removing connections: ", err)
		} else {
			pubKeysMap[pubKey] = true
		}
	}

	// Removing connections requires a lock on the connections manager.
	s.connMgr.mu.Lock()
	defer s.connMgr.mu.Unlock()

	for k, conn := range s.connMgr.conns {
		if _, ok := pubKeysMap[k]; !ok {
			conn.Close()
			s.connMgr.removeConnection(k)
		}
	}
}

// Ensure there is only a single connection per public key by checking the
// certificate's public key against the list of registered connections.
func (s *Server) ensureSingleClientConnection(cert *x509.Certificate) ([ed25519.PublicKeySize]byte, error) {
	pubKey, err := credentials.PubKeyFromCert(cert)
	if err != nil {
		return pubKey, errors.New("could not extracting public key from certificate")
	}

	_, err = s.connMgr.getTransport(pubKey)
	if err == nil {
		return pubKey, errors.New("only one connection allowed per client")
	}

	return pubKey, nil
}

// registerMethodCall registers a method call to the method call map.
//
// This requires a lock on cc.mu.
func (s *Server) registerMethodCall(id string) <-chan *message.Response {
	wait := make(chan *message.Response)
	s.methodCalls[id] = wait

	return wait
}

// removeMethodCall deregisters a method call to the method call map.
//
// This requires a lock on cc.mu.
func (s *Server) removeMethodCall(id string) {
	delete(s.methodCalls, id)
}

// connectionsManager manages the active clients connections.
type connectionsManager struct {
	mu sync.Mutex
	// Holds a list of the open connections mapped to a buffered channel of
	// outbound messages.
	conns map[credentials.StaticSizedPublicKey]transport.ServerTransport
	// Notifies receivers on this channel when the list of connections change
	notifyChan chan struct{}
}

func newConnectionsManager() *connectionsManager {
	return &connectionsManager{
		conns: map[credentials.StaticSizedPublicKey]transport.ServerTransport{},
	}
}

// getTransport fetches the transport which matches the public key.
func (cm *connectionsManager) getTransport(key credentials.StaticSizedPublicKey) (transport.ServerTransport, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	tr, ok := cm.conns[key]
	if !ok {
		return nil, ErrNotConnected
	}

	return tr, nil
}

// registerConnection registers a new transport mapped to a public key.
func (cm *connectionsManager) registerConnection(key credentials.StaticSizedPublicKey, value transport.ServerTransport) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.conns[key] = value

	if cm.notifyChan != nil {
		// There are other goroutines waiting on this channel.
		close(cm.notifyChan)
		cm.notifyChan = nil
	}
}

// removeConnection removes a transport from the registered list.
//
// This requires a lock on cm.mu.
func (cm *connectionsManager) removeConnection(key credentials.StaticSizedPublicKey) {
	delete(cm.conns, key)

	if cm.notifyChan != nil {
		// There are other goroutines waiting on this channel.
		close(cm.notifyChan)
		cm.notifyChan = nil
	}
}

// getConnectionPublicKeys gets the public keys of the active connections.
func (cm *connectionsManager) getConnectionPublicKeys() []credentials.StaticSizedPublicKey {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	keys := []credentials.StaticSizedPublicKey{}
	for k := range cm.conns {
		keys = append(keys, k)
	}

	return keys
}

// getNotifyChan returns a channel used to notify the watcher when the list
// of connections changes.
func (cm *connectionsManager) getNotifyChan() <-chan struct{} {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.notifyChan == nil {
		cm.notifyChan = make(chan struct{})
	}

	return cm.notifyChan
}

// close closes all registered connections.
func (cm *connectionsManager) close() {
	for _, conn := range cm.conns {
		conn.Close()
	}
}
