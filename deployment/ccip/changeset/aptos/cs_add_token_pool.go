package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

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
		return fmt.Errorf("failed to load existing Aptos onchain state: %w", err)
	}
	// Validate supported chain
	supportedChains := state.SupportedChains()
	if _, ok := supportedChains[cfg.ChainSelector]; !ok {
		return fmt.Errorf("chain is not a supported: %d", cfg.ChainSelector)
	}
	// Validate CCIP deployed
	if state.AptosChains[cfg.ChainSelector].CCIPAddress == (aptos.AccountAddress{}) {
		return fmt.Errorf("CCIP is not deployed on Aptos chain %d", cfg.ChainSelector)
	}
	// Validate MCMS config
	if cfg.MCMSConfig == nil {
		return errors.New("MCMS config is required for AddTokenPool changeset")
	}
	// Validate config.TokenParams
	err = cfg.TokenParams.Validate()
	if err != nil {
		return fmt.Errorf("invalid token parameters: %w", err)
	}
	// Validate if token address is provided if pool address is specified
	if cfg.TokenObjAddress == (aptos.AccountAddress{}) && cfg.TokenPoolAddress != (aptos.AccountAddress{}) {
		return errors.New("token object address must be provided if token pool address is specified")
	}
	// No token pool address provided, so no need to validate token address
	if cfg.TokenObjAddress == (aptos.AccountAddress{}) && cfg.TokenPoolAddress == (aptos.AccountAddress{}) {
		return nil
	}
	// Validate if token already exists with different pool address
	for token, pool := range state.AptosChains[cfg.ChainSelector].AptosManagedTokenPools {
		if (token == cfg.TokenAddress) && (pool != cfg.TokenPoolAddress) {
			return fmt.Errorf("token %s already exists with a different pool address %s", token, pool)
		}
		if (pool == cfg.TokenPoolAddress) && (token != cfg.TokenAddress) {
			return fmt.Errorf("pool %s already exists with a different token address %s", pool, token)
		}
	}
	for token, pool := range state.AptosChains[cfg.ChainSelector].BurnMintTokenPools {
		if (token == cfg.TokenAddress) && (pool != cfg.TokenPoolAddress) {
			return fmt.Errorf("token %s already exists with a different pool address %s", token, pool)
		}
		if (pool == cfg.TokenPoolAddress) && (token != cfg.TokenAddress) {
			return fmt.Errorf("pool %s already exists with a different token address %s", pool, token)
		}
	}
	for token, pool := range state.AptosChains[cfg.ChainSelector].LockReleaseTokenPools {
		if (token == cfg.TokenAddress) && (pool != cfg.TokenPoolAddress) {
			return fmt.Errorf("token %s already exists with a different pool address %s", token, pool)
		}
		if (pool == cfg.TokenPoolAddress) && (token != cfg.TokenAddress) {
			return fmt.Errorf("pool %s already exists with a different token address %s", pool, token)
		}
	}
	return nil
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
	tokenObjectAddress := cfg.TokenObjAddress
	tokenAddress := cfg.TokenAddress
	tokenOwnerAddress := aptos.AccountAddress{}
	if cfg.TokenObjAddress == (aptos.AccountAddress{}) {
		deployTokenIn := seq.DeployTokenSeqInput{
			TokenParams: cfg.TokenParams,
			MCMSAddress: state.AptosChains[cfg.ChainSelector].MCMSAddress,
		}
		deploySeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployAptosTokenSequence, deps, deployTokenIn)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tokenObjectAddress = deploySeq.Output.TokenObjAddress
		tokenAddress = deploySeq.Output.TokenAddress
		tokenOwnerAddress = deploySeq.Output.TokenOwnerAddress
		seqReports = append(seqReports, deploySeq.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, deploySeq.Output.MCMSOperations...)
		// Save token object address in address book
		typeAndVersion := cldf.NewTypeAndVersion(shared.AptosManagedTokenType, deployment.Version1_6_0)
		typeAndVersion.AddLabel(string(cfg.TokenParams.Symbol))
		err = deps.AB.Save(deps.AptosChain.Selector, deploySeq.Output.TokenObjAddress.String(), typeAndVersion)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save token object address %s: %w", deploySeq.Output.TokenObjAddress, err)
		}
		// Save token address in address book
		typeAndVersion = cldf.NewTypeAndVersion(cldf.ContractType(cfg.TokenParams.Symbol), deployment.Version1_6_0)
		err = deps.AB.Save(deps.AptosChain.Selector, deploySeq.Output.TokenAddress.String(), typeAndVersion)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save token address %s: %w", deploySeq.Output.TokenAddress, err)
		}
	}

	// Deploy Aptos token pool
	tokenPoolAddress := cfg.TokenPoolAddress
	if cfg.TokenPoolAddress == (aptos.AccountAddress{}) {
		depInput := seq.DeployTokenPoolSeqInput{
			TokenObjAddress:   tokenObjectAddress,
			TokenAddress:      tokenAddress,
			TokenOwnerAddress: tokenOwnerAddress,
			PoolType:          cfg.PoolType,
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
		typeAndVersion.AddLabel(tokenAddress.String())
		err = deps.AB.Save(deps.AptosChain.Selector, deploySeq.Output.TokenPoolAddress.String(), typeAndVersion)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save token pool address %s: %w", deploySeq.Output.TokenPoolAddress, err)
		}
	}

	// Connect token pools EVM -> Aptos
	connInput := seq.ConnectTokenPoolSeqInput{
		TokenPoolAddress: tokenPoolAddress,
		RemotePools:      toRemotePools(cfg.EVMRemoteConfigs),
	}
	connectSeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.ConnectTokenPoolSequence, deps, connInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	seqReports = append(seqReports, connectSeq.ExecutionReports...)
	mcmsOperations = append(mcmsOperations, connectSeq.Output)

	// Generate Aptos MCMS proposals
	proposal, err := utils.GenerateProposal(
		aptosChain.Client,
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
			RemoteTokenAddress: remoteConfig.TokenAddress.Bytes(),
			RateLimiterConfig:  remoteConfig.RateLimiterConfig,
		}
	}
	return remotePools
}
