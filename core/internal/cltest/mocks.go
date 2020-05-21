package cltest

import (
	"context"
	"encoding"
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

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LenientEthMock flag prevents the mock eth client from panicking if an unexpected call is made
const LenientEthMock = "lenient"

// EthMockRegisterChainID registers the common case of calling eth_chainId
// and returns the store.config.ChainID
const EthMockRegisterChainID = "eth_mock_register_chain_id"

// MockCallerSubscriberClient create new EthMock Client
func (ta *TestApplication) MockCallerSubscriberClient(flags ...string) *EthMock {
	if ta.ChainlinkApplication.HeadTracker.Connected() {
		logger.Panic("Cannot mock eth client after being connected")
	}
	return MockEthOnStore(ta.t, ta.Store, flags...)
}

// MockEthOnStore given store return new EthMock Client
func MockEthOnStore(t testing.TB, s *store.Store, flags ...string) *EthMock {
	mock := &EthMock{t: t, strict: true}
	for _, flag := range flags {
		if flag == LenientEthMock {
			mock.strict = false
		} else if flag == EthMockRegisterChainID {
			mock.Register("eth_chainId", s.Config.ChainID())
		}
	}
	eth := &eth.CallerSubscriberClient{CallerSubscriber: mock}
	if txm, ok := s.TxManager.(*store.EthTxManager); ok {
		txm.Client = eth
	} else {
		log.Panic("MockEthOnStore only works on EthTxManager")
	}
	return mock
}

// SimpleGethWrapper offers an easy way to mock the eth client
type SimpleGethWrapper struct {
	c eth.GethClient
}

func NewSimpleGethWrapper(c eth.GethClient) *SimpleGethWrapper {
	wrapper := SimpleGethWrapper{c: c}
	return &wrapper
}

func (wrapper *SimpleGethWrapper) GethClient(f func(c eth.GethClient) error) error {
	return f(wrapper.c)
}

// EthMock is a mock ethereum client
type EthMock struct {
	Responses      []MockResponse
	Subscriptions  []*MockSubscription
	newHeadsCalled bool
	logsCalled     bool
	mutex          sync.RWMutex
	context        string
	strict         bool
	t              testing.TB
}

// GethClient is a noop, solely needed to conform to GethClientWrapper interface
func (mock *EthMock) GethClient(f func(c eth.GethClient) error) error {
	return nil
}

// Dial mock dial
func (mock *EthMock) Dial(url string) (eth.CallerSubscriber, error) {
	return mock, nil
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

			if err := assignResult(result, resp.response); err != nil {
				return err
			}

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
	}
	return err
}

// assignResult attempts to mimick more closely how go-ethereum actually does
// Call, falling back to reflection if the values dont support the required
// encoding interfaces
func assignResult(result, response interface{}) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			switch perr := perr.(type) {
			case string:
				err = errors.New(perr)
			case error:
				err = perr
			}
		}
	}()
	if unmarshaler, ok := result.(encoding.TextUnmarshaler); ok {
		switch resp := response.(type) {
		case encoding.TextMarshaler:
			bytes, err := resp.MarshalText()
			if err != nil {
				return err
			}
			return unmarshaler.UnmarshalText(bytes)
		case string:
			return unmarshaler.UnmarshalText([]byte(resp))
		case []byte:
			return unmarshaler.UnmarshalText(resp)
		}
	}

	ref := reflect.ValueOf(result)
	reflect.Indirect(ref).Set(reflect.ValueOf(response))
	return nil
}

// RegisterSubscription register a mock subscription to the given name and channels
func (mock *EthMock) RegisterSubscription(name string, channels ...interface{}) *MockSubscription {
	var channel interface{}
	if len(channels) > 0 {
		channel = channels[0]
	} else {
		channel = channelFromSubscriptionName(name)
	}

	sub := &MockSubscription{
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
		return make(chan eth.Log)
	case "newHeads":
		return make(chan eth.BlockHeader)
	default:
		return make(chan struct{})
	}
}

// Subscribe registers a subscription to the channel
func (mock *EthMock) Subscribe(
	ctx context.Context,
	channel interface{},
	args ...interface{},
) (eth.Subscription, error) {
	mock.mutex.Lock()
	defer mock.mutex.Unlock()
	for i, sub := range mock.Subscriptions {
		if sub.name == args[0] {
			mock.Subscriptions = append(mock.Subscriptions[:i], mock.Subscriptions[i+1:]...)
			switch channel.(type) {
			case chan<- eth.Log:
				fwdLogs(channel, sub.channel)
			case chan<- eth.BlockHeader:
				fwdHeaders(channel, sub.channel)
			default:
				return nil, errors.New("channel type not supported by ethMock")
			}
			return sub, nil
		}
	}
	if args[0] == "newHeads" && !mock.newHeadsCalled {
		mock.newHeadsCalled = true
		return EmptyMockSubscription(), nil
	} else if args[0] == "logs" && !mock.logsCalled {
		mock.logsCalled = true
		return &MockSubscription{
			channel: make(chan eth.Log),
			Errors:  make(chan error),
		}, nil
	} else if args[0] == "newHeads" {
		return nil, errors.New("newHeads subscription only expected once, please register another mock subscription if more are needed")
	}
	return nil, errors.New("must RegisterSubscription before Subscribe")
}

// RegisterNewHeads registers a newheads subscription
func (mock *EthMock) RegisterNewHeads() chan eth.BlockHeader {
	newHeads := make(chan eth.BlockHeader, 10)
	mock.RegisterSubscription("newHeads", newHeads)
	return newHeads
}

// RegisterNewHead register new head at given blocknumber
func (mock *EthMock) RegisterNewHead(blockNumber int64) chan eth.BlockHeader {
	newHeads := mock.RegisterNewHeads()
	newHeads <- eth.BlockHeader{Number: BigHexInt(blockNumber)}
	return newHeads
}

func fwdLogs(actual, mock interface{}) {
	logChan := actual.(chan<- eth.Log)
	mockChan := mock.(chan eth.Log)
	go func() {
		for e := range mockChan {
			logChan <- e
		}
	}()
}

func fwdHeaders(actual, mock interface{}) {
	logChan := actual.(chan<- eth.BlockHeader)
	mockChan := mock.(chan eth.BlockHeader)
	go func() {
		for e := range mockChan {
			logChan <- e
		}
	}()
}

// MockSubscription a mock subscription
type MockSubscription struct {
	mut          sync.Mutex
	name         string
	channel      interface{}
	unsubscribed bool
	Errors       chan error
}

// EmptyMockSubscription return empty MockSubscription
func EmptyMockSubscription() *MockSubscription {
	return &MockSubscription{Errors: make(chan error, 1), channel: make(chan struct{})}
}

// Err returns error channel from mes
func (mes *MockSubscription) Err() <-chan error { return mes.Errors }

// Unsubscribe closes the subscription
func (mes *MockSubscription) Unsubscribe() {
	mes.mut.Lock()
	defer mes.mut.Unlock()

	if mes.unsubscribed {
		return
	}
	mes.unsubscribed = true
	switch mes.channel.(type) {
	case chan struct{}:
		close(mes.channel.(chan struct{}))
	case chan eth.Log:
		close(mes.channel.(chan eth.Log))
	case chan eth.BlockHeader:
		close(mes.channel.(chan eth.BlockHeader))
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

// InstantClock an InstantClock
type InstantClock struct{}

// Now returns the current local time
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
	t        testing.TB
}

// NewTriggerClock returns a new TriggerClock, that a test can manually fire
// to continue processing in a Clock dependency.
func NewTriggerClock(t testing.TB) *TriggerClock {
	return &TriggerClock{
		triggers: make(chan time.Time),
		t:        t,
	}
}

// Trigger sends a time to unblock the After call.
func (t *TriggerClock) Trigger() {
	select {
	case t.triggers <- time.Now():
	case <-time.After(60 * time.Second):
		t.t.Error("timed out while trying to trigger clock")
	}
}

// TriggerWithoutTimeout is a special case where we know the trigger might
// block but don't care
func (t *TriggerClock) TriggerWithoutTimeout() {
	t.triggers <- time.Now()
}

// Now returns the current local time
func (t TriggerClock) Now() time.Time {
	return time.Now()
}

// After waits on a manual trigger.
func (t *TriggerClock) After(_ time.Duration) <-chan time.Time {
	return t.triggers
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
	App chainlink.Application
}

// NewApplication creates a new application with specified config
func (f InstanceAppFactory) NewApplication(config *orm.Config, onConnectCallbacks ...func(chainlink.Application)) chainlink.Application {
	return f.App
}

type seededAppFactory struct {
	Application chainlink.Application
}

func (s seededAppFactory) NewApplication(config *orm.Config, onConnectCallbacks ...func(chainlink.Application)) chainlink.Application {
	return noopStopApplication{s.Application}
}

type noopStopApplication struct {
	chainlink.Application
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

func (a CallbackAuthenticator) AuthenticateVRFKey(*store.Store, string) error {
	return nil
}

var _ cmd.KeyStoreAuthenticator = CallbackAuthenticator{}

// BlockedRunner is a Runner that blocks until its channel is posted to
type BlockedRunner struct {
	Done chan struct{}
}

// Run runs the blocked runner, doesn't return until the channel is signalled
func (r BlockedRunner) Run(app chainlink.Application) error {
	<-r.Done
	return nil
}

// EmptyRunner is an EmptyRunner
type EmptyRunner struct{}

// Run runs the empty runner
func (r EmptyRunner) Run(app chainlink.Application) error {
	return nil
}

// MockCountingPrompter is a mock counting prompt
type MockCountingPrompter struct {
	T              *testing.T
	EnteredStrings []string
	Count          int
	NotTerminal    bool
}

// Prompt returns an entered string
func (p *MockCountingPrompter) Prompt(string) string {
	i := p.Count
	p.Count++
	if len(p.EnteredStrings)-1 < i {
		p.T.Errorf("Not enough passwords supplied to MockCountingPrompter, wanted %d", i)
		p.T.FailNow()
	}
	return p.EnteredStrings[i]
}

// PasswordPrompt returns an entered string
func (p *MockCountingPrompter) PasswordPrompt(string) string {
	i := p.Count
	p.Count++
	if len(p.EnteredStrings)-1 < i {
		p.T.Errorf("Not enough passwords supplied to MockCountingPrompter, wanted %d", i)
		p.T.FailNow()
	}
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
		assert.True(t, called, "expected call Mock HTTP endpoint '%s'", server.URL)
	}
}

func NewHTTPMockServerWithAlterableResponse(
	t *testing.T, response func() string) (server *httptest.Server) {
	server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response())
		}))
	return server
}

// MockCron represents a mock cron
type MockCron struct {
	Entries []MockCronEntry
	nextID  cron.EntryID
}

// NewMockCron returns a new mock cron
func NewMockCron() *MockCron {
	return &MockCron{}
}

// Start starts the mockcron
func (*MockCron) Start() {}

// Stop stops the mockcron
func (*MockCron) Stop() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

// AddFunc appends a schedule to mockcron entries
func (mc *MockCron) AddFunc(schd string, fn func()) (cron.EntryID, error) {
	mc.Entries = append(mc.Entries, MockCronEntry{
		Schedule: schd,
		Function: fn,
	})
	mc.nextID++
	return mc.nextID, nil
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

func MustRandomUser() models.User {
	email := fmt.Sprintf("user-%v@chainlink.test", NewRandomInt64())
	r, err := models.NewUser(email, Password)
	if err != nil {
		logger.Panic(err)
	}
	return r
}

func MustNewUser(t *testing.T, email, password string) models.User {
	r, err := models.NewUser(email, password)
	if err != nil {
		t.Fatal(err)
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
	user := MustRandomUser()
	return user, store.SaveUser(&user)
}

func NewMockAuthenticatedHTTPClient(cfg orm.ConfigReader, sessionID string) cmd.HTTPClient {
	return cmd.NewAuthenticatedHTTPClient(cfg, MockCookieAuthenticator{SessionID: sessionID})
}

type MockCookieAuthenticator struct {
	SessionID string
	Error     error
}

func (m MockCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return MustGenerateSessionCookie(m.SessionID), m.Error
}

func (m MockCookieAuthenticator) Authenticate(models.SessionRequest) (*http.Cookie, error) {
	return MustGenerateSessionCookie(m.SessionID), m.Error
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

func (m mockSecretGenerator) Generate(orm.Config) ([]byte, error) {
	return []byte(SessionSecret), nil
}

// extractERC20BalanceTargetAddress returns the address whose balance is being
// queried by the message in the given call to an ERC20 contract, which is
// interpreted as a callArgs.
func extractERC20BalanceTargetAddress(args interface{}) (common.Address, bool) {
	call, ok := (args).(eth.CallArgs)
	if !ok {
		return common.Address{}, false
	}
	message := call.Data
	return common.BytesToAddress(([]byte)(message)[len(message)-20:]), true
}

// ExtractTargetAddressFromERC20EthEthCallMock extracts the contract address and the
// method data, for checking in a test.
func ExtractTargetAddressFromERC20EthEthCallMock(
	t *testing.T, arg ...interface{}) common.Address {
	ethMockCallArgs, ethMockCallArgsOk := (arg[0]).([]interface{})
	require.True(t, ethMockCallArgsOk)
	actualCallArgs, actualCallArgsOk := (ethMockCallArgs[0]).([]interface{})
	require.True(t, actualCallArgsOk)
	address, ok := extractERC20BalanceTargetAddress(actualCallArgs[0])
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
