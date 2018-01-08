package cltest

import (
	"fmt"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink-go/store"
)

func (self *TestApplication) MockEthClient() *EthMock {
	mock := NewMockGethRpc()
	eth := &store.EthClient{mock}
	self.Store.Eth.EthClient = eth
	return mock
}

func NewMockGethRpc() *EthMock {
	return &EthMock{}
}

type EthMock struct {
	Responses []MockResponse
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
	return len(self.Responses) == 0
}

func copyWithoutIndex(s []MockResponse, index int) []MockResponse {
	return append(s[:index], s[index+1:]...)
}

func (self *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	for i, resp := range self.Responses {
		if resp.methodName == method {
			self.Responses = copyWithoutIndex(self.Responses, i)
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

func (self *TestApplication) InstantClock() *InstantClock {
	clock := &InstantClock{}
	self.Scheduler.OneTime.Clock = clock
	return clock
}

type InstantClock struct{}

func (self *InstantClock) Sleep(_ time.Duration) {
}
