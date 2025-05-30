package v1_6

import (
	"fmt"
	"reflect"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
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

type DeployRMNRemoteInput struct {
	RMNLegacyAddr common.Address `json:"rmnLegacyAddr"`
	ChainSelector uint64         `json:"chainSelector"`
}

var (
	DeployRMNRemoteOp = operations.NewOperation(
		"DeployRMNRemote",
		semver.MustParse("1.0.0"),
		"Deploys RMNRemote 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.DeployContractDependencies, input DeployRMNRemoteInput) (common.Address, error) {
			ab := deps.AddressBook
			chain := deps.Chain
			contract, err := cldf.DeployContract(b.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*rmn_remote.RMNRemote] {
					rmnRemoteAddr, tx, rmnRemote, err2 := rmn_remote.DeployRMNRemote(
						chain.DeployerKey,
						chain.Client,
						chain.Selector,
						input.RMNLegacyAddr,
					)
					return cldf.ContractDeploy[*rmn_remote.RMNRemote]{
						Address: rmnRemoteAddr, Contract: rmnRemote, Tx: tx, Tv: cldf.NewTypeAndVersion(shared.RMNRemote, deployment.Version1_6_0), Err: err2,
					}
				})
			if err != nil {
				b.Logger.Errorw("Failed to deploy RMNRemote", "chain", chain.String(), "err", err)
				return common.Address{}, err
			}
			return contract.Address, nil
		})

	SetRMNRemoteConfigOp = operations.NewOperation(
		"SetRMNRemoteConfigOp",
		semver.MustParse("1.0.0"),
		"Setting RMNRemoteConfig based on ActiveDigest from RMNHome",
		func(b operations.Bundle, deps opsutil.ConfigureDependencies, input SetRMNRemoteConfig) (opsutil.OpOutput, error) {
			state := deps.CurrentState
			b.Logger.Infow("Setting RMNRemoteConfig based on ActiveDigest from RMNHome", "chain", input.ChainSelector)
			e := deps.Env
			chain := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
			homeChainSel, err := state.HomeChainSelector()
			if err != nil {
				return opsutil.OpOutput{}, err
			}

			rmnHome := state.MustGetEVMChainState(homeChainSel).RMNHome
			if rmnHome == nil {
				return opsutil.OpOutput{}, fmt.Errorf("RMNHome not found for chain %d", homeChainSel)
			}

			activeConfig, err := rmnHome.GetActiveDigest(nil)
			if err != nil {
				return opsutil.OpOutput{}, fmt.Errorf("failed to get RMNHome active digest for chain %d: %w", homeChainSel, err)
			}
			rmnRemote := state.MustGetEVMChainState(input.ChainSelector).RMNRemote
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
				b.Logger.Infow("RMNRemote config already up to date, it's a no-op", "chain", chain.String())
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
