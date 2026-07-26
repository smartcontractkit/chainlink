package bridgeconn

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/bridgeconn/streamspb"
)

// grpcStreamClient adapts streamspb's generated Subscribe stream to eaStreamClient.
type grpcStreamClient struct {
	conn   *grpc.ClientConn
	stream streamspb.StreamService_SubscribeClient
}

// dialGRPCStream opens a gRPC connection to target and starts the bidirectional
// Subscribe stream. When useTLS is true (bridge URL scheme is https) the connection
// uses TLS with the host's system root CAs; otherwise it is plaintext. Neither mode
// uses application tokens or client certificates.
func dialGRPCStream(ctx context.Context, target string, useTLS bool) (eaStreamClient, error) {
	creds := insecure.NewCredentials()
	if useTLS {
		creds = credentials.NewClientTLSFromCert(nil, "")
	}
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	stream, err := streamspb.NewStreamServiceClient(conn).Subscribe(ctx)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &grpcStreamClient{conn: conn, stream: stream}, nil
}

func (c *grpcStreamClient) Send(req *streamspb.SubscribeRequest) error {
	return c.stream.Send(req)
}

func (c *grpcStreamClient) Recv() (*streamspb.SubscribeResponse, error) {
	return c.stream.Recv()
}

func (c *grpcStreamClient) Close() error {
	_ = c.stream.CloseSend()
	return c.conn.Close()
}
