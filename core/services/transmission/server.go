package transmission

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/transmission/wsrpc/proto"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/wsrpc"
	"golang.org/x/sync/errgroup"
)

type Config interface {
	RPCPort() int
}

type server struct {
	utils.StartStopOnce
	serverConn *wsrpc.Server
	config     Config
	lggr       logger.Logger
}

func (s *server) Start(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	err := s.StartOnce("WSRPC Server", func() error {
		addr := fmt.Sprintf(":%d", s.config.RPCPort())
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return errors.Wrap(err, "failed to listen to rpc server")
		}

		s.serverConn = wsrpc.NewServer()

		go s.serverConn.Serve(lis)
		proto.RegisterTransmissionServer(s.serverConn, NewHandler(s.lggr))

		g.Go(func() error {
			<-gCtx.Done()

			s.serverConn.Stop()

			return nil
		})

		s.lggr.Infof("Listening and serving RPC on port %d", s.config.RPCPort())

		return nil
	})

	if err != nil {
		return err
	}

	return errors.WithStack(g.Wait())
}
