package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
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
	InitializeGlobalTokenPoolConfig       []TokenPoolConfigWithMCM
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
	err := ProcessConfig(&e, cfg.InitializeGlobalTokenPoolConfig, InitGlobalConfigTokenPoolProgram, &finalOutput, addressBookToRemove)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize global config for token pool: %w", err)
	}
	err = ProcessConfig(&e, cfg.AddTokenPoolAndLookupTable, AddTokenPoolAndLookupTable, &finalOutput, addressBookToRemove)
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

func E2ETokenPoolv2(env cldf.Environment, cfg E2ETokenPoolConfigv2) (cldf.ChangesetOutput, error) {
	// use a clone of env to avoid modifying the original env
	// if you modify the original env, the in memory tests will fail
	// because after the changeset is complete and the ApplyChangesets function is called,
	// it will try to add addresses from the cs output which already exist in the env
	e := env.Clone()
	finalCSOut := &cldf.ChangesetOutput{
		AddressBook: cldf.NewMemoryAddressBook(),
	}

	// token pool and lookup table
	tokenPoolAndLookupTableCfg := AddTokenPoolAndLookupTableConfig{
		ChainSelector:    cfg.ChainSelector,
		TokenPoolConfigs: make([]TokenPoolConfig, 0),
	}

	// register token admin registry
	registerTokenAdminRegistryCfg := RegisterTokenAdminRegistryConfig{
		ChainSelector:        cfg.ChainSelector,
		MCMS:                 cfg.MCMS,
		RegisterTokenConfigs: make([]RegisterTokenConfig, 0),
	}
	var tokenAdminRegistryAdmin solana.PublicKey
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	if cfg.MCMS != nil {
		tokenAdminRegistryAdmin = timelockSignerPDA
	} else {
		tokenAdminRegistryAdmin = e.BlockChains.SolanaChains()[cfg.ChainSelector].DeployerKey.PublicKey()
	}

	// accept admin role token admin registry
	acceptAdminRoleTokenAdminRegistryCfg := AcceptAdminRoleTokenAdminRegistryConfig{
		ChainSelector:               cfg.ChainSelector,
		MCMS:                        cfg.MCMS,
		AcceptAdminRoleTokenConfigs: make([]AcceptAdminRoleTokenConfig, 0),
	}

	// set pool
	setPoolCfg := SetPoolConfig{
		ChainSelector:       cfg.ChainSelector,
		SetPoolTokenConfigs: make([]SetPoolTokenConfig, 0),
		MCMS:                cfg.MCMS,
	}

	// solana to evm remote pool setup
	remotePoolConfig := SetupTokenPoolForRemoteChainConfig{
		SolChainSelector:       cfg.ChainSelector,
		RemoteTokenPoolConfigs: make([]RemoteChainTokenPoolConfig, 0),
		MCMS:                   cfg.MCMS,
	}

	// evm to solana remote pool setup
	evmToSolanaRemotePoolCfg := v1_5_1.ConfigureMultipleTokenPoolsConfig{
		MCMS: cfg.MCMS,
	}

	// transfer away the pool to timelock
	var transferPoolToTimelockConfig TransferCCIPToMCMSWithTimelockSolanaConfig
	if cfg.MCMS != nil {
		transferPoolToTimelockConfig = TransferCCIPToMCMSWithTimelockSolanaConfig{
			MCMSCfg: *cfg.MCMS,
			ContractsByChain: map[uint64]CCIPContractsToTransfer{
				cfg.ChainSelector: {
					BurnMintTokenPools:    map[string][]solana.PublicKey{},
					LockReleaseTokenPools: map[string][]solana.PublicKey{},
					CCTPTokenPoolMints:    []solana.PublicKey{},
				},
			},
		}
	}
	poolsByType := transferPoolToTimelockConfig.ContractsByChain[cfg.ChainSelector]

	var uniquePoolTypeConfigs []E2ETokenConfig
	for _, tokenCfg := range cfg.E2ETokens {
		if err := tokenCfg.Validate(); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate token config: %w, token cfg: %+v", err, tokenCfg)
		}
		// add token pool and lookup table
		tokenPoolAndLookupTableCfg.TokenPoolConfigs = append(tokenPoolAndLookupTableCfg.TokenPoolConfigs, TokenPoolConfig{
			PoolType:    tokenCfg.PoolType,
			Metadata:    tokenCfg.Metadata,
			TokenPubKey: tokenCfg.TokenPubKey,
		})
		// register token admin registry
		registerTokenAdminRegistryCfg.RegisterTokenConfigs = append(registerTokenAdminRegistryCfg.RegisterTokenConfigs, RegisterTokenConfig{
			TokenPubKey:             tokenCfg.TokenPubKey,
			RegisterType:            ViaGetCcipAdminInstruction,
			TokenAdminRegistryAdmin: tokenAdminRegistryAdmin,
		})
		// accept admin role token admin registry
		acceptAdminRoleTokenAdminRegistryCfg.AcceptAdminRoleTokenConfigs = append(acceptAdminRoleTokenAdminRegistryCfg.AcceptAdminRoleTokenConfigs, AcceptAdminRoleTokenConfig{
			TokenPubKey: tokenCfg.TokenPubKey,
			// registering in the same changeset so skip registry check
			SkipRegistryCheck: true,
		})
		// set pool
		setPoolCfg.SetPoolTokenConfigs = append(setPoolCfg.SetPoolTokenConfigs, SetPoolTokenConfig{
			TokenPubKey: tokenCfg.TokenPubKey,
			PoolType:    tokenCfg.PoolType,
			Metadata:    tokenCfg.Metadata,
			// registering in the same changeset so skip registry check
			SkipRegistryCheck: true,
		})
		// setup evm remote pool on solana
		if len(tokenCfg.SolanaToEVMRemoteConfigs) > 0 {
			remotePoolConfig.RemoteTokenPoolConfigs = append(remotePoolConfig.RemoteTokenPoolConfigs, RemoteChainTokenPoolConfig{
				SolTokenPubKey:   tokenCfg.TokenPubKey,
				SolPoolType:      tokenCfg.PoolType,
				Metadata:         tokenCfg.Metadata,
				EVMRemoteConfigs: tokenCfg.SolanaToEVMRemoteConfigs,
			})
		}
		// setup solana remote pool on evm
		if len(tokenCfg.EVMToSolanaRemoteConfigs.PoolUpdates) > 0 {
			evmToSolanaRemotePoolCfg.Tokens = append(evmToSolanaRemotePoolCfg.Tokens, &tokenCfg.EVMToSolanaRemoteConfigs)
		}
		// transfer pool to timelock
		if cfg.MCMS != nil {
			switch tokenCfg.PoolType {
			case shared.BurnMintTokenPool:
				poolsByType.BurnMintTokenPools[tokenCfg.Metadata] = append(
					poolsByType.BurnMintTokenPools[tokenCfg.Metadata],
					tokenCfg.TokenPubKey,
				)
			case shared.LockReleaseTokenPool:
				poolsByType.LockReleaseTokenPools[tokenCfg.Metadata] = append(
					poolsByType.LockReleaseTokenPools[tokenCfg.Metadata],
					tokenCfg.TokenPubKey,
				)
			case shared.CCTPTokenPool:
				poolsByType.CCTPTokenPoolMints = append(poolsByType.CCTPTokenPoolMints, tokenCfg.TokenPubKey)
			}
		}
		isUniquePoolType := true
		for _, uniqueCfg := range uniquePoolTypeConfigs {
			if uniqueCfg.PoolType == tokenCfg.PoolType {
				isUniquePoolType = false
			}
		}
		if isUniquePoolType {
			uniquePoolTypeConfigs = append(uniquePoolTypeConfigs, tokenCfg)
		}
	}
	// Initialize global configs once for each unique token pool
	for _, tokenCfg := range uniquePoolTypeConfigs {
		output, err := InitGlobalConfigTokenPoolProgram(e, TokenPoolConfigWithMCM{
			ChainSelector: cfg.ChainSelector,
			TokenPoolConfigs: []TokenPoolConfig{
				{
					PoolType:    tokenCfg.PoolType,
					TokenPubKey: tokenCfg.TokenPubKey,
					Metadata:    tokenCfg.Metadata,
				},
			},
			MCMS: cfg.MCMS,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize global config for token pool: %w", err)
		}
		if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running InitGlobalConfigTokenPoolProgram: %w", err)
		}
	}
	output, err := AddTokenPoolAndLookupTable(e, tokenPoolAndLookupTableCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to add token pool and lookup table: %w", err)
	}
	if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running AddTokenPoolAndLookupTable: %w", err)
	}
	output, err = SetupTokenPoolForRemoteChain(e, remotePoolConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to setup token pool for remote chain: %w", err)
	}
	if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running SetupTokenPoolForRemoteChain: %w", err)
	}
	output, err = RegisterTokenAdminRegistry(e, registerTokenAdminRegistryCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register token admin registry: %w", err)
	}
	if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running RegisterTokenAdminRegistry: %w", err)
	}
	output, err = AcceptAdminRoleTokenAdminRegistry(e, acceptAdminRoleTokenAdminRegistryCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to accept admin role: %w", err)
	}
	if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running AcceptAdminRoleTokenAdminRegistry: %w", err)
	}
	output, err = SetPool(e, setPoolCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set pool: %w", err)
	}
	if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running SetPool: %w", err)
	}
	output, err = v1_5_1.ConfigureMultiplePoolLogic(e, evmToSolanaRemotePoolCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure token pool contracts: %w", err)
	}
	if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running ConfigureTokenPoolContractsChangeset: %w", err)
	}
	// and finally lets transfer away the pool to timelock
	if cfg.MCMS != nil {
		output, err = TransferCCIPToMCMSWithTimelockSolana(e, transferPoolToTimelockConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ccip to mcms with timelock: %w", err)
		}
		if err = cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after running TransferCCIPToMCMSWithTimelockSolana: %w", err)
		}
	}

	if len(finalCSOut.MCMSTimelockProposals) > 1 {
		state, err := stateview.LoadOnchainState(e)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
		}
		proposal, err := proposalutils.AggregateProposalsV2(
			e, proposalutils.MCMSStates{
				MCMSEVMState:    state.EVMMCMSStateByChain(),
				MCMSSolanaState: state.SolanaMCMSStateByChain(e),
			},
			finalCSOut.MCMSTimelockProposals, "E2ETokenPoolv2 changeset", cfg.MCMS,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to aggregate proposals: %w", err)
		}
		if proposal != nil {
			finalCSOut.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
		}
	}

	return *finalCSOut, nil
}
