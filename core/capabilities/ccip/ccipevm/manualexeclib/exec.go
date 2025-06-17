package manualexeclib

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

func GetMessageHashes(
	ctx context.Context,
	lggr logger.Logger,
	onrampAddress common.Address,
	ccipMessageSentEvents []onramp.OnRampCCIPMessageSent,
	extraDataCodec ccipcommon.ExtraDataCodec,
) ([][32]byte, error) {
	msgHasher := ccipevm.NewMessageHasherV1(
		lggr,
		extraDataCodec,
	)
	var ret [][32]byte
	for _, event := range ccipMessageSentEvents {
		ccipMsg := ccipevm.EVM2AnyToCCIPMsg(onrampAddress, event.Message)
		msgHash, err := msgHasher.Hash(ctx, ccipMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to hash message (evm2any: %+v, generic: %+v): %w", event.Message, ccipMsg, err)
		}
		ret = append(ret, msgHash)
	}

	return ret, nil
}

// GetMerkleProof returns the merkle proof of inclusion for the given sequence number
// in the given merkleRoot.
//
// In the event that:
// 1. the calculated merkle root does not match the committed merkle root
// 2. the sequence number is not found in the merkle root struct
// an error is returned.
func GetMerkleProof(
	lggr logger.Logger,
	merkleRoot offramp.InternalMerkleRoot,
	messageHashes [][32]byte,
	msgSeqNr uint64,
) ([][32]byte, *big.Int, error) {
	mtree, err := merklemulti.NewTree(hashutil.NewKeccak(), messageHashes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create merkle tree: %w", err)
	}

	// calculate merkle root from tree, should match the committed root
	root := mtree.Root()
	if root != merkleRoot.MerkleRoot {
		return nil, nil, fmt.Errorf(
			"merkle root mismatch, calculated != committed: %x != %x",
			root, merkleRoot.MerkleRoot)
	}

	lggr.Debugw("merkle roots match", "calculated", hexutil.Encode(root[:]), "committed", hexutil.Encode(merkleRoot.MerkleRoot[:]))

	// get the index of the msgSeqNr in the messageHashes
	var idx int = -1
	var j int
	for i := merkleRoot.MinSeqNr; i <= merkleRoot.MaxSeqNr; i++ {
		if i == msgSeqNr {
			idx = j
			break
		}
		j++
	}

	if idx == -1 {
		return nil, nil, fmt.Errorf("msgSeqNr %d not found in merkle root struct, range: [%d, %d]",
			msgSeqNr, merkleRoot.MinSeqNr, merkleRoot.MaxSeqNr)
	}

	proof, err := mtree.Prove([]int{idx})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to prove: %w", err)
	}

	bitFlag := ccipevm.BoolsToBitFlags(proof.SourceFlags)

	return proof.Hashes, bitFlag, nil
}

// CreateExecutionReport creates an offramp.InternalExecutionReport from the given parameters.
//
// It converts the given ccipMessageSentEvents to offramp.InternalAny2EVMRampMessage
// and then to offramp.InternalExecutionReport.
//
// The offchainTokenData is not currently supported and is set to an empty slice.
func CreateExecutionReport(
	srcChainSel uint64,
	onrampAddress common.Address,
	ccipMessageSentEvents []onramp.OnRampCCIPMessageSent,
	hashes [][32]byte,
	flags *big.Int,
	extraDataCodec ccipcommon.ExtraDataCodec,
) (offramp.InternalExecutionReport, error) {
	var any2EVMs []offramp.InternalAny2EVMRampMessage
	for _, event := range ccipMessageSentEvents {
		ccipMsg := ccipevm.EVM2AnyToCCIPMsg(onrampAddress, event.Message)
		any2EVM, err := ccipevm.CCIPMsgToAny2EVMMessage(ccipMsg, extraDataCodec)
		if err != nil {
			return offramp.InternalExecutionReport{}, fmt.Errorf("failed to convert ccip message to any2evm message: %w", err)
		}
		any2EVMs = append(any2EVMs, any2EVM)
	}

	return offramp.InternalExecutionReport{
		SourceChainSelector: srcChainSel,
		Messages:            any2EVMs,
		// not currently supported
		OffchainTokenData: [][][]byte{
			{},
		},
		Proofs:        hashes,
		ProofFlagBits: flags,
	}, nil
}
