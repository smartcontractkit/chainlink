package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	httptypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http"
)

type JSONRPCRequest struct {
	Input json.RawMessage `json:"input"`
}

type Config struct {
	Port uint16
}

type LocalGateway struct {
	config Config
}

func NewLocalGateway(config Config) *LocalGateway {
	return &LocalGateway{config: config}
}

func (g *LocalGateway) ListenForTriggerPayload(ctx context.Context) (*httptypedapi.Payload, error) {
	payloadCh := make(chan *httptypedapi.Payload, 1)
	errorCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) {
		input, err := parseRequest(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("error processing request: %v", err), http.StatusBadRequest)
			return
		}

		payloadCh <- &httptypedapi.Payload{
			Input: input,
		}
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", g.config.Port),
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}
	defer server.Close()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorCh <- err
		}
	}()

	select {
	case payload := <-payloadCh:
		return payload, nil
	case err := <-errorCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func parseRequest(req *http.Request) ([]byte, error) {
	if req.Method != http.MethodPost {
		return nil, errors.New("gateway expects POST request")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var rpcRequest JSONRPCRequest
	if err := json.Unmarshal(body, &rpcRequest); err != nil {
		return nil, fmt.Errorf("failed to parse request body: %w", err)
	}

	return rpcRequest.Input, nil
}
