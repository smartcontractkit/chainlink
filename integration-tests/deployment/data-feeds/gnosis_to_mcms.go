package data_feeds

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type GnosisProposal struct {
	chainSel       uint64 // Which chain
	to             common.Address
	value          *big.Int
	data           []byte
	operation      uint8
	safeTxGas      *big.Int
	baseGas        *big.Int
	gasPrice       *big.Int
	gasToken       common.Address
	refundReceiver common.Address
	nonce          *big.Int
	// Populated via BuildGnosisProposals
	hashToSign [32]byte
	// Signatures populated offchain
	signatures []byte
	// TODO: Serialize this to/from a file for gnosis signing.
}

func BuildGnosisProposals(e deployment.Environment, ab deployment.AddressBook) ([]GnosisProposal, error) {
	chains, err := LoadOnchainState(e, ab)
	if err != nil {
		return nil, err
	}
	var proposals []GnosisProposal
	for chainSel, chain := range chains.Chains {
		// We want to call a transferOwnership on an aggregator from gnosis safe.
		// where the target is feeds timelock on the chain.
		transferOwnershipData, err1 := chain.Aggregator.TransferOwnership(deployment.SimTransactOpts(), chain.Timelock.Address())
		if err1 != nil {
			return nil, err1
		}
		proposal := GnosisProposal{
			chainSel:       chainSel,
			to:             chain.AggregatorAddr,
			value:          big.NewInt(0),
			data:           transferOwnershipData.Data(),
			operation:      0,   // Call?
			safeTxGas:      nil, // TODO
			baseGas:        nil, // TODO
			gasPrice:       nil, // TODO
			gasToken:       common.Address{},
			refundReceiver: common.Address{},
			nonce:          nil,
		}
		hashToSign, err1 := chain.Safe.GetTransactionHash(nil,
			proposal.to,
			proposal.value,
			proposal.data,
			proposal.operation,
			proposal.safeTxGas,
			proposal.baseGas,
			proposal.gasPrice,
			proposal.gasToken,
			proposal.refundReceiver,
			proposal.nonce,
		)
		if err1 != nil {
			return nil, err1
		}
		proposal.hashToSign = hashToSign
	}
	return proposals, nil
}

func ApplySignedGnosisProposals(e deployment.Environment, ab deployment.AddressBook, proposals []GnosisProposal) error {
	chains, err := LoadOnchainState(e, ab)
	if err != nil {
		return err
	}
	for _, proposal := range proposals {
		chainState := chains.Chains[proposal.chainSel]
		_, err1 := chainState.Safe.ExecTransaction(
			e.Chains[proposal.chainSel].DeployerKey,
			proposal.to,
			proposal.value,
			proposal.data,
			proposal.operation,
			proposal.safeTxGas,
			proposal.baseGas,
			proposal.gasPrice,
			proposal.gasToken,
			proposal.refundReceiver,
			proposal.signatures,
		)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func GenerateAcceptOwnershipProposal() timelock.MCMSWithTimelockProposal {
	// See similar proposals in CCIP.
	return timelock.MCMSWithTimelockProposal{}
}
