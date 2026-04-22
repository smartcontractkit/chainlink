package vault

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// TestLinkingService is a local gRPC fixture that resolves workflow owners to org IDs.
type TestLinkingService struct {
	linkingclient.UnimplementedLinkingServiceServer

	server   *grpc.Server
	listener net.Listener

	mu         sync.RWMutex
	ownerToOrg map[string]string
}

const SharedTestLinkingServiceAddr = "0.0.0.0:18124"

var (
	sharedTestLinkingServiceOnce sync.Once
	sharedTestLinkingService     *TestLinkingService
	errSharedTestLinkingService  error
)

// NewTestLinkingService starts a mock linking service immediately.
func NewTestLinkingService(ownerToOrg map[string]string) (*TestLinkingService, error) {
	return NewTestLinkingServiceOnAddr(ownerToOrg, "0.0.0.0:0")
}

// EnsureSharedTestLinkingServiceStarted starts the shared local linking service once
// on the fixed host port that Dockerized nodes use during local CRE runs.
func EnsureSharedTestLinkingServiceStarted() (*TestLinkingService, error) {
	sharedTestLinkingServiceOnce.Do(func() {
		sharedTestLinkingService, errSharedTestLinkingService = NewTestLinkingServiceOnAddr(nil, SharedTestLinkingServiceAddr)
	})

	return sharedTestLinkingService, errSharedTestLinkingService
}

// NewTestLinkingServiceOnAddr starts a mock linking service on the provided TCP address.
func NewTestLinkingServiceOnAddr(ownerToOrg map[string]string, listenAddr string) (*TestLinkingService, error) {
	if listenAddr == "" {
		listenAddr = "0.0.0.0:0"
	}

	// #nosec G102 -- test-only linking server must be reachable from Dockerized nodes during local CRE runs.
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate linking service listener: %w", err)
	}

	svc := &TestLinkingService{
		listener:   listener,
		ownerToOrg: make(map[string]string, len(ownerToOrg)),
	}
	for owner, orgID := range ownerToOrg {
		svc.ownerToOrg[normalizeWorkflowOwner(owner)] = orgID
	}

	server := grpc.NewServer()
	linkingclient.RegisterLinkingServiceServer(server, svc)
	svc.server = server

	go func() {
		_ = server.Serve(listener)
	}()

	return svc, nil
}

// Close stops the mock gRPC service.
func (s *TestLinkingService) Close() error {
	if s == nil || s.server == nil {
		return nil
	}

	s.server.Stop()
	if s.listener != nil {
		if err := s.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
	}

	return nil
}

// LocalURL returns the host-local gRPC address.
func (s *TestLinkingService) LocalURL() string {
	return s.baseURL("127.0.0.1")
}

// DockerURL returns the address Dockerized nodes can use.
func (s *TestLinkingService) DockerURL() string {
	return s.baseURL(strings.TrimPrefix(framework.HostDockerInternal(), "http://"))
}

// SetOwnerOrg installs or updates a workflow-owner mapping.
func (s *TestLinkingService) SetOwnerOrg(owner, orgID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ownerToOrg[normalizeWorkflowOwner(owner)] = orgID
}

// GetOrganizationFromWorkflowOwner implements the linking-service gRPC API.
func (s *TestLinkingService) GetOrganizationFromWorkflowOwner(_ context.Context, req *linkingclient.GetOrganizationFromWorkflowOwnerRequest) (*linkingclient.GetOrganizationFromWorkflowOwnerResponse, error) {
	owner := normalizeWorkflowOwner(req.GetWorkflowOwner())

	s.mu.RLock()
	orgID, ok := s.ownerToOrg[owner]
	s.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "workflow owner %q not linked", req.GetWorkflowOwner())
	}

	return &linkingclient.GetOrganizationFromWorkflowOwnerResponse{
		OrganizationId: orgID,
	}, nil
}

func (s *TestLinkingService) baseURL(host string) string {
	if s == nil || s.listener == nil {
		return ""
	}

	tcpAddr, ok := s.listener.Addr().(*net.TCPAddr)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%d", host, tcpAddr.Port)
}

func normalizeWorkflowOwner(owner string) string {
	owner = strings.ToLower(strings.TrimSpace(owner))
	return strings.TrimPrefix(owner, "0x")
}
