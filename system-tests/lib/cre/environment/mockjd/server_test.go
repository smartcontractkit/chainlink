package mockjd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

func TestMockJDServer(t *testing.T) {
	ctx := context.Background()
	csaKey := "test-csa-key-12345"

	server, err := NewServer(csaKey)
	require.NoError(t, err)
	require.NoError(t, server.Start())
	defer server.Stop()

	conn, err := grpc.NewClient(server.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	t.Run("CSA keypairs", func(t *testing.T) {
		client := csav1.NewCSAServiceClient(conn)
		resp, err := client.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Keypairs, 1)
		assert.Equal(t, csaKey, resp.Keypairs[0].PublicKey)
	})

	t.Run("register and get node", func(t *testing.T) {
		client := nodev1.NewNodeServiceClient(conn)

		regResp, err := client.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
			Name:      "test-node",
			PublicKey: "node-csa-key",
		})
		require.NoError(t, err)
		require.NotEmpty(t, regResp.Node.Id)
		assert.Equal(t, "test-node", regResp.Node.Name)
		assert.True(t, regResp.Node.IsConnected)

		getResp, err := client.GetNode(ctx, &nodev1.GetNodeRequest{Id: regResp.Node.Id})
		require.NoError(t, err)
		assert.Equal(t, regResp.Node.Id, getResp.Node.Id)
	})

	t.Run("propose job", func(t *testing.T) {
		nodeClient := nodev1.NewNodeServiceClient(conn)
		regResp, err := nodeClient.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
			Name:      "job-node",
			PublicKey: "job-node-csa",
		})
		require.NoError(t, err)

		jobClient := jobv1.NewJobServiceClient(conn)
		propResp, err := jobClient.ProposeJob(ctx, &jobv1.ProposeJobRequest{
			NodeId: regResp.Node.Id,
			Spec:   "test-job-spec",
		})
		require.NoError(t, err)
		require.NotEmpty(t, propResp.Proposal.Id)
		assert.Equal(t, "test-job-spec", propResp.Proposal.Spec)
	})
}
