package ccip

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// TokenDataReaderGRPCClient implements [cciptypes.TokenDataReaderReader] by wrapping a
// [ccippb.TokenDataReaderReaderGRPCClient] grpc client.
// It is used by a ReportingPlugin to call theTokenDataReaderReader service, which
// is hosted by the relayer
type TokenDataReaderGRPCClient struct {
	client ccippb.TokenDataReaderClient
	conn   grpc.ClientConnInterface
}

func NewTokenDataReaderGRPCClient(cc grpc.ClientConnInterface) *TokenDataReaderGRPCClient {
	return &TokenDataReaderGRPCClient{client: ccippb.NewTokenDataReaderClient(cc), conn: cc}
}

// TokenDataReaderGRPCServer implements [ccippb.TokenDataReaderReaderServer] by wrapping a
// [cciptypes.TokenDataReaderReader] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [TokenDataReaderGRPCClient]
type TokenDataReaderGRPCServer struct {
	ccippb.UnimplementedTokenDataReaderServer

	impl cciptypes.TokenDataReader

	deps []io.Closer
}

func NewTokenDataReaderGRPCServer(impl cciptypes.TokenDataReader) *TokenDataReaderGRPCServer {
	return &TokenDataReaderGRPCServer{impl: impl, deps: []io.Closer{impl}}
}

// ensure interface is implemented
var _ ccippb.TokenDataReaderServer = (*TokenDataReaderGRPCServer)(nil)
var _ cciptypes.TokenDataReader = (*TokenDataReaderGRPCClient)(nil)

func (t *TokenDataReaderGRPCClient) ClientConn() grpc.ClientConnInterface {
	return t.conn
}

// Close implements ccip.TokenDataReader.
func (t *TokenDataReaderGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), t.client)
}

// ReadTokenData implements ccip.TokenDataReader.
func (t *TokenDataReaderGRPCClient) ReadTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) (tokenData []byte, err error) {
	resp, err := t.client.ReadTokenData(ctx, &ccippb.TokenDataRequest{
		Msg: &ccippb.EVM2EVMOnRampCCIPSendRequestedWithMeta{
			EvmToEvmMsg:    evm2EVMMessagePB(msg.EVM2EVMMessage),
			BlockTimestamp: timestamppb.New(msg.BlockTimestamp),
			Executed:       msg.Executed,
			Finalized:      msg.Finalized,
			LogIndex:       uint64(msg.LogIndex),
			TxHash:         msg.TxHash,
		},
		TokenIndex: uint64(tokenIndex),
	})
	if err != nil {
		return nil, err
	}
	return resp.TokenData, nil
}

// Server implementation

// Close implements ccippb.TokenDataReaderServer.
func (t *TokenDataReaderGRPCServer) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(t.deps).Close()
}

// ReadTokenData implements ccippb.TokenDataReaderServer.
func (t *TokenDataReaderGRPCServer) ReadTokenData(ctx context.Context, req *ccippb.TokenDataRequest) (*ccippb.TokenDataResponse, error) {
	evmMsg, err := evm2EVMMessage(req.Msg.EvmToEvmMsg)
	if err != nil {
		return nil, err
	}
	msg := cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: evmMsg,
		BlockTimestamp: req.Msg.BlockTimestamp.AsTime(),
		Executed:       req.Msg.Executed,
		Finalized:      req.Msg.Finalized,
		LogIndex:       uint(req.Msg.LogIndex),
		TxHash:         req.Msg.TxHash,
	}

	tokenData, err := t.impl.ReadTokenData(ctx, msg, int(req.TokenIndex))
	if err != nil {
		return nil, err
	}
	return &ccippb.TokenDataResponse{TokenData: tokenData}, nil
}

// AddDep adds a dependency to the server that will be closed when the server is closed.
func (t *TokenDataReaderGRPCServer) AddDep(dep io.Closer) *TokenDataReaderGRPCServer {
	t.deps = append(t.deps, dep)
	return t
}
