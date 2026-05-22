package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	opsevm "github.com/smartcontractkit/cld-changesets/pkg/family/evm/operations"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

var (
	DeployNonceManagerOp = opsevm.NewEVMDeployOperation(
		"DeployNonceManager",
		semver.MustParse("1.0.0"),
		"Deploys NonceManager 1.6 contract on the specified evm chain",
		shared.NonceManager,
		nonce_manager.NonceManagerMetaData,
		&opsevm.ContractOpts{
			Version:          &deployment.Version1_6_0,
			EVMBytecode:      common.FromHex(nonce_manager.NonceManagerBin),
			ZkSyncVMBytecode: nonce_manager.ZkBytecode,
		},
		func(input []common.Address) []any {
			return []any{input}
		},
	)

	NonceManagerUpdateAuthorizedCallerOp = opsevm.NewEVMCallOperation(
		"NonceManagerUpdateAuthorizedCallerOp",
		semver.MustParse("1.0.0"),
		"Updates authorized callers in NonceManager 1.6 contract on the specified evm chain",
		nonce_manager.NonceManagerABI,
		shared.NonceManager,
		nonce_manager.NewNonceManager,
		func(nonceManager *nonce_manager.NonceManager, opts *bind.TransactOpts, input nonce_manager.AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
			return nonceManager.ApplyAuthorizedCallerUpdates(opts, input)
		},
	)

	NonceManagerPreviousRampsUpdatesOp = opsevm.NewEVMCallOperation(
		"NonceManagerPreviousRampsUpdatesOp",
		semver.MustParse("1.0.0"),
		"Applies previous ramps updates in NonceManager 1.6 contract on the specified evm chain",
		nonce_manager.NonceManagerABI,
		shared.NonceManager,
		nonce_manager.NewNonceManager,
		func(nonceManager *nonce_manager.NonceManager, opts *bind.TransactOpts, input []nonce_manager.NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
			return nonceManager.ApplyPreviousRampsUpdates(opts, input)
		},
	)
)
