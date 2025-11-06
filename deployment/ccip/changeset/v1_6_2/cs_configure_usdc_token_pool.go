package v1_6_2

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	utp "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/usdc_token_pool"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

var (
	ConfigUSDCTokenPoolChangeSet = cldf.CreateChangeSet(configUSDCTokenPoolLogic, configUSDCTokenPoolPrecondition)

	USDCTokenPoolConfigOp = opsutil.NewEVMCallOperation(
		"USDCTokenPoolConfigOp",
		semver.MustParse("1.6.2"),
		"Setting USDC Token Pool config",
		utp.USDCTokenPoolABI,
		shared.USDCTokenPool,
		utp.NewUSDCTokenPool,
		func(tokenPool *utp.USDCTokenPool, opts *bind.TransactOpts, input []utp.USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
			return tokenPool.SetDomains(opts, input)
		})

	USDCTokenPoolConfigSequence = operations.NewSequence(
		"USDCTokenPoolConfigSequence",
		semver.MustParse("1.6.2"),
		"Setting USDC Token Pool config across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, inputs map[uint64]opsutil.EVMCallInput[[]utp.USDCTokenPoolDomainUpdate]) (map[uint64][]opsutil.EVMCallOutput, error) {
			out := make(map[uint64][]opsutil.EVMCallOutput, len(inputs))

			for chainSelector, input := range inputs {
				if _, ok := chains[chainSelector]; !ok {
					return nil, fmt.Errorf("chain with selector %d not defined in dependencies", chainSelector)
				}

				report, err := operations.ExecuteOperation(b, USDCTokenPoolConfigOp, chains[chainSelector], input)
				if err != nil {
					return map[uint64][]opsutil.EVMCallOutput{}, fmt.Errorf("failed to set USDC token pool config for chain %d: %w", chainSelector, err)
				}
				out[chainSelector] = []opsutil.EVMCallOutput{report.Output}
			}

			return out, nil
		})
)

type DomainUpdateInput struct {
	AllowedCaller    string
	MintRecipient    string
	DomainIdentifier uint32
	Enabled          bool
}
type ConfigUSDCTokenPoolInput struct {
	DestinationUpdates map[uint64]DomainUpdateInput
}

func (i ConfigUSDCTokenPoolInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {
	if _, ok := state.USDCTokenPoolsV1_6[deployment.Version1_6_2]; !ok {
		return fmt.Errorf("no USDC token pool with version %s found on %s", deployment.Version1_6_2, chain.Name())
	}
	for destSelector, update := range i.DestinationUpdates {
		err := cldf.IsValidChainSelector(destSelector)
		if err != nil {
			return fmt.Errorf("invalid destination chain selector %d: %w", destSelector, err)
		}

		fam, err := chain_selectors.GetSelectorFamily(destSelector)
		if err != nil {
			return fmt.Errorf("failed to get selector family for destination chain selector %d: %w", destSelector, err)
		}
		switch fam {
		case chain_selectors.FamilyEVM:
			allowedCallerAddr := common.HexToAddress(update.AllowedCaller)
			if allowedCallerAddr == utils.ZeroAddress {
				return fmt.Errorf("allowed caller must be defined for EVM destination chain selector %d", destSelector)
			}
		case chain_selectors.FamilySolana:
			allowedCallerAddr, err := solana.PublicKeyFromBase58(update.AllowedCaller)
			if err != nil {
				return fmt.Errorf("invalid allowed caller format %s for chain family %s", update.AllowedCaller, fam)
			}
			mintRecipientAddr, err := solana.PublicKeyFromBase58(update.MintRecipient)
			if err != nil {
				return fmt.Errorf("invalid mint recipient format %s for chain family %s", update.AllowedCaller, fam)
			}
			if mintRecipientAddr.IsZero() {
				return fmt.Errorf("mint recipient must be defined for Solana destination chain selector %d", destSelector)
			}
			if allowedCallerAddr.IsZero() {
				return fmt.Errorf("allowed caller must be defined for Solana destination chain selector %d", destSelector)
			}
		default:
			return fmt.Errorf("unsupported chain family: %s", fam)
		}

		// TODO: Other validations? Domain's are defined in chainlink-deployments so they can't be verified here...
	}
	return nil
}

type ConfigUSDCTokenPoolConfig struct {
	USDCPools map[uint64]ConfigUSDCTokenPoolInput

	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func configUSDCTokenPoolPrecondition(env cldf.Environment, c ConfigUSDCTokenPoolConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, poolConfig := range c.USDCPools {
		chain, chainState, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
		}

		err = poolConfig.Validate(env.GetContext(), chain, chainState)
		if err != nil {
			return fmt.Errorf("failed to validate USDC token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

func configUSDCTokenPoolLogic(env cldf.Environment, c ConfigUSDCTokenPoolConfig) (cldf.ChangesetOutput, error) {
	if err := configUSDCTokenPoolPrecondition(env, c); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid ConfigUSDCTokenPoolConfig: %w", err)
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Convert CLD/migrations inputs to onchain inputs.
	input := make(map[uint64]opsutil.EVMCallInput[[]utp.USDCTokenPoolDomainUpdate], len(c.USDCPools))
	for sourceChainSelector, poolConfig := range c.USDCPools {
		_, chainState, err := state.GetEVMChainState(env, sourceChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get EVM chain state for chain selector %d: %w",
				sourceChainSelector, err)
		}

		var domainUpdates []utp.USDCTokenPoolDomainUpdate
		for destSelector, update := range poolConfig.DestinationUpdates {
			fam, err := chain_selectors.GetSelectorFamily(destSelector)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get selector family for destination chain selector %d: %w", destSelector, err)
			}
			var allowedCaller [32]byte
			var mintRecipient [32]byte
			switch fam {
			case chain_selectors.FamilyEVM:
				allowedCallerAddr := common.HexToAddress(update.AllowedCaller)
				allowedCaller = [32]byte(common.LeftPadBytes(allowedCallerAddr.Bytes(), 32))
			case chain_selectors.FamilySolana:
				allowedCallerAddr, err := solana.PublicKeyFromBase58(update.AllowedCaller)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("invalid allowed caller format %s for chain family %s", update.AllowedCaller, fam)
				}
				allowedCaller = [32]byte(common.LeftPadBytes(allowedCallerAddr.Bytes(), 32))
				mintRecipientAddr, err := solana.PublicKeyFromBase58(update.MintRecipient)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("invalid mint recipient format %s for chain family %s", update.AllowedCaller, fam)
				}
				mintRecipient = [32]byte(common.LeftPadBytes(mintRecipientAddr.Bytes(), 32))
			default:
				return cldf.ChangesetOutput{}, fmt.Errorf("unsupported chain family: %s", fam)
			}
			domainUpdates = append(domainUpdates, utp.USDCTokenPoolDomainUpdate{
				AllowedCaller:     allowedCaller,
				MintRecipient:     mintRecipient,
				DomainIdentifier:  update.DomainIdentifier,
				DestChainSelector: destSelector,
				Enabled:           update.Enabled,
			})
		}

		input[sourceChainSelector] = opsutil.EVMCallInput[[]utp.USDCTokenPoolDomainUpdate]{
			ChainSelector: sourceChainSelector,
			NoSend:        c.MCMS != nil,
			Address:       chainState.USDCTokenPoolsV1_6[deployment.Version1_6_2].Address(),
			CallInput:     domainUpdates,
		}
	}

	// Configure sequence.
	seqReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		USDCTokenPoolConfigSequence,
		env.BlockChains.EVMChains(),
		input,
	)
	return opsutil.AddEVMCallSequenceToCSOutput(
		env,
		cldf.ChangesetOutput{},
		seqReport,
		err,
		state.EVMMCMSStateByChain(),
		c.MCMS,
		USDCTokenPoolConfigSequence.Description(),
	)
}
