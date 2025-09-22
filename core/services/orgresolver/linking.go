package orgresolver

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/smartcontractkit/chainlink-common/pkg/logger"
	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type Config struct {
	URL                           string
	TLSEnabled                    bool
	WorkflowRegistryAddress       string
	WorkflowRegistryChainSelector uint64
}

type OrgResolver interface {
	services.ServiceCtx
	Get(ctx context.Context, owner string) (string, error)
}

// orgResolver makes direct calls to the linking service to resolve organization IDs from workflow owners.
// This simplified implementation makes a network call for each Get() request.
type orgResolver struct {
	workflowRegistryAddress       string
	workflowRegistryChainSelector uint64

	client linkingclient.LinkingServiceClient
	conn   *grpc.ClientConn // nil if client was injected
	logger log.Logger
}

// NewOrgResolver creates a new org resolver with the specified configuration
// If client is nil, it will create a gRPC connection using the config
func NewOrgResolver(cfg Config, logger log.Logger) (*orgResolver, error) {
	return NewOrgResolverWithClient(cfg, nil, logger)
}

// NewOrgResolverWithClient creates a new org resolver with an optional injected client (for testing)
func NewOrgResolverWithClient(cfg Config, client linkingclient.LinkingServiceClient, logger log.Logger) (*orgResolver, error) {
	resolver := &orgResolver{
		workflowRegistryAddress:       cfg.WorkflowRegistryAddress,
		workflowRegistryChainSelector: cfg.WorkflowRegistryChainSelector,
		logger:                        logger,
	}

	if client != nil {
		// Use injected client (for testing)
		resolver.client = client
	} else {
		// Create gRPC connection and client
		if cfg.URL == "" {
			return nil, errors.New("URL is required when client is not provided")
		}

		var opts []grpc.DialOption
		if cfg.TLSEnabled {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		conn, err := grpc.NewClient(cfg.URL, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to linking service at %s: %w", cfg.URL, err)
		}

		resolver.conn = conn
		resolver.client = linkingclient.NewLinkingServiceClient(conn)
	}

	return resolver, nil
}

func (o *orgResolver) Get(ctx context.Context, owner string) (string, error) {
	req := linkingclient.GetOrganizationFromWorkflowOwnerRequest{
		WorkflowOwner:           owner,
		WorkflowRegistryAddress: o.workflowRegistryAddress,
		ChainSelector:           o.workflowRegistryChainSelector,
	}

	resp, err := o.client.GetOrganizationFromWorkflowOwner(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch organization from workflow owner: %w", err)
	}

	return resp.OrganizationId, nil
}

func (o *orgResolver) Start(_ context.Context) error {
	return nil
}

func (o *orgResolver) HealthReport() map[string]error {
	return map[string]error{o.Name(): nil}
}

func (o *orgResolver) Close() error {
	if o.conn != nil {
		return o.conn.Close()
	}
	return nil
}

func (o *orgResolver) Name() string {
	return "OrgResolver"
}

func (o *orgResolver) Ready() error {
	return nil
}
