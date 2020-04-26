package bulletin_board

import "regexp"

// BoardUpdateHandler is used to communicate values on a Board.
type BoardUpdateHandler func(BoardKey, BoardValue)

type Timeout error
type TimeoutHandler func(Timeout)

// Subscription represents a node's subscription to messages on another board
type Subscription struct {
	flatMatch bool // True for flatKey match, false for regexpKey match
	flatKey   BoardKey
	regexpKey regexp.Regexp
	//handler   SubscriptionHandler
}

var AllMessages = Subscription{}

func (s Subscription) Match(key BoardKey) bool {
	return (s.flatKey == "" && s.flatMatch) || // AllMessages
		(s.flatMatch && key == s.flatKey) || // Flat string match
		(!s.flatMatch && s.regexpKey.Match([]byte(key))) // Regexp match
}
