package ton

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	tonaddress "github.com/xssnick/tonutils-go/address"

	tonstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/ton"
)

type DeployCCIPContractsCfg struct {
	TonChainSelector uint64
}

func (c DeployCCIPContractsCfg) Validate() error {
	// TODO: implement chain selector validation, contract parameters validation
	return nil
}

var _ cldf.ChangeSetV2[DeployCCIPContractsCfg] = DeployCCIPContracts{}

// DeployCCIPContracts deploys Ton chain packages and modules
type DeployCCIPContracts struct{}

func (cs DeployCCIPContracts) VerifyPreconditions(_ cldf.Environment, _ DeployCCIPContractsCfg) error {
	// TODO: Implement precondition checks for contract deployment
	return nil
}

func (cs DeployCCIPContracts) Apply(env cldf.Environment, config DeployCCIPContractsCfg) (cldf.ChangesetOutput, error) {
	// TODO: Implement logic of deploying Ton chain packages and modules
	// - once all contracts are deployed, we can remove the hardcoded addresses from the TonTestDeployPrerequisitesChangeSet
	// - Deploy TON MCMS, https://smartcontract-it.atlassian.net/browse/NONEVM-1939
	// - Deploy and initialize TON CCIP Offramp, Router, Onramp, Dummy Receiver and set the contract address https://smartcontract-it.atlassian.net/browse/NONEVM-1938
	// - Replace with actual TON addresses after contracts are supported, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
	env.Logger.Infof("TON_E2E: Deploying contracts for TON chains: %v", config.TonChainSelector)
	selector := config.TonChainSelector

	chains, err := tonstate.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	state := chains[selector]

	address := tonaddress.MustParseAddr("EQDtFpEwcFAEcRe5mLVh2N6C0x-_hJEM7W61_JLnSF74p4q2")
	state.OffRamp = *address
	address = tonaddress.MustParseAddr("UQCfQRaJr2vxgZr5NHc0CTx6tAb0jverj9QQFirNfoCkGcUy")
	state.Router = *address
	address = tonaddress.MustParseAddr("EQADa3W6G0nSiTV4a6euRA42fU9QxSEnb-WeDpcrtWzA2jM8")
	state.LinkTokenAddress = *address
	address = tonaddress.MustParseAddr("UQDgFwiokL1ojVwXa3Ac7xCLfGB0Ti0foSw5NZ48Aj_vhs_6")
	state.CCIPAddress = *address
	address = tonaddress.MustParseAddr("UQCk4967vNM_V46Dn8I0x-gB_QE2KkdW1GQ7mWz1DtYGLEd8")
	state.ReceiverAddress = *address
	// update chain state
	err = tonstate.SaveOnchainState(selector, state, env)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{}, nil
}
