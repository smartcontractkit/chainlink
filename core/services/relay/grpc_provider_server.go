package relay

import (
	"context"
	"net"

	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ProviderServer struct {
	s     *grpc.Server
	lis   net.Listener
	lggr  logger.Logger
	conns []*grpc.ClientConn
}

func (p *ProviderServer) Start(ctx context.Context) error {
	p.serve()
	return nil
}

func (p *ProviderServer) Close() error {
	var err error
	for _, c := range p.conns {
		err = multierr.Combine(err, c.Close())
	}
	p.s.Stop()
	return err
}

func (p *ProviderServer) GetConn() (*grpc.ClientConn, error) {
	cc, err := grpc.Dial(p.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	p.conns = append(p.conns, cc)
	return cc, err
}

// NewProviderServer creates a GRPC server that will wrap a provider, this is a workaround to test the Node API PoC until the EVM relayer is loopifyed
func NewProviderServer(p types.PluginProvider, pType types.OCR2PluginType, lggr logger.Logger) (*ProviderServer, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	ps := ProviderServer{
		s:    grpc.NewServer(),
		lis:  lis,
		lggr: lggr.Named("EVM.ProviderServer"),
	}
	err = loop.RegisterStandAloneProvider(ps.s, p, pType)
	if err != nil {
		return nil, err
	}

	return &ps, nil
}

func (p *ProviderServer) serve() {
	go func() {
		if err := p.s.Serve(p.lis); err != nil {
			p.lggr.Errorf("Failed to serve EVM provider server: %v", err)
		}
	}()
}
