package withdrawprover

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l2_to_l1_message_passer"
)

var (
	// MessagePassedTopic is the topic of the MessagePassed event from the L2ToL1MessagePasser contract.
	MessagePassedTopic = optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed{}.Topic()
)

// OutputRootProof contains the elements that are hashed together to generate an output root
// which itself represents a snapshot of the L2 state.
type OutputRootProof struct {
	// Version of the output root.
	Version [32]byte

	// StateRoot is the root of the state trie at the block of this output.
	StateRoot [32]byte

	// MessagePasserStorageRoot is the root of the L2ToL1MessagePasser contract's storage trie.
	MessagePasserStorageRoot [32]byte

	// LatestBlockHash is the hash of the block this output was generated from.
	LatestBlockHash [32]byte
}

// BedrockMessageProof contains the elements that are needed to prove a withdrawal from L2 to L1.
type BedrockMessageProof struct {
	// LowLevelMessage is the MessagePassed event from the L2ToL1MessagePasser contract.
	// The jargon used throughout the optimism sdk is "low level message" to refer to this log,
	// so we use the same jargon here to reduce confusion.
	LowLevelMessage *optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed

	// WithdrawalProof is the merkle proof of inclusion of the withdrawal transaction in the L2 state trie.
	// The withdrawal transaction is stored in a mapping in the L2ToL1MessagePasser contract by hashing
	// various constituents of it together.
	// See https://github.com/ethereum-optimism/optimism/blob/005be54bde97747b6f1669030721cd4e0c14bc69/packages/contracts-bedrock/src/L2/L2ToL1MessagePasser.sol#L73
	// for more details.
	WithdrawalProof [][]byte

	// L2OutputIndex is the index of the L2 output.
	// Pre-FPAC, this contains the index of the output in the L2OutputOracle.
	// Post-FPAC, this will be the index of the game in the DisputeGameFactory.
	L2OutputIndex *big.Int

	// OutputRootProof contains the elements that are hashed together to generate an output root.
	// See the docstring of OutputRootProof for more details.
	OutputRootProof OutputRootProof
}
