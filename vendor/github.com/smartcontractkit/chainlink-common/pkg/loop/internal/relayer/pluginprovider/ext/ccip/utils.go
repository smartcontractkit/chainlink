package ccip

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type GRPCResourceCloser interface {
	Close(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// shutdownGRPCServer is a helper function to release server resources
// created by a grpc client.
func shutdownGRPCServer(ctx context.Context, rc GRPCResourceCloser) error {
	_, err := rc.Close(ctx, &emptypb.Empty{})
	// due to the handler in the server, it may shutdown before it sends a response to client
	// in that case, we expect the client to receive an Unavailable or Internal error
	if status.Code(err) == codes.Unavailable || status.Code(err) == codes.Internal {
		return nil
	}
	return err
}

func txMetaToPB(in cciptypes.TxMeta) *ccippb.TxMeta {
	return &ccippb.TxMeta{
		BlockTimestampUnixMilli: in.BlockTimestampUnixMilli,
		BlockNumber:             in.BlockNumber,
		TxHash:                  in.TxHash,
		LogIndex:                in.LogIndex,
		Finalized:               finalityStatusToPB(in.Finalized),
	}
}

func txMeta(meta *ccippb.TxMeta) cciptypes.TxMeta {
	return cciptypes.TxMeta{
		BlockTimestampUnixMilli: meta.BlockTimestampUnixMilli,
		BlockNumber:             meta.BlockNumber,
		TxHash:                  meta.TxHash,
		LogIndex:                meta.LogIndex,
		Finalized:               finalityStatus(meta.Finalized),
	}
}

func finalityStatus(finalized ccippb.FinalityStatus) cciptypes.FinalizedStatus {
	switch finalized {
	case ccippb.FinalityStatus_Finalized:
		return cciptypes.FinalizedStatusFinalized
	case ccippb.FinalityStatus_NotFinalized:
		return cciptypes.FinalizedStatusNotFinalized
	default:
		return cciptypes.FinalizedStatusUnknown
	}
}

func finalityStatusToPB(finalized cciptypes.FinalizedStatus) ccippb.FinalityStatus {
	switch finalized {
	case cciptypes.FinalizedStatusFinalized:
		return ccippb.FinalityStatus_Finalized
	case cciptypes.FinalizedStatusNotFinalized:
		return ccippb.FinalityStatus_NotFinalized
	default:
		return ccippb.FinalityStatus_Unknown
	}
}
