package solana

import (
	"fmt"

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
	AddTokenPoolAndLookupTable            []TokenPoolConfig
	RegisterTokenAdminRegistry            []RegisterTokenAdminRegistryConfig
	AcceptAdminRoleTokenAdminRegistry     []AcceptAdminRoleTokenAdminRegistryConfig
	SetPool                               []SetPoolConfig
	RemoteChainTokenPool                  []RemoteChainTokenPoolConfig
	ConfigureTokenPoolContractsChangesets []v1_5_1.ConfigureTokenPoolContractsConfig
	MCMS                                  *proposalutils.TimelockConfig // set it to aggregate all the proposals
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

type E2ETokenConfig struct {
	DeploySolanaToken   []DeploySolanaTokenConfig
	UploadTokenMetadata []UploadTokenMetadataConfig
	SetTokenAuthority   []SetTokenAuthorityConfig
}

func E2EToken(e cldf.Environment, cfg E2ETokenConfig) (cldf.ChangesetOutput, error) {
	finalOutput := cldf.ChangesetOutput{}
	finalOutput.AddressBook = cldf.NewMemoryAddressBook() //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	addressBookToRemove := cldf.NewMemoryAddressBook()
	defer func(e cldf.Environment) {
		e.Logger.Info("E2EToken changeset completed")
		e.Logger.Info("Final output: ", finalOutput.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	}(e)
	err := ProcessConfig(&e, cfg.DeploySolanaToken, DeploySolanaToken, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy solana token: %w", err)
	}
	err = ProcessConfig(&e, cfg.UploadTokenMetadata, UploadTokenMetadata, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to upload token metadata: %w", err)
	}
	err = ProcessConfig(&e, cfg.SetTokenAuthority, SetTokenAuthority, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set token authority: %w", err)
	}
	err = AggregateAndCleanup(e, &finalOutput, addressBookToRemove, nil, "E2EToken changeset")
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
		proposal, err := proposalutils.AggregateProposals(
			e, state.EVMMCMSStateByChain(), state.SolanaMCMSStateByChain(e),
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
