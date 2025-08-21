package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMSetConfigMCMInput struct {
	SignerAddresses []common.Address `json:"signerAddresses"`
	SignerGroups    []uint8          `json:"signerGroups"` // Signer 1 is int group 0 (root group) with quorum 1.
	GroupQuorums    [32]uint8        `json:"groupQuorums"`
	GroupParents    [32]uint8        `json:"groupParents"`
}

var OpEVMSetConfigMCM = opsutils.NewEVMCallOperation(
	"evm-mcm-set-config",
	semver.MustParse("1.0.0"),
	"Sets Config on the deployed MCM contract",
	bindings.ManyChainMultiSigABI,
	commontypes.ManyChainMultisig,
	bindings.NewManyChainMultiSig,
	func(mcm *bindings.ManyChainMultiSig, opts *bind.TransactOpts, input OpEVMSetConfigMCMInput) (*types.Transaction, error) {
		return mcm.SetConfig(
			opts,
			input.SignerAddresses,
			input.SignerGroups,
			input.GroupQuorums,
			input.GroupParents,
			false,
		)
	})
