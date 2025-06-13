package ops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	mcms_shared "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/shared"
)

type OpEVMOwnershipDeps struct {
	Chain    cldf_evm.Chain
	OwnableC mcms_shared.Ownable // Ownable contract interface
}

type OpEVMTransferOwnershipInput struct {
	ChainSelector   uint64         `json:"chainSelector"`   // Chain selector for the EVM chain
	TimelockAddress common.Address `json:"timelockAddress"` // Address of the EVM Timelock contract
	Address         common.Address `json:"address"`         // Address to Contract for which ownership is being transferred
}

type OpEVMOwnershipOutput struct {
	Tx *gethtypes.Transaction `json:"tx"`
}

var OpEVMTransferOwnership = operations.NewOperation(
	"evm-transfer-ownership",
	semver.MustParse("1.0.0"),
	"Transfer ownership of an ownable contract to the specified address",
	func(b operations.Bundle, deps OpEVMOwnershipDeps, in OpEVMTransferOwnershipInput) (OpEVMOwnershipOutput, error) {
		tx, err := deps.OwnableC.TransferOwnership(deps.Chain.DeployerKey, common.HexToAddress(in.TimelockAddress.Hex()))
		_, err = cldf.ConfirmIfNoError(deps.Chain, tx, err)
		if err != nil {
			return OpEVMOwnershipOutput{Tx: tx}, fmt.Errorf(
				"failed to transfer ownership of contract %T: %w",
				in.Address.Hex(),
				err,
			)
		}
		return OpEVMOwnershipOutput{
			Tx: tx,
		}, nil
	})

var OpEVMAcceptOwnership = operations.NewOperation(
	"evm-accept-ownership",
	semver.MustParse("1.0.0"),
	"Accepts ownership of an ownable contract Via the Timelock contract",
	func(b operations.Bundle, deps OpEVMOwnershipDeps, in OpEVMTransferOwnershipInput) (OpEVMOwnershipOutput, error) {
		tx, err := deps.OwnableC.AcceptOwnership(cldf.SimTransactOpts())
		if err != nil {
			return OpEVMOwnershipOutput{Tx: tx}, fmt.Errorf(
				"failed to Accept ownership of contract %T: %w",
				in.Address.Hex(),
				err,
			)
		}
		return OpEVMOwnershipOutput{
			Tx: tx,
		}, nil
	})
