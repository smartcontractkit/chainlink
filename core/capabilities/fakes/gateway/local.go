package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	httptypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http"
)

const (
	triggerPath       = "/trigger"
	maxRequestBytes   = 1 << 20 // 1 MiB
	readHeaderTimeout = time.Second
	shutdownTimeout   = time.Second
)

type triggerRequest struct {
	Input json.RawMessage `json:"input"`
}

// Config holds settings for a LocalGateway test server.
type Config struct {
	Port uint16
}

// LocalGateway is a minimal HTTP server that accepts a single trigger POST
// and returns the parsed payload to the caller.
type LocalGateway struct {
	config Config
}

// NewLocalGateway returns a LocalGateway bound to the port in config.
func NewLocalGateway(config Config) *LocalGateway {
	return &LocalGateway{config: config}
}

// ListenForTriggerPayload starts an HTTP server on the configured port and
// blocks until a POST /trigger request arrives or ctx is cancelled.
func (g *LocalGateway) ListenForTriggerPayload(ctx context.Context) (*httptypedapi.Payload, error) {
	type result struct {
		payload *httptypedapi.Payload
		err     error
	}
	resultCh := make(chan result, 1)

	mux := http.NewServeMux()
	mux.HandleFunc(triggerPath, func(w http.ResponseWriter, r *http.Request) {
		input, err := parseRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		select {
		case resultCh <- result{payload: &httptypedapi.Payload{Input: input}}:
			w.WriteHeader(http.StatusOK)
		case <-r.Context().Done():
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		default:
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("", strconv.Itoa(int(g.config.Port))),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		BaseContext:       func(net.Listener) context.Context { return ctx },
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case resultCh <- result{err: err}:
			default:
			}
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	select {
	case res := <-resultCh:
		return res.payload, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func parseRequest(req *http.Request) ([]byte, error) {
	if req.Method != http.MethodPost {
		return nil, errors.New("gateway expects POST request")
	}
	defer req.Body.Close()

	body, err := io.ReadAll(http.MaxBytesReader(nil, req.Body, maxRequestBytes))
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	var request triggerRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, fmt.Errorf("parse request body: %w", err)
	}

	return request.Input, nil
}
