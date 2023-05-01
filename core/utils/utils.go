// Package utils is used for common functions and tools used across the codebase.
package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	mrand "math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	"golang.org/x/exp/constraints"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/jpillora/backoff"
	"github.com/libp2p/go-libp2p-core/peer"
	pkgerrors "github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

const (
	// DefaultSecretSize is the entropy in bytes to generate a base64 string of 64 characters.
	DefaultSecretSize = 48
	// EVMWordByteLen the length of an EVM Word Byte
	EVMWordByteLen = 32

	// defaultErrorBufferCap is the default cap on the errors an error buffer can store at any time
	defaultErrorBufferCap = 50
)

// ZeroAddress is an address of all zeroes, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

func RandomAddress() common.Address {
	b := make([]byte, 20)
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return common.BytesToAddress(b)
}

func RandomBytes32() (r [32]byte) {
	b := make([]byte, 32)
	_, _ = rand.Read(b[:]) // Assignment for errcheck. Only used in tests so we can ignore.
	copy(r[:], b)
	return
}

func Bytes32ToSlice(a [32]byte) (r []byte) {
	r = append(r, a[:]...)
	return
}

func MustNewPeerID() string {
	_, pubKey, err := cryptop2p.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		panic(err)
	}
	return peerID.String()
}

// EmptyHash is a hash of all zeroes, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000000000000000000000000000
var EmptyHash = common.Hash{}

// Uint256ToBytes is x represented as the bytes of a uint256
func Uint256ToBytes(x *big.Int) (uint256 []byte, err error) {
	if x.Cmp(MaxUint256) > 0 {
		return nil, fmt.Errorf("too large to convert to uint256")
	}
	uint256 = common.LeftPadBytes(x.Bytes(), EVMWordByteLen)
	if x.Cmp(big.NewInt(0).SetBytes(uint256)) != 0 {
		panic("failed to round-trip uint256 back to source big.Int")
	}
	return uint256, err
}

// ISO8601UTC formats given time to ISO8601.
func ISO8601UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// DurationFromNow returns the amount of time since the Time
// field was last updated.
func DurationFromNow(t time.Time) time.Duration {
	return time.Until(t)
}

// FormatJSON applies indent to format a JSON response.
func FormatJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// NewBytes32ID returns a randomly generated UUID that conforms to
// Ethereum bytes32.
func NewBytes32ID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// NewSecret returns a new securely random sequence of n bytes of entropy.  The
// result is a base64 encoded string.
//
// Panics on failed attempts to read from system's PRNG.
func NewSecret(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(pkgerrors.Wrap(err, "generating secret failed"))
	}
	return base64.StdEncoding.EncodeToString(b)
}

// RemoveHexPrefix removes the prefix (0x) of a given hex string.
func RemoveHexPrefix(str string) string {
	if HasHexPrefix(str) {
		return str[2:]
	}
	return str
}

// HasHexPrefix returns true if the string starts with 0x.
func HasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// IsEmptyAddress checks that the address is empty, synonymous with the zero
// account/address. No logs can come from this address, as there is no contract
// present there.
//
// See https://stackoverflow.com/questions/48219716/what-is-address0-in-solidity
// for the more info on the zero address.
func IsEmptyAddress(addr common.Address) bool {
	return addr == ZeroAddress
}

// StringToHex converts a standard string to a hex encoded string.
func StringToHex(in string) string {
	return AddHexPrefix(hex.EncodeToString([]byte(in)))
}

// AddHexPrefix adds the prefix (0x) to a given hex string.
func AddHexPrefix(str string) string {
	if len(str) < 2 || len(str) > 1 && strings.ToLower(str[0:2]) != "0x" {
		str = "0x" + str
	}
	return str
}

// IsEmpty returns true if bytes contains only zero values, or has len 0.
func IsEmpty(bytes []byte) bool {
	for _, b := range bytes {
		if b != 0 {
			return false
		}
	}
	return true
}

// Sleeper interface is used for tasks that need to be done on some
// interval, excluding Cron, like reconnecting.
type Sleeper interface {
	Reset()
	Sleep()
	After() time.Duration
	Duration() time.Duration
}

// BackoffSleeper is a sleeper that backs off on subsequent attempts.
type BackoffSleeper struct {
	backoff.Backoff
	beenRun atomic.Bool
}

// NewBackoffSleeper returns a BackoffSleeper that is configured to
// sleep for 0 seconds initially, then backs off from 1 second minimum
// to 10 seconds maximum.
func NewBackoffSleeper() *BackoffSleeper {
	return &BackoffSleeper{
		Backoff: backoff.Backoff{
			Min: 1 * time.Second,
			Max: 10 * time.Second,
		},
	}
}

// Sleep waits for the given duration, incrementing the back off.
func (bs *BackoffSleeper) Sleep() {
	if bs.beenRun.CompareAndSwap(false, true) {
		return
	}
	time.Sleep(bs.Backoff.Duration())
}

// After returns the duration for the next stop, and increments the backoff.
func (bs *BackoffSleeper) After() time.Duration {
	if bs.beenRun.CompareAndSwap(false, true) {
		return 0
	}
	return bs.Backoff.Duration()
}

// Duration returns the current duration value.
func (bs *BackoffSleeper) Duration() time.Duration {
	if !bs.beenRun.Load() {
		return 0
	}
	return bs.ForAttempt(bs.Attempt())
}

// Reset resets the backoff intervals.
func (bs *BackoffSleeper) Reset() {
	bs.beenRun.Store(false)
	bs.Backoff.Reset()
}

// RetryWithBackoff retries the sleeper and backs off if not Done
func RetryWithBackoff(ctx context.Context, fn func() (retry bool)) {
	sleeper := NewBackoffSleeper()
	sleeper.Reset()
	for {
		retry := fn()
		if !retry {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleeper.After()):
			continue
		}
	}
}

// UnmarshalToMap takes an input json string and returns a map[string]interface i.e. a raw object
func UnmarshalToMap(input string) (map[string]interface{}, error) {
	var output map[string]interface{}
	err := json.Unmarshal([]byte(input), &output)
	return output, err
}

// MustUnmarshalToMap performs UnmarshalToMap, panics upon failure
func MustUnmarshalToMap(input string) map[string]interface{} {
	output, err := UnmarshalToMap(input)
	if err != nil {
		panic(err)
	}
	return output
}

// HashPassword wraps around bcrypt.GenerateFromPassword for a friendlier API.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash wraps around bcrypt.CompareHashAndPassword for a friendlier API.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Keccak256 is a simplified interface for the legacy SHA3 implementation that
// Ethereum uses.
func Keccak256(in []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(in)
	return hash.Sum(nil), err
}

func Keccak256Fixed(in []byte) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	hash.Write(in)
	var h [32]byte
	copy(h[:], hash.Sum(nil))
	return h
}

// Sha256 returns a hexadecimal encoded string of a hashed input
func Sha256(in string) (string, error) {
	hasher := sha3.New256()
	_, err := hasher.Write([]byte(in))
	if err != nil {
		return "", pkgerrors.Wrap(err, "sha256 write error")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// IsQuoted checks if the first and last characters are either " or '.
func IsQuoted(input []byte) bool {
	return len(input) >= 2 &&
		((input[0] == '"' && input[len(input)-1] == '"') ||
			(input[0] == '\'' && input[len(input)-1] == '\''))
}

// RemoveQuotes removes the first and last character if they are both either
// " or ', otherwise it is a noop.
func RemoveQuotes(input []byte) []byte {
	if IsQuoted(input) {
		return input[1 : len(input)-1]
	}
	return input
}

// EIP55CapitalizedAddress returns true iff possibleAddressString has the correct
// capitalization for an Ethereum address, per EIP 55
func EIP55CapitalizedAddress(possibleAddressString string) bool {
	if !HasHexPrefix(possibleAddressString) {
		possibleAddressString = "0x" + possibleAddressString
	}
	EIP55Capitalized := common.HexToAddress(possibleAddressString).Hex()
	return possibleAddressString == EIP55Capitalized
}

// ParseEthereumAddress returns addressString as a go-ethereum Address, or an
// error if it's invalid, e.g. if EIP 55 capitalization check fails
func ParseEthereumAddress(addressString string) (common.Address, error) {
	if !common.IsHexAddress(addressString) {
		return common.Address{}, fmt.Errorf(
			"not a valid Ethereum address: %s", addressString)
	}
	address := common.HexToAddress(addressString)
	if !EIP55CapitalizedAddress(addressString) {
		return common.Address{}, fmt.Errorf(
			"%s treated as Ethereum address, but it has an invalid capitalization! "+
				"The correctly-capitalized address would be %s, but "+
				"check carefully before copying and pasting! ",
			addressString, address.Hex())
	}
	return address, nil
}

// MustHash returns the keccak256 hash, or panics on failure.
func MustHash(in string) common.Hash {
	out, err := Keccak256([]byte(in))
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(out)
}

// JustError takes a tuple and returns the last entry, the error.
func JustError(_ interface{}, err error) error {
	return err
}

var zero = big.NewInt(0)

// CheckUint256 returns an error if n is out of bounds for a uint256
func CheckUint256(n *big.Int) error {
	if n.Cmp(zero) < 0 || n.Cmp(MaxUint256) >= 0 {
		return fmt.Errorf("number out of range for uint256")
	}
	return nil
}

// HexToUint256 returns the uint256 represented by s, or an error if it doesn't
// represent one.
func HexToUint256(s string) (*big.Int, error) {
	rawNum, err := hexutil.Decode(s)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "while parsing %s as hex: ", s)
	}
	rv := big.NewInt(0).SetBytes(rawNum) // can't be negative number
	if err := CheckUint256(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

// HexToBig parses the given hex string or panics if it is invalid.
func HexToBig(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, s))
	}
	return n
}

// Uint256ToBytes32 returns the bytes32 encoding of the big int provided
func Uint256ToBytes32(n *big.Int) []byte {
	if n.BitLen() > 256 {
		panic("vrf.uint256ToBytes32: too big to marshal to uint256")
	}
	return common.LeftPadBytes(n.Bytes(), 32)
}

// WaitGroupChan creates a channel that closes when the provided sync.WaitGroup is done.
func WaitGroupChan(wg *sync.WaitGroup) <-chan struct{} {
	chAwait := make(chan struct{})
	go func() {
		defer close(chAwait)
		wg.Wait()
	}()
	return chAwait
}

// WithCloseChan wraps a context so that it is canceled if the passed in channel is closed.
// Deprecated: Call StopChan.Ctx directly
func WithCloseChan(parentCtx context.Context, chStop chan struct{}) (context.Context, context.CancelFunc) {
	return StopChan(chStop).Ctx(parentCtx)
}

// ContextFromChan creates a context that finishes when the provided channel receives or is closed.
// Deprecated: Call StopChan.NewCtx directly.
func ContextFromChan(chStop chan struct{}) (context.Context, context.CancelFunc) {
	return StopChan(chStop).NewCtx()
}

// ContextFromChanWithTimeout creates a context with a timeout that finishes when the provided channel receives or is closed.
// Deprecated: Call StopChan.CtxCancel directly
func ContextFromChanWithTimeout(chStop chan struct{}, timeout time.Duration) (context.Context, context.CancelFunc) {
	return StopChan(chStop).CtxCancel(context.WithTimeout(context.Background(), timeout))
}

// A StopChan signals when some work should stop.
type StopChan chan struct{}

// NewCtx returns a background [context.Context] that is cancelled when StopChan is closed.
func (s StopChan) NewCtx() (context.Context, context.CancelFunc) {
	return StopRChan((<-chan struct{})(s)).NewCtx()
}

// Ctx cancels a [context.Context] when StopChan is closed.
func (s StopChan) Ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	return StopRChan((<-chan struct{})(s)).Ctx(ctx)
}

// CtxCancel cancels a [context.Context] when StopChan is closed.
// Returns ctx and cancel unmodified, for convenience.
func (s StopChan) CtxCancel(ctx context.Context, cancel context.CancelFunc) (context.Context, context.CancelFunc) {
	return StopRChan((<-chan struct{})(s)).CtxCancel(ctx, cancel)
}

// A StopRChan signals when some work should stop.
// This version is receive-only.
type StopRChan <-chan struct{}

// NewCtx returns a background [context.Context] that is cancelled when StopChan is closed.
func (s StopRChan) NewCtx() (context.Context, context.CancelFunc) {
	return s.Ctx(context.Background())
}

// Ctx cancels a [context.Context] when StopChan is closed.
func (s StopRChan) Ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	return s.CtxCancel(context.WithCancel(ctx))
}

// CtxCancel cancels a [context.Context] when StopChan is closed.
// Returns ctx and cancel unmodified, for convenience.
func (s StopRChan) CtxCancel(ctx context.Context, cancel context.CancelFunc) (context.Context, context.CancelFunc) {
	go func() {
		select {
		case <-s:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

// DependentAwaiter contains Dependent funcs
type DependentAwaiter interface {
	AwaitDependents() <-chan struct{}
	AddDependents(n int)
	DependentReady()
}

type dependentAwaiter struct {
	wg *sync.WaitGroup
	ch <-chan struct{}
}

// NewDependentAwaiter creates a new DependentAwaiter
func NewDependentAwaiter() DependentAwaiter {
	return &dependentAwaiter{
		wg: &sync.WaitGroup{},
	}
}

func (da *dependentAwaiter) AwaitDependents() <-chan struct{} {
	if da.ch == nil {
		da.ch = WaitGroupChan(da.wg)
	}
	return da.ch
}

func (da *dependentAwaiter) AddDependents(n int) {
	da.wg.Add(n)
}

func (da *dependentAwaiter) DependentReady() {
	da.wg.Done()
}

// BoundedQueue is a FIFO queue that discards older items when it reaches its capacity.
type BoundedQueue[T any] struct {
	capacity int
	items    []T
	mu       sync.RWMutex
}

// NewBoundedQueue creates a new BoundedQueue instance
func NewBoundedQueue[T any](capacity int) *BoundedQueue[T] {
	var bq BoundedQueue[T]
	bq.capacity = capacity
	return &bq
}

// Add appends items to a BoundedQueue
func (q *BoundedQueue[T]) Add(x T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, x)
	if len(q.items) > q.capacity {
		excess := len(q.items) - q.capacity
		q.items = q.items[excess:]
	}
}

// Take pulls the first item from the array and removes it
func (q *BoundedQueue[T]) Take() (t T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return
	}
	t = q.items[0]
	q.items = q.items[1:]
	return
}

// Empty check is a BoundedQueue is empty
func (q *BoundedQueue[T]) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items) == 0
}

// Full checks if a BoundedQueue is over capacity.
func (q *BoundedQueue[T]) Full() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items) >= q.capacity
}

// BoundedPriorityQueue stores a series of BoundedQueues
// with associated priorities and capacities
type BoundedPriorityQueue[T any] struct {
	queues     map[uint]*BoundedQueue[T]
	priorities []uint
	capacities map[uint]int
	mu         sync.RWMutex
}

// NewBoundedPriorityQueue creates a new BoundedPriorityQueue
func NewBoundedPriorityQueue[T any](capacities map[uint]int) *BoundedPriorityQueue[T] {
	queues := make(map[uint]*BoundedQueue[T])
	var priorities []uint
	for priority, capacity := range capacities {
		priorities = append(priorities, priority)
		queues[priority] = NewBoundedQueue[T](capacity)
	}
	sort.Slice(priorities, func(i, j int) bool { return priorities[i] < priorities[j] })
	bpq := BoundedPriorityQueue[T]{
		queues:     queues,
		priorities: priorities,
		capacities: capacities,
	}
	return &bpq
}

// Add pushes an item into a subque within a BoundedPriorityQueue
func (q *BoundedPriorityQueue[T]) Add(priority uint, x T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	subqueue, exists := q.queues[priority]
	if !exists {
		panic(fmt.Sprintf("nonexistent priority: %v", priority))
	}

	subqueue.Add(x)
}

// Take takes from the BoundedPriorityQueue's subque
func (q *BoundedPriorityQueue[T]) Take() (t T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, priority := range q.priorities {
		queue := q.queues[priority]
		if queue.Empty() {
			continue
		}
		return queue.Take()
	}
	return
}

// Empty checks the BoundedPriorityQueue
// if all subqueues are empty
func (q *BoundedPriorityQueue[T]) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, priority := range q.priorities {
		queue := q.queues[priority]
		if !queue.Empty() {
			return false
		}
	}
	return true
}

// WrapIfError decorates an error with the given message.  It is intended to
// be used with `defer` statements, like so:
//
//	func SomeFunction() (err error) {
//	    defer WrapIfError(&err, "error in SomeFunction:")
//
//	    ...
//	}
func WrapIfError(err *error, msg string) {
	if *err != nil {
		*err = pkgerrors.Wrap(*err, msg)
	}
}

// TickerBase is an interface for pausable tickers.
type TickerBase interface {
	Resume()
	Pause()
	Destroy()
	Ticks() <-chan time.Time
}

// PausableTicker stores a ticker with a duration
type PausableTicker struct {
	ticker   *time.Ticker
	duration time.Duration
	mu       *sync.RWMutex
}

// NewPausableTicker creates a new PausableTicker
func NewPausableTicker(duration time.Duration) PausableTicker {
	return PausableTicker{
		duration: duration,
		mu:       &sync.RWMutex{},
	}
}

// Ticks retrieves the ticks from a PausableTicker
func (t *PausableTicker) Ticks() <-chan time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.ticker == nil {
		return nil
	}
	return t.ticker.C
}

// Pause pauses a PausableTicker
func (t *PausableTicker) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

// Resume resumes a Ticker
// using a PausibleTicker's duration
func (t *PausableTicker) Resume() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ticker == nil {
		t.ticker = time.NewTicker(t.duration)
	}
}

// Destroy pauses the PausibleTicker
func (t *PausableTicker) Destroy() {
	t.Pause()
}

// CronTicker is like a time.Ticker but for a cron schedule.
type CronTicker struct {
	*cron.Cron
	ch      chan time.Time
	beenRun atomic.Bool
}

// NewCronTicker returns a new CrontTicker for the given schedule.
func NewCronTicker(schedule string) (CronTicker, error) {
	cron := cron.New(cron.WithSeconds())
	ch := make(chan time.Time, 1)
	_, err := cron.AddFunc(schedule, func() {
		select {
		case ch <- time.Now():
		default:
		}
	})
	if err != nil {
		return CronTicker{}, err
	}
	return CronTicker{Cron: cron, ch: ch}, nil
}

// Start - returns true if the CronTicker was actually started, false otherwise
func (t *CronTicker) Start() bool {
	if t.Cron != nil {
		if t.beenRun.CompareAndSwap(false, true) {
			t.Cron.Start()
			return true
		}
	}
	return false
}

// Stop - returns true if the CronTicker was actually stopped, false otherwise
func (t *CronTicker) Stop() bool {
	if t.Cron != nil {
		if t.beenRun.CompareAndSwap(true, false) {
			t.Cron.Stop()
			return true
		}
	}
	return false
}

// Ticks returns the underlying chanel.
func (t *CronTicker) Ticks() <-chan time.Time {
	return t.ch
}

// ValidateCronSchedule returns an error if the given schedule is invalid.
func ValidateCronSchedule(schedule string) error {
	if !(strings.HasPrefix(schedule, "CRON_TZ=") || strings.HasPrefix(schedule, "@every ")) {
		return errors.New("cron schedule must specify a time zone using CRON_TZ, e.g. 'CRON_TZ=UTC 5 * * * *', or use the @every syntax, e.g. '@every 1h30m'")
	}
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(schedule)
	return pkgerrors.Wrapf(err, "invalid cron schedule '%v'", schedule)
}

// ResettableTimer stores a timer
type ResettableTimer struct {
	timer *time.Timer
	mu    *sync.RWMutex
}

// NewResettableTimer creates a new ResettableTimer
func NewResettableTimer() ResettableTimer {
	return ResettableTimer{
		mu: &sync.RWMutex{},
	}
}

// Ticks retrieves the ticks from a ResettableTimer
func (t *ResettableTimer) Ticks() <-chan time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.timer == nil {
		return nil
	}
	return t.timer.C
}

// Stop stops a ResettableTimer
func (t *ResettableTimer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
}

// Reset stops a ResettableTimer
// and resets it with a new duration
func (t *ResettableTimer) Reset(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.timer != nil {
		t.timer.Stop()
	}
	t.timer = time.NewTimer(duration)
}

// EVMBytesToUint64 converts
// a bytebuffer to uint64
func EVMBytesToUint64(buf []byte) uint64 {
	var result uint64
	for _, b := range buf {
		result = result<<8 + uint64(b)
	}
	return result
}

type errNotStarted struct {
	state StartStopOnceState
}

func (e *errNotStarted) Error() string {
	return fmt.Sprintf("service is %q, not started", e.state)
}

var (
	ErrAlreadyStopped      = errors.New("already stopped")
	ErrCannotStopUnstarted = errors.New("cannot stop unstarted service")
)

// StartStopOnce contains a StartStopOnceState integer
type StartStopOnce struct {
	state        atomic.Int32
	sync.RWMutex // lock is held during startup/shutdown, RLock is held while executing functions dependent on a particular state

	// SvcErrBuffer is an ErrorBuffer that let service owners track critical errors happening in the service.
	//
	// SvcErrBuffer.SetCap(int) Overrides buffer limit from defaultErrorBufferCap
	// SvcErrBuffer.Append(error) Appends an error to the buffer
	// SvcErrBuffer.Flush() error returns all tracked errors as a single joined error
	SvcErrBuffer ErrorBuffer
}

// StartStopOnceState holds the state for StartStopOnce
type StartStopOnceState int32

// nolint
const (
	StartStopOnce_Unstarted StartStopOnceState = iota
	StartStopOnce_Started
	StartStopOnce_Starting
	StartStopOnce_StartFailed
	StartStopOnce_Stopping
	StartStopOnce_Stopped
	StartStopOnce_StopFailed
)

func (s StartStopOnceState) String() string {
	switch s {
	case StartStopOnce_Unstarted:
		return "Unstarted"
	case StartStopOnce_Started:
		return "Started"
	case StartStopOnce_Starting:
		return "Starting"
	case StartStopOnce_StartFailed:
		return "StartFailed"
	case StartStopOnce_Stopping:
		return "Stopping"
	case StartStopOnce_Stopped:
		return "Stopped"
	case StartStopOnce_StopFailed:
		return "StopFailed"
	default:
		return fmt.Sprintf("unrecognized state: %d", s)
	}
}

// StartOnce sets the state to Started
func (once *StartStopOnce) StartOnce(name string, fn func() error) error {
	// SAFETY: We do this compare-and-swap outside of the lock so that
	// concurrent StartOnce() calls return immediately.
	success := once.state.CompareAndSwap(int32(StartStopOnce_Unstarted), int32(StartStopOnce_Starting))

	if !success {
		return pkgerrors.Errorf("%v has already been started once; state=%v", name, StartStopOnceState(once.state.Load()))
	}

	once.Lock()
	defer once.Unlock()

	// Setting cap before calling startup fn in case of crits in startup
	once.SvcErrBuffer.SetCap(defaultErrorBufferCap)
	err := fn()

	if err == nil {
		success = once.state.CompareAndSwap(int32(StartStopOnce_Starting), int32(StartStopOnce_Started))
	} else {
		success = once.state.CompareAndSwap(int32(StartStopOnce_Starting), int32(StartStopOnce_StartFailed))
	}

	if !success {
		// SAFETY: If this is reached, something must be very wrong: once.state
		// was tampered with outside of the lock.
		panic(fmt.Sprintf("%v entered unreachable state, unable to set state to started", name))
	}

	return err
}

// StopOnce sets the state to Stopped
func (once *StartStopOnce) StopOnce(name string, fn func() error) error {
	// SAFETY: We hold the lock here so that Stop blocks until StartOnce
	// executes. This ensures that a very fast call to Stop will wait for the
	// code to finish starting up before teardown.
	once.Lock()
	defer once.Unlock()

	success := once.state.CompareAndSwap(int32(StartStopOnce_Started), int32(StartStopOnce_Stopping))

	if !success {
		state := once.state.Load()
		switch state {
		case int32(StartStopOnce_Stopped):
			return pkgerrors.Wrapf(ErrAlreadyStopped, "%s has already been stopped", name)
		case int32(StartStopOnce_Unstarted):
			return pkgerrors.Wrapf(ErrCannotStopUnstarted, "%s has not been started", name)
		default:
			return pkgerrors.Errorf("%v cannot be stopped from this state; state=%v", name, StartStopOnceState(state))
		}
	}

	err := fn()

	if err == nil {
		success = once.state.CompareAndSwap(int32(StartStopOnce_Stopping), int32(StartStopOnce_Stopped))
	} else {
		success = once.state.CompareAndSwap(int32(StartStopOnce_Stopping), int32(StartStopOnce_StopFailed))
	}

	if !success {
		// SAFETY: If this is reached, something must be very wrong: once.state
		// was tampered with outside of the lock.
		panic(fmt.Sprintf("%v entered unreachable state, unable to set state to stopped", name))
	}

	return err
}

// State retrieves the current state
func (once *StartStopOnce) State() StartStopOnceState {
	state := once.state.Load()
	return StartStopOnceState(state)
}

// IfStarted runs the func and returns true only if started, otherwise returns false
func (once *StartStopOnce) IfStarted(f func()) (ok bool) {
	once.RLock()
	defer once.RUnlock()

	state := once.state.Load()

	if StartStopOnceState(state) == StartStopOnce_Started {
		f()
		return true
	}
	return false
}

// IfNotStopped runs the func and returns true if in any state other than Stopped
func (once *StartStopOnce) IfNotStopped(f func()) (ok bool) {
	once.RLock()
	defer once.RUnlock()

	state := once.state.Load()

	if StartStopOnceState(state) == StartStopOnce_Stopped {
		return false
	}
	f()
	return true
}

// Ready returns ErrNotStarted if the state is not started.
func (once *StartStopOnce) Ready() error {
	state := once.State()
	if state == StartStopOnce_Started {
		return nil
	}
	return &errNotStarted{state: state}
}

// Healthy returns ErrNotStarted if the state is not started.
// Override this per-service with more specific implementations.
func (once *StartStopOnce) Healthy() error {
	state := once.State()
	if state == StartStopOnce_Started {
		return once.SvcErrBuffer.Flush()
	}
	return &errNotStarted{state: state}
}

// EnsureClosed closes the io.Closer, returning nil if it was already
// closed or not started yet
func EnsureClosed(c io.Closer) error {
	err := c.Close()
	if errors.Is(err, ErrAlreadyStopped) || errors.Is(err, ErrCannotStopUnstarted) {
		return nil
	}
	return err
}

// WithJitter adds +/- 10% to a duration
func WithJitter(d time.Duration) time.Duration {
	// #nosec
	if d == 0 {
		return 0
	}
	jitter := mrand.Intn(int(d) / 5)
	jitter = jitter - (jitter / 2)
	return time.Duration(int(d) + jitter)
}

// NewRedialBackoff is a standard backoff to use for redialling or reconnecting to
// unreachable network endpoints
func NewRedialBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}

}

// KeyedMutex allows to lock based on particular values
type KeyedMutex struct {
	mutexes sync.Map
}

// LockInt64 locks the value for read/write
func (m *KeyedMutex) LockInt64(key int64) func() {
	value, _ := m.mutexes.LoadOrStore(key, new(sync.Mutex))
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return mtx.Unlock
}

// BoxOutput formats its arguments as fmt.Printf, and encloses them in a box of
// arrows pointing at their content, in order to better highlight it. See
// ExampleBoxOutput
func BoxOutput(errorMsgTemplate string, errorMsgValues ...interface{}) string {
	errorMsgTemplate = fmt.Sprintf(errorMsgTemplate, errorMsgValues...)
	lines := strings.Split(errorMsgTemplate, "\n")
	maxlen := 0
	for _, line := range lines {
		if len(line) > maxlen {
			maxlen = len(line)
		}
	}
	internalLength := maxlen + 4
	output := "↘" + strings.Repeat("↓", internalLength) + "↙\n" // top line
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	readme := strings.Repeat("README ", maxlen/7)
	output += "→  " + readme + strings.Repeat(" ", maxlen-len(readme)) + "  ←\n"
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	for _, line := range lines {
		output += "→  " + line + strings.Repeat(" ", maxlen-len(line)) + "  ←\n"
	}
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	output += "→  " + readme + strings.Repeat(" ", maxlen-len(readme)) + "  ←\n"
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	return "\n" + output + "↗" + strings.Repeat("↑", internalLength) + "↖" + // bottom line
		"\n\n"
}

// AllEqual returns true iff all the provided elements are equal to each other.
func AllEqual[T comparable](elems ...T) bool {
	for i := 1; i < len(elems); i++ {
		if elems[i] != elems[0] {
			return false
		}
	}
	return true
}

// RandUint256 generates a random bigNum up to 2 ** 256 - 1
func RandUint256() *big.Int {
	n, err := rand.Int(rand.Reader, MaxUint256)
	if err != nil {
		panic(err)
	}
	return n
}

func LeftPadBitString(input string, length int) string {
	if len(input) >= length {
		return input
	}
	return strings.Repeat("0", length-len(input)) + input
}

// TryParseHex parses the given hex string to bytes,
// it can return error if the hex string is invalid.
// Follows the semantic of ethereum's FromHex.
func TryParseHex(s string) (b []byte, err error) {
	if !HasHexPrefix(s) {
		err = errors.New("hex string must have 0x prefix")
	} else {
		s = s[2:]
		if len(s)%2 == 1 {
			s = "0" + s
		}
		b, err = hex.DecodeString(s)
	}
	return
}

// MinKey returns the minimum value of the given element array with respect
// to the given key function. In the event U is not a compound type (e.g a
// struct) an identity function can be provided.
func MinKey[U any, T constraints.Ordered](elems []U, key func(U) T) T {
	var min T
	if len(elems) == 0 {
		return min
	}

	min = key(elems[0])
	for i := 1; i < len(elems); i++ {
		v := key(elems[i])
		if v < min {
			min = v
		}
	}

	return min
}

// ErrorBuffer uses joinedErrors interface to join multiple errors into a single error.
// This is useful to track the most recent N errors in a service and flush them as a single error.
type ErrorBuffer struct {
	// buffer is a slice of errors
	buffer []error

	// cap is the maximum number of errors that the buffer can hold.
	// Exceeding the cap results in discarding the oldest error
	cap int

	mu sync.RWMutex
}

func (eb *ErrorBuffer) Flush() (err error) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	err = errors.Join(eb.buffer...)
	eb.buffer = nil
	return
}

func (eb *ErrorBuffer) Append(incoming error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if len(eb.buffer) == eb.cap && eb.cap != 0 {
		eb.buffer = append(eb.buffer[1:], incoming)
		return
	}
	eb.buffer = append(eb.buffer, incoming)
}

func (eb *ErrorBuffer) SetCap(cap int) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if len(eb.buffer) > cap {
		eb.buffer = eb.buffer[len(eb.buffer)-cap:]
	}
	eb.cap = cap
}

// UnwrapError returns a list of underlying errors if passed error implements joinedError or return the err in a single-element list otherwise.
//
//nolint:errorlint // error type checks will fail on wrapped errors. Disabled since we are not doing checks on error types.
func UnwrapError(err error) []error {
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return []error{err}
	}
	return joined.Unwrap()
}

// DeleteUnstable destructively removes slice element at index i
// It does no bounds checking and may re-order the slice
func DeleteUnstable[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	return s
}
