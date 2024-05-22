package monitoring

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// HTTPServer is the HTTP interface exposed by every monitoring.
// It's used to export metrics to prometheus, to query the node for configurations, etc.
type HTTPServer interface {
	Handle(path string, handler http.Handler)
	Run(ctx context.Context)
}

func NewHTTPServer(stopCh services.StopRChan, addr string, log Logger) HTTPServer {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(_ net.Listener) (ctx context.Context) {
			ctx, _ = stopCh.NewCtx()
			return
		},
		ReadHeaderTimeout: 60 * time.Second,
	}
	return &httpServer{srv, mux, log, addr}
}

type httpServer struct {
	server *http.Server
	mux    *http.ServeMux
	log    Logger
	addr   string
}

func (h *httpServer) Handle(path string, handler http.Handler) {
	h.mux.Handle(path, handler)
}

// Run should be executed as a goroutine
func (h *httpServer) Run(ctx context.Context) {
	var subs utils.Subprocesses
	defer subs.Wait()
	subs.Go(func() {
		h.log.Debugw("starting HTTP server")
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Fatalw("failed to start HTTP server", "address", h.addr, "error", err)
		} else {
			h.log.Infow("HTTP server stopped")
		}
	})
	subs.Go(func() {
		<-ctx.Done()
		if err := h.server.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
			h.log.Errorw("failed to shut HTTP server down", "error", err)
		}
	})
}
