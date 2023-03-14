package transmission

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/transmission/wsrpc/proto"
)

type handler struct {
	lggr logger.Logger
}

type HandlerInterface interface {
	SendUserOperation(ctx context.Context, req *proto.SendUserOperationRequest) (*proto.SendUserOperationResponse, error)
}

func NewHandler(lggr logger.Logger) HandlerInterface {
	return &handler{
		lggr: lggr,
	}
}

func (h *handler) SendUserOperation(ctx context.Context, req *proto.SendUserOperationRequest) (*proto.SendUserOperationResponse, error) {
	return nil, nil
}
