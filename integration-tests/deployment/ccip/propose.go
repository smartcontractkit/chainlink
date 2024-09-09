package ccipdeployment

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

// TODO: Pull up to deploy
func SimTransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		return transaction, nil
	}, From: common.HexToAddress("0x0"), NoSend: true, GasLimit: 1_000_000}
}

func GenerateAcceptOwnershipProposal(
	e deployment.Environment,
	chains []uint64,
	ab deployment.AddressBook,
) (timelock.MCMSWithTimelockProposal, error) {
	state, err := LoadOnchainState(e, ab)
	if err != nil {
		return timelock.MCMSWithTimelockProposal{}, err
	}
	// TODO: Just onramp as an example
	var batches []timelock.BatchChainOperation
	metaDataPerChain := make(map[mcms.ChainIdentifier]timelock.MCMSWithTimelockChainMetadata)
	for _, sel := range chains {
		chain, _ := chainsel.ChainBySelector(sel)
		acceptOnRamp, err := state.Chains[sel].OnRamp.AcceptOwnership(SimTransactOpts())
		if err != nil {
			return timelock.MCMSWithTimelockProposal{}, err
		}
		chainSel := mcms.ChainIdentifier(chain.Selector)
		metaDataPerChain[chainSel] = timelock.MCMSWithTimelockChainMetadata{
			ChainMetadata: mcms.ChainMetadata{
				NonceOffset: 0,
				MCMAddress:  state.Chains[sel].McmAddr,
			},
			TimelockAddress: state.Chains[sel].TimelockAddr,
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: chainSel,
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[sel].OnRamp.Address(),
					Data:  acceptOnRamp.Data(),
					Value: big.NewInt(0),
				},
			},
		})
	}
	// TODO: Real valid until.
	return timelock.MCMSWithTimelockProposal{
		Operation:     timelock.Schedule,
		MinDelay:      "1h",
		ChainMetadata: metaDataPerChain,
		Transactions:  batches,
	}, nil
}
