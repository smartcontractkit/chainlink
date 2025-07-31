package utils

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
)

type BillingService struct {
	services.Service
	eng *services.Engine

	lggr   logger.Logger
	server *grpc.Server

	billing.UnimplementedCreditReservationServiceServer
}

var _ services.Service = (*BillingService)(nil)

func NewBillingService(lggr logger.Logger) *BillingService {
	b := &BillingService{
		lggr: lggr,
	}
	b.Service, b.eng = services.Config{
		Name:  "fakeBillingService",
		Start: b.start,
		Close: b.close,
	}.NewServiceEngine(lggr)
	return b
}

func (s *BillingService) ReserveCredits(
	_ context.Context,
	request *billing.ReserveCreditsRequest,
) (*billing.ReserveCreditsResponse, error) {
	s.lggr.Infof("ReserveCredits: %v", request)

	return &billing.ReserveCreditsResponse{
		Success: true,
		RateCards: []*billing.RateCard{
			{
				ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
				MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
				UnitsPerCredit:  "0.0001",
			},
		},
		Credits: "10000",
	}, nil
}

func (s *BillingService) SubmitWorkflowReceipt(
	_ context.Context,
	request *billing.SubmitWorkflowReceiptRequest,
) (*emptypb.Empty, error) {
	s.lggr.Infof("WorkflowReceipt: %v", request.Metering)

	return &emptypb.Empty{}, nil
}

func (s *BillingService) start(ctx context.Context) error {
	lis, err := net.Listen("tcp", "localhost:4319")
	if err != nil {
		log.Fatalf("billing failed to listen: %v", err)
		return err
	}

	server := grpc.NewServer()

	billing.RegisterCreditReservationServiceServer(server, &BillingService{lggr: s.lggr})

	go func() {
		err = server.Serve(lis)
		if err != nil {
			log.Fatalf("billing failed to serve: %v", err)
			return
		}
	}()

	s.server = server

	return nil
}

func (s *BillingService) close() error {
	if s.server != nil {
		s.server.Stop()
	}
	return nil
}

func SetupBeholder(lggr logger.Logger) error {
	writer := &lggrWriter{lggr: lggr}

	client, err := beholder.NewWriterClient(writer)
	if err != nil {
		return err
	}

	beholder.SetClient(client)

	return nil
}

func cleanupBeholder() error {
	client := beholder.GetClient()
	if client != nil {
		return client.Close()
	}

	return nil
}

type lggrWriter struct {
	lggr logger.Logger
}

func (w lggrWriter) Write(bts []byte) (int, error) {
	w.lggr.Info(string(bts))

	return len(bts), nil
}
