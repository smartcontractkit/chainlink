package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

var (
	DeployNonceManagerOp = operations.NewOperation(
		"DeployNonceManager",
		semver.MustParse("1.0.0"),
		"Deploys NonceManager 1.6 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.DeployContractDependencies, input uint64) (common.Address, error) {
			ab := deps.AddressBook
			chain := deps.Chain
			nonceManager, err := cldf.DeployContract(b.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*nonce_manager.NonceManager] {
					var (
						nonceManagerAddr common.Address
						tx2              *types.Transaction
						nonceManager     *nonce_manager.NonceManager
						err2             error
					)
					if chain.IsZkSyncVM {
						nonceManagerAddr, _, nonceManager, err2 = nonce_manager.DeployNonceManagerZk(
							nil,
							chain.ClientZkSyncVM,
							chain.DeployerKeyZkSyncVM,
							chain.Client,
							[]common.Address{},
						)
					} else {
						nonceManagerAddr, tx2, nonceManager, err2 = nonce_manager.DeployNonceManager(
							chain.DeployerKey,
							chain.Client,
							[]common.Address{}, // Need to add onRamp after
						)
					}
					return cldf.ContractDeploy[*nonce_manager.NonceManager]{
						Address: nonceManagerAddr, Contract: nonceManager, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.NonceManager, deployment.Version1_6_0), Err: err2,
					}
				})
			if err != nil {
				b.Logger.Errorw("Failed to deploy nonce manager", "chain", chain.String(), "err", err)
				return common.Address{}, err
			}
			return nonceManager.Address, nil
		})

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
