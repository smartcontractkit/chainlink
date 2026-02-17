package testhelpers

import (
	"context"
	"math/big"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// Based on unfinished implementation at https://github.com/smartcontractkit/chainlink/blob/f0432dc777d33b83a621da2b042657601d5db8b6/integration-tests/smoke/ccip/canonical/types/types.go

type TokenAmount struct {
	// Will be encoded in the source-native format, so EIP-55 for Ethereum,
	// base58 for Solana, etc.
	Token  string
	Amount *big.Int
}

// MessageComponents is a struct that contains the makeup for a general CCIP message
// irrespective of the chain family it originates from.
type MessageComponents struct {
	DestChainSelector uint64
	// Receiver is the receiver on the destination chain.
	// Must be appropriately dest-chain-family encoded, so abi.encode for Ethereum,
	// 32 bytes for Solana, etc.
	Receiver []byte
	// Data is the data to be sent to the destination chain.
	Data []byte
	// Will be encoded in the source-native format, so EIP-55 for Ethereum,
	// base58 for Solana, etc.
	FeeToken string
	// ExtraArgs are the message extra args which tune message semantics and behavior.
	// For example, out of order execution can be specified here.
	ExtraArgs []byte
	// TokenAmounts are the tokens and their respective amounts to be sent to the
	// destination chain.
	// Note that the tokens must be "approved" to the router for the message send to work.
	TokenAmounts []TokenAmount
}

// ExtraArgOpt is a generic representation of an extra arg that can be applied
// to any kind of ccip message.
// We use this to make it possible to specify extra args in a chain-agnostic way.
type ExtraArgOpt struct {
	Name  string
	Value any
}

func NewOutOfOrderExtraArg(outOfOrder bool) ExtraArgOpt {
	return ExtraArgOpt{
		Name:  "outOfOrderExecutionEnabled",
		Value: outOfOrder,
	}
}

func NewGasLimitExtraArg(gasLimit *big.Int) ExtraArgOpt {
	return ExtraArgOpt{
		Name:  "gasLimit|computeUnits",
		Value: gasLimit,
	}
}

// Adapter is our interface for interacting with a specific chain, scoped to a family.
// An adapter instance is an instance of a concrete chain.
// So if there are e.g 3 source chains that are EVM and a dest that is Solana,
// we would have 3 EVM adapters and 1 Solana adapter.
type Adapter interface {
	// ChainSelector returns the selector of the chain for the given adapter.
	ChainSelector() uint64

	// ChainFamily returns the family of the chain for the given adapter.
	Family() string

	// BuildMessage builds a message from the given components,
	// with the overall message type being ChainFamily2Any, where
	// ChainFamily is the family of the adapter.
	// As a concrete example, for EVM, the message type is router.ClientEVM2AnyMessage,
	// and for Solana, the message type is ccip_router.SVM2AnyMessage.
	BuildMessage(components MessageComponents) (any, error)

	// // RandomReceiver returns a random receiver for the given chain family.
	// RandomReceiver() []byte

	// // CCIPReceiver returns a CCIP receiver for the given chain family.
	// CCIPReceiver() []byte

	// NativeFeeToken returns the native fee token for the given chain family.
	NativeFeeToken() string

	// GetExtraArgs returns the default extra args for sending messages to this
	// chain family from the given source family.
	// Therefore the extra args are source-family encoded, so abi.encode for EVM,
	// borsch for Solana, etc.
	GetExtraArgs(receiver []byte, sourceFamily string, opts ...ExtraArgOpt) ([]byte, error)

	// GetInboundNonce returns the inbound nonce for the given sender and source chain selector.
	// For chains that don't have the concept of nonces, this will always return 0.
	GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error)

	// ValidateCommit validates that the message specified by the given send event was committed.
	ValidateCommit(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNumRange ccipocr3.SeqNumRange)

	// ValidateExec validates that the message specified by the given send event was executed.
	ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (execStates map[uint64]int)
}

type AdapterFactory = func(cldf.BlockChain, deployment.Environment) Adapter

var Adapters = map[string]AdapterFactory{
	chain_selectors.FamilyEVM:    NewEVMAdapter,
	chain_selectors.FamilyTron:   NewEVMAdapter, // TODO: is this right?
	chain_selectors.FamilySolana: NewSVMAdapter,
	chain_selectors.FamilyAptos:  NewAptosAdapter,
	chain_selectors.FamilySui:    NewSuiAdapter,
	chain_selectors.FamilyTon:    NewTonAdapter,
}
