package cltest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/sasha-s/go-deadlock"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Strict flag makes the mock eth client panic if an unexpected call is made
const Strict = "strict"

// MockEthClient create new EthMock Client
func (ta *TestApplication) MockEthClient(flags ...string) *EthMock {
	if ta.ChainlinkApplication.HeadTracker.Connected() {
		logger.Panic("Cannot mock eth client after being connected")
	}
	return MockEthOnStore(ta.t, ta.Store, flags...)
}

// MockEthOnStore given store return new EthMock Client
func MockEthOnStore(t testing.TB, s *store.Store, flags ...string) *EthMock {
	mock := &EthMock{t: t}
	for _, flag := range flags {
		if flag == Strict {
			mock.strict = true
		}
	}
	eth := &store.EthClient{CallerSubscriber: mock}
	if txm, ok := s.TxManager.(*store.EthTxManager); ok {
		txm.EthClient = eth
	} else {
		log.Panic("MockEthOnStore only works on EthTxManager")
	}
	return mock
}

// EthMock is a mock ethereum client
type EthMock struct {
	Responses      []MockResponse
	Subscriptions  []MockSubscription
	newHeadsCalled bool
	logsCalled     bool
	mutex          sync.RWMutex
	context        string
	strict         bool
	t              testing.TB
}

// Dial mock dial
func (mock *EthMock) Dial(url string) (store.CallerSubscriber, error) {
	return mock, nil
}

// Clear all stubs/mocks/expectations
func (mock *EthMock) Clear() {
	mock.mutex.Lock()
	defer mock.mutex.Unlock()

	mock.Responses = nil
	mock.Subscriptions = nil
	mock.newHeadsCalled = false
	mock.logsCalled = false
}

// Context adds helpful context to EthMock values set in the callback function.
func (mock *EthMock) Context(context string, callback func(*EthMock)) {
	mock.context = context
	callback(mock)
	mock.context = ""
}

func (mock *EthMock) ShouldCall(setup func(mock *EthMock)) ethMockDuring {
	if !mock.AllCalled() {
		mock.t.Errorf("Remaining ethMockCalls: %v", mock.Remaining())
	}
	setup(mock)
	return ethMockDuring{mock: mock}
}

type ethMockDuring struct {
	mock *EthMock
	t    testing.TB
}

func (emd ethMockDuring) During(action func()) {
	action()
	if !emd.mock.AllCalled() {
		emd.mock.t.Errorf("Remaining ethMockCalls: %v", emd.mock.Remaining())
	}
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

func (mock *EthMock) Remaining() string {
	mock.mutex.RLock()
	defer mock.mutex.RUnlock()
	rvals := []string{}
	for _, r := range mock.Responses {
		rvals = append(rvals, fmt.Sprintf("Response %s#%s not called", r.context, r.methodName))
	}
	for _, s := range mock.Subscriptions {
		rvals = append(rvals, fmt.Sprintf("Subscription %s not called", s.name))
	}
	return strings.Join(rvals, ",")
}

// EventuallyAllCalled eventually will return after all the mock subscriptions and responses are called
func (mock *EthMock) EventuallyAllCalled(t *testing.T) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	g.Eventually(mock.Remaining).Should(gomega.HaveLen(0))
}

// AssertAllCalled immediately checks that all calls have been made
func (mock *EthMock) AssertAllCalled() {
	assert.Empty(mock.t, mock.Remaining())
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

	err := fmt.Errorf("EthMock: Method %v not registered", method)
	if mock.strict {
		mock.t.Errorf("%s\n%s", err, debug.Stack())
	} else {
		mock.t.Logf("%s\n%s", err, debug.Stack())
	}
	return err
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
		return make(chan models.Log)
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
			case chan<- models.Log:
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
			channel: make(chan models.Log),
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

func (mock *EthMock) NoMagic() {
	mock.newHeadsCalled = true
}

func fwdLogs(actual, mock interface{}) {
	logChan := actual.(chan<- models.Log)
	mockChan := mock.(chan models.Log)
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
	case chan models.Log:
		close(mes.channel.(chan models.Log))
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
	mutex deadlock.Mutex
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

// TriggerClock implements the AfterNower interface, but must be manually triggered
// to resume computation on After.
type TriggerClock struct {
	triggers chan time.Time
}

// NewTriggerClock returns a new TriggerClock, that a test can manually fire
// to continue processing in a Clock dependency.
func NewTriggerClock() *TriggerClock {
	return &TriggerClock{
		triggers: make(chan time.Time),
	}
}

// Trigger sends a time to unblock the After call.
func (t *TriggerClock) Trigger() {
	t.triggers <- time.Now()
}

// After waits on a manual trigger.
func (t *TriggerClock) After(_ time.Duration) <-chan time.Time {
	return t.triggers
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
func (f InstanceAppFactory) NewApplication(config store.Config, onConnectCallbacks ...func(services.Application)) services.Application {
	return f.App
}

type seededAppFactory struct {
	Application services.Application
}

func (s seededAppFactory) NewApplication(config store.Config, onConnectCallbacks ...func(services.Application)) services.Application {
	return noopStopApplication{s.Application}
}

type noopStopApplication struct {
	services.Application
}

func (a noopStopApplication) Stop() error {
	return nil
}

// CallbackAuthenticator contains a call back authenticator method
type CallbackAuthenticator struct {
	Callback func(*store.Store, string) (string, error)
}

// Authenticate authenticates store and pwd with the callback authenticator
func (a CallbackAuthenticator) Authenticate(store *store.Store, pwd string) (string, error) {
	return a.Callback(store, pwd)
}

// BlockedRunner is a Runner that blocks until its channel is posted to
type BlockedRunner struct {
	Done chan struct{}
}

// Run runs the blocked runner, doesn't return until the channel is signalled
func (r BlockedRunner) Run(app services.Application) error {
	<-r.Done
	return nil
}

// EmptyRunner is an EmptyRunner
type EmptyRunner struct{}

// Run runs the empty runner
func (r EmptyRunner) Run(app services.Application) error {
	return nil
}

// MockCountingPrompter is a mock counting prompt
type MockCountingPrompter struct {
	EnteredStrings []string
	Count          int
	NotTerminal    bool
}

// Prompt returns an entered string
func (p *MockCountingPrompter) Prompt(string) string {
	i := p.Count
	p.Count++
	return p.EnteredStrings[i]
}

// PasswordPrompt returns an entered string
func (p *MockCountingPrompter) PasswordPrompt(string) string {
	i := p.Count
	p.Count++
	return p.EnteredStrings[i]
}

// IsTerminal always returns true in tests
func (p *MockCountingPrompter) IsTerminal() bool {
	return !p.NotTerminal
}

// NewHTTPMockServer create http test server with passed in parameters
func NewHTTPMockServer(
	t *testing.T,
	status int,
	wantMethod string,
	response string,
	callback ...func(http.Header, string),
) (*httptest.Server, func()) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, wantMethod, r.Method)
		if len(callback) > 0 {
			callback[0](r.Header, string(b))
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
	ConnectedCallback func(bn *models.Head)
	disconnectedCount int32
	onNewHeadCount    int32
}

// Connect increases the connected count by one
func (m *MockHeadTrackable) Connect(bn *models.Head) error {
	atomic.AddInt32(&m.connectedCount, 1)
	if m.ConnectedCallback != nil {
		m.ConnectedCallback(bn)
	}
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
func (m *MockHeadTrackable) OnNewHead(*models.Head) { atomic.AddInt32(&m.onNewHeadCount, 1) }

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

// After returns a duration
func (ns NeverSleeper) After() time.Duration { return 0 * time.Microsecond }

// Duration returns a duration
func (ns NeverSleeper) Duration() time.Duration { return 0 * time.Microsecond }

func MustUser(email, pwd string) models.User {
	r, err := models.NewUser(email, pwd)
	if err != nil {
		logger.Panic(err)
	}
	return r
}

type MockAPIInitializer struct {
	Count int
}

func (m *MockAPIInitializer) Initialize(store *store.Store) (models.User, error) {
	if user, err := store.FindUser(); err == nil {
		return user, err
	}
	m.Count += 1
	user := MustUser(APIEmail, Password)
	return user, store.SaveUser(&user)
}

func NewMockAuthenticatedHTTPClient(cfg store.Config) cmd.HTTPClient {
	return cmd.NewAuthenticatedHTTPClient(cfg, MockCookieAuthenticator{})
}

type MockCookieAuthenticator struct {
	Error error
}

func (m MockCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return MustGenerateSessionCookie(APISessionID), m.Error
}

func (m MockCookieAuthenticator) Authenticate(models.SessionRequest) (*http.Cookie, error) {
	return MustGenerateSessionCookie(APISessionID), m.Error
}

type MockSessionRequestBuilder struct {
	Count int
	Error error
}

func (m *MockSessionRequestBuilder) Build(string) (models.SessionRequest, error) {
	m.Count += 1
	if m.Error != nil {
		return models.SessionRequest{}, m.Error
	}
	return models.SessionRequest{Email: APIEmail, Password: Password}, nil
}

type mockSecretGenerator struct{}

func (m mockSecretGenerator) Generate(store.Config) ([]byte, error) {
	return []byte(SessionSecret), nil
}

type MockRunChannel struct {
	Runs               []models.RunResult
	neverReturningChan chan store.RunRequest
}

func NewMockRunChannel() *MockRunChannel {
	return &MockRunChannel{
		neverReturningChan: make(chan store.RunRequest, 1),
	}
}

func (m *MockRunChannel) Send(jobRunID string) error {
	m.Runs = append(m.Runs, models.RunResult{})
	return nil
}

func (m *MockRunChannel) Receive() <-chan store.RunRequest {
	return m.neverReturningChan
}

func (m *MockRunChannel) Close() {}

// ExtractTargetAddressFromERC20EthEthCallMock extracts the contract address and the
// method data, for checking in a test.
func ExtractTargetAddressFromERC20EthEthCallMock(
	t *testing.T, arg ...interface{}) common.Address {
	ethMockCallArgs, ethMockCallArgsOk := (arg[0]).([]interface{})
	require.True(t, ethMockCallArgsOk)
	actualCallArgs, actualCallArgsOk := (ethMockCallArgs[0]).([]interface{})
	require.True(t, actualCallArgsOk)
	address, ok := store.ExtractERC20BalanceTargetAddress(actualCallArgs[0])
	require.True(t, ok)
	return address
}

type MockChangePasswordPrompter struct {
	models.ChangePasswordRequest
	err error
}

func (m MockChangePasswordPrompter) Prompt() (models.ChangePasswordRequest, error) {
	return m.ChangePasswordRequest, m.err
}

type MockPasswordPrompter struct {
	Password string
}

func (m MockPasswordPrompter) Prompt() string {
	return m.Password
}
