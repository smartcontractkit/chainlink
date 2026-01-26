package testsinkstateful

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
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

	// Started optionally receives a signal once both listeners are bound.
	Started chan<- struct{}
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
	mux.HandleFunc("/events", s.handleEvents)

	s.httpServer = &http.Server{
		Addr:    cfg.HTTPListen,
		Handler: mux,
	}

	return s, nil
}

// Run starts both gRPC and HTTP servers and blocks until one of them exits.
func (s *Server) Run() error {
	grpcLis, err := net.Listen("tcp", s.cfg.GRPCListen)
	if err != nil {
		return fmt.Errorf("gRPC listen: %w", err)
	}

	httpLis, err := net.Listen("tcp", s.cfg.HTTPListen)
	if err != nil {
		grpcLis.Close()
		return fmt.Errorf("HTTP listen: %w", err)
	}

	notifyStarted(s.cfg.Started)
	log.Printf("[chip-testsink] gRPC listening on %s", s.cfg.GRPCListen)
	log.Printf("[chip-testsink] HTTP API listening on %s", s.cfg.HTTPListen)

	errCh := make(chan error, 2)

	// gRPC
	go func() {
		errCh <- s.grpcServer.Serve(grpcLis)
	}()

	// HTTP
	go func() {
		errCh <- s.httpServer.Serve(httpLis)
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

func notifyStarted(ch chan<- struct{}) {
	if ch == nil {
		return
	}

	select {
	case ch <- struct{}{}:
	default:
	}
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

type ServedEvent struct {
	Timestamp int            `json:"timestamp"`
	Domain    string         `json:"domain"`
	Entity    string         `json:"entity"`
	Schema    string         `json:"schema"`
	Body      string         `json:"body"` // base64
	BodyRaw   string         `json:"bodyRaw"`
	Attrs     map[string]any `json:"attrs"`
}

// GET /events[?entity=Foo.Bar][&timestamp=1769427180][&limit=100]
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	entityFilter := r.URL.Query().Get("entity")

	var timestampFilter time.Time
	if tsStr := r.URL.Query().Get("timestamp"); tsStr != "" {
		var tsInt int64
		if _, err := fmt.Sscanf(tsStr, "%d", &tsInt); err != nil {
			http.Error(w, "invalid timestamp parameter, expected Unix timestamp", http.StatusBadRequest)
			return
		}
		timestampFilter = time.Unix(tsInt, 0)
	}

	var limitFilter int
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if _, err := fmt.Sscanf(limitStr, "%d", &limitFilter); err != nil {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	all := s.store.All()
	out := make([]ServedEvent, 0, len(all))

	for _, ev := range all {
		// Filter by timestamp (only events with timestamp > timestampFilter)
		if !timestampFilter.IsZero() && ev.Timestamp.Before(timestampFilter) {
			continue
		}

		if entityFilter != "" && ev.Entity != entityFilter {
			continue
		}

		out = append(out, ServedEvent{
			Timestamp: int(ev.Timestamp.Unix()),
			Domain:    ev.Domain,
			Entity:    ev.Entity,
			Schema:    ev.Schema,
			Body:      base64.StdEncoding.EncodeToString(ev.Body),
			BodyRaw:   string(ev.Body),
			Attrs:     ev.Attrs,
		})

		// Apply limit if specified
		if limitFilter > 0 && len(out) >= limitFilter {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(out)

	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
