package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.AddTokenPoolConfig] = AddTokenPool{}

// AddTokenPool deploys token pools and sets up tokens on lanes
type AddTokenPool struct{}

func (cs AddTokenPool) VerifyPreconditions(env cldf.Environment, cfg config.AddTokenPoolConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	// Validate supported chain
	supportedChains := state.SupportedChains()
	if _, ok := supportedChains[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	// Validate CCIP deployed
	if state.AptosChains[cfg.ChainSelector].CCIPAddress == (aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("CCIP is not deployed on Aptos chain %d", cfg.ChainSelector))
	}
	// Validate MCMS config
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required for AddTokenPool changeset"))
	}
	// Validate config.TokenParams
	if cfg.TokenCodeObjAddress == (aptos.AccountAddress{}) {
		err = cfg.TokenParams.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid token parameters: %w", err))
		}
	}
	// Validate config.EVMRemoteConfigs
	for chainSelector, remoteConfig := range cfg.EVMRemoteConfigs {
		if err := remoteConfig.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid EVM remote config for chain %d: %w", chainSelector, err))
		}
	}
	// Validate if token address is provided if pool address is specified
	if cfg.TokenCodeObjAddress == (aptos.AccountAddress{}) && cfg.TokenPoolAddress != (aptos.AccountAddress{}) {
		errs = append(errs, errors.New("token object address must be provided if token pool address is specified"))
	}
	// No token pool address provided, so no need to validate token address
	if cfg.TokenCodeObjAddress == (aptos.AccountAddress{}) && cfg.TokenPoolAddress == (aptos.AccountAddress{}) {
		return errors.Join(errs...)
	}
	// Validate if token already exists with different pool address
	for token, pool := range state.AptosChains[cfg.ChainSelector].AptosManagedTokenPools {
		if (token == cfg.TokenAddress) && (pool != cfg.TokenPoolAddress) {
			errs = append(errs, fmt.Errorf("token %s already exists with a different pool address %s", token, pool))
		}
		if (pool == cfg.TokenPoolAddress) && (token != cfg.TokenAddress) {
			errs = append(errs, fmt.Errorf("pool %s already exists with a different token address %s", pool, token))
		}
	}
	for token, pool := range state.AptosChains[cfg.ChainSelector].BurnMintTokenPools {
		if (token == cfg.TokenAddress) && (pool != cfg.TokenPoolAddress) {
			errs = append(errs, fmt.Errorf("token %s already exists with a different pool address %s", token, pool))
		}
		if (pool == cfg.TokenPoolAddress) && (token != cfg.TokenAddress) {
			errs = append(errs, fmt.Errorf("pool %s already exists with a different token address %s", pool, token))
		}
	}
	for token, pool := range state.AptosChains[cfg.ChainSelector].LockReleaseTokenPools {
		if (token == cfg.TokenAddress) && (pool != cfg.TokenPoolAddress) {
			errs = append(errs, fmt.Errorf("token %s already exists with a different pool address %s", token, pool))
		}
		if (pool == cfg.TokenPoolAddress) && (token != cfg.TokenAddress) {
			errs = append(errs, fmt.Errorf("pool %s already exists with a different token address %s", pool, token))
		}
	}
	return errors.Join(errs...)
}

func (cs AddTokenPool) Apply(env cldf.Environment, cfg config.AddTokenPoolConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChain := env.BlockChains.AptosChains()[cfg.ChainSelector]
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)
	proposals := make([]mcms.TimelockProposal, 0)
	mcmsOperations := []mcmstypes.BatchOperation{}

	deps := operation.AptosDeps{
		AB:               ab,
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
	}

	// Deploy Aptos Token
	tokenCodeObjAddress := cfg.TokenCodeObjAddress
	tokenAddress := cfg.TokenAddress
	if cfg.TokenCodeObjAddress == (aptos.AccountAddress{}) {
		deployTokenIn := seq.DeployTokenSeqInput{
			TokenParams: cfg.TokenParams,
			MCMSAddress: state.AptosChains[cfg.ChainSelector].MCMSAddress,
			TokenMint:   cfg.TokenMint,
		}
		deploySeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployAptosTokenSequence, deps, deployTokenIn)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tokenCodeObjAddress = deploySeq.Output.TokenCodeObjAddress
		tokenAddress = deploySeq.Output.TokenAddress
		seqReports = append(seqReports, deploySeq.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, deploySeq.Output.MCMSOperations...)
		// Save token object address in address book
		typeAndVersion := cldf.NewTypeAndVersion(shared.AptosManagedTokenType, deployment.Version1_6_0)
		typeAndVersion.AddLabel(string(cfg.TokenParams.Symbol))
		err = deps.AB.Save(deps.AptosChain.Selector, deploySeq.Output.TokenCodeObjAddress.StringLong(), typeAndVersion)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save token object address %s: %w", deploySeq.Output.TokenCodeObjAddress, err)
		}
		// Save token address in address book
		typeAndVersion = cldf.NewTypeAndVersion(cldf.ContractType(cfg.TokenParams.Symbol), deployment.Version1_6_0)
		err = deps.AB.Save(deps.AptosChain.Selector, deploySeq.Output.TokenAddress.StringLong(), typeAndVersion)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save token address %s: %w", deploySeq.Output.TokenAddress, err)
		}
	}

	// Deploy Aptos token pool
	tokenPoolAddress := cfg.TokenPoolAddress
	if cfg.TokenPoolAddress == (aptos.AccountAddress{}) {
		depInput := seq.DeployTokenPoolSeqInput{
			TokenCodeObjAddress: tokenCodeObjAddress,
			TokenAddress:        tokenAddress,
			PoolType:            cfg.PoolType,
		}
		deploySeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployAptosTokenPoolSequence, deps, depInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		seqReports = append(seqReports, deploySeq.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, deploySeq.Output.MCMSOps...)
		tokenPoolAddress = deploySeq.Output.TokenPoolAddress
		// Save token pool address in address book
		typeAndVersion := cldf.NewTypeAndVersion(cfg.PoolType, deployment.Version1_6_0)
		typeAndVersion.AddLabel(tokenAddress.StringLong())
		err = deps.AB.Save(deps.AptosChain.Selector, deploySeq.Output.TokenPoolAddress.StringLong(), typeAndVersion)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save token pool address %s: %w", deploySeq.Output.TokenPoolAddress, err)
		}
	}

	// Connect token pools EVM -> Aptos
	connInput := seq.ConnectTokenPoolSeqInput{
		TokenPoolAddress:                    tokenPoolAddress,
		RemotePools:                         toRemotePools(cfg.EVMRemoteConfigs),
		TokenAddress:                        tokenAddress,
		TokenTransferFeeByRemoteChainConfig: cfg.TokenTransferFeeByRemoteChainConfig,
	}
	connectSeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.ConnectTokenPoolSequence, deps, connInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	seqReports = append(seqReports, connectSeq.ExecutionReports...)
	mcmsOperations = append(mcmsOperations, connectSeq.Output)

	// Generate Aptos MCMS proposals
	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		mcmsOperations,
		"Deploy and configure token pool on Aptos chain",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}
	proposals = append(proposals, *proposal)

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		MCMSTimelockProposals: proposals,
		Reports:               seqReports,
	}, nil
}

func toRemotePools(evmRemoteCfg map[uint64]config.EVMRemoteConfig) map[uint64]seq.RemotePool {
	remotePools := make(map[uint64]seq.RemotePool)
	for chainSelector, remoteConfig := range evmRemoteCfg {
		remotePools[chainSelector] = seq.RemotePool{
			RemotePoolAddress:  remoteConfig.TokenPoolAddress.Bytes(),
			RemoteTokenAddress: common.LeftPadBytes(remoteConfig.TokenAddress.Bytes(), 32),
			RateLimiterConfig:  remoteConfig.RateLimiterConfig,
		}
	}
	return remotePools
}
