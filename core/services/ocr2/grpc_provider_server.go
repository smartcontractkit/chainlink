package ocr2

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ProviderServer struct {
	s      *grpc.Server
	lis    net.Listener
	lggr   logger.Logger
	conns  []*grpc.ClientConn
	doneCh chan any
}

// checkConnectionsLoop will check if we have active connections,
// if no active connections are found the server will be closed
func (p *ProviderServer) checkConnectionsLoop() {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-p.doneCh:
				return
			case <-ticker.C:
				hasActiveClients := false
				for _, c := range p.conns {
					if c.GetState() < 3 {
						hasActiveClients = true
						break
					}
				}
				if !hasActiveClients {
					ticker.Stop()
					p.Close()
				}
			}
		}
	}()
}

// Start is a NOOP as the server is started at creation
func (p *ProviderServer) Start(ctx context.Context) error {
	return nil
}

func (p *ProviderServer) Close() error {
	p.doneCh <- true
	p.s.Stop()
	return nil
}

func (p *ProviderServer) GetConn() (*grpc.ClientConn, error) {
	cc, err := grpc.Dial(p.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	p.conns = append(p.conns, cc)
	return cc, err
}

// NewProviderServer creates a GRPC server that will wrap a provider, this is a workaround to test the Node API PoC until the EVM relayer is loopifyed
func NewProviderServer(p types.PluginProvider, lggr logger.Logger) (*ProviderServer, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	ps := ProviderServer{
		s:      grpc.NewServer(),
		lis:    lis,
		lggr:   lggr.Named("EVM.ProviderServer"),
		doneCh: make(chan any),
	}
	loop.RegisterStandAloneProvider(ps.s, p)

	ps.serve()
	ps.checkConnectionsLoop()
	return &ps, nil
}

func (p *ProviderServer) serve() {
	go func() {
		if err := p.s.Serve(p.lis); err != nil {
			p.lggr.Errorf("Failed to serve EVM provider server: %v", err)
		}
	}()
}
