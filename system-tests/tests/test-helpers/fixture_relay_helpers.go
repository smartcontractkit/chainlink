package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	envRemoteAgentURL  = "CRE_REMOTE_AGENT_URL"
	envRemoteAgentPort = "CRE_REMOTE_AGENT_PORT"
)

type relayOpenResponse struct {
	RelayID string `json:"relayId"`
}

type relayCloseResponse struct {
	Found bool `json:"found"`
}

type fixtureRelayHandle struct {
	relayID string
	cancel  context.CancelFunc
}

var (
	fixtureRelayMu      sync.Mutex
	fixtureRelayHandles = make(map[string]*fixtureRelayHandle)
)

// EnsureFixtureRelayForPort ensures a local fixture port is reachable from remote components.
// It is a no-op when no remote NodeSets are configured.
func EnsureFixtureRelayForPort(t *testing.T, testEnv *ttypes.TestEnvironment, relayName string, localPort int) {
	t.Helper()
	require.Greater(t, localPort, 0, "fixture relay local port must be > 0")

	cfg := resolveEnvConfigForRelay(t, testEnv)
	if !hasRemoteNodeSets(cfg) {
		return
	}

	agentBaseURL, err := resolveAgentBaseURLForRelay()
	require.NoError(t, err, "failed to resolve agent base URL for fixture relay")

	key := fmt.Sprintf("%s|%s|%d", strings.TrimSpace(relayName), agentBaseURL, localPort)
	fixtureRelayMu.Lock()
	if _, exists := fixtureRelayHandles[key]; exists {
		fixtureRelayMu.Unlock()
		return
	}
	fixtureRelayMu.Unlock()

	relayID, err := openRelay(context.Background(), agentBaseURL, relayName, localPort)
	require.NoError(t, err, "failed to open fixture relay on agent")

	ctx, cancel := context.WithCancel(context.Background())
	localFixtureAddr := net.JoinHostPort("127.0.0.1", strconv.Itoa(localPort))
	for i := 0; i < 4; i++ {
		go relayWorker(ctx, agentBaseURL, relayID, localFixtureAddr)
	}

	fixtureRelayMu.Lock()
	fixtureRelayHandles[key] = &fixtureRelayHandle{relayID: relayID, cancel: cancel}
	fixtureRelayMu.Unlock()

	t.Cleanup(func() {
		fixtureRelayMu.Lock()
		handle, ok := fixtureRelayHandles[key]
		if ok {
			delete(fixtureRelayHandles, key)
		}
		fixtureRelayMu.Unlock()
		if !ok {
			return
		}
		handle.cancel()
		_, _ = closeRelay(context.Background(), agentBaseURL, handle.relayID)
	})
}

func resolveEnvConfigForRelay(t *testing.T, testEnv *ttypes.TestEnvironment) *envconfig.Config {
	t.Helper()
	if testEnv != nil && testEnv.Config != nil {
		return testEnv.Config
	}
	configPath := strings.TrimSpace(os.Getenv("CTF_CONFIGS"))
	if configPath == "" {
		return nil
	}
	cfg := &envconfig.Config{}
	if err := cfg.Load(configPath); err != nil {
		return nil
	}
	return cfg
}

func hasRemoteNodeSets(cfg *envconfig.Config) bool {
	if cfg == nil {
		return false
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.EqualFold(strings.TrimSpace(nodeSet.Placement), string(envconfig.PlacementRemote)) {
			return true
		}
	}
	return false
}

func resolveAgentBaseURLForRelay() (string, error) {
	if v := strings.TrimSpace(os.Getenv(envRemoteAgentURL)); v != "" {
		return v, nil
	}
	hostIP, err := runtimecfg.DirectHostIP()
	if err != nil {
		return "", err
	}
	port := 8080
	if rawPort := strings.TrimSpace(os.Getenv(envRemoteAgentPort)); rawPort != "" {
		parsed, err := strconv.Atoi(rawPort)
		if err != nil || parsed <= 0 || parsed > 65535 {
			return "", fmt.Errorf("invalid %s: %q", envRemoteAgentPort, rawPort)
		}
		port = parsed
	}
	return fmt.Sprintf("http://%s:%d", hostIP, port), nil
}

func openRelay(ctx context.Context, agentBaseURL, name string, requestedPort int) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"name":          name,
		"requestedPort": requestedPort,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(agentBaseURL, "/")+"/v1/relay/open", bytes.NewReader(body))
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

func closeRelay(ctx context.Context, agentBaseURL, relayID string) (*relayCloseResponse, error) {
	body, _ := json.Marshal(map[string]any{"relayId": relayID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(agentBaseURL, "/")+"/v1/relay/close", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("close relay failed: status %s body %s", resp.Status, strings.TrimSpace(string(respBody)))
	}
	out := &relayCloseResponse{}
	if len(respBody) > 0 {
		_ = json.Unmarshal(respBody, out)
	}
	return out, nil
}

func relayWorker(ctx context.Context, agentBaseURL, relayID, localFixtureAddr string) {
	backoff := 250 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		wsURL, err := relayConnectWSURL(agentBaseURL, relayID)
		if err != nil {
			time.Sleep(backoff)
			continue
		}
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			time.Sleep(backoff)
			continue
		}

		localConn, err := net.DialTimeout("tcp", localFixtureAddr, 2*time.Second)
		if err != nil {
			_ = ws.Close()
			time.Sleep(backoff)
			continue
		}

		_ = bridgeFixtureRelayStream(ctx, ws, localConn)
		_ = localConn.Close()
		_ = ws.Close()

		if backoff < 2*time.Second {
			backoff *= 2
		}
	}
}

func relayConnectWSURL(agentBaseURL, relayID string) (string, error) {
	base := strings.TrimRight(agentBaseURL, "/")
	u, err := url.Parse(base)
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

func bridgeFixtureRelayStream(ctx context.Context, ws *websocket.Conn, localConn net.Conn) error {
	errCh := make(chan error, 2)

	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := localConn.Read(buf)
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
			if _, wErr := localConn.Write(payload); wErr != nil {
				errCh <- wErr
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err == nil || errors.Is(err, io.EOF) || websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return nil
		}
		return err
	}
}
