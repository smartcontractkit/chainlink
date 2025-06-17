package changeset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	modulefeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"

	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/view"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/view/v1_0"
)

var (
	DataFeedsCache            cldf.ContractType      = "DataFeedsCache"
	ChainlinkDataFeedsPackage datastore.ContractType = "ChainlinkDataFeeds"
)

type DataFeedsChainState struct {
	ABIByAddress map[string]string
	commonchangeset.MCMSWithTimelockState
	DataFeedsCache  map[common.Address]*cache.DataFeedsCache
	AggregatorProxy map[common.Address]*proxy.AggregatorProxy
}

type DataFeedsAptosChainState struct {
	DataFeeds map[aptos.AccountAddress]*modulefeeds.DataFeeds
}

type DataFeedsOnChainState struct {
	Chains      map[uint64]DataFeedsChainState
	AptosChains map[uint64]DataFeedsAptosChainState
}

func LoadAptosOnchainState(e cldf.Environment) (DataFeedsOnChainState, error) {
	state := DataFeedsOnChainState{
		AptosChains: make(map[uint64]DataFeedsAptosChainState),
	}

	for chainSelector, chain := range e.BlockChains.AptosChains() {
		records := e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
		chainState, err := LoadAptosChainState(e.Logger, chain, records)
		if err != nil {
			return state, err
		}
		state.AptosChains[chainSelector] = *chainState
	}
	return state, nil
}

func LoadOnchainState(e cldf.Environment) (DataFeedsOnChainState, error) {
	state := DataFeedsOnChainState{
		Chains: make(map[uint64]DataFeedsChainState),
	}
	for chainSelector, chain := range e.BlockChains.EVMChains() {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := LoadChainState(e.Logger, chain, addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = *chainState
	}
	return state, nil
}

// LoadChainState Loads all state for a chain into state
func LoadChainState(logger logger.Logger, chain cldf_evm.Chain, addresses map[string]cldf.TypeAndVersion) (*DataFeedsChainState, error) {
	var state DataFeedsChainState

	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	state.MCMSWithTimelockState = *mcmsWithTimelock

	dfCacheTV := cldf.NewTypeAndVersion(DataFeedsCache, deployment.Version1_0_0)
	dfCacheTV.Labels.Add("data-feeds")

	devPlatformCacheTV := cldf.NewTypeAndVersion(DataFeedsCache, deployment.Version1_0_0)
	devPlatformCacheTV.Labels.Add("dev-platform")

	state.DataFeedsCache = make(map[common.Address]*cache.DataFeedsCache)
	state.AggregatorProxy = make(map[common.Address]*proxy.AggregatorProxy)
	state.ABIByAddress = make(map[string]string)

	for address, tv := range addresses {
		switch {
		case tv.String() == dfCacheTV.String() || tv.String() == devPlatformCacheTV.String():
			contract, err := cache.NewDataFeedsCache(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &state, err
			}
			state.DataFeedsCache[common.HexToAddress(address)] = contract
			state.ABIByAddress[address] = cache.DataFeedsCacheABI
		case strings.Contains(tv.String(), "AggregatorProxy"):
			contract, err := proxy.NewAggregatorProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &state, err
			}
			state.AggregatorProxy[common.HexToAddress(address)] = contract
			state.ABIByAddress[address] = proxy.AggregatorProxyABI
		case tv.String() == cldf.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.RBACTimelockABI
		case tv.String() == cldf.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.CallProxyABI
		case tv.String() == cldf.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0).String() || tv.String() == cldf.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0).String() || tv.String() == cldf.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
		default:
			logger.Warnw("unknown contract type", "type", tv.Type)
		}
	}
	return &state, nil
}

// LoadAptosChainState Loads all state for aptos chain into state
func LoadAptosChainState(logger logger.Logger, chain cldf_aptos.Chain, addresses []datastore.AddressRef) (*DataFeedsAptosChainState, error) {
	var state DataFeedsAptosChainState

	state.DataFeeds = make(map[aptos.AccountAddress]*modulefeeds.DataFeeds)

	for _, address := range addresses {
		if address.Type == ChainlinkDataFeedsPackage {
			feedsAddress := aptos.AccountAddress{}
			err := feedsAddress.ParseStringRelaxed(address.Address)
			if err != nil {
				return &state, fmt.Errorf("failed to parse address %s: %w", address.Address, err)
			}

			bindContract := modulefeeds.Bind(feedsAddress, chain.Client)
			state.DataFeeds[feedsAddress] = &bindContract
		}
	}
	return &state, nil
}

func (s DataFeedsOnChainState) View(chains []uint64, e cldf.Environment) (map[string]view.ChainView, error) {
	m := make(map[string]view.ChainView)
	for _, chainSelector := range chains {
		chainInfo, err := cldf_chain_utils.ChainInfo(chainSelector)
		if err != nil {
			return m, err
		}
		if _, ok := s.Chains[chainSelector]; !ok {
			return m, fmt.Errorf("chain not supported %d", chainSelector)
		}
		chainState := s.Chains[chainSelector]
		chainView, err := chainState.GenerateView()
		if err != nil {
			return m, err
		}
		name := chainInfo.ChainName
		if chainInfo.ChainName == "" {
			name = strconv.FormatUint(chainSelector, 10)
		}
		chainView.FeedConfig = *GenerateFeedConfigView(e, name)
		m[name] = chainView
	}
	return m, nil
}

func GenerateFeedConfigView(e cldf.Environment, chainName string) *v1_0.FeedState {
	baseDir := ".."
	envName := e.Name

	filePath := filepath.Join(baseDir, envName, "inputs", "feeds", chainName+".json")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		e.Logger.Warnf("File %s does not exist", filePath)
		return &v1_0.FeedState{}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		e.Logger.Warnf("Error reading file %s: %v", filePath, err)
		return &v1_0.FeedState{}
	}

	var feedsView v1_0.FeedState

	err = json.Unmarshal(content, &feedsView)
	if err != nil {
		e.Logger.Warnf("Error unmarshalling file %s: %v", filePath, err)
	}
	return &feedsView
}

func (c DataFeedsChainState) GenerateView() (view.ChainView, error) {
	chainView := view.NewChain()
	if c.DataFeedsCache != nil {
		for _, cacheContract := range c.DataFeedsCache {
			cacheView, err := v1_0.GenerateDataFeedsCacheView(cacheContract)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate cache view %s", cacheContract.Address().String())
			}
			chainView.DataFeedsCache[cacheContract.Address().Hex()] = cacheView
		}
	}
	if c.AggregatorProxy != nil {
		for _, proxyContract := range c.AggregatorProxy {
			proxyView, err := v1_0.GenerateAggregatorProxyView(proxyContract)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate proxy view %s", proxyContract.Address().String())
			}
			chainView.AggregatorProxy[proxyContract.Address().Hex()] = proxyView
		}
	}
	return chainView, nil
}
