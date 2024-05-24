// Package utils is used for common functions and tools used across the codebase.
package utils

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	mrand "math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jpillora/backoff"
	pkgerrors "github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// DefaultSecretSize is the entropy in bytes to generate a base64 string of 64 characters.
const DefaultSecretSize = 48

func MustNewPeerID() string {
	pubKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	peerID, err := ragep2ptypes.PeerIDFromPublicKey(pubKey)
	if err != nil {
		panic(err)
	}
	return peerID.String()
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

// Sha256 returns a hexadecimal encoded string of a hashed input
func Sha256(in string) (string, error) {
	hasher := sha3.New256()
	_, err := hasher.Write([]byte(in))
	if err != nil {
		return "", pkgerrors.Wrap(err, "sha256 write error")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// WithCloseChan wraps a context so that it is canceled if the passed in channel is closed.
// Deprecated: Call [services.StopChan.Ctx] directly
func WithCloseChan(parentCtx context.Context, chStop chan struct{}) (context.Context, context.CancelFunc) {
	return services.StopChan(chStop).Ctx(parentCtx)
}

// ContextFromChan creates a context that finishes when the provided channel receives or is closed.
// Deprecated: Call [services.StopChan.NewCtx] directly.
func ContextFromChan(chStop chan struct{}) (context.Context, context.CancelFunc) {
	return services.StopChan(chStop).NewCtx()
}

// ContextFromChanWithTimeout creates a context with a timeout that finishes when the provided channel receives or is closed.
// Deprecated: Call [services.StopChan.CtxCancel] directly
func ContextFromChanWithTimeout(chStop chan struct{}, timeout time.Duration) (context.Context, context.CancelFunc) {
	return services.StopChan(chStop).CtxCancel(context.WithTimeout(context.Background(), timeout))
}

// Deprecated: use services.StopChan
type StopChan = services.StopChan

// Deprecated: use services.StopRChan
type StopRChan = services.StopRChan

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

var (
	ErrAlreadyStopped      = errors.New("already stopped")
	ErrCannotStopUnstarted = errors.New("cannot stop unstarted service")
)

// StartStopOnce contains a StartStopOnceState integer
// Deprecated: use services.StateMachine
type StartStopOnce = services.StateMachine

// WithJitter adds +/- 10% to a duration
func WithJitter(d time.Duration) time.Duration {
	// #nosec
	if d == 0 {
		return 0
	}
	// ensure non-zero arg to Intn to avoid panic
	max := math.Max(float64(d.Abs())/5.0, 1.)
	// #nosec - non critical randomness
	jitter := mrand.Intn(int(max))
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

// ConcatBytes appends a bunch of byte arrays into a single byte array
func ConcatBytes(bufs ...[]byte) []byte {
	return bytes.Join(bufs, []byte{})
}

func LeftPadBitString(input string, length int) string {
	if len(input) >= length {
		return input
	}
	return strings.Repeat("0", length-len(input)) + input
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
