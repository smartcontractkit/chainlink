package ton

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	tonConfig "github.com/smartcontractkit/chainlink-ton/ops/ccip/config"
	tonOps "github.com/smartcontractkit/chainlink-ton/ops/ccip/operation"
	"github.com/smartcontractkit/chainlink-ton/ops/ccip/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	tonstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/ton"
	tonaddress "github.com/xssnick/tonutils-go/address"
)

type DeployCCIPContractsCfg struct {
	TonChainSelector    uint64
	ChainContractParams tonConfig.ChainContractParams
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
	tonState := chains[selector]

	ab := cldf.NewMemoryAddressBook()
	tonChains := env.BlockChains.TonChains()
	tonChain := tonChains[config.TonChainSelector]

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state after deployment: %w", err)
	}

	deps := tonOps.TonDeps{
		AB:               ab,
		TonChain:         tonChain,
		CCIPOnChainState: state,
	}

	deployInput := sequence.DeployCCIPSeqInput{
		CCIPConfig: config.ChainContractParams,
	}

	ccipDeployOutput, err := operations.ExecuteSequence(env.OperationsBundle, sequence.DeployCCIPSequence, deps, deployInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute MCMS deploy sequence: %w", err)
	}

	tonState.Router = *ccipDeployOutput.Output.RouterAddress
	tonState.OnRamp = *ccipDeployOutput.Output.OnRampAddress
	tonState.FeeQuoter = *ccipDeployOutput.Output.FeeQuoterAddress

	// TODO replace the dummy addresses with actual deployed contracts
	address := tonaddress.MustParseAddr("EQDtFpEwcFAEcRe5mLVh2N6C0x-_hJEM7W61_JLnSF74p4q2")
	tonState.OffRamp = *address
	address = tonaddress.MustParseAddr("EQADa3W6G0nSiTV4a6euRA42fU9QxSEnb-WeDpcrtWzA2jM8")
	tonState.LinkTokenAddress = *address
	address = tonaddress.MustParseAddr("UQCk4967vNM_V46Dn8I0x-gB_QE2KkdW1GQ7mWz1DtYGLEd8")
	tonState.ReceiverAddress = *address

	// update chain state
	err = tonstate.SaveOnchainState(selector, tonState, env)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{}, nil
}
