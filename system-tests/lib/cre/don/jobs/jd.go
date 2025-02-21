package jobs

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

const (
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
)

func ReinitialiseJDClients(ctfEnv *deployment.Environment, jdOutput *jd.Output, nodeOutputs ...*types.WrappedNodeOutput) (*deployment.Environment, error) {
	offchainClients := make([]deployment.OffchainClient, len(nodeOutputs))

	for i, nodeOutput := range nodeOutputs {
		nodeInfo, err := node.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get node info")
		}

		jdConfig := devenv.JDConfig{
			GRPC:     jdOutput.HostGRPCUrl,
			WSRPC:    jdOutput.DockerWSRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		}

		offChain, err := devenv.NewJDClient(context.Background(), jdConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create JD client")
		}

		offchainClients[i] = offChain
	}

	// we don't really care, which instance we set here, since there's only one
	// what's important is that we create a new JD client for each DON, because
	// that authenticates JD with each node
	ctfEnv.Offchain = offchainClients[0]

	return ctfEnv, nil
}
