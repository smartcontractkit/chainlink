package cltest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

// MockSubscription a mock subscription
type MockSubscription struct {
	t            testing.TB
	mut          sync.Mutex
	channel      interface{}
	unsubscribed bool
	Errors       chan error
}

// EmptyMockSubscription return empty MockSubscription
func EmptyMockSubscription(t testing.TB) *MockSubscription {
	return &MockSubscription{t: t, Errors: make(chan error, 1), channel: make(chan struct{})}
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
	case chan *evmtypes.Head:
		close(mes.channel.(chan *evmtypes.Head))
	default:
		logger.TestLogger(mes.t).Fatalf("Unable to close MockSubscription channel of type %T", mes.channel)
	}
	close(mes.Errors)
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

type InstanceAppFactoryWithKeystoreMock struct {
	App chainlink.Application
}

// NewApplication creates a new application with specified config and calls the authenticate function of the keystore
func (f InstanceAppFactoryWithKeystoreMock) NewApplication(ctx context.Context, cfg chainlink.GeneralConfig, lggr logger.Logger, registerer prometheus.Registerer, db *sqlx.DB, ks cmd.TerminalKeyStoreAuthenticator) (chainlink.Application, error) {
	keyStore := f.App.GetKeyStore()
	err := ks.Authenticate(ctx, keyStore, cfg.Password())
	if err != nil {
		return nil, fmt.Errorf("error authenticating keystore: %w", err)
	}
	return f.App, nil
}

// InstanceAppFactory is an InstanceAppFactory
type InstanceAppFactory struct {
	App chainlink.Application
}

// NewApplication creates a new application with specified config
func (f InstanceAppFactory) NewApplication(context.Context, chainlink.GeneralConfig, logger.Logger, prometheus.Registerer, *sqlx.DB, cmd.TerminalKeyStoreAuthenticator) (chainlink.Application, error) {
	return f.App, nil
}

type seededAppFactory struct {
	Application chainlink.Application
}

func (s seededAppFactory) NewApplication(context.Context, chainlink.GeneralConfig, logger.Logger, prometheus.Registerer, *sqlx.DB, cmd.TerminalKeyStoreAuthenticator) (chainlink.Application, error) {
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
func (r BlockedRunner) Run(context.Context, chainlink.Application) error {
	<-r.Done
	return nil
}

// EmptyRunner is an EmptyRunner
type EmptyRunner struct{}

// Run runs the empty runner
func (r EmptyRunner) Run(context.Context, chainlink.Application) error {
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
func (p *MockCountingPrompter) Prompt(string) string { return p.prompt() }

func (p *MockCountingPrompter) prompt() string {
	i := p.Count
	p.Count++
	if len(p.EnteredStrings)-1 < i {
		p.T.Errorf("Not enough passwords supplied to MockCountingPrompter, wanted %d", i)
		p.T.FailNow()
	}
	return p.EnteredStrings[i]
}

// PasswordPrompt returns an entered string
func (p *MockCountingPrompter) PasswordPrompt(string) string { return p.prompt() }

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
		b, err := io.ReadAll(r.Body)
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

// MustRandomUser inserts a new admin user with a random email into the test DB
func MustRandomUser(t testing.TB) sessions.User {
	email := fmt.Sprintf("user-%v@chainlink.test", NewRandomPositiveInt64())
	r, err := sessions.NewUser(email, Password, sessions.UserRoleAdmin)
	if err != nil {
		logger.TestLogger(t).Panic(err)
	}
	return r
}

func NewUserWithSession(t testing.TB, orm sessions.AuthenticationProvider) sessions.User {
	ctx := testutils.Context(t)
	u := MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &u))

	_, err := orm.CreateSession(ctx, sessions.SessionRequest{
		Email:    u.Email,
		Password: Password,
	})
	require.NoError(t, err)
	return u
}

type MockAPIInitializer struct {
	t     testing.TB
	Count int
}

func NewMockAPIInitializer(t testing.TB) *MockAPIInitializer {
	return &MockAPIInitializer{t: t}
}

func (m *MockAPIInitializer) Initialize(ctx context.Context, orm sessions.BasicAdminUsersORM, lggr logger.Logger) (sessions.User, error) {
	if user, err := orm.FindUser(ctx, APIEmailAdmin); err == nil {
		return user, err
	}
	m.Count++
	user := MustRandomUser(m.t)
	return user, orm.CreateUser(ctx, &user)
}

func NewMockAuthenticatedHTTPClient(lggr logger.Logger, cfg cmd.ClientOpts, sessionID string) cmd.HTTPClient {
	return cmd.NewAuthenticatedHTTPClient(lggr, cfg, MockCookieAuthenticator{SessionID: sessionID}, sessions.SessionRequest{})
}

type MockCookieAuthenticator struct {
	t         testing.TB
	SessionID string
	Error     error
}

func (m MockCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return MustGenerateSessionCookie(m.t, m.SessionID), m.Error
}

func (m MockCookieAuthenticator) Authenticate(context.Context, sessions.SessionRequest) (*http.Cookie, error) {
	return MustGenerateSessionCookie(m.t, m.SessionID), m.Error
}

func (m MockCookieAuthenticator) Logout() error {
	return nil
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
	return sessions.SessionRequest{Email: APIEmailAdmin, Password: Password}, nil
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

func NewLegacyChainsWithMockChain(t testing.TB, ethClient evmclient.Client, cfg toml.HasEVMConfigs) legacyevm.LegacyChainContainer {
	ch := new(evmmocks.Chain)
	ch.On("Client").Return(ethClient)
	ch.On("Logger").Return(logger.TestLogger(t))
	scopedCfg := evmtest.NewChainScopedConfig(t, cfg)
	ch.On("ID").Return(scopedCfg.EVM().ChainID())
	ch.On("Config").Return(scopedCfg)
	ch.On("HeadTracker").Return(nil)

	return NewLegacyChainsWithChain(ch, cfg)
}

func NewLegacyChainsWithMockChainAndTxManager(t testing.TB, ethClient evmclient.Client, cfg toml.HasEVMConfigs, txm txmgr.TxManager) legacyevm.LegacyChainContainer {
	ch := new(evmmocks.Chain)
	ch.On("Client").Return(ethClient)
	ch.On("Logger").Return(logger.TestLogger(t))
	scopedCfg := evmtest.NewChainScopedConfig(t, cfg)
	ch.On("ID").Return(scopedCfg.EVM().ChainID())
	ch.On("Config").Return(scopedCfg)
	ch.On("TxManager").Return(txm)

	return NewLegacyChainsWithChain(ch, cfg)
}

func NewLegacyChainsWithChain(ch legacyevm.Chain, cfg toml.HasEVMConfigs) legacyevm.LegacyChainContainer {
	m := map[string]legacyevm.Chain{ch.ID().String(): ch}
	return legacyevm.NewLegacyChains(m, cfg.EVMConfigs())
}
