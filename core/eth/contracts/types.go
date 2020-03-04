package contracts

import (
	"chainlink/core/eth"
)

type MaybeDecodedLog struct {
	Log   interface{}
	Error error
}

//go:generate mockery -name LogSubscription -output ../../internal/mocks/ -case=underscore

type LogSubscription interface {
	Logs() <-chan MaybeDecodedLog
	Unsubscribe()
}

type logSubscription struct {
	subscription eth.Subscription
	chLogs       chan MaybeDecodedLog
}

func (s *logSubscription) Logs() <-chan MaybeDecodedLog {
	return s.chLogs
}

func (s *logSubscription) Unsubscribe() {
	s.subscription.Unsubscribe()
}
