package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this changeset to
// add a token pool and lookup table
// register the deployer key as the token admin to the token admin registry
// accept the admin role as the deployer key
// call setPool on the token admin registry
// configure evm pools on the solana side
// configure solana pools on the evm side
var _ cldf.ChangeSet[E2ETokenPoolConfig] = E2ETokenPool

type E2ETokenPoolConfig struct {
	AddTokenPoolAndLookupTable            []AddTokenPoolAndLookupTableConfig
	RegisterTokenAdminRegistry            []RegisterTokenAdminRegistryConfig
	AcceptAdminRoleTokenAdminRegistry     []AcceptAdminRoleTokenAdminRegistryConfig
	SetPool                               []SetPoolConfig
	RemoteChainTokenPool                  []SetupTokenPoolForRemoteChainConfig       // setup evm remote pools on solana
	ConfigureTokenPoolContractsChangesets []v1_5_1.ConfigureTokenPoolContractsConfig // setup evm/solana remote pools on evm
	MCMS                                  *proposalutils.TimelockConfig              // set it to aggregate all the proposals
}

func E2ETokenPool(e cldf.Environment, cfg E2ETokenPoolConfig) (cldf.ChangesetOutput, error) {
	finalOutput := cldf.ChangesetOutput{}
	finalOutput.AddressBook = cldf.NewMemoryAddressBook() //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	addressBookToRemove := cldf.NewMemoryAddressBook()
	defer func(e cldf.Environment) {
		e.Logger.Info("SolanaE2ETokenPool changeset completed")
		e.Logger.Info("Final output: ", finalOutput.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	}(e)
	// if mcms config is not provided, use the mcms config from one of the other configs
	if cfg.MCMS == nil {
		switch {
		case len(cfg.RegisterTokenAdminRegistry) > 0 && cfg.RegisterTokenAdminRegistry[0].MCMS != nil:
			cfg.MCMS = cfg.RegisterTokenAdminRegistry[0].MCMS
		case len(cfg.AcceptAdminRoleTokenAdminRegistry) > 0 && cfg.AcceptAdminRoleTokenAdminRegistry[0].MCMS != nil:
			cfg.MCMS = cfg.AcceptAdminRoleTokenAdminRegistry[0].MCMS
		case len(cfg.SetPool) > 0 && cfg.SetPool[0].MCMS != nil:
			cfg.MCMS = cfg.SetPool[0].MCMS
		}
	}
	err := ProcessConfig(&e, cfg.AddTokenPoolAndLookupTable, AddTokenPoolAndLookupTable, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to add token pool and lookup table: %w", err)
	}
	err = ProcessConfig(&e, cfg.RemoteChainTokenPool, SetupTokenPoolForRemoteChain, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure remote chain token pool: %w", err)
	}
	err = ProcessConfig(&e, cfg.RegisterTokenAdminRegistry, RegisterTokenAdminRegistry, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register token admin registry: %w", err)
	}
	err = ProcessConfig(&e, cfg.AcceptAdminRoleTokenAdminRegistry, AcceptAdminRoleTokenAdminRegistry, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to accept admin role: %w", err)
	}
	err = ProcessConfig(&e, cfg.SetPool, SetPool, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set pool: %w", err)
	}
	err = ProcessConfig(&e, cfg.ConfigureTokenPoolContractsChangesets, v1_5_1.ConfigureTokenPoolContractsChangeset, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure token pool contracts: %w", err)
	}
	err = AggregateAndCleanup(e, &finalOutput, addressBookToRemove, cfg.MCMS, "E2ETokenPool changeset")
	if err != nil {
		e.Logger.Error("failed to aggregate and cleanup: ", err)
	}

	return finalOutput, nil
}

func ProcessConfig[T any](
	e *cldf.Environment,
	configs []T,
	handler func(cldf.Environment, T) (cldf.ChangesetOutput, error),
	finalOutput *cldf.ChangesetOutput,
	tempRemoveBook cldf.AddressBook,
) error {
	for _, cfg := range configs {
		output, err := handler(*e, cfg)
		if err != nil {
			return err
		}
		err = cldf.MergeChangesetOutput(*e, finalOutput, output)
		if err != nil {
			return fmt.Errorf("failed to merge changeset output: %w", err)
		}

		if ab := output.AddressBook; ab != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err := tempRemoveBook.Merge(ab); err != nil {
				return fmt.Errorf("failed to merge into temp: %w", err)
			}
		}
	}
	return nil
}

func AggregateAndCleanup(e cldf.Environment, finalOutput *cldf.ChangesetOutput, abToRemove cldf.AddressBook, cfg *proposalutils.TimelockConfig, proposalDesc string) error {
	allProposals := finalOutput.MCMSTimelockProposals
	if len(allProposals) > 0 {
		state, err := stateview.LoadOnchainState(e)
		if err != nil {
			return fmt.Errorf("failed to load onchain state: %w", err)
		}
		proposal, err := proposalutils.AggregateProposalsV2(
			e, proposalutils.MCMSStates{
				MCMSEVMState:    state.EVMMCMSStateByChain(),
				MCMSSolanaState: state.SolanaMCMSStateByChain(e),
			},
			allProposals, proposalDesc, cfg,
		)
		if err != nil {
			return fmt.Errorf("failed to aggregate proposals: %w", err)
		}
		if proposal != nil {
			finalOutput.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
		}
	}
	if addresses, err := abToRemove.Addresses(); err == nil && len(addresses) > 0 {
		if err := e.ExistingAddresses.Remove(abToRemove); err != nil {
			return fmt.Errorf("failed to remove temp address book: %w", err)
		}
	}
	return nil
}

type E2ETokenConfig struct {
	TokenPubKey solana.PublicKey
	Metadata    string
	PoolType    cldf.ContractType
	// evm chain id -> evm remote config
	SolanaToEVMRemoteConfigs map[uint64]EVMRemoteConfig
	// solana remote config for evm pool
	EVMToSolanaRemoteConfigs v1_5_1.ConfigureTokenPoolContractsConfig
}

func (cfg E2ETokenConfig) Validate() error {
	if cfg.PoolType == "" {
		return errors.New("pool type is required")
	}
	if cfg.TokenPubKey.IsZero() {
		return errors.New("token pubkey is required")
	}
	if cfg.Metadata == "" {
		return errors.New("metadata is required")
	}
	return nil
}

type E2ETokenPoolConfigv2 struct {
	ChainSelector uint64
	E2ETokens     []E2ETokenConfig
	// this determines whether we want to set timelock as token admin or not
	// this is also required if router is owned by timelock
	// so you cannot really have a case where router is owned by timelock but you want to
	// set deployer key as token admin
	MCMS *proposalutils.TimelockConfig
}
