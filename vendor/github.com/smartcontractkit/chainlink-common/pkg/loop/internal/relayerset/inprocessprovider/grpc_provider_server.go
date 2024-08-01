package inprocessprovider

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// ProviderServer is a GRPC server that will wrap a provider, this is a workaround, copied from core, until all
// providers are in their own loops.
type ProviderServer struct {
	s     *grpc.Server
	lis   net.Listener
	lggr  logger.Logger
	conns []*grpc.ClientConn

	provider types.PluginProvider
}

func (p *ProviderServer) Start(ctx context.Context) error {
	p.serve()
	return nil
}

func (p *ProviderServer) Close() error {
	var err error
	for _, c := range p.conns {
		err = errors.Join(err, c.Close())
	}
	p.s.Stop()
	return err
}

func (p *ProviderServer) GetConn() (grpc.ClientConnInterface, error) {
	//TODO https://smartcontract-it.atlassian.net/browse/BCF-3290
	cc, err := grpc.Dial(p.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:staticcheck
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
		lggr: lggr,
	}
	err = RegisterStandAloneProvider(ps.s, p, pType)
	if err != nil {
		return nil, err
	}

	return &ps, nil
}

func (p *ProviderServer) serve() {
	go func() {
		if err := p.s.Serve(p.lis); err != nil {
			p.lggr.Error("Failed to serve in process provider server", "providerName", p.provider.Name(), "error", err)
		}
	}()
}
