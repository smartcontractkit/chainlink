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

func (ta *TestApplication) MockEthClient() *EthMock {
	return MockEthOnStore(ta.Store)
}

func MockEthOnStore(s *store.Store) *EthMock {
	mock := &EthMock{}
	eth := &store.EthClient{mock}
	s.TxManager.EthClient = eth
	return mock
}

type EthMock struct {
	Responses      []MockResponse
	Subscriptions  []MockSubscription
	newHeadsCalled bool
}

func (mock *EthMock) Register(
	method string,
	response interface{},
	callback ...func(interface{}, ...interface{}) error,
) {
	res := MockResponse{
		methodName: method,
		response:   response,
	}
	if len(callback) > 0 {
		res.callback = callback[0]
	}
	mock.Responses = append(mock.Responses, res)
}

func (mock *EthMock) RegisterError(method, errMsg string) {
	res := MockResponse{
		methodName: method,
		errMsg:     errMsg,
		hasError:   true,
	}
	mock.Responses = append(mock.Responses, res)
}

func (mock *EthMock) AllCalled() bool {
	return (len(mock.Responses) == 0) && (len(mock.Subscriptions) == 0)
}

func (mock *EthMock) EventuallyAllCalled(t *testing.T) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	g.Eventually(mock.AllCalled).Should(gomega.BeTrue())
}

func (mock *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	for i, resp := range mock.Responses {
		if resp.methodName == method {
			mock.Responses = append(mock.Responses[:i], mock.Responses[i+1:]...)
			if resp.hasError {
				return fmt.Errorf(resp.errMsg)
			} else {
				ref := reflect.ValueOf(result)
				reflect.Indirect(ref).Set(reflect.ValueOf(resp.response))
				if resp.callback != nil {
					if err := resp.callback(result, args); err != nil {
						return fmt.Errorf("ethMock Error: %v", err)
					}
				}
				return nil
			}
		}
	}
	return fmt.Errorf("EthMock: Method %v not registered", method)
}

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

func (mock *EthMock) EthSubscribe(
	ctx context.Context,
	channel interface{},
	args ...interface{},
) (models.EthSubscription, error) {
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
	} else if args[0] == "newHeads" {
		return nil, errors.New("newHeads subscription only expected once, please register another mock subscription if more are needed.")
	}
	return nil, errors.New("Must RegisterSubscription before EthSubscribe")
}

func (mock *EthMock) RegisterNewHeads() chan models.BlockHeader {
	newHeads := make(chan models.BlockHeader, 10)
	mock.RegisterSubscription("newHeads", newHeads)
	return newHeads
}

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

type MockSubscription struct {
	name    string
	channel interface{}
	Errors  chan error
}

func EmptyMockSubscription() MockSubscription {
	return MockSubscription{Errors: make(chan error, 1), channel: make(chan struct{})}
}

func (mes MockSubscription) Err() <-chan error { return mes.Errors }
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

type MockResponse struct {
	methodName string
	response   interface{}
	errMsg     string
	hasError   bool
	callback   func(interface{}, ...interface{}) error
}

func (ta *TestApplication) InstantClock() InstantClock {
	clock := InstantClock{}
	ta.Scheduler.OneTime.Clock = clock
	return clock
}

func UseSettableClock(s *store.Store) *SettableClock {
	clock := &SettableClock{}
	s.Clock = clock
	return clock
}

type SettableClock struct {
	time time.Time
}

func (clock *SettableClock) Now() time.Time {
	if clock.time.IsZero() {
		return time.Now()
	}
	return clock.time
}

func (clock *SettableClock) SetTime(t time.Time) {
	clock.time = t
}

func (*SettableClock) After(_ time.Duration) <-chan time.Time {
	channel := make(chan time.Time, 1)
	channel <- time.Now()
	return channel
}

type InstantClock struct{}

func (InstantClock) Now() time.Time {
	return time.Now()
}

func (InstantClock) After(_ time.Duration) <-chan time.Time {
	c := make(chan time.Time, 100)
	c <- time.Now()
	return c
}

type NeverClock struct{}

func (NeverClock) After(_ time.Duration) <-chan time.Time {
	return make(chan time.Time)
}

func (NeverClock) Now() time.Time {
	return time.Now()
}

type RendererMock struct {
	Renders []interface{}
}

func (rm *RendererMock) Render(v interface{}) error {
	rm.Renders = append(rm.Renders, v)
	return nil
}

type InstanceAppFactory struct {
	App services.Application
}

func (f InstanceAppFactory) NewApplication(config store.Config) services.Application {
	return f.App
}

type EmptyAppFactory struct{}

func (f EmptyAppFactory) NewApplication(config store.Config) services.Application {
	return &EmptyApplication{}
}

type EmptyApplication struct{}

func (a *EmptyApplication) Start() error {
	return nil
}

func (a *EmptyApplication) Stop() error {
	return nil
}

func (a *EmptyApplication) GetStore() *store.Store {
	return nil
}

type CallbackAuthenticator struct {
	Callback func(*store.Store, string)
}

func (a CallbackAuthenticator) Authenticate(store *store.Store, pwd string) {
	a.Callback(store, pwd)
}

type EmptyRunner struct{}

func (r EmptyRunner) Run(app services.Application) error {
	return nil
}

type MockCountingPrompt struct {
	EnteredStrings []string
	Count          int
}

func (p *MockCountingPrompt) Prompt(string) string {
	i := p.Count
	p.Count++
	return p.EnteredStrings[i]
}

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
		assert.Nil(t, err)
		assert.Equal(t, wantMethod, r.Method)
		if len(callback) > 0 {
			callback[0](string(b))
		}
		called = true

		w.WriteHeader(status)
		io.WriteString(w, response)
	})

	return httptest.NewServer(handler), func() {
		assert.True(t, called)
	}
}

type MockCron struct {
	Entries []MockCronEntry
}

func NewMockCron() *MockCron {
	return &MockCron{}
}

func (*MockCron) Start() {}
func (*MockCron) Stop()  {}

func (mc *MockCron) AddFunc(schd string, fn func()) error {
	mc.Entries = append(mc.Entries, MockCronEntry{
		Schedule: schd,
		Function: fn,
	})
	return nil
}

func (mc *MockCron) RunEntries() {
	for _, entry := range mc.Entries {
		entry.Function()
	}
}

type MockCronEntry struct {
	Schedule string
	Function func()
}

type MockHeadTrackable struct {
	ConnectedCount    int
	DisconnectedCount int
	OnNewHeadCount    int
}

func (m *MockHeadTrackable) Connect(*models.IndexableBlockNumber) error {
	m.ConnectedCount += 1
	return nil
}

func (m *MockHeadTrackable) Disconnect()                   { m.DisconnectedCount += 1 }
func (m *MockHeadTrackable) OnNewHead(*models.BlockHeader) { m.OnNewHeadCount += 1 }

type NeverSleeper struct{}

func (ns NeverSleeper) Reset()                  {}
func (ns NeverSleeper) Sleep()                  {}
func (ns NeverSleeper) Duration() time.Duration { return 0 * time.Microsecond }
