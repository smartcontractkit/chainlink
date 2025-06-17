package devenv

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type GAPConfig struct {
	Token      string
	Repository string
}

type JDConfig struct {
	GRPC     string
	WSRPC    string
	Creds    credentials.TransportCredentials
	Auth     oauth2.TokenSource
	GAP      *GAPConfig
	NodeInfo []NodeInfo
}

func authTokenInterceptor(source oauth2.TokenSource) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		token, err := source.Token()
		if err != nil {
			return err
		}

		return invoker(
			metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken),
			method, req, reply, cc, opts...,
		)
	}
}

func gapTokenInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(
			metadata.AppendToOutgoingContext(ctx, "x-authorization-github-jwt", "Bearer "+token),
			method, req, reply, cc, opts...,
		)
	}
}

func gapRepositoryInterceptor(repository string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(
			metadata.AppendToOutgoingContext(ctx, "x-repository", repository),
			method, req, reply, cc, opts...,
		)
	}
}

func NewJDConnection(cfg JDConfig) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	interceptors := []grpc.UnaryClientInterceptor{}

	if cfg.Creds != nil {
		opts = append(opts, grpc.WithTransportCredentials(cfg.Creds))
	}
	if cfg.Auth != nil {
		interceptors = append(interceptors, authTokenInterceptor(cfg.Auth))
	}
	if cfg.GAP != nil && cfg.GAP.Token != "" && cfg.GAP.Repository != "" {
		interceptors = append(interceptors, gapTokenInterceptor(cfg.GAP.Token), gapRepositoryInterceptor(cfg.GAP.Repository))
	}

	if len(interceptors) > 0 {
		opts = append(opts, grpc.WithChainUnaryInterceptor(interceptors...))
	}

	conn, err := grpc.NewClient(cfg.GRPC, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}

	return conn, nil
}

type JobDistributor struct {
	WSRPC string
	nodev1.NodeServiceClient
	jobv1.JobServiceClient
	csav1.CSAServiceClient
	don *DON
}

func NewJDClient(ctx context.Context, cfg JDConfig) (deployment.OffchainClient, error) {
	conn, err := NewJDConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}
	jd := &JobDistributor{
		WSRPC:             cfg.WSRPC,
		NodeServiceClient: nodev1.NewNodeServiceClient(conn),
		JobServiceClient:  jobv1.NewJobServiceClient(conn),
		CSAServiceClient:  csav1.NewCSAServiceClient(conn),
	}
	if cfg.NodeInfo != nil && len(cfg.NodeInfo) > 0 {
		jd.don, err = NewRegisteredDON(ctx, cfg.NodeInfo, *jd)
		if err != nil {
			return nil, fmt.Errorf("failed to create registered DON: %w", err)
		}
	}
	return jd, err
}

func (jd JobDistributor) GetCSAPublicKey(ctx context.Context) (string, error) {
	keypairs, err := jd.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return "", err
	}
	if keypairs == nil || len(keypairs.Keypairs) == 0 {
		return "", errors.New("no keypairs found")
	}
	csakey := keypairs.Keypairs[0].PublicKey
	return csakey, nil
}

func (jd JobDistributor) ReplayLogs(selectorToBlock map[uint64]uint64) error {
	return jd.don.ReplayAllLogs(selectorToBlock)
}

// ProposeJob proposes jobs through the jobService and accepts the proposed job on selected node based on ProposeJobRequest.NodeId
func (jd JobDistributor) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	res, err := jd.JobServiceClient.ProposeJob(ctx, in, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to propose job. err: %w", err)
	}
	if res.Proposal == nil {
		return nil, errors.New("failed to propose job. err: proposal is nil")
	}
	if jd.don == nil || len(jd.don.Nodes) == 0 {
		return res, nil
	}
	for _, node := range jd.don.Nodes {
		if node.NodeID != in.NodeId {
			continue
		}
		// TODO : is there a way to accept the job with proposal id?
		if err := node.AcceptJob(ctx, res.Proposal.Spec); err != nil {
			return nil, fmt.Errorf("failed to accept job. err: %w", err)
		}
	}
	return res, nil
}
