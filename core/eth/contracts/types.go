package contracts

import (
	"chainlink/core/eth"
)

type MaybeDecodedLog struct {
	Log   interface{}
	Error error
}

type LogSubscription struct {
	subscription eth.Subscription
	chLogs       chan MaybeDecodedLog
}

func (s *LogSubscription) Logs() <-chan MaybeDecodedLog {
	return s.chLogs
}

func (s *LogSubscription) Unsubscribe() {
	s.subscription.Unsubscribe()
}
