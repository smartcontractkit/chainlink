package changeset

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_remote"
)

// RMN REMOTE

type DeployRMNRemoteInput struct {
	LocalChainSelector uint64
	LegacyRMN          common.Address
}

func wrapRMNRemoteDeploy(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input DeployRMNRemoteInput,
) (common.Address, *types.Transaction, *rmn_remote.RMNRemote, error) {
	return rmn_remote.DeployRMNRemote(
		auth,
		backend,
		input.LocalChainSelector,
		input.LegacyRMN,
	)
}

var DeployRMNRemoteOp = deployment_ethereum.NewEthDeployOperationFromBinding(wrapRMNRemoteDeploy, "v1")

type SetConfigRMNRemoteInput struct {
	contractAddress common.Address
	rmn_remote.RMNRemoteConfig
}

func (i SetConfigRMNRemoteInput) GetOrderedParams() []any {
	return []any{i.RMNRemoteConfig}
}

func (i SetConfigRMNRemoteInput) Address() common.Address {
	return i.contractAddress
}

var RMNRemoteSetConfigOp = deployment_ethereum.NewEthOperationFromBinding[SetConfigRMNRemoteInput](rmn_remote.RMNRemoteMetaData, "v1", "setConfig")

// ROUTER

type DeployRouterInput struct {
	wrappedNative common.Address
	armProxy      common.Address
}

func wrapRouterDeploy(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input DeployRouterInput,
) (common.Address, *types.Transaction, *router.Router, error) {
	return router.DeployRouter(
		auth,
		backend,
		input.wrappedNative,
		input.armProxy,
	)
}

var DeployRouterOp = deployment_ethereum.NewEthDeployOperationFromBinding(wrapRouterDeploy, "v1")

// NONCE MANAGER

type DeployNonceManagerInput struct {
	AuthorizedCallers []common.Address
}

func wrapNonceManagerDeploy(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input DeployNonceManagerInput,
) (common.Address, *types.Transaction, *nonce_manager.NonceManager, error) {
	return nonce_manager.DeployNonceManager(
		auth,
		backend,
		input.AuthorizedCallers,
	)
}

var DeployNonceManagerOp = deployment_ethereum.NewEthDeployOperationFromBinding(wrapNonceManagerDeploy, "v1")

type NMApplyAuthorizedCallerUpdatesInput struct {
	contractAddress common.Address
	nonce_manager.AuthorizedCallersAuthorizedCallerArgs
}

func (i NMApplyAuthorizedCallerUpdatesInput) GetOrderedParams() []any {
	return []any{i.AuthorizedCallersAuthorizedCallerArgs}
}

func (i NMApplyAuthorizedCallerUpdatesInput) Address() common.Address {
	return i.contractAddress
}

var NonceManagerApplyCallerUpdatesOp = deployment_ethereum.NewEthOperationFromBinding[NMApplyAuthorizedCallerUpdatesInput](nonce_manager.NonceManagerMetaData, "v1", "applyAuthorizedCallerUpdates")

// FEE QUOTER

type DeployFeeQuoterInput struct {
	staticConfig                   fee_quoter.FeeQuoterStaticConfig
	priceUpdaters                  []common.Address
	feeTokens                      []common.Address
	tokenPriceFeeds                []fee_quoter.FeeQuoterTokenPriceFeedUpdate
	tokenTransferFeeConfigArgs     []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
	premiumMultiplierWeiPerEthArgs []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
	destChainConfigArgs            []fee_quoter.FeeQuoterDestChainConfigArgs
}

func wrapFeeQuoterDeploy(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input DeployFeeQuoterInput,
) (common.Address, *types.Transaction, *fee_quoter.FeeQuoter, error) {
	return fee_quoter.DeployFeeQuoter(
		auth,
		backend,
		input.staticConfig,
		input.priceUpdaters,
		input.feeTokens,
		input.tokenPriceFeeds,
		input.tokenTransferFeeConfigArgs,
		input.premiumMultiplierWeiPerEthArgs,
		input.destChainConfigArgs,
	)
}

var DeployFeeQuoterOp = deployment_ethereum.NewEthDeployOperationFromBinding(wrapFeeQuoterDeploy, "v1")

type ApplyAuthorizedCallerUpdatesInput struct {
	contractAddress common.Address
	fee_quoter.AuthorizedCallersAuthorizedCallerArgs
}

func (i ApplyAuthorizedCallerUpdatesInput) GetOrderedParams() []any {
	return []any{i.AuthorizedCallersAuthorizedCallerArgs}
}

func (i ApplyAuthorizedCallerUpdatesInput) Address() common.Address {
	return i.contractAddress
}

var FeeQuoterApplyCallerUpdatesOp = deployment_ethereum.NewEthOperationFromBinding[ApplyAuthorizedCallerUpdatesInput](fee_quoter.FeeQuoterMetaData, "v1", "applyAuthorizedCallerUpdates")

// ONRAMP

type DeployOnRampInput struct {
	staticConfig        onramp.OnRampStaticConfig
	dynamicConfig       onramp.OnRampDynamicConfig
	destChainConfigArgs []onramp.OnRampDestChainConfigArgs
}

func wrapOnRampDeploy(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input DeployOnRampInput,
) (common.Address, *types.Transaction, *onramp.OnRamp, error) {
	return onramp.DeployOnRamp(
		auth,
		backend,
		input.staticConfig,
		input.dynamicConfig,
		input.destChainConfigArgs,
	)
}

var DeployOnRampOp = deployment_ethereum.NewEthDeployOperationFromBinding(wrapOnRampDeploy, "v1")

// OFFRAMP

type DeployOffRampInput struct {
	staticConfig       offramp.OffRampStaticConfig
	dynamicConfig      offramp.OffRampDynamicConfig
	sourceChainConfigs []offramp.OffRampSourceChainConfigArgs
}

func wrapOffRampDeploy(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input DeployOffRampInput,
) (common.Address, *types.Transaction, *offramp.OffRamp, error) {
	return offramp.DeployOffRamp(
		auth,
		backend,
		input.staticConfig,
		input.dynamicConfig,
		input.sourceChainConfigs,
	)
}

var DeployOffRampOp = deployment_ethereum.NewEthDeployOperationFromBinding(wrapOffRampDeploy, "v1")
