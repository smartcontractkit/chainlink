package changeset

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
)

type DeployChainInput struct {
	// Targeted chain
	ChainID uint64

	// Existing contract addresses
	Weth9ContractAddress      common.Address
	RMNHomeAddress            common.Address
	LINKAddress               common.Address
	TokenAdminRegistryAddress common.Address
	RMNProxyAddress           common.Address
	MockRMNAddress            common.Address
	RMNAddress                common.Address
	TimelockAddress           common.Address

	// Params needed for new deployments and configs
	FeeQuoterParams FeeQuoterParams
	OffRampParams   OffRampParams
}

/*
Replication of deployChainContractsEVM changeset in cs_deploy_chain.go

Process:

- Verify some contracts exist in the state: LINK...
- rmn_remote.DeployRMNRemote()
- activeDigest, err := rmnHome.GetActiveDigest(&bind.CallOpts{})
- rmnRemoteContract.SetConfig()
- router.DeployRouter()
- nonce_manager.DeployNonceManager()
- fee_quoter.DeployFeeQuoter()
- onramp.DeployOnRamp()
- offramp.DeployOffRamp()
- feeQuoterContract.ApplyAuthorizedCallerUpdates()
*/
var DeployChain = func(env deployment.Environment, input DeployChainInput) (deployment.ChangesetOutput, error) {
	// Prepare operation context
	auth := env.Chains[input.ChainID].DeployerKey
	client := env.Chains[input.ChainID].Client
	ethCtx := deployment_ethereum.EthereumDeps{
		Auth:    auth,
		Client:  client,
		Confirm: env.Chains[input.ChainID].ConfirmByHash,
	}

	// Deploy RMNRemote
	deployRMNRemoteInput := DeployRMNRemoteInput{
		LocalChainSelector: input.ChainID,
		LegacyRMN:          input.RMNAddress,
	}

	deployRMNRemoteRep, err := deployment.ExecuteOp(env.OpEnv, DeployRMNRemoteOp, ethCtx, deployRMNRemoteInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Get active digest
	rmnHomeContract, err := rmn_home.NewRMNHome(input.RMNHomeAddress, client)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	activeDigest, err := rmnHomeContract.GetActiveDigest(nil)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Set RMNRemote config
	setConfigInput := SetConfigRMNRemoteInput{
		contractAddress: deployRMNRemoteRep.Output.Address,
		RMNRemoteConfig: rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeDigest,
			Signers: []rmn_remote.RMNRemoteSigner{
				{NodeIndex: 0, OnchainPublicKey: common.Address{1}},
			},
			FSign: 0, // TODO: update when we have signers
		},
	}

	_, err = deployment.ExecuteOp(env.OpEnv, RMNRemoteSetConfigOp, ethCtx, setConfigInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Deploy Router
	deployRouterInput := DeployRouterInput{
		wrappedNative: input.Weth9ContractAddress,
		armProxy:      input.RMNProxyAddress,
	}
	_, err = deployment.ExecuteOp(env.OpEnv, DeployRouterOp, ethCtx, deployRouterInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Deploy NonceManager
	deployNonceManagerInput := DeployNonceManagerInput{
		AuthorizedCallers: []common.Address{},
	}
	deployNonceManagerRep, err := deployment.ExecuteOp(env.OpEnv, DeployNonceManagerOp, ethCtx, deployNonceManagerInput)

	// Deploy FeeQuoter
	deployFeeQuoterInput := DeployFeeQuoterInput{
		fee_quoter.FeeQuoterStaticConfig{
			MaxFeeJuelsPerMsg:            input.FeeQuoterParams.MaxFeeJuelsPerMsg,
			LinkToken:                    input.LINKAddress,
			TokenPriceStalenessThreshold: input.FeeQuoterParams.TokenPriceStalenessThreshold,
		},
		[]common.Address{input.TimelockAddress},                         // timelock should be able to update, ramps added after
		[]common.Address{input.Weth9ContractAddress, input.LINKAddress}, // fee tokens
		input.FeeQuoterParams.TokenPriceFeedUpdates,
		input.FeeQuoterParams.TokenTransferFeeConfigArgs,
		append([]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
			{
				PremiumMultiplierWeiPerEth: input.FeeQuoterParams.LinkPremiumMultiplierWeiPerEth,
				Token:                      input.LINKAddress,
			},
			{
				PremiumMultiplierWeiPerEth: input.FeeQuoterParams.WethPremiumMultiplierWeiPerEth,
				Token:                      input.Weth9ContractAddress,
			},
		}, input.FeeQuoterParams.MorePremiumMultiplierWeiPerEth...),
		input.FeeQuoterParams.DestChainConfigArgs,
	}

	deployFeeQuoterRep, err := deployment.ExecuteOp(env.OpEnv, DeployFeeQuoterOp, ethCtx, deployFeeQuoterInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Deploy Onramp

	deployOnrampInput := DeployOnRampInput{
		onramp.OnRampStaticConfig{
			ChainSelector:      input.ChainID,
			RmnRemote:          input.RMNProxyAddress,
			NonceManager:       deployNonceManagerRep.Output.Address,
			TokenAdminRegistry: input.TokenAdminRegistryAddress,
		},
		onramp.OnRampDynamicConfig{
			FeeQuoter:     deployFeeQuoterRep.Output.Address,
			FeeAggregator: auth.From, // TODO real fee aggregator, using deployer key for now
		},
		[]onramp.OnRampDestChainConfigArgs{},
	}
	deployOnrampRep, err := deployment.ExecuteOp(env.OpEnv, DeployOnRampOp, ethCtx, deployOnrampInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Deploy Offramp
	deployOfframpInput := DeployOffRampInput{
		offramp.OffRampStaticConfig{
			ChainSelector:        input.ChainID,
			GasForCallExactCheck: input.OffRampParams.GasForCallExactCheck,
			RmnRemote:            input.RMNProxyAddress,
			NonceManager:         deployNonceManagerRep.Output.Address,
			TokenAdminRegistry:   input.TokenAdminRegistryAddress,
		},
		offramp.OffRampDynamicConfig{
			FeeQuoter:                               deployFeeQuoterRep.Output.Address,
			PermissionLessExecutionThresholdSeconds: input.OffRampParams.PermissionLessExecutionThresholdSeconds,
			IsRMNVerificationDisabled:               input.OffRampParams.IsRMNVerificationDisabled,
		},
		[]offramp.OffRampSourceChainConfigArgs{},
	}
	deployOfframpRep, err := deployment.ExecuteOp(env.OpEnv, DeployOffRampOp, ethCtx, deployOfframpInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Apply authorized caller updates on FeeQuoter and NonceManager
	applyAuthorizedCallerUpdatesInput := ApplyAuthorizedCallerUpdatesInput{
		contractAddress: deployFeeQuoterRep.Output.Address,
		AuthorizedCallersAuthorizedCallerArgs: fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
			AddedCallers: []common.Address{deployOfframpRep.Output.Address, auth.From},
		},
	}

	_, err = deployment.ExecuteOp(env.OpEnv, FeeQuoterApplyCallerUpdatesOp, ethCtx, applyAuthorizedCallerUpdatesInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	nmApplyAuthorizedCallerUpdatesInput := NMApplyAuthorizedCallerUpdatesInput{
		contractAddress: deployNonceManagerRep.Output.Address,
		AuthorizedCallersAuthorizedCallerArgs: nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
			AddedCallers: []common.Address{deployOfframpRep.Output.Address, deployOnrampRep.Output.Address},
		},
	}

	_, err = deployment.ExecuteOp(env.OpEnv, NonceManagerApplyCallerUpdatesOp, ethCtx, nmApplyAuthorizedCallerUpdatesInput)

	// We can compile every address deployed in the address book
	ab := deployment.NewMemoryAddressBook()

	// TV should come from the operation output
	ab.Save(input.ChainID, deployRMNRemoteRep.Output.Address.String(), deployment.TypeAndVersion{
		Type:    "RMNRemote",
		Version: deployment.Version1_6_0_dev,
	})

	return deployment.ChangesetOutput{
		AddressBook: ab,
		Reports:     env.OpEnv.Reporter.GetReports(),
	}, nil

}
