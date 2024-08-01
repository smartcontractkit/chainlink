package ocr3

import (
	"context"
	"math"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr3pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ocr3"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*contractTransmitterClient)(nil)

type contractTransmitterClient struct {
	*net.BrokerExt
	grpc ocr3pb.ContractTransmitterClient
}

func NewContractTransmitterClient(broker *net.BrokerExt, cc grpc.ClientConnInterface) *contractTransmitterClient {
	return &contractTransmitterClient{
		BrokerExt: broker,
		grpc:      ocr3pb.NewContractTransmitterClient(cc),
	}
}

func (c *contractTransmitterClient) Transmit(ctx context.Context, configDigest libocr.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte], signatures []libocr.AttributedOnchainSignature) error {
	cd := [32]byte(configDigest)
	req := &ocr3pb.TransmitRequest{
		ConfigDigest: cd[:],
		SeqNr:        seqNr,
		ReportWithInfo: &ocr3pb.ReportWithInfo{
			Report: reportWithInfo.Report,
			Info:   reportWithInfo.Info,
		},
	}
	for _, s := range signatures {
		req.Signatures = append(req.Signatures, &ocr3pb.Signature{
			Signature: s.Signature,
			Signer:    uint32(s.Signer),
		})
	}

	_, err := c.grpc.Transmit(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *contractTransmitterClient) FromAccount() (libocr.Account, error) {
	reply, err := c.grpc.FromAccount(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return libocr.Account(reply.Account), nil
}

var _ ocr3pb.ContractTransmitterServer = (*contractTransmitterServer)(nil)

type contractTransmitterServer struct {
	ocr3pb.UnimplementedContractTransmitterServer
	impl ocr3types.ContractTransmitter[[]byte]
}

func (c *contractTransmitterServer) Transmit(ctx context.Context, request *ocr3pb.TransmitRequest) (*emptypb.Empty, error) {
	if l := len(request.ConfigDigest); l != 32 {
		return nil, pb.ErrConfigDigestLen(l)
	}
	cd := libocr.ConfigDigest([32]byte(request.ConfigDigest))

	info := ocr3types.ReportWithInfo[[]byte]{}
	if request.ReportWithInfo != nil {
		info.Report = request.ReportWithInfo.Report
		info.Info = request.ReportWithInfo.Info
	}

	signatures := []libocr.AttributedOnchainSignature{}
	for _, s := range request.Signatures {
		if s.Signer > math.MaxUint8 {
			return nil, pb.ErrUint8Bounds{Name: "Signer", U: s.Signer}
		}
		signatures = append(signatures, libocr.AttributedOnchainSignature{
			Signature: s.Signature,
			Signer:    commontypes.OracleID(s.Signer),
		})
	}

	return &emptypb.Empty{}, c.impl.Transmit(ctx, cd, request.SeqNr, info, signatures)
}

func (c *contractTransmitterServer) FromAccount(ctx context.Context, request *emptypb.Empty) (*ocr3pb.FromAccountReply, error) {
	a, err := c.impl.FromAccount()
	if err != nil {
		return nil, err
	}
	return &ocr3pb.FromAccountReply{Account: string(a)}, nil
}

func NewContractTransmitterServer(transmitter ocr3types.ContractTransmitter[[]byte]) *contractTransmitterServer {
	return &contractTransmitterServer{impl: transmitter}
}
