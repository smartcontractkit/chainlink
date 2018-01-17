package cltest

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/store"
)

func (ta *TestApplication) MockEthClient() *EthMock {
	return MockEthOnStore(ta.Store)
}

func MockEthOnStore(s *store.Store) *EthMock {
	mock := NewMockGethRpc()
	eth := &store.EthClient{mock}
	s.TxManager.EthClient = eth
	return mock
}

func NewMockGethRpc() *EthMock {
	return &EthMock{}
}

type EthMock struct {
	Responses     []MockResponse
	Subscriptions []MockSubscription
}

type MockResponse struct {
	methodName string
	response   interface{}
	errMsg     string
	hasError   bool
}

func (mock *EthMock) Register(method string, response interface{}) {
	res := MockResponse{
		methodName: method,
		response:   response,
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

func (mock *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	for i, resp := range mock.Responses {
		if resp.methodName == method {
			mock.Responses = append(mock.Responses[:i], mock.Responses[i+1:]...)
			if resp.hasError {
				return fmt.Errorf(resp.errMsg)
			} else {
				ref := reflect.ValueOf(result)
				reflect.Indirect(ref).Set(reflect.ValueOf(resp.response))
				return nil
			}
		}
	}
	return fmt.Errorf("EthMock: Method %v not registered", method)
}

type MockSubscription struct {
	name    string
	channel interface{}
}

func (mock *EthMock) RegisterSubscription(name string, channel interface{}) {
	res := MockSubscription{
		name:    name,
		channel: channel,
	}
	mock.Subscriptions = append(mock.Subscriptions, res)
}

func (mock *EthMock) EthSubscribe(
	ctx context.Context,
	channel interface{},
	args ...interface{},
) (*rpc.ClientSubscription, error) {
	for i, sub := range mock.Subscriptions {
		if sub.name == args[0] {
			mock.Subscriptions = append(mock.Subscriptions[:i], mock.Subscriptions[i+1:]...)
			mockChan := sub.channel.(chan store.EventLog)
			logChan := channel.(chan store.EventLog)
			go func() {
				for e := range mockChan {
					logChan <- e
				}
			}()
			return &rpc.ClientSubscription{}, nil
		}
	}
	return &rpc.ClientSubscription{}, nil
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
