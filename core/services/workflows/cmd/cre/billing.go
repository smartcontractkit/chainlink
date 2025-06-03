package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

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

	billing.UnimplementedWorkflowServiceServer
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

	return &billing.ReserveCreditsResponse{Success: true}, nil
}

func (s *BillingService) WorkflowReceipt(
	_ context.Context,
	request *billing.SubmitWorkflowReceiptRequest,
) (*billing.SubmitWorkflowReceiptResponse, error) {
	s.lggr.Infof("WorkflowReceipt: %v", request.Metering)

	return &billing.SubmitWorkflowReceiptResponse{Success: true}, nil
}

func (s *BillingService) start(ctx context.Context) error {
	lis, err := net.Listen("tcp", "localhost:4319")
	if err != nil {
		log.Fatalf("billing failed to listen: %v", err)
		return err
	}

	server := grpc.NewServer()

	billing.RegisterWorkflowServiceServer(server, &BillingService{lggr: s.lggr})

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("billing failed to serve: %v", err)
		return err
	}

	s.server = server

	return nil
}

func (s *BillingService) close() error {
	s.server.Stop()
	return nil
}

func setupBeholder(lggr logger.Logger) error {
	writer := &lggrWriter{lggr: lggr}

	client, err := beholder.NewWriterClient(writer)
	if err != nil {
		return err
	}

	beholder.SetClient(client)

	return nil
}

type lggrWriter struct {
	lggr logger.Logger
}

func (w lggrWriter) Write(bts []byte) (int, error) {
	w.lggr.Info(string(bts))

	return len(bts), nil
}
