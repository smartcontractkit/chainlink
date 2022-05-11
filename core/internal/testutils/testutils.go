package testutils

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	mrand "math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// NOTE: To avoid circular dependencies, this package MUST NOT import
	// anything from "github.com/smartcontractkit/chainlink/core"
)

// FixtureChainID matches the chain always added by fixtures.sql
// It is set to 0 since no real chain ever has this ID and allows a virtual
// "test" chain ID to be used without clashes
var FixtureChainID = big.NewInt(0)

// SimulatedChainID is the chain ID for the go-ethereum simulated backend
var SimulatedChainID = big.NewInt(1337)

// MustNewSimTransactor returns a transactor for interacting with the
// geth simulated backend.
func MustNewSimTransactor(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	transactor, err := bind.NewKeyedTransactorWithChainID(key, SimulatedChainID)
	require.NoError(t, err)
	return transactor
}

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

// NewRandomInt64 returns a (non-cryptographically secure) random positive int64
func NewRandomInt64() int64 {
	id := mrand.Int63()
	return id
}

// NewRandomEVMChainID returns a suitable random chain ID that will not conflict
// with fixtures
func NewRandomEVMChainID() *big.Int {
	id := mrand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

// TestCtx is a context that will expire on test timeout
func TestCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), WaitTimeout(t))
	t.Cleanup(cancel)
	return ctx
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = mrand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return b
}

// Random32Byte returns a random [32]byte
func Random32Byte() (b [32]byte) {
	copy(b[:], randomBytes(32))
	return b
}

// DefaultWaitTimeout is the default wait timeout. If you have a *testing.T, use WaitTimeout instead.
const DefaultWaitTimeout = 30 * time.Second

// WaitTimeout returns a timeout based on the test's Deadline, if available.
// Especially important to use in parallel tests, as their individual execution
// can get paused for arbitrary amounts of time.
func WaitTimeout(t *testing.T) time.Duration {
	if d, ok := t.Deadline(); ok {
		// 10% buffer for cleanup and scheduling delay
		return time.Until(d) * 9 / 10
	}
	return DefaultWaitTimeout
}

// Context returns a context with the test's deadline, if available.
func Context(t *testing.T) (ctx context.Context) {
	ctx = context.Background()
	if d, ok := t.Deadline(); ok {
		var cancel func()
		ctx, cancel = context.WithDeadline(ctx, d)
		t.Cleanup(cancel)
	}
	return ctx
}

// MustParseURL parses the URL or fails the test
func MustParseURL(t testing.TB, input string) *url.URL {
	u, err := url.Parse(input)
	require.NoError(t, err)
	return u
}

// MustParseBigInt parses a big int value from string or fails the test
func MustParseBigInt(t *testing.T, input string) *big.Int {
	i := new(big.Int)
	_, err := fmt.Sscan(input, i)
	require.NoError(t, err)
	return i
}

// JSONRPCHandler is called with the method and request param(s).
// respResult will be sent immediately. notifyResult is optional, and sent after a short delay.
type JSONRPCHandler func(reqMethod string, reqParams gjson.Result) (respResult, notifyResult string)

type testWSServer struct {
	t       *testing.T
	s       *httptest.Server
	mu      sync.RWMutex
	wsconns []*websocket.Conn
}

// NewWSServer starts a websocket server which invokes callback for each message received.
// If chainID is set, then eth_chainId calls will be automatically handled.
func NewWSServer(t *testing.T, chainID *big.Int, callback JSONRPCHandler) (ts *testWSServer) {
	ts = new(testWSServer)
	ts.t = t
	ts.wsconns = make([]*websocket.Conn, 0)
	handler := ts.newWSHandler(chainID, callback)
	ts.s = httptest.NewServer(handler)
	return
}

func (ts *testWSServer) Close() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.wsconns == nil {
		ts.t.Log("Test WS server already closed")
		return
	}
	ts.s.CloseClientConnections()
	ts.s.Close()
	for _, ws := range ts.wsconns {
		ws.Close()
	}
	ts.wsconns = nil // nil indicates server closed
}

func (ts *testWSServer) WSURL() *url.URL {
	return WSServerURL(ts.t, ts.s)
}

func (ts *testWSServer) GetConns(t *testing.T) (conns []*websocket.Conn) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	if ts.wsconns == nil {
		t.Fatal("cannot get conns from closed server")
	}
	conns = append(conns, ts.wsconns...)
	return
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
	t := ts.t
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err, "Failed to upgrade WS connection")
		defer conn.Close()
		ts.mu.Lock()
		if ts.wsconns == nil {
			log.Println("Server closed")
			ts.mu.Unlock()
			return
		}
		ts.wsconns = append(ts.wsconns, conn)
		ts.mu.Unlock()
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
					log.Println("Websocket closing")
					return
				}
				log.Printf("Failed to read message: %v", err)
				return
			}
			log.Println("Received message", string(data))
			req := gjson.ParseBytes(data)
			if !req.IsObject() {
				log.Printf("Request must be object: %v", req.Type)
				return
			}
			if e := req.Get("error"); e.Exists() {
				log.Printf("Received jsonrpc error message: %v", e)
				break
			}
			m := req.Get("method")
			if m.Type != gjson.String {
				log.Printf("Method must be string: %v", m.Type)
				return
			}

			var resp, notify string
			if chainID != nil && m.String() == "eth_chainId" {
				resp = `"0x` + chainID.Text(16) + `"`
			} else {
				resp, notify = callback(m.String(), req.Get("params"))
			}
			id := req.Get("id")
			msg := fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":%s}`, id, resp)
			log.Printf("Sending message: %v", msg)
			ts.mu.Lock()
			err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
			ts.mu.Unlock()
			if err != nil {
				log.Printf("Failed to write message: %v", err)
				return
			}

			if notify != "" {
				time.Sleep(100 * time.Millisecond)
				msg := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"0x00","result":%s}}`, notify)
				log.Println("Sending message", msg)
				ts.mu.Lock()
				err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
				ts.mu.Unlock()
				if err != nil {
					log.Printf("Failed to write message: %v", err)
					return
				}
			}
		}
	})
	return handler
}

// WaitWithTimeout waits for the channel to close (or receive anything) and
// fatals the test if the default wait timeout is exceeded
func WaitWithTimeout(t *testing.T, ch <-chan struct{}, failMsg string) {
	select {
	case <-ch:
	case <-time.After(WaitTimeout(t)):
		t.Fatal(failMsg)
	}
}

// WSServerURL returns a ws:// url for the server
func WSServerURL(t *testing.T, s *httptest.Server) *url.URL {
	u, err := url.Parse(s.URL)
	require.NoError(t, err, "Failed to parse url")
	u.Scheme = "ws"
	return u
}

// IntToHex converts int to geth-compatible hex
func IntToHex(n int) string {
	return hexutil.EncodeBig(big.NewInt(int64(n)))
}

// TestInterval is just a sensible poll interval that gives fast tests without
// risk of spamming
const TestInterval = 10 * time.Millisecond

// AssertEventually waits for f to return true
func AssertEventually(t *testing.T, f func() bool) {
	assert.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}

// RequireLogMessage fails the test if emitted logs don't contain the given message
func RequireLogMessage(t *testing.T, observedLogs *observer.ObservedLogs, msg string) {
	for _, l := range observedLogs.All() {
		if strings.Contains(l.Message, msg) {
			return
		}
	}
	t.Log("observed logs", observedLogs.All())
	t.Fatalf("expected observed logs to contain msg %q, but it didn't", msg)
}

// WaitForLogMessage waits until at least one log message containing the
// specified msg is emitted.
// NOTE: This does not "pop" messages so it cannot be used multiple times to
// check for new instances of the same msg. See WaitForLogMessageCount instead.
//
// Get a *observer.ObservedLogs like so:
//
// 		observedZapCore, observedLogs := observer.New(zap.DebugLevel)
// 		lggr := logger.TestLogger(t, observedZapCore)
func WaitForLogMessage(t *testing.T, observedLogs *observer.ObservedLogs, msg string) {
	AssertEventually(t, func() bool {
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) {
				return true
			}
		}
		return false
	})
}

// WaitForLogMessageCount waits until at least count log message containing the
// specified msg is emitted
func WaitForLogMessageCount(t *testing.T, observedLogs *observer.ObservedLogs, msg string, count int) {
	i := 0
	AssertEventually(t, func() bool {
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) {
				i++
				if i >= count {
					return true
				}
			}
		}
		return false
	})
}

// SkipShort skips tb during -short runs, and notes why.
func SkipShort(tb testing.TB, why string) {
	if testing.Short() {
		tb.Skipf("skipping: %s", why)
	}
}

// SkipShortDB skips tb during -short runs, and notes the DB dependency.
func SkipShortDB(tb testing.TB) {
	SkipShort(tb, "DB dependency")
}
