package grpcsourcemock

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/privateregistry"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// TestContainer wraps the mock gRPC server for use in integration tests
// It runs the server in-process and provides the URL that Docker containers
// can use to connect to it via host.docker.internal
type TestContainer struct {
	server       *Server
	authProvider *MockNodeAuthProvider
	mu           sync.Mutex
	started      bool
}

// TestContainerConfig contains configuration for the test container
type TestContainerConfig struct {
	// SourcePort is the port for the WorkflowMetadataSourceService (default: 8544)
	SourcePort int
	// PrivateRegistryPort is the port for the private registry API (default: 8545)
	PrivateRegistryPort int
	// TrustedKeys is the initial set of trusted public keys
	TrustedKeys []ed25519.PublicKey
	// RejectAllAuth if true, will reject all authentication attempts
	RejectAllAuth bool
}

// NewTestContainer creates a new test container for the mock gRPC server
func NewTestContainer(config TestContainerConfig) *TestContainer {
	if config.SourcePort == 0 {
		config.SourcePort = DefaultSourcePort
	}
	if config.PrivateRegistryPort == 0 {
		config.PrivateRegistryPort = DefaultPrivateRegistryPort
	}

	var authProvider NodeAuthProvider
	var mockAuthProvider *MockNodeAuthProvider

	switch {
	case config.RejectAllAuth:
		authProvider = &RejectAllAuthProvider{}
	case len(config.TrustedKeys) > 0:
		// Use MockNodeAuthProvider with specific trusted keys
		mockAuthProvider = NewMockNodeAuthProvider()
		for _, key := range config.TrustedKeys {
			mockAuthProvider.AddTrustedKey(key)
		}
		authProvider = mockAuthProvider
	default:
		// Accept all valid JWTs when no specific keys are provided
		// This is useful for tests where we don't know node keys ahead of time
		authProvider = &AcceptAllAuthProvider{}
	}

	server := NewServer(ServerConfig{
		SourcePort:          config.SourcePort,
		PrivateRegistryPort: config.PrivateRegistryPort,
		AuthProvider:        authProvider,
	})

	return &TestContainer{
		server:       server,
		authProvider: mockAuthProvider,
	}
}

// Start starts the mock server
func (tc *TestContainer) Start(ctx context.Context) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.started {
		return errors.New("test container already started")
	}

	if err := tc.server.Start(); err != nil {
		return fmt.Errorf("failed to start mock server: %w", err)
	}

	tc.started = true
	return nil
}

// Stop stops the mock server
func (tc *TestContainer) Stop(ctx context.Context) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if !tc.started {
		return nil
	}

	tc.server.Stop()
	tc.started = false
	return nil
}

// SourceURL returns the URL that Docker containers should use to connect to the source service.
// Uses framework.HostDockerInternal() which handles both local (Mac) and CI (Linux) environments.
func (tc *TestContainer) SourceURL() string {
	// Same pattern as telemetry endpoint in lib/cre/don/config/config.go:279
	host := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
	return fmt.Sprintf("%s:%d", host, tc.server.config.SourcePort)
}

// PrivateRegistryURL returns the URL that can be used to connect to the private registry service
// This is typically used from the test process, not from Docker containers
func (tc *TestContainer) PrivateRegistryURL() string {
	return tc.server.PrivateRegistryAddr()
}

// InternalSourceURL returns the source URL for use within Docker containers
// This is an alias for SourceURL for clarity
func (tc *TestContainer) InternalSourceURL() string {
	return tc.SourceURL()
}

// PrivateRegistryService returns the private registry service for direct manipulation in tests
func (tc *TestContainer) PrivateRegistryService() privateregistry.WorkflowDeploymentAction {
	return tc.server.PrivateRegistryService()
}

// Store returns the underlying workflow store for direct inspection in tests
func (tc *TestContainer) Store() *WorkflowStore {
	return tc.server.Store()
}

// AuthProvider returns the mock auth provider for managing trusted keys
// Returns nil if RejectAllAuth was set in the config
func (tc *TestContainer) AuthProvider() *MockNodeAuthProvider {
	return tc.authProvider
}

// AddTrustedKey adds a public key to the trusted list
// This is a no-op if RejectAllAuth was set in the config
func (tc *TestContainer) AddTrustedKey(publicKey ed25519.PublicKey) {
	if tc.authProvider != nil {
		tc.authProvider.AddTrustedKey(publicKey)
	}
}

// SetTrustedKeys replaces all trusted keys with the provided list
// This is a no-op if RejectAllAuth was set in the config
func (tc *TestContainer) SetTrustedKeys(publicKeys []ed25519.PublicKey) {
	if tc.authProvider != nil {
		tc.authProvider.SetTrustedKeys(publicKeys)
	}
}
