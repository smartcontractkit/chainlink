package client

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Types extracted from testhelpers to avoid import cycle

type AnyMsgSentEvent struct {
	SequenceNumber uint64
	// RawEvent contains the raw event depending on the chain:
	//  EVM:   *onramp.OnRampCCIPMessageSent
	//  Aptos: module_onramp.CCIPMessageSent
	RawEvent any
}

type CCIPSendReqConfig struct {
	SourceChain  uint64
	DestChain    uint64
	IsTestRouter bool
	Sender       *bind.TransactOpts
	Message      any
	MaxRetries   int // Number of retries for errors (excluding insufficient fee errors)
}

type SendReqOpts func(*CCIPSendReqConfig)

// WithMaxRetries sets the maximum number of retries for the CCIP send request.
func WithMaxRetries(maxRetries int) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.MaxRetries = maxRetries
	}
}

func WithSender(sender *bind.TransactOpts) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.Sender = sender
	}
}

func WithMessage(msg any) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.Message = msg
	}
}

func WithTestRouter(isTestRouter bool) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.IsTestRouter = isTestRouter
	}
}

func WithSourceChain(sourceChain uint64) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.SourceChain = sourceChain
	}
}

func WithDestChain(destChain uint64) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.DestChain = destChain
	}
}
