package internal

import (
	"context"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.Codec = (*codecClient)(nil)

// NewCodecTestClient is a test client for [types.Codec]
// internal users should instantiate a client directly and set all private fields
func NewCodecTestClient(conn *grpc.ClientConn) types.Codec {
	return &codecClient{grpc: pb.NewCodecClient(conn)}
}

type codecClient struct {
	*BrokerExt
	grpc pb.CodecClient
}

func (c *codecClient) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	versionedParams, err := EncodeVersionedBytes(item, CurrentEncodingVersion)
	if err != nil {
		return nil, err
	}

	reply, err := c.grpc.GetEncoding(ctx, &pb.GetEncodingRequest{
		Params:   versionedParams,
		ItemType: itemType,
	})

	if err != nil {
		return nil, wrapRPCErr(err)
	}

	return reply.RetVal, nil
}

func (c *codecClient) Decode(ctx context.Context, raw []byte, into any, itemType string) error {
	request := &pb.GetDecodingRequest{
		Encoded:             raw,
		ItemType:            itemType,
		WireEncodingVersion: CurrentEncodingVersion,
	}
	resp, err := c.grpc.GetDecoding(ctx, request)
	if err != nil {
		return wrapRPCErr(err)
	}

	return DecodeVersionedBytes(into, resp.RetVal)
}

func (c *codecClient) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	res, err := c.grpc.GetMaxSize(ctx, &pb.GetMaxSizeRequest{N: int32(n), ItemType: itemType, ForEncoding: true})
	if err != nil {
		return 0, wrapRPCErr(err)
	}

	return int(res.SizeInBytes), nil
}

func (c *codecClient) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	res, err := c.grpc.GetMaxSize(ctx, &pb.GetMaxSizeRequest{N: int32(n), ItemType: itemType, ForEncoding: false})
	if err != nil {
		return 0, wrapRPCErr(err)
	}

	return int(res.SizeInBytes), nil
}

var _ pb.CodecServer = (*codecServer)(nil)

func NewCodecServer(impl types.Codec) pb.CodecServer {
	return &codecServer{impl: impl}
}

type codecServer struct {
	pb.UnimplementedCodecServer
	impl types.Codec
}

func (c *codecServer) GetEncoding(ctx context.Context, req *pb.GetEncodingRequest) (*pb.GetEncodingResponse, error) {
	encodedType, err := getEncodedType(req.ItemType, c.impl, true)
	if err != nil {
		return nil, err
	}

	if err = DecodeVersionedBytes(encodedType, req.Params); err != nil {
		return nil, err
	}

	encoded, err := c.impl.Encode(ctx, encodedType, req.ItemType)
	return &pb.GetEncodingResponse{RetVal: encoded}, err
}

func (c *codecServer) GetDecoding(ctx context.Context, req *pb.GetDecodingRequest) (*pb.GetDecodingResponse, error) {
	encodedType, err := getEncodedType(req.ItemType, c.impl, false)
	if err != nil {
		return nil, err
	}

	err = c.impl.Decode(ctx, req.Encoded, encodedType, req.ItemType)
	if err != nil {
		return nil, err
	}

	versionBytes, err := EncodeVersionedBytes(encodedType, req.WireEncodingVersion)
	return &pb.GetDecodingResponse{RetVal: versionBytes}, err
}

func (c *codecServer) GetMaxSize(ctx context.Context, req *pb.GetMaxSizeRequest) (*pb.GetMaxSizeResponse, error) {
	var sizeFn func(context.Context, int, string) (int, error)
	if req.ForEncoding {
		sizeFn = c.impl.GetMaxEncodingSize
	} else {
		sizeFn = c.impl.GetMaxDecodingSize
	}

	maxSize, err := sizeFn(ctx, int(req.N), req.ItemType)
	if err != nil {
		return nil, err
	}
	return &pb.GetMaxSizeResponse{SizeInBytes: int32(maxSize)}, nil
}

func getEncodedType(itemType string, possibleTypeProvider any, forEncoding bool) (any, error) {
	if tp, ok := possibleTypeProvider.(types.TypeProvider); ok {
		return tp.CreateType(itemType, forEncoding)
	}

	return &map[string]any{}, nil
}
