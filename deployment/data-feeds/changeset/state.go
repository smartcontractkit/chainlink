package changeset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

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
	DataFeedsCache deployment.ContractType = "DataFeedsCache"
)

type DataFeedsChainState struct {
	ABIByAddress map[string]string
	commonchangeset.MCMSWithTimelockState
	DataFeedsCache  map[common.Address]*cache.DataFeedsCache
	AggregatorProxy map[common.Address]*proxy.AggregatorProxy
}

type DataFeedsOnChainState struct {
	Chains map[uint64]DataFeedsChainState
}

func LoadOnchainState(e deployment.Environment) (DataFeedsOnChainState, error) {
	state := DataFeedsOnChainState{
		Chains: make(map[uint64]DataFeedsChainState),
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
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
func LoadChainState(logger logger.Logger, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*DataFeedsChainState, error) {
	var state DataFeedsChainState

	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	state.MCMSWithTimelockState = *mcmsWithTimelock

	dfCacheTV := deployment.NewTypeAndVersion(DataFeedsCache, deployment.Version1_0_0)
	dfCacheTV.Labels.Add("data-feeds")

	devPlatformCacheTV := deployment.NewTypeAndVersion(DataFeedsCache, deployment.Version1_0_0)
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
		case tv.String() == deployment.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.RBACTimelockABI
		case tv.String() == deployment.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.CallProxyABI
		case tv.String() == deployment.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0).String() || tv.String() == deployment.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0).String() || tv.String() == deployment.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
		default:
			logger.Warnw("unknown contract type", "type", tv.Type)
		}
	}
	return &state, nil
}

func (s DataFeedsOnChainState) View(chains []uint64, e deployment.Environment) (map[string]view.ChainView, error) {
	m := make(map[string]view.ChainView)
	for _, chainSelector := range chains {
		chainInfo, err := deployment.ChainInfo(chainSelector)
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

func GenerateFeedConfigView(e deployment.Environment, chainName string) *v1_0.FeedState {
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
