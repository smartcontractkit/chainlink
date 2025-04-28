package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

func ValidateCacheForChain(env deployment.Environment, chainSelector uint64, cacheAddress common.Address) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load on chain state %w", err)
	}
	_, ok := env.Chains[chainSelector]
	if !ok {
		return errors.New("chain not found in environment")
	}
	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return errors.New("chain not found in on chain state")
	}
	if chainState.DataFeedsCache == nil {
		return errors.New("DataFeedsCache not found in on chain state")
	}
	_, ok = chainState.DataFeedsCache[cacheAddress]
	if !ok {
		return errors.New("contract not found in on chain state")
	}
	return nil
}

func ValidateMCMSAddresses(ab deployment.AddressBook, chainSelector uint64) error {
	if _, err := deployment.SearchAddressBook(ab, chainSelector, commonTypes.RBACTimelock); err != nil {
		return fmt.Errorf("timelock not present on the chain %w", err)
	}
	if _, err := deployment.SearchAddressBook(ab, chainSelector, commonTypes.ProposerManyChainMultisig); err != nil {
		return fmt.Errorf("mcms proposer not present on the chain %w", err)
	}
	return nil
}
