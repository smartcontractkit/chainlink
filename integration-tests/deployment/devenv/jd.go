package devenv

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type JDConfig struct {
	GRPC     string
	WSRPC    string
	Creds    credentials.TransportCredentials
	nodeInfo []NodeInfo
}

func NewJDConnection(cfg JDConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(cfg.GRPC, grpc.WithTransportCredentials(cfg.Creds))
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
	if cfg.nodeInfo != nil && len(cfg.nodeInfo) > 0 {
		jd.don, err = NewRegisteredDON(ctx, cfg.nodeInfo, *jd)
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
		return "", fmt.Errorf("no keypairs found")
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
		return nil, fmt.Errorf("failed to propose job. err: proposal is nil")
	}
	if jd.don == nil || len(jd.don.Nodes) == 0 {
		return res, nil
	}
	for _, node := range jd.don.Nodes {
		if node.NodeId != in.NodeId {
			continue
		}
		// TODO : is there a way to accept the job with proposal id?
		if err := node.AcceptJob(ctx, res.Proposal.Spec); err != nil {
			return nil, fmt.Errorf("failed to accept job. err: %w", err)
		}
	}
	return res, nil
}
