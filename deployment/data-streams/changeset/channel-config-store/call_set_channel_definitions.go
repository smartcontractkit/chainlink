package channel_config_store

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

type (
	SetChannelDefinitionsConfig struct {
		// DefinitionsByChain is a map of chain selectors -> ChannelDefinitions to deploy.
		DefinitionsByChain map[uint64][]ChannelDefinition
		MCMSConfig         *types.MCMSConfig
	}

	ChannelDefinition struct {
		ChannelConfigStore common.Address

		DonID uint32
		S3URL string
		Hash  [32]byte
	}

	ChannelConfigStoreState struct {
		ChannelConfigStore *channel_config_store.ChannelConfigStore
	}
)

func (a ChannelDefinition) GetContractAddress() common.Address {
	return a.ChannelConfigStore
}

// SetChannelDefinitionChangeset sets the channel definitions on the ChannelConfigStore contract
var SetChannelDefinitionChangeset = cldf.CreateChangeSet(callSetChannelDefinitions, callSetChannelDefinitionsPrecondition)

func (cfg SetChannelDefinitionsConfig) Validate() error {
	if len(cfg.DefinitionsByChain) == 0 {
		return errors.New("DefinitionsByChain cannot be empty")
	}
	return nil
}

func callSetChannelDefinitionsPrecondition(e cldf.Environment, cfg SetChannelDefinitionsConfig) error {
	if len(cfg.DefinitionsByChain) == 0 {
		return errors.New("DefinitionsByChain cannot be empty")
	}
	for chainSel := range cfg.DefinitionsByChain {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chainSel, err)
		}
	}
	return nil
}

func callSetChannelDefinitions(e cldf.Environment, cfg SetChannelDefinitionsConfig) (cldf.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.ChannelConfigStore.String(),
		cfg.DefinitionsByChain,
		maybeLoadChannelConfigStoreState,
		doSetChannelDefinitions,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building SetNativeSurcharge txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetNativeSurcharge proposal")
}

func maybeLoadChannelConfigStoreState(e cldf.Environment, chainSel uint64, contractAddr string) (*channel_config_store.ChannelConfigStore, error) {
	if err := utils.ValidateContract(e, chainSel, contractAddr, types.ChannelConfigStore, deployment.Version1_0_0); err != nil {
		return nil, err
	}
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}
	ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to load ChannelConfigStore contract: %w", err)
	}

	return ccs, nil
}

func doSetChannelDefinitions(
	ccs *channel_config_store.ChannelConfigStore,
	c ChannelDefinition,
) (*ethTypes.Transaction, error) {
	return ccs.SetChannelDefinitions(
		cldf.SimTransactOpts(),
		c.DonID,
		c.S3URL,
		c.Hash,
	)
}
