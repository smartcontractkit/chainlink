package changeset

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[TransferLinkTokenConfig] = TransferLinkToken

type Transfer struct {
	To     common.Address
	Amount *big.Int
}

type TransferLinkTokenConfig struct {
	Transfers map[uint64]Transfer
}

func (c TransferLinkTokenConfig) Validate() error {
	for k, v := range c.Transfers {
		if err := deployment.IsValidChainSelector(k); err != nil {
			return err
		}

		if v.To == (common.Address{}) {
			return errors.New("to address must be set")
		}
		if v.Amount == nil || v.Amount.Sign() == -1 {
			return errors.New("amount must be set and positive")
		}
	}
	return nil
}

// TransferLinkToken transfers link token to the to address on the chain identified by the chainSelector.
func TransferLinkToken(e deployment.Environment, config TransferLinkTokenConfig) (deployment.ChangesetOutput, error) {
	if err := config.Validate(); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSelector, transferConfig := range config.Transfers {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		chain, ok := e.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
		}

		linkState, err := LoadLinkTokenState(chain, addresses)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		tx, err := linkState.LinkToken.Transfer(chain.DeployerKey, transferConfig.To, transferConfig.Amount)
		if _, err = deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transfer link token to %s: %v", transferConfig.To, err)
		}
		e.Logger.Infow("Transferred LINK",
			"to", transferConfig.To,
			"amount", transferConfig.Amount,
			"txHash", tx.Hash().Hex(),
			"chainSelector", chainSelector)
	}
	return deployment.ChangesetOutput{}, nil
}
