package chainreader

import (
	"context"
	jsonv1 "encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	jsonv2 "github.com/go-json-experiment/json"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.ChainReader = (*Client)(nil)

// NewChainReaderTestClient is a test client for [types.ChainReader]
// internal users should instantiate a client directly and set all private fields.
func NewChainReaderTestClient(conn *grpc.ClientConn) types.ChainReader {
	return &Client{grpc: pb.NewChainReaderClient(conn)}
}

type Client struct {
	*net.BrokerExt
	grpc pb.ChainReaderClient
}

func NewClient(b *net.BrokerExt, cc grpc.ClientConnInterface) *Client {
	return &Client{BrokerExt: b, grpc: pb.NewChainReaderClient(cc)}
}

// enum of all known encoding formats for versioned data.
const (
	JSONEncodingVersion1 = iota
	JSONEncodingVersion2
	CBOREncodingVersion
)

// Version to be used for encoding (version used for decoding is determined by data received).
const CurrentEncodingVersion = CBOREncodingVersion

func EncodeVersionedBytes(data any, version uint32) (*pb.VersionedBytes, error) {
	var bytes []byte
	var err error

	switch version {
	case JSONEncodingVersion1:
		bytes, err = jsonv1.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", types.ErrInvalidType, err)
		}
	case JSONEncodingVersion2:
		bytes, err = jsonv2.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", types.ErrInvalidType, err)
		}
	case CBOREncodingVersion:
		enco := cbor.CoreDetEncOptions()
		enco.Time = cbor.TimeRFC3339Nano
		var enc cbor.EncMode
		enc, err = enco.EncMode()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", types.ErrInternal, err)
		}
		bytes, err = enc.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", types.ErrInvalidType, err)
		}
	default:
		return nil, fmt.Errorf("%w: unsupported encoding version %d for data %v", types.ErrInvalidEncoding, version, data)
	}

	return &pb.VersionedBytes{Version: version, Data: bytes}, nil
}

func DecodeVersionedBytes(res any, vData *pb.VersionedBytes) error {
	var err error
	switch vData.Version {
	case JSONEncodingVersion1:
		err = jsonv1.Unmarshal(vData.Data, res)
	case JSONEncodingVersion2:
		err = jsonv2.Unmarshal(vData.Data, res)
	case CBOREncodingVersion:
		decopt := cbor.DecOptions{UTF8: cbor.UTF8DecodeInvalid}
		var dec cbor.DecMode
		dec, err = decopt.DecMode()
		if err != nil {
			return fmt.Errorf("%w: %w", types.ErrInternal, err)
		}
		err = dec.Unmarshal(vData.Data, res)
	default:
		return fmt.Errorf("unsupported encoding version %d for versionedData %v", vData.Version, vData.Data)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", types.ErrInvalidType, err)
	}
	return nil
}

func (c *Client) GetLatestValue(ctx context.Context, contractName, method string, params, retVal any) error {
	versionedParams, err := EncodeVersionedBytes(params, CurrentEncodingVersion)
	if err != nil {
		return err
	}

	reply, err := c.grpc.GetLatestValue(ctx, &pb.GetLatestValueRequest{ContractName: contractName, Method: method, Params: versionedParams})
	if err != nil {
		return net.WrapRPCErr(err)
	}

	return DecodeVersionedBytes(retVal, reply.RetVal)
}

func (c *Client) Bind(ctx context.Context, bindings []types.BoundContract) error {
	pbBindings := make([]*pb.BoundContract, len(bindings))
	for i, b := range bindings {
		pbBindings[i] = &pb.BoundContract{Address: b.Address, Name: b.Name, Pending: b.Pending}
	}
	_, err := c.grpc.Bind(ctx, &pb.BindRequest{Bindings: pbBindings})
	return net.WrapRPCErr(err)
}

var _ pb.ChainReaderServer = (*Server)(nil)

func NewServer(impl types.ChainReader) pb.ChainReaderServer {
	return &Server{impl: impl}
}

type Server struct {
	pb.UnimplementedChainReaderServer
	impl types.ChainReader
}

func (c *Server) GetLatestValue(ctx context.Context, request *pb.GetLatestValueRequest) (*pb.GetLatestValueReply, error) {
	contractName := request.ContractName
	params, err := getContractEncodedType(contractName, request.Method, c.impl, true)
	if err != nil {
		return nil, err
	}

	if err = DecodeVersionedBytes(params, request.Params); err != nil {
		return nil, err
	}

	retVal, err := getContractEncodedType(contractName, request.Method, c.impl, false)
	if err != nil {
		return nil, err
	}
	err = c.impl.GetLatestValue(ctx, contractName, request.Method, params, retVal)
	if err != nil {
		return nil, err
	}

	encodedRetVal, err := EncodeVersionedBytes(retVal, request.Params.Version)
	if err != nil {
		return nil, err
	}

	return &pb.GetLatestValueReply{RetVal: encodedRetVal}, nil
}

func (c *Server) Bind(ctx context.Context, bindings *pb.BindRequest) (*emptypb.Empty, error) {
	tBindings := make([]types.BoundContract, len(bindings.Bindings))
	for i, b := range bindings.Bindings {
		tBindings[i] = types.BoundContract{Address: b.Address, Name: b.Name, Pending: b.Pending}
	}

	return &emptypb.Empty{}, c.impl.Bind(ctx, tBindings)
}

func getContractEncodedType(contractName, itemType string, possibleTypeProvider any, forEncoding bool) (any, error) {
	if ctp, ok := possibleTypeProvider.(types.ContractTypeProvider); ok {
		return ctp.CreateContractType(contractName, itemType, forEncoding)
	}

	return &map[string]any{}, nil
}
