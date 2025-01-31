package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// MCMSWithTimelockProgramsSolana holds the solana publick keys
// for a MCMSWithTimelock deployed program IDs.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockProgramsSolana struct {
	CancellerMcm solana.PublicKey
	BypasserMcm  solana.PublicKey
	ProposerMcm  solana.PublicKey
	Timelock     solana.PublicKey
	CallProxy    solana.PublicKey
}

func (state MCMSWithTimelockProgramsSolana) GetStateFromType(programType deployment.ContractType) solana.PublicKey {
	switch programType {
	case types.RBACTimelock:
		return state.Timelock
	case types.CallProxy:
		return state.CallProxy
	case types.ProposerManyChainMultisig:
		return state.ProposerMcm
	case types.BypasserManyChainMultisig:
		return state.BypasserMcm
	case types.CancellerManyChainMultisig:
		return state.CancellerMcm
	}
	return solana.PublicKey{}

}

// Validate checks that all fields are non-nil, ensuring it's ready
// for use generating views or interactions.
func (state MCMSWithTimelockProgramsSolana) Validate() error {
	if state.Timelock.IsZero() {
		return errors.New("timelock not found")
	}
	if state.CancellerMcm.IsZero() {
		return errors.New("canceller not found")
	}
	if state.ProposerMcm.IsZero() {
		return errors.New("proposer not found")
	}
	if state.BypasserMcm.IsZero() {
		return errors.New("bypasser not found")
	}
	if state.CallProxy.IsZero() {
		return errors.New("call proxy not found")
	}
	return nil
}

// MaybeLoadMCMSSolanaWithTimelockContracts looks for the program IDs / seeds corresponding to
// contracts deployed with DeployMCMSWithTimelock and loads them into a
// MCMSWithTimelockState struct. If none of the contracts are found, the state struct will be nil.
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
func MaybeLoadMCMSSolanaWithTimelockContracts(solChain deployment.SolChain, addresses map[string]deployment.TypeAndVersion) (*MCMSWithTimelockProgramsSolana, error) {
	state := MCMSWithTimelockProgramsSolana{}
	// We expect one of each contract on the chain.
	timelock := deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
	callProxy := deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
	proposer := deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
	canceller := deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
	bypasser := deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)

	// Ensure we either have the bundle or not.
	_, err := deployment.AddressesContainBundle(addresses,
		map[deployment.TypeAndVersion]struct{}{
			timelock: {}, proposer: {}, canceller: {}, bypasser: {}, callProxy: {},
		})
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", solChain.Name(), err)
	}

	for address, tvStr := range addresses {
		switch tvStr {
		case timelock:
			pub := solana.MustPublicKeyFromBase58(address)
			state.Timelock = pub
		case callProxy:
			pub := solana.MustPublicKeyFromBase58(address)
			state.CallProxy = pub
		case proposer:
			pub := solana.MustPublicKeyFromBase58(address)
			state.ProposerMcm = pub
		case bypasser:
			pub := solana.MustPublicKeyFromBase58(address)
			state.BypasserMcm = pub
		case canceller:
			pub := solana.MustPublicKeyFromBase58(address)
			state.CancellerMcm = pub
		}
	}
	return &state, nil
}
