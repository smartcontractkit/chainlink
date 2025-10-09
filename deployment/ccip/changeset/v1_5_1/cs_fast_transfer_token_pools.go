package v1_5_1

import (
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5_1"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	FastTransferUpdateLaneConfigChangeset = cldf.CreateChangeSet(fastTransferUpdateLaneConfigLogic, fastTransferUpdateLaneConfigPrecondition)
	FastTransferFillerAllowlistChangeset  = cldf.CreateChangeSet(fastTransferUpdateFillerAllowlistLogic, fastTransferUpdateFillerAllowlistPrecondition)
	FastTransferWithdrawPoolFeesChangeset = cldf.CreateChangeSet(fastTransferWithdrawPoolFeesLogic, fastTransferWithdrawPoolFeesPrecondition)
)

var (
	MaxFastTransferFillerFeeBps  = uint16(10000)
	DefaultSettlementOverheadGas = uint32(0)
	ChainFamilySelectorEVM       = [4]byte{0x28, 0x12, 0xd5, 0x2c}
)

type UpdateLaneConfig struct {
	FastTransferFillerFeeBps uint16
	FastTransferPoolFeeBps   uint16
	FillAmountMaxRequest     *big.Int
	FillerAllowlistEnabled   bool
	SkipAllowlistValidation  bool
	SettlementOverheadGas    *uint32
	CustomExtraArgs          []byte
	// Optional: Different contract type for the destination pool (if not set, uses the root config's ContractType)
	DestinationContractType *cldf.ContractType
	// Optional: Different contract version for the destination pool (if not set, uses the root config's ContractVersion)
	DestinationContractVersion *semver.Version
}

func (u UpdateLaneConfig) Validate(contract *bindings.FastTransferTokenPoolWrapper) error {
	if u.FastTransferFillerFeeBps > MaxFastTransferFillerFeeBps {
		return fmt.Errorf("fast transfer filler fee bps %d is greater than %d", u.FastTransferFillerFeeBps, MaxFastTransferFillerFeeBps)
	}
	if u.FastTransferPoolFeeBps > MaxFastTransferFillerFeeBps {
		return fmt.Errorf("fast transfer pool fee bps %d is greater than %d", u.FastTransferPoolFeeBps, MaxFastTransferFillerFeeBps)
	}
	if u.FillAmountMaxRequest == nil || u.FillAmountMaxRequest.Sign() <= 0 {
		return errors.New("fill amount max request must be a positive integer")
	}

	allowedFiller, err := contract.GetAllowedFillers(nil)
	if err != nil {
		return fmt.Errorf("failed to get allowed fillers: %w", err)
	}

	if !u.SkipAllowlistValidation && u.FillerAllowlistEnabled && len(allowedFiller) == 0 {
		return errors.New("filler allowlist is enabled but no fillers are allowed")
	}

	return nil
}

type FillerAllowlistConfig struct {
	AddFillers    []common.Address
	RemoveFillers []common.Address
}

func (f FillerAllowlistConfig) Validate(contract *bindings.FastTransferTokenPoolWrapper) error {
	if len(f.AddFillers) == 0 && len(f.RemoveFillers) == 0 {
		return errors.New("at least one filler must be added or removed")
	}
	if slices.Contains(f.AddFillers, (common.Address{})) {
		return errors.New("filler address cannot be empty")
	}
	if slices.Contains(f.RemoveFillers, (common.Address{})) {
		return errors.New("filler address cannot be empty")
	}

	allowedFillers, err := contract.GetAllowedFillers(nil)
	if err != nil {
		return fmt.Errorf("failed to get allowed fillers: %w", err)
	}
	for _, filler := range f.AddFillers {
		if slices.Contains(allowedFillers, filler) {
			return fmt.Errorf("filler %s is already in the allowlist", filler.Hex())
		}
	}
	for _, filler := range f.RemoveFillers {
		found := slices.Contains(allowedFillers, filler)
		if !found {
			return fmt.Errorf("filler %s is not in the allowlist", filler.Hex())
		}
	}

	return nil
}

type FastTransferWithdrawPoolFeesConfig struct {
	TokenSymbol     shared.TokenSymbol
	ContractType    cldf.ContractType
	ContractVersion semver.Version
	Withdrawals     map[uint64]common.Address // chainSelector -> recipient address
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func (c FastTransferWithdrawPoolFeesConfig) Validate(env cldf.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	if len(c.Withdrawals) == 0 {
		return errors.New("at least one withdrawal must be specified")
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, recipient := range c.Withdrawals {
		err := cldf.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain.String())
		}

		if err := validateFastTransferTokenPoolExists(chainState, c.TokenSymbol, c.ContractType, c.ContractVersion, chain.String()); err != nil {
			return err
		}

		if recipient == (common.Address{}) {
			return fmt.Errorf("recipient address cannot be empty for chain %d", chainSelector)
		}

		if c.MCMS != nil {
			if timelock := chainState.Timelock; timelock == nil {
				return fmt.Errorf("missing timelock on %s", chain.String())
			}
			if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
				return fmt.Errorf("missing proposerMcm on %s", chain.String())
			}
		}
	}
	return nil
}

type FastTransferUpdateLaneConfigConfig struct {
	TokenSymbol     shared.TokenSymbol
	ContractType    cldf.ContractType
	ContractVersion semver.Version
	Updates         map[uint64](map[uint64]UpdateLaneConfig)
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func (c FastTransferUpdateLaneConfigConfig) Validate(env cldf.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, poolUpdate := range c.Updates {
		err := cldf.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain.String())
		}

		if err := validateFastTransferTokenPoolExists(chainState, c.TokenSymbol, c.ContractType, c.ContractVersion, chain.String()); err != nil {
			return err
		}

		if c.MCMS != nil {
			if timelock := chainState.Timelock; timelock == nil {
				return fmt.Errorf("missing timelock on %s", chain.String())
			}
			if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
				return fmt.Errorf("missing proposerMcm on %s", chain.String())
			}
		}

		pool, err := bindings.GetFastTransferTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get fast transfer token pool contract for %s token on chain %d: %w", c.TokenSymbol, chainSelector, err)
		}

		for destinationChainSelector, update := range poolUpdate {
			err := update.Validate(pool)
			if err != nil {
				return fmt.Errorf("failed to validate update for chain selector %d: %w", chainSelector, err)
			}

			// Validate destination pool exists if custom type/version is specified
			if update.DestinationContractType != nil || update.DestinationContractVersion != nil {
				destContractType := c.ContractType
				destContractVersion := c.ContractVersion

				if update.DestinationContractType != nil {
					destContractType = *update.DestinationContractType
				}
				if update.DestinationContractVersion != nil {
					destContractVersion = *update.DestinationContractVersion
				}

				// Validate destination chain exists
				err := cldf.IsValidChainSelector(destinationChainSelector)
				if err != nil {
					return fmt.Errorf("failed to validate destination chain selector %d: %w", destinationChainSelector, err)
				}

				destChain, ok := env.BlockChains.EVMChains()[destinationChainSelector]
				if !ok {
					return fmt.Errorf("destination chain with selector %d does not exist in environment", destinationChainSelector)
				}

				destChainState, ok := state.Chains[destinationChainSelector]
				if !ok {
					return fmt.Errorf("destination %s does not exist in state", destChain.String())
				}

				// Validate destination pool exists with the specified type/version
				if err := validateFastTransferTokenPoolExists(destChainState, c.TokenSymbol, destContractType, destContractVersion, destChain.String()); err != nil {
					return fmt.Errorf("destination pool validation failed: %w", err)
				}
			}
		}
	}
	return nil
}

type FastTransferFillerAllowlistConfig struct {
	TokenSymbol     shared.TokenSymbol
	ContractType    cldf.ContractType
	ContractVersion semver.Version
	Updates         map[uint64]FillerAllowlistConfig
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func (c FastTransferFillerAllowlistConfig) Validate(env cldf.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, update := range c.Updates {
		err := cldf.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain.String())
		}

		if err := validateFastTransferTokenPoolExists(chainState, c.TokenSymbol, c.ContractType, c.ContractVersion, chain.String()); err != nil {
			return err
		}

		if c.MCMS != nil {
			if timelock := chainState.Timelock; timelock == nil {
				return fmt.Errorf("missing timelock on %s", chain.String())
			}
			if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
				return fmt.Errorf("missing proposerMcm on %s", chain.String())
			}
		}

		pool, err := bindings.GetFastTransferTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get fast transfer token pool contract for %s token on chain %d: %w", c.TokenSymbol, chainSelector, err)
		}

		err = update.Validate(pool)
		if err != nil {
			return fmt.Errorf("failed to validate filler allowlist update for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

func validateFastTransferTokenPoolExists(chainState evm.CCIPChainState, tokenSymbol shared.TokenSymbol, contractType cldf.ContractType, contractVersion semver.Version, chainString string) error {
	switch contractType {
	case shared.BurnMintFastTransferTokenPool:
		if _, ok := chainState.BurnMintFastTransferTokenPools[tokenSymbol]; !ok {
			return fmt.Errorf("token %s does not have a fast transfer token pool on %s", tokenSymbol, chainString)
		}
		if _, ok := chainState.BurnMintFastTransferTokenPools[tokenSymbol][contractVersion]; !ok {
			return fmt.Errorf("token %s does not have a fast transfer token pool with version %s on %s", tokenSymbol, contractVersion.String(), chainString)
		}
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		if _, ok := chainState.BurnMintWithExternalMinterFastTransferTokenPools[tokenSymbol]; !ok {
			return fmt.Errorf("token %s does not have a fast transfer token pool on %s", tokenSymbol, chainString)
		}
		if _, ok := chainState.BurnMintWithExternalMinterFastTransferTokenPools[tokenSymbol][contractVersion]; !ok {
			return fmt.Errorf("token %s does not have a fast transfer token pool with version %s on %s", tokenSymbol, contractVersion.String(), chainString)
		}
	case shared.HybridWithExternalMinterFastTransferTokenPool:
		if _, ok := chainState.HybridWithExternalMinterFastTransferTokenPools[tokenSymbol]; !ok {
			return fmt.Errorf("token %s does not have a fast transfer token pool on %s", tokenSymbol, chainString)
		}
		if _, ok := chainState.HybridWithExternalMinterFastTransferTokenPools[tokenSymbol][contractVersion]; !ok {
			return fmt.Errorf("token %s does not have a fast transfer token pool with version %s on %s", tokenSymbol, contractVersion.String(), chainString)
		}
	default:
		return fmt.Errorf("unsupported contract type %s for fast transfer token pools", contractType)
	}
	return nil
}

func fastTransferUpdateLaneConfigPrecondition(env cldf.Environment, c FastTransferUpdateLaneConfigConfig) error {
	return c.Validate(env)
}

func fastTransferUpdateLaneConfigLogic(env cldf.Environment, c FastTransferUpdateLaneConfigConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid FastTransferUpdateLaneConfigConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Build the sequence input for multi-chain updates
	updatesByChain := make(map[uint64]opsutil.EVMCallInput[ccipops.UpdateDestChainConfigInput])

	for sourceChainSelector, updates := range c.Updates {
		pool, err := bindings.GetFastTransferTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, sourceChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fast transfer token pool contract for %s token on chain %d: %w", c.TokenSymbol, sourceChainSelector, err)
		}

		laneConfigs := make([]bindings.DestChainConfigUpdateArgs, 0)
		for destinationChainSelector, update := range updates {
			// Determine destination pool contract type and version
			destContractType := c.ContractType
			destContractVersion := c.ContractVersion

			if update.DestinationContractType != nil {
				destContractType = *update.DestinationContractType
			}
			if update.DestinationContractVersion != nil {
				destContractVersion = *update.DestinationContractVersion
			}

			destinationPool, err := bindings.GetFastTransferTokenPoolContract(env, c.TokenSymbol, destContractType, destContractVersion, destinationChainSelector)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fast transfer token pool contract for %s token on chain %d: %w", c.TokenSymbol, destinationChainSelector, err)
			}
			settlementOverheadGas := DefaultSettlementOverheadGas
			if update.SettlementOverheadGas != nil {
				settlementOverheadGas = *update.SettlementOverheadGas
			}

			customExtraArgs := update.CustomExtraArgs
			if customExtraArgs == nil {
				customExtraArgs = []byte{}
			}

			laneConfigs = append(laneConfigs, bindings.DestChainConfigUpdateArgs{
				MaxFillAmountPerRequest:  update.FillAmountMaxRequest,
				FastTransferFillerFeeBps: update.FastTransferFillerFeeBps,
				FastTransferPoolFeeBps:   update.FastTransferPoolFeeBps,
				RemoteChainSelector:      destinationChainSelector,
				DestinationPool:          common.LeftPadBytes(destinationPool.Address().Bytes(), 32),
				FillerAllowlistEnabled:   update.FillerAllowlistEnabled,
				SettlementOverheadGas:    settlementOverheadGas,
				ChainFamilySelector:      ChainFamilySelectorEVM, // Only EVM chains supported
				CustomExtraArgs:          customExtraArgs,
			})
		}

		updatesByChain[sourceChainSelector] = opsutil.EVMCallInput[ccipops.UpdateDestChainConfigInput]{
			Address:       pool.Address(),
			ChainSelector: sourceChainSelector,
			CallInput: ccipops.UpdateDestChainConfigInput{
				Updates: laneConfigs,
			},
			NoSend: c.MCMS != nil, // Use NoSend for MCMS proposals
		}
	}

	// Execute the sequence
	seqInput := ccipseq.FastTransferTokenPoolUpdateDestChainConfigSequenceInput{
		ContractType:   c.ContractType,
		UpdatesByChain: updatesByChain,
	}

	seqReport, err := operations.ExecuteSequence(env.OperationsBundle, ccipseq.FastTransferTokenPoolUpdateDestChainConfigSequence, env.BlockChains.EVMChains(), seqInput)
	return opsutil.AddEVMCallSequenceToCSOutput(
		env,
		cldf.ChangesetOutput{},
		seqReport,
		err,
		state.EVMMCMSStateByChain(),
		c.MCMS,
		fmt.Sprintf("Update %s fast transfer token pool destination chain configurations", c.TokenSymbol),
	)
}

func fastTransferUpdateFillerAllowlistPrecondition(env cldf.Environment, c FastTransferFillerAllowlistConfig) error {
	return c.Validate(env)
}

func fastTransferUpdateFillerAllowlistLogic(env cldf.Environment, c FastTransferFillerAllowlistConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid FastTransferFillerAllowlistConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Build the sequence input for multi-chain updates
	updatesByChain := make(map[uint64]opsutil.EVMCallInput[ccipops.UpdateFillerAllowlistInput])

	for sourceChainSelector, update := range c.Updates {
		pool, err := bindings.GetFastTransferTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, sourceChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fast transfer token pool contract for %s token on chain %d: %w", c.TokenSymbol, sourceChainSelector, err)
		}

		updatesByChain[sourceChainSelector] = opsutil.EVMCallInput[ccipops.UpdateFillerAllowlistInput]{
			Address:       pool.Address(),
			ChainSelector: sourceChainSelector,
			CallInput: ccipops.UpdateFillerAllowlistInput{
				AddFillers:    update.AddFillers,
				RemoveFillers: update.RemoveFillers,
			},
			NoSend: c.MCMS != nil, // Use NoSend for MCMS proposals
		}
	}

	// Execute the sequence
	seqInput := ccipseq.FastTransferTokenPoolUpdateFillerAllowlistSequenceInput{
		ContractType:   c.ContractType,
		UpdatesByChain: updatesByChain,
	}

	seqReport, err := operations.ExecuteSequence(env.OperationsBundle, ccipseq.FastTransferTokenPoolUpdateFillerAllowlistSequence, env.BlockChains.EVMChains(), seqInput)
	return opsutil.AddEVMCallSequenceToCSOutput(
		env,
		cldf.ChangesetOutput{},
		seqReport,
		err,
		state.EVMMCMSStateByChain(),
		c.MCMS,
		fmt.Sprintf("Update %s fast transfer token pool filler allowlists", c.TokenSymbol),
	)
}

func fastTransferWithdrawPoolFeesPrecondition(env cldf.Environment, c FastTransferWithdrawPoolFeesConfig) error {
	return c.Validate(env)
}

func fastTransferWithdrawPoolFeesLogic(env cldf.Environment, c FastTransferWithdrawPoolFeesConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid FastTransferWithdrawPoolFeesConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Build the sequence input for multi-chain withdrawals
	withdrawalsByChain := make(map[uint64]opsutil.EVMCallInput[ccipops.WithdrawPoolFeesInput])

	for chainSelector, recipient := range c.Withdrawals {
		pool, err := bindings.GetFastTransferTokenPoolContract(env, c.TokenSymbol, c.ContractType, c.ContractVersion, chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fast transfer token pool contract for %s token on chain %d: %w", c.TokenSymbol, chainSelector, err)
		}

		withdrawalsByChain[chainSelector] = opsutil.EVMCallInput[ccipops.WithdrawPoolFeesInput]{
			Address:       pool.Address(),
			ChainSelector: chainSelector,
			CallInput: ccipops.WithdrawPoolFeesInput{
				Recipient: recipient,
			},
			NoSend: c.MCMS != nil, // Use NoSend for MCMS proposals
		}
	}

	// Execute the sequence
	seqInput := ccipseq.FastTransferTokenPoolWithdrawPoolFeesSequenceInput{
		ContractType:       c.ContractType,
		WithdrawalsByChain: withdrawalsByChain,
	}

	seqReport, err := operations.ExecuteSequence(env.OperationsBundle, ccipseq.FastTransferTokenPoolWithdrawPoolFeesSequence, env.BlockChains.EVMChains(), seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute fast transfer token pool withdraw pool fees sequence: %w", err)
	}

	return opsutil.AddEVMCallSequenceToCSOutput(
		env,
		cldf.ChangesetOutput{},
		seqReport,
		nil, // no error since we already handled it above
		state.EVMMCMSStateByChain(),
		c.MCMS,
		fmt.Sprintf("Withdraw pool fees from %s fast transfer token pools", c.TokenSymbol),
	)
}
