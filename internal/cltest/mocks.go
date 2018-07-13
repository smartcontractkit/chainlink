package cltest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

// MockEthClient create new EthMock Client
func (ta *TestApplication) MockEthClient() *EthMock {
	return MockEthOnStore(ta.Store)
}

// MockEthOnStore given store return new EthMock Client
func MockEthOnStore(s *store.Store) *EthMock {
	mock := &EthMock{}
	eth := &store.EthClient{CallerSubscriber: mock}
	s.TxManager.EthClient = eth
	return mock
}

// EthMock is a mock etheruem client
type EthMock struct {
	Responses      []MockResponse
	Subscriptions  []MockSubscription
	newHeadsCalled bool
	logsCalled     bool
	mutex          sync.RWMutex
	context        string
}

// Dial mock dial
func (mock *EthMock) Dial(url string) (store.CallerSubscriber, error) {
	return mock, nil
}

// Context adds helpful context to EthMock values set in the callback function.
func (mock *EthMock) Context(context string, callback func(*EthMock)) {
	mock.context = context
	callback(mock)
	mock.context = ""
}

// Register register mock responses and append to Ethmock
func (mock *EthMock) Register(
	method string,
	response interface{},
	callback ...func(interface{}, ...interface{}) error,
) {
	res := MockResponse{
		methodName: method,
		response:   response,
		context:    mock.context,
	}
	if len(callback) > 0 {
		res.callback = callback[0]
	}

	mock.mutex.Lock()
	defer mock.mutex.Unlock()
	mock.Responses = append(mock.Responses, res)
}

// RegisterError register mock errors to EthMock
func (mock *EthMock) RegisterError(method, errMsg string) {
	res := MockResponse{
		methodName: method,
		errMsg:     errMsg,
		hasError:   true,
		context:    mock.context,
	}

	mock.mutex.Lock()
	defer mock.mutex.Unlock()
	mock.Responses = append(mock.Responses, res)
}

// AllCalled return true if all mocks have been mocked
func (mock *EthMock) AllCalled() bool {
	mock.mutex.RLock()
	defer mock.mutex.RUnlock()
	return (len(mock.Responses) == 0) && (len(mock.Subscriptions) == 0)
}

// EventuallyAllCalled eventually will return after all the mock subscriptions and responses are called
func (mock *EthMock) EventuallyAllCalled(t *testing.T) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	g.Eventually(mock.AllCalled).Should(gomega.BeTrue())
}

// Call will call given method and set the result
func (mock *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	mock.mutex.Lock()
	defer mock.mutex.Unlock()

	for i, resp := range mock.Responses {
		if resp.methodName == method {
			mock.Responses = append(mock.Responses[:i], mock.Responses[i+1:]...)
			if resp.hasError {
				return fmt.Errorf(resp.errMsg)
			}
			ref := reflect.ValueOf(result)
			reflect.Indirect(ref).Set(reflect.ValueOf(resp.response))
			if resp.callback != nil {
				if err := resp.callback(result, args); err != nil {
					return fmt.Errorf("ethMock Error: %v\ncontext: %v", err, resp.context)
				}
			}
			return nil
		}
	}
	return fmt.Errorf("EthMock: Method %v not registered", method)
}

// RegisterSubscription register a mock subscription to the given name and channels
func (mock *EthMock) RegisterSubscription(name string, channels ...interface{}) MockSubscription {
	var channel interface{}
	if len(channels) > 0 {
		channel = channels[0]
	} else {
		channel = channelFromSubscriptionName(name)
	}

	sub := MockSubscription{
		name:    name,
		channel: channel,
		Errors:  make(chan error, 1),
	}
	mock.mutex.Lock()
	defer mock.mutex.Unlock()
	mock.Subscriptions = append(mock.Subscriptions, sub)
	return sub
}

func channelFromSubscriptionName(name string) interface{} {
	switch name {
	case "logs":
		return make(chan types.Log)
	case "newHeads":
		return make(chan models.BlockHeader)
	default:
		return make(chan struct{})
	}
}

// EthSubscribe registers a subscription to the channel
func (mock *EthMock) EthSubscribe(
	ctx context.Context,
	channel interface{},
	args ...interface{},
) (models.EthSubscription, error) {
	mock.mutex.Lock()
	defer mock.mutex.Unlock()
	for i, sub := range mock.Subscriptions {
		if sub.name == args[0] {
			mock.Subscriptions = append(mock.Subscriptions[:i], mock.Subscriptions[i+1:]...)
			switch channel.(type) {
			case chan<- types.Log:
				fwdLogs(channel, sub.channel)
			case chan<- models.BlockHeader:
				fwdHeaders(channel, sub.channel)
			default:
				return nil, errors.New("Channel type not supported by ethMock")
			}
			return sub, nil
		}
	}
	if args[0] == "newHeads" && !mock.newHeadsCalled {
		mock.newHeadsCalled = true
		return EmptyMockSubscription(), nil
	} else if args[0] == "logs" && !mock.logsCalled {
		mock.logsCalled = true
		return MockSubscription{
			channel: make(chan types.Log),
			Errors:  make(chan error),
		}, nil
	} else if args[0] == "newHeads" {
		return nil, errors.New("newHeads subscription only expected once, please register another mock subscription if more are needed")
	}
	return nil, errors.New("Must RegisterSubscription before EthSubscribe")
}

// RegisterNewHeads registers a newheads subscription
func (mock *EthMock) RegisterNewHeads() chan models.BlockHeader {
	newHeads := make(chan models.BlockHeader, 10)
	mock.RegisterSubscription("newHeads", newHeads)
	return newHeads
}

// RegisterNewHead register new head at given blocknumber
func (mock *EthMock) RegisterNewHead(blockNumber int64) chan models.BlockHeader {
	newHeads := mock.RegisterNewHeads()
	newHeads <- models.BlockHeader{Number: BigHexInt(blockNumber)}
	return newHeads
}

func fwdLogs(actual, mock interface{}) {
	logChan := actual.(chan<- types.Log)
	mockChan := mock.(chan types.Log)
	go func() {
		for e := range mockChan {
			logChan <- e
		}
	}()
}

func fwdHeaders(actual, mock interface{}) {
	logChan := actual.(chan<- models.BlockHeader)
	mockChan := mock.(chan models.BlockHeader)
	go func() {
		for e := range mockChan {
			logChan <- e
		}
	}()
}

// MockSubscription a mock subscription
type MockSubscription struct {
	name    string
	channel interface{}
	Errors  chan error
}

// EmptyMockSubscription return empty MockSubscription
func EmptyMockSubscription() MockSubscription {
	return MockSubscription{Errors: make(chan error, 1), channel: make(chan struct{})}
}

// Err returns error channel from mes
func (mes MockSubscription) Err() <-chan error { return mes.Errors }

// Unsubscribe closes the subscription
func (mes MockSubscription) Unsubscribe() {
	switch mes.channel.(type) {
	case chan struct{}:
		close(mes.channel.(chan struct{}))
	case chan types.Log:
		close(mes.channel.(chan types.Log))
	case chan models.BlockHeader:
		close(mes.channel.(chan models.BlockHeader))
	default:
		logger.Fatal(fmt.Sprintf("Unable to close MockSubscription channel of type %T", mes.channel))
	}
	close(mes.Errors)
}

// MockResponse a mock response
type MockResponse struct {
	methodName string
	context    string
	response   interface{}
	errMsg     string
	hasError   bool
	callback   func(interface{}, ...interface{}) error
}

// InstantClock create InstantClock
func (ta *TestApplication) InstantClock() InstantClock {
	clock := InstantClock{}
	ta.Scheduler.OneTime.Clock = clock
	return clock
}

// UseSettableClock creates a SettableClock on the store
func UseSettableClock(s *store.Store) *SettableClock {
	clock := &SettableClock{}
	s.Clock = clock
	return clock
}

// SettableClock a settable clock
type SettableClock struct {
	mutex sync.Mutex
	time  time.Time
}

// Now get the current time
func (clock *SettableClock) Now() time.Time {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	if clock.time.IsZero() {
		return time.Now()
	}
	return clock.time
}

// SetTime set the current time
func (clock *SettableClock) SetTime(t time.Time) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.time = t
}

// After return channel of time
func (*SettableClock) After(_ time.Duration) <-chan time.Time {
	channel := make(chan time.Time, 1)
	channel <- time.Now()
	return channel
}

// InstantClock an InstantClock
type InstantClock struct{}

// Now current local time
func (InstantClock) Now() time.Time {
	return time.Now()
}

// After return channel of time
func (InstantClock) After(_ time.Duration) <-chan time.Time {
	c := make(chan time.Time, 100)
	c <- time.Now()
	return c
}

// NeverClock a never clock
type NeverClock struct{}

// After return channel of time
func (NeverClock) After(_ time.Duration) <-chan time.Time {
	return make(chan time.Time)
}

// Now returns current local time
func (NeverClock) Now() time.Time {
	return time.Now()
}

// RendererMock a mock renderer
type RendererMock struct {
	Renders []interface{}
}

// Render appends values to renderer mock
func (rm *RendererMock) Render(v interface{}) error {
	rm.Renders = append(rm.Renders, v)
	return nil
}

// InstanceAppFactory is an InstanceAppFactory
type InstanceAppFactory struct {
	App services.Application
}

// NewApplication creates a new application with specified config
func (f InstanceAppFactory) NewApplication(config store.Config) services.Application {
	return f.App
}

// EmptyAppFactory an empty application factory
type EmptyAppFactory struct{}

// NewApplication creates a new empty application with specified config
func (f EmptyAppFactory) NewApplication(config store.Config) services.Application {
	return &EmptyApplication{}
}

// EmptyApplication an empty application
type EmptyApplication struct{}

// Start starts the empty application
func (a *EmptyApplication) Start() error {
	return nil
}

// Stop stopts the empty application
func (a *EmptyApplication) Stop() error {
	return nil
}

// GetStore retrieves the store of the empty application
func (a *EmptyApplication) GetStore() *store.Store {
	return nil
}

// CallbackAuthenticator contains a call back authenticator method
type CallbackAuthenticator struct {
	Callback func(*store.Store, string) error
}

// Authenticate authenticates store and pwd with the callback authenticator
func (a CallbackAuthenticator) Authenticate(store *store.Store, pwd string) error {
	return a.Callback(store, pwd)
}

// EmptyRunner is an EmptyRunner
type EmptyRunner struct{}

// Run runs the empty runner
func (r EmptyRunner) Run(app services.Application) error {
	return nil
}

// MockCountingPrompt is a mock counting prompt
type MockCountingPrompt struct {
	EnteredStrings []string
	Count          int
}

// Prompt returns an entered string
func (p *MockCountingPrompt) Prompt(string) string {
	i := p.Count
	p.Count++
	return p.EnteredStrings[i]
}

// PasswordPrompt returns an entered string
func (p *MockCountingPrompt) PasswordPrompt(string) string {
	i := p.Count
	p.Count++
	return p.EnteredStrings[i]
}

// IsTerminal always returns true in tests
func (p *MockCountingPrompt) IsTerminal() bool {
	return true
}

type MockUserInitializer struct{}

func (m MockUserInitializer) Prompt(string) string         { return "" }
func (m MockUserInitializer) PasswordPrompt(string) string { return "" }
func (m MockUserInitializer) IsTerminal() bool             { return true }

// NewHTTPMockServer create http test server with passed in parameters
func NewHTTPMockServer(
	t *testing.T,
	status int,
	wantMethod string,
	response string,
	callback ...func(string),
) (*httptest.Server, func()) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, wantMethod, r.Method)
		if len(callback) > 0 {
			callback[0](string(b))
		}
		called = true

		w.WriteHeader(status)
		io.WriteString(w, response)
	})

	server := httptest.NewServer(handler)
	return server, func() {
		server.Close()
		assert.True(t, called)
	}
}

// MockCron represents a mock cron
type MockCron struct {
	Entries []MockCronEntry
}

// NewMockCron returns a new mock cron
func NewMockCron() *MockCron {
	return &MockCron{}
}

// Start starts the mockcron
func (*MockCron) Start() {}

// Stop stops the mockcron
func (*MockCron) Stop() {}

// AddFunc appends a schedule to mockcron entries
func (mc *MockCron) AddFunc(schd string, fn func()) error {
	mc.Entries = append(mc.Entries, MockCronEntry{
		Schedule: schd,
		Function: fn,
	})
	return nil
}

// RunEntries run every function for each mockcron entry
func (mc *MockCron) RunEntries() {
	for _, entry := range mc.Entries {
		entry.Function()
	}
}

// MockCronEntry a cron schedule and function
type MockCronEntry struct {
	Schedule string
	Function func()
}

// MockHeadTrackable allows you to mock HeadTrackable
type MockHeadTrackable struct {
	connectedCount    int32
	disconnectedCount int32
	onNewHeadCount    int32
}

// Connect increases the connected count by one
func (m *MockHeadTrackable) Connect(*models.IndexableBlockNumber) error {
	atomic.AddInt32(&m.connectedCount, 1)
	return nil
}

// ConnectedCount returns the count of connections made, safely.
func (m *MockHeadTrackable) ConnectedCount() int32 {
	return atomic.LoadInt32(&m.connectedCount)
}

// Disconnect increases the disconnected count by one
func (m *MockHeadTrackable) Disconnect() { atomic.AddInt32(&m.disconnectedCount, 1) }

// DisconnectedCount returns the count of disconnections made, safely.
func (m *MockHeadTrackable) DisconnectedCount() int32 {
	return atomic.LoadInt32(&m.disconnectedCount)
}

// OnNewHead increases the OnNewHeadCount count by one
func (m *MockHeadTrackable) OnNewHead(*models.BlockHeader) { atomic.AddInt32(&m.onNewHeadCount, 1) }

// OnNewHeadCount returns the count of new heads, safely.
func (m *MockHeadTrackable) OnNewHeadCount() int32 {
	return atomic.LoadInt32(&m.onNewHeadCount)
}

// NeverSleeper is a struct that never sleeps
type NeverSleeper struct{}

// Reset resets the never sleeper
func (ns NeverSleeper) Reset() {}

// Sleep puts the never sleeper to sleep
func (ns NeverSleeper) Sleep() {}

// Duration returns a duration
func (ns NeverSleeper) Duration() time.Duration { return 0 * time.Microsecond }
