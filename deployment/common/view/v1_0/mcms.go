package v1_0

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
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

type MCMSView struct {
	types.ContractMetaData
	// Note config is json marshallable.
	Config config.Config `json:"config"`
}

func GenerateMCMSView(mcms owner_helpers.ManyChainMultiSig) (MCMSView, error) {
	owner, err := mcms.Owner(nil)
	if err != nil {
		return MCMSView{}, nil
	}
	c, err := mcms.GetConfig(nil)
	if err != nil {
		return MCMSView{}, nil
	}
	parsedConfig, err := config.NewConfigFromRaw(c)
	if err != nil {
		return MCMSView{}, nil
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
) (MCMSWithTimelockView, error) {
	timelockView, err := GenerateTimelockView(timelock)
	if err != nil {
		return MCMSWithTimelockView{}, nil
	}
	callProxyView, err := GenerateCallProxyView(owner_helpers.CallProxy{})
	if err != nil {
		return MCMSWithTimelockView{}, nil
	}
	bypasserView, err := GenerateMCMSView(bypasser)
	if err != nil {
		return MCMSWithTimelockView{}, nil
	}
	proposerView, err := GenerateMCMSView(proposer)
	if err != nil {
		return MCMSWithTimelockView{}, nil
	}
	cancellerView, err := GenerateMCMSView(canceller)
	if err != nil {
		return MCMSWithTimelockView{}, nil
	}

	return MCMSWithTimelockView{
		Timelock:  timelockView,
		Bypasser:  bypasserView,
		Proposer:  proposerView,
		Canceller: cancellerView,
		CallProxy: callProxyView,
	}, nil
}
