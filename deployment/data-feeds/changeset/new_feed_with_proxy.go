package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// NewFeedWithProxyChangeset configures a new feed with a proxy
// 1. Deploys AggregatorProxy contract for given chainselector
// 2. Proposes and confirms DataFeedsCache contract as an aggregator on AggregatorProxy
// 3. Creates an MCMS proposal to transfer the ownership of AggregatorProxy contract to timelock
// 4. Creates a proposal to set a feed config on DataFeedsCache contract
// 5. Creates a proposal to set a feed proxy mapping on DataFeedsCache contract
// Returns a new addressbook with the new AggregatorProxy contract address and 3 MCMS proposals
var NewFeedWithProxyChangeset = deployment.CreateChangeSet(newFeedWithProxyLogic, newFeedWithProxyPrecondition)

func newFeedWithProxyLogic(env deployment.Environment, c types.NewFeedWithProxyConfig) (deployment.ChangesetOutput, error) {
	chain := env.Chains[c.ChainSelector]
	state, _ := LoadOnchainState(env)
	chainState := state.Chains[c.ChainSelector]

	// Deploy AggregatorProxy contract with deployer key
	proxyConfig := types.DeployAggregatorProxyConfig{
		ChainsToDeploy:   []uint64{c.ChainSelector},
		AccessController: []common.Address{c.AccessController},
		Labels:           c.Labels,
	}
	newEnv, err := DeployAggregatorProxyChangeset.Apply(env, proxyConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute DeployAggregatorProxyChangeset: %w", err)
	}

	proxyAddress, err := deployment.SearchAddressBook(newEnv.AddressBook, c.ChainSelector, "AggregatorProxy")
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("AggregatorProxy not present in addressbook: %w", err)
	}

	addressMap, _ := env.ExistingAddresses.AddressesForChain(c.ChainSelector)
	var dataFeedsCacheAddress string
	cacheTV := deployment.NewTypeAndVersion(DataFeedsCache, deployment.Version1_0_0)
	cacheTV.Labels.Add("data-feeds")
	for addr, tv := range addressMap {
		if tv.String() == cacheTV.String() {
			dataFeedsCacheAddress = addr
		}
	}

	dataFeedsCache := chainState.DataFeedsCache[common.HexToAddress(dataFeedsCacheAddress)]
	if dataFeedsCache == nil {
		return deployment.ChangesetOutput{}, errors.New("DataFeedsCache contract not found in onchain state")
	}

	// Propose and confirm DataFeedsCache contract as an aggregator on AggregatorProxy
	proposeAggregatorConfig := types.ProposeConfirmAggregatorConfig{
		ChainSelector:        c.ChainSelector,
		ProxyAddress:         common.HexToAddress(proxyAddress),
		NewAggregatorAddress: common.HexToAddress(dataFeedsCacheAddress),
	}

	_, err = ProposeAggregatorChangeset.Apply(env, proposeAggregatorConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute ProposeAggregatorChangeset: %w", err)
	}

	_, err = ConfirmAggregatorChangeset.Apply(env, proposeAggregatorConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute ConfirmAggregatorChangeset: %w", err)
	}

	// Create an MCMS proposal to transfer the ownership of AggregatorProxy contract to timelock and set the feed configs
	// We don't use the existing changesets so that we can batch the transactions into a single MCMS proposal

	// transfer proxy ownership
	timelockAddr, _ := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, commonTypes.RBACTimelock)
	_, proxyContract, err := changeset.LoadOwnableContract(common.HexToAddress(proxyAddress), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load proxy contract %w", err)
	}
	tx, err := proxyContract.TransferOwnership(chain.DeployerKey, common.HexToAddress(timelockAddr))
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	// accept proxy ownership proposal
	acceptProxyOwnerShipTx, err := proxyContract.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create accept transfer ownership tx %w", err)
	}

	// set feed config proposal
	setFeedConfigTx, err := dataFeedsCache.SetDecimalFeedConfigs(deployment.SimTransactOpts(), [][16]byte{c.DataID}, []string{c.Description}, c.WorkflowMetadata)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set feed config %w", err)
	}

	// set feed proxy mapping proposal
	setProxyMappingTx, err := dataFeedsCache.UpdateDataIdMappingsForProxies(deployment.SimTransactOpts(), []common.Address{common.HexToAddress(proxyAddress)}, [][16]byte{c.DataID})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set proxy-dataId mapping %w", err)
	}

	proposalConfig := MultiChainProposalConfig{
		c.ChainSelector: []ProposalData{
			{
				contract: proxyContract.Address().Hex(),
				tx:       acceptProxyOwnerShipTx,
			},
			{
				contract: dataFeedsCache.Address().Hex(),
				tx:       setFeedConfigTx,
			},
			{
				contract: dataFeedsCache.Address().Hex(),
				tx:       setProxyMappingTx,
			},
		},
	}

	proposals, err := BuildMultiChainProposals(env, "accept AggregatorProxy ownership to timelock. set feed config and proxy mapping on cache", proposalConfig, c.McmsConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: newEnv.AddressBook, MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposals}}, nil
}

func newFeedWithProxyPrecondition(env deployment.Environment, c types.NewFeedWithProxyConfig) error {
	if c.McmsConfig == nil {
		return errors.New("mcms config is required")
	}

	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	return ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector)
}
