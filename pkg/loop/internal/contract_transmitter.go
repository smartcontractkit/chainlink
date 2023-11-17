package internal

import (
	"context"
	"fmt"
	"math"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
)

var _ libocr.ContractTransmitter = (*contractTransmitterClient)(nil)

type contractTransmitterClient struct {
	*brokerExt
	grpc pb.ContractTransmitterClient
}

func (c *contractTransmitterClient) Transmit(ctx context.Context, reportContext libocr.ReportContext, report libocr.Report, signatures []libocr.AttributedOnchainSignature) error {
	req := &pb.TransmitRequest{
		ReportContext: &pb.ReportContext{
			ReportTimestamp: &pb.ReportTimestamp{
				ConfigDigest: reportContext.ReportTimestamp.ConfigDigest[:],
				Epoch:        reportContext.ReportTimestamp.Epoch,
				Round:        uint32(reportContext.ReportTimestamp.Round),
			},
			ExtraHash: reportContext.ExtraHash[:],
		},
		Report: report,
	}
	for _, s := range signatures {
		req.AttributedOnchainSignatures = append(req.AttributedOnchainSignatures,
			&pb.AttributedOnchainSignature{
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

func (c *contractTransmitterClient) LatestConfigDigestAndEpoch(ctx context.Context) (configDigest libocr.ConfigDigest, epoch uint32, err error) {
	var reply *pb.LatestConfigDigestAndEpochReply
	reply, err = c.grpc.LatestConfigDigestAndEpoch(ctx, &pb.LatestConfigDigestAndEpochRequest{})
	if err != nil {
		return
	}
	if l := len(reply.ConfigDigest); l != 32 {
		err = ErrConfigDigestLen(l)
		return
	}
	copy(configDigest[:], reply.ConfigDigest)
	epoch = reply.Epoch
	return
}

func (c *contractTransmitterClient) FromAccount() (libocr.Account, error) {
	ctx, cancel := c.stopCtx()
	defer cancel()

	reply, err := c.grpc.FromAccount(ctx, &pb.FromAccountRequest{})
	if err != nil {
		return "", err
	}
	return libocr.Account(reply.Account), nil
}

var _ pb.ContractTransmitterServer = (*contractTransmitterServer)(nil)

type contractTransmitterServer struct {
	pb.UnimplementedContractTransmitterServer
	impl libocr.ContractTransmitter
}

func (c *contractTransmitterServer) Transmit(ctx context.Context, request *pb.TransmitRequest) (*pb.TransmitReply, error) {
	var reportCtx libocr.ReportContext
	if l := len(request.ReportContext.ReportTimestamp.ConfigDigest); l != 32 {
		return nil, ErrConfigDigestLen(l)
	}
	copy(reportCtx.ConfigDigest[:], request.ReportContext.ReportTimestamp.ConfigDigest)
	reportCtx.Epoch = request.ReportContext.ReportTimestamp.Epoch
	if request.ReportContext.ReportTimestamp.Round > math.MaxUint8 {
		return nil, ErrUint8Bounds{Name: "Round", U: request.ReportContext.ReportTimestamp.Round}
	}
	reportCtx.Round = uint8(request.ReportContext.ReportTimestamp.Round)
	if l := len(request.ReportContext.ExtraHash); l != 32 {
		return nil, fmt.Errorf("invalid ExtraHash len %d; must be 32", l)
	}
	copy(reportCtx.ExtraHash[:], request.ReportContext.ExtraHash)
	var sigs []libocr.AttributedOnchainSignature
	for _, s := range request.AttributedOnchainSignatures {
		if s.Signer > math.MaxUint8 {
			return nil, ErrUint8Bounds{Name: "Signer", U: s.Signer}
		}
		sigs = append(sigs, libocr.AttributedOnchainSignature{
			Signature: s.Signature,
			Signer:    commontypes.OracleID(s.Signer),
		})
	}
	return &pb.TransmitReply{}, c.impl.Transmit(ctx, reportCtx, request.Report, sigs)
}

func (c *contractTransmitterServer) LatestConfigDigestAndEpoch(ctx context.Context, request *pb.LatestConfigDigestAndEpochRequest) (*pb.LatestConfigDigestAndEpochReply, error) {
	digest, epoch, err := c.impl.LatestConfigDigestAndEpoch(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.LatestConfigDigestAndEpochReply{ConfigDigest: digest[:], Epoch: epoch}, nil
}

func (c *contractTransmitterServer) FromAccount(ctx context.Context, request *pb.FromAccountRequest) (*pb.FromAccountReply, error) {
	a, err := c.impl.FromAccount()
	if err != nil {
		return nil, err
	}
	return &pb.FromAccountReply{Account: string(a)}, nil
}
