package v1_6_2

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cmtp "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/cctp_message_transmitter_proxy"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var (
	ConfigureCCTPMessageTransmitterProxy = cldf.CreateChangeSet(configureCCTPMessageTransmitterProxyContractLogic, configureCCTPMessageTransmitterProxyContractPrecondition)

	CCTPMessageTransmitterProxyConfigOp = opsutil.NewEVMCallOperation(
		"CCTPMessageTransmitterProxyConfigOp",
		semver.MustParse("1.6.2"),
		"Setting CCTP message transmitter proxy config",
		cmtp.CCTPMessageTransmitterProxyABI,
		shared.CCTPMessageTransmitterProxy,
		cmtp.NewCCTPMessageTransmitterProxy,
		func(proxy *cmtp.CCTPMessageTransmitterProxy, opts *bind.TransactOpts, input []cmtp.CCTPMessageTransmitterProxyAllowedCallerConfigArgs) (*types.Transaction, error) {
			return proxy.ConfigureAllowedCallers(opts, input)
		})

	CCTPMessageTransmitterProxyConfigSequence = operations.NewSequence(
		"CCTPMessageTransmitterProxyConfigSequence",
		semver.MustParse("1.6.2"),
		"Setting CCTP message transmitter proxy config across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, inputs map[uint64]opsutil.EVMCallInput[[]cmtp.CCTPMessageTransmitterProxyAllowedCallerConfigArgs]) (map[uint64][]opsutil.EVMCallOutput, error) {
			out := make(map[uint64][]opsutil.EVMCallOutput, len(inputs))

			for chainSelector, input := range inputs {
				if _, ok := chains[chainSelector]; !ok {
					return nil, fmt.Errorf("chain with selector %d not defined in dependencies", chainSelector)
				}

				report, err := operations.ExecuteOperation(b, CCTPMessageTransmitterProxyConfigOp, chains[chainSelector], input)
				if err != nil {
					return map[uint64][]opsutil.EVMCallOutput{}, fmt.Errorf("failed to set CCTP message transmitt proxy config for chain %d: %w", chainSelector, err)
				}
				out[chainSelector] = []opsutil.EVMCallOutput{report.Output}
			}

			return out, nil
		})
)

type AllowedCallerUpdate struct {
	AllowedCaller common.Address
	Enabled       bool
}

// ConfigureCCTPMessageTransmitterProxyInput defines all information required of the user to configure a new CCTP message transmitter proxy contract.
type ConfigureCCTPMessageTransmitterProxyInput struct {
	// AllowedCaller is the address of the USDC token messenger contract.
	AllowedCallerUpdates []AllowedCallerUpdate
}

func (i ConfigureCCTPMessageTransmitterProxyInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {
	if _, ok := state.CCTPMessageTransmitterProxies[deployment.Version1_6_2]; !ok {
		return fmt.Errorf("no CCTP proxy with version %s found on %s", deployment.Version1_6_2, chain.Name())
	}
	for _, allowedCalleUpdate := range i.AllowedCallerUpdates {
		if allowedCalleUpdate.AllowedCaller == utils.ZeroAddress {
			return fmt.Errorf("allowed caller must be defined for chain %s", chain.Name())
		}

		// Skip allowed caller validation against USDC token pools when removing callers
		if !allowedCalleUpdate.Enabled {
			continue
		}

		matchedPool := false
		for _, usdcTokenPool := range state.USDCTokenPoolsV1_6 {
			if usdcTokenPool.Address().Cmp(allowedCalleUpdate.AllowedCaller) == 0 {
				matchedPool = true
			}
		}
		if !matchedPool {
			return fmt.Errorf("allowed caller does not match any existing 1.6 USDC token pools (%s)", allowedCalleUpdate.AllowedCaller.Hex())
		}
	}

	return nil
}

// ConfigureCCTPMessageTransmitterProxyContractConfig defines the configuration for configuring the CCTP message transmitter proxy contracts.
type ConfigureCCTPMessageTransmitterProxyContractConfig struct {
	CCTPProxies map[uint64]ConfigureCCTPMessageTransmitterProxyInput

	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func configureCCTPMessageTransmitterProxyContractPrecondition(env cldf.Environment, c ConfigureCCTPMessageTransmitterProxyContractConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, proxyConfig := range c.CCTPProxies {
		chain, chainState, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
		}
		err = proxyConfig.Validate(env.GetContext(), chain, chainState)
		if err != nil {
			return fmt.Errorf("failed to validate USDC token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

// configureCCTPMessageTransmitterProxyContractLogic sets the configurations in the new CCTP message transmitter proxy across multiple chains.
func configureCCTPMessageTransmitterProxyContractLogic(env cldf.Environment, c ConfigureCCTPMessageTransmitterProxyContractConfig) (cldf.ChangesetOutput, error) {
	if err := configureCCTPMessageTransmitterProxyContractPrecondition(env, c); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid ConfigureCCTPMessageTransmitterProxyContractConfig: %w", err)
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Convert CLD/migrations inputs to onchain inputs.
	input := make(map[uint64]opsutil.EVMCallInput[[]cmtp.CCTPMessageTransmitterProxyAllowedCallerConfigArgs], len(c.CCTPProxies))
	for chainSelector, proxyConfig := range c.CCTPProxies {
		_, chainState, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get EVM chain state for chain selector %d: %w",
				chainSelector, err)
		}

		allowedCallerInputs := make([]cmtp.CCTPMessageTransmitterProxyAllowedCallerConfigArgs, len(proxyConfig.AllowedCallerUpdates))
		for _, allowedCallerUpdate := range proxyConfig.AllowedCallerUpdates {
			allowedCallerInputs = append(allowedCallerInputs, cmtp.CCTPMessageTransmitterProxyAllowedCallerConfigArgs{
				Allowed: true,
				Caller:  allowedCallerUpdate.AllowedCaller,
			})
		}

		input[chainSelector] = opsutil.EVMCallInput[[]cmtp.CCTPMessageTransmitterProxyAllowedCallerConfigArgs]{
			ChainSelector: chainSelector,
			NoSend:        c.MCMS != nil,
			Address:       chainState.CCTPMessageTransmitterProxies[deployment.Version1_6_2].Address(),
			CallInput:     allowedCallerInputs,
		}
	}

	// Configure sequence.
	seqReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		CCTPMessageTransmitterProxyConfigSequence,
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
		CCTPMessageTransmitterProxyConfigSequence.Description(),
	)
}
