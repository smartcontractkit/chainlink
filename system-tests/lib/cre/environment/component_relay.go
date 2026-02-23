package environment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

type componentRelayManager struct {
	lggr    zerolog.Logger
	baseURL string

	mu      sync.Mutex
	handles map[string]*componentRelayHandle
}

type componentRelayHandle struct {
	relayID string
	cancel  context.CancelFunc
}

type relayOpenResponse struct {
	RelayID string `json:"relayId"`
}

func newComponentRelayManager(lggr zerolog.Logger) (*componentRelayManager, error) {
	baseURL, err := resolveAgentBaseURLForRelay()
	if err != nil {
		return nil, err
	}
	return &componentRelayManager{
		lggr:    lggr,
		baseURL: baseURL,
		handles: make(map[string]*componentRelayHandle),
	}, nil
}

func (m *componentRelayManager) EnsurePort(ctx context.Context, relayName string, localPort int) error {
	if m == nil || localPort <= 0 {
		return nil
	}
	// Deduplicate by port. HTTP and WS for the same endpoint can share one listener.
	key := strconv.Itoa(localPort)

	m.mu.Lock()
	if _, ok := m.handles[key]; ok {
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	relayID, err := openRelay(ctx, m.baseURL, relayName, localPort)
	if err != nil {
		return err
	}

	workerCtx, cancel := context.WithCancel(context.Background())
	localAddr := net.JoinHostPort("127.0.0.1", strconv.Itoa(localPort))
	for i := 0; i < 4; i++ {
		go relayWorker(workerCtx, m.baseURL, relayID, localAddr)
	}

	m.mu.Lock()
	m.handles[key] = &componentRelayHandle{relayID: relayID, cancel: cancel}
	m.mu.Unlock()
	m.lggr.Info().Str("relayName", relayName).Int("port", localPort).Msg("ensured mixed component relay")
	return nil
}

func (m *componentRelayManager) Close(ctx context.Context) error {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	handles := make([]*componentRelayHandle, 0, len(m.handles))
	for _, h := range m.handles {
		handles = append(handles, h)
	}
	m.handles = map[string]*componentRelayHandle{}
	m.mu.Unlock()

	var firstErr error
	for _, h := range handles {
		h.cancel()
		if err := closeRelay(ctx, m.baseURL, h.relayID); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func resolveAgentBaseURLForRelay() (string, error) {
	if v := strings.TrimSpace(os.Getenv(envEC2AgentURL)); v != "" {
		return v, nil
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv(envAgentMode)), "ec2") && runtimecfg.IsDirectMode() {
		hostIP, err := runtimecfg.DirectHostIP()
		if err != nil {
			return "", err
		}
		port, err := resolveEC2AgentPort()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("http://%s:%d", hostIP, port), nil
	}
	if v := strings.TrimSpace(os.Getenv(envLocalAgentURL)); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("cannot resolve agent base URL for relay; set %s (or use direct mode with EC2 host resolution)", envEC2AgentURL)
}

func openRelay(ctx context.Context, baseURL, name string, requestedPort int) (string, error) {
	body, _ := json.Marshal(map[string]any{"name": name, "requestedPort": requestedPort})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/v1/relay/open", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("open relay failed: status %s body %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	var out relayOpenResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", err
	}
	if strings.TrimSpace(out.RelayID) == "" {
		return "", fmt.Errorf("open relay returned empty relayId")
	}
	return out.RelayID, nil
}

func closeRelay(ctx context.Context, baseURL, relayID string) error {
	body, _ := json.Marshal(map[string]any{"relayId": relayID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/v1/relay/close", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("close relay failed: status %s body %s", resp.Status, strings.TrimSpace(string(respBody)))
	}
	return nil
}

func relayWorker(ctx context.Context, baseURL, relayID, localAddr string) {
	backoff := 250 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		wsURL, err := relayConnectWSURL(baseURL, relayID)
		if err != nil {
			time.Sleep(backoff)
			continue
		}
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			time.Sleep(backoff)
			continue
		}
		_ = bridgeRelayStream(ctx, ws, localAddr)
		_ = ws.Close()
		if backoff < 2*time.Second {
			backoff *= 2
		}
	}
}

func relayConnectWSURL(baseURL, relayID string) (string, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported agent url scheme: %s", u.Scheme)
	}
	u.Path = "/v1/relay/connect"
	q := u.Query()
	q.Set("relayId", relayID)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func bridgeRelayStream(ctx context.Context, ws *websocket.Conn, localAddr string) error {
	errCh := make(chan error, 2)
	localReady := make(chan net.Conn, 1)
	var localConn net.Conn
	var localConnMu sync.Mutex
	getLocalConn := func() net.Conn {
		localConnMu.Lock()
		defer localConnMu.Unlock()
		return localConn
	}
	setLocalConn := func(conn net.Conn) {
		localConnMu.Lock()
		localConn = conn
		localConnMu.Unlock()
	}
	ensureLocalConn := func() (net.Conn, error) {
		if existing := getLocalConn(); existing != nil {
			return existing, nil
		}
		conn, err := net.DialTimeout("tcp", localAddr, 2*time.Second)
		if err != nil {
			return nil, err
		}
		setLocalConn(conn)
		select {
		case localReady <- conn:
		default:
		}
		return conn, nil
	}
	defer func() {
		if conn := getLocalConn(); conn != nil {
			_ = conn.Close()
		}
	}()
	go func() {
		var conn net.Conn
		select {
		case conn = <-localReady:
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		}
		if conn == nil {
			errCh <- fmt.Errorf("local relay connection was nil")
			return
		}
		buf := make([]byte, 32*1024)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				if wErr := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); wErr != nil {
					errCh <- wErr
					return
				}
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()
	go func() {
		for {
			msgType, payload, err := ws.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			if msgType != websocket.BinaryMessage && msgType != websocket.TextMessage {
				continue
			}
			if len(payload) == 0 {
				continue
			}
			conn, dialErr := ensureLocalConn()
			if dialErr != nil {
				errCh <- dialErr
				return
			}
			if _, wErr := conn.Write(payload); wErr != nil {
				errCh <- wErr
				return
			}
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-errCh:
		return nil
	}
}
