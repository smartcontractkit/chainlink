package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
)

// Config defines how the test sink listens.
type Config struct {
	// gRPC listen address for ChipIngress, e.g. ":9090".
	GRPCListen string

	// HTTP listen address for test API, e.g. ":8080".
	HTTPListen string

	// Optional upstream Chip Ingress endpoint to forward to.
	// If empty, no pass-through is performed.
	UpstreamEndpoint string

	// Maximum number of events to cache.
	CacheSize int
}

// Server implements the ChipIngress gRPC service + a tiny HTTP API.
type Server struct {
	cfg   Config
	store *Store

	grpcServer *grpc.Server
	httpServer *http.Server

	chippb.UnimplementedChipIngressServer

	// Optional pass-through client.
	upstream chippb.ChipIngressClient
}

// NewServer constructs a new test sink.
func NewServer(cfg Config) (*Server, error) {
	store, err := NewStore(cfg.CacheSize)
	if err != nil {
		return nil, fmt.Errorf("create store: %w", err)
	}

	s := &Server{
		cfg:   cfg,
		store: store,
	}

	// gRPC server
	s.grpcServer = grpc.NewServer()

	// Register ChipIngress service implementation on this server.
	chippb.RegisterChipIngressServer(s.grpcServer, s)

	// Optional upstream (pass-through)
	if cfg.UpstreamEndpoint != "" {
		conn, err := grpc.NewClient(cfg.UpstreamEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("dial upstream chip ingress: %w", err)
		}
		s.upstream = chippb.NewChipIngressClient(conn)
	}

	// HTTP API
	mux := http.NewServeMux()
	mux.HandleFunc("/reset", s.handleReset)
	mux.HandleFunc("/events", s.handleEvents)
	mux.HandleFunc("/sequence", s.handleSequence)

	s.httpServer = &http.Server{
		Addr:    cfg.HTTPListen,
		Handler: mux,
	}

	return s, nil
}

// Run starts both gRPC and HTTP servers and blocks until one of them exits.
func (s *Server) Run() error {
	errCh := make(chan error, 2)

	// gRPC
	go func() {
		lis, err := net.Listen("tcp", s.cfg.GRPCListen)
		if err != nil {
			errCh <- fmt.Errorf("gRPC listen: %w", err)
			return
		}
		log.Printf("[chip-testsink] gRPC listening on %s", s.cfg.GRPCListen)
		errCh <- s.grpcServer.Serve(lis)
	}()

	// HTTP
	go func() {
		log.Printf("[chip-testsink] HTTP API listening on %s", s.cfg.HTTPListen)
		errCh <- s.httpServer.ListenAndServe()
	}()

	if s.cfg.UpstreamEndpoint != "" {
		log.Printf("[chip-testsink] Forwarding to upstream Chip Ingress endpoint: %s", s.cfg.UpstreamEndpoint)
	}

	// Wait for first error.
	return <-errCh
}

// Shutdown gracefully stops both servers.
func (s *Server) Shutdown(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return s.httpServer.Shutdown(ctx)
}

//
// ===== gRPC: ChipIngressServer implementation =====
//

// Publish implements chippb.ChipIngressServer.Publish.
//
// Adjust the signature if your generated interface differs.
func (s *Server) Publish(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
	if event == nil {
		return &chippb.PublishResponse{}, nil
	}

	attrs := make(map[string]any, len(event.Attributes))
	for k, v := range event.Attributes {
		attrs[k] = v
	}

	fmt.Printf("Received event with type %s\n", event.Type)

	ce := CapturedEvent{
		Domain: event.Source,
		Entity: event.Type,
		Body:   event.GetProtoData().Value, // raw proto bytes
		Attrs:  attrs,
	}

	if s, ok := attrs["beholder_domain"].(string); ok {
		ce.Domain = s
	}
	if s, ok := attrs["beholder_entity"].(string); ok {
		ce.Entity = s
	}
	if s, ok := attrs["beholder_data_schema"].(string); ok {
		ce.Schema = s
	}

	s.store.Add(ce)

	// 2) Pass-through
	if s.upstream != nil {
		// Forward the exact same batch upstream.
		return s.upstream.Publish(ctx, event)
	}

	// 3) No upstream: just pretend success.
	return &chippb.PublishResponse{}, nil
}

//
// ===== HTTP handlers =====
//

// POST /reset
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	s.store.Reset()
	w.WriteHeader(http.StatusNoContent)
}

// GET /events[?entity=Foo.Bar][&domain=my.domain][&orphans=1][&sequence=123][&limit=100]
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	entityFilter := r.URL.Query().Get("entity")
	domainFilter := r.URL.Query().Get("domain")
	orphansOnly := r.URL.Query().Get("orphans") == "1"

	var sequenceFilter uint64
	if seqStr := r.URL.Query().Get("sequence"); seqStr != "" {
		if _, err := fmt.Sscanf(seqStr, "%d", &sequenceFilter); err != nil {
			http.Error(w, "invalid sequence parameter", http.StatusBadRequest)
			return
		}
	}

	var limitFilter int
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if _, err := fmt.Sscanf(limitStr, "%d", &limitFilter); err != nil {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	type outEvent struct {
		Sequence uint64         `json:"sequence"`
		Domain   string         `json:"domain"`
		Entity   string         `json:"entity"`
		Schema   string         `json:"schema"`
		Body     string         `json:"body"` // base64
		BodyRaw  string         `json:"bodyRaw"`
		Attrs    map[string]any `json:"attrs"`
	}

	all := s.store.All()
	out := make([]outEvent, 0, len(all))

	for _, ev := range all {
		// Filter by sequence (only events with sequence > sequenceFilter)
		if sequenceFilter > 0 && ev.Sequence < sequenceFilter {
			continue
		}

		if entityFilter != "" && ev.Entity != entityFilter {
			continue
		}
		if domainFilter != "" && ev.Domain != domainFilter {
			continue
		}
		if orphansOnly && ev.Entity != "" {
			continue
		}

		out = append(out, outEvent{
			Sequence: ev.Sequence,
			Domain:   ev.Domain,
			Entity:   ev.Entity,
			Schema:   ev.Schema,
			Body:     base64.StdEncoding.EncodeToString(ev.Body),
			BodyRaw:  string(ev.Body),
			Attrs:    ev.Attrs,
		})

		// Apply limit if specified
		if limitFilter > 0 && len(out) >= limitFilter {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// GET /sequence - returns the current (latest) sequence number
func (s *Server) handleSequence(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}

	currentSeq := s.store.CurrentSequence()

	response := map[string]uint64{
		"sequence": currentSeq,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
