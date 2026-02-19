package grpcsourcemock

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"

	sourcesv1 "github.com/smartcontractkit/chainlink-protos/workflows/go/sources"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/privateregistry"
)

const (
	// DefaultSourcePort is the default port for the WorkflowMetadataSourceService
	DefaultSourcePort = 8544
	// DefaultPrivateRegistryPort is the default port for the private registry API
	// Uses 8547 to avoid conflicts with anvil chains (8545 for chain 1337, 8546 for chain 2337)
	DefaultPrivateRegistryPort = 8547
)

// NodeAuthProvider is the interface for validating node public keys
type NodeAuthProvider interface {
	IsNodePubKeyTrusted(ctx context.Context, publicKey ed25519.PublicKey) (bool, error)
}

// ServerConfig contains configuration for the mock server
type ServerConfig struct {
	// SourcePort is the port for the WorkflowMetadataSourceService (default: 8544)
	SourcePort int
	// PrivateRegistryPort is the port for the private registry API (default: 8545)
	PrivateRegistryPort int
	// AuthProvider is the provider for validating node public keys
	// If nil, all requests are allowed (no auth)
	AuthProvider NodeAuthProvider
}

// Server is the mock gRPC workflow source server
type Server struct {
	config                 ServerConfig
	store                  *WorkflowStore
	sourceServer           *grpc.Server
	privateRegistryServer  *grpc.Server
	privateRegistryService *PrivateRegistryService

	sourceListener          net.Listener
	privateRegistryListener net.Listener

	mu      sync.Mutex
	started bool
}

// NewServer creates a new mock gRPC workflow source server
func NewServer(config ServerConfig) *Server {
	if config.SourcePort == 0 {
		config.SourcePort = DefaultSourcePort
	}
	if config.PrivateRegistryPort == 0 {
		config.PrivateRegistryPort = DefaultPrivateRegistryPort
	}

	store := NewWorkflowStore()

	// Create source server with optional auth interceptor
	var sourceOpts []grpc.ServerOption
	if config.AuthProvider != nil {
		sourceOpts = append(sourceOpts, grpc.UnaryInterceptor(
			NewJWTAuthInterceptor(config.AuthProvider),
		))
	}
	sourceServer := grpc.NewServer(sourceOpts...)
	sourcesv1.RegisterWorkflowMetadataSourceServiceServer(sourceServer, NewSourceService(store))

	// Create private registry server (no auth needed for tests)
	privateRegistryServer := grpc.NewServer()
	privateRegistryService := NewPrivateRegistryService(store)

	return &Server{
		config:                 config,
		store:                  store,
		sourceServer:           sourceServer,
		privateRegistryServer:  privateRegistryServer,
		privateRegistryService: privateRegistryService,
	}
}

// Start starts both gRPC servers
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return errors.New("server already started")
	}

	// Start source server
	sourceAddr := fmt.Sprintf(":%d", s.config.SourcePort)
	lc := &net.ListenConfig{}
	sourceListener, err := lc.Listen(context.Background(), "tcp", sourceAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on source port %d: %w", s.config.SourcePort, err)
	}
	s.sourceListener = sourceListener

	// Start private registry server
	privateRegistryAddr := fmt.Sprintf(":%d", s.config.PrivateRegistryPort)
	privateRegistryListener, err := lc.Listen(context.Background(), "tcp", privateRegistryAddr)
	if err != nil {
		sourceListener.Close()
		return fmt.Errorf("failed to listen on private registry port %d: %w", s.config.PrivateRegistryPort, err)
	}
	s.privateRegistryListener = privateRegistryListener

	// Serve source requests
	go func() {
		_ = s.sourceServer.Serve(sourceListener)
		// Error is expected when server is stopped gracefully
	}()

	// Serve private registry requests
	go func() {
		_ = s.privateRegistryServer.Serve(privateRegistryListener)
		// Error is expected when server is stopped gracefully
	}()

	s.started = true
	return nil
}

// Stop stops both gRPC servers
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return
	}

	s.sourceServer.GracefulStop()
	s.privateRegistryServer.GracefulStop()

	if s.sourceListener != nil {
		s.sourceListener.Close()
	}
	if s.privateRegistryListener != nil {
		s.privateRegistryListener.Close()
	}

	s.started = false
}

// SourceAddr returns the address of the source service
func (s *Server) SourceAddr() string {
	return fmt.Sprintf("localhost:%d", s.config.SourcePort)
}

// PrivateRegistryAddr returns the address of the private registry service
func (s *Server) PrivateRegistryAddr() string {
	return fmt.Sprintf("localhost:%d", s.config.PrivateRegistryPort)
}

// PrivateRegistryService returns the private registry service for direct manipulation in tests
func (s *Server) PrivateRegistryService() privateregistry.WorkflowDeploymentAction {
	return s.privateRegistryService
}

// Store returns the underlying workflow store for direct inspection in tests
func (s *Server) Store() *WorkflowStore {
	return s.store
}
