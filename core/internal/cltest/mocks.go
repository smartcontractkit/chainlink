package cltest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/web"

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
	case chan *models.Head:
		close(mes.channel.(chan *models.Head))
	default:
		logger.Fatal(fmt.Sprintf("Unable to close MockSubscription channel of type %T", mes.channel))
	}
	close(mes.Errors)
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
func (rm *RendererMock) Render(v interface{}, headers ...string) error {
	rm.Renders = append(rm.Renders, v)
	return nil
}

// InstanceAppFactory is an InstanceAppFactory
type InstanceAppFactory struct {
	App chainlink.Application
}

// NewApplication creates a new application with specified config
func (f InstanceAppFactory) NewApplication(config *orm.Config, onConnectCallbacks ...func(chainlink.Application)) (chainlink.Application, error) {
	return f.App, nil
}

type seededAppFactory struct {
	Application chainlink.Application
}

func (s seededAppFactory) NewApplication(config *orm.Config, onConnectCallbacks ...func(chainlink.Application)) (chainlink.Application, error) {
	return noopStopApplication{s.Application}, nil
}

type noopStopApplication struct {
	chainlink.Application
}

func (a noopStopApplication) Stop() error {
	return nil
}

// CallbackAuthenticator contains a call back authenticator method
type CallbackAuthenticator struct {
	Callback func(*keystore.Eth, string) (string, error)
}

// Authenticate authenticates store and pwd with the callback authenticator
func (a CallbackAuthenticator) AuthenticateEthKey(ethKeyStore *keystore.Eth, pwd string) (string, error) {
	return a.Callback(ethKeyStore, pwd)
}

func (a CallbackAuthenticator) AuthenticateVRFKey(*store.Store, string) error {
	return nil
}

func (a CallbackAuthenticator) AuthenticateOCRKey(*keystore.OCR, *orm.Config, string) error {
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

// NewHTTPMockServerWithRequest creates http test server that makes the request
// available in the callback
func NewHTTPMockServerWithRequest(
	t *testing.T,
	status int,
	response string,
	callback func(r *http.Request),
) (*httptest.Server, func()) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback(r)
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

func NewHTTPMockServerWithAlterableResponseAndRequest(t *testing.T, response func() string, callback func(r *http.Request)) (server *httptest.Server) {
	server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callback(r)
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

// OnNewLongestChain increases the OnNewLongestChainCount count by one
func (m *MockHeadTrackable) OnNewLongestChain(context.Context, models.Head) {
	atomic.AddInt32(&m.onNewHeadCount, 1)
}

// OnNewLongestChainCount returns the count of new heads, safely.
func (m *MockHeadTrackable) OnNewLongestChainCount() int32 {
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
	return cmd.NewAuthenticatedHTTPClient(cfg, MockCookieAuthenticator{SessionID: sessionID}, models.SessionRequest{})
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
