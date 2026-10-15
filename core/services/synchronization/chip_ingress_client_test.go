package synchronization

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	chipingresspb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
)

type capturingChipIngressServer struct {
	chipingresspb.UnimplementedChipIngressServer
	lastMD metadata.MD
}

func (s *capturingChipIngressServer) Ping(ctx context.Context, _ *chipingresspb.EmptyRequest) (*chipingresspb.PingResponse, error) {
	s.lastMD, _ = metadata.FromIncomingContext(ctx)
	return &chipingresspb.PingResponse{}, nil
}

func TestChipIngressClient_SendAuthHeaders(t *testing.T) {
	t.Parallel()
	lis, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lis.Close()

	srv := grpc.NewServer()
	capture := &capturingChipIngressServer{}
	chipingresspb.RegisterChipIngressServer(srv, capture)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	client, err := NewChipIngressClient(ChipIngressClientConfig{
		Endpoint:           lis.Addr().String(),
		InsecureConnection: true,
		AuthHeaders:        map[string]string{"X-Beholder-Node-Auth-Token": "x"},
	})
	require.NoError(t, err)
	defer client.Close()

	_, err = client.Ping(t.Context(), &chipingresspb.EmptyRequest{})
	require.NoError(t, err)

	assert.NotNil(t, capture.lastMD)
	assert.Equal(t, []string{"true"}, capture.lastMD.Get("x-include-nop-info"))
	assert.Equal(t, []string{"x"}, capture.lastMD.Get("X-Beholder-Node-Auth-Token"))
}
