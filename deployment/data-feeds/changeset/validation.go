package changeset

import (
	"errors"
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

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

func ValidateMCMSAddresses(addressStore datastore.AddressRefStore, chainSelector uint64) error {
	records := addressStore.Filter(datastore.AddressRefByType(datastore.ContractType(commonTypes.RBACTimelock)))
	if len(records) == 0 {
		return fmt.Errorf("timelock not present on the chain %d", chainSelector)
	}

	records = addressStore.Filter(datastore.AddressRefByType(datastore.ContractType(commonTypes.ProposerManyChainMultisig)))
	if len(records) == 0 {
		return fmt.Errorf("mcms proposer not present on the chain %d", chainSelector)
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

func ValidateCacheForTronChain(env cldf.Environment, chainSelector uint64, cacheAddress address.Address) error {
	state, err := LoadTronOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load on chain state %w", err)
	}
	_, ok := env.BlockChains.TronChains()[chainSelector]
	if !ok {
		return errors.New("chain not found in environment")
	}
	chainState, ok := state.TronChains[chainSelector]
	if !ok {
		return errors.New("chain not found in on chain state")
	}
	if chainState.DataFeeds == nil {
		return errors.New("DataFeeds not found in on chain state")
	}
	addr := cacheAddress.String()
	isEvm, _ := chain_selectors.IsEvm(chainSelector)
	if isEvm {
		addr = cacheAddress.EthAddress().Hex()
	}
	exists := chainState.DataFeeds[addr]
	if !exists {
		return errors.New("contract not found in on chain state")
	}
	return nil
}
