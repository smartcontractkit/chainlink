package v1_6_2

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cmtp "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/cctp_message_transmitter_proxy"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var DeployCCTPMessageTransmitterProxyNew = cldf.CreateChangeSet(deployCCTPMessageTransmitterProxyContractLogic, deployCCTPMessageTransmitterProxyContractPrecondition)

// DeployCCTPMessageTransmitterProxyInput defines all information required of the user to deploy a new CCTP message transmitter proxy contract.
type DeployCCTPMessageTransmitterProxyInput struct {
	// TokenMessenger is the address of the USDC token messenger contract.
	TokenMessenger common.Address
}

func (i DeployCCTPMessageTransmitterProxyInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {
	// The message transmitter consts are defined in the chainlink-deployments project so we can't validate them here.
	if i.TokenMessenger == utils.ZeroAddress {
		return fmt.Errorf("token messenger must be defined for chain %s", chain.Name())
	}

	return nil
}

// DeployCCTPMessageTransmitterProxyContractConfig defines the configuration for deploying CCTP message transmitter proxy contracts.
type DeployCCTPMessageTransmitterProxyContractConfig struct {
	USDCProxies map[uint64]DeployCCTPMessageTransmitterProxyInput
}

func deployCCTPMessageTransmitterProxyContractPrecondition(env cldf.Environment, c DeployCCTPMessageTransmitterProxyContractConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, proxyConfig := range c.USDCProxies {
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

// DeployCCTPMessageTransmitterProxyContractChangeset deploys new CCTP message transmitter proxies across multiple chains.
func deployCCTPMessageTransmitterProxyContractLogic(env cldf.Environment, c DeployCCTPMessageTransmitterProxyContractConfig) (cldf.ChangesetOutput, error) {
	if err := deployCCTPMessageTransmitterProxyContractPrecondition(env, c); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployCCTPMessageTransmitterProxyContractConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, proxyConfig := range c.USDCProxies {
		chain, _, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
		}
		_, err = cldf.DeployContract(env.Logger, chain, newAddresses,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*cmtp.CCTPMessageTransmitterProxy] {
				proxyAddress, tx, proxy, err := cmtp.DeployCCTPMessageTransmitterProxy(
					chain.DeployerKey,          // auth
					chain.Client,               // backend
					proxyConfig.TokenMessenger, // tokenMessenger
				)
				return cldf.ContractDeploy[*cmtp.CCTPMessageTransmitterProxy]{
					Address:  proxyAddress,
					Contract: proxy,
					Tv:       cldf.NewTypeAndVersion(shared.CCTPMessageTransmitterProxy, deployment.Version1_6_2),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCTPMessageTransmitterProxy on %s: %w", chain, err)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: newAddresses, // TODO: this is deprecated, how do I use the DataStore instead?
	}, nil
}
