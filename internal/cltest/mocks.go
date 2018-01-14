package cltest

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/store"
)

func (self *TestApplication) MockEthClient() *EthMock {
	return MockEthOnStore(self.Store)
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

func (self *EthMock) Register(method string, response interface{}) {
	res := MockResponse{
		methodName: method,
		response:   response,
	}
	self.Responses = append(self.Responses, res)
}

func (self *EthMock) RegisterError(method, errMsg string) {
	res := MockResponse{
		methodName: method,
		errMsg:     errMsg,
		hasError:   true,
	}
	self.Responses = append(self.Responses, res)
}

func (self *EthMock) AllCalled() bool {
	return (len(self.Responses) == 0) && (len(self.Subscriptions) == 0)
}

func (self *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	for i, resp := range self.Responses {
		if resp.methodName == method {
			self.Responses = append(self.Responses[:i], self.Responses[i+1:]...)
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

func (self *EthMock) RegisterSubscription(name string, channel interface{}) {
	res := MockSubscription{
		name:    name,
		channel: channel,
	}
	self.Subscriptions = append(self.Subscriptions, res)
}

func (self *EthMock) EthSubscribe(
	ctx context.Context,
	channel interface{},
	args ...interface{},
) (*rpc.ClientSubscription, error) {
	for i, sub := range self.Subscriptions {
		if sub.name == args[0] {
			self.Subscriptions = append(self.Subscriptions[:i], self.Subscriptions[i+1:]...)
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

func (self *TestApplication) InstantClock() InstantClock {
	clock := InstantClock{}
	self.Scheduler.OneTime.Clock = clock
	return clock
}

type InstantClock struct{}

func (self InstantClock) After(_ time.Duration) <-chan time.Time {
	c := make(chan time.Time, 100)
	c <- time.Now()
	return c
}

type NeverClock struct{}

func (self NeverClock) After(_ time.Duration) <-chan time.Time {
	return make(chan time.Time)
}
