package pipeline

import (
	"context"
	"crypto/sha256"
	stdErrors "errors"
	"fmt"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type BridgeConnManager interface {
	GetObservation(bridge bridges.BridgeType, requestData MapParam) ([]byte, error)
	PutObservation(key [32]byte, observation []byte)
}

var (
	ErrBridgeObservationNotFound = stdErrors.New("bridge observation not found")
)

// bridgeConnManager is a package-level singleton: one observation cache plus one
// EAConn registry shared by every pipeline run in the process. It self-initializes
// lazily as bridges are first used; there is no explicit start/close lifecycle.
type bridgeConnManager struct {
	mu    sync.RWMutex
	cache map[[32]byte][]byte

	connsMu sync.Mutex
	conns   map[string]*eaConn // bridge name -> EAConn
	lggr    logger.Logger      // guarded by connsMu; set at most once, from NewBridgeConnManager

	dial eaStreamDialer
}

var defaultBridgeConnManager BridgeConnManager = &bridgeConnManager{
	cache: make(map[[32]byte][]byte),
	conns: make(map[string]*eaConn),
	lggr:  logger.Nop(),
	dial:  dialGRPCStream,
}

// NewBridgeConnManager returns the package-level singleton. Passing a logger sets
// it on the singleton for use by lazily-created EAConns; it's expected to be
// called once, from PipelineRunner startup, with all other call sites (fallback
// construction, tests) using the zero-arg form and getting whatever logger (or
// the Nop default) is already set.
func NewBridgeConnManager(lggr ...logger.Logger) BridgeConnManager {
	m := defaultBridgeConnManager.(*bridgeConnManager)
	if len(lggr) > 0 && lggr[0] != nil {
		m.connsMu.Lock()
		m.lggr = lggr[0]
		m.connsMu.Unlock()
	}
	return m
}

func (m *bridgeConnManager) GetObservation(bridge bridges.BridgeType, requestData MapParam) ([]byte, error) {
	bridgeName := strings.TrimPrefix(bridge.Name.String(), "bridge-")
	data, err := subscriptionData(requestData)
	if err != nil {
		return nil, fmt.Errorf("bridge %q: %w", bridgeName, err)
	}
	key, err := bridgeObservationCacheKey(bridgeName, data)
	if err != nil {
		return nil, err
	}
	m.lggr.Debugw("cache key generated", "key", fmt.Sprintf("%x", key), "bridge", bridgeName, "data", data)
	subscription, err := structpb.NewStruct(data)
	if err != nil {
		return nil, fmt.Errorf("failed to build subscription payload for bridge %q: %w", bridgeName, err)
	}
	m.getOrCreateConn(bridgeName, bridge.URL).registerAsset(key, subscription)

	m.mu.RLock()
	entry, ok := m.cache[key]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w for bridge %q", ErrBridgeObservationNotFound, bridgeName)
	}
	payload := make([]byte, len(entry))
	copy(payload, entry)
	return payload, nil
}

// putObservation is called by an EAConn's receiver loop when it accepts a
// streamed observation for a currently registered asset.
func (m *bridgeConnManager) PutObservation(key [32]byte, observation []byte) {
	payload := make([]byte, len(observation))
	copy(payload, observation)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[key] = payload
}

// getOrCreateConn returns the bridge's persistent EAConn, lazily creating and
// starting it on first use.
func (m *bridgeConnManager) getOrCreateConn(bridgeName string, bridgeURL models.WebURL) *eaConn {

	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	if conn, ok := m.conns[bridgeName]; ok {
		return conn
	}

	conn := newEAConn(bridgeName, bridgeURL, m)
	m.conns[bridgeName] = conn
	conn.start()
	return conn
}

var errStreamDialingDisabledForTest = stdErrors.New("EAConn stream dialing disabled for test")

// DisableEAConnDialingForTest replaces the manager's stream dialer with one that
// fails immediately without any network I/O, for tests that seed the observation
// cache directly and must not depend on a real streams-adapter connection. It
// mutates the shared package-level singleton and is intended for test setup only.
func (m *bridgeConnManager) DisableEAConnDialingForTest() {
	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	m.dial = func(_ context.Context, _ string) (eaStreamClient, error) {
		return nil, errStreamDialingDisabledForTest
	}
}

// subscriptionData extracts the inner "data" object from a bridge task's request
// payload — the only part sent as Subscription.Data and the only part the adapter
// hashes (see ObservationPayloadHash on the streams-adapter side).
func subscriptionData(requestData MapParam) (map[string]any, error) {
	data, ok := requestData["data"].(map[string]any)
	if !ok || len(data) == 0 {
		return nil, fmt.Errorf("request data is missing a non-empty \"data\" field required for subscription")
	}
	return data, nil
}

// bridgeObservationCacheKey mirrors the streams-adapter's own ObservationPayloadHash.
// The adapter is configured with its own adapterName equal to this bridge's name,
// so payload_hash on an accepted observation equals this same key.
func bridgeObservationCacheKey(bridgeName string, data map[string]any) ([32]byte, error) {
	lookupBytes, err := json.Marshal(data)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal bridge lookup payload: %w", err)
	}
	b := make([]byte, 0, len(bridgeName)+len(lookupBytes))
	b = append(b, bridgeName...)
	b = append(b, lookupBytes...)
	return sha256.Sum256(b), nil
}
