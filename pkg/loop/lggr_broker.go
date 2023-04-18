package loop

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

type lggrBroker struct {
	lggr   logger.Logger
	broker *plugin.GRPCBroker
}

func (lb *lggrBroker) named(name string) *lggrBroker {
	return &lggrBroker{
		lggr:   logger.Named(lb.lggr, name),
		broker: lb.broker,
	}
}

func (lb *lggrBroker) serve(server *grpc.Server, name string, deps ...resource) (uint32, error) {
	id := lb.broker.NextId()
	lb.lggr.Debugf("Serving %s on connection %d", name, id)
	lis, err := lb.broker.Accept(id)
	if err != nil {
		lb.closeAll(deps...)
		return 0, ErrConnAccept{Name: name, ID: id, Err: err}
	}
	go func() {
		defer lb.closeAll(deps...)
		if err := server.Serve(lis); err != nil {
			lb.lggr.Errorw(fmt.Sprintf("Failed to serve %s on connection %d", name, id), "err", err)
		}
	}()
	return id, nil
}

func (lb *lggrBroker) closeAll(deps ...resource) {
	for _, d := range deps {
		if err := d.Close(); err != nil {
			lb.lggr.Error(fmt.Sprintf("Error closing %s", d.name), "err", err)
		}
	}
}

type resource struct {
	io.Closer
	name string
}
