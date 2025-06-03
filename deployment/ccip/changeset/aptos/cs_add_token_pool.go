package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var _ cldf.ChangeSetV2[config.AddTokenPoolConfig] = AddTokenPool{}

// AddTokenPool deploys token pools and sets up tokens on lanes
type AddTokenPool struct{}

func (cs AddTokenPool) VerifyPreconditions(env cldf.Environment, cfg config.AddTokenPoolConfig) error {
	return nil
}

func (cs AddTokenPool) Apply(env cldf.Environment, cfg config.AddTokenPoolConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChain := env.AptosChains[cfg.ChainSelector]
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
	if cfg.TokenObjAddress == (aptos.AccountAddress{}) {
		deploySeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployAptosTokenSequence, deps, cfg.TokenParams)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tokenObjectAddress = deploySeq.Output.TokenObjAddress
		seqReports = append(seqReports, deploySeq.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, deploySeq.Output.MCMSOperations...)
		// TODO: Save token address and object address to address book
		// tokenAddress = deploySeq.Output.TokenAddress
	}

	// Deploy Aptos token pool
	tokenPoolAddress := cfg.TokenPoolAddress
	if cfg.TokenPoolAddress != (aptos.AccountAddress{}) {
		depInput := seq.DeployTokenPoolSeqInput{
			TokenObjAddress: tokenObjectAddress,
			PoolType:        cfg.PoolType,
		}
		deploySeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployAptosTokenPoolSequence, deps, depInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		seqReports = append(seqReports, deploySeq.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, deploySeq.Output...)
		tokenPoolAddress = tokenObjectAddress
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
