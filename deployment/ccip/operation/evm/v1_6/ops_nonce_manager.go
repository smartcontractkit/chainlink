package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

var (
	DeployNonceManagerOp = opsutil.NewEVMDeployOperation(
		"DeployNonceManager",
		semver.MustParse("1.0.0"),
		"Deploys NonceManager 1.6 contract on the specified evm chain",
		cldf.NewTypeAndVersion(shared.NonceManager, deployment.Version1_6_0),
		opsutil.VMDeployers[[]common.Address]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, input []common.Address) (common.Address, *types.Transaction, error) {
				addr, tx, _, err := nonce_manager.DeployNonceManager(
					opts,
					backend,
					input,
				)
				return addr, tx, err
			},
			DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, input []common.Address) (common.Address, error) {
				addr, _, _, err := nonce_manager.DeployNonceManagerZk(
					opts,
					client,
					wallet,
					backend,
					input,
				)
				return addr, err
			},
		},
	)

	NonceManagerUpdateAuthorizedCallerOp = opsutil.NewEVMCallOperation(
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

	NonceManagerPreviousRampsUpdatesOp = opsutil.NewEVMCallOperation(
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
