package solana

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
)

// use this changeset to add a token pool and lookup table
var _ cldf.ChangeSet[E2ETokenPoolConfig] = E2ETokenPool

type E2ETokenPoolConfig struct {
	AddTokenPoolAndLookupTable            []TokenPoolConfig
	AddTokenPoolLookupTable               []TokenPoolLookupTableConfig
	RegisterTokenAdminRegistry            []RegisterTokenAdminRegistryConfig
	AcceptAdminRoleTokenAdminRegistry     []AcceptAdminRoleTokenAdminRegistryConfig
	SetPool                               []SetPoolConfig
	RemoteChainTokenPool                  []RemoteChainTokenPoolConfig
	ConfigureTokenPoolContractsChangesets []v1_5_1.ConfigureTokenPoolContractsConfig
}

func E2ETokenPool(e cldf.Environment, cfg E2ETokenPoolConfig) (cldf.ChangesetOutput, error) {
	finalOutput := cldf.ChangesetOutput{}
	finalOutput.AddressBook = cldf.NewMemoryAddressBook() //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	defer func(e cldf.Environment) {
		e.Logger.Info("SolanaE2ETokenPool changeset completed")
		e.Logger.Info("Final output: ", finalOutput.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	}(e)

	for _, tokenPoolConfig := range cfg.AddTokenPoolAndLookupTable {
		output, err := AddTokenPoolAndLookupTable(e, tokenPoolConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add token pool lookup table: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, tokenPoolConfig := range cfg.AddTokenPoolLookupTable {
		output, err := AddTokenPoolLookupTable(e, tokenPoolConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add token pool lookup table: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, registerConfig := range cfg.RegisterTokenAdminRegistry {
		output, err := RegisterTokenAdminRegistry(e, registerConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to register token admin registry: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, acceptConfig := range cfg.AcceptAdminRoleTokenAdminRegistry {
		output, err := AcceptAdminRoleTokenAdminRegistry(e, acceptConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to accept admin role token admin registry: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, setPoolConfig := range cfg.SetPool {
		output, err := SetPool(e, setPoolConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set pool: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, remoteChainConfig := range cfg.RemoteChainTokenPool {
		output, err := SetupTokenPoolForRemoteChain(e, remoteChainConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to remote chain token pool config: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, configureTokenPoolConfig := range cfg.ConfigureTokenPoolContractsChangesets {
		output, err := v1_5_1.ConfigureTokenPoolContractsChangeset(e, configureTokenPoolConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure token pool contracts: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
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
	defer func(e cldf.Environment) {
		e.Logger.Info("E2EToken changeset completed")
		e.Logger.Info("Final output: ", finalOutput.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	}(e)

	for _, config := range cfg.DeploySolanaToken {
		output, err := DeploySolanaToken(e, config)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy solana token: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, config := range cfg.UploadTokenMetadata {
		output, err := UploadTokenMetadata(e, config)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to upload token metadata: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}
	for _, config := range cfg.SetTokenAuthority {
		output, err := SetTokenAuthority(e, config)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to register token admin registry: %w", err)
		}
		if output.AddressBook != nil { //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			err = finalOutput.AddressBook.Merge(output.AddressBook) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
			}
		}
		if len(output.MCMSTimelockProposals) > 0 {
			finalOutput.MCMSTimelockProposals = append(finalOutput.MCMSTimelockProposals, output.MCMSTimelockProposals...)
		}
	}

	return finalOutput, nil
}
