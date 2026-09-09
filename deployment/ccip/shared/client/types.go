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

	// SkipSuiFeeQuoterPriceUpdate skips the legacy EOA FeeQuoter price update that
	// SendRequestSui runs before ccip_send. That update uses the deployer's
	// CCIPOwnerCapObjectId, which is consumed (moved into the MCMS registry) once Sui
	// CCIP ownership is transferred to MCMS — as the lanes-based Sui<->Solana setup
	// does. Lanes setup seeds Sui fee-quoter prices via an MCMS proposal instead, so
	// the EOA update is both redundant and broken on that path. Legacy deployer-owned
	// Sui tests leave this false and keep the EOA update.
	SkipSuiFeeQuoterPriceUpdate bool
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

// WithSkipSuiFeeQuoterPriceUpdate marks the Sui source send as having its fee-quoter
// prices already seeded by the lane setup, so SendRequestSui must not run its legacy
// EOA price update (whose OwnerCap is consumed under MCMS ownership). See
// CCIPSendReqConfig.SkipSuiFeeQuoterPriceUpdate.
func WithSkipSuiFeeQuoterPriceUpdate() SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.SkipSuiFeeQuoterPriceUpdate = true
	}
}
