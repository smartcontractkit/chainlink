package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// NewFeedWithProxyChangeset configures new feeds with a proxy addresses
// 1. Deploys AggregatorProxy contracts for given chainselector with DataFeedsCache as an aggregator
// 2. Creates an MCMS proposal to transfer the ownership of AggregatorProxy contracts to timelock
// 3. Creates a proposal to set a feed configs on DataFeedsCache contract
// 4. Creates a proposal to set a feed proxy mappings on DataFeedsCache contract
// Returns a new addressbook with the new AggregatorProxy contracts address and MCMS proposal
var NewFeedWithProxyChangeset = cldf.CreateChangeSet(newFeedWithProxyLogic, newFeedWithProxyPrecondition)

func newFeedWithProxyLogic(env deployment.Environment, c types.NewFeedWithProxyConfig) (deployment.ChangesetOutput, error) {
	chain := env.Chains[c.ChainSelector]
	state, _ := LoadOnchainState(env)
	chainState := state.Chains[c.ChainSelector]
	ab := deployment.NewMemoryAddressBook()

	dataFeedsCacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, c.ChainSelector, nil)
	if dataFeedsCacheAddress == "" {
		return deployment.ChangesetOutput{}, fmt.Errorf("DataFeedsCache contract address not found in addressbook for chain %d", c.ChainSelector)
	}

	dataFeedsCache := chainState.DataFeedsCache[common.HexToAddress(dataFeedsCacheAddress)]
	if dataFeedsCache == nil {
		return deployment.ChangesetOutput{}, errors.New("DataFeedsCache contract not found in onchain state")
	}

	var proxyAddresses []common.Address
	var acceptProxyOwnerShipProposals []ProposalData

	// For each Data ID, deploy an AggregatorProxy contract with DataFeedsCache as an aggregator on it.
	// Transfer ownership to timelock and create accept ownership proposal
	for index := range c.DataIDs {
		// Deploy AggregatorProxy contract with deployer key
		proxyConfig := types.DeployAggregatorProxyConfig{
			ChainsToDeploy:   []uint64{c.ChainSelector},
			AccessController: []common.Address{c.AccessController},
			Labels:           append([]string{c.Descriptions[index]}, c.Labels...),
		}
		newEnv, err := RunChangeset(DeployAggregatorProxyChangeset, env, proxyConfig)

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute DeployAggregatorProxyChangeset: %w", err)
		}
		proxyAddress, err := deployment.SearchAddressBook(newEnv.AddressBook, c.ChainSelector, "AggregatorProxy")
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("AggregatorProxy not present in addressbook: %w", err)
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

		err = ab.Merge(newEnv.AddressBook)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge addressbooks: %w", err)
		}

		// accept proxy ownership proposal
		acceptProxyOwnerShipTx, err := proxyContract.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create accept transfer ownership tx %w", err)
		}
		acceptProxyOwnerShipProposals = append(acceptProxyOwnerShipProposals, ProposalData{
			contract: proxyContract.Address().Hex(),
			tx:       acceptProxyOwnerShipTx,
		})

		proxyAddresses = append(proxyAddresses, proxyContract.Address())
	}

	dataIDs, _ := FeedIDsToBytes16(c.DataIDs)

	// set feed config proposal
	setFeedConfigTx, err := dataFeedsCache.SetDecimalFeedConfigs(deployment.SimTransactOpts(), dataIDs, c.Descriptions, c.WorkflowMetadata)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set feed config %w", err)
	}

	// set feed proxy mapping proposal
	setProxyMappingTx, err := dataFeedsCache.UpdateDataIdMappingsForProxies(deployment.SimTransactOpts(), proxyAddresses, dataIDs)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set proxy-dataId mapping %w", err)
	}

	proposalConfig := MultiChainProposalConfig{
		c.ChainSelector: []ProposalData{
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
	proposalConfig[c.ChainSelector] = append(proposalConfig[c.ChainSelector], acceptProxyOwnerShipProposals...)

	proposals, err := BuildMultiChainProposals(env, "accept AggregatorProxies ownership to timelock. set feed config and proxy mapping on cache", proposalConfig, c.McmsConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: ab, MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposals}}, nil
}

func newFeedWithProxyPrecondition(env deployment.Environment, c types.NewFeedWithProxyConfig) error {
	if c.McmsConfig == nil {
		return errors.New("mcms config is required")
	}
	if len(c.DataIDs) != len(c.Descriptions) {
		return errors.New("data ids and descriptions length mismatch")
	}
	_, err := FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	return ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector)
}
