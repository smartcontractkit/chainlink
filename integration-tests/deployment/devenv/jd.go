package devenv

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	csav1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/csa/v1"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
)

type JDConfig struct {
	GRPC  string
	WSRPC string
	creds credentials.TransportCredentials
}

func NewJDConnection(cfg JDConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	// TODO: add auth details
	if cfg.creds != nil {
		opts = append(opts, grpc.WithTransportCredentials(cfg.creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

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
}

func NewJDClient(cfg JDConfig) (deployment.OffchainClient, error) {
	conn, err := NewJDConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}
	return JobDistributor{
		WSRPC:             cfg.WSRPC,
		NodeServiceClient: nodev1.NewNodeServiceClient(conn),
		JobServiceClient:  jobv1.NewJobServiceClient(conn),
		CSAServiceClient:  csav1.NewCSAServiceClient(conn),
	}, err
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
