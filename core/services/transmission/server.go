package transmission

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/transmission/handler"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Server struct {
	utils.StartStopOnce
	lggr       logger.Logger
	httpServer http.Server
	rpcPort    uint16
}

func NewServer(handler handler.HttpHandler, rpcPort uint16, lggr logger.Logger) *Server {
	r := mux.NewRouter()

	r.HandleFunc("/get_nonce", handler.GetNonce).Methods("POST")

	addr := fmt.Sprintf(":%d", rpcPort)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{
		httpServer: *httpServer,
		lggr:       lggr,
		rpcPort:    rpcPort,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.StartOnce("TransmissionServer", func() error {
		go func() {
			err := s.httpServer.ListenAndServe()
			if !errors.Is(err, http.ErrServerClosed) {
				s.lggr.Warnf("server closed with err: %v", err)
			}
			s.lggr.Info("Stopped serving new connections.")
		}()

		shutdownServer := func() {
			<-ctx.Done()
			err := s.httpServer.Close()
			if err != nil {
				s.lggr.Errorf("error while closing server: %v", err)
			}
		}
		go shutdownServer()

		s.lggr.Infof("Listening and serving RPC on port %d", s.rpcPort)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() error {
	return s.StopOnce("TransmissionServer", func() error {
		err := s.httpServer.Close()
		if err != nil {
			return err
		}
		s.lggr.Infof("closed http server")
		return nil
	})
}
