package cltest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/web"
	"go.uber.org/atomic"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

// MockSubscription a mock subscription
type MockSubscription struct {
	mut          sync.Mutex
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
	case chan gethTypes.Log:
		close(mes.channel.(chan gethTypes.Log))
	case chan *eth.Head:
		close(mes.channel.(chan *eth.Head))
	default:
		logger.Fatal(fmt.Sprintf("Unable to close MockSubscription channel of type %T", mes.channel))
	}
	close(mes.Errors)
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
func (rm *RendererMock) Render(v interface{}, headers ...string) error {
	rm.Renders = append(rm.Renders, v)
	return nil
}

// InstanceAppFactory is an InstanceAppFactory
type InstanceAppFactory struct {
	App chainlink.Application
}

// NewApplication creates a new application with specified config
func (f InstanceAppFactory) NewApplication(config config.GeneralConfig) (chainlink.Application, error) {
	return f.App, nil
}

type seededAppFactory struct {
	Application chainlink.Application
}

func (s seededAppFactory) NewApplication(config config.GeneralConfig) (chainlink.Application, error) {
	return noopStopApplication{s.Application}, nil
}

type noopStopApplication struct {
	chainlink.Application
}

// FIXME: Why bother with this wrapper?
func (a noopStopApplication) Stop() error {
	return nil
}

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
) *httptest.Server {
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
		_, _ = io.WriteString(w, response) // Assignment for errcheck. Only used in tests so we can ignore.
	})

	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
		assert.True(t, called, "expected call Mock HTTP endpoint '%s'", server.URL)
	})
	return server
}

// NewHTTPMockServerWithRequest creates http test server that makes the request
// available in the callback
func NewHTTPMockServerWithRequest(
	t *testing.T,
	status int,
	response string,
	callback func(r *http.Request),
) *httptest.Server {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback(r)
		called = true

		w.WriteHeader(status)
		_, _ = io.WriteString(w, response) // Assignment for errcheck. Only used in tests so we can ignore.
	})

	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
		assert.True(t, called, "expected call Mock HTTP endpoint '%s'", server.URL)
	})
	return server
}

func NewHTTPMockServerWithAlterableResponse(
	t *testing.T, response func() string) (server *httptest.Server) {
	server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, response())
		}))
	return server
}

func NewHTTPMockServerWithAlterableResponseAndRequest(t *testing.T, response func() string, callback func(r *http.Request)) (server *httptest.Server) {
	server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callback(r)
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, response())
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
	onNewHeadCount atomic.Int32
}

// OnNewLongestChain increases the OnNewLongestChainCount count by one
func (m *MockHeadTrackable) OnNewLongestChain(context.Context, eth.Head) {
	m.onNewHeadCount.Inc()
}

// OnNewLongestChainCount returns the count of new heads, safely.
func (m *MockHeadTrackable) OnNewLongestChainCount() int32 {
	return m.onNewHeadCount.Load()
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

func MustRandomUser() sessions.User {
	email := fmt.Sprintf("user-%v@chainlink.test", NewRandomInt64())
	r, err := sessions.NewUser(email, Password)
	if err != nil {
		logger.Panic(err)
	}
	return r
}

func MustNewUser(t *testing.T, email, password string) sessions.User {
	r, err := sessions.NewUser(email, password)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

type MockAPIInitializer struct {
	Count int
}

func (m *MockAPIInitializer) Initialize(orm sessions.ORM) (sessions.User, error) {
	if user, err := orm.FindUser(); err == nil {
		return user, err
	}
	m.Count++
	user := MustRandomUser()
	return user, orm.CreateUser(&user)
}

func NewMockAuthenticatedHTTPClient(cfg cmd.HTTPClientConfig, sessionID string) cmd.HTTPClient {
	return cmd.NewAuthenticatedHTTPClient(cfg, MockCookieAuthenticator{SessionID: sessionID}, sessions.SessionRequest{})
}

type MockCookieAuthenticator struct {
	SessionID string
	Error     error
}

func (m MockCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return MustGenerateSessionCookie(m.SessionID), m.Error
}

func (m MockCookieAuthenticator) Authenticate(sessions.SessionRequest) (*http.Cookie, error) {
	return MustGenerateSessionCookie(m.SessionID), m.Error
}

type MockSessionRequestBuilder struct {
	Count int
	Error error
}

func (m *MockSessionRequestBuilder) Build(string) (sessions.SessionRequest, error) {
	m.Count++
	if m.Error != nil {
		return sessions.SessionRequest{}, m.Error
	}
	return sessions.SessionRequest{Email: APIEmail, Password: Password}, nil
}

type MockSecretGenerator struct{}

func (m MockSecretGenerator) Generate(string) ([]byte, error) {
	return []byte(SessionSecret), nil
}

type MockChangePasswordPrompter struct {
	web.UpdatePasswordRequest
	err error
}

func (m MockChangePasswordPrompter) Prompt() (web.UpdatePasswordRequest, error) {
	return m.UpdatePasswordRequest, m.err
}

type MockPasswordPrompter struct {
	Password string
}

func (m MockPasswordPrompter) Prompt() string {
	return m.Password
}

var _ gracefulpanic.Signal = &testShutdownSignal{}

type testShutdownSignal struct {
	t testing.TB
}

func (tss *testShutdownSignal) Panic() {
	tss.t.Errorf("panic: %s", debug.Stack())
	panic("panic")
}

func (tss *testShutdownSignal) Wait() <-chan struct{} {
	return make(chan struct{})
}

func NewChainSetMockWithOneChain(t testing.TB, ethClient eth.Client, cfg evmconfig.ChainScopedConfig) evm.ChainSet {
	cc := new(evmmocks.ChainSet)
	ch := new(evmmocks.Chain)
	ch.On("Client").Return(ethClient)
	ch.On("Config").Return(cfg)
	ch.On("Logger").Return(logger.Default)
	ch.On("ID").Return(cfg.ChainID())
	cc.On("Default").Return(ch, nil)
	cc.On("Get", (*big.Int)(nil)).Return(ch, nil)
	cc.On("Chains").Return([]evm.Chain{ch})
	return cc
}
