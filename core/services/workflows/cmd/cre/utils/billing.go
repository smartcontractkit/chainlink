package utils

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
)

const defaultBillingPort = 4319

type BillingService struct {
	services.Service
	eng *services.Engine

	port   int
	lggr   logger.Logger
	server *grpc.Server
	addr   string

	billing.UnimplementedCreditReservationServiceServer
}

// Addr returns the address the billing service is listening on. Only valid after Start.
func (s *BillingService) Addr() string { return s.addr }

var _ services.Service = (*BillingService)(nil)

// NewBillingServiceOnPort specify a port to listen on for the billing service.
// Use 0 for a random port.
func NewBillingServiceOnPort(lggr logger.Logger, port int) *BillingService {
	b := &BillingService{
		lggr: lggr,
		port: port,
	}
	b.Service, b.eng = services.Config{
		Name:  "fakeBillingService",
		Start: b.start,
		Close: b.close,
	}.NewServiceEngine(lggr)
	return b
}

// NewBillingService listens on the default billing port.
func NewBillingService(lggr logger.Logger) *BillingService {
	b := &BillingService{
		lggr: lggr,
		port: defaultBillingPort,
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
	lc := &net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", fmt.Sprintf("localhost:%d", s.port))
	if err != nil {
		return err
	}
	s.addr = lis.Addr().String()

	server := grpc.NewServer()

	billing.RegisterCreditReservationServiceServer(server, s)

	go func() {
		err := server.Serve(lis)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.lggr.Errorf("billing gRPC serve ended unexpectedly: %v", err)
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
