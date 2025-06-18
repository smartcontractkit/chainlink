package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	view "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

type PDASeed [32]byte

// MCMSWithTimelockProgramsSolana holds the solana publick keys
// and seeds for the MCM, AccessController and Timelock programs.
// It is public for use in product specific packages.
type MCMSWithTimelockProgramsSolana struct {
	McmProgram                       solana.PublicKey
	ProposerMcmSeed                  PDASeed
	CancellerMcmSeed                 PDASeed
	BypasserMcmSeed                  PDASeed
	TimelockProgram                  solana.PublicKey
	TimelockSeed                     PDASeed
	AccessControllerProgram          solana.PublicKey
	ProposerAccessControllerAccount  solana.PublicKey
	ExecutorAccessControllerAccount  solana.PublicKey
	CancellerAccessControllerAccount solana.PublicKey
	BypasserAccessControllerAccount  solana.PublicKey
}

func (s *MCMSWithTimelockProgramsSolana) GetStateFromType(programType cldf.ContractType) (solana.PublicKey, PDASeed, error) {
	switch programType {
	case types.ManyChainMultisigProgram:
		return s.McmProgram, PDASeed{}, nil
	case types.ProposerManyChainMultisig:
		return s.McmProgram, s.ProposerMcmSeed, nil
	case types.BypasserManyChainMultisig:
		return s.McmProgram, s.BypasserMcmSeed, nil
	case types.CancellerManyChainMultisig:
		return s.McmProgram, s.CancellerMcmSeed, nil
	case types.RBACTimelockProgram:
		return s.TimelockProgram, PDASeed{}, nil
	case types.RBACTimelock:
		return s.TimelockProgram, s.TimelockSeed, nil
	case types.AccessControllerProgram:
		return s.AccessControllerProgram, PDASeed{}, nil
	case types.ProposerAccessControllerAccount:
		return s.AccessControllerProgram, PDASeed(s.ProposerAccessControllerAccount), nil
	case types.ExecutorAccessControllerAccount:
		return s.AccessControllerProgram, PDASeed(s.ExecutorAccessControllerAccount), nil
	case types.CancellerAccessControllerAccount:
		return s.AccessControllerProgram, PDASeed(s.CancellerAccessControllerAccount), nil
	case types.BypasserAccessControllerAccount:
		return s.AccessControllerProgram, PDASeed(s.BypasserAccessControllerAccount), nil
	default:
		return solana.PublicKey{}, PDASeed{}, fmt.Errorf("unknown program type: %s", programType)
	}
}

func (s *MCMSWithTimelockProgramsSolana) SetState(contractType cldf.ContractType, program solana.PublicKey, seed PDASeed) error {
	switch contractType {
	case types.ManyChainMultisigProgram:
		s.McmProgram = program
	case types.ProposerManyChainMultisig:
		s.McmProgram = program
		s.ProposerMcmSeed = seed
	case types.BypasserManyChainMultisig:
		s.McmProgram = program
		s.BypasserMcmSeed = seed
	case types.CancellerManyChainMultisig:
		s.McmProgram = program
		s.CancellerMcmSeed = seed
	case types.RBACTimelockProgram:
		s.TimelockProgram = program
	case types.RBACTimelock:
		s.TimelockProgram = program
		s.TimelockSeed = seed
	case types.AccessControllerProgram:
		s.AccessControllerProgram = program
	case types.ProposerAccessControllerAccount:
		s.ProposerAccessControllerAccount = program
	case types.ExecutorAccessControllerAccount:
		s.ExecutorAccessControllerAccount = program
	case types.CancellerAccessControllerAccount:
		s.CancellerAccessControllerAccount = program
	case types.BypasserAccessControllerAccount:
		s.BypasserAccessControllerAccount = program
	default:
		return fmt.Errorf("unknown contract type: %s", contractType)
	}
	return nil
}

// Validate checks that all fields are non-nil, ensuring it's ready
// for use generating views or interactions.
func (s *MCMSWithTimelockProgramsSolana) Validate() error {
	if s.McmProgram.IsZero() {
		return errors.New("mcm program not found")
	}
	if s.TimelockProgram.IsZero() {
		return errors.New("timelock program not found")
	}
	if s.AccessControllerProgram.IsZero() {
		return errors.New("access controller program not found")
	}
	if s.ProposerAccessControllerAccount.IsZero() {
		return errors.New("proposer access controller account not found")
	}
	if s.ExecutorAccessControllerAccount.IsZero() {
		return errors.New("executor access controller account not found")
	}
	if s.CancellerAccessControllerAccount.IsZero() {
		return errors.New("canceller access controller account not found")
	}
	if s.BypasserAccessControllerAccount.IsZero() {
		return errors.New("bypasser access controller account not found")
	}
	return nil
}

func (s *MCMSWithTimelockProgramsSolana) RoleAccount(role timelockBindings.Role) solana.PublicKey {
	switch role {
	case timelockBindings.Proposer_Role:
		return s.ProposerAccessControllerAccount
	case timelockBindings.Executor_Role:
		return s.ExecutorAccessControllerAccount
	case timelockBindings.Canceller_Role:
		return s.CancellerAccessControllerAccount
	case timelockBindings.Bypasser_Role:
		return s.BypasserAccessControllerAccount
	default:
		return solana.PublicKey{}
	}
}

func (s *MCMSWithTimelockProgramsSolana) GenerateView(
	ctx context.Context, chain cldf_solana.Chain,
) (view.MCMSWithTimelockViewSolana, error) {
	if err := s.Validate(); err != nil {
		return view.MCMSWithTimelockViewSolana{}, fmt.Errorf("unable to validate state: %w", err)
	}

	inspector := mcmssolanasdk.NewInspector(chain.Client)
	timelockInspector := mcmssolanasdk.NewTimelockInspector(chain.Client)

	return view.GenerateMCMSWithTimelockViewSolana(ctx, inspector, timelockInspector, s.McmProgram,
		s.ProposerMcmSeed, s.CancellerMcmSeed, s.BypasserMcmSeed, s.TimelockProgram, s.TimelockSeed)
}

// MCMSWithTimelockStateStateSolana holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockStateSolana struct {
	*MCMSWithTimelockProgramsSolana
}

// MaybeLoadMCMSWithTimelockState loads the MCMSWithTimelockState state for each chain in the given environment.
func MaybeLoadMCMSWithTimelockStateSolana(env cldf.Environment, chainSelectors []uint64) (map[uint64]*MCMSWithTimelockStateSolana, error) {
	result := map[uint64]*MCMSWithTimelockStateSolana{}
	solChains := env.BlockChains.SolanaChains()
	for _, chainSelector := range chainSelectors {
		chain, ok := solChains[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chainSelector)
		}
		addressesChain, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return nil, fmt.Errorf("unable to get addresses for chain %v: %w", chainSelector, err)
			}
			// chain not found in address book, initialize empty
			addressesChain = make(map[string]cldf.TypeAndVersion)
		}
		state, err := MaybeLoadMCMSWithTimelockChainStateSolana(chain, addressesChain)
		if err != nil {
			return nil, fmt.Errorf("unable to load mcms and timelock solana chain state: %w", err)
		}
		result[chainSelector] = state
	}
	return result, nil
}

// MaybeLoadMCMSWithTimelockChainState looks for the addresses corresponding to
// contracts deployed with DeployMCMSWithTimelock and loads them into a
// MCMSWithTimelockStateSolana struct. If none of the contracts are found, the
// state struct will be nil.
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
func MaybeLoadMCMSWithTimelockChainStateSolana(chain cldf_solana.Chain, addresses map[string]cldf.TypeAndVersion) (*MCMSWithTimelockStateSolana, error) {
	state := MCMSWithTimelockStateSolana{MCMSWithTimelockProgramsSolana: &MCMSWithTimelockProgramsSolana{}}

	mcmProgram := cldf.NewTypeAndVersion(types.ManyChainMultisigProgram, deployment.Version1_0_0)
	timelockProgram := cldf.NewTypeAndVersion(types.RBACTimelockProgram, deployment.Version1_0_0)
	accessControllerProgram := cldf.NewTypeAndVersion(types.AccessControllerProgram, deployment.Version1_0_0)
	proposerMCM := cldf.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
	cancellerMCM := cldf.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
	bypasserMCM := cldf.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)
	timelock := cldf.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
	proposerAccessControllerAccount := cldf.NewTypeAndVersion(types.ProposerAccessControllerAccount, deployment.Version1_0_0)
	executorAccessControllerAccount := cldf.NewTypeAndVersion(types.ExecutorAccessControllerAccount, deployment.Version1_0_0)
	cancellerAccessControllerAccount := cldf.NewTypeAndVersion(types.CancellerAccessControllerAccount, deployment.Version1_0_0)
	bypasserAccessControllerAccount := cldf.NewTypeAndVersion(types.BypasserAccessControllerAccount, deployment.Version1_0_0)

	// Convert map keys to a slice
	wantTypes := []cldf.TypeAndVersion{
		mcmProgram, timelockProgram, accessControllerProgram, proposerMCM, cancellerMCM, bypasserMCM, timelock,
		proposerAccessControllerAccount, executorAccessControllerAccount, cancellerAccessControllerAccount,
		bypasserAccessControllerAccount,
	}

	// Ensure we either have the bundle or not.
	_, err := cldf.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", chain.Name(), err)
	}

	for address, tvStr := range addresses {
		switch {
		case tvStr.Type == timelockProgram.Type && tvStr.Version.String() == timelockProgram.Version.String():
			programID, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode timelock program address (%s): %w", address, err)
			}
			state.TimelockProgram = programID

		case tvStr.Type == timelock.Type && tvStr.Version.String() == timelock.Version.String():
			programID, seed, err := DecodeAddressWithSeed(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode timelock address (%s): %w", address, err)
			}
			state.TimelockProgram = programID
			state.TimelockSeed = seed

		case tvStr.Type == accessControllerProgram.Type && tvStr.Version.String() == accessControllerProgram.Version.String():
			programID, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to parse public key from access controller address: %s", address)
			}
			state.AccessControllerProgram = programID

		case tvStr.Type == proposerAccessControllerAccount.Type && tvStr.Version.String() == proposerAccessControllerAccount.Version.String():
			account, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode proposer access controller address (%s): %w", address, err)
			}
			state.ProposerAccessControllerAccount = account

		case tvStr.Type == executorAccessControllerAccount.Type && tvStr.Version.String() == executorAccessControllerAccount.Version.String():
			account, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode executor access controller address (%s): %w", address, err)
			}
			state.ExecutorAccessControllerAccount = account

		case tvStr.Type == cancellerAccessControllerAccount.Type && tvStr.Version.String() == cancellerAccessControllerAccount.Version.String():
			account, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode canceller access controller address (%s): %w", address, err)
			}
			state.CancellerAccessControllerAccount = account

		case tvStr.Type == bypasserAccessControllerAccount.Type && tvStr.Version.String() == bypasserAccessControllerAccount.Version.String():
			account, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode bypasser access controller address (%s): %w", address, err)
			}
			state.BypasserAccessControllerAccount = account

		case tvStr.Type == mcmProgram.Type && tvStr.Version.String() == mcmProgram.Version.String():
			programID, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return nil, fmt.Errorf("unable to parse public key from mcm address: %s", address)
			}
			state.McmProgram = programID

		case tvStr.Type == proposerMCM.Type && tvStr.Version.String() == proposerMCM.Version.String():
			programID, seed, err := DecodeAddressWithSeed(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode proposer address (%s): %w", address, err)
			}
			state.McmProgram = programID
			state.ProposerMcmSeed = seed

		case tvStr.Type == bypasserMCM.Type && tvStr.Version.String() == bypasserMCM.Version.String():
			programID, seed, err := DecodeAddressWithSeed(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode bypasser address (%s): %w", address, err)
			}
			state.McmProgram = programID
			state.BypasserMcmSeed = seed

		case tvStr.Type == cancellerMCM.Type && tvStr.Version.String() == cancellerMCM.Version.String():
			programID, seed, err := DecodeAddressWithSeed(address)
			if err != nil {
				return nil, fmt.Errorf("unable to decode canceller address (%s): %w", address, err)
			}
			state.McmProgram = programID
			state.CancellerMcmSeed = seed
		}
	}
	return &state, nil
}

func EncodeAddressWithSeed(programID solana.PublicKey, seed PDASeed) string {
	return mcmssolanasdk.ContractAddress(programID, mcmssolanasdk.PDASeed(seed))
}

func DecodeAddressWithSeed(address string) (solana.PublicKey, PDASeed, error) {
	programID, seed, err := mcmssolanasdk.ParseContractAddress(address)
	if err != nil {
		return solana.PublicKey{}, PDASeed{}, fmt.Errorf("unable to parse address: %s", address)
	}

	return programID, PDASeed(seed), nil
}
