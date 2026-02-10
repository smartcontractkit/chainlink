package grpcsourcemock

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"log/slog"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	nodeauthgrpc "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/grpc"
	"github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
)

// MockNodeAuthProvider is a mock implementation of NodeAuthProvider for testing
type MockNodeAuthProvider struct {
	mu             sync.RWMutex
	trustedPubKeys map[string]bool
}

// NewMockNodeAuthProvider creates a new MockNodeAuthProvider
func NewMockNodeAuthProvider() *MockNodeAuthProvider {
	return &MockNodeAuthProvider{
		trustedPubKeys: make(map[string]bool),
	}
}

// AddTrustedKey adds a public key to the trusted list
func (m *MockNodeAuthProvider) AddTrustedKey(publicKey ed25519.PublicKey) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.trustedPubKeys[hex.EncodeToString(publicKey)] = true
}

// RemoveTrustedKey removes a public key from the trusted list
func (m *MockNodeAuthProvider) RemoveTrustedKey(publicKey ed25519.PublicKey) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.trustedPubKeys, hex.EncodeToString(publicKey))
}

// ClearTrustedKeys removes all trusted keys
func (m *MockNodeAuthProvider) ClearTrustedKeys() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.trustedPubKeys = make(map[string]bool)
}

// SetTrustedKeys replaces all trusted keys with the provided list
func (m *MockNodeAuthProvider) SetTrustedKeys(publicKeys []ed25519.PublicKey) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.trustedPubKeys = make(map[string]bool)
	for _, pk := range publicKeys {
		m.trustedPubKeys[hex.EncodeToString(pk)] = true
	}
}

// IsNodePubKeyTrusted checks if a node's public key is trusted
func (m *MockNodeAuthProvider) IsNodePubKeyTrusted(ctx context.Context, publicKey ed25519.PublicKey) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.trustedPubKeys[hex.EncodeToString(publicKey)], nil
}

// RejectAllAuthProvider is an implementation that rejects all public keys
// Used for testing graceful auth failure handling
type RejectAllAuthProvider struct{}

// IsNodePubKeyTrusted always returns false for RejectAllAuthProvider
func (r *RejectAllAuthProvider) IsNodePubKeyTrusted(ctx context.Context, publicKey ed25519.PublicKey) (bool, error) {
	return false, nil
}

// AcceptAllAuthProvider is an implementation that accepts all public keys
// Used for testing when we don't know node keys ahead of time
type AcceptAllAuthProvider struct{}

// IsNodePubKeyTrusted always returns true for AcceptAllAuthProvider
func (a *AcceptAllAuthProvider) IsNodePubKeyTrusted(ctx context.Context, publicKey ed25519.PublicKey) (bool, error) {
	return true, nil
}

// NewJWTAuthInterceptor creates a gRPC unary interceptor that validates JWT tokens.
// Uses the nodeauth token extractor from chainlink-common for consistent token extraction.
func NewJWTAuthInterceptor(authProvider NodeAuthProvider) grpc.UnaryServerInterceptor {
	// Create the JWT authenticator with the provided auth provider
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})).With("logger", "grpc_source_mock.JWTAuthInterceptor")
	authenticator := jwt.NewNodeJWTAuthenticator(authProvider, logger)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract token from metadata using the shared token extractor
		token, err := nodeauthgrpc.ExtractBearerToken(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "missing auth: %v", err)
		}

		// Validate the JWT token
		valid, _, err := authenticator.AuthenticateJWT(ctx, token, req)
		if err != nil {
			// Return unauthenticated error without panicking
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
		}

		if !valid {
			return nil, status.Error(codes.Unauthenticated, "invalid authentication")
		}

		// Continue to the handler if authenticated
		return handler(ctx, req)
	}
}
