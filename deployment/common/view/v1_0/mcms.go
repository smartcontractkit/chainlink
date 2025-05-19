package v1_0

import (
	"context"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/evm/bindings"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type Role struct {
	ID   common.Hash
	Name string
}

const (
	EXECUTOR_ROLE_STR  = "EXECUTOR_ROLE"
	BYPASSER_ROLE_STR  = "BYPASSER_ROLE"
	CANCELLER_ROLE_STR = "CANCELLER_ROLE"
	PROPOSER_ROLE_STR  = "PROPOSER_ROLE"
	ADMIN_ROLE_STR     = "ADMIN_ROLE"
)

// https://github.com/smartcontractkit/ccip-owner-contracts/blob/9d81692b324ce7ea2ef8a75e683889edbc7e2dd0/src/RBACTimelock.sol#L71
// Just to avoid invoking the Go binding to get these.
var (
	ADMIN_ROLE = Role{
		ID:   utils.MustHash(ADMIN_ROLE_STR),
		Name: ADMIN_ROLE_STR,
	}
	PROPOSER_ROLE = Role{
		ID:   utils.MustHash(PROPOSER_ROLE_STR),
		Name: PROPOSER_ROLE_STR,
	}
	BYPASSER_ROLE = Role{
		ID:   utils.MustHash(BYPASSER_ROLE_STR),
		Name: BYPASSER_ROLE_STR,
	}
	CANCELLER_ROLE = Role{
		ID:   utils.MustHash(CANCELLER_ROLE_STR),
		Name: CANCELLER_ROLE_STR,
	}
	EXECUTOR_ROLE = Role{
		ID:   utils.MustHash(EXECUTOR_ROLE_STR),
		Name: EXECUTOR_ROLE_STR,
	}
)

// --- evm ---

type MCMSView struct {
	types.ContractMetaData
	// Note config is json marshallable.
	Config mcmstypes.Config `json:"config"`
}

func GenerateMCMSView(mcms owner_helpers.ManyChainMultiSig) (MCMSView, error) {
	owner, err := mcms.Owner(nil)
	if err != nil {
		return MCMSView{}, err
	}
	mcmsConfig, err := mcms.GetConfig(nil)
	if err != nil {
		return MCMSView{}, err
	}

	mapSigners := func(in []owner_helpers.ManyChainMultiSigSigner) []bindings.ManyChainMultiSigSigner {
		out := make([]bindings.ManyChainMultiSigSigner, len(in))
		for i, s := range in {
			out[i] = bindings.ManyChainMultiSigSigner{Addr: s.Addr, Index: s.Index, Group: s.Group}
		}
		return out
	}

	parsedConfig, err := mcmsevmsdk.NewConfigTransformer().ToConfig(bindings.ManyChainMultiSigConfig{
		Signers:      mapSigners(mcmsConfig.Signers),
		GroupQuorums: mcmsConfig.GroupQuorums,
		GroupParents: mcmsConfig.GroupParents,
	})
	if err != nil {
		return MCMSView{}, err
	}
	return MCMSView{
		// Has no type and version on the contract
		ContractMetaData: types.ContractMetaData{
			Owner:   owner,
			Address: mcms.Address(),
		},
		Config: *parsedConfig,
	}, nil
}

type TimelockView struct {
	types.ContractMetaData
	MembersByRole map[string][]common.Address `json:"membersByRole"`
}

func GenerateTimelockView(tl owner_helpers.RBACTimelock) (TimelockView, error) {
	membersByRole := make(map[string][]common.Address)
	for _, role := range []Role{ADMIN_ROLE, PROPOSER_ROLE, BYPASSER_ROLE, CANCELLER_ROLE, EXECUTOR_ROLE} {
		numMembers, err := tl.GetRoleMemberCount(nil, role.ID)
		if err != nil {
			return TimelockView{}, nil
		}
		for i := int64(0); i < numMembers.Int64(); i++ {
			member, err2 := tl.GetRoleMember(nil, role.ID, big.NewInt(i))
			if err2 != nil {
				return TimelockView{}, nil
			}
			membersByRole[role.Name] = append(membersByRole[role.Name], member)
		}
	}
	return TimelockView{
		// Has no type and version or owner.
		ContractMetaData: types.ContractMetaData{
			Address: tl.Address(),
		},
		MembersByRole: membersByRole,
	}, nil
}

type CallProxyView struct {
	types.ContractMetaData
}

func GenerateCallProxyView(cp owner_helpers.CallProxy) (CallProxyView, error) {
	return CallProxyView{
		ContractMetaData: types.ContractMetaData{
			Address: cp.Address(),
		},
	}, nil
}

type MCMSWithTimelockView struct {
	Bypasser  MCMSView      `json:"bypasser"`
	Canceller MCMSView      `json:"canceller"`
	Proposer  MCMSView      `json:"proposer"`
	Timelock  TimelockView  `json:"timelock"`
	CallProxy CallProxyView `json:"callProxy"`
}

func GenerateMCMSWithTimelockView(
	bypasser owner_helpers.ManyChainMultiSig,
	canceller owner_helpers.ManyChainMultiSig,
	proposer owner_helpers.ManyChainMultiSig,
	timelock owner_helpers.RBACTimelock,
	callProxy owner_helpers.CallProxy,
) (MCMSWithTimelockView, error) {
	timelockView, err := GenerateTimelockView(timelock)
	if err != nil {
		return MCMSWithTimelockView{}, err
	}
	callProxyView, err := GenerateCallProxyView(callProxy)
	if err != nil {
		return MCMSWithTimelockView{}, err
	}
	bypasserView, err := GenerateMCMSView(bypasser)
	if err != nil {
		return MCMSWithTimelockView{}, err
	}
	proposerView, err := GenerateMCMSView(proposer)
	if err != nil {
		return MCMSWithTimelockView{}, err
	}
	cancellerView, err := GenerateMCMSView(canceller)
	if err != nil {
		return MCMSWithTimelockView{}, err
	}

	return MCMSWithTimelockView{
		Timelock:  timelockView,
		Bypasser:  bypasserView,
		Proposer:  proposerView,
		Canceller: cancellerView,
		CallProxy: callProxyView,
	}, nil
}

// --- solana ---

type MCMSWithTimelockViewSolana struct {
	Bypasser  MCMViewSolana      `json:"bypasser"`
	Canceller MCMViewSolana      `json:"canceller"`
	Proposer  MCMViewSolana      `json:"proposer"`
	Timelock  TimelockViewSolana `json:"timelock"`
}

func GenerateMCMSWithTimelockViewSolana(
	ctx context.Context,
	inspector *mcmssolanasdk.Inspector,
	timelockInspector *mcmssolanasdk.TimelockInspector,
	mcmProgram solana.PublicKey,
	proposerMcmSeed [32]byte,
	cancellerMcmSeed [32]byte,
	bypasserMcmSeed [32]byte,
	timelockProgram solana.PublicKey,
	timelockSeed [32]byte,
) (MCMSWithTimelockViewSolana, error) {
	timelockView, err := GenerateTimelockViewSolana(ctx, timelockInspector, timelockProgram, timelockSeed)
	if err != nil {
		return MCMSWithTimelockViewSolana{}, fmt.Errorf("unable to generate timelock view: %w", err)
	}
	bypasserView, err := GenerateMCMViewSolana(ctx, inspector, mcmProgram, bypasserMcmSeed)
	if err != nil {
		return MCMSWithTimelockViewSolana{}, fmt.Errorf("unable to generate bypasser mcm view: %w", err)
	}
	proposerView, err := GenerateMCMViewSolana(ctx, inspector, mcmProgram, proposerMcmSeed)
	if err != nil {
		return MCMSWithTimelockViewSolana{}, fmt.Errorf("unable to generate proposer mcm view: %w", err)
	}
	cancellerView, err := GenerateMCMViewSolana(ctx, inspector, mcmProgram, cancellerMcmSeed)
	if err != nil {
		return MCMSWithTimelockViewSolana{}, fmt.Errorf("unable to generate canceller mcm view: %w", err)
	}

	return MCMSWithTimelockViewSolana{
		Timelock:  timelockView,
		Bypasser:  bypasserView,
		Proposer:  proposerView,
		Canceller: cancellerView,
	}, nil
}

type MCMViewSolana struct {
	ProgramID solana.PublicKey `json:"programID"`
	Seed      string           `json:"seed"`
	Owner     solana.PublicKey `json:"owner"`
	Config    mcmstypes.Config `json:"config"`
}

func GenerateMCMViewSolana(
	ctx context.Context, inspector *mcmssolanasdk.Inspector, programID solana.PublicKey, seed [32]byte,
) (MCMViewSolana, error) {
	address := mcmssolanasdk.ContractAddress(programID, mcmssolanasdk.PDASeed(seed))
	config, err := inspector.GetConfig(ctx, address)
	if err != nil {
		return MCMViewSolana{}, fmt.Errorf("unable to get config from mcm (%v): %w", address, err)
	}

	return MCMViewSolana{
		ProgramID: programID,
		Seed:      string(seed[:]),
		Owner:     solana.PublicKey{}, // FIXME: needs inspector.GetOwner() in mcms solana sdk
		Config:    *config,
	}, nil
}

type TimelockViewSolana struct {
	ProgramID  solana.PublicKey `json:"programID"`
	Seed       string           `json:"seed"`
	Owner      solana.PublicKey `json:"owner"`
	Proposers  []string         `json:"proposers"`
	Executors  []string         `json:"executors"`
	Cancellers []string         `json:"cancellers"`
	Bypassers  []string         `json:"bypassers"`
}

func GenerateTimelockViewSolana(
	ctx context.Context, inspector *mcmssolanasdk.TimelockInspector, programID solana.PublicKey, seed [32]byte,
) (TimelockViewSolana, error) {
	address := mcmssolanasdk.ContractAddress(programID, mcmssolanasdk.PDASeed(seed))

	proposers, err := inspector.GetProposers(ctx, address)
	if err != nil {
		return TimelockViewSolana{}, fmt.Errorf("unable to get proposers from timelock (%v): %w", address, err)
	}
	executors, err := inspector.GetExecutors(ctx, address)
	if err != nil {
		return TimelockViewSolana{}, fmt.Errorf("unable to get executors from timelock (%v): %w", address, err)
	}
	cancellers, err := inspector.GetCancellers(ctx, address)
	if err != nil {
		return TimelockViewSolana{}, fmt.Errorf("unable to get cancellers from timelock (%v): %w", address, err)
	}
	bypassers, err := inspector.GetBypassers(ctx, address)
	if err != nil {
		return TimelockViewSolana{}, fmt.Errorf("unable to get bypassers from timelock (%v): %w", address, err)
	}

	return TimelockViewSolana{
		ProgramID:  programID,
		Seed:       string(seed[:]),
		Owner:      solana.PublicKey{}, // FIXME: needs inspector.GetOwner() in mcms solana sdk
		Proposers:  slices.Sorted(slices.Values(proposers)),
		Executors:  slices.Sorted(slices.Values(executors)),
		Cancellers: slices.Sorted(slices.Values(cancellers)),
		Bypassers:  slices.Sorted(slices.Values(bypassers)),
	}, nil
}
