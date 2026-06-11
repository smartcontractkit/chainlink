package mockchip

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

// EventSummary is the JSON shape returned by the events listing endpoint. The
// CloudEvent payload bytes are returned base64-encoded (Go's default JSON
// encoding for []byte) so the endpoint is also useful from `curl | jq`.
type EventSummary struct {
	ReceivedAt time.Time         `json:"received_at"`
	ID         string            `json:"id"`
	Source     string            `json:"source"`
	Type       string            `json:"type"`
	Subject    string            `json:"subject,omitempty"`
	Time       string            `json:"time,omitempty"`
	Data       []byte            `json:"data,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// HTTPController exposes a small inspection / control API around a Server:
//
//	GET  /events       — JSON list of captured events
//	GET  /events/count — plaintext count
//	GET  /stats        — JSON stats blob
//	POST /outage/on    — start dropping every gRPC publish call
//	POST /outage/off   — restore normal behaviour
//	POST /reset        — clear captured events and counters
//	GET  /healthz      — 200 OK
type HTTPController struct {
	server   *Server
	httpSrv  *http.Server
	listener net.Listener
}

// NewHTTPController wires up routes for s.
func NewHTTPController(s *Server) *HTTPController {
	c := &HTTPController{server: s}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", c.handleHealth)
	mux.HandleFunc("/events", c.handleEvents)
	mux.HandleFunc("/events/count", c.handleCount)
	mux.HandleFunc("/stats", c.handleStats)
	mux.HandleFunc("/outage/on", c.handleOutageOn)
	mux.HandleFunc("/outage/off", c.handleOutageOff)
	mux.HandleFunc("/reset", c.handleReset)
	c.httpSrv = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return c
}

// Start binds the HTTP listener on listenAddr and serves in a goroutine.
func (c *HTTPController) Start(listenAddr string) (string, error) {
	lc := &net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", listenAddr)
	if err != nil {
		return "", fmt.Errorf("mockchip: http listen %s: %w", listenAddr, err)
	}
	c.listener = lis
	go func() {
		if err := c.httpSrv.Serve(lis); err != nil && err != http.ErrServerClosed {
			fmt.Printf("mockchip: HTTP Serve returned: %v\n", err)
		}
	}()
	return lis.Addr().String(), nil
}

// Stop gracefully shuts down the HTTP server.
func (c *HTTPController) Stop(ctx context.Context) error {
	return c.httpSrv.Shutdown(ctx)
}

func (c *HTTPController) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (c *HTTPController) handleEvents(w http.ResponseWriter, _ *http.Request) {
	events := c.server.Captured()
	out := make([]EventSummary, 0, len(events))
	for _, ce := range events {
		out = append(out, summarize(ce))
	}
	writeJSON(w, http.StatusOK, out)
}

func (c *HTTPController) handleCount(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "%d\n", c.server.CapturedCount())
}

func (c *HTTPController) handleStats(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, c.server.Stats())
}

func (c *HTTPController) handleOutageOn(w http.ResponseWriter, r *http.Request) {
	if !mustBePOST(w, r) {
		return
	}
	c.server.SetOutage(true)
	writeJSON(w, http.StatusOK, map[string]any{"outage_active": true})
}

func (c *HTTPController) handleOutageOff(w http.ResponseWriter, r *http.Request) {
	if !mustBePOST(w, r) {
		return
	}
	c.server.SetOutage(false)
	writeJSON(w, http.StatusOK, map[string]any{"outage_active": false})
}

func (c *HTTPController) handleReset(w http.ResponseWriter, r *http.Request) {
	if !mustBePOST(w, r) {
		return
	}
	c.server.Reset()
	writeJSON(w, http.StatusOK, map[string]any{"reset": true})
}

func mustBePOST(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func summarize(ce CapturedEvent) EventSummary {
	ev := ce.Event
	s := EventSummary{ReceivedAt: ce.ReceivedAt}
	if ev == nil {
		return s
	}
	s.ID = ev.GetId()
	s.Source = ev.GetSource()
	s.Type = ev.GetType()
	if attrs := ev.GetAttributes(); len(attrs) > 0 {
		s.Attributes = make(map[string]string, len(attrs))
		for k, v := range attrs {
			str := attributeToString(v)
			s.Attributes[k] = str
			switch k {
			case "subject":
				s.Subject = str
			case "time":
				s.Time = str
			}
		}
	}
	if raw := ev.GetBinaryData(); raw != nil {
		s.Data = raw
	}
	return s
}

// attributeToString renders a CloudEvent attribute value as a human-readable
// string for the inspection endpoint.
func attributeToString(v *cepb.CloudEventAttributeValue) string {
	if v == nil {
		return ""
	}
	switch a := v.GetAttr().(type) {
	case *cepb.CloudEventAttributeValue_CeBoolean:
		return strconv.FormatBool(a.CeBoolean)
	case *cepb.CloudEventAttributeValue_CeInteger:
		return strconv.FormatInt(int64(a.CeInteger), 10)
	case *cepb.CloudEventAttributeValue_CeString:
		return a.CeString
	case *cepb.CloudEventAttributeValue_CeBytes:
		return base64.StdEncoding.EncodeToString(a.CeBytes)
	case *cepb.CloudEventAttributeValue_CeUri:
		return a.CeUri
	case *cepb.CloudEventAttributeValue_CeUriRef:
		return a.CeUriRef
	case *cepb.CloudEventAttributeValue_CeTimestamp:
		if ts := a.CeTimestamp; ts != nil {
			return ts.AsTime().UTC().Format(time.RFC3339Nano)
		}
	}
	return ""
}
