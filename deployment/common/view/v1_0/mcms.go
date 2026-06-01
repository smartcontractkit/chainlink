package v1_0

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/evm/bindings"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	commoncldchangesets "github.com/smartcontractkit/cld-changesets/pkg/cldfutil"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
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
	commoncldchangesets.ContractMetaData
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
		ContractMetaData: commoncldchangesets.ContractMetaData{
			Owner:   owner,
			Address: mcms.Address(),
		},
		Config: *parsedConfig,
	}, nil
}

type TimelockView struct {
	commoncldchangesets.ContractMetaData
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
		ContractMetaData: commoncldchangesets.ContractMetaData{
			Address: tl.Address(),
		},
		MembersByRole: membersByRole,
	}, nil
}

type CallProxyView struct {
	commoncldchangesets.ContractMetaData
}

func GenerateCallProxyView(cp owner_helpers.CallProxy) (CallProxyView, error) {
	return CallProxyView{
		ContractMetaData: commoncldchangesets.ContractMetaData{
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
