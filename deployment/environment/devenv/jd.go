package devenv

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
)

type JDConfig struct {
	GRPC     string
	WSRPC    string
	Creds    credentials.TransportCredentials
	Auth     oauth2.TokenSource
	NodeInfo []NodeInfo
}

// JobDistributor implements the OffchainClient interface in CLDF and wraps the CLDF JD client and add DON functionality.
// The CLDF JD client does not have the DON functionality, so we wrap it here.
type JobDistributor struct {
	*jd.JobDistributor
	don *DON
}

// NewJDClient creates a new Job Distributor client with the provided configuration.
func NewJDClient(ctx context.Context, cfg JDConfig) (cldf_offchain.Client, error) {
	jdConfig := jd.JDConfig{
		GRPC:  cfg.GRPC,
		WSRPC: cfg.WSRPC,
		Creds: cfg.Creds,
		Auth:  cfg.Auth,
	}
	jdClient, err := jd.NewJDClient(jdConfig)
	if err != nil {
		return nil, err
	}
	donJDClient := &JobDistributor{
		JobDistributor: jdClient,
	}
	if cfg.NodeInfo != nil && len(cfg.NodeInfo) > 0 {
		donJDClient.don, err = NewRegisteredDON(ctx, cfg.NodeInfo, *donJDClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create registered DON: %w", err)
		}
	}
	return donJDClient, err
}

func (jd JobDistributor) ReplayLogs(selectorToBlock map[uint64]uint64) error {
	return jd.don.ReplayAllLogs(selectorToBlock)
}

// ProposeJob proposes jobs through the jobService and accepts the proposed job on selected node based on ProposeJobRequest.NodeId
func (jd JobDistributor) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	res, err := jd.JobDistributor.ProposeJob(ctx, in, opts...)
	if err != nil {
		return nil, err
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
