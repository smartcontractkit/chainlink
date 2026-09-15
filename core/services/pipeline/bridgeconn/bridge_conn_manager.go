package bridgeconn

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	stdErrors "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

//nolint:revive // Interface name matches existing project convention.
type BridgeConnManager interface {
	GetObservation(bridge bridges.BridgeType, requestData map[string]any) ([]byte, error)
}

var (
	ErrBridgeObservationNotFound = stdErrors.New("bridge observation not found")
	ErrBridgeObservationExpired  = stdErrors.New("bridge observation expired")
)

// observationTTL bounds how long a cached observation may be served before it is
// treated as stale. Hardcoded for now; may become configurable later.
const observationTTL = 60 * time.Second

// cacheEntry pairs a cached observation with the time it was stored, so
// GetObservation can reject entries older than observationTTL.
type cacheEntry struct {
	payload  []byte
	storedAt time.Time
}

// bridgeConnManager is a package-level singleton: one observation cache plus one
// EAConn registry shared by every pipeline run in the process. It self-initializes
// lazily as bridges are first used; there is no explicit start/close lifecycle.
type bridgeConnManager struct {
	mu    sync.RWMutex
	cache map[[32]byte]cacheEntry

	connsMu sync.Mutex
	conns   map[string]*eaConn // bridge name -> EAConn
	lggr    logger.Logger      // immutable after singleton creation

	dial  eaStreamDialer
	clock clockwork.Clock
}

var (
	defaultBridgeConnManager     BridgeConnManager
	defaultBridgeConnManagerOnce sync.Once
)

// NewBridgeConnManager returns the package-level singleton, creating it on the
// first call with the supplied logger. Subsequent calls return the same manager
// and ignore the logger. The logger is immutable after creation.
func NewBridgeConnManager(lggr logger.Logger) BridgeConnManager {
	defaultBridgeConnManagerOnce.Do(func() {
		if lggr == nil {
			lggr = logger.Nop()
		}
		defaultBridgeConnManager = &bridgeConnManager{
			cache: make(map[[32]byte]cacheEntry),
			conns: make(map[string]*eaConn),
			lggr:  lggr,
			dial:  dialGRPCStream,
			clock: clockwork.NewRealClock(),
		}
	})
	return defaultBridgeConnManager
}

func (m *bridgeConnManager) GetObservation(bridge bridges.BridgeType, requestData map[string]any) ([]byte, error) {
	bridgeName := strings.TrimPrefix(bridge.Name.String(), "bridge-")
	data, err := subscriptionData(requestData)
	if err != nil {
		return nil, fmt.Errorf("bridge %q: %w", bridgeName, err)
	}
	key, err := bridgeObservationCacheKey(bridgeName, data)
	if err != nil {
		return nil, err
	}
	m.lggr.Debugw("cache key generated", "key", hex.EncodeToString(key[:]), "bridge", bridgeName, "data", data)
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
	if m.clock.Now().Sub(entry.storedAt) > observationTTL {
		return nil, fmt.Errorf("%w for bridge %q", ErrBridgeObservationExpired, bridgeName)
	}
	payload := make([]byte, len(entry.payload))
	copy(payload, entry.payload)
	return payload, nil
}

// PutObservation stores observation bytes under the given payload hash key.
// It is used by EAConn receiver loops and white-box tests in this package.
func (m *bridgeConnManager) PutObservation(key [32]byte, observation []byte) {
	payload := make([]byte, len(observation))
	copy(payload, observation)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[key] = cacheEntry{payload: payload, storedAt: m.clock.Now()}
}

// SeedObservation computes the bridge observation key from request data and
// stores a cache entry. This is intended for tests.
func (m *bridgeConnManager) SeedObservation(bridge bridges.BridgeType, requestData map[string]any, observation []byte) error {
	data, err := subscriptionData(requestData)
	if err != nil {
		return err
	}
	key, err := bridgeObservationCacheKey(strings.TrimPrefix(bridge.Name.String(), "bridge-"), data)
	if err != nil {
		return err
	}
	m.PutObservation(key, observation)
	return nil
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
	m.dial = func(_ context.Context, _ string, _ bool) (eaStreamClient, error) {
		return nil, errStreamDialingDisabledForTest
	}
}

// subscriptionData extracts the inner "data" object from a bridge task's request
// payload: the only part sent as Subscription.Data and the only part the adapter
// hashes (see ObservationPayloadHash on the streams-adapter side).
func subscriptionData(requestData map[string]any) (map[string]any, error) {
	data, ok := requestData["data"].(map[string]any)
	if !ok || len(data) == 0 {
		return nil, stdErrors.New("request data is missing a non-empty \"data\" field required for subscription")
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
