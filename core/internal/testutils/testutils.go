package testutils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest/observer"

	// NOTE: To avoid circular dependencies, this package MUST NOT import
	// anything from "github.com/smartcontractkit/chainlink/v2/core"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
)

const (
	// Password just a password we use everywhere for testing
	Password = "16charlengthp4SsW0rD1!@#_"
)

var FixtureChainID = evmtestutils.FixtureChainID

// SimulatedChainID is the chain ID for the go-ethereum simulated backend
var SimulatedChainID = big.NewInt(1337)

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

// NewPrivateKeyAndAddress returns a new private key and the corresponding address
func NewPrivateKeyAndAddress(t testing.TB) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, address
}

// NewRandomEVMChainID returns a suitable random chain ID that will not conflict
// with fixtures
func NewRandomEVMChainID() *big.Int {
	id := mrand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Random32Byte returns a random [32]byte
func Random32Byte() (b [32]byte) {
	copy(b[:], randomBytes(32))
	return b
}

// RandomizeName appends a random UUID to the provided name
func RandomizeName(n string) string {
	id := uuid.New().String()
	return n + id
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
func Context(tb testing.TB) context.Context {
	return tb.Context()
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

// TestInterval is just a sensible poll interval that gives fast tests without
// risk of spamming
const TestInterval = 100 * time.Millisecond

// AssertEventually calls assert.Eventually with default wait and tick durations.
func AssertEventually(t *testing.T, f func() bool) bool {
	return assert.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}

// RequireEventually calls assert.Eventually with default wait and tick durations.
func RequireEventually(t *testing.T, f func() bool) {
	require.Eventually(t, f, WaitTimeout(t), TestInterval/2)
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
//	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
//	lggr := logger.TestLogger(t, observedZapCore)
func WaitForLogMessage(t *testing.T, observedLogs *observer.ObservedLogs, msg string) (le observer.LoggedEntry) {
	RequireEventually(t, func() bool {
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) {
				le = l
				return true
			}
		}
		return false
	})
	return
}

func WaitForLogMessageWithField(t *testing.T, observedLogs *observer.ObservedLogs, msg, field, value string) (le observer.LoggedEntry) {
	RequireEventually(t, func() bool {
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) && strings.Contains(l.ContextMap()[field].(string), value) {
				le = l
				return true
			}
		}
		return false
	})
	return
}

// WaitForLogMessageCount waits until at least count log message containing the
// specified msg is emitted
func WaitForLogMessageCount(t *testing.T, observedLogs *observer.ObservedLogs, msg string, count int) {
	RequireEventually(t, func() bool {
		i := 0
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

// SkipShortDB skips tb during -short runs, and notes the DB dependency.
func SkipShortDB(tb testing.TB) {
	tests.SkipShort(tb, "DB dependency")
}

func AssertCount(t testing.TB, ds sqlutil.DataSource, tableName string, expected int64) {
	t.Helper()
	ctx := Context(t)
	var count int64
	err := ds.GetContext(ctx, &count, fmt.Sprintf(`SELECT count(*) FROM %s;`, tableName))
	require.NoError(t, err)
	require.Equal(t, expected, count)
}

// Ptr takes pointer of anything
func Ptr[T any](v T) *T {
	return &v
}

func MustRandBytes(n int) (b []byte) {
	b = make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return
}
