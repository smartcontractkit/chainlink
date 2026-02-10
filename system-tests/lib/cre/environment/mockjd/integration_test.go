package mockjd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cldf_jd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

func TestMockJDWithCLDFClient(t *testing.T) {
	ctx := context.Background()
	csaKey := "abcd1234567890abcdef"

	server, err := NewServer(csaKey)
	require.NoError(t, err)
	require.NoError(t, server.Start())
	defer server.Stop()

	jdConfig := cldf_jd.JDConfig{
		GRPC:  server.Addr(),
		WSRPC: server.Addr(),
		Creds: insecure.NewCredentials(),
	}

	jdClient, err := cldf_jd.NewJDClient(jdConfig)
	require.NoError(t, err)

	csaPubKey, err := jdClient.GetCSAPublicKey(ctx)
	require.NoError(t, err)
	assert.Equal(t, csaKey, csaPubKey)

	regResp, err := jdClient.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		Name:      "cldf-test-node",
		PublicKey: "cldf-node-csa",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.Node.Id)

	getResp, err := jdClient.GetNode(ctx, &nodev1.GetNodeRequest{Id: regResp.Node.Id})
	require.NoError(t, err)
	assert.Equal(t, "cldf-test-node", getResp.Node.Name)
}

func TestMockJDChainConfigs(t *testing.T) {
	ctx := context.Background()
	csaKey := "test-key"

	server, err := NewServer(csaKey)
	require.NoError(t, err)
	require.NoError(t, server.Start())
	defer server.Stop()

	conn, err := grpc.NewClient(server.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := nodev1.NewNodeServiceClient(conn)

	regResp, err := client.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		Name:      "chain-node",
		PublicKey: "chain-node-csa",
	})
	require.NoError(t, err)

	listResp, err := client.ListNodeChainConfigs(ctx, &nodev1.ListNodeChainConfigsRequest{
		Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeIds: []string{regResp.Node.Id},
		},
	})
	require.NoError(t, err)
	assert.Empty(t, listResp.ChainConfigs)
}
