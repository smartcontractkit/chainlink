package ccipdeployment

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

// TODO: Pull up to deploy
func SimTransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		return transaction, nil
	}, From: common.HexToAddress("0x0"), NoSend: true, GasLimit: 200_000}
}

func GenerateAcceptOwnershipProposal(
	e deployment.Environment,
	chains []uint64,
	ab deployment.AddressBook,
) (deployment.Proposal, error) {
	state, err := LoadOnchainState(e, ab)
	if err != nil {
		return deployment.Proposal{}, err
	}
	// TODO: Just onramp as an example
	var ops []owner_helpers.ManyChainMultiSigOp
	for _, sel := range chains {
		opCount, err := state.Chains[sel].Mcm.GetOpCount(nil)
		if err != nil {
			return deployment.Proposal{}, err
		}

		txData, err := state.Chains[sel].EvmOnRampV160.AcceptOwnership(SimTransactOpts())
		if err != nil {
			return deployment.Proposal{}, err
		}
		evmID, err := chainsel.ChainIdFromSelector(sel)
		if err != nil {
			return deployment.Proposal{}, err
		}
		ops = append(ops, owner_helpers.ManyChainMultiSigOp{
			ChainId:  big.NewInt(int64(evmID)),
			MultiSig: state.Chains[sel].McmsAddr,
			Nonce:    opCount,
			To:       state.Chains[sel].EvmOnRampV160.Address(),
			Value:    big.NewInt(0),
			Data:     txData.Data(),
		})
	}
	// TODO: Real valid until.
	return deployment.Proposal{ValidUntil: uint32(time.Now().Unix()), Ops: ops}, nil
}

func ApplyProposal(env deployment.Environment, p deployment.Proposal, state CCIPOnChainState) error {
	// TODO
	return nil
}
