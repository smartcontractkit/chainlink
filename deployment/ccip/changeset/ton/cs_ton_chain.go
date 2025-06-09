package ton

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	tonstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/ton"
	tonaddress "github.com/xssnick/tonutils-go/address"
)

// TODO: This file should be split into multiple files, one for each change set domain
type DeployPrerequisitesTONChainConfig struct {
	// MCMSDeployConfigPerChain   map[uint64]types.MCMSWithTimelockConfigV2
	// MCMSTimelockConfigPerChain map[uint64]proposalutils.TimelockConfig
	// ContractParamsPerChain     map[uint64]ChainContractParams
}

func (c DeployPrerequisitesTONChainConfig) Validate() error {
	// TODO: implement chain selector validation, contract parameters validation
	return nil
}

//////////////////////////////////////////////////////////////////////////////
// DeployPrerequisitesTONChain is a change set for deploying Prerequisite Ton chain packages and modules.
//////////////////////////////////////////////////////////////////////////////

var _ cldf.ChangeSetV2[DeployPrerequisitesTONChainConfig] = DeployPrerequisitesTONChain{}

type DeployPrerequisitesTONChain struct {
	TonChainSelectors []uint64 // List of Ton chain selectors to deploy
}

func (cs DeployPrerequisitesTONChain) VerifyPreconditions(env cldf.Environment, config DeployPrerequisitesTONChainConfig) error {
	// TODO: Implement precondition checks for contract deployment
	return nil
}

func (cs DeployPrerequisitesTONChain) Apply(env cldf.Environment, config DeployPrerequisitesTONChainConfig) (cldf.ChangesetOutput, error) {
	env.Logger.Infof("TON_E2E: Deploying prerequisites for TON chains: %v", cs.TonChainSelectors)
	tonChains, err := tonstate.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if len(tonChains) == 0 {
		return cldf.ChangesetOutput{}, nil // No Ton chains to deploy
	}

	for _, chainSelector := range cs.TonChainSelectors {
		tonChainState := tonChains[chainSelector]
		// TODO: these have to be in prerequisites, but we don't have prerequisites for TON yet????????

		// TODO(ton): once all contracts are deployed, we can remove the hardcoded addresses from the TonTestDeployPrerequisitesChangeSet
		// TODO(ton): Deploy TON MCMS, https://smartcontract-it.atlassian.net/browse/NONEVM-1939
		// TODO(ton): Deploy TON CCIP, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
		// TODO(ton): Deploy TON CCIP Offramp, Router, Onramp, Dummy Receiver and set the contract address
		// TODO(ton): Initialize Onramp, Offramp, FeeQuoter, RMNRemote

		// TODO: replace with actual TON addresses after contracts are supported, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
		address := tonaddress.MustParseAddr("EQDtFpEwcFAEcRe5mLVh2N6C0x-_hJEM7W61_JLnSF74p4q2")
		tonChainState.OffRamp = *address
		address = tonaddress.MustParseAddr("UQCfQRaJr2vxgZr5NHc0CTx6tAb0jverj9QQFirNfoCkGcUy")
		tonChainState.Router = *address
		address = tonaddress.MustParseAddr("EQADa3W6G0nSiTV4a6euRA42fU9QxSEnb-WeDpcrtWzA2jM8")
		tonChainState.LinkTokenAddress = *address
		address = tonaddress.MustParseAddr("UQDgFwiokL1ojVwXa3Ac7xCLfGB0Ti0foSw5NZ48Aj_vhs_6")
		tonChainState.CCIPAddress = *address
		address = tonaddress.MustParseAddr("UQCk4967vNM_V46Dn8I0x-gB_QE2KkdW1GQ7mWz1DtYGLEd8")
		tonChainState.ReceiverAddress = *address

		err = tonstate.SaveOnchainState(chainSelector, tonChainState, env)
	}

	return cldf.ChangesetOutput{}, nil
}

//////////////////////////////////////////////////////////////////////////////
// DeployTONChain is a change set for deploying Ton chain packages and modules.
//////////////////////////////////////////////////////////////////////////////

type DeployTONChainConfig struct{}

func (c DeployTONChainConfig) Validate() error {
	// TODO: implement chain selector validation, contract parameters validation
	return nil
}

var _ cldf.ChangeSetV2[DeployTONChainConfig] = DeployTONChain{}

// DeployTONChain deploys Ton chain packages and modules
type DeployTONChain struct {
	TonChainSelectors []uint64 // List of Ton chain selectors to deploy
}

func (cs DeployTONChain) VerifyPreconditions(env cldf.Environment, config DeployTONChainConfig) error {
	// TODO: Implement precondition checks for contract deployment
	return nil
}

func (cs DeployTONChain) Apply(env cldf.Environment, config DeployTONChainConfig) (cldf.ChangesetOutput, error) {
	env.Logger.Infof("TON_E2E: Deploying contracts for TON chains: %v", cs.TonChainSelectors)
	// TODO: Implement logic of deploying Ton chain packages and modules
	return cldf.ChangesetOutput{}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Add/UpdateLaneTONChain is a change set for adding or updating a lane on Ton chain.
//////////////////////////////////////////////////////////////////////////////

type AddLaneTONChainConfig struct {
	FromChainSelector uint64
	ToChainSelector   uint64
	FromFamily        string
	ToFamily          string
}

type AddLaneTONChain struct {
}

var _ cldf.ChangeSetV2[AddLaneTONChainConfig] = AddLaneTONChain{}

func (c AddLaneTONChain) VerifyPreconditions(env cldf.Environment, config AddLaneTONChainConfig) error {
	// TODO: Implement precondition checks for adding or updating a lane on Ton chain
	return nil
}

func (cs AddLaneTONChain) Apply(env cldf.Environment, config AddLaneTONChainConfig) (cldf.ChangesetOutput, error) {
	// TODO: Implement logic of adding or updating a lane on Ton chain
	return cldf.ChangesetOutput{}, nil
}
