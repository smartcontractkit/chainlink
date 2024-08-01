package chainwriter

import (
	"context"
	"math/big"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.ChainWriter = (*Client)(nil)

type ClientOpt func(*Client)

type Client struct {
	*goplugin.ServiceClient
	grpc       pb.ChainWriterClient
	encodeWith chainreader.EncodingVersion
}

func NewClient(b *net.BrokerExt, cc grpc.ClientConnInterface, opts ...ClientOpt) *Client {
	client := &Client{
		ServiceClient: goplugin.NewServiceClient(b, cc),
		grpc:          pb.NewChainWriterClient(cc),
		encodeWith:    chainreader.DefaultEncodingVersion,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func WithClientEncoding(version chainreader.EncodingVersion) ClientOpt {
	return func(client *Client) {
		client.encodeWith = version
	}
}

func (c *Client) SubmitTransaction(ctx context.Context, contractName, method string, params any, transactionID, toAddress string, meta *types.TxMeta, value *big.Int) error {
	versionedParams, err := chainreader.EncodeVersionedBytes(params, c.encodeWith)
	if err != nil {
		return err
	}

	req := pb.SubmitTransactionRequest{
		ContractName:  contractName,
		Method:        method,
		Params:        versionedParams,
		TransactionId: transactionID,
		ToAddress:     toAddress,
		Meta:          TxMetaToProto(meta),
		Value:         pb.NewBigIntFromInt(value),
	}

	_, err = c.grpc.SubmitTransaction(ctx, &req)
	if err != nil {
		return net.WrapRPCErr(err)
	}

	return nil
}

func (c *Client) GetTransactionStatus(ctx context.Context, transactionID string) (types.TransactionStatus, error) {
	reply, err := c.grpc.GetTransactionStatus(ctx, &pb.GetTransactionStatusRequest{TransactionId: transactionID})
	if err != nil {
		return types.Unknown, net.WrapRPCErr(err)
	}

	return types.TransactionStatus(reply.TransactionStatus), nil
}

func (c *Client) GetFeeComponents(ctx context.Context) (*types.ChainFeeComponents, error) {
	reply, err := c.grpc.GetFeeComponents(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, net.WrapRPCErr(err)
	}

	return &types.ChainFeeComponents{
		ExecutionFee:        reply.ExecutionFee.Int(),
		DataAvailabilityFee: reply.DataAvailabilityFee.Int(),
	}, nil
}

// Server.

var _ pb.ChainWriterServer = (*Server)(nil)

type ServerOpt func(*Server)

type Server struct {
	pb.UnimplementedChainWriterServer
	impl       types.ChainWriter
	encodeWith chainreader.EncodingVersion
}

func NewServer(impl types.ChainWriter, opts ...ServerOpt) pb.ChainWriterServer {
	server := &Server{
		impl:       impl,
		encodeWith: chainreader.DefaultEncodingVersion,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func WithServerEncoding(version chainreader.EncodingVersion) ServerOpt {
	return func(server *Server) {
		server.encodeWith = version
	}
}

func (s *Server) SubmitTransaction(ctx context.Context, req *pb.SubmitTransactionRequest) (*emptypb.Empty, error) {
	err := s.impl.SubmitTransaction(ctx, req.ContractName, req.Method, req.Params, req.TransactionId, req.ToAddress, TxMetaFromProto(req.Meta), req.Value.Int())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetTransactionStatus(ctx context.Context, req *pb.GetTransactionStatusRequest) (*pb.GetTransactionStatusReply, error) {
	status, err := s.impl.GetTransactionStatus(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	return &pb.GetTransactionStatusReply{TransactionStatus: pb.TransactionStatus(status)}, nil
}

func (s *Server) GetFeeComponents(ctx context.Context, _ *emptypb.Empty) (*pb.GetFeeComponentsReply, error) {
	feeComponents, err := s.impl.GetFeeComponents(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetFeeComponentsReply{
		ExecutionFee:        pb.NewBigIntFromInt(feeComponents.ExecutionFee),
		DataAvailabilityFee: pb.NewBigIntFromInt(feeComponents.DataAvailabilityFee),
	}, nil
}

func RegisterChainWriterService(s *grpc.Server, chainWriter types.ChainWriter) {
	pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: chainWriter})
}
