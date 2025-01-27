package testutils

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// JSONRPCHandler is called with the method and request param(s).
// respResult will be sent immediately. notifyResult is optional, and sent after a short delay.
type JSONRPCHandler func(reqMethod string, reqParams gjson.Result) JSONRPCResponse

type JSONRPCResponse struct {
	Result, Notify string // raw JSON (i.e. quoted strings etc.)

	Error struct {
		Code    int
		Message string
	}
}

type testWSServer struct {
	t       *testing.T
	s       *httptest.Server
	mu      sync.RWMutex
	wsconns []*websocket.Conn
	wg      sync.WaitGroup
}

// NewWSServer starts a websocket server which invokes callback for each message received.
// If chainID is set, then eth_chainId calls will be automatically handled.
func NewWSServer(t *testing.T, chainID *big.Int, callback JSONRPCHandler) (ts *testWSServer) {
	ts = new(testWSServer)
	ts.t = t
	ts.wsconns = make([]*websocket.Conn, 0)
	handler := ts.newWSHandler(chainID, callback)
	ts.s = httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return
}

func (ts *testWSServer) Close() {
	if func() bool {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		if ts.wsconns == nil {
			ts.t.Log("Test WS server already closed")
			return false
		}
		ts.s.CloseClientConnections()
		ts.s.Close()
		for _, ws := range ts.wsconns {
			ws.Close()
		}
		ts.wsconns = nil // nil indicates server closed
		return true
	}() {
		ts.wg.Wait()
	}
}

func (ts *testWSServer) WSURL() *url.URL {
	return WSServerURL(ts.t, ts.s)
}

// WSServerURL returns a ws:// url for the server
func WSServerURL(t *testing.T, s *httptest.Server) *url.URL {
	u, err := url.Parse(s.URL)
	require.NoError(t, err, "Failed to parse url")
	u.Scheme = "ws"
	return u
}

func (ts *testWSServer) MustWriteBinaryMessageSync(t *testing.T, msg string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	conns := ts.wsconns
	if len(conns) != 1 {
		t.Fatalf("expected 1 conn, got %d", len(conns))
	}
	conn := conns[0]
	err := conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
	require.NoError(t, err)
}

func (ts *testWSServer) newWSHandler(chainID *big.Int, callback JSONRPCHandler) (handler http.HandlerFunc) {
	if callback == nil {
		callback = func(method string, params gjson.Result) (resp JSONRPCResponse) { return }
	}
	t := ts.t
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ts.mu.Lock()
		if ts.wsconns == nil { // closed
			ts.mu.Unlock()
			return
		}
		ts.wg.Add(1)
		defer ts.wg.Done()
		conn, err := upgrader.Upgrade(w, r, nil)
		if !assert.NoError(t, err, "Failed to upgrade WS connection") {
			ts.mu.Unlock()
			return
		}
		defer conn.Close()
		ts.wsconns = append(ts.wsconns, conn)
		ts.mu.Unlock()

		for {
			err := ts.handleNewMsg(chainID, conn, callback)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
					ts.t.Log("Websocket closing")
					return
				}
				ts.t.Logf("Failed to read message: %v", err)
				return
			}
		}
	}
}

func (ts *testWSServer) handleNewMsg(chainID *big.Int, conn *websocket.Conn, callback JSONRPCHandler) error {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	ts.t.Log("Received message", string(data))

	req := gjson.ParseBytes(data)

	if req.IsArray() { // Handle batch request
		ts.t.Log("Received batch request")
		var responses []string
		for i, reqElem := range req.Array() {
			var response string
			response, _, err = ts.handleRequest(chainID, callback, reqElem)
			if err != nil {
				return fmt.Errorf("failed to handle elem %d of batch request: %w", i, err)
			}
			responses = append(responses, response)
		}

		return ts.writeMsg(conn, fmt.Sprintf("[%s]", strings.Join(responses, ",")))
	}
	// Handle single request
	response, asyncResponse, err := ts.handleRequest(chainID, callback, req)
	if err != nil {
		return fmt.Errorf("failed to handle request: %w", err)
	}

	if response != "" {
		err = ts.writeMsg(conn, response)
		if err != nil {
			return err
		}
	}

	if asyncResponse != "" {
		time.Sleep(100 * time.Millisecond)
		return ts.writeMsg(conn, asyncResponse)
	}

	return nil
}

func (ts *testWSServer) handleRequest(chainID *big.Int, callback JSONRPCHandler, req gjson.Result) (response, asyncResponse string, err error) {
	if e := req.Get("error"); e.Exists() {
		ts.t.Logf("Received jsonrpc error: %v", e)
		return
	}

	m := req.Get("method")
	if m.Type != gjson.String {
		err = fmt.Errorf("method must be string: %v", m.Type)
		return
	}

	var resp JSONRPCResponse
	if chainID != nil && m.String() == "eth_chainId" {
		resp.Result = `"0x` + chainID.Text(16) + `"`
	} else if m.String() == "eth_syncing" {
		resp.Result = "false"
	} else {
		resp = callback(m.String(), req.Get("params"))
	}
	id := req.Get("id")
	if resp.Error.Message != "" {
		response = fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"error":{"code":%d,"message":"%s"}}`, id, resp.Error.Code, resp.Error.Message)
	} else {
		response = fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":%s}`, id, resp.Result)
	}

	if resp.Notify != "" {
		asyncResponse = fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"0x00","result":%s}}`, resp.Notify)
	}

	return
}

func (ts *testWSServer) writeMsg(conn *websocket.Conn, msg string) error {
	ts.t.Logf("Sending message: %v", msg)
	ts.mu.Lock()
	err := conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
	ts.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to write msg: %w", err)
	}
	return nil
}

type RawSub[T any] struct {
	ch  chan<- T
	err <-chan error
}

func NewRawSub[T any](ch chan<- T, err <-chan error) RawSub[T] {
	return RawSub[T]{ch: ch, err: err}
}

func (r *RawSub[T]) CloseCh() {
	close(r.ch)
}

func (r *RawSub[T]) TrySend(t T) {
	select {
	case <-r.err:
	case r.ch <- t:
	}
}
