package v1_6

import (
	"fmt"
	"reflect"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type SetRMNRemoteConfig struct {
	ChainSelector   uint64                        `json:"chainSelector"`
	RMNRemoteConfig RMNRemoteConfig               `json:"rmnRemoteConfigs"`
	MCMSConfig      *proposalutils.TimelockConfig `json:"mcmsConfig,omitempty"`
}

type RMNRemoteConfig struct {
	Signers []rmn_remote.RMNRemoteSigner `json:"signers"`
	F       uint64                       `json:"f"`
}

var (
	DeployRMNRemoteOp = operations.NewOperation(
		"DeployRMNRemote",
		semver.MustParse("1.0.0"),
		"Deploys RMNRemote 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.OpDependencies, input uint64) (common.Address, error) {
			state := deps.CurrentState
			e := deps.Env
			ab := deps.AddressBook
			chain := e.Chains[input]
			chainState, chainExists := state.Chains[chain.Selector]
			if !chainExists {
				return common.Address{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", chain.String())
			}
			var rmnLegacyAddr common.Address
			if chainState.MockRMN != nil {
				rmnLegacyAddr = chainState.MockRMN.Address()
			}
			// If RMN is deployed, set rmnLegacyAddr to the RMN address
			if chainState.RMN != nil {
				rmnLegacyAddr = chainState.RMN.Address()
			}
			if rmnLegacyAddr == (common.Address{}) {
				e.Logger.Warnf("No legacy RMN contract found for chain %s, will not setRMN in RMNRemote", chain.String())
			}
			if chainState.RMNRemote == nil {
				contract, err := cldf.DeployContract(e.Logger, chain, ab,
					func(chain cldf.Chain) cldf.ContractDeploy[*rmn_remote.RMNRemote] {
						rmnRemoteAddr, tx, rmnRemote, err2 := rmn_remote.DeployRMNRemote(
							chain.DeployerKey,
							chain.Client,
							chain.Selector,
							rmnLegacyAddr,
						)
						return cldf.ContractDeploy[*rmn_remote.RMNRemote]{
							Address: rmnRemoteAddr, Contract: rmnRemote, Tx: tx, Tv: cldf.NewTypeAndVersion(shared.RMNRemote, deployment.Version1_6_0), Err: err2,
						}
					})
				if err != nil {
					e.Logger.Errorw("Failed to deploy RMNRemote", "chain", chain.String(), "err", err)
					return common.Address{}, err
				}
				return contract.Address, nil
			}
			e.Logger.Infow("rmn remote already deployed, no-op", "chain", chain.String(), "addr", chainState.RMNRemote.Address)
			return common.Address{}, nil
		})

	SetRMNRemoteConfigOp = operations.NewOperation(
		"SetRMNRemoteConfigOp",
		semver.MustParse("1.0.0"),
		"Setting RMNRemoteConfig based on ActiveDigest from RMNHome",
		func(b operations.Bundle, deps opsutil.OpDependencies, input SetRMNRemoteConfig) (opsutil.OpOutput, error) {
			e := deps.Env
			state := deps.CurrentState
			lggr := e.Logger
			lggr.Infow("Setting RMNRemoteConfig based on ActiveDigest from RMNHome", "chain", input.ChainSelector)
			chain, ok := e.Chains[input.ChainSelector]
			if !ok {
				return opsutil.OpOutput{}, fmt.Errorf("chain %d not found in environment", input.ChainSelector)
			}
			homeChainSel, err := state.HomeChainSelector()
			if err != nil {
				return opsutil.OpOutput{}, err
			}
			homeChain, ok := e.Chains[homeChainSel]
			if !ok {
				return opsutil.OpOutput{}, fmt.Errorf("chain %d not found in ", homeChainSel)
			}
			rmnHome := state.Chains[homeChainSel].RMNHome
			if rmnHome == nil {
				return opsutil.OpOutput{}, fmt.Errorf("RMNHome not found for chain %s", homeChain.String())
			}

			activeConfig, err := rmnHome.GetActiveDigest(nil)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to get RMNHome active digest for chain %s: %w", homeChain.String(), err)
			}
			rmnRemote := state.Chains[input.ChainSelector].RMNRemote
			currentVersionConfig, err := rmnRemote.GetVersionedConfig(nil)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to get RMNRemote config for chain %s: %w", chain, err)
			}
			newConfig := rmn_remote.RMNRemoteConfig{
				RmnHomeContractConfigDigest: activeConfig,
				Signers:                     input.RMNRemoteConfig.Signers,
				FSign:                       input.RMNRemoteConfig.F,
			}

			if reflect.DeepEqual(currentVersionConfig.Config, newConfig) {
				lggr.Infow("RMNRemote config already up to date, it's a no-op", "chain", chain.String())
				return opsutil.OpOutput{}, nil
			}
			deployerGroup := deployergroup.NewDeployerGroup(e, state, input.MCMSConfig).
				WithDeploymentContext("set RMNRemote config for chain " + chain.String())
			opts, err := deployerGroup.GetDeployer(input.ChainSelector)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
			}
			_, err = rmnRemote.SetConfig(opts, newConfig)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("build call data to set RMNRemote config for chain %s: %w", chain.String(), err)
			}
			csOutput, err := deployerGroup.Enact()
			if err != nil {
				return opsutil.OpOutput{}, err
			}
			return opsutil.OpOutput{
				Proposals:                  csOutput.MCMSTimelockProposals,
				DescribedTimelockProposals: csOutput.DescribedTimelockProposals,
			}, nil
		})
)
