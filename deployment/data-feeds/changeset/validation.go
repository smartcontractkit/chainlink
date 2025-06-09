package changeset

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func ValidateCacheForChain(env cldf.Environment, chainSelector uint64, cacheAddress common.Address) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load on chain state %w", err)
	}
	_, ok := env.BlockChains.EVMChains()[chainSelector]
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

func ValidateMCMSAddresses(ab cldf.AddressBook, chainSelector uint64) error {
	if _, err := cldf.SearchAddressBook(ab, chainSelector, commonTypes.RBACTimelock); err != nil {
		return fmt.Errorf("timelock not present on the chain %w", err)
	}
	if _, err := cldf.SearchAddressBook(ab, chainSelector, commonTypes.ProposerManyChainMultisig); err != nil {
		return fmt.Errorf("mcms proposer not present on the chain %w", err)
	}
	return nil
}

func ValidateCacheForAptosChain(env cldf.Environment, chainSelector uint64, cacheAddress string) error {
	state, err := LoadAptosOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load on chain state %w", err)
	}
	_, ok := env.BlockChains.AptosChains()[chainSelector]
	if !ok {
		return errors.New("chain not found in environment")
	}
	chainState, ok := state.AptosChains[chainSelector]
	if !ok {
		return errors.New("chain not found in on chain state")
	}
	if chainState.DataFeeds == nil {
		return errors.New("DataFeeds not found in on chain state")
	}
	cacheAccountAddress := aptos.AccountAddress{}
	err = cacheAccountAddress.ParseStringRelaxed(cacheAddress)
	if err != nil {
		return fmt.Errorf("failed to parse cache address %w", err)
	}
	_, ok = chainState.DataFeeds[cacheAccountAddress]
	if !ok {
		return errors.New("contract not found in on chain state")
	}
	return nil
}
